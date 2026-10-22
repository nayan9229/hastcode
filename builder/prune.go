package builder

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"

	"hastcode.com/site/internal/imagestore"
)

// Prune (spec "Image Store", research §4) finds the Image Store objects
// the current content would not produce: it builds the content tree with
// the local store into a temporary directory (nothing is uploaded) to get
// Result.ImageKeys, lists the bucket, and reports every other object as an
// orphan. It deletes them only when PruneOptions.Delete is set.

// PruneOptions are the prune command's switches.
type PruneOptions struct {
	// Delete removes the orphans; without it Prune is a dry run.
	Delete bool
	// Log receives the orphan list and totals (nil discards them).
	Log io.Writer
}

// PruneReport is what Prune found (and deleted).
type PruneReport struct {
	// Orphans are the bucket's objects no build of the content tree
	// references, in key order.
	Orphans []imagestore.ObjectInfo
	// Bytes is the orphans' total size.
	Bytes int64
	// Kept is how many listed objects the content references.
	Kept int
	// Deleted reports whether the orphans were deleted.
	Deleted bool
}

// Prune runs the prune command for contentDir against the Image Store
// cfg.Images names: images.bucket and images.base_url, and the credentials
// read through cfg.Getenv, whatever images.store says.
func Prune(ctx context.Context, cfg Config, contentDir string, opts PruneOptions) (PruneReport, error) {
	log := opts.Log
	if log == nil {
		log = io.Discard
	}
	if cfg.Getenv == nil {
		cfg.Getenv = os.Getenv
	}
	r2 := cfg
	r2.Images.Store = StoreR2
	client, err := openImageStore(r2)
	if err != nil {
		return PruneReport{}, Problem{Rule: ruleImageStore, Message: err.Error()}
	}

	keep, err := referencedImageKeys(cfg, contentDir)
	if err != nil {
		return PruneReport{}, err
	}
	listed, err := client.List(ctx, "")
	if err != nil {
		return PruneReport{}, fmt.Errorf("listing bucket %s: %w", cfg.Images.Bucket, err)
	}

	var rep PruneReport
	for _, o := range listed {
		if _, ok := slices.BinarySearch(keep, o.Key); ok {
			rep.Kept++
			continue
		}
		rep.Orphans = append(rep.Orphans, o)
		rep.Bytes += o.Size
		fmt.Fprintf(log, "%12d B  %s\n", o.Size, o.Key)
	}
	fmt.Fprintf(log, "%d orphan(s), %d B in bucket %s; %d referenced object(s) kept\n", len(rep.Orphans), rep.Bytes, cfg.Images.Bucket, rep.Kept)
	if !opts.Delete {
		fmt.Fprintln(log, "dry run: nothing deleted (pass --delete to delete the orphans)")
		return rep, nil
	}
	keys := make([]string, len(rep.Orphans))
	for i, o := range rep.Orphans {
		keys[i] = o.Key
	}
	if err := client.Delete(ctx, keys); err != nil {
		return rep, err
	}
	rep.Deleted = true
	fmt.Fprintf(log, "deleted %d object(s)\n", len(keys))
	return rep, nil
}

// referencedImageKeys are the keys the content tree's builds publish,
// sorted: a Production build's, plus a Preview build's, so the images a
// Preview uploaded for a draft still in the tree are kept. Both builds use
// the local store in a temporary directory, so nothing is uploaded and no
// search index is built. A Production build that fails (beyond the
// skipped index) is an error: a failed piece drops its images from
// ImageKeys, and they would be pruned as orphans. A failing Preview build
// (an unfinished draft) is tolerated; the drafts it could not build keep
// nothing.
func referencedImageKeys(cfg Config, contentDir string) ([]string, error) {
	tmp, err := os.MkdirTemp("", "hastcode-prune-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmp)

	cfg.Images.Store = StoreLocal
	cfg.SkipSearchIndex = true
	var keys []string
	for _, mode := range []Mode{Production, Preview} {
		cfg.Mode = mode
		res := Build(cfg, contentDir, filepath.Join(tmp, string(mode)))
		if mode == Production {
			if failed := slices.DeleteFunc(res.Errors, func(p Problem) bool { return p.Rule == ruleSearchIndex }); len(failed) > 0 {
				return nil, fmt.Errorf("the production build of %s failed, so its image keys are incomplete; fix it before pruning:\n%s", contentDir, Result{Errors: failed})
			}
		}
		keys = append(keys, res.ImageKeys...)
	}
	slices.Sort(keys)
	return slices.Compact(keys), nil
}

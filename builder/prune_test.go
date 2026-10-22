package builder_test

import (
	"bytes"
	"context"
	"fmt"
	"image/color"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildtest"
)

// orphanObjects are objects no build of r2Tree produces: an old hero
// variant, a removed Author's photo and a previous OG badge.
var orphanObjects = map[string]string{
	"heroes/retired-piece-1200x675.0123456789.jpg": "12345",
	"authors/gone-400x400.abcdef0123.jpg":          "1234567",
	"og-badge.000000000000.png":                    "123",
}

// publishedStore is a fake Image Store holding what an r2 build of tree
// uploaded, plus orphanObjects; it returns the config pointing at it and
// the keys the build uploaded.
func publishedStore(t *testing.T, tree buildtest.Tree) (*buildtest.ImageStore, builder.Config, []string) {
	t.Helper()
	store := buildtest.NewImageStore(t)
	cfg := store.Configure(buildtest.Config())
	buildtest.MustBuild(t, cfg, tree)
	referenced := store.Keys()
	for k, data := range orphanObjects {
		store.Seed(k, buildtest.StoredObject{Data: []byte(data), ContentType: "image/jpeg"})
	}
	store.ResetRequests()
	cfg.CacheDir = t.TempDir()
	return store, cfg, referenced
}

func sortedKeys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func orphanKeys(rep builder.PruneReport) []string {
	var keys []string
	for _, o := range rep.Orphans {
		keys = append(keys, o.Key)
	}
	return keys
}

func TestPrune_DryRunListsExactlyTheOrphansWithSizesAndDeletesNothing(t *testing.T) {
	tree := r2Tree()
	store, cfg, referenced := publishedStore(t, tree)
	before := store.Keys()

	var log bytes.Buffer
	rep, err := builder.Prune(context.Background(), cfg, tree.Write(t), builder.PruneOptions{Log: &log})
	if err != nil {
		t.Fatal(err)
	}

	if got, want := orphanKeys(rep), sortedKeys(orphanObjects); !slices.Equal(got, want) {
		t.Errorf("orphans:\n got %v\nwant %v", got, want)
	}
	for _, o := range rep.Orphans {
		if o.Size != int64(len(orphanObjects[o.Key])) {
			t.Errorf("%s size = %d, want %d", o.Key, o.Size, len(orphanObjects[o.Key]))
		}
	}
	if rep.Bytes != 15 || rep.Deleted {
		t.Errorf("report Bytes = %d, Deleted = %v; want 15, false", rep.Bytes, rep.Deleted)
	}
	if rep.Kept != len(referenced) {
		t.Errorf("report Kept = %d, want the %d referenced keys", rep.Kept, len(referenced))
	}
	if got := store.Keys(); !slices.Equal(got, before) {
		t.Errorf("dry run changed the bucket:\n got %v\nwant %v", got, before)
	}
	for _, r := range store.Requests("") {
		if r.Op != "ListObjectsV2" {
			t.Errorf("dry run sent %s %s", r.Op, r.Key)
		}
	}
	out := log.String()
	for k, data := range orphanObjects {
		if !strings.Contains(out, k) || !strings.Contains(out, " "+strconv.Itoa(len(data))+" B") {
			t.Errorf("output does not list %s with its size:\n%s", k, out)
		}
	}
	if !strings.Contains(out, "3 orphan(s), 15 B") || !strings.Contains(out, "--delete") {
		t.Errorf("output lacks the total or the dry-run note:\n%s", out)
	}
}

func TestPrune_KeepsAPreviewsDraftImagesOnlyWhileTheDraftIsInTheTree(t *testing.T) {
	d := buildtest.ValidExplainer()
	d.Slug = "raft-membership-changes"
	d.Title = "How Raft changes cluster membership"
	d.Hero = "raft-membership-changes.png"
	d.HeroImage = buildtest.HeroPNGFilled(1600, 900, color.RGBA{R: 40, G: 90, B: 160, A: 255})
	d.Draft = true
	withDraft := r2Tree().WithExplainer(d)

	store, cfg, published := publishedStore(t, r2Tree())
	preview := cfg
	preview.Mode = builder.Preview
	preview.CacheDir = ""
	buildtest.MustBuild(t, preview, withDraft)
	var draftKeys []string
	for _, k := range store.Keys() {
		if _, orphan := orphanObjects[k]; !orphan && !slices.Contains(published, k) {
			draftKeys = append(draftKeys, k)
		}
	}
	if len(draftKeys) == 0 || !slices.ContainsFunc(draftKeys, func(k string) bool { return strings.HasPrefix(k, "heroes/raft-membership-changes-") }) {
		t.Fatalf("the Preview uploaded no images for the draft: %v", draftKeys)
	}

	rep, err := builder.Prune(context.Background(), cfg, withDraft.Write(t), builder.PruneOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := orphanKeys(rep), sortedKeys(orphanObjects); !slices.Equal(got, want) {
		t.Errorf("draft in the tree: orphans\n got %v\nwant %v", got, want)
	}

	rep, err = builder.Prune(context.Background(), cfg, r2Tree().Write(t), builder.PruneOptions{})
	if err != nil {
		t.Fatal(err)
	}
	want := append(sortedKeys(orphanObjects), draftKeys...)
	slices.Sort(want)
	if got := orphanKeys(rep); !slices.Equal(got, want) {
		t.Errorf("draft gone from the tree: orphans\n got %v\nwant %v", got, want)
	}
}

func TestPrune_WithoutACredentialFailsWithTheBuildersMessage(t *testing.T) {
	for _, missing := range []string{"CLOUDFLARE_ACCOUNT_ID", "R2_ACCESS_KEY_ID", "R2_SECRET_ACCESS_KEY"} {
		t.Run(missing, func(t *testing.T) {
			store := buildtest.NewImageStore(t)
			store.Seed("og-badge.000000000000.png", buildtest.StoredObject{Data: []byte("123")})
			env := store.Env()
			delete(env, missing)
			cfg := store.ConfigureEnv(buildtest.Config(), env)
			cfg.Images.CacheDir = t.TempDir()
			tree := r2Tree()
			_, built := buildtest.Build(t, cfg, tree)
			if len(built.Errors) != 1 {
				t.Fatalf("builder errors: %v", built.Errors)
			}

			_, err := builder.Prune(context.Background(), cfg, tree.Write(t), builder.PruneOptions{Delete: true})

			if err == nil || err.Error() != built.Errors[0].Error() {
				t.Errorf("prune error = %v\nwant the builder's %q", err, built.Errors[0].Error())
			}
			if reqs := store.Requests(""); len(reqs) != 0 {
				t.Errorf("store received %d requests", len(reqs))
			}
			if entries, _ := os.ReadDir(cfg.Images.CacheDir); len(entries) != 0 {
				t.Errorf("image cache has %d entries: prune built before checking credentials", len(entries))
			}
		})
	}
}

func TestPrune_RefusesWhenTheProductionBuildFails(t *testing.T) {
	store, cfg, _ := publishedStore(t, r2Tree())
	before := store.Keys()
	// The published Explainer loses its TL;DR: its page (and so its hero)
	// drops out of the build, and its uploaded variants would look orphaned.
	broken := buildtest.ValidExplainer()
	broken.Body = "Raft splits consensus into leader election, log replication and safety.\n"
	tree := r2Tree().WithExplainer(broken)

	rep, err := builder.Prune(context.Background(), cfg, tree.Write(t), builder.PruneOptions{Delete: true})

	if err == nil || !strings.Contains(err.Error(), "tldr") {
		t.Errorf("prune error = %v, want the production build's errors (rule tldr)", err)
	}
	if len(rep.Orphans) != 0 {
		t.Errorf("reported %d orphans", len(rep.Orphans))
	}
	if got := store.Keys(); !slices.Equal(got, before) {
		t.Errorf("bucket changed:\n got %v\nwant %v", got, before)
	}
	if reqs := store.Requests(""); len(reqs) != 0 {
		t.Errorf("store received %d requests", len(reqs))
	}
}

func TestPrune_DeleteRemovesExactlyTheOrphansInBatchesOfAtMost1000(t *testing.T) {
	tree := r2Tree()
	store, cfg, referenced := publishedStore(t, tree)
	// 1,500 more orphans: the listing spans two pages and the deletes two
	// DeleteObjects requests.
	for i := range 1500 {
		store.Seed(fmt.Sprintf("heroes/old-%04d-480x270.0123456789.avif", i), buildtest.StoredObject{Data: []byte("x")})
	}

	var log bytes.Buffer
	rep, err := builder.Prune(context.Background(), cfg, tree.Write(t), builder.PruneOptions{Delete: true, Log: &log})
	if err != nil {
		t.Fatal(err)
	}

	if len(rep.Orphans) != 1503 || !rep.Deleted {
		t.Errorf("report: %d orphans, Deleted = %v; want 1503, true", len(rep.Orphans), rep.Deleted)
	}
	if got := store.Keys(); !slices.Equal(got, referenced) {
		t.Errorf("bucket after --delete holds %d keys, want exactly the %d referenced ones", len(got), len(referenced))
	}
	deletes := 0
	for _, r := range store.Requests("") {
		switch r.Op {
		case "ListObjectsV2":
		case "DeleteObjects":
			deletes++
			if n := bytes.Count(r.Body, []byte("<Object>")); n > 1000 {
				t.Errorf("DeleteObjects request names %d keys, want at most 1000", n)
			}
		default:
			t.Errorf("prune sent %s %s", r.Op, r.Key)
		}
	}
	if deletes != 2 {
		t.Errorf("%d DeleteObjects requests, want 2", deletes)
	}
	if !strings.Contains(log.String(), "deleted 1503 object(s)") {
		t.Errorf("output does not report the deletion:\n%s", log.String())
	}
}

package factpack

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"hastcode.com/site/internal/content"
)

// RefreshReport is what a Refresh found.
type RefreshReport struct {
	// Profiles is how many Repo Profiles the content tree holds.
	Profiles int
	// Changes lists every repository whose Fact Pack was rewritten, in
	// content-tree order. Empty means nothing changed and no file moved.
	Changes []Change
}

// The fetch command's exit statuses for a Refresh that ran to the end
// (cmd/fetch; .github/workflows/refresh.yml acts on them).
const (
	// ExitNothingChanged: no Fact Pack was rewritten.
	ExitNothingChanged = 0
	// ExitChanged: at least one Fact Pack was rewritten.
	ExitChanged = 3
)

// ExitStatus is the fetch command's exit status for this Refresh:
// ExitChanged when any Fact Pack was rewritten, else ExitNothingChanged.
func (r *RefreshReport) ExitStatus() int {
	if len(r.Changes) > 0 {
		return ExitChanged
	}
	return ExitNothingChanged
}

// Change is one rewritten Fact Pack and what changed in it.
type Change struct {
	// Repo is owner/name; File the Fact Pack's content-relative path.
	Repo, File string
	// Facts says what changed, one phrase each ("stars 2468 -> 2600",
	// "latest release server-0.5.0 -> server-0.6.0"); "new Fact Pack" when
	// there was none.
	Facts []string
}

// String is the Refresh's one line per changed repository:
// "apache/iggy: stars 2468 -> 2600, latest release … (explainers/…/apache-iggy.json)".
func (c Change) String() string {
	return c.Repo + ": " + strings.Join(c.Facts, ", ") + " (" + c.File + ")"
}

// Refresh re-fetches the Fact Pack of every Repo Profile under
// contentDir/explainers (GLOSSARY "Refresh") and rewrites only those whose
// facts changed, fetched_at aside, so an unchanged file keeps its bytes
// and its fetched_at. A Profile without a Fact Pack gets one. Each
// rewritten repository is logged as one Change line; the README text goes
// to CacheDir as on a Fetch. A fetch that fails (a rate limit, say) stops
// the Refresh with the error; the Fact Packs rewritten before it stand.
func (f *Fetcher) Refresh(ctx context.Context, contentDir string) (*RefreshReport, error) {
	profiles, err := repoProfiles(contentDir)
	if err != nil {
		return nil, err
	}
	report := &RefreshReport{Profiles: len(profiles)}
	for _, p := range profiles {
		file := content.FactPackFile(p.file)
		path := filepath.Join(contentDir, filepath.FromSlash(file))
		old, err := readFactPack(path)
		if err != nil {
			return report, fmt.Errorf("%s: %w (fetch it again: -repo %s -markdown %s)", file, err, p.repo, filepath.Join(contentDir, filepath.FromSlash(p.file)))
		}

		got, err := f.fetch(ctx, p.repo)
		if err != nil {
			return report, fmt.Errorf("%s: %w", p.repo, err)
		}
		facts := []string{"new Fact Pack"}
		if old != nil {
			if !changed(old, got.pack) {
				continue
			}
			facts = diff(old, got.pack)
		}
		if err := writeFactPack(path, got.pack); err != nil {
			return report, err
		}
		if got.readme != nil {
			if _, err := f.writeREADME(p.repo, got.readme); err != nil {
				return report, err
			}
		}
		change := Change{Repo: p.repo, File: file, Facts: facts}
		report.Changes = append(report.Changes, change)
		f.logf("%s", change)
		for _, n := range got.notes {
			f.logf("note: %s", n)
		}
	}
	f.logf("refresh: %d Repo Profiles checked, %d changed", report.Profiles, len(report.Changes))
	return report, nil
}

// profile is one Repo Profile found in the content tree.
type profile struct {
	file string // content-relative Markdown path, "/"-separated
	repo string
}

// repoProfiles lists the Repo Profiles under contentDir/explainers in path
// order. A content tree without explainers has none. A Markdown file whose
// frontmatter does not parse is an error naming it, as it would be in the
// build.
func repoProfiles(contentDir string) ([]profile, error) {
	root := filepath.Join(contentDir, "explainers")
	if _, err := os.Stat(root); errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	var profiles []profile
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".md" {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(contentDir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		e, err := content.ParseExplainer(rel, data)
		if err != nil {
			return fmt.Errorf("%s: %w", rel, err)
		}
		if e.Kind == content.RepoProfile {
			profiles = append(profiles, profile{file: rel, repo: e.Repo})
		}
		return nil
	})
	return profiles, err
}

// readFactPack reads an existing Fact Pack; nil when there is none.
func readFactPack(path string) (*content.FactPack, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return content.ParseFactPack(data)
}

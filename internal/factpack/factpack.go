// Package factpack fetches a Repo Profile's Fact Pack from the GitHub REST
// API and writes it beside the Profile's Markdown (GLOSSARY "Fact Pack";
// ADR-0005), and refreshes every Fact Pack under a content tree (GLOSSARY
// "Refresh"). It is the library behind cmd/fetch, built against the GitHub
// interface so tests serve recorded JSON from an in-process server.
//
// One fetch is at most ten REST calls (spec "Content model"), in this
// order. Each that answers cites, in the Fact Pack's sources, the
// github.com page a Reader opens to check it (the third column; {repo} is
// the repository's page), each page once:
//
//	GET /repos/{owner}/{name}                            repository facts          {repo}
//	GET /repos/{owner}/{name}/languages                  bytes by language         {repo}
//	GET /repos/{owner}/{name}/releases?per_page=10       latest ten releases       {repo}/releases
//	GET /repos/{owner}/{name}/contributors?per_page=10   top ten contributors      {repo}/graphs/contributors
//	GET /repos/{owner}/{name}/contributors?per_page=1    total: Link's last page   {repo}/graphs/contributors
//	GET /repos/{owner}/{name}/git/trees/{branch}         root files and dirs       {repo}
//	GET /repos/{owner}/{name}/contents/.github?ref=…     .github files (if any)    {repo}/tree/{branch}/.github
//	GET /repos/{owner}/{name}/readme                     README blob SHA and text  the README's blob page
//	GET /repos/{owner}/{name}/stargazers/count           stargazer count           {repo}
package factpack

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"hastcode.com/site/internal/content"
)

// GitHub is what a fetch needs from the API: one authenticated JSON GET.
// *github.Client implements it.
type GitHub interface {
	Get(ctx context.Context, path string, v any) (http.Header, error)
}

// Fetcher fetches and writes Fact Packs.
type Fetcher struct {
	// GitHub makes the requests.
	GitHub GitHub
	// CacheDir is the builder's cache directory; README text is written
	// under it at readmes/<owner>/<name>/<file> for drafting, never into
	// the content tree.
	CacheDir string
	// Now is the clock fetched_at records; default time.Now.
	Now func() time.Time
	// Log receives one line per file written, every note (a NOASSERTION
	// license, a fallback) and, on a Refresh, one line per changed
	// repository. Default: discarded.
	Log io.Writer
}

// Result is what one Fetch wrote.
type Result struct {
	Pack *content.FactPack
	// FactPackFile is the Fact Pack's path, beside the Markdown.
	FactPackFile string
	// READMEFile is where the README text went under CacheDir; "" when the
	// repository has no README.
	READMEFile string
	// Notes are the things the Author must look at by hand, also logged.
	Notes []string
}

// Fetch fetches repo (owner/name) and writes its Fact Pack beside
// markdownFile, named after it (content.FactPackFile), and the README text
// under CacheDir. When the Markdown already exists it must be a Repo
// Profile of that repository, so a Fact Pack never lands beside the wrong
// piece. A NOASSERTION license sets needs_review and is noted, never an
// error.
func (f *Fetcher) Fetch(ctx context.Context, repo, markdownFile string) (*Result, error) {
	if err := checkRepo(repo); err != nil {
		return nil, err
	}
	if filepath.Ext(markdownFile) != ".md" {
		return nil, fmt.Errorf("%s: a Repo Profile's Markdown file ends in .md; the Fact Pack is written beside it with .json", markdownFile)
	}
	if err := checkProfile(markdownFile, repo); err != nil {
		return nil, err
	}

	got, err := f.fetch(ctx, repo)
	if err != nil {
		return nil, err
	}
	res := &Result{Pack: got.pack, Notes: got.notes}
	res.FactPackFile = filepath.FromSlash(content.FactPackFile(filepath.ToSlash(markdownFile)))
	if err := writeFactPack(res.FactPackFile, got.pack); err != nil {
		return nil, err
	}
	f.logf("wrote %s (%s)", res.FactPackFile, summary(got.pack))
	if got.readme != nil {
		res.READMEFile, err = f.writeREADME(repo, got.readme)
		if err != nil {
			return nil, err
		}
		f.logf("wrote README text to %s", res.READMEFile)
	}
	for _, n := range got.notes {
		f.logf("note: %s", n)
	}
	return res, nil
}

// checkRepo requires owner/name: GitHub's owner (alphanumerics and hyphens)
// and repository (alphanumerics, hyphens, underscores, dots) names.
func checkRepo(repo string) error {
	owner, name, ok := strings.Cut(repo, "/")
	if !ok || owner == "" || name == "" || strings.Contains(name, "/") {
		return fmt.Errorf("repository %q must be owner/name", repo)
	}
	for _, r := range owner {
		if !isAlnum(r) && r != '-' {
			return fmt.Errorf("repository %q must be owner/name", repo)
		}
	}
	for _, r := range name {
		if !isAlnum(r) && r != '-' && r != '_' && r != '.' {
			return fmt.Errorf("repository %q must be owner/name", repo)
		}
	}
	return nil
}

func isAlnum(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9'
}

// checkProfile refuses to write beside a Markdown file that exists and is
// not a Repo Profile of repo.
func checkProfile(markdownFile, repo string) error {
	data, err := os.ReadFile(markdownFile)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	e, err := content.ParseExplainer(filepath.ToSlash(markdownFile), data)
	if err != nil {
		return fmt.Errorf("%s: %w", markdownFile, err)
	}
	if e.Kind != content.RepoProfile {
		return fmt.Errorf("%s is a %s, not a Repo Profile (kind: repo-profile); a Fact Pack belongs beside a Repo Profile", markdownFile, e.Kind)
	}
	if e.Repo != repo {
		return fmt.Errorf("%s profiles %s, not %s: fetch the repository its frontmatter names", markdownFile, e.Repo, repo)
	}
	return nil
}

// readmePath is where repo's README text goes under CacheDir.
func (f *Fetcher) readmePath(repo, name string) string {
	owner, repoName, _ := strings.Cut(repo, "/")
	return filepath.Join(f.CacheDir, "readmes", owner, repoName, path.Base(name))
}

func (f *Fetcher) writeREADME(repo string, r *readmeText) (string, error) {
	file := f.readmePath(repo, r.name)
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(file, r.text, 0o644); err != nil {
		return "", err
	}
	return file, nil
}

func (f *Fetcher) now() time.Time {
	if f.Now == nil {
		return time.Now()
	}
	return f.Now()
}

func (f *Fetcher) logf(format string, args ...any) {
	if f.Log != nil {
		fmt.Fprintf(f.Log, format+"\n", args...)
	}
}

// summary is the one-phrase description of a Fact Pack for the log.
func summary(fp *content.FactPack) string {
	latest := "no release yet"
	if fp.LatestRelease != nil {
		latest = "latest release " + fp.LatestRelease.Tag
	}
	return fmt.Sprintf("%s: %d stars, %s, %d contributors", fp.Slug, fp.Stars, latest, fp.Contributors.Total)
}

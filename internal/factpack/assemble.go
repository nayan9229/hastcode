package factpack

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/github"
)

// fetched is one repository's answers mapped to a Fact Pack.
type fetched struct {
	pack   *content.FactPack
	readme *readmeText // nil when the repository has no README
	notes  []string
}

// readmeText is the README's file name and decoded text.
type readmeText struct {
	name string
	text []byte
}

// GitHub's answer shapes, only the fields the Fact Pack needs
// (https://docs.github.com/en/rest).

type repository struct {
	Name          string    `json:"name"`
	FullName      string    `json:"full_name"`
	HTMLURL       string    `json:"html_url"`
	Description   string    `json:"description"`
	Homepage      string    `json:"homepage"`
	Topics        []string  `json:"topics"`
	License       *license  `json:"license"`
	DefaultBranch string    `json:"default_branch"`
	Language      string    `json:"language"`
	Stars         int       `json:"stargazers_count"`
	Forks         int       `json:"forks_count"`
	Subscribers   int       `json:"subscribers_count"` // watchers; watchers_count is stars
	OpenIssues    int       `json:"open_issues_count"`
	CreatedAt     time.Time `json:"created_at"`
	PushedAt      time.Time `json:"pushed_at"`
}

type license struct {
	Name   string `json:"name"`
	SPDXID string `json:"spdx_id"`
}

type release struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	HTMLURL     string    `json:"html_url"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
}

type contributor struct {
	Login         string `json:"login"`
	HTMLURL       string `json:"html_url"`
	Contributions int    `json:"contributions"`
}

type gitTree struct {
	Tree []struct {
		Path string `json:"path"`
		Type string `json:"type"` // blob, tree, commit (submodule)
	} `json:"tree"`
}

// contentEntry is one entry of the contents API's directory listing.
type contentEntry struct {
	Path string `json:"path"`
	Type string `json:"type"` // file, dir, symlink, submodule
}

type readmeAnswer struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	SHA      string `json:"sha"`
	HTMLURL  string `json:"html_url"`
	Encoding string `json:"encoding"` // base64; "none" above 1 MB
	Content  string `json:"content"`
}

// noAssertion is GitHub's SPDX id when it could not classify the license.
const noAssertion = "NOASSERTION"

// fetch makes the calls the package doc lists and maps the answers.
func (f *Fetcher) fetch(ctx context.Context, repo string) (*fetched, error) {
	out := &fetched{pack: &content.FactPack{FetchedAt: f.now().UTC().Truncate(time.Second)}}
	fp := out.pack
	base := "/repos/" + repo
	var r repository
	// cite records a source: the github.com page a Reader opens to check
	// what a call answered, once (several calls read the repository page).
	cite := func(page string) {
		if page != "" && !slices.Contains(fp.Sources, page) {
			fp.Sources = append(fp.Sources, page)
		}
	}
	// get makes one call and, when it answers, cites its page; page is
	// read after the answer is decoded, so it may use the answer.
	get := func(path string, v any, page func() string) error {
		if _, err := f.GitHub.Get(ctx, path, v); err != nil {
			return err
		}
		cite(page())
		return nil
	}
	// onGitHub is the repository's page plus suffix, e.g. "/releases".
	onGitHub := func(suffix string) func() string {
		return func() string { return r.HTMLURL + suffix }
	}

	if err := get(base, &r, onGitHub("")); err != nil {
		if errors.Is(err, github.ErrNotFound) {
			return nil, fmt.Errorf("repository %s is not on GitHub (or the token cannot see it): %w", repo, err)
		}
		return nil, err
	}
	fp.Slug, fp.Name, fp.URL = r.FullName, r.Name, r.HTMLURL
	fp.Description, fp.Homepage, fp.Topics = r.Description, r.Homepage, r.Topics
	fp.DefaultBranch, fp.Language = r.DefaultBranch, r.Language
	fp.Stars, fp.Forks, fp.Watchers, fp.OpenIssues = r.Stars, r.Forks, r.Subscribers, r.OpenIssues
	fp.CreatedAt, fp.PushedAt = r.CreatedAt.UTC(), r.PushedAt.UTC()
	if r.FullName != repo {
		out.notes = append(out.notes, fmt.Sprintf("GitHub names this repository %s (asked for %s): set repo: %s in the Profile's frontmatter", r.FullName, repo, r.FullName))
	}
	if r.License != nil && r.License.SPDXID != "" {
		fp.License = &content.License{SPDXID: r.License.SPDXID, Name: r.License.Name}
		if r.License.SPDXID == noAssertion {
			fp.License.NeedsReview = true
		} else {
			fp.License.URL = "https://spdx.org/licenses/" + r.License.SPDXID + ".html"
		}
	}

	fp.Languages = map[string]int64{}
	if err := get(base+"/languages", &fp.Languages, onGitHub("")); err != nil {
		return nil, err
	}

	var releases []release
	if err := get(base+"/releases?per_page=10", &releases, onGitHub("/releases")); err != nil {
		return nil, err
	}
	fp.Releases = []content.Release{}
	for _, rel := range releases {
		if rel.Draft {
			continue
		}
		date := rel.PublishedAt
		if date.IsZero() {
			date = rel.CreatedAt
		}
		cr := content.Release{Tag: rel.TagName, Name: rel.Name, Date: date.UTC(), URL: rel.HTMLURL, Prerelease: rel.Prerelease}
		fp.Releases = append(fp.Releases, cr)
		if fp.LatestRelease == nil && !rel.Prerelease {
			latest := cr
			fp.LatestRelease = &latest
		}
	}

	var top []contributor
	if err := get(base+"/contributors?per_page=10", &top, onGitHub("/graphs/contributors")); err != nil {
		return nil, err
	}
	fp.Contributors.Top = []content.Contributor{}
	for _, c := range top {
		fp.Contributors.Top = append(fp.Contributors.Top, content.Contributor{Login: c.Login, URL: c.HTMLURL, Contributions: c.Contributions})
	}
	var one []contributor
	countPath := base + "/contributors?per_page=1"
	header, err := f.GitHub.Get(ctx, countPath, &one)
	if err != nil {
		return nil, err
	}
	cite(onGitHub("/graphs/contributors")())
	fp.Contributors.Total = len(one)
	if last, ok := github.LastPage(header); ok {
		fp.Contributors.Total = last
	}

	var tree gitTree
	entries := []entry{}
	err = get(base+"/git/trees/"+r.DefaultBranch, &tree, onGitHub(""))
	switch {
	case errors.Is(err, github.ErrNotFound): // an empty repository has no tree
		out.notes = append(out.notes, fmt.Sprintf("%s has no files on %s yet", repo, r.DefaultBranch))
	case err != nil:
		return nil, err
	}
	hasGithubDir := false
	for _, t := range tree.Tree {
		entries = append(entries, entry{path: t.Path, dir: t.Type == "tree"})
		hasGithubDir = hasGithubDir || t.Path == ".github" && t.Type == "tree"
	}
	if hasGithubDir {
		var dir []contentEntry
		if err := get(base+"/contents/.github?ref="+r.DefaultBranch, &dir, onGitHub("/tree/"+r.DefaultBranch+"/.github")); err != nil {
			return nil, err
		}
		for _, d := range dir {
			entries = append(entries, entry{path: d.Path, dir: d.Type == "dir"})
		}
	}
	fp.Files = detectFiles(entries)

	var readme readmeAnswer
	err = get(base+"/readme", &readme, func() string { return readme.HTMLURL })
	switch {
	case errors.Is(err, github.ErrNotFound):
		out.notes = append(out.notes, repo+" has no README")
	case err != nil:
		return nil, err
	default:
		fp.README = content.READMERef{URL: readme.HTMLURL, SHA: readme.SHA}
		fp.Files.README = readme.Path
		if readme.Encoding == "base64" {
			text, err := base64.StdEncoding.DecodeString(strings.Join(strings.Fields(readme.Content), ""))
			if err != nil {
				return nil, fmt.Errorf("README of %s: decoding: %w", repo, err)
			}
			out.readme = &readmeText{name: readme.Name, text: text}
		} else {
			out.notes = append(out.notes, fmt.Sprintf("README of %s is too large for the API to return (encoding %q); read it at %s", repo, readme.Encoding, readme.HTMLURL))
		}
	}

	var count struct {
		Count int `json:"count"`
	}
	countPath = base + "/stargazers/count"
	_, err = f.GitHub.Get(ctx, countPath, &count)
	switch {
	case errors.Is(err, github.ErrNotFound):
		fp.Stargazers.Count = r.Stars
		out.notes = append(out.notes, "stargazers/count is not available for "+repo+" (404); stargazers.count is the repository's stargazers_count")
	case err != nil:
		return nil, err
	case count.Count == 0 && r.Stars > 0:
		fp.Stargazers.Count = r.Stars
		out.notes = append(out.notes, "stargazers/count answered without a count for "+repo+"; stargazers.count is the repository's stargazers_count")
	default:
		fp.Stargazers.Count = count.Count
		cite(onGitHub("")())
	}

	if fp.License != nil && fp.License.NeedsReview {
		where := "the LICENSE file"
		if fp.Files.License != "" {
			where = fp.Files.License
		}
		out.notes = append(out.notes, fmt.Sprintf("GitHub could not classify the license of %s (NOASSERTION): needs_review is set; read %s by hand before publishing", repo, where))
	}
	return out, nil
}

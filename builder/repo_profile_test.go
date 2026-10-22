package builder_test

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/internal/buildtest"
)

// Repo Profiles (ADR-0005): an Explainer of kind repo-profile whose fact
// sections render from a Fact Pack beside the Markdown.

const (
	repoProfileFile = "explainers/distributed-systems/apache-iggy.md"
	repoProfilePath = "/distributed-systems/apache-iggy/"
)

func repoProfileTree() buildtest.Tree {
	return buildtest.NewTree().WithRepoProfile(buildtest.ValidRepoProfile())
}

func TestRepoProfile_IsPublishedAtTheVerticalURLWithBylineDatesAndNote(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), repoProfileTree())
	doc := buildtest.Page(t, out, repoProfilePath)

	article := doc.Find("article.h-entry")
	if got := strings.TrimSpace(article.Find("h1.p-name").Text()); got != "Apache Iggy: a persistent message stream in Rust" {
		t.Errorf("h1 = %q", got)
	}
	if got := strings.TrimSpace(article.Find(".p-category").Text()); got != "Distributed Systems" {
		t.Errorf("Vertical = %q, want Distributed Systems", got)
	}
	if got := strings.TrimSpace(article.Find(".byline a.p-author").Text()); got != "Ada Lovelace" {
		t.Errorf("byline = %q, want the Author", got)
	}
	if got := article.Find("time.dt-published").AttrOr("datetime", ""); got != "2026-10-21T09:00:00Z" {
		t.Errorf("published = %q", got)
	}
	if got := article.Find("time.dt-updated").AttrOr("datetime", ""); got != "2026-10-21T09:00:00Z" {
		t.Errorf("last reviewed = %q", got)
	}
	if got := article.Find(".e-content blockquote").Text(); !strings.Contains(got, "Drafted with AI assistance") {
		t.Errorf("AI-assistance note = %q, want the Editorial Policy's note", got)
	}
}

func TestRepoProfile_RepoIsRequiredAndRepoAndWhoForAreRejectedOnOtherKinds(t *testing.T) {
	for name, tc := range map[string]struct {
		mutate func(e *buildtest.Explainer)
		file   string
		want   string
	}{
		"repo missing on a Repo Profile": {
			func(e *buildtest.Explainer) { e.Repo = "" },
			repoProfileFile, "repo is required for a Repo Profile",
		},
		"repo not owner/name": {
			func(e *buildtest.Explainer) { e.Repo = "iggy" },
			repoProfileFile, `repo "iggy" must be owner/name`,
		},
		"repo on a Concept": {
			func(e *buildtest.Explainer) { *e = buildtest.ValidExplainer(); e.Repo = "apache/iggy" },
			"explainers/distributed-systems/raft-leader-election.md", "repo is only for a Repo Profile",
		},
		"who_for on a Concept": {
			func(e *buildtest.Explainer) { *e = buildtest.ValidExplainer(); e.WhoFor = "Anyone" },
			"explainers/distributed-systems/raft-leader-election.md", "who_for is only for a Repo Profile",
		},
	} {
		t.Run(name, func(t *testing.T) {
			e := buildtest.ValidRepoProfile()
			tc.mutate(&e)

			out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithRepoProfile(e))

			buildtest.RequireErrorRule(t, res, tc.file, "frontmatter", tc.want)
			if buildtest.Exists(out, "/"+e.Vertical+"/"+e.Slug+"/") {
				t.Error("page was written despite the frontmatter error")
			}
		})
	}
}

const factPackFile = "explainers/distributed-systems/apache-iggy.json"

func TestRepoProfile_FactPackMissingUnparsableOrForAnotherRepositoryIsABuildError(t *testing.T) {
	for name, tc := range map[string]struct {
		tree func() buildtest.Tree
		rule string
		want string
	}{
		"missing": {
			func() buildtest.Tree { return repoProfileTree().Without(factPackFile) },
			"fact-pack", "missing Fact Pack " + factPackFile,
		},
		"unparsable": {
			func() buildtest.Tree { return repoProfileTree().With(factPackFile, []byte("{\"slug\": ")) },
			"fact-pack", "Fact Pack " + factPackFile,
		},
		"names another repository": {
			func() buildtest.Tree {
				p := buildtest.ValidRepoProfile()
				p.FactPack.Slug = "apache/kafka"
				return buildtest.NewTree().WithRepoProfile(p)
			},
			"fact-pack-repo", `Fact Pack ` + factPackFile + ` names repository "apache/kafka", but repo is "apache/iggy"`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			out, res := buildtest.Build(t, buildtest.Config(), tc.tree())

			buildtest.RequireErrorRule(t, res, repoProfileFile, tc.rule, tc.want)
			if buildtest.Exists(out, repoProfilePath) {
				t.Error("page was written without a usable Fact Pack")
			}
		})
	}
}

func TestRepoProfile_FactPackFetchedOver30DaysAgoIsAWarningNotAnError(t *testing.T) {
	// ValidFactPack was fetched 2026-10-19.
	cfg := buildtest.Config()
	cfg.Now = time.Date(2026, 11, 20, 12, 0, 0, 0, time.UTC) // 32 days later

	out, res := buildtest.Build(t, cfg, repoProfileTree())

	buildtest.RequireWarningRule(t, res, repoProfileFile, "fact-pack-stale", "Fact Pack "+factPackFile+" was fetched 2026-10-19, 32 days ago")
	if len(res.Errors) != 0 || !buildtest.Exists(out, repoProfilePath) {
		t.Errorf("a stale Fact Pack must not fail the build; errors:\n%s", res)
	}

	cfg.Now = time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC) // 2 days later
	_, res = buildtest.Build(t, cfg, repoProfileTree())
	buildtest.RequireNoProblemRule(t, res, repoProfileFile, "fact-pack-stale")
}

// fact returns the text of one header-card fact by its data-fact name.
func fact(card *goquery.Selection, name string) string {
	return strings.Join(strings.Fields(card.Find(`[data-fact="`+name+`"] dd`).Text()), " ")
}

func TestRepoProfile_HeaderCardShowsWhatWhoStarsLicenseLanguageReleaseAndLastCommit(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), repoProfileTree())
	doc := buildtest.Page(t, out, repoProfilePath)

	card := doc.Find("article.h-entry > header .repo-profile-card")
	if card.Length() != 1 {
		t.Fatalf("found %d header cards in the article header, want 1", card.Length())
	}
	repo := card.Find(`a[href="https://github.com/apache/iggy"]`)
	if strings.TrimSpace(repo.Text()) != "apache/iggy" {
		t.Errorf("card names the repository %q, want apache/iggy linking its GitHub page", repo.Text())
	}
	for name, want := range map[string]string{
		"what":     "Iggy is a persistent message streaming platform in Rust that trades Kafka's ecosystem for a smaller, faster single binary.",
		"who-for":  "Teams that need Kafka-style streams without running a JVM cluster.",
		"stars":    "2,468",
		"license":  "Apache-2.0",
		"language": "Rust",
	} {
		if got := fact(card, name); got != want {
			t.Errorf("%s = %q, want %q", name, got, want)
		}
	}
	release := card.Find(`[data-fact="latest-release"] dd`)
	if got := release.Find(`a[href="https://github.com/apache/iggy/releases/tag/server-0.5.0"]`).Text(); got != "server-0.5.0" {
		t.Errorf("latest release link = %q, want the tag linking the release page", got)
	}
	if got := release.Find("time").AttrOr("datetime", ""); got != "2026-10-10T12:00:00Z" {
		t.Errorf("latest release date = %q", got)
	}
	commit := card.Find(`[data-fact="last-commit"] dd time`)
	if got, text := commit.AttrOr("datetime", ""), strings.TrimSpace(commit.Text()); got != "2026-10-18T14:03:00Z" || text != "18 October 2026" {
		t.Errorf("last commit = %q (%q), want 2026-10-18T14:03:00Z shown as 18 October 2026", got, text)
	}
	if n := doc.Find("article.h-entry > header .p-summary").Length(); n != 1 {
		t.Errorf("the header carries the description %d times, want once (the card's \"what it is\" is the p-summary)", n)
	}

	// The card is a fixed-height box so nothing below it moves as it paints.
	css := stylesheet(t, out, doc)
	h := declaration(findRule(css, func(p string) bool { return p == ".repo-profile-card" }), "height")
	if m := regexp.MustCompile(`^var\((--[a-z0-9-]+)\)$`).FindStringSubmatch(h); m != nil {
		h = declaration(allRootVars(css), m[1])
	}
	if !regexp.MustCompile(`^\d+px$`).MatchString(h) {
		t.Errorf(".repo-profile-card height = %q, want a fixed pixel height (no layout shift)", h)
	}
}

func TestRepoProfile_HeaderCardOmitsWhoForWhenNotGiven(t *testing.T) {
	p := buildtest.ValidRepoProfile()
	p.WhoFor = ""
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithRepoProfile(p))

	card := buildtest.Page(t, out, repoProfilePath).Find(".repo-profile-card")
	if card.Length() != 1 || card.Find(`[data-fact="who-for"]`).Length() != 0 {
		t.Errorf("card without who_for: %d cards, %d who-for rows; want 1 card and no row", card.Length(), card.Find(`[data-fact="who-for"]`).Length())
	}
	if got := fact(card, "stars"); got != "2,468" {
		t.Errorf("stars = %q", got)
	}
}

// secondRepoProfile is a second valid Repo Profile (php/frankenphp) in
// another Vertical, published after ValidRepoProfile.
func secondRepoProfile() buildtest.Explainer {
	p := buildtest.ValidRepoProfile()
	p.Vertical = "languages-runtimes"
	p.Slug = "php-frankenphp"
	p.Title = "FrankenPHP: PHP on a Go server"
	p.Repo = "php/frankenphp"
	p.Published = "2026-10-22T09:00:00Z"
	fp := buildtest.ValidFactPack()
	fp.Slug = "php/frankenphp"
	fp.Name = "frankenphp"
	fp.URL = "https://github.com/php/frankenphp"
	p.FactPack = &fp
	return p
}

func TestRepoProfile_KindPageListsProfilesAcrossVerticalsNewestFirstWithBadge(t *testing.T) {
	tree := buildtest.ValidTree().WithRepoProfile(buildtest.ValidRepoProfile()).WithRepoProfile(secondRepoProfile())
	out := buildtest.MustBuild(t, buildtest.Config(), tree)

	kind := buildtest.Page(t, out, "/explainers/repo-profile/")
	if got := kind.Find("h1").Text(); got != "Repo Profiles" {
		t.Errorf("kind page h1 = %q, want Repo Profiles", got)
	}
	if got := kind.Find(".listing > .p-summary").Text(); got != "Explainers that evaluate one open-source repository: what it is, how it is built, and whether you should adopt it." {
		t.Errorf("kind description = %q, want the glossary sentence", got)
	}
	got := listedURLs(kind.Find("main"))
	want := []string{"/languages-runtimes/php-frankenphp/", repoProfilePath}
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("kind page lists %v, want %v (newest first, no Concept)", got, want)
	}
	if got := strings.TrimSpace(kind.Find("main .h-entry .badge").First().Text()); got != "Repo Profile" {
		t.Errorf("badge = %q, want Repo Profile", got)
	}

	hub := buildtest.Page(t, out, "/distributed-systems/")
	item := hub.Find(`main .h-entry:has(a[href="` + repoProfilePath + `"])`)
	if got := strings.TrimSpace(item.Find(".badge").Text()); got != "Repo Profile" {
		t.Errorf("hub badge = %q, want Repo Profile", got)
	}
	if n := buildtest.Page(t, out, "/").Find(`main a[href="/explainers/repo-profile/"]`).Length(); n != 1 {
		t.Errorf("home links the Repo Profiles kind page %d times, want 1", n)
	}
}

func TestRepoProfile_AppearsInItsHubTheHomePageFeedsAndSearch(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree().WithRepoProfile(buildtest.ValidRepoProfile()))

	for _, page := range []string{"/distributed-systems/", "/"} {
		if n := buildtest.Page(t, out, page).Find(`main a.u-url[href="` + repoProfilePath + `"]`).Length(); n != 1 {
			t.Errorf("%s lists the Repo Profile %d times, want 1", page, n)
		}
	}
	for _, feed := range []string{"/feed.xml", "/distributed-systems/feed.xml"} {
		if got := itemLinks(readFeed(t, out, feed)); len(got) == 0 || got[0] != siteURL+repoProfilePath {
			t.Errorf("%s items = %v, want the Repo Profile first (it is the newest piece)", feed, got)
		}
	}
	if text := buildtest.SearchIndex(t, out)[repoProfilePath]; !strings.Contains(text, "Apache Iggy") {
		t.Errorf("search index for %s = %q, want the Profile's title", repoProfilePath, text)
	}
}

// The Fact Pack's key names are the spec's ("Content model"), including
// the exchange-JSON names where they overlap, and `needs_review` on the
// license. This literal is the independent spelling the builder must
// read; a key it does not know fails the build.
func TestRepoProfile_FactPackKeysFollowTheSpecAndAbsentFactsAreNeverInvented(t *testing.T) {
	const spelled = `{
  "fetched_at": "2026-10-19T06:00:00Z",
  "slug": "apache/iggy",
  "name": "iggy",
  "url": "https://github.com/apache/iggy",
  "description": "Persistent message streaming platform written in Rust",
  "homepage": "https://iggy.apache.org",
  "topics": ["streaming"],
  "license": {"spdx_id": "NOASSERTION", "name": "Other", "url": "", "needs_review": true},
  "default_branch": "main",
  "language": "Go",
  "languages": {"Go": 1000},
  "stars": 1234567,
  "forks": 1,
  "watchers": 2,
  "open_issues": 3,
  "created_at": "2023-07-24T10:00:00Z",
  "pushed_at": "2026-10-18T14:03:00Z",
  "latest_release": {"tag": "v1.2.0", "name": "1.2.0", "date": "2026-10-01T00:00:00Z", "url": "https://github.com/apache/iggy/releases/tag/v1.2.0", "prerelease": false},
  "releases": [],
  "contributors": {"total": 1, "top": [{"login": "spetz", "url": "https://github.com/spetz", "contributions": 1}]},
  "files": {"readme": "README.md", "license": "", "security": "", "contributing": "", "code_of_conduct": "", "ci": [], "manifests": ["Cargo.toml"], "dockerfile": "", "docs": "", "examples": ""},
  "readme": {"url": "https://github.com/apache/iggy/blob/main/README.md", "sha": "abc"},
  "stargazers": {"count": 1234567},
  "sources": ["https://api.github.com/repos/apache/iggy"]
}`
	out := buildtest.MustBuild(t, buildtest.Config(), repoProfileTree().With(factPackFile, []byte(spelled)))
	card := buildtest.Page(t, out, repoProfilePath).Find(".repo-profile-card")
	for name, want := range map[string]string{
		"stars":          "1,234,567",
		"language":       "Go",
		"license":        "not classified", // NOASSERTION is for the Author, never shown as a license
		"latest-release": "v1.2.0 1 October 2026",
	} {
		if got := fact(card, name); got != want {
			t.Errorf("%s = %q, want %q", name, got, want)
		}
	}

	// Facts GitHub does not report stay absent; nothing is made up.
	const absent = `{"fetched_at": "2026-10-19T06:00:00Z", "slug": "apache/iggy", "url": "https://github.com/apache/iggy",
  "license": null, "language": "", "latest_release": null, "pushed_at": "2026-10-18T14:03:00Z"}`
	out = buildtest.MustBuild(t, buildtest.Config(), repoProfileTree().With(factPackFile, []byte(absent)))
	card = buildtest.Page(t, out, repoProfilePath).Find(".repo-profile-card")
	for name, want := range map[string]string{"license": "not found", "language": "not found", "latest-release": "none yet", "stars": "0"} {
		if got := fact(card, name); got != want {
			t.Errorf("absent %s = %q, want %q", name, got, want)
		}
	}
}

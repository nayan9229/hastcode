package builder_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/internal/buildtest"
)

// A Repo Profile's fact sections (spec "Page composition and gates"):
// Signals after Architecture and Changelog after Alternatives, rendered
// from the Fact Pack, never written by the Author.

// headings returns the text of every H2 in sel, in document order.
func headings(sel *goquery.Selection) []string {
	var out []string
	sel.Find("h2").Each(func(_ int, h *goquery.Selection) { out = append(out, strings.TrimSpace(h.Text())) })
	return out
}

// signal returns the text of one Signals row by its data-signal name.
func signal(doc *goquery.Document, name string) string {
	return strings.Join(strings.Fields(doc.Find(`.signals [data-signal="`+name+`"] dd`).Text()), " ")
}

func TestRepoProfile_SignalsFollowsArchitectureWithFactsFromTheFactPack(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), repoProfileTree())
	doc := buildtest.Page(t, out, repoProfilePath)

	got := strings.Join(headings(doc.Find("article .e-content")), " | ")
	if !strings.Contains(got, "Architecture | Signals | Use it if") {
		t.Errorf("headings = %s; want Signals right after Architecture", got)
	}
	if n := doc.Find(".signals").Length(); n != 1 {
		t.Fatalf("found %d Signals sections, want 1", n)
	}
	for name, want := range map[string]string{
		"last-commit":        "18 October 2026",
		"releases-last-year": "3",
		"latest-release":     "server-0.5.0 10 October 2026",
		"contributors":       "61",
		"top-contributors":   "spetz holds 81% of the 2,220 contributions by the top 2",
		"open-issues":        "57",
		"security-policy":    ".github/SECURITY.md",
		"license":            "Apache-2.0",
		"stars":              "2,468",
		"forks":              "135",
		"watchers":           "42",
		"documentation":      "README.md, CONTRIBUTING.md, CODE_OF_CONDUCT.md, docs/, examples/",
	} {
		if got := signal(doc, name); got != want {
			t.Errorf("signal %s = %q, want %q", name, got, want)
		}
	}
	if got := doc.Find(`.signals [data-signal="last-commit"] time`).AttrOr("datetime", ""); got != "2026-10-18T14:03:00Z" {
		t.Errorf("last commit datetime = %q", got)
	}
	if got := doc.Find(`.signals [data-signal="latest-release"] a`).AttrOr("href", ""); got != "https://github.com/apache/iggy/releases/tag/server-0.5.0" {
		t.Errorf("latest release links %q, want its release page", got)
	}
	// Facts only: no score, no grade.
	if text := strings.ToLower(doc.Find(".signals").Text()); strings.Contains(text, "score") || strings.Contains(text, "grade") {
		t.Errorf("Signals shows a score or grade: %q", text)
	}
}

func TestRepoProfile_SignalsShowsAbsentFactsAsNotFoundNeverInvented(t *testing.T) {
	// A Fact Pack with only what the builder requires: no releases, no
	// contributors, no files, no license, no last commit.
	const bare = `{"fetched_at": "2026-10-19T06:00:00Z", "slug": "apache/iggy", "url": "https://github.com/apache/iggy"}`
	out := buildtest.MustBuild(t, buildtest.Config(), repoProfileTree().With(factPackFile, []byte(bare)))
	doc := buildtest.Page(t, out, repoProfilePath)

	for _, name := range []string{"last-commit", "releases-last-year", "contributors", "top-contributors", "security-policy", "license", "documentation"} {
		if got := signal(doc, name); got != "not found" {
			t.Errorf("absent %s = %q, want \"not found\"", name, got)
		}
	}
	if got := signal(doc, "latest-release"); got != "none yet" {
		t.Errorf("absent latest release = %q, want \"none yet\"", got)
	}
	if n := doc.Find(".signals time").Length(); n != 0 {
		t.Errorf("Signals renders %d dates for a Fact Pack without any", n)
	}
}

func TestRepoProfile_ChangelogFollowsAlternativesWithReleasesNewestFirstPrereleasesMarked(t *testing.T) {
	p := buildtest.ValidRepoProfile()
	// Out of date order on purpose: the page sorts newest first.
	p.FactPack.Releases = []buildtest.Release{
		{Tag: "server-0.5.0", Name: "Server 0.5.0", Date: "2026-10-10T12:00:00Z", URL: "https://github.com/apache/iggy/releases/tag/server-0.5.0"},
		{Tag: "server-0.6.0-rc1", Name: "Server 0.6.0 RC1", Date: "2026-10-16T12:00:00Z", URL: "https://github.com/apache/iggy/releases/tag/server-0.6.0-rc1", Prerelease: true},
		{Tag: "server-0.4.0", Name: "Server 0.4.0", Date: "2026-06-02T12:00:00Z", URL: "https://github.com/apache/iggy/releases/tag/server-0.4.0"},
	}
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithRepoProfile(p))
	doc := buildtest.Page(t, out, repoProfilePath)

	got := strings.Join(headings(doc.Find("article .e-content")), " | ")
	if !strings.Contains(got, "Alternatives | Changelog | FAQ") {
		t.Errorf("headings = %s; want Changelog right after Alternatives", got)
	}
	items := doc.Find(".changelog li")
	if items.Length() != 3 {
		t.Fatalf("Changelog lists %d releases, want 3", items.Length())
	}
	want := []struct{ tag, name, date, url string }{
		{"server-0.6.0-rc1", "Server 0.6.0 RC1", "2026-10-16T12:00:00Z", "https://github.com/apache/iggy/releases/tag/server-0.6.0-rc1"},
		{"server-0.5.0", "Server 0.5.0", "2026-10-10T12:00:00Z", "https://github.com/apache/iggy/releases/tag/server-0.5.0"},
		{"server-0.4.0", "Server 0.4.0", "2026-06-02T12:00:00Z", "https://github.com/apache/iggy/releases/tag/server-0.4.0"},
	}
	items.Each(func(i int, li *goquery.Selection) {
		w := want[i]
		a := li.Find("a")
		if a.Text() != w.tag || a.AttrOr("href", "") != w.url {
			t.Errorf("release %d links %q to %q, want %q to %q", i, a.Text(), a.AttrOr("href", ""), w.tag, w.url)
		}
		if !strings.Contains(li.Text(), w.name) {
			t.Errorf("release %d = %q, want its name %q", i, li.Text(), w.name)
		}
		if got := li.Find("time").AttrOr("datetime", ""); got != w.date {
			t.Errorf("release %d date = %q, want %q", i, got, w.date)
		}
		_, marked := li.Attr("data-prerelease")
		if marked != (i == 0) || strings.Contains(li.Text(), "pre-release") != (i == 0) {
			t.Errorf("release %d pre-release mark = %v (text %q); only the RC is a pre-release", i, marked, li.Text())
		}
	})
}

func TestRepoProfile_ChangelogWithoutReleasesSaysSo(t *testing.T) {
	none := buildtest.ValidRepoProfile()
	none.FactPack.Releases = []buildtest.Release{}
	none.FactPack.LatestRelease = nil
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithRepoProfile(none))
	log := buildtest.Page(t, out, repoProfilePath).Find(".changelog")
	if log.Find("li").Length() != 0 || !strings.Contains(log.Text(), "none yet") {
		t.Errorf("Changelog without releases = %q, want no list and \"none yet\"", strings.TrimSpace(log.Text()))
	}

	const bare = `{"fetched_at": "2026-10-19T06:00:00Z", "slug": "apache/iggy", "url": "https://github.com/apache/iggy"}`
	out = buildtest.MustBuild(t, buildtest.Config(), repoProfileTree().With(factPackFile, []byte(bare)))
	log = buildtest.Page(t, out, repoProfilePath).Find(".changelog")
	if log.Find("li").Length() != 0 || !strings.Contains(log.Text(), "not found") {
		t.Errorf("Changelog without a release list = %q, want \"not found\"", strings.TrimSpace(log.Text()))
	}
}

func TestRepoProfile_FactPackSourcesJoinTheSourcesSectionOncePerURL(t *testing.T) {
	p := buildtest.ValidRepoProfile()
	// The body's ## Sources already links the repository page; the
	// frontmatter cites the releases page; the Fact Pack's pages repeat both
	// (one with a trailing slash) and add the contributors graph.
	p.Sources = []buildtest.Source{{Title: "Release notes", URL: "https://github.com/apache/iggy/releases"}}
	p.FactPack.Sources = []string{
		"https://github.com/apache/iggy/",
		"https://github.com/apache/iggy/releases",
		"https://github.com/apache/iggy/graphs/contributors",
	}
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithRepoProfile(p))
	sources := buildtest.Page(t, out, repoProfilePath).Find("article .sources")

	for href, text := range map[string]string{
		"https://github.com/apache/iggy":                     "apache/iggy on GitHub", // the body link, not repeated
		"https://github.com/apache/iggy/releases":            "Release notes",         // the frontmatter title wins
		"https://github.com/apache/iggy/graphs/contributors": "https://github.com/apache/iggy/graphs/contributors",
	} {
		links := sources.Find(`a[href="` + href + `"]`)
		if links.Length() != 1 {
			t.Errorf("Sources links %s %d times, want 1", href, links.Length())
		}
		if got := links.Text(); got != text {
			t.Errorf("Sources link text for %s = %q, want %q", href, got, text)
		}
	}
	if n := sources.Find(`a[href="https://github.com/apache/iggy/"]`).Length(); n != 0 {
		t.Errorf("Sources repeats the repository with a trailing slash (%d links)", n)
	}
	feed := readFeed(t, out, "/feed.xml")
	if len(feed.Channel.Items) != 1 || !strings.Contains(feed.Channel.Items[0].Content, `href="https://github.com/apache/iggy/graphs/contributors"`) {
		t.Errorf("the full-text feed item lacks the Fact Pack source")
	}
}

func TestRepoProfileJSONLD_AddsSoftwareSourceCodeBesideTheArticleAndNoFAQPage(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), repoProfileTree())
	doc := buildtest.Page(t, out, repoProfilePath)

	if buildtest.JSONLDOfType(t, doc, "Article") == nil || buildtest.JSONLDOfType(t, doc, "BreadcrumbList") == nil {
		t.Error("the Article graph every Explainer emits is missing")
	}
	code := buildtest.JSONLDOfType(t, doc, "SoftwareSourceCode")
	if code == nil {
		t.Fatal("no SoftwareSourceCode JSON-LD block")
	}
	for key, want := range map[string]string{
		"@context":            "https://schema.org",
		"name":                "apache/iggy",
		"codeRepository":      "https://github.com/apache/iggy",
		"programmingLanguage": "Rust",
		"license":             "https://spdx.org/licenses/Apache-2.0.html",
		"description":         "Persistent message streaming platform written in Rust",
	} {
		if got := code[key]; got != want {
			t.Errorf("SoftwareSourceCode %s = %#v, want %q", key, got, want)
		}
	}

	if faq := buildtest.JSONLDOfType(t, doc, "FAQPage"); faq != nil {
		t.Errorf("FAQPage JSON-LD = %#v; the Author dropped it (Google retired FAQ results); the FAQ stays visible prose only", faq)
	}
	if n := doc.Find("article .e-content h3").Length(); n != 5 {
		t.Errorf("FAQ renders %d questions on the page, want the fixture's 5", n)
	}
}

func TestRepoProfileJSONLD_OmitsFactsTheFactPackLacks(t *testing.T) {
	const bare = `{"fetched_at": "2026-10-19T06:00:00Z", "slug": "apache/iggy", "url": "https://github.com/apache/iggy",
  "license": {"spdx_id": "NOASSERTION", "name": "Other", "url": "", "needs_review": true}}`
	out := buildtest.MustBuild(t, buildtest.Config(), repoProfileTree().With(factPackFile, []byte(bare)))
	code := buildtest.JSONLDOfType(t, buildtest.Page(t, out, repoProfilePath), "SoftwareSourceCode")
	if code == nil {
		t.Fatal("no SoftwareSourceCode JSON-LD block")
	}
	if code["name"] != "apache/iggy" || code["codeRepository"] != "https://github.com/apache/iggy" {
		t.Errorf("SoftwareSourceCode = %#v, want the repository name and URL", code)
	}
	for _, key := range []string{"programmingLanguage", "license", "description"} {
		if got, ok := code[key]; ok {
			t.Errorf("SoftwareSourceCode %s = %#v, want the key absent (never invented)", key, got)
		}
	}
}

// componentRules returns the body of every CSS rule (at any nesting, so
// rules inside @media count) whose prelude contains selector.
func componentRules(css, selector string) []string {
	css = regexp.MustCompile(`(?s)/\*.*?\*/`).ReplaceAllString(css, "")
	var bodies []string
	for _, m := range regexp.MustCompile(`([^{}]+)\{([^{}]*)\}`).FindAllStringSubmatch(css, -1) {
		if strings.Contains(m[1], selector) {
			bodies = append(bodies, m[2])
		}
	}
	return bodies
}

// The header card, Signals and Changelog paint with the brand tokens only
// (the palette rules TestStylesheet_DefinesBrandPaletteForLightDarkAndOverride
// verifies), so both themes apply; a literal colour would be wrong in one.
func TestRepoProfile_HeaderCardSignalsAndChangelogUseBrandTokensOnly(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), repoProfileTree())
	css := stylesheet(t, out, buildtest.Page(t, out, repoProfilePath))

	literal := regexp.MustCompile(`#[0-9a-f]{3,8}\b|\b(rgba?|hsla?|oklch|color)\(`)
	token := regexp.MustCompile(`var\(--(bg|fg|accent|link|muted|border|surface)\)`)
	for _, sel := range []string{".repo-profile-card", ".signals", ".changelog"} {
		rules := componentRules(css, sel)
		if len(rules) == 0 {
			t.Errorf("no %s rules in the stylesheet", sel)
			continue
		}
		all := strings.Join(rules, "\n")
		if m := literal.FindString(all); m != "" {
			t.Errorf("%s rules hard-code a colour (%q); use the brand tokens", sel, m)
		}
		if !token.MatchString(all) {
			t.Errorf("%s rules use no brand colour token", sel)
		}
	}
}

func TestRepoProfile_FixtureProfileStaysUnderThe50KBGate(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), repoProfileTree()) // MustBuild fails on html-size
	if size := buildtest.FileSize(t, out, repoProfilePath); size >= 51200 {
		t.Errorf("the fixture Repo Profile is %d bytes, want under the 51,200-byte gate", size)
	}
}

func TestRepoProfile_SignalsCountsReleasesInTheYearBeforeTheFetch(t *testing.T) {
	// ValidFactPack was fetched 2026-10-19 with releases dated 2026-10-16,
	// 2026-10-10 and 2026-06-02.
	older := buildtest.ValidRepoProfile()
	older.FactPack.Releases = append(older.FactPack.Releases, buildtest.Release{
		Tag: "server-0.3.0", Name: "Server 0.3.0", Date: "2025-10-01T12:00:00Z", URL: "https://github.com/apache/iggy/releases/tag/server-0.3.0",
	})
	full := buildtest.ValidRepoProfile()
	full.FactPack.Releases = nil
	for i := range 10 {
		full.FactPack.Releases = append(full.FactPack.Releases, buildtest.Release{
			Tag: "v0." + string(rune('9'-i)), Date: "2026-09-01T12:00:00Z", URL: "https://github.com/apache/iggy/releases/tag/v0." + string(rune('9'-i)),
		})
	}
	empty := buildtest.ValidRepoProfile()
	empty.FactPack.Releases = []buildtest.Release{}
	empty.FactPack.LatestRelease = nil

	for name, tc := range map[string]struct {
		profile buildtest.Explainer
		want    string
	}{
		"a release older than a year is not counted": {older, "3"},
		"a full list within the year is a floor":     {full, "10 or more"},
		"no releases":                                {empty, "0"},
	} {
		t.Run(name, func(t *testing.T) {
			out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithRepoProfile(tc.profile))
			if got := signal(buildtest.Page(t, out, repoProfilePath), "releases-last-year"); got != tc.want {
				t.Errorf("releases in the last 12 months = %q, want %q", got, tc.want)
			}
		})
	}
}

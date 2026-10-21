package builder_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildtest"
)

const searchPath = "/search/"

// pagefindEntry is the part of Pagefind's pagefind-entry.json the tests
// read: how many pages each language index holds.
type pagefindEntry struct {
	Version   string `json:"version"`
	Languages map[string]struct {
		PageCount int `json:"page_count"`
	} `json:"languages"`
}

func readPagefindEntry(t *testing.T, outDir string) pagefindEntry {
	t.Helper()
	data, err := os.ReadFile(buildtest.StaticPath(outDir, "/pagefind/pagefind-entry.json"))
	if err != nil {
		t.Fatalf("search index not in output: %v", err)
	}
	var entry pagefindEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		t.Fatalf("pagefind-entry.json: %v", err)
	}
	return entry
}

func TestSearch_ProductionBuildIndexesOnlyExplainersAndBriefs(t *testing.T) {
	tree := buildtest.ValidTree().WithBrief(buildtest.ValidBrief())
	out := buildtest.MustBuild(t, buildtest.Config(), tree)

	entry := readPagefindEntry(t, out)
	if got := entry.Languages["en"].PageCount; got != 2 {
		t.Errorf("English index holds %d pages, want 2 (the Explainer and the Brief; listings, static and search pages stay out)", got)
	}

	// The search page points the lazy loader at the bundle the index
	// step wrote.
	form := buildtest.Page(t, out, searchPath).Find("form[data-pagefind-bundle]")
	bundle := form.AttrOr("data-pagefind-bundle", "")
	if bundle != "/pagefind/" {
		t.Fatalf("search form data-pagefind-bundle = %q, want /pagefind/", bundle)
	}
	for _, file := range []string{"pagefind.js", "pagefind-modular-ui.js"} {
		if !buildtest.Exists(out, bundle+file) {
			t.Errorf("search page loads %s%s, which is not in the output", bundle, file)
		}
	}
}

// The index holds what a Reader would search for (title, summary, body)
// and none of the chrome around it: header, footer, breadcrumbs, ad slots,
// the newsletter box, diagrams (SVG labels) and the Sources list.
func TestSearch_IndexHoldsTheArticleBodyWithoutChrome(t *testing.T) {
	tree := buildtest.ValidTree().WithBrief(buildtest.ValidBrief())
	index := buildtest.SearchIndex(t, buildtest.MustBuild(t, buildtest.Config(), tree))

	explainer, ok := index[raftURL]
	if !ok {
		t.Fatalf("Explainer %s not in the index; indexed: %v", raftURL, keys(index))
	}
	for _, want := range []string{"How Raft elects a leader", "Each term has at most one leader.", "150–300 ms"} {
		if !strings.Contains(explainer, want) {
			t.Errorf("indexed Explainer text lacks %q:\n%s", want, explainer)
		}
	}
	chrome := []string{
		"Skip to content", "How real systems are built", "Editorial Policy", // header, footer
		"Get new Explainers by email", // newsletter box
		"A Raft election",             // diagram caption
		"Understandable Consensus",    // Sources list
	}
	for url, text := range index {
		for _, unwanted := range chrome {
			if strings.Contains(text, unwanted) {
				t.Errorf("indexed text of %s contains %q (chrome, not content):\n%s", url, unwanted, text)
			}
		}
	}
	if _, ok := index["/briefs/"+buildtest.ValidBrief().Slug+"/"]; !ok {
		t.Errorf("Brief not in the index; indexed: %v", keys(index))
	}
}

func keys(m map[string]string) []string {
	var ks []string
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

func TestSearch_ProductionBuildWithoutTheIndexFails(t *testing.T) {
	cfg := buildtest.Config()
	cfg.SkipSearchIndex = true
	out, res := buildtest.Build(t, cfg, buildtest.ValidTree())

	buildtest.RequireErrorRule(t, res, "", "search-index", "production builds must build the search index")
	if buildtest.Exists(out, "/pagefind/pagefind-entry.json") {
		t.Errorf("index written although the post-step was skipped")
	}
}

func TestSearch_PreviewBuildMaySkipTheIndex(t *testing.T) {
	cfg := buildtest.Config()
	cfg.Mode = builder.Preview
	cfg.SkipSearchIndex = true
	out, res := buildtest.Build(t, cfg, buildtest.ValidTree())

	if res.Failed() {
		t.Fatalf("preview build without the index failed:\n%s", res)
	}
	buildtest.RequireWarningRule(t, res, "", "search-index", "skipped")
	if !buildtest.Exists(out, searchPath) {
		t.Errorf("search page missing from a preview build")
	}
	if buildtest.Exists(out, "/pagefind/pagefind-entry.json") {
		t.Errorf("index written although the post-step was skipped")
	}
}

func TestSearchPage_IsNoindexAndExplainsItselfWithoutJavaScript(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	doc := buildtest.Page(t, out, searchPath)

	if got := doc.Find(`meta[name="robots"]`).AttrOr("content", ""); got != "noindex" {
		t.Errorf("robots = %q, want noindex", got)
	}
	if got := doc.Find("h1").Text(); got != "Search" {
		t.Errorf("h1 = %q, want Search", got)
	}

	input := doc.Find(`form[role="search"] input[type="search"]`)
	if input.Length() != 1 {
		t.Fatalf("found %d search inputs in a role=search form, want 1", input.Length())
	}
	if _, hidden := input.Attr("hidden"); hidden || input.Closest("[hidden], noscript").Length() > 0 {
		t.Errorf("search input must be visible before any script runs")
	}
	if id := input.AttrOr("id", ""); id == "" || doc.Find(`label[for="`+id+`"]`).Length() != 1 {
		t.Errorf("search input %q needs a <label for>", id)
	}

	if note := doc.Find("noscript").Text(); !strings.Contains(note, "JavaScript") {
		t.Errorf("noscript = %q, want a note that search needs JavaScript", note)
	}
	// Without JS the Reader can still browse: every Vertical hub and the
	// Briefs are linked from the page itself.
	for _, href := range []string{"/distributed-systems/", "/briefs/"} {
		if doc.Find(`main a[href="`+href+`"]`).Not(`nav a`).Length() == 0 {
			t.Errorf("search page does not link %s for browsing", href)
		}
	}
}

// Pagefind loads only on interaction: the page itself ships one small
// inline loader and no script file.
func TestSearchPage_LoadsPagefindLazily(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	doc := buildtest.Page(t, out, searchPath)

	var loader string
	doc.Find("script").Each(func(_ int, s *goquery.Selection) {
		if src, ok := s.Attr("src"); ok {
			t.Errorf("search page loads %s up front; Pagefind must load on first focus", src)
		}
		if strings.Contains(s.Text(), "data-pagefind-bundle") {
			loader = s.Text()
		}
	})
	if loader == "" {
		t.Fatal("no inline search loader reading data-pagefind-bundle on the page")
	}
	if n := len(loader); n > 1536 {
		t.Errorf("search loader is %d bytes, want ≤ 1536", n)
	}
	if !strings.Contains(loader, `"focus"`) {
		t.Errorf("search loader should start loading Pagefind on focus")
	}
}

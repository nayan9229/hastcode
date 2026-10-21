package builder_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"hastcode.com/site/internal/buildtest"
)

// Spec: the build fails on a page lacking JSON-LD. writePage gates every
// indexable page (rule jsonld-missing); this proves every page type the
// build writes passes it.
func TestJSONLD_EveryIndexablePageCarriesAtLeastOneBlock(t *testing.T) {
	tree := buildtest.ValidTree().WithExplainer(buildtest.ValidComparison()).WithBrief(buildtest.ValidBrief())
	out := buildtest.MustBuild(t, buildtest.Config(), tree)

	checked := 0
	for _, u := range buildtest.URLPaths(t, out) {
		if !strings.HasSuffix(u, "/") {
			continue
		}
		doc := buildtest.Page(t, out, u)
		if strings.Contains(doc.Find(`meta[name="robots"]`).AttrOr("content", ""), "noindex") {
			continue
		}
		checked++
		if doc.Find(`script[type="application/ld+json"]`).Length() == 0 {
			t.Errorf("indexable page %s has no JSON-LD", u)
		}
	}
	if checked < 10 {
		t.Fatalf("checked only %d indexable pages; the fixture should produce home, hubs, kind pages, articles, Author and static pages", checked)
	}
}

func TestExplainerJSONLD_ArticleCarriesHeadlineDescriptionAndFrontmatterDates(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	article := buildtest.JSONLDOfType(t, buildtest.Page(t, out, raftURL), "Article")
	if article == nil {
		t.Fatal("no Article JSON-LD block")
	}

	for key, want := range map[string]string{
		"@context":      "https://schema.org",
		"headline":      "How Raft elects a leader",
		"description":   "Raft keeps a replicated log consistent by electing one leader per term with randomised timeouts.",
		"datePublished": "2026-10-01T09:00:00+05:30",
		"dateModified":  "2026-10-15T18:30:00Z",
	} {
		if got := article[key]; got != want {
			t.Errorf("Article %s = %#v, want %q", key, got, want)
		}
	}
}

func TestExplainerJSONLD_ArticleNamesAuthorPublisherPageSectionAndKeywords(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	article := buildtest.JSONLDOfType(t, buildtest.Page(t, out, raftURL), "Article")
	if article == nil {
		t.Fatal("no Article JSON-LD block")
	}

	author, _ := article["author"].(map[string]any)
	if author["@type"] != "Person" || author["name"] != "Ada Lovelace" || author["url"] != "https://hastcode.test/authors/ada/" {
		t.Errorf("author = %#v, want Person Ada Lovelace with url https://hastcode.test/authors/ada/", author)
	}
	if got, want := fmt.Sprint(author["sameAs"]), "[https://github.com/ada-example https://www.linkedin.com/in/ada-example]"; got != want {
		t.Errorf("author.sameAs = %s, want %s", got, want)
	}

	publisher, _ := article["publisher"].(map[string]any)
	if publisher["@type"] != "Organization" || publisher["name"] != "Hast Code" || publisher["url"] != "https://hastcode.test/" {
		t.Errorf("publisher = %#v, want Organization Hast Code at https://hastcode.test/", publisher)
	}
	if logo, _ := publisher["logo"].(map[string]any); logo["url"] != "https://hastcode.test/mark.svg" {
		t.Errorf("publisher.logo = %#v, want url https://hastcode.test/mark.svg", publisher["logo"])
	}

	page, _ := article["mainEntityOfPage"].(map[string]any)
	if page["@type"] != "WebPage" || page["@id"] != "https://hastcode.test"+raftURL {
		t.Errorf("mainEntityOfPage = %#v, want WebPage @id https://hastcode.test%s", page, raftURL)
	}
	if got := article["articleSection"]; got != "Distributed Systems" {
		t.Errorf("articleSection = %#v, want the Vertical name", got)
	}
	if got := fmt.Sprint(article["keywords"]); got != "[consensus raft]" {
		t.Errorf("keywords = %s, want the tags [consensus raft]", got)
	}
}

func TestExplainerJSONLD_ArticleImagesAreTheHeroIn16x9_4x3And1x1(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	article := buildtest.JSONLDOfType(t, buildtest.Page(t, out, raftURL), "Article")
	if article == nil {
		t.Fatal("no Article JSON-LD block")
	}
	images, _ := article["image"].([]any)
	if len(images) != 3 {
		t.Fatalf("image = %#v, want 3 URLs", article["image"])
	}

	for i, want := range []struct{ w, h int }{{16, 9}, {4, 3}, {1, 1}} {
		url, _ := images[i].(string)
		path, ok := strings.CutPrefix(url, "https://hastcode.test/")
		if !ok {
			t.Errorf("image[%d] = %q, want an absolute URL on the site", i, url)
			continue
		}
		img := buildtest.DecodeImage(t, out, "/"+path)
		if img.Format != "jpeg" || img.Width < 1200 || img.Width*want.h != img.Height*want.w {
			t.Errorf("image[%d] %s is a %dx%d %s, want a JPEG ≥1200 px wide at %d:%d", i, url, img.Width, img.Height, img.Format, want.w, want.h)
		}
	}
}

// rfc3339WithOffset is RFC 3339 with an explicit timezone: Z or ±hh:mm.
var rfc3339WithOffset = regexp.MustCompile(`^\d{4}-\d\d-\d\dT\d\d:\d\d:\d\d(\.\d+)?(Z|[+-]\d\d:\d\d)$`)

func TestExplainerJSONLD_DatesAreRFC3339WithOffsetAndMatchTheVisibleDates(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	doc := buildtest.Page(t, out, raftURL)
	article := buildtest.JSONLDOfType(t, doc, "Article")
	if article == nil {
		t.Fatal("no Article JSON-LD block")
	}

	for key, visible := range map[string]string{
		"datePublished": "time.dt-published",
		"dateModified":  "time.dt-updated",
	} {
		value, _ := article[key].(string)
		if !rfc3339WithOffset.MatchString(value) {
			t.Errorf("%s = %q, want RFC 3339 with a timezone offset", key, value)
			continue
		}
		got, _ := time.Parse(time.RFC3339, value)
		shownAttr := doc.Find(visible).AttrOr("datetime", "")
		shown, err := time.Parse(time.RFC3339, shownAttr)
		if err != nil || !got.Equal(shown) {
			t.Errorf("%s = %s, but the visible %s is %q", key, value, visible, shownAttr)
		}
	}
}

func TestExplainerJSONLD_DateModifiedIsTheReviewDateNotTheBuildTime(t *testing.T) {
	cfg := buildtest.Config()
	cfg.Now = time.Date(2026, 12, 25, 12, 0, 0, 0, time.UTC)
	out := buildtest.MustBuild(t, cfg, buildtest.ValidTree())
	article := buildtest.JSONLDOfType(t, buildtest.Page(t, out, raftURL), "Article")
	if article == nil {
		t.Fatal("no Article JSON-LD block")
	}
	if got := article["dateModified"]; got != "2026-10-15T18:30:00Z" {
		t.Errorf("dateModified = %#v, want the frontmatter reviewed 2026-10-15T18:30:00Z", got)
	}
}

func TestExplainerJSONLD_BreadcrumbListRunsHomeVerticalExplainer(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	list := buildtest.JSONLDOfType(t, buildtest.Page(t, out, raftURL), "BreadcrumbList")
	if list == nil {
		t.Fatal("no BreadcrumbList JSON-LD block")
	}
	items, _ := list["itemListElement"].([]any)
	want := []struct{ name, url string }{
		{"Home", "https://hastcode.test/"},
		{"Distributed Systems", "https://hastcode.test/distributed-systems/"},
		{"How Raft elects a leader", "https://hastcode.test" + raftURL},
	}
	if len(items) != len(want) {
		t.Fatalf("itemListElement has %d items, want %d: %#v", len(items), len(want), items)
	}
	for i, w := range want {
		item, _ := items[i].(map[string]any)
		if item["@type"] != "ListItem" || item["position"] != float64(i+1) || item["name"] != w.name || item["item"] != w.url {
			t.Errorf("item %d = %#v, want ListItem %d %q %s", i, item, i+1, w.name, w.url)
		}
	}
}

func TestExplainerJSONLD_NoUnsupportedTypes(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	for _, typ := range []string{"FAQPage", "HowTo", "Speakable", "TechArticle"} {
		if hits := buildtest.FilesContaining(t, out, `"`+typ+`"`); len(hits) > 0 {
			t.Errorf("%s appears in %v; the spec rules it out", typ, hits)
		}
	}
}

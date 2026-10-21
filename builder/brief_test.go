package builder_test

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/internal/buildtest"
)

const briefURL = "/briefs/etcd-raft-rewrite/"

func briefTree() buildtest.Tree {
	return buildtest.ValidTree().WithBrief(buildtest.ValidBrief())
}

func TestBrief_ValidBriefIsWrittenAtBriefsSlug(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), briefTree())

	doc := buildtest.Page(t, out, briefURL)
	if got := doc.Find("article.h-entry h1.p-name").Text(); got != "etcd ships a rewritten Raft library" {
		t.Errorf("h1 = %q, want the Brief title", got)
	}
	if got := doc.Find("article.h-entry .e-content h2").Length(); got != 2 {
		t.Errorf("found %d section headings in the body, want What happened and Why it matters", got)
	}
}

func TestBriefJSONLD_NewsArticleCarriesHeadlineDatesAuthorPublisherAndHeroCrops(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), briefTree())
	doc := buildtest.Page(t, out, briefURL)
	news := buildtest.JSONLDOfType(t, doc, "NewsArticle")
	if news == nil {
		t.Fatal("no NewsArticle JSON-LD block")
	}

	for key, want := range map[string]string{
		"@context":      "https://schema.org",
		"headline":      "etcd ships a rewritten Raft library",
		"description":   "etcd's new Raft library cuts election latency; here is why leader election design matters.",
		"datePublished": "2026-10-20T08:15:00-04:00",
		"dateModified":  "2026-10-20T08:15:00-04:00",
	} {
		if got := news[key]; got != want {
			t.Errorf("NewsArticle %s = %#v, want %q", key, got, want)
		}
	}
	shown := doc.Find("article.h-entry time.dt-published").AttrOr("datetime", "")
	if shown != "2026-10-20T08:15:00-04:00" {
		t.Errorf("visible dt-published datetime = %q, want the JSON-LD datePublished", shown)
	}

	author, _ := news["author"].(map[string]any)
	if author["@type"] != "Person" || author["name"] != "Ada Lovelace" || author["url"] != "https://hastcode.test/authors/ada/" {
		t.Errorf("author = %#v, want Person Ada Lovelace with url https://hastcode.test/authors/ada/", author)
	}
	publisher, _ := news["publisher"].(map[string]any)
	if publisher["@type"] != "Organization" || publisher["name"] != "Hast Code" {
		t.Errorf("publisher = %#v, want Organization Hast Code", publisher)
	}
	page, _ := news["mainEntityOfPage"].(map[string]any)
	if page["@id"] != "https://hastcode.test"+briefURL {
		t.Errorf("mainEntityOfPage = %#v, want @id https://hastcode.test%s", page, briefURL)
	}

	images, _ := news["image"].([]any)
	if len(images) != 3 {
		t.Fatalf("image = %#v, want 3 URLs", news["image"])
	}
	for i, want := range []struct{ w, h int }{{16, 9}, {4, 3}, {1, 1}} {
		url, _ := images[i].(string)
		path, ok := strings.CutPrefix(url, "https://hastcode.test/")
		if !ok {
			t.Errorf("image[%d] = %q, want an absolute URL on the site", i, url)
			continue
		}
		img := buildtest.DecodeImage(t, out, "/"+path)
		if img.Width < 1200 || img.Width*want.h != img.Height*want.w {
			t.Errorf("image[%d] %s is %dx%d, want ≥1200 px wide at %d:%d", i, url, img.Width, img.Height, want.w, want.h)
		}
	}
	if buildtest.JSONLDOfType(t, doc, "Article") != nil {
		t.Error("a Brief carries an Article block; it must be a NewsArticle only")
	}
}

func TestBrief_BreadcrumbRunsHomeBriefsBrief(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), briefTree())
	doc := buildtest.Page(t, out, briefURL)

	var visible []string
	doc.Find("nav.breadcrumbs li").Each(func(_ int, li *goquery.Selection) { visible = append(visible, li.Text()) })
	if got := strings.Join(visible, " > "); got != "Home > Briefs > etcd ships a rewritten Raft library" {
		t.Errorf("visible breadcrumb = %q", got)
	}
	if href := doc.Find(`nav.breadcrumbs a`).Eq(1).AttrOr("href", ""); href != "/briefs/" {
		t.Errorf("Briefs crumb links %q, want /briefs/", href)
	}

	list := buildtest.JSONLDOfType(t, doc, "BreadcrumbList")
	items, _ := list["itemListElement"].([]any)
	want := []string{"https://hastcode.test/", "https://hastcode.test/briefs/", "https://hastcode.test" + briefURL}
	if len(items) != len(want) {
		t.Fatalf("BreadcrumbList has %d items, want %d", len(items), len(want))
	}
	for i, w := range want {
		if item, _ := items[i].(map[string]any); item["item"] != w {
			t.Errorf("BreadcrumbList item %d = %#v, want %s", i, item, w)
		}
	}
}

func TestBrief_ExplainerLinksBecomeTrailingSlashInternalLinks(t *testing.T) {
	b := buildtest.ValidBrief()
	b.Body = strings.Replace(buildtest.ValidBriefBody, "## Why it matters",
		"Compare [B-trees and LSM-trees](/databases-storage/btree-vs-lsm-tree#trade-offs).\n\n## Why it matters", 1)
	b.Explainers = []string{"raft-leader-election", "btree-vs-lsm-tree"}
	tree := buildtest.ValidTree().WithExplainer(buildtest.ValidComparison()).WithBrief(b)

	out := buildtest.MustBuild(t, buildtest.Config(), tree)
	doc := buildtest.Page(t, out, briefURL)

	for text, want := range map[string]string{
		"how Raft elects a leader": "/distributed-systems/raft-leader-election/",
		"B-trees and LSM-trees":    "/databases-storage/btree-vs-lsm-tree/#trade-offs",
	} {
		link := doc.Find(".e-content a").FilterFunction(func(_ int, a *goquery.Selection) bool { return a.Text() == text })
		if got := link.AttrOr("href", ""); got != want {
			t.Errorf("link %q href = %q, want %s", text, got, want)
		}
	}
	if got := doc.Find(".brief-source a").AttrOr("href", ""); got != "https://etcd.io/blog/2026/raft-rewrite/" {
		t.Errorf("source link = %q, want the source_url", got)
	}
}

func TestBriefByline_LinksTheAuthorPage(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), briefTree())
	author := buildtest.Page(t, out, briefURL).Find("article.h-entry .byline a.p-author")

	if author.Length() != 1 || author.AttrOr("href", "") != "/authors/ada/" || author.Text() != "Ada Lovelace" {
		t.Fatalf("byline = %d links, href %q, text %q; want one link to /authors/ada/ named Ada Lovelace",
			author.Length(), author.AttrOr("href", ""), author.Text())
	}
}

func TestAuthorPage_ListsTheAuthorsBriefsAmongTheirPieces(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), briefTree())

	var got []string
	buildtest.Page(t, out, "/authors/ada/").Find("main .latest li a").Each(func(_ int, a *goquery.Selection) {
		got = append(got, a.AttrOr("href", ""))
	})
	// The Brief (20 Oct) is newer than the Explainer (1 Oct).
	if want := briefURL + "|/distributed-systems/raft-leader-election/"; strings.Join(got, "|") != want {
		t.Errorf("Author page pieces = %q, want %q", got, want)
	}
}

func TestBrief_NewsletterFormEndsTheBriefWithItsPathAsSource(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), briefTree())
	doc := buildtest.Page(t, out, briefURL)

	form := doc.Find("article.h-entry form.newsletter-form")
	if form.Length() != 1 {
		t.Fatalf("found %d newsletter forms in the Brief, want 1", form.Length())
	}
	if doc.Find("article.h-entry .e-content").NextAllFiltered(".newsletter").Length() != 1 {
		t.Error("newsletter form does not follow the Brief body")
	}
	if got := form.Find(`input[name="source"]`).AttrOr("value", ""); got != briefURL {
		t.Errorf("hidden source = %q, want %s", got, briefURL)
	}
}

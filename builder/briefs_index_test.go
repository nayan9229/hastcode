package builder_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/internal/buildtest"
)

// briefsAtHours is ValidTree plus one valid Brief per hour h in [0, n),
// slug "brief-hNN", published at 2026-10-01T00:00 + h hours (so higher
// hours are newer).
func briefsAtHours(n int) buildtest.Tree {
	tree := buildtest.ValidTree()
	for h := range n {
		b := buildtest.ValidBrief()
		b.Slug = fmt.Sprintf("brief-h%02d", h)
		b.Title = fmt.Sprintf("Brief at hour %02d", h)
		b.Hero = b.Slug + ".png"
		b.Published = fmt.Sprintf("2026-10-01T%02d:00:00Z", h%24)
		if h >= 24 {
			b.Published = fmt.Sprintf("2026-10-02T%02d:00:00Z", h-24)
		}
		tree.WithBrief(b)
	}
	return tree
}

// indexSlugs lists the Brief URLs a Briefs index page links, in order.
func indexSlugs(doc *goquery.Document) []string {
	var urls []string
	doc.Find(".brief-list a.u-url").Each(func(_ int, a *goquery.Selection) {
		urls = append(urls, strings.TrimSuffix(strings.TrimPrefix(a.AttrOr("href", ""), "/briefs/"), "/"))
	})
	return urls
}

func TestBriefsIndex_ListsNewestFirstTwentyPerPage(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), briefsAtHours(41))

	first := indexSlugs(buildtest.Page(t, out, "/briefs/"))
	second := indexSlugs(buildtest.Page(t, out, "/briefs/page/2/"))
	third := indexSlugs(buildtest.Page(t, out, "/briefs/page/3/"))

	if len(first) != 20 || first[0] != "brief-h40" || first[19] != "brief-h21" {
		t.Errorf("page 1 = %v, want 20 Briefs from brief-h40 down to brief-h21", first)
	}
	if len(second) != 20 || second[0] != "brief-h20" || second[19] != "brief-h01" {
		t.Errorf("page 2 = %v, want 20 Briefs from brief-h20 down to brief-h01", second)
	}
	if strings.Join(third, ",") != "brief-h00" {
		t.Errorf("page 3 = %v, want only brief-h00", third)
	}
	for _, url := range []string{"/briefs/page/1/", "/briefs/page/4/"} {
		if buildtest.Exists(out, url) {
			t.Errorf("%s was written; page 1 is /briefs/ and there are only 3 pages", url)
		}
	}
}

func TestBriefsIndex_PagesLinkRelPrevAndNext(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), briefsAtHours(41))

	for url, want := range map[string]struct{ prev, next string }{
		"/briefs/":        {"", "/briefs/page/2/"},
		"/briefs/page/2/": {"/briefs/", "/briefs/page/3/"},
		"/briefs/page/3/": {"/briefs/page/2/", ""},
	} {
		doc := buildtest.Page(t, out, url)
		if got := doc.Find(`a[rel="prev"]`).AttrOr("href", ""); got != want.prev {
			t.Errorf("%s rel=prev = %q, want %q", url, got, want.prev)
		}
		if got := doc.Find(`a[rel="next"]`).AttrOr("href", ""); got != want.next {
			t.Errorf("%s rel=next = %q, want %q", url, got, want.next)
		}
		if got := doc.Find(`link[rel="canonical"]`).AttrOr("href", ""); got != "https://hastcode.test"+url {
			t.Errorf("%s canonical = %q, want itself", url, got)
		}
	}
}

func TestBriefsIndex_ItemShowsTitleDescriptionAndDate(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), briefTree())
	item := buildtest.Page(t, out, "/briefs/").Find(".brief-list .h-entry").First()

	if got := item.Find("a.u-url").AttrOr("href", ""); got != briefURL {
		t.Errorf("item link = %q, want %s", got, briefURL)
	}
	if got := strings.TrimSpace(item.Find(".p-name").Text()); got != "etcd ships a rewritten Raft library" {
		t.Errorf("item title = %q", got)
	}
	if got := item.Find(".p-summary").Text(); !strings.HasPrefix(got, "etcd's new Raft library") {
		t.Errorf("item description = %q", got)
	}
	if got := item.Find("time.dt-published").AttrOr("datetime", ""); got != "2026-10-20T08:15:00-04:00" {
		t.Errorf("item date = %q", got)
	}
}

func TestBriefsIndex_BreadcrumbRunsHomeBriefs(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), briefTree())
	doc := buildtest.Page(t, out, "/briefs/")

	var visible []string
	doc.Find("nav.breadcrumbs li").Each(func(_ int, li *goquery.Selection) { visible = append(visible, li.Text()) })
	if got := strings.Join(visible, " > "); got != "Home > Briefs" {
		t.Errorf("visible breadcrumb = %q, want Home > Briefs", got)
	}
	if buildtest.JSONLDOfType(t, doc, "BreadcrumbList") == nil {
		t.Error("no BreadcrumbList JSON-LD on the Briefs index")
	}
}

func TestBriefsIndex_WrittenWithNoBriefsAndNoPagination(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	doc := buildtest.Page(t, out, "/briefs/")

	if n := doc.Find(".brief-list .h-entry").Length(); n != 0 {
		t.Errorf("found %d Briefs, want none", n)
	}
	if n := doc.Find(`a[rel="prev"], a[rel="next"]`).Length(); n != 0 {
		t.Errorf("found %d pagination links on a single page, want none", n)
	}
	if buildtest.Exists(out, "/briefs/page/2/") {
		t.Error("an empty /briefs/page/2/ was written")
	}
}

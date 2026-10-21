package builder_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildtest"
)

// explainerFixture is a valid Explainer of the given kind in vertical,
// published on day (2026-09-<day>, so a higher day is newer), with tags.
// Comparisons get the Trade-offs table their skeleton needs.
func explainerFixture(vertical, slug, kind string, day int, tags ...string) buildtest.Explainer {
	e := buildtest.ValidExplainer()
	if kind == "comparison" {
		e = buildtest.ValidComparison()
	}
	e.Vertical = vertical
	e.Slug = slug
	e.Kind = kind
	e.Title = "Title of " + slug
	e.Description = "Description of " + slug + "."
	e.Hero = slug + ".png"
	e.Published = fmt.Sprintf("2026-09-%02dT09:00:00Z", day)
	e.Reviewed = "2026-10-15T09:00:00Z"
	e.Tags = tags
	return e
}

// listingTree has Explainers across two Verticals and all three kinds,
// plus one Brief (which links raft-leader-election, so that Explainer
// is in the tree too, published 2026-10-01).
//
//	distributed-systems: raft-leader-election (concept, 10-01), paxos (concept, 09-05),
//	                     spanner-truetime (case-study, 09-20), gossip (concept, 09-10)
//	databases-storage:   btree-vs-lsm (comparison, 09-15), dynamo (case-study, 09-25)
func listingTree() buildtest.Tree {
	return buildtest.ValidTree().
		WithExplainer(explainerFixture("distributed-systems", "paxos", "concept", 5, "consensus")).
		WithExplainer(explainerFixture("distributed-systems", "spanner-truetime", "case-study", 20, "clocks")).
		WithExplainer(explainerFixture("distributed-systems", "gossip", "concept", 10, "membership")).
		WithExplainer(explainerFixture("databases-storage", "btree-vs-lsm", "comparison", 15, "storage-engines")).
		WithExplainer(explainerFixture("databases-storage", "dynamo", "case-study", 25, "consensus", "storage-engines")).
		WithBrief(buildtest.ValidBrief())
}

// listedURLs lists the piece URLs a listing section links, in order.
func listedURLs(sel *goquery.Selection) []string {
	var urls []string
	sel.Find(".h-entry a.u-url").Each(func(_ int, a *goquery.Selection) {
		urls = append(urls, a.AttrOr("href", ""))
	})
	return urls
}

func TestVerticalHub_ListsItsExplainersNewestFirst(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), listingTree())

	got := listedURLs(buildtest.Page(t, out, "/distributed-systems/").Find("main"))
	want := []string{
		"/distributed-systems/raft-leader-election/",
		"/distributed-systems/spanner-truetime/",
		"/distributed-systems/gossip/",
		"/distributed-systems/paxos/",
	}
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("hub lists %v, want %v", got, want)
	}

	got = listedURLs(buildtest.Page(t, out, "/databases-storage/").Find("main"))
	want = []string{"/databases-storage/dynamo/", "/databases-storage/btree-vs-lsm/"}
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("hub lists %v, want %v", got, want)
	}
}

// manyExplainers is ValidTree plus n Concepts in databases-storage,
// slug "db-NN" published on 2026-09-NN (higher is newer; n ≤ 28).
func manyExplainers(n int) buildtest.Tree {
	tree := buildtest.ValidTree()
	for i := 1; i <= n; i++ {
		tree.WithExplainer(explainerFixture("databases-storage", fmt.Sprintf("db-%02d", i), "concept", i))
	}
	return tree
}

func TestVerticalHub_PaginatesTwentyPerPageWithRelPrevNext(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), manyExplainers(21))

	first := listedURLs(buildtest.Page(t, out, "/databases-storage/").Find("main"))
	second := listedURLs(buildtest.Page(t, out, "/databases-storage/page/2/").Find("main"))
	if len(first) != 20 || first[0] != "/databases-storage/db-21/" || first[19] != "/databases-storage/db-02/" {
		t.Errorf("page 1 = %v, want db-21 down to db-02", first)
	}
	if strings.Join(second, " ") != "/databases-storage/db-01/" {
		t.Errorf("page 2 = %v, want only db-01", second)
	}
	for _, url := range []string{"/databases-storage/page/1/", "/databases-storage/page/3/"} {
		if buildtest.Exists(out, url) {
			t.Errorf("%s was written; page 1 is the hub and there are only 2 pages", url)
		}
	}
	for url, want := range map[string]struct{ prev, next string }{
		"/databases-storage/":        {"", "/databases-storage/page/2/"},
		"/databases-storage/page/2/": {"/databases-storage/", ""},
	} {
		doc := buildtest.Page(t, out, url)
		if got := doc.Find(`a[rel="prev"]`).AttrOr("href", ""); got != want.prev {
			t.Errorf("%s rel=prev = %q, want %q", url, got, want.prev)
		}
		if got := doc.Find(`a[rel="next"]`).AttrOr("href", ""); got != want.next {
			t.Errorf("%s rel=next = %q, want %q", url, got, want.next)
		}
	}
}

func TestVerticalHub_WrittenForEveryVerticalEvenWhenEmpty(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())

	doc := buildtest.Page(t, out, "/languages-runtimes/")
	if n := doc.Find("main .h-entry").Length(); n != 0 {
		t.Errorf("empty hub lists %d items", n)
	}
	if got := doc.Find("h1").Text(); got != "Languages & Runtimes" {
		t.Errorf("empty hub h1 = %q", got)
	}
	if n := doc.Find(`a[rel="prev"], a[rel="next"]`).Length(); n != 0 {
		t.Errorf("single-page hub has %d pagination links", n)
	}
	if buildtest.Exists(out, "/languages-runtimes/page/2/") {
		t.Error("an empty page 2 was written")
	}
}

func TestKindPage_ListsThatKindAcrossVerticalsNewestFirst(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), listingTree())

	for url, want := range map[string][]string{
		"/explainers/case-study/": {"/databases-storage/dynamo/", "/distributed-systems/spanner-truetime/"},
		"/explainers/concept/": {
			"/distributed-systems/raft-leader-election/",
			"/distributed-systems/gossip/",
			"/distributed-systems/paxos/",
		},
		"/explainers/comparison/": {"/databases-storage/btree-vs-lsm/"},
	} {
		got := listedURLs(buildtest.Page(t, out, url).Find("main"))
		if strings.Join(got, " ") != strings.Join(want, " ") {
			t.Errorf("%s lists %v, want %v", url, got, want)
		}
	}
}

func TestListingItem_ShowsTitleDescriptionKindDateAndThumbnail(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), listingTree())
	item := buildtest.Page(t, out, "/databases-storage/").Find("main .h-entry").First()

	if got := item.Find("a.u-url").AttrOr("href", ""); got != "/databases-storage/dynamo/" {
		t.Fatalf("first item links %q, want the newest Explainer", got)
	}
	if got := strings.TrimSpace(item.Find(".p-name").Text()); got != "Title of dynamo" {
		t.Errorf("title = %q", got)
	}
	if got := item.Find(".p-summary").Text(); got != "Description of dynamo." {
		t.Errorf("description = %q", got)
	}
	if got := strings.TrimSpace(item.Find(".badge").Text()); got != "Case Study" {
		t.Errorf("kind badge = %q, want Case Study", got)
	}
	if got := item.Find("time.dt-published").AttrOr("datetime", ""); got != "2026-09-25T09:00:00Z" {
		t.Errorf("date = %q", got)
	}
	assertThumbnail480(t, out, item)
}

func TestBriefsIndex_ItemShowsHeroThumbnailAndNoKindBadge(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), listingTree())
	item := buildtest.Page(t, out, "/briefs/").Find("main .h-entry").First()

	assertThumbnail480(t, out, item)
	if n := item.Find(".badge").Length(); n != 0 {
		t.Errorf("a Brief shows %d kind badges, want none", n)
	}
}

// assertThumbnail480 checks a listing item's thumbnail is the hero's
// 480 px variant in AVIF, WebP and JPEG, with its size reserved.
func assertThumbnail480(t *testing.T, out string, item *goquery.Selection) {
	t.Helper()
	img := item.Find("picture img")
	src := img.AttrOr("src", "")
	if !strings.Contains(src, "-480x") || !strings.HasSuffix(src, ".jpg") {
		t.Errorf("thumbnail img src = %q, want the 480 px JPEG", src)
	}
	if w, h := img.AttrOr("width", ""), img.AttrOr("height", ""); w != "480" || h != "270" {
		t.Errorf("thumbnail size = %sx%s, want 480x270", w, h)
	}
	for _, typ := range []string{"image/avif", "image/webp"} {
		srcset := item.Find(`picture source[type="`+typ+`"]`).AttrOr("srcset", "")
		if !strings.Contains(srcset, "-480x270.") {
			t.Errorf("%s thumbnail = %q, want the 480 px variant", typ, srcset)
		}
		if !buildtest.Exists(out, srcset) {
			t.Errorf("%s thumbnail %s was not written", typ, srcset)
		}
	}
	if img.AttrOr("loading", "") != "lazy" {
		t.Error("thumbnail is not lazy-loaded")
	}
}

func TestHome_ShowsLatestExplainersAndBriefsNewestFirst(t *testing.T) {
	tree := listingTree()
	older := buildtest.ValidBrief()
	older.Slug = "older-brief"
	older.Hero = "older-brief.png"
	older.Published = "2026-10-19T08:00:00Z"
	tree.WithBrief(older)

	doc := buildtest.Page(t, buildtest.MustBuild(t, buildtest.Config(), tree), "/")

	got := listedURLs(doc.Find("section.latest-explainers"))
	want := []string{
		"/distributed-systems/raft-leader-election/",
		"/databases-storage/dynamo/",
		"/distributed-systems/spanner-truetime/",
		"/databases-storage/btree-vs-lsm/",
		"/distributed-systems/gossip/",
		"/distributed-systems/paxos/",
	}
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("latest Explainers = %v, want %v", got, want)
	}
	got = listedURLs(doc.Find("section.latest-briefs"))
	if strings.Join(got, " ") != "/briefs/etcd-raft-rewrite/ /briefs/older-brief/" {
		t.Errorf("latest Briefs = %v, want etcd-raft-rewrite then older-brief", got)
	}
	if got := doc.Find(`section.latest-briefs a[href="/briefs/"]`).Length(); got != 1 {
		t.Errorf("latest Briefs links the Briefs index %d times, want 1", got)
	}
}

func TestHome_ShowsOnlyTheSixNewestExplainers(t *testing.T) {
	doc := buildtest.Page(t, buildtest.MustBuild(t, buildtest.Config(), manyExplainers(10)), "/")

	got := listedURLs(doc.Find("section.latest-explainers"))
	if len(got) != 6 || got[0] != "/distributed-systems/raft-leader-election/" || got[5] != "/databases-storage/db-06/" {
		t.Errorf("latest Explainers = %v, want the 6 newest (raft-leader-election, db-10 … db-06)", got)
	}
}

func TestHome_ShowsEveryVerticalWithItsDescriptionAndTheKindPages(t *testing.T) {
	doc := buildtest.Page(t, buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree()), "/")

	var got []string
	doc.Find("section.verticals-list li").Each(func(_ int, li *goquery.Selection) {
		got = append(got, li.Find("a").AttrOr("href", "")+" "+strings.TrimSpace(li.Find("p").Text()))
	})
	var want []string
	for _, v := range buildtest.Config().Verticals {
		want = append(want, "/"+v.Slug+"/ "+v.Description)
	}
	if len(want) != 6 || strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("home Verticals = %q, want the six in order with descriptions %q", got, want)
	}
	for _, url := range []string{"/explainers/case-study/", "/explainers/concept/", "/explainers/comparison/"} {
		if doc.Find(`main a[href="`+url+`"]`).Length() == 0 {
			t.Errorf("home does not link kind page %s", url)
		}
	}
	if got := doc.Find(`link[rel="canonical"]`).AttrOr("href", ""); got != "https://hastcode.test/" {
		t.Errorf("home canonical = %q", got)
	}
}

func TestVerticals_AddingOneToConfigAddsItsHubNavAndHomeEntry(t *testing.T) {
	cfg := buildtest.Config()
	cfg.Verticals = append(cfg.Verticals, builder.Vertical{Slug: "security", Name: "Security", Description: "How systems keep attackers out.", Order: 7})
	tree := buildtest.ValidTree().WithExplainer(explainerFixture("security", "zero-trust", "concept", 3))

	out := buildtest.MustBuild(t, cfg, tree)

	hub := buildtest.Page(t, out, "/security/")
	if got := listedURLs(hub.Find("main")); strings.Join(got, " ") != "/security/zero-trust/" {
		t.Errorf("/security/ lists %v, want zero-trust", got)
	}
	if got := hub.Find("h1").Text(); got != "Security" {
		t.Errorf("hub h1 = %q", got)
	}
	nav := buildtest.Page(t, out, raftURL).Find(`nav[aria-label="Verticals"] a`)
	if last := nav.Last(); nav.Length() != 7 || last.AttrOr("href", "") != "/security/" || last.Text() != "Security" {
		t.Errorf("header Verticals nav has %d links, last %q; want 7 ending with Security", nav.Length(), last.AttrOr("href", ""))
	}
	if buildtest.Page(t, out, "/").Find(`section.verticals-list a[href="/security/"]`).Length() != 1 {
		t.Error("home does not list the new Vertical")
	}
}

// breadcrumbTrail is a page's visible breadcrumb as "Name href" entries
// (the current page last, unlinked).
func breadcrumbTrail(doc *goquery.Document) []string {
	var got []string
	doc.Find(`nav[aria-label="Breadcrumb"] li`).Each(func(_ int, li *goquery.Selection) {
		got = append(got, strings.TrimSpace(strings.TrimSpace(li.Text())+" "+li.Find("a").AttrOr("href", "")))
	})
	return got
}

func TestListingPages_HaveBreadcrumbsAndBreadcrumbList(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), manyExplainers(21))

	for url, want := range map[string][]string{
		"/databases-storage/":        {"Home /", "Databases & Storage"},
		"/databases-storage/page/2/": {"Home /", "Databases & Storage /databases-storage/", "Page 2"},
		"/explainers/case-study/":    {"Home /", "Case Studies"},
		"/explainers/concept/":       {"Home /", "Concepts"},
		"/explainers/comparison/":    {"Home /", "Comparisons"},
	} {
		doc := buildtest.Page(t, out, url)
		if got := breadcrumbTrail(doc); strings.Join(got, "|") != strings.Join(want, "|") {
			t.Errorf("%s breadcrumb = %q, want %q", url, got, want)
		}
		list := buildtest.JSONLDOfType(t, doc, "BreadcrumbList")
		if items, _ := list["itemListElement"].([]any); len(items) != len(want) {
			t.Errorf("%s BreadcrumbList has %d items, want %d", url, len(items), len(want))
		}
	}
}

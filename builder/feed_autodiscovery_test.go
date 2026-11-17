package builder_test

import (
	"fmt"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/internal/buildtest"
)

// Feed autodiscovery: every page's <head> advertises the site feed as
// <link rel="alternate" type="application/rss+xml">, so a feed reader,
// browser or feed-intake service given any page URL finds the feed.
// Vertical hub and Briefs listings also advertise their own feed.

// feedLink is one RSS autodiscovery link in a page's head.
type feedLink struct{ href, title string }

// feedLinks returns the page's RSS autodiscovery links, in head order.
func feedLinks(doc *goquery.Document) []feedLink {
	var links []feedLink
	doc.Find(`head link[rel="alternate"][type="application/rss+xml"]`).Each(func(_ int, l *goquery.Selection) {
		links = append(links, feedLink{href: l.AttrOr("href", ""), title: l.AttrOr("title", "")})
	})
	return links
}

var siteFeedLink = feedLink{href: "/feed.xml", title: "HastCode"}

func TestFeedAutodiscovery_EveryKindOfPageAdvertisesTheSiteFeedOnce(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree().WithBrief(buildtest.ValidBrief()))

	for _, urlPath := range []string{
		"/", // home
		"/distributed-systems/raft-leader-election/", // Explainer
		"/briefs/etcd-raft-rewrite/",                 // Brief
		"/about/",                                    // static page
		notFoundURL,                                  // 404
	} {
		got := feedLinks(buildtest.Page(t, out, urlPath))
		if len(got) != 1 || got[0] != siteFeedLink {
			t.Errorf("%s head feed links = %+v, want exactly %+v", urlPath, got, siteFeedLink)
		}
	}
}

// htmlPages lists the URL path of every HTML page the build wrote
// ("/", "/briefs/page/2/", "/404.html").
func htmlPages(t *testing.T, out string) []string {
	t.Helper()
	var pages []string
	for _, f := range buildtest.OutputFiles(t, out) {
		rel, ok := strings.CutPrefix(f, "static/")
		if !ok || !strings.HasSuffix(rel, ".html") {
			continue
		}
		if path.Base(rel) == "index.html" {
			pages = append(pages, "/"+strings.TrimSuffix(rel, "index.html"))
		} else {
			pages = append(pages, "/"+rel)
		}
	}
	if len(pages) == 0 {
		t.Fatal("the build wrote no HTML pages")
	}
	return pages
}

// listingFeedTree has two pages of databases-storage Explainers, two
// pages of Briefs, and Explainers in distributed-systems.
func listingFeedTree() buildtest.Tree {
	tree := manyExplainers(21)
	for h := range 21 {
		b := buildtest.ValidBrief()
		b.Slug = fmt.Sprintf("brief-h%02d", h)
		b.Title = fmt.Sprintf("Brief at hour %02d", h)
		b.Hero = b.Slug + ".png"
		b.Published = fmt.Sprintf("2026-10-01T%02d:00:00Z", h)
		tree.WithBrief(b)
	}
	return tree
}

func TestFeedAutodiscovery_VerticalAndBriefsListingsAddTheirOwnFeedAfterTheSiteFeed(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), listingFeedTree())

	databases := feedLink{href: "/databases-storage/feed.xml", title: "HastCode: Databases & Storage"}
	briefs := feedLink{href: "/briefs/feed.xml", title: "HastCode: Briefs"}
	want := map[string][]feedLink{
		"/databases-storage/":        {siteFeedLink, databases},
		"/databases-storage/page/2/": {siteFeedLink, databases},
		"/distributed-systems/":      {siteFeedLink, {href: "/distributed-systems/feed.xml", title: "HastCode: Distributed Systems"}},
		// Hubs with no Explainers are still hubs.
		"/apis-networking/":      {siteFeedLink, {href: "/apis-networking/feed.xml", title: "HastCode: APIs & Networking"}},
		"/infrastructure-cloud/": {siteFeedLink, {href: "/infrastructure-cloud/feed.xml", title: "HastCode: Infrastructure & Cloud"}},
		"/ai-ml-systems/":        {siteFeedLink, {href: "/ai-ml-systems/feed.xml", title: "HastCode: AI/ML Systems"}},
		"/languages-runtimes/":   {siteFeedLink, {href: "/languages-runtimes/feed.xml", title: "HastCode: Languages & Runtimes"}},
		"/briefs/":               {siteFeedLink, briefs},
		"/briefs/page/2/":        {siteFeedLink, briefs},
	}
	for urlPath, w := range want {
		if got := feedLinks(buildtest.Page(t, out, urlPath)); !reflect.DeepEqual(got, w) {
			t.Errorf("%s head feed links = %+v, want %+v", urlPath, got, w)
		}
	}

	// Every other page (Pieces, kind listings, home, Authors, static
	// pages, search, 404) advertises only the site feed.
	for _, urlPath := range htmlPages(t, out) {
		if _, listing := want[urlPath]; listing {
			continue
		}
		if got := feedLinks(buildtest.Page(t, out, urlPath)); !reflect.DeepEqual(got, []feedLink{siteFeedLink}) {
			t.Errorf("%s head feed links = %+v, want only the site feed", urlPath, got)
		}
	}
}

func TestFeedAutodiscovery_EveryAdvertisedFeedIsOneTheBuildWroteUnderThatTitle(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), listingFeedTree())

	checked := map[feedLink]bool{}
	for _, urlPath := range htmlPages(t, out) {
		for _, l := range feedLinks(buildtest.Page(t, out, urlPath)) {
			if checked[l] {
				continue
			}
			checked[l] = true
			if !buildtest.Exists(out, l.href) {
				t.Errorf("%s advertises %s, which the build did not write", urlPath, l.href)
				continue
			}
			f := readFeed(t, out, l.href)
			if _, self := f.link(); self != siteURL+l.href {
				t.Errorf("%s advertises %s, but that feed's self link is %q", urlPath, l.href, self)
			}
			if f.Channel.Title != l.title {
				t.Errorf("%s advertises %s titled %q; the feed's channel title is %q", urlPath, l.href, l.title, f.Channel.Title)
			}
		}
	}
	// The site feed, the Briefs feed and all six Vertical feeds.
	if len(checked) != 8 {
		t.Errorf("pages advertise %d distinct feeds, want 8: %+v", len(checked), checked)
	}
}

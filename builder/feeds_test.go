package builder_test

import (
	"encoding/xml"
	"os"
	"strings"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

// rss is the part of an RSS 2.0 document the feed tests read.
type rss struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Title string `xml:"title"`
		// Links holds <link> and <atom:link rel="self">.
		Links []struct {
			XMLName xml.Name
			Href    string `xml:"href,attr"`
			Rel     string `xml:"rel,attr"`
			Value   string `xml:",chardata"`
		} `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Items       []rssItem `xml:"item"`
	} `xml:"channel"`
}

// link is the channel's RSS <link> (the site), and self its Atom
// rel="self" link (the feed's own URL).
func (f rss) link() (link, self string) {
	for _, l := range f.Channel.Links {
		switch {
		case l.XMLName.Space == "":
			link = l.Value
		case l.XMLName.Space == "http://www.w3.org/2005/Atom" && l.Rel == "self":
			self = l.Href
		}
	}
	return link, self
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
	// Content is <content:encoded> (namespace http://purl.org/rss/1.0/modules/content/).
	Content string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
}

func readFeed(t *testing.T, out, urlPath string) rss {
	t.Helper()
	data, err := os.ReadFile(buildtest.StaticPath(out, urlPath))
	if err != nil {
		t.Fatalf("feed %s: %v", urlPath, err)
	}
	var f rss
	if err := xml.Unmarshal(data, &f); err != nil {
		t.Fatalf("feed %s does not parse as XML: %v\n%s", urlPath, err, data)
	}
	if f.Version != "2.0" {
		t.Errorf("feed %s: rss version = %q, want 2.0", urlPath, f.Version)
	}
	return f
}

func itemLinks(f rss) []string {
	var links []string
	for _, it := range f.Channel.Items {
		links = append(links, it.Link)
	}
	return links
}

func feedTree() buildtest.Tree {
	return buildtest.ValidTree().
		WithExplainer(buildtest.ValidComparison()).
		WithBrief(buildtest.ValidBrief())
}

const siteURL = "https://hastcode.test"

func TestSiteFeed_ListsEveryPieceNewestFirstWithAbsoluteLinks(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), feedTree())
	f := readFeed(t, out, "/feed.xml")

	want := []string{
		siteURL + "/briefs/etcd-raft-rewrite/",            // 2026-10-20
		siteURL + "/databases-storage/btree-vs-lsm-tree/", // 2026-10-02
		siteURL + "/distributed-systems/raft-leader-election/",
	}
	if got := itemLinks(f); strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("item links =\n%v\nwant (newest first)\n%v", got, want)
	}
	if got := f.Channel.Items[0]; got.Title != "etcd ships a rewritten Raft library" || got.GUID != want[0] || got.PubDate != "Tue, 20 Oct 2026 08:15:00 -0400" {
		t.Errorf("first item title/guid/pubDate = %q / %q / %q", got.Title, got.GUID, got.PubDate)
	}
}

func TestSiteFeed_ChannelIsHastCodeWithTheTaglineAsDescription(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), feedTree())
	f := readFeed(t, out, "/feed.xml")
	ch := f.Channel
	link, self := f.link()

	if ch.Title != "Hast Code" || link != siteURL+"/" || ch.Description != "How real systems are built." || ch.Language != "en" {
		t.Errorf("channel title/link/description/language = %q / %q / %q / %q", ch.Title, link, ch.Description, ch.Language)
	}
	if self != siteURL+"/feed.xml" {
		t.Errorf("atom:link rel=self = %q, want the feed's own absolute URL", self)
	}
}

func TestSiteFeed_ItemsCarryTheFullHTMLBodyWithAbsoluteLinks(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), feedTree())
	f := readFeed(t, out, "/feed.xml")

	var brief, raft rssItem
	for _, it := range f.Channel.Items {
		switch it.Link {
		case siteURL + "/briefs/etcd-raft-rewrite/":
			brief = it
		case siteURL + "/distributed-systems/raft-leader-election/":
			raft = it
		}
	}
	// Full text: every section, not a summary.
	for _, want := range []string{"<h2", "What happened", "Why it matters", "halves election latency"} {
		if !strings.Contains(brief.Content, want) {
			t.Errorf("Brief content:encoded lacks %q:\n%s", want, brief.Content)
		}
	}
	// The Brief's link to the Explainer is absolute in the feed.
	if !strings.Contains(brief.Content, `href="`+siteURL+`/distributed-systems/raft-leader-election/"`) {
		t.Errorf("Brief content:encoded does not link the Explainer absolutely:\n%s", brief.Content)
	}
	for _, want := range []string{"TL;DR", "Each term has at most one leader.", "<svg", "Sources", "https://raft.github.io/raft.pdf"} {
		if !strings.Contains(raft.Content, want) {
			t.Errorf("Explainer content:encoded lacks %q", want)
		}
	}
	if raft.Description != buildtest.ValidExplainer().Description {
		t.Errorf("Explainer description = %q, want the frontmatter description", raft.Description)
	}
	for _, it := range f.Channel.Items {
		for _, rel := range []string{`href="/`, `src="/`} {
			if strings.Contains(it.Content, rel) {
				t.Errorf("%s content:encoded has a site-relative %s…", it.Link, rel)
			}
		}
	}
}

// Feed readers load no site CSS, so a diagram painted through CSS
// variables would come out black boxes with invisible edges. Feeds carry
// the diagrams with the light palette resolved (spec brand values: fg
// #111827, accent #0D9488); pages keep the themable variables.
func TestSiteFeed_DiagramsAreReadableWithoutTheSiteStylesheet(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	f := readFeed(t, out, "/feed.xml")
	if len(f.Channel.Items) != 1 {
		t.Fatalf("feed has %d items, want 1", len(f.Channel.Items))
	}
	content := f.Channel.Items[0].Content
	svg := content[strings.Index(content, "<svg"):strings.Index(content, "</svg>")]
	if strings.Contains(svg, "var(--") {
		t.Errorf("feed diagram still paints through CSS variables, which feed readers do not define:\n%.600s", svg)
	}
	for _, want := range []string{`fill="#111827"`, `fill="#0D9488"`} {
		if !strings.Contains(strings.ToUpper(svg), strings.ToUpper(want)) {
			t.Errorf("feed diagram lacks %s (text in the brand fg, arrows in the brand accent)", want)
		}
	}

	page := buildtest.Page(t, out, explainerPath)
	pageSVG, _ := page.Find("figure.diagram svg").Html()
	if !strings.Contains(pageSVG, "var(--_text)") {
		t.Errorf("page diagram lost its CSS variables; it must stay themable (light, dark, toggle)")
	}
}

func TestVerticalAndBriefsFeeds_CarryOnlyTheirOwnPieces(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), feedTree())

	for path, want := range map[string][]string{
		"/distributed-systems/feed.xml": {siteURL + "/distributed-systems/raft-leader-election/"},
		"/databases-storage/feed.xml":   {siteURL + "/databases-storage/btree-vs-lsm-tree/"},
		"/apis-networking/feed.xml":     nil, // a Vertical with no Explainers yet still has a valid, empty feed
		"/briefs/feed.xml":              {siteURL + "/briefs/etcd-raft-rewrite/"},
	} {
		f := readFeed(t, out, path)
		if got := itemLinks(f); strings.Join(got, " ") != strings.Join(want, " ") {
			t.Errorf("%s items = %v, want %v", path, got, want)
		}
		if f.Channel.Description != "How real systems are built." {
			t.Errorf("%s description = %q, want the tagline", path, f.Channel.Description)
		}
		if _, self := f.link(); self != siteURL+path {
			t.Errorf("%s atom:link self = %q", path, self)
		}
	}
}

func TestFeeds_LeaveOutAPieceWhosePageWasNotWritten(t *testing.T) {
	broken := buildtest.ValidComparison()
	broken.Body += "\nSee [the missing page](/databases-storage/no-such-explainer/).\n"
	out, res := buildtest.Build(t, buildtest.Config(), feedTree().WithExplainer(broken))
	buildtest.RequireErrorRule(t, res, "explainers/databases-storage/btree-vs-lsm-tree.md", "broken-link", "no-such-explainer")

	for _, path := range []string{"/feed.xml", "/databases-storage/feed.xml"} {
		for _, link := range itemLinks(readFeed(t, out, path)) {
			if strings.Contains(link, "btree-vs-lsm-tree") {
				t.Errorf("%s lists %s, whose page the link gate removed", path, link)
			}
		}
	}
}

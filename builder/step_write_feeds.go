package builder

import (
	"cmp"
	"encoding/xml"
	"regexp"
	"slices"
	"time"

	"hastcode.com/site/internal/markdown"
	"hastcode.com/site/internal/views"
)

func init() {
	registerStep(phaseWrite, 20, "write feeds", writeFeeds)
}

// webSubHub is the WebSub hub every feed names (atom:link rel="hub"):
// Google's public hub. Its own page (https://pubsubhubbub.appspot.com/,
// read 2026-11-17) asks publishers for exactly this link under
// rss/channel, and a form POST hub.mode=publish&hub.url=<feed> after each
// update. The rel="self" link beside it is the WebSub topic. The build
// never pings: step_write_websub.go lists the changed feeds for the
// post-deploy step. No live publish request has been sent yet: the
// Author sends one before turning WEBSUB_ENABLED on (README "WebSub hub
// ping", "Live check").
const webSubHub = "https://pubsubhubbub.appspot.com/"

// feedLength is how many of the newest pieces each feed carries. Feeds
// are full text, so the cap keeps /feed.xml a reasonable download.
const feedLength = 50

// writeFeeds writes the full-text RSS 2.0 feeds (spec "URLs"): /feed.xml
// (every piece: Explainers and Briefs), /{vertical}/feed.xml (that
// Vertical's Explainers) and /briefs/feed.xml. Feeds run in the write
// phase, after the link gate, and list only pieces whose page was written.
func writeFeeds(s *site) error {
	var all, briefs []feedItem
	byVertical := map[string][]feedItem{}
	for _, e := range s.publishedExplainers() {
		if !s.isWritten(explainerURL(e)) {
			continue
		}
		it, err := s.withFullText(feedItem{
			url: explainerURL(e), title: e.Title, description: e.Description,
			category: s.verticals[e.Vertical].Name, author: s.authors[e.Author].Name,
			published: e.Published,
		}, e.Body)
		if err != nil {
			return err
		}
		all = append(all, it)
		byVertical[e.Vertical] = append(byVertical[e.Vertical], it)
	}
	for _, b := range s.publishedBriefs() {
		if !s.isWritten(briefURL(b)) {
			continue
		}
		it, err := s.withFullText(feedItem{
			url: briefURL(b), title: b.Title, description: b.Description,
			category: "Briefs", author: s.authors[b.Author].Name,
			published: b.Published,
		}, b.Body)
		if err != nil {
			return err
		}
		all = append(all, it)
		briefs = append(briefs, it)
	}

	if err := s.writeFeed(siteFeed(), homeURL, views.Tagline, all); err != nil {
		return err
	}
	if err := s.writeFeed(briefsFeed(), briefsIndexURL, views.Tagline, briefs); err != nil {
		return err
	}
	for _, v := range s.sortedVerticals() {
		if err := s.writeFeed(verticalFeed(v), verticalURL(v.Slug), views.Tagline, byVertical[v.Slug]); err != nil {
			return err
		}
	}
	return nil
}

// siteFeed, briefsFeed and verticalFeed name each feed: Name is its
// channel title, URL its path. writeFeeds writes them and page heads
// advertise them (s.page, the Vertical hub and Briefs listings), so a
// head link always names a feed the build wrote, under its own title.
func siteFeed() views.Link { return views.Link{Name: views.PublicationName, URL: siteFeedURL} }

func briefsFeed() views.Link {
	return views.Link{Name: views.PublicationName + ": Briefs", URL: briefsFeedURL}
}

func verticalFeed(v Vertical) views.Link {
	return views.Link{Name: views.PublicationName + ": " + v.Name, URL: verticalFeedURL(v.Slug)}
}

// feedItem is one piece in a feed, before encoding.
type feedItem struct {
	url, title, description, category, author string
	published                                 time.Time
	html                                      string // full body, absolute links
}

// withFullText completes it with the piece's full body rendered for feed
// readers: root-relative links made absolute.
func (s *site) withFullText(it feedItem, body *markdown.Document) (feedItem, error) {
	html, err := markdown.Render(body)
	if err != nil {
		return feedItem{}, err
	}
	it.html = views.FeedSafeDiagrams(s.absoluteLinks(html))
	return it, nil
}

// siteRelativeRef matches an href or src attribute holding a root-relative
// path ("/a/", not the protocol-relative "//host/").
var siteRelativeRef = regexp.MustCompile(`(\s(?:href|src)=")/([^/])`)

// absoluteLinks makes root-relative href and src values in rendered HTML
// absolute on SiteURL: a feed reader resolves links against the feed, not
// the page.
func (s *site) absoluteLinks(html string) string {
	return siteRelativeRef.ReplaceAllString(html, "${1}"+s.absURL("/")+"${2}")
}

// RSS 2.0 with the content, Dublin Core and Atom extensions. encoding/xml
// writes the prefixed element names verbatim; the namespaces are declared
// on <rss>.
type rssDoc struct {
	XMLName   xml.Name   `xml:"rss"`
	Version   string     `xml:"version,attr"`
	ContentNS string     `xml:"xmlns:content,attr"`
	DCNS      string     `xml:"xmlns:dc,attr"`
	AtomNS    string     `xml:"xmlns:atom,attr"`
	Channel   rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title         string        `xml:"title"`
	Link          string        `xml:"link"`
	Description   string        `xml:"description"`
	Language      string        `xml:"language"`
	AtomLinks     []rssAtomLink `xml:"atom:link"` // rel="self", then rel="hub"
	LastBuildDate string        `xml:"lastBuildDate,omitempty"`
	Items         []rssItem     `xml:"item"`
}

type rssAtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr,omitempty"`
}

type rssItem struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	GUID        rssGUID  `xml:"guid"`
	PubDate     string   `xml:"pubDate"`
	Creator     string   `xml:"dc:creator,omitempty"`
	Category    string   `xml:"category,omitempty"`
	Description string   `xml:"description"`
	Content     rssCDATA `xml:"content:encoded"`
}

type rssGUID struct {
	IsPermaLink bool   `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

type rssCDATA struct {
	Value string `xml:",cdata"`
}

// writeFeed writes one feed of items (feed.Name the channel title,
// feed.URL where it is written), newest first, at most feedLength.
// lastBuildDate is the newest item's date (not the build clock), so an
// unchanged site writes an unchanged feed.
func (s *site) writeFeed(feed views.Link, link, description string, items []feedItem) error {
	urlPath, title := feed.URL, feed.Name
	items = slices.Clone(items)
	slices.SortStableFunc(items, func(a, b feedItem) int {
		return cmp.Or(b.published.Compare(a.published), cmp.Compare(a.url, b.url))
	})
	items = items[:min(len(items), feedLength)]

	doc := rssDoc{
		Version:   "2.0",
		ContentNS: "http://purl.org/rss/1.0/modules/content/",
		DCNS:      "http://purl.org/dc/elements/1.1/",
		AtomNS:    "http://www.w3.org/2005/Atom",
		Channel: rssChannel{
			Title:       title,
			Link:        s.absURL(link),
			Description: description,
			Language:    "en",
			AtomLinks: []rssAtomLink{
				{Href: s.absURL(urlPath), Rel: "self", Type: "application/rss+xml"},
				{Href: webSubHub, Rel: "hub"},
			},
		},
	}
	if len(items) > 0 {
		doc.Channel.LastBuildDate = items[0].published.Format(time.RFC1123Z)
	}
	for _, it := range items {
		u := s.absURL(it.url)
		doc.Channel.Items = append(doc.Channel.Items, rssItem{
			Title:       it.title,
			Link:        u,
			GUID:        rssGUID{IsPermaLink: true, Value: u},
			PubDate:     it.published.Format(time.RFC1123Z),
			Creator:     it.author,
			Category:    it.category,
			Description: it.description,
			Content:     rssCDATA{Value: it.html},
		})
	}
	data, err := encodeXML(doc)
	if err != nil {
		return err
	}
	if s.index != nil {
		s.index.recordFeed(s.absURL(urlPath), data)
	}
	return s.out.WriteStatic(urlPath, data)
}

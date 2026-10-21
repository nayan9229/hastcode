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

	if err := s.writeFeed(siteFeedURL, views.PublicationName, homeURL, views.Tagline, all); err != nil {
		return err
	}
	if err := s.writeFeed(briefsFeedURL, views.PublicationName+": Briefs", briefsIndexURL, views.Tagline, briefs); err != nil {
		return err
	}
	for _, v := range s.sortedVerticals() {
		if err := s.writeFeed(verticalFeedURL(v.Slug), views.PublicationName+": "+v.Name, verticalURL(v.Slug), views.Tagline, byVertical[v.Slug]); err != nil {
			return err
		}
	}
	return nil
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
	it.html = s.absoluteLinks(html)
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
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language"`
	Self          rssSelf   `xml:"atom:link"`
	LastBuildDate string    `xml:"lastBuildDate,omitempty"`
	Items         []rssItem `xml:"item"`
}

type rssSelf struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
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

// writeFeed writes one feed of items, newest first, at most feedLength.
// lastBuildDate is the newest item's date (not the build clock), so an
// unchanged site writes an unchanged feed.
func (s *site) writeFeed(urlPath, title, link, description string, items []feedItem) error {
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
			Self:        rssSelf{Href: s.absURL(urlPath), Rel: "self", Type: "application/rss+xml"},
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
	return s.writeXML(urlPath, doc)
}

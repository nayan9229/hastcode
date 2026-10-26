package intake

import (
	"bytes"
	"cmp"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

func init() { registerFetcher(FeedSource, feedFetcher{}) }

// feedFetcher fetches a feed Source: RSS 2.0 or Atom 1.0, told apart by the
// root element, never the Content-Type (publishers mislabel both; research
// §5.1). RSS items are keyed on guid, else link; Atom entries on id, else
// their link. A relative link resolves against the feed URL.
type feedFetcher struct{}

const atomNS = "http://www.w3.org/2005/Atom"

func (feedFetcher) fetch(ctx context.Context, in *Intake, src Source, prev SourceState) (Fetched, error) {
	resp, err := in.get(ctx, src.Feed, prev)
	if err != nil {
		return Fetched{}, err
	}
	if resp.NotModified {
		return Fetched{State: resp.State}, nil
	}
	entries, err := parseFeed(resp.Body, src.Feed)
	if err != nil {
		return Fetched{}, fmt.Errorf("%s: %w", src.Feed, err)
	}
	return Fetched{Entries: entries, State: resp.State}, nil
}

func parseFeed(body []byte, feedURL string) ([]Entry, error) {
	dec := xml.NewDecoder(bytes.NewReader(body))
	dec.CharsetReader = charset.NewReaderLabel
	dec.Strict = false
	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			return nil, errors.New("not a feed: no root element")
		}
		if err != nil {
			return nil, fmt.Errorf("not a feed: %w", err)
		}
		root, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		base, _ := url.Parse(feedURL)
		switch {
		case root.Name.Local == "rss":
			var doc rssDoc
			if err := dec.DecodeElement(&doc, &root); err != nil {
				return nil, fmt.Errorf("RSS: %w", err)
			}
			return doc.entries(base), nil
		case root.Name.Local == "feed" && root.Name.Space == atomNS:
			var doc atomFeed
			if err := dec.DecodeElement(&doc, &root); err != nil {
				return nil, fmt.Errorf("Atom: %w", err)
			}
			return doc.entries(base), nil
		default:
			return nil, fmt.Errorf("not an RSS 2.0 or Atom feed (root element <%s>)", root.Name.Local)
		}
	}
}

type rssDoc struct {
	Items []struct {
		Title   string    `xml:"title"`
		Links   []xmlText `xml:"link"`
		GUID    string    `xml:"guid"`
		PubDate string    `xml:"pubDate"`
	} `xml:"channel>item"`
}

// xmlText is an element's name and text; an RSS item's <link> is told
// apart from an <atom:link> some feeds add beside it by its namespace.
type xmlText struct {
	XMLName xml.Name
	Text    string `xml:",chardata"`
}

func (d *rssDoc) entries(base *url.URL) []Entry {
	var entries []Entry
	for _, it := range d.Items {
		link := ""
		for _, l := range it.Links {
			if l.XMLName.Space == "" && strings.TrimSpace(l.Text) != "" {
				link = resolve(base, strings.TrimSpace(l.Text))
				break
			}
		}
		key := cmp.Or(strings.TrimSpace(it.GUID), link)
		if key == "" {
			continue
		}
		entries = append(entries, Entry{
			Key:       key,
			Title:     oneLine(it.Title),
			Link:      link,
			Published: parseDate(it.PubDate, rssDateLayouts),
		})
	}
	return entries
}

// rssDateLayouts are RFC 822 dates as feeds write them: four-digit years,
// numeric or named zones, with or without the weekday.
var rssDateLayouts = []string{
	time.RFC1123Z, time.RFC1123,
	"Mon, 2 Jan 2006 15:04:05 -0700", "Mon, 2 Jan 2006 15:04:05 MST",
	"2 Jan 2006 15:04:05 -0700", "2 Jan 2006 15:04:05 MST",
	"Mon, 2 Jan 2006 15:04 -0700", time.RFC3339,
}

type atomFeed struct {
	Entries []struct {
		ID        string `xml:"http://www.w3.org/2005/Atom id"`
		Title     string `xml:"http://www.w3.org/2005/Atom title"`
		Published string `xml:"http://www.w3.org/2005/Atom published"`
		Updated   string `xml:"http://www.w3.org/2005/Atom updated"`
		Links     []struct {
			Rel  string `xml:"rel,attr"`
			Href string `xml:"href,attr"`
		} `xml:"http://www.w3.org/2005/Atom link"`
	} `xml:"http://www.w3.org/2005/Atom entry"`
}

func (f *atomFeed) entries(base *url.URL) []Entry {
	var entries []Entry
	for _, e := range f.Entries {
		link := ""
		for _, l := range e.Links {
			if l.Rel == "" || l.Rel == "alternate" {
				link = resolve(base, strings.TrimSpace(l.Href))
				break
			}
		}
		key := cmp.Or(strings.TrimSpace(e.ID), link)
		if key == "" {
			continue
		}
		// A changed <updated> is an edit, not a new item (research §5.1):
		// the key stays the id, and the date is <published> when given.
		entries = append(entries, Entry{
			Key:       key,
			Title:     oneLine(e.Title),
			Link:      link,
			Published: parseDate(cmp.Or(e.Published, e.Updated), []string{time.RFC3339}),
		})
	}
	return entries
}

func parseDate(s string, layouts []string) time.Time {
	s = strings.TrimSpace(s)
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}

func resolve(base *url.URL, ref string) string {
	if base == nil || ref == "" {
		return ref
	}
	u, err := base.Parse(ref)
	if err != nil {
		return ref
	}
	return u.String()
}

// oneLine collapses a title's whitespace (feeds wrap long titles).
func oneLine(s string) string { return strings.Join(strings.Fields(s), " ") }

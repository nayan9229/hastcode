package builder

import (
	"encoding/xml"
	"time"

	"hastcode.com/site/internal/views"
)

func init() {
	registerStep(phaseWrite, 30, "write sitemap", writeSitemap)
	registerStep(phaseWrite, 35, "write news sitemap", writeNewsSitemap)
}

const (
	sitemapURL     = "/sitemap.xml"
	newsSitemapURL = "/news-sitemap.xml"

	sitemapNS     = "http://www.sitemaps.org/schemas/sitemap/0.9"
	newsSitemapNS = "http://www.google.com/schemas/sitemap-news/0.9"

	// Google's news sitemap rules
	// (developers.google.com/search/docs/crawling-indexing/sitemaps/news-sitemap,
	// checked 2026-10-21): only articles created in the last two days, at
	// most 1,000 news:news entries per sitemap.
	newsWindow     = 48 * time.Hour
	newsMaxEntries = 1000
	newsLanguage   = "en" // ISO 639
)

type sitemapDoc struct {
	XMLName xml.Name       `xml:"urlset"`
	NS      string         `xml:"xmlns,attr"`
	NewsNS  string         `xml:"xmlns:news,attr,omitempty"`
	URLs    []sitemapEntry `xml:"url"`
}

type sitemapEntry struct {
	Loc     string       `xml:"loc"`
	LastMod string       `xml:"lastmod,omitempty"`
	News    *sitemapNews `xml:"news:news"`
}

type sitemapNews struct {
	Name            string `xml:"news:publication>news:name"`
	Language        string `xml:"news:publication>news:language"`
	PublicationDate string `xml:"news:publication_date"`
	Title           string `xml:"news:title"`
}

// writeSitemap writes /sitemap.xml: every indexable page the build wrote
// (its canonical URL, once) with lastmod. See indexWrittenPages for what
// counts and where lastmod comes from.
func writeSitemap(s *site) error {
	doc := sitemapDoc{NS: sitemapNS}
	for _, p := range s.index.pages {
		doc.URLs = append(doc.URLs, sitemapEntry{Loc: p.canonical, LastMod: p.lastmod})
	}
	return s.writeXML(sitemapURL, doc)
}

// writeNewsSitemap writes /news-sitemap.xml: the Briefs published in the
// 48 hours up to the build clock (s.cfg.Now), newest first, at most
// newsMaxEntries. Explainers are evergreen, not news, and never appear.
// With no recent Brief it is an empty urlset (Google: a warning, no harm).
func writeNewsSitemap(s *site) error {
	doc := sitemapDoc{NS: sitemapNS, NewsNS: newsSitemapNS}
	since := s.cfg.Now.Add(-newsWindow)
	for _, b := range newestBriefs(s.publishedBriefs()) {
		if len(doc.URLs) == newsMaxEntries {
			break
		}
		if b.Published.Before(since) || b.Published.After(s.cfg.Now) || !s.isWritten(briefURL(b)) {
			continue
		}
		doc.URLs = append(doc.URLs, sitemapEntry{
			Loc: s.absURL(briefURL(b)),
			News: &sitemapNews{
				Name:            views.PublicationName,
				Language:        newsLanguage,
				PublicationDate: b.Published.Format(time.RFC3339),
				Title:           b.Title,
			},
		})
	}
	return s.writeXML(newsSitemapURL, doc)
}

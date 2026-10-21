package builder

import "encoding/xml"

func init() {
	registerStep(phaseWrite, 30, "write sitemap", writeSitemap)
}

// sitemapNS is the sitemap protocol's namespace (the news sitemap uses it
// too; see step_write_news_sitemap.go).
const sitemapNS = "http://www.sitemaps.org/schemas/sitemap/0.9"

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

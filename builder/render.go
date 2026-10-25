package builder

import (
	"bytes"
	"context"
	"encoding/xml"
	"hastcode.com/site/internal/views"
	"os"

	"github.com/a-h/templ"
)

// Output helpers shared by render and write steps: writePage (every HTML
// page, with the size gate), writeXML (feeds, sitemaps) and isWritten.

// HTML size gate (spec "Non-functional": HTML ≤ 50 KB uncompressed). It
// lives in writePage, the one path every page type writes through, so new
// page types are gated without opting in. The shared inline stylesheet
// (ADR-0006) does not count; inline JSON-LD, SVG and scripts do.
const (
	ruleHTMLSize  = "html-size"
	maxHTMLBytes  = 50 * 1024
	maxHTMLBudget = "50 KB"
)

// ruleMarkdown is a content body goldmark could not render (render steps
// for Explainers, Briefs, Author bios and static pages).
const ruleMarkdown = "markdown"

// JSON-LD presence gate (spec: the build fails on a page lacking JSON-LD).
// Every indexable page carries at least one block (Article/NewsArticle,
// BreadcrumbList, WebSite, Person…); noindex pages (404, /search/,
// /subscribed/) are exempt. Like the size gate it lives in writePage, so a
// new page type cannot ship without structured data.
const ruleJSONLDMissing = "jsonld-missing"

var (
	noindexMeta = []byte(`<meta name="robots" content="noindex"`)
	jsonLDBlock = []byte(`<script type="application/ld+json">`)
)

// writePage renders a templ component and writes it at urlPath
// ("/a/b/" -> static/a/b/index.html). A page over the size gate, or an
// indexable page without JSON-LD, is reported as a build error and not
// written.
func (s *site) writePage(urlPath string, c templ.Component) error {
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		return err
	}
	html := buf.Bytes()
	n := len(html)
	if bytes.Contains(html, views.InlineStyle()) {
		n -= len(views.InlineStyle())
	}
	if n > maxHTMLBytes {
		s.fail("", ruleHTMLSize, "page %s is %d bytes of HTML without the inline stylesheet, over the %s (%d bytes) uncompressed limit; trim it or split it",
			urlPath, n, maxHTMLBudget, maxHTMLBytes)
		return nil
	}
	if !bytes.Contains(html, noindexMeta) && !bytes.Contains(html, jsonLDBlock) {
		s.fail("", ruleJSONLDMissing, "page %s is indexable but has no JSON-LD; append at least s.breadcrumbs(page.Breadcrumbs) to page.StructuredData", urlPath)
		return nil
	}
	return s.out.WriteStatic(urlPath, html)
}

// writeXML writes v as an indented UTF-8 XML document at urlPath.
func (s *site) writeXML(urlPath string, v any) error {
	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(v); err != nil {
		return err
	}
	buf.WriteByte('\n')
	return s.out.WriteStatic(urlPath, buf.Bytes())
}

// isWritten reports whether the build wrote a file for urlPath ("/a/" is
// the directory index). Steps after the link gate use it so manifests
// list only pages that really exist.
func (s *site) isWritten(urlPath string) bool {
	p, err := s.out.StaticFile(urlPath)
	if err != nil {
		return false
	}
	_, err = os.Stat(p)
	return err == nil
}

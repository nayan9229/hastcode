package builder

import (
	"bytes"
	"context"
	"encoding/xml"
	"os"

	"github.com/a-h/templ"
)

// Output helpers shared by render and write steps: writePage (every HTML
// page, with the size gate), writeXML (feeds, sitemaps) and isWritten.

// HTML size gate (spec "Non-functional": HTML ≤ 50 KB uncompressed). It
// lives in writePage, the one path every page type writes through, so new
// page types are gated without opting in. CSS is an external file
// (ADR-0003) and does not count; inline JSON-LD, SVG and scripts do.
const (
	ruleHTMLSize  = "html-size"
	maxHTMLBytes  = 50 * 1024
	maxHTMLBudget = "50 KB"
)

// ruleMarkdown is a content body goldmark could not render (render steps
// for Explainers, Briefs, Author bios and static pages).
const ruleMarkdown = "markdown"

// writePage renders a templ component and writes it at urlPath
// ("/a/b/" -> static/a/b/index.html). A page over the size gate is
// reported as a build error and not written.
func (s *site) writePage(urlPath string, c templ.Component) error {
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		return err
	}
	if n := buf.Len(); n > maxHTMLBytes {
		s.fail("", ruleHTMLSize, "page %s is %d bytes of HTML, over the %s (%d bytes) uncompressed limit; trim it or split it",
			urlPath, n, maxHTMLBudget, maxHTMLBytes)
		return nil
	}
	return s.out.WriteStatic(urlPath, buf.Bytes())
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

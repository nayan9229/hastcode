package builder

import (
	"bytes"
	"context"

	"github.com/a-h/templ"
)

func init() {
	registerStep(phaseRender, 0, "reset output", func(s *site) error { return s.out.Reset() })
}

// HTML size gate (spec "Non-functional": HTML ≤ 50 KB uncompressed). It
// lives in writePage, the one path every page type writes through, so new
// page types are gated without opting in. CSS is an external file
// (ADR-0003) and does not count; inline JSON-LD, SVG and scripts do.
const (
	ruleHTMLSize  = "html-size"
	maxHTMLBytes  = 50 * 1024
	maxHTMLBudget = "50 KB"
)

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

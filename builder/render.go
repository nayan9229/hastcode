package builder

import (
	"bytes"
	"context"

	"github.com/a-h/templ"
)

func init() {
	registerStep(phaseRender, 0, "reset output", func(s *site) error { return s.out.Reset() })
}

// writePage renders a templ component and writes it at urlPath
// ("/a/b/" -> static/a/b/index.html).
func (s *site) writePage(urlPath string, c templ.Component) error {
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		return err
	}
	return s.out.WriteStatic(urlPath, buf.Bytes())
}

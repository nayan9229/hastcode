package builder

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/markdown"
	"hastcode.com/site/internal/mermaid"
	"hastcode.com/site/internal/views"
)

func init() {
	registerStep(phaseTransform, 10, "render diagrams", renderDiagrams)
}

// ruleDiagram: a Mermaid fence the renderer cannot draw (malformed, empty,
// or an unsupported type such as gantt or pie).
const ruleDiagram = "diagram"

// renderDiagrams replaces every ```mermaid fence in published Explainers
// with an inline-SVG <figure>. All fences in the build go to the renderer
// in one batch; cached diagrams never reach it.
func renderDiagrams(s *site) error {
	type fence struct {
		e *content.Explainer
		n int // 1-based position on the page
		d *markdown.Diagram
	}
	var fences []fence
	var sources []string
	for _, e := range s.publishedExplainers() {
		for i, d := range e.Body.Diagrams() {
			fences = append(fences, fence{e, i + 1, d})
			sources = append(sources, d.Code)
		}
	}
	if len(sources) == 0 {
		return nil
	}

	r := mermaid.Renderer{Node: s.cfg.Node, ToolsDir: s.cfg.ToolsDir, CacheDir: filepath.Join(s.cfg.CacheDir, "mermaid")}
	results, err := r.Render(context.Background(), sources)
	if err != nil {
		return err
	}
	for i, f := range fences {
		if err := results[i].Err; err != nil {
			s.fail(f.e.File, ruleDiagram, "diagram %d %q on %s: %v", f.n, f.d.Title, explainerURL(f.e), err)
			continue
		}
		fig := views.Diagram(views.DiagramFigure{
			SVG:     mermaid.PrefixIDs(results[i].SVG, fmt.Sprintf("d%d-", f.n)),
			Caption: f.d.Title,
		})
		var buf bytes.Buffer
		if err := fig.Render(context.Background(), &buf); err != nil {
			return err
		}
		f.e.Body.ReplaceDiagram(f.d, buf.String())
	}
	return nil
}

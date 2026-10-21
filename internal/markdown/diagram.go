package markdown

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// Diagram is one ```mermaid fence in a Document. Gates count them (the
// "diagram" skeleton rule) and the diagram transform replaces them with
// rendered figures.
type Diagram struct {
	// Title is the caption, from a leading "%% title: …" line ("" if none).
	Title string
	// Code is the Mermaid source with the title line removed.
	Code string

	node *ast.FencedCodeBlock
}

const titlePrefix = "%% title:"

// Diagrams returns the document's Mermaid fences in document order.
func (d *Document) Diagrams() []*Diagram {
	var out []*Diagram
	_ = ast.Walk(d.Root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		fence, ok := n.(*ast.FencedCodeBlock)
		if !entering || !ok {
			return ast.WalkContinue, nil
		}
		if string(fence.Language(d.Source)) != "mermaid" {
			return ast.WalkSkipChildren, nil
		}
		var code bytes.Buffer
		lines := fence.Lines()
		for i := 0; i < lines.Len(); i++ {
			seg := lines.At(i)
			code.Write(seg.Value(d.Source))
		}
		dg := &Diagram{node: fence}
		src := strings.TrimSpace(code.String())
		first, rest, _ := strings.Cut(src, "\n")
		if title, ok := strings.CutPrefix(strings.TrimSpace(first), titlePrefix); ok {
			dg.Title = strings.TrimSpace(title)
			src = rest
		}
		dg.Code = strings.TrimSpace(src)
		out = append(out, dg)
		return ast.WalkSkipChildren, nil
	})
	return out
}

// ReplaceDiagram swaps the fence for html, which Render emits verbatim.
func (d *Document) ReplaceDiagram(dg *Diagram, html string) {
	parent := dg.node.Parent()
	if parent == nil {
		return // already replaced
	}
	parent.ReplaceChild(parent, dg.node, &rawHTML{html: html})
}

// rawHTML is a block the builder generated (not Author-written HTML).
type rawHTML struct {
	ast.BaseBlock
	html string
}

var kindRawHTML = ast.NewNodeKind("BuilderHTML")

func (n *rawHTML) Kind() ast.NodeKind { return kindRawHTML }

func (n *rawHTML) Dump(source []byte, level int) { ast.DumpHelper(n, source, level, nil, nil) }

type rawHTMLRenderer struct{}

func (rawHTMLRenderer) RegisterFuncs(r renderer.NodeRendererFuncRegisterer) {
	r.Register(kindRawHTML, func(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			_, _ = w.WriteString(n.(*rawHTML).html)
			_ = w.WriteByte('\n')
		}
		return ast.WalkSkipChildren, nil
	})
}

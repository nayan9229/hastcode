package markdown

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
)

// Prose is one block of running prose: a paragraph, or the text of a list
// item, as the Reader reads it. Literal marks the bytes of Text that came
// from inline code or an autolink's URL: the Reader reads them, but their
// words are not ours to check.
type Prose struct {
	Text    string
	Literal []bool
}

// Prose returns the running prose in nodes, in order: every paragraph and
// list item text, skipping headings, tables, blockquotes, code blocks
// (Mermaid included), raw HTML and images. A link keeps its text and drops
// its URL.
func (d *Document) Prose(nodes []ast.Node) []Prose {
	var out []Prose
	for _, n := range nodes {
		_ = ast.Walk(n, func(c ast.Node, entering bool) (ast.WalkStatus, error) {
			if !entering {
				return ast.WalkContinue, nil
			}
			switch c.(type) {
			case *ast.Heading, *ast.Blockquote, *east.Table, *ast.FencedCodeBlock, *ast.CodeBlock, *ast.HTMLBlock:
				return ast.WalkSkipChildren, nil
			case *ast.Paragraph, *ast.TextBlock:
				var p proseBuilder
				p.inline(d, c)
				if t := strings.TrimSpace(p.text.String()); t != "" {
					out = append(out, p.prose())
				}
				return ast.WalkSkipChildren, nil
			}
			return ast.WalkContinue, nil
		})
	}
	return out
}

type proseBuilder struct {
	text    strings.Builder
	literal []bool
}

func (p *proseBuilder) write(s string, literal bool) {
	p.text.WriteString(s)
	for range len(s) {
		p.literal = append(p.literal, literal)
	}
}

// inline appends the text of n's inline children.
func (p *proseBuilder) inline(d *Document, n ast.Node) {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		switch c := c.(type) {
		case *ast.Text:
			p.write(string(c.Segment.Value(d.Source)), false)
			if c.SoftLineBreak() || c.HardLineBreak() {
				p.write(" ", false)
			}
		case *ast.String:
			p.write(string(c.Value), false)
		case *ast.CodeSpan:
			var code proseBuilder
			code.inline(d, c)
			p.write(code.text.String(), true)
		case *ast.AutoLink:
			p.write(string(c.URL(d.Source)), true)
		case *ast.Image, *ast.RawHTML:
		default:
			p.inline(d, c)
		}
	}
}

// prose trims surrounding space, keeping Literal aligned with Text.
func (p *proseBuilder) prose() Prose {
	s := p.text.String()
	start := len(s) - len(strings.TrimLeft(s, " \t\n"))
	end := len(strings.TrimRight(s, " \t\n"))
	return Prose{Text: s[start:end], Literal: p.literal[start:end]}
}

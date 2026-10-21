package markdown

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
)

// Section is an H2 heading and the top-level blocks under it, up to the
// next H1 or H2. H3+ headings stay inside their section.
type Section struct {
	Heading string
	Blocks  []ast.Node
}

// Sections splits the document at its H2 headings. intro is the blocks
// before the first H2 (or H1); sections are the H2 sections in order.
func (d *Document) Sections() (intro []ast.Node, sections []Section) {
	var cur *Section
	for n := d.Root.FirstChild(); n != nil; n = n.NextSibling() {
		if h, ok := n.(*ast.Heading); ok && h.Level <= 2 {
			if cur != nil {
				sections = append(sections, *cur)
			}
			cur = nil
			if h.Level == 2 {
				cur = &Section{Heading: d.Text(h)}
			}
			continue
		}
		if cur != nil {
			cur.Blocks = append(cur.Blocks, n)
		} else if len(sections) == 0 {
			intro = append(intro, n)
		}
	}
	if cur != nil {
		sections = append(sections, *cur)
	}
	return intro, sections
}

// Text returns the plain text of n's inline content (code blocks and raw
// HTML excluded), with blocks separated by spaces.
func (d *Document) Text(n ast.Node) string {
	var b strings.Builder
	_ = ast.Walk(n, func(c ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch c := c.(type) {
		case *ast.FencedCodeBlock, *ast.CodeBlock, *ast.HTMLBlock, *ast.RawHTML:
			return ast.WalkSkipChildren, nil
		case *ast.Text:
			b.Write(c.Segment.Value(d.Source))
			if c.SoftLineBreak() || c.HardLineBreak() {
				b.WriteByte(' ')
			}
		case *ast.String:
			b.Write(c.Value)
		case *ast.AutoLink:
			b.Write(c.URL(d.Source))
		default:
			if c.Type() == ast.TypeBlock && b.Len() > 0 {
				b.WriteByte(' ')
			}
		}
		return ast.WalkContinue, nil
	})
	return strings.TrimSpace(b.String())
}

// Words counts the words of prose in nodes (code blocks excluded).
func (d *Document) Words(nodes []ast.Node) int {
	total := 0
	for _, n := range nodes {
		total += len(strings.Fields(d.Text(n)))
	}
	return total
}

// Find reports whether any node in nodes (or their descendants) satisfies
// match.
func Find(nodes []ast.Node, match func(ast.Node) bool) bool {
	found := false
	for _, n := range nodes {
		_ = ast.Walk(n, func(c ast.Node, entering bool) (ast.WalkStatus, error) {
			if entering && match(c) {
				found = true
				return ast.WalkStop, nil
			}
			return ast.WalkContinue, nil
		})
		if found {
			return true
		}
	}
	return false
}

// IsLink matches inline links and autolinks.
func IsLink(n ast.Node) bool {
	switch n.(type) {
	case *ast.Link, *ast.AutoLink:
		return true
	}
	return false
}

// IsTable matches GFM tables.
func IsTable(n ast.Node) bool {
	_, ok := n.(*east.Table)
	return ok
}

// FenceLanguage returns the info-string language of a fenced code block,
// or "" if n is not one.
func (d *Document) FenceLanguage(n ast.Node) string {
	if f, ok := n.(*ast.FencedCodeBlock); ok {
		return string(f.Language(d.Source))
	}
	return ""
}

// SplitAfterSection splits the top-level blocks at the end of the H2
// section named heading: before runs from the start through that section,
// after is everything from the next H1/H2 on. Without such a section it
// splits after the intro. Layouts use it to place content (an ad slot)
// between the two halves.
func (d *Document) SplitAfterSection(heading string) (before, after []ast.Node) {
	var all []ast.Node
	for n := d.Root.FirstChild(); n != nil; n = n.NextSibling() {
		all = append(all, n)
	}
	cut, in := -1, false
	for i, n := range all {
		h, ok := n.(*ast.Heading)
		if !ok || h.Level > 2 {
			continue
		}
		if cut < 0 {
			cut = i // end of the intro: the fallback
		}
		if in {
			return all[:i], all[i:]
		}
		in = h.Level == 2 && d.Text(h) == heading
	}
	if in {
		return all, nil
	}
	if cut < 0 {
		cut = len(all)
	}
	return all[:cut], all[cut:]
}

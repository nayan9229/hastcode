package markdown

import "github.com/yuin/goldmark/ast"

// SplitAtSection splits top-level blocks (from SplitAfterSection or a
// walk of the root) at the H2 named heading: before is everything up to
// it, from starts with that heading. Without such a heading, from is nil.
// Layouts use it to wrap a trailing section (Sources) in its own element.
func (d *Document) SplitAtSection(blocks []ast.Node, heading string) (before, from []ast.Node) {
	for i, n := range blocks {
		if h, ok := n.(*ast.Heading); ok && h.Level == 2 && d.Text(h) == heading {
			return blocks[:i], blocks[i:]
		}
	}
	return blocks, nil
}

package markdown

import "github.com/yuin/goldmark/ast"

// InsertHTMLAfterSection places html (a block the builder generated, emitted
// verbatim like a replaced diagram) at the end of the H2 section named
// heading: after its last block, before the next H1/H2. Page, feeds and
// the search index all see it. It reports whether the section exists; the
// document is unchanged when it does not. Sections() keeps treating the
// inserted block as part of the section it ends, so a builder-added H2 in
// html never counts as an Author-written heading.
func (d *Document) InsertHTMLAfterSection(heading, html string) bool {
	var last ast.Node
	in := false
	for n := d.Root.FirstChild(); n != nil; n = n.NextSibling() {
		if h, ok := n.(*ast.Heading); ok && h.Level <= 2 {
			if in {
				break
			}
			in = h.Level == 2 && d.Text(h) == heading
			if in {
				last = n
			}
			continue
		}
		if in {
			last = n
		}
	}
	if last == nil {
		return false
	}
	d.Root.InsertAfter(d.Root, last, &rawHTML{html: html})
	return true
}

package markdown

import "github.com/yuin/goldmark/ast"

// Link is a link the builder adds to a document.
type Link struct {
	Text, URL string
}

// LinkDestinations returns the destination of every inline link and
// autolink in nodes (and their descendants), in document order.
func (d *Document) LinkDestinations(nodes []ast.Node) []string {
	var out []string
	for _, n := range nodes {
		_ = ast.Walk(n, func(c ast.Node, entering bool) (ast.WalkStatus, error) {
			if !entering {
				return ast.WalkContinue, nil
			}
			switch l := c.(type) {
			case *ast.Link:
				out = append(out, string(l.Destination))
			case *ast.AutoLink:
				out = append(out, string(l.URL(d.Source)))
			}
			return ast.WalkContinue, nil
		})
	}
	return out
}

// AppendLinkList appends a bulleted list of links as the document's last
// block (link text is escaped when rendered, like any Markdown text).
func (d *Document) AppendLinkList(links []Link) {
	if len(links) == 0 {
		return
	}
	list := ast.NewList('-')
	list.IsTight = true
	for _, l := range links {
		a := ast.NewLink()
		a.Destination = []byte(l.URL)
		a.AppendChild(a, ast.NewString([]byte(l.Text)))
		tb := ast.NewTextBlock()
		tb.AppendChild(tb, a)
		item := ast.NewListItem(2)
		item.AppendChild(item, tb)
		list.AppendChild(list, item)
	}
	d.Root.AppendChild(d.Root, list)
}

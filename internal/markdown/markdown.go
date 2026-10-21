// Package markdown parses piece bodies into a goldmark AST and renders that
// AST to HTML. Parsing and rendering are separate so build steps can inspect
// (gates) and rewrite (diagrams) the tree in between.
package markdown

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

// Document is a parsed Markdown body. Root nodes reference Source by
// segment, so keep them together.
type Document struct {
	Source []byte
	Root   ast.Node
}

var md = goldmark.New(goldmark.WithExtensions(extension.GFM))

// Parse parses a Markdown body (frontmatter already removed).
func Parse(src []byte) *Document {
	return &Document{Source: src, Root: md.Parser().Parse(text.NewReader(src))}
}

// Render renders the (possibly rewritten) document to HTML.
func Render(doc *Document) (string, error) {
	var buf bytes.Buffer
	if err := md.Renderer().Render(&buf, doc.Source, doc.Root); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderBlocks renders some of doc's top-level blocks (e.g. one half of
// SplitAfterSection) to HTML. Rendering every top-level block in order
// gives the same HTML as Render.
func RenderBlocks(doc *Document, blocks []ast.Node) (string, error) {
	var buf bytes.Buffer
	for _, b := range blocks {
		if err := md.Renderer().Render(&buf, doc.Source, b); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

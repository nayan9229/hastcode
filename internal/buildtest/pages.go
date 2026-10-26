package buildtest

import (
	"path/filepath"
	"strings"
)

// StaticPage is a fixture static page (About, Contact, Editorial Policy),
// written to pages/<slug>.md and published at /<slug>/. Zero-value fields
// are omitted from the frontmatter.
type StaticPage struct {
	Slug        string
	Title       string
	Description string
	Body        string
}

// TransparencyPages are the static pages the site chrome links from every
// page: About, Contact, Editorial Policy and Privacy.
func TransparencyPages() []StaticPage {
	return []StaticPage{
		{Slug: "about", Title: "About HastCode", Description: "Who makes HastCode and why.", Body: "HastCode explains how real systems are built.\n"},
		{Slug: "contact", Title: "Contact", Description: "How to reach HastCode.", Body: "Email [the Author](mailto:author@hastcode.test).\n"},
		{Slug: "editorial-policy", Title: "Editorial Policy", Description: "How HastCode researches, writes and corrects.", Body: "Every claim links a first-party source.\n"},
		{Slug: "privacy", Title: "Privacy", Description: "What HastCode and its providers collect.", Body: "Google Analytics measures visits.\n"},
	}
}

// WithPage adds a static page at pages/<slug>.md.
func (t Tree) WithPage(p StaticPage) Tree {
	t[filepath.Join("pages", p.Slug+".md")] = []byte(p.Markdown())
	return t
}

// Markdown renders the fixture as a content file.
func (p StaticPage) Markdown() string {
	var b strings.Builder
	b.WriteString("---\n")
	field(&b, "title", p.Title)
	field(&b, "description", p.Description)
	b.WriteString("---\n\n")
	b.WriteString(p.Body)
	return b.String()
}

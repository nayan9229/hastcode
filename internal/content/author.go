package content

import (
	"path/filepath"
	"strings"

	"hastcode.com/site/internal/markdown"
)

// Author is a real, named person accountable for Explainers and Briefs.
// File is relative to the content directory; Slug is its basename.
type Author struct {
	File string
	Slug string
	Name string
	Bio  *markdown.Document
}

type authorFrontmatter struct {
	Name string `yaml:"name"`
}

// ParseAuthor parses an Author file (content/authors/<slug>.md).
func ParseAuthor(file string, data []byte) (*Author, error) {
	front, body, err := splitFrontmatter(data)
	if err != nil {
		return nil, err
	}
	var fm authorFrontmatter
	if err := decodeStrict(front, &fm); err != nil {
		return nil, err
	}
	return &Author{
		File: file,
		Slug: strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)),
		Name: fm.Name,
		Bio:  markdown.Parse(body),
	}, nil
}

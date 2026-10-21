package content

import (
	"errors"
	"path/filepath"
	"strings"

	"hastcode.com/site/internal/markdown"
)

// Page is a static page (About, Contact, Editorial Policy) from
// content/pages/<slug>.md, published at /<slug>/. File is relative to the
// content directory; Slug is its basename.
type Page struct {
	File        string
	Slug        string
	Title       string
	Description string
	Body        *markdown.Document
}

type pageFrontmatter struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

// ParsePage parses a static page file. title and description are required.
func ParsePage(file string, data []byte) (*Page, error) {
	front, body, err := splitFrontmatter(data)
	if err != nil {
		return nil, err
	}
	var fm pageFrontmatter
	if err := decodeStrict(front, &fm); err != nil {
		return nil, err
	}
	if strings.TrimSpace(fm.Title) == "" {
		return nil, errors.New("frontmatter: title is required")
	}
	if strings.TrimSpace(fm.Description) == "" {
		return nil, errors.New("frontmatter: description is required")
	}
	return &Page{
		File:        file,
		Slug:        strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)),
		Title:       fm.Title,
		Description: fm.Description,
		Body:        markdown.Parse(body),
	}, nil
}

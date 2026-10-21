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
	// SameAs are the Author's profiles elsewhere (GitHub, LinkedIn, …):
	// absolute URLs, the Person's sameAs in JSON-LD.
	SameAs []string
	// Photo is a raster beside the Markdown (authors/<photo>); "" = none.
	// PhotoAlt is its alt text.
	Photo    string
	PhotoAlt string
	Bio      *markdown.Document
}

type authorFrontmatter struct {
	Name     string   `yaml:"name"`
	SameAs   []string `yaml:"same_as"`
	Photo    string   `yaml:"photo"`
	PhotoAlt string   `yaml:"photo_alt"`
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
		File:     file,
		Slug:     strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)),
		Name:     fm.Name,
		SameAs:   fm.SameAs,
		Photo:    fm.Photo,
		PhotoAlt: fm.PhotoAlt,
		Bio:      markdown.Parse(body),
	}, nil
}

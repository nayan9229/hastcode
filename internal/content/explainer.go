package content

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"hastcode.com/site/internal/markdown"
)

// Kind is the Explainer kind.
type Kind string

const (
	CaseStudy  Kind = "case-study"
	Concept    Kind = "concept"
	Comparison Kind = "comparison"
)

// Source is a frontmatter source entry (in addition to inline links).
type Source struct {
	Title string `yaml:"title"`
	URL   string `yaml:"url"`
}

// Explainer is an evergreen, long-form article teaching one architectural
// idea. File is relative to the content directory.
type Explainer struct {
	File        string
	Slug        string
	Title       string
	Description string
	Kind        Kind
	Vertical    string
	Tags        []string
	Published   time.Time
	Reviewed    time.Time
	Author      string
	Hero        string
	HeroAlt     string
	Related     []string
	Draft       bool
	Sources     []Source
	Body        *markdown.Document
}

type explainerFrontmatter struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Kind        Kind     `yaml:"kind"`
	Vertical    string   `yaml:"vertical"`
	Tags        []string `yaml:"tags"`
	Published   string   `yaml:"published"`
	Reviewed    string   `yaml:"reviewed"`
	Author      string   `yaml:"author"`
	Hero        string   `yaml:"hero"`
	HeroAlt     string   `yaml:"hero_alt"`
	Related     []string `yaml:"related"`
	Draft       bool     `yaml:"draft"`
	Sources     []Source `yaml:"sources"`
}

// ParseExplainer parses an Explainer file. file is its content-relative
// path; the slug is the file's basename without extension.
func ParseExplainer(file string, data []byte) (*Explainer, error) {
	front, body, err := splitFrontmatter(data)
	if err != nil {
		return nil, err
	}
	var fm explainerFrontmatter
	if err := decodeStrict(front, &fm); err != nil {
		return nil, err
	}
	switch fm.Kind {
	case CaseStudy, Concept, Comparison:
	default:
		return nil, fmt.Errorf("frontmatter: kind %q must be one of case-study, concept, comparison", fm.Kind)
	}
	published, err := parseTime("published", fm.Published)
	if err != nil {
		return nil, err
	}
	reviewed, err := parseTime("reviewed", fm.Reviewed)
	if err != nil {
		return nil, err
	}
	return &Explainer{
		File:        file,
		Slug:        strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)),
		Title:       fm.Title,
		Description: fm.Description,
		Kind:        fm.Kind,
		Vertical:    fm.Vertical,
		Tags:        fm.Tags,
		Published:   published,
		Reviewed:    reviewed,
		Author:      fm.Author,
		Hero:        fm.Hero,
		HeroAlt:     fm.HeroAlt,
		Related:     fm.Related,
		Draft:       fm.Draft,
		Sources:     fm.Sources,
		Body:        markdown.Parse(body),
	}, nil
}

package content

import (
	"path/filepath"
	"strings"
	"time"

	"hastcode.com/site/internal/markdown"
)

// Brief is a short, dated commentary on something that just happened in
// tech, with the architectural "why it matters", linking into Explainers.
// It has no kind, Vertical or review date. File is relative to the content
// directory; the slug is its basename.
type Brief struct {
	File        string
	Slug        string
	Title       string
	Description string
	Tags        []string
	Published   time.Time
	Author      string
	Hero        string
	HeroAlt     string
	Related     []string
	Draft       bool
	Sources     []Source
	// SourceURL is the public source the Brief is curated from.
	SourceURL string
	// Explainers are the slugs of the Explainers the body links to.
	Explainers []string
	Body       *markdown.Document
}

type briefFrontmatter struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags"`
	Published   string   `yaml:"published"`
	Author      string   `yaml:"author"`
	Hero        string   `yaml:"hero"`
	HeroAlt     string   `yaml:"hero_alt"`
	Related     []string `yaml:"related"`
	Draft       bool     `yaml:"draft"`
	Sources     []Source `yaml:"sources"`
	SourceURL   string   `yaml:"source_url"`
	Explainers  []string `yaml:"explainers"`
}

// ParseBrief parses a Brief file (content/briefs/<slug>.md). Explainer-only
// keys (kind, vertical, reviewed) are unknown keys and rejected.
func ParseBrief(file string, data []byte) (*Brief, error) {
	front, body, err := splitFrontmatter(data)
	if err != nil {
		return nil, err
	}
	var fm briefFrontmatter
	if err := decodeStrict(front, &fm); err != nil {
		return nil, err
	}
	published, err := parseTime("published", fm.Published)
	if err != nil {
		return nil, err
	}
	return &Brief{
		File:        file,
		Slug:        strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)),
		Title:       fm.Title,
		Description: fm.Description,
		Tags:        fm.Tags,
		Published:   published,
		Author:      fm.Author,
		Hero:        fm.Hero,
		HeroAlt:     fm.HeroAlt,
		Related:     fm.Related,
		Draft:       fm.Draft,
		Sources:     fm.Sources,
		SourceURL:   fm.SourceURL,
		Explainers:  fm.Explainers,
		Body:        markdown.Parse(body),
	}, nil
}

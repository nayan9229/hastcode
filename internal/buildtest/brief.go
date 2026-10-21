package buildtest

import (
	"path/filepath"
	"strings"
)

// Brief is a fixture Brief. Zero-value fields are omitted from the
// frontmatter, so tests can blank a field to make it invalid.
type Brief struct {
	Slug        string
	Title       string
	Description string
	Tags        []string
	Published   string
	Author      string
	Hero        string
	HeroAlt     string
	HeroImage   []byte
	SourceURL   string
	Explainers  []string
	// Sources is the optional `sources` frontmatter list.
	Sources []Source
	Draft   bool
	Body    string
}

// ValidBrief returns a Brief that satisfies the Brief skeleton: under 250
// words, exactly "What happened" and "Why it matters", a source_url, a
// 1600x900 hero, and a body link to ValidExplainer() (so add that
// Explainer to the tree too; ValidTree() has it).
func ValidBrief() Brief {
	return Brief{
		Slug:        "etcd-raft-rewrite",
		Title:       "etcd ships a rewritten Raft library",
		Description: "etcd's new Raft library cuts election latency; here is why leader election design matters.",
		Tags:        []string{"raft", "etcd"},
		Published:   "2026-10-20T08:15:00-04:00",
		Author:      DefaultAuthor().Slug,
		Hero:        "etcd-raft-rewrite.png",
		HeroAlt:     "The etcd logo above a Raft cluster",
		HeroImage:   HeroPNG(1600, 900),
		SourceURL:   "https://etcd.io/blog/2026/raft-rewrite/",
		Explainers:  []string{"raft-leader-election"},
		Body:        ValidBriefBody,
	}
}

// ValidBriefBody is the Markdown body of ValidBrief. It links the
// Explainer by bare slug.
const ValidBriefBody = `## What happened

etcd released a rewritten Raft library that halves election latency under load.

## Why it matters

Faster elections shrink the window in which a cluster cannot accept writes. See [how Raft elects a leader](raft-leader-election) for the mechanism.
`

// WithBrief adds a Brief at briefs/<slug>.md and, when b.Hero is set and
// HeroImage is non-nil, its hero raster beside it.
func (t Tree) WithBrief(b Brief) Tree {
	t[filepath.Join("briefs", b.Slug+".md")] = []byte(b.Markdown())
	if b.Hero != "" && b.HeroImage != nil {
		t[filepath.Join("briefs", b.Hero)] = b.HeroImage
	}
	return t
}

// Markdown renders the fixture as a content file.
func (b Brief) Markdown() string {
	var s strings.Builder
	s.WriteString("---\n")
	field(&s, "title", b.Title)
	field(&s, "description", b.Description)
	list(&s, "tags", b.Tags)
	field(&s, "published", b.Published)
	field(&s, "author", b.Author)
	field(&s, "hero", b.Hero)
	field(&s, "hero_alt", b.HeroAlt)
	field(&s, "source_url", b.SourceURL)
	list(&s, "explainers", b.Explainers)
	sources(&s, b.Sources)
	if b.Draft {
		s.WriteString("draft: true\n")
	}
	s.WriteString("---\n\n")
	s.WriteString(b.Body)
	return s.String()
}

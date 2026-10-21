package builder

import (
	"net/url"
	"slices"
	"strings"

	"github.com/yuin/goldmark/ast"

	"hastcode.com/site/internal/content"
)

func init() {
	registerStep(phaseValidate, 12, "validate brief references", validateBriefReferences)
	registerStep(phaseValidate, 35, "validate brief skeleton", validateBriefSkeleton)
}

// validateBriefReferences checks that every Brief's Author has a file
// under authors/. (Explainer links are checked with the skeleton.)
func validateBriefReferences(s *site) error {
	for _, b := range s.briefs {
		if _, ok := s.authors[b.Author]; !ok {
			s.fail(b.File, ruleUnknownAuthor, "unknown Author %q (no authors/%s.md)", b.Author, b.Author)
		}
	}
	return nil
}

// Rules for the Brief skeleton (spec "Content model").
const (
	ruleBriefLength        = "brief-length"
	ruleBriefSections      = "brief-sections"
	ruleBriefSource        = "brief-source"
	ruleBriefExplainerLink = "brief-explainer-link"
)

// Brief skeleton: the H2 sections a Brief has, exactly these and in this
// order, and its word limit (headings included).
const (
	headingWhatHappened = "What happened"
	headingWhyItMatters = "Why it matters"
	maxBriefWords       = 250
)

// validateBriefSkeleton enforces the Brief editorial floor on every Brief.
// Every broken rule is reported, not just the first.
func validateBriefSkeleton(s *site) error {
	for _, b := range s.briefs {
		checkBriefLength(s, b)
		checkBriefSections(s, b)
		checkBriefSource(s, b)
		checkBriefExplainerLinks(s, b)
	}
	return nil
}

// checkBriefExplainerLinks: every Explainer link in the body names a loaded
// Explainer (otherwise it is a broken internal link), there is at least
// one, and the `explainers` frontmatter lists exactly the linked slugs.
func checkBriefExplainerLinks(s *site, b *content.Brief) {
	var linked []string
	for _, l := range s.explainerLinks(b) {
		if l.target == nil {
			s.fail(b.File, ruleBrokenLink, "link %s points at no Explainer: there is no Explainer %q", l.node.Destination, l.slug)
			continue
		}
		if !slices.Contains(linked, l.slug) {
			linked = append(linked, l.slug)
		}
	}
	for _, slug := range b.Explainers {
		switch {
		case !slices.ContainsFunc(s.explainers, func(e *content.Explainer) bool { return e.Slug == slug }):
			s.fail(b.File, ruleBrokenLink, "explainers lists %q, but there is no such Explainer", slug)
		case !slices.Contains(linked, slug):
			s.fail(b.File, ruleBriefExplainerLink, "explainers lists %q, but the body does not link to it", slug)
		}
	}
	for _, slug := range linked {
		if !slices.Contains(b.Explainers, slug) {
			s.fail(b.File, ruleBriefExplainerLink, "the body links to Explainer %q; add it to explainers in frontmatter", slug)
		}
	}
	if len(linked) == 0 {
		s.fail(b.File, ruleBriefExplainerLink, "the body must link to at least one Explainer (by slug, e.g. [text](raft-leader-election))")
	}
}

// checkBriefSource: source_url is an absolute http(s) URL (Briefs are
// curated from public sources, never original reporting).
func checkBriefSource(s *site, b *content.Brief) {
	if b.SourceURL == "" {
		s.fail(b.File, ruleBriefSource, "missing source_url: a Brief links the public source it is curated from")
		return
	}
	u, err := url.Parse(b.SourceURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		s.fail(b.File, ruleBriefSource, "source_url %q must be an absolute http(s) URL", b.SourceURL)
	}
}

// checkBriefSections: the H2 sections are exactly "What happened" then
// "Why it matters". An intro before them is allowed.
func checkBriefSections(s *site, b *content.Brief) {
	_, sections := b.Body.Sections()
	var got []string
	for _, sec := range sections {
		got = append(got, sec.Heading)
	}
	want := []string{headingWhatHappened, headingWhyItMatters}
	if slices.Equal(got, want) {
		return
	}
	const shape = "a Brief has exactly ## " + headingWhatHappened + " then ## " + headingWhyItMatters
	reported := false
	for _, h := range want {
		if !slices.Contains(got, h) {
			s.fail(b.File, ruleBriefSections, "missing ## %s; %s", h, shape)
			reported = true
		}
	}
	for _, h := range got {
		if !slices.Contains(want, h) {
			s.fail(b.File, ruleBriefSections, "unexpected section ## %s; %s", h, shape)
			reported = true
		}
	}
	if !reported {
		s.fail(b.File, ruleBriefSections, "sections are %s; %s", strings.Join(got, ", "), shape)
	}
}

// checkBriefLength: the body is at most maxBriefWords words.
func checkBriefLength(s *site, b *content.Brief) {
	if n := b.Body.Words([]ast.Node{b.Body.Root}); n > maxBriefWords {
		s.fail(b.File, ruleBriefLength, "body is %d words; a Brief is at most %d words", n, maxBriefWords)
	}
}

package builder

import (
	"fmt"
	"slices"
	"strings"

	"hastcode.com/site/internal/content"
)

func init() {
	registerStep(phaseValidate, 32, "validate repo profile skeleton", validateRepoProfileSkeleton)
}

// Rules for the Repo Profile skeleton (spec "Page composition and gates").
const (
	// ruleRepoProfileSections is a missing, extra or reordered H2 heading:
	// every Profile has exactly repoProfileHeadings, in that order.
	ruleRepoProfileSections = "repo-profile-sections"
	// ruleRepoProfileFactSection is a hand-written Signals or Changelog
	// section: fact sections render from the Fact Pack, never from prose.
	ruleRepoProfileFactSection = "repo-profile-fact-section"
	// ruleRepoProfileFAQ is a FAQ with fewer than minFAQQuestions or more
	// than maxFAQQuestions ### questions.
	ruleRepoProfileFAQ = "repo-profile-faq"
)

// A Repo Profile's FAQ answers five to eight questions phrased the way a
// Reader searches, one ### heading each (spec User Story 10).
const (
	minFAQQuestions = 5
	maxFAQQuestions = 8
)

// The Repo Profile's H2 headings. The opinion sections are written in the
// Markdown body in this order; the fact sections (Signals, Changelog) are
// never written: the builder renders them from the Fact Pack.
const (
	headingIn60Seconds  = "In 60 seconds"
	headingArchitecture = "Architecture"
	headingUseItIf      = "Use it if"
	headingDontUseItIf  = "Don't use it if"
	headingAlternatives = "Alternatives"
	headingFAQ          = "FAQ"
	headingVerdict      = "Verdict"
	headingSignals      = "Signals"
	headingChangelog    = "Changelog"
)

// repoProfileFactSections are the H2 sections the build renders from the
// Fact Pack, each with the written section it follows on the page.
var repoProfileFactSections = map[string]string{
	headingSignals:   headingArchitecture,
	headingChangelog: headingAlternatives,
}

// repoProfileHeadings is the Markdown body of a Repo Profile: exactly these
// H2 sections, in this order (TL;DR and Sources keep their Explainer rules).
var repoProfileHeadings = []string{
	headingTLDR,
	headingIn60Seconds,
	headingArchitecture,
	headingUseItIf,
	headingDontUseItIf,
	headingAlternatives,
	headingFAQ,
	headingVerdict,
	headingSources,
}

// validateRepoProfileSkeleton enforces the Repo Profile shape on every
// Repo Profile, so every Profile reads the same. Other kinds are untouched.
// Every broken rule is reported, not just the first.
func validateRepoProfileSkeleton(s *site) error {
	for _, e := range s.explainers {
		if e.Kind != content.RepoProfile {
			continue
		}
		checkRepoProfileFactSections(s, e)
		checkRepoProfileSections(s, e)
		checkRepoProfileFAQ(s, e)
	}
	return nil
}

// checkRepoProfileFAQ: the ## FAQ section holds minFAQQuestions to
// maxFAQQuestions ### questions. A missing FAQ is the sections rule's.
func checkRepoProfileFAQ(s *site, e *content.Explainer) {
	questions, ok := faqQuestions(e.Body)
	if !ok {
		return
	}
	n := len(questions)
	if n < minFAQQuestions || n > maxFAQQuestions {
		s.fail(e.File, ruleRepoProfileFAQ, "## FAQ has %d questions, want %d–%d (one ### heading per question, phrased the way a Reader searches)", n, minFAQQuestions, maxFAQQuestions)
	}
}

// checkRepoProfileFactSections: no Signals or Changelog section is written
// by hand; the build renders them from the Fact Pack.
func checkRepoProfileFactSections(s *site, e *content.Explainer) {
	_, sections := e.Body.Sections()
	for _, sec := range sections {
		if after, ok := repoProfileFactSections[sec.Heading]; ok {
			s.fail(e.File, ruleRepoProfileFactSection, "## %s is rendered from the Fact Pack (the build places it after ## %s): delete the hand-written section", sec.Heading, after)
		}
	}
}

// checkRepoProfileSections: the written H2 sections (fact sections aside,
// they have their own rule) are exactly repoProfileHeadings in that order.
// Each missing and each unexpected heading (one not in the shape, or a
// second copy of one) is reported with its position; when the headings are
// all there but out of order, the first one out of place is.
func checkRepoProfileSections(s *site, e *content.Explainer) {
	_, sections := e.Body.Sections()
	var got []string
	for _, sec := range sections {
		if _, fact := repoProfileFactSections[sec.Heading]; !fact {
			got = append(got, sec.Heading)
		}
	}
	want := repoProfileHeadings
	if slices.Equal(got, want) {
		return
	}
	shape := "a Repo Profile has exactly these H2 sections in this order: " + strings.Join(want, ", ")
	reported := false
	for i, h := range want {
		if !slices.Contains(got, h) {
			s.fail(e.File, ruleRepoProfileSections, "missing ## %s (%s); %s", h, sectionPosition(want, i), shape)
			reported = true
		}
	}
	for i, h := range got {
		if !slices.Contains(want, h) || slices.Contains(got[:i], h) {
			s.fail(e.File, ruleRepoProfileSections, "unexpected ## %s (section %d); %s", h, i+1, shape)
			reported = true
		}
	}
	if reported {
		return
	}
	// Same headings, same count: the order differs.
	for i := range want {
		if got[i] != want[i] {
			s.fail(e.File, ruleRepoProfileSections, "## %s is section %d, where ## %s belongs; %s", got[i], i+1, want[i], shape)
			return
		}
	}
}

// sectionPosition describes where section i of want sits, for a message:
// "section 1" or "section 3, after ## In 60 seconds".
func sectionPosition(want []string, i int) string {
	pos := fmt.Sprintf("section %d", i+1)
	if i > 0 {
		pos += ", after ## " + want[i-1]
	}
	return pos
}

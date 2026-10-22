package builder

import (
	"bytes"
	"context"

	"github.com/a-h/templ"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/views"
)

// A Repo Profile's fact sections (spec "Page composition and gates"): the
// Author writes the opinion sections; the builder inserts Signals after
// Architecture and Changelog after Alternatives, rendered from the Fact
// Pack. They go into the body as builder HTML (like a rendered diagram),
// so the page, the full-text feeds and the search index all carry them.
// Runs after the Sources merge and before the heroes encode.
func init() {
	registerStep(phaseTransform, 35, "insert repo profile fact sections", insertFactSections)
}

// ruleFactSections: a published Repo Profile lacks the section a fact
// section is inserted after (the Repo Profile skeleton rule reports the
// missing heading itself; this is the safety net behind it).
const ruleFactSections = "fact-sections"

func insertFactSections(s *site) error {
	for _, e := range s.publishedExplainers() {
		if e.Kind != content.RepoProfile {
			continue
		}
		// Every published Repo Profile has a Fact Pack (loadFactPacks fails
		// the Profile otherwise).
		fp := s.factPacks[e.File]
		for heading, section := range map[string]templ.Component{
			headingSignals:   views.RepoProfileSignalsSection(repoProfileSignals(fp)),
			headingChangelog: views.RepoProfileChangelogSection(repoProfileChangelog(fp)),
		} {
			html, err := renderHTML(section)
			if err != nil {
				return err
			}
			// repoProfileFactSections (the skeleton rule) names the written
			// section each fact section follows.
			after := repoProfileFactSections[heading]
			if !e.Body.InsertHTMLAfterSection(after, html) {
				s.fail(e.File, ruleFactSections, "no ## %s section to insert %s after", after, heading)
			}
		}
	}
	return nil
}

// renderHTML renders a component to a string for insertion into a body.
func renderHTML(c templ.Component) (string, error) {
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

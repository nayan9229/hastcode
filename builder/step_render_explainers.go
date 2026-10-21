package builder

import (
	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/markdown"
	"hastcode.com/site/internal/views"
)

func init() {
	registerStep(phaseRender, 20, "render explainers", renderExplainers)
}

func renderExplainers(s *site) error {
	for _, e := range s.publishedExplainers() {
		lead, rest := e.Body.SplitAfterSection(headingTLDR)
		leadHTML, err := markdown.RenderBlocks(e.Body, lead)
		if err != nil {
			s.fail(e.File, "markdown", "%v", err)
			continue
		}
		body, err := markdown.RenderBlocks(e.Body, rest)
		if err != nil {
			s.fail(e.File, "markdown", "%v", err)
			continue
		}
		if err := s.writePage(explainerURL(e), views.ExplainerArticle(explainerPage(s, e, leadHTML, body))); err != nil {
			return err
		}
	}
	return nil
}

// explainerPage maps an Explainer to its view model. lead is the intro and
// TL;DR (the first ad slot follows it), body the rest.
func explainerPage(s *site, e *content.Explainer, lead, body string) views.ExplainerPage {
	return views.ExplainerPage{
		Page:         s.page(e.Title+" | Hast Code", e.Description),
		Title:        e.Title,
		Description:  e.Description,
		VerticalName: s.verticals[e.Vertical].Name,
		AuthorName:   s.authors[e.Author].Name,
		Published:    e.Published,
		Reviewed:     e.Reviewed,
		LeadHTML:     lead,
		BodyHTML:     body,
	}
}

package builder

import (
	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/markdown"
	"hastcode.com/site/internal/schema"
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
			s.fail(e.File, ruleMarkdown, "%v", err)
			continue
		}
		rest, sources := e.Body.SplitAtSection(rest, headingSources)
		body, err := markdown.RenderBlocks(e.Body, rest)
		if err != nil {
			s.fail(e.File, ruleMarkdown, "%v", err)
			continue
		}
		sourcesHTML, err := markdown.RenderBlocks(e.Body, sources)
		if err != nil {
			s.fail(e.File, ruleMarkdown, "%v", err)
			continue
		}
		page := explainerPage(s, e, leadHTML, body)
		page.SourcesHTML = sourcesHTML
		if err := s.writePage(explainerURL(e), views.ExplainerArticle(page)); err != nil {
			return err
		}
	}
	return nil
}

// explainerPage maps an Explainer to its view model. lead is the intro and
// TL;DR (the first ad slot follows it), body the rest.
func explainerPage(s *site, e *content.Explainer, lead, body string) views.ExplainerPage {
	hero := s.heroView(e.File)
	author := s.authors[e.Author]
	vertical := s.verticals[e.Vertical]
	page := s.articlePage(articleHead{
		Type:        schema.TypeArticle,
		URL:         explainerURL(e),
		Headline:    e.Title,
		Description: e.Description,
		Hero:        hero,
		Published:   e.Published,
		Modified:    e.Reviewed, // dateModified is the last Review, never the build time
		Author:      author,
		Section:     vertical.Name,
		Keywords:    e.Tags,
		Breadcrumbs: []views.Link{
			{Name: "Home", URL: homeURL},
			{Name: vertical.Name, URL: verticalURL(vertical.Slug)},
			{Name: e.Title, URL: explainerURL(e)},
		},
	})
	ep := views.ExplainerPage{
		Page:         page,
		Title:        e.Title,
		Description:  e.Description,
		VerticalName: vertical.Name,
		Author:       views.Link{Name: author.Name, URL: authorURL(author.Slug)},
		Published:    e.Published,
		Reviewed:     e.Reviewed,
		Hero:         hero,
		LeadHTML:     lead,
		BodyHTML:     body,
		Newsletter:   views.Newsletter{Source: explainerURL(e)},
		ReadNext:     s.readingLists(e),
	}
	if e.Kind == content.RepoProfile {
		// Every published Repo Profile has a Fact Pack (loadFactPacks fails
		// the Profile otherwise).
		ep.RepoProfile = repoProfileCard(e, s.factPacks[e.File])
	}
	return ep
}

package builder

import (
	"net/url"
	"strings"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/markdown"
	"hastcode.com/site/internal/schema"
	"hastcode.com/site/internal/views"
)

func init() {
	registerStep(phaseRender, 30, "render briefs", renderBriefs)
}

func renderBriefs(s *site) error {
	for _, b := range s.publishedBriefs() {
		// The first ad slot follows "What happened" (a Brief has no TL;DR).
		lead, rest := b.Body.SplitAfterSection(headingWhatHappened)
		leadHTML, err := markdown.RenderBlocks(b.Body, lead)
		if err != nil {
			s.fail(b.File, "markdown", "%v", err)
			continue
		}
		body, err := markdown.RenderBlocks(b.Body, rest)
		if err != nil {
			s.fail(b.File, "markdown", "%v", err)
			continue
		}
		if err := s.writePage(briefURL(b), views.BriefArticle(briefPage(s, b, leadHTML, body))); err != nil {
			return err
		}
	}
	return nil
}

// briefPage maps a Brief to its view model.
func briefPage(s *site, b *content.Brief, lead, body string) views.BriefPage {
	var hero views.Hero // every published Brief has one (validateHeroes)
	if hr := s.heroes[b.File]; hr != nil {
		hero = hr.view
	}
	author := s.authors[b.Author]
	page := s.articlePage(articleHead{
		Type:        schema.TypeNewsArticle,
		URL:         briefURL(b),
		Headline:    b.Title,
		Description: b.Description,
		Hero:        hero,
		Published:   b.Published,
		Modified:    b.Published, // a Brief is never revised in place
		Author:      author,
		Keywords:    b.Tags,
		Breadcrumbs: []views.Link{
			{Name: "Home", URL: homeURL},
			{Name: "Briefs", URL: briefsIndexURL},
			{Name: b.Title, URL: briefURL(b)},
		},
	})
	return views.BriefPage{
		Page:        page,
		Title:       b.Title,
		Description: b.Description,
		Author:      views.Link{Name: author.Name, URL: authorURL(author.Slug)},
		Published:   b.Published,
		Hero:        hero,
		LeadHTML:    lead,
		BodyHTML:    body,
		Source:      b.SourceURL,
		SourceName:  sourceName(b.SourceURL),
		Newsletter:  views.Newsletter{Source: briefURL(b)},
	}
}

// sourceName is the host a source URL points at, without "www.".
func sourceName(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return raw
	}
	return strings.TrimPrefix(u.Host, "www.")
}

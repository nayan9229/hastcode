package builder

import (
	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/markdown"
	"hastcode.com/site/internal/views"
)

func init() {
	registerStep(phaseRender, 45, "render static pages", renderStaticPages)
}

func renderStaticPages(s *site) error {
	for _, p := range s.publishedPages() {
		body, err := markdown.Render(p.Body)
		if err != nil {
			s.fail(p.File, ruleMarkdown, "%v", err)
			continue
		}
		if err := s.writePage(staticPageURL(p), views.StaticPage(staticPage(s, p, body))); err != nil {
			return err
		}
	}
	return nil
}

// staticPage maps a static page to its view model: s.page plus the visible
// breadcrumb and its BreadcrumbList.
func staticPage(s *site, p *content.Page, body string) views.StaticPageView {
	url := staticPageURL(p)
	page := s.page(url, p.Title+" | "+views.PublicationName, p.Description)
	page.Breadcrumbs = []views.Link{
		{Name: "Home", URL: homeURL},
		{Name: p.Title, URL: url},
	}
	page.StructuredData = []any{s.breadcrumbs(page.Breadcrumbs)}
	return views.StaticPageView{Page: page, Title: p.Title, BodyHTML: body}
}

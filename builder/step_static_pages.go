package builder

import (
	"strings"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/markdown"
	"hastcode.com/site/internal/views"
)

// Static pages: every content/pages/<slug>.md is published at /<slug>/
// (About, Contact and Editorial Policy, which the site chrome links from
// every page; Google News transparency).
func init() {
	registerStep(phaseLoad, 25, "load static pages", loadStaticPages)
	registerStep(phaseValidate, 25, "validate static page paths", validateStaticPagePaths)
	registerStep(phaseRender, 45, "render static pages", renderStaticPages)
}

const rulePagePath = "page-path"

func loadStaticPages(s *site) error {
	return s.eachMarkdown("pages", func(rel string, data []byte) {
		p, err := content.ParsePage(rel, data)
		if err != nil {
			s.fail(rel, ruleFrontmatter, "%v", err)
			return
		}
		s.pages = append(s.pages, p)
	})
}

// reservedTopLevel are the first path segments the site (or a Vercel
// route) already uses, besides the Vertical hubs.
func reservedTopLevel() []string {
	var reserved []string
	for _, u := range []string{authorURL(""), briefsIndexURL, searchURL, subscribedURL(), heroAssetDir, "/explainers/", "/api/", "/pagefind/"} {
		reserved = append(reserved, strings.Split(strings.Trim(u, "/"), "/")[0])
	}
	// Retired community URLs (410) and /home (301 to /), ADR-0001.
	return append(reserved, "video", "community", "group", "home")
}

// validateStaticPagePaths fails a page whose /<slug>/ is a Vertical hub or
// another reserved path.
func validateStaticPagePaths(s *site) error {
	reserved := map[string]bool{}
	for _, r := range reservedTopLevel() {
		reserved[r] = true
	}
	for slug := range s.verticals {
		reserved[slug] = true
	}
	for _, p := range s.pages {
		if reserved[p.Slug] {
			s.fail(p.File, rulePagePath, "static page path %s is already used by the site; rename the file", staticPageURL(p))
		}
	}
	return nil
}

func renderStaticPages(s *site) error {
	for _, p := range s.pages {
		if s.failed[p.File] {
			continue
		}
		body, err := markdown.Render(p.Body)
		if err != nil {
			s.fail(p.File, "markdown", "%v", err)
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

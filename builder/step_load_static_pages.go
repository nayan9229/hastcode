package builder

import "hastcode.com/site/internal/content"

// Static pages: every content/pages/<slug>.md is published at /<slug>/
// (About, Contact and Editorial Policy, which the site chrome links from
// every page; Google News transparency). Paths are checked in
// step_validate_static_page_paths.go, pages written in
// step_render_static_pages.go.
func init() {
	registerStep(phaseLoad, 25, "load static pages", loadStaticPages)
}

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

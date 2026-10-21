package builder

import "hastcode.com/site/internal/views"

func init() {
	// Before the link gate (1000), which then checks the header's Search
	// link against this file.
	registerStep(phaseRender, 62, "render search page", renderSearch)
}

// renderSearch writes /search/: noindex (a results page has nothing of its
// own to rank), and never indexed by Pagefind (no data-pagefind-body).
// The Pagefind bundle it loads is written by the index step (write phase).
func renderSearch(s *site) error {
	p := s.page(searchURL, "Search | "+views.PublicationName, "Search every Explainer and Brief on "+views.PublicationName+".")
	p.NoIndex = true
	p.Breadcrumbs = []views.Link{{Name: "Home", URL: homeURL}, {Name: "Search", URL: searchURL}}
	p.StructuredData = append(p.StructuredData, s.breadcrumbs(p.Breadcrumbs))
	return s.writePage(searchURL, views.Search(views.SearchPage{Page: p, Bundle: pagefindBundleURL}))
}

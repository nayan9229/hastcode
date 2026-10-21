package builder

import "hastcode.com/site/internal/views"

func init() {
	registerStep(phaseRender, 65, "render 404 page", renderNotFound)
}

// renderNotFound writes the static 404 page. It is noindex and has no
// canonical: it is served at whatever unknown path was requested.
// Ticket 10 points Vercel's fallback route at notFoundURL.
func renderNotFound(s *site) error {
	p := s.page(notFoundURL, "Page not found | "+views.PublicationName, "This page does not exist on "+views.PublicationName+".")
	p.Canonical = ""
	p.NoIndex = true
	return s.writePage(notFoundURL, views.NotFound(p))
}

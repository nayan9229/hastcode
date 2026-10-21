package builder

import "hastcode.com/site/internal/views"

func init() {
	registerStep(phaseRender, 60, "render subscribed page", renderSubscribed)
}

// renderSubscribed writes the noindex landing page for no-JS newsletter
// signups (the subscribe function redirects there).
func renderSubscribed(s *site) error {
	return s.writePage(subscribedURL(), views.SubscribedPage(views.Page{
		Title:       "Check your inbox | Hast Code",
		Description: "Confirm your Hast Code newsletter subscription.",
		NoIndex:     true,
	}))
}

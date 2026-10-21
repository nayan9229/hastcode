package builder

import (
	"hastcode.com/site/internal/views"
	"hastcode.com/site/subscribe"
)

func init() {
	registerStep(phaseRender, 60, "render subscribed page", renderSubscribed)
}

// renderSubscribed writes the noindex landing page for no-JS newsletter
// signups (the subscribe function redirects there).
func renderSubscribed(s *site) error {
	p := s.page(subscribe.SubscribedPath, "Check your inbox | "+views.PublicationName,
		"Confirm your "+views.PublicationName+" newsletter subscription.")
	p.NoIndex = true
	return s.writePage(subscribe.SubscribedPath, views.SubscribedPage(p))
}

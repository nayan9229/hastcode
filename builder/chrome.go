package builder

import (
	"cmp"
	"slices"

	"hastcode.com/site/internal/views"
)

// page returns the view model every page passes to views.Layout: head
// fields (canonical from urlPath, default Open Graph) plus the site chrome
// (header nav, footer links), ads and analytics. Every render step builds
// its views.Page through this (or s.articlePage, which calls it) so the
// canonical and chrome are never missing.
func (s *site) page(urlPath, title, description string) views.Page {
	return views.Page{
		Title:           title,
		Description:     description,
		Canonical:       s.absURL(urlPath),
		Social:          views.Social{Type: "website", Title: title},
		Chrome:          s.chrome(),
		Ads:             views.Ads{Enabled: s.cfg.AdsEnabled, Client: s.cfg.AdSenseClient},
		PlausibleDomain: s.cfg.PlausibleDomain,
	}
}

// sortedVerticals are the configured Verticals in their configured order.
func (s *site) sortedVerticals() []Vertical {
	verticals := slices.Clone(s.cfg.Verticals)
	slices.SortStableFunc(verticals, func(a, b Vertical) int { return cmp.Compare(a.Order, b.Order) })
	return verticals
}

func (s *site) chrome() views.Chrome {
	c := views.Chrome{
		Briefs: views.Link{Name: "Briefs", URL: briefsIndexURL},
		Search: views.Link{Name: "Search", URL: searchURL},
		Pages: []views.Link{
			{Name: "About", URL: aboutURL},
			{Name: "Contact", URL: contactURL},
			{Name: "Editorial Policy", URL: editorialPolicyURL},
		},
		Feeds: []views.Link{
			{Name: "RSS (all)", URL: siteFeedURL},
			{Name: "RSS (Briefs)", URL: briefsFeedURL},
		},
	}
	for _, v := range s.sortedVerticals() {
		c.Verticals = append(c.Verticals, views.Link{Name: v.Name, URL: verticalURL(v.Slug)})
	}
	return c
}

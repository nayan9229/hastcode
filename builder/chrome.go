package builder

import (
	"cmp"
	"slices"

	"hastcode.com/site/internal/views"
)

// page returns the view model every page passes to views.Layout: head
// fields plus the site chrome (header nav, footer links), ads and analytics. Every render step
// builds its views.Page through this so the chrome is never missing.
func (s *site) page(title, description string) views.Page {
	return views.Page{
		Title:           title,
		Description:     description,
		Chrome:          s.chrome(),
		Ads:             views.Ads{Enabled: s.cfg.AdsEnabled, Client: s.cfg.AdSenseClient},
		PlausibleDomain: s.cfg.PlausibleDomain,
	}
}

func (s *site) chrome() views.Chrome {
	verticals := slices.Clone(s.cfg.Verticals)
	slices.SortStableFunc(verticals, func(a, b Vertical) int { return cmp.Compare(a.Order, b.Order) })
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
	for _, v := range verticals {
		c.Verticals = append(c.Verticals, views.Link{Name: v.Name, URL: verticalURL(v.Slug)})
	}
	return c
}

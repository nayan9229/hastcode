package builder

import (
	"cmp"
	"net/url"
	"slices"

	"hastcode.com/site/internal/views"
)

// page returns the view model every page passes to views.Layout: head
// fields (canonical from urlPath, default Open Graph) plus the site chrome
// (header nav, footer links), ads and analytics. Every render step builds
// its views.Page through this (or s.articlePage, which calls it) so the
// canonical and chrome are never missing.
func (s *site) page(urlPath, title, description string) views.Page {
	badge := views.OGBadge()
	badge.URL = s.absURL(s.imageURL(badge.URL))
	return views.Page{
		Title:       title,
		Description: description,
		Canonical:   s.absURL(urlPath),
		Social:      views.Social{Type: "website", Title: title, Image: &badge, ImageAlt: views.OGBadgeAlt},
		Chrome:      s.chrome(),
		Ads:         views.Ads{Enabled: s.adsOn(), Client: s.cfg.AdSenseClient, Slots: s.cfg.AdSenseSlots},
		Analytics:   s.analytics(),
	}
}

// adsOn reports whether pages carry ad units: never in preview builds
// (Previews and make dev), whose impressions would be the Author's own and
// which Lighthouse CI measures without third-party scripts (ADR-0006).
func (s *site) adsOn() bool {
	return s.cfg.AdsEnabled && s.cfg.Mode != Preview
}

// analytics is the Google Analytics 4 configuration pages carry (ADR-0008):
// none in preview builds (Previews and make dev), whose visits are not
// Readers' and which Lighthouse CI measures without third-party scripts
// (ADR-0006). Host is site_url's host, the only one that reports.
func (s *site) analytics() views.Analytics {
	if s.cfg.Mode == Preview || s.cfg.GAMeasurementID == "" {
		return views.Analytics{}
	}
	host := ""
	if u, err := url.Parse(s.cfg.SiteURL); err == nil {
		host = u.Hostname()
	}
	return views.Analytics{MeasurementID: s.cfg.GAMeasurementID, Host: host}
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
			{Name: "Privacy", URL: privacyURL},
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

package builder

import "hastcode.com/site/internal/content"

// Site URL scheme (spec "URLs"). Paths are site-relative with a trailing
// slash; prefix s.cfg.SiteURL for absolute URLs. Add new page types here.

func explainerURL(e *content.Explainer) string {
	return "/" + e.Vertical + "/" + e.Slug + "/"
}

func verticalURL(slug string) string { return "/" + slug + "/" }

// Fixed pages and feeds. Some are written by later tickets; the site
// chrome links them from every page already.
const (
	briefsIndexURL     = "/briefs/"
	searchURL          = "/search/"
	aboutURL           = "/about/"
	contactURL         = "/contact/"
	editorialPolicyURL = "/editorial-policy/"
	siteFeedURL        = "/feed.xml"
	briefsFeedURL      = "/briefs/feed.xml"
)

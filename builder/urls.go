package builder

import (
	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/subscribe"
)

// Site URL scheme (spec "URLs"). Paths are site-relative with a trailing
// slash; prefix s.cfg.SiteURL for absolute URLs. Add new page types here.

func explainerURL(e *content.Explainer) string {
	return "/" + e.Vertical + "/" + e.Slug + "/"
}

func briefURL(b *content.Brief) string { return briefsIndexURL + b.Slug + "/" }

// subscribedURL is where the subscribe function sends no-JS signups.
func subscribedURL() string { return subscribe.SubscribedPath }

func verticalURL(slug string) string { return "/" + slug + "/" }

// authorURL is the Author page (written by ticket 08); Person JSON-LD
// uses it as the author's url.
func authorURL(slug string) string { return "/authors/" + slug + "/" }

// staticPageURL is a static page from content/pages/<slug>.md.
func staticPageURL(p *content.Page) string { return "/" + p.Slug + "/" }

// notFoundURL is the static 404 page; Vercel's fallback route (ticket 10)
// serves it, with status 404, for every unknown path.
const notFoundURL = "/404.html"

// Fixed pages and feeds. Some are written by later tickets; the site
// chrome links them from every page already.
const (
	homeURL            = "/"
	briefsIndexURL     = "/briefs/"
	searchURL          = "/search/"
	aboutURL           = "/about/"
	contactURL         = "/contact/"
	editorialPolicyURL = "/editorial-policy/"
	siteFeedURL        = "/feed.xml"
	briefsFeedURL      = "/briefs/feed.xml"
)

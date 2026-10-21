package builder

import "hastcode.com/site/internal/content"

// Site URL scheme (spec "URLs"). Paths are site-relative with a trailing
// slash; prefix s.cfg.SiteURL for absolute URLs. Add new page types here.

func explainerURL(e *content.Explainer) string {
	return "/" + e.Vertical + "/" + e.Slug + "/"
}

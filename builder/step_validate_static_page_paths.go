package builder

import (
	"strings"

	"hastcode.com/site/subscribe"
)

func init() {
	registerStep(phaseValidate, 25, "validate static page paths", validateStaticPagePaths)
}

const rulePagePath = "page-path"

// reservedTopLevel are the first path segments the site (or a Vercel
// route) already uses, besides the Vertical hubs.
func reservedTopLevel() []string {
	var reserved []string
	for _, u := range []string{authorURL(""), briefsIndexURL, searchURL, subscribe.SubscribedPath, assetsDir, kindPagesURL, subscribe.Path, pagefindBundleURL, homeAliasURL} {
		reserved = append(reserved, strings.Split(strings.Trim(u, "/"), "/")[0])
	}
	return append(reserved, retiredSections...)
}

// validateStaticPagePaths fails a page whose /<slug>/ is a Vertical hub or
// another reserved path.
func validateStaticPagePaths(s *site) error {
	reserved := map[string]bool{}
	for _, r := range reservedTopLevel() {
		reserved[r] = true
	}
	for slug := range s.verticals {
		reserved[slug] = true
	}
	for _, p := range s.pages {
		if reserved[p.Slug] {
			s.fail(p.File, rulePagePath, "static page path %s is already used by the site; rename the file", staticPageURL(p))
		}
	}
	return nil
}

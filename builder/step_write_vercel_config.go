package builder

import (
	"strings"

	"hastcode.com/site/internal/vercel"
)

func init() {
	registerStep(phaseWrite, 100, "write config.json", writeVercelConfig)
}

// retiredSections are the Genuin community's top-level paths, gone for good
// (ADR-0001): every URL under them answers 410.
var retiredSections = []string{"video", "community", "group"}

const (
	// assetsDir holds only content-hashed files (heroAssetDir,
	// authorPhotoDir, the stylesheet): safe to cache forever.
	assetsDir = "/assets/"

	htmlCacheControl      = "public, max-age=0, must-revalidate"
	htmlCDNCacheControl   = "s-maxage=31536000, stale-while-revalidate=86400"
	immutableCacheControl = "public, max-age=31536000, immutable"
)

// writeVercelConfig writes the Build Output API v3 config.json.
func writeVercelConfig(s *site) error {
	return s.out.WriteConfig(vercel.Config{Version: 3, Routes: routes()})
}

// routes is the deployment's routing table (spec "URLs", "Deployment").
// Vercel's route src is PCRE; keep every pattern RE2-compatible too (no
// lookaround) so the local router (vercel.NewRouter) applies the same file.
func routes() []vercel.Route {
	return []vercel.Route{
		// Cache headers first ("continue": they add headers, then routing
		// goes on). Static output only changes with a deployment, and
		// Vercel scopes its CDN cache to the deployment, so the CDN may
		// hold HTML for long while browsers revalidate every time.
		{Src: "^/(?:.*/|.*\\.html)?$", Continue: true, Headers: map[string]string{
			"Cache-Control":            htmlCacheControl,
			"Vercel-CDN-Cache-Control": htmlCDNCacheControl,
		}},
		// Everything under /assets/ is named by a hash of its bytes
		// (stylesheet, hero variants, author photos). Pagefind's files keep
		// stable names, so /pagefind/ keeps Vercel's revalidating default.
		{Src: "^" + assetsDir + ".+$", Continue: true, Headers: map[string]string{
			"Cache-Control": immutableCacheControl,
		}},

		// Retired community (ADR-0001): 410 with the 404 page as the body.
		{Src: "^/(?:" + strings.Join(retiredSections, "|") + ")(?:/.*)?$", Status: 410, Dest: notFoundURL},
		// The old app's home alias (ADR-0001).
		{Src: "^/home/?$", Status: 301, Headers: map[string]string{"Location": "/"}},
		// Functions and /.well-known keep their exact path: no slash
		// redirect may turn POST /api/subscribe into /api/subscribe/.
		{Src: "^/(?:api|\\.well-known)(?:/.*)?$"},
		// One URL per page (trailingSlash: true, as @vercel/routing-utils
		// writes it, plus index.html): pages live at "/dir/", files at
		// their own name. Everything else 308s there.
		{Src: "^/((?:[^/]+/)*)index\\.html$", Status: 308, Headers: map[string]string{"Location": "/$1"}},
		{Src: "^/((?:[^/]+/)*[^/.]+)$", Status: 308, Headers: map[string]string{"Location": "/$1/"}},
		{Src: "^/((?:[^/]+/)*[^/]+\\.\\w+)/$", Status: 308, Headers: map[string]string{"Location": "/$1"}},

		// Static files and functions (/api/subscribe) are served here.
		{Handle: vercel.HandleFilesystem},
		// Nothing matched: the static 404 page, status 404.
		{Src: "^/.*$", Status: 404, Dest: notFoundURL},
	}
}

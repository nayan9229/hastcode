package builder

import (
	"fmt"
	"path"
	"strings"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/images"
	"hastcode.com/site/internal/pagefind"
)

// Site URL scheme (spec "URLs"). Paths are site-relative with a trailing
// slash; prefix s.cfg.SiteURL for absolute URLs. Add new page types here.
// The subscribe function's paths (subscribe.Path, subscribe.SubscribedPath)
// are its form contract and live in package subscribe.

func explainerURL(e *content.Explainer) string {
	return "/" + e.Vertical + "/" + e.Slug + "/"
}

func briefURL(b *content.Brief) string { return briefsIndexURL + b.Slug + "/" }

func verticalURL(slug string) string { return "/" + slug + "/" }

// verticalFeedURL is a Vertical's full-text RSS feed.
func verticalFeedURL(slug string) string { return verticalURL(slug) + "feed.xml" }

// kindURL is the kind page listing every Explainer of kind k.
func kindURL(k content.Kind) string { return kindPagesURL + string(k) + "/" }

// authorURL is the Author page; Person JSON-LD uses it as the author's url.
func authorURL(slug string) string { return "/authors/" + slug + "/" }

// staticPageURL is a static page from content/pages/<slug>.md.
func staticPageURL(p *content.Page) string { return "/" + p.Slug + "/" }

// notFoundURL is the static 404 page; Vercel's fallback route serves it,
// with status 404, for every unknown path.
const notFoundURL = "/404.html"

// Fixed pages, feeds and sitemaps. The site chrome links the pages and
// feeds from every page.
const (
	homeURL            = "/"
	briefsIndexURL     = "/briefs/"
	kindPagesURL       = "/explainers/" // kind pages live under it; no page of its own
	searchURL          = "/search/"
	aboutURL           = "/about/"
	contactURL         = "/contact/"
	editorialPolicyURL = "/editorial-policy/"
	privacyURL         = "/privacy/"
	siteFeedURL        = "/feed.xml"
	briefsFeedURL      = "/briefs/feed.xml"
	sitemapURL         = "/sitemap.xml"
	newsSitemapURL     = "/news-sitemap.xml"
)

// Retired URLs (ADR-0001): every path under a retired section (the Genuin
// community's top-level paths) answers 410, and homeAliasURL 301s to /.
// No page may take these paths.
var retiredSections = []string{"video", "community", "group"}

const homeAliasURL = "/home"

// encodedImageURL is the site path of an encoded image variant:
// <dir><source name>-<w>x<h>.<hash>.<ext>, named by a hash of its bytes so
// it can be cached immutably (heroes, Author photos). Pages use
// s.imageURL(path), which maps it onto the Image Store in r2 mode.
func encodedImageURL(dir, source string, f images.File) string {
	base := strings.TrimSuffix(path.Base(source), path.Ext(source))
	return fmt.Sprintf("%s%s-%dx%d.%s%s", dir, base, f.Width, f.Height, f.Hash[:10], f.Format.Ext())
}

// Asset directories. assetsDir holds only content-hashed files (the
// stylesheet, the OG badge, heroAssetDir, authorPhotoDir), so the routes
// cache it forever. pagefindBundleURL is where the browser loads Pagefind
// from (the search page's lazy loader reads it, the index step writes it);
// its files keep stable names and are revalidated.
const (
	assetsDir         = "/assets/"
	heroAssetDir      = assetsDir + "heroes/"
	authorPhotoDir    = assetsDir + "authors/"
	pagefindBundleURL = "/" + pagefind.BundleDir + "/"
)

// adsTxtURL is where ad systems look for the Authorized Digital Sellers
// list: the root of the domain.
const adsTxtURL = "/ads.txt"

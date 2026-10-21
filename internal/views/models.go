// Package views holds the templ components that render every page, plus the
// plain view-model structs they take. Components never see content.* types;
// the builder maps content to these models.
//
// File naming = component family:
//
//	layout*.templ   page shell, head, site header, site footer
//	article*.templ  Explainer and Brief pages
//	listing*.templ  home, Vertical hub, Briefs index, kind filter
//	page*.templ     Author, About/Contact/Editorial Policy, 404
//	partial_*.templ one file per partial (dates, hero, TL;DR, newsletter…)
//
// Generated *_templ.go files are not committed; run `make generate`.
package views

//go:generate go tool templ generate

import (
	"fmt"
	"strings"
	"time"
)

// Page is what every page passes to Layout: the head and site chrome.
// The builder makes it with s.page(title, description) so Chrome is
// always filled.
type Page struct {
	// Title is the full <title> text.
	Title       string
	Description string
	// Canonical is the page's absolute canonical URL (trailing slash).
	Canonical string
	// Social is the Open Graph / Twitter card preview.
	Social Social
	// Breadcrumbs is the trail from Home to this page (site-relative
	// URLs, last = this page), shown above the content. Empty = none.
	Breadcrumbs []Link
	// StructuredData are JSON-LD blocks (internal/schema values), each
	// written to its own script in <head>.
	StructuredData []any
	Chrome         Chrome
	// NoIndex keeps the page out of search indexes (meta robots noindex).
	NoIndex bool
	// HeroPreload, when set, adds a preload hint for the hero to the head.
	HeroPreload *Hero
	// Ads is the site-wide ad configuration every ad slot reads.
	Ads Ads
	// PlausibleDomain is the Plausible analytics site domain; "" omits
	// the analytics script (tests, local builds).
	PlausibleDomain string
}

// Social is a page's Open Graph and Twitter card preview.
type Social struct {
	// Type is the og:type: "website" (default) or "article".
	Type string
	// Title is og:title: the page's own name, without the site suffix
	// <title> carries.
	Title string
	// Image is the preview image with an absolute URL; nil = none (a
	// "summary" card). Article pages use the 1200 px 16:9 hero crop.
	Image    *Image
	ImageAlt string
	// Published and Modified are RFC 3339 article:published_time and
	// article:modified_time ("" = omitted).
	Published, Modified string
}

// Ads is the one flag (plus account) that turns every AdSlot from an
// empty placeholder into an AdSense unit.
type Ads struct {
	Enabled bool
	// Client is the AdSense publisher id ("ca-pub-..."); required when Enabled.
	Client string
}

// Link is a named site-relative URL.
type Link struct {
	Name string
	URL  string
}

// Chrome is the site-wide navigation in every page's header and footer.
// The builder fills it from config and its URL scheme (builder/urls.go).
type Chrome struct {
	// Verticals are the Vertical hubs, in configured order.
	Verticals []Link
	Briefs    Link
	Search    Link
	// Pages are the transparency pages (About, Contact, Editorial Policy).
	Pages []Link
	// Feeds are the RSS feeds linked from the footer.
	Feeds []Link
}

// ExplainerPage is the view model for an Explainer article page.
type ExplainerPage struct {
	Page         Page
	Title        string
	Description  string
	VerticalName string
	// Author is the Author's name and Author page (the byline link).
	Author    Link
	Published time.Time
	Reviewed  time.Time
	Hero      Hero
	// LeadHTML is the rendered intro and TL;DR, before the first ad slot;
	// BodyHTML is the rest of the Markdown body (both trusted: built from
	// content).
	LeadHTML string
	BodyHTML string
	// SourcesHTML is the rendered "## Sources" section that ends the body
	// (kept out of the search index); "" when there is none.
	SourcesHTML string
	// Newsletter is the subscribe form after the body.
	Newsletter Newsletter
	// ReadNext are the "More in {Vertical}" and "Related" lists that end
	// the article.
	ReadNext RelatedLinks
}

// BriefPage is the view model for a Brief article page.
type BriefPage struct {
	Page        Page
	Title       string
	Description string
	// Author is the byline: the Author's name and Author page.
	Author    Link
	Published time.Time
	Hero      Hero
	// LeadHTML is the rendered body through "What happened", before the
	// first ad slot; BodyHTML is the rest (both trusted: built from content).
	LeadHTML string
	BodyHTML string
	// Sources are the public sources the Brief is curated from (absolute
	// URLs): source_url first, then the frontmatter `sources`.
	Sources []Link
	// Newsletter is the subscribe form that ends the Brief.
	Newsletter Newsletter
}

// BriefsIndexPage is one page of the Briefs index (newest first).
type BriefsIndexPage struct {
	Page       Page
	Heading    string
	Items      []ListingItem
	Pagination Pagination
}

// Pagination is a listing page's place in its sequence: Prev and Next are
// site-relative URLs ("" at either end). Rendered by PaginationNav.
type Pagination struct {
	Current, Total int
	Prev, Next     string
}

// Newsletter is the view model for the newsletter form partial.
type Newsletter struct {
	// Source is the site-relative path of the page the form is on; it is
	// posted as `source` and becomes the subscriber's Buttondown tag.
	Source string
}

// DiagramFigure is a rendered Mermaid diagram placed in an article body.
type DiagramFigure struct {
	// SVG is the inline SVG (trusted: built by internal/mermaid).
	SVG string
	// Caption is from the fence's "%% title:" line; "" means no figcaption.
	Caption string
}

// Image is one image file in the output: its site path and pixel size.
type Image struct {
	URL    string // site-relative, e.g. "/assets/heroes/slug-1200x675.3f9a2c1b7d.jpg"
	Width  int
	Height int
}

// Hero is an article's hero image: responsive variants for the <picture>
// and fixed crops for structured data and Open Graph.
type Hero struct {
	Alt string
	// Sizes is the <img sizes> value for the responsive variants.
	Sizes string
	// AVIF, WebP and JPEG hold the same widths in ascending order, at the
	// source's aspect ratio.
	AVIF, WebP, JPEG []Image
	// Img is the <img> fallback (a JPEG); its width and height reserve
	// layout space.
	Img Image
	// Crops, 1200 px wide JPEGs, for JSON-LD image and og:image.
	Crop16x9, Crop4x3, Crop1x1 Image
}

// Srcset formats images as a srcset attribute with width descriptors.
func Srcset(imgs []Image) string {
	parts := make([]string, len(imgs))
	for i, img := range imgs {
		parts[i] = fmt.Sprintf("%s %dw", img.URL, img.Width)
	}
	return strings.Join(parts, ", ")
}

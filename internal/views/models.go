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
	// Analytics is the Google Analytics 4 configuration; a zero value
	// omits the analytics scripts (preview builds).
	Analytics Analytics
}

// Analytics configures Google Analytics 4 (ADR-0008).
type Analytics struct {
	// MeasurementID is the GA4 web data stream id ("G-..."); "" omits
	// analytics entirely.
	MeasurementID string
	// Host is the only hostname that reports (the Production Target's);
	// on any other host (the QA Target) the snippet opts out.
	Host string
}

// Social is a page's Open Graph and Twitter card preview.
type Social struct {
	// Type is the og:type: "website" (default) or "article".
	Type string
	// Title is og:title: the page's own name, without the site suffix
	// <title> carries.
	Title string
	// Image is the preview image with an absolute URL; nil = none (a
	// "summary" card). Article pages use the 1200 px 16:9 hero crop,
	// every other page the generated OG badge (OGBadge).
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
	// Slots are the AdSense ad unit ids (data-ad-slot) by slot position.
	Slots map[string]string
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
	// RepoProfile is the header card of a Repo Profile, rendered from its
	// Fact Pack; nil for every other kind.
	RepoProfile *RepoProfileCard
}

// RepoProfileCard is a Repo Profile's header card
// (partial_repo_profile_card.templ): what the repository is, who it is
// for, and the facts a Reader judges it by in five seconds.
type RepoProfileCard struct {
	// Repo is the repository's full name (owner/name) linking its GitHub
	// page.
	Repo Link
	// WhatItIs is the Profile's description (it is the h-entry p-summary);
	// WhoFor is its `who_for` line ("" omits the row).
	WhatItIs string
	WhoFor   string
	Stars    int
	// License is the SPDX id, or what the Reader should know instead
	// ("not found", "not classified").
	License string
	// Language is the primary language ("not found" when GitHub reports
	// none).
	Language string
	// LatestRelease is the release tag linking its release page (zero =
	// no release yet) and LatestReleaseDate its date.
	LatestRelease     Link
	LatestReleaseDate time.Time
	// LastCommit is the Fact Pack's pushed date.
	LastCommit time.Time
}

// FactNotFound is what a Repo Profile's fact sections show for a fact the
// Fact Pack does not carry; a fact is never invented.
const FactNotFound = "not found"

// RepoProfileSignals is a Repo Profile's Signals section
// (partial_repo_profile_signals.templ): plain facts from the Fact Pack,
// no score and no grade. Zero dates and empty strings render as
// FactNotFound; the builder decides what counts as absent.
type RepoProfileSignals struct {
	// LastCommit is the Fact Pack's pushed date (zero = not found).
	LastCommit time.Time
	// ReleasesLastYear is the number of listed releases in the twelve
	// months before the fetch, already formatted ("3", "10 or more", or
	// FactNotFound): the Fact Pack lists at most ten releases.
	ReleasesLastYear string
	// LatestRelease is the release tag linking its release page (zero =
	// no release yet) and LatestReleaseDate its date.
	LatestRelease     Link
	LatestReleaseDate time.Time
	// Contributors is the contributor count (0 = not found: a repository
	// has at least one contributor).
	Contributors int
	// TopContributor holds TopShare percent of the ListedContributions
	// made by the Listed contributors the Fact Pack names (the top ten);
	// Listed 0 = not found. The Fact Pack does not carry every
	// contribution, so the share is of the listed ones.
	TopContributor      Link
	TopShare            int
	ListedContributions int
	Listed              int
	OpenIssues          int
	// SecurityPolicy is the security policy's path ("" = not found).
	SecurityPolicy string
	// License is the SPDX id, or what the Reader should know instead
	// ("not found", "not classified").
	License  string
	Stars    int
	Forks    int
	Watchers int
	// Documentation lists the documentation files and directories found
	// (README, CONTRIBUTING, CODE_OF_CONDUCT, docs/, examples/); empty =
	// none found.
	Documentation []string
}

// RepoProfileChangelog is a Repo Profile's Changelog section
// (partial_repo_profile_changelog.templ): the Fact Pack's releases,
// newest first. nil = the Fact Pack lists no releases (not found); empty
// = the repository has none yet.
type RepoProfileChangelog struct {
	Releases []ReleaseItem
}

// ReleaseItem is one release in a Changelog: its tag linking the release
// page, its name, date and whether it is a pre-release.
type ReleaseItem struct {
	Tag        Link
	Name       string
	Date       time.Time
	Prerelease bool
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
	// AVIF and JPEG hold the same widths in ascending order, at the
	// source's aspect ratio.
	AVIF, JPEG []Image
	// Img is the <img> fallback (a JPEG); its width and height reserve
	// layout space.
	Img Image
	// Crops, 1200 px wide JPEGs: 16:9, 4:3 and 1:1 for JSON-LD image,
	// 1200x630 for og:image and twitter:image.
	Crop16x9, Crop4x3, Crop1x1, CropOG Image
}

// Srcset formats images as a srcset attribute with width descriptors.
func Srcset(imgs []Image) string {
	parts := make([]string, len(imgs))
	for i, img := range imgs {
		parts[i] = fmt.Sprintf("%s %dw", img.URL, img.Width)
	}
	return strings.Join(parts, ", ")
}

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
	Chrome      Chrome
	// NoIndex keeps the page out of search indexes (meta robots noindex).
	NoIndex bool
	// HeroPreload, when set, adds a preload hint for the hero to the head.
	HeroPreload *Hero
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
	AuthorName   string
	Published    time.Time
	Reviewed     time.Time
	Hero         Hero
	// BodyHTML is the rendered Markdown body (trusted: built from content).
	BodyHTML string
	// Newsletter is the subscribe form that ends the article.
	Newsletter Newsletter
}

// Newsletter is the view model for the newsletter form partial.
type Newsletter struct {
	// Source is the site-relative path of the page the form is on; it is
	// posted as `source` and becomes the subscriber's Buttondown tag.
	Source string
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

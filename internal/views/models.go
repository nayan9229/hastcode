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

import "time"

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

// DiagramFigure is a rendered Mermaid diagram placed in an article body.
type DiagramFigure struct {
	// SVG is the inline SVG (trusted: built by internal/mermaid).
	SVG string
	// Caption is from the fence's "%% title:" line; "" means no figcaption.
	Caption string
}

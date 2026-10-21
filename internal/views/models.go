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
type Page struct {
	// Title is the full <title> text.
	Title       string
	Description string
	// NoIndex keeps the page out of search indexes (meta robots noindex).
	NoIndex bool
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

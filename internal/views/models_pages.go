package views

// View models for the page family (page*.templ): Author, static pages
// (About, Contact, Editorial Policy) and the 404 page.

import "time"

// AuthorPage is the view model for an Author page (/authors/{slug}/).
type AuthorPage struct {
	Page Page
	Name string
	// Photo is the square Author photo; nil = none.
	Photo    *Image
	PhotoAlt string
	// BioHTML is the Author file's rendered Markdown body (trusted: built
	// from content).
	BioHTML string
	// Profiles are the Author's profiles elsewhere (the same_as list):
	// absolute URLs with a short label, linked rel="me".
	Profiles []Link
	// Latest are the Author's latest pieces, newest first.
	Latest []AuthorPiece
}

// AuthorPiece is one piece in an Author page's latest list.
type AuthorPiece struct {
	// Type is the piece type as Readers see it: "Explainer" or "Brief".
	Type      string
	Title     string
	URL       string
	Published time.Time
}

// StaticPageView is the view model for a static page (/about/, /contact/,
// /editorial-policy/).
type StaticPageView struct {
	Page  Page
	Title string
	// BodyHTML is the page's rendered Markdown (trusted: built from content).
	BodyHTML string
}

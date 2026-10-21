package views

// View models for the listing family (listing*.templ): home, Vertical
// hubs, kind pages, the Briefs index, and the "More in" / "Related"
// lists that end an Explainer.

import "time"

// ListingItem is one piece (Explainer or Brief) in a listing.
type ListingItem struct {
	Title       string
	URL         string
	Description string
	// Kind is the Explainer kind badge ("Case Study", "Concept",
	// "Comparison"); "" for a Brief.
	Kind      string
	Published time.Time
	// Thumb is the hero's 480 px variant.
	Thumb Thumbnail
}

// Thumbnail is a hero's smallest (480 px) variant in each format, for
// listing items. It is decorative: the item's title says what it is.
type Thumbnail struct {
	AVIF, WebP, JPEG Image
}

// ListingPage is one page of a paginated listing of Explainers: a
// Vertical hub or a kind page.
type ListingPage struct {
	Page    Page
	Heading string
	// Intro is shown under the heading (the Vertical or kind description).
	Intro string
	Items []ListingItem
	// Empty is shown instead of the list when there are no items.
	Empty      string
	Pagination Pagination
}

// HomePage is the home page: latest Explainers and Briefs, the Verticals
// and the kind pages.
type HomePage struct {
	Page       Page
	Explainers []ListingItem
	Briefs     []ListingItem
	Verticals  []VerticalCard
	// Kinds are the kind pages (Case Studies, Concepts, Comparisons).
	Kinds []Link
}

// VerticalCard is a Vertical on the home page: its hub and description.
type VerticalCard struct {
	Name        string
	URL         string
	Description string
}

// RelatedLinks are the reading lists that end an Explainer.
type RelatedLinks struct {
	// MoreHeading is "More in {Vertical}"; More are the latest other
	// Explainers in the Vertical.
	MoreHeading string
	More        []ListingItem
	// Related are Explainers named in `related` or sharing tags.
	Related []ListingItem
}

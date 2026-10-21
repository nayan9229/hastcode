// Package schema holds the schema.org JSON-LD node types the site emits
// (spec "Structured data and head"): Article/NewsArticle, Person,
// Organization, BreadcrumbList (and later WebSite). They are plain data:
// the builder fills them with absolute URLs and formatted dates, and the
// views marshal each top-level node into its own
// <script type="application/ld+json"> block.
//
// Deliberately absent: FAQPage, HowTo, Speakable and TechArticle (spec:
// removed, deprecated, no evidence, unlisted by Google).
package schema

// Context is the @context of every top-level block.
const Context = "https://schema.org"

// ArticleType is one of the Article subtypes Google's Article structured
// data supports.
type ArticleType string

const (
	TypeArticle     ArticleType = "Article"     // Explainers
	TypeNewsArticle ArticleType = "NewsArticle" // Briefs
)

// Article is an Article or NewsArticle. Dates are RFC 3339 with the
// offset the page shows (match the visible dates exactly); Image lists
// absolute URLs of the 16:9, 4:3 and 1:1 crops.
type Article struct {
	Context          string       `json:"@context"`
	Type             ArticleType  `json:"@type"`
	Headline         string       `json:"headline"`
	Description      string       `json:"description,omitempty"`
	Image            []string     `json:"image,omitempty"`
	DatePublished    string       `json:"datePublished"`
	DateModified     string       `json:"dateModified"`
	Author           Person       `json:"author"`
	Publisher        Organization `json:"publisher"`
	MainEntityOfPage WebPageRef   `json:"mainEntityOfPage"`
	ArticleSection   string       `json:"articleSection,omitempty"`
	Keywords         []string     `json:"keywords,omitempty"`
}

// WebPageRef points at the page an article is the main entity of.
type WebPageRef struct {
	Type string `json:"@type"` // "WebPage"
	ID   string `json:"@id"`   // the canonical URL
}

// Page returns the mainEntityOfPage reference for a canonical URL.
func Page(canonical string) WebPageRef { return WebPageRef{Type: "WebPage", ID: canonical} }

// Person is an Author. URL is the Author page (it uniquely identifies the
// author); SameAs are the Author's profiles elsewhere. Use it nested (as
// an article's author) or top-level on the Author page with Context set.
type Person struct {
	Context string `json:"@context,omitempty"`
	Type    string `json:"@type"` // "Person"
	Name    string `json:"name"`
	URL     string `json:"url,omitempty"`
	// Image is the Author's photo (absolute URL); set on the Author
	// page's top-level Person only.
	Image  string   `json:"image,omitempty"`
	SameAs []string `json:"sameAs,omitempty"`
}

// Organization is the publisher, Hast Code. Top-level on the home page
// (with Context set), nested as an article's publisher.
type Organization struct {
	Context     string       `json:"@context,omitempty"`
	Type        string       `json:"@type"` // "Organization"
	Name        string       `json:"name"`
	URL         string       `json:"url"`
	Logo        *ImageObject `json:"logo,omitempty"`
	Description string       `json:"description,omitempty"`
	SameAs      []string     `json:"sameAs,omitempty"`
}

// ImageObject is an image by absolute URL.
type ImageObject struct {
	Type string `json:"@type"` // "ImageObject"
	URL  string `json:"url"`
}

// BreadcrumbList is a page's trail from Home to the page itself.
type BreadcrumbList struct {
	Context string     `json:"@context"`
	Type    string     `json:"@type"` // "BreadcrumbList"
	Items   []ListItem `json:"itemListElement"`
}

// ListItem is one breadcrumb: 1-based Position, Name and absolute URL.
type ListItem struct {
	Type     string `json:"@type"` // "ListItem"
	Position int    `json:"position"`
	Name     string `json:"name"`
	Item     string `json:"item"`
}

// Crumb is a breadcrumb's name and absolute URL, for NewBreadcrumbList.
type Crumb struct{ Name, URL string }

// NewBreadcrumbList numbers crumbs from 1 in the order given.
func NewBreadcrumbList(crumbs ...Crumb) BreadcrumbList {
	b := BreadcrumbList{Context: Context, Type: "BreadcrumbList"}
	for i, c := range crumbs {
		b.Items = append(b.Items, ListItem{Type: "ListItem", Position: i + 1, Name: c.Name, Item: c.URL})
	}
	return b
}

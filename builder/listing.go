package builder

import (
	"cmp"
	"slices"
	"time"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/views"
)

// Shared by every listing (home, Vertical hubs, kind pages, Briefs index,
// an Explainer's "More in" and "Related"): order and list items.

// newestFirst returns a copy of pieces, newest published first (ties by
// slug, so output is stable). key reads a piece's publication and slug.
func newestFirst[T any](pieces []T, key func(T) (time.Time, string)) []T {
	pieces = slices.Clone(pieces)
	slices.SortStableFunc(pieces, func(a, b T) int {
		ta, sa := key(a)
		tb, sb := key(b)
		return cmp.Or(tb.Compare(ta), cmp.Compare(sa, sb))
	})
	return pieces
}

// newestExplainers returns a copy of es, newest published first.
func newestExplainers(es []*content.Explainer) []*content.Explainer {
	return newestFirst(es, func(e *content.Explainer) (time.Time, string) { return e.Published, e.Slug })
}

// newestBriefs returns a copy of bs, newest published first.
func newestBriefs(bs []*content.Brief) []*content.Brief {
	return newestFirst(bs, func(b *content.Brief) (time.Time, string) { return b.Published, b.Slug })
}

// explainerItems maps Explainers, in order, to listing items.
func (s *site) explainerItems(es []*content.Explainer) []views.ListingItem {
	items := make([]views.ListingItem, len(es))
	for i, e := range es {
		items[i] = views.ListingItem{
			Title:       e.Title,
			URL:         explainerURL(e),
			Description: e.Description,
			Kind:        kindOf(e.Kind).name,
			Published:   e.Published,
			Thumb:       s.thumbnail(e.File),
		}
	}
	return items
}

// briefItems maps Briefs, in order, to listing items.
func (s *site) briefItems(bs []*content.Brief) []views.ListingItem {
	items := make([]views.ListingItem, len(bs))
	for i, b := range bs {
		items[i] = views.ListingItem{
			Title:       b.Title,
			URL:         briefURL(b),
			Description: b.Description,
			Published:   b.Published,
			Thumb:       s.thumbnail(b.File),
		}
	}
	return items
}

// thumbnail is the smallest (480 px) hero variant of a published piece.
// Every published piece has an encoded hero (validateHeroes).
func (s *site) thumbnail(file string) views.Thumbnail {
	h := s.heroView(file)
	if len(h.JPEG) == 0 {
		return views.Thumbnail{}
	}
	return views.Thumbnail{AVIF: h.AVIF[0], JPEG: h.JPEG[0]}
}

// explainerKind is an Explainer kind as Readers see it, with its kind
// page.
type explainerKind struct {
	kind content.Kind
	// name is the badge ("Case Study"); plural heads the kind page.
	name, plural string
	// description introduces the kind page (from GLOSSARY.md).
	description string
}

// explainerKinds are the kinds in the order Readers see them.
var explainerKinds = []explainerKind{
	{content.CaseStudy, "Case Study", "Case Studies", "Explainers that reconstruct how one real, named system is built, from that company's public material."},
	{content.Concept, "Concept", "Concepts", "Explainers that define and teach one named architectural idea or pattern."},
	{content.Comparison, "Comparison", "Comparisons", "Explainers that set two or more alternatives side by side to help you choose."},
	{content.RepoProfile, "Repo Profile", "Repo Profiles", "Explainers that evaluate one open-source repository: what it is, how it is built, and whether you should adopt it."},
}

func kindOf(k content.Kind) explainerKind {
	for _, ek := range explainerKinds {
		if ek.kind == k {
			return ek
		}
	}
	return explainerKind{kind: k, name: string(k)}
}

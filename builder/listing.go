package builder

import (
	"cmp"
	"slices"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/views"
)

// Shared by every listing (home, Vertical hubs, kind pages, Briefs index,
// an Explainer's "More in" and "Related"): order and list items.

// newestExplainers returns a copy of es, newest published first (ties by
// slug, so output is stable).
func newestExplainers(es []*content.Explainer) []*content.Explainer {
	es = slices.Clone(es)
	slices.SortStableFunc(es, func(a, b *content.Explainer) int {
		return cmp.Or(b.Published.Compare(a.Published), cmp.Compare(a.Slug, b.Slug))
	})
	return es
}

// newestBriefs returns a copy of bs, newest published first (ties by slug).
func newestBriefs(bs []*content.Brief) []*content.Brief {
	bs = slices.Clone(bs)
	slices.SortStableFunc(bs, func(a, b *content.Brief) int {
		return cmp.Or(b.Published.Compare(a.Published), cmp.Compare(a.Slug, b.Slug))
	})
	return bs
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
	hr := s.heroes[file]
	if hr == nil || len(hr.view.JPEG) == 0 {
		return views.Thumbnail{}
	}
	return views.Thumbnail{AVIF: hr.view.AVIF[0], WebP: hr.view.WebP[0], JPEG: hr.view.JPEG[0]}
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
}

func kindOf(k content.Kind) explainerKind {
	for _, ek := range explainerKinds {
		if ek.kind == k {
			return ek
		}
	}
	return explainerKind{kind: k, name: string(k)}
}

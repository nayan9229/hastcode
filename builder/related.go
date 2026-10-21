package builder

import (
	"cmp"
	"slices"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/views"
)

// readingListSize is how many Explainers "More in {Vertical}" and
// "Related" each show at the end of an Explainer.
const readingListSize = 3

// readingLists are the lists that end Explainer e, chosen from the
// published Explainers (never e itself):
//
//   - More in {Vertical}: the latest others in e's Vertical.
//   - Related: e's `related` slugs in frontmatter order, then others by
//     number of shared tags (most first; ties newest first).
func (s *site) readingLists(e *content.Explainer) views.RelatedLinks {
	published := newestExplainers(s.publishedExplainers())

	var more []*content.Explainer
	for _, o := range published {
		if o != e && o.Vertical == e.Vertical && len(more) < readingListSize {
			more = append(more, o)
		}
	}

	var related []*content.Explainer
	add := func(o *content.Explainer) {
		if o != e && len(related) < readingListSize && !slices.Contains(related, o) {
			related = append(related, o)
		}
	}
	for _, slug := range e.Related {
		if i := slices.IndexFunc(published, func(o *content.Explainer) bool { return o.Slug == slug }); i >= 0 {
			add(published[i])
		}
	}
	var byTags []*content.Explainer
	for _, o := range published {
		if sharedTags(e, o) > 0 {
			byTags = append(byTags, o)
		}
	}
	// Stable sort keeps newest first within the same number of shared tags.
	slices.SortStableFunc(byTags, func(a, b *content.Explainer) int {
		return cmp.Compare(sharedTags(e, b), sharedTags(e, a))
	})
	for _, o := range byTags {
		add(o)
	}

	return views.RelatedLinks{
		MoreHeading: "More in " + s.verticals[e.Vertical].Name,
		More:        s.explainerItems(more),
		Related:     s.explainerItems(related),
	}
}

// sharedTags counts the tags a and b have in common.
func sharedTags(a, b *content.Explainer) int {
	n := 0
	for _, t := range a.Tags {
		if slices.Contains(b.Tags, t) {
			n++
		}
	}
	return n
}

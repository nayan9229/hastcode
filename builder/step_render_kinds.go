package builder

import (
	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/views"
)

func init() {
	registerStep(phaseRender, 52, "render kind pages", renderKindPages)
}

// renderKindPages writes /explainers/{kind}/ (and /explainers/{kind}/page/{n}/)
// for every Explainer kind: its published Explainers across all Verticals,
// newest first, listingPageSize per page.
func renderKindPages(s *site) error {
	all := newestExplainers(s.publishedExplainers())
	for _, k := range explainerKinds {
		var es []*content.Explainer
		for _, e := range all {
			if e.Kind == k.kind {
				es = append(es, e)
			}
		}
		for _, lp := range paginate(es, listingPageSize, kindURL(k.kind)) {
			page := s.listingPage(lp, k.plural, k.description, kindURL(k.kind), "No "+k.plural+" yet.")
			if err := s.writePage(lp.URL, views.Listing(page)); err != nil {
				return err
			}
		}
	}
	return nil
}

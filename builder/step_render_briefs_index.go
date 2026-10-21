package builder

import (
	"fmt"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/views"
)

func init() {
	registerStep(phaseRender, 40, "render briefs index", renderBriefsIndex)
}

const briefsIndexDescription = "Short, dated commentary on what just happened in tech and why it matters architecturally."

// renderBriefsIndex writes /briefs/ and /briefs/page/{n}/: every published
// Brief, newest first, listingPageSize per page. /briefs/ is written even
// with no Briefs, since every page's header links it.
func renderBriefsIndex(s *site) error {
	briefs := newestBriefs(s.publishedBriefs())
	for _, lp := range paginate(briefs, listingPageSize, briefsIndexURL) {
		if err := s.writePage(lp.URL, views.BriefsIndex(briefsIndexPage(s, lp))); err != nil {
			return err
		}
	}
	return nil
}

func briefsIndexPage(s *site, lp listingPage[*content.Brief]) views.BriefsIndexPage {
	heading := "Briefs"
	trail := []views.Link{{Name: "Home", URL: homeURL}, {Name: "Briefs", URL: briefsIndexURL}}
	if n := lp.Pagination.Current; n > 1 {
		heading = fmt.Sprintf("Briefs, page %d", n)
		trail = append(trail, views.Link{Name: fmt.Sprintf("Page %d", n), URL: lp.URL})
	}
	page := s.page(lp.URL, heading+" | "+views.PublicationName, briefsIndexDescription)
	page.Breadcrumbs = trail
	page.StructuredData = append(page.StructuredData, s.breadcrumbs(trail))
	return views.BriefsIndexPage{Page: page, Heading: heading, Items: s.briefItems(lp.Items), Pagination: lp.Pagination}
}

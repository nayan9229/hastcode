package builder

import (
	"fmt"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/views"
)

func init() {
	registerStep(phaseRender, 50, "render vertical hubs", renderVerticalHubs)
}

// renderVerticalHubs writes /{vertical}/ and /{vertical}/page/{n}/ for
// every configured Vertical: its published Explainers, newest first,
// listingPageSize per page. A hub is written even with no Explainers,
// since every page's header links it.
func renderVerticalHubs(s *site) error {
	all := newestExplainers(s.publishedExplainers())
	for _, v := range s.sortedVerticals() {
		var es []*content.Explainer
		for _, e := range all {
			if e.Vertical == v.Slug {
				es = append(es, e)
			}
		}
		for _, lp := range paginate(es, listingPageSize, verticalURL(v.Slug)) {
			if err := s.writePage(lp.URL, views.Listing(verticalHubPage(s, v, lp))); err != nil {
				return err
			}
		}
	}
	return nil
}

func verticalHubPage(s *site, v Vertical, lp listingPage[*content.Explainer]) views.ListingPage {
	description := v.Description
	if description == "" {
		description = "Explainers on " + v.Name + ": how real systems are built."
	}
	return s.listingPage(lp, v.Name, description, verticalURL(v.Slug),
		"No Explainers in "+v.Name+" yet.")
}

// listingPage is the view model for one page of an Explainer listing
// headed name (a Vertical hub or kind page at base): the visible
// breadcrumb Home > name [> Page n], its BreadcrumbList, and the items.
func (s *site) listingPage(lp listingPage[*content.Explainer], name, description, base, empty string) views.ListingPage {
	heading := name
	trail := []views.Link{{Name: "Home", URL: homeURL}, {Name: name, URL: base}}
	if n := lp.Pagination.Current; n > 1 {
		heading = fmt.Sprintf("%s, page %d", name, n)
		trail = append(trail, views.Link{Name: fmt.Sprintf("Page %d", n), URL: lp.URL})
	}
	page := s.page(lp.URL, heading+" | "+views.PublicationName, description)
	page.Breadcrumbs = trail
	page.StructuredData = append(page.StructuredData, s.breadcrumbs(trail))
	return views.ListingPage{
		Page:       page,
		Heading:    heading,
		Intro:      description,
		Items:      s.explainerItems(lp.Items),
		Empty:      empty,
		Pagination: lp.Pagination,
	}
}

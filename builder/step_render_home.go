package builder

import "hastcode.com/site/internal/views"

func init() {
	registerStep(phaseRender, 54, "render home", renderHome)
}

// How many of the latest pieces the home page shows.
const (
	homeExplainers = 6
	homeBriefs     = 5
)

// renderHome writes / : the latest Explainers and Briefs, newest first,
// every Vertical (configured order) with its description, and the kind
// pages, with WebSite and Organization JSON-LD.
func renderHome(s *site) error {
	explainers := newestExplainers(s.publishedExplainers())
	briefs := newestBriefs(s.publishedBriefs())
	page := s.page(homeURL, views.PublicationName+": "+views.Tagline, homeDescription)
	page.StructuredData = append(page.StructuredData, s.webSite(), s.homeOrganization())
	home := views.HomePage{
		Page:       page,
		Explainers: s.explainerItems(explainers[:min(len(explainers), homeExplainers)]),
		Briefs:     s.briefItems(briefs[:min(len(briefs), homeBriefs)]),
	}
	for _, v := range s.sortedVerticals() {
		home.Verticals = append(home.Verticals, views.VerticalCard{Name: v.Name, URL: verticalURL(v.Slug), Description: v.Description})
	}
	for _, k := range explainerKinds {
		home.Kinds = append(home.Kinds, views.Link{Name: k.plural, URL: kindURL(k.kind)})
	}
	return s.writePage(homeURL, views.Home(home))
}

const homeDescription = "Evergreen explainers on how real software systems are built, and short Briefs on what just happened in tech and why it matters architecturally."

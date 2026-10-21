package builder

import (
	"strings"
	"time"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/schema"
	"hastcode.com/site/internal/views"
)

// Head and structured data shared by every page type (spec "Structured
// data and head"). s.page fills what every page has (canonical, default
// Open Graph); s.articlePage adds what article pages carry (Article or
// NewsArticle JSON-LD, BreadcrumbList, og:type article with the hero).
// The schema builders (s.person, s.organization, s.breadcrumbs) are the
// pieces Author, home and listing pages reuse.

// absURL turns a site-relative path into an absolute URL on SiteURL.
func (s *site) absURL(path string) string {
	return strings.TrimSuffix(s.cfg.SiteURL, "/") + path
}

// articleHead is what an article page's head is built from.
type articleHead struct {
	Type        schema.ArticleType // TypeArticle (Explainer), TypeNewsArticle (Brief)
	URL         string             // site-relative
	Headline    string
	Description string
	Hero        views.Hero
	// Published and Modified keep the frontmatter's offset; they are
	// written as RFC 3339 so they match the visible <time> values.
	Published, Modified time.Time
	Author              *content.Author
	Section             string   // articleSection: the Vertical name for Explainers
	Keywords            []string // tags
	// Breadcrumbs is the trail from Home to this page (site-relative).
	Breadcrumbs []views.Link
}

// articlePage returns an article's views.Page: s.page plus hero preload,
// og:type article with the 16:9 hero crop, the visible breadcrumb trail,
// and the Article/NewsArticle and BreadcrumbList JSON-LD.
func (s *site) articlePage(a articleHead) views.Page {
	p := s.page(a.URL, a.Headline+" | "+views.PublicationName, a.Description)
	hero := a.Hero
	p.HeroPreload = &hero
	og := hero.Crop16x9
	og.URL = s.absURL(og.URL)
	p.Social = views.Social{
		Type:      "article",
		Title:     a.Headline,
		Image:     &og,
		ImageAlt:  hero.Alt,
		Published: a.Published.Format(time.RFC3339),
		Modified:  a.Modified.Format(time.RFC3339),
	}
	p.Breadcrumbs = a.Breadcrumbs
	p.StructuredData = []any{
		schema.Article{
			Context:          schema.Context,
			Type:             a.Type,
			Headline:         a.Headline,
			Description:      a.Description,
			Image:            []string{s.absURL(hero.Crop16x9.URL), s.absURL(hero.Crop4x3.URL), s.absURL(hero.Crop1x1.URL)},
			DatePublished:    a.Published.Format(time.RFC3339),
			DateModified:     a.Modified.Format(time.RFC3339),
			Author:           s.person(a.Author),
			Publisher:        s.organization(),
			MainEntityOfPage: schema.Page(p.Canonical),
			ArticleSection:   a.Section,
			Keywords:         a.Keywords,
		},
		s.breadcrumbs(a.Breadcrumbs),
	}
	return p
}

// person is an Author as a schema.org Person: url is the Author page,
// sameAs the Author's profiles from content/authors/<slug>.md. Set
// Context on the result to use it as a top-level block (Author page).
func (s *site) person(a *content.Author) schema.Person {
	return schema.Person{
		Type:   "Person",
		Name:   a.Name,
		URL:    s.absURL(authorURL(a.Slug)),
		SameAs: a.SameAs,
	}
}

// organization is the publisher, Hast Code. Set Context (and add
// Description/SameAs) to use it as the home page's top-level block.
func (s *site) organization() schema.Organization {
	return schema.Organization{
		Type: "Organization",
		Name: views.PublicationName,
		URL:  s.absURL(homeURL),
		Logo: &schema.ImageObject{Type: "ImageObject", URL: s.absURL(views.MarkPath)},
	}
}

// homeOrganization is the home page's top-level Organization: the
// publisher plus the tagline and the publication's profiles (site.yaml
// same_as).
func (s *site) homeOrganization() schema.Organization {
	org := s.organization()
	org.Context = schema.Context
	org.Description = views.Tagline
	org.SameAs = s.cfg.SameAs
	return org
}

// webSite is the home page's WebSite: the site name and a SearchAction on
// /search/?q= (the Pagefind page reads q).
func (s *site) webSite() schema.WebSite {
	return schema.WebSite{
		Context:         schema.Context,
		Type:            "WebSite",
		Name:            views.PublicationName,
		URL:             s.absURL(homeURL),
		Description:     views.Tagline,
		PotentialAction: schema.NewSearchAction(s.absURL(searchURL) + "?q="),
	}
}

// breadcrumbs is the BreadcrumbList for a trail of site-relative links.
func (s *site) breadcrumbs(trail []views.Link) schema.BreadcrumbList {
	crumbs := make([]schema.Crumb, len(trail))
	for i, l := range trail {
		crumbs[i] = schema.Crumb{Name: l.Name, URL: s.absURL(l.URL)}
	}
	return schema.NewBreadcrumbList(crumbs...)
}

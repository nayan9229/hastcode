package builder

import (
	"cmp"
	"net/url"
	"slices"
	"strings"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/markdown"
	"hastcode.com/site/internal/schema"
	"hastcode.com/site/internal/views"
)

func init() {
	registerStep(phaseRender, 40, "render author pages", renderAuthors)
}

// sortedAuthors are the loaded Authors in slug order.
func (s *site) sortedAuthors() []*content.Author {
	slugs := make([]string, 0, len(s.authors))
	for slug := range s.authors {
		slugs = append(slugs, slug)
	}
	slices.Sort(slugs)
	authors := make([]*content.Author, len(slugs))
	for i, slug := range slugs {
		authors[i] = s.authors[slug]
	}
	return authors
}

// renderAuthors writes an Author page for every Author file that loaded
// without errors, in slug order.
func renderAuthors(s *site) error {
	for _, a := range s.sortedAuthors() {
		if s.failed[a.File] {
			continue
		}
		bio, err := markdown.Render(a.Bio)
		if err != nil {
			s.fail(a.File, "markdown", "%v", err)
			continue
		}
		if err := s.writePage(authorURL(a.Slug), views.AuthorProfile(authorPage(s, a, bio))); err != nil {
			return err
		}
	}
	return nil
}

// authorPage maps an Author to its view model: s.page plus the visible
// breadcrumb, the BreadcrumbList and a top-level Person (the same url and
// sameAs that articles name as their author).
func authorPage(s *site, a *content.Author, bio string) views.AuthorPage {
	page := s.page(authorURL(a.Slug), a.Name+" | "+views.PublicationName, "About "+a.Name+", "+views.PublicationName+" Author.")
	page.Social.Type = "profile"
	page.Breadcrumbs = []views.Link{
		{Name: "Home", URL: homeURL},
		{Name: a.Name, URL: authorURL(a.Slug)},
	}
	person := s.person(a)
	person.Context = schema.Context
	var photo *views.Image
	if ph := s.authorPhotos[a.Slug]; ph != nil && ph.view.URL != "" {
		photo = &ph.view
		person.Image = s.absURL(ph.view.URL)
	}
	page.StructuredData = []any{person, s.breadcrumbs(page.Breadcrumbs)}

	profiles := make([]views.Link, len(a.SameAs))
	for i, u := range a.SameAs {
		profiles[i] = views.Link{Name: profileLabel(u), URL: u}
	}
	return views.AuthorPage{
		Page:     page,
		Name:     a.Name,
		Photo:    photo,
		PhotoAlt: cmp.Or(a.PhotoAlt, a.Name),
		BioHTML:  bio,
		Profiles: profiles,
		Latest:   s.authorPieces(a.Slug),
	}
}

// authorLatestPieces caps the Author page's list of latest pieces.
const authorLatestPieces = 20

// authorPieces are the Author's latest published pieces, newest first.
// Every piece type with a byline adds its own loop here (Briefs: one more
// loop over the published Briefs with Type "Brief").
func (s *site) authorPieces(slug string) []views.AuthorPiece {
	var pieces []views.AuthorPiece
	for _, e := range s.publishedExplainers() {
		if e.Author == slug {
			pieces = append(pieces, views.AuthorPiece{Type: "Explainer", Title: e.Title, URL: explainerURL(e), Published: e.Published})
		}
	}
	slices.SortStableFunc(pieces, func(a, b views.AuthorPiece) int { return b.Published.Compare(a.Published) })
	return pieces[:min(len(pieces), authorLatestPieces)]
}

// profileLabel is a profile URL's host without "www." (github.com,
// linkedin.com): the visible text of the link.
func profileLabel(u string) string {
	p, err := url.Parse(u)
	if err != nil || p.Host == "" {
		return u
	}
	return strings.TrimPrefix(p.Host, "www.")
}

package builder

import (
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/yuin/goldmark/ast"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/markdown"
	"hastcode.com/site/internal/schema"
	"hastcode.com/site/internal/views"
)

// Shared by the Repo Profile steps: mapping a Profile and its Fact Pack
// (loaded by step_load_fact_packs.go) to what the page shows.

// factNotFound is shown for a fact the Fact Pack does not carry; a fact
// is never invented.
const factNotFound = views.FactNotFound

// repoProfileSignals maps a Fact Pack to the Signals section: the plain
// facts the spec lists ("Page composition and gates"), no score.
func repoProfileSignals(fp *content.FactPack) views.RepoProfileSignals {
	s := views.RepoProfileSignals{
		LastCommit:       fp.PushedAt,
		ReleasesLastYear: releasesLastYear(fp),
		Contributors:     fp.Contributors.Total,
		OpenIssues:       fp.OpenIssues,
		SecurityPolicy:   fp.Files.Security,
		License:          licenseLabel(fp.License),
		Stars:            fp.Stars,
		Forks:            fp.Forks,
		Watchers:         fp.Watchers,
		Documentation:    documentation(fp.Files),
	}
	if r := fp.LatestRelease; r != nil {
		s.LatestRelease = views.Link{Name: r.Tag, URL: r.URL}
		s.LatestReleaseDate = r.Date
	}
	if top := fp.Contributors.Top; len(top) > 0 {
		sum := 0
		for _, c := range top {
			sum += c.Contributions
		}
		lead := slices.MaxFunc(top, func(a, b content.Contributor) int { return a.Contributions - b.Contributions })
		if sum > 0 {
			s.TopContributor = views.Link{Name: lead.Login, URL: lead.URL}
			s.TopShare = int(math.Round(float64(lead.Contributions) * 100 / float64(sum)))
			s.ListedContributions = sum
			s.Listed = len(top)
		}
	}
	return s
}

// releasesLastYear counts the listed releases dated within the twelve
// months before the fetch. Without a release list or a fetch date the
// count is not known; when every listed release falls in the window and
// the list is full, the count is a floor ("10 or more").
func releasesLastYear(fp *content.FactPack) string {
	if fp.Releases == nil || fp.FetchedAt.IsZero() {
		return factNotFound
	}
	since := fp.FetchedAt.AddDate(-1, 0, 0)
	n := 0
	for _, r := range fp.Releases {
		if !r.Date.Before(since) {
			n++
		}
	}
	if n == len(fp.Releases) && n >= content.FactPackReleaseCap {
		return strconv.Itoa(n) + " or more"
	}
	return strconv.Itoa(n)
}

// repoProfileChangelog maps a Fact Pack to the Changelog section: its
// releases newest first (the Fact Pack promises that order; sorting keeps
// it true). nil releases stay nil (not found).
func repoProfileChangelog(fp *content.FactPack) views.RepoProfileChangelog {
	var c views.RepoProfileChangelog
	if fp.Releases == nil {
		return c
	}
	c.Releases = []views.ReleaseItem{}
	for _, r := range fp.Releases {
		c.Releases = append(c.Releases, views.ReleaseItem{
			Tag:        views.Link{Name: r.Tag, URL: r.URL},
			Name:       r.Name,
			Date:       r.Date,
			Prerelease: r.Prerelease,
		})
	}
	slices.SortStableFunc(c.Releases, func(a, b views.ReleaseItem) int { return b.Date.Compare(a.Date) })
	return c
}

// softwareSourceCode is the repository as JSON-LD, beside the Article:
// name, repository URL, and only those of language, license URL and
// description the Fact Pack carries (an unclassified license has no URL).
func softwareSourceCode(fp *content.FactPack) schema.SoftwareSourceCode {
	code := schema.SoftwareSourceCode{
		Context:             schema.Context,
		Type:                "SoftwareSourceCode",
		Name:                fp.Slug,
		CodeRepository:      fp.URL,
		ProgrammingLanguage: fp.Language,
		Description:         fp.Description,
	}
	if l := fp.License; l != nil && !l.NeedsReview {
		code.License = l.URL
	}
	return code
}

// faqPage is the Profile's FAQ as JSON-LD: one Question per H3 under ##
// FAQ (headingFAQ, whose count the skeleton rule checks), in order,
// answered by the blocks between it and the next H3 (as HTML). nil when
// the body has no FAQ questions.
func faqPage(e *content.Explainer) (*schema.FAQPage, error) {
	_, sections := e.Body.Sections()
	i := slices.IndexFunc(sections, func(sec markdown.Section) bool { return sec.Heading == headingFAQ })
	if i < 0 {
		return nil, nil
	}
	faq := &schema.FAQPage{Context: schema.Context, Type: "FAQPage"}
	var answer []ast.Node
	flush := func() error {
		if n := len(faq.MainEntity); n > 0 {
			html, err := markdown.RenderBlocks(e.Body, answer)
			if err != nil {
				return err
			}
			faq.MainEntity[n-1].AcceptedAnswer = schema.Answer{Type: "Answer", Text: strings.TrimSpace(html)}
		}
		answer = nil
		return nil
	}
	for _, b := range sections[i].Blocks {
		if h, ok := b.(*ast.Heading); ok && h.Level == 3 {
			if err := flush(); err != nil {
				return nil, err
			}
			faq.MainEntity = append(faq.MainEntity, schema.Question{Type: "Question", Name: e.Body.Text(h)})
			continue
		}
		answer = append(answer, b)
	}
	if err := flush(); err != nil {
		return nil, err
	}
	if len(faq.MainEntity) == 0 {
		return nil, nil
	}
	return faq, nil
}

// documentation lists the documentation files and directories the Fact
// Pack found, in a fixed order; directories end with a slash.
func documentation(f content.RepoFiles) []string {
	var out []string
	for _, p := range []string{f.README, f.Contributing, f.CodeOfConduct} {
		if p != "" {
			out = append(out, p)
		}
	}
	for _, d := range []string{f.Docs, f.Examples} {
		if d != "" {
			out = append(out, d+"/")
		}
	}
	return out
}

// repoProfileCard maps a Repo Profile and its Fact Pack to the header
// card: what it is and who it is for come from the frontmatter, the rest
// from the Fact Pack.
func repoProfileCard(e *content.Explainer, fp *content.FactPack) *views.RepoProfileCard {
	c := &views.RepoProfileCard{
		Repo:       views.Link{Name: fp.Slug, URL: fp.URL},
		WhatItIs:   e.Description,
		WhoFor:     e.WhoFor,
		Stars:      fp.Stars,
		License:    licenseLabel(fp.License),
		Language:   fp.Language,
		LastCommit: fp.PushedAt,
	}
	if c.Language == "" {
		c.Language = factNotFound
	}
	if r := fp.LatestRelease; r != nil {
		c.LatestRelease = views.Link{Name: r.Tag, URL: r.URL}
		c.LatestReleaseDate = r.Date
	}
	return c
}

// licenseLabel is the license as the card shows it: the SPDX id, "not
// found" when GitHub reports none, "not classified" when GitHub could not
// classify it (NOASSERTION; the Author reads the LICENSE file by hand).
func licenseLabel(l *content.License) string {
	switch {
	case l == nil:
		return factNotFound
	case l.NeedsReview:
		return "not classified"
	}
	return l.SPDXID
}

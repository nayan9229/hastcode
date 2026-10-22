package builder

import (
	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/views"
)

// Shared by the Repo Profile steps: mapping a Profile and its Fact Pack
// (loaded by step_load_fact_packs.go) to what the page shows.

// factNotFound is shown for a fact the Fact Pack does not carry; a fact
// is never invented.
const factNotFound = "not found"

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

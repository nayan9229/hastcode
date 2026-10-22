package builder_test

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/internal/buildtest"
)

// The Friday roundup (spec "Editorial docs", User Story 29) is an ordinary
// Brief over the week's Repo Profiles: these tests prove the existing Brief
// gates accept it, so it needs no new piece type.

const (
	roundupFile = "briefs/trending-repositories-2026-10-23.md"
	roundupURL  = "/briefs/trending-repositories-2026-10-23/"
)

// roundupTree holds the five Repo Profiles the roundup links, each with its
// Fact Pack, and the roundup Brief.
func roundupTree() buildtest.Tree {
	tree := buildtest.NewTree()
	for _, p := range buildtest.RoundupProfiles() {
		tree.WithRepoProfile(p)
	}
	return tree.WithBrief(buildtest.RoundupBrief())
}

var briefRules = []string{"brief-length", "brief-sections", "brief-source", "brief-explainer-link", "broken-link", "frontmatter"}

func TestRoundupBrief_PassesTheBriefGatesAndIsWrittenAtItsBriefsURL(t *testing.T) {
	out, res := buildtest.Build(t, buildtest.Config(), roundupTree())

	for _, rule := range briefRules {
		buildtest.RequireNoProblemRule(t, res, roundupFile, rule)
	}
	if len(res.Errors) != 0 {
		t.Fatalf("the roundup must build clean:\n%s", res)
	}

	doc := buildtest.Page(t, out, roundupURL)
	if got := doc.Find("article.h-entry h1.p-name").Text(); got != "Five repositories worth a look this week" {
		t.Errorf("h1 = %q, want the roundup's title", got)
	}
	var headings []string
	doc.Find("article.h-entry .e-content h2").Each(func(_ int, h *goquery.Selection) {
		headings = append(headings, strings.TrimSpace(h.Text()))
	})
	if strings.Join(headings, " | ") != "What happened | Why it matters" {
		t.Errorf("roundup h2s = %q, want What happened then Why it matters", headings)
	}
	if got := buildtest.JSONLDOfType(t, doc, "NewsArticle")["headline"]; got != "Five repositories worth a look this week" {
		t.Errorf("NewsArticle headline = %v", got)
	}
}

func TestRoundupBrief_WhatHappenedLinksEachOfTheFiveProfiles(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), roundupTree())

	var links []string
	buildtest.Page(t, out, roundupURL).Find(".e-content li a").Each(func(_ int, a *goquery.Selection) {
		links = append(links, a.AttrOr("href", ""))
	})
	var want []string
	for _, p := range buildtest.RoundupProfiles() {
		url := "/" + p.Vertical + "/" + p.Slug + "/"
		want = append(want, url)
		if !buildtest.Exists(out, url) {
			t.Errorf("Profile %s was not written", url)
		}
	}
	if strings.Join(links, " ") != strings.Join(want, " ") {
		t.Errorf("What happened links = %q, want the five Profiles %q", links, want)
	}
}

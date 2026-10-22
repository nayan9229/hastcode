package views_test

import (
	"strings"
	"testing"

	"hastcode.com/site/internal/views"
)

// The hero card renderer's text seam (README "Tests"). The builder tests
// prove everything a build shows of the card: the hero pipeline, alt text,
// the 1600x900 size, the crop-safe layout and the same bytes for the same
// Fact Pack. Which words are drawn cannot be read back from pixels, so
// these check the text the renderer lays out (HeroCardText, the same
// layout HeroCardPNG draws): the lines, the cut description, absent facts.

func iggy() views.HeroCard {
	return views.HeroCard{
		Repo:        "apache/iggy",
		Description: "Persistent message streaming platform written in Rust",
		Language:    "Rust",
		Stars:       2468,
		License:     "Apache-2.0",
		Release:     "server-0.5.0",
	}
}

func TestHeroCard_DrawsTheWordmarkRepositoryDescriptionLanguageStarsLicenseAndRelease(t *testing.T) {
	lines := views.HeroCardText(iggy())

	if len(lines) < 4 {
		t.Fatalf("card lines = %q, want the wordmark, repository, description and facts", lines)
	}
	if lines[0] != "hastcode" {
		t.Errorf("first line = %q, want the wordmark", lines[0])
	}
	if lines[1] != "apache/iggy" {
		t.Errorf("second line = %q, want the repository's full name", lines[1])
	}
	if got := strings.Join(lines[2:len(lines)-1], " "); got != "Persistent message streaming platform written in Rust" {
		t.Errorf("description lines = %q, want the whole description (it fits)", got)
	}
	if got := lines[len(lines)-1]; got != "Rust · 2,468 stars · Apache-2.0 · server-0.5.0" {
		t.Errorf("facts line = %q, want language, stars, license and latest release", got)
	}
}

func TestHeroCard_LongDescriptionIsCutWithAnEllipsisAfterTwoLines(t *testing.T) {
	c := iggy()
	c.Description = "Kubernetes Cluster API Provider AWS provides consistent deployment and day 2 operations of self-managed and EKS Kubernetes clusters on AWS, with more words than any card could hold."

	lines := views.HeroCardText(c)

	desc := lines[2 : len(lines)-1]
	if len(desc) != 2 {
		t.Fatalf("description takes %d lines %q, want 2", len(desc), desc)
	}
	if !strings.HasPrefix(desc[0], "Kubernetes Cluster API Provider AWS") || !strings.HasSuffix(desc[1], "…") {
		t.Errorf("description lines = %q, want the start of the description ending in an ellipsis", desc)
	}
	if strings.Contains(strings.Join(desc, " "), "any card could hold") {
		t.Errorf("description lines = %q, want the end cut", desc)
	}
}

func TestHeroCard_LeavesOutAbsentFactsAndRunesTheFontCannotDraw(t *testing.T) {
	c := views.HeroCard{Repo: "php/frankenphp", Description: "🚀 The modern PHP app server ★", Stars: 9001}

	lines := views.HeroCardText(c)

	want := []string{"hastcode", "php/frankenphp", "The modern PHP app server", "9,001 stars"}
	if strings.Join(lines, "|") != strings.Join(want, "|") {
		t.Errorf("card lines = %q, want %q (no language, license or release; no emoji)", lines, want)
	}
	if lines := views.HeroCardText(views.HeroCard{Repo: "php/frankenphp"}); strings.Join(lines, "|") != "hastcode|php/frankenphp" {
		t.Errorf("card lines with no facts = %q, want only the wordmark and repository", lines)
	}
}

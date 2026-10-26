package views

import "strings"

// The hero card (spec "Page composition and gates": generated hero) is the
// hero of a Repo Profile that names no `hero`: drawn at build time from its
// Fact Pack on the light brand palette with the Go fonts, like the OG badge
// (og_badge.go), so it needs no image file and no system font. The builder
// feeds the PNG through the hero pipeline as if it were a sourced file.
//
// Only facts are drawn, never the fetch time, so the same Fact Pack always
// gives the same bytes (the image cache hits) and a changed star count or
// release gives new bytes (every variant is renamed, caches bust). It is
// drawn on the shared card canvas (card.go).

// heroCardDescriptionLines is how many lines the description may take
// before it is cut with an ellipsis.
const heroCardDescriptionLines = 2

// HeroCard is what a hero card shows. Absent facts ("" or 0) are left
// out, never invented.
type HeroCard struct {
	// Repo is the repository's full name, owner/name.
	Repo string
	// Description is the repository's description; it is wrapped and
	// truncated to fit.
	Description string
	// Language is the primary language.
	Language string
	Stars    int
	// License is the SPDX id; "" when GitHub reports none or could not
	// classify it.
	License string
	// Release is the latest release's tag; "" when there is none.
	Release string
}

// Alt is the card's alt text, "{repo}: {description}" (the full
// description, even where the card cuts it).
func (c HeroCard) Alt() string {
	if c.Description == "" {
		return c.Repo
	}
	return c.Repo + ": " + c.Description
}

// HeroCardPNG draws the card: the wordmark, the repository name over an
// accent rule, the description and a line of facts (language, stars,
// license, latest release), with the accent band the OG badge has along
// the bottom edge.
func HeroCardPNG(c HeroCard) []byte {
	return drawCard(heroCardRuns(c), centredRule(heroCardRuleY))
}

// HeroCardText is the text HeroCardPNG draws, one string per line in
// reading order: the wordmark, the repository, the description's lines
// (as wrapped and cut) and the facts line. It is the renderer's text seam
// (README "Tests"; pixels cannot say which words were drawn):
// the same layout feeds both.
func HeroCardText(c HeroCard) []string {
	runs := heroCardRuns(c)
	lines := make([]string, len(runs))
	for i, r := range runs {
		lines[i] = r.text
	}
	return lines
}

// Layout, in pixels at 1600x900. Baselines stay inside the 30 px the
// 1200x630 crop trims top and bottom; the band is the only thing it cuts.
const (
	heroCardWordmarkY = 236

	heroCardRepoMaxSize = 100
	heroCardRepoMinSize = 48
	heroCardRepoY       = 420

	heroCardRuleY = 472 // under the name's descenders

	heroCardDescriptionSize = 42
	heroCardDescriptionY    = 548
	heroCardLineHeight      = 58

	heroCardFactsMaxSize = 34
	heroCardFactsMinSize = 26
	heroCardFactsY       = 720
)

// heroCardRuns lays the card out: every line fitted to cardSafeWidth,
// the repository name and facts by shrinking then cutting, the description
// by wrapping then cutting.
func heroCardRuns(c HeroCard) []cardRun {
	mono, text := goMonoBold(), goRegular()
	fg, muted := hexColour(lightPalette.fg), hexColour(lightPalette.muted)
	runs := []cardRun{{text: wordmark, y: heroCardWordmarkY, wordmark: true}}

	repo, face := fitText(drawable(c.Repo, mono), mono, heroCardRepoMaxSize, heroCardRepoMinSize)
	runs = append(runs, cardRun{text: repo, face: face, colour: fg, y: heroCardRepoY})

	descFace := newFace(text, heroCardDescriptionSize)
	for i, line := range wrapText(drawable(c.Description, text), descFace, cardSafeWidth, heroCardDescriptionLines) {
		runs = append(runs, cardRun{text: line, face: descFace, colour: muted, y: heroCardDescriptionY + i*heroCardLineHeight})
	}

	if facts := heroCardFacts(c); facts != "" {
		line, face := fitText(drawable(facts, text), text, heroCardFactsMaxSize, heroCardFactsMinSize)
		runs = append(runs, cardRun{text: line, face: face, colour: fg, y: heroCardFactsY})
	}
	return runs
}

// heroCardFacts joins the facts the card has: "Rust · 2,468 stars ·
// Apache-2.0 · v1.2.0".
func heroCardFacts(c HeroCard) string {
	var parts []string
	if c.Language != "" {
		parts = append(parts, c.Language)
	}
	if c.Stars > 0 {
		parts = append(parts, humanInt(c.Stars)+" stars")
	}
	if c.License != "" {
		parts = append(parts, c.License)
	}
	if c.Release != "" {
		parts = append(parts, c.Release)
	}
	return strings.Join(parts, " · ")
}

package views

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// The hero card (spec "Page composition and gates": generated hero) is the
// hero of a Repo Profile that names no `hero`: drawn at build time from its
// Fact Pack on the light brand palette with the Go fonts, like the OG badge
// (og_badge.go), so it needs no image file and no system font. The builder
// feeds the PNG through the hero pipeline as if it were a sourced file.
//
// Only facts are drawn, never the fetch time, so the same Fact Pack always
// gives the same bytes (the image cache hits) and a changed star count or
// release gives new bytes (every variant is renamed, caches bust).
//
// Every line is centred within the middle heroCardSafeWidth px and kept off
// the top and bottom edges, so the 16:9, 4:3, 1:1 and 1200x630 crops the
// hero pipeline cuts from it all keep the whole text.
const (
	heroCardWidth  = 1600
	heroCardHeight = 900

	// heroCardSafeWidth is the widest a line may be: the 1:1 crop keeps
	// only the middle 900 px.
	heroCardSafeWidth = 900
	// heroCardDescriptionLines is how many lines the description may take
	// before it is cut with an ellipsis.
	heroCardDescriptionLines = 2

	heroCardEllipsis = "…"
)

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
	rule := image.Rect((heroCardWidth-heroCardRuleWidth)/2, heroCardRuleY, (heroCardWidth+heroCardRuleWidth)/2, heroCardRuleY+heroCardRuleHeight)
	return drawCard(heroCardRuns(c), rule)
}

// drawCard paints a 1600x900 card: the light background, the accent band
// along the bottom edge, the accent rule and each run centred on its
// baseline. The hero card and the title card (title_card.go) share it.
func drawCard(runs []heroCardRun, rule image.Rectangle) []byte {
	img := image.NewRGBA(image.Rect(0, 0, heroCardWidth, heroCardHeight))
	fill(img, img.Bounds(), hexColour(lightPalette.bg))
	fill(img, image.Rect(0, heroCardHeight-heroCardBand, heroCardWidth, heroCardHeight), hexColour(lightPalette.accent))
	fill(img, rule, hexColour(lightPalette.accent))

	for _, r := range runs {
		if r.wordmark {
			drawWordmark(img, heroCardWordmarkSize, (heroCardWidth-wordmarkWidth(heroCardWordmarkSize))/2, r.y)
			continue
		}
		d := &font.Drawer{Dst: img, Src: image.NewUniform(r.colour), Face: r.face, Dot: fixed.P((heroCardWidth-textWidth(r.face, r.text))/2, r.y)}
		d.DrawString(r.text)
	}

	var buf bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := enc.Encode(&buf, img); err != nil {
		panic("views: encode card: " + err.Error())
	}
	return buf.Bytes()
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
	heroCardWordmarkSize = 52
	heroCardWordmarkY    = 236

	heroCardRepoMaxSize = 100
	heroCardRepoMinSize = 48
	heroCardRepoY       = 420

	heroCardRuleY      = 472 // under the name's descenders
	heroCardRuleWidth  = 120
	heroCardRuleHeight = 8

	heroCardDescriptionSize = 42
	heroCardDescriptionY    = 548
	heroCardLineHeight      = 58

	heroCardFactsMaxSize = 34
	heroCardFactsMinSize = 26
	heroCardFactsY       = 720

	heroCardBand = 40
)

// heroCardRun is one centred line of text.
type heroCardRun struct {
	text     string
	face     font.Face
	colour   color.RGBA
	y        int  // baseline
	wordmark bool // drawn by drawWordmark (og_badge.go), dot included
}

// heroCardRuns lays the card out: every line fitted to heroCardSafeWidth,
// the repository name and facts by shrinking then cutting, the description
// by wrapping then cutting.
func heroCardRuns(c HeroCard) []heroCardRun {
	mono, text := goMonoBold(), goRegular()
	fg, muted := hexColour(lightPalette.fg), hexColour(lightPalette.muted)
	runs := []heroCardRun{{text: wordmark, y: heroCardWordmarkY, wordmark: true}}

	repo, face := fitText(drawable(c.Repo, mono), mono, heroCardRepoMaxSize, heroCardRepoMinSize)
	runs = append(runs, heroCardRun{text: repo, face: face, colour: fg, y: heroCardRepoY})

	descFace := newFace(text, heroCardDescriptionSize)
	for i, line := range wrapText(drawable(c.Description, text), descFace, heroCardSafeWidth, heroCardDescriptionLines) {
		runs = append(runs, heroCardRun{text: line, face: descFace, colour: muted, y: heroCardDescriptionY + i*heroCardLineHeight})
	}

	if facts := heroCardFacts(c); facts != "" {
		line, face := fitText(drawable(facts, text), text, heroCardFactsMaxSize, heroCardFactsMinSize)
		runs = append(runs, heroCardRun{text: line, face: face, colour: fg, y: heroCardFactsY})
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

// drawable is s with its whitespace collapsed and every rune the font has
// no glyph for (emoji, symbols) dropped, so nothing renders as a box.
func drawable(s string, f *opentype.Font) string {
	var buf sfnt.Buffer
	var b strings.Builder
	for _, r := range strings.Join(strings.Fields(s), " ") {
		if idx, err := f.GlyphIndex(&buf, r); err != nil || idx == 0 {
			continue
		}
		b.WriteRune(r)
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

// fitText shrinks s from maxSize in 4 px steps until it fits
// heroCardSafeWidth; at minSize it is cut with an ellipsis instead.
func fitText(s string, f *opentype.Font, maxSize, minSize float64) (string, font.Face) {
	size := maxSize
	face := newFace(f, size)
	for textWidth(face, s) > heroCardSafeWidth && size-4 >= minSize {
		size -= 4
		face = newFace(f, size)
	}
	return truncateText(s, face, heroCardSafeWidth), face
}

// wrapText breaks text into at most maxLines lines no wider than maxWidth,
// on spaces; what does not fit is cut from the last line with an ellipsis.
func wrapText(text string, face font.Face, maxWidth, maxLines int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}
	var lines []string
	line := words[0]
	for i := 1; i < len(words); i++ {
		if next := line + " " + words[i]; textWidth(face, next) <= maxWidth {
			line = next
			continue
		}
		if len(lines) == maxLines-1 {
			return append(lines, truncateText(line+" "+strings.Join(words[i:], " "), face, maxWidth))
		}
		lines = append(lines, truncateText(line, face, maxWidth))
		line = words[i]
	}
	return append(lines, truncateText(line, face, maxWidth))
}

// truncateText cuts s to fit maxWidth, dropping whole words from the end
// (runes, for a single word) and ending it with an ellipsis.
func truncateText(s string, face font.Face, maxWidth int) string {
	if textWidth(face, s) <= maxWidth {
		return s
	}
	words := strings.Fields(s)
	for len(words) > 1 {
		words = words[:len(words)-1]
		if t := strings.Join(words, " ") + heroCardEllipsis; textWidth(face, t) <= maxWidth {
			return t
		}
	}
	runes := []rune(strings.Join(words, ""))
	for len(runes) > 0 {
		runes = runes[:len(runes)-1]
		if t := string(runes) + heroCardEllipsis; textWidth(face, t) <= maxWidth {
			return t
		}
	}
	return heroCardEllipsis
}

package views

import (
	"image"
	"strings"
)

// The title card (ticket 10) is the hero of a Case Study, Concept,
// Comparison or Brief that names no `hero`: drawn at build time from the
// Piece's title, kind and Vertical on the light brand palette with the Go
// fonts, like the hero card and the OG badge, so no Piece needs a
// hand-made image. The builder feeds the PNG through the hero pipeline as
// if it were a sourced file.
//
// It is the hero card's canvas (1600x900, accent band along the bottom,
// wordmark at the top) with the title in its middle: the kind and Vertical
// as an eyebrow, the title under it, the accent rule under the title. The
// eyebrow, title and rule are centred as one block below the wordmark,
// and every line keeps within heroCardSafeWidth and off the edges the
// 1200x630 crop trims, so every crop the hero pipeline cuts keeps the
// whole text. Only the title, kind and Vertical are drawn, so the same
// Piece always gives the same bytes and a retitled one new bytes.

// TitleCard is what a title card shows.
type TitleCard struct {
	// Title is the Piece's title; it is wrapped, shrunk and cut to fit.
	Title string
	// Kind is the kind label ("Case Study", "Brief").
	Kind string
	// Vertical is the Vertical's name; "" for a Brief.
	Vertical string
}

// Alt is the card's alt text: the title, which is what the card shows.
func (c TitleCard) Alt() string { return c.Title }

// TitleCardPNG draws the card.
func TitleCardPNG(c TitleCard) []byte {
	runs, rule := titleCardRuns(c)
	return drawCard(runs, rule)
}

// Layout, in pixels at 1600x900.
const (
	titleCardWordmarkY = 170

	// The eyebrow, title and rule are centred as one block between
	// titleCardBlockTop and titleCardBlockBottom.
	titleCardBlockTop    = 250
	titleCardBlockBottom = 820

	titleCardEyebrowMaxSize = 30
	titleCardEyebrowMinSize = 22
	titleCardEyebrowCap     = 22 // the eyebrow's cap height above its baseline, at most

	titleCardTitleMaxSize = 76
	titleCardTitleMinSize = 48
	titleCardTitleLines   = 4
	titleCardEyebrowGap   = 40 // eyebrow baseline to the title's cap height
)

// titleCardRuns lays the card out: the wordmark, the eyebrow (kind and
// Vertical, upper case, shrunk then cut), the title (the largest size at
// which it wraps into titleCardTitleLines lines, cut with an ellipsis at
// the smallest) and the rule under it.
func titleCardRuns(c TitleCard) ([]heroCardRun, image.Rectangle) {
	mono := goMonoBold()
	fg, accent := hexColour(lightPalette.fg), hexColour(lightPalette.link)

	eyebrowText, eyebrowFace := fitText(drawable(strings.ToUpper(titleCardEyebrow(c)), mono), mono, titleCardEyebrowMaxSize, titleCardEyebrowMinSize)

	title := drawable(c.Title, mono)
	var lines []string
	size := float64(titleCardTitleMaxSize)
	for ; ; size -= 4 {
		lines = wrapText(title, newFace(mono, size), heroCardSafeWidth, titleCardTitleLines)
		if strings.Join(lines, " ") == title || size-4 < titleCardTitleMinSize {
			break
		}
	}
	titleFace := newFace(mono, size)
	lineHeight := int(size * 1.22)
	capHeight := int(size * 0.75)
	ruleGap := int(size * 0.55) // below the last baseline, clear of descenders

	height := titleCardEyebrowCap + titleCardEyebrowGap + capHeight + (len(lines)-1)*lineHeight + ruleGap + heroCardRuleHeight
	top := titleCardBlockTop + (titleCardBlockBottom-titleCardBlockTop-height)/2

	runs := []heroCardRun{{text: wordmark, y: titleCardWordmarkY, wordmark: true}}
	y := top + titleCardEyebrowCap
	runs = append(runs, heroCardRun{text: eyebrowText, face: eyebrowFace, colour: accent, y: y})
	y += titleCardEyebrowGap + capHeight
	for i, line := range lines {
		if i > 0 {
			y += lineHeight
		}
		runs = append(runs, heroCardRun{text: line, face: titleFace, colour: fg, y: y})
	}
	ruleY := y + ruleGap
	rule := image.Rect((heroCardWidth-heroCardRuleWidth)/2, ruleY, (heroCardWidth+heroCardRuleWidth)/2, ruleY+heroCardRuleHeight)
	return runs, rule
}

// titleCardEyebrow is "Case Study · Distributed Systems", or the kind
// alone ("Brief").
func titleCardEyebrow(c TitleCard) string {
	if c.Vertical == "" {
		return c.Kind
	}
	return c.Kind + " · " + c.Vertical
}

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

// The card canvas the hero card (hero_card.go) and the title card
// (title_card.go) share: 1600x900 on the light brand palette, the wordmark
// and every line centred, an accent rule, and the accent band the OG badge
// has along the bottom edge. Every line keeps within the middle
// cardSafeWidth px and off the top and bottom edges, so the 16:9, 4:3, 1:1
// and 1200x630 crops the hero pipeline cuts from it all keep the whole
// text.
const (
	cardWidth  = 1600
	cardHeight = 900

	// cardSafeWidth is the widest a line may be: the 1:1 crop keeps only
	// the middle 900 px.
	cardSafeWidth = 900

	cardWordmarkSize = 52

	cardRuleWidth  = 120
	cardRuleHeight = 8

	// cardBand is the accent band's height; the 1200x630 crop trims top
	// and bottom, and the band is the only thing it cuts.
	cardBand = 40

	cardEllipsis = "…"
)

// cardRun is one centred line of text.
type cardRun struct {
	text     string
	face     font.Face
	colour   color.RGBA
	y        int  // baseline
	wordmark bool // drawn by drawWordmark (og_badge.go), dot included
}

// centredRule is the accent rule, cardRuleWidth wide and centred, with
// its top edge at y.
func centredRule(y int) image.Rectangle {
	return image.Rect((cardWidth-cardRuleWidth)/2, y, (cardWidth+cardRuleWidth)/2, y+cardRuleHeight)
}

// drawCard paints a 1600x900 card: the light background, the accent band
// along the bottom edge, the accent rule and each run centred on its
// baseline. The hero card and the title card (title_card.go) share it.
func drawCard(runs []cardRun, rule image.Rectangle) []byte {
	img := image.NewRGBA(image.Rect(0, 0, cardWidth, cardHeight))
	fill(img, img.Bounds(), hexColour(lightPalette.bg))
	fill(img, image.Rect(0, cardHeight-cardBand, cardWidth, cardHeight), hexColour(lightPalette.accent))
	fill(img, rule, hexColour(lightPalette.accent))

	for _, r := range runs {
		if r.wordmark {
			drawWordmark(img, cardWordmarkSize, (cardWidth-wordmarkWidth(cardWordmarkSize))/2, r.y)
			continue
		}
		d := &font.Drawer{Dst: img, Src: image.NewUniform(r.colour), Face: r.face, Dot: fixed.P((cardWidth-textWidth(r.face, r.text))/2, r.y)}
		d.DrawString(r.text)
	}

	var buf bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := enc.Encode(&buf, img); err != nil {
		panic("views: encode card: " + err.Error())
	}
	return buf.Bytes()
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
// cardSafeWidth; at minSize it is cut with an ellipsis instead.
func fitText(s string, f *opentype.Font, maxSize, minSize float64) (string, font.Face) {
	size := maxSize
	face := newFace(f, size)
	for textWidth(face, s) > cardSafeWidth && size-4 >= minSize {
		size -= 4
		face = newFace(f, size)
	}
	return truncateText(s, face, cardSafeWidth), face
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
		if t := strings.Join(words, " ") + cardEllipsis; textWidth(face, t) <= maxWidth {
			return t
		}
	}
	runes := []rune(strings.Join(words, ""))
	for len(runes) > 0 {
		runes = runes[:len(runes)-1]
		if t := string(runes) + cardEllipsis; textWidth(face, t) <= maxWidth {
			return t
		}
	}
	return cardEllipsis
}

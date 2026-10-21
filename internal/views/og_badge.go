package views

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strconv"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// The OG badge (spec "Brand": generated OG badge, no bitmap logo assets):
// the default Open Graph / Twitter large-card image for every page without
// a hero. It is drawn at build time from the brand constants (wordmark,
// tagline, light palette) with the Go fonts, so it needs no image file and
// no system font, and the same build always writes the same bytes.
const (
	ogBadgeWidth  = 1200 // Discover and Open Graph want ≥ 1200 px
	ogBadgeHeight = 630  // the 1.91:1 large-card ratio

	// OGBadgeAlt describes the badge (og:image:alt).
	OGBadgeAlt = PublicationName + ": " + Tagline
)

// OGBadge is the badge's view model; URL is site-relative and
// content-hashed under /assets/ (cached immutably).
func OGBadge() Image {
	return Image{URL: ogBadgeFile().path, Width: ogBadgeWidth, Height: ogBadgeHeight}
}

type badgeFile struct {
	path string
	png  []byte
}

var ogBadgeFile = sync.OnceValue(func() badgeFile {
	data := drawOGBadge()
	return badgeFile{path: "/assets/og-badge." + contentHash(data) + ".png", png: data}
})

// drawOGBadge lays out the wordmark (monospace bold with the teal dot the
// header draws in CSS) and the tagline on the light background, with an
// accent band along the bottom edge.
func drawOGBadge() []byte {
	img := image.NewRGBA(image.Rect(0, 0, ogBadgeWidth, ogBadgeHeight))
	fill(img, img.Bounds(), hexColour(lightPalette.bg))
	fill(img, image.Rect(0, ogBadgeHeight-24, ogBadgeWidth, ogBadgeHeight), hexColour(lightPalette.accent))

	const left = 96
	word := &font.Drawer{Dst: img, Src: image.NewUniform(hexColour(lightPalette.fg)), Face: goFace(gomonobold.TTF, 144), Dot: fixed.P(left, 320)}
	word.DrawString("hastcode")
	// The dot: .34em wide, .06em after the word, sitting on the baseline.
	r := 144 * 0.17
	cx := float64(word.Dot.X.Round()) + 144*0.06 + r
	cy := 320 - r
	disc(img, cx, cy, r, hexColour(lightPalette.accent))

	tag := &font.Drawer{Dst: img, Src: image.NewUniform(hexColour(lightPalette.muted)), Face: goFace(goregular.TTF, 56), Dot: fixed.P(left, 420)}
	tag.DrawString(Tagline)

	var buf bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := enc.Encode(&buf, img); err != nil {
		panic("views: encode OG badge: " + err.Error())
	}
	return buf.Bytes()
}

func goFace(ttf []byte, size float64) font.Face {
	f, err := opentype.Parse(ttf)
	if err != nil {
		panic("views: parse Go font: " + err.Error())
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		panic("views: Go font face: " + err.Error())
	}
	return face
}

func fill(img *image.RGBA, r image.Rectangle, c color.Color) {
	draw.Draw(img, r, image.NewUniform(c), image.Point{}, draw.Src)
}

// disc paints a filled circle, anti-aliased by 4x4 supersampling.
func disc(img *image.RGBA, cx, cy, r float64, c color.RGBA) {
	for y := int(cy - r - 1); y <= int(cy+r+1); y++ {
		for x := int(cx - r - 1); x <= int(cx+r+1); x++ {
			in := 0
			for sy := range 4 {
				for sx := range 4 {
					dx := float64(x) + (float64(sx)+0.5)/4 - cx
					dy := float64(y) + (float64(sy)+0.5)/4 - cy
					if dx*dx+dy*dy <= r*r {
						in++
					}
				}
			}
			if in == 0 {
				continue
			}
			bg := img.RGBAAt(x, y)
			blend := func(a, b uint8) uint8 { return uint8((int(a)*in + int(b)*(16-in) + 8) / 16) }
			img.SetRGBA(x, y, color.RGBA{blend(c.R, bg.R), blend(c.G, bg.G), blend(c.B, bg.B), 0xFF})
		}
	}
}

// hexColour parses a palette value ("#RRGGBB").
func hexColour(h string) color.RGBA {
	v, err := strconv.ParseUint(h[1:], 16, 32)
	if err != nil || len(h) != 7 {
		panic("views: bad palette colour " + h)
	}
	return color.RGBA{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v), A: 0xFF}
}

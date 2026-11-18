package views

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
	"strconv"
	"strings"
	"sync"
)

// TouchIconPath is the crawlable icon (ADR-0013, amending ADR-0006): the
// mark as a square PNG at a stable path (no content hash, so outside
// /assets/), linked from the home page only with rel="apple-touch-icon".
// Google's "Favicon in Search" accepts that rel and wants a file
// Googlebot-Image can crawl; every page keeps the data: URI favicon for
// browsers. iPhones request this path on their own, linked or not.
const TouchIconPath = "/apple-touch-icon.png"

// touchIconSize is the icon's side in pixels: a multiple of 48 (Google)
// and the 192 px Android and iOS home-screen size.
const touchIconSize = 192

// The mark's geometry, shared by markSVG and drawTouchIcon so the "h"
// never drifts between the two files: a 32-unit square, the stroke's path
// data and width. The SVG's corner radius (rx="7") is SVG-only on
// purpose: the PNG stays square because iOS and Google mask its corners.
const (
	markUnits       = 32
	markStrokePath  = "M10.5 7v18M10.5 17.5c0-3.3 2.4-5.5 5.5-5.5s5.5 2.2 5.5 5.5V25"
	markStrokeWidth = 3.4
)

// TouchIconPNG is the icon's bytes, the same on every build.
var TouchIconPNG = sync.OnceValue(drawTouchIcon)

// drawTouchIcon rasterises the mark in the light palette: the white "h"
// (round caps and joins) on the brand teal. The ground fills the whole
// square, opaque, with no rounded corners: iOS fills transparent pixels
// with black and masks the corners itself, and Google crops its own shape.
// Each pixel is 4x4 supersampled for anti-aliased edges.
func drawTouchIcon() []byte {
	const n = touchIconSize
	img := image.NewRGBA(image.Rect(0, 0, n, n))
	ground, ink := hexColour(lightPalette.accent), hexColour(lightPalette.bg)
	lines := flattenPath(markStrokePath)
	half := markStrokeWidth / 2
	scale := float64(markUnits) / n
	for y := range n {
		for x := range n {
			in := 0
			for sy := range 4 {
				for sx := range 4 {
					px := (float64(x) + (float64(sx)+0.5)/4) * scale
					py := (float64(y) + (float64(sy)+0.5)/4) * scale
					if nearPolylines(lines, px, py, half) {
						in++
					}
				}
			}
			blend := func(a, b uint8) uint8 { return uint8((int(a)*in + int(b)*(16-in) + 8) / 16) }
			img.SetRGBA(x, y, color.RGBA{blend(ink.R, ground.R), blend(ink.G, ground.G), blend(ink.B, ground.B), 0xFF})
		}
	}
	var buf bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	if err := enc.Encode(&buf, img); err != nil {
		panic("views: encode touch icon: " + err.Error())
	}
	return buf.Bytes()
}

type point struct{ x, y float64 }

// nearPolylines reports whether (x, y) lies within r of any polyline: the
// area a stroke of width 2r with round caps and joins paints.
func nearPolylines(lines [][]point, x, y, r float64) bool {
	for _, l := range lines {
		for i := 1; i < len(l); i++ {
			if segmentDistance(l[i-1], l[i], x, y) <= r {
				return true
			}
		}
	}
	return false
}

func segmentDistance(a, b point, x, y float64) float64 {
	dx, dy := b.x-a.x, b.y-a.y
	t := 0.0
	if l := dx*dx + dy*dy; l > 0 {
		t = math.Max(0, math.Min(1, ((x-a.x)*dx+(y-a.y)*dy)/l))
	}
	return math.Hypot(x-(a.x+t*dx), y-(a.y+t*dy))
}

// flattenPath turns SVG path data into polylines, one per subpath. It
// reads only the commands markStrokePath uses (M, V, v, c, s) and panics
// on any other. The input is a constant, so this is a programming error
// like a bad regexp.MustCompile: TestTouchIcon_* in `make test` panics
// first, and so would any build, before a page ships.
func flattenPath(d string) [][]point {
	var (
		lines [][]point
		cur   point
		ctrl  point // the last cubic's second control point
	)
	toks := pathTokens(d)
	num := func(i *int) float64 {
		v, err := strconv.ParseFloat(toks[*i], 64)
		if err != nil {
			panic("views: mark path: " + err.Error())
		}
		*i++
		return v
	}
	cubic := func(c1, c2, end point) {
		l := &lines[len(lines)-1]
		const steps = 32
		for k := 1; k <= steps; k++ {
			t := float64(k) / steps
			u := 1 - t
			*l = append(*l, point{
				u*u*u*cur.x + 3*u*u*t*c1.x + 3*u*t*t*c2.x + t*t*t*end.x,
				u*u*u*cur.y + 3*u*u*t*c1.y + 3*u*t*t*c2.y + t*t*t*end.y,
			})
		}
		cur, ctrl = end, c2
	}
	for i := 0; i < len(toks); {
		cmd := toks[i]
		i++
		switch cmd {
		case "M":
			cur = point{num(&i), num(&i)}
			ctrl = cur
			lines = append(lines, []point{cur})
		case "V", "v":
			y := num(&i)
			if cmd == "v" {
				y += cur.y
			}
			cur = point{cur.x, y}
			ctrl = cur
			lines[len(lines)-1] = append(lines[len(lines)-1], cur)
		case "c":
			c1 := point{cur.x + num(&i), cur.y + num(&i)}
			c2 := point{cur.x + num(&i), cur.y + num(&i)}
			cubic(c1, c2, point{cur.x + num(&i), cur.y + num(&i)})
		case "s":
			c1 := point{2*cur.x - ctrl.x, 2*cur.y - ctrl.y}
			c2 := point{cur.x + num(&i), cur.y + num(&i)}
			cubic(c1, c2, point{cur.x + num(&i), cur.y + num(&i)})
		default:
			panic("views: mark path command " + strconv.Quote(cmd) + " is not supported")
		}
	}
	return lines
}

// pathTokens splits path data into commands and numbers ("c0-3.3" is
// "c", "0", "-3.3").
func pathTokens(d string) []string {
	var toks []string
	var b strings.Builder
	flush := func() {
		if b.Len() > 0 {
			toks = append(toks, b.String())
			b.Reset()
		}
	}
	for _, r := range d {
		switch {
		case r == ' ' || r == ',':
			flush()
		case r == '-':
			flush()
			b.WriteRune(r)
		case (r >= '0' && r <= '9') || r == '.':
			b.WriteRune(r)
		default:
			flush()
			toks = append(toks, string(r))
		}
	}
	flush()
	return toks
}

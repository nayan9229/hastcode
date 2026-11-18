package views

import (
	"fmt"
	"net/url"
)

// Brand constants (spec "Brand"). Use these for OG defaults, feed and
// Organization descriptions rather than retyping the strings.
const (
	PublicationName = "HastCode"
	Tagline         = "How real systems are built."
)

// palette is one colour scheme. bg, fg and accent are the spec's brand
// values; the rest are derived tokens for chrome and content.
type palette struct {
	bg, fg, accent string
	// link is the ink for link text. In light mode the brand teal
	// (#0D9488, 3.7:1 on white) is too light for AA body text, so links use
	// a darker teal; in dark mode the accent itself passes comfortably.
	link    string
	muted   string // secondary text (dates, captions, footer)
	border  string // hairlines, table rules
	surface string // raised panels, code blocks, placeholders
}

var (
	lightPalette = palette{
		bg: "#FFFFFF", fg: "#111827", accent: "#0D9488",
		link: "#0F766E", muted: "#4B5563", border: "#E5E7EB", surface: "#F3F4F6",
	}
	darkPalette = palette{
		bg: "#0B0F19", fg: "#F3F4F6", accent: "#2DD4BF",
		link: "#2DD4BF", muted: "#9CA3AF", border: "#1F2937", surface: "#111827",
	}
)

// MarkPath is the site path of the generated SVG mark: the favicon and the
// Organization logo (use SiteURL + MarkPath for the absolute logo URL).
const MarkPath = "/mark.svg"

// markDataURI is the SVG mark embedded as the favicon's data: URI, so the
// browser makes no request for it before first paint (ADR-0006). The file
// at MarkPath stays: it is the Organization logo.
var markDataURI = "data:image/svg+xml," + url.PathEscape(string(markSVG()))

// markSVG is the brand mark: a lowercase "h" stroke on a rounded teal
// square, inverting to the dark palette when the viewer prefers dark.
// Drawn as paths, not text, so it never depends on an installed font. The
// square and the "h" stroke are shared with the PNG icon (drawTouchIcon);
// the rounded corner (rx) is the SVG's own, as the PNG stays square.
func markSVG() []byte {
	return fmt.Appendf(nil, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %[1]d %[1]d" width="%[1]d" height="%[1]d">`+
		`<title>%[2]s</title>`+
		`<style>.b{fill:%[3]s}.h{stroke:%[4]s}@media (prefers-color-scheme:dark){.b{fill:%[5]s}.h{stroke:%[6]s}}</style>`+
		`<rect class="b" width="%[1]d" height="%[1]d" rx="7"/>`+
		`<path class="h" d="%[7]s" fill="none" stroke-width="%[8]g" stroke-linecap="round" stroke-linejoin="round"/>`+
		`</svg>`+"\n",
		markUnits, PublicationName, lightPalette.accent, lightPalette.bg, darkPalette.accent, darkPalette.bg,
		markStrokePath, markStrokeWidth)
}

// StaticAssets are the non-page files every build writes into the static
// tree, by site path: the mark and the crawlable icon drawn from it, both
// at stable names. The OG badge is content-hashed under /assets/, so it
// goes through the Image Store instead (OGBadgePNG).
func StaticAssets() map[string][]byte {
	return map[string][]byte{
		MarkPath:      markSVG(),
		TouchIconPath: TouchIconPNG(),
	}
}

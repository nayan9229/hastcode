package views

import "fmt"

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

// markSVG is the brand mark: a lowercase "h" stroke on a rounded teal
// square, inverting to the dark palette when the viewer prefers dark.
// Drawn as paths, not text, so it never depends on an installed font.
func markSVG() []byte {
	return fmt.Appendf(nil, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="32" height="32">`+
		`<title>%s</title>`+
		`<style>.b{fill:%s}.h{stroke:%s}@media (prefers-color-scheme:dark){.b{fill:%s}.h{stroke:%s}}</style>`+
		`<rect class="b" width="32" height="32" rx="7"/>`+
		`<path class="h" d="M10.5 7v18M10.5 17.5c0-3.3 2.4-5.5 5.5-5.5s5.5 2.2 5.5 5.5V25" fill="none" stroke-width="3.4" stroke-linecap="round" stroke-linejoin="round"/>`+
		`</svg>`+"\n",
		PublicationName, lightPalette.accent, lightPalette.bg, darkPalette.accent, darkPalette.bg)
}

// StaticAssets are the non-page files every build writes into the static
// tree, by site path. The OG badge is an image, so it goes through the
// Image Store instead (OGBadgePNG).
func StaticAssets() map[string][]byte {
	return map[string][]byte{
		MarkPath:       markSVG(),
		stylesheetPath: stylesheet,
	}
}

package views

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

// CSS delivery (ADR-0006, superseding ADR-0003): the stylesheet is inlined
// once in every page's <head>, so text paints without a second request. The
// 50 KB per-page HTML gate does not count it (InlineStyle).

//go:embed site.css
var baseCSS string

var (
	stylesheet  = buildStylesheet()
	inlineStyle = []byte("<style>" + string(stylesheet) + "</style>")
)

// InlineStyle is the exact <style> element every page's <head> carries. The
// size gate subtracts it once, so the gate measures content, not design.
func InlineStyle() []byte { return inlineStyle }

// buildStylesheet is the palette tokens (from brand.go) followed by
// site.css, lightly minified.
func buildStylesheet() []byte {
	return []byte(paletteCSS() + minifyLines(baseCSS))
}

// paletteCSS defines the colour tokens: light by default, dark when the
// viewer prefers dark (unless the toggle forced light), and dark when the
// toggle forced it via <html data-theme="dark">.
func paletteCSS() string {
	return fmt.Sprintf(":root{color-scheme:light;%s}\n"+
		"@media (prefers-color-scheme:dark){:root:not([data-theme=light]){color-scheme:dark;%s}}\n"+
		":root[data-theme=dark]{color-scheme:dark;%s}\n",
		lightPalette.vars(), darkPalette.vars(), darkPalette.vars())
}

func (p palette) vars() string {
	return fmt.Sprintf("--bg:%s;--fg:%s;--accent:%s;--link:%s;--muted:%s;--border:%s;--surface:%s",
		p.bg, p.fg, p.accent, p.link, p.muted, p.border, p.surface)
}

//go:embed theme.js
var themeJS string

// themeScript is the inline theme toggle (layout_head.templ).
var themeScript = strings.TrimSuffix(minifyLines(themeJS), "\n")

var blockComment = regexp.MustCompile(`(?s)/\*.*?\*/`)

// minifyLines drops /* */ comments, indentation and blank lines from CSS
// or JS. It never joins two lines without a separator, so it cannot fuse
// tokens. Sources must not contain "/*" inside strings.
func minifyLines(src string) string {
	var lines []string
	for _, l := range strings.Split(blockComment.ReplaceAllString(src, ""), "\n") {
		if l = strings.TrimSpace(l); l != "" {
			lines = append(lines, l)
		}
	}
	return strings.Join(lines, "\n") + "\n"
}

func contentHash(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])[:12]
}

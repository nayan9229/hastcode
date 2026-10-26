package builder_test

import (
	"encoding/xml"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildtest"
)

const explainerPath = "/distributed-systems/raft-leader-election/"

func TestSiteHeader_ShowsWordmarkTaglineAndNav(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	header := buildtest.Page(t, out, explainerPath).Find("body > header").First()

	wordmark := header.Find(`a[href="/"]`)
	if got := strings.TrimSpace(wordmark.Text()); got != "hastcode" {
		t.Errorf("wordmark linking home = %q, want hastcode", got)
	}
	if !strings.Contains(header.Text(), "How real systems are built.") {
		t.Errorf("header does not contain the tagline; text = %q", header.Text())
	}

	links := map[string]string{}
	header.Find("nav a").Each(func(_ int, a *goquery.Selection) {
		links[a.AttrOr("href", "")] = strings.TrimSpace(a.Text())
	})
	want := map[string]string{
		"/distributed-systems/":  "Distributed Systems",
		"/databases-storage/":    "Databases & Storage",
		"/apis-networking/":      "APIs & Networking",
		"/infrastructure-cloud/": "Infrastructure & Cloud",
		"/ai-ml-systems/":        "AI/ML Systems",
		"/languages-runtimes/":   "Languages & Runtimes",
		"/briefs/":               "Briefs",
		"/search/":               "Search",
	}
	for href, text := range want {
		if got, ok := links[href]; !ok || got != text {
			t.Errorf("nav link %s = %q (present %v), want %q", href, got, ok, text)
		}
	}
}

func TestHead_FaviconIsTheGeneratedSVGMark(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	icon := buildtest.Page(t, out, explainerPath).Find(`head link[rel="icon"]`)

	if icon.Length() != 1 {
		t.Fatalf("found %d <link rel=icon>, want 1", icon.Length())
	}
	if got := icon.AttrOr("type", ""); got != "image/svg+xml" {
		t.Errorf("favicon type = %q, want image/svg+xml", got)
	}
	// Embedded as a data: URI (ADR-0006): no request before first paint.
	href := icon.AttrOr("href", "")
	const prefix = "data:image/svg+xml,"
	if !strings.HasPrefix(href, prefix) {
		t.Fatalf("favicon href = %.60q…, want an embedded %s URI", href, prefix)
	}
	unescaped, err := url.PathUnescape(strings.TrimPrefix(href, prefix))
	if err != nil {
		t.Fatalf("favicon data URI does not unescape: %v", err)
	}
	data := []byte(unescaped)
	if file, err := os.ReadFile(buildtest.StaticPath(out, "/mark.svg")); err != nil || string(file) != unescaped {
		t.Errorf("embedded favicon differs from /mark.svg (the Organization logo file): %v", err)
	}
	var svg struct {
		XMLName xml.Name
		ViewBox string `xml:"viewBox,attr"`
	}
	if err := xml.Unmarshal(data, &svg); err != nil {
		t.Fatalf("favicon is not well-formed XML: %v", err)
	}
	if svg.XMLName.Local != "svg" || svg.XMLName.Space != "http://www.w3.org/2000/svg" || svg.ViewBox == "" {
		t.Errorf("favicon root = %v viewBox=%q, want an <svg> in the SVG namespace with a viewBox", svg.XMLName, svg.ViewBox)
	}
	if !strings.Contains(strings.ToUpper(string(data)), "#0D9488") {
		t.Errorf("favicon does not use the brand teal #0D9488")
	}
}

func TestHead_ThemeColorForBothColourSchemes(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	doc := buildtest.Page(t, out, explainerPath)

	for media, want := range map[string]string{
		"(prefers-color-scheme: light)": "#FFFFFF",
		"(prefers-color-scheme: dark)":  "#0B0F19",
	} {
		m := doc.Find(`head meta[name="theme-color"][media="` + media + `"]`)
		if got := m.AttrOr("content", ""); !strings.EqualFold(got, want) {
			t.Errorf("theme-color for %s = %q, want %s", media, got, want)
		}
	}
}

// ADR-0006: the stylesheet is inlined once in every page's <head>, so text
// paints without waiting for a second request.
func TestStylesheet_IsInlinedOnceInTheHeadWithSystemFontsOnly(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	doc := buildtest.Page(t, out, explainerPath)

	if n := doc.Find(`link[rel="stylesheet"]`).Length(); n != 0 {
		t.Errorf("found %d stylesheet links, want none (CSS is inlined)", n)
	}
	if n := doc.Find("head style").Length(); n != 1 {
		t.Fatalf("found %d <style> elements in <head>, want exactly 1", n)
	}
	css := stylesheet(t, out, doc)
	if len(css) < 1000 {
		t.Fatalf("inline stylesheet is %d bytes; want the whole site stylesheet", len(css))
	}
	for _, banned := range []string{"@font-face", "@import", "fonts.googleapis", "url("} {
		if strings.Contains(css, banned) {
			t.Errorf("stylesheet contains %q; want system fonts and no external resources", banned)
		}
	}
	// Every page carries the same stylesheet.
	if home := stylesheet(t, out, buildtest.Page(t, out, "/")); home != css {
		t.Errorf("home page stylesheet differs from the Explainer's (%d vs %d bytes)", len(home), len(css))
	}
	// No stylesheet file is written any more.
	if matches, _ := filepath.Glob(buildtest.StaticPath(out, "/assets/site.*.css")); len(matches) != 0 {
		t.Errorf("found %v; want no external stylesheet in the output", matches)
	}
}

func TestStylesheet_DefinesBrandPaletteForLightDarkAndOverride(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	css := stylesheet(t, out, buildtest.Page(t, out, explainerPath))

	light := map[string]string{"--bg": "#ffffff", "--fg": "#111827", "--accent": "#0d9488"}
	dark := map[string]string{"--bg": "#0b0f19", "--fg": "#f3f4f6", "--accent": "#2dd4bf"}

	requireVars(t, "light :root", findRule(css, func(p string) bool {
		return p == ":root"
	}), light)
	requireVars(t, "@media (prefers-color-scheme: dark)", findRule(css, func(p string) bool {
		return strings.HasPrefix(p, "@media") && strings.Contains(p, "prefers-color-scheme:dark")
	}), dark)
	requireVars(t, "[data-theme=dark] override", findRule(css, func(p string) bool {
		return strings.Contains(p, "[data-theme=dark]") && !strings.HasPrefix(p, "@")
	}), dark)
}

func TestThemeToggle_IsTheOnlyScriptAndUnder1KB(t *testing.T) {
	preview := buildtest.Config()
	preview.Mode = builder.Preview
	ads := buildtest.Config()
	ads.AdsEnabled = true
	ads.AdSenseClient = "ca-pub-1234567890123456"
	ads.AdSenseSlots = map[string]string{"after-tldr": "1", "after-what-happened": "2", "after-body": "3"}

	for name, cfg := range map[string]builder.Config{"production": buildtest.Config(), "preview": preview, "ads": ads} {
		t.Run(name, func(t *testing.T) {
			out := buildtest.MustBuild(t, cfg, buildtest.ValidTree())
			requireOnlyAllowedScripts(t, buildtest.Page(t, out, explainerPath), cfg.AdsEnabled && cfg.Mode != builder.Preview, cfg.Mode != builder.Preview)
		})
	}
}

// requireOnlyAllowedScripts: the theme toggle is the one inline script of
// our own; the only other scripts allowed are Pagefind (search), in
// Production builds the Google Analytics 4 snippet and gtag.js (ADR-0008),
// and, only when ads are on, the AdSense library and its one-line
// per-unit push.
func requireOnlyAllowedScripts(t *testing.T, doc *goquery.Document, adsEnabled, analytics bool) {
	t.Helper()
	const adsensePush = "(adsbygoogle=window.adsbygoogle||[]).push({});"

	// Executable scripts only: JSON-LD blocks are data, not code.
	var inline []*goquery.Selection
	doc.Find("script").Each(func(_ int, s *goquery.Selection) {
		switch typ := s.AttrOr("type", ""); typ {
		case "application/ld+json":
			return
		case "", "text/javascript", "module":
			if s.Closest(".newsletter").Length() > 0 {
				// The newsletter form's enhancement (spec allows it); its
				// budget is checked in newsletter_test.go.
				return
			}
		default:
			t.Errorf("unexpected script type %q", typ)
		}
		if src, ok := s.Attr("src"); ok {
			switch {
			case analytics && strings.HasPrefix(src, "https://www.googletagmanager.com/gtag/js?id="):
			case strings.HasPrefix(src, "/pagefind/"):
			case adsEnabled && strings.HasPrefix(src, "https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js?client="):
			default:
				t.Errorf("unexpected external script %q", src)
			}
			return
		}
		if adsEnabled && strings.TrimSpace(s.Text()) == adsensePush {
			return
		}
		if analytics && strings.HasPrefix(s.Text(), "window.dataLayer=") && s.Next().Is(`script[src^="https://www.googletagmanager.com/"]`) {
			return
		}
		inline = append(inline, s)
	})
	if len(inline) != 1 {
		t.Fatalf("found %d inline scripts, want exactly 1 (the theme toggle)", len(inline))
	}
	script := inline[0]
	if n := len(script.Text()); n > 1024 {
		t.Errorf("theme script is %d bytes, want ≤ 1024", n)
	}
	if script.ParentsFiltered("head").Length() != 1 {
		t.Errorf("theme script must be in <head> so a stored override applies before first paint")
	}
	if !strings.Contains(script.Text(), "localStorage") || !strings.Contains(script.Text(), "data-theme-toggle") {
		t.Errorf("theme script should persist the override in localStorage and drive the [data-theme-toggle] button")
	}
}

func TestThemeToggle_IsHiddenUntilJavaScriptRuns(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	toggle := buildtest.Page(t, out, explainerPath).Find("button[data-theme-toggle]")

	if toggle.Length() != 1 {
		t.Fatalf("found %d theme toggles, want 1", toggle.Length())
	}
	if _, hidden := toggle.Attr("hidden"); !hidden {
		t.Errorf("theme toggle must carry `hidden` so readers without JS never see a dead control")
	}
	if toggle.AttrOr("aria-label", "") == "" {
		t.Errorf("theme toggle needs an accessible name")
	}
}

// Revealing the toggle must not move the header (CLS 0 budget, ticket 15:
// Lighthouse measured a 0.0008 shift on nav.site-nav when `hidden` meant
// display:none). The hidden toggle keeps its box but is invisible, so it is
// still no dead control for readers without JS.
func TestThemeToggle_RevealCausesNoLayoutShift(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	css := stylesheet(t, out, buildtest.Page(t, out, explainerPath))

	rule := strings.NewReplacer(" ", "", "\n", "", "\t", "").Replace(findRule(css, func(p string) bool {
		return p == ".theme-toggle[hidden]"
	}))
	if rule == "" {
		t.Fatalf("no .theme-toggle[hidden] rule: the UA's [hidden]{display:none} collapses the toggle, so revealing it shifts the nav")
	}
	if !strings.Contains(rule, "visibility:hidden") {
		t.Errorf(".theme-toggle[hidden] = {%s}, want visibility:hidden", rule)
	}
	if !strings.Contains(rule, "display:") || strings.Contains(rule, "display:none") {
		t.Errorf(".theme-toggle[hidden] = {%s}, want a display that keeps its box (not none)", rule)
	}
	// A global [hidden]{display:none !important} beats a plain declaration
	// whatever its specificity.
	global := strings.ReplaceAll(findRule(css, func(p string) bool { return p == "[hidden]" }), " ", "")
	if strings.Contains(global, "!important") && !regexp.MustCompile(`display:[a-z-]+!important`).MatchString(rule) {
		t.Errorf(".theme-toggle[hidden] = {%s}: the stylesheet's [hidden] rule is !important, so the toggle's display must be too", rule)
	}
}

// stylesheet returns the (lower-cased) text of the page's inline
// stylesheet, the single <style> in its <head> (ADR-0006).
func stylesheet(t *testing.T, _ string, doc *goquery.Document) string {
	t.Helper()
	style := doc.Find("head style")
	if style.Length() != 1 {
		t.Fatalf("found %d <style> elements in <head>, want 1", style.Length())
	}
	return strings.ToLower(style.Text())
}

// findRule returns the body of the first top-level CSS rule whose prelude
// (whitespace and quotes removed) satisfies match, or "" if none does.
func findRule(css string, match func(prelude string) bool) string {
	css = regexp.MustCompile(`(?s)/\*.*?\*/`).ReplaceAllString(css, "")
	norm := strings.NewReplacer(" ", "", "\n", "", "\t", "", `"`, "", "'", "")
	for i := 0; i < len(css); {
		open := strings.IndexByte(css[i:], '{')
		if open < 0 {
			break
		}
		open += i
		depth, end := 0, open
		for ; end < len(css); end++ {
			if css[end] == '{' {
				depth++
			} else if css[end] == '}' {
				if depth--; depth == 0 {
					break
				}
			}
		}
		if match(norm.Replace(css[i:open])) {
			return css[open+1 : end]
		}
		i = end + 1
	}
	return ""
}

func requireVars(t *testing.T, where, body string, want map[string]string) {
	t.Helper()
	body = strings.NewReplacer(" ", "", "\n", "", "\t", "").Replace(body)
	if body == "" {
		t.Errorf("%s: no such rule in the stylesheet", where)
		return
	}
	for name, value := range want {
		if !strings.Contains(body, name+":"+value) {
			t.Errorf("%s: want %s: %s", where, name, value)
		}
	}
}

func TestSiteFooter_LinksTransparencyPagesAndFeeds(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	footer := buildtest.Page(t, out, explainerPath).Find("body > footer").First()

	links := map[string]string{}
	footer.Find("a").Each(func(_ int, a *goquery.Selection) {
		links[a.AttrOr("href", "")] = strings.TrimSpace(a.Text())
	})
	for href, text := range map[string]string{
		"/about/":            "About",
		"/contact/":          "Contact",
		"/editorial-policy/": "Editorial Policy",
	} {
		if got, ok := links[href]; !ok || got != text {
			t.Errorf("footer link %s = %q (present %v), want %q", href, got, ok, text)
		}
	}
	for _, feed := range []string{"/feed.xml", "/briefs/feed.xml"} {
		a := footer.Find(`a[href="` + feed + `"]`)
		if a.Length() != 1 {
			t.Errorf("footer has %d links to %s, want 1", a.Length(), feed)
			continue
		}
		if got := a.AttrOr("type", ""); got != "application/rss+xml" {
			t.Errorf("feed link %s type = %q, want application/rss+xml", feed, got)
		}
	}
}

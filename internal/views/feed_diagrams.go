package views

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Feed readers load no site CSS, so a Mermaid SVG painted through CSS
// variables (var(--_node-fill), var(--_text)…) renders as black boxes with
// invisible edges there. FeedSafeDiagrams returns HTML whose diagram
// colours are resolved to the light palette. The resolution reads the
// .diagram / .diagram svg declarations in site.css and the palette tokens
// from brand.go, so pages and feeds share one source of truth; pages keep
// the variables (themable: light, dark, toggle).
func FeedSafeDiagrams(html string) string {
	return cssVarRef.ReplaceAllStringFunc(html, func(ref string) string {
		if v, ok := feedDiagramColours[cssVarRef.FindStringSubmatch(ref)[1]]; ok {
			return v
		}
		return ref
	})
}

// cssVarRef is a var() reference without a fallback, as the renderer
// writes them in SVG attributes.
var cssVarRef = regexp.MustCompile(`var\((--[a-zA-Z0-9_-]+)\)`)

// feedDiagramColours maps every diagram variable (and the palette tokens
// it builds on) to a concrete light-palette colour.
var feedDiagramColours = resolveDiagramVars(lightPalette, baseCSS)

func resolveDiagramVars(p palette, css string) map[string]string {
	vars := map[string]string{
		"--bg": p.bg, "--fg": p.fg, "--accent": p.accent, "--link": p.link,
		"--muted": p.muted, "--border": p.border, "--surface": p.surface,
	}
	for _, sel := range []string{".diagram", ".diagram svg"} {
		for _, d := range ruleDeclarations(css, sel) {
			name, expr, ok := strings.Cut(d, ":")
			name = strings.TrimSpace(name)
			if !ok || !strings.HasPrefix(name, "--") {
				continue
			}
			v, err := evalColour(strings.TrimSpace(expr), vars)
			if err != nil {
				panic(fmt.Sprintf("views: site.css %s %s: %v", sel, name, err))
			}
			vars[name] = v
		}
	}
	return vars
}

// ruleDeclarations are the declarations of the first one-line rule in css
// whose selector is exactly sel (site.css keeps one rule per line).
func ruleDeclarations(css, sel string) []string {
	for _, line := range strings.Split(css, "\n") {
		before, body, ok := strings.Cut(line, "{")
		if !ok || strings.TrimSpace(before) != sel {
			continue
		}
		body, _, _ = strings.Cut(body, "}")
		return splitTopLevel(body, ';')
	}
	panic("views: site.css has no " + sel + " rule")
}

// evalColour evaluates the colour expressions the diagram rules use:
// a hex colour, currentColor (the text colour, --fg), var(--x) or
// var(--x, fallback), and color-mix(in srgb, A p%, B).
func evalColour(expr string, vars map[string]string) (string, error) {
	switch {
	case strings.HasPrefix(expr, "#"):
		return expr, nil
	case expr == "currentColor":
		return vars["--fg"], nil
	case strings.HasPrefix(expr, "var(") && strings.HasSuffix(expr, ")"):
		args := splitTopLevel(expr[len("var("):len(expr)-1], ',')
		if v, ok := vars[strings.TrimSpace(args[0])]; ok {
			return v, nil
		}
		if len(args) < 2 {
			return "", fmt.Errorf("undefined %s", args[0])
		}
		return evalColour(strings.TrimSpace(strings.Join(args[1:], ",")), vars)
	case strings.HasPrefix(expr, "color-mix(in srgb,") && strings.HasSuffix(expr, ")"):
		args := splitTopLevel(expr[len("color-mix(in srgb,"):len(expr)-1], ',')
		if len(args) != 2 {
			return "", fmt.Errorf("color-mix wants two colours: %s", expr)
		}
		aExpr, pct, _ := strings.Cut(strings.TrimSpace(args[0]), " ")
		weight, err := strconv.ParseFloat(strings.TrimSuffix(strings.TrimSpace(pct), "%"), 64)
		if err != nil {
			return "", fmt.Errorf("color-mix percentage in %s: %v", expr, err)
		}
		a, err := evalColour(aExpr, vars)
		if err != nil {
			return "", err
		}
		b, err := evalColour(strings.TrimSpace(args[1]), vars)
		if err != nil {
			return "", err
		}
		return mixHex(a, b, weight/100)
	}
	return "", fmt.Errorf("unsupported colour expression %q", expr)
}

// mixHex is color-mix(in srgb, a w, b): per-channel a*w + b*(1-w).
func mixHex(a, b string, w float64) (string, error) {
	var ca, cb [3]int64
	for i := range 3 {
		var err error
		if ca[i], err = strconv.ParseInt(a[1+2*i:3+2*i], 16, 0); err != nil {
			return "", err
		}
		if cb[i], err = strconv.ParseInt(b[1+2*i:3+2*i], 16, 0); err != nil {
			return "", err
		}
	}
	out := "#"
	for i := range 3 {
		out += fmt.Sprintf("%02X", int(float64(ca[i])*w+float64(cb[i])*(1-w)+0.5))
	}
	return out, nil
}

// splitTopLevel splits s at sep outside parentheses.
func splitTopLevel(s string, sep rune) []string {
	var parts []string
	depth, start := 0, 0
	for i, r := range s {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
		case sep:
			if depth == 0 {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	if rest := strings.TrimSpace(s[start:]); rest != "" {
		parts = append(parts, s[start:])
	}
	return parts
}

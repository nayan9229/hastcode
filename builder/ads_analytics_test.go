package builder_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildtest"
)

const testAdSenseClient = "ca-pub-1234567890123456"

func adsOnConfig() builder.Config {
	cfg := buildtest.Config()
	cfg.AdsEnabled = true
	cfg.AdSenseClient = testAdSenseClient
	return cfg
}

func TestAdSlots_EnabledRenderTheAdSenseScriptAndTwoUnits(t *testing.T) {
	out := buildtest.MustBuild(t, adsOnConfig(), buildtest.ValidTree())
	doc := buildtest.Page(t, out, explainerPath)

	loaders := doc.Find(`script[src^="https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js"]`)
	if loaders.Length() != 1 {
		t.Fatalf("found %d AdSense library scripts, want exactly 1", loaders.Length())
	}
	if src := loaders.AttrOr("src", ""); !strings.Contains(src, "client="+testAdSenseClient) {
		t.Errorf("AdSense script src = %q, want client=%s", src, testAdSenseClient)
	}
	if _, async := loaders.Attr("async"); !async {
		t.Errorf("AdSense script must be async")
	}

	slots := doc.Find("aside.ad-slot")
	if slots.Length() != 2 {
		t.Fatalf("found %d ad slots, want 2", slots.Length())
	}
	slots.Each(func(i int, slot *goquery.Selection) {
		unit := slot.Find("ins.adsbygoogle")
		if unit.Length() != 1 {
			t.Errorf("slot %d holds %d AdSense units, want 1 (inside the fixed-height box)", i, unit.Length())
		}
		if got := unit.AttrOr("data-ad-client", ""); got != testAdSenseClient {
			t.Errorf("slot %d data-ad-client = %q, want %s", i, got, testAdSenseClient)
		}
		if got := slot.AttrOr("data-pagefind-ignore", ""); got != "all" {
			t.Errorf("slot %d data-pagefind-ignore = %q, want all", i, got)
		}
	})
}

func TestAdSlots_EnabledWithoutAdSenseClientFailsTheBuild(t *testing.T) {
	cfg := adsOnConfig()
	cfg.AdSenseClient = ""

	_, res := buildtest.Build(t, cfg, buildtest.ValidTree())

	buildtest.RequireErrorRule(t, res, "", "ads-config", "adsense_client")
}

func TestAdSlots_DisabledRenderTwoEmptyFixedHeightPlaceholders(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	doc := buildtest.Page(t, out, explainerPath)

	slots := doc.Find("aside.ad-slot")
	if slots.Length() != 2 {
		t.Fatalf("found %d ad slots, want 2 (after the TL;DR, after the body)", slots.Length())
	}
	slots.Each(func(i int, slot *goquery.Selection) {
		if got := slot.AttrOr("aria-hidden", ""); got != "true" {
			t.Errorf("slot %d aria-hidden = %q, want true", i, got)
		}
		if got := slot.AttrOr("data-pagefind-ignore", ""); got != "all" {
			t.Errorf("slot %d data-pagefind-ignore = %q, want all", i, got)
		}
		if html, _ := slot.Html(); strings.TrimSpace(html) != "" {
			t.Errorf("slot %d is not empty: %q", i, html)
		}
	})

	// Positions: the first follows the TL;DR list, the second follows the body.
	tldr := doc.Find(".e-content h2").FilterFunction(func(_ int, h *goquery.Selection) bool { return h.Text() == "TL;DR" })
	if !tldr.Next().Next().Is("aside.ad-slot") {
		t.Errorf("first ad slot does not directly follow the TL;DR list")
	}
	if !doc.Find(".e-content").Next().Is("aside.ad-slot") {
		t.Errorf("second ad slot does not directly follow the article body (.e-content)")
	}

	if h := adSlotHeight(t, stylesheet(t, out, doc)); !regexp.MustCompile(`^\d+px$`).MatchString(h) {
		t.Errorf(".ad-slot height = %q, want a fixed pixel height so enabling ads cannot shift layout", h)
	}

	if n := doc.Find("ins.adsbygoogle").Length(); n != 0 {
		t.Errorf("found %d AdSense units with ads off, want 0", n)
	}
	for _, third := range []string{"googlesyndication.com", "plausible.io"} {
		if hits := buildtest.FilesContaining(t, out, third); len(hits) != 0 {
			t.Errorf("%q appears in %v with ads and analytics off", third, hits)
		}
	}
}

func TestAnalytics_PlausibleScriptWhenDomainConfigured(t *testing.T) {
	cfg := buildtest.Config()
	cfg.PlausibleDomain = "hastcode.com"
	out := buildtest.MustBuild(t, cfg, buildtest.ValidTree())
	doc := buildtest.Page(t, out, explainerPath)

	s := doc.Find(`script[src="https://plausible.io/js/script.js"]`)
	if s.Length() != 1 {
		t.Fatalf("found %d Plausible scripts, want 1", s.Length())
	}
	if got := s.AttrOr("data-domain", ""); got != "hastcode.com" {
		t.Errorf("Plausible data-domain = %q, want hastcode.com", got)
	}
	if _, deferred := s.Attr("defer"); !deferred {
		t.Errorf("Plausible script must be defer")
	}
	// Plausible's standard script sets no cookies; nothing else may either.
	if hits := buildtest.FilesContaining(t, out, "cookie"); len(hits) != 0 {
		t.Errorf("output mentions cookies in %v; analytics must be cookieless", hits)
	}
}

func TestAnalytics_OmittedWhenDomainBlank(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())

	if hits := buildtest.FilesContaining(t, out, "plausible"); len(hits) != 0 {
		t.Errorf("Plausible appears in %v with no domain configured", hits)
	}
}

// adSlotHeight returns the .ad-slot rule's height, resolving one var()
// token, from the lower-cased stylesheet text.
func adSlotHeight(t *testing.T, css string) string {
	t.Helper()
	body := findRule(css, func(p string) bool { return p == ".ad-slot" })
	h := declaration(body, "height")
	if m := regexp.MustCompile(`^var\((--[a-z0-9-]+)\)$`).FindStringSubmatch(h); m != nil {
		h = declaration(allRootVars(css), m[1])
	}
	return h
}

// declaration returns the value of property prop in a CSS rule body.
func declaration(body, prop string) string {
	for _, d := range strings.Split(body, ";") {
		name, value, ok := strings.Cut(d, ":")
		if ok && strings.TrimSpace(name) == prop {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

// allRootVars concatenates the bodies of every top-level :root rule.
func allRootVars(css string) string {
	var b strings.Builder
	for _, m := range regexp.MustCompile(`(?m)^:root\s*\{([^}]*)\}`).FindAllStringSubmatch(css, -1) {
		b.WriteString(m[1] + ";")
	}
	return b.String()
}

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

// testAdSlots are the AdSense ad unit ids per slot position.
var testAdSlots = map[string]string{
	"after-tldr":          "1111111111",
	"after-what-happened": "2222222222",
	"after-body":          "3333333333",
}

func adsOnConfig() builder.Config {
	cfg := buildtest.Config()
	cfg.AdsEnabled = true
	cfg.AdSenseClient = testAdSenseClient
	cfg.AdSenseSlots = testAdSlots
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
		// An AdSense display unit fills only with its ad unit id.
		position := slot.AttrOr("data-ad-position", "")
		if got, want := unit.AttrOr("data-ad-slot", ""), testAdSlots[position]; want == "" || got != want {
			t.Errorf("slot %d (%s) data-ad-slot = %q, want %q", i, position, got, want)
		}
		if got := unit.AttrOr("data-ad-format", ""); got != "auto" {
			t.Errorf("slot %d data-ad-format = %q, want auto (responsive inside the fixed-height box)", i, got)
		}
		if got := slot.AttrOr("data-pagefind-ignore", ""); got != "all" {
			t.Errorf("slot %d data-pagefind-ignore = %q, want all", i, got)
		}
	})
}

// Without an ad unit id a slot would render an empty box on every page,
// so each position needs one before ads can go on.
func TestAdSlots_EnabledWithoutAnAdUnitIDFailsTheBuildNamingThePosition(t *testing.T) {
	cfg := adsOnConfig()
	cfg.AdSenseSlots = map[string]string{"after-tldr": "1111111111"}

	_, res := buildtest.Build(t, cfg, buildtest.ValidTree())

	buildtest.RequireErrorRule(t, res, "", "ads-config", "after-what-happened")
	buildtest.RequireErrorRule(t, res, "", "ads-config", "after-body")
}

func TestAdSlots_AdUnitIDMustBeDigits(t *testing.T) {
	cfg := adsOnConfig()
	cfg.AdSenseSlots = map[string]string{"after-tldr": "ca-pub-1", "after-what-happened": "2222222222", "after-body": "3333333333"}

	_, res := buildtest.Build(t, cfg, buildtest.ValidTree())

	buildtest.RequireErrorRule(t, res, "", "ads-config", "after-tldr")
}

// Previews (and make dev) never carry ads: their impressions would be the
// Author's own, and Lighthouse CI measures Previews (ADR-0006).
func TestAdSlots_OmittedFromPreviewBuilds(t *testing.T) {
	cfg := adsOnConfig()
	cfg.Mode = builder.Preview
	out := buildtest.MustBuild(t, cfg, buildtest.ValidTree())

	if hits := buildtest.FilesContaining(t, out, "googlesyndication.com"); len(hits) != 0 {
		t.Errorf("AdSense appears in %v in a preview build", hits)
	}
	if n := buildtest.Page(t, out, explainerPath).Find(".ad-slot").Length(); n != 0 {
		t.Errorf("found %d ad slots in a preview build, want none", n)
	}
}

func TestAdSlots_EnabledWithoutAdSenseClientFailsTheBuild(t *testing.T) {
	cfg := adsOnConfig()
	cfg.AdSenseClient = ""

	_, res := buildtest.Build(t, cfg, buildtest.ValidTree())

	buildtest.RequireErrorRule(t, res, "", "ads-config", "adsense_client")
}

func TestAdSlots_EnabledSitInFixedHeightBoxesAfterTLDRAndAfterBody(t *testing.T) {
	out := buildtest.MustBuild(t, adsOnConfig(), buildtest.ValidTree())
	doc := buildtest.Page(t, out, explainerPath)

	tldr := doc.Find(".e-content h2").FilterFunction(func(_ int, h *goquery.Selection) bool { return h.Text() == "TL;DR" })
	first := tldr.Next().Next()
	if !first.Is("aside.ad-slot") || first.AttrOr("data-ad-position", "") != "after-tldr" {
		t.Errorf("first ad slot does not directly follow the TL;DR list with data-ad-position=after-tldr")
	}
	second := doc.Find(".e-content").Next()
	if !second.Is("aside.ad-slot") || second.AttrOr("data-ad-position", "") != "after-body" {
		t.Errorf("second ad slot does not directly follow the article body (.e-content) with data-ad-position=after-body")
	}

	if h := adSlotHeight(t, stylesheet(t, out, doc)); !regexp.MustCompile(`^\d+px$`).MatchString(h) {
		t.Errorf(".ad-slot height = %q, want a fixed pixel height so an AdSense unit filling in cannot shift layout", h)
	}
}

func TestAdSlots_EnabledOnABriefTheFirstFollowsWhatHappened(t *testing.T) {
	tree := buildtest.ValidTree().WithBrief(buildtest.ValidBrief())
	out := buildtest.MustBuild(t, adsOnConfig(), tree)
	doc := buildtest.Page(t, out, "/briefs/etcd-raft-rewrite/")

	slots := doc.Find("aside.ad-slot")
	if slots.Length() != 2 {
		t.Fatalf("found %d ad slots on the Brief, want 2", slots.Length())
	}
	first := slots.First()
	if got := first.AttrOr("data-ad-position", ""); got != "after-what-happened" {
		t.Errorf("first Brief ad slot data-ad-position = %q, want after-what-happened (a Brief has no TL;DR)", got)
	}
	if !first.PrevAllFiltered("h2").First().Is("h2") || first.PrevAllFiltered("h2").First().Text() != "What happened" {
		t.Errorf("first Brief ad slot does not follow the What happened section")
	}
	if !first.NextAllFiltered("h2").First().Is("h2") || first.NextAllFiltered("h2").First().Text() != "Why it matters" {
		t.Errorf("first Brief ad slot is not before the Why it matters section")
	}
}

func TestAdSlots_DisabledRenderNoSlotAndReserveNoSpace(t *testing.T) {
	tree := buildtest.ValidTree().WithBrief(buildtest.ValidBrief())
	out := buildtest.MustBuild(t, buildtest.Config(), tree)

	for _, page := range []string{explainerPath, "/briefs/etcd-raft-rewrite/"} {
		doc := buildtest.Page(t, out, page)
		if n := doc.Find(".ad-slot, [data-ad-position]").Length(); n != 0 {
			t.Errorf("%s: found %d ad-slot elements with ads off, want none (no reserved blank boxes)", page, n)
		}
		if n := doc.Find("ins.adsbygoogle").Length(); n != 0 {
			t.Errorf("%s: found %d AdSense units with ads off, want 0", page, n)
		}
	}
	if hits := buildtest.FilesContaining(t, out, "googlesyndication.com"); len(hits) != 0 {
		t.Errorf("AdSense appears in %v with ads off", hits)
	}
}

// gtagSrc is gtag.js as every fixture Production page loads it.
const gtagSrc = "https://www.googletagmanager.com/gtag/js?id=" + buildtest.GAMeasurementID

// Every Production page reports to Google Analytics 4 (ADR-0008): an inline
// snippet that queues consent defaults and the config, then gtag.js
// deferred so it runs after parsing.
func TestAnalytics_ProductionPagesCarryGA4(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())

	for _, page := range []string{"/", explainerPath, "/search/", "/privacy/"} {
		doc := buildtest.Page(t, out, page)
		head := doc.Find("head")
		loaders := head.Find(`script[src^="https://www.googletagmanager.com/"]`)
		if loaders.Length() != 1 {
			t.Fatalf("%s: found %d gtag.js scripts in <head>, want 1", page, loaders.Length())
		}
		if got := loaders.AttrOr("src", ""); got != gtagSrc {
			t.Errorf("%s: gtag.js src = %q, want %q", page, got, gtagSrc)
		}
		if _, deferred := loaders.Attr("defer"); !deferred {
			t.Errorf("%s: gtag.js must be defer", page)
		}
		snippet := loaders.Prev()
		if !snippet.Is("script") || snippet.AttrOr("src", "") != "" {
			t.Fatalf("%s: gtag.js is not preceded by the inline analytics snippet", page)
		}
		code := snippet.Text()
		for _, want := range []string{
			`window.dataLayer=window.dataLayer||[]`,
			`if(location.hostname!=="hastcode.test"){window["ga-disable-"+"` + buildtest.GAMeasurementID + `"]=true;}`,
			`gtag("consent","default",{analytics_storage:"denied",wait_for_update:500,region:[`,
			`"DE"`, `"NO"`, `"GB"`, `"CH"`,
			`gtag("consent","default",{analytics_storage:"granted"});`,
			`gtag("config","` + buildtest.GAMeasurementID + `")`,
		} {
			if !strings.Contains(code, want) {
				t.Errorf("%s: analytics snippet lacks %s:\n%s", page, want, code)
			}
		}
		// Ad consent belongs to AdSense's own consent message (ADR-0008).
		if strings.Contains(code, "ad_storage") || strings.Contains(code, "ad_personalization") {
			t.Errorf("%s: analytics snippet sets ad consent; it may set analytics_storage only", page)
		}
		// Consent defaults only apply when set before the config command.
		if strings.LastIndex(code, `gtag("consent","default"`) > strings.Index(code, `gtag("config"`) {
			t.Errorf("%s: consent defaults must precede the config", page)
		}
	}
}

// Previews (and make dev) never carry analytics: their visits are not
// Readers', and Lighthouse CI, which only ever tests Previews, then measures
// the page's own code (ADR-0006).
func TestAnalytics_OmittedFromPreviewBuilds(t *testing.T) {
	cfg := buildtest.Config()
	cfg.Mode = builder.Preview
	out := buildtest.MustBuild(t, cfg, buildtest.ValidTree())

	for _, needle := range []string{"googletagmanager", buildtest.GAMeasurementID, "dataLayer"} {
		if hits := buildtest.FilesContaining(t, out, needle); len(hits) != 0 {
			t.Errorf("%q appears in %v in a preview build", needle, hits)
		}
	}
}

// A preview build needs no Measurement ID (local builds, branch Previews).
func TestAnalytics_PreviewBuildPassesWithoutMeasurementID(t *testing.T) {
	cfg := buildtest.Config()
	cfg.Mode = builder.Preview
	cfg.GAMeasurementID = ""
	buildtest.MustBuild(t, cfg, buildtest.ValidTree())
}

// A Production build (Production and QA Targets) without a Measurement ID
// would silently record nothing, so it fails (ADR-0008).
func TestAnalytics_ProductionBuildWithoutMeasurementIDFails(t *testing.T) {
	cfg := buildtest.Config()
	cfg.GAMeasurementID = ""

	_, res := buildtest.Build(t, cfg, buildtest.ValidTree())

	buildtest.RequireErrorRule(t, res, "", "analytics-config", "ga_measurement_id is blank")
}

func TestAnalytics_MeasurementIDMustBeGA4(t *testing.T) {
	for _, id := range []string{"UA-12345678-1", "g-abc12345", "G-", `G-ABC"};alert(1)//`} {
		t.Run(id, func(t *testing.T) {
			cfg := buildtest.Config()
			cfg.GAMeasurementID = id

			_, res := buildtest.Build(t, cfg, buildtest.ValidTree())

			buildtest.RequireErrorRule(t, res, "", "analytics-config", "not a GA4 Measurement ID")
		})
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

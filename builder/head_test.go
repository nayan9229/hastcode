package builder_test

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/internal/buildtest"
)

func TestExplainerHead_CanonicalOpenGraphAndTwitterCardUseTheHero(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	doc := buildtest.Page(t, out, raftURL)
	meta := func(attr, name string) string {
		return doc.Find(`head meta[`+attr+`="`+name+`"]`).AttrOr("content", "")
	}

	if got := doc.Find(`head link[rel="canonical"]`).AttrOr("href", ""); got != "https://hastcode.test"+raftURL {
		t.Errorf("canonical = %q, want https://hastcode.test%s (absolute, trailing slash)", got, raftURL)
	}
	for name, want := range map[string]string{
		"og:type":        "article",
		"og:url":         "https://hastcode.test" + raftURL,
		"og:title":       "How Raft elects a leader",
		"og:description": "Raft keeps a replicated log consistent by electing one leader per term with randomised timeouts.",
		"og:site_name":   "HastCode",
		"og:image:alt":   "Five Raft servers voting for a leader",
	} {
		if got := meta("property", name); got != want {
			t.Errorf("%s = %q, want %q", name, got, want)
		}
	}
	if got := meta("name", "twitter:card"); got != "summary_large_image" {
		t.Errorf("twitter:card = %q, want summary_large_image", got)
	}

	image := meta("property", "og:image")
	path, ok := strings.CutPrefix(image, "https://hastcode.test/")
	if !ok {
		t.Fatalf("og:image = %q, want an absolute URL on the site", image)
	}
	img := buildtest.DecodeImage(t, out, "/"+path)
	if img.Width < 1200 || img.Width*9 != img.Height*16 {
		t.Errorf("og:image is %dx%d, want the ≥1200 px 16:9 hero crop", img.Width, img.Height)
	}
	if meta("property", "og:image:width") != "1200" || meta("property", "og:image:height") != "675" {
		t.Errorf("og:image:width/height = %q/%q, want 1200/675", meta("property", "og:image:width"), meta("property", "og:image:height"))
	}
}

// Spec: og:image at least 1200 px wide on every page, and a generated OG
// badge as the default. Pages without a hero share one 1200x630 PNG under
// the immutable /assets/ path.
func TestPagesWithoutAHero_UseTheGeneratedOGBadgeAsALargeImageCard(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree().WithBrief(buildtest.ValidBrief()))

	var badge string
	for _, page := range []string{"/", "/distributed-systems/", "/briefs/", "/about/", "/authors/ada/", "/search/", "/explainers/concept/"} {
		doc := buildtest.Page(t, out, page)
		meta := func(attr, name string) string {
			return doc.Find(`head meta[`+attr+`="`+name+`"]`).AttrOr("content", "")
		}
		image := meta("property", "og:image")
		if badge == "" {
			badge = image
		}
		if image == "" || image != badge {
			t.Errorf("%s og:image = %q, want the shared badge %q", page, image, badge)
		}
		if got := meta("name", "twitter:card"); got != "summary_large_image" {
			t.Errorf("%s twitter:card = %q, want summary_large_image", page, got)
		}
		if meta("property", "og:image:width") != "1200" || meta("property", "og:image:height") != "630" {
			t.Errorf("%s og:image:width/height = %q/%q, want 1200/630", page, meta("property", "og:image:width"), meta("property", "og:image:height"))
		}
		if meta("property", "og:image:alt") == "" {
			t.Errorf("%s has no og:image:alt", page)
		}
	}

	path, ok := strings.CutPrefix(badge, "https://hastcode.test/assets/")
	if !ok || !strings.HasSuffix(path, ".png") {
		t.Fatalf("badge URL = %q, want an absolute .png under /assets/ (content-hashed, cached immutably)", badge)
	}
	img := buildtest.DecodeImage(t, out, "/assets/"+path)
	if img.Format != "png" || img.Width != 1200 || img.Height != 630 {
		t.Errorf("badge is a %dx%d %s, want a 1200x630 png", img.Width, img.Height, img.Format)
	}

	// Article pages keep their hero.
	if got := buildtest.Page(t, out, raftURL).Find(`head meta[property="og:image"]`).AttrOr("content", ""); got == badge {
		t.Errorf("the Explainer uses the badge, want its hero crop")
	}
}

func TestExplainerHead_RobotsAllowLargeImagePreviews(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	doc := buildtest.Page(t, out, raftURL)

	robots := doc.Find(`meta[name="robots"]`)
	if robots.Length() != 1 {
		t.Fatalf("found %d robots meta tags, want 1", robots.Length())
	}
	if got := robots.AttrOr("content", ""); !strings.Contains(got, "max-image-preview:large") || strings.Contains(got, "noindex") {
		t.Errorf("robots = %q, want max-image-preview:large and no noindex", got)
	}
	if got := doc.Find("html").AttrOr("lang", ""); got != "en" {
		t.Errorf(`html lang = %q, want "en"`, got)
	}
	for _, scheme := range []string{"light", "dark"} {
		sel := `meta[name="theme-color"][media="(prefers-color-scheme: ` + scheme + `)"]`
		if doc.Find(sel).AttrOr("content", "") == "" {
			t.Errorf("no theme-color for the %s scheme", scheme)
		}
	}
}

// The only production pages kept out of search are the newsletter's
// "check your inbox" page, the 404 page and the site search page; every
// content page (articles, Author pages, static pages) is indexable. Add a
// new noindex page here deliberately.
func TestProductionOutput_NothingIsNoindexButTheSubscribed404AndSearchPages(t *testing.T) {
	tree := buildtest.ValidTree().WithExplainer(buildtest.ValidComparison())
	for _, p := range buildtest.TransparencyPages() {
		tree.WithPage(p)
	}
	out := buildtest.MustBuild(t, buildtest.Config(), tree)

	noindex := map[string]bool{
		"static/subscribed/index.html": true, "static/404.html": true, "static/search/index.html": true,
		// Cloudflare's header rules: noindex only on *.workers.dev hosts,
		// never hastcode.com (TestCloudflare_WorkersDevHostnamesAreNoindexAndHastcodeIsNot).
		"static/_headers": true,
	}
	for _, file := range buildtest.FilesContaining(t, out, "noindex") {
		if !noindex[file] {
			t.Errorf("%s contains noindex", file)
		}
	}
}

func TestExplainerPage_ShowsTheBreadcrumbTrail(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	nav := buildtest.Page(t, out, raftURL).Find(`nav[aria-label="Breadcrumb"]`)
	if nav.Length() != 1 {
		t.Fatalf("found %d breadcrumb navs, want 1", nav.Length())
	}

	var got []string
	nav.Find("li").Each(func(_ int, li *goquery.Selection) {
		href, _ := li.Find("a").Attr("href")
		got = append(got, strings.TrimSpace(li.Text())+" "+href)
	})
	want := []string{"Home /", "Distributed Systems /distributed-systems/", "How Raft elects a leader "}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("breadcrumb = %q, want %q (last item is the current page, unlinked)", got, want)
	}
	if cur := nav.Find(`[aria-current="page"]`).Text(); strings.TrimSpace(cur) != "How Raft elects a leader" {
		t.Errorf(`aria-current="page" = %q, want the Explainer title`, cur)
	}
}

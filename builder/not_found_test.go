package builder_test

import (
	"strings"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

// notFoundURL is the static 404 page Vercel's fallback route serves for
// unknown paths.
const notFoundURL = "/404.html"

func TestNotFoundPage_IsNoindexAndLinksHomeVerticalsAndSearch(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree())
	if !buildtest.Exists(out, notFoundURL) {
		t.Fatalf("no %s in output", notFoundURL)
	}
	doc := buildtest.Page(t, out, notFoundURL)

	if got := doc.Find(`meta[name="robots"]`).AttrOr("content", ""); !strings.Contains(got, "noindex") {
		t.Errorf("meta robots = %q, want noindex", got)
	}
	if n := doc.Find(`link[rel="canonical"]`).Length(); n != 0 {
		t.Errorf("found %d canonical links, want none on the 404 page", n)
	}

	links := map[string]bool{}
	for _, n := range doc.Find("main a").Nodes {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				links[attr.Val] = true
			}
		}
	}
	want := []string{"/", "/search/"}
	for _, v := range buildtest.Config().Verticals {
		want = append(want, "/"+v.Slug+"/")
	}
	for _, href := range want {
		if !links[href] {
			t.Errorf("404 page main content does not link %s", href)
		}
	}
}

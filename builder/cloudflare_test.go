package builder_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

// The Cloudflare manifest (ticket 05): the Production Target is Workers
// Static Assets, configured by <out>/wrangler.json plus _headers and
// _redirects inside the assets directory.

func TestCloudflare_WorkersConfigServesTheStaticTreeWithTheWorkerOnlyOnTheSignupPath(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree())
	wc := buildtest.WorkersConfig(t, out)

	if wc["name"] != "hastcode" {
		t.Errorf("name = %v, want hastcode", wc["name"])
	}
	if d, _ := wc["compatibility_date"].(string); len(d) != len("2006-01-02") {
		t.Errorf("compatibility_date = %v, want a yyyy-mm-dd date", wc["compatibility_date"])
	}
	// Relative paths in a Wrangler config resolve against the config
	// file's directory (wrangler deploy --config <out>/wrangler.json).
	main, _ := wc["main"].(string)
	if got, want := filepath.Join(out, main), filepath.Join(buildtest.RepoRoot(), "worker", "index.js"); got != want {
		t.Errorf("main = %q resolves to %s, want the Worker module %s", main, got, want)
	}

	assets, _ := wc["assets"].(map[string]any)
	dir, _ := assets["directory"].(string)
	if got := filepath.Join(out, dir); got != filepath.Join(out, "static") {
		t.Errorf("assets.directory = %q resolves to %s, want the static tree", dir, got)
	}
	for key, want := range map[string]any{
		"binding":            "ASSETS",
		"html_handling":      "force-trailing-slash",
		"not_found_handling": "404-page",
		"run_worker_first":   []any{"/api/subscribe"},
	} {
		if !reflect.DeepEqual(assets[key], want) {
			t.Errorf("assets.%s = %#v, want %#v", key, assets[key], want)
		}
	}
	// not_found_handling "404-page" serves the nearest 404.html with
	// status 404: the one at the assets root covers every unknown path.
	if _, err := os.Stat(filepath.Join(out, "static", "404.html")); err != nil {
		t.Errorf("no 404.html at the assets root: %v", err)
	}
}

func TestCloudflare_HeadersCacheHashedAssetsForeverAndRevalidateHTML(t *testing.T) {
	tree := buildtest.ValidTree()
	for _, p := range buildtest.TransparencyPages() {
		tree.WithPage(p)
	}
	out := buildtest.MustBuild(t, buildtest.Config(), tree)

	var assets, pages int
	for _, path := range buildtest.URLPaths(t, out) {
		cc := buildtest.CloudflareHeaders(t, out, path).Get("Cache-Control")
		switch {
		case strings.HasPrefix(path, "/assets/"):
			assets++
			if cc != "public, max-age=31536000, immutable" {
				t.Errorf("%s Cache-Control = %q, want public, max-age=31536000, immutable", path, cc)
			}
		case strings.HasSuffix(path, "/") || strings.HasSuffix(path, ".html"):
			pages++
			if cc != "public, max-age=0, must-revalidate" {
				t.Errorf("%s Cache-Control = %q, want public, max-age=0, must-revalidate", path, cc)
			}
		case strings.Contains(cc, "immutable"):
			t.Errorf("%s Cache-Control = %q; only content-hashed /assets/ may be immutable", path, cc)
		}
	}
	if assets < 10 || pages < 10 {
		t.Errorf("checked %d /assets/ files and %d pages; want the stylesheet, hero variants and every page", assets, pages)
	}
	// An unknown path is answered with the 404 page; it revalidates too.
	if cc := buildtest.CloudflareHeaders(t, out, "/no-such-page/").Get("Cache-Control"); cc != "public, max-age=0, must-revalidate" {
		t.Errorf("/no-such-page/ Cache-Control = %q, want revalidation", cc)
	}
}

func TestCloudflare_RedirectsAreExactlyTheHomeAlias(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree())

	got := buildtest.CloudflareRedirects(t, out)
	want := []string{"/home / 301", "/home/ / 301"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("_redirects rules = %q, want %q", got, want)
	}
}

// ADR-0001 as amended by ADR-0004: on the Production Target the retired
// community sections are plain 404s. No redirect, header rule or Worker
// route mentions them, so they fall through to not_found_handling.
func TestCloudflare_RetiredCommunityPathsHaveNoRuleAndFallThroughToNotFound(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree())

	for _, file := range []string{filepath.Join(out, "wrangler.json"), filepath.Join(out, "static", "_headers"), filepath.Join(out, "static", "_redirects")} {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		for _, section := range []string{"video", "community", "group"} {
			if strings.Contains(string(data), "/"+section) {
				t.Errorf("%s mentions /%s; retired paths must fall through to the 404 page", filepath.Base(file), section)
			}
		}
	}
	for _, path := range []string{"/video/", "/community/", "/group/"} {
		if buildtest.Exists(out, path) {
			t.Errorf("an asset exists at retired path %s", path)
		}
	}
	assets, _ := buildtest.WorkersConfig(t, out)["assets"].(map[string]any)
	if assets["not_found_handling"] != "404-page" {
		t.Errorf("not_found_handling = %v, want 404-page so unknown paths get /404.html with status 404", assets["not_found_handling"])
	}
	if page := buildtest.Page(t, out, "/404.html"); page.Find(`meta[name="robots"][content*="noindex"]`).Length() == 0 {
		t.Error("/404.html is not the noindex 404 page")
	}
}

func TestCloudflare_WorkersDevHostnamesAreNoindexAndHastcodeIsNot(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree())

	for _, u := range []string{
		"https://hastcode.someaccount.workers.dev/",
		"https://hastcode.someaccount.workers.dev/about/",
		"https://my-branch-hastcode.someaccount.workers.dev/about/", // a Preview
	} {
		if got := buildtest.CloudflareHeaders(t, out, u).Get("X-Robots-Tag"); got != "noindex" {
			t.Errorf("%s X-Robots-Tag = %q, want noindex", u, got)
		}
	}
	for _, u := range []string{"https://hastcode.com/", "https://hastcode.com/about/"} {
		if got := buildtest.CloudflareHeaders(t, out, u).Get("X-Robots-Tag"); got != "" {
			t.Errorf("%s X-Robots-Tag = %q, want none on the production host", u, got)
		}
	}
}

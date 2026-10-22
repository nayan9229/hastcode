package builder_test

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildtest"
)

// The QA Target (ticket 07, ADR-0004): Vercel only ever serves QA, so
// nothing it serves may be indexed or mistaken for Production.

// Every response on Vercel, whatever its status, carries
// X-Robots-Tag: noindex. The local router applies config.json the way
// Vercel does, so `make dev` sends the header too.
func TestQATarget_EveryResponseIsNoindex(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	srv := buildtest.Serve(t, out)

	for path, status := range map[string]int{
		"/":                    http.StatusOK,
		raftURL:                http.StatusOK,
		"/robots.txt":          http.StatusOK,
		"/sitemap.xml":         http.StatusOK,
		"/no-such-page/":       http.StatusNotFound,
		"/video/retired":       http.StatusGone,
		"/home":                http.StatusMovedPermanently,
		raftURL + "index.html": http.StatusPermanentRedirect,
		"/api/subscribe":       http.StatusNoContent,
	} {
		res := srv.Get(t, path)
		if res.StatusCode != status {
			t.Errorf("GET %s = %d, want %d", path, res.StatusCode, status)
		}
		if got := res.Header.Get("X-Robots-Tag"); got != "noindex" {
			t.Errorf("GET %s: X-Robots-Tag = %q, want noindex", path, got)
		}
	}
}

// A QA build is the Production build plus the Vercel adapter: site_url is
// the one in site.yaml, so canonical, Open Graph and JSON-LD URLs on QA
// point at https://hastcode.com, never at the QA host.
func TestQATarget_CanonicalAndOpenGraphURLsAreProductions(t *testing.T) {
	site, err := builder.LoadConfig(filepath.Join(buildtest.RepoRoot(), "site.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	cfg := buildtest.Config()
	cfg.SiteURL = site.SiteURL
	out := buildtest.MustBuild(t, cfg, buildtest.ValidTree())
	buildtest.Adapt(t, out, "vercel")

	const want = "https://hastcode.com" + raftURL
	doc := buildtest.Page(t, out, raftURL)
	if got := doc.Find(`head link[rel="canonical"]`).AttrOr("href", ""); got != want {
		t.Errorf("canonical = %q, want %q", got, want)
	}
	if got := doc.Find(`head meta[property="og:url"]`).AttrOr("content", ""); got != want {
		t.Errorf("og:url = %q, want %q", got, want)
	}
	page, _ := buildtest.JSONLDOfType(t, doc, "Article")["mainEntityOfPage"].(map[string]any)
	if got := page["@id"]; got != want {
		t.Errorf("JSON-LD Article mainEntityOfPage @id = %v, want %q", got, want)
	}
}

// The Vercel adapter overlays a robots.txt that disallows everything in
// place of the real one; the Cloudflare adapter keeps the real one.
func TestQATarget_VercelAdapterOverlaysADisallowAllRobotsTxt(t *testing.T) {
	const disallowAll = "User-agent: *\nDisallow: /\n"
	for target, wantQA := range map[string]bool{"vercel": true, "cloudflare": false} {
		t.Run(target, func(t *testing.T) {
			out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
			buildtest.Adapt(t, out, target)
			data, err := os.ReadFile(buildtest.StaticPath(out, "/robots.txt"))
			if err != nil {
				t.Fatal(err)
			}
			got := string(data)
			switch {
			case wantQA && got != disallowAll:
				t.Errorf("Vercel robots.txt =\n%s\nwant\n%s", got, disallowAll)
			case !wantQA && (strings.Contains(got, "Disallow") || !strings.Contains(got, "Sitemap: https://hastcode.test/sitemap.xml")):
				t.Errorf("Cloudflare robots.txt is not the real one:\n%s", got)
			}
		})
	}
}

package builder_test

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

// Routes are asserted by serving the build output through the local router,
// which applies config.json's routes the way Vercel does (ticket 10).

func TestRoutes_RetiredCommunityURLsAreGone(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree())
	srv := buildtest.Serve(t, out)

	for _, path := range []string{
		"/video/8f2a-some-clip", "/video/8f2a-some-clip/", "/video", "/video/",
		"/community/golang", "/community/", "/group/abc/members",
		"/Video/UPPER-case",
	} {
		res := srv.Get(t, path)
		if res.StatusCode != http.StatusGone {
			t.Errorf("GET %s = %d, want 410", path, res.StatusCode)
		}
	}
	for _, path := range []string{"/videos/", "/groups/", "/communitynotes/"} {
		if res := srv.Get(t, path); res.StatusCode == http.StatusGone {
			t.Errorf("GET %s = 410; only /video, /community and /group are retired", path)
		}
	}
}

func TestRoutes_EveryOutputFileHasExactlyOneURL(t *testing.T) {
	tree := buildtest.ValidTree()
	for _, p := range buildtest.TransparencyPages() {
		tree.WithPage(p)
	}
	out := buildtest.MustBuild(t, buildtest.Config(), tree)
	srv := buildtest.Serve(t, out)

	paths := buildtest.URLPaths(t, out)
	if len(paths) < 10 {
		t.Fatalf("only %d output URLs: %v", len(paths), paths)
	}
	for _, path := range paths {
		if res := srv.Get(t, path); res.StatusCode != http.StatusOK {
			t.Errorf("GET %s = %d, want 200", path, res.StatusCode)
		}
		var others map[string]string // other spelling -> want Location
		if dir, ok := strings.CutSuffix(path, "/"); ok {
			others = map[string]string{path + "index.html": path}
			if dir != "" {
				others[dir] = path
			}
		} else {
			others = map[string]string{path + "/": path}
		}
		for other, want := range others {
			res := srv.Get(t, other)
			if res.StatusCode != http.StatusPermanentRedirect || res.Header.Get("Location") != want {
				t.Errorf("GET %s = %d Location %q, want 308 to %s", other, res.StatusCode, res.Header.Get("Location"), want)
			}
		}
	}
}

func TestRoutes_UnknownPathsGetTheNotFoundPageWithStatus404(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree())
	srv := buildtest.Serve(t, out)

	for _, path := range []string{"/no-such-page/", "/distributed-systems/no-such-slug/", "/no-such-file.xml", "/assets/missing.0123456789.css"} {
		res := srv.Get(t, path)
		if res.StatusCode != http.StatusNotFound {
			t.Errorf("GET %s = %d, want 404", path, res.StatusCode)
		}
		if !strings.Contains(res.Body, "Page not found") || !strings.Contains(res.Body, "noindex") {
			t.Errorf("GET %s body is not the 404 page:\n%.300s", path, res.Body)
		}
	}
}

func TestRoutes_SubscribeFunctionIsNeverRedirectedOrShadowed(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree())
	srv := buildtest.Serve(t, out)

	res := srv.Do(t, http.MethodPost, "/api/subscribe", "email=reader%40example.com&source=%2F")
	if res.StatusCode != http.StatusNoContent || srv.FunctionCalls != 1 {
		t.Errorf("POST /api/subscribe = %d (Location %q), function calls %d; want the function to answer",
			res.StatusCode, res.Header.Get("Location"), srv.FunctionCalls)
	}
}

func TestRoutes_HTMLRevalidatesInBrowsersAndIsCachedLongOnTheCDN(t *testing.T) {
	tree := buildtest.ValidTree()
	for _, p := range buildtest.TransparencyPages() {
		tree.WithPage(p)
	}
	out := buildtest.MustBuild(t, buildtest.Config(), tree)
	srv := buildtest.Serve(t, out)

	for _, path := range []string{"/distributed-systems/raft-leader-election/", "/about/", "/authors/ada/", "/404.html", "/no-such-page/"} {
		res := srv.Get(t, path)
		if got := res.Header.Get("Cache-Control"); got != "public, max-age=0, must-revalidate" {
			t.Errorf("GET %s Cache-Control = %q, want browsers to revalidate every time", path, got)
		}
		cdn := res.Header.Get("Vercel-CDN-Cache-Control")
		if maxAge := directive(cdn, "s-maxage"); maxAge < 86400 {
			t.Errorf("GET %s Vercel-CDN-Cache-Control = %q, want s-maxage of at least a day", path, cdn)
		}
		if swr := directive(cdn, "stale-while-revalidate"); swr <= 0 {
			t.Errorf("GET %s Vercel-CDN-Cache-Control = %q, want stale-while-revalidate", path, cdn)
		}
	}
}

func TestRoutes_HashedAssetsAreImmutableAndNothingElseIs(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	// Pagefind's post-step writes stable file names (pagefind.js,
	// pagefind-entry.json…) that change content between builds.
	writeStatic(t, out, "/pagefind/pagefind.js", "export {}")
	srv := buildtest.Serve(t, out)

	var assets int
	for _, path := range buildtest.URLPaths(t, out) {
		cc := srv.Get(t, path).Header.Get("Cache-Control")
		immutable := strings.Contains(cc, "immutable")
		if strings.HasPrefix(path, "/assets/") {
			assets++
			if !immutable || directive(cc, "max-age") < 31536000 || !strings.Contains(cc, "public") {
				t.Errorf("GET %s Cache-Control = %q, want public, max-age=31536000, immutable", path, cc)
			}
		} else if immutable {
			t.Errorf("GET %s Cache-Control = %q; only content-hashed /assets/ may be immutable", path, cc)
		}
	}
	if assets < 10 {
		t.Errorf("only %d /assets/ files checked; want the hero variants", assets)
	}
}

func TestRoutes_GoneAndRedirectsComeBeforeTheFilesystemAndTheNotFoundFallbackLast(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree())
	routes, _ := buildtest.VercelConfig(t, out)["routes"].([]any)

	fs := -1
	for i, r := range routes {
		if r.(map[string]any)["handle"] == "filesystem" {
			if fs != -1 {
				t.Errorf("routes[%d]: second filesystem handler", i)
			}
			fs = i
		}
	}
	if fs == -1 {
		t.Fatalf("config.json routes have no {\"handle\":\"filesystem\"}: %v", routes)
	}
	seen := map[float64]bool{}
	for i, r := range routes {
		status, _ := r.(map[string]any)["status"].(float64)
		seen[status] = true
		switch status {
		case 410, 301, 308:
			if i > fs {
				t.Errorf("routes[%d] (status %v) comes after the filesystem handler; static files would win", i, status)
			}
		case 404:
			if i != len(routes)-1 || i < fs {
				t.Errorf("routes[%d] (status 404) is not the last route after the filesystem handler", i)
			}
		}
	}
	for _, s := range []float64{410, 301, 308, 404} {
		if !seen[s] {
			t.Errorf("no route with status %v", s)
		}
	}
	last := routes[len(routes)-1].(map[string]any)
	if last["dest"] != notFoundURL || last["status"] != float64(404) {
		t.Errorf("last route = %v, want the %s page with status 404", last, notFoundURL)
	}
}

// directive returns the seconds of a Cache-Control directive, or -1.
func directive(header, name string) int {
	for _, d := range strings.Split(header, ",") {
		k, v, _ := strings.Cut(strings.TrimSpace(d), "=")
		if strings.EqualFold(k, name) {
			n, err := strconv.Atoi(v)
			if err != nil {
				return -1
			}
			return n
		}
	}
	return -1
}

// writeStatic adds a file to the build output, as a post-build step would.
func writeStatic(t *testing.T, outDir, urlPath, body string) {
	t.Helper()
	p := buildtest.StaticPath(outDir, urlPath)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRoutes_HomeIsPermanentlyRedirectedToTheRoot(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree())
	srv := buildtest.Serve(t, out)

	for _, path := range []string{"/home", "/home/", "/Home"} {
		res := srv.Get(t, path)
		if res.StatusCode != http.StatusMovedPermanently || res.Header.Get("Location") != "/" {
			t.Errorf("GET %s = %d Location %q, want 301 to /", path, res.StatusCode, res.Header.Get("Location"))
		}
	}
	if res := srv.Get(t, "/homes/"); res.StatusCode == http.StatusMovedPermanently {
		t.Errorf("GET /homes/ = 301; only /home redirects to /")
	}
}

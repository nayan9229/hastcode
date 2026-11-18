package builder_test

import (
	"bytes"
	"image/color"
	"image/png"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

// The crawlable icon (ADR-0013, amending ADR-0006): a square PNG of the
// site mark at a stable path, linked from the home page only with
// rel="apple-touch-icon", so Googlebot-Image can fetch a file for the
// favicon beside our results. Every page keeps the data: URI favicon.

const (
	touchIconURL = "/apple-touch-icon.png"
	markURL      = "/mark.svg"
)

func TestTouchIcon_IsASquare192PxOpaquePNGOfTheMark(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())

	data, err := os.ReadFile(buildtest.StaticPath(out, touchIconURL))
	if err != nil {
		t.Fatalf("%s was not written: %v", touchIconURL, err)
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("%s is not a PNG: %v", touchIconURL, err)
	}
	// Google: square, larger than 48x48, ideally a multiple of 48.
	if b := img.Bounds(); b.Dx() != 192 || b.Dy() != 192 {
		t.Fatalf("%s is %dx%d, want 192x192", touchIconURL, b.Dx(), b.Dy())
	}
	// iOS fills transparent pixels with black: every pixel is opaque.
	for y := range 192 {
		for x := range 192 {
			if _, _, _, a := img.At(x, y).RGBA(); a != 0xFFFF {
				t.Fatalf("pixel (%d,%d) has alpha %d, want opaque", x, y, a>>8)
			}
		}
	}
	// The mark: brand teal ground (#0D9488), the white "h" on it. The
	// corner is ground (no transparent rounding: iOS and Google mask it),
	// the stem's middle (x 10.5, y 16 of the 32-unit mark) is white.
	teal := color.RGBA{0x0D, 0x94, 0x88, 0xFF}
	white := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	for _, p := range []struct {
		x, y int
		want color.RGBA
	}{
		{2, 2, teal}, {189, 189, teal}, {96, 40, teal},
		{63, 96, white}, {129, 140, white},
	} {
		if got := color.RGBAModel.Convert(img.At(p.x, p.y)).(color.RGBA); got != p.want {
			t.Errorf("pixel (%d,%d) = %v, want %v", p.x, p.y, got, p.want)
		}
	}
}

func TestTouchIcon_SameBytesOnEveryBuild(t *testing.T) {
	a := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	b := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree())
	first, err1 := os.ReadFile(buildtest.StaticPath(a, touchIconURL))
	second, err2 := os.ReadFile(buildtest.StaticPath(b, touchIconURL))
	if err1 != nil || err2 != nil {
		t.Fatalf("read icons: %v, %v", err1, err2)
	}
	if !bytes.Equal(first, second) {
		t.Error("two builds wrote different icon bytes; the stable URL must serve one file")
	}
}

// The icon's name is stable, so it must not be cached immutably: on each
// Target it gets exactly the headers /mark.svg gets (the revalidating
// default), and a request for it answers 200 with a PNG.
func TestTouchIcon_CachedLikeTheMarkOnEachTarget(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())

	// Production: Cloudflare Workers Static Assets, _headers.
	mark := buildtest.CloudflareHeaders(t, out, "https://hastcode.com"+markURL)
	icon := buildtest.CloudflareHeaders(t, out, "https://hastcode.com"+touchIconURL)
	if !reflect.DeepEqual(icon, mark) {
		t.Errorf("Cloudflare headers for %s = %v, want the same as %s: %v", touchIconURL, icon, markURL, mark)
	}
	if cc := icon.Get("Cache-Control"); strings.Contains(cc, "immutable") {
		t.Errorf("Cloudflare Cache-Control for %s = %q; only content-hashed /assets/ files are immutable", touchIconURL, cc)
	}

	// QA: Vercel's config.json routes, through the local router.
	srv := buildtest.Serve(t, out)
	markRes, iconRes := srv.Get(t, markURL), srv.Get(t, touchIconURL)
	if iconRes.StatusCode != http.StatusOK || !strings.HasPrefix(iconRes.Body, "\x89PNG") {
		t.Fatalf("GET %s = %d, want 200 with the PNG", touchIconURL, iconRes.StatusCode)
	}
	for _, h := range []string{"Cache-Control", "Vercel-CDN-Cache-Control", "X-Robots-Tag"} {
		if got, want := iconRes.Header.Get(h), markRes.Header.Get(h); got != want {
			t.Errorf("Vercel %s for %s = %q, want %q as for %s", h, touchIconURL, got, want, markURL)
		}
	}
}

// The icon is a file, not a page: never in the sitemap, and robots.txt
// lets Googlebot-Image fetch it (Google's favicon requirement).
func TestTouchIcon_NotInTheSitemapAndNotBlockedByRobots(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())

	for _, loc := range readSitemap(t, out, "/sitemap.xml").locs() {
		if strings.HasSuffix(loc, touchIconURL) {
			t.Errorf("sitemap.xml lists %s", loc)
		}
	}
	robots, err := os.ReadFile(buildtest.StaticPath(out, "/robots.txt"))
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(string(robots), "\n") {
		rule, ok := strings.CutPrefix(strings.TrimSpace(line), "Disallow:")
		if rule = strings.TrimSpace(rule); ok && rule != "" && strings.HasPrefix(touchIconURL, rule) {
			t.Errorf("robots.txt %q blocks %s", line, touchIconURL)
		}
	}
	if strings.HasPrefix(touchIconURL, "/assets/") {
		t.Errorf("%s is under /assets/, whose files are cached immutably; its name is stable", touchIconURL)
	}
}

func TestTouchIcon_LinkedFromTheHomePageOnlyAndEveryPageKeepsTheDataFavicon(t *testing.T) {
	tree := buildtest.ValidTree().WithBrief(buildtest.ValidBrief())
	for _, p := range buildtest.TransparencyPages() {
		tree.WithPage(p)
	}
	out := buildtest.MustBuild(t, buildtest.Config(), tree)

	pages := htmlPages(t, out)
	if len(pages) < 10 {
		t.Fatalf("only %d pages: %v", len(pages), pages)
	}
	for _, urlPath := range pages {
		doc := buildtest.Page(t, out, urlPath)
		touch := doc.Find(`link[rel="apple-touch-icon"]`)
		if urlPath == "/" {
			if touch.Length() != 1 || touch.AttrOr("href", "") != touchIconURL || !touch.Parent().Is("head") {
				t.Errorf("home page apple-touch-icon links = %d (href %q), want one <head> link to %s",
					touch.Length(), touch.AttrOr("href", ""), touchIconURL)
			}
		} else if touch.Length() != 0 {
			t.Errorf("%s links the apple-touch-icon; only the home page does (ADR-0013)", urlPath)
		}
		icon := doc.Find(`head link[rel="icon"]`)
		if icon.Length() != 1 || !strings.HasPrefix(icon.AttrOr("href", ""), "data:image/svg+xml,") {
			t.Errorf("%s has %d rel=icon links (first href %.30q), want the one data: favicon", urlPath, icon.Length(), icon.AttrOr("href", ""))
		}
	}
}

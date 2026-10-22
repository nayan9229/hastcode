package builder_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/internal/buildtest"
)

func TestHero_ImgIsFetchedAtHighPriorityAndNeverLazyOnExplainerAndBriefPages(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), briefTree())

	for _, page := range []string{raftPage, briefURL} {
		doc := buildtest.Page(t, out, page)
		img := doc.Find("article picture img")
		if img.Length() != 1 {
			t.Fatalf("%s: <picture><img>: found %d, want 1", page, img.Length())
		}
		if got := img.AttrOr("fetchpriority", ""); got != "high" {
			t.Errorf("%s: hero <img> fetchpriority = %q, want high", page, got)
		}
		if got, ok := img.Attr("loading"); ok {
			t.Errorf("%s: hero <img> has loading=%q; the LCP image must not carry a loading attribute", page, got)
		}
		if n := doc.Find(`head link[rel="preload"][as="image"]`).Length(); n != 1 {
			t.Errorf("%s: hero preload links: found %d, want 1", page, n)
		}
	}
}

func TestHero_JPEGWithEXIFOrientationIsRotatedUprightBeforeResizing(t *testing.T) {
	// Upright the hero is 1600x900 with a red top-left quadrant. Stored for
	// orientation 6 ("0th row is the visual right side, 0th column the
	// visual top"), it is 900x1600: stored(x, y) = upright(1599-y, x).
	upright := buildtest.Quadrants(1600, 900)
	stored := image.NewRGBA(image.Rect(0, 0, 900, 1600))
	for y := range 1600 {
		for x := range 900 {
			stored.Set(x, y, upright.At(1599-y, x))
		}
	}
	for name, bigEndian := range map[string]bool{"little-endian EXIF": false, "big-endian EXIF": true} {
		t.Run(name, func(t *testing.T) {
			e := buildtest.ValidExplainer()
			e.Hero = "raft-leader-election.jpg"
			e.HeroImage = buildtest.WithEXIFOrientation(buildtest.HeroJPEG(stored), 6, bigEndian)

			out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

			picture := buildtest.Page(t, out, raftPage).Find("article picture")
			var urls []string
			for _, sel := range []string{"source", "img"} {
				picture.Find(sel).Each(func(_ int, s *goquery.Selection) {
					for _, c := range buildtest.Srcset(t, s.AttrOr("srcset", "")) {
						urls = append(urls, c.URL)
					}
				})
			}
			urls = append(urls, picture.Find("img").AttrOr("src", ""))
			if len(urls) < 2 {
				t.Fatalf("hero <picture> lists %d images", len(urls))
			}
			for _, u := range urls {
				requireUprightQuadrants(t, buildtest.DecodePixels(t, out, u), u)
			}
		})
	}
}

// requireUprightQuadrants checks a decoded variant of buildtest.Quadrants
// (16:9): landscape, red in the top-left quadrant, blue in the other three.
func requireUprightQuadrants(t *testing.T, img image.Image, url string) {
	t.Helper()
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= h || h != (w*9+8)/16 && h != w*9/16 {
		t.Errorf("%s is %dx%d, want an upright 16:9 landscape", url, w, h)
		return
	}
	isRed := func(x, y int) bool {
		r, g, bl, _ := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
		return r>>8 > g>>8+80 && r>>8 > bl>>8+80
	}
	isBlue := func(x, y int) bool {
		r, _, bl, _ := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
		return bl>>8 > r>>8+80
	}
	if !isRed(w/4, h/4) {
		t.Errorf("%s: top-left quadrant is not red: the hero was not rotated upright", url)
	}
	for _, p := range [][2]int{{3 * w / 4, h / 4}, {w / 4, 3 * h / 4}, {3 * w / 4, 3 * h / 4}} {
		if !isBlue(p[0], p[1]) {
			t.Errorf("%s: pixel at %v is not blue: the hero was not rotated upright", url, p)
		}
	}
}

func TestHero_NonSRGBColourProfileIsAHeroColourWarning(t *testing.T) {
	jpg := buildtest.HeroJPEG(buildtest.Quadrants(1600, 900))
	for name, c := range map[string]struct {
		hero, profile string
		image         []byte
	}{
		"PNG, Display P3 (v4)":    {"raft-leader-election.png", "Display P3", buildtest.WithPNGICCProfile(buildtest.HeroPNG(1600, 900), buildtest.ICCProfile("Display P3", true))},
		"JPEG, Adobe RGB (v2)":    {"raft-leader-election.jpg", "Adobe RGB (1998)", buildtest.WithJPEGICCProfile(jpg, buildtest.ICCProfile("Adobe RGB (1998)", false))},
		"JPEG, EXIF + Display P3": {"raft-leader-election.jpg", "Display P3", buildtest.WithJPEGICCProfile(buildtest.WithEXIFOrientation(jpg, 1, false), buildtest.ICCProfile("Display P3", true))},
	} {
		t.Run(name, func(t *testing.T) {
			e := buildtest.ValidExplainer()
			e.Hero, e.HeroImage = c.hero, c.image

			out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

			buildtest.RequireWarningRule(t, res, raftFile, "hero-colour", c.profile)
			buildtest.RequireWarningRule(t, res, raftFile, "hero-colour", "export it as sRGB")
			if len(res.Errors) > 0 || !buildtest.Exists(out, raftPage) {
				t.Errorf("a non-sRGB hero must only warn; got:\n%s", res)
			}
		})
	}
}

func TestHero_SRGBOrUntaggedSourceRaisesNoColourWarning(t *testing.T) {
	jpg := buildtest.HeroJPEG(buildtest.Quadrants(1600, 900))
	for name, c := range map[string]struct {
		hero  string
		image []byte
	}{
		"PNG, no profile":  {"raft-leader-election.png", buildtest.HeroPNG(1600, 900)},
		"JPEG, no profile": {"raft-leader-election.jpg", jpg},
		"JPEG, sRGB (v2)":  {"raft-leader-election.jpg", buildtest.WithJPEGICCProfile(jpg, buildtest.ICCProfile("sRGB IEC61966-2.1", false))},
		"PNG, sRGB (v4)":   {"raft-leader-election.png", buildtest.WithPNGICCProfile(buildtest.HeroPNG(1600, 900), buildtest.ICCProfile("sRGB", true))},
	} {
		t.Run(name, func(t *testing.T) {
			e := buildtest.ValidExplainer()
			e.Hero, e.HeroImage = c.hero, c.image

			_, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

			if res.Failed() {
				t.Fatalf("build failed:\n%s", res)
			}
			buildtest.RequireNoProblemRule(t, res, raftFile, "hero-colour")
		})
	}
}

func TestHero_WebPSourceIsAHeroFormatError(t *testing.T) {
	teal := color.RGBA{R: 0x0D, G: 0x94, B: 0x88, A: 0xFF}
	for name, hero := range map[string]string{
		"webp extension": "raft-leader-election.webp",
		"png extension":  "raft-leader-election.png", // a WebP file renamed .png
	} {
		t.Run(name, func(t *testing.T) {
			e := buildtest.ValidExplainer()
			e.Hero = hero
			e.HeroImage = buildtest.HeroWebP(1600, 900, teal)

			out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

			buildtest.RequireErrorRule(t, res, raftFile, "hero-format", "WebP")
			if buildtest.Exists(out, raftPage) {
				t.Error("page was written for an Explainer with a WebP hero")
			}
		})
	}
}

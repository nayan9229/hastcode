package builder_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"maps"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/internal/buildtest"
)

// Generated hero card (spec "Page composition and gates": generated hero).
// A Repo Profile needs no sourced image: without `hero` the build draws a
// 1600x900 card from the Fact Pack and feeds it through the hero pipeline
// like a sourced file, so the page gets everything the hero tests require.

func TestHeroCard_ProfileWithoutAHeroGetsTheFiveWidthsTheOGCropAndThePreload(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), repoProfileTree()) // ValidRepoProfile names no hero
	doc := buildtest.Page(t, out, repoProfilePath)
	picture := doc.Find("article picture")

	assertHeroSrcset(t, out, picture.Find(`source[type="image/avif"]`), "avif", []int{480, 800, 1000, 1200, 1472})
	assertHeroSrcset(t, out, picture.Find("img"), "jpeg", []int{480, 800, 1000, 1200, 1472})

	img := picture.Find("img")
	if got := img.AttrOr("fetchpriority", ""); got != "high" {
		t.Errorf("hero <img> fetchpriority = %q, want high", got)
	}
	if got, ok := img.Attr("loading"); ok {
		t.Errorf("hero <img> has loading=%q; the LCP image must not carry a loading attribute", got)
	}
	src := buildtest.DecodeImage(t, out, img.AttrOr("src", ""))
	if src.Format != "jpeg" || src.Width != 1200 || src.Height != 675 {
		t.Errorf("img src is %s %dx%d, want jpeg 1200x675", src.Format, src.Width, src.Height)
	}

	preload := doc.Find(`head link[rel="preload"][as="image"]`)
	if preload.Length() != 1 {
		t.Fatalf("hero preload links: found %d, want 1", preload.Length())
	}
	avif := picture.Find(`source[type="image/avif"]`)
	if got, want := preload.AttrOr("imagesrcset", ""), avif.AttrOr("srcset", ""); got == "" || got != want {
		t.Errorf("preload imagesrcset = %q, want the AVIF srcset %q", got, want)
	}

	og := doc.Find(`head meta[property="og:image"]`).AttrOr("content", "")
	path, ok := strings.CutPrefix(og, siteURL+"/")
	if !ok {
		t.Fatalf("og:image = %q, want an absolute URL on the site", og)
	}
	if crop := buildtest.DecodeImage(t, out, "/"+path); crop.Format != "jpeg" || crop.Width != 1200 || crop.Height != 630 {
		t.Errorf("og:image is a %dx%d %s, want the 1200x630 JPEG crop of the card", crop.Width, crop.Height, crop.Format)
	}
	article := buildtest.JSONLDOfType(t, doc, "Article")
	images, _ := article["image"].([]any)
	if len(images) != 3 {
		t.Fatalf("Article image = %#v, want the 16:9, 4:3 and 1:1 crops", article["image"])
	}
	for i, want := range [][2]int{{1200, 675}, {1200, 900}, {1200, 1200}} {
		url, _ := images[i].(string)
		p, _ := strings.CutPrefix(url, siteURL)
		if crop := buildtest.DecodeImage(t, out, p); crop.Format != "jpeg" || crop.Width != want[0] || crop.Height != want[1] {
			t.Errorf("Article image[%d] %s is %s %dx%d, want jpeg %dx%d", i, url, crop.Format, crop.Width, crop.Height, want[0], want[1])
		}
	}
}

func TestHeroCard_AltTextIsTheRepositoryAndItsDescriptionOnTheImgAndOpenGraph(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), repoProfileTree())
	doc := buildtest.Page(t, out, repoProfilePath)

	const want = "apache/iggy: Persistent message streaming platform written in Rust"
	if got := doc.Find("article picture img").AttrOr("alt", ""); got != want {
		t.Errorf("hero alt = %q, want %q (the Fact Pack's repository and description)", got, want)
	}
	if got := doc.Find(`head meta[property="og:image:alt"]`).AttrOr("content", ""); got != want {
		t.Errorf("og:image:alt = %q, want %q", got, want)
	}

	// GitHub has no description: the card and its alt fall back to the
	// Profile's own.
	p := buildtest.ValidRepoProfile()
	p.FactPack.Description = ""
	out = buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithRepoProfile(p))
	if got := buildtest.Page(t, out, repoProfilePath).Find("article picture img").AttrOr("alt", ""); got != "apache/iggy: "+p.Description {
		t.Errorf("hero alt without a GitHub description = %q, want the repository and the Profile's description", got)
	}
}

// heroCardFile is the card the build draws for the fixture Profile.
func heroCardFile(cacheDir string) string {
	return buildtest.HeroCardFile(cacheDir, "explainers/distributed-systems/apache-iggy.md")
}

// Spec: the hero pipeline cuts 16:9, 4:3, 1:1 and 1200x630 crops from the
// 1600x900 card. The 1:1 crop keeps the middle 900 px and the OG crop trims
// 30 px top and bottom, so everything but the bottom band must sit inside
// that, for a short Fact Pack and for one whose every line is too long.
func TestHeroCard_IsA1600x900PNGWithEveryLineInsideTheSquareAndOpenGraphCrops(t *testing.T) {
	long := buildtest.ValidRepoProfile()
	long.Repo = "kubernetes-sigs/cluster-api-provider-aws"
	long.FactPack.FullName = long.Repo
	long.FactPack.Description = "Kubernetes Cluster API Provider AWS provides consistent deployment and day 2 operations of self-managed and EKS Kubernetes clusters on AWS."
	long.FactPack.Language = "Jupyter Notebook"
	long.FactPack.Stars = 1234567
	long.FactPack.License = &buildtest.License{SPDXID: "GPL-3.0-or-later", Name: "GNU General Public License v3.0 or later", URL: "https://spdx.org/licenses/GPL-3.0-or-later.html"}
	long.FactPack.LatestRelease.Tag = "v2.10.0-beta.1"

	for name, p := range map[string]buildtest.Explainer{"short": buildtest.ValidRepoProfile(), "long": long} {
		t.Run(name, func(t *testing.T) {
			cfg := buildtest.Config()
			cfg.CacheDir = t.TempDir()
			buildtest.MustBuild(t, cfg, buildtest.NewTree().WithRepoProfile(p))

			data, err := os.ReadFile(heroCardFile(cfg.CacheDir))
			if err != nil {
				t.Fatalf("the build drew no card: %v", err)
			}
			img, err := png.Decode(bytes.NewReader(data))
			if err != nil {
				t.Fatal(err)
			}
			if b := img.Bounds(); b.Dx() != 1600 || b.Dy() != 900 {
				t.Fatalf("card is %dx%d, want 1600x900", b.Dx(), b.Dy())
			}
			minX, minY, maxX, maxY := 1600, 900, -1, -1
			for y := 0; y < 860; y++ { // above the accent band
				for x := range 1600 {
					if !isWhite(img, x, y) {
						minX, minY, maxX, maxY = min(minX, x), min(minY, y), max(maxX, x), max(maxY, y)
					}
				}
			}
			if maxX < 0 {
				t.Fatal("nothing is drawn above the band")
			}
			if minX < 350 || maxX >= 1250 {
				t.Errorf("ink spans x %d..%d, want within 350..1249 (the 1:1 crop)", minX, maxX)
			}
			if minY < 30 || maxY >= 870 {
				t.Errorf("ink spans y %d..%d, want within 30..869 (the 1200x630 crop)", minY, maxY)
			}
		})
	}
}

func isWhite(img image.Image, x, y int) bool {
	r, g, b, _ := img.At(x, y).RGBA()
	return r>>8 == 0xFF && g>>8 == 0xFF && b>>8 == 0xFF
}

func TestHeroCard_ProfileWithAHeroUsesItUnchanged(t *testing.T) {
	p := buildtest.ValidRepoProfile()
	p.Hero, p.HeroAlt = "apache-iggy.png", "The Iggy server fanning a stream out to consumers"
	amber := color.RGBA{R: 0xF5, G: 0x9E, B: 0x0B, A: 0xFF}
	p.HeroImage = buildtest.HeroPNGFilled(1600, 900, amber)

	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithRepoProfile(p))

	img := buildtest.Page(t, out, repoProfilePath).Find("article picture img")
	if got := img.AttrOr("alt", ""); got != p.HeroAlt {
		t.Errorf("hero alt = %q, want hero_alt from frontmatter %q", got, p.HeroAlt)
	}
	src := img.AttrOr("src", "")
	if !strings.HasPrefix(src, "/assets/heroes/apache-iggy-1200x675.") {
		t.Errorf("img src = %q, want the sourced hero's variant, named after apache-iggy.png", src)
	}
	px := buildtest.DecodePixels(t, out, src)
	requireColour(t, px, px.Bounds().Dx()/2, px.Bounds().Dy()/2, amber, "the sourced hero's amber")
}

func TestHeroCard_FileNamesChangeWhenStarsOrTheLatestReleaseChangeButNotWhenOnlyTheFetchTimeDoes(t *testing.T) {
	base := pictureURLs(t, buildtest.MustBuild(t, buildtest.Config(), repoProfileTree()), repoProfilePath)
	if len(base) == 0 {
		t.Fatal("the Profile's <picture> references no hero files")
	}
	for name, tc := range map[string]struct {
		change  func(fp *buildtest.FactPack)
		renamed bool
	}{
		"stars": {func(fp *buildtest.FactPack) { fp.Stars++; fp.Stargazers++ }, true},
		"latest release": {func(fp *buildtest.FactPack) {
			fp.LatestRelease.Tag, fp.LatestRelease.Name = "server-0.6.0", "Server 0.6.0"
		}, true},
		"fetched_at": {func(fp *buildtest.FactPack) { fp.FetchedAt = "2026-10-26T06:00:00Z" }, false},
	} {
		t.Run(name, func(t *testing.T) {
			p := buildtest.ValidRepoProfile()
			tc.change(p.FactPack)

			got := pictureURLs(t, buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithRepoProfile(p)), repoProfilePath)

			if len(got) != len(base) {
				t.Fatalf("hero URL counts: %d before, %d after", len(base), len(got))
			}
			for i := range base {
				if tc.renamed && got[i] == base[i] {
					t.Errorf("hero URL %s did not change when the %s changed", base[i], name)
				}
				if !tc.renamed && got[i] != base[i] {
					t.Errorf("hero URL %s changed to %s when only %s changed (every Refresh would bust every cache)", base[i], got[i], name)
				}
			}
		})
	}
}

func TestHeroCard_SecondBuildOfAnUnchangedFactPackDrawsTheSameBytesAndHitsTheImageCache(t *testing.T) {
	cfg := buildtest.Config()
	cfg.Images.CacheDir = t.TempDir()

	cfg.CacheDir = t.TempDir()
	first := buildtest.MustBuild(t, cfg, repoProfileTree())
	firstCard := heroCardFile(cfg.CacheDir)
	before := cacheSnapshot(t, cfg.Images.CacheDir)
	cfg.CacheDir = t.TempDir() // fresh content, output and build cache: the card is drawn again
	second := buildtest.MustBuild(t, cfg, repoProfileTree())
	after := cacheSnapshot(t, cfg.Images.CacheDir)

	a, errA := os.ReadFile(firstCard)
	b, errB := os.ReadFile(heroCardFile(cfg.CacheDir))
	if errA != nil || errB != nil || !bytes.Equal(a, b) {
		t.Errorf("two builds of the same Fact Pack drew different cards (%v, %v)", errA, errB)
	}

	if !maps.Equal(before, after) {
		t.Errorf("image cache changed on a rebuild of the same Fact Pack: the card was re-encoded, so it was not drawn to the same bytes:\nbefore %v\nafter  %v", before, after)
	}
	if a, b := pictureURLs(t, first, repoProfilePath), pictureURLs(t, second, repoProfilePath); !slices.Equal(a, b) {
		t.Errorf("hero URLs changed between builds of the same Fact Pack:\n%v\n%v", a, b)
	}
}

func TestHeroCard_IsDrawnOnTheLightBrandPalette(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), repoProfileTree())
	src := buildtest.Page(t, out, repoProfilePath).Find("article picture img").AttrOr("src", "")
	px := buildtest.DecodePixels(t, out, src)
	w, h := px.Bounds().Dx(), px.Bounds().Dy()

	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	teal := color.RGBA{R: 0x0D, G: 0x94, B: 0x88, A: 0xFF}
	requireColour(t, px, 8, 8, white, "the light background")
	requireColour(t, px, w-9, 8, white, "the light background")
	requireColour(t, px, 8, h/2, white, "the light background")
	requireColour(t, px, w/2, h-6, teal, "the accent band along the bottom")
}

// requireColour checks a pixel is within JPEG tolerance of want.
func requireColour(t *testing.T, img image.Image, x, y int, want color.RGBA, what string) {
	t.Helper()
	b := img.Bounds()
	r, g, bl, _ := img.At(b.Min.X+x, b.Min.Y+y).RGBA()
	got := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(bl >> 8), A: 0xFF}
	off := func(a, b uint8) int { return max(int(a)-int(b), int(b)-int(a)) }
	if off(got.R, want.R) > 24 || off(got.G, want.G) > 24 || off(got.B, want.B) > 24 {
		t.Errorf("%s: pixel (%d, %d) = #%02X%02X%02X, want about #%02X%02X%02X", what, x, y, got.R, got.G, got.B, want.R, want.G, want.B)
	}
}

// pictureURLs is every hero file URL a page's <picture> references, in
// document order (like heroURLs, for any page).
func pictureURLs(t *testing.T, out, page string) []string {
	t.Helper()
	var urls []string
	buildtest.Page(t, out, page).Find("article picture source, article picture img").Each(func(_ int, s *goquery.Selection) {
		for _, e := range buildtest.Srcset(t, s.AttrOr("srcset", "")) {
			urls = append(urls, e.URL)
		}
	})
	return urls
}

package builder_test

import (
	"bytes"
	"cmp"
	"image/color"
	"image/png"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

// Generated title card (ticket 10). Every Piece without `hero` gets its
// hero from the build, the way a Repo Profile does: a Case Study, Concept,
// Comparison or Brief gets a card drawn from its title, kind and Vertical
// (or "Brief") and fed through the hero pipeline like a sourced file.

// titleCardKinds are the kinds that get a title card.
var titleCardKinds = []string{"case study", "concept", "comparison", "brief"}

// titleCardPiece is a Piece without `hero`: its tree, page and content file.
type titleCardPiece struct {
	tree       buildtest.Tree
	page, file string
	title      string
}

// noHeroPiece is the fixture of kind (one of titleCardKinds) without
// `hero`, titled title ("" keeps the fixture's).
func noHeroPiece(kind, title string) titleCardPiece {
	if kind == "brief" {
		b := buildtest.ValidBrief()
		b.Hero, b.HeroAlt, b.HeroImage = "", "", nil
		b.Title = cmp.Or(title, b.Title)
		return titleCardPiece{buildtest.ValidTree().WithBrief(b), briefURL, briefFile, b.Title}
	}
	e := buildtest.ValidExplainer()
	page, file := raftURL, raftFile
	switch kind {
	case "case study":
		e.Kind = "case-study"
	case "comparison":
		e = buildtest.ValidComparison()
		page, file = "/databases-storage/btree-vs-lsm-tree/", "explainers/databases-storage/btree-vs-lsm-tree.md"
	}
	e.Hero, e.HeroAlt, e.HeroImage = "", "", nil
	e.Title = cmp.Or(title, e.Title)
	return titleCardPiece{buildtest.NewTree().WithExplainer(e), page, file, e.Title}
}

func TestTitleCard_EveryKindWithoutAHeroGetsACardThroughTheHeroPipeline(t *testing.T) {
	for _, name := range titleCardKinds {
		p := noHeroPiece(name, "")
		t.Run(name, func(t *testing.T) {
			out, res := buildtest.Build(t, buildtest.Config(), p.tree)
			buildtest.RequireNoProblemRule(t, res, p.file, "hero-missing")
			if !buildtest.Exists(out, p.page) {
				t.Fatalf("page %s was not written:\n%s", p.page, res)
			}
			doc := buildtest.Page(t, out, p.page)
			img := doc.Find("article picture img")

			src := buildtest.DecodeImage(t, out, img.AttrOr("src", ""))
			if src.Width < 1200 {
				t.Errorf("hero img src is %dx%d, want at least 1200 px wide", src.Width, src.Height)
			}
			widest := 0
			for _, e := range buildtest.Srcset(t, doc.Find(`article picture source[type="image/avif"]`).AttrOr("srcset", "")) {
				widest = max(widest, buildtest.DecodeImage(t, out, e.URL).Width)
			}
			if widest < 1200 {
				t.Errorf("widest AVIF variant is %d px, want at least 1200", widest)
			}
			if got := img.AttrOr("fetchpriority", ""); got != "high" {
				t.Errorf("hero <img> fetchpriority = %q, want high (the card is the LCP image)", got)
			}
			if n := doc.Find(`head link[rel="preload"][as="image"]`).Length(); n != 1 {
				t.Errorf("hero preload links: %d, want 1", n)
			}

			if got := img.AttrOr("alt", ""); got != p.title {
				t.Errorf("hero alt = %q, want the title %q", got, p.title)
			}
			if got := doc.Find(`head meta[property="og:image:alt"]`).AttrOr("content", ""); got != p.title {
				t.Errorf("og:image:alt = %q, want the title %q", got, p.title)
			}

			og := doc.Find(`head meta[property="og:image"]`).AttrOr("content", "")
			path, ok := strings.CutPrefix(og, siteURL+"/")
			if !ok {
				t.Fatalf("og:image = %q, want an absolute URL on the site", og)
			}
			if crop := buildtest.DecodeImage(t, out, "/"+path); crop.Format != "jpeg" || crop.Width != 1200 || crop.Height != 630 {
				t.Errorf("og:image is a %dx%d %s, want the 1200x630 JPEG crop of the card", crop.Width, crop.Height, crop.Format)
			}
			article := buildtest.JSONLDOfType(t, doc, articleType(name))
			images, _ := article["image"].([]any)
			if len(images) != 3 {
				t.Fatalf("%s image = %#v, want the 16:9, 4:3 and 1:1 crops", articleType(name), article["image"])
			}
			for i, want := range [][2]int{{1200, 675}, {1200, 900}, {1200, 1200}} {
				url, _ := images[i].(string)
				p, _ := strings.CutPrefix(url, siteURL)
				if crop := buildtest.DecodeImage(t, out, p); crop.Format != "jpeg" || crop.Width != want[0] || crop.Height != want[1] {
					t.Errorf("image[%d] %s is %s %dx%d, want jpeg %dx%d", i, url, crop.Format, crop.Width, crop.Height, want[0], want[1])
				}
			}
		})
	}
}

// articleType is the JSON-LD type of a titleCardKinds page.
func articleType(name string) string {
	if name == "brief" {
		return "NewsArticle"
	}
	return "Article"
}

// titleCardFile is the card the build draws for a Piece: written under the
// build cache, mirroring the content path (README "generated hero").
func titleCardFile(cacheDir, contentFile string) string {
	return filepath.Join(cacheDir, "hero-cards", filepath.FromSlash(strings.TrimSuffix(contentFile, ".md")+".png"))
}

// drawnTitleCard builds p into a fresh build cache and returns the card it
// drew.
func drawnTitleCard(t *testing.T, p titleCardPiece) []byte {
	t.Helper()
	cfg := buildtest.Config()
	cfg.CacheDir = t.TempDir()
	buildtest.MustBuild(t, cfg, p.tree)
	data, err := os.ReadFile(titleCardFile(cfg.CacheDir, p.file))
	if err != nil {
		t.Fatalf("the build drew no card: %v", err)
	}
	return data
}

// The hero pipeline cuts 16:9, 4:3, 1:1 and 1200x630 crops from the
// 1600x900 card, so everything but the bottom band must sit inside the
// middle 900 px and off the 30 px the Open Graph crop trims top and
// bottom: for a short title, a typical one and one too long for any size.
func TestTitleCard_IsA1600x900PNGWithEveryLineInsideTheSquareAndOpenGraphCrops(t *testing.T) {
	titles := map[string]string{
		"short":   "Raft",
		"typical": "Connection reuse in HTTP/1.1, HTTP/2 and HTTP/3: how each works and where each breaks",
		"long":    "Connection reuse in HTTP/1.1, HTTP/2 and HTTP/3: how each works, where each breaks, and which one to use between services in a mesh that spans three regions and two clouds",
		"word":    "Supercalifragilisticexpialidociousconnectionpoolingthatneverends",
	}
	for name, title := range titles {
		for _, kind := range []string{"comparison", "brief"} {
			t.Run(name+" "+kind, func(t *testing.T) {
				img, err := png.Decode(bytes.NewReader(drawnTitleCard(t, noHeroPiece(kind, title))))
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
				for y := 860; y < 900; y++ { // the band is solid accent: no text runs into it
					for x := range 1600 {
						if r, g, b, _ := img.At(x, y).RGBA(); r>>8 != 0x0D || g>>8 != 0x94 || b>>8 != 0x88 {
							t.Fatalf("band pixel (%d, %d) is not the accent: the text runs into the band", x, y)
						}
					}
				}
			})
		}
	}
}

func TestTitleCard_IsDrawnOnTheLightBrandPalette(t *testing.T) {
	p := noHeroPiece("case study", "")
	out := buildtest.MustBuild(t, buildtest.Config(), p.tree)
	src := buildtest.Page(t, out, p.page).Find("article picture img").AttrOr("src", "")
	px := buildtest.DecodePixels(t, out, src)
	w, h := px.Bounds().Dx(), px.Bounds().Dy()

	white := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	teal := color.RGBA{R: 0x0D, G: 0x94, B: 0x88, A: 0xFF}
	requireColour(t, px, 8, 8, white, "the light background")
	requireColour(t, px, w-9, 8, white, "the light background")
	requireColour(t, px, 8, h/2, white, "the light background")
	requireColour(t, px, w/2, h-6, teal, "the accent band along the bottom")
}

// The card shows the title, the kind and the Vertical: change any one and
// the card is drawn anew, so its variants get new names.
func TestTitleCard_DiffersWhenTheTitleKindOrVerticalDiffers(t *testing.T) {
	concept := noHeroPiece("concept", "")
	base := drawnTitleCard(t, concept)

	otherVertical := buildtest.ValidExplainer()
	otherVertical.Hero, otherVertical.HeroAlt, otherVertical.HeroImage = "", "", nil
	otherVertical.Vertical = "infrastructure-cloud"

	for name, p := range map[string]titleCardPiece{
		"title":    noHeroPiece("concept", "How Raft elects a leader when the network splits"),
		"kind":     noHeroPiece("case study", ""),
		"vertical": {tree: buildtest.NewTree().WithExplainer(otherVertical), file: "explainers/infrastructure-cloud/raft-leader-election.md"},
	} {
		t.Run(name, func(t *testing.T) {
			if bytes.Equal(drawnTitleCard(t, p), base) {
				t.Errorf("the card did not change when the %s did", name)
			}
		})
	}

	retitled := noHeroPiece("concept", "How Raft elects a leader when the network splits")
	a := pictureURLs(t, buildtest.MustBuild(t, buildtest.Config(), concept.tree), concept.page)
	b := pictureURLs(t, buildtest.MustBuild(t, buildtest.Config(), retitled.tree), retitled.page)
	if len(a) == 0 || len(a) != len(b) {
		t.Fatalf("hero URLs: %d before, %d after", len(a), len(b))
	}
	for i := range a {
		if a[i] == b[i] {
			t.Errorf("hero URL %s did not change when the title changed", a[i])
		}
	}
}

func TestTitleCard_SecondBuildOfAnUnchangedPieceDrawsTheSameBytesAndHitsTheImageCache(t *testing.T) {
	for _, kind := range []string{"case study", "brief"} {
		t.Run(kind, func(t *testing.T) {
			p := noHeroPiece(kind, "")
			cfg := buildtest.Config()
			cfg.Images.CacheDir = t.TempDir()

			cfg.CacheDir = t.TempDir()
			first := buildtest.MustBuild(t, cfg, p.tree)
			firstCard := titleCardFile(cfg.CacheDir, p.file)
			before := cacheSnapshot(t, cfg.Images.CacheDir)
			cfg.CacheDir = t.TempDir() // fresh output and build cache: the card is drawn again
			second := buildtest.MustBuild(t, cfg, p.tree)
			after := cacheSnapshot(t, cfg.Images.CacheDir)

			a, errA := os.ReadFile(firstCard)
			b, errB := os.ReadFile(titleCardFile(cfg.CacheDir, p.file))
			if errA != nil || errB != nil || !bytes.Equal(a, b) {
				t.Errorf("two builds of the same Piece drew different cards (%v, %v)", errA, errB)
			}
			if !maps.Equal(before, after) {
				t.Errorf("image cache changed on a rebuild of the same Piece: the card was re-encoded:\nbefore %v\nafter  %v", before, after)
			}
			if a, b := pictureURLs(t, first, p.page), pictureURLs(t, second, p.page); len(a) == 0 || !slices.Equal(a, b) {
				t.Errorf("hero URLs changed between builds of the same Piece:\n%v\n%v", a, b)
			}
		})
	}
}

// The Repo Profile rule, for every kind: `hero` is optional, and a Piece
// that names one keeps it (with its hero_alt) instead of the card
// (TestHeroCard_ProfileWithAHeroUsesItUnchanged); a named hero that is not
// there is still an error, never silently replaced by the card.
func TestTitleCard_APieceWithAHeroUsesItUnchanged(t *testing.T) {
	amber := color.RGBA{R: 0xF5, G: 0x9E, B: 0x0B, A: 0xFF}
	caseStudy := buildtest.ValidExplainer()
	caseStudy.Kind = "case-study"
	caseStudy.HeroImage = buildtest.HeroPNGFilled(1600, 900, amber)
	brief := buildtest.ValidBrief()
	brief.HeroImage = buildtest.HeroPNGFilled(1600, 900, amber)

	for name, tc := range map[string]struct {
		tree            buildtest.Tree
		page, hero, alt string
	}{
		"case study": {buildtest.NewTree().WithExplainer(caseStudy), raftURL, caseStudy.Hero, caseStudy.HeroAlt},
		"brief":      {buildtest.ValidTree().WithBrief(brief), briefURL, brief.Hero, brief.HeroAlt},
	} {
		t.Run(name, func(t *testing.T) {
			cfg := buildtest.Config()
			cfg.CacheDir = t.TempDir()
			out := buildtest.MustBuild(t, cfg, tc.tree)

			img := buildtest.Page(t, out, tc.page).Find("article picture img")
			if got := img.AttrOr("alt", ""); got != tc.alt {
				t.Errorf("hero alt = %q, want hero_alt from frontmatter %q", got, tc.alt)
			}
			src := img.AttrOr("src", "")
			if want := "/assets/heroes/" + strings.TrimSuffix(tc.hero, ".png") + "-1200x675."; !strings.HasPrefix(src, want) {
				t.Errorf("img src = %q, want the sourced hero's variant %s…", src, want)
			}
			px := buildtest.DecodePixels(t, out, src)
			requireColour(t, px, px.Bounds().Dx()/2, px.Bounds().Dy()/2, amber, "the sourced hero's amber")
			if _, err := os.Stat(filepath.Join(cfg.CacheDir, "hero-cards")); err == nil {
				t.Error("the build drew a card for a Piece that names its hero")
			}
		})
	}

	t.Run("named but missing", func(t *testing.T) {
		e := buildtest.ValidExplainer()
		e.Kind = "case-study"
		e.HeroImage = nil

		out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

		buildtest.RequireErrorRule(t, res, raftFile, "hero-missing", "raft-leader-election.png not found")
		if buildtest.Exists(out, raftURL) {
			t.Error("page was written for a Case Study whose named hero is missing")
		}
	})
}

package builder_test

import (
	"fmt"
	"image/color"
	"io/fs"
	"maps"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"hastcode.com/site/internal/buildtest"
)

const raftPage = "/distributed-systems/raft-leader-election/"

func TestHero_PictureOffersAVIFAndWebPAtFourWidths(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree()) // 1600x900 hero
	picture := buildtest.Page(t, out, raftPage).Find("article picture")

	for _, format := range []string{"avif", "webp"} {
		source := picture.Find(`source[type="image/` + format + `"]`)
		if source.Length() != 1 {
			t.Fatalf("<source type=image/%s>: found %d, want 1", format, source.Length())
		}
		entries := buildtest.Srcset(t, source.AttrOr("srcset", ""))
		var widths []int
		for _, e := range entries {
			widths = append(widths, e.Width)
			img := buildtest.DecodeImage(t, out, e.URL)
			if img.Format != format {
				t.Errorf("%s decodes as %s, want %s", e.URL, img.Format, format)
			}
			if img.Width != e.Width || img.Height != e.Width*9/16 {
				t.Errorf("%s is %dx%d, want %dx%d", e.URL, img.Width, img.Height, e.Width, e.Width*9/16)
			}
		}
		if want := []int{480, 800, 1200, 1600}; !slices.Equal(widths, want) {
			t.Errorf("%s srcset widths = %v, want %v", format, widths, want)
		}
	}
}

func TestHero_CropsForStructuredDataAreJPEGAt1200PxWideIn16x9And4x3And1x1(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree()) // 1600x900 hero

	found := map[[2]int]bool{}
	for _, img := range buildtest.Images(t, out) {
		if img.Format == "jpeg" {
			found[[2]int{img.Width, img.Height}] = true
		}
	}
	for _, want := range [][2]int{{1200, 675}, {1200, 900}, {1200, 1200}} {
		if !found[want] {
			t.Errorf("no %dx%d JPEG crop in output; JPEGs found: %v", want[0], want[1], found)
		}
	}
}

func TestHero_PictureHasJPEGFallbackWithSizeAndAltText(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	img := buildtest.Page(t, out, raftPage).Find("article picture img")
	if img.Length() != 1 {
		t.Fatalf("<picture><img>: found %d, want 1", img.Length())
	}

	if got := img.AttrOr("alt", ""); got != "Five Raft servers voting for a leader" {
		t.Errorf("alt = %q, want hero_alt from frontmatter", got)
	}
	width, height := img.AttrOr("width", ""), img.AttrOr("height", "")
	if width != "1200" || height != "675" {
		t.Errorf("width x height = %s x %s, want 1200 x 675", width, height)
	}
	src := buildtest.DecodeImage(t, out, img.AttrOr("src", ""))
	if src.Format != "jpeg" || src.Width != 1200 || src.Height != 675 {
		t.Errorf("img src is %s %dx%d, want jpeg 1200x675", src.Format, src.Width, src.Height)
	}
}

func TestHero_IsPreloadedFromTheHead(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	doc := buildtest.Page(t, out, raftPage)

	preload := doc.Find(`head link[rel="preload"][as="image"]`)
	if preload.Length() != 1 {
		t.Fatalf("hero preload links: found %d, want 1", preload.Length())
	}
	avif := doc.Find(`article picture source[type="image/avif"]`)
	if got, want := preload.AttrOr("imagesrcset", ""), avif.AttrOr("srcset", ""); got == "" || got != want {
		t.Errorf("preload imagesrcset = %q, want the AVIF srcset %q", got, want)
	}
	if got, want := preload.AttrOr("imagesizes", ""), avif.AttrOr("sizes", ""); got == "" || got != want {
		t.Errorf("preload imagesizes = %q, want the picture's sizes %q", got, want)
	}
}

const raftFile = "explainers/distributed-systems/raft-leader-election.md"

func TestHero_NarrowerThan1200PxIsABuildError(t *testing.T) {
	e := buildtest.ValidExplainer()
	e.HeroImage = buildtest.HeroPNG(1199, 800)

	out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

	buildtest.RequireError(t, res, raftFile, "1199 px wide")
	if buildtest.Exists(out, raftPage) {
		t.Error("page was written for an Explainer with a narrow hero")
	}
}

func TestHero_MissingFileIsABuildError(t *testing.T) {
	e := buildtest.ValidExplainer()
	e.HeroImage = nil // frontmatter names the hero, but no file beside the Markdown

	out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

	buildtest.RequireError(t, res, raftFile, "raft-leader-election.png not found")
	if buildtest.Exists(out, raftPage) {
		t.Error("page was written for an Explainer whose hero file is missing")
	}
}

func TestHero_MissingFromFrontmatterIsABuildError(t *testing.T) {
	e := buildtest.ValidExplainer()
	e.Hero = ""

	out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

	buildtest.RequireError(t, res, raftFile, "no hero")
	if buildtest.Exists(out, raftPage) {
		t.Error("page was written for an Explainer without a hero")
	}
}

func TestHero_FileThatIsNotAnImageIsABuildError(t *testing.T) {
	e := buildtest.ValidExplainer()
	e.HeroImage = []byte("not a png")

	out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

	buildtest.RequireError(t, res, raftFile, "not a supported raster image")
	if buildtest.Exists(out, raftPage) {
		t.Error("page was written for an Explainer whose hero does not decode")
	}
}

func TestHero_VariantsAreCachedInTheConfiguredDirectory(t *testing.T) {
	cfg := buildtest.Config()
	cfg.Images.CacheDir = filepath.Join(t.TempDir(), "image-cache")

	buildtest.MustBuild(t, cfg, smallHeroTree())

	if files := cacheSnapshot(t, cfg.Images.CacheDir); len(files) == 0 {
		t.Fatalf("image cache %s is empty after a build", cfg.Images.CacheDir)
	}
}

func TestHero_SecondBuildOfAnUnchangedHeroDoesNotReEncode(t *testing.T) {
	cfg := buildtest.Config()
	cfg.Images.CacheDir = t.TempDir()

	first := buildtest.MustBuild(t, cfg, smallHeroTree())
	before := cacheSnapshot(t, cfg.Images.CacheDir)
	second := buildtest.MustBuild(t, cfg, smallHeroTree()) // fresh content and output dirs
	after := cacheSnapshot(t, cfg.Images.CacheDir)

	if !maps.Equal(before, after) {
		t.Errorf("image cache changed on a rebuild with the same hero (re-encoded):\nbefore %v\nafter  %v", before, after)
	}
	if a, b := heroSrcsets(t, first), heroSrcsets(t, second); a != b {
		t.Errorf("hero URLs changed between builds of the same hero:\n%s\n%s", a, b)
	}
}

func TestHero_FileNamesChangeWhenTheHeroChanges(t *testing.T) {
	e := buildtest.ValidExplainer()
	teal := heroURLs(t, buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e)))
	e.HeroImage = buildtest.HeroPNGFilled(1600, 900, color.RGBA{R: 0xF5, G: 0x9E, B: 0x0B, A: 0xFF})
	amber := heroURLs(t, buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e)))

	if len(teal) == 0 || len(teal) != len(amber) {
		t.Fatalf("hero URL counts: %d before, %d after", len(teal), len(amber))
	}
	for i := range teal {
		if teal[i] == amber[i] {
			t.Errorf("hero URL %s did not change when the hero image changed", teal[i])
		}
	}
}

// smallHeroTree is the valid tree with a 1200x675 hero: the smallest the
// gate allows, so cold encodes stay quick.
func smallHeroTree() buildtest.Tree {
	e := buildtest.ValidExplainer()
	e.HeroImage = buildtest.HeroPNG(1200, 675)
	return buildtest.NewTree().WithExplainer(e)
}

// cacheSnapshot maps every file under dir to its size and modification time.
func cacheSnapshot(t *testing.T, dir string) map[string]string {
	t.Helper()
	files := map[string]string{}
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		files[p] = fmt.Sprintf("%d bytes @ %s", info.Size(), info.ModTime().Format(time.RFC3339Nano))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

// heroSrcsets is every srcset in the Explainer's <picture>, joined.
func heroSrcsets(t *testing.T, out string) string {
	t.Helper()
	var all []string
	sel := buildtest.Page(t, out, raftPage).Find("article picture source, article picture img")
	for i := range sel.Length() {
		all = append(all, sel.Eq(i).AttrOr("srcset", ""))
	}
	return strings.Join(all, "\n")
}

// heroURLs is every hero file URL the Explainer's <picture> references.
func heroURLs(t *testing.T, out string) []string {
	t.Helper()
	var urls []string
	for _, line := range strings.Split(heroSrcsets(t, out), "\n") {
		for _, e := range buildtest.Srcset(t, line) {
			urls = append(urls, e.URL)
		}
	}
	return urls
}

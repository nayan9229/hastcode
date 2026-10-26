package builder_test

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"maps"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gen2brain/avif"
	"golang.org/x/image/draw"
	"hastcode.com/site/internal/buildtest"
)

const raftPage = "/distributed-systems/raft-leader-election/"

func TestHero_PictureOffersAVIFAndJPEGAtFiveWidths(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree()) // 1600x900 hero
	picture := buildtest.Page(t, out, raftPage).Find("article picture")

	assertHeroSrcset(t, out, picture.Find(`source[type="image/avif"]`), "avif", []int{480, 800, 1000, 1200, 1472})
	assertHeroSrcset(t, out, picture.Find("img"), "jpeg", []int{480, 800, 1000, 1200, 1472})
}

func TestHero_NarrowSourceIsNeverUpscaled(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), smallHeroTree()) // 1200x675 hero
	picture := buildtest.Page(t, out, raftPage).Find("article picture")

	assertHeroSrcset(t, out, picture.Find(`source[type="image/avif"]`), "avif", []int{480, 800, 1000, 1200})
	assertHeroSrcset(t, out, picture.Find("img"), "jpeg", []int{480, 800, 1000, 1200})
}

func TestHero_NoWebPIsWrittenOrOffered(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())

	for _, img := range buildtest.Images(t, out) {
		if img.Format == "webp" || strings.HasSuffix(img.URL, ".webp") {
			t.Errorf("WebP file written: %s", img.URL)
		}
	}
	for _, page := range []string{raftPage, "/"} {
		if n := buildtest.Page(t, out, page).Find(`source[type="image/webp"]`).Length(); n != 0 {
			t.Errorf("%s has %d <source type=image/webp>, want none", page, n)
		}
	}
}

func TestHero_AVIFKeepsFullChroma(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	avif := buildtest.Page(t, out, raftPage).Find(`article picture source[type="image/avif"]`)

	for _, e := range buildtest.Srcset(t, avif.AttrOr("srcset", "")) {
		if got := buildtest.AVIFChroma(t, out, e.URL); got != "4:4:4" {
			t.Errorf("%s chroma = %s, want 4:4:4", e.URL, got)
		}
	}
}

// Spec: AVIF quality 60 below 1200 px, 55 at 1200 px and above. AVIF files
// do not record their quality, so each variant is measured against a
// reference Q60 encode of the same fixture at the same size: the 1000 px
// variant matches it, the 1200 px one is clearly smaller.
func TestHero_AVIFQualityDropsAt1200Px(t *testing.T) {
	e := buildtest.ValidExplainer()
	e.HeroImage = buildtest.HeroPNGTextured(1200, 675, 24)
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))
	avif := buildtest.Page(t, out, raftPage).Find(`article picture source[type="image/avif"]`)

	ratio := map[int]float64{}
	for _, entry := range buildtest.Srcset(t, avif.AttrOr("srcset", "")) {
		got := buildtest.FileSize(t, out, entry.URL)
		img := buildtest.DecodeImage(t, out, entry.URL)
		ref := referenceAVIFSize(t, e.HeroImage, img.Width, img.Height, 60)
		ratio[entry.Width] = float64(got) / float64(ref)
		t.Logf("%d px: %d bytes, Q60 reference %d bytes (%.2f)", entry.Width, got, ref, ratio[entry.Width])
	}
	if r := ratio[1000]; r < 0.9 || r > 1.1 {
		t.Errorf("1000 px AVIF is %.2fx a Q60 encode, want about 1 (quality 60 below 1200 px)", r)
	}
	if r := ratio[1200]; r > 0.92*ratio[1000] {
		t.Errorf("1200 px AVIF is %.2fx a Q60 encode (1000 px: %.2fx), want clearly smaller (quality 55 at 1200 px)", r, ratio[1000])
	}
}

func TestHero_HeavyHeroIsAWeightWarningNamingThePiece(t *testing.T) {
	e := buildtest.ValidExplainer()
	e.HeroImage = buildtest.HeroPNGTextured(1200, 675, 255) // static: a huge AVIF
	out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

	buildtest.RequireWarningRule(t, res, raftFile, "hero-weight", "raft-leader-election.png")
	if res.Failed() {
		t.Errorf("a heavy hero failed the build: %v", res.Errors)
	}
	if !buildtest.Exists(out, raftPage) {
		t.Error("page was not written for an Explainer with a heavy hero")
	}
}

func TestHero_NormalHeroRaisesNoWeightWarning(t *testing.T) {
	_, res := buildtest.Build(t, buildtest.Config(), buildtest.ValidTree())

	buildtest.RequireNoProblemRule(t, res, raftFile, "hero-weight")
}

// referenceAVIFSize scales a PNG source to w x h and returns the byte
// length of its AVIF encode at quality q (4:4:4, the fast test speed).
func referenceAVIFSize(t *testing.T, src []byte, w, h, q int) int {
	t.Helper()
	img, err := png.Decode(bytes.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Src, nil)
	var buf bytes.Buffer
	if err := avif.Encode(&buf, dst, avif.Options{Quality: q, Speed: 10, ChromaSubsampling: image.YCbCrSubsampleRatio444}); err != nil {
		t.Fatal(err)
	}
	return buf.Len()
}

// assertHeroSrcset checks that sel (a <source> or the <img>) has exactly
// one element whose srcset lists the given widths, each a file of format
// that decodes at that width and the 16:9 height.
func assertHeroSrcset(t *testing.T, out string, sel *goquery.Selection, format string, want []int) {
	t.Helper()
	if sel.Length() != 1 {
		t.Fatalf("%s srcset element: found %d, want 1", format, sel.Length())
	}
	var widths []int
	for _, e := range buildtest.Srcset(t, sel.AttrOr("srcset", "")) {
		widths = append(widths, e.Width)
		img := buildtest.DecodeImage(t, out, e.URL)
		if img.Format != format {
			t.Errorf("%s decodes as %s, want %s", e.URL, img.Format, format)
		}
		if h := (e.Width*9 + 8) / 16; img.Width != e.Width || img.Height != h {
			t.Errorf("%s is %dx%d, want %dx%d", e.URL, img.Width, img.Height, e.Width, h)
		}
	}
	if !slices.Equal(widths, want) {
		t.Errorf("%s srcset widths = %v, want %v", format, widths, want)
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

// Spec: the cache key carries the encoder settings and their version, so a
// retune re-encodes every cached hero exactly once. Fast vs full settings
// stand in for a settings change at the seam.
func TestHero_ChangedEncoderSettingsReEncodeOnceThenHitTheCache(t *testing.T) {
	cfg := buildtest.Config()
	cfg.Images.CacheDir = t.TempDir()

	buildtest.MustBuild(t, cfg, smallHeroTree())
	old := cacheSnapshot(t, cfg.Images.CacheDir)

	cfg.Images.Fast = false
	retuned := buildtest.MustBuild(t, cfg, smallHeroTree())
	afterRetune := cacheSnapshot(t, cfg.Images.CacheDir)
	var added int
	for p := range afterRetune {
		if _, ok := old[p]; !ok {
			added++
		}
	}
	if added != len(old) {
		t.Errorf("retuned build added %d cache files, want %d (every variant re-encoded once)", added, len(old))
	}

	again := buildtest.MustBuild(t, cfg, smallHeroTree())
	if after := cacheSnapshot(t, cfg.Images.CacheDir); !maps.Equal(afterRetune, after) {
		t.Errorf("third build changed the image cache (re-encoded again):\nbefore %v\nafter  %v", afterRetune, after)
	}
	if a, b := heroSrcsets(t, retuned), heroSrcsets(t, again); a != b {
		t.Errorf("hero URLs changed between builds with the same settings:\n%s\n%s", a, b)
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

package builder

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"hastcode.com/site/internal/images"
	"hastcode.com/site/internal/views"
)

// Hero pipeline (spec "Images and diagrams"). The Author drops one raster
// beside the Markdown and names it in `hero`. The build validates it,
// encodes every variant the page, Discover, Open Graph and JSON-LD need into
// the image cache (Config.Images.CacheDir), and copies them to
// content-hashed paths under heroAssetDir.
func init() {
	registerStep(phaseValidate, 50, "validate heroes", validateHeroes)
	registerStep(phaseTransform, 50, "encode heroes", encodeHeroes)
	registerStep(phaseRender, 15, "write hero assets", writeHeroAssets)
}

const (
	ruleHeroMissing   = "hero-missing"
	ruleHeroFormat    = "hero-format"
	ruleHeroTooNarrow = "hero-too-narrow"
)

const (
	// heroMinWidth is Google Discover's minimum for large images.
	heroMinWidth = 1200
	// heroCropWidth is the width of the 16:9, 4:3 and 1:1 crops.
	heroCropWidth = 1200
	// heroFallbackWidth is the width of the <img src> JPEG.
	heroFallbackWidth = 1200
	// heroAssetDir is where hero files are served. Names carry a hash of
	// the file's bytes, so these paths can be cached immutably.
	heroAssetDir = "/assets/heroes/"
	// heroSizes is the <img sizes> for the article column.
	heroSizes = "(min-width: 48rem) 46rem, 100vw"
)

// heroWidths are the responsive widths, each written as AVIF, WebP and JPEG.
var heroWidths = []int{480, 800, 1200, 1600}

// hero is a validated hero source and, after encoding, its view model and
// output files.
type hero struct {
	source        string // path of the raster
	width, height int

	view  views.Hero
	files map[string]string // site path -> cached file
}

// heroRef is a piece that must carry a hero (an Explainer or a Brief):
// its content file and its `hero` / `hero_alt` frontmatter.
type heroRef struct {
	file, hero, alt string
}

// heroRefs lists every Explainer and Brief; with published set, only those
// that passed every gate so far.
func (s *site) heroRefs(published bool) []heroRef {
	var refs []heroRef
	for _, e := range s.explainers {
		if !published || !s.failed[e.File] {
			refs = append(refs, heroRef{file: e.File, hero: e.Hero, alt: e.HeroAlt})
		}
	}
	for _, b := range s.briefs {
		if !published || !s.failed[b.File] {
			refs = append(refs, heroRef{file: b.File, hero: b.Hero, alt: b.HeroAlt})
		}
	}
	return refs
}

// validateHeroes checks that every Explainer and Brief names a hero that
// exists beside its Markdown, decodes, and is at least heroMinWidth wide.
func validateHeroes(s *site) error {
	for _, e := range s.heroRefs(false) {
		if e.hero == "" {
			s.fail(e.file, ruleHeroMissing, "no hero: set `hero` in frontmatter to a raster image beside the Markdown, at least %d px wide", heroMinWidth)
			continue
		}
		rel := path.Join(path.Dir(e.file), e.hero)
		src := filepath.Join(s.contentDir, filepath.FromSlash(rel))
		if _, err := os.Stat(src); err != nil {
			s.fail(e.file, ruleHeroMissing, "hero %s not found", rel)
			continue
		}
		w, h, _, err := images.Probe(src)
		if err != nil {
			s.fail(e.file, ruleHeroFormat, "hero %s: %v", rel, err)
			continue
		}
		if w < heroMinWidth {
			s.fail(e.file, ruleHeroTooNarrow, "hero %s is %d px wide; heroes must be at least %d px wide", rel, w, heroMinWidth)
			continue
		}
		s.heroes[e.file] = &hero{source: src, width: w, height: h}
	}
	return nil
}

// heroVariants is every image a hero needs, in this order: the responsive
// widths (never upscaled past the source) as AVIF, then WebP, then JPEG;
// then the 16:9, 4:3 and 1:1 crops.
func heroVariants(w, h int) []images.Variant {
	var vs []images.Variant
	for _, f := range []images.Format{images.AVIF, images.WebP, images.JPEG} {
		for _, rw := range heroWidths {
			if rw <= w {
				vs = append(vs, images.Variant{Format: f, Width: rw, Height: (rw*h + w/2) / w})
			}
		}
	}
	return append(vs,
		images.Variant{Format: images.JPEG, Width: heroCropWidth, Height: heroCropWidth * 9 / 16},
		images.Variant{Format: images.JPEG, Width: heroCropWidth, Height: heroCropWidth * 3 / 4},
		images.Variant{Format: images.JPEG, Width: heroCropWidth, Height: heroCropWidth},
	)
}

// encodeHeroes encodes the heroes of every published Explainer and Brief in
// one parallel run, reusing the image cache, and builds their view models.
func encodeHeroes(s *site) error {
	var (
		jobs   []images.Job
		owners []heroRef
	)
	for _, e := range s.heroRefs(true) {
		hr := s.heroes[e.file]
		if hr == nil {
			continue
		}
		jobs = append(jobs, images.Job{Source: hr.source, Variants: heroVariants(hr.width, hr.height)})
		owners = append(owners, e)
	}
	if len(jobs) == 0 {
		return nil
	}

	enc := &images.Encoder{CacheDir: s.cfg.Images.CacheDir, Fast: s.cfg.Images.Fast}
	results, err := enc.Run(jobs)
	if err != nil {
		return err
	}
	for i, res := range results {
		e := owners[i]
		if res.Err != nil {
			s.fail(e.file, ruleHeroFormat, "hero: %v", res.Err)
			continue
		}
		hr := s.heroes[e.file]
		hr.files = map[string]string{}
		base := strings.TrimSuffix(path.Base(e.hero), path.Ext(e.hero))
		imgs := make([]views.Image, len(res.Files))
		for j, f := range res.Files {
			url := fmt.Sprintf("%s%s-%dx%d.%s%s", heroAssetDir, base, f.Width, f.Height, f.Hash[:10], f.Format.Ext())
			hr.files[url] = f.Path
			imgs[j] = views.Image{URL: url, Width: f.Width, Height: f.Height}
		}
		// Split in heroVariants order: AVIF, WebP, JPEG widths, then crops.
		n := (len(imgs) - 3) / 3
		v := views.Hero{
			Alt:      e.alt,
			Sizes:    heroSizes,
			AVIF:     imgs[:n],
			WebP:     imgs[n : 2*n],
			JPEG:     imgs[2*n : 3*n],
			Crop16x9: imgs[3*n],
			Crop4x3:  imgs[3*n+1],
			Crop1x1:  imgs[3*n+2],
		}
		for _, im := range v.JPEG {
			if im.Width == heroFallbackWidth {
				v.Img = im
			}
		}
		hr.view = v
	}
	return nil
}

// writeHeroAssets copies encoded hero files from the cache into the output.
func writeHeroAssets(s *site) error {
	for _, e := range s.heroRefs(true) {
		hr := s.heroes[e.file]
		if hr == nil {
			continue
		}
		for url, cached := range hr.files {
			data, err := os.ReadFile(cached)
			if err != nil {
				return err
			}
			if err := s.out.WriteStatic(url, data); err != nil {
				return err
			}
		}
	}
	return nil
}

package builder

import (
	"fmt"
	"path"
	"strings"

	"hastcode.com/site/internal/images"
	"hastcode.com/site/internal/views"
)

// encodeHeroes is the hero pipeline's encode step: validateHeroes
// (step_validate_heroes.go) finds the sources, writeHeroAssets
// (step_render_hero_assets.go) copies the encoded files to the output.
func init() {
	registerStep(phaseTransform, 50, "encode heroes", encodeHeroes)
}

const (
	// heroCropWidth is the width of the 16:9, 4:3 and 1:1 crops.
	heroCropWidth = 1200
	// heroFallbackWidth is the width of the <img src> JPEG.
	heroFallbackWidth = 1200
	// heroSizes is the <img sizes> for the article column.
	heroSizes = "(min-width: 48rem) 46rem, 100vw"
)

// heroWidths are the responsive widths, each written as AVIF, WebP and JPEG.
var heroWidths = []int{480, 800, 1200, 1600}

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

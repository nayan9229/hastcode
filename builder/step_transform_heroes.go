package builder

import (
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

// heroWidths are the responsive widths, each written in every heroFormat.
var heroWidths = []int{480, 800, 1200, 1600}

// heroFormats are the responsive formats, in <picture> source order.
var heroFormats = []images.Format{images.AVIF, images.WebP, images.JPEG}

// heroCropRatios are the JPEG crops (width:height) for Discover, Open
// Graph and JSON-LD: 16:9, 4:3, 1:1, in views.Hero field order.
var heroCropRatios = [][2]int{{16, 9}, {4, 3}, {1, 1}}

// heroVariants is every image a hero needs, in this order: the responsive
// widths (never upscaled past the source) in each of heroFormats, then
// heroCropRatios. heroViewOf splits encoder output by the same layout.
func heroVariants(w, h int) []images.Variant {
	var vs []images.Variant
	for _, f := range heroFormats {
		for _, rw := range heroWidths {
			if rw <= w {
				vs = append(vs, images.Variant{Format: f, Width: rw, Height: (rw*h + w/2) / w})
			}
		}
	}
	for _, r := range heroCropRatios {
		vs = append(vs, images.Variant{Format: images.JPEG, Width: heroCropWidth, Height: heroCropWidth * r[1] / r[0]})
	}
	return vs
}

// heroViewOf builds a hero's view model from its encoded images, laid out
// as heroVariants returns them.
func heroViewOf(alt string, imgs []views.Image) views.Hero {
	nCrops := len(heroCropRatios)
	perFormat := (len(imgs) - nCrops) / len(heroFormats)
	widths := func(i int) []views.Image { return imgs[i*perFormat : (i+1)*perFormat] }
	crops := imgs[len(imgs)-nCrops:]
	v := views.Hero{
		Alt:      alt,
		Sizes:    heroSizes,
		AVIF:     widths(0),
		WebP:     widths(1),
		JPEG:     widths(2),
		Crop16x9: crops[0],
		Crop4x3:  crops[1],
		Crop1x1:  crops[2],
	}
	for _, im := range v.JPEG {
		if im.Width == heroFallbackWidth {
			v.Img = im
		}
	}
	return v
}

// encodeHeroes encodes the heroes of every published Explainer and Brief in
// one parallel run, reusing the image cache, and builds their view models.
func encodeHeroes(s *site) error {
	var (
		jobs   []images.Job
		owners []heroRef
	)
	for _, ref := range s.publishedHeroRefs() {
		hr := s.heroes[ref.file]
		if hr == nil {
			continue
		}
		jobs = append(jobs, images.Job{Source: hr.source, Variants: heroVariants(hr.width, hr.height)})
		owners = append(owners, ref)
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
		ref := owners[i]
		if res.Err != nil {
			s.fail(ref.file, ruleHeroFormat, "hero: %v", res.Err)
			continue
		}
		hr := s.heroes[ref.file]
		hr.files = map[string]string{}
		imgs := make([]views.Image, len(res.Files))
		for j, f := range res.Files {
			url := encodedImageURL(heroAssetDir, ref.hero, f)
			hr.files[url] = f.Path
			imgs[j] = views.Image{URL: url, Width: f.Width, Height: f.Height}
		}
		hr.view = heroViewOf(ref.alt, imgs)
	}
	return nil
}

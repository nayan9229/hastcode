package builder

import (
	"hastcode.com/site/internal/images"
	"hastcode.com/site/internal/views"
)

// encodeHeroes is the hero pipeline's encode step: validateHeroes
// (step_validate_heroes.go) finds the sources, writeHeroAssets
// (step_render_hero_assets.go) publishes the encoded files through the
// Image Store.
func init() {
	registerStep(phaseTransform, 50, "encode heroes", encodeHeroes)
}

const (
	// heroFallbackWidth is the width of the <img src> JPEG.
	heroFallbackWidth = 1200
	// heroSizes is the <img sizes> for the article column.
	heroSizes = "(min-width: 48rem) 46rem, 100vw"

	// AVIF quality tiers: the 1200 px variant is the mobile LCP image on
	// most phones, so it and wider ones trade a little quality for bytes
	// on slow 4G (docs/research/image-optimisation-targets.md).
	heroAVIFQuality      = 60
	heroAVIFQualityLarge = 55
	heroAVIFLargeFrom    = 1200

	// heroWeightBudget is the most the 1200 px AVIF (the mobile LCP
	// image) should weigh, in bytes (decimal KB: 150 KB); over it the build
	// warns (rule hero-weight).
	heroWeightBudget = 150_000
	heroWeightWidth  = 1200
)

const ruleHeroWeight = "hero-weight"

// heroWidths are the responsive widths, each written in every heroFormat.
var heroWidths = []int{480, 800, 1000, 1200, 1472}

// heroFormats are the responsive formats, in <picture> source order.
var heroFormats = []images.Format{images.AVIF, images.JPEG}

// cropSize is a hero crop's pixel size.
type cropSize struct{ width, height int }

// heroCrops are the 1200 px wide JPEG crops, in views.Hero field order:
// 16:9, 4:3 and 1:1 for Discover and JSON-LD, then 1200x630 for Open Graph
// and the Twitter card.
var heroCrops = []cropSize{{1200, 675}, {1200, 900}, {1200, 1200}, {1200, 630}}

// heroVariants is every image a hero needs, in this order: the responsive
// widths (never upscaled past the source) in each of heroFormats, then
// heroCrops. heroViewOf splits encoder output by the same layout.
func heroVariants(w, h int) []images.Variant {
	var vs []images.Variant
	for _, f := range heroFormats {
		for _, rw := range heroWidths {
			if rw <= w {
				vs = append(vs, images.Variant{Format: f, Width: rw, Height: (rw*h + w/2) / w, Quality: heroQuality(f, rw)})
			}
		}
	}
	for _, c := range heroCrops {
		vs = append(vs, images.Variant{Format: images.JPEG, Width: c.width, Height: c.height})
	}
	return vs
}

// heroQuality is the encoder quality of a responsive variant: the AVIF
// tiers, or 0 (the encoder's default) for JPEG.
func heroQuality(f images.Format, width int) int {
	switch {
	case f != images.AVIF:
		return 0
	case width >= heroAVIFLargeFrom:
		return heroAVIFQualityLarge
	default:
		return heroAVIFQuality
	}
}

// heroViewOf builds a hero's view model from its encoded images, laid out
// as heroVariants returns them.
func heroViewOf(alt string, imgs []views.Image) views.Hero {
	nCrops := len(heroCrops)
	perFormat := (len(imgs) - nCrops) / len(heroFormats)
	widths := func(i int) []views.Image { return imgs[i*perFormat : (i+1)*perFormat] }
	crops := imgs[len(imgs)-nCrops:]
	v := views.Hero{
		Alt:      alt,
		Sizes:    heroSizes,
		AVIF:     widths(0),
		JPEG:     widths(1),
		Crop16x9: crops[0],
		Crop4x3:  crops[1],
		Crop1x1:  crops[2],
		CropOG:   crops[3],
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
			sitePath := encodedImageURL(heroAssetDir, hr.name, f)
			hr.files[sitePath] = f.Path
			imgs[j] = views.Image{URL: s.imageURL(sitePath), Width: f.Width, Height: f.Height}
		}
		hr.view = heroViewOf(hr.alt, imgs)
		for _, f := range res.Files {
			if f.Format == images.AVIF && f.Width == heroWeightWidth && f.Size > heroWeightBudget {
				s.warn(ref.file, ruleHeroWeight, "hero %s: the %d px AVIF is %d KB, over the %d KB budget for the LCP image; re-export it smaller or simpler",
					hr.name, heroWeightWidth, f.Size/1000, heroWeightBudget/1000)
			}
		}
	}
	return nil
}

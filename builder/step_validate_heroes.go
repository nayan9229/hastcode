package builder

import (
	"os"
	"path"
	"path/filepath"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/images"
	"hastcode.com/site/internal/views"
)

// Hero pipeline (spec "Images and diagrams"). The Author drops one raster
// beside the Markdown and names it in `hero`. The build validates it (this
// step), encodes every variant the page, Discover, Open Graph and JSON-LD
// need into the image cache (step_transform_heroes.go), and copies them to
// content-hashed paths under heroAssetDir (step_render_hero_assets.go).
func init() {
	registerStep(phaseValidate, 50, "validate heroes", validateHeroes)
}

const (
	ruleHeroMissing   = "hero-missing"
	ruleHeroFormat    = "hero-format"
	ruleHeroTooNarrow = "hero-too-narrow"
)

// heroMinWidth is Google Discover's minimum for large images.
const heroMinWidth = 1200

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

// heroRefsOf lists the given Explainers and Briefs as heroRefs.
func heroRefsOf(es []*content.Explainer, bs []*content.Brief) []heroRef {
	var refs []heroRef
	for _, e := range es {
		refs = append(refs, heroRef{file: e.File, hero: e.Hero, alt: e.HeroAlt})
	}
	for _, b := range bs {
		refs = append(refs, heroRef{file: b.File, hero: b.Hero, alt: b.HeroAlt})
	}
	return refs
}

// publishedHeroRefs are the heroRefs of the pieces that passed every gate
// so far.
func (s *site) publishedHeroRefs() []heroRef {
	return heroRefsOf(s.publishedExplainers(), s.publishedBriefs())
}

// heroView is the encoded hero of a published piece's content file (every
// published Explainer and Brief has one, see validateHeroes), or the zero
// Hero.
func (s *site) heroView(file string) views.Hero {
	if hr := s.heroes[file]; hr != nil {
		return hr.view
	}
	return views.Hero{}
}

// validateHeroes checks that every Explainer and Brief names a hero that
// exists beside its Markdown, decodes, and is at least heroMinWidth wide.
func validateHeroes(s *site) error {
	for _, ref := range heroRefsOf(s.explainers, s.briefs) {
		if ref.hero == "" {
			s.fail(ref.file, ruleHeroMissing, "no hero: set `hero` in frontmatter to a raster image beside the Markdown, at least %d px wide", heroMinWidth)
			continue
		}
		rel := path.Join(path.Dir(ref.file), ref.hero)
		src := filepath.Join(s.contentDir, filepath.FromSlash(rel))
		if _, err := os.Stat(src); err != nil {
			s.fail(ref.file, ruleHeroMissing, "hero %s not found", rel)
			continue
		}
		w, h, _, err := images.Probe(src)
		if err != nil {
			s.fail(ref.file, ruleHeroFormat, "hero %s: %v", rel, err)
			continue
		}
		if w < heroMinWidth {
			s.fail(ref.file, ruleHeroTooNarrow, "hero %s is %d px wide; heroes must be at least %d px wide", rel, w, heroMinWidth)
			continue
		}
		s.heroes[ref.file] = &hero{source: src, width: w, height: h}
	}
	return nil
}

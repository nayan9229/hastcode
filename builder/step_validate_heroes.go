package builder

import (
	"fmt"
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
// need into the image cache (step_transform_heroes.go), and publishes them
// at content-hashed paths under heroAssetDir through the Image Store
// (step_render_hero_assets.go, image_store.go). A Repo Profile without a
// `hero` gets a card drawn from its Fact Pack instead, registered in the
// load phase (step_load_hero_cards.go).
func init() {
	registerStep(phaseValidate, 50, "validate heroes", validateHeroes)
}

const (
	ruleHeroMissing   = "hero-missing"
	ruleHeroFormat    = "hero-format"
	ruleHeroTooNarrow = "hero-too-narrow"
	// ruleHeroColour warns on a non-sRGB ICC profile (nothing is converted).
	ruleHeroColour = "hero-colour"
)

// profileName quotes an ICC profile description for a message.
func profileName(desc string) string {
	if desc == "" {
		return "an unnamed"
	}
	return fmt.Sprintf("%q", desc)
}

// heroMinWidth is Google Discover's minimum for large images.
const heroMinWidth = 1200

// hero is a validated hero source and, after encoding, its view model and
// output files.
type hero struct {
	source        string // path of the raster
	width, height int
	// name is the file name the encoded variants are named after (the
	// `hero` file, or heroCardName for a drawn card); alt is the alt text
	// (`hero_alt`, or the card's).
	name, alt string

	view  views.Hero
	files map[string]string // site path -> cached file
}

// heroRef is a piece that must carry a hero (an Explainer or a Brief):
// its content file and its `hero` / `hero_alt` frontmatter. card marks a
// Repo Profile without `hero`, whose hero is drawn from its Fact Pack
// (step_load_hero_cards.go) rather than sourced.
type heroRef struct {
	file, hero, alt string
	card            bool
}

// heroRefsOf lists the given Explainers and Briefs as heroRefs.
func heroRefsOf(es []*content.Explainer, bs []*content.Brief) []heroRef {
	var refs []heroRef
	for _, e := range es {
		refs = append(refs, heroRef{file: e.File, hero: e.Hero, alt: e.HeroAlt, card: e.Kind == content.RepoProfile && e.Hero == ""})
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
// A Repo Profile without `hero` is skipped: its card was drawn and
// registered already, or its Fact Pack failed and so did the Profile.
func validateHeroes(s *site) error {
	for _, ref := range heroRefsOf(s.explainers, s.briefs) {
		if ref.card {
			continue
		}
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
		info, err := images.Inspect(src)
		if err != nil {
			s.fail(ref.file, ruleHeroFormat, "hero %s: %v", rel, err)
			continue
		}
		w, h := info.Width, info.Height
		if !info.SRGB {
			s.warn(ref.file, ruleHeroColour, "hero %s carries the %s colour profile, not sRGB; the build ignores profiles, so it will look desaturated or hue-shifted: export it as sRGB", rel, profileName(info.Profile))
		}
		if w < heroMinWidth {
			s.fail(ref.file, ruleHeroTooNarrow, "hero %s is %d px wide; heroes must be at least %d px wide", rel, w, heroMinWidth)
			continue
		}
		s.heroes[ref.file] = &hero{source: src, width: w, height: h, name: ref.hero, alt: ref.alt}
	}
	return nil
}

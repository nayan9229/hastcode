package builder

import (
	"bytes"
	"cmp"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/images"
	"hastcode.com/site/internal/views"
)

// Generated heroes (spec "Page composition and gates": generated hero).
// A Piece that names no `hero` gets one drawn by the build: a
// Repo Profile the hero card, from its Fact Pack (views.HeroCardPNG); a
// Case Study, Concept, Comparison or Brief the title card, from its title,
// kind and Vertical (views.TitleCardPNG). The card is drawn, cached under
// the build cache and registered as the Piece's hero source in the load
// phase, after the Fact Packs (22) and the draft exclusion (30) and so
// before every validate step, so validate heroes, encode and write treat
// it exactly like a sourced file: the responsive widths in AVIF and JPEG,
// the JSON-LD and Open Graph crops, the preload, the weight warning and
// the image cache, which is keyed on the card's bytes: an unchanged Piece
// hits it, a changed title (or star count, or release) renames every
// variant. A Piece of any kind that names a `hero` keeps it and skips the
// card.
func init() {
	registerStep(phaseLoad, 35, "draw hero cards", drawHeroCards)
}

// heroCardCacheDir is the directory under Config.CacheDir the drawn cards
// are written to, mirroring the content path
// (hero-cards/explainers/<vertical>/<slug>.png, hero-cards/briefs/<slug>.png).
const heroCardCacheDir = "hero-cards"

// drawHeroCards draws the card of every Explainer and Brief without a
// `hero` (a Repo Profile's only when its Fact Pack loaded) and registers
// it in s.heroes.
func drawHeroCards(s *site) error {
	for _, e := range s.explainers {
		if e.Hero != "" || s.failed[e.File] {
			continue
		}
		if e.Kind != content.RepoProfile {
			c := views.TitleCard{Title: e.Title, Kind: kindOf(e.Kind).name, Vertical: s.verticals[e.Vertical].Name}
			if err := s.registerHeroCard(e.File, e.Slug, views.TitleCardPNG(c), c.Alt()); err != nil {
				return err
			}
			continue
		}
		fp := s.factPacks[e.File]
		if fp == nil {
			continue
		}
		c := heroCardOf(e, fp)
		if err := s.registerHeroCard(e.File, e.Slug, views.HeroCardPNG(c), c.Alt()); err != nil {
			return err
		}
	}
	for _, b := range s.briefs {
		if b.Hero != "" || s.failed[b.File] {
			continue
		}
		c := views.TitleCard{Title: b.Title, Kind: briefLabel}
		if err := s.registerHeroCard(b.File, b.Slug, views.TitleCardPNG(c), c.Alt()); err != nil {
			return err
		}
	}
	return nil
}

// registerHeroCard writes a drawn card under the build cache (unless it
// already holds those bytes) and registers it as file's hero source, its
// variants named after slug (heroCardName).
func (s *site) registerHeroCard(file, slug string, png []byte, alt string) error {
	src := filepath.Join(s.cfg.CacheDir, heroCardCacheDir, filepath.FromSlash(strings.TrimSuffix(file, path.Ext(file))+".png"))
	if err := writeIfChanged(src, png); err != nil {
		return err
	}
	info, err := images.Inspect(src)
	if err != nil {
		return fmt.Errorf("hero card of %s: %w", file, err)
	}
	s.heroes[file] = &hero{source: src, width: info.Width, height: info.Height, name: heroCardName(slug), alt: alt}
	return nil
}

// heroCardName is the file name a Piece's card variants are named after
// (encodedImageURL), in place of the `hero` file name a sourced hero has.
func heroCardName(slug string) string { return slug + "-card.png" }

// heroCardOf maps a Profile and its Fact Pack to the card: the
// repository's own description (the Profile's when GitHub has none), its
// language, stars, license and latest release. A license GitHub could not
// classify is left off, as is anything the Fact Pack lacks.
func heroCardOf(e *content.Explainer, fp *content.FactPack) views.HeroCard {
	c := views.HeroCard{
		Repo:        fp.FullName,
		Description: cmp.Or(fp.Description, e.Description),
		Language:    fp.Language,
		Stars:       fp.Stars,
	}
	if l := fp.License; l != nil && !l.NeedsReview {
		c.License = l.SPDXID
	}
	if r := fp.LatestRelease; r != nil {
		c.Release = r.Tag
	}
	return c
}

// writeIfChanged writes data to p unless p already holds it, through a
// temp file and a rename so a concurrent build never reads a partial card.
func writeIfChanged(p string, data []byte) error {
	if old, err := os.ReadFile(p); err == nil && bytes.Equal(old, data) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(p), ".tmp-*")
	if err != nil {
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return os.Rename(tmp.Name(), p)
}

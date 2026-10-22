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

// Hero cards (spec "Page composition and gates": generated hero). A Repo
// Profile that names no `hero` gets one drawn from its Fact Pack
// (views.HeroCardPNG). The card is written under the build cache and
// registered as the Profile's hero source here, before validate heroes
// (50), so the encode and write steps treat it exactly like a sourced
// file: the responsive widths in AVIF and JPEG, the JSON-LD and Open Graph
// crops, the preload, the weight warning and the image cache, which is
// keyed on the card's bytes: an unchanged Fact Pack hits it, a changed
// star count or release renames every variant. A Profile that names a
// `hero` keeps it and skips the card; every other kind still needs one.
func init() {
	registerStep(phaseValidate, 45, "draw hero cards", drawHeroCards)
}

// heroCardCacheDir is the directory under Config.CacheDir the drawn cards
// are written to, mirroring the content path
// (hero-cards/explainers/<vertical>/<slug>.png).
const heroCardCacheDir = "hero-cards"

// drawHeroCards draws the card of every Repo Profile without a `hero`
// whose Fact Pack loaded, and registers it in s.heroes.
func drawHeroCards(s *site) error {
	for _, e := range s.explainers {
		if e.Kind != content.RepoProfile || e.Hero != "" || s.failed[e.File] {
			continue
		}
		fp := s.factPacks[e.File]
		if fp == nil {
			continue
		}
		card := heroCardOf(e, fp)
		src := filepath.Join(s.cfg.CacheDir, heroCardCacheDir, filepath.FromSlash(strings.TrimSuffix(e.File, path.Ext(e.File))+".png"))
		if err := writeIfChanged(src, views.HeroCardPNG(card)); err != nil {
			return err
		}
		info, err := images.Inspect(src)
		if err != nil {
			return fmt.Errorf("hero card of %s: %w", e.File, err)
		}
		s.heroes[e.File] = &hero{source: src, width: info.Width, height: info.Height, name: heroCardName(e), alt: card.Alt()}
	}
	return nil
}

// heroCardName is the file name a Profile's card variants are named after
// (encodedImageURL), in place of the `hero` file name a sourced hero has.
func heroCardName(e *content.Explainer) string { return e.Slug + "-card.png" }

// heroCardOf maps a Profile and its Fact Pack to the card: the
// repository's own description (the Profile's when GitHub has none), its
// language, stars, license and latest release. A license GitHub could not
// classify is left off, as is anything the Fact Pack lacks.
func heroCardOf(e *content.Explainer, fp *content.FactPack) views.HeroCard {
	c := views.HeroCard{
		Repo:        fp.Slug,
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

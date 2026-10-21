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

// Author photos. The Author drops a raster beside their Markdown and names
// it in `photo`; the build encodes one square JPEG (center crop) through
// the hero image cache and serves it at a content-hashed path. A photo is
// optional (Google News transparency only needs bylines that lead to
// author information) but its absence warns: a real face is the strongest
// sign the profile is not fabricated.
func init() {
	registerStep(phaseValidate, 55, "validate author photos", validateAuthorPhotos)
	registerStep(phaseTransform, 55, "encode author photos", encodeAuthorPhotos)
	registerStep(phaseRender, 16, "write author photos", writeAuthorPhotos)
}

const ruleAuthorPhoto = "author-photo"

const (
	// authorPhotoSize is the side of the square photo, 2x its CSS size.
	authorPhotoSize = 400
	// authorPhotoDir is where photos are served; names carry a hash of
	// the file's bytes, so these paths can be cached immutably.
	authorPhotoDir = "/assets/authors/"
)

// authorPhoto is a validated photo source and, after encoding, its view
// model and cached output file.
type authorPhoto struct {
	source string
	view   views.Image
	cached string
}

// validateAuthorPhotos warns on Authors without a photo and fails an
// Author whose photo is missing or does not decode.
func validateAuthorPhotos(s *site) error {
	s.authorPhotos = map[string]*authorPhoto{}
	for _, a := range s.sortedAuthors() {
		if a.Photo == "" {
			s.warn(a.File, ruleAuthorPhoto, "no photo: set `photo` to a square-ish raster beside the Markdown (at least %d px)", authorPhotoSize)
			continue
		}
		rel := path.Join(path.Dir(a.File), a.Photo)
		src := filepath.Join(s.contentDir, filepath.FromSlash(rel))
		if _, err := os.Stat(src); err != nil {
			s.fail(a.File, ruleAuthorPhoto, "photo %s not found", rel)
			continue
		}
		if _, _, _, err := images.Probe(src); err != nil {
			s.fail(a.File, ruleAuthorPhoto, "photo %s: %v", rel, err)
			continue
		}
		s.authorPhotos[a.Slug] = &authorPhoto{source: src}
	}
	return nil
}

// encodeAuthorPhotos encodes every valid photo as one square JPEG.
func encodeAuthorPhotos(s *site) error {
	var (
		jobs   []images.Job
		owners []string
	)
	for _, a := range s.sortedAuthors() {
		ph := s.authorPhotos[a.Slug]
		if ph == nil || s.failed[a.File] {
			continue
		}
		jobs = append(jobs, images.Job{Source: ph.source, Variants: []images.Variant{
			{Format: images.JPEG, Width: authorPhotoSize, Height: authorPhotoSize},
		}})
		owners = append(owners, a.Slug)
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
		a := s.authors[owners[i]]
		if res.Err != nil {
			s.fail(a.File, ruleAuthorPhoto, "photo: %v", res.Err)
			continue
		}
		f := res.Files[0]
		base := strings.TrimSuffix(path.Base(a.Photo), path.Ext(a.Photo))
		ph := s.authorPhotos[a.Slug]
		ph.cached = f.Path
		ph.view = views.Image{
			URL:    fmt.Sprintf("%s%s-%dx%d.%s%s", authorPhotoDir, base, f.Width, f.Height, f.Hash[:10], f.Format.Ext()),
			Width:  f.Width,
			Height: f.Height,
		}
	}
	return nil
}

// writeAuthorPhotos copies encoded photos of Authors that will get a page.
func writeAuthorPhotos(s *site) error {
	for _, a := range s.sortedAuthors() {
		ph := s.authorPhotos[a.Slug]
		if ph == nil || ph.cached == "" || s.failed[a.File] {
			continue
		}
		data, err := os.ReadFile(ph.cached)
		if err != nil {
			return err
		}
		if err := s.out.WriteStatic(ph.view.URL, data); err != nil {
			return err
		}
	}
	return nil
}

package builder

import (
	"os"
	"path"
	"path/filepath"

	"hastcode.com/site/internal/images"
	"hastcode.com/site/internal/views"
)

// Author photos. The Author drops a raster beside their Markdown and names
// it in `photo`; the build encodes one square JPEG (center crop) through
// the hero image cache (step_transform_author_photos.go) and serves it at a
// content-hashed path (step_render_author_photos.go). A photo is optional
// (Google News transparency only needs bylines that lead to author
// information) but its absence warns: a real face is the strongest sign
// the profile is not fabricated.
func init() {
	registerStep(phaseValidate, 55, "validate author photos", validateAuthorPhotos)
}

const ruleAuthorPhoto = "author-photo"

// authorPhotoSize is the side of the square photo, 2x its CSS size.
const authorPhotoSize = 400

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

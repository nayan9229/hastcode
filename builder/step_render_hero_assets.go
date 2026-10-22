package builder

import (
	"maps"
	"path"
	"slices"

	"hastcode.com/site/internal/images"
)

func init() {
	registerStep(phaseRender, 15, "write hero assets", writeHeroAssets)
}

// writeHeroAssets publishes the encoded hero files of published pieces
// through the Image Store (image_store.go): written to the output in local
// mode, queued for upload in r2 mode.
func writeHeroAssets(s *site) error {
	for _, ref := range s.publishedHeroRefs() {
		hr := s.heroes[ref.file]
		if hr == nil {
			continue
		}
		for _, sitePath := range slices.Sorted(maps.Keys(hr.files)) {
			if err := s.storeImage(sitePath, hr.files[sitePath], nil, mimeOf(sitePath)); err != nil {
				return err
			}
		}
	}
	return nil
}

// mimeOf is the Content-Type of an encoded image, from its extension.
func mimeOf(sitePath string) string {
	for _, f := range []images.Format{images.AVIF, images.JPEG} {
		if path.Ext(sitePath) == f.Ext() {
			return f.MIME()
		}
	}
	return "image/png" // the OG badge
}

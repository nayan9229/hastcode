package builder

import "hastcode.com/site/internal/images"

func init() {
	registerStep(phaseRender, 16, "write author photos", writeAuthorPhotos)
}

// writeAuthorPhotos publishes the encoded photos of Authors that will get a
// page through the Image Store.
func writeAuthorPhotos(s *site) error {
	for _, a := range s.sortedAuthors() {
		ph := s.authorPhotos[a.Slug]
		if ph == nil || ph.cached == "" || s.failed[a.File] {
			continue
		}
		if err := s.storeImage(ph.sitePath, ph.cached, nil, images.JPEG.MIME()); err != nil {
			return err
		}
	}
	return nil
}

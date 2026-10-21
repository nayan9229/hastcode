package builder

import "os"

func init() {
	registerStep(phaseRender, 16, "write author photos", writeAuthorPhotos)
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

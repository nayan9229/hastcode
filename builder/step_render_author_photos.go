package builder

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
		if err := s.copyCachedFile(ph.view.URL, ph.cached); err != nil {
			return err
		}
	}
	return nil
}

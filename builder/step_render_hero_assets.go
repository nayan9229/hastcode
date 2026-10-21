package builder

import "os"

func init() {
	registerStep(phaseRender, 15, "write hero assets", writeHeroAssets)
}

// writeHeroAssets copies encoded hero files from the cache into the output.
func writeHeroAssets(s *site) error {
	for _, ref := range s.publishedHeroRefs() {
		hr := s.heroes[ref.file]
		if hr == nil {
			continue
		}
		for url, cached := range hr.files {
			if err := s.copyCachedFile(url, cached); err != nil {
				return err
			}
		}
	}
	return nil
}

// copyCachedFile writes an encoded image from the image cache to urlPath.
func (s *site) copyCachedFile(urlPath, cached string) error {
	data, err := os.ReadFile(cached)
	if err != nil {
		return err
	}
	return s.out.WriteStatic(urlPath, data)
}

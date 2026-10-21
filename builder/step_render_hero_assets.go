package builder

import "os"

func init() {
	registerStep(phaseRender, 15, "write hero assets", writeHeroAssets)
}

// writeHeroAssets copies encoded hero files from the cache into the output.
func writeHeroAssets(s *site) error {
	for _, e := range s.heroRefs(true) {
		hr := s.heroes[e.file]
		if hr == nil {
			continue
		}
		for url, cached := range hr.files {
			data, err := os.ReadFile(cached)
			if err != nil {
				return err
			}
			if err := s.out.WriteStatic(url, data); err != nil {
				return err
			}
		}
	}
	return nil
}

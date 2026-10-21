package builder

import (
	"maps"
	"slices"

	"hastcode.com/site/internal/views"
)

func init() {
	registerStep(phaseRender, 10, "write site assets", writeSiteAssets)
}

// writeSiteAssets writes the design system's static files (stylesheet,
// SVG mark) that every page references.
func writeSiteAssets(s *site) error {
	assets := views.StaticAssets()
	for _, p := range slices.Sorted(maps.Keys(assets)) {
		if err := s.out.WriteStatic(p, assets[p]); err != nil {
			return err
		}
	}
	return nil
}

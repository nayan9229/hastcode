package builder

import (
	"maps"
	"slices"

	"hastcode.com/site/internal/views"
)

func init() {
	registerStep(phaseRender, 10, "write site assets", writeSiteAssets)
}

// writeSiteAssets writes the static files at stable paths (StaticAssets:
// the SVG mark and the crawlable PNG icon, ADR-0013; the stylesheet is
// inline since ADR-0006), and publishes the OG badge through the Image
// Store.
func writeSiteAssets(s *site) error {
	if err := s.storeImage(views.OGBadge().URL, "", views.OGBadgePNG(), "image/png"); err != nil {
		return err
	}
	assets := views.StaticAssets()
	for _, p := range slices.Sorted(maps.Keys(assets)) {
		if err := s.out.WriteStatic(p, assets[p]); err != nil {
			return err
		}
	}
	return nil
}

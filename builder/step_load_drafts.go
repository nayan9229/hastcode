package builder

import (
	"slices"

	"hastcode.com/site/internal/content"
)

func init() {
	registerStep(phaseLoad, 30, "exclude drafts from production", excludeDrafts)
}

// excludeDrafts drops `draft: true` pieces in production builds, before any
// gate runs, so a draft is absent from every page, listing, feed and
// manifest, and an unfinished draft never fails a production build. Preview
// builds keep drafts and gate them like any other piece.
func excludeDrafts(s *site) error {
	if s.cfg.Mode != Production {
		return nil
	}
	s.explainers = slices.DeleteFunc(s.explainers, func(e *content.Explainer) bool { return e.Draft })
	s.briefs = slices.DeleteFunc(s.briefs, func(b *content.Brief) bool { return b.Draft })
	return nil
}

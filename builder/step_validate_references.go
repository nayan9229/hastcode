package builder

func init() {
	registerStep(phaseValidate, 10, "validate explainer references", validateExplainerReferences)
}

const (
	ruleUnknownVertical = "unknown-vertical"
	ruleUnknownAuthor   = "unknown-author"
)

// validateExplainerReferences checks that every Explainer's Vertical is in
// the configured list, its Author has a file under authors/, and every
// `related` slug names an Explainer (otherwise its Related link would be a
// broken internal link).
func validateExplainerReferences(s *site) error {
	slugs := map[string]bool{}
	for _, e := range s.explainers {
		slugs[e.Slug] = true
	}
	for _, e := range s.explainers {
		if _, ok := s.verticals[e.Vertical]; !ok {
			s.fail(e.File, ruleUnknownVertical, "unknown Vertical %q (add it to the Verticals list in the site config)", e.Vertical)
		}
		if _, ok := s.authors[e.Author]; !ok {
			s.fail(e.File, ruleUnknownAuthor, "unknown Author %q (no authors/%s.md)", e.Author, e.Author)
		}
		for _, slug := range e.Related {
			if !slugs[slug] {
				s.fail(e.File, ruleBrokenLink, "related lists %q, but there is no such Explainer", slug)
			}
		}
	}
	return nil
}

package builder

func init() {
	registerStep(phaseValidate, 10, "validate explainer references", validateExplainerReferences)
}

const (
	ruleUnknownVertical = "unknown-vertical"
	ruleUnknownAuthor   = "unknown-author"
)

// validateExplainerReferences checks that every Explainer's Vertical is in
// the configured list and its Author has a file under authors/.
func validateExplainerReferences(s *site) error {
	for _, e := range s.explainers {
		if _, ok := s.verticals[e.Vertical]; !ok {
			s.fail(e.File, ruleUnknownVertical, "unknown Vertical %q (add it to the Verticals list in the site config)", e.Vertical)
		}
		if _, ok := s.authors[e.Author]; !ok {
			s.fail(e.File, ruleUnknownAuthor, "unknown Author %q (no authors/%s.md)", e.Author, e.Author)
		}
	}
	return nil
}

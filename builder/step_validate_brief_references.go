package builder

func init() {
	registerStep(phaseValidate, 12, "validate brief references", validateBriefReferences)
}

// validateBriefReferences checks that every Brief's Author has a file
// under authors/. (Explainer links are checked with the skeleton.)
func validateBriefReferences(s *site) error {
	for _, b := range s.briefs {
		if _, ok := s.authors[b.Author]; !ok {
			s.fail(b.File, ruleUnknownAuthor, "unknown Author %q (no authors/%s.md)", b.Author, b.Author)
		}
	}
	return nil
}

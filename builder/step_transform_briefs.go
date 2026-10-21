package builder

func init() {
	registerStep(phaseTransform, 20, "rewrite brief explainer links", rewriteBriefExplainerLinks)
}

// rewriteBriefExplainerLinks points every Explainer link in a published
// Brief at the Explainer's canonical trailing-slash URL. Broken links have
// already failed the Brief (validateBriefLinks).
func rewriteBriefExplainerLinks(s *site) error {
	for _, b := range s.publishedBriefs() {
		for _, l := range s.explainerLinks(b) {
			if l.target != nil {
				l.node.Destination = []byte(l.href())
			}
		}
	}
	return nil
}

package builder

import "hastcode.com/site/internal/content"

func init() {
	registerStep(phaseLoad, 20, "load explainers", loadExplainers)
}

func loadExplainers(s *site) error {
	return s.eachMarkdown("explainers", func(rel string, data []byte) {
		e, err := content.ParseExplainer(rel, data)
		if err != nil {
			s.fail(rel, ruleFrontmatter, "%v", err)
			return
		}
		s.explainers = append(s.explainers, e)
	})
}

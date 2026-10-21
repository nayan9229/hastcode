package builder

import "hastcode.com/site/internal/content"

func init() {
	registerStep(phaseLoad, 25, "load briefs", loadBriefs)
}

func loadBriefs(s *site) error {
	return s.eachMarkdown("briefs", func(rel string, data []byte) {
		b, err := content.ParseBrief(rel, data)
		if err != nil {
			s.fail(rel, ruleFrontmatter, "%v", err)
			return
		}
		s.briefs = append(s.briefs, b)
	})
}

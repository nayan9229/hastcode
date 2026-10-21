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

// publishedBriefs are the Briefs that passed every gate so far, in load
// order. Render and write steps must use this, not s.briefs.
func (s *site) publishedBriefs() []*content.Brief {
	var ok []*content.Brief
	for _, b := range s.briefs {
		if !s.failed[b.File] {
			ok = append(ok, b)
		}
	}
	return ok
}

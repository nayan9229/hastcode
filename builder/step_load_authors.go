package builder

import "hastcode.com/site/internal/content"

func init() {
	registerStep(phaseLoad, 10, "load authors", loadAuthors)
}

func loadAuthors(s *site) error {
	return s.eachMarkdown("authors", func(rel string, data []byte) {
		a, err := content.ParseAuthor(rel, data)
		if err != nil {
			s.fail(rel, ruleFrontmatter, "%v", err)
			return
		}
		s.authors[a.Slug] = a
	})
}

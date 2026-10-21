package builder

import (
	"io/fs"
	"os"
	"path/filepath"

	"hastcode.com/site/internal/content"
)

func init() {
	registerStep(phaseLoad, 10, "load authors", loadAuthors)
	registerStep(phaseLoad, 20, "load explainers", loadExplainers)
}

// Rule names for load-time problems.
const ruleFrontmatter = "frontmatter"

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

// eachMarkdown calls fn for every .md file under contentDir/sub (in lexical
// order) with its content-relative path. A missing sub directory is not an
// error: the content tree simply has none of that piece yet.
func (s *site) eachMarkdown(sub string, fn func(rel string, data []byte)) error {
	root := filepath.Join(s.contentDir, sub)
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil
	}
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(s.contentDir, path)
		if err != nil {
			return err
		}
		fn(filepath.ToSlash(rel), data)
		return nil
	})
}

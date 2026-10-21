package builder

import (
	"io/fs"
	"os"
	"path/filepath"
)

// Shared by the load steps (step_load_*.go): reading content files.

// ruleFrontmatter is a content file whose frontmatter does not parse or
// lacks a required field.
const ruleFrontmatter = "frontmatter"

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

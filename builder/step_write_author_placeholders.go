package builder

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func init() {
	// After every page, feed and sitemap is written (render phase, feeds
	// 20, sitemaps 30/35), before the search index (800) and the page
	// manifest (900), which a failed build skips.
	registerStep(phaseWrite, 60, "check author placeholders", checkAuthorPlaceholders)
}

// ruleAuthorPlaceholder: published output still holds a TODO(author)
// placeholder, a fact only the Author can supply (README "Content
// layout"). An error in production, a warning in preview, so drafts of the
// transparency pages can still be read on a preview.
const ruleAuthorPlaceholder = "author-placeholder"

// authorPlaceholder is the marker content uses for facts only the Author
// can supply.
const authorPlaceholder = "TODO(author)"

// checkAuthorPlaceholders reads every written page and feed and reports
// each one containing authorPlaceholder, attributed to the content file to
// fix: the page's own source when it holds a placeholder, else the first
// marked Author file (names and profiles spread to bylines, JSON-LD and
// feeds), else the first marked content file.
func checkAuthorPlaceholders(s *site) error {
	marked, err := s.markedContentFiles()
	if err != nil {
		return err
	}
	report := s.fail
	if s.cfg.Mode == Preview {
		report = s.warn
	}
	sources := s.pageSources()
	root := filepath.Join(s.out.Dir, "static")
	return filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || (filepath.Ext(p) != ".html" && filepath.Ext(p) != ".xml") {
			return err
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if !bytes.Contains(data, []byte(authorPlaceholder)) {
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		u := "/" + filepath.ToSlash(rel)
		pageURL := strings.TrimSuffix(u, "index.html")
		file := sources[pageURL]
		if !slices.Contains(marked, file) {
			file = ""
			if len(marked) > 0 {
				file = marked[0]
			}
		}
		report(file, ruleAuthorPlaceholder, "%s contains %q placeholder text; replace every %s line in %s before publishing",
			pageURL, authorPlaceholder, authorPlaceholder, strings.Join(marked, ", "))
		return nil
	})
}

// markedContentFiles are the content files holding authorPlaceholder:
// Author files first, then the rest, each group in path order.
func (s *site) markedContentFiles() ([]string, error) {
	var marked []string
	for _, sub := range []string{"authors", "explainers", "briefs", "pages"} {
		err := s.eachMarkdown(sub, func(rel string, data []byte) {
			if bytes.Contains(data, []byte(authorPlaceholder)) {
				marked = append(marked, rel)
			}
		})
		if err != nil {
			return nil, err
		}
	}
	return marked, nil
}

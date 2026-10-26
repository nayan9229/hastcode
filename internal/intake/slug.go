package intake

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// slugify makes a slug from a title: lower-case ASCII letters and digits,
// every other run of characters one hyphen, at most 80 characters, cut at
// a hyphen.
func slugify(title string) string {
	var b strings.Builder
	hyphen := false
	for _, r := range strings.ToLower(title) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			if hyphen && b.Len() > 0 {
				b.WriteByte('-')
			}
			hyphen = false
			b.WriteRune(r)
		case r == '\'' || r == '’':
			// "Netflix's" → "netflixs", not "netflix-s"
		default:
			hyphen = true
		}
	}
	slug := b.String()
	const max = 80
	if len(slug) > max {
		slug = slug[:max]
		if i := strings.LastIndexByte(slug, '-'); i > 0 {
			slug = slug[:i]
		}
	}
	return slug
}

// isSlug reports whether s has the shape slugify gives: lower-case ASCII
// letters and digits in runs joined by single hyphens.
func isSlug(s string) bool {
	if s == "" || strings.HasPrefix(s, "-") || strings.HasSuffix(s, "-") || strings.Contains(s, "--") {
		return false
	}
	for _, r := range s {
		if !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-') {
			return false
		}
	}
	return true
}

// pieceFile is the content-relative file of the Piece using slug
// (explainers/<vertical>/<slug>.md or briefs/<slug>.md; slugs are unique
// site-wide), or "" when none does. slug must be a slug (isSlug): it is
// joined into exact paths, never matched as a pattern.
func (in *Intake) pieceFile(slug string) (string, error) {
	var candidates []string
	verticals, err := os.ReadDir(filepath.Join(in.ContentDir, "explainers"))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return "", err
	}
	for _, v := range verticals {
		if v.IsDir() {
			candidates = append(candidates, path.Join("explainers", v.Name(), slug+".md"))
		}
	}
	candidates = append(candidates, path.Join("briefs", slug+".md"))
	for _, rel := range candidates {
		info, err := os.Stat(filepath.Join(in.ContentDir, filepath.FromSlash(rel)))
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return "", err
		}
		if !info.IsDir() {
			return rel, nil
		}
	}
	return "", nil
}

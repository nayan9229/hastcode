package content

import (
	"fmt"
	"strings"
)

// Repo is the GitHub repository a Repo Profile evaluates: owner/name.
// The zero Repo is "no repository" (every kind but a Repo Profile).
type Repo struct {
	Owner, Name string
}

// ParseRepo reads owner/name, GitHub's spelling of a repository: an owner
// of letters, digits and hyphens, a name of letters, digits, hyphens,
// underscores and dots. It is the one check the frontmatter `repo` key and
// the fetch command's -repo both go through.
func ParseRepo(s string) (Repo, error) {
	owner, name, ok := strings.Cut(s, "/")
	if !ok || owner == "" || name == "" || strings.Contains(name, "/") ||
		strings.ContainsFunc(owner, func(r rune) bool { return !isAlnum(r) && r != '-' }) ||
		strings.ContainsFunc(name, func(r rune) bool { return !isAlnum(r) && r != '-' && r != '_' && r != '.' }) {
		return Repo{}, fmt.Errorf("%q must be owner/name", s)
	}
	return Repo{Owner: owner, Name: name}, nil
}

// String is owner/name, or "" for the zero Repo.
func (r Repo) String() string {
	if r == (Repo{}) {
		return ""
	}
	return r.Owner + "/" + r.Name
}

func isAlnum(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9'
}

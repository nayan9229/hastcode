package builder

import "hastcode.com/site/internal/content"

// An Author is a named person or, for a team byline, an organization
// (content/authors/<slug>.md `kind`). The kind decides the schema.org type
// articles name as their author, so an unknown kind fails the build.
func init() {
	registerStep(phaseValidate, 12, "validate author kinds", validateAuthorKinds)
}

const ruleAuthorKind = "author-kind"

func validateAuthorKinds(s *site) error {
	for _, a := range s.sortedAuthors() {
		if a.Kind != content.AuthorPerson && a.Kind != content.AuthorOrganization {
			s.fail(a.File, ruleAuthorKind, "unknown kind %q: use %q or %q", a.Kind, content.AuthorPerson, content.AuthorOrganization)
		}
	}
	return nil
}

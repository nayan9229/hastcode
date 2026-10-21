package builder

import "strings"

func init() {
	registerStep(phaseValidate, 20, "validate unique slugs", validateUniqueSlugs)
}

const ruleDuplicateSlug = "duplicate-slug"

// validateUniqueSlugs fails every piece whose slug is shared with another
// piece, across Verticals and across the Explainer and Brief namespaces.
// Related links and Briefs reference pieces by slug
// alone, so a slug must be unambiguous site-wide.
func validateUniqueSlugs(s *site) error {
	files := map[string][]string{} // slug -> content files, in load order
	var order []string
	add := func(slug, file string) {
		if _, seen := files[slug]; !seen {
			order = append(order, slug)
		}
		files[slug] = append(files[slug], file)
	}
	for _, e := range s.explainers {
		add(e.Slug, e.File)
	}
	for _, b := range s.briefs {
		add(b.Slug, b.File)
	}
	for _, slug := range order {
		group := files[slug]
		if len(group) < 2 {
			continue
		}
		for _, file := range group {
			s.fail(file, ruleDuplicateSlug, "slug %q is also used by %s", slug, strings.Join(others(group, file), ", "))
		}
	}
	return nil
}

func others(all []string, except string) []string {
	var out []string
	for _, f := range all {
		if f != except {
			out = append(out, f)
		}
	}
	return out
}

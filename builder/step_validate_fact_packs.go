package builder

import (
	"time"

	"hastcode.com/site/internal/content"
)

func init() {
	registerStep(phaseValidate, 15, "validate fact packs", validateFactPacks)
}

const (
	// ruleFactPackRepo is a Fact Pack that names a different repository
	// than the Profile's `repo` frontmatter: the page would render another
	// repository's facts.
	ruleFactPackRepo = "fact-pack-repo"
	// ruleFactPackStale warns (never fails) when a Fact Pack was fetched
	// more than factPackMaxAge before the build clock: the weekly Refresh
	// has not landed, so the facts may be behind.
	ruleFactPackStale = "fact-pack-stale"
)

// factPackMaxAge is how old a Fact Pack may be before the build warns.
const factPackMaxAge = 30 * 24 * time.Hour

// validateFactPacks checks every loaded Fact Pack against its Repo Profile
// and the build clock.
func validateFactPacks(s *site) error {
	for _, e := range s.explainers {
		fp, ok := s.factPacks[e.File]
		if !ok {
			continue
		}
		file := content.FactPackFile(e.File)
		if fp.Slug != e.Repo {
			s.fail(e.File, ruleFactPackRepo, "Fact Pack %s names repository %q, but repo is %q: fetch the Fact Pack for %s", file, fp.Slug, e.Repo, e.Repo)
		}
		if age := s.cfg.Now.Sub(fp.FetchedAt); age > factPackMaxAge {
			s.warn(e.File, ruleFactPackStale, "Fact Pack %s was fetched %s, %d days ago; Refresh it (every week)",
				file, fp.FetchedAt.Format("2006-01-02"), int(age/(24*time.Hour)))
		}
	}
	return nil
}

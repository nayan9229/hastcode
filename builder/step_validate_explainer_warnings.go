package builder

import "time"

func init() {
	registerStep(phaseValidate, 40, "warn on explainer section length", warnSectionLength)
	registerStep(phaseValidate, 50, "warn on overdue explainer review", warnReviewOverdue)
}

const (
	ruleSectionLength = "section-length"
	ruleReviewOverdue = "review-overdue"

	// reviewEvery is the Review cadence (GLOSSARY "Review").
	reviewEvery = 90 * 24 * time.Hour

	// maxSectionWords is the soft cap on an H2 section's prose (code
	// excluded). Answer engines extract 120–180-word sections best.
	maxSectionWords = 250
)

// warnSectionLength warns (never fails) on any H2 section over
// maxSectionWords words.
func warnSectionLength(s *site) error {
	for _, e := range s.explainers {
		_, sections := e.Body.Sections()
		for _, sec := range sections {
			if n := e.Body.Words(sec.Blocks); n > maxSectionWords {
				s.warn(e.File, ruleSectionLength, "## %s has %d words, over the %d-word guide; split it", sec.Heading, n, maxSectionWords)
			}
		}
	}
	return nil
}

// warnReviewOverdue warns (never fails) when an Explainer's last Review is
// more than 90 days before the build clock.
func warnReviewOverdue(s *site) error {
	for _, e := range s.explainers {
		if age := s.cfg.Now.Sub(e.Reviewed); age > reviewEvery {
			s.warn(e.File, ruleReviewOverdue, "last Reviewed %s, %d days ago; Review it (every 90 days)",
				e.Reviewed.Format("2006-01-02"), int(age/(24*time.Hour)))
		}
	}
	return nil
}

package builder

import "time"

func init() {
	registerStep(phaseValidate, 50, "warn on overdue explainer review", warnReviewOverdue)
}

const ruleReviewOverdue = "review-overdue"

// reviewEvery is the Review cadence (GLOSSARY "Review").
const reviewEvery = 90 * 24 * time.Hour

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

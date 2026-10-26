package builder

import (
	"time"
)

func init() {
	registerStep(phaseValidate, 45, "check verified dates", checkVerified)
}

// ruleUnverified: a non-draft Piece published after verifiedCutoff has no
// valid `verified` date, the Author's record that every claim was checked
// against its Source (ADR-0007). An error in production, a warning in
// preview, so a Preview still shows the Piece while the Author checks it.
const ruleUnverified = "unverified"

// verifiedCutoff is the last publication day outside the rule: the ten
// Repo Profiles published before ADR-0007 keep publishing until their
// review queue is cleared.
var verifiedCutoff = time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC)

// checkVerified applies the verified gate to every Explainer and Brief.
func checkVerified(s *site) error {
	report := s.fail
	if s.cfg.Mode == Preview {
		report = s.warn
	}
	for _, e := range s.explainers {
		if e.Draft {
			continue
		}
		if msg := verifiedProblem(e.Published, e.Verified, s.cfg.Now); msg != "" {
			report(e.File, ruleUnverified, "%s", msg)
		}
	}
	for _, b := range s.briefs {
		if b.Draft {
			continue
		}
		if msg := verifiedProblem(b.Published, b.Verified, s.cfg.Now); msg != "" {
			report(b.File, ruleUnverified, "%s", msg)
		}
	}
	return nil
}

// verifiedProblem returns why a non-draft Piece's `verified` date fails the
// gate, or "" when it passes or the Piece predates the rule. Days, the
// build clock's (now) included, are read in published's own offset.
func verifiedProblem(published, verified, now time.Time) string {
	pubDay := day(published)
	if !pubDay.After(verifiedCutoff) {
		return ""
	}
	switch {
	case verified.IsZero():
		return "no `verified` date: only the Author sets it, after checking every claim against its Source"
	case verified.Before(pubDay):
		return "verified " + verified.Format(time.DateOnly) + " is before published " + pubDay.Format(time.DateOnly)
	case verified.After(day(now.In(published.Location()))):
		return "verified " + verified.Format(time.DateOnly) + " is in the future (build clock " + now.In(published.Location()).Format(time.DateOnly) + ")"
	}
	return ""
}

// day is t's calendar day in its own offset, as midnight UTC (the form
// content.Explainer.Verified takes).
func day(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

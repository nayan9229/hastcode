package builder_test

import (
	"strings"
	"testing"
	"time"

	"hastcode.com/site/internal/buildtest"
)

// withSection returns ValidExplainer with an extra H2 section of exactly
// words words of prose (split over two paragraphs and a list) inserted
// before Sources.
func withSection(words int) buildtest.Explainer {
	half := words / 2
	prose := strings.TrimSpace(strings.Repeat("leader ", half)) + "\n\n- " + strings.TrimSpace(strings.Repeat("term ", words-half))
	e := buildtest.ValidExplainer()
	e.Body = strings.Replace(buildtest.ValidExplainerBody, "## Sources", "## Terms advance monotonically\n\n"+prose+"\n\n```go\nvar term = 1 // code is not counted\n```\n\n## Sources", 1)
	return e
}

func TestEditorialWarnings_SectionOver250WordsWarnsButStillPublishes(t *testing.T) {
	out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(withSection(251)))

	if res.Failed() {
		t.Fatalf("a long section failed the build; it must only warn:\n%s", res)
	}
	buildtest.RequireWarningRule(t, res, raftFile, "section-length", "## Terms advance monotonically")
	buildtest.RequireWarningRule(t, res, raftFile, "section-length", "251 words")
	if !buildtest.Exists(out, raftURL) {
		t.Error("page was not written for an Explainer that only has warnings")
	}
}

func TestEditorialWarnings_ReviewOverdueAfter90Days(t *testing.T) {
	// ValidExplainer was last Reviewed 2026-10-15T18:30:00Z.
	reviewed := time.Date(2026, 10, 15, 18, 30, 0, 0, time.UTC)
	for _, tc := range []struct {
		name string
		now  time.Time
		warn bool
	}{
		{"89 days", reviewed.AddDate(0, 0, 89), false},
		{"90 days", reviewed.AddDate(0, 0, 90), false},
		{"91 days", reviewed.AddDate(0, 0, 91), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := buildtest.Config()
			cfg.Now = tc.now

			out, res := buildtest.Build(t, cfg, buildtest.ValidTree())

			if res.Failed() {
				t.Fatalf("an overdue Review failed the build; it must only warn:\n%s", res)
			}
			if tc.warn {
				buildtest.RequireWarningRule(t, res, raftFile, "review-overdue", "91 days")
			} else {
				buildtest.RequireNoProblemRule(t, res, raftFile, "review-overdue")
			}
			if !buildtest.Exists(out, raftURL) {
				t.Error("page was not written for an Explainer with an overdue Review")
			}
		})
	}
}

func TestEditorialWarnings_SectionOf250WordsDoesNotWarn(t *testing.T) {
	_, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(withSection(250)))

	buildtest.RequireNoProblemRule(t, res, raftFile, "section-length")
}

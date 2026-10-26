package builder_test

import (
	"strings"
	"testing"
	"time"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildtest"
)

// verifiedClock is the build clock of the verified-gate tests: a week
// after the rule's cutoff (2026-10-26).
var verifiedClock = time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)

// newExplainer is a non-draft Explainer published after the cutoff, with
// no `verified` date.
func newExplainer() buildtest.Explainer {
	e := buildtest.ValidExplainer()
	e.Slug = "raft-log-compaction"
	e.Title = "How Raft compacts its log"
	e.Published = "2026-10-28T09:00:00Z"
	e.Reviewed = "2026-10-28T09:00:00Z"
	return e
}

const newExplainerFile = "explainers/distributed-systems/raft-log-compaction.md"

func verifiedConfig() builder.Config {
	cfg := buildtest.Config()
	cfg.Now = verifiedClock
	return cfg
}

func TestVerified_NewPieceWithoutItFailsProduction(t *testing.T) {
	tree := buildtest.ValidTree().WithExplainer(newExplainer())

	_, res := buildtest.Build(t, verifiedConfig(), tree)

	buildtest.RequireErrorRule(t, res, newExplainerFile, "unverified", "verified")
}

func TestVerified_ValidDatePasses(t *testing.T) {
	for _, verified := range []string{"2026-10-28", "2026-10-30", "2026-11-02"} {
		t.Run(verified, func(t *testing.T) {
			e := newExplainer()
			e.Verified = verified
			tree := buildtest.ValidTree().WithExplainer(e)

			out := buildtest.MustBuild(t, verifiedConfig(), tree)

			if !buildtest.Exists(out, "/distributed-systems/raft-log-compaction/") {
				t.Error("the verified Explainer was not written")
			}
		})
	}
}

func TestVerified_UnquotedYAMLDatePasses(t *testing.T) {
	e := newExplainer()
	e.Verified = "2026-10-29"
	md := strings.Replace(e.Markdown(), `verified: "2026-10-29"`, "verified: 2026-10-29", 1)
	tree := buildtest.ValidTree().With(newExplainerFile, []byte(md))

	out := buildtest.MustBuild(t, verifiedConfig(), tree)

	if !buildtest.Exists(out, "/distributed-systems/raft-log-compaction/") {
		t.Error("the verified Explainer was not written")
	}
}

func TestVerified_NotADateFails(t *testing.T) {
	e := newExplainer()
	e.Verified = "yesterday"
	tree := buildtest.ValidTree().WithExplainer(e)

	_, res := buildtest.Build(t, verifiedConfig(), tree)

	buildtest.RequireError(t, res, newExplainerFile, `verified "yesterday" is not a date (YYYY-MM-DD)`)
}

func TestVerified_DateBeforePublishedFails(t *testing.T) {
	e := newExplainer()
	e.Verified = "2026-10-27"
	tree := buildtest.ValidTree().WithExplainer(e)

	_, res := buildtest.Build(t, verifiedConfig(), tree)

	buildtest.RequireErrorRule(t, res, newExplainerFile, "unverified", "before published")
}

func TestVerified_FutureDateFails(t *testing.T) {
	e := newExplainer()
	e.Verified = "2026-11-03" // the day after the build clock
	tree := buildtest.ValidTree().WithExplainer(e)

	_, res := buildtest.Build(t, verifiedConfig(), tree)

	buildtest.RequireErrorRule(t, res, newExplainerFile, "unverified", "in the future")
}

func TestVerified_PreviewWarnsInsteadOfFailing(t *testing.T) {
	tree := buildtest.ValidTree().WithExplainer(newExplainer())
	cfg := verifiedConfig()
	cfg.Mode = builder.Preview

	out, res := buildtest.Build(t, cfg, tree)

	buildtest.RequireWarningRule(t, res, newExplainerFile, "unverified", "verified")
	if res.Failed() {
		t.Fatalf("preview build failed:\n%s", res)
	}
	if !buildtest.Exists(out, "/distributed-systems/raft-log-compaction/") {
		t.Error("the unverified Explainer is missing from the preview")
	}
}

func TestVerified_DraftIsOutsideTheRule(t *testing.T) {
	e := newExplainer()
	e.Draft = true
	tree := buildtest.ValidTree().WithExplainer(e)
	for _, mode := range []builder.Mode{builder.Production, builder.Preview} {
		cfg := verifiedConfig()
		cfg.Mode = mode

		_, res := buildtest.Build(t, cfg, tree)

		buildtest.RequireNoProblemRule(t, res, newExplainerFile, "unverified")
	}
}

func TestVerified_AppliesToEveryKindAndToBriefs(t *testing.T) {
	caseStudy := newExplainer()
	caseStudy.Kind = "case-study"
	comparison := buildtest.ValidComparison()
	profile := buildtest.ValidRepoProfile()
	brief := buildtest.ValidBrief()
	pieces := []struct {
		name, file string
		tree       func(verified string) buildtest.Tree
	}{
		{"concept", newExplainerFile, func(v string) buildtest.Tree {
			e := newExplainer()
			e.Verified = v
			return buildtest.ValidTree().WithExplainer(e)
		}},
		{"case-study", newExplainerFile, func(v string) buildtest.Tree {
			e := caseStudy
			e.Verified = v
			return buildtest.ValidTree().WithExplainer(e)
		}},
		{"comparison", "explainers/databases-storage/btree-vs-lsm-tree.md", func(v string) buildtest.Tree {
			e := comparison
			e.Published, e.Reviewed, e.Verified = "2026-10-28T09:00:00Z", "2026-10-28T09:00:00Z", v
			return buildtest.ValidTree().WithExplainer(e)
		}},
		{"repo-profile", "explainers/distributed-systems/apache-iggy.md", func(v string) buildtest.Tree {
			e := profile
			e.Published, e.Reviewed, e.Verified = "2026-10-28T09:00:00Z", "2026-10-28T09:00:00Z", v
			return buildtest.ValidTree().WithRepoProfile(e)
		}},
		{"brief", "briefs/etcd-raft-rewrite.md", func(v string) buildtest.Tree {
			b := brief
			b.Published, b.Verified = "2026-10-28T08:15:00-04:00", v
			return buildtest.ValidTree().WithBrief(b)
		}},
	}
	for _, p := range pieces {
		t.Run(p.name, func(t *testing.T) {
			_, res := buildtest.Build(t, verifiedConfig(), p.tree(""))
			buildtest.RequireErrorRule(t, res, p.file, "unverified", "verified")

			_, res = buildtest.Build(t, verifiedConfig(), p.tree("2026-10-29"))
			buildtest.RequireNoProblemRule(t, res, p.file, "unverified")
			if res.Failed() {
				t.Errorf("verified %s failed the build:\n%s", p.name, res)
			}
		})
	}
}

func TestVerified_PiecePublishedOnOrBeforeCutoffIsOutsideTheRule(t *testing.T) {
	// The last is 2026-10-27 in UTC but 2026-10-26 in its own offset.
	for _, published := range []string{"2026-10-22T09:00:00Z", "2026-10-26T23:59:00Z", "2026-10-26T23:30:00-05:00"} {
		t.Run(published, func(t *testing.T) {
			e := newExplainer()
			e.Published = published
			e.Reviewed = "2026-10-28T09:00:00Z"
			tree := buildtest.ValidTree().WithExplainer(e)

			out, res := buildtest.Build(t, verifiedConfig(), tree)

			buildtest.RequireNoProblemRule(t, res, newExplainerFile, "unverified")
			if !buildtest.Exists(out, "/distributed-systems/raft-log-compaction/") {
				t.Errorf("the Explainer was not written:\n%s", res)
			}
		})
	}
}

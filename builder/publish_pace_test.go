package builder_test

import (
	"fmt"
	"testing"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildtest"
)

// pacedExplainer is a verified, non-draft Explainer published at the given
// time, with a slug of its own.
func pacedExplainer(n int, published string) buildtest.Explainer {
	e := newExplainer()
	e.Slug = fmt.Sprintf("raft-log-compaction-%d", n)
	e.Title = fmt.Sprintf("How Raft compacts its log, part %d", n)
	e.Published, e.Reviewed = published, published
	e.Verified = published[:10]
	return e
}

func pacedTree(published ...string) buildtest.Tree {
	tree := buildtest.ValidTree()
	for i, p := range published {
		tree = tree.WithExplainer(pacedExplainer(i, p))
	}
	return tree
}

func TestPublishPace_ThreeOnOneDayFailsProductionNamingTheDay(t *testing.T) {
	tree := pacedTree("2026-10-31T08:00:00Z", "2026-10-31T12:00:00Z", "2026-10-31T20:00:00Z")

	_, res := buildtest.Build(t, verifiedConfig(), tree)

	buildtest.RequireErrorRule(t, res, "", "publish-pace", "2026-10-31")
	buildtest.RequireErrorRule(t, res, "", "publish-pace", "raft-log-compaction-2.md")
}

func TestPublishPace_TwoADayPasses(t *testing.T) {
	tree := pacedTree("2026-10-31T08:00:00Z", "2026-10-31T20:00:00Z", "2026-11-01T08:00:00Z", "2026-11-01T09:00:00Z")

	buildtest.MustBuild(t, verifiedConfig(), tree)
}

// The cap starts on 2026-10-30: Pieces published before it stay.
func TestPublishPace_DaysBeforeTheCapAreExempt(t *testing.T) {
	tree := pacedTree("2026-10-28T08:00:00Z", "2026-10-28T09:00:00Z", "2026-10-28T10:00:00Z")

	buildtest.MustBuild(t, verifiedConfig(), tree)
}

func TestPublishPace_DraftsDoNotCount(t *testing.T) {
	tree := buildtest.ValidTree()
	for i := 0; i < 4; i++ {
		e := pacedExplainer(i, "2026-10-31T08:00:00Z")
		e.Draft = i >= 2
		tree = tree.WithExplainer(e)
	}

	buildtest.MustBuild(t, verifiedConfig(), tree)
}

// Briefs count toward the same cap as Explainers.
func TestPublishPace_BriefsCountToo(t *testing.T) {
	b := buildtest.ValidBrief()
	b.Published = "2026-10-31T10:00:00Z"
	b.Verified = "2026-10-31"
	tree := pacedTree("2026-10-31T08:00:00Z", "2026-10-31T20:00:00Z").WithBrief(b)

	_, res := buildtest.Build(t, verifiedConfig(), tree)

	buildtest.RequireErrorRule(t, res, "", "publish-pace", "3 Pieces are published on 2026-10-31")
}

func TestPublishPace_PreviewOnlyWarns(t *testing.T) {
	cfg := verifiedConfig()
	cfg.Mode = builder.Preview
	tree := pacedTree("2026-10-31T08:00:00Z", "2026-10-31T12:00:00Z", "2026-10-31T20:00:00Z")

	_, res := buildtest.Build(t, cfg, tree)

	if len(res.Errors) != 0 {
		t.Fatalf("preview build errors: %v", res.Errors)
	}
	buildtest.RequireWarningRule(t, res, "", "publish-pace", "2026-10-31")
}

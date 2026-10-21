package builder_test

import (
	"testing"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildtest"
)

func draftTree() buildtest.Tree {
	d := buildtest.ValidExplainer()
	d.Slug = "raft-membership-changes"
	d.Title = "How Raft changes cluster membership"
	d.Draft = true
	return buildtest.ValidTree().WithExplainer(d)
}

const draftURL = "/distributed-systems/raft-membership-changes/"

func TestDraft_AbsentFromProductionOutput(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), draftTree())

	if buildtest.Exists(out, draftURL) {
		t.Errorf("draft page %s was written in a production build", draftURL)
	}
	for _, needle := range []string{"raft-membership-changes", "How Raft changes cluster membership"} {
		if hits := buildtest.FilesContaining(t, out, needle); len(hits) > 0 {
			t.Errorf("draft %q appears in production output: %v", needle, hits)
		}
	}
	if !buildtest.Exists(out, raftURL) {
		t.Error("the published Explainer beside the draft was not written")
	}
}

func TestDraft_PresentInPreviewOutput(t *testing.T) {
	cfg := buildtest.Config()
	cfg.Mode = builder.Preview

	out := buildtest.MustBuild(t, cfg, draftTree())

	doc := buildtest.Page(t, out, draftURL)
	if got := doc.Find("h1.p-name").Text(); got != "How Raft changes cluster membership" {
		t.Errorf("draft h1 = %q", got)
	}
}

func TestDraft_UnfinishedDraftDoesNotFailProduction(t *testing.T) {
	d := buildtest.ValidExplainer()
	d.Slug = "raft-membership-changes"
	d.Draft = true
	d.Body = "Notes so far: joint consensus.\n"
	tree := buildtest.ValidTree().WithExplainer(d)

	buildtest.MustBuild(t, buildtest.Config(), tree)

	cfg := buildtest.Config()
	cfg.Mode = builder.Preview
	_, res := buildtest.Build(t, cfg, tree)
	buildtest.RequireErrorRule(t, res, "explainers/distributed-systems/raft-membership-changes.md", "tldr", "")
}

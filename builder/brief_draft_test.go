package builder_test

import (
	"testing"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildtest"
)

func draftBriefTree() buildtest.Tree {
	d := buildtest.ValidBrief()
	d.Title = "A Brief still being drafted"
	d.Draft = true
	return buildtest.ValidTree().WithBrief(d)
}

func TestBriefDraft_AbsentFromProductionOutput(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), draftBriefTree())

	if buildtest.Exists(out, briefURL) {
		t.Errorf("draft Brief %s was written in a production build", briefURL)
	}
	for _, needle := range []string{"etcd-raft-rewrite", "A Brief still being drafted"} {
		if hits := buildtest.FilesContaining(t, out, needle); len(hits) > 0 {
			t.Errorf("draft Brief %q appears in production output: %v", needle, hits)
		}
	}
}

func TestBriefDraft_PresentInPreviewOutput(t *testing.T) {
	cfg := buildtest.Config()
	cfg.Mode = builder.Preview

	out := buildtest.MustBuild(t, cfg, draftBriefTree())

	if got := buildtest.Page(t, out, briefURL).Find("h1.p-name").Text(); got != "A Brief still being drafted" {
		t.Errorf("draft Brief h1 = %q", got)
	}
}

func TestBriefDraft_UnfinishedDraftDoesNotFailProduction(t *testing.T) {
	d := buildtest.ValidBrief()
	d.Draft = true
	d.Body = "Notes so far.\n"
	// A draft may even reuse a published slug while it is being written.
	d.Slug = buildtest.ValidExplainer().Slug

	buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree().WithBrief(d))
}

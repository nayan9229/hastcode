package builder_test

import (
	"strings"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

func TestFrontmatter_UnknownVerticalIsABuildError(t *testing.T) {
	e := buildtest.ValidExplainer()
	e.Vertical = "security"

	out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

	buildtest.RequireError(t, res, "explainers/security/raft-leader-election.md", `unknown Vertical "security"`)
	if buildtest.Exists(out, "/security/raft-leader-election/") {
		t.Error("page was written for an Explainer with an unknown Vertical")
	}
}

func TestFrontmatter_UnknownAuthorIsABuildError(t *testing.T) {
	e := buildtest.ValidExplainer()
	e.Author = "nobody"

	out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

	buildtest.RequireError(t, res, "explainers/distributed-systems/raft-leader-election.md", `unknown Author "nobody"`)
	if buildtest.Exists(out, "/distributed-systems/raft-leader-election/") {
		t.Error("page was written for an Explainer with an unknown Author")
	}
}

func TestFrontmatter_PublishedWithoutTimezoneIsABuildError(t *testing.T) {
	e := buildtest.ValidExplainer()
	e.Published = "2026-10-01T09:00:00"

	out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

	buildtest.RequireError(t, res, "explainers/distributed-systems/raft-leader-election.md", "published")
	if buildtest.Exists(out, "/distributed-systems/raft-leader-election/") {
		t.Error("page was written for an Explainer with an invalid published date")
	}
}

func TestFrontmatter_KindOutsideTheFourKindsIsABuildErrorNamingThem(t *testing.T) {
	e := buildtest.ValidExplainer()
	e.Kind = "tutorial"

	_, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

	buildtest.RequireError(t, res, "explainers/distributed-systems/raft-leader-election.md", `kind "tutorial" must be one of case-study, concept, comparison, repo-profile`)
}

func TestFrontmatter_UnknownKeyIsABuildError(t *testing.T) {
	file := "explainers/distributed-systems/raft-leader-election.md"
	tree := buildtest.NewTree().WithExplainer(buildtest.ValidExplainer())
	tree.With(file, []byte("---\ntitel: typo\n"+strings.TrimPrefix(string(tree[file]), "---\n")))

	_, res := buildtest.Build(t, buildtest.Config(), tree)

	buildtest.RequireError(t, res, file, "titel")
}

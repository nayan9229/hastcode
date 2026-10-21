package builder_test

import (
	"strings"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

const raftFile = "explainers/distributed-systems/raft-leader-election.md"
const raftURL = "/distributed-systems/raft-leader-election/"

// buildWithBody builds a tree holding ValidExplainer with its body replaced.
func buildWithBody(t *testing.T, body string) (string, func(rule, want string)) {
	t.Helper()
	e := buildtest.ValidExplainer()
	e.Body = body
	out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))
	return out, func(rule, want string) {
		t.Helper()
		buildtest.RequireErrorRule(t, res, raftFile, rule, want)
	}
}

func requireNoPage(t *testing.T, out, url string) {
	t.Helper()
	if buildtest.Exists(out, url) {
		t.Errorf("page %s was written for an Explainer that failed a gate", url)
	}
}

func TestEditorialGates_ExplainerWithoutTLDRFails(t *testing.T) {
	body := strings.Replace(buildtest.ValidExplainerBody, "## TL;DR\n\n- Each term has at most one leader.\n- Followers become candidates when their election timeout fires.\n- A candidate needs votes from a majority of the cluster.\n\n", "", 1)

	out, requireError := buildWithBody(t, body)

	requireError("tldr", "## TL;DR")
	requireNoPage(t, out, raftURL)
}

const validTLDR = "## TL;DR\n\n- Each term has at most one leader.\n- Followers become candidates when their election timeout fires.\n- A candidate needs votes from a majority of the cluster.\n"

func TestEditorialGates_TLDRMustHaveThreeToFiveBullets(t *testing.T) {
	for _, tc := range []struct {
		name string
		tldr string
	}{
		{"two bullets", "## TL;DR\n\n- One.\n- Two.\n"},
		{"six bullets", "## TL;DR\n\n- One.\n- Two.\n- Three.\n- Four.\n- Five.\n- Six.\n"},
		{"prose instead of bullets", "## TL;DR\n\nRaft elects one leader per term.\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, requireError := buildWithBody(t, strings.Replace(buildtest.ValidExplainerBody, validTLDR, tc.tldr, 1))

			requireError("tldr", "3–5 bullets")
			requireNoPage(t, out, raftURL)
		})
	}
}

func TestEditorialGates_FiveBulletTLDRPasses(t *testing.T) {
	e := buildtest.ValidExplainer()
	e.Body = strings.Replace(buildtest.ValidExplainerBody, "- A candidate needs", "- Four.\n- Five.\n- A candidate needs", 1)

	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

	if got := buildtest.Page(t, out, raftURL).Find(".e-content h2").First().Next().Find("li").Length(); got != 5 {
		t.Errorf("TL;DR bullets = %d, want 5", got)
	}
}

const validSources = "## Sources\n\n- [In Search of an Understandable Consensus Algorithm](https://raft.github.io/raft.pdf)\n"

func TestEditorialGates_ExplainerWithoutSourcesFails(t *testing.T) {
	out, requireError := buildWithBody(t, strings.Replace(buildtest.ValidExplainerBody, validSources, "", 1))

	requireError("sources", "missing a ## Sources section")
	requireNoPage(t, out, raftURL)
}

func TestEditorialGates_SourcesWithoutALinkFails(t *testing.T) {
	body := strings.Replace(buildtest.ValidExplainerBody, validSources, "## Sources\n\n- The Raft paper, by Ongaro and Ousterhout.\n", 1)

	out, requireError := buildWithBody(t, body)

	requireError("sources", "at least one link")
	requireNoPage(t, out, raftURL)
}

func TestEditorialGates_SourcesNotAtTheEndFails(t *testing.T) {
	body := buildtest.ValidExplainerBody + "\n## Further reading\n\nSee the Raft website.\n"

	out, requireError := buildWithBody(t, body)

	requireError("sources", "at the end")
	requireNoPage(t, out, raftURL)
}

const validMermaid = "```mermaid\n%% title: A Raft election\nflowchart LR\n  Follower --> Candidate --> Leader\n```"

func TestEditorialGates_ExplainerWithoutMermaidDiagramFails(t *testing.T) {
	for _, tc := range []struct {
		name    string
		diagram string
	}{
		{"no fence", ""},
		{"non-mermaid fence", "```text\nFollower -> Candidate -> Leader\n```"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, requireError := buildWithBody(t, strings.Replace(buildtest.ValidExplainerBody, validMermaid, tc.diagram, 1))

			requireError("diagram", "Mermaid")
			requireNoPage(t, out, raftURL)
		})
	}
}

func TestEditorialGates_ValidComparisonPasses(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithExplainer(buildtest.ValidComparison()))

	if got := buildtest.Page(t, out, "/databases-storage/btree-vs-lsm-tree/").Find(".e-content table").Length(); got != 1 {
		t.Errorf("tables = %d, want 1", got)
	}
}

func TestEditorialGates_ComparisonWithoutTradeOffsTableFails(t *testing.T) {
	table := "| | B-tree | LSM-tree |\n|---|---|---|\n| Writes | Random, in place | Sequential, batched |\n| Reads | One tree walk | Several levels plus Bloom filters |\n"
	for _, tc := range []struct {
		name      string
		tradeOffs string
	}{
		{"no trade-offs section", ""},
		{"trade-offs section without a table", "## Trade-offs\n\nB-trees favour reads; LSM-trees favour writes.\n"},
		{"table outside the trade-offs section", "## Write costs\n\n" + table + "\n## Trade-offs\n\nB-trees favour reads.\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e := buildtest.ValidComparison()
			e.Body = strings.Replace(buildtest.ValidComparisonBody, buildtest.ValidComparisonTradeOffs, tc.tradeOffs, 1)

			out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

			buildtest.RequireErrorRule(t, res, "explainers/databases-storage/btree-vs-lsm-tree.md", "trade-offs-table", "## Trade-offs")
			requireNoPage(t, out, "/databases-storage/btree-vs-lsm-tree/")
		})
	}
}

func TestEditorialGates_ConceptNeedsNoTradeOffsTable(t *testing.T) {
	_, res := buildtest.Build(t, buildtest.Config(), buildtest.ValidTree())

	buildtest.RequireNoProblemRule(t, res, raftFile, "trade-offs-table")
}

func TestEditorialGates_AllFailuresAreReportedTogether(t *testing.T) {
	bare := buildtest.ValidExplainer()
	bare.Body = "Raft is a consensus algorithm.\n\n## Elections\n\nFollowers time out and stand for election.\n"
	noTable := buildtest.ValidComparison()
	noTable.Body = strings.Replace(buildtest.ValidComparisonBody, buildtest.ValidComparisonTradeOffs, "", 1)
	valid := buildtest.ValidExplainer()
	valid.Slug = "raft-log-replication"

	out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(bare).WithExplainer(noTable).WithExplainer(valid))

	buildtest.RequireErrorRule(t, res, raftFile, "tldr", "")
	buildtest.RequireErrorRule(t, res, raftFile, "sources", "")
	buildtest.RequireErrorRule(t, res, raftFile, "diagram", "")
	buildtest.RequireErrorRule(t, res, "explainers/databases-storage/btree-vs-lsm-tree.md", "trade-offs-table", "")
	if got := len(res.Errors); got != 4 {
		t.Errorf("errors = %d, want 4:\n%s", got, res)
	}
	requireNoPage(t, out, raftURL)
	requireNoPage(t, out, "/databases-storage/btree-vs-lsm-tree/")
	if !buildtest.Exists(out, "/distributed-systems/raft-log-replication/") {
		t.Error("the valid Explainer's page was not written alongside the failed ones")
	}
}

func TestEditorialGates_TLDRAfterAnotherSectionFails(t *testing.T) {
	body := strings.Replace(buildtest.ValidExplainerBody, "## TL;DR", "## Background\n\nRaft is a consensus algorithm.\n\n## TL;DR", 1)

	out, requireError := buildWithBody(t, body)

	requireError("tldr", "immediately after the intro")
	requireNoPage(t, out, raftURL)
}

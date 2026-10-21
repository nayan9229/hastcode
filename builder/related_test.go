package builder_test

import (
	"strings"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

// relatedTree is listingTree with raft-leader-election (tags consensus,
// raft) naming related slugs in frontmatter.
func relatedTree(related ...string) buildtest.Tree {
	raft := buildtest.ValidExplainer()
	raft.Related = related
	return listingTree().WithExplainer(raft)
}

func TestExplainer_MoreInVerticalListsTheThreeLatestOthers(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), listingTree())
	section := buildtest.Page(t, out, raftURL).Find("section.more-in-vertical")

	if got := strings.TrimSpace(section.Find("h2").Text()); got != "More in Distributed Systems" {
		t.Errorf("heading = %q", got)
	}
	got := listedURLs(section)
	want := []string{
		"/distributed-systems/spanner-truetime/",
		"/distributed-systems/gossip/",
		"/distributed-systems/paxos/",
	}
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("More in Vertical = %v, want %v", got, want)
	}

	// paxos is the oldest: the three latest others are raft, spanner, gossip.
	got = listedURLs(buildtest.Page(t, out, "/distributed-systems/paxos/").Find("section.more-in-vertical"))
	want = []string{
		"/distributed-systems/raft-leader-election/",
		"/distributed-systems/spanner-truetime/",
		"/distributed-systems/gossip/",
	}
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("paxos More in Vertical = %v, want %v", got, want)
	}
}

func TestExplainer_RelatedRanksBySharedTagsThenNewest(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), listingTree())

	// raft shares "consensus" with dynamo (newer) and paxos.
	got := listedURLs(buildtest.Page(t, out, raftURL).Find("section.related"))
	if want := "/databases-storage/dynamo/ /distributed-systems/paxos/"; strings.Join(got, " ") != want {
		t.Errorf("Related = %v, want %s", got, want)
	}

	// dynamo shares one tag each with raft, paxos and btree-vs-lsm; newest first.
	got = listedURLs(buildtest.Page(t, out, "/databases-storage/dynamo/").Find("section.related"))
	want := "/distributed-systems/raft-leader-election/ /databases-storage/btree-vs-lsm/ /distributed-systems/paxos/"
	if strings.Join(got, " ") != want {
		t.Errorf("dynamo Related = %v, want %s", got, want)
	}
}

func TestExplainer_RelatedPutsFrontmatterSlugsFirst(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), relatedTree("gossip", "btree-vs-lsm"))

	got := listedURLs(buildtest.Page(t, out, raftURL).Find("section.related"))
	want := "/distributed-systems/gossip/ /databases-storage/btree-vs-lsm/ /databases-storage/dynamo/"
	if strings.Join(got, " ") != want {
		t.Errorf("Related = %v, want %s", got, want)
	}
}

func TestExplainer_NeverRelatesToItself(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), relatedTree("raft-leader-election"))
	doc := buildtest.Page(t, out, raftURL)

	for _, section := range []string{"section.related", "section.more-in-vertical"} {
		for _, url := range listedURLs(doc.Find(section)) {
			if url == raftURL {
				t.Errorf("%s lists the Explainer itself", section)
			}
		}
	}
	if got := listedURLs(doc.Find("section.related")); len(got) != 2 {
		t.Errorf("Related = %v, want the two tag matches only", got)
	}
}

func TestExplainer_NoRelatedSectionsWhenNothingQualifies(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	doc := buildtest.Page(t, out, raftURL)

	if n := doc.Find("section.more-in-vertical, section.related").Length(); n != 0 {
		t.Errorf("a lone Explainer has %d reading-list sections, want none", n)
	}
}

func TestExplainer_UnknownRelatedSlugIsABrokenLink(t *testing.T) {
	_, res := buildtest.Build(t, buildtest.Config(), relatedTree("raft-log-compaction"))

	buildtest.RequireErrorRule(t, res, "explainers/distributed-systems/raft-leader-election.md", "broken-link", `"raft-log-compaction"`)
}

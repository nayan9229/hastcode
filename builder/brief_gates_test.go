package builder_test

import (
	"strings"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

const briefFile = "briefs/etcd-raft-rewrite.md"

// requireBriefFails builds a tree holding ValidTree plus b, and requires a
// build error on the Brief with the given rule and message fragment, and
// no Brief page.
func requireBriefFails(t *testing.T, b buildtest.Brief, rule, substring string) {
	t.Helper()
	out, res := buildtest.Build(t, buildtest.Config(), buildtest.ValidTree().WithBrief(b))
	buildtest.RequireErrorRule(t, res, briefFile, rule, substring)
	if buildtest.Exists(out, briefURL) {
		t.Errorf("page %s was written for a Brief that failed %s", briefURL, rule)
	}
}

// validBriefWords is the word count of ValidBriefBody, headings included:
// 2 + 12 + 3 + 12 + 9.
const validBriefWords = 38

// paddedBrief is ValidBrief with its body padded to n words.
func paddedBrief(n int) buildtest.Brief {
	b := buildtest.ValidBrief()
	b.Body = strings.Replace(buildtest.ValidBriefBody, "## Why it matters",
		strings.Repeat("word ", n-validBriefWords)+"\n\n## Why it matters", 1)
	return b
}

func TestBriefGate_BodyOver250WordsFails(t *testing.T) {
	requireBriefFails(t, paddedBrief(251), "brief-length", "251 words")
}

func TestBriefGate_MissingWhyItMattersFails(t *testing.T) {
	b := buildtest.ValidBrief()
	b.Body = strings.Replace(buildtest.ValidBriefBody, "## Why it matters\n", "", 1)

	requireBriefFails(t, b, "brief-sections", `missing ## Why it matters`)
}

func TestBriefGate_ExtraOrReorderedSectionsFail(t *testing.T) {
	cases := map[string]struct{ body, message string }{
		"extra section": {
			body:    buildtest.ValidBriefBody + "\n## What next\n\nWatch the next release.\n",
			message: `## What next`,
		},
		"reordered": {
			body:    "## Why it matters\n\nSee [Raft](raft-leader-election).\n\n## What happened\n\netcd shipped.\n",
			message: `exactly ## What happened then ## Why it matters`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			b := buildtest.ValidBrief()
			b.Body = c.body
			requireBriefFails(t, b, "brief-sections", c.message)
		})
	}
}

func TestBriefGate_MissingOrRelativeSourceURLFails(t *testing.T) {
	for name, source := range map[string]string{"missing": "", "relative": "/briefs/elsewhere/"} {
		t.Run(name, func(t *testing.T) {
			b := buildtest.ValidBrief()
			b.SourceURL = source
			requireBriefFails(t, b, "brief-source", "source_url")
		})
	}
}

func TestBriefGate_NoExplainerLinkFails(t *testing.T) {
	b := buildtest.ValidBrief()
	b.Body = strings.Replace(buildtest.ValidBriefBody, "[how Raft elects a leader](raft-leader-election)", "how Raft elects a leader", 1)
	b.Explainers = nil

	requireBriefFails(t, b, "brief-explainer-link", "at least one Explainer")
}

func TestBriefGate_LinkToUnknownExplainerIsABrokenLink(t *testing.T) {
	cases := map[string]string{
		"bare slug": "raft-log-compaction",
		"site path": "/distributed-systems/raft-log-compaction/",
		// The slug exists, but not in that Vertical: the URL would 404.
		"wrong vertical": "/databases-storage/raft-leader-election/",
	}
	for name, dest := range cases {
		t.Run(name, func(t *testing.T) {
			b := buildtest.ValidBrief()
			b.Body = buildtest.ValidBriefBody + "\nAlso read [compaction](" + dest + ").\n"
			requireBriefFails(t, b, "broken-link", dest)
		})
	}
}

func TestBriefGate_ExplainersFrontmatterMustMatchTheBodyLinks(t *testing.T) {
	t.Run("listed but unknown", func(t *testing.T) {
		b := buildtest.ValidBrief()
		b.Explainers = []string{"raft-leader-election", "raft-log-compaction"}
		requireBriefFails(t, b, "broken-link", `"raft-log-compaction"`)
	})
	t.Run("listed but not linked", func(t *testing.T) {
		b := buildtest.ValidBrief()
		b.Explainers = []string{"raft-leader-election", "btree-vs-lsm-tree"}
		out, res := buildtest.Build(t, buildtest.Config(), buildtest.ValidTree().WithExplainer(buildtest.ValidComparison()).WithBrief(b))
		buildtest.RequireErrorRule(t, res, briefFile, "brief-explainer-link", `"btree-vs-lsm-tree"`)
		if buildtest.Exists(out, briefURL) {
			t.Errorf("page %s was written", briefURL)
		}
	})
	t.Run("linked but not listed", func(t *testing.T) {
		b := buildtest.ValidBrief()
		b.Explainers = nil
		requireBriefFails(t, b, "brief-explainer-link", `"raft-leader-election"`)
	})
}

func TestBriefGate_UnknownAuthorFails(t *testing.T) {
	b := buildtest.ValidBrief()
	b.Author = "grace"
	requireBriefFails(t, b, "unknown-author", `"grace"`)
}

func TestBriefGate_MissingHeroFails(t *testing.T) {
	b := buildtest.ValidBrief()
	b.Hero = ""
	requireBriefFails(t, b, "hero-missing", "hero")
}

func TestBriefFrontmatter_ExplainerOnlyKeysAreRejected(t *testing.T) {
	for _, key := range []string{`kind: "concept"`, `vertical: "distributed-systems"`, `reviewed: "2026-10-21T09:00:00Z"`} {
		t.Run(key, func(t *testing.T) {
			md := strings.Replace(buildtest.ValidBrief().Markdown(), "---\n\n", key+"\n---\n\n", 1)
			tree := buildtest.ValidTree().WithBrief(buildtest.ValidBrief()).With(briefFile, []byte(md))
			out, res := buildtest.Build(t, buildtest.Config(), tree)
			buildtest.RequireErrorRule(t, res, briefFile, "frontmatter", strings.SplitN(key, ":", 2)[0])
			if buildtest.Exists(out, briefURL) {
				t.Errorf("page %s was written", briefURL)
			}
		})
	}
}

func TestBriefGate_BodyOf250WordsPasses(t *testing.T) {
	b := paddedBrief(250)

	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree().WithBrief(b))
	if !buildtest.Exists(out, briefURL) {
		t.Errorf("a 250-word Brief was not written")
	}
}

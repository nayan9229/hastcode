package builder_test

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/internal/buildtest"
)

// The Repo Profile skeleton (spec "Page composition and gates"): the
// Markdown body carries exactly the nine opinion headings in a fixed order,
// never a hand-written Signals or Changelog (they render from the Fact
// Pack), and the FAQ holds five to eight questions.

// requireRepoProfileFails builds ValidRepoProfile with its body replaced
// and requires a build error on the Profile with the given rule and
// message fragment, and no Profile page.
func requireRepoProfileFails(t *testing.T, body, rule, substring string) {
	t.Helper()
	p := buildtest.ValidRepoProfile()
	p.Body = body
	out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithRepoProfile(p))
	buildtest.RequireErrorRule(t, res, repoProfileFile, rule, substring)
	if buildtest.Exists(out, repoProfilePath) {
		t.Errorf("page %s was written for a Repo Profile that failed %s", repoProfilePath, rule)
	}
}

const validVerdict = "## Verdict\n\nAdopt it for single-node streaming where a small binary matters; wait for clustering otherwise.\n\n"

func TestRepoProfileSkeleton_MissingHeadingFailsNamingTheHeadingAndItsPosition(t *testing.T) {
	body := strings.Replace(buildtest.ValidRepoProfileBody, validVerdict, "", 1)

	requireRepoProfileFails(t, body, "repo-profile-sections", "missing ## Verdict (section 8, after ## FAQ)")
}

func TestRepoProfileSkeleton_ExtraHeadingFailsIncludingAnH2UnderFAQ(t *testing.T) {
	cases := map[string]struct{ body, message string }{
		"extra section before FAQ": {
			body:    strings.Replace(buildtest.ValidRepoProfileBody, "## FAQ\n", "## Roadmap\n\nClustering is planned.\n\n## FAQ\n", 1),
			message: "unexpected ## Roadmap (section 7)",
		},
		"H2 under FAQ": {
			body:    strings.Replace(buildtest.ValidRepoProfileBody, "## Verdict\n", "## Pricing\n\nFree.\n\n## Verdict\n", 1),
			message: "unexpected ## Pricing (section 8)",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			requireRepoProfileFails(t, c.body, "repo-profile-sections", c.message)
		})
	}
}

const (
	validUseItIf     = "## Use it if\n\n- You want a single binary and a small footprint.\n\n"
	validDontUseItIf = "## Don't use it if\n\n- You need multi-node replication today.\n\n"
)

func TestRepoProfileSkeleton_ReorderedHeadingsFailNamingThePositionAndTheHeadingThatBelongsThere(t *testing.T) {
	body := strings.Replace(buildtest.ValidRepoProfileBody, validUseItIf+validDontUseItIf, validDontUseItIf+validUseItIf, 1)

	requireRepoProfileFails(t, body, "repo-profile-sections", "## Don't use it if is section 4, where ## Use it if belongs")
}

func TestRepoProfileSkeleton_HandWrittenSignalsOrChangelogFailsAsRenderedFromTheFactPack(t *testing.T) {
	cases := map[string]struct{ body, message string }{
		"Signals after Architecture": {
			body:    strings.Replace(buildtest.ValidRepoProfileBody, "## Use it if\n", "## Signals\n\nStars: 2,468. Last commit: 18 October 2026.\n\n## Use it if\n", 1),
			message: "## Signals is rendered from the Fact Pack",
		},
		"Changelog after Alternatives": {
			body:    strings.Replace(buildtest.ValidRepoProfileBody, "## FAQ\n", "## Changelog\n\n- server-0.5.0, 10 October 2026.\n\n## FAQ\n", 1),
			message: "## Changelog is rendered from the Fact Pack",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			p := buildtest.ValidRepoProfile()
			p.Body = c.body
			out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithRepoProfile(p))

			buildtest.RequireErrorRule(t, res, repoProfileFile, "repo-profile-fact-section", c.message)
			// A fact section is not an "extra" opinion section: one error, the right one.
			buildtest.RequireNoProblemRule(t, res, repoProfileFile, "repo-profile-sections")
			if buildtest.Exists(out, repoProfilePath) {
				t.Errorf("page %s was written with a hand-written fact section", repoProfilePath)
			}
		})
	}
}

const validLastFAQ = "### How is Iggy licensed?\n\nApache 2.0.\n\n"

// faqBody is ValidRepoProfileBody with its FAQ holding n questions
// (ValidRepoProfileBody has five).
func faqBody(n int) string {
	extra := ""
	for i := 6; i <= n; i++ {
		extra += "### Question " + strings.Repeat("I", i-5) + "?\n\nAn answer.\n\n"
	}
	if n < 5 {
		return strings.Replace(buildtest.ValidRepoProfileBody, validLastFAQ, "", 1)
	}
	return strings.Replace(buildtest.ValidRepoProfileBody, validLastFAQ, validLastFAQ+extra, 1)
}

func TestRepoProfileSkeleton_FAQHoldsFiveToEightQuestions(t *testing.T) {
	for n, message := range map[int]string{4: "## FAQ has 4 questions, want 5–8", 9: "## FAQ has 9 questions, want 5–8"} {
		t.Run(message, func(t *testing.T) {
			requireRepoProfileFails(t, faqBody(n), "repo-profile-faq", message)
		})
	}
	for _, n := range []int{5, 8} {
		t.Run(strings.Repeat("I", n)+" questions pass", func(t *testing.T) {
			p := buildtest.ValidRepoProfile()
			p.Body = faqBody(n)
			out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithRepoProfile(p))

			if got := buildtest.Page(t, out, repoProfilePath).Find(".e-content h3").Length(); got != n {
				t.Errorf("FAQ questions on the page = %d, want %d", got, n)
			}
		})
	}
}

var repoProfileRules = []string{"repo-profile-sections", "repo-profile-fact-section", "repo-profile-faq"}

func TestRepoProfileSkeleton_LeavesCaseStudyConceptAndComparisonBodiesAlone(t *testing.T) {
	caseStudy := buildtest.ValidExplainer()
	caseStudy.Kind = "case-study"
	caseStudy.Slug = "etcd-raft-in-production"
	tree := buildtest.ValidTree().WithExplainer(caseStudy).WithExplainer(buildtest.ValidComparison())

	_, res := buildtest.Build(t, buildtest.Config(), tree)

	for _, file := range []string{raftFile, "explainers/distributed-systems/etcd-raft-in-production.md", "explainers/databases-storage/btree-vs-lsm-tree.md"} {
		for _, rule := range repoProfileRules {
			buildtest.RequireNoProblemRule(t, res, file, rule)
		}
	}
	if len(res.Errors) != 0 {
		t.Errorf("the Repo Profile skeleton touched other kinds:\n%s", res)
	}
}

func TestRepoProfileSkeleton_ExplainerTLDRAndSourcesRulesStillApplyToAProfile(t *testing.T) {
	cases := map[string]struct{ body, rule, message string }{
		"two-bullet TL;DR": {
			body:    strings.Replace(buildtest.ValidRepoProfileBody, "- Clustering is on the roadmap, not in a release.\n", "", 1),
			rule:    "tldr",
			message: "3–5 bullets",
		},
		"Sources without a link": {
			body:    strings.Replace(buildtest.ValidRepoProfileBody, "- [apache/iggy on GitHub](https://github.com/apache/iggy)\n", "- The Iggy repository.\n", 1),
			rule:    "sources",
			message: "at least one link",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			requireRepoProfileFails(t, c.body, c.rule, c.message)
		})
	}
}

func TestRepoProfileSkeleton_ValidRepoProfileFixturePassesWithItsHeadingsInOrder(t *testing.T) {
	out, res := buildtest.Build(t, buildtest.Config(), repoProfileTree())

	for _, rule := range repoProfileRules {
		buildtest.RequireNoProblemRule(t, res, repoProfileFile, rule)
	}
	if len(res.Errors) != 0 {
		t.Fatalf("the Repo Profile fixture must build clean:\n%s", res)
	}
	// The nine written headings reach the page in the spec's order (the
	// build may add fact sections between them).
	want := []string{"TL;DR", "In 60 seconds", "Architecture", "Use it if", "Don't use it if", "Alternatives", "FAQ", "Verdict", "Sources"}
	var got []string
	buildtest.Page(t, out, repoProfilePath).Find(".e-content h2").Each(func(_ int, h *goquery.Selection) {
		got = append(got, strings.TrimSpace(h.Text()))
	})
	next := 0
	for _, h := range got {
		if next < len(want) && h == want[next] {
			next++
		}
	}
	if next != len(want) {
		t.Errorf("page h2s = %v, want %v in that order", got, want)
	}
}

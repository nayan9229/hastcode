package builder_test

import (
	"slices"
	"strings"
	"testing"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildtest"
)

// The Author decided (ADR-0014) that no Reader sees an AI-assistance note
// or AI-assistance wording anywhere on the site. The build enforces it
// with the ai-note error rule.

// oldAINote is the note every Piece carried until 2026-11-18.
const oldAINote = "> Drafted with AI assistance, researched, verified and edited by the HastCode editorial team. We publish only as many pieces as the Author can verify this way.\n"

// withOpening puts text right after the hook of ValidExplainerBody, where
// the note used to sit.
func withOpening(text string) string {
	hook, rest, _ := strings.Cut(buildtest.ValidExplainerBody, "\n\n")
	return hook + "\n\n" + text + "\n" + rest
}

func TestAINote_ExplainerCarryingTheOldNoteFails(t *testing.T) {
	out, requireError := buildWithBody(t, withOpening(oldAINote))

	requireError("ai-note", "AI-assistance")
	requireNoPage(t, out, raftURL)
}

func TestAINote_BriefOpeningWithTheOldNoteFails(t *testing.T) {
	b := buildtest.ValidBrief()
	b.Body = oldAINote + "\n" + buildtest.ValidBriefBody

	out, res := buildtest.Build(t, buildtest.Config(), buildtest.ValidTree().WithBrief(b))

	buildtest.RequireErrorRule(t, res, briefFile, "ai-note", "ADR-0014")
	if buildtest.Exists(out, briefURL) {
		t.Errorf("the Brief carrying the AI note was published")
	}
}

func TestAINote_StaticPageQuotingTheOldNoteFails(t *testing.T) {
	policy := buildtest.TransparencyPages()[2]
	policy.Body = "Pieces carry this note:\n\n" + oldAINote

	_, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithPage(policy))

	buildtest.RequireErrorRule(t, res, "pages/editorial-policy.md", "ai-note", "ADR-0014")
}

// The rewritten Editorial Policy keeps the true process statements:
// claim-by-claim verification and the publishing pace it sets.
func TestAINote_StaticPageDescribingVerificationBuildsClean(t *testing.T) {
	policy := buildtest.TransparencyPages()[2]
	policy.Body = "## Verification\n\nA draft is only the start: the Author researches every claim against its first-party source. We publish only as many pieces as the Author can verify this way.\n"

	out, res := buildtest.Build(t, buildtest.Config(), treeWithTransparencyPages().WithPage(policy))

	buildtest.RequireNoProblemRule(t, res, "pages/editorial-policy.md", "ai-note")
	if !buildtest.Exists(out, "/editorial-policy/") {
		t.Errorf("the Editorial Policy was not published: %v", res.Errors)
	}
}

func TestAINote_StaticPageDescriptionFails(t *testing.T) {
	policy := buildtest.TransparencyPages()[2]
	policy.Description = "How HastCode sources, writes, discloses AI assistance and corrects."

	_, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithPage(policy))

	buildtest.RequireErrorRule(t, res, "pages/editorial-policy.md", "ai-note", "description")
}

func TestAINote_AuthorBioFails(t *testing.T) {
	a := buildtest.DefaultAuthor()
	a.Body = "Pieces under this name are drafted with AI assistance, then verified.\n"

	_, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithAuthor(a))

	buildtest.RequireErrorRule(t, res, "authors/ada.md", "ai-note", "ADR-0014")
}

func TestAINote_PieceAboutAIAsATopicBuildsClean(t *testing.T) {
	e := buildtest.ValidExplainer()
	e.Title = "How AI agents route a tool call"
	e.Description = "An AI assistant, an LLM and AI-assisted development: what each one runs."
	e.Body = withOpening("AI-assisted development and AI agents turn this inside out: an AI assistant drafts a plan, and [pgrust was rebuilt with AI](https://malisper.me/pgrust-rebuilding-postgres-in-rust-with-ai/).\n")

	out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

	buildtest.RequireNoProblemRule(t, res, raftFile, "ai-note")
	if !buildtest.Exists(out, raftURL) {
		t.Errorf("the Explainer about AI was not published: %v", res.Errors)
	}
}

func TestAINote_EveryWordingAndFieldFails(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit func(*buildtest.Explainer)
	}{
		{"note split over two lines", func(e *buildtest.Explainer) {
			e.Body = withOpening("> Drafted with AI\n> assistance, researched and edited by the team.\n")
		}},
		{"AI assistance in prose", func(e *buildtest.Explainer) {
			e.Body = withOpening("This piece was written with AI assistance.\n")
		}},
		{"drafted with AI, lower case", func(e *buildtest.Explainer) { e.Body = withOpening("drafted with ai, then checked.\n") }},
		{"hyphenated in the description", func(e *buildtest.Explainer) {
			e.Description = "Raft elections, with an AI-assistance note for the Reader."
		}},
		{"in the title", func(e *buildtest.Explainer) { e.Title = "Raft, drafted with AI" }},
		{"in an HTML comment the twin carries", func(e *buildtest.Explainer) {
			e.Body = withOpening("<!-- Drafted with AI assistance -->\n")
		}},
		{"written with the help of AI", func(e *buildtest.Explainer) {
			e.Body = withOpening("This piece was written with the help of AI, then checked.\n")
		}},
		{"the old Editorial Policy sentence", func(e *buildtest.Explainer) {
			e.Body = withOpening("Drafts may be written with the help of AI tools.\n")
		}},
		{"AI-generated about this article", func(e *buildtest.Explainer) {
			e.Body = withOpening("This article is AI-generated and reviewed by a human.\n")
		}},
		{"an AI-generated label on its own", func(e *buildtest.Explainer) { e.Body = withOpening("*AI-generated.*\n") }},
		{"generated with Claude, opening the sentence", func(e *buildtest.Explainer) {
			e.Body = withOpening("> Generated with Claude and edited by the team.\n")
		}},
		{"written using an LLM, opening the sentence", func(e *buildtest.Explainer) {
			e.Body = withOpening("Written using an LLM, verified by the Author.\n")
		}},
		{"AI-drafted with a non-breaking hyphen", func(e *buildtest.Explainer) {
			e.Description = "This explainer was AI\u2011drafted from the Raft paper."
		}},
		{"AI assistance with an en dash", func(e *buildtest.Explainer) {
			e.Body = withOpening("This piece was made with AI\u2013assistance.\n")
		}},
		{"drafted with AI across a non-breaking space", func(e *buildtest.Explainer) {
			e.Title = "Raft, drafted with\u00a0AI"
		}},
		{"HastCode drafting with ChatGPT", func(e *buildtest.Explainer) {
			e.Body = withOpening("HastCode produces its pieces with ChatGPT.\n")
		}},
		{"in the hero alt text", func(e *buildtest.Explainer) {
			e.HeroAlt = "A Raft cluster diagram. This page was generated by AI."
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e := buildtest.ValidExplainer()
			tc.edit(&e)

			out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

			buildtest.RequireErrorRule(t, res, raftFile, "ai-note", "ADR-0014")
			requireNoPage(t, out, raftURL)
		})
	}
}

// Wording about AI as a topic, including "AI assistance" itself, never
// fails: only wording about how this Piece was made does. Source titles
// are other publishers' titles, so they are not read at all.
func TestAINote_TopicWordingBuildsClean(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit func(*buildtest.Explainer)
	}{
		{"AI assistance as a product feature", func(e *buildtest.Explainer) {
			e.Title = "Copilot's AI assistance in code review"
			e.Body = withOpening("GitHub ships AI assistance in code review: Copilot drafts a comment on each diff.\n")
		}},
		{"AI-generated tests in a quote", func(e *buildtest.Explainer) {
			e.Body = withOpening("AI-generated test suites, Fowler writes, \"miss the edge cases a human would.\"\n")
		}},
		{"a question about another project", func(e *buildtest.Explainer) {
			e.Body = withOpening("Was pgrust written by AI? Its author says Claude wrote most of the port.\n")
		}},
		{"another project's process", func(e *buildtest.Explainer) {
			e.Body = withOpening("An account named claude opened 40 pull requests, so we infer part of the work is AI-assisted.\n")
		}},
		{"a Source title", func(e *buildtest.Explainer) {
			e.Sources = []buildtest.Source{{Title: "Drafted with AI assistance: how Copilot reviews code", URL: "https://github.blog/copilot-review/"}}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e := buildtest.ValidExplainer()
			tc.edit(&e)

			out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

			buildtest.RequireNoProblemRule(t, res, raftFile, "ai-note")
			if !buildtest.Exists(out, raftURL) {
				t.Errorf("the Explainer was not published: %v", res.Errors)
			}
		})
	}
}

// site.yaml's Vertical descriptions show on the home page, the hubs and
// llms.txt, so they are read too.
func TestAINote_VerticalDescriptionFails(t *testing.T) {
	cfg := buildtest.Config()
	cfg.Verticals[0].Description = "Consensus and replication, drafted with AI assistance."

	_, res := buildtest.Build(t, cfg, buildtest.ValidTree())

	if !slices.ContainsFunc(res.Errors, func(p builder.Problem) bool {
		return p.Rule == "ai-note" && strings.Contains(p.Message, "site.yaml") && strings.Contains(p.Message, "distributed-systems")
	}) {
		t.Errorf("no ai-note error naming site.yaml and the Vertical:\n%s", res)
	}
}

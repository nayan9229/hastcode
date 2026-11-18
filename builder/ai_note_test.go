package builder_test

import (
	"strings"
	"testing"

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

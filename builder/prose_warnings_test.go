package builder_test

import (
	"strings"
	"testing"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildtest"
)

// The prose warnings (house-byline-voice.md rule 5) turn the voice rules
// into numbers the build checks: sentence-length, banned-word and
// question-heading. Like section-length they only ever warn.

// proseRules are the three prose warnings. The filled kind templates and
// the Repo Profile template must fire none of them.
var proseRules = []string{"sentence-length", "banned-word", "question-heading"}

// proseHeading is the H2 the prose fixtures sit under.
const proseHeading = "Terms advance monotonically"

// withProse returns ValidExplainer with an extra H2 section, ## heading,
// holding body, inserted before Sources.
func withProse(heading, body string) buildtest.Explainer {
	e := buildtest.ValidExplainer()
	e.Body = strings.Replace(buildtest.ValidExplainerBody, "## Sources", "## "+heading+"\n\n"+body+"\n\n## Sources", 1)
	return e
}

// sentence returns one sentence of exactly n words, starting "Leaders keep".
func sentence(n int) string {
	return "Leaders keep" + strings.Repeat(" term", n-2) + "."
}

// buildProse builds ValidTree with the Explainer replaced by withProse.
func buildProse(t *testing.T, heading, body string) builder.Result {
	t.Helper()
	_, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(withProse(heading, body)))
	requireOnlyWarns(t, res)
	return res
}

// requireOnlyWarns fails the test if the build failed: the prose rules
// warn, never fail.
func requireOnlyWarns(t *testing.T, res builder.Result) {
	t.Helper()
	if res.Failed() {
		t.Fatalf("a prose warning failed the build; it must only warn:\n%s", res)
	}
}

func TestProseWarnings_SentenceOver30WordsWarns(t *testing.T) {
	res := buildProse(t, proseHeading, "Raft is short. "+sentence(31)+" Then it is short again.")

	buildtest.RequireWarningRule(t, res, raftFile, "sentence-length", "## "+proseHeading)
	buildtest.RequireWarningRule(t, res, raftFile, "sentence-length", `"Leaders keep term term`)
	buildtest.RequireWarningRule(t, res, raftFile, "sentence-length", "31 words")
}

func TestProseWarnings_SentenceOf30WordsDoesNotWarn(t *testing.T) {
	res := buildProse(t, proseHeading, "Raft is short. "+sentence(30)+" Then it is short again.")

	buildtest.RequireNoProblemRule(t, res, raftFile, "sentence-length")
}

func TestProseWarnings_EachBannedWordWarns(t *testing.T) {
	for _, tc := range []struct{ prose, word string }{
		{"Raft was a revolutionary design.", "revolutionary"},
		{"Raft was a game-changing design.", "game-changing"},
		{"Raft has no magic in it.", "magic"},
		{"Raft simply elects a leader.", "simply"},
		{"Raft makes elections easy to follow.", "easy"},
		{"Raft is an easy-to-read protocol.", "easy"},
		{"Raft obviously needs a majority.", "obviously"},
		{"It is important to note that Raft needs a majority.", "it is important to note"},
		{"Please note that Raft needs a majority.", "please note"},
		{"To fix it, just add a node.", "just add"},
		{"Just run the election again.", "just run"},
	} {
		t.Run(tc.word, func(t *testing.T) {
			res := buildProse(t, proseHeading, tc.prose)

			buildtest.RequireWarningRule(t, res, raftFile, "banned-word", "## "+proseHeading)
			buildtest.RequireWarningRule(t, res, raftFile, "banned-word", `"`+tc.word+`"`)
		})
	}
}

func TestProseWarnings_QuestionH2WarnsInAConcept(t *testing.T) {
	res := buildProse(t, "Why do terms only go up?", "Raft is short.")

	buildtest.RequireWarningRule(t, res, raftFile, "question-heading", "## Why do terms only go up?")
}

func TestProseWarnings_QuestionH2WarnsInABrief(t *testing.T) {
	b := buildtest.ValidBrief()
	b.Body = strings.Replace(buildtest.ValidBriefBody, "## Why it matters", "## Why does it matter?", 1)

	_, res := buildtest.Build(t, buildtest.Config(), buildtest.ValidTree().WithBrief(b))

	buildtest.RequireWarningRule(t, res, briefFile, "question-heading", "## Why does it matter?")
	// The Brief skeleton fails any H2 but its two literal ones; the
	// question itself is never an error.
	for _, p := range res.Errors {
		if p.Rule != "brief-sections" {
			t.Errorf("want only the Brief skeleton's brief-sections error, got %s", p.Error())
		}
	}
}

func TestProseWarnings_RepoProfileFAQQuestionsDoNotWarn(t *testing.T) {
	_, res := buildtest.Build(t, buildtest.Config(), repoProfileTree())

	requireOnlyWarns(t, res)
	for _, rule := range proseRules {
		buildtest.RequireNoProblemRule(t, res, repoProfileFile, rule)
	}
}

func TestProseWarnings_CodeTablesBlockquotesAndLinkURLsAreNeverChecked(t *testing.T) {
	long := strings.TrimSuffix(sentence(40), ".") + " simply and easy."
	for name, body := range map[string]string{
		"code fence":    "```go\n// " + long + "\n```",
		"mermaid fence": "```mermaid\n%% title: " + long + "\nflowchart LR\n  A[magic] --> B[just add]\n```",
		"inline code":   "Run `" + long + "` once.",
		"table":         "| Step | What |\n|---|---|\n| One | " + long + " |",
		"blockquote":    "> " + long,
		"link URL":      "See [the paper](https://example.com/simply/easy/magic/just-add/" + strings.ReplaceAll(long, " ", "-") + ").",
	} {
		t.Run(name, func(t *testing.T) {
			res := buildProse(t, proseHeading, "Raft is short.\n\n"+body)

			for _, rule := range proseRules {
				buildtest.RequireNoProblemRule(t, res, raftFile, rule)
			}
		})
	}
}

func TestProseWarnings_JustInItsOtherSensesAndQuotedWordsDoNotWarn(t *testing.T) {
	for _, prose := range []string{
		"A follower is just as likely to win.",
		"The timeout is just under 300 ms.",
		"The vote arrives just before the timeout.",
		"The leader steps in just in time.",
		`The paper calls it "simply a log" and moves on.`,
		"The README says “just add a node” and stops there.",
		"The paper calls it ‘simply a log’ and moves on.",
		"The README’s author calls it ‘simply a log’ and moves on.",
		"The job is just run once a day.",
		"Each task was just added to the queue.",
		"The file starts with a magic number.",
		"Each segment opens with four magic bytes.",
	} {
		t.Run(prose, func(t *testing.T) {
			res := buildProse(t, proseHeading, prose)

			buildtest.RequireNoProblemRule(t, res, raftFile, "banned-word")
		})
	}
}

func TestProseWarnings_ApostrophesAreNotQuotes(t *testing.T) {
	// A curly apostrophe never opens a quotation, so a banned word after
	// one still warns.
	res := buildProse(t, proseHeading, "Raft’s election is simply a vote.")

	buildtest.RequireWarningRule(t, res, raftFile, "banned-word", `"simply"`)
}

func TestProseWarnings_SentenceEnds(t *testing.T) {
	for _, tc := range []struct {
		name, prose string
		warn        bool
	}{
		// "Fig. 3" and "No. 5" are abbreviations before a number, not
		// sentence ends: 16 + 16 words stay one 32-word sentence.
		{"Fig. before a number", sentenceNoStop(16) + " in Fig. 3 then" + strings.Repeat(" term", 12) + ".", true},
		{"No. before a number", sentenceNoStop(16) + " in No. 5 then" + strings.Repeat(" term", 12) + ".", true},
		// "no." before a capital still ends a sentence.
		{"no. ends a sentence", sentenceNoStop(15) + " no. Then" + strings.Repeat(" term", 15) + ".", false},
		// A next sentence that opens with a lower-case brand (iOS, gRPC)
		// is still a new sentence.
		{"lower-case brand starts a sentence", sentence(20) + " iOS" + strings.Repeat(" term", 18) + ".", false},
		{"gRPC starts a sentence", sentence(20) + " gRPC" + strings.Repeat(" term", 18) + ".", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res := buildProse(t, proseHeading, tc.prose)

			if tc.warn {
				buildtest.RequireWarningRule(t, res, raftFile, "sentence-length", "## "+proseHeading)
			} else {
				buildtest.RequireNoProblemRule(t, res, raftFile, "sentence-length")
			}
		})
	}
}

// sentenceNoStop is sentence(n) without its full stop.
func sentenceNoStop(n int) string {
	return strings.TrimSuffix(sentence(n), ".")
}

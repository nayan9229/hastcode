// Package daily tests the Author's check list that `make daily` writes
// (scripts/daily/report.py, ADR-0009). report.py is Python, so this test
// runs it with python3 on a fixture run directory (testdata/sample-run):
// no Claude session, no network. `make test` runs it with the Go tests,
// in GitHub Actions too.
package daily

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const day = "2026-11-17"

// writeQueue runs report.py on the sample run (one Concept and one Brief
// drafted, one Case Study not viable) from a scratch repo root and returns
// the review queue it wrote.
func writeQueue(t *testing.T) string {
	t.Helper()
	py, err := exec.LookPath("python3")
	if err != nil {
		t.Fatal("make daily's check list is written by python3, which is not on PATH")
	}
	report, err := filepath.Abs("report.py")
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	run := filepath.Join(root, ".intake", "daily", "sample-run")
	if err := os.MkdirAll(run, 0o755); err != nil {
		t.Fatal(err)
	}
	fixtures, err := filepath.Glob(filepath.Join("testdata", "sample-run", "*.json"))
	if err != nil || len(fixtures) == 0 {
		t.Fatalf("no fixtures: %v", err)
	}
	for _, f := range fixtures {
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(run, filepath.Base(f)), b, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// The two drafted Pieces exist on disk; the not-viable one was removed.
	for _, p := range []string{"content/explainers/distributed-systems/retry-storms.md", "content/briefs/acme-queue-outage.md"} {
		full := filepath.Join(root, p)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("---\ndraft: true\n---\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command(py, report, filepath.Join(".intake", "daily", "sample-run"), day)
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("report.py: %v\n%s", err, out)
	}
	b, err := os.ReadFile(filepath.Join(root, ".scratch", "daily", day, "review-queue.md"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// entry returns the queue's text from the "### <title prefix>" heading to
// the next "### " heading or the end of the run.
func entry(t *testing.T, queue, title string) string {
	t.Helper()
	i := strings.Index(queue, "\n### "+title)
	if i < 0 {
		t.Fatalf("no entry %q in the queue:\n%s", title, queue)
	}
	rest := queue[i+1:]
	if j := strings.Index(rest[4:], "\n### "); j >= 0 {
		rest = rest[:4+j]
	}
	return rest
}

// The self-edit lines a test looks for, as docs/research/article-writing-craft.md
// §4 words them, unticked.
const (
	pointBox   = "- [ ] Point: I can say the point in one sentence"
	titleBox   = "- [ ] Top of the page: Title: key term in the first two words"
	hookBox    = "- [ ] Top of the page: Hook: one real, linked moment"
	nutGrafBox = "- [ ] Top of the page: Nut graf:"
	tldrBox    = "- [ ] Top of the page: TL;DR:"
	perH2Box   = "- [ ] Structure: Each H2 is under 250 words"
	workedBox  = "- [ ] Teaching: At least one worked example with real numbers"
	diagramBox = "- [ ] Teaching: The text walks the diagram box by box"
	readAloud  = "- [ ] Last pass: Read it aloud once."
	link       = "](../../../docs/editorial/self-edit-checklist.md)"
	// Brief-only lines (the Brief writing guide's rules).
	briefCapBox  = "- [ ] Structure: The whole Brief is 250 words at most"
	briefOpenBox = `- [ ] Top of the page: Nothing comes before "What happened".`
)

func TestBriefEntryLeavesOutLinesThatDoNotApplyToBriefs(t *testing.T) {
	e := entry(t, writeQueue(t), "Acme's queue outage")
	for _, want := range []string{pointBox, titleBox, readAloud, briefCapBox, briefOpenBox, link} {
		if !strings.Contains(e, want) {
			t.Errorf("Brief entry lacks %q:\n%s", want, e)
		}
	}
	for _, not := range []string{hookBox, nutGrafBox, tldrBox, perH2Box, workedBox, diagramBox} {
		if strings.Contains(e, not) {
			t.Errorf("Brief entry carries %q, which does not apply to a Brief:\n%s", not, e)
		}
	}
}

func TestNotDraftedLeadGetsNoSelfEditBoxes(t *testing.T) {
	q := writeQueue(t)
	e := entry(t, q, "Not drafted")
	if !strings.Contains(e, "thin-source") {
		t.Fatalf("the not-viable Lead is not listed under Not drafted:\n%s", e)
	}
	if strings.Contains(e, "- [ ]") || strings.Contains(e, "Self-edit") {
		t.Errorf("the Not drafted entry carries self-edit boxes:\n%s", e)
	}
	// One self-edit block per drafted Piece, two in all.
	if n := strings.Count(q, "**Self-edit before you verify**"); n != 2 {
		t.Errorf("want 2 self-edit blocks (Concept and Brief), got %d:\n%s", n, q)
	}
}

func TestOpeningSaysSelfEditComesBeforeVerificationAndDoesNotReplaceIt(t *testing.T) {
	q := writeQueue(t)
	open := q[:strings.Index(q, "\n## Run ")]
	for _, want := range []string{
		"self-edit boxes",
		"before you verify",
		"does not replace verification",
		"self-edit-checklist.md",
	} {
		if !strings.Contains(open, want) {
			t.Errorf("opening paragraph lacks %q:\n%s", want, open)
		}
	}
}

// make daily's drafter self-edits with the same checklist before the
// fact-check session runs.
func TestDrafterInstructionsNameTheChecklist(t *testing.T) {
	b, err := os.ReadFile("draft.md")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "docs/editorial/self-edit-checklist.md") {
		t.Errorf("scripts/daily/draft.md does not name docs/editorial/self-edit-checklist.md")
	}
}

// The boxes are the Author's: make daily writes every one unticked.
func TestNoBoxIsTicked(t *testing.T) {
	q := writeQueue(t)
	for _, ticked := range []string{"- [x]", "- [X]"} {
		if strings.Contains(q, ticked) {
			t.Errorf("the queue carries a ticked box %q:\n%s", ticked, q)
		}
	}
}

func TestExplainerEntryLeavesOutBriefOnlyLinesAndNotes(t *testing.T) {
	e := entry(t, writeQueue(t), "Retry storms")
	for _, not := range []string{briefCapBox, briefOpenBox, "not in a Brief", "Repo Profile:", "*("} {
		if strings.Contains(e, not) {
			t.Errorf("Concept entry carries %q:\n%s", not, e)
		}
	}
}

func TestExplainerEntryCarriesUntickedSelfEditBoxes(t *testing.T) {
	e := entry(t, writeQueue(t), "Retry storms")
	for _, want := range []string{pointBox, titleBox, hookBox, nutGrafBox, tldrBox, perH2Box, workedBox, diagramBox, readAloud, link} {
		if !strings.Contains(e, want) {
			t.Errorf("Explainer entry lacks %q:\n%s", want, e)
		}
	}
}

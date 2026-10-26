package intake_test

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"hastcode.com/site/internal/intake"
)

func TestRun_ListsNewLeadsWithTheSourcesKindAndVertical(t *testing.T) {
	feeds := newFeedServer(t)
	in := newIntake(t, exampleWatchlist(t, feeds))

	report := mustRun(t, in)

	want := []queuedLead{
		{Vertical: "Distributed Systems", Title: "How we lost a region for 40 minutes",
			Link: "https://engineering.example.com/posts/region-outage", Source: "Example Engineering",
			Kind: "case-study", Date: "2026-10-24 (1 day old)", New: true},
		{Vertical: "Distributed Systems", Title: "Raft in production: five years on",
			Link: "https://engineering.example.com/posts/raft-five-years", Source: "Example Engineering",
			Kind: "case-study", Date: "2026-10-19 (7 days old)", New: true},
		{Vertical: "Databases & Storage", Title: "Inside the write-ahead log",
			Link: "https://db.example.org/blog/wal", Source: "Example Database Blog",
			Kind: "concept", Date: "2026-10-25 (1 day old)", New: true},
		{Vertical: "Databases & Storage", Title: "Why our B-tree pages are 8 KiB",
			Link: "https://db.example.org/blog/btree-pages", Source: "Example Database Blog",
			Kind: "concept", Date: "2026-10-12 (13 days old)", New: true},
	}
	got := readQueue(t, in)
	if len(got) != len(want) {
		t.Fatalf("queue lists %d Leads %q, want %d", len(got), titles(got), len(want))
	}
	for i := range want {
		want[i].ID = got[i].ID // IDs are checked by their own test
		if got[i] != want[i] {
			t.Errorf("queue Lead %d:\n got %+v\nwant %+v", i, got[i], want[i])
		}
	}
	if len(report.New) != 4 {
		t.Errorf("report lists %d new Leads, want 4", len(report.New))
	}
	for _, l := range report.New {
		if !slices.ContainsFunc(got, func(q queuedLead) bool { return q.ID == l.ID }) {
			t.Errorf("new Lead %s (%q) is not in the queue", l.ID, l.Title)
		}
	}
}

func TestRun_SecondRunListsNothingNew(t *testing.T) {
	feeds := newFeedServer(t)
	in := newIntake(t, exampleWatchlist(t, feeds))
	first := readQueueAfterRun(t, in)

	report := mustRun(t, in)

	if len(report.New) != 0 {
		t.Errorf("second run lists %d new Leads %+v, want none", len(report.New), report.New)
	}
	second := readQueue(t, in)
	for _, l := range second {
		if l.New {
			t.Errorf("second run's queue marks %q new", l.Title)
		}
	}
	if len(second) != len(first) {
		t.Errorf("second queue holds %d open Leads, want the first run's %d", len(second), len(first))
	}
}

func TestRun_ListsOnlyTheItemPublishedSinceTheLastRun(t *testing.T) {
	feeds := newFeedServer(t)
	in := newIntake(t, exampleWatchlist(t, feeds))
	mustRun(t, in)
	feeds.setBody("/rss.xml", []byte(`<?xml version="1.0"?>
<rss version="2.0"><channel><title>Example Engineering</title>
  <item>
    <title>Cell-based architecture, one year in</title>
    <link>https://engineering.example.com/posts/cells</link>
    <guid isPermaLink="false">example-eng-1043</guid>
    <pubDate>Mon, 26 Oct 2026 07:00:00 +0000</pubDate>
  </item>
  <item>
    <title>How we lost a region for 40 minutes</title>
    <link>https://engineering.example.com/posts/region-outage</link>
    <guid isPermaLink="false">example-eng-1042</guid>
    <pubDate>Sat, 24 Oct 2026 14:30:00 +0000</pubDate>
  </item>
</channel></rss>`))

	report := mustRun(t, in)

	if len(report.New) != 1 || report.New[0].Title != "Cell-based architecture, one year in" {
		t.Fatalf("second run's new Leads = %+v, want only the cell-based architecture post", report.New)
	}
	q := readQueue(t, in)
	if q[0].Title != "Cell-based architecture, one year in" || !q[0].New || q[0].Date != "2026-10-26 (today)" {
		t.Errorf("queue starts with %+v, want the new post first, marked new, today", q[0])
	}
	if len(q) != 5 {
		t.Errorf("queue holds %d Leads %q, want 5 (four from the first run, one new)", len(q), titles(q))
	}
}

func TestRun_LeadIDsAreStableAcrossRunsAndShortEnoughToType(t *testing.T) {
	feeds := newFeedServer(t)
	watchlist := exampleWatchlist(t, feeds)
	a := readQueueAfterRun(t, newIntake(t, watchlist))
	later := newIntake(t, watchlist) // a fresh state directory, a week on
	later.Now = func() time.Time { return testNow.AddDate(0, 0, 7) }
	b := readQueueAfterRun(t, later)

	if len(a) != 4 || len(b) != 4 {
		t.Fatalf("queues list %d and %d Leads, want 4 each", len(a), len(b))
	}
	ids := map[string]bool{}
	for i := range a {
		if a[i].ID != b[i].ID {
			t.Errorf("%q has ID %q in one run and %q in another", a[i].Title, a[i].ID, b[i].ID)
		}
		if !regexp.MustCompile(`^[a-z0-9]{1,8}$`).MatchString(a[i].ID) {
			t.Errorf("%q has ID %q, want at most 8 lowercase letters and digits", a[i].Title, a[i].ID)
		}
		if ids[a[i].ID] {
			t.Errorf("ID %q is listed twice", a[i].ID)
		}
		ids[a[i].ID] = true
	}
}

func TestRun_WritesNothingOutsideTheStateDirectory(t *testing.T) {
	feeds := newFeedServer(t)
	root := t.TempDir()
	contentDir := filepath.Join(root, "content")
	piece := filepath.Join(contentDir, "explainers", "distributed-systems", "raft.md")
	if err := os.MkdirAll(filepath.Dir(piece), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(piece, []byte("---\ntitle: Raft\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	in := newIntake(t, exampleWatchlist(t, feeds))
	in.ContentDir = contentDir
	in.StateDir = filepath.Join(root, ".intake")

	mustRun(t, in)

	if got, want := filesUnder(t, root), []string{
		".intake/" + intake.QueueFile,
		".intake/" + intake.StateFile,
		"content/explainers/distributed-systems/raft.md",
	}; !slices.Equal(got, want) {
		t.Errorf("files after the run = %q, want %q", got, want)
	}
}

func TestRun_RefusesAStateDirectoryInsideTheContentTree(t *testing.T) {
	feeds := newFeedServer(t)
	contentDir := t.TempDir()
	in := newIntake(t, exampleWatchlist(t, feeds))
	in.ContentDir = contentDir
	in.StateDir = filepath.Join(contentDir, ".intake")

	_, err := in.Run(context.Background())

	if err == nil || !strings.Contains(err.Error(), "inside the content tree") {
		t.Fatalf("Run = %v, want a refusal naming the content tree", err)
	}
	if files := filesUnder(t, contentDir); len(files) != 0 {
		t.Errorf("the refused run wrote %q under the content tree", files)
	}
	if n := feeds.requestCount(); n != 0 {
		t.Errorf("the refused run made %d requests", n)
	}
}

func readQueueAfterRun(t *testing.T, in *intake.Intake) []queuedLead {
	t.Helper()
	mustRun(t, in)
	return readQueue(t, in)
}

// filesUnder lists every file under dir, dir-relative with / separators,
// sorted.
func filesUnder(t *testing.T, dir string) []string {
	t.Helper()
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		files = append(files, filepath.ToSlash(rel))
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	slices.Sort(files)
	return files
}

package intake_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"hastcode.com/site/internal/intake"
)

// at moves in's clock to t.
func at(in *intake.Intake, t time.Time) {
	in.Now = func() time.Time { return t }
}

func TestRun_OpenLeadIsKeptAtSixDays(t *testing.T) {
	feeds := newFeedServer(t)
	in := newIntake(t, exampleWatchlist(t, feeds))
	first := readQueueAfterRun(t, in)

	at(in, testNow.Add(7*24*time.Hour-time.Minute))
	report := mustRun(t, in)

	if got := readQueue(t, in); len(got) != len(first) {
		t.Errorf("a run 6 days 23:59 after the first lists %d Leads %q, want all %d still open", len(got), titles(got), len(first))
	}
	if report.Expired != 0 {
		t.Errorf("report: %d expired, want 0", report.Expired)
	}
}

func TestRun_OpenLeadExpiresAtSevenDays(t *testing.T) {
	feeds := newFeedServer(t)
	in := newIntake(t, exampleWatchlist(t, feeds))
	first := readQueueAfterRun(t, in)
	if err := in.Dismiss(context.Background(), first[0].ID); err != nil {
		t.Fatal(err)
	}

	at(in, testNow.Add(7*24*time.Hour))
	report := mustRun(t, in)

	if got := readQueue(t, in); len(got) != 0 {
		t.Errorf("a run 7 days after the first lists %q, want every undecided Lead expired", titles(got))
	}
	if report.Expired != len(first)-1 {
		t.Errorf("report: %d expired, want the %d undecided Leads", report.Expired, len(first)-1)
	}
	if report.Open != 0 {
		t.Errorf("report: %d open, want 0", report.Open)
	}
	st := readState(t, in)
	if d := st[first[0].ID].Decision; d != "dismissed" {
		t.Errorf("dismissed Lead %s: decision %q after expiry, want it still dismissed", first[0].ID, d)
	}
	for _, l := range first[1:] {
		if d := st[l.ID].Decision; d != "expired" {
			t.Errorf("Lead %s (%q): decision %q, want expired", l.ID, l.Title, d)
		}
	}
}

func TestRun_ANewLeadIsNotExpiredByTheRunThatListsIt(t *testing.T) {
	feeds := newFeedServer(t)
	in := newIntake(t, writeWatchlist(t, feeds, `
sources:
  - name: Example Engineering
    site: https://engineering.example.com/
    feed: "{{server}}/rss.xml"
    vertical: distributed-systems
    kind: case-study
    reason: Postmortems.
`))
	mustRun(t, in)
	feeds.setBody("/rss.xml", []byte(`<?xml version="1.0"?>
<rss version="2.0"><channel><title>Example Engineering</title>
<item><title>Retries considered harmful</title><link>https://engineering.example.com/posts/retries</link></item>
</channel></rss>`))

	at(in, testNow.Add(8*24*time.Hour))
	report := mustRun(t, in)

	got := readQueue(t, in)
	if len(got) != 1 || got[0].Title != "Retries considered harmful" || !got[0].New {
		t.Errorf("queue %q, want only the new Lead, marked new", titles(got))
	}
	if report.Expired != 2 {
		t.Errorf("report: %d expired, want the first run's 2", report.Expired)
	}
}

func TestRun_StateRecordsEachLeadsDecisionAndFirstSeenDate(t *testing.T) {
	feeds := newFeedServer(t)
	in := newIntake(t, exampleWatchlist(t, feeds))
	first := readQueueAfterRun(t, in)
	if err := in.Dismiss(context.Background(), first[0].ID); err != nil {
		t.Fatal(err)
	}

	at(in, testNow.Add(24*time.Hour))
	mustRun(t, in)

	st := readState(t, in)
	if len(st) != len(first) {
		t.Fatalf("state holds %d Leads, want %d", len(st), len(first))
	}
	for i, l := range first {
		want := "open"
		if i == 0 {
			want = "dismissed"
		}
		got := st[l.ID]
		if got.Decision != want {
			t.Errorf("Lead %s: decision %q, want %q", l.ID, got.Decision, want)
		}
		if got.FirstSeen != "2026-10-26T09:00:00Z" {
			t.Errorf("Lead %s: first seen %q, want the first run's 2026-10-26T09:00:00Z", l.ID, got.FirstSeen)
		}
	}
}

// stateLead is a Lead's record in the seen-state file (README "Intake").
type stateLead struct {
	Decision  string `json:"decision"`
	FirstSeen string `json:"first_seen"`
}

// readState reads the Leads of in's seen-state file, by ID.
func readState(t *testing.T, in *intake.Intake) map[string]stateLead {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(in.StateDir, intake.StateFile))
	if err != nil {
		t.Fatal(err)
	}
	var st struct {
		Leads map[string]stateLead `json:"leads"`
	}
	if err := json.Unmarshal(data, &st); err != nil {
		t.Fatal(err)
	}
	return st.Leads
}

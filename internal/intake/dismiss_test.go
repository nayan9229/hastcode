package intake_test

import (
	"context"
	"strings"
	"testing"
)

func TestDismiss_DismissedLeadIsAbsentFromLaterQueues(t *testing.T) {
	feeds := newFeedServer(t)
	in := newIntake(t, exampleWatchlist(t, feeds))
	first := readQueueAfterRun(t, in)
	dismissed := first[0]

	if err := in.Dismiss(context.Background(), dismissed.ID); err != nil {
		t.Fatalf("Dismiss(%s): %v", dismissed.ID, err)
	}

	for run := 1; run <= 2; run++ {
		report := mustRun(t, in)
		for _, l := range readQueue(t, in) {
			if l.ID == dismissed.ID {
				t.Errorf("run %d after dismissing: queue still lists %s (%q)", run, l.ID, l.Title)
			}
		}
		if report.Open != len(first)-1 {
			t.Errorf("run %d after dismissing: %d open Leads, want %d", run, report.Open, len(first)-1)
		}
		if len(report.New) != 0 {
			t.Errorf("run %d after dismissing lists the dismissed item again: %+v", run, report.New)
		}
	}
}

func TestDismiss_RefusesAnUnknownID(t *testing.T) {
	feeds := newFeedServer(t)
	in := newIntake(t, exampleWatchlist(t, feeds))
	mustRun(t, in)

	err := in.Dismiss(context.Background(), "zzzzzzzz")

	if err == nil || !strings.Contains(err.Error(), "zzzzzzzz") {
		t.Errorf("Dismiss of an unknown ID: error %v, want one naming the ID", err)
	}
}

func TestDismiss_RefusesALeadAlreadyDecided(t *testing.T) {
	feeds := newFeedServer(t)
	in := newIntake(t, exampleWatchlist(t, feeds))
	id := readQueueAfterRun(t, in)[0].ID
	if err := in.Dismiss(context.Background(), id); err != nil {
		t.Fatal(err)
	}

	err := in.Dismiss(context.Background(), id)

	if err == nil || !strings.Contains(err.Error(), "dismissed") {
		t.Errorf("second Dismiss(%s): error %v, want one saying it is already dismissed", id, err)
	}
}

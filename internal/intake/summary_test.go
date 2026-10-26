package intake_test

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestRun_SummaryCountsNewLeadsPerVerticalAndKind(t *testing.T) {
	feeds := newFeedServer(t)
	in := newIntake(t, exampleWatchlist(t, feeds))

	var out strings.Builder
	mustRun(t, in).WriteSummary(&out)

	// Every Vertical and kind is listed, zeros included, so a starving
	// Vertical shows.
	for _, want := range []string{
		"new by Vertical: Distributed Systems 2, Databases & Storage 2, Languages & Runtimes 0\n",
		"new by kind: case-study 2, concept 2, comparison 0, brief 0\n",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("summary lacks %q:\n%s", want, out.String())
		}
	}
}

func TestRun_SummaryCountsOnlyThisRunsNewLeadsAndReportsExpired(t *testing.T) {
	feeds := newFeedServer(t)
	in := newIntake(t, exampleWatchlist(t, feeds))
	first := readQueueAfterRun(t, in)
	if err := in.Dismiss(context.Background(), first[0].ID); err != nil {
		t.Fatal(err)
	}

	at(in, testNow.Add(7*24*time.Hour))
	var out strings.Builder
	mustRun(t, in).WriteSummary(&out)

	for _, want := range []string{
		"new by Vertical: Distributed Systems 0, Databases & Storage 0, Languages & Runtimes 0\n",
		"new by kind: case-study 0, concept 0, comparison 0, brief 0\n",
		"expired 3 Leads left open 7 days\n",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("summary lacks %q:\n%s", want, out.String())
		}
	}
}

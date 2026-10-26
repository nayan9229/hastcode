package intake_test

import (
	"slices"
	"testing"
	"time"
)

// The age cutoff: a run lists only entries published within MaxAge (14
// days by default) before its clock, so a first run over a whole Watchlist
// does not flood the queue with every Source's back catalogue.

// agedFeed is a feed whose entries are 1, 14 (exactly), 15 and 200 days
// older than testNow (2026-10-26 09:00 UTC), plus one with no date.
const agedFeed = `<?xml version="1.0"?>
<rss version="2.0"><channel><title>Aged Blog</title>
  <item><title>A day old</title><link>https://aged.example.com/1</link>
    <pubDate>Sun, 25 Oct 2026 09:00:00 +0000</pubDate></item>
  <item><title>Exactly two weeks old</title><link>https://aged.example.com/14</link>
    <pubDate>Mon, 12 Oct 2026 09:00:00 +0000</pubDate></item>
  <item><title>Fifteen days old</title><link>https://aged.example.com/15</link>
    <pubDate>Sun, 11 Oct 2026 09:00:00 +0000</pubDate></item>
  <item><title>From the archive</title><link>https://aged.example.com/200</link>
    <pubDate>Fri, 10 Apr 2026 09:00:00 +0000</pubDate></item>
  <item><title>Undated</title><link>https://aged.example.com/undated</link></item>
</channel></rss>`

const agedWatchlist = `sources:
  - name: Aged Blog
    site: https://aged.example.com/
    feed: "{{server}}/aged.xml"
    vertical: distributed-systems
    kind: concept
    reason: Has a long back catalogue.
`

func TestRun_ListsOnlyEntriesPublishedWithinMaxAge(t *testing.T) {
	feeds := newFeedServer(t)
	feeds.setBody("/aged.xml", []byte(agedFeed))
	in := newIntake(t, writeWatchlist(t, feeds, agedWatchlist))

	report := mustRun(t, in)

	want := []string{"A day old", "Exactly two weeks old", "Undated"}
	if got := titles(readQueue(t, in)); !slices.Equal(sorted(got), want) {
		t.Errorf("first run's queue = %q, want only %q (older entries are not Leads)", got, want)
	}
	if len(report.New) != len(want) {
		t.Errorf("first run: %d new Leads, want %d", len(report.New), len(want))
	}

	report = mustRun(t, in)

	if len(report.New) != 0 {
		t.Errorf("second run lists %+v, want nothing: old entries were seen, the undated one listed once", report.New)
	}
	if got := titles(readQueue(t, in)); !slices.Equal(sorted(got), want) {
		t.Errorf("second run's queue = %q, want still only %q", got, want)
	}
}

func TestRun_AWiderMaxAgeListsOlderEntriesOnce(t *testing.T) {
	feeds := newFeedServer(t)
	feeds.setBody("/aged.xml", []byte(agedFeed))
	in := newIntake(t, writeWatchlist(t, feeds, agedWatchlist))
	mustRun(t, in) // the default 14 days: the two old entries are seen, not listed

	in.MaxAge = 365 * 24 * time.Hour // a one-off wider run
	report := mustRun(t, in)

	var got []string
	for _, l := range report.New {
		got = append(got, l.Title)
	}
	if want := []string{"Fifteen days old", "From the archive"}; !slices.Equal(sorted(got), want) {
		t.Errorf("wider run's new Leads = %q, want %q", got, want)
	}
	if q := readQueue(t, in); len(q) != 5 {
		t.Errorf("wider run's queue lists %q, want all 5 entries open", titles(q))
	}

	in.MaxAge = 0 // back to the default
	if report := mustRun(t, in); len(report.New) != 0 {
		t.Errorf("the run after the wider one lists %+v, want nothing new", report.New)
	}
}

func sorted(s []string) []string {
	s = slices.Clone(s)
	slices.Sort(s)
	return s
}

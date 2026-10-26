package intake_test

import (
	"path/filepath"
	"slices"
	"testing"
)

// fixtureContent is a small content tree (testdata/content) whose Pieces
// cite three of the four example feed items: one in a Sources entry, one in
// a Brief's body, one in a Repo Profile's Fact Pack Sources. A page (not a
// Piece) cites the fourth.
var fixtureContent = filepath.Join("testdata", "content")

func TestRun_MarksLeadsAPieceAlreadyCitesAsCovered(t *testing.T) {
	feeds := newFeedServer(t)
	in := newIntake(t, exampleWatchlist(t, feeds))
	in.ContentDir = fixtureContent
	before := filesUnder(t, fixtureContent)

	got := readQueueAfterRun(t, in)

	want := map[string]bool{
		"How we lost a region for 40 minutes": true,  // Sources entry, normalised
		"Raft in production: five years on":   false, // only a near miss and a page cite it
		"Inside the write-ahead log":          true,  // a Brief's body
		"Why our B-tree pages are 8 KiB":      true,  // a Fact Pack's Sources
	}
	if len(got) != len(want) {
		t.Fatalf("queue lists %q, want the 4 example Leads", titles(got))
	}
	for _, l := range got {
		if l.Covered != want[l.Title] {
			t.Errorf("%q: covered %v, want %v", l.Title, l.Covered, want[l.Title])
		}
		if !l.New {
			t.Errorf("%q: not marked new; covered goes after new", l.Title)
		}
	}
	if after := filesUnder(t, fixtureContent); !slices.Equal(before, after) {
		t.Errorf("content tree changed by a run:\nbefore %q\n after %q", before, after)
	}
}

func TestRun_MarksNothingCoveredWithoutAContentTree(t *testing.T) {
	feeds := newFeedServer(t)
	in := newIntake(t, exampleWatchlist(t, feeds))

	for _, l := range readQueueAfterRun(t, in) {
		if l.Covered {
			t.Errorf("%q marked covered with no content tree", l.Title)
		}
	}
}

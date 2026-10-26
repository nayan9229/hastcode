package intake

import (
	"context"
	"time"
)

// Entry is what a Source published, as its fetcher read it (an RSS item,
// an Atom entry, a GitHub release), before Intake decides it is a Lead:
// Run makes a Lead of an entry not seen before and within the age cutoff
// (GLOSSARY "Lead").
type Entry struct {
	// Key identifies the entry within its Source over time: an RSS guid
	// (else its link), an Atom id, a release's id. With the Source's name
	// it makes the Lead ID (lead.go).
	Key       string
	Title     string
	Link      string
	Published time.Time // zero when the Source gives no date
}

// Fetched is what one fetch of a Source returned.
type Fetched struct {
	// Entries are every entry the Source currently lists; Run keeps the
	// ones not seen before.
	Entries []Entry
	// State is the Source's state to keep for the next run (validators,
	// …); Run stores it, with Newest carried forward (source_state.go).
	State SourceState
}

// sourceFetcher fetches one type of Source. prev is the state its last
// fetch returned (zero on the first run).
type sourceFetcher interface {
	fetch(ctx context.Context, in *Intake, src Source, prev SourceState) (Fetched, error)
}

// fetchers is the registry Run looks a Source's type up in: every
// SourceType has one.
var fetchers = map[SourceType]sourceFetcher{}

// registerFetcher adds the fetcher for one Source type, from that type's
// own file's init (fetch_feed.go, fetch_github.go). A new Source type also
// needs its Watchlist field, Source.Type and validation (watchlist.go).
func registerFetcher(t SourceType, f sourceFetcher) {
	if _, dup := fetchers[t]; dup {
		panic("intake: two fetchers for Source type " + string(t))
	}
	fetchers[t] = f
}

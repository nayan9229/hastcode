package intake

import (
	"context"
	"time"
)

// Item is one thing a Source published, as its fetcher read it.
type Item struct {
	// Key identifies the item within its Source over time: an RSS guid
	// (else its link), an Atom id, a release's id. With the Source's name
	// it makes the Lead ID (lead.go).
	Key       string
	Title     string
	Link      string
	Published time.Time // zero when the Source gives no date
}

// Fetched is what one fetch of a Source returned.
type Fetched struct {
	// Items are every item the Source currently lists; Run keeps the ones
	// not seen before.
	Items []Item
	// State is the Source's state to keep for the next run (validators,
	// …); Run stores it whole, replacing the previous one.
	State SourceState
}

// sourceFetcher fetches one type of Source. prev is the state its last
// fetch returned (zero on the first run).
type sourceFetcher interface {
	fetch(ctx context.Context, in *Intake, src Source, prev SourceState) (Fetched, error)
}

// fetchers is the registry Run looks a Source's type up in. A type with no
// fetcher is skipped with a note in the run summary.
var fetchers = map[SourceType]sourceFetcher{}

// registerFetcher adds the fetcher for one Source type. Each type
// registers itself from its own file's init (fetch_feed.go), so a new type
// is a new file, not an edit to a central list.
func registerFetcher(t SourceType, f sourceFetcher) {
	if _, dup := fetchers[t]; dup {
		panic("intake: two fetchers for Source type " + string(t))
	}
	fetchers[t] = f
}

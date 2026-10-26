package intake

import (
	"fmt"
	"io"
	"time"
)

// SilentAfter is how old a Source's newest item may be before the run
// summary flags the Source silent, a candidate to retire (spec "Intake").
const SilentAfter = 90 * 24 * time.Hour

// FailedSource is a Source whose fetch failed in a run, with why. One
// broken Source never stops the run (spec "Intake"); only a Watchlist or
// state error does.
type FailedSource struct {
	Source string
	Reason string
}

// SilentSource is a Source whose newest item is more than SilentAfter old.
type SilentSource struct {
	Source string
	Newest time.Time
	Days   int // whole days from Newest to the run
}

// fail records src's failed fetch.
func (r *Report) fail(in *Intake, src Source, err error) {
	r.Failed = append(r.Failed, FailedSource{Source: src.Name, Reason: err.Error()})
	in.logf("failed %s: %v", src.Name, err)
}

// fetchedState is the state to keep for src after a successful fetch: the
// fetcher's, with Newest carried over from prev and moved on by any newer
// dated item (a 304 lists none, so it keeps prev's). It flags src silent
// when that newest item is more than SilentAfter before now.
func (r *Report) fetchedState(src Source, prev SourceState, fetched Fetched, now time.Time) SourceState {
	st := fetched.State
	st.Newest = prev.Newest
	for _, it := range fetched.Items {
		if it.Published.After(st.Newest) {
			st.Newest = it.Published
		}
	}
	if !st.Newest.IsZero() && now.Sub(st.Newest) > SilentAfter {
		r.Silent = append(r.Silent, SilentSource{Source: src.Name, Newest: st.Newest, Days: int(now.Sub(st.Newest) / (24 * time.Hour))})
	}
	return st
}

// writeSourceHealth prints the failed and the silent Sources, one line
// each, for the run summary.
func (r *Report) writeSourceHealth(w io.Writer) {
	for _, f := range r.Failed {
		fmt.Fprintf(w, "failed %s: %s\n", f.Source, f.Reason)
	}
	for _, s := range r.Silent {
		fmt.Fprintf(w, "silent %s: newest item %s (%s ago)\n", s.Source, s.Newest.Format(time.DateOnly), plural(s.Days, "day"))
	}
}

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

// fail records src's failed fetch and logs it.
func (r *Report) fail(logf func(format string, args ...any), src Source, err error) {
	r.Failed = append(r.Failed, FailedSource{Source: src.Name, Reason: err.Error()})
	logf("failed %s: %v", src.Name, err)
}

// checkSilent flags src silent when its newest entry, published at newest
// (zero: none dated yet), is more than SilentAfter before now.
func (r *Report) checkSilent(src Source, newest, now time.Time) {
	if !newest.IsZero() && now.Sub(newest) > SilentAfter {
		r.Silent = append(r.Silent, SilentSource{Source: src.Name, Newest: newest, Days: wholeDays(now.Sub(newest))})
	}
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

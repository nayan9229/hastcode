package intake

import (
	"cmp"
	"net/http"
	"time"
)

// SourceState is what Intake keeps about one Source between runs, keyed by
// the Source's name in state.json "sources". Its fetcher returns it
// (Fetched.State) and gets it back on the next run; a Source whose fetch
// failed keeps the state it had.
type SourceState struct {
	// ETag and LastModified are the validators the Source's last full
	// response carried, kept verbatim and sent back as If-None-Match and
	// If-Modified-Since (RFC 9110 §13.1; research §5.2). Empty: not sent.
	ETag         string `json:"etag,omitempty"`
	LastModified string `json:"last_modified,omitempty"`
	// Newest is the publication date of the newest entry the Source has
	// listed (zero while none had a date). Run carries it forward
	// (withNewest), whatever the fetcher returns, and flags the Source
	// silent when it is more than SilentAfter old.
	Newest time.Time `json:"newest,omitzero"`
}

// withNewest is s (a fetcher's state) with Newest carried forward from
// the stored newest and moved on by any later dated entry; a 304 lists no
// entries, so it keeps the stored date.
func (s SourceState) withNewest(stored time.Time, entries []Entry) SourceState {
	s.Newest = stored
	for _, e := range entries {
		if e.Published.After(s.Newest) {
			s.Newest = e.Published
		}
	}
	return s
}

// revalidated is a stored validator after a 304: the value the 304
// repeats under name when it sends one (it may be fresher, RFC 9110
// §15.4.5), else the stored one.
func revalidated(h http.Header, name, stored string) string {
	return cmp.Or(h.Get(name), stored)
}

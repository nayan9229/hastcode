package intake

import "time"

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
	// Newest is the publication date of the newest item the Source has
	// listed (zero while none had a date). Run keeps it, whatever the
	// fetcher returns, and flags the Source silent when it is more than
	// SilentAfter old.
	Newest time.Time `json:"newest,omitzero"`
}

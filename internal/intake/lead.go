package intake

import (
	"crypto/sha256"
	"encoding/base32"
	"strings"
	"time"
)

// Lead is one item a Source published that the Author may turn into a
// Piece (GLOSSARY "Lead"). Its Kind and Vertical are the Source's, copied,
// never guessed. It is also the seen-state record (state.go), so its JSON
// names are the state file's format.
type Lead struct {
	ID        string    `json:"id"`
	Source    string    `json:"source"`
	Title     string    `json:"title"`
	Link      string    `json:"link"`
	Published time.Time `json:"published,omitzero"` // zero when the Source gives no date
	Kind      string    `json:"kind"`
	Vertical  string    `json:"vertical"`
	// FirstSeen is the date of the run that first listed it.
	FirstSeen time.Time `json:"first_seen"`
	// Decision is what the Author decided: open until picked or
	// dismissed, or expired.
	Decision Decision `json:"decision"`
}

// leadIDLength is the length of a Lead ID: eight base32 characters (40
// bits), short enough to type, with collisions negligible for the tens of
// thousands of Leads a Watchlist yields in years.
const leadIDLength = 8

var leadIDEncoding = base32.StdEncoding.WithPadding(base32.NoPadding)

// leadID is the stable ID of the entry key from the Source named source:
// the same entry gets the same ID on every run and in every state
// directory.
func leadID(source, key string) string {
	sum := sha256.Sum256([]byte(source + "\x00" + key))
	return strings.ToLower(leadIDEncoding.EncodeToString(sum[:]))[:leadIDLength]
}

func newLead(src Source, e Entry, now time.Time) Lead {
	return Lead{
		ID:        leadID(src.Name, e.Key),
		Source:    src.Name,
		Title:     e.Title,
		Link:      e.Link,
		Published: e.Published,
		Kind:      src.Kind,
		Vertical:  src.Vertical,
		FirstSeen: now.UTC().Truncate(time.Second),
		Decision:  Open,
	}
}

// date is the Lead's date for ordering and age: published, else first
// seen.
func (l Lead) date() time.Time {
	if l.Published.IsZero() {
		return l.FirstSeen
	}
	return l.Published
}

package intake_test

import (
	"context"
	"os"
	"strings"
	"testing"
)

// validSource is a Watchlist entry that passes; each case breaks it.
const validSource = `
  - name: Example Engineering
    site: https://engineering.example.com/
    feed: "{{server}}/rss.xml"
    vertical: distributed-systems
    kind: case-study
    reason: Postmortems with timelines.
`

func TestRun_RejectsAnInvalidWatchlistBeforeAnyRequest(t *testing.T) {
	cases := []struct {
		name    string
		sources string
		want    []string // each must appear in the error
	}{
		{"unknown vertical", strings.Replace(validSource, "distributed-systems", "distributed-sytems", 1),
			[]string{`source "Example Engineering"`, `unknown vertical "distributed-sytems"`}},
		{"unknown kind", strings.Replace(validSource, "kind: case-study", "kind: casestudy", 1),
			[]string{`source "Example Engineering"`, `unknown kind "casestudy"`}},
		{"repo-profile is not a Lead kind", strings.Replace(validSource, "kind: case-study", "kind: repo-profile", 1),
			[]string{`unknown kind "repo-profile"`}},
		{"both feed and repository", strings.Replace(validSource, "    vertical:", "    repository: example/runtime\n    vertical:", 1),
			[]string{`source "Example Engineering"`, "both feed and repository"}},
		{"neither feed nor repository", strings.Replace(validSource, `    feed: "{{server}}/rss.xml"`+"\n", "", 1),
			[]string{`source "Example Engineering"`, "neither feed nor repository"}},
		{"duplicate names", validSource + strings.Replace(validSource, "rss.xml", "atom.xml", 1),
			[]string{`source "Example Engineering"`, "duplicate name"}},
		{"missing reason", strings.Replace(validSource, "    reason: Postmortems with timelines.\n", "", 1),
			[]string{`source "Example Engineering"`, "reason is required"}},
		{"unknown key", strings.Replace(validSource, "    reason:", "    verticle: databases-storage\n    reason:", 1),
			[]string{"verticle"}},
		{"every problem reported", strings.Replace(strings.Replace(validSource, "kind: case-study", "kind: casestudy", 1), "distributed-systems", "nope", 1),
			[]string{`unknown kind "casestudy"`, `unknown vertical "nope"`}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			feeds := newFeedServer(t)
			in := newIntake(t, writeWatchlist(t, feeds, "sources:\n"+tc.sources))

			_, err := in.Run(context.Background())

			if err == nil {
				t.Fatal("Run passed an invalid Watchlist")
			}
			for _, w := range tc.want {
				if !strings.Contains(err.Error(), w) {
					t.Errorf("error does not say %q:\n%v", w, err)
				}
			}
			if n := feeds.requestCount(); n != 0 {
				t.Errorf("Run made %d requests before failing on the Watchlist", n)
			}
			if _, err := os.Stat(in.StateDir); !os.IsNotExist(err) {
				t.Errorf("Run wrote the state directory before failing on the Watchlist (%v)", err)
			}
		})
	}
}

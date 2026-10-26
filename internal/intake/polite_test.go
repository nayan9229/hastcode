package intake_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

// Polite fetching (ticket 03): conditional requests, an identifying
// User-Agent, a per-request timeout, and broken or silent Sources reported
// without stopping the run.

func TestRun_IdentifiesItselfAsHastCodeWithAContactAddress(t *testing.T) {
	feeds := newFeedServer(t)
	in := newIntake(t, exampleWatchlist(t, feeds))

	mustRun(t, in)

	for _, path := range []string{"/rss.xml", "/atom.xml"} {
		reqs := feeds.requestsTo(path)
		if len(reqs) != 1 {
			t.Fatalf("%s requested %d times, want once", path, len(reqs))
		}
		ua := reqs[0].UserAgent()
		if !strings.Contains(ua, "HastCode") || !strings.Contains(ua, "mira-vance@hastcode.com") {
			t.Errorf("%s User-Agent = %q, want one naming HastCode and mira-vance@hastcode.com", path, ua)
		}
	}
}

func TestRun_SendsTheStoredValidatorsBackOnTheNextRun(t *testing.T) {
	feeds := newFeedServer(t)
	feeds.handle("/rss.xml", conditionalFeed(t, "rss.xml", `W/"rss-v1"`, "Sat, 24 Oct 2026 14:30:00 GMT"))
	feeds.handle("/atom.xml", conditionalFeed(t, "atom.xml", "", "Sun, 25 Oct 2026 08:00:00 GMT"))
	in := newIntake(t, exampleWatchlist(t, feeds))

	mustRun(t, in)
	mustRun(t, in)

	for _, c := range []struct{ path, ifNoneMatch, ifModifiedSince string }{
		{"/rss.xml", `W/"rss-v1"`, "Sat, 24 Oct 2026 14:30:00 GMT"},
		{"/atom.xml", "", "Sun, 25 Oct 2026 08:00:00 GMT"},
	} {
		reqs := feeds.requestsTo(c.path)
		if len(reqs) != 2 {
			t.Fatalf("%s requested %d times, want twice", c.path, len(reqs))
		}
		if h := reqs[0].Header; h.Get("If-None-Match") != "" || h.Get("If-Modified-Since") != "" {
			t.Errorf("%s first request is conditional (%q, %q), want a plain GET", c.path, h.Get("If-None-Match"), h.Get("If-Modified-Since"))
		}
		h := reqs[1].Header
		if got := h.Get("If-None-Match"); got != c.ifNoneMatch {
			t.Errorf("%s second request If-None-Match = %q, want %q", c.path, got, c.ifNoneMatch)
		}
		if got := h.Get("If-Modified-Since"); got != c.ifModifiedSince {
			t.Errorf("%s second request If-Modified-Since = %q, want the stored Last-Modified %q verbatim", c.path, got, c.ifModifiedSince)
		}
	}
}

// conditionalFeed serves testdata/feeds/<name> with the validators etag
// and lastModified (either may be empty: not sent), and answers 304 Not
// Modified to a request that presents them (RFC 9110 §13.1: If-None-Match
// wins over If-Modified-Since when both are sent).
func conditionalFeed(t *testing.T, name, etag, lastModified string) http.HandlerFunc {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("testdata", "feeds", name))
	if err != nil {
		t.Fatal(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if etag != "" {
			w.Header().Set("ETag", etag)
		}
		if lastModified != "" {
			w.Header().Set("Last-Modified", lastModified)
		}
		inm, ims := r.Header.Get("If-None-Match"), r.Header.Get("If-Modified-Since")
		if (inm != "" && inm == etag) || (inm == "" && ims != "" && ims == lastModified) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Content-Type", "text/xml")
		w.Write(body)
	}
}

func TestRun_ANotModifiedSourceListsNothingAndKeepsItsLeads(t *testing.T) {
	feeds := newFeedServer(t)
	notModified := 0
	cond := conditionalFeed(t, "rss.xml", `"rss-v1"`, "")
	feeds.handle("/rss.xml", func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w}
		cond(rec, r)
		if rec.status == http.StatusNotModified {
			notModified++
		}
	})
	in := newIntake(t, exampleWatchlist(t, feeds))
	first := readQueueAfterRun(t, in)

	report := mustRun(t, in)

	if notModified != 1 {
		t.Fatalf("the publisher answered 304 %d times, want once (the second run)", notModified)
	}
	if len(report.New) != 0 {
		t.Errorf("a run answered 304 lists new Leads %+v, want none", report.New)
	}
	if len(report.Failed) != 0 {
		t.Errorf("a 304 is reported as a failure: %+v", report.Failed)
	}
	if got := readQueue(t, in); len(got) != len(first) {
		t.Errorf("after a 304 the queue holds %d Leads %q, want the first run's %d", len(got), titles(got), len(first))
	}
	// The validators survive the 304: a third run is conditional too.
	mustRun(t, in)
	reqs := feeds.requestsTo("/rss.xml")
	if got := reqs[len(reqs)-1].Header.Get("If-None-Match"); got != `"rss-v1"` {
		t.Errorf("the run after a 304 sent If-None-Match %q, want %q", got, `"rss-v1"`)
	}
}

func TestRun_ASourceWithoutValidatorsIsFetchedInFull(t *testing.T) {
	feeds := newFeedServer(t) // sends neither ETag nor Last-Modified, like the Go blog
	in := newIntake(t, exampleWatchlist(t, feeds))
	mustRun(t, in)
	feeds.setBody("/rss.xml", []byte(`<?xml version="1.0"?>
<rss version="2.0"><channel><title>Example Engineering</title>
  <item>
    <title>Cell-based architecture, one year in</title>
    <link>https://engineering.example.com/posts/cells</link>
    <pubDate>Mon, 26 Oct 2026 07:00:00 +0000</pubDate>
  </item>
</channel></rss>`))

	report := mustRun(t, in)

	reqs := feeds.requestsTo("/rss.xml")
	if h := reqs[len(reqs)-1].Header; h.Get("If-None-Match") != "" || h.Get("If-Modified-Since") != "" {
		t.Errorf("second request to a Source without validators is conditional (%q, %q), want a plain GET",
			h.Get("If-None-Match"), h.Get("If-Modified-Since"))
	}
	if len(report.New) != 1 || report.New[0].Title != "Cell-based architecture, one year in" {
		t.Errorf("new Leads = %+v, want the cell-based architecture post", report.New)
	}
}

func TestRun_ReportsEachFailingSourceByNameAndListsTheOthers(t *testing.T) {
	feeds := newFeedServer(t)
	feeds.handle("/down.xml", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "upstream unavailable", http.StatusServiceUnavailable)
	})
	feeds.setBody("/moved.html", []byte("<!doctype html><html><body>We moved our blog</body></html>"))
	in := newIntake(t, writeWatchlist(t, feeds, `sources:
  - name: Example Engineering
    site: https://engineering.example.com/
    feed: "{{server}}/rss.xml"
    vertical: distributed-systems
    kind: case-study
    reason: Postmortems.
  - name: Down Blog
    site: https://down.example.com/
    feed: "{{server}}/down.xml"
    vertical: distributed-systems
    kind: concept
    reason: Answers 503.
  - name: Gone Blog
    site: https://gone.example.com/
    feed: "{{server}}/gone.xml"
    vertical: distributed-systems
    kind: concept
    reason: Answers 404.
  - name: Moved Blog
    site: https://moved.example.com/
    feed: "{{server}}/moved.html"
    vertical: databases-storage
    kind: concept
    reason: Its feed URL now serves an HTML page with a 200.
  - name: Unreachable Blog
    site: https://unreachable.example.com/
    feed: "http://127.0.0.1:1/feed.xml"
    vertical: databases-storage
    kind: concept
    reason: Nothing listens there.
  - name: Example Database Blog
    site: https://db.example.org/blog/
    feed: "{{server}}/atom.xml"
    vertical: databases-storage
    kind: concept
    reason: Storage-engine internals.
`))

	report, err := in.Run(context.Background())

	if err != nil {
		t.Fatalf("Run = %v, want the run to finish despite failing Sources", err)
	}
	want := map[string]string{ // Source -> what its reason must mention
		"Down Blog":        "503",
		"Gone Blog":        "404",
		"Moved Blog":       "not an RSS 2.0 or Atom feed",
		"Unreachable Blog": "127.0.0.1:1",
	}
	var failed []string
	for _, f := range report.Failed {
		failed = append(failed, f.Source)
		if !strings.Contains(f.Reason, want[f.Source]) {
			t.Errorf("%s failed with %q, want a reason mentioning %q", f.Source, f.Reason, want[f.Source])
		}
	}
	if wantOrder := []string{"Down Blog", "Gone Blog", "Moved Blog", "Unreachable Blog"}; !slices.Equal(failed, wantOrder) {
		t.Errorf("failed Sources = %q, want %q (Watchlist order)", failed, wantOrder)
	}
	if got := titles(readQueue(t, in)); len(got) != 4 {
		t.Errorf("queue lists %q, want the four Leads of the two working Sources", got)
	}
	var summary strings.Builder
	report.WriteSummary(&summary)
	for name, reason := range want {
		if !strings.Contains(summary.String(), "failed "+name+": ") || !strings.Contains(summary.String(), reason) {
			t.Errorf("summary does not report %s failing (%s):\n%s", name, reason, summary.String())
		}
	}
}

func TestRun_AHangingSourceTimesOutAndTheRunGoesOn(t *testing.T) {
	feeds := newFeedServer(t)
	feeds.handle("/rss.xml", func(w http.ResponseWriter, r *http.Request) {
		select { // never answers; gives up when the client does
		case <-r.Context().Done():
		case <-time.After(10 * time.Second):
		}
	})
	in := newIntake(t, exampleWatchlist(t, feeds))
	in.Timeout = 200 * time.Millisecond

	start := time.Now()
	report := mustRun(t, in)

	if took := time.Since(start); took > 5*time.Second {
		t.Errorf("the run took %s, want the hanging request cut off after its 200ms timeout", took)
	}
	if len(report.Failed) != 1 || report.Failed[0].Source != "Example Engineering" {
		t.Fatalf("failed = %+v, want the hanging Example Engineering", report.Failed)
	}
	if got := readQueue(t, in); len(got) != 2 || got[0].Source != "Example Database Blog" {
		t.Errorf("queue lists %q, want the two Leads of Example Database Blog", titles(got))
	}
}

func TestRun_FlagsASourceWhoseNewestItemIsOver90DaysOld(t *testing.T) {
	feeds := newFeedServer(t)
	quiet := []byte(`<?xml version="1.0"?>
<rss version="2.0"><channel><title>Quiet Blog</title>
  <item><title>Our last post</title><link>https://quiet.example.com/last</link>
    <pubDate>Mon, 20 Jul 2026 10:00:00 +0000</pubDate></item>
  <item><title>The one before</title><link>https://quiet.example.com/before</link>
    <pubDate>Mon, 01 Jun 2026 10:00:00 +0000</pubDate></item>
</channel></rss>`)
	slow := []byte(`<?xml version="1.0"?>
<rss version="2.0"><channel><title>Slow Blog</title>
  <item><title>A post 77 days ago</title><link>https://slow.example.com/post</link>
    <pubDate>Mon, 10 Aug 2026 10:00:00 +0000</pubDate></item>
</channel></rss>`)
	feeds.handle("/quiet.xml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `"quiet-1"`)
		if r.Header.Get("If-None-Match") == `"quiet-1"` {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Write(quiet)
	})
	feeds.setBody("/slow.xml", slow)
	in := newIntake(t, writeWatchlist(t, feeds, `sources:
  - name: Quiet Blog
    site: https://quiet.example.com/
    feed: "{{server}}/quiet.xml"
    vertical: distributed-systems
    kind: concept
    reason: Last posted in July.
  - name: Slow Blog
    site: https://slow.example.com/
    feed: "{{server}}/slow.xml"
    vertical: distributed-systems
    kind: concept
    reason: Posts every couple of months.
`))

	for run := 1; run <= 2; run++ { // the second run gets a 304 from Quiet Blog
		report := mustRun(t, in)

		if len(report.Silent) != 1 || report.Silent[0].Source != "Quiet Blog" ||
			!report.Silent[0].Newest.Equal(time.Date(2026, 7, 20, 10, 0, 0, 0, time.UTC)) {
			t.Fatalf("run %d: silent = %+v, want Quiet Blog, newest item 2026-07-20", run, report.Silent)
		}
		var summary strings.Builder
		report.WriteSummary(&summary)
		if want := "silent Quiet Blog: newest item 2026-07-20 (97 days ago)"; !strings.Contains(summary.String(), want) {
			t.Errorf("run %d: summary does not contain %q:\n%s", run, want, summary.String())
		}
		if strings.Contains(summary.String(), "silent Slow Blog") {
			t.Errorf("run %d: summary flags Slow Blog (newest item 77 days old) silent:\n%s", run, summary.String())
		}
	}
}

// statusRecorder notes the status a handler wrote.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

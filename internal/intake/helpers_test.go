package intake_test

import (
	"bufio"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"hastcode.com/site/internal/intake"
)

// testNow is the fixed clock of every run unless a test moves it.
var testNow = time.Date(2026, 10, 26, 9, 0, 0, 0, time.UTC)

// testVerticals are the Verticals a test Intake knows, in display order (in
// the command they come from site.yaml).
var testVerticals = []intake.Vertical{
	{Slug: "distributed-systems", Name: "Distributed Systems"},
	{Slug: "databases-storage", Name: "Databases & Storage"},
	{Slug: "languages-runtimes", Name: "Languages & Runtimes"},
}

// feedServer is an in-process publisher: an httptest server answering each
// path with a fixture from testdata/feeds/ (or a body a test set), 404 for
// anything else. It records every request.
type feedServer struct {
	t   *testing.T
	srv *httptest.Server

	mu       sync.Mutex
	bodies   map[string][]byte // "/rss.xml" -> body
	handlers map[string]http.HandlerFunc
	requests []*http.Request
}

func newFeedServer(t *testing.T) *feedServer {
	t.Helper()
	f := &feedServer{t: t, bodies: map[string][]byte{}, handlers: map[string]http.HandlerFunc{}}
	for _, name := range []string{"rss.xml", "atom.xml"} {
		f.serveFile("/"+name, filepath.Join("testdata", "feeds", name))
	}
	f.srv = httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(f.srv.Close)
	return f
}

// serveFile answers path with the bytes of file.
func (f *feedServer) serveFile(path, file string) {
	f.t.Helper()
	data, err := os.ReadFile(file)
	if err != nil {
		f.t.Fatal(err)
	}
	f.setBody(path, data)
}

// setBody answers path with body from now on.
func (f *feedServer) setBody(path string, body []byte) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.bodies[path] = body
}

func (f *feedServer) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	f.requests = append(f.requests, r.Clone(context.Background()))
	body, ok := f.bodies[r.URL.Path]
	h := f.handlers[r.URL.Path]
	f.mu.Unlock()
	if h != nil {
		h(w, r)
		return
	}
	if !ok {
		http.NotFound(w, r)
		return
	}
	// Publishers label feeds inconsistently (research §5.1); the format
	// comes from the root element.
	w.Header().Set("Content-Type", "text/xml")
	w.Write(body)
}

// handle answers path with h from now on, in place of its body (a test's
// own publisher behaviour: validators, errors, hangs).
func (f *feedServer) handle(path string, h http.HandlerFunc) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.handlers[path] = h
}

// requestsTo are the requests made for path so far, in order.
func (f *feedServer) requestsTo(path string) []*http.Request {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []*http.Request
	for _, r := range f.requests {
		if r.URL.Path == path {
			out = append(out, r)
		}
	}
	return out
}

func (f *feedServer) requestCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.requests)
}

// writeWatchlist writes yaml (with {{server}} replaced by the feed server's
// URL) as a Watchlist file in a temp directory and returns its path.
func writeWatchlist(t *testing.T, f *feedServer, yaml string) string {
	t.Helper()
	file := filepath.Join(t.TempDir(), "sources.yaml")
	if err := os.WriteFile(file, []byte(strings.ReplaceAll(yaml, "{{server}}", f.srv.URL)), 0o644); err != nil {
		t.Fatal(err)
	}
	return file
}

// exampleWatchlist is testdata/sources.yaml pointed at the feed server.
func exampleWatchlist(t *testing.T, f *feedServer) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "sources.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return writeWatchlist(t, f, string(data))
}

// newIntake is an Intake over watchlist with a temp state directory, the
// fixed clock and the test Verticals. Its repository Sources go to a fake
// GitHub where every repository has no releases, with a test token
// (newIntakeWithGitHub gives a test the fake).
func newIntake(t *testing.T, watchlist string) *intake.Intake {
	t.Helper()
	return &intake.Intake{
		Watchlist: watchlist,
		Verticals: testVerticals,
		StateDir:  filepath.Join(t.TempDir(), ".intake"),
		Now:       func() time.Time { return testNow },
		GitHubURL: newFakeGitHub(t).srv.URL,
		Getenv:    tokenEnv,
	}
}

func mustRun(t *testing.T, in *intake.Intake) *intake.Report {
	t.Helper()
	r, err := in.Run(context.Background())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	return r
}

// queuedLead is one Lead line of the queue file, as the Author reads it:
//
//   - `id` [Title](link) · Source · kind · 2026-10-24 (1 day old) · new
type queuedLead struct {
	Vertical string // the heading the Lead sits under
	ID       string
	Title    string
	Link     string
	Source   string
	Kind     string
	Date     string // "2026-10-24 (1 day old)"
	New      bool
	Covered  bool // marked covered: a Piece already cites its link
}

// readQueue parses the queue file of in's state directory.
func readQueue(t *testing.T, in *intake.Intake) []queuedLead {
	t.Helper()
	file, err := os.Open(filepath.Join(in.StateDir, intake.QueueFile))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	var leads []queuedLead
	vertical := ""
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		if h, ok := strings.CutPrefix(line, "## "); ok {
			vertical = h
			continue
		}
		rest, ok := strings.CutPrefix(line, "- `")
		if !ok {
			continue
		}
		id, rest, _ := strings.Cut(rest, "` [")
		title, rest, _ := strings.Cut(rest, "](")
		link, rest, _ := strings.Cut(rest, ") · ")
		fields := strings.Split(rest, " · ")
		if len(fields) < 3 {
			t.Fatalf("queue line %q: want Source · kind · date", line)
		}
		leads = append(leads, queuedLead{
			Vertical: vertical, ID: id, Title: title, Link: link,
			Source: fields[0], Kind: fields[1], Date: fields[2],
			New:     len(fields) > 3 && fields[3] == "new",
			Covered: slices.Contains(fields[3:], "covered"),
		})
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	return leads
}

func titles(leads []queuedLead) []string {
	var out []string
	for _, l := range leads {
		out = append(out, l.Title)
	}
	return out
}

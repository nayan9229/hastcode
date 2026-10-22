package factpack_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"hastcode.com/site/internal/factpack"
	"hastcode.com/site/internal/github"
)

// fakeGitHub is an in-process GitHub REST API: an httptest server answering
// the paths the fetch command asks for with recorded JSON from
// testdata/github/<dir>/ (one file per endpoint, see serveRepo), 404 with
// GitHub's "Not Found" body for anything else, and GitHub's rate-limit
// headers on every answer. It records every request so a test can check the
// headers sent and the number of calls made.
type fakeGitHub struct {
	t   *testing.T
	srv *httptest.Server

	mu       sync.Mutex
	routes   map[string]*route // "/repos/apache/iggy/releases?per_page=10"
	requests []fakeRequest
	sleeps   []time.Duration // what the client asked to wait
}

// route serves its answers in order; the last one repeats.
type route struct {
	answers []answer
}

type answer struct {
	status int
	header http.Header
	body   []byte
}

// fakeRequest is one request the fake received: the path with its query,
// and the headers.
type fakeRequest struct {
	Path   string
	Header http.Header
}

// fixtureFiles maps a fixture file in testdata/github/<dir>/ to the path
// it answers, with {branch} replaced by the repository's default branch.
var fixtureFiles = map[string]string{
	"repo.json":               "",
	"languages.json":          "/languages",
	"releases.json":           "/releases?per_page=10",
	"contributors.json":       "/contributors?per_page=10",
	"contributors-count.json": "/contributors?per_page=1",
	"tree.json":               "/git/trees/{branch}",
	"github-dir.json":         "/contents/.github?ref={branch}",
	"readme.json":             "/readme",
	"stargazers-count.json":   "/stargazers/count",
}

const notFoundBody = `{"message":"Not Found","documentation_url":"https://docs.github.com/rest","status":"404"}`

func newFakeGitHub(t *testing.T) *fakeGitHub {
	t.Helper()
	g := &fakeGitHub{t: t, routes: map[string]*route{}}
	g.srv = httptest.NewServer(http.HandlerFunc(g.serve))
	t.Cleanup(g.srv.Close)
	return g
}

// serveRepo loads testdata/github/<dir>/: repo.json names the repository
// (full_name) and its default branch; every other fixture file present
// answers its endpoint (fixtureFiles). A <name>.headers file beside a
// fixture adds response headers (one "Name: value" per line). Endpoints
// without a fixture answer 404, which is how a repository without a README,
// a .github directory or the stargazers/count endpoint is recorded.
func (g *fakeGitHub) serveRepo(dir string) {
	g.t.Helper()
	root := filepath.Join("testdata", "github", dir)
	repoJSON := readFile(g.t, filepath.Join(root, "repo.json"))
	var repo struct {
		FullName      string `json:"full_name"`
		DefaultBranch string `json:"default_branch"`
	}
	if err := json.Unmarshal(repoJSON, &repo); err != nil {
		g.t.Fatalf("%s/repo.json: %v", root, err)
	}
	for file, suffix := range fixtureFiles {
		body, err := os.ReadFile(filepath.Join(root, file))
		if err != nil {
			continue
		}
		path := "/repos/" + repo.FullName + strings.ReplaceAll(suffix, "{branch}", repo.DefaultBranch)
		h := http.Header{}
		if extra, err := os.ReadFile(filepath.Join(root, strings.TrimSuffix(file, ".json")+".headers")); err == nil {
			sc := bufio.NewScanner(bytes.NewReader(extra))
			for sc.Scan() {
				if name, value, ok := strings.Cut(sc.Text(), ":"); ok {
					h.Set(strings.TrimSpace(name), strings.TrimSpace(value))
				}
			}
		}
		g.answer(path, answer{status: http.StatusOK, header: h, body: body})
	}
}

// answer sets path's only answer.
func (g *fakeGitHub) answer(path string, a answer) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.routes[path] = &route{answers: []answer{a}}
}

// answerSequence makes path answer each answer in turn, repeating the last.
func (g *fakeGitHub) answerSequence(path string, answers ...answer) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.routes[path] = &route{answers: answers}
}

// requestsTo returns the requests whose path (with query) has prefix.
func (g *fakeGitHub) requestsTo(prefix string) []fakeRequest {
	g.mu.Lock()
	defer g.mu.Unlock()
	var out []fakeRequest
	for _, r := range g.requests {
		if strings.HasPrefix(r.Path, prefix) {
			out = append(out, r)
		}
	}
	return out
}

// sleptFor returns every wait the client asked for, in order.
func (g *fakeGitHub) sleptFor() []time.Duration {
	g.mu.Lock()
	defer g.mu.Unlock()
	return append([]time.Duration(nil), g.sleeps...)
}

// client is a GitHub client pointed at the fake, holding the token "test-token"
// and recording waits instead of sleeping.
func (g *fakeGitHub) client() *github.Client {
	g.t.Helper()
	c, err := github.NewClient(func(key string) string {
		if key == github.EnvToken {
			return "test-token"
		}
		return ""
	})
	if err != nil {
		g.t.Fatal(err)
	}
	c.BaseURL = g.srv.URL
	c.Sleep = func(_ context.Context, d time.Duration) error {
		g.mu.Lock()
		defer g.mu.Unlock()
		g.sleeps = append(g.sleeps, d)
		return nil
	}
	return c
}

func (g *fakeGitHub) serve(w http.ResponseWriter, r *http.Request) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.requests = append(g.requests, fakeRequest{Path: r.URL.RequestURI(), Header: r.Header.Clone()})

	a := answer{status: http.StatusNotFound, body: []byte(notFoundBody)}
	if rt, ok := g.routes[r.URL.RequestURI()]; ok {
		a = rt.answers[0]
		if len(rt.answers) > 1 {
			rt.answers = rt.answers[1:]
		}
	}
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-RateLimit-Limit", "5000")
	h.Set("X-RateLimit-Remaining", "4990")
	h.Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10))
	h.Set("X-RateLimit-Resource", "core")
	for name, values := range a.header {
		h[name] = values
	}
	w.WriteHeader(a.status)
	_, _ = w.Write(a.body)
}

// harness is a Fetcher against a fake GitHub with its own content tree and
// cache directory, a fixed clock, and the log captured.
type harness struct {
	gh         *fakeGitHub
	fetcher    *factpack.Fetcher
	contentDir string
	cacheDir   string
	log        *bytes.Buffer
}

// fetchedAt is the harness clock: the golden Fact Pack's fetched_at.
var fetchedAt = time.Date(2026, 10, 19, 6, 0, 0, 0, time.UTC)

func newHarness(t *testing.T) *harness {
	t.Helper()
	root := t.TempDir()
	h := &harness{
		gh:         newFakeGitHub(t),
		contentDir: filepath.Join(root, "content"),
		cacheDir:   filepath.Join(root, "node_modules", ".cache", "hastcode"),
		log:        &bytes.Buffer{},
	}
	h.fetcher = &factpack.Fetcher{
		GitHub:   h.gh.client(),
		CacheDir: h.cacheDir,
		Now:      func() time.Time { return fetchedAt },
		Log:      h.log,
	}
	return h
}

// markdownPath is where a Profile's Markdown would live in the harness
// content tree (the directory is created; the file is not).
func (h *harness) markdownPath(vertical, slug string) string {
	h.gh.t.Helper()
	dir := filepath.Join(h.contentDir, "explainers", vertical)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		h.gh.t.Fatal(err)
	}
	return filepath.Join(dir, slug+".md")
}

// writeProfile writes a minimal Repo Profile Markdown for repo at vertical/slug
// and returns the Markdown path.
func (h *harness) writeProfile(vertical, slug, repo string) string {
	h.gh.t.Helper()
	md := h.markdownPath(vertical, slug)
	fm := "---\n" +
		"title: \"" + slug + "\"\n" +
		"description: \"A Repo Profile of " + repo + ".\"\n" +
		"kind: repo-profile\n" +
		"repo: \"" + repo + "\"\n" +
		"vertical: " + vertical + "\n" +
		"published: 2026-10-21T09:00:00Z\n" +
		"reviewed: 2026-10-21T09:00:00Z\n" +
		"author: ada\n" +
		"---\n\nBody.\n"
	if err := os.WriteFile(md, []byte(fm), 0o644); err != nil {
		h.gh.t.Fatal(err)
	}
	return md
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// firstDifference reports the first line where got differs from want.
func firstDifference(want, got []byte) string {
	w, g := strings.Split(string(want), "\n"), strings.Split(string(got), "\n")
	for i := 0; i < len(w) || i < len(g); i++ {
		var wl, gl string
		if i < len(w) {
			wl = w[i]
		}
		if i < len(g) {
			gl = g[i]
		}
		if wl != gl {
			return "line " + strconv.Itoa(i+1) + ":\n  want: " + wl + "\n  got:  " + gl
		}
	}
	return "identical"
}

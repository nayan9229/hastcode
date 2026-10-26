package intake_test

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"hastcode.com/site/internal/intake"
)

// testToken is the GitHub token a test Intake reads from its environment.
const testToken = "test-token"

// fakeGitHub is an in-process GitHub REST API (the Fact Pack fetcher's
// approach, internal/factpack/fake_test.go) answering only
// GET /repos/{owner}/{name}/releases: the JSON a test set for the
// repository, else [] (a repository with no releases), paginated by
// per_page and page with a Link rel="next" header. It sends a weak ETag
// derived from the body and answers 304 to a matching If-None-Match, as
// GitHub does, with GitHub's rate-limit headers on every answer. Anything
// else is 404. It records every request.
type fakeGitHub struct {
	t   *testing.T
	srv *httptest.Server

	mu       sync.Mutex
	releases map[string][]byte // "example/runtime" -> JSON body
	requests []githubRequest
}

// githubRequest is one request the fake received and how it answered.
type githubRequest struct {
	URI    string // path and query
	Header http.Header
	Status int
}

func newFakeGitHub(t *testing.T) *fakeGitHub {
	t.Helper()
	g := &fakeGitHub{t: t, releases: map[string][]byte{}}
	g.srv = httptest.NewServer(http.HandlerFunc(g.serve))
	t.Cleanup(g.srv.Close)
	return g
}

// serveReleases answers repo's releases with the JSON in file
// (testdata/github/…).
func (g *fakeGitHub) serveReleases(repo, file string) {
	g.t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "github", file))
	if err != nil {
		g.t.Fatal(err)
	}
	g.setReleases(repo, data)
}

// setReleases answers repo's releases with body from now on.
func (g *fakeGitHub) setReleases(repo string, body []byte) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.releases[repo] = body
}

func (g *fakeGitHub) serve(w http.ResponseWriter, r *http.Request) {
	g.mu.Lock()
	defer g.mu.Unlock()
	status := http.StatusNotFound
	defer func() {
		g.requests = append(g.requests, githubRequest{URI: r.URL.RequestURI(), Header: r.Header.Clone(), Status: status})
	}()
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-RateLimit-Limit", "5000")
	h.Set("X-RateLimit-Remaining", "4990")
	h.Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Hour).Unix(), 10))
	h.Set("X-RateLimit-Resource", "core")

	repo, ok := strings.CutPrefix(r.URL.Path, "/repos/")
	repo, ok2 := strings.CutSuffix(repo, "/releases")
	if !ok || !ok2 || strings.Count(repo, "/") != 1 {
		w.WriteHeader(status)
		fmt.Fprint(w, `{"message":"Not Found","status":"404"}`)
		return
	}
	all, ok := g.releases[repo]
	if !ok {
		all = []byte("[]")
	}
	body, next := g.page(all, r.URL.Query())
	if next != 0 {
		q := r.URL.Query()
		q.Set("page", strconv.Itoa(next))
		h.Set("Link", fmt.Sprintf(`<%s%s?%s>; rel="next"`, g.srv.URL, r.URL.Path, q.Encode()))
	}
	etag := fmt.Sprintf(`W/"%x"`, sha256.Sum256(body))
	h.Set("ETag", etag)
	h.Set("Cache-Control", "private, max-age=60, s-maxage=60")
	if r.Header.Get("If-None-Match") == etag {
		status = http.StatusNotModified
		w.WriteHeader(status)
		return
	}
	status = http.StatusOK
	w.WriteHeader(status)
	w.Write(body)
}

// page is the page of the releases JSON array all that query asks for
// (per_page, default 30; page, default 1), as GitHub paginates, and the
// next page's number (0 on the last page).
func (g *fakeGitHub) page(all []byte, query url.Values) (body []byte, next int) {
	var releases []json.RawMessage
	if err := json.Unmarshal(all, &releases); err != nil {
		g.t.Errorf("fake GitHub: releases are not a JSON array: %v", err)
		return all, 0
	}
	perPage, page := 30, 1
	if n, err := strconv.Atoi(query.Get("per_page")); err == nil && n > 0 {
		perPage = n
	}
	if n, err := strconv.Atoi(query.Get("page")); err == nil && n > 0 {
		page = n
	}
	start := min((page-1)*perPage, len(releases))
	end := min(start+perPage, len(releases))
	if end < len(releases) {
		next = page + 1
	}
	body, err := json.Marshal(releases[start:end])
	if err != nil {
		g.t.Errorf("fake GitHub: %v", err)
	}
	return body, next
}

// requestLog returns every request received so far.
func (g *fakeGitHub) requestLog() []githubRequest {
	g.mu.Lock()
	defer g.mu.Unlock()
	return append([]githubRequest(nil), g.requests...)
}

// tokenEnv is an environment holding only the test token.
func tokenEnv(key string) string {
	if key == "GITHUB_TOKEN" {
		return testToken
	}
	return ""
}

// newIntakeWithGitHub is newIntake with a fake GitHub of its own, returned
// so the test can set releases and read the requests.
func newIntakeWithGitHub(t *testing.T, watchlist string) (*intake.Intake, *fakeGitHub) {
	t.Helper()
	gh := newFakeGitHub(t)
	in := newIntake(t, watchlist)
	in.GitHubURL = gh.srv.URL
	return in, gh
}

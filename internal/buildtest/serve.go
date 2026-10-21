package buildtest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"hastcode.com/site/internal/vercel"
	"hastcode.com/site/subscribe"
)

// Server serves a build output the way the local run command does: the
// Build Output config.json routes in front of static/ and a stand-in for
// the subscribe function (which answers 204 and records each request).
type Server struct {
	URL string
	// FunctionCalls counts requests that reached the subscribe function.
	FunctionCalls int
	client        *http.Client
}

// Response is a fetched response with its body read.
type Response struct {
	*http.Response
	Body string
}

// Serve starts a local server over outDir (config.json + static/). It fails
// the test if config.json has routes the local router cannot apply.
func Serve(tb testing.TB, outDir string) *Server {
	tb.Helper()
	s := &Server{}
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.FunctionCalls++
		w.WriteHeader(http.StatusNoContent)
	})
	h, err := vercel.NewRouter(outDir, map[string]http.Handler{subscribe.Path: fn})
	if err != nil {
		tb.Fatalf("local router over %s: %v", outDir, err)
	}
	ts := httptest.NewServer(h)
	tb.Cleanup(ts.Close)
	s.URL = ts.URL
	s.client = &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	return s
}

// URLPaths lists the site path of every file under static/: a directory
// index ("a/index.html") is "/a/", any other file its own path. Sorted.
func URLPaths(tb testing.TB, outDir string) []string {
	tb.Helper()
	root := filepath.Join(outDir, "static")
	var paths []string
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		u := "/" + filepath.ToSlash(rel)
		if dir, ok := strings.CutSuffix(u, "/index.html"); ok {
			u = dir + "/"
		}
		paths = append(paths, u)
		return nil
	})
	if err != nil {
		tb.Fatalf("walk %s: %v", root, err)
	}
	sort.Strings(paths)
	return paths
}

// Get requests path without following redirects.
func (s *Server) Get(tb testing.TB, path string) Response {
	tb.Helper()
	return s.Do(tb, http.MethodGet, path, "")
}

// Do sends a request (form-encoded when body is non-empty) without following
// redirects.
func (s *Server) Do(tb testing.TB, method, path, body string) Response {
	tb.Helper()
	req, err := http.NewRequest(method, s.URL+path, strings.NewReader(body))
	if err != nil {
		tb.Fatalf("%s %s: %v", method, path, err)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	res, err := s.client.Do(req)
	if err != nil {
		tb.Fatalf("%s %s: %v", method, path, err)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		tb.Fatalf("%s %s: read body: %v", method, path, err)
	}
	return Response{Response: res, Body: string(data)}
}

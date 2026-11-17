package crawlcheck_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"hastcode.com/site/internal/crawlcheck"
)

func TestUnreachableHostFailsInPlainWords(t *testing.T) {
	dead := httptest.NewServer(http.NotFoundHandler())
	addr := dead.URL
	dead.Close()
	r, out := run(t, crawlcheck.Checker{Site: addr, HTTP: addr, WWW: addr, Timeout: time.Second})
	requireFailed(t, r, out, "robots.txt", "could not reach "+addr+"/robots.txt")
	requireFailed(t, r, out, "Googlebot", "could not reach")
	if !strings.Contains(out, "checks failed") {
		t.Errorf("no failure summary; output:\n%s", out)
	}
}

func TestTimeoutFails(t *testing.T) {
	release := make(chan struct{})
	slow := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-release:
		case <-r.Context().Done():
		}
	}))
	t.Cleanup(func() { close(release); slow.Close() })
	r, out := run(t, crawlcheck.Checker{Site: slow.URL, HTTP: slow.URL, WWW: slow.URL, Timeout: 50 * time.Millisecond})
	requireFailed(t, r, out, "robots.txt", "no answer within 50ms")
}

// recordingTransport records every request the Checker sends.
type recordingTransport struct {
	mu       sync.Mutex
	requests []*http.Request
}

func (rt *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.mu.Lock()
	rt.requests = append(rt.requests, req)
	rt.mu.Unlock()
	return http.DefaultTransport.RoundTrip(req)
}

func TestOnlyReadsFromTheStandIn(t *testing.T) {
	s := newStandIn(t)
	rt := &recordingTransport{}
	c := s.checker()
	c.Client = &http.Client{Transport: rt}
	if r, out := run(t, c); r.ExitCode() != 0 {
		t.Fatalf("exit code %d, want 0; output:\n%s", r.ExitCode(), out)
	}
	hosts := map[string]bool{}
	for _, srv := range []*httptest.Server{s.site, s.http, s.www} {
		u, _ := url.Parse(srv.URL)
		hosts[u.Host] = true
	}
	if len(rt.requests) == 0 {
		t.Fatal("no request went through the client")
	}
	for _, req := range rt.requests {
		if !hosts[req.URL.Host] {
			t.Errorf("request to %s, which is not the stand-in", req.URL)
		}
		if req.Method != http.MethodGet && req.Method != http.MethodHead {
			t.Errorf("%s %s: only GET or HEAD may be sent", req.Method, req.URL)
		}
	}
}

func TestRedirectsAreNotFollowed(t *testing.T) {
	// The http and www stand-ins redirect to the site's /; following them
	// would ask the site for / more than once per crawler.
	s := newStandIn(t)
	rt := &recordingTransport{}
	c := s.checker()
	c.Client = &http.Client{Transport: rt}
	run(t, c)
	home := 0
	for _, req := range rt.requests {
		if "http://"+req.URL.Host == s.site.URL && req.URL.Path == "/" {
			home++
		}
	}
	if want := len(crawlcheck.Crawlers); home != want {
		t.Errorf("%d requests for the site's /, want %d (one per crawler)", home, want)
	}
}

func TestNonCrawlerRequestsSayWhoIsAsking(t *testing.T) {
	// robots.txt and the redirect checks are not sent as a crawler: they
	// name the check, so Cloudflare's logs show what it was.
	s := newStandIn(t)
	rt := &recordingTransport{}
	c := s.checker()
	c.Client = &http.Client{Transport: rt}
	run(t, c)
	for _, req := range rt.requests {
		if req.URL.Path == "/robots.txt" || "http://"+req.URL.Host != s.site.URL {
			if req.UserAgent() != crawlcheck.UserAgent {
				t.Errorf("%s sent User-Agent %q, want %q", req.URL, req.UserAgent(), crawlcheck.UserAgent)
			}
		}
	}
}

func TestNoWorkflowRunsIt(t *testing.T) {
	// Local only: from CI's IP addresses it would test a different path
	// through Cloudflare's edge (README "Crawl check").
	files, err := filepath.Glob("../../.github/workflows/*.yml")
	if err != nil || len(files) == 0 {
		t.Fatalf("no workflows found (%v)", err)
	}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(data), "crawl-check") || strings.Contains(string(data), "crawlcheck") {
			t.Errorf("%s runs the crawl check; it is local only", f)
		}
	}
}

func TestNewDerivesHTTPAndWWWFromSiteURL(t *testing.T) {
	c, err := crawlcheck.New("https://hastcode.com")
	if err != nil {
		t.Fatal(err)
	}
	if c.Site != "https://hastcode.com" || c.HTTP != "http://hastcode.com" || c.WWW != "https://www.hastcode.com" {
		t.Errorf("New = %+v", c)
	}
	if _, err := crawlcheck.New("hastcode.com"); err == nil {
		t.Error("a site_url without https:// was accepted")
	}
}

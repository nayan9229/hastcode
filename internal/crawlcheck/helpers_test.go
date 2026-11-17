package crawlcheck_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"hastcode.com/site/internal/crawlcheck"
)

// cloudflareBlock is a shortened copy of the managed block Cloudflare
// prepends to the live /robots.txt (docs/research/indexing-submission-channels.md
// §0): a combined * group with the Content-Signal line, then Disallow
// groups for training crawlers and two search engines we do not check.
// Applebot-Extended sits here on purpose: it is not Applebot.
const cloudflareBlock = `# As a condition of accessing this website, you agree to abide by the
# following content signals:
User-Agent: *
Content-Signal: search=yes,ai-train=no,use=reference
Allow: /

User-agent: Amazonbot
Disallow: /

User-agent: Applebot-Extended
Disallow: /

User-agent: ClaudeBot
Disallow: /

User-agent: GPTBot
Disallow: /

User-agent: Baiduspider
Disallow: /

# END Cloudflare Managed Content

`

// ourRobots is what builder/step_write_robots.go writes, for a site at
// origin.
func ourRobots(origin string) string {
	return "User-agent: *\nAllow: /\n\n" +
		"User-agent: Googlebot\nAllow: /\n\n" +
		"User-agent: Bingbot\nAllow: /\n\n" +
		"User-agent: OAI-SearchBot\nAllow: /\n\n" +
		"User-agent: PerplexityBot\nAllow: /\n\n" +
		"Sitemap: " + origin + "/sitemap.xml\n" +
		"Sitemap: " + origin + "/news-sitemap.xml\n"
}

// standIn is a local Production: three httptest servers in place of
// https://hastcode.com, http://hastcode.com and https://www.hastcode.com.
// Out of the box every check passes: the site serves the live-shaped
// robots.txt and 200 for every crawler, http and www 301 to the site.
type standIn struct {
	site, http, www *httptest.Server

	mu sync.Mutex
	// robots is the /robots.txt body; "" means the live shape.
	robots string
	// blocked maps a user-agent substring to the status the site answers
	// it with (Cloudflare's edge 403).
	blocked map[string]int
	// httpStatus and httpLocation are what http answers; zero values
	// mean 301 to the site.
	httpStatus   int
	httpLocation string
}

func newStandIn(t *testing.T) *standIn {
	t.Helper()
	s := &standIn{blocked: map[string]int{}}
	s.site = httptest.NewServer(http.HandlerFunc(s.serveSite))
	s.http = httptest.NewServer(http.HandlerFunc(s.serveHTTP))
	s.www = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, s.site.URL+r.URL.Path, http.StatusMovedPermanently)
	}))
	t.Cleanup(func() { s.site.Close(); s.http.Close(); s.www.Close() })
	return s
}

func (s *standIn) serveSite(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for ua, status := range s.blocked {
		if strings.Contains(r.UserAgent(), ua) {
			w.WriteHeader(status)
			_, _ = w.Write([]byte("Your request was blocked."))
			return
		}
	}
	switch r.URL.Path {
	case "/robots.txt":
		body := s.robots
		if body == "" {
			body = cloudflareBlock + ourRobots(s.site.URL)
		}
		_, _ = w.Write([]byte(body))
	case "/", "/sitemap.xml":
		_, _ = w.Write([]byte("ok"))
	default:
		http.NotFound(w, r)
	}
}

func (s *standIn) serveHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	status, location := s.httpStatus, s.httpLocation
	s.mu.Unlock()
	if status == 0 {
		status, location = http.StatusMovedPermanently, s.site.URL+"/"
	}
	if location != "" {
		w.Header().Set("Location", location)
	}
	w.WriteHeader(status)
}

func (s *standIn) block(uaSubstring string, status int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blocked[uaSubstring] = status
}

func (s *standIn) setRobots(body string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.robots = body
}

func (s *standIn) setHTTP(status int, location string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.httpStatus, s.httpLocation = status, location
}

// checker points a Checker at the stand-in through its three overrides.
func (s *standIn) checker() crawlcheck.Checker {
	return crawlcheck.Checker{
		Site:    s.site.URL,
		HTTP:    s.http.URL,
		WWW:     s.www.URL,
		Timeout: 2 * time.Second,
	}
}

// run runs c and returns its report and printed output.
func run(t *testing.T, c crawlcheck.Checker) (crawlcheck.Report, string) {
	t.Helper()
	var out bytes.Buffer
	c.Out = &out
	r := c.Run(context.Background())
	return r, out.String()
}

// requireFailed asserts the report failed and some FAIL line contains
// every one of want.
func requireFailed(t *testing.T, r crawlcheck.Report, out string, want ...string) {
	t.Helper()
	if r.ExitCode() == 0 {
		t.Fatalf("exit code 0, want non-zero; output:\n%s", out)
	}
	for _, line := range strings.Split(out, "\n") {
		if !strings.HasPrefix(line, "FAIL") {
			continue
		}
		ok := true
		for _, w := range want {
			ok = ok && strings.Contains(line, w)
		}
		if ok {
			return
		}
	}
	t.Fatalf("no FAIL line contains all of %q; output:\n%s", want, out)
}

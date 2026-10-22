package factpack_test

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"hastcode.com/site/internal/github"
)

const releasesPath = "/repos/apache/iggy/releases?per_page=10"

func TestFetch_StopsOnARateLimitedForbiddenWithTheDocumentedMessage(t *testing.T) {
	h := newHarness(t)
	h.gh.serveRepo("apache-iggy")
	reset := time.Now().Add(47 * time.Minute).Truncate(time.Second)
	h.gh.answer(releasesPath, answer{
		status: http.StatusForbidden,
		header: http.Header{
			"X-Ratelimit-Limit":     {"1000"},
			"X-Ratelimit-Remaining": {"0"},
			"X-Ratelimit-Reset":     {strconv.FormatInt(reset.Unix(), 10)},
			"X-Ratelimit-Resource":  {"core"},
		},
		body: []byte(`{"message":"API rate limit exceeded for installation ID 1.","documentation_url":"https://docs.github.com/rest/overview/rate-limits-for-the-rest-api"}`),
	})
	md := h.markdownPath("distributed-systems", "apache-iggy")

	_, err := h.fetcher.Fetch(context.Background(), "apache/iggy", md)

	var limit *github.RateLimitError
	if !errors.As(err, &limit) {
		t.Fatalf("expected a RateLimitError, got %v", err)
	}
	msg := err.Error()
	for _, want := range []string{"GitHub rate limit reached", "0 of 1000", reset.UTC().Format(time.RFC3339), "never add a second token"} {
		if !strings.Contains(msg, want) {
			t.Errorf("message lacks %q:\n%s", want, msg)
		}
	}
	if n := len(h.gh.requestsTo(releasesPath)); n != 1 {
		t.Errorf("releases requested %d times; a limit that resets in 47m must not be retried", n)
	}
	if slept := h.gh.sleptFor(); len(slept) != 0 {
		t.Errorf("the client waited %v instead of stopping", slept)
	}
	if _, statErr := os.Stat(filepath.Join(filepath.Dir(md), "apache-iggy.json")); statErr == nil {
		t.Error("a Fact Pack was written from a half-fetched repository")
	}
}

func TestFetch_HonoursRetryAfterWithinTheBound(t *testing.T) {
	h := newHarness(t)
	h.gh.serveRepo("apache-iggy")
	ok := h.gh.routes[releasesPath].answers[0]
	h.gh.answerSequence(releasesPath,
		answer{
			status: http.StatusForbidden,
			header: http.Header{"Retry-After": {"30"}, "X-Ratelimit-Remaining": {"900"}},
			body:   []byte(`{"message":"You have exceeded a secondary rate limit. Please wait a few minutes before you try again."}`),
		},
		ok,
	)
	md := h.markdownPath("distributed-systems", "apache-iggy")

	if _, err := h.fetcher.Fetch(context.Background(), "apache/iggy", md); err != nil {
		t.Fatalf("a 30 s Retry-After is within the bound and should have been waited out: %v", err)
	}

	if slept := h.gh.sleptFor(); len(slept) != 1 || slept[0] != 30*time.Second {
		t.Errorf("waits = %v, want [30s]", slept)
	}
	if n := len(h.gh.requestsTo(releasesPath)); n != 2 {
		t.Errorf("releases requested %d times, want 2 (one retry)", n)
	}
}

func TestFetch_RetriesA202ThenSucceeds(t *testing.T) {
	h := newHarness(t)
	h.gh.serveRepo("apache-iggy")
	const contributors = "/repos/apache/iggy/contributors?per_page=10"
	ok := h.gh.routes[contributors].answers[0]
	computing := answer{status: http.StatusAccepted, body: []byte(`{}`)}
	h.gh.answerSequence(contributors, computing, computing, ok)
	md := h.markdownPath("distributed-systems", "apache-iggy")

	if _, err := h.fetcher.Fetch(context.Background(), "apache/iggy", md); err != nil {
		t.Fatalf("two 202s then 200 should succeed: %v", err)
	}

	fp := parseFactPack(t, filepath.Join(filepath.Dir(md), "apache-iggy.json"))
	if len(fp.Contributors.Top) != 4 {
		t.Errorf("top contributors = %d, want the 4 of the final answer", len(fp.Contributors.Top))
	}
	if slept := h.gh.sleptFor(); len(slept) != 2 {
		t.Errorf("waits = %v, want two (one per 202)", slept)
	}
}

func TestFetch_GivesUpOnAnEndpointThatKeepsAnswering202(t *testing.T) {
	h := newHarness(t)
	h.gh.serveRepo("apache-iggy")
	h.gh.answer(releasesPath, answer{status: http.StatusAccepted, body: []byte(`{}`)})
	md := h.markdownPath("distributed-systems", "apache-iggy")

	_, err := h.fetcher.Fetch(context.Background(), "apache/iggy", md)

	if err == nil || !strings.Contains(err.Error(), "202") {
		t.Fatalf("expected an error naming the 202, got %v", err)
	}
	if n := len(h.gh.requestsTo(releasesPath)); n != github.DefaultRetries+1 {
		t.Errorf("releases requested %d times, want %d (bounded retries)", n, github.DefaultRetries+1)
	}
}

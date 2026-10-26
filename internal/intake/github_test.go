package intake_test

import (
	"net/http"
	"strings"
	"testing"
)

// languagesLeads are the queue's Leads under Languages & Runtimes, where
// the example repository Source (Example Runtime, example/runtime, brief)
// files its releases.
func languagesLeads(t *testing.T, q []queuedLead) []queuedLead {
	t.Helper()
	var out []queuedLead
	for _, l := range q {
		if l.Vertical == "Languages & Runtimes" {
			out = append(out, l)
		}
	}
	return out
}

func TestRun_ListsNewReleasesOfARepositorySourceAsLeads(t *testing.T) {
	feeds := newFeedServer(t)
	in, gh := newIntakeWithGitHub(t, exampleWatchlist(t, feeds))
	gh.serveReleases("example/runtime", "releases.json")

	report := mustRun(t, in)

	want := []queuedLead{
		{Vertical: "Languages & Runtimes", Title: "Runtime 2.1.0",
			Link: "https://github.com/example/runtime/releases/tag/v2.1.0", Source: "Example Runtime",
			Kind: "brief", Date: "2026-10-25 (today)", New: true},
		{Vertical: "Languages & Runtimes", Title: "v2.0.3", // no release name: the tag
			Link: "https://github.com/example/runtime/releases/tag/v2.0.3", Source: "Example Runtime",
			Kind: "brief", Date: "2026-10-14 (11 days old)", New: true},
	}
	got := languagesLeads(t, readQueue(t, in))
	if len(got) != len(want) {
		t.Fatalf("queue lists %d releases %q, want %d", len(got), titles(got), len(want))
	}
	for i := range want {
		want[i].ID = got[i].ID
		if got[i] != want[i] {
			t.Errorf("release Lead %d:\n got %+v\nwant %+v", i, got[i], want[i])
		}
	}
	if len(report.New) != 6 {
		t.Errorf("report lists %d new Leads, want 6 (four feed items, two releases)", len(report.New))
	}
	if len(report.Skipped) != 0 {
		t.Errorf("skipped = %+v, want none: repository Sources are fetched", report.Skipped)
	}
}

func TestRun_SecondRunAsksGitHubConditionallyAndListsNoReleaseAgain(t *testing.T) {
	feeds := newFeedServer(t)
	in, gh := newIntakeWithGitHub(t, exampleWatchlist(t, feeds))
	gh.serveReleases("example/runtime", "releases.json")
	mustRun(t, in)

	report := mustRun(t, in)

	if len(report.New) != 0 {
		t.Errorf("second run lists %d new Leads %+v, want none", len(report.New), report.New)
	}
	if got := languagesLeads(t, readQueue(t, in)); len(got) != 2 || got[0].New || got[1].New {
		t.Errorf("second queue's releases = %+v, want the first run's two, not new", got)
	}
	reqs := gh.requestLog()
	if len(reqs) != 2 {
		t.Fatalf("GitHub got %d requests %+v, want one per run", len(reqs), reqs)
	}
	first, second := reqs[0], reqs[1]
	if first.Header.Get("If-None-Match") != "" {
		t.Errorf("first request sent If-None-Match %q, want none: nothing stored yet", first.Header.Get("If-None-Match"))
	}
	if second.Header.Get("If-None-Match") == "" || second.Status != http.StatusNotModified {
		t.Errorf("second request: If-None-Match %q, answered %d; want the stored ETag, answered 304", second.Header.Get("If-None-Match"), second.Status)
	}
	// An authenticated 304 does not count against GitHub's rate limit.
	if second.Header.Get("Authorization") != "Bearer "+testToken {
		t.Errorf("second request's Authorization = %q, want the token", second.Header.Get("Authorization"))
	}
}

func TestRun_WithoutAGitHubTokenFailsRepositorySourcesOnly(t *testing.T) {
	feeds := newFeedServer(t)
	noToken := func(string) string { return "" }

	feedOnly, gh := newIntakeWithGitHub(t, writeWatchlist(t, feeds, `sources:
  - name: Example Engineering
    site: https://engineering.example.com/
    feed: "{{server}}/rss.xml"
    vertical: distributed-systems
    kind: case-study
    reason: Postmortems with timelines.
`))
	feedOnly.Getenv = noToken
	if report := mustRun(t, feedOnly); len(report.New) != 2 {
		t.Errorf("a feed-only Watchlist without a token lists %d new Leads, want its 2", len(report.New))
	}

	withRepo := newIntake(t, exampleWatchlist(t, feeds))
	withRepo.GitHubURL = gh.srv.URL
	withRepo.Getenv = noToken
	report := mustRun(t, withRepo)

	if len(report.Failed) != 1 || report.Failed[0].Source != "Example Runtime" ||
		!strings.Contains(report.Failed[0].Reason, "GITHUB_TOKEN is not set") {
		t.Fatalf("failed = %+v, want only Example Runtime, naming GITHUB_TOKEN", report.Failed)
	}
	if len(report.New) != 4 {
		t.Errorf("the feeds beside it list %d new Leads, want their 4", len(report.New))
	}
	if reqs := gh.requestLog(); len(reqs) != 0 {
		t.Errorf("GitHub got %d requests without a token, want none", len(reqs))
	}
}

// The prerelease rule (fetch_github.go): prereleases and drafts are not
// Leads; a prerelease promoted to a full release is listed when promoted.
func TestRun_ListsAPrereleaseOnlyOncePromotedToAFullRelease(t *testing.T) {
	feeds := newFeedServer(t)
	in, gh := newIntakeWithGitHub(t, exampleWatchlist(t, feeds))
	gh.serveReleases("example/runtime", "releases.json")
	for _, l := range languagesLeads(t, readQueueAfterRun(t, in)) {
		if l.Title == "Runtime 2.1.0 RC 1" || l.Title == "Runtime 2.2.0" {
			t.Errorf("queue lists %q, a prerelease or draft", l.Title)
		}
	}
	gh.setReleases("example/runtime", []byte(`[
  {"id": 300, "tag_name": "v2.1.0-rc.1", "name": "Runtime 2.1.0 RC 1", "html_url": "https://github.com/example/runtime/releases/tag/v2.1.0-rc.1",
   "draft": false, "prerelease": false, "published_at": "2026-10-20T10:00:00Z"}
]`))

	report := mustRun(t, in)

	if len(report.New) != 1 || report.New[0].Title != "Runtime 2.1.0 RC 1" {
		t.Errorf("run after the promotion lists %+v, want the promoted release", report.New)
	}
}

func TestRun_NeverRequestsGitHubWebPages(t *testing.T) {
	feeds := newFeedServer(t)
	in, gh := newIntakeWithGitHub(t, exampleWatchlist(t, feeds))
	gh.serveReleases("example/runtime", "releases.json")
	var hosts []string
	in.HTTP = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		hosts = append(hosts, r.URL.Host)
		return http.DefaultTransport.RoundTrip(r)
	})}

	mustRun(t, in)
	mustRun(t, in)

	for _, h := range hosts {
		if h != strings.TrimPrefix(feeds.srv.URL, "http://") && h != strings.TrimPrefix(gh.srv.URL, "http://") {
			t.Errorf("Intake requested %s; want only the feed server and the GitHub REST API", h)
		}
	}
	for _, r := range gh.requestLog() {
		if r.URI != "/repos/example/runtime/releases?per_page=10" {
			t.Errorf("GitHub got GET %s; want only the REST releases endpoint", r.URI)
		}
	}
	if len(hosts) != 6 {
		t.Errorf("Intake made %d requests, want 6 (three Sources, two runs)", len(hosts))
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestRun_ListsOnlyTheReleasePublishedSinceTheLastRun(t *testing.T) {
	feeds := newFeedServer(t)
	in, gh := newIntakeWithGitHub(t, exampleWatchlist(t, feeds))
	gh.serveReleases("example/runtime", "releases.json")
	mustRun(t, in)
	gh.setReleases("example/runtime", []byte(`[
  {"id": 303, "tag_name": "v2.1.1", "name": "Runtime 2.1.1", "html_url": "https://github.com/example/runtime/releases/tag/v2.1.1",
   "draft": false, "prerelease": false, "published_at": "2026-10-26T08:00:00Z"},
  {"id": 301, "tag_name": "v2.1.0", "name": "Runtime 2.1.0", "html_url": "https://github.com/example/runtime/releases/tag/v2.1.0",
   "draft": false, "prerelease": false, "published_at": "2026-10-25T16:00:00Z"}
]`))

	report := mustRun(t, in)

	if len(report.New) != 1 || report.New[0].Title != "Runtime 2.1.1" {
		t.Fatalf("second run's new Leads = %+v, want only Runtime 2.1.1", report.New)
	}
	got := languagesLeads(t, readQueue(t, in))
	if len(got) != 3 || got[0].Title != "Runtime 2.1.1" || !got[0].New || got[0].Date != "2026-10-26 (today)" {
		t.Errorf("queue's releases = %+v, want Runtime 2.1.1 first, new, today, then the first run's two", got)
	}
}

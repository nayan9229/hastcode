package intake_test

import (
	"fmt"
	"net/http"
	"slices"
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
	if len(report.Failed) != 0 {
		t.Errorf("failed = %+v, want none: repository Sources are fetched", report.Failed)
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

// A busy release-candidate stream can fill the newest page of releases:
// the run reads further pages until one's newest non-draft release is
// older than the age cutoff, so the full release behind them is listed.
func TestRun_PagesPastPrereleasesToTheFullReleaseBehindThem(t *testing.T) {
	feeds := newFeedServer(t)
	in, gh := newIntakeWithGitHub(t, exampleWatchlist(t, feeds))
	var releases []string
	add := func(id int, name string, prerelease bool, published string) {
		releases = append(releases, fmt.Sprintf(`{"id": %d, "tag_name": %q, "name": %q,
  "html_url": "https://github.com/example/runtime/releases/tag/%s", "draft": false,
  "prerelease": %t, "published_at": %q}`, id, name, name, name, prerelease, published))
	}
	for i := range 10 { // page 1: ten release candidates this week
		add(500-i, fmt.Sprintf("v3.1.0-rc.%d", 10-i), true, fmt.Sprintf("2026-10-%02dT10:00:00Z", 25-i/2))
	}
	add(480, "v3.0.0", false, "2026-10-19T10:00:00Z") // page 2: the full release, then older ones
	for i := range 19 { // the rest of page 2 and all of page 3: last spring's
		add(470-i, fmt.Sprintf("v2.%d.0", 19-i), false, fmt.Sprintf("2026-04-%02dT10:00:00Z", 28-i))
	}
	for i := range 10 { // page 4: never needed
		add(400-i, fmt.Sprintf("v1.%d.0", 10-i), false, "2025-06-01T10:00:00Z")
	}
	gh.setReleases("example/runtime", []byte("["+strings.Join(releases, ",\n")+"]"))

	mustRun(t, in)

	if got := titles(languagesLeads(t, readQueue(t, in))); !slices.Equal(got, []string{"v3.0.0"}) {
		t.Errorf("queue's releases = %q, want only v3.0.0 (candidates are not Leads, last spring's are past the cutoff)", got)
	}
	var uris []string
	for _, r := range gh.requestLog() {
		uris = append(uris, r.URI)
	}
	want := []string{
		"/repos/example/runtime/releases?per_page=10",
		"/repos/example/runtime/releases?per_page=10&page=2",
		"/repos/example/runtime/releases?per_page=10&page=3", // page 2's newest full release is recent
	}
	if !slices.Equal(uris, want) {
		t.Errorf("GitHub requests = %q, want %q (page 3's newest release is past the cutoff: stop)", uris, want)
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

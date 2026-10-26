package intake

import (
	"cmp"
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"hastcode.com/site/internal/github"
)

func init() { registerFetcher(RepositorySource, githubFetcher{}) }

// githubFetcher fetches a repository Source: its releases, through the
// GitHub REST API with the Fact Pack fetcher's client (internal/github:
// the one GITHUB_TOKEN, its rate-limit handling), sending Intake's own
// User-Agent (http.go). Never the github.com web pages, releases.atom
// included (research §5.3; ADR-0005 calls the REST API the sanctioned
// path).
//
// Each published full release is an entry keyed on its release id, titled
// by its name (else its tag), linked to its release page, dated by
// published_at.
//
// Paging: the first page (releasesPerPage, newest first) is conditional on
// the stored ETag; a 304 means nothing changed. Otherwise the fetcher
// reads further pages while the page it just read was full and its newest
// non-draft release (prereleases included) is within the run's age cutoff
// (Intake.MaxAge), at most maxReleasePages: a busy release-candidate
// stream cannot push a full release off the first page, and an old
// repository costs one request.
//
// Prerelease rule: prereleases (and drafts) are not Leads. A release
// candidate is not news a Reader can act on, and projects that cut them
// (Kubernetes' alphas, betas and rcs) would bury the full releases. A
// skipped prerelease is not recorded as seen, so if its maintainers later
// promote that same release to a full one, the next run lists it then.
type githubFetcher struct{}

const (
	// releasesPerPage is how many releases one request reads.
	releasesPerPage = 10
	// maxReleasePages bounds the requests per repository Source per run.
	maxReleasePages = 5
)

// release is the part of GitHub's release object a Lead needs.
type release struct {
	ID          int64      `json:"id"`
	TagName     string     `json:"tag_name"`
	Name        string     `json:"name"`
	HTMLURL     string     `json:"html_url"`
	Draft       bool       `json:"draft"`
	Prerelease  bool       `json:"prerelease"`
	PublishedAt *time.Time `json:"published_at"`
}

func (githubFetcher) fetch(ctx context.Context, in *Intake, src Source, prev SourceState) (Fetched, error) {
	client, err := in.githubClient()
	if err != nil {
		return Fetched{}, err
	}
	cutoff := in.now().Add(-in.maxAge())
	state := prev
	var entries []Entry
	for page := 1; page <= maxReleasePages; page++ {
		path := fmt.Sprintf("/repos/%s/releases?per_page=%d", src.Repository, releasesPerPage)
		etag := prev.ETag // only the first page is conditional
		if page > 1 {
			path += fmt.Sprintf("&page=%d", page)
			etag = ""
		}
		var releases []release
		header, notModified, err := client.GetIfNoneMatch(ctx, path, etag, &releases)
		if err != nil {
			return Fetched{}, err
		}
		if notModified {
			state.ETag = revalidated(header, "ETag", prev.ETag)
			return Fetched{State: state}, nil
		}
		if page == 1 {
			state.ETag = header.Get("ETag")
		}
		for _, r := range releases {
			if !r.Draft && !r.Prerelease {
				entries = append(entries, r.entry())
			}
		}
		if newest := newestRelease(releases); len(releases) < releasesPerPage || newest.Before(cutoff) {
			break
		}
	}
	return Fetched{Entries: entries, State: state}, nil
}

func (r release) entry() Entry {
	e := Entry{
		Key:   strconv.FormatInt(r.ID, 10),
		Title: oneLine(cmp.Or(r.Name, r.TagName)),
		Link:  r.HTMLURL,
	}
	if r.PublishedAt != nil {
		e.Published = r.PublishedAt.UTC()
	}
	return e
}

// newestRelease is the latest published_at of a page's non-draft releases,
// prereleases included; zero when none has one.
func newestRelease(releases []release) time.Time {
	var newest time.Time
	for _, r := range releases {
		if !r.Draft && r.PublishedAt != nil && r.PublishedAt.After(newest) {
			newest = *r.PublishedAt
		}
	}
	return newest
}

// githubClient is the Fact Pack fetcher's GitHub client, pointed at
// GitHubURL, with the token read through Getenv: a missing token is an
// error for repository Sources only (Run reports each as failed). It sends
// Intake's UserAgent, not the fetch command's, so GitHub sees who polls.
func (in *Intake) githubClient() (*github.Client, error) {
	getenv := in.Getenv
	if getenv == nil {
		getenv = os.Getenv
	}
	client, err := github.NewClient(getenv)
	if err != nil {
		return nil, err
	}
	client.BaseURL = cmp.Or(in.GitHubURL, github.APIURL)
	client.UserAgent = UserAgent
	// Each request gets Intake's timeout (http.go) over Intake's transport;
	// the client's own retries and rate-limit waits sit between requests.
	client.HTTP = &http.Client{Transport: in.httpClient().Transport, Timeout: in.timeout()}
	return client, nil
}

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
// the one GITHUB_TOKEN, its User-Agent, its rate-limit handling). Never
// the github.com web pages, releases.atom included (research §5.3; ADR-0005
// calls the REST API the sanctioned path).
//
// One request per Source per run: the newest releasesPerPage releases.
// Each published full release is an item keyed on its release id, titled
// by its name (else its tag), linked to its release page, dated by
// published_at.
//
// Prerelease rule: prereleases (and drafts) are not Leads. A release
// candidate is not news a Reader can act on, and projects that cut them
// (Kubernetes' alphas, betas and rcs) would bury the full releases. A
// skipped prerelease is not recorded as seen, so if its maintainers later
// promote that same release to a full one, the next run lists it then.
type githubFetcher struct{}

// releasesPerPage is how many releases one request reads: more than any
// watched project publishes between two daily runs.
const releasesPerPage = 10

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
	var releases []release
	path := fmt.Sprintf("/repos/%s/releases?per_page=%d", src.Repository, releasesPerPage)
	header, notModified, err := client.GetIfNoneMatch(ctx, path, prev.ETag, &releases)
	if err != nil {
		return Fetched{}, err
	}
	if notModified {
		// Nothing new; a 304 may carry a fresher ETag (RFC 9110 §15.4.5).
		prev.ETag = cmp.Or(header.Get("ETag"), prev.ETag)
		return Fetched{State: prev}, nil
	}
	state := prev
	state.ETag = header.Get("ETag")
	var items []Item
	for _, r := range releases {
		if r.Draft || r.Prerelease {
			continue
		}
		item := Item{
			Key:   strconv.FormatInt(r.ID, 10),
			Title: oneLine(cmp.Or(r.Name, r.TagName)),
			Link:  r.HTMLURL,
		}
		if r.PublishedAt != nil {
			item.Published = r.PublishedAt.UTC()
		}
		items = append(items, item)
	}
	return Fetched{Items: items, State: state}, nil
}

// githubClient is the Fact Pack fetcher's GitHub client, pointed at
// GitHubURL, with the token read through Getenv: a missing token is an
// error for repository Sources only (Run reports each as failed). It keeps
// the client's User-Agent, the one make fetch sends.
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
	// Each request gets Intake's timeout (http.go) over Intake's transport;
	// the client's own retries and rate-limit waits sit between requests.
	client.HTTP = &http.Client{Transport: in.httpClient().Transport, Timeout: in.timeout()}
	return client, nil
}

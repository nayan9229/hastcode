package content

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"
)

// FactPackReleaseCap is how many releases a Fact Pack lists at most: the
// fetch command asks GitHub for the latest this many, so a count of
// releases that reaches it is a floor ("10 or more").
const FactPackReleaseCap = 10

// FactPack is the machine-fetched facts about a repository that a Repo
// Profile's fact sections (header card, Signals, Changelog) render from
// (GLOSSARY "Fact Pack"). It is a JSON file beside the Profile's Markdown,
// named after it (FactPackFile), written by the fetch command and never by
// hand; a Refresh rewrites it. Fields are the spec's Fact Pack v1 (one
// GitHub REST call each); where they overlap the exchange JSON of the
// original brief they keep its names (slug, url, name, language, license,
// stars, latest_release, manifests). No key is omitted when empty, so a
// Refresh diff is a change of values, never of shape.
type FactPack struct {
	// FetchedAt is when the fetch command ran; the builder warns when it
	// is more than 30 days before the build clock.
	FetchedAt time.Time `json:"fetched_at"`
	// FullName is the repository's full name, owner/name. It must equal
	// the Profile's `repo` frontmatter. Its JSON key is `slug`, the
	// original brief's exchange name.
	FullName string `json:"slug"`
	// Name is the repository name without the owner.
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Homepage    string   `json:"homepage"`
	Topics      []string `json:"topics"`
	// License is nil when GitHub reports no license.
	License       *License `json:"license"`
	DefaultBranch string   `json:"default_branch"`
	// Language is the primary language; Languages are bytes by language.
	Language   string           `json:"language"`
	Languages  map[string]int64 `json:"languages"`
	Stars      int              `json:"stars"`
	Forks      int              `json:"forks"`
	Watchers   int              `json:"watchers"`
	OpenIssues int              `json:"open_issues"`
	CreatedAt  time.Time        `json:"created_at"`
	// PushedAt is the last commit on any branch.
	PushedAt time.Time `json:"pushed_at"`
	// LatestRelease is the newest release that is not a prerelease; nil
	// when there is none. Releases are the latest FactPackReleaseCap,
	// newest first.
	LatestRelease *Release     `json:"latest_release"`
	Releases      []Release    `json:"releases"`
	Contributors  Contributors `json:"contributors"`
	Files         RepoFiles    `json:"files"`
	README        READMERef    `json:"readme"`
	// Stargazers is the stargazers/count endpoint's answer; it reserves the
	// place for star history, which v1 does not store.
	Stargazers Stargazers `json:"stargazers"`
	// Sources are the github.com pages behind the calls the fetch made (never
	// API URLs), each once, merged into the Profile's Sources.
	Sources []string `json:"sources"`
}

// License is the repository's license as GitHub classifies it.
type License struct {
	// SPDXID is the SPDX identifier ("Apache-2.0"); "NOASSERTION" when
	// GitHub could not classify the license.
	SPDXID string `json:"spdx_id"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	// NeedsReview is set for NOASSERTION: the Author reads the LICENSE file
	// by hand (ADR-0005: it means "check", not "fails the floor").
	NeedsReview bool `json:"needs_review"`
}

// Release is one GitHub release.
type Release struct {
	Tag        string    `json:"tag"`
	Name       string    `json:"name"`
	Date       time.Time `json:"date"`
	URL        string    `json:"url"`
	Prerelease bool      `json:"prerelease"`
}

// Contributors is the contributor count and the top ten by contributions,
// from the contributors list (not the stats endpoint).
type Contributors struct {
	Total int           `json:"total"`
	Top   []Contributor `json:"top"`
}

// Contributor is one contributor and their contribution count.
type Contributor struct {
	Login         string `json:"login"`
	URL           string `json:"url"`
	Contributions int    `json:"contributions"`
}

// RepoFiles records which well-known files and directories exist at the
// repository root or under .github: the path found, "" when absent. CI and
// Manifests list every path found.
type RepoFiles struct {
	README        string   `json:"readme"`
	License       string   `json:"license"`
	Security      string   `json:"security"`
	Contributing  string   `json:"contributing"`
	CodeOfConduct string   `json:"code_of_conduct"`
	CI            []string `json:"ci"`
	Manifests     []string `json:"manifests"`
	Dockerfile    string   `json:"dockerfile"`
	Docs          string   `json:"docs"`
	Examples      string   `json:"examples"`
}

// READMERef locates the README the Profile was drafted against: its page
// and blob SHA. The text itself is never committed.
type READMERef struct {
	URL string `json:"url"`
	SHA string `json:"sha"`
}

// Stargazers is the stargazer count from the dedicated endpoint.
type Stargazers struct {
	Count int `json:"count"`
}

// FactPackFile is the Fact Pack's content-relative path for a Repo
// Profile's Markdown file: the same name with a .json extension, beside it
// ("explainers/distributed-systems/apache-iggy.md" ->
// "explainers/distributed-systems/apache-iggy.json").
func FactPackFile(markdownFile string) string {
	return strings.TrimSuffix(markdownFile, path.Ext(markdownFile)) + ".json"
}

// ParseFactPack parses a Fact Pack file. Unknown keys are an error, so a
// Fact Pack written for a newer shape fails loudly rather than rendering
// without the facts it carries.
func ParseFactPack(data []byte) (*FactPack, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var fp FactPack
	if err := dec.Decode(&fp); err != nil {
		return nil, fmt.Errorf("not a Fact Pack: %w", err)
	}
	return &fp, nil
}

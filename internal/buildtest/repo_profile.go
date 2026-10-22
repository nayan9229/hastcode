package buildtest

import (
	"encoding/json"
	"path/filepath"
)

// Repo Profile fixtures: a valid Repo Profile (an Explainer of kind
// repo-profile with a `repo` key) and the Fact Pack written beside it.
// FactPack.JSON spells every key by hand, independently of the parser in
// internal/content, so a drift between the two fails the Build-seam tests.

// WithRepoProfile adds a Repo Profile like WithExplainer and, when
// e.FactPack is set, its Fact Pack as explainers/<vertical>/<slug>.json
// (the JSON is named after the Markdown file and sits beside it).
func (t Tree) WithRepoProfile(e Explainer) Tree {
	t.WithExplainer(e)
	if e.FactPack != nil {
		t[filepath.Join("explainers", e.Vertical, e.Slug+".json")] = e.FactPack.JSON()
	}
	return t
}

// ValidRepoProfile returns a Repo Profile of apache/iggy that satisfies the
// Explainer skeleton (TL;DR, Mermaid fence, Sources), carries the Repo
// Profile headings in the spec's order with five FAQ questions, the
// AI-assistance note from the Editorial Policy, a `who_for` line, and a
// valid Fact Pack (ValidFactPack) beside it. It names no hero: the build
// draws its hero card from the Fact Pack. Set Hero, HeroAlt and HeroImage
// (e.g. HeroPNG(1600, 900)) for a Profile with a sourced hero.
func ValidRepoProfile() Explainer {
	fp := ValidFactPack()
	return Explainer{
		Vertical:    "distributed-systems",
		Slug:        "apache-iggy",
		Title:       "Apache Iggy: a persistent message stream in Rust",
		Description: "Iggy is a persistent message streaming platform in Rust that trades Kafka's ecosystem for a smaller, faster single binary.",
		Kind:        "repo-profile",
		Repo:        "apache/iggy",
		WhoFor:      "Teams that need Kafka-style streams without running a JVM cluster.",
		Tags:        []string{"streaming", "rust"},
		Published:   "2026-10-21T09:00:00Z",
		Reviewed:    "2026-10-21T09:00:00Z",
		Author:      DefaultAuthor().Slug,
		FactPack:    &fp,
		Body:        ValidRepoProfileBody,
	}
}

// ValidRepoProfileBody is the Markdown body of ValidRepoProfile: the
// opinion sections in the spec's fixed order (Signals and Changelog are
// never written; they render from the Fact Pack).
const ValidRepoProfileBody = `Iggy is a message streaming server that keeps Kafka's partitioned-log model and drops the JVM.

> Drafted with AI assistance, researched, verified and edited by the HastCode editorial team.

## TL;DR

- One Rust binary serves QUIC, TCP and HTTP clients.
- Streams are append-only logs partitioned like Kafka topics.
- Clustering is on the roadmap, not in a release.

## In 60 seconds

` + "```sh\ndocker run -p 8090:8090 apache/iggy:latest\n```" + `

## Architecture

The server owns the log; clients talk to it over three transports.

` + "```mermaid\n%% title: Iggy's transports\nflowchart LR\n  Producer --> Server --> Consumer\n```" + `

## Use it if

- You want a single binary and a small footprint.

## Don't use it if

- You need multi-node replication today.

## Alternatives

- Apache Kafka: the incumbent, with the ecosystem.
- NATS JetStream: lighter, less log-shaped.
- Redpanda: Kafka-compatible, also no JVM.

## FAQ

### Is Iggy production ready?

Not for multi-node deployments yet.

### Does Iggy support Kafka clients?

No; it has its own protocol and SDKs.

### What transports does Iggy speak?

QUIC, TCP and HTTP.

### Is Iggy an Apache project?

Yes, it entered the Apache Incubator in 2025.

### How is Iggy licensed?

Apache 2.0.

## Verdict

Adopt it for single-node streaming where a small binary matters; wait for clustering otherwise.

## Sources

- [apache/iggy on GitHub](https://github.com/apache/iggy)
`

// FactPack is a fixture Fact Pack: the machine-fetched facts a Repo
// Profile's fact sections render from. Field names mirror the spec
// ("Content model"); JSON writes them with the exchange-JSON spellings.
// Dates are RFC 3339 strings so a test can blank or backdate one.
type FactPack struct {
	FetchedAt     string
	FullName      string // owner/name; JSON key "slug"
	Name          string
	URL           string
	Description   string
	Homepage      string
	Topics        []string
	License       *License // nil = GitHub reports no license
	DefaultBranch string
	Language      string
	Languages     map[string]int
	Stars         int
	Forks         int
	Watchers      int
	OpenIssues    int
	CreatedAt     string
	PushedAt      string
	LatestRelease *Release // nil = no release yet
	Releases      []Release
	Contributors  Contributors
	Files         RepoFiles
	README        README
	Stargazers    int
	Sources       []string
}

// License is the Fact Pack's license block. NeedsReview is set when
// GitHub could not classify it (NOASSERTION).
type License struct {
	SPDXID      string
	Name        string
	URL         string
	NeedsReview bool
}

// Release is one of the Fact Pack's releases.
type Release struct {
	Tag        string
	Name       string
	Date       string
	URL        string
	Prerelease bool
}

// Contributors is the contributor count and the top contributors.
type Contributors struct {
	Total int
	Top   []Contributor
}

// Contributor is one of the top contributors.
type Contributor struct {
	Login         string
	URL           string
	Contributions int
}

// RepoFiles records which well-known files and directories exist: the path
// found, or "" when absent.
type RepoFiles struct {
	README, License, Security, Contributing, CodeOfConduct string
	CI, Manifests                                          []string
	Dockerfile, Docs, Examples                             string
}

// README is the README's blob reference.
type README struct {
	URL string
	SHA string
}

// ValidFactPack returns the Fact Pack of apache/iggy, fetched two days
// before ValidRepoProfile's reviewed date.
func ValidFactPack() FactPack {
	return FactPack{
		FetchedAt:     "2026-10-19T06:00:00Z",
		FullName:      "apache/iggy",
		Name:          "iggy",
		URL:           "https://github.com/apache/iggy",
		Description:   "Persistent message streaming platform written in Rust",
		Homepage:      "https://iggy.apache.org",
		Topics:        []string{"streaming", "rust", "message-broker"},
		License:       &License{SPDXID: "Apache-2.0", Name: "Apache License 2.0", URL: "https://spdx.org/licenses/Apache-2.0.html"},
		DefaultBranch: "main",
		Language:      "Rust",
		Languages:     map[string]int{"Rust": 2345678, "Shell": 12345},
		Stars:         2468,
		Forks:         135,
		Watchers:      42,
		OpenIssues:    57,
		CreatedAt:     "2023-07-24T10:00:00Z",
		PushedAt:      "2026-10-18T14:03:00Z",
		LatestRelease: &Release{Tag: "server-0.5.0", Name: "Server 0.5.0", Date: "2026-10-10T12:00:00Z", URL: "https://github.com/apache/iggy/releases/tag/server-0.5.0"},
		Releases: []Release{
			{Tag: "server-0.6.0-rc1", Name: "Server 0.6.0 RC1", Date: "2026-10-16T12:00:00Z", URL: "https://github.com/apache/iggy/releases/tag/server-0.6.0-rc1", Prerelease: true},
			{Tag: "server-0.5.0", Name: "Server 0.5.0", Date: "2026-10-10T12:00:00Z", URL: "https://github.com/apache/iggy/releases/tag/server-0.5.0"},
			{Tag: "server-0.4.0", Name: "Server 0.4.0", Date: "2026-06-02T12:00:00Z", URL: "https://github.com/apache/iggy/releases/tag/server-0.4.0"},
		},
		Contributors: Contributors{
			Total: 61,
			Top: []Contributor{
				{Login: "spetz", URL: "https://github.com/spetz", Contributions: 1800},
				{Login: "hubcio", URL: "https://github.com/hubcio", Contributions: 420},
			},
		},
		Files: RepoFiles{
			README: "README.md", License: "LICENSE", Security: ".github/SECURITY.md",
			Contributing: "CONTRIBUTING.md", CodeOfConduct: "CODE_OF_CONDUCT.md",
			CI: []string{".github/workflows"}, Manifests: []string{"Cargo.toml"},
			Dockerfile: "Dockerfile", Docs: "docs", Examples: "examples",
		},
		README:     README{URL: "https://github.com/apache/iggy/blob/main/README.md", SHA: "7c6b2a0f9e8d4c3b1a0f9e8d7c6b5a4f3e2d1c0b"},
		Stargazers: 2468,
		Sources: []string{
			"https://github.com/apache/iggy",
			"https://github.com/apache/iggy/releases",
			"https://github.com/apache/iggy/graphs/contributors",
			"https://github.com/apache/iggy/tree/main/.github",
			"https://github.com/apache/iggy/blob/main/README.md",
		},
	}
}

// JSON renders the fixture as the Fact Pack file. The key names are the
// spec's, spelled here by hand (not by the parser's struct tags).
func (f FactPack) JSON() []byte {
	var license, latest any
	if f.License != nil {
		license = map[string]any{"spdx_id": f.License.SPDXID, "name": f.License.Name, "url": f.License.URL, "needs_review": f.License.NeedsReview}
	}
	if f.LatestRelease != nil {
		latest = release(*f.LatestRelease)
	}
	releases := []any{}
	for _, r := range f.Releases {
		releases = append(releases, release(r))
	}
	top := []any{}
	for _, c := range f.Contributors.Top {
		top = append(top, map[string]any{"login": c.Login, "url": c.URL, "contributions": c.Contributions})
	}
	doc := map[string]any{
		"fetched_at":     f.FetchedAt,
		"slug":           f.FullName,
		"name":           f.Name,
		"url":            f.URL,
		"description":    f.Description,
		"homepage":       f.Homepage,
		"topics":         f.Topics,
		"license":        license,
		"default_branch": f.DefaultBranch,
		"language":       f.Language,
		"languages":      f.Languages,
		"stars":          f.Stars,
		"forks":          f.Forks,
		"watchers":       f.Watchers,
		"open_issues":    f.OpenIssues,
		"created_at":     f.CreatedAt,
		"pushed_at":      f.PushedAt,
		"latest_release": latest,
		"releases":       releases,
		"contributors":   map[string]any{"total": f.Contributors.Total, "top": top},
		"files": map[string]any{
			"readme": f.Files.README, "license": f.Files.License, "security": f.Files.Security,
			"contributing": f.Files.Contributing, "code_of_conduct": f.Files.CodeOfConduct,
			"ci": orEmpty(f.Files.CI), "manifests": orEmpty(f.Files.Manifests),
			"dockerfile": f.Files.Dockerfile, "docs": f.Files.Docs, "examples": f.Files.Examples,
		},
		"readme":     map[string]any{"url": f.README.URL, "sha": f.README.SHA},
		"stargazers": map[string]any{"count": f.Stargazers},
		"sources":    orEmpty(f.Sources),
	}
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		panic(err)
	}
	return append(data, '\n')
}

func release(r Release) map[string]any {
	return map[string]any{"tag": r.Tag, "name": r.Name, "date": r.Date, "url": r.URL, "prerelease": r.Prerelease}
}

func orEmpty(ss []string) []string {
	if ss == nil {
		return []string{}
	}
	return ss
}

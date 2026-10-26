package intake

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"

	"hastcode.com/site/internal/content"
)

// DefaultWatchlist is the Source Watchlist's file name at the repo root.
const DefaultWatchlist = "sources.yaml"

// Kinds a Source may suggest for its Leads: the three Explainer kinds a
// Lead can become, and a Brief. Repo Profiles are chosen by other means
// (docs/editorial/repo-profile-selection.md).
var Kinds = []string{string(content.CaseStudy), string(content.Concept), string(content.Comparison), "brief"}

// SourceType is how a Source is fetched: the key of the fetcher registry
// (fetch.go).
type SourceType string

const (
	// FeedSource has a feed URL (RSS 2.0 or Atom).
	FeedSource SourceType = "feed"
	// RepositorySource has a GitHub repository (owner/name) whose releases
	// are its items.
	RepositorySource SourceType = "repository"
)

// Watchlist is the Source Watchlist file.
type Watchlist struct {
	Sources []Source `yaml:"sources"`
}

// Source is one Watchlist entry (GLOSSARY "Source"). Exactly one of Feed
// and Repository is set; Vertical and Kind are copied onto every Lead it
// yields, never guessed.
type Source struct {
	// Name identifies the Source in the queue and the state; it is part of
	// every Lead ID, so renaming a Source lists its items again.
	Name string `yaml:"name"`
	// Site is the publisher's home page.
	Site string `yaml:"site"`
	// Feed is an RSS 2.0 or Atom feed URL.
	Feed string `yaml:"feed"`
	// Repository is a GitHub repository, owner/name.
	Repository string `yaml:"repository"`
	// Vertical is a Vertical slug from site.yaml.
	Vertical string `yaml:"vertical"`
	// Kind is the kind its Leads usually become (Kinds).
	Kind string `yaml:"kind"`
	// Reason is the one line saying why the Source is on the Watchlist.
	Reason string `yaml:"reason"`
}

// Type is how the Source is fetched.
func (s Source) Type() SourceType {
	if s.Repository != "" {
		return RepositorySource
	}
	return FeedSource
}

// LoadWatchlist reads and validates the Watchlist file against the known
// Verticals. Every problem is reported, each naming its Source: an unknown
// Vertical or kind, both or neither of feed and repository, a duplicate
// name, a missing name, site or reason, a malformed URL or repository, an
// unknown key.
func LoadWatchlist(file string, verticals []Vertical) (*Watchlist, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("watchlist: %w", err)
	}
	var wl Watchlist
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&wl); err != nil {
		return nil, fmt.Errorf("watchlist %s: %w", file, err)
	}
	if err := wl.validate(verticals); err != nil {
		return nil, fmt.Errorf("watchlist %s:\n%w", file, err)
	}
	return &wl, nil
}

func (wl *Watchlist) validate(verticals []Vertical) error {
	var errs []error
	if len(wl.Sources) == 0 {
		errs = append(errs, errors.New("no sources"))
	}
	seen := map[string]bool{}
	for i, s := range wl.Sources {
		name := s.Name
		if name == "" {
			name = fmt.Sprintf("#%d", i+1)
		}
		bad := func(format string, args ...any) {
			errs = append(errs, fmt.Errorf("source %q: %s", name, fmt.Sprintf(format, args...)))
		}
		switch {
		case s.Name == "":
			bad("name is required")
		case seen[s.Name]:
			bad("duplicate name")
		}
		seen[s.Name] = true
		if !isHTTPURL(s.Site) {
			bad("site %q must be an absolute http(s) URL", s.Site)
		}
		switch {
		case s.Feed != "" && s.Repository != "":
			bad("has both feed and repository; give one")
		case s.Feed == "" && s.Repository == "":
			bad("has neither feed nor repository; give one")
		case s.Feed != "" && !isHTTPURL(s.Feed):
			bad("feed %q must be an absolute http(s) URL", s.Feed)
		case s.Repository != "":
			if _, err := content.ParseRepo(s.Repository); err != nil {
				bad("repository %v", err)
			}
		}
		if !knownVertical(verticals, s.Vertical) {
			bad("unknown vertical %q (site.yaml has %s)", s.Vertical, verticalSlugs(verticals))
		}
		if !slices.Contains(Kinds, s.Kind) {
			bad("unknown kind %q (one of %s)", s.Kind, strings.Join(Kinds, ", "))
		}
		if strings.TrimSpace(s.Reason) == "" {
			bad("reason is required")
		}
	}
	return errors.Join(errs...)
}

func isHTTPURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

func knownVertical(verticals []Vertical, slug string) bool {
	return slices.ContainsFunc(verticals, func(v Vertical) bool { return v.Slug == slug })
}

func verticalSlugs(verticals []Vertical) string {
	var s []string
	for _, v := range verticals {
		s = append(s, v.Slug)
	}
	return strings.Join(s, ", ")
}

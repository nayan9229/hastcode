package factpack

import (
	"bytes"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	"hastcode.com/site/internal/content"
)

// changed reports whether two Fact Packs differ in anything but
// fetched_at: their files, written with fetched_at blanked, differ.
func changed(old, new *content.FactPack) bool {
	return !bytes.Equal(withoutFetchedAt(old), withoutFetchedAt(new))
}

func withoutFetchedAt(fp *content.FactPack) []byte {
	cp := *fp
	cp.FetchedAt = time.Time{}
	data, err := Marshal(&cp)
	if err != nil {
		panic(err) // a FactPack always marshals
	}
	return data
}

// diff says what changed from old to new, fetched_at aside, as short
// phrases for the Refresh line and the pull request: the counts with their
// values, the rest by name. Facts the phrases do not cover fall back to
// "other facts changed", so a line is never empty when changed is true.
func diff(old, new *content.FactPack) []string {
	var out []string
	count := func(label string, a, b int) {
		if a != b {
			out = append(out, fmt.Sprintf("%s %d -> %d", label, a, b))
		}
	}
	text := func(label, a, b string) {
		if a != b {
			out = append(out, label+" changed")
		}
	}

	count("stars", old.Stars, new.Stars)
	count("forks", old.Forks, new.Forks)
	count("watchers", old.Watchers, new.Watchers)
	count("open issues", old.OpenIssues, new.OpenIssues)
	if !old.PushedAt.Equal(new.PushedAt) {
		out = append(out, fmt.Sprintf("last commit %s -> %s", old.PushedAt.Format(time.DateOnly), new.PushedAt.Format(time.DateOnly)))
	}
	if a, b := releaseTag(old.LatestRelease), releaseTag(new.LatestRelease); a != b {
		out = append(out, fmt.Sprintf("latest release %s -> %s", a, b))
	}
	if added := newTags(old.Releases, new.Releases); len(added) > 0 {
		out = append(out, "new releases "+strings.Join(added, ", "))
	} else if !slices.Equal(old.Releases, new.Releases) {
		out = append(out, "releases changed")
	}
	count("contributors", old.Contributors.Total, new.Contributors.Total)
	if !slices.Equal(old.Contributors.Top, new.Contributors.Top) {
		out = append(out, "top contributors changed")
	}
	out = append(out, fileChanges(old.Files, new.Files)...)
	text("README", old.README.SHA, new.README.SHA)
	count("stargazer count", old.Stargazers.Count, new.Stargazers.Count)
	text("description", old.Description, new.Description)
	text("homepage", old.Homepage, new.Homepage)
	if !slices.Equal(old.Topics, new.Topics) {
		out = append(out, "topics changed")
	}
	text("license", licenseID(old.License), licenseID(new.License))
	text("default branch", old.DefaultBranch, new.DefaultBranch)
	text("language", old.Language, new.Language)
	if !maps.Equal(old.Languages, new.Languages) {
		out = append(out, "languages changed")
	}
	if old.Slug != new.Slug || old.Name != new.Name || old.URL != new.URL {
		out = append(out, "repository renamed "+old.Slug+" -> "+new.Slug)
	}
	if !slices.Equal(old.Sources, new.Sources) {
		out = append(out, "sources changed")
	}
	if len(out) == 0 {
		out = append(out, "other facts changed")
	}
	return out
}

func releaseTag(r *content.Release) string {
	if r == nil {
		return "none"
	}
	return r.Tag
}

// newTags are the tags in new that old lacks, in new's order.
func newTags(old, new []content.Release) []string {
	known := map[string]bool{}
	for _, r := range old {
		known[r.Tag] = true
	}
	var tags []string
	for _, r := range new {
		if !known[r.Tag] {
			tags = append(tags, r.Tag)
		}
	}
	return tags
}

func licenseID(l *content.License) string {
	if l == nil {
		return ""
	}
	return l.SPDXID
}

// fileChanges describes each well-known file that appeared, went or moved.
func fileChanges(old, new content.RepoFiles) []string {
	var out []string
	one := func(label, a, b string) {
		switch {
		case a == b:
		case a == "":
			out = append(out, "files: "+label+" added ("+b+")")
		case b == "":
			out = append(out, "files: "+label+" removed ("+a+")")
		default:
			out = append(out, "files: "+label+" "+a+" -> "+b)
		}
	}
	one("README", old.README, new.README)
	one("license", old.License, new.License)
	one("security policy", old.Security, new.Security)
	one("contributing guide", old.Contributing, new.Contributing)
	one("code of conduct", old.CodeOfConduct, new.CodeOfConduct)
	if !slices.Equal(old.CI, new.CI) {
		out = append(out, "files: CI configuration changed")
	}
	if !slices.Equal(old.Manifests, new.Manifests) {
		out = append(out, "files: manifests changed")
	}
	one("Dockerfile", old.Dockerfile, new.Dockerfile)
	one("docs", old.Docs, new.Docs)
	one("examples", old.Examples, new.Examples)
	return out
}

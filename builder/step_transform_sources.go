package builder

import (
	"cmp"
	"strings"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/markdown"
)

// Frontmatter `sources` (spec "Content model": used in addition to inline
// links). An Explainer's are appended to its ## Sources section (the last
// section, so the end of the body), which the page and the full-text feed
// both render; a Brief's are listed with its source_url (briefPage). Each
// URL appears once: a source already linked there is not repeated.
func init() {
	registerStep(phaseTransform, 30, "merge explainer frontmatter sources", mergeExplainerSources)
}

func mergeExplainerSources(s *site) error {
	for _, e := range s.publishedExplainers() {
		_, sections := e.Body.Sections()
		var linked []string
		for _, sec := range sections {
			if sec.Heading == headingSources {
				linked = e.Body.LinkDestinations(sec.Blocks)
			}
		}
		var extra []markdown.Link
		for _, src := range newSources(linked, e.Sources) {
			extra = append(extra, markdown.Link{Text: cmp.Or(src.Title, src.URL), URL: src.URL})
		}
		e.Body.AppendLinkList(extra)
	}
	return nil
}

// newSources are the frontmatter sources whose URL is neither in linked
// nor earlier in the list, in order. Sources without a URL are dropped.
func newSources(linked []string, sources []content.Source) []content.Source {
	seen := map[string]bool{}
	for _, u := range linked {
		seen[sourceKey(u)] = true
	}
	var out []content.Source
	for _, src := range sources {
		k := sourceKey(src.URL)
		if k == "" || seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, src)
	}
	return out
}

// sourceKey compares source URLs: surrounding space and a trailing slash
// do not make a different source.
func sourceKey(u string) string {
	return strings.TrimSuffix(strings.TrimSpace(u), "/")
}

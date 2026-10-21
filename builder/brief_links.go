package builder

import (
	"net/url"
	"strings"

	"github.com/yuin/goldmark/ast"

	"hastcode.com/site/internal/content"
)

// A Brief links into Explainers from its body. The Author may write the
// link as the Explainer's bare slug ("raft-leader-election") or its site
// path ("/distributed-systems/raft-leader-election", trailing slash and
// #fragment optional); the build checks the target exists and renders the
// canonical trailing-slash URL.

// explainerLink is a Brief body link that points at an Explainer.
type explainerLink struct {
	node *ast.Link
	slug string
	// target is the Explainer linked to; nil when no such Explainer is
	// loaded (a broken internal link).
	target   *content.Explainer
	fragment string
}

// explainerLinks returns the links in b's body that point at an
// Explainer, in document order. External URLs, fragment-only links and
// links to other site pages are not Explainer links.
func (s *site) explainerLinks(b *content.Brief) []explainerLink {
	var links []explainerLink
	_ = ast.Walk(b.Body.Root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		link, ok := n.(*ast.Link)
		if !entering || !ok {
			return ast.WalkContinue, nil
		}
		if l, ok := s.resolveExplainerLink(string(link.Destination)); ok {
			l.node = link
			links = append(links, l)
		}
		return ast.WalkContinue, nil
	})
	return links
}

// resolveExplainerLink reports whether dest has the shape of an Explainer
// link (bare slug, or /{vertical}/{slug} for a configured Vertical) and,
// if so, which loaded Explainer it names.
func (s *site) resolveExplainerLink(dest string) (explainerLink, bool) {
	u, err := url.Parse(dest)
	if err != nil || u.Scheme != "" || u.Host != "" || u.Path == "" {
		return explainerLink{}, false
	}
	var vertical, slug string
	switch parts := strings.Split(strings.Trim(u.Path, "/"), "/"); {
	case !strings.HasPrefix(u.Path, "/") && len(parts) == 1:
		slug = parts[0]
	case strings.HasPrefix(u.Path, "/") && len(parts) == 2:
		if _, ok := s.verticals[parts[0]]; !ok {
			return explainerLink{}, false
		}
		vertical, slug = parts[0], parts[1]
	default:
		return explainerLink{}, false
	}
	l := explainerLink{slug: slug, fragment: u.Fragment}
	for _, e := range s.explainers {
		if e.Slug == slug && (vertical == "" || e.Vertical == vertical) {
			l.target = e
			break
		}
	}
	return l, true
}

// href is the canonical URL the link renders with.
func (l explainerLink) href() string {
	h := explainerURL(l.target)
	if l.fragment != "" {
		h += "#" + l.fragment
	}
	return h
}

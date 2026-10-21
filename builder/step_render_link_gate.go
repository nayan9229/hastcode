package builder

import (
	"bytes"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/net/html"

	"hastcode.com/site/subscribe"
)

func init() {
	// Last in the render phase: every page and asset is written by now.
	registerStep(phaseRender, 1000, "check internal links", checkInternalLinks)
}

// ruleBrokenLink is an internal link to a page, asset or feed that does
// not exist (spec user story 31). Every rendered page is checked here;
// Brief Explainer links and `related` slugs are also checked earlier, at
// the source, because the build rewrites or resolves them.
const ruleBrokenLink = "broken-link"

// checkInternalLinks reads every HTML page written to the output and
// checks that each internal href (site-relative, page-relative, or
// absolute on the site's own origin) resolves to a written file, to a
// directory index ("/a/" -> /a/index.html), or to a path the build does
// not write itself (linkTargetsNotWritten). A page with a broken link is a
// build error attributed to its content file ("" for pages without one,
// such as the home page or a listing) and is removed from the output.
func checkInternalLinks(s *site) error {
	root := filepath.Join(s.out.Dir, "static")
	written := map[string]bool{}
	var pages []string
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		u := "/" + filepath.ToSlash(rel)
		written[u] = true
		if strings.HasSuffix(u, ".html") {
			pages = append(pages, u)
		}
		return nil
	})
	if err != nil {
		return err
	}
	slices.Sort(pages)

	allowed := s.linkTargetsNotWritten()
	sources := s.pageSources()
	site, _ := url.Parse(s.cfg.SiteURL)
	for _, file := range pages {
		pageURL := strings.TrimSuffix(file, "index.html")
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(file[1:])))
		if err != nil {
			return err
		}
		var broken []string
		for _, href := range hrefs(data) {
			target, internal := internalTarget(site, pageURL, href)
			if !internal || allowed[target] || written[target] || (strings.HasSuffix(target, "/") && written[target+"index.html"]) {
				continue
			}
			if !slices.Contains(broken, href) {
				broken = append(broken, href)
				s.fail(sources[pageURL], ruleBrokenLink, "page %s links %s (%s), which is not a page, asset or feed on this site", pageURL, href, target)
			}
		}
		if len(broken) > 0 {
			if err := os.Remove(filepath.Join(root, filepath.FromSlash(file[1:]))); err != nil {
				return err
			}
		}
	}
	return nil
}

// internalTarget resolves href, found on the page at pageURL, to a site
// path. internal is false for links off the site (other origins, mailto:)
// and for links to a fragment of the same page.
func internalTarget(site *url.URL, pageURL, href string) (target string, internal bool) {
	href = strings.TrimSpace(href)
	if href == "" || strings.HasPrefix(href, "#") {
		return "", false
	}
	u, err := url.Parse(href)
	if err != nil {
		return href, true // unparseable: broken
	}
	if u.Scheme != "" || u.Host != "" {
		if (u.Scheme != "" && u.Scheme != "http" && u.Scheme != "https") || site == nil || u.Host != site.Host {
			return "", false
		}
	}
	if u.Path == "" { // "?q=…" on the same page
		return pageURL, true
	}
	// Resolves "../x/" and "x/" against the page; keeps a trailing slash.
	return (&url.URL{Path: pageURL}).ResolveReference(&url.URL{Path: u.Path}).Path, true
}

// hrefs returns every href attribute value in an HTML document.
func hrefs(doc []byte) []string {
	var out []string
	z := html.NewTokenizer(bytes.NewReader(doc))
	for {
		switch z.Next() {
		case html.ErrorToken:
			return out
		case html.StartTagToken, html.SelfClosingTagToken:
			for {
				key, val, more := z.TagAttr()
				if string(key) == "href" {
					out = append(out, string(val))
				}
				if !more {
					break
				}
			}
		}
	}
}

// linkTargetsNotWritten are site paths pages may link although this build
// does not write them: the subscribe function (a Vercel function), and
// paths written after this gate (feeds and sitemaps: write phase, ticket
// 09). The search page is rendered before the gate, so it is checked like
// any page.
func (s *site) linkTargetsNotWritten() map[string]bool {
	allowed := map[string]bool{
		subscribe.Path:      true,
		siteFeedURL:         true,
		briefsFeedURL:       true,
		"/sitemap.xml":      true,
		"/news-sitemap.xml": true,
	}
	for slug := range s.verticals {
		allowed[verticalURL(slug)+"feed.xml"] = true
	}
	return allowed
}

// pageSources maps the URL of every page built from a content file to that
// file, so a broken link is reported against the file to fix.
func (s *site) pageSources() map[string]string {
	sources := map[string]string{}
	for _, e := range s.explainers {
		sources[explainerURL(e)] = e.File
	}
	for _, b := range s.briefs {
		sources[briefURL(b)] = b.File
	}
	for _, p := range s.pages {
		sources[staticPageURL(p)] = p.File
	}
	for _, a := range s.authors {
		sources[authorURL(a.Slug)] = a.File
	}
	return sources
}

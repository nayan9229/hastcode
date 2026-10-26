package intake

import (
	"errors"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"hastcode.com/site/internal/content"
)

// pieceDirs are the content directories holding Pieces: Explainers (Repo
// Profiles among them, with their Fact Packs beside) and Briefs. Pages are
// not Pieces and cover nothing.
var pieceDirs = []string{"explainers", "briefs"}

// coveredLinks reads the content tree and returns the normalised link of
// every page a Piece already cites: each URL in a Piece's Markdown (its
// frontmatter Sources and its body) and each of the Sources of the Fact
// Pack beside it (content.FactPackFile), if any. A Fact Pack that does not
// parse is logged and skipped (the build reports it, rule fact-pack); any
// other JSON file is not read. It only reads the content tree. No content
// tree covers nothing.
func (in *Intake) coveredLinks() (map[string]bool, error) {
	covered := map[string]bool{}
	if in.ContentDir == "" {
		return covered, nil
	}
	for _, dir := range pieceDirs {
		err := filepath.WalkDir(filepath.Join(in.ContentDir, dir), func(file string, d fs.DirEntry, err error) error {
			if errors.Is(err, fs.ErrNotExist) && d == nil {
				return fs.SkipDir
			}
			if err != nil || d.IsDir() || filepath.Ext(file) != ".md" {
				return err
			}
			links, err := markdownLinks(file)
			if err != nil {
				return err
			}
			factPackLinks, err := in.factPackSources(content.FactPackFile(file))
			if err != nil {
				return err
			}
			for _, link := range append(links, factPackLinks...) {
				covered[normaliseLink(link)] = true
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return covered, nil
}

// urlPattern finds absolute URLs in Markdown and YAML: it stops at
// whitespace, quotes, angle brackets, parentheses and square brackets, so
// `[text](url)`, `<url>` and `url: "…"` all yield the bare URL.
var urlPattern = regexp.MustCompile("https?://[^\\s<>\"'`()\\[\\]]+")

// markdownLinks are the URLs a Piece's Markdown file cites, frontmatter
// and body.
func markdownLinks(file string) ([]string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var links []string
	for _, m := range urlPattern.FindAllString(string(data), -1) {
		links = append(links, strings.TrimRight(m, ".,;:!?*_"))
	}
	return links, nil
}

// factPackSources are the Sources of the Fact Pack file (merged into its
// Repo Profile's), none when there is no such file. One that does not
// parse is logged and skipped.
func (in *Intake) factPackSources(file string) ([]string, error) {
	data, err := os.ReadFile(file)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	pack, err := content.ParseFactPack(data)
	if err != nil {
		in.logf("covered: skipping %s: %v", file, err)
		return nil, nil
	}
	return pack.Sources, nil
}

// trackingParams are query parameters that name where a reader came from,
// not which page: utm_* and these. `source` is one only as source=rss
// (Medium-hosted feeds tag their links so); any other value is kept, as it
// may select the page.
var trackingParams = []string{"fbclid", "gclid", "mc_cid", "mc_eid", "ref"}

// normaliseLink is the form two links to the same page share: http and
// https alike, host lower-case without "www." or a default port, no
// fragment, no tracking parameters, the rest of the query sorted, no
// trailing slash. A link that does not parse is only trimmed.
func normaliseLink(link string) string {
	u, err := url.Parse(strings.TrimSpace(link))
	if err != nil || u.Host == "" {
		return strings.TrimSpace(link)
	}
	host := strings.TrimPrefix(strings.ToLower(u.Hostname()), "www.")
	if port := u.Port(); port != "" && port != "80" && port != "443" {
		host += ":" + port
	}
	q := u.Query()
	for k, v := range q {
		key := strings.ToLower(k)
		feedTag := key == "source" && len(v) == 1 && strings.EqualFold(v[0], "rss")
		if strings.HasPrefix(key, "utm_") || slices.Contains(trackingParams, key) || feedTag {
			q.Del(k)
		}
	}
	norm := host + strings.TrimRight(u.EscapedPath(), "/")
	if len(q) > 0 {
		norm += "?" + q.Encode() // Encode sorts by key
	}
	return norm
}

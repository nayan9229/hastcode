package intake

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

// pieceDirs are the content directories holding Pieces: Explainers (Repo
// Profiles among them, with their Fact Packs beside) and Briefs. Pages are
// not Pieces and cover nothing.
var pieceDirs = []string{"explainers", "briefs"}

// coveredLinks reads the content tree and returns the normalised link of
// every page a Piece already cites: each URL in a Piece's Markdown (its
// frontmatter Sources and its body) and each of a Fact Pack's Sources.
// It only reads; Intake never writes under the content tree. No content
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
			if err != nil || d.IsDir() {
				return err
			}
			var links []string
			switch filepath.Ext(file) {
			case ".md":
				links, err = markdownLinks(file)
			case ".json":
				links, err = factPackSources(file)
			}
			if err != nil {
				return err
			}
			for _, link := range links {
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

// factPackSources are a Fact Pack's Sources (merged into its Repo
// Profile's). A JSON file that is not a Fact Pack yields none.
func factPackSources(file string) ([]string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var pack struct {
		Sources []string `json:"sources"`
	}
	if json.Unmarshal(data, &pack) != nil {
		return nil, nil
	}
	return pack.Sources, nil
}

// trackingParams are query parameters that name where a reader came from,
// not which page: utm_* and these.
var trackingParams = []string{"fbclid", "gclid", "mc_cid", "mc_eid", "ref", "source"}

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
	for k := range q {
		if strings.HasPrefix(strings.ToLower(k), "utm_") || slices.Contains(trackingParams, strings.ToLower(k)) {
			q.Del(k)
		}
	}
	norm := host + strings.TrimRight(u.EscapedPath(), "/")
	if len(q) > 0 {
		norm += "?" + q.Encode() // Encode sorts by key
	}
	return norm
}

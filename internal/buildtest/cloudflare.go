package buildtest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// WorkersConfig reads and decodes the Cloudflare Workers config the builder
// writes beside config.json (<out>/wrangler.json).
func WorkersConfig(tb testing.TB, outDir string) map[string]any {
	tb.Helper()
	data, err := os.ReadFile(filepath.Join(outDir, "wrangler.json"))
	if err != nil {
		tb.Fatalf("wrangler.json: %v", err)
	}
	var v map[string]any
	if err := json.Unmarshal(data, &v); err != nil {
		tb.Fatalf("wrangler.json does not decode: %v", err)
	}
	return v
}

// CloudflareHeaders returns the headers the static/_headers file adds to a
// response for target, a site path ("/about/") or an absolute URL. It
// applies Workers' documented rules
// (https://developers.cloudflare.com/workers/static-assets/headers/): a
// rule line is a path or an https URL with at most one "*" splat (any
// characters) and ":name" placeholders (anything but "/" in the path, "."
// in the host); indented "Name: value" lines follow; every matching rule
// contributes its headers, and a header set by two rules is joined with
// ", ". A path rule matches on any host.
func CloudflareHeaders(tb testing.TB, outDir, target string) http.Header {
	tb.Helper()
	u, err := url.Parse(target)
	if err != nil {
		tb.Fatalf("CloudflareHeaders(%q): %v", target, err)
	}
	got := http.Header{}
	for _, r := range cloudflareHeaderRules(tb, outDir) {
		if !r.matches(u) {
			continue
		}
		for name, value := range r.headers {
			if prev := got.Get(name); prev != "" {
				value = prev + ", " + value
			}
			got.Set(name, value)
		}
	}
	return got
}

type headerRule struct {
	host, path *regexp.Regexp // host nil: any host
	headers    map[string]string
}

func (r headerRule) matches(u *url.URL) bool {
	if r.host != nil && !r.host.MatchString(u.Host) {
		return false
	}
	return r.path.MatchString(u.Path)
}

func cloudflareHeaderRules(tb testing.TB, outDir string) []headerRule {
	tb.Helper()
	data, err := os.ReadFile(filepath.Join(outDir, "static", "_headers"))
	if err != nil {
		tb.Fatalf("_headers: %v", err)
	}
	var rules []headerRule
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		line := sc.Text()
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "" || strings.HasPrefix(trimmed, "#"):
		case line[0] == ' ' || line[0] == '\t':
			if len(rules) == 0 {
				tb.Fatalf("_headers: header line before any URL: %q", line)
			}
			name, value, ok := strings.Cut(trimmed, ":")
			if !ok {
				tb.Fatalf("_headers: header line without a colon: %q", line)
			}
			rules[len(rules)-1].headers[strings.TrimSpace(name)] = strings.TrimSpace(value)
		default:
			rules = append(rules, parseHeaderRule(tb, trimmed))
		}
	}
	return rules
}

func parseHeaderRule(tb testing.TB, pattern string) headerRule {
	tb.Helper()
	if strings.Count(pattern, "*") > 1 {
		tb.Fatalf("_headers: %q has more than one splat", pattern)
	}
	r := headerRule{headers: map[string]string{}}
	p := pattern
	if rest, ok := strings.CutPrefix(pattern, "https://"); ok {
		host, path, _ := strings.Cut(rest, "/")
		r.host = patternRegexp(host, "[^.]+")
		p = "/" + path
	} else if !strings.HasPrefix(pattern, "/") {
		tb.Fatalf("_headers: rule %q is neither a path nor an https URL", pattern)
	}
	r.path = patternRegexp(p, "[^/]+")
	return r
}

var placeholder = regexp.MustCompile(`:[A-Za-z]\w*`)

func patternRegexp(p, placeholderRE string) *regexp.Regexp {
	var b strings.Builder
	b.WriteString("^")
	for i, part := range strings.Split(p, "*") {
		if i > 0 {
			b.WriteString(".*")
		}
		lits := placeholder.Split(part, -1)
		for j, lit := range lits {
			if j > 0 {
				b.WriteString(placeholderRE)
			}
			b.WriteString(regexp.QuoteMeta(lit))
		}
	}
	b.WriteString("$")
	return regexp.MustCompile(b.String())
}

// CloudflareRedirects returns the rule lines of static/_redirects
// ("source destination [code]"), comments and blank lines dropped, fields
// separated by single spaces.
func CloudflareRedirects(tb testing.TB, outDir string) []string {
	tb.Helper()
	data, err := os.ReadFile(filepath.Join(outDir, "static", "_redirects"))
	if err != nil {
		tb.Fatalf("_redirects: %v", err)
	}
	var rules []string
	for _, line := range strings.Split(string(data), "\n") {
		if f := strings.Fields(line); len(f) > 0 && !strings.HasPrefix(f[0], "#") {
			rules = append(rules, strings.Join(f, " "))
		}
	}
	return rules
}

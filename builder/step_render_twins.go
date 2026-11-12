package builder

import (
	"bytes"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark/ast"

	"hastcode.com/site/internal/markdown"
)

func init() {
	// After the link gate (1000): a twin follows its page, so a page the
	// gate removed gets no twin. The gate lets a page link its own twin.
	registerStep(phaseRender, 1010, "write Markdown twins", writeMarkdownTwins)
}

// The Agent Surface (ADR-0011): every Piece page gets a Markdown twin at the
// page URL with ".md" in place of the trailing slash ("/briefs/x/" ->
// "/briefs/x.md"): the Piece's written Markdown under a generated header,
// for Agents rather than Readers. A twin is written only when its page was,
// so a draft's twin appears in Previews and never in Production, and a
// page a gate dropped has no twin.

// twinURL is the twin's site path for a page URL.
func twinURL(pageURL string) string { return strings.TrimSuffix(pageURL, "/") + ".md" }

// twin is one Piece's Markdown twin before it is written.
type twin struct {
	pageURL string
	title   string
	// lines are the header lines between Canonical and the attribution:
	// kind, dates, author, a Brief's Source.
	lines []string
	body  *markdown.Document
}

func writeMarkdownTwins(s *site) error {
	var twins []twin
	for _, e := range s.publishedExplainers() {
		lines := []string{"Kind: " + kindOf(e.Kind).name + " · Vertical: " + s.verticals[e.Vertical].Name}
		lines = append(lines, twinDates(e.Published, e.Reviewed, e.Verified)...)
		lines = append(lines, "Author: "+s.authors[e.Author].Name)
		twins = append(twins, twin{pageURL: explainerURL(e), title: e.Title, lines: lines, body: e.Body})
	}
	for _, b := range s.publishedBriefs() {
		lines := []string{"Kind: " + briefLabel}
		lines = append(lines, twinDates(b.Published, time.Time{}, b.Verified)...)
		lines = append(lines, "Author: "+s.authors[b.Author].Name, "Source: "+b.SourceURL)
		twins = append(twins, twin{pageURL: briefURL(b), title: b.Title, lines: lines, body: b.Body})
	}
	for _, t := range twins {
		if !s.isWritten(t.pageURL) {
			continue
		}
		if err := s.writeTwin(t); err != nil {
			return err
		}
	}
	return nil
}

// twinDates is the Published line and, when set, Reviewed and Verified,
// each a calendar day in the Piece's own offset.
func twinDates(published, reviewed, verified time.Time) []string {
	lines := []string{"Published: " + published.Format(time.DateOnly)}
	if !reviewed.IsZero() {
		lines = append(lines, "Reviewed: "+reviewed.Format(time.DateOnly))
	}
	if !verified.IsZero() {
		lines = append(lines, "Verified: "+verified.Format(time.DateOnly))
	}
	return lines
}

// writeTwin writes t: the title, the canonical URL, the header lines, the
// attribution line, then the body as the Author wrote it (frontmatter was
// removed at load) with its internal links made absolute.
func (s *site) writeTwin(t twin) error {
	canonical := s.absURL(t.pageURL)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "# %s\n\nCanonical: %s\n", t.title, canonical)
	for _, l := range t.lines {
		buf.WriteString(l)
		buf.WriteByte('\n')
	}
	fmt.Fprintf(&buf, "© HastCode. Quote with a link to %s.\n\n", canonical)
	buf.Write(bytes.TrimSpace(s.absoluteTwinLinks(t.body, t.pageURL)))
	buf.WriteByte('\n')
	return s.out.WriteStatic(twinURL(t.pageURL), buf.Bytes())
}

// absoluteTwinLinks returns body's source with every internal inline link
// or image destination made absolute, so the twin reads correctly fetched
// on its own. Links are located through the parsed tree, never by pattern
// over the text, so "](" inside a code span or fence is never touched: for
// each Link or Image node the destination is the bytes after "](" that
// follow the node's last text segment (closing emphasis or backticks may
// sit between), up to the closing ")" or a title. A destination written
// another way (a reference link, an "<angle>" destination) is left as
// written; the writing guides use inline links.
func (s *site) absoluteTwinLinks(body *markdown.Document, pageURL string) []byte {
	base, err := url.Parse(s.absURL(pageURL))
	if err != nil {
		return body.Source
	}
	src := body.Source
	type edit struct {
		start, stop int
		abs         string
	}
	var edits []edit
	_ = ast.Walk(body.Root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n.(type) {
		case *ast.Link, *ast.Image:
		default:
			return ast.WalkContinue, nil
		}
		// The destination is read from the source, not the node: the
		// transform phase may already have rewritten the node's (a Brief's
		// bare Explainer slug becomes its canonical path).
		stop := lastTextStop(n)
		if stop < 0 || stop > len(src) {
			return ast.WalkContinue, nil
		}
		i := bytes.Index(src[stop:], []byte("]("))
		if i < 0 || strings.Trim(string(src[stop:stop+i]), "`*_~") != "" {
			return ast.WalkContinue, nil
		}
		start := stop + i + 2
		end := start
		for end < len(src) && !strings.ContainsRune(") \t\n", rune(src[end])) {
			end++
		}
		if end == start || src[start] == '<' {
			return ast.WalkContinue, nil
		}
		abs, ok := s.absoluteTwinLink(string(src[start:end]), base)
		if !ok {
			return ast.WalkContinue, nil
		}
		edits = append(edits, edit{start: start, stop: end, abs: abs})
		return ast.WalkContinue, nil
	})
	sort.Slice(edits, func(i, j int) bool { return edits[i].start < edits[j].start })
	var out bytes.Buffer
	pos := 0
	for _, e := range edits {
		if e.start < pos {
			continue
		}
		out.Write(src[pos:e.start])
		out.WriteString(e.abs)
		pos = e.stop
	}
	out.Write(src[pos:])
	return out.Bytes()
}

// absoluteTwinLink is the absolute form of an internal link destination:
// an Explainer link as a Brief writes it (bare slug or /vertical/slug,
// resolveExplainerLink), else a site- or page-relative path resolved
// against the page. External, mailto and fragment-only links are not
// internal.
func (s *site) absoluteTwinLink(dest string, base *url.URL) (string, bool) {
	u, err := url.Parse(dest)
	if err != nil || u.Scheme != "" || u.Host != "" || u.Path == "" {
		return "", false
	}
	if l, ok := s.resolveExplainerLink(dest); ok && l.target != nil {
		return s.absURL(l.href()), true
	}
	return base.ResolveReference(u).String(), true
}

// lastTextStop is the end offset in the source of the last text segment
// under n (its link text or alt text), or -1 when n has none.
func lastTextStop(n ast.Node) int {
	stop := -1
	_ = ast.Walk(n, func(c ast.Node, entering bool) (ast.WalkStatus, error) {
		if t, ok := c.(*ast.Text); ok && entering && t.Segment.Stop > stop {
			stop = t.Segment.Stop
		}
		return ast.WalkContinue, nil
	})
	return stop
}

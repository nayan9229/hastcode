package intake

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// DefaultTemplates is the directory holding the per-kind templates
// (<kind>-template.md) when Intake.Templates is empty: docs/editorial, from
// the repo root where make runs.
const DefaultTemplates = "docs/editorial"

// BriefMaxAge is how old a Lead may be, by its date (published, else first
// seen) against the pick clock, and still be picked as a Brief: 3 days
// (docs/editorial/brief-writing-guide.md, "A fresh Lead").
const BriefMaxAge = 72 * time.Hour

// brief is the kind a Lead becomes when it is picked as a Brief.
const brief = "brief"

// PickOptions override what a pick copies from the Lead's Source.
type PickOptions struct {
	// Kind replaces the Lead's suggested kind when set (one of Kinds).
	Kind string
	// Vertical replaces the Lead's Vertical when set (a site.yaml slug).
	// An Explainer only: a Brief has no Vertical.
	Vertical string
	// Explainers are the slugs of the existing Explainers a Brief links
	// (its `explainers`), at least one; required for a Brief, refused for
	// an Explainer.
	Explainers []string
}

// Pick turns the open Lead id into a draft skeleton under ContentDir and
// marks it Picked; it returns the skeleton's path. `make lead ID=…` runs
// it.
//
// The skeleton is the kind's template (Templates/<kind>-template.md) with
// the Lead filled in and nothing written for the Author: the title from
// the Lead, the kind, the Vertical, a provisional `published` (and, for an
// Explainer, `reviewed`) at the pick time, the Lead's link in Sources (a
// Brief: `source_url`, and its `explainers` linked in "Why it matters"),
// and every other line the template's, its TODO(author) hints included,
// so the build refuses to publish it until the Author has replaced them.
// It is `draft: true` and never sets `verified` or `hero`.
//
// It refuses a Lead that is not open, an unknown kind or Vertical, a slug
// that a Piece already uses, a Brief from a Lead older than BriefMaxAge,
// and a Brief without Explainers or naming one the content tree lacks.
// Nothing is written when it refuses.
func (in *Intake) Pick(ctx context.Context, id string, opts PickOptions) (string, error) {
	id = strings.ToLower(strings.TrimSpace(id))
	if err := in.checkStateDir(); err != nil {
		return "", err
	}
	if in.ContentDir == "" {
		return "", errors.New("intake: no content directory to write the skeleton under")
	}
	st, err := loadState(in.StateDir)
	if err != nil {
		return "", err
	}
	lead, ok := st.Leads[id]
	if !ok || lead.Decision == TooOld {
		return "", fmt.Errorf("intake: no Lead %q in %s (IDs are in the queue, %s)", id, in.StateDir, QueueFile)
	}
	if lead.Decision != Open {
		return "", fmt.Errorf("intake: Lead %s (%q) is already %s; only an open Lead can be picked", id, lead.Title, lead.Decision)
	}
	now := in.now().UTC().Truncate(time.Second)
	p, err := in.planPick(lead, opts, now)
	if err != nil {
		return "", err
	}
	tmpl, err := os.ReadFile(filepath.Join(in.templates(), p.kind+"-template.md"))
	if err != nil {
		return "", fmt.Errorf("intake: reading the %s template: %w", p.kind, err)
	}
	data, err := p.skeleton(string(tmpl))
	if err != nil {
		return "", fmt.Errorf("intake: %s template: %w", p.kind, err)
	}
	file := filepath.Join(in.ContentDir, filepath.FromSlash(p.path))
	if err := writeNewFile(file, data); err != nil {
		return "", err
	}
	lead.Decision = Picked
	st.Leads[id] = lead
	if err := st.save(in.StateDir); err != nil {
		os.Remove(file) // the Lead stays open, so leave no skeleton
		return "", err
	}
	in.logf("picked %s: %s → %s (%s; replace every TODO(author), keep draft: true)", id, lead.Title, file, p.kind)
	return file, nil
}

// pick is what a pick writes: the Lead, its kind, Vertical, slug, the
// Explainers a Brief links and the pick time.
type pick struct {
	lead       Lead
	kind       string
	vertical   string // empty for a Brief
	slug       string
	path       string // content-relative, slash-separated
	explainers []string
	now        time.Time
}

// planPick applies the options to the Lead and checks every refusal that
// needs no template.
func (in *Intake) planPick(lead Lead, opts PickOptions, now time.Time) (pick, error) {
	p := pick{lead: lead, kind: lead.Kind, vertical: lead.Vertical, now: now}
	if opts.Kind != "" {
		p.kind = opts.Kind
	}
	if !slices.Contains(Kinds, p.kind) {
		return p, fmt.Errorf("intake: unknown kind %q (want one of %s)", p.kind, strings.Join(Kinds, ", "))
	}
	p.slug = slugify(lead.Title)
	if p.slug == "" {
		return p, fmt.Errorf("intake: Lead %s's title %q gives no slug", lead.ID, lead.Title)
	}
	if existing, err := in.pieceFile(p.slug); err != nil {
		return p, err
	} else if existing != "" {
		return p, fmt.Errorf("intake: slug %q is already used by %s; a pick never overwrites a Piece", p.slug, existing)
	}

	if p.kind != brief {
		if len(opts.Explainers) > 0 {
			return p, fmt.Errorf("intake: only a Brief names the Explainers it links (EXPLAINER=), not a %s", p.kind)
		}
		if opts.Vertical != "" {
			p.vertical = opts.Vertical
		}
		if !slices.ContainsFunc(in.Verticals, func(v Vertical) bool { return v.Slug == p.vertical }) {
			return p, fmt.Errorf("intake: unknown Vertical %q (want a site.yaml Vertical slug)", p.vertical)
		}
		p.path = "explainers/" + p.vertical + "/" + p.slug + ".md"
		return p, nil
	}

	if opts.Vertical != "" {
		return p, errors.New("intake: a Brief has no Vertical; leave VERTICAL= out")
	}
	p.vertical = ""
	if age := now.Sub(lead.date()); age > BriefMaxAge {
		return p, fmt.Errorf("intake: Lead %s is %s old (dated %s); a Brief needs a Lead at most 3 days (72 hours) old. Pick it as an Explainer (KIND=case-study, concept or comparison) or dismiss it",
			lead.ID, plural(int(age.Hours()), "hour"), lead.date().UTC().Format(time.RFC3339))
	}
	if len(opts.Explainers) == 0 {
		return p, errors.New("intake: a Brief links the Explainer that teaches its lesson: name it (EXPLAINER=slug[,slug]); if none exists yet, the Lead is not a Brief today")
	}
	for _, slug := range opts.Explainers {
		file, err := in.pieceFile(slug)
		if err != nil {
			return p, err
		}
		if !strings.HasPrefix(filepath.ToSlash(file), "explainers/") {
			return p, fmt.Errorf("intake: no Explainer %q in %s; a Brief links an existing Explainer", slug, in.ContentDir)
		}
		if !slices.Contains(p.explainers, slug) {
			p.explainers = append(p.explainers, slug)
		}
	}
	p.path = "briefs/" + p.slug + ".md"
	return p, nil
}

// pieceFile is the content-relative file of the Piece using slug
// (explainers/<vertical>/<slug>.md or briefs/<slug>.md; slugs are unique
// site-wide), or "" when none does.
func (in *Intake) pieceFile(slug string) (string, error) {
	patterns := []string{filepath.Join("explainers", "*", slug+".md"), filepath.Join("briefs", slug+".md")}
	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(in.ContentDir, pattern))
		if err != nil {
			return "", err
		}
		if len(matches) > 0 {
			rel, err := filepath.Rel(in.ContentDir, matches[0])
			return filepath.ToSlash(rel), err
		}
	}
	return "", nil
}

// skeleton fills the template tmpl for p.
func (p pick) skeleton(tmpl string) ([]byte, error) {
	rest, ok := strings.CutPrefix(tmpl, "---\n")
	fm, body, ok2 := strings.Cut(rest, "\n---\n")
	if !ok || !ok2 {
		return nil, errors.New("no frontmatter block")
	}
	frontmatter, err := p.frontmatter(fm)
	if err != nil {
		return nil, err
	}
	if p.kind == brief {
		body, err = p.briefBody(body)
	} else {
		body, err = p.explainerBody(body)
	}
	if err != nil {
		return nil, err
	}
	return []byte("---\n" + frontmatter + "\n---\n" + body), nil
}

// frontmatter is the template's frontmatter with its opening comment
// replaced by the skeleton's own and the Lead's values set; every other
// line, TODO(author) values and commented options included, is kept.
func (p pick) frontmatter(fm string) (string, error) {
	published := quote(p.now.Format(time.RFC3339))
	values := map[string]string{
		"title":     quote(p.lead.Title),
		"published": published,
	}
	if p.kind == brief {
		values["source_url"] = quote(p.lead.Link)
		values["explainers"] = "[" + strings.Join(p.explainers, ", ") + "]"
	} else {
		values["vertical"] = p.vertical
		values["reviewed"] = published
	}
	dates := "published is the pick time, provisional: set it"
	if p.kind != brief {
		dates = "published and reviewed are the pick time, provisional: set them"
	}
	lines := []string{
		fmt.Sprintf("# Skeleton from Lead %s (%s), picked %s.", p.lead.ID, p.lead.Source, p.now.Format(time.DateOnly)),
		fmt.Sprintf("# Follow docs/editorial/%s-writing-guide.md and replace every", p.kind),
		"# TODO(author) line. " + dates,
		"# when you publish. verified: is the Author's alone (ADR-0007).",
	}
	opening := true
	set := map[string]bool{}
	for _, line := range strings.Split(fm, "\n") {
		comment := strings.HasPrefix(line, "#")
		if opening && comment {
			continue // the template's own note on how to use it
		}
		opening = false
		key, _, isKey := strings.Cut(line, ":")
		if comment || !isKey {
			lines = append(lines, line)
			continue
		}
		switch value, known := values[key]; {
		case known:
			line = key + ": " + value
			set[key] = true
		case key == "kind" && strings.TrimSpace(line) != "kind: "+p.kind:
			return "", fmt.Errorf("its kind line %q is not kind: %s", line, p.kind)
		case key == "verified" || key == "hero" || key == "hero_alt":
			return "", fmt.Errorf("it sets %s, which a skeleton never does", key)
		}
		lines = append(lines, line)
	}
	for key := range values {
		if !set[key] {
			return "", fmt.Errorf("it has no %s line to fill", key)
		}
	}
	return strings.Join(lines, "\n"), nil
}

// explainerBody is the template body with ## Sources holding only the
// Lead's link.
func (p pick) explainerBody(body string) (string, error) {
	before, sources, ok := strings.Cut(body, "\n## Sources\n")
	if !ok {
		return "", errors.New("no ## Sources section")
	}
	if strings.Contains(sources, "\n## ") {
		return "", errors.New("## Sources is not the last section")
	}
	lead := fmt.Sprintf("- [%s](%s), %s\n", escapeLinkText(p.lead.Title), p.lead.Link, p.lead.Source)
	return before + "\n## Sources\n\n" + lead, nil
}

// explainerToken is where the Brief template links its Explainer.
const explainerToken = "{explainer}"

// briefBody is the template body with "Why it matters" linking each of
// the Brief's Explainers: the template's link names the first, and a
// TODO(author) line links each other one.
func (p pick) briefBody(body string) (string, error) {
	if strings.Count(body, explainerToken) != 1 {
		return "", fmt.Errorf("want one %s link in the body", explainerToken)
	}
	body = strings.Replace(body, explainerToken, p.explainers[0], 1)
	var extra strings.Builder
	for _, slug := range p.explainers[1:] {
		fmt.Fprintf(&extra, "\nTODO(author): link [the idea](%s) where the analysis uses it.\n", slug)
	}
	return strings.TrimRight(body, "\n") + "\n" + extra.String(), nil
}

// writeNewFile creates file with data, refusing one that exists.
func writeNewFile(file string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if errors.Is(err, fs.ErrExist) {
		return fmt.Errorf("intake: %s exists; a pick never overwrites a Piece", file)
	}
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		os.Remove(file)
		return err
	}
	return f.Close()
}

// slugify makes a slug from a title: lower-case ASCII letters and digits,
// every other run of characters one hyphen, at most 80 characters, cut at
// a hyphen.
func slugify(title string) string {
	var b strings.Builder
	hyphen := false
	for _, r := range strings.ToLower(title) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			if hyphen && b.Len() > 0 {
				b.WriteByte('-')
			}
			hyphen = false
			b.WriteRune(r)
		case r == '\'' || r == '’':
			// "Netflix's" → "netflixs", not "netflix-s"
		default:
			hyphen = true
		}
	}
	slug := b.String()
	const max = 80
	if len(slug) > max {
		slug = slug[:max]
		if i := strings.LastIndexByte(slug, '-'); i > 0 {
			slug = slug[:i]
		}
	}
	return slug
}

// quote is s as a YAML double-quoted scalar (JSON's string syntax is a
// subset of it).
func quote(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (in *Intake) templates() string {
	if in.Templates == "" {
		return DefaultTemplates
	}
	return in.Templates
}

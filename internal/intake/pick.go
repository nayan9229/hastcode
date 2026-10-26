package intake

import (
	"context"
	"errors"
	"fmt"
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
		if !isSlug(slug) {
			return p, fmt.Errorf("intake: EXPLAINER=%q is not a slug; name each Explainer by its exact slug (lower-case letters and digits, single hyphens between), no wildcards or paths", slug)
		}
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

func (in *Intake) templates() string {
	if in.Templates == "" {
		return DefaultTemplates
	}
	return in.Templates
}

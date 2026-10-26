package intake

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// QueueFile is the queue's file name in the state directory: the Markdown
// the Author reads each morning.
const QueueFile = "queue.md"

// queue is the open Leads grouped by Vertical, each group newest first,
// rendered as the queue file. It is rebuilt from the state on every run,
// never edited in place.
type queue struct {
	now    time.Time
	groups []queueGroup
	isNew  map[string]bool // Lead IDs this run listed first
	// covered are the normalised links Pieces already cite (covered.go).
	covered map[string]bool
}

type queueGroup struct {
	heading string // the Vertical's name
	leads   []Lead
}

// queue builds the queue of st's open Leads: groups in Vertical order (a
// Vertical no longer in site.yaml last, under its slug), newest first,
// ties by title then ID so the file is stable.
func (in *Intake) queue(st *state, isNew map[string]bool, now time.Time) *queue {
	byVertical := map[string][]Lead{}
	for _, l := range st.Leads {
		if l.Decision == Open {
			byVertical[l.Vertical] = append(byVertical[l.Vertical], l)
		}
	}
	q := &queue{now: now, isNew: isNew}
	add := func(heading string, leads []Lead) {
		if len(leads) == 0 {
			return
		}
		slices.SortFunc(leads, func(a, b Lead) int {
			if c := b.date().Compare(a.date()); c != 0 {
				return c
			}
			return strings.Compare(a.Title+a.ID, b.Title+b.ID)
		})
		q.groups = append(q.groups, queueGroup{heading, leads})
	}
	for _, v := range in.Verticals {
		add(v.Name, byVertical[v.Slug])
		delete(byVertical, v.Slug)
	}
	var unknown []string
	for slug := range byVertical {
		unknown = append(unknown, slug)
	}
	slices.Sort(unknown)
	for _, slug := range unknown {
		add(slug, byVertical[slug])
	}
	return q
}

// leads are the queue's Leads in order.
func (q *queue) leads() []Lead {
	var all []Lead
	for _, g := range q.groups {
		all = append(all, g.leads...)
	}
	return all
}

// write renders the queue into dir and returns the file's path.
func (q *queue) write(dir string) (string, error) {
	file := filepath.Join(dir, QueueFile)
	return file, writeFileAtomic(file, []byte(q.render()))
}

// render is the queue file. One line per Lead, under its Vertical:
//
//   - `id` [Title](link) · Source · kind · 2026-10-24 (1 day old) · new · covered
func (q *queue) render() string {
	var b strings.Builder
	all := q.leads()
	newCount := 0
	for _, l := range all {
		if q.isNew[l.ID] {
			newCount++
		}
	}
	fmt.Fprintf(&b, "# Leads\n\n")
	fmt.Fprintf(&b, "Intake %s: %s open, %d new.\n", q.now.UTC().Format("2006-01-02 15:04 UTC"), plural(len(all), "Lead"), newCount)
	fmt.Fprintf(&b, "Pick one with `make lead ID=<id>`; dismiss one with `make lead-dismiss ID=<id>`.\n")
	for _, g := range q.groups {
		fmt.Fprintf(&b, "\n## %s\n\n", g.heading)
		for _, l := range g.leads {
			fmt.Fprintf(&b, "- `%s` [%s](%s) · %s · %s · %s", l.ID, escapeLinkText(l.Title), l.Link, l.Source, l.Kind, q.dateAndAge(l))
			if q.isNew[l.ID] {
				b.WriteString(" · new")
			}
			if q.covered[normaliseLink(l.Link)] {
				b.WriteString(" · covered")
			}
			b.WriteString("\n")
		}
	}
	return b.String()
}

// dateAndAge is "2026-10-24 (1 day old)", by the publication date or, when
// the Source gives none, "first seen 2026-10-24 (…)".
func (q *queue) dateAndAge(l Lead) string {
	days := int(q.now.Sub(l.date()) / (24 * time.Hour))
	age := "today"
	if days > 0 {
		age = plural(days, "day") + " old"
	}
	date := l.date().UTC().Format("2006-01-02")
	if l.Published.IsZero() {
		date = "first seen " + date
	}
	return fmt.Sprintf("%s (%s)", date, age)
}

func plural(n int, noun string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", noun)
	}
	return fmt.Sprintf("%d %ss", n, noun)
}

// escapeLinkText keeps a title's brackets from ending the link text.
func escapeLinkText(s string) string {
	return strings.NewReplacer(`\`, `\\`, `[`, `\[`, `]`, `\]`).Replace(s)
}

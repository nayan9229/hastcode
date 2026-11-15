// Package release is the Author's release step and the Gap Day arithmetic
// behind it (ADR-0012): which Publication Days since the cap started are
// short of two Pieces, re-dating a draft to such a day, and releasing
// drafts while keeping the Publication Day each one carries. It edits
// frontmatter lines in place and knows nothing about rendering; the
// builder's gates remain the final check.
package release

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"hastcode.com/site/internal/content"
)

// MaxPerDay is the Publish cap (ADR-0010), Explainers and Briefs together.
// The build enforces the same number as builder.maxPublishedPerDay
// (step_validate_publish_pace.go); change both together.
const MaxPerDay = 2

// CapStart is the first Publication Day the cap applies to (ADR-0010), the
// build's builder.paceStart; change both together.
var CapStart = time.Date(2026, 10, 30, 0, 0, 0, 0, time.UTC)

// Piece is one content file's release-relevant frontmatter.
type Piece struct {
	Path      string // relative to the content directory, slash-separated
	Title     string
	Published time.Time
	Draft     bool
	Explainer bool // under explainers/ (carries kind and reviewed); else a Brief
}

// Day is the Piece's Publication Day: the calendar day of published in its
// own offset, as YYYY-MM-DD.
func (p Piece) Day() string { return p.Published.Format(time.DateOnly) }

// Gap is a Gap Day: a Publication Day short of the cap, and by how much.
type Gap struct {
	Day   string
	Slots int
}

// GapDays lists the Gap Days inside the look-back: the lookBack days before
// now's own calendar day (never today, never before CapStart), oldest
// first, each with the Pieces, live or draft, that already carry it
// counted against the cap. now's day is read in now's own location, the
// same clock make lead stamps a skeleton with. An empty list means nothing
// to fill.
func GapDays(contentDir string, now time.Time, lookBack int) ([]Gap, error) {
	pieces, err := load(contentDir)
	if err != nil {
		return nil, err
	}
	count := map[string]int{}
	for _, p := range pieces {
		count[p.Day()]++
	}
	y, m, d := now.Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	var gaps []Gap
	for i := lookBack; i >= 1; i-- {
		day := today.AddDate(0, 0, -i)
		if day.Before(CapStart) {
			continue
		}
		key := day.Format(time.DateOnly)
		if slots := MaxPerDay - count[key]; slots > 0 {
			gaps = append(gaps, Gap{Day: key, Slots: slots})
		}
	}
	return gaps, nil
}

// Redate moves a draft's Publication Day to day (YYYY-MM-DD), keeping the
// clock time and offset of its published (and, on an Explainer, reviewed).
// It refuses a Brief: a Brief is never backfilled (ADR-0012).
func Redate(contentDir, file, day string) error {
	if _, err := time.Parse(time.DateOnly, day); err != nil {
		return fmt.Errorf("%q is not a day (YYYY-MM-DD)", day)
	}
	rel, err := contentRelative(contentDir, file)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(rel, "explainers/") {
		return fmt.Errorf("%s: only an Explainer may be re-dated to a Gap Day, never a Brief (ADR-0012)", rel)
	}
	path := filepath.Join(contentDir, filepath.FromSlash(rel))
	return editFrontmatter(path, func(lines []string) []string {
		for i, l := range lines {
			if key, value := keyValue(l); (key == "published" || key == "reviewed") && len(value) > len(day) {
				lines[i] = key + ": " + fmt.Sprintf("%q", day+value[len(day):])
			}
		}
		return lines
	})
}

// load reads every Explainer and Brief under contentDir (explainers/*/*.md
// and briefs/*.md), sorted by path.
func load(contentDir string) ([]Piece, error) {
	var paths []string
	for _, pattern := range []string{"explainers/*/*.md", "briefs/*.md"} {
		m, err := filepath.Glob(filepath.Join(contentDir, filepath.FromSlash(pattern)))
		if err != nil {
			return nil, err
		}
		paths = append(paths, m...)
	}
	sort.Strings(paths)
	var pieces []Piece
	for _, p := range paths {
		rel, err := filepath.Rel(contentDir, p)
		if err != nil {
			return nil, err
		}
		piece, err := read(p, filepath.ToSlash(rel))
		if err != nil {
			return nil, err
		}
		pieces = append(pieces, piece)
	}
	return pieces, nil
}

var keyLine = regexp.MustCompile(`^([a-z_]+):\s*(.*?)\s*$`)

// keyValue is a frontmatter line's top-level key and unquoted value, or
// "" for a line that is not one (a comment, a list item).
func keyValue(line string) (key, value string) {
	m := keyLine.FindStringSubmatch(line)
	if m == nil {
		return "", ""
	}
	return m[1], strings.Trim(m[2], `"`)
}

// frontLines splits a file into its frontmatter lines and the body after
// the closing marker.
func frontLines(data []byte) (lines []string, body []byte, err error) {
	front, body, err := content.SplitFrontmatter(data)
	if err != nil {
		return nil, nil, err
	}
	return strings.Split(strings.TrimSuffix(string(front), "\n"), "\n"), body, nil
}

// editFrontmatter rewrites path's frontmatter lines through edit, keeping
// the body byte for byte.
func editFrontmatter(path string, edit func([]string) []string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines, body, err := frontLines(data)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	out := "---\n" + strings.Join(edit(lines), "\n") + "\n---\n" + string(body)
	return os.WriteFile(path, []byte(out), 0o644)
}

func read(path, rel string) (Piece, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Piece{}, err
	}
	lines, _, err := frontLines(data)
	if err != nil {
		return Piece{}, fmt.Errorf("%s: %w", rel, err)
	}
	values := map[string]string{}
	for _, l := range lines {
		if key, value := keyValue(l); key != "" {
			values[key] = value
		}
	}
	published, err := time.Parse(time.RFC3339, values["published"])
	if err != nil {
		return Piece{}, fmt.Errorf("%s: published %q is not RFC 3339 with timezone", rel, values["published"])
	}
	return Piece{
		Path:      rel,
		Title:     values["title"],
		Published: published,
		Draft:     values["draft"] == "true",
		Explainer: strings.HasPrefix(rel, "explainers/"),
	}, nil
}

// A Plan is a release the Author asked for, checked against the cap and
// not yet applied.
type Plan struct {
	contentDir string
	now        time.Time
	// Releases are the drafts to release, by path.
	Releases []Piece
}

// Verified is the verified day a released Piece gets: the plan's clock in
// the Piece's own offset, the rule the build's verified gate reads it by.
func (p *Plan) Verified(piece Piece) string {
	return p.now.In(piece.Published.Location()).Format(time.DateOnly)
}

// NewPlan plans a release of the drafts named by files (paths relative to
// the content directory, the repo root or absolute; every one must be a
// draft), or of every draft when files is empty. It refuses, naming the
// files, when a Publication Day since CapStart would hold more than
// MaxPerDay Pieces live plus released.
func NewPlan(contentDir string, files []string, now time.Time) (*Plan, error) {
	pieces, err := load(contentDir)
	if err != nil {
		return nil, err
	}
	byPath := map[string]Piece{}
	for _, p := range pieces {
		byPath[p.Path] = p
	}
	var releases []Piece
	if len(files) == 0 {
		for _, p := range pieces {
			if p.Draft {
				releases = append(releases, p)
			}
		}
	} else {
		for _, f := range files {
			rel, err := contentRelative(contentDir, f)
			if err != nil {
				return nil, err
			}
			p, ok := byPath[rel]
			switch {
			case !ok:
				return nil, fmt.Errorf("%s: not an Explainer or Brief under %s", f, contentDir)
			case !p.Draft:
				return nil, fmt.Errorf("%s: not a draft (no draft: true), nothing to release", p.Path)
			}
			releases = append(releases, p)
		}
	}
	if len(releases) == 0 {
		return nil, errors.New("no drafts to release")
	}
	count := map[string]int{}
	for _, p := range pieces {
		if !p.Draft {
			count[p.Day()]++
		}
	}
	onDay := map[string][]string{}
	for _, p := range releases {
		count[p.Day()]++
		onDay[p.Day()] = append(onDay[p.Day()], p.Path)
	}
	var over []string
	for day, n := range count {
		if n > MaxPerDay && !dayBefore(day, CapStart) && len(onDay[day]) > 0 {
			over = append(over, fmt.Sprintf("%s would hold %d Pieces (releasing %s)", day, n, strings.Join(onDay[day], ", ")))
		}
	}
	if len(over) > 0 {
		sort.Strings(over)
		return nil, fmt.Errorf("over the cap of %d a Publication Day (ADR-0010): %s. Release fewer with FILES=, or delete a stale draft", MaxPerDay, strings.Join(over, "; "))
	}
	return &Plan{contentDir: contentDir, now: now, Releases: releases}, nil
}

func dayBefore(day string, t time.Time) bool {
	d, err := time.Parse(time.DateOnly, day)
	return err == nil && d.Before(t)
}

// contentRelative turns a file argument into a path relative to contentDir.
func contentRelative(contentDir, f string) (string, error) {
	if filepath.IsAbs(f) {
		abs, err := filepath.Abs(contentDir)
		if err != nil {
			return "", err
		}
		rel, err := filepath.Rel(abs, f)
		if err != nil {
			return "", err
		}
		return filepath.ToSlash(rel), nil
	}
	f = filepath.ToSlash(filepath.Clean(f))
	prefix := filepath.ToSlash(filepath.Clean(contentDir)) + "/"
	return strings.TrimPrefix(f, prefix), nil
}

// Apply releases every Piece in the plan: removes draft: true, sets
// verified to Verified (one line, right after published, replacing any
// earlier value), sets reviewed to now on an Explainer, and keeps every
// other line, published included, exactly as written.
func (p *Plan) Apply() error {
	for _, piece := range p.Releases {
		verified := "verified: " + p.Verified(piece)
		reviewed := "reviewed: " + fmt.Sprintf("%q", p.now.In(piece.Published.Location()).Format(time.RFC3339))
		err := editFrontmatter(filepath.Join(p.contentDir, filepath.FromSlash(piece.Path)), func(in []string) []string {
			var out []string
			hadReviewed := false
			for _, l := range in {
				key, _ := keyValue(l)
				switch key {
				case "draft", "verified":
					continue // dropped; verified is re-added after published
				case "reviewed":
					if piece.Explainer {
						out = append(out, reviewed)
						hadReviewed = true
						continue
					}
				}
				out = append(out, l)
				if key == "published" {
					out = append(out, verified)
				}
			}
			if piece.Explainer && !hadReviewed {
				out = append(out, reviewed)
			}
			return out
		})
		if err != nil {
			return err
		}
	}
	return nil
}

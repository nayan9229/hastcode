// Package intake is the Intake (GLOSSARY "Sourcing"; ADR-0007): the run
// that checks the Source Watchlist and lists the new Leads for the Author.
// It lists Leads; it never writes prose.
//
// What writes where: Run and Dismiss write only under StateDir (the
// seen-state and the queue) and only read the content tree; they never
// write outside StateDir. Pick is the only entry point that
// writes outside StateDir, and only under ContentDir: one draft skeleton
// per pick.
//
// Entry points, one file each, wrapped by cmd/intake:
//
//	Run      intake.go     make intake                validate the Watchlist, fetch every Source, record new Leads, rewrite the queue
//	Dismiss  dismiss.go    make lead-dismiss ID=…     hide a Lead for good
//	Pick     pick.go       make lead ID=…             turn a Lead into a draft skeleton from its kind's template
//
// The rest, also one concern per file:
//
//	verticals.go   Vertical and LoadVerticals (site.yaml's Verticals in display order)
//	watchlist.go   the Watchlist file: Source, LoadWatchlist and its validation
//	fetch.go       the per-Source-type fetcher registry (registerFetcher) and Entry
//	fetch_feed.go  feed Sources: RSS 2.0 and Atom, the format read from the root element
//	fetch_github.go  repository Sources: releases through the GitHub REST API (internal/github), never github.com pages
//	http.go        the HTTP GET a feed fetch makes
//	lead.go        Lead and its ID
//	state.go       the seen-state file (state.json): every entry ever seen, by Lead ID, and per-Source state
//	source_state.go   SourceState: a Source's validators (ETag, Last-Modified) and newest entry date
//	source_health.go  failed and silent Sources in the Report and the run summary
//	expire.go      expiry: open Leads undecided for ExpireAfter (7 days) become Expired
//	covered.go     covered Leads: links a Piece already cites, read from the content tree, normalised
//	queue.go       the queue file (queue.md) the Author reads, regenerated from state each run
//	summary.go     the run summary printed after a run
package intake

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Intake holds what every entry point needs. Tests build one directly;
// cmd/intake fills it from flags and site.yaml.
type Intake struct {
	// Watchlist is the Source Watchlist file (sources.yaml at the repo
	// root), loaded and validated at the start of every run.
	Watchlist string
	// Verticals are the known Verticals in display order: a Source must
	// name one, and the queue groups Leads in this order.
	Verticals []Vertical
	// StateDir is the local, gitignored directory holding the seen-state
	// (StateFile) and the queue (QueueFile). Run and Dismiss write nowhere
	// else.
	StateDir string
	// ContentDir is the content tree. Run reads it (covered Leads); Run
	// and Dismiss refuse a StateDir inside it, so neither writes under it.
	// Pick is its only writer: one draft skeleton per pick.
	ContentDir string
	// Templates is the directory of the per-kind templates Pick fills
	// (<kind>-template.md); default DefaultTemplates.
	Templates string
	// HTTP makes the feed requests; default http.DefaultClient. When set,
	// the GitHub requests use it too.
	HTTP *http.Client
	// Timeout bounds each request (http.go); default DefaultTimeout.
	Timeout time.Duration
	// MaxAge is the age cutoff: a run lists only entries published at most
	// this long before its clock; default DefaultMaxAge. An older entry is
	// recorded as seen (TooOld) and never listed, unless a later run with
	// a wider MaxAge reaches it. An entry with no date is listed once.
	MaxAge time.Duration
	// GitHubURL is the GitHub REST API repository Sources are fetched
	// from (fetch_github.go); default github.APIURL. Tests point it at a
	// fake.
	GitHubURL string
	// Getenv reads the GitHub token (github.EnvToken), needed only when
	// the Watchlist has a repository Source; default os.Getenv.
	Getenv func(string) string
	// Now is the run's clock (Lead ages, first-seen dates); default
	// time.Now.
	Now func() time.Time
	// Log receives progress lines (one per Source fetched or failed, one
	// per Fact Pack skipped). Default: discarded.
	Log io.Writer
}

// Report is what one Run found.
type Report struct {
	// New are the Leads listed for the first time by this run, in queue
	// order (by Vertical, newest first).
	New []Lead
	// Open is the number of open Leads in the queue, new or not.
	Open int
	// QueueFile is the queue's path.
	QueueFile string
	// Failed are the Sources whose fetch failed this run (network, a
	// status other than 200 or 304, an unparseable body), with why; the
	// run went on without them (source_health.go).
	Failed []FailedSource
	// Silent are the Sources fetched this run whose newest item is more
	// than SilentAfter old (source_health.go).
	Silent []SilentSource
	// Expired is how many open Leads this run expired (ExpireAfter).
	Expired int
	// NewByVertical counts New per Vertical, every Vertical in display
	// order, zeros included, so a starving Vertical shows.
	NewByVertical []Count
	// NewByKind counts New per suggested kind, every kind (Kinds).
	NewByKind []Count
}

// Count is how many new Leads one Vertical or kind got.
type Count struct {
	Name string
	New  int
}

// Run validates the Watchlist, fetches every Source in Watchlist order,
// records each entry not seen before as an open Lead (an entry published
// before the age cutoff as TooOld: seen, not listed), saves the state and
// rewrites the queue. A Watchlist error stops the run before any request.
// Requests go one at a time; a Source that fails (network, timeout, a
// status other than 200 or 304, unparseable) is reported in Report.Failed
// and the run goes on.
func (in *Intake) Run(ctx context.Context) (*Report, error) {
	wl, err := LoadWatchlist(in.Watchlist, in.Verticals)
	if err != nil {
		return nil, err
	}
	if err := in.checkStateDir(); err != nil {
		return nil, err
	}
	st, err := loadState(in.StateDir)
	if err != nil {
		return nil, err
	}
	now := in.now()
	cutoff := now.Add(-in.maxAge())
	report := &Report{}
	newIDs := map[string]bool{}
	for _, src := range wl.Sources {
		prev := st.Sources[src.Name]
		fetched, err := fetchers[src.Type()].fetch(ctx, in, src, prev)
		if err != nil {
			if ctx.Err() != nil { // interrupted: stop, not one Source's fault
				return nil, ctx.Err()
			}
			report.fail(in.logf, src, err) // the Source keeps its state; the run goes on
			continue
		}
		next := fetched.State.withNewest(prev.Newest, fetched.Entries)
		st.Sources[src.Name] = next
		report.checkSilent(src, next.Newest, now)
		added := 0
		for _, entry := range fetched.Entries {
			lead := newLead(src, entry, now)
			old := !entry.Published.IsZero() && entry.Published.Before(cutoff)
			if prev, seen := st.Leads[lead.ID]; seen && (prev.Decision != TooOld || old) {
				continue // listed before, or still too old
			}
			if old {
				lead.Decision = TooOld // seen, never listed
				st.Leads[lead.ID] = lead
				continue
			}
			st.Leads[lead.ID] = lead
			newIDs[lead.ID] = true
			added++
		}
		in.logf("fetched %s: %d entries, %d new", src.Name, len(fetched.Entries), added)
	}
	report.Expired = st.expire(now)
	if err := st.save(in.StateDir); err != nil {
		return nil, err
	}
	covered, err := in.coveredLinks()
	if err != nil {
		return nil, err
	}
	q := in.queue(st, newIDs, covered, now)
	open := q.leads()
	for _, l := range open {
		if newIDs[l.ID] {
			report.New = append(report.New, l)
		}
	}
	report.Open = len(open)
	report.NewByVertical, report.NewByKind = in.countNew(report.New)
	report.QueueFile, err = q.write(in.StateDir)
	if err != nil {
		return nil, err
	}
	return report, nil
}

// checkStateDir refuses a state directory inside the content tree.
func (in *Intake) checkStateDir() error {
	if in.StateDir == "" {
		return fmt.Errorf("intake: no state directory")
	}
	if in.ContentDir == "" {
		return nil
	}
	state, err1 := filepath.Abs(in.StateDir)
	content, err2 := filepath.Abs(in.ContentDir)
	if err1 != nil || err2 != nil {
		return fmt.Errorf("intake: resolving %s and %s", in.StateDir, in.ContentDir)
	}
	if rel, err := filepath.Rel(content, state); err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("intake: state directory %s is inside the content tree %s; Intake never writes under it", in.StateDir, in.ContentDir)
	}
	return nil
}

// DefaultMaxAge is the age cutoff when Intake.MaxAge is zero: two weeks,
// so a first run over the whole Watchlist lists what is current, not every
// Source's back catalogue.
const DefaultMaxAge = 14 * 24 * time.Hour

// ParseMaxAge reads an age cutoff as the command takes it: whole days
// ("30d") or a Go duration ("720h").
func ParseMaxAge(s string) (time.Duration, error) {
	if days, ok := strings.CutSuffix(s, "d"); ok {
		n, err := strconv.Atoi(days)
		if err != nil || n <= 0 {
			return 0, fmt.Errorf("max age %q: want whole days (30d) or a duration (720h)", s)
		}
		return time.Duration(n) * 24 * time.Hour, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		return 0, fmt.Errorf("max age %q: want whole days (30d) or a duration (720h)", s)
	}
	return d, nil
}

func (in *Intake) maxAge() time.Duration {
	if in.MaxAge <= 0 {
		return DefaultMaxAge
	}
	return in.MaxAge
}

func (in *Intake) now() time.Time {
	if in.Now == nil {
		return time.Now()
	}
	return in.Now()
}

func (in *Intake) httpClient() *http.Client {
	if in.HTTP == nil {
		return http.DefaultClient
	}
	return in.HTTP
}

func (in *Intake) logf(format string, args ...any) {
	if in.Log != nil {
		fmt.Fprintf(in.Log, format+"\n", args...)
	}
}

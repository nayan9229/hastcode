// Package intake is the Intake (GLOSSARY "Sourcing"; ADR-0007): the run
// that checks the Source Watchlist and lists the new Leads for the Author.
// It lists Leads; it never writes a Piece, and it writes nothing outside
// its local state directory (never under the content tree).
//
// Entry points, one file each, wrapped by cmd/intake:
//
//	Run      intake.go     make intake                validate the Watchlist, fetch every Source, record new Leads, rewrite the queue
//	Dismiss  dismiss.go    make lead-dismiss ID=…     hide a Lead for good (ticket 05)
//	Pick     pick.go       make lead ID=…             turn a Lead into a draft skeleton (ticket 09)
//
// The rest, also one concern per file:
//
//	watchlist.go   the Watchlist file: Source, LoadWatchlist and its validation
//	fetch.go       the per-Source-type fetcher registry (registerFetcher) and Item
//	fetch_feed.go  feed Sources: RSS 2.0 and Atom, the format read from the root element
//	fetch_github.go  repository Sources: releases through the GitHub REST API (internal/github), never github.com pages
//	http.go        the HTTP GET a feed fetch makes
//	lead.go        Lead and its ID
//	state.go       the seen-state file (state.json): every Lead ever listed and per-Source state
//	source_state.go   SourceState: a Source's validators (ETag, Last-Modified) and newest item date
//	source_health.go  failed and silent Sources in the Report and the run summary
//	queue.go       the queue file (queue.md) the Author reads, regenerated from state each run
//	summary.go     the run summary printed after a run
package intake

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// Vertical is a Vertical Intake files Leads under (site.yaml "verticals").
type Vertical struct {
	Slug string
	Name string
}

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
	// (StateFile) and the queue (QueueFile). Intake writes nowhere else.
	StateDir string
	// ContentDir is the content tree. Run refuses a StateDir inside it,
	// so Intake can never write under it.
	ContentDir string
	// HTTP makes the feed requests; default http.DefaultClient. When set,
	// the GitHub requests use it too.
	HTTP *http.Client
	// Timeout bounds each request (http.go); default DefaultTimeout.
	Timeout time.Duration
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
	// Log receives progress lines (one per Source fetched or skipped).
	// Default: discarded.
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
	// Skipped are the Sources this run did not fetch, with why (a Source
	// type with no registered fetcher).
	Skipped []SkippedSource
	// Failed are the Sources whose fetch failed this run (network, a
	// status other than 200 or 304, an unparseable body), with why; the
	// run went on without them (source_health.go).
	Failed []FailedSource
	// Silent are the Sources fetched this run whose newest item is more
	// than SilentAfter old (source_health.go).
	Silent []SilentSource
}

// SkippedSource is a Source a run did not fetch.
type SkippedSource struct {
	Source string
	Reason string
}

// Run validates the Watchlist, fetches every Source in Watchlist order,
// records each item not seen before as an open Lead, saves the state and
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
	report := &Report{}
	newIDs := map[string]bool{}
	for _, src := range wl.Sources {
		f, ok := fetchers[src.Type()]
		if !ok {
			report.Skipped = append(report.Skipped, SkippedSource{src.Name, fmt.Sprintf("%s Sources are not fetched yet", src.Type())})
			continue
		}
		fetched, err := f.fetch(ctx, in, src, st.Sources[src.Name])
		if err != nil {
			if ctx.Err() != nil { // interrupted: stop, not one Source's fault
				return nil, ctx.Err()
			}
			report.fail(in, src, err) // the Source keeps its state; the run goes on
			continue
		}
		st.Sources[src.Name] = report.fetchedState(src, st.Sources[src.Name], fetched, now)
		added := 0
		for _, item := range fetched.Items {
			lead := newLead(src, item, now)
			if _, seen := st.Leads[lead.ID]; seen {
				continue
			}
			st.Leads[lead.ID] = lead
			newIDs[lead.ID] = true
			added++
		}
		in.logf("fetched %s: %d items, %d new", src.Name, len(fetched.Items), added)
	}
	if err := st.save(in.StateDir); err != nil {
		return nil, err
	}
	q := in.queue(st, newIDs, now)
	for _, l := range q.leads() {
		if newIDs[l.ID] {
			report.New = append(report.New, l)
		}
	}
	report.Open = len(q.leads())
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

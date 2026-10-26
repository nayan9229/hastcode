// Command intake is the Intake (GLOSSARY "Sourcing"; ADR-0007): it checks
// the Source Watchlist and lists the new Leads in a local queue. `make
// intake`, `make lead-dismiss ID=…` and `make lead ID=…` run it.
//
//	go run ./cmd/intake [-max-age 30d]          run: list new Leads in .intake/queue.md
//	go run ./cmd/intake -dismiss <id>           hide a Lead for good
//	go run ./cmd/intake -pick <id> [-kind k] [-vertical v] [-explainer slug[,slug]]
//	                                            start a draft skeleton from a Lead (a Brief names its Explainers)
//
// The Verticals come from -config (site.yaml). The queue and the
// seen-state live in -state, a local gitignored directory; a run and a
// dismiss write nothing else, and nothing under -content. A pick writes
// one draft skeleton under -content, from the kind's template in
// -templates. Repository
// Sources read GITHUB_TOKEN, as the fetch command does; a feed-only
// Watchlist needs none. All logic is in internal/intake; main is flags and
// exit codes:
//
//	0  done, even when Sources failed (each is named in the summary)
//	1  the run failed: an invalid Watchlist, the state, I/O, interrupted
//	2  usage: bad flags (an empty -dismiss or -pick ID included), unreadable site config
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"hastcode.com/site/internal/intake"
)

const (
	exitOK     = 0
	exitFailed = 1
	exitUsage  = 2
)

// options are the command's flags.
type options struct {
	watchlist, config, contentDir, stateDir, templates string
	// dismiss and pick are the Lead IDs given with -dismiss and -pick;
	// set records which of the two flags appeared at all, so an empty ID
	// (`make lead-dismiss` without ID=) is a usage error, never a run.
	dismiss, pick string
	set           map[string]bool
	pickOptions   intake.PickOptions
	// explainers is -explainer, comma-separated slugs.
	explainers string
	// maxAge is -max-age: empty for the default (intake.DefaultMaxAge).
	maxAge string
}

func main() {
	var o options
	flag.StringVar(&o.watchlist, "watchlist", intake.DefaultWatchlist, "the Source Watchlist file")
	flag.StringVar(&o.config, "config", "site.yaml", "site config; its verticals are the ones a Source may name")
	flag.StringVar(&o.contentDir, "content", "content", "content directory; a run and a dismiss never write under it, a pick writes its skeleton there")
	flag.StringVar(&o.templates, "templates", intake.DefaultTemplates, "with -pick: the directory of the per-kind templates (<kind>-template.md)")
	flag.StringVar(&o.stateDir, "state", ".intake", "local state directory (gitignored): the queue and the seen-state")
	flag.StringVar(&o.dismiss, "dismiss", "", "dismiss the Lead with this ID for good")
	flag.StringVar(&o.pick, "pick", "", "pick the Lead with this ID: write its draft skeleton")
	flag.StringVar(&o.pickOptions.Kind, "kind", "", "with -pick: override the Lead's suggested kind")
	flag.StringVar(&o.pickOptions.Vertical, "vertical", "", "with -pick: override the Lead's Vertical")
	flag.StringVar(&o.explainers, "explainer", "", "with -pick of a Brief: the slug(s) of the Explainer(s) it links, comma-separated")
	flag.StringVar(&o.maxAge, "max-age", "", "list only entries published this recently: whole days (30d) or a duration; default 14d")
	flag.Parse()
	o.set = map[string]bool{}
	flag.Visit(func(f *flag.Flag) { o.set[f.Name] = true })

	os.Exit(run(o))
}

func run(o options) int {
	switch {
	case flag.NArg() > 0:
		return usage(fmt.Sprintf("unexpected argument %q", flag.Arg(0)))
	case o.set["dismiss"] && o.set["pick"]:
		return usage("-dismiss and -pick are exclusive")
	case o.set["dismiss"] && o.dismiss == "":
		return usage("-dismiss needs a Lead ID (make lead-dismiss ID=<id>)")
	case o.set["pick"] && o.pick == "":
		return usage("-pick needs a Lead ID (make lead ID=<id>)")
	case !o.set["pick"] && (o.pickOptions.Kind != "" || o.pickOptions.Vertical != "" || o.explainers != ""):
		return usage("-kind, -vertical and -explainer go with -pick")
	case (o.set["dismiss"] || o.set["pick"]) && o.maxAge != "":
		return usage("-max-age goes with a run, not -dismiss or -pick")
	}
	var maxAge time.Duration
	if o.maxAge != "" {
		var err error
		if maxAge, err = intake.ParseMaxAge(o.maxAge); err != nil {
			return usage(err.Error())
		}
	}
	for _, slug := range strings.Split(o.explainers, ",") {
		if slug = strings.TrimSpace(slug); slug != "" {
			o.pickOptions.Explainers = append(o.pickOptions.Explainers, slug)
		}
	}
	verticals, err := intake.LoadVerticals(o.config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return exitUsage
	}
	in := &intake.Intake{
		Watchlist:  o.watchlist,
		Verticals:  verticals,
		StateDir:   o.stateDir,
		ContentDir: o.contentDir,
		Templates:  o.templates,
		MaxAge:     maxAge,
		Log:        os.Stdout,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	switch {
	case o.set["dismiss"]:
		err = in.Dismiss(ctx, o.dismiss)
	case o.set["pick"]:
		var file string
		if file, err = in.Pick(ctx, o.pick, o.pickOptions); err == nil {
			fmt.Println("wrote", file)
		}
	default:
		var report *intake.Report
		if report, err = in.Run(ctx); err == nil {
			report.WriteSummary(os.Stdout)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return exitFailed
	}
	return exitOK
}

func usage(msg string) int {
	fmt.Fprintln(os.Stderr, "error:", msg)
	flag.Usage()
	return exitUsage
}

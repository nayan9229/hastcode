// Command intake is the Intake (GLOSSARY "Sourcing"; ADR-0007): it checks
// the Source Watchlist and lists the new Leads in a local queue. `make
// intake`, `make lead-dismiss ID=…` and `make lead ID=…` run it.
//
//	go run ./cmd/intake                         run: list new Leads in .intake/queue.md
//	go run ./cmd/intake -dismiss <id>           hide a Lead for good
//	go run ./cmd/intake -pick <id> [-kind k] [-vertical v]
//	                                            start a draft skeleton from a Lead (ticket 09)
//
// The Verticals come from -config (site.yaml). The queue and the
// seen-state live in -state, a local gitignored directory; Intake writes
// nothing under -content. All logic is in internal/intake; main is flags
// and exit codes:
//
//	0  done
//	1  the run failed: an invalid Watchlist, a Source that failed, I/O
//	2  usage: bad flags, unreadable site config
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"slices"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/intake"
)

const (
	exitOK     = 0
	exitFailed = 1
	exitUsage  = 2
)

func main() {
	watchlist := flag.String("watchlist", intake.DefaultWatchlist, "the Source Watchlist file")
	config := flag.String("config", "site.yaml", "site config; its verticals are the ones a Source may name")
	contentDir := flag.String("content", "content", "content directory; Intake never writes under it")
	stateDir := flag.String("state", ".intake", "local state directory (gitignored): the queue and the seen-state")
	dismiss := flag.String("dismiss", "", "dismiss the Lead with this ID for good")
	pick := flag.String("pick", "", "pick the Lead with this ID: write its draft skeleton")
	kind := flag.String("kind", "", "with -pick: override the Lead's suggested kind")
	vertical := flag.String("vertical", "", "with -pick: override the Lead's Vertical")
	flag.Parse()

	os.Exit(run(*watchlist, *config, *contentDir, *stateDir, *dismiss, *pick, intake.PickOptions{Kind: *kind, Vertical: *vertical}))
}

func run(watchlist, config, contentDir, stateDir, dismiss, pick string, opts intake.PickOptions) int {
	switch {
	case flag.NArg() > 0:
		return usage(fmt.Sprintf("unexpected argument %q", flag.Arg(0)))
	case dismiss != "" && pick != "":
		return usage("-dismiss and -pick are exclusive")
	case pick == "" && (opts.Kind != "" || opts.Vertical != ""):
		return usage("-kind and -vertical go with -pick")
	}
	cfg, err := builder.LoadConfig(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return exitUsage
	}
	in := &intake.Intake{
		Watchlist:  watchlist,
		Verticals:  verticals(cfg),
		StateDir:   stateDir,
		ContentDir: contentDir,
		Log:        os.Stdout,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	switch {
	case dismiss != "":
		err = in.Dismiss(ctx, dismiss)
	case pick != "":
		var file string
		if file, err = in.Pick(ctx, pick, opts); err == nil {
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

// verticals are site.yaml's Verticals in their display order.
func verticals(cfg builder.Config) []intake.Vertical {
	vs := slices.Clone(cfg.Verticals)
	slices.SortStableFunc(vs, func(a, b builder.Vertical) int { return a.Order - b.Order })
	var out []intake.Vertical
	for _, v := range vs {
		out = append(out, intake.Vertical{Slug: v.Slug, Name: v.Name})
	}
	return out
}

func usage(msg string) int {
	fmt.Fprintln(os.Stderr, "error:", msg)
	flag.Usage()
	return exitUsage
}

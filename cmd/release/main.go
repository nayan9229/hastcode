// Command release is the Author's release step (ADR-0012), run by
// scripts/release.sh (make release): flags and exit codes only; the work is
// internal/release.
//
//	release -content content                 # release every draft
//	release -content content -files "a b"    # only these (repo-relative or absolute)
//	release -content content -dry            # print the plan, write nothing
//	release -content content -gaps -lookback 3   # Gap Days, one "day slots" line each, for make daily
//	release -content content -redate 2026-11-13 -files path   # move a draft Explainer to a Gap Day
//
// Prints one line per Piece, tab-separated: "release" (dry) or "released",
// the content-relative path, the Publication Day, the verified day it gets,
// the title. Exit codes: 0 done; 1 the plan was refused or a file could not
// be read or written (reason on stderr); 2 bad flags.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"hastcode.com/site/internal/release"
)

func main() {
	content := flag.String("content", "content", "content directory")
	files := flag.String("files", "", "drafts to release, space- or comma-separated; empty = every draft")
	dry := flag.Bool("dry", false, "print the plan and write nothing")
	gaps := flag.Bool("gaps", false, "print the Gap Days inside -lookback, one \"day slots\" line each, and exit")
	lookBack := flag.Int("lookback", 3, "days before today to check for Gap Days")
	redate := flag.String("redate", "", "move the draft Explainer in -files to this Publication Day (YYYY-MM-DD) and exit")
	nowFlag := flag.String("now", "", "the clock, RFC 3339 (default: now)")
	flag.Parse()

	now := time.Now()
	if *nowFlag != "" {
		t, err := time.Parse(time.RFC3339, *nowFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "release: -now %q is not RFC 3339\n", *nowFlag)
			os.Exit(2)
		}
		now = t
	}
	if *gaps {
		g, err := release.GapDays(*content, now, *lookBack)
		if err != nil {
			fmt.Fprintln(os.Stderr, "release:", err)
			os.Exit(1)
		}
		for _, gap := range g {
			fmt.Printf("%s %d\n", gap.Day, gap.Slots)
		}
		return
	}
	if *redate != "" {
		files := strings.Fields(strings.ReplaceAll(*files, ",", " "))
		if len(files) != 1 {
			fmt.Fprintln(os.Stderr, "release: -redate takes exactly one -files path")
			os.Exit(2)
		}
		if err := release.Redate(*content, files[0], *redate); err != nil {
			fmt.Fprintln(os.Stderr, "release:", err)
			os.Exit(1)
		}
		return
	}
	plan, err := release.NewPlan(*content, strings.Fields(strings.ReplaceAll(*files, ",", " ")), now)
	if err != nil {
		fmt.Fprintln(os.Stderr, "release:", err)
		os.Exit(1)
	}
	verb := "released"
	if *dry {
		verb = "release"
	} else if err := plan.Apply(); err != nil {
		fmt.Fprintln(os.Stderr, "release:", err)
		os.Exit(1)
	}
	for _, p := range plan.Releases {
		fmt.Printf("%s\t%s\t%s\t%s\t%s\n", verb, p.Path, p.Day(), plan.Verified(p), p.Title)
	}
}

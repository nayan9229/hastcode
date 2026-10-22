// Command fetch writes a Repo Profile's Fact Pack (GLOSSARY; ADR-0005)
// from the GitHub REST API, or refreshes every Fact Pack under the content
// tree (GLOSSARY "Refresh"). `make fetch REPO=owner/name MD=<markdown>` and
// `make refresh` run it.
//
//	go run ./cmd/fetch -repo apache/iggy -markdown content/explainers/distributed-systems/apache-iggy.md
//	go run ./cmd/fetch -refresh -content content
//
// It reads one token from GITHUB_TOKEN (github.EnvToken) and refuses to run
// without it. The Fact Pack is written beside the Markdown with .json; the
// README text under -cache (the builder's build cache, default
// node_modules/.cache/hastcode) at readmes/<owner>/<name>/<file>, never in
// the content tree. In refresh mode only the Fact Packs whose facts changed
// are rewritten, one line per changed repository is printed, and the exit
// code says whether anything changed. All logic is in internal/factpack and
// internal/github; main is flags and exit codes:
//
//	0  done (refresh: nothing changed)
//	1  a fetch or write failed: GitHub rate limit reached, network, I/O
//	2  usage: bad flags, GITHUB_TOKEN unset
//	3  refresh: at least one Fact Pack was rewritten
//
// `go run` turns exit 3 into 1 ("exit status 3"): the Makefile builds the
// binary first, and a workflow that needs the code must do the same.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"hastcode.com/site/internal/factpack"
	"hastcode.com/site/internal/github"
)

const (
	exitOK      = 0
	exitFailed  = 1
	exitUsage   = 2
	exitChanged = 3
)

func main() {
	repo := flag.String("repo", "", "repository to fetch, owner/name (with -markdown)")
	markdown := flag.String("markdown", "", "the Repo Profile's Markdown file; the Fact Pack is written beside it")
	refresh := flag.Bool("refresh", false, "re-fetch every Repo Profile's Fact Pack under -content and rewrite the changed ones (exit 3 when any changed)")
	contentDir := flag.String("content", "content", "content directory (-refresh walks its explainers/)")
	cacheDir := flag.String("cache", "node_modules/.cache/hastcode", "build cache directory; README text is written under it (readmes/<owner>/<name>/)")
	flag.Parse()

	os.Exit(run(*repo, *markdown, *refresh, *contentDir, *cacheDir))
}

func run(repo, markdown string, refresh bool, contentDir, cacheDir string) int {
	switch {
	case refresh && (repo != "" || markdown != ""):
		return usage("-refresh takes neither -repo nor -markdown")
	case !refresh && (repo == "" || markdown == ""):
		return usage("-repo owner/name and -markdown <file> are required (or -refresh)")
	}
	client, err := github.NewClient(os.Getenv)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return exitUsage
	}
	f := &factpack.Fetcher{GitHub: client, CacheDir: cacheDir, Log: os.Stdout}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if refresh {
		report, err := f.Refresh(ctx, contentDir)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return exitFailed
		}
		if len(report.Changes) > 0 {
			return exitChanged
		}
		return exitOK
	}
	if _, err := f.Fetch(ctx, repo, markdown); err != nil {
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

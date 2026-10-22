// Command prune lists the Image Store (R2) objects the current content
// would not produce, with their sizes and a total, and deletes them only
// with --delete. Manual, dry run by default (`make prune-images`).
//
//	go run ./cmd/prune -content content -config site.yaml            # dry run
//	go run ./cmd/prune -content content -config site.yaml --delete   # delete the orphans
//
// It reads the same site.yaml (images.bucket, images.base_url), the same
// flags (internal/buildflags) and the same environment as an r2 build
// (imagestore.EnvAccountID, EnvAccessKeyID, EnvSecretAccessKey, optional
// EnvEndpoint), builds the content locally into a temporary
// directory to learn the referenced keys (Production plus Preview, so
// drafts in the tree keep their images), lists the bucket and diffs. Exit 2
// on a bad config, 1 on any other failure.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildflags"
)

func main() {
	shared := buildflags.Register(flag.CommandLine)
	del := flag.Bool("delete", false, "delete the orphans (default: dry run, list only)")
	flag.Parse()

	cfg, err := shared.LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}

	if _, err := builder.Prune(context.Background(), cfg, shared.Content, builder.PruneOptions{Delete: *del, Log: os.Stdout}); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

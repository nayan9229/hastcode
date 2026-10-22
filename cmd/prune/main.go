// Command prune lists the Image Store (R2) objects the current content
// would not produce, with their sizes and a total, and deletes them only
// with --delete. Manual, dry run by default (`make prune-images`).
//
//	go run ./cmd/prune -content content -config site.yaml            # dry run
//	go run ./cmd/prune -content content -config site.yaml --delete   # delete the orphans
//
// It reads the same site.yaml (images.bucket, images.base_url) and the same
// environment as an r2 build (CLOUDFLARE_ACCOUNT_ID, R2_ACCESS_KEY_ID,
// R2_SECRET_ACCESS_KEY, optional R2_ENDPOINT), builds the content locally
// into a temporary directory to learn the referenced keys (Production plus
// Preview, so drafts in the tree keep their images), lists the bucket and
// diffs. Exit 2 on a bad config, 1 on any other failure.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"hastcode.com/site/builder"
)

func main() {
	contentDir := flag.String("content", "content", "content directory")
	configPath := flag.String("config", "site.yaml", "site config YAML")
	toolsDir := flag.String("tools", ".", "directory with package.json and node_modules (run npm ci there)")
	cacheDir := flag.String("cache", "", "build cache directory (default <tools>/node_modules/.cache/hastcode)")
	imageCache := flag.String("image-cache", "", "hero image cache dir (overrides images.cache_dir; default "+builder.DefaultImageCacheDir+")")
	del := flag.Bool("delete", false, "delete the orphans (default: dry run, list only)")
	flag.Parse()

	cfg, err := builder.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}
	cfg.ToolsDir, cfg.CacheDir = *toolsDir, *cacheDir
	if *imageCache != "" {
		cfg.Images.CacheDir = *imageCache
	}

	if _, err := builder.Prune(context.Background(), cfg, *contentDir, builder.PruneOptions{Delete: *del, Log: os.Stdout}); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

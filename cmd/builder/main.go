// Command builder builds the static site: content/ -> .vercel/output/.
//
//	go run ./cmd/builder -content content -config site.yaml -out .vercel/output -mode production
//
// Build with CGO_ENABLED=0 -tags nodynamic (the Makefile does) so image
// encoders run as WASM and output is identical on every machine.
//
// It prints every error and warning and exits 1 if the build failed.
// Run `make generate` (templ) before building from a fresh checkout.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"hastcode.com/site/builder"
)

func main() {
	contentDir := flag.String("content", "content", "content directory")
	configPath := flag.String("config", "site.yaml", "site config YAML")
	outDir := flag.String("out", ".vercel/output", "Build Output API directory")
	mode := flag.String("mode", string(builder.Production), "production or preview")
	imageCache := flag.String("image-cache", "", "hero image cache dir (overrides images.cache_dir; default "+builder.DefaultImageCacheDir+")")
	flag.Parse()

	os.Exit(run(*contentDir, *configPath, *outDir, *mode, *imageCache))
}

func run(contentDir, configPath, outDir, modeFlag, imageCache string) int {
	cfg, err := builder.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	if cfg.Mode, err = builder.ParseMode(modeFlag); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	if imageCache != "" {
		cfg.Images.CacheDir = imageCache
	}

	start := time.Now()
	res := builder.Build(cfg, contentDir, outDir)
	fmt.Fprint(os.Stderr, res.String())
	if res.Failed() {
		fmt.Fprintf(os.Stderr, "build failed: %d error(s), %d warning(s)\n", len(res.Errors), len(res.Warnings))
		return 1
	}
	fmt.Fprintf(os.Stderr, "built %s in %s (%d warning(s))\n", outDir, time.Since(start).Round(time.Millisecond), len(res.Warnings))
	return 0
}

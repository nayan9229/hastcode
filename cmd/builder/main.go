// Command builder builds the static site: content/ -> .vercel/output/.
//
//	go run ./cmd/builder -content content -config site.yaml -out .vercel/output -mode production
//
// Build with CGO_ENABLED=0 -tags nodynamic (the Makefile does) so image
// encoders run as WASM and output is identical on every machine.
//
// It prints every error and warning and exits 1 if the build failed.
// Run `make generate` (templ) and `npm ci` (Mermaid renderer, Pagefind)
// before building from a fresh checkout. The build ends by indexing the
// output with Pagefind (static/pagefind/); a production build fails
// without that index, and -skip-search-index is for preview builds only.
// With -image-store r2 (CI) every image is uploaded to R2 and referenced
// on images.base_url; the R2 credentials come from the environment.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/imagestore"
)

func main() {
	contentDir := flag.String("content", "content", "content directory")
	configPath := flag.String("config", "site.yaml", "site config YAML")
	outDir := flag.String("out", ".vercel/output", "Build Output API directory")
	mode := flag.String("mode", string(builder.Production), "production or preview")
	toolsDir := flag.String("tools", ".", "directory with package.json and node_modules (run npm ci there)")
	cacheDir := flag.String("cache", "", "build cache directory (default <tools>/node_modules/.cache/hastcode)")
	imageCache := flag.String("image-cache", "", "hero image cache dir (overrides images.cache_dir; default "+builder.DefaultImageCacheDir+")")
	imageStore := flag.String("image-store", "", "image store: local or r2 (overrides images.store; r2 reads "+imagestore.EnvAccountID+", "+imagestore.EnvAccessKeyID+", "+imagestore.EnvSecretAccessKey+" from the environment)")
	skipSearch := flag.Bool("skip-search-index", false, "do not run Pagefind (preview mode only; a production build fails without the index)")
	flag.Parse()

	os.Exit(run(*contentDir, *configPath, *outDir, *mode, *toolsDir, *cacheDir, *imageCache, *imageStore, *skipSearch))
}

func run(contentDir, configPath, outDir, modeFlag, toolsDir, cacheDir, imageCache, imageStore string, skipSearch bool) int {
	cfg, err := builder.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	if cfg.Mode, err = builder.ParseMode(modeFlag); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	cfg.ToolsDir, cfg.CacheDir = toolsDir, cacheDir
	cfg.SkipSearchIndex = skipSearch
	if imageCache != "" {
		cfg.Images.CacheDir = imageCache
	}
	if imageStore != "" {
		cfg.Images.Store = builder.ImageStore(imageStore)
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

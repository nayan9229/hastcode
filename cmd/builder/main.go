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
	"hastcode.com/site/internal/buildflags"
	"hastcode.com/site/internal/imagestore"
)

func main() {
	shared := buildflags.Register(flag.CommandLine)
	outDir := flag.String("out", ".vercel/output", "Build Output API directory")
	mode := flag.String("mode", string(builder.Production), "production or preview")
	imageStore := flag.String("image-store", "", "image store: local or r2 (overrides images.store; r2 reads "+imagestore.EnvAccountID+", "+imagestore.EnvAccessKeyID+", "+imagestore.EnvSecretAccessKey+" from the environment)")
	skipSearch := flag.Bool("skip-search-index", false, "do not run Pagefind (preview mode only; a production build fails without the index)")
	flag.Parse()

	os.Exit(run(shared, *outDir, *mode, *imageStore, *skipSearch))
}

func run(shared *buildflags.Flags, outDir, modeFlag, imageStore string, skipSearch bool) int {
	cfg, err := shared.LoadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	if cfg.Mode, err = builder.ParseMode(modeFlag); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 2
	}
	cfg.SkipSearchIndex = skipSearch
	if imageStore != "" {
		cfg.Images.Store = builder.ImageStoreMode(imageStore)
	}

	start := time.Now()
	res := builder.Build(cfg, shared.Content, outDir)
	fmt.Fprint(os.Stderr, res.String())
	if res.Failed() {
		fmt.Fprintf(os.Stderr, "build failed: %d error(s), %d warning(s)\n", len(res.Errors), len(res.Warnings))
		return 1
	}
	fmt.Fprintf(os.Stderr, "built %s in %s (%d warning(s))\n", outDir, time.Since(start).Round(time.Millisecond), len(res.Warnings))
	return 0
}

// Command adapt applies the deploy-adapter rule (internal/deploy) to a build
// output in place: it removes every file only the other Target reads and
// writes the Target's overlays (Vercel: a disallow-all robots.txt), so the
// upload that follows ships exactly the Target's file set.
//
//	go run ./cmd/adapt -target cloudflare -out .vercel/output
//	go run ./cmd/adapt -target vercel -out .vercel/output
//
// Run it once per deploy, after the build and before the upload. It prints
// each removed file; exit 2 on a bad target, 1 on an I/O error.
package main

import (
	"flag"
	"fmt"
	"os"

	"hastcode.com/site/internal/deploy"
)

func main() {
	target := flag.String("target", "", `"cloudflare" (Production Target) or "vercel" (QA Target)`)
	outDir := flag.String("out", ".vercel/output", "build output directory")
	flag.Parse()

	t, err := deploy.ParseTarget(*target)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}
	removed, err := deploy.Apply(*outDir, t)
	for _, f := range removed {
		fmt.Fprintf(os.Stderr, "removed %s (not part of the %s file set)\n", f, t)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

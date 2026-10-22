// Command target names a Target from any of its spellings (deploy.ParseTarget:
// "cloudflare", "cf", "PUBCF-…"; the same for Vercel) and prints all of
// them, one key=value per line, so the publish scripts never spell a Target
// themselves:
//
//	$ go run ./cmd/target PUBCF-20261022-101500
//	target=cloudflare
//	code=cf
//	tag_prefix=PUBCF-
//
// Exit 2 when the argument names no Target.
package main

import (
	"fmt"
	"os"

	"hastcode.com/site/internal/deploy"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: target NAME|CODE|TAG")
		os.Exit(2)
	}
	t, err := deploy.ParseTarget(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}
	fmt.Printf("target=%s\ncode=%s\ntag_prefix=%s\n", t, t.Code(), t.TagPrefix())
}

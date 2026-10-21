// Command subscribe runs the newsletter subscribe function locally, the
// same handler Vercel serves from api/subscribe.go, behind the built site's
// routes:
//
//	make build
//	BUTTONDOWN_API_KEY=... go run ./cmd/subscribe   # http://localhost:3000
//
// Requests go through the routes in <out>/config.json (410s, the /home
// redirect, trailing-slash redirects, cache headers, the 404 fallback)
// applied by vercel.NewRouter, then to POST /api/subscribe or the files
// under <out>/static. Point BUTTONDOWN_BASE_URL at a fake to avoid touching
// the real list.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"hastcode.com/site/internal/subscribe"
	"hastcode.com/site/internal/vercel"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	addr := flag.String("addr", "localhost:"+port, "listen address")
	out := flag.String("out", ".vercel/output", "Build Output dir to serve beside the function (empty: function only)")
	flag.Parse()

	cfg := subscribe.ConfigFromEnv()
	if cfg.APIKey == "" {
		log.Printf("warning: %s is not set; subscribe requests will fail with 502", subscribe.EnvAPIKey)
	}
	fn := subscribe.New(cfg)

	mux := http.NewServeMux()
	mux.Handle(subscribe.Path, fn)
	var h http.Handler = mux
	if *out != "" {
		var err error
		h, err = vercel.NewRouter(*out, map[string]http.Handler{subscribe.Path: fn})
		if err != nil {
			log.Fatalf("serving %s (run make build first): %v", *out, err)
		}
	}
	log.Printf("serving %s on http://%s (build output: %q)", subscribe.Path, *addr, *out)
	log.Fatal(http.ListenAndServe(*addr, h))
}

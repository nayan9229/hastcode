// Command subscribe runs the newsletter subscribe function locally, the
// same handler Vercel serves from api/subscribe.go, next to the built site:
//
//	make build
//	BUTTONDOWN_API_KEY=... go run ./cmd/subscribe   # http://localhost:3000
//
// It serves POST /api/subscribe and the static files under -static, so the
// newsletter form on any built page works end to end. Point
// BUTTONDOWN_BASE_URL at a fake to avoid touching the real list.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"hastcode.com/site/internal/subscribe"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	addr := flag.String("addr", "localhost:"+port, "listen address")
	static := flag.String("static", ".vercel/output/static", "built site to serve beside the function (empty: function only)")
	flag.Parse()

	cfg := subscribe.ConfigFromEnv()
	if cfg.APIKey == "" {
		log.Printf("warning: %s is not set; subscribe requests will fail with 502", subscribe.EnvAPIKey)
	}

	mux := http.NewServeMux()
	mux.Handle(subscribe.Path, subscribe.New(cfg))
	if *static != "" {
		mux.Handle("/", http.FileServer(http.Dir(*static)))
	}
	log.Printf("serving %s on http://%s (static: %q)", subscribe.Path, *addr, *static)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

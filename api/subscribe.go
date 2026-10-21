// Package handler is the Vercel Go function entry for POST /api/subscribe.
//
// Vercel's Go runtime turns each api/*.go file that exports an
// http.HandlerFunc into one function served at the file's path
// (vercel.com/docs/functions/runtimes/go). All behaviour lives in
// internal/subscribe; this file only wires it to the environment. Keep this
// directory free of other .go files (each would become a function).
package handler

import (
	"net/http"
	"sync"

	"hastcode.com/site/internal/subscribe"
)

// The handler (and its in-memory rate limiter) is built once per warm
// function instance.
var subscribeHandler = sync.OnceValue(func() http.Handler {
	return subscribe.New(subscribe.ConfigFromEnv())
})

// Handler serves POST /api/subscribe.
func Handler(w http.ResponseWriter, r *http.Request) {
	subscribeHandler().ServeHTTP(w, r)
}

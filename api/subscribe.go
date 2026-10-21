// Package handler is the Vercel Go function entry for POST /api/subscribe.
//
// Vercel's Go runtime turns each api/*.go file that exports an
// http.HandlerFunc into one function served at the file's path
// (vercel.com/docs/functions/runtimes/go). All behaviour lives in
// subscribe; this file only wires it to the environment. Keep this
// directory free of other .go files (each would become a function).
//
// Import nothing under internal/: @vercel/go renames the root module to
// "handler" while compiling this file and requires hastcode.com/site
// through a replace, so an internal/ import fails with "use of internal
// package ... not allowed" (found with a local `vercel build`, ticket 15).
// That is why the function's package lives at subscribe/, not
// internal/subscribe/.
package handler

import (
	"net/http"
	"sync"

	"hastcode.com/site/subscribe"
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

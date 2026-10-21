// Package subscribe is the newsletter subscribe function: the one dynamic
// endpoint on the site (ADR-0002). New returns an http.Handler for
// POST /api/subscribe that forwards a Reader's email to Buttondown.
package subscribe

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config configures the handler. Only APIKey is required.
type Config struct {
	// APIKey is the Buttondown API key. It is sent upstream only and never
	// appears in a response.
	APIKey string
	// BaseURL is Buttondown's API origin; default DefaultBaseURL. Tests
	// point it at a fake server.
	BaseURL string
	// RateLimit is how many requests one IP may make per RateWindow;
	// defaults DefaultRateLimit per DefaultRateWindow.
	RateLimit  int
	RateWindow time.Duration
}

// Environment variables read by ConfigFromEnv.
const (
	EnvAPIKey  = "BUTTONDOWN_API_KEY"  // required, secret
	EnvBaseURL = "BUTTONDOWN_BASE_URL" // optional, default DefaultBaseURL
)

// ConfigFromEnv builds the production Config from the environment.
func ConfigFromEnv() Config {
	return Config{APIKey: os.Getenv(EnvAPIKey), BaseURL: os.Getenv(EnvBaseURL)}
}

// Rate limit defaults: enough for a Reader who mistypes a few times.
const (
	DefaultRateLimit  = 5
	DefaultRateWindow = 10 * time.Minute
)

// DefaultBaseURL is Buttondown's API origin.
const DefaultBaseURL = "https://api.buttondown.com"

// Form contract shared with the newsletter form partial
// (internal/views/partial_newsletter.templ).
const (
	// Path is where the form posts.
	Path = "/api/subscribe"
	// SubscribedPath is where no-JS form posts land on success.
	SubscribedPath = "/subscribed/"
	// HoneypotField is the visually hidden field only bots fill in.
	HoneypotField = "website"
)

type handler struct {
	cfg     Config
	client  *http.Client
	limiter *limiter
}

// New returns the subscribe handler.
func New(cfg Config) http.Handler {
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}
	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")
	if cfg.RateLimit <= 0 {
		cfg.RateLimit = DefaultRateLimit
	}
	if cfg.RateWindow <= 0 {
		cfg.RateWindow = DefaultRateWindow
	}
	return &handler{
		cfg:     cfg,
		client:  &http.Client{Timeout: 8 * time.Second},
		limiter: newLimiter(cfg.RateLimit, cfg.RateWindow),
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	if r.PostFormValue(HoneypotField) != "" {
		// A bot filled the field humans never see: pretend it worked.
		succeed(w, r)
		return
	}
	email := strings.TrimSpace(r.PostFormValue("email"))
	source := cleanSource(r.PostFormValue("source"))
	if !h.limiter.allow(clientIP(r), time.Now()) {
		w.Header().Set("Retry-After", strconv.Itoa(int(h.cfg.RateWindow.Seconds())))
		fail(w, r, http.StatusTooManyRequests, msgRateLimited, source)
		return
	}
	if !validEmail(email) {
		fail(w, r, http.StatusBadRequest, msgInvalidEmail, source)
		return
	}
	if h.cfg.APIKey == "" {
		log.Printf("subscribe: %s is not set", EnvAPIKey)
		fail(w, r, http.StatusBadGateway, msgUpstream, source)
		return
	}
	switch err := h.forward(r.Context(), email, source, clientIP(r)); {
	case err == nil, errors.Is(err, errAlreadySubscribed):
		// Same answer either way, so the form never reveals who is subscribed.
		succeed(w, r)
	case errors.Is(err, errRejectedEmail):
		fail(w, r, http.StatusBadRequest, msgInvalidEmail, source)
	default:
		// Log for the Author (Vercel runtime logs); never echo upstream
		// detail to the Reader.
		log.Printf("subscribe: %v", err)
		fail(w, r, http.StatusBadGateway, msgUpstream, source)
	}
}

const (
	msgSubscribed   = "Thanks! Check your inbox to confirm your subscription."
	msgInvalidEmail = "Please enter a valid email address."
	msgRateLimited  = "Too many attempts from your network. Please try again in a few minutes."
	msgUpstream     = "We couldn't reach our newsletter service. Please try again in a few minutes."
)

// validEmail accepts a bare addr-spec (no display name) whose domain has a
// dot. Deliverability is Buttondown's problem; this only stops obvious junk.
func validEmail(s string) bool {
	if len(s) > 254 {
		return false
	}
	a, err := mail.ParseAddress(s)
	if err != nil || a.Address != s {
		return false
	}
	at := strings.LastIndexByte(s, '@')
	domain := s[at+1:]
	return at > 0 && strings.Contains(domain, ".") && !strings.HasPrefix(domain, ".") && !strings.HasSuffix(domain, ".")
}

// cleanSource keeps source only if it is a site-relative path, so it is
// safe to tag with and to link back to (no open redirect).
func cleanSource(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 200 || !strings.HasPrefix(s, "/") || strings.HasPrefix(s, "//") || strings.ContainsAny(s, "\\\r\n\t <>\"'") {
		return "/"
	}
	return s
}

// fail answers with an error: JSON for the script, a minimal HTML page
// linking back to the source page for a no-JS form post.
func fail(w http.ResponseWriter, r *http.Request, status int, message, source string) {
	if wantsJSON(r) {
		writeJSON(w, status, false, message)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><meta name="robots" content="noindex"><title>Subscribe | HastCode</title></head>
<body><main><h1>Subscription not completed</h1><p>%s</p><p><a href="%s">Go back</a></p></main></body></html>
`, html.EscapeString(message), html.EscapeString(source))
}

// wantsJSON reports whether the caller is the form's script (Accept JSON)
// rather than a plain no-JS form post.
func wantsJSON(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "application/json")
}

func succeed(w http.ResponseWriter, r *http.Request) {
	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, true, msgSubscribed)
		return
	}
	http.Redirect(w, r, SubscribedPath, http.StatusSeeOther)
}

func writeJSON(w http.ResponseWriter, status int, ok bool, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(struct {
		OK      bool   `json:"ok"`
		Message string `json:"message"`
	}{ok, message})
}

package subscribe

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// limiter is a fixed-window, in-memory, per-IP counter. On Vercel each warm
// function instance has its own copy and a cold start resets it, so this is
// best-effort abuse damping (spec "Subscribe function"), not a quota.
type limiter struct {
	max    int
	window time.Duration

	mu   sync.Mutex
	hits map[string]*window
}

type window struct {
	start time.Time
	count int
}

// maxTracked bounds memory: past it, expired windows are dropped.
const maxTracked = 10_000

func newLimiter(max int, w time.Duration) *limiter {
	return &limiter{max: max, window: w, hits: map[string]*window{}}
}

// allow records one request from ip and reports whether it is within the limit.
func (l *limiter) allow(ip string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.hits) >= maxTracked {
		for k, w := range l.hits {
			if now.Sub(w.start) >= l.window {
				delete(l.hits, k)
			}
		}
	}
	w, ok := l.hits[ip]
	if !ok || now.Sub(w.start) >= l.window {
		w = &window{start: now}
		l.hits[ip] = w
	}
	w.count++
	return w.count <= l.max
}

// clientIP is the Reader's address. Vercel overwrites X-Forwarded-For with
// the real client IP (so it cannot be spoofed there); locally it falls back
// to the socket address.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		first, _, _ := strings.Cut(xff, ",")
		if ip := strings.TrimSpace(first); ip != "" {
			return ip
		}
	}
	if ip := strings.TrimSpace(r.Header.Get("X-Real-IP")); ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

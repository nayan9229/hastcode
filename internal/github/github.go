// Package github is the small GitHub REST client behind the fetch command
// (ADR-0005; spec "Fetch command"): one authenticated JSON GET at a time
// over net/http. It reads the one token from the environment (EnvToken,
// the variable GitHub Actions provides) and refuses to run without it,
// sends the identifying User-Agent, honours Retry-After and the primary
// rate-limit headers with a bounded backoff, retries a 202 ("still
// computing") a bounded number of times, and stops with a RateLimitError
// rather than wait past the bound
// (docs/research/github-trending-geography-and-data-access.md §2.3).
// Requests are serial, as GitHub's best practices ask; the fetch never
// pools tokens.
package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	// EnvToken is the environment variable holding the token: GITHUB_TOKEN,
	// which GitHub Actions provides to every job (1,000 requests an hour
	// per repository) and which a fine-grained personal access token with
	// read access to public repositories fills locally.
	EnvToken = "GITHUB_TOKEN"
	// APIURL is the public REST API. A Fact Pack's sources are spelled on
	// it whatever BaseURL the Client was pointed at.
	APIURL = "https://api.github.com"
	// UserAgent identifies the fetch command on every request, as GitHub's
	// terms ask.
	UserAgent = "hastcode-fetch/1 (+https://hastcode.com; Repo Profile Fact Packs)"
	// DefaultMaxWait bounds one backoff. A longer wait (a primary limit
	// that resets in an hour) stops the command with a RateLimitError
	// instead.
	DefaultMaxWait = 2 * time.Minute
	// DefaultRetries bounds the retries of one request: 202 answers,
	// rate-limit waits and 5xx answers each count.
	DefaultRetries = 3

	apiVersion = "2022-11-28"
	// computingDelay is the first wait after a 202; it doubles per retry.
	computingDelay = 2 * time.Second
)

// ErrNotFound is the error a 404 wraps: the repository, file or endpoint
// does not exist (or the token cannot see it).
var ErrNotFound = errors.New("not found (404)")

// Client makes authenticated GET requests against BaseURL. The zero value
// is unusable; NewClient fills every field from the environment and the
// defaults. Tests point BaseURL at an httptest server and replace Sleep.
type Client struct {
	// BaseURL is where requests go; APIURL unless a test overrides it.
	BaseURL string
	// Token is sent as a Bearer token on every request.
	Token string
	// HTTP makes the requests.
	HTTP *http.Client
	// MaxWait bounds a single backoff (DefaultMaxWait).
	MaxWait time.Duration
	// Retries bounds the retries of one request (DefaultRetries).
	Retries int
	// Sleep waits between retries; the default sleeps, honouring ctx.
	Sleep func(ctx context.Context, d time.Duration) error
	// Now is the clock the rate-limit reset is measured against.
	Now func() time.Time
}

// NewClient returns a Client holding the token read through getenv. It
// refuses, naming EnvToken, when the variable is unset: unauthenticated
// GitHub allows 60 requests an hour, which one fetch would nearly spend.
func NewClient(getenv func(string) string) (*Client, error) {
	token := getenv(EnvToken)
	if token == "" {
		return nil, fmt.Errorf("%s is not set: the fetch command reads one GitHub token from that variable (GitHub Actions provides it; locally, create a fine-grained personal access token with read access to public repositories) and refuses to run unauthenticated, where GitHub allows 60 requests an hour", EnvToken)
	}
	return &Client{
		BaseURL: APIURL,
		Token:   token,
		HTTP:    &http.Client{Timeout: 30 * time.Second},
		MaxWait: DefaultMaxWait,
		Retries: DefaultRetries,
		Sleep:   sleep,
		Now:     time.Now,
	}, nil
}

// RateLimitError is the documented stop: GitHub answered 403 or 429 with
// rate-limit headers and the wait they ask for exceeds the Client's bound
// (or its retries are spent). The command prints it and exits non-zero.
type RateLimitError struct {
	// Path is the request that hit the limit.
	Path string
	// Resource, Limit and Remaining are the X-RateLimit-* headers ("core",
	// 1000, 0); Reset is when the limit resets, zero when GitHub sent only
	// Retry-After.
	Resource         string
	Limit, Remaining int
	Reset            time.Time
	// Wait is how long GitHub asked the caller to wait.
	Wait time.Duration
}

func (e *RateLimitError) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "GitHub rate limit reached: GET %s (%d of %d requests left on the %s limit", e.Path, e.Remaining, e.Limit, cmpOr(e.Resource, "primary"))
	if !e.Reset.IsZero() {
		fmt.Fprintf(&b, "; resets at %s, in %s", e.Reset.UTC().Format(time.RFC3339), e.Wait.Round(time.Second))
	} else {
		fmt.Fprintf(&b, "; GitHub asks for a %s wait", e.Wait.Round(time.Second))
	}
	b.WriteString("); stopping instead of waiting. Run the command again after the reset; never add a second token, GitHub's terms forbid pooling them")
	return b.String()
}

// StatusError is any other non-success answer, with GitHub's message.
type StatusError struct {
	Path    string
	Status  int
	Message string
}

func (e *StatusError) Error() string {
	s := fmt.Sprintf("GitHub answered %d for GET %s", e.Status, e.Path)
	if e.Message != "" {
		s += ": " + e.Message
	}
	if e.Status == http.StatusUnauthorized {
		s += " (check " + EnvToken + ")"
	}
	return s
}

// Get requests path (relative to BaseURL, query included) and decodes the
// JSON answer into v (nil to discard it). It returns the answer's headers
// for pagination (LastPage). A 404 wraps ErrNotFound; a 204 leaves v
// untouched; a 202 is retried after a growing delay, a rate-limited 403 or
// 429 after the wait GitHub asks for, and a 5xx after a short backoff,
// each at most Retries times and never longer than MaxWait at once.
func (c *Client) Get(ctx context.Context, path string, v any) (http.Header, error) {
	for attempt := 0; ; attempt++ {
		status, header, body, err := c.do(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("GET %s: %w", path, err)
		}
		var wait time.Duration
		switch {
		case status == http.StatusOK:
			if v != nil {
				if err := json.Unmarshal(body, v); err != nil {
					return header, fmt.Errorf("GET %s: decoding the answer: %w", path, err)
				}
			}
			return header, nil
		case status == http.StatusNoContent:
			return header, nil
		case status == http.StatusNotFound:
			return header, fmt.Errorf("GET %s: %w", path, ErrNotFound)
		case status == http.StatusAccepted:
			if attempt >= c.retries() {
				return header, fmt.Errorf("GET %s: GitHub is still computing the answer (202) after %d attempts; run the command again in a minute", path, attempt+1)
			}
			wait = computingDelay << attempt
		case status == http.StatusForbidden || status == http.StatusTooManyRequests:
			limit, ok := c.rateLimit(path, header, body)
			if !ok {
				return header, statusError(path, status, body)
			}
			if limit.Wait > c.maxWait() || attempt >= c.retries() {
				return header, limit
			}
			wait = limit.Wait
		case status >= 500:
			if attempt >= c.retries() {
				return header, statusError(path, status, body)
			}
			wait = time.Second << attempt
		default:
			return header, statusError(path, status, body)
		}
		if err := c.sleep(ctx, wait); err != nil {
			return header, fmt.Errorf("GET %s: %w", path, err)
		}
	}
}

func (c *Client) do(ctx context.Context, path string) (int, http.Header, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cmpOr(c.BaseURL, APIURL)+path, nil)
	if err != nil {
		return 0, nil, nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", apiVersion)
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("User-Agent", UserAgent)
	httpClient := c.HTTP
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, nil, err
	}
	return resp.StatusCode, resp.Header, body, nil
}

// rateLimit reads a 403/429 answer as GitHub's rate-limit rules say: a
// Retry-After header is the wait; else X-RateLimit-Remaining 0 means wait
// for X-RateLimit-Reset; else a message naming the rate limit (a secondary
// limit) means wait at least a minute. Anything else is not a rate limit.
func (c *Client) rateLimit(path string, h http.Header, body []byte) (*RateLimitError, bool) {
	e := &RateLimitError{Path: path, Resource: h.Get("X-RateLimit-Resource")}
	e.Limit, _ = strconv.Atoi(h.Get("X-RateLimit-Limit"))
	e.Remaining, _ = strconv.Atoi(h.Get("X-RateLimit-Remaining"))
	if n, err := strconv.ParseInt(h.Get("X-RateLimit-Reset"), 10, 64); err == nil {
		e.Reset = time.Unix(n, 0).UTC()
	}
	switch {
	case h.Get("Retry-After") != "":
		secs, _ := strconv.Atoi(h.Get("Retry-After"))
		e.Wait = time.Duration(secs) * time.Second
	case h.Get("X-RateLimit-Remaining") == "0" && !e.Reset.IsZero():
		e.Wait = e.Reset.Sub(c.now()) + time.Second
	case strings.Contains(strings.ToLower(message(body)), "rate limit"):
		e.Wait = time.Minute
	default:
		return nil, false
	}
	e.Wait = max(e.Wait, time.Second)
	return e, true
}

// LastPage reads the page number of the rel="last" link in a Link header,
// GitHub's pagination: with per_page=1 it is the number of items.
func LastPage(h http.Header) (int, bool) {
	for part := range strings.SplitSeq(h.Get("Link"), ",") {
		link, rel, ok := strings.Cut(part, ";")
		if !ok || !strings.Contains(rel, `rel="last"`) {
			continue
		}
		u, err := url.Parse(strings.Trim(strings.TrimSpace(link), "<>"))
		if err != nil {
			return 0, false
		}
		n, err := strconv.Atoi(u.Query().Get("page"))
		return n, err == nil
	}
	return 0, false
}

func statusError(path string, status int, body []byte) error {
	return &StatusError{Path: path, Status: status, Message: message(body)}
}

// message is the "message" of a GitHub error body, "" when there is none.
func message(body []byte) string {
	var m struct {
		Message string `json:"message"`
	}
	_ = json.Unmarshal(body, &m)
	return m.Message
}

func (c *Client) retries() int {
	if c.Retries <= 0 {
		return DefaultRetries
	}
	return c.Retries
}

func (c *Client) maxWait() time.Duration {
	if c.MaxWait <= 0 {
		return DefaultMaxWait
	}
	return c.MaxWait
}

func (c *Client) now() time.Time {
	if c.Now == nil {
		return time.Now()
	}
	return c.Now()
}

func (c *Client) sleep(ctx context.Context, d time.Duration) error {
	if c.Sleep == nil {
		return sleep(ctx, d)
	}
	return c.Sleep(ctx, d)
}

func sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

func cmpOr(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

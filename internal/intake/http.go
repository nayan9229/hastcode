package intake

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// maxFeedBytes bounds one feed body; the largest feed in the research was
// well under a megabyte.
const maxFeedBytes = 8 << 20

// UserAgent is sent with every request Intake makes, so a publisher sees
// who is polling and how to reach us (research §5.2; content/pages/contact.md).
const UserAgent = "HastCode-Intake/1 (+https://hastcode.com/contact/; mira-vance@hastcode.com)"

// DefaultTimeout bounds one request when Intake.Timeout is zero: a Source
// that hangs fails its own fetch, not the run.
const DefaultTimeout = 30 * time.Second

// response is one conditional GET's outcome.
type response struct {
	// NotModified: the server answered 304 to the validators sent; Body
	// is empty and the Source has nothing new.
	NotModified bool
	Body        []byte
	// State is prev with the validators this response carried (a 200
	// replaces both; a 304 updates whichever it repeats).
	State SourceState
}

// get fetches url with the identifying User-Agent and a per-request
// timeout, conditional on prev's validators (If-None-Match and
// If-Modified-Since together; a server that knows both uses the former,
// RFC 9110 §13.1.3). It returns a 200's body or a 304; any other status
// is an error.
func (in *Intake) get(ctx context.Context, url string, prev SourceState) (response, error) {
	ctx, cancel := context.WithTimeout(ctx, in.timeout())
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return response{}, err
	}
	req.Header.Set("User-Agent", UserAgent)
	if prev.ETag != "" {
		req.Header.Set("If-None-Match", prev.ETag)
	}
	if prev.LastModified != "" {
		req.Header.Set("If-Modified-Since", prev.LastModified)
	}
	resp, err := in.httpClient().Do(req)
	if err != nil {
		return response{}, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusNotModified:
		st := prev
		st.ETag = revalidated(resp.Header, "ETag", prev.ETag)
		st.LastModified = revalidated(resp.Header, "Last-Modified", prev.LastModified)
		return response{NotModified: true, State: st}, nil
	case http.StatusOK:
	default:
		return response{}, fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxFeedBytes+1))
	if err != nil {
		return response{}, fmt.Errorf("GET %s: %w", url, err)
	}
	if len(body) > maxFeedBytes {
		return response{}, fmt.Errorf("GET %s: body over %d bytes", url, maxFeedBytes)
	}
	st := prev
	st.ETag, st.LastModified = resp.Header.Get("ETag"), resp.Header.Get("Last-Modified")
	return response{Body: body, State: st}, nil
}

func (in *Intake) timeout() time.Duration {
	if in.Timeout <= 0 {
		return DefaultTimeout
	}
	return in.Timeout
}

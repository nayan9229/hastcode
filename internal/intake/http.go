package intake

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// maxFeedBytes bounds one feed body; the largest feed in the research was
// well under a megabyte.
const maxFeedBytes = 8 << 20

// get fetches url with a plain GET and returns the body of a 200.
func (in *Intake) get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := in.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxFeedBytes+1))
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", url, err)
	}
	if len(body) > maxFeedBytes {
		return nil, fmt.Errorf("GET %s: body over %d bytes", url, maxFeedBytes)
	}
	return body, nil
}

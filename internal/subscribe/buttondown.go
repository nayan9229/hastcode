package subscribe

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

// Buttondown subscriber API, checked 2026-10-21 against
// https://docs.buttondown.com/api-subscribers-create and
// https://docs.buttondown.com/api-error-codes:
//
//	POST {base}/v1/subscribers
//	Authorization: Token <key>
//	{"email_address": "...", "tags": ["..."], "ip_address": "..."}
//	-> 201 created
//	-> 4xx {"code": "...", "detail": "...", "metadata": {...}}
//
// A duplicate email is a 400 (code subscriber_already_exists; older API
// versions used email_already_exists). We send no collision header, so an
// existing subscriber is left untouched. Tags are created on the fly only on
// plans that include tags; otherwise Buttondown answers 403 feature_disabled.

var (
	errAlreadySubscribed = errors.New("buttondown: subscriber already exists")
	errRejectedEmail     = errors.New("buttondown: email rejected")
)

// forward creates the subscriber in Buttondown, tagged with the source page.
func (h *handler) forward(ctx context.Context, email, source, ip string) error {
	sub := map[string]any{
		"email_address": email,
		"tags":          []string{source},
	}
	if net.ParseIP(ip) != nil {
		// Buttondown uses it for location and its spam firewall.
		sub["ip_address"] = ip
	}
	body, err := json.Marshal(sub)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.cfg.BaseURL+"/v1/subscribers", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Token "+h.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("buttondown: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 == 2 {
		return nil
	}

	var apiErr struct {
		Code   string `json:"code"`
		Detail any    `json:"detail"`
	}
	_ = json.NewDecoder(io.LimitReader(resp.Body, 64<<10)).Decode(&apiErr)
	switch {
	case resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusConflict:
		switch {
		case strings.HasSuffix(apiErr.Code, "already_exists"):
			return errAlreadySubscribed
		case strings.HasPrefix(apiErr.Code, "email_"):
			return fmt.Errorf("%w: %s", errRejectedEmail, apiErr.Code)
		}
	case resp.StatusCode == http.StatusUnprocessableEntity:
		return fmt.Errorf("%w: validation %v", errRejectedEmail, apiErr.Detail)
	}
	return fmt.Errorf("buttondown: status %d code %q detail %v", resp.StatusCode, apiErr.Code, apiErr.Detail)
}

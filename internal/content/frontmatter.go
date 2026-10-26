// Package content parses content files (Markdown with YAML frontmatter)
// into domain types: Explainer, Brief, Author. It knows nothing about Verticals
// config or output; the builder validates and renders.
package content

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

var delim = []byte("---")

// splitFrontmatter separates a leading "---" YAML block from the body.
func splitFrontmatter(data []byte) (front, body []byte, err error) {
	data = bytes.TrimPrefix(data, []byte("\ufeff"))
	first, rest, ok := bytes.Cut(data, []byte("\n"))
	if !ok || !bytes.Equal(bytes.TrimSpace(first), delim) {
		return nil, nil, errors.New("missing frontmatter: file must start with ---")
	}
	for off := 0; off < len(rest); {
		line, _, _ := bytes.Cut(rest[off:], []byte("\n"))
		if bytes.Equal(bytes.TrimSpace(line), delim) {
			end := min(off+len(line)+1, len(rest))
			return rest[:off], rest[end:], nil
		}
		off += len(line) + 1
	}
	return nil, nil, errors.New("unterminated frontmatter: no closing ---")
}

// decodeStrict decodes YAML, rejecting unknown keys so typos fail loudly.
func decodeStrict(front []byte, v any) error {
	dec := yaml.NewDecoder(bytes.NewReader(front))
	dec.KnownFields(true)
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("frontmatter: %w", err)
	}
	return nil
}

// parseTime parses an RFC 3339 timestamp, which must carry a timezone.
func parseTime(key, value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("frontmatter: %s is required", key)
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("frontmatter: %s %q is not RFC 3339 with timezone", key, value)
	}
	return t, nil
}

// parseDate parses an optional calendar date (YYYY-MM-DD) as midnight UTC;
// an empty value is the zero time.
func parseDate(key, value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("frontmatter: %s %q is not a date (YYYY-MM-DD)", key, value)
	}
	return t, nil
}

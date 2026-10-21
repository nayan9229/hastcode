// Package vercel writes the Vercel Build Output API v3 directory layout:
// static files under <out>/static/ and the deployment config at
// <out>/config.json. See docs/research/vercel-go-ssr-viability.md §4 and
// https://vercel.com/docs/build-output-api.
package vercel

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config is the Build Output API config.json.
type Config struct {
	Version int     `json:"version"`
	Routes  []Route `json:"routes,omitempty"`
}

// Route is one entry of config.json "routes": a Source route (Src set) or a
// Handler route (Handle set) that starts a routing phase. See
// https://vercel.com/docs/build-output-api/configuration#routes. Only the
// fields the builder uses are modelled.
type Route struct {
	Src           string            `json:"src,omitempty"`
	Dest          string            `json:"dest,omitempty"`
	Headers       map[string]string `json:"headers,omitempty"`
	Status        int               `json:"status,omitempty"`
	Continue      bool              `json:"continue,omitempty"`
	CaseSensitive bool              `json:"caseSensitive,omitempty"`
	Handle        string            `json:"handle,omitempty"`
}

// HandleFilesystem starts the phase that runs after static files and
// functions have been checked and none matched.
const HandleFilesystem = "filesystem"

// Output is a Build Output API directory being written.
type Output struct {
	Dir string
}

// NewOutput returns an Output rooted at dir.
func NewOutput(dir string) *Output { return &Output{Dir: dir} }

// Reset removes previously built static files and config so stale pages
// never survive a rebuild. Other entries (e.g. functions/) are left alone.
func (o *Output) Reset() error {
	if err := os.RemoveAll(filepath.Join(o.Dir, "static")); err != nil {
		return err
	}
	if err := os.Remove(filepath.Join(o.Dir, "config.json")); err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(filepath.Join(o.Dir, "static"), 0o755)
}

// StaticFile returns the absolute file path for a site path. A path ending
// in "/" is a directory index: "/a/b/" -> static/a/b/index.html.
func (o *Output) StaticFile(urlPath string) (string, error) {
	if !strings.HasPrefix(urlPath, "/") || strings.Contains(urlPath, "..") {
		return "", fmt.Errorf("invalid site path %q", urlPath)
	}
	p := filepath.Join(o.Dir, "static", filepath.FromSlash(urlPath[1:]))
	if strings.HasSuffix(urlPath, "/") {
		p = filepath.Join(p, "index.html")
	}
	return p, nil
}

// WriteStatic writes data to be served at urlPath.
func (o *Output) WriteStatic(urlPath string, data []byte) error {
	p, err := o.StaticFile(urlPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

// WriteConfig writes config.json.
func (o *Output) WriteConfig(c Config) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(o.Dir, "config.json"), append(data, '\n'), 0o644)
}

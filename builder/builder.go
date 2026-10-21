// Package builder turns a content tree into a Vercel Build Output API v3
// directory. It is the site's single deep module: callers pass a Config, a
// content directory and an output directory, and get back a Result listing
// errors (the build failed) and warnings. Everything else — Markdown,
// templ rendering, images, feeds — is internal.
//
// See ADR-0002: Go is a build tool here, not a server.
package builder

import (
	"fmt"
	"strings"
	"time"
)

// Mode selects production or preview output (drafts differ between them).
type Mode string

const (
	Production Mode = "production"
	Preview    Mode = "preview"
)

// Vertical is one top-level subject area every Explainer belongs to.
type Vertical struct {
	Slug        string `yaml:"slug"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Order       int    `yaml:"order"`
}

// Config is everything the builder needs besides the content and output
// directories. LoadConfig reads it from the site YAML file.
type Config struct {
	// SiteURL is the absolute origin, e.g. "https://hastcode.com".
	SiteURL   string     `yaml:"site_url"`
	Verticals []Vertical `yaml:"verticals"`
	// Mode is set by the caller (CLI flag), not the config file.
	Mode Mode `yaml:"-"`
	// Now is the build's clock for date rules (Review age, news windows).
	// Zero means time.Now(); tests pin it.
	Now time.Time `yaml:"-"`
}

// Problem is one build error or warning, attributed to a content file
// (relative to the content directory) and the rule it broke.
type Problem struct {
	File    string
	Rule    string
	Message string
}

func (p Problem) Error() string {
	if p.File == "" {
		return fmt.Sprintf("%s: %s", p.Rule, p.Message)
	}
	return fmt.Sprintf("%s: %s: %s", p.File, p.Rule, p.Message)
}

// Result is what a build reports. A build with any Errors has failed; pages
// whose content produced an error are not written.
type Result struct {
	Errors   []Problem
	Warnings []Problem
}

// Failed reports whether the build produced any error.
func (r Result) Failed() bool { return len(r.Errors) > 0 }

func (r Result) String() string {
	var b strings.Builder
	for _, p := range r.Errors {
		fmt.Fprintf(&b, "error: %s\n", p.Error())
	}
	for _, p := range r.Warnings {
		fmt.Fprintf(&b, "warning: %s\n", p.Error())
	}
	return b.String()
}

// Build reads contentDir, renders the site per cfg, and writes it to outDir
// in Build Output API v3 layout (outDir/static/..., outDir/config.json).
func Build(cfg Config, contentDir, outDir string) Result {
	if cfg.Mode == "" {
		cfg.Mode = Production
	}
	if cfg.Now.IsZero() {
		cfg.Now = time.Now()
	}
	s := newSite(cfg, contentDir, outDir)
	s.run()
	return s.result
}

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
	"os"
	"path/filepath"
	"slices"
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
	// Images configures hero image encoding.
	Images ImageConfig `yaml:"images"`
	// AdsEnabled turns every ad slot from an empty placeholder into an
	// AdSense unit. Off in v1.
	AdsEnabled bool `yaml:"ads_enabled"`
	// AdSenseClient is the AdSense publisher id ("ca-pub-..."), required
	// when AdsEnabled.
	AdSenseClient string `yaml:"adsense_client"`
	// PlausibleDomain is the Plausible analytics site; blank omits the
	// analytics script.
	PlausibleDomain string `yaml:"plausible_domain"`
	// IndexNowKey is the site's IndexNow key (indexnow.org: 8 to 128
	// characters of a-z, A-Z, 0-9 and "-"). It is public by design (the
	// build serves it at /<key>.txt), so it lives in site.yaml, and it
	// must stay the same from build to build. Blank turns IndexNow off.
	IndexNowKey string `yaml:"indexnow_key"`
	// SameAs are the publication's own profiles elsewhere (the home
	// page's Organization sameAs).
	SameAs []string `yaml:"same_as"`
	// Mode is set by the caller (CLI flag), not the config file.
	Mode Mode `yaml:"-"`
	// Now is the build's clock for date rules (Review age, news windows).
	// Zero means time.Now(); tests pin it.
	Now time.Time `yaml:"-"`

	// Build tooling, set by the caller (CLI flags, tests), not site.yaml.

	// ToolsDir holds package.json and node_modules (Node build tools such
	// as the Mermaid renderer). Default ".".
	ToolsDir string `yaml:"-"`
	// Node is the node executable. Default "node" (from PATH).
	Node string `yaml:"-"`
	// CacheDir is the build cache root (rendered diagrams, …). Default
	// <ToolsDir>/node_modules/.cache/hastcode: Vercel persists only
	// node_modules/** between builds.
	CacheDir string `yaml:"-"`
	// SkipSearchIndex leaves out the Pagefind post-step (search index under
	// /pagefind/). Allowed in preview builds only, for quick local runs;
	// a production build that skips it fails (rule search-index).
	SkipSearchIndex bool `yaml:"-"`
	// Getenv reads the Image Store credentials (imagestore.Env*) in r2
	// mode. Default os.Getenv; tests pass a fake environment.
	Getenv func(key string) string `yaml:"-"`
}

// DefaultImageCacheDir is where encoded hero variants are cached between
// builds when ImageConfig.CacheDir is empty, relative to the working
// directory. Vercel restores only node_modules/** between builds, so the
// cache lives there.
const DefaultImageCacheDir = "node_modules/.cache/hastcode/images"

// ImageConfig configures the hero image pipeline.
type ImageConfig struct {
	// CacheDir holds encoded variants keyed on the source bytes and encoder
	// settings, so unchanged heroes are not re-encoded. Deployment must
	// persist it between builds. Default DefaultImageCacheDir.
	CacheDir string `yaml:"cache_dir"`
	// Fast uses low-effort encoder settings (bigger files, much quicker
	// encodes) for tests and local previews. Never set it for production.
	Fast bool `yaml:"fast"`
	// Store is where encoded images (hero variants and crops, Author
	// photos, the OG badge) go: StoreLocal (default) writes them under
	// /assets/ in the output; StoreR2 uploads them to the Image Store and
	// every page references them by absolute URL on BaseURL (ADR-0004).
	Store ImageStore `yaml:"store"`
	// BaseURL is the Image Store's public URL, e.g.
	// "https://images.hastcode.com/". Required in r2 mode.
	BaseURL string `yaml:"base_url"`
	// Bucket is the R2 bucket. Required in r2 mode.
	Bucket string `yaml:"bucket"`
}

// ImageStore is ImageConfig.Store.
type ImageStore string

const (
	StoreLocal ImageStore = "local"
	StoreR2    ImageStore = "r2"
)

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
	// ImageKeys are the Image Store keys of every image the build
	// published (sorted, unique), the same in local and r2 mode: the set
	// the prune command keeps.
	ImageKeys []string
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
	if cfg.Images.CacheDir == "" {
		cfg.Images.CacheDir = DefaultImageCacheDir
	}
	if cfg.Now.IsZero() {
		cfg.Now = time.Now()
	}
	if cfg.ToolsDir == "" {
		cfg.ToolsDir = "."
	}
	if cfg.Node == "" {
		cfg.Node = "node"
	}
	if cfg.CacheDir == "" {
		cfg.CacheDir = filepath.Join(cfg.ToolsDir, "node_modules", ".cache", "hastcode")
	}
	if cfg.Images.Store == "" {
		cfg.Images.Store = StoreLocal
	}
	if cfg.Getenv == nil {
		cfg.Getenv = os.Getenv
	}
	store, err := openImageStore(cfg)
	if err != nil {
		return Result{Errors: []Problem{{Rule: ruleImageStore, Message: err.Error()}}}
	}
	s := newSite(cfg, contentDir, outDir, store)
	s.run()
	slices.Sort(s.result.ImageKeys)
	s.result.ImageKeys = slices.Compact(s.result.ImageKeys)
	return s.result
}

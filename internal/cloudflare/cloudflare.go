// Package cloudflare writes the Production Target's manifest (ADR-0004):
// a Workers Static Assets config, wrangler.json, at the root of the build
// output beside Vercel's config.json, and the _headers and _redirects files
// Workers reads from the assets directory. The static tree itself is
// written by package vercel; this package only adds Cloudflare's files.
//
// Deploy with `wrangler deploy --config <out>/wrangler.json`: Wrangler only
// finds a config by itself in the working directory and its parents, and
// resolves relative paths in it (main, assets.directory) against the
// config file's own directory. See docs/research/cloudflare-r2-and-deploy.md
// §5 and https://developers.cloudflare.com/workers/static-assets/.
package cloudflare

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// File names. ConfigFile sits at the output root; HeadersFile and
// RedirectsFile sit at the root of the assets directory, where Workers
// parses them and never serves them.
const (
	ConfigFile    = "wrangler.json"
	HeadersFile   = "_headers"
	RedirectsFile = "_redirects"
)

// Config is wrangler.json. Only the keys the site uses are modelled; see
// https://developers.cloudflare.com/workers/wrangler/configuration/.
type Config struct {
	Name              string `json:"name"`
	Main              string `json:"main"`
	CompatibilityDate string `json:"compatibility_date"`
	// WorkersDev keeps <name>.<account>.workers.dev; PreviewURLs enables
	// per-version and per-branch preview URLs.
	WorkersDev  bool   `json:"workers_dev"`
	PreviewURLs bool   `json:"preview_urls"`
	Assets      Assets `json:"assets"`
}

// Assets is wrangler.json "assets": Workers Static Assets.
type Assets struct {
	// Directory is relative to the config file.
	Directory string `json:"directory"`
	// Binding names env.<Binding>, the Worker's handle on the assets.
	Binding string `json:"binding"`
	// HTMLHandling: "auto-trailing-slash", "force-trailing-slash",
	// "drop-trailing-slash" or "none".
	HTMLHandling string `json:"html_handling"`
	// NotFoundHandling: "none", "404-page" or "single-page-application".
	NotFoundHandling string `json:"not_found_handling"`
	// RunWorkerFirst lists the only paths that invoke the Worker before
	// (instead of) a static asset. Everything else is a free asset hit.
	RunWorkerFirst []string `json:"run_worker_first"`
}

// WriteConfig writes c to <outDir>/wrangler.json.
func WriteConfig(outDir string, c Config) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outDir, ConfigFile), append(data, '\n'), 0o644)
}

// RemoveConfig deletes <outDir>/wrangler.json if present, so a stale config
// never survives a rebuild.
func RemoveConfig(outDir string) error {
	err := os.Remove(filepath.Join(outDir, ConfigFile))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

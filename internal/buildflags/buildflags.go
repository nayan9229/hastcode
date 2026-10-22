// Package buildflags holds the command-line flags cmd/builder and cmd/prune
// share: where the content, the site config and the Node tools are, and the
// build and image caches. Both commands build the same site, so they must
// read these the same way.
package buildflags

import (
	"flag"

	"hastcode.com/site/builder"
)

// Flags are the shared flags' values after parsing.
type Flags struct {
	Content, Config, Tools, Cache, ImageCache string
}

// Register defines -content, -config, -tools, -cache and -image-cache on fs.
func Register(fs *flag.FlagSet) *Flags {
	f := &Flags{}
	fs.StringVar(&f.Content, "content", "content", "content directory")
	fs.StringVar(&f.Config, "config", "site.yaml", "site config YAML")
	fs.StringVar(&f.Tools, "tools", ".", "directory with package.json and node_modules (run npm ci there)")
	fs.StringVar(&f.Cache, "cache", "", "build cache directory (default <tools>/node_modules/.cache/hastcode)")
	fs.StringVar(&f.ImageCache, "image-cache", "", "hero image cache dir (overrides images.cache_dir; default "+builder.DefaultImageCacheDir+")")
	return f
}

// LoadConfig loads -config and applies -tools, -cache and -image-cache.
func (f *Flags) LoadConfig() (builder.Config, error) {
	cfg, err := builder.LoadConfig(f.Config)
	if err != nil {
		return cfg, err
	}
	cfg.ToolsDir, cfg.CacheDir = f.Tools, f.Cache
	if f.ImageCache != "" {
		cfg.Images.CacheDir = f.ImageCache
	}
	return cfg, nil
}

package builder

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig reads the site config YAML (site.yaml). Unknown keys are an
// error. Mode is not read from the file; the caller sets it.
func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("%s: %w", path, err)
	}
	if cfg.SiteURL == "" {
		return Config{}, fmt.Errorf("%s: site_url is required", path)
	}
	if len(cfg.Verticals) == 0 {
		return Config{}, fmt.Errorf("%s: verticals list is empty", path)
	}
	return cfg, nil
}

// ParseMode parses a mode flag value.
func ParseMode(s string) (Mode, error) {
	switch m := Mode(s); m {
	case Production, Preview:
		return m, nil
	}
	return "", fmt.Errorf("mode %q must be %q or %q", s, Production, Preview)
}

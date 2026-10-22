// Package deploy is the deploy-adapter rule (ADR-0004): which files of the
// one build output belong to which Target. The builder writes a single
// output directory (default .vercel/output) holding the static tree plus
// both Targets' manifests; before a deploy, the Target's adapter removes
// the other Target's files so neither host ever serves the other's
// configuration.
//
//	Vercel (QA Target):            everything except wrangler.json,
//	                               static/_headers, static/_redirects;
//	                               static/robots.txt overlaid to disallow all
//	Cloudflare (Production Target): everything except config.json
//
// Everything else (the site under static/) is shared. An overlay replaces a
// shared file's bytes for one Target: Vercel is only ever the QA Target, so
// its robots.txt disallows every crawler (the manifest adds
// X-Robots-Tag: noindex too). Deploy steps call
// Apply through cmd/adapt (`make adapt TARGET=…`) once, after the build and
// before the upload; after that the output belongs to that Target only.
package deploy

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"hastcode.com/site/internal/cloudflare"
	"hastcode.com/site/internal/vercel"
)

// Target is a host a build can be deployed to (GLOSSARY "Target").
type Target string

const (
	// Cloudflare is the Production Target: Workers Static Assets, deployed
	// with `wrangler deploy --config <out>/wrangler.json`. Tag PUBCF-*.
	Cloudflare Target = "cloudflare"
	// Vercel is the QA Target: Build Output API, deployed with
	// `vercel deploy --prebuilt`. Tag PUBVL-*.
	Vercel Target = "vercel"
)

// Targets lists every Target.
var Targets = []Target{Cloudflare, Vercel}

// own lists, per Target, the output-relative (slash-separated) files only
// that Target reads. A file belongs to every Target except those that do
// not own it while another Target does.
var own = map[Target][]string{
	Cloudflare: {
		cloudflare.ConfigFile,
		"static/" + cloudflare.HeadersFile,
		"static/" + cloudflare.RedirectsFile,
	},
	Vercel: {
		vercel.ConfigFile,
	},
}

// RobotsFile is the site's robots.txt, output-relative.
const RobotsFile = "static/robots.txt"

// DisallowAllRobots is the QA Target's robots.txt: no crawler may fetch
// anything.
const DisallowAllRobots = "User-agent: *\nDisallow: /\n"

// overlays lists, per Target, files (output-relative, slash-separated)
// whose bytes Apply replaces for that Target. Each path is part of the
// Target's file set whether or not the build wrote it.
var overlays = map[Target]map[string]string{
	Vercel: {RobotsFile: DisallowAllRobots},
}

// ParseTarget parses a Target name ("cloudflare" or "vercel").
func ParseTarget(s string) (Target, error) {
	if t := Target(s); slices.Contains(Targets, t) {
		return t, nil
	}
	return "", fmt.Errorf("target %q must be %q or %q", s, Cloudflare, Vercel)
}

// foreign reports whether rel (output-relative, slash-separated) belongs
// only to Targets other than t.
func foreign(t Target, rel string) bool {
	for other, files := range own {
		if other != t && slices.Contains(files, rel) {
			return true
		}
	}
	return false
}

// Files returns t's file set: every file under outDir (output-relative,
// slash-separated, sorted) except those only another Target reads.
func Files(outDir string, t Target) ([]string, error) {
	if _, err := ParseTarget(string(t)); err != nil {
		return nil, err
	}
	var files []string
	err := filepath.WalkDir(outDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(outDir, p)
		if err != nil {
			return err
		}
		if rel = filepath.ToSlash(rel); !foreign(t, rel) && overlays[t][rel] == "" {
			files = append(files, rel)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	for rel := range overlays[t] {
		files = append(files, rel)
	}
	slices.Sort(files)
	return files, nil
}

// Apply removes from outDir every file only another Target reads and writes
// t's overlays, leaving exactly Files(outDir, t). It returns the removed
// paths. Applying twice is harmless.
func Apply(outDir string, t Target) (removed []string, err error) {
	if _, err := ParseTarget(string(t)); err != nil {
		return nil, err
	}
	for other, files := range own {
		if other == t {
			continue
		}
		for _, rel := range files {
			if slices.Contains(own[t], rel) {
				continue
			}
			err := os.Remove(filepath.Join(outDir, filepath.FromSlash(rel)))
			switch {
			case err == nil:
				removed = append(removed, rel)
			case !os.IsNotExist(err):
				return removed, err
			}
		}
	}
	slices.Sort(removed)
	for rel, data := range overlays[t] {
		p := filepath.Join(outDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return removed, err
		}
		if err := os.WriteFile(p, []byte(data), 0o644); err != nil {
			return removed, err
		}
	}
	return removed, nil
}

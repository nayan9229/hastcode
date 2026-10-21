package buildtest

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/builder"
)

// Build writes tree to a temp content dir, runs the builder with cfg into a
// fresh temp output dir, and returns the output dir and the builder's Result.
//
// Builds get a fresh diagram cache unless cfg.CacheDir is set. A tree with a
// Mermaid fence skips the test when the renderer (node + npm packages) is
// not installed.
func Build(tb testing.TB, cfg builder.Config, tree Tree) (outDir string, res builder.Result) {
	tb.Helper()
	if tree.hasDiagram() {
		requireDiagramRenderer(tb, cfg)
	}
	if cfg.CacheDir == "" {
		cfg.CacheDir = tb.TempDir()
	}
	contentDir := tree.Write(tb)
	outDir = filepath.Join(tb.TempDir(), "output")
	res = builder.Build(cfg, contentDir, outDir)
	return outDir, res
}

// MustBuild is Build that fails the test if the builder reports any error.
func MustBuild(tb testing.TB, cfg builder.Config, tree Tree) string {
	tb.Helper()
	outDir, res := Build(tb, cfg, tree)
	if res.Failed() {
		tb.Fatalf("build failed:\n%s", res)
	}
	return outDir
}

// StaticPath maps a site URL path ("/vertical/slug/") to the file the
// builder must write under the Build Output API static/ dir.
func StaticPath(outDir, urlPath string) string {
	p := filepath.Join(outDir, "static", filepath.FromSlash(strings.TrimPrefix(urlPath, "/")))
	if strings.HasSuffix(urlPath, "/") {
		p = filepath.Join(p, "index.html")
	}
	return p
}

// Exists reports whether a file was written for urlPath.
func Exists(outDir, urlPath string) bool {
	_, err := os.Stat(StaticPath(outDir, urlPath))
	return err == nil
}

// Page loads and parses the HTML written for urlPath.
func Page(tb testing.TB, outDir, urlPath string) *goquery.Document {
	tb.Helper()
	data, err := os.ReadFile(StaticPath(outDir, urlPath))
	if err != nil {
		tb.Fatalf("page %s: %v", urlPath, err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		tb.Fatalf("parse %s: %v", urlPath, err)
	}
	return doc
}

// JSONLD decodes every <script type="application/ld+json"> block in doc.
func JSONLD(tb testing.TB, doc *goquery.Document) []map[string]any {
	tb.Helper()
	var blocks []map[string]any
	doc.Find(`script[type="application/ld+json"]`).Each(func(_ int, s *goquery.Selection) {
		var v map[string]any
		if err := json.Unmarshal([]byte(s.Text()), &v); err != nil {
			tb.Fatalf("JSON-LD block does not decode: %v\n%s", err, s.Text())
		}
		blocks = append(blocks, v)
	})
	return blocks
}

// JSONLDOfType returns the first JSON-LD block whose @type is typ, or nil.
func JSONLDOfType(tb testing.TB, doc *goquery.Document, typ string) map[string]any {
	tb.Helper()
	for _, b := range JSONLD(tb, doc) {
		if b["@type"] == typ {
			return b
		}
	}
	return nil
}

// FilesContaining returns the output files (relative to outDir, slash
// separated) whose bytes contain s. Use it to prove a piece appears nowhere:
// no page, listing, feed, sitemap or manifest.
func FilesContaining(tb testing.TB, outDir, s string) []string {
	tb.Helper()
	var hits []string
	err := filepath.WalkDir(outDir, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if bytes.Contains(data, []byte(s)) {
			rel, _ := filepath.Rel(outDir, p)
			hits = append(hits, filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		tb.Fatalf("walk %s: %v", outDir, err)
	}
	return hits
}

// VercelConfig reads and decodes the Build Output API config.json.
func VercelConfig(tb testing.TB, outDir string) map[string]any {
	tb.Helper()
	data, err := os.ReadFile(filepath.Join(outDir, "config.json"))
	if err != nil {
		tb.Fatalf("config.json: %v", err)
	}
	var v map[string]any
	if err := json.Unmarshal(data, &v); err != nil {
		tb.Fatalf("config.json does not decode: %v", err)
	}
	return v
}

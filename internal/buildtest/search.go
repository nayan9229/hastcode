package buildtest

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"hastcode.com/site/builder"
)

// SearchIndex decodes the Pagefind index in the output and returns the
// indexed text of every page by its URL ("/vertical/slug/"). Each page is
// one gzipped fragment file holding "pagefind_dcd" + JSON.
func SearchIndex(tb testing.TB, outDir string) map[string]string {
	tb.Helper()
	files, err := filepath.Glob(filepath.Join(outDir, "static", "pagefind", "fragment", "*.pf_fragment"))
	if err != nil {
		tb.Fatal(err)
	}
	if len(files) == 0 {
		tb.Fatalf("no Pagefind fragments in %s", outDir)
	}
	pages := map[string]string{}
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			tb.Fatal(err)
		}
		zr, err := gzip.NewReader(bytes.NewReader(raw))
		if err != nil {
			tb.Fatalf("%s: %v", f, err)
		}
		data, err := io.ReadAll(zr)
		if err != nil {
			tb.Fatalf("%s: %v", f, err)
		}
		var frag struct {
			URL     string `json:"url"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal(bytes.TrimPrefix(data, []byte("pagefind_dcd")), &frag); err != nil {
			tb.Fatalf("%s: %v", f, err)
		}
		pages[frag.URL] = frag.Content
	}
	return pages
}

// requireSearchIndexer skips the test when Pagefind cannot run here: no
// node on PATH, or `npm ci` not run in cfg.ToolsDir. Every production
// build runs it (Config.SkipSearchIndex is for preview builds only).
func requireSearchIndexer(tb testing.TB, cfg builder.Config) {
	tb.Helper()
	if _, err := exec.LookPath("node"); err != nil {
		tb.Skip("skipping: production builds index search with Pagefind and `node` is not on PATH (install Node 22+, then `npm ci` at the repo root)")
	}
	pkg := filepath.Join(cfg.ToolsDir, "node_modules", "pagefind", "package.json")
	if _, err := os.Stat(pkg); err != nil {
		tb.Skipf("skipping: production builds index search with Pagefind and %s is missing (run `npm ci` at the repo root)", pkg)
	}
}

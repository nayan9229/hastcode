// Package pagefind builds the site's search index by running the Pagefind
// CLI (npm package `pagefind`, pinned in package.json) over a directory of
// built HTML. It knows nothing of pages or URLs: which parts of a page are
// indexed is decided by data-pagefind-* attributes in the markup.
package pagefind

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Package is the npm package that provides the CLI.
const Package = "pagefind"

// BundleDir is the directory Pagefind writes into, relative to the site
// root; the browser loads the bundle from "/" + BundleDir + "/".
const BundleDir = "pagefind"

// EntryFile is the bundle's manifest, written last: it exists only when
// indexing finished.
const EntryFile = "pagefind-entry.json"

// Indexer runs the Pagefind CLI through Node.
type Indexer struct {
	// Node is the node executable.
	Node string
	// ToolsDir holds node_modules/pagefind (run `npm ci` there); the CLI
	// runs with it as the working directory.
	ToolsDir string
}

// Index indexes every HTML file under siteDir, writing the bundle to
// siteDir/BundleDir, and returns how many pages the index holds (all
// languages).
func (ix Indexer) Index(ctx context.Context, siteDir string) (pages int, err error) {
	bin := filepath.Join("node_modules", Package, "lib", "runner", "bin.cjs")
	if _, err := os.Stat(filepath.Join(ix.ToolsDir, bin)); err != nil {
		return 0, fmt.Errorf("pagefind is not installed in %s (run `npm ci` there): %w", ix.ToolsDir, err)
	}
	site, err := filepath.Abs(siteDir)
	if err != nil {
		return 0, err
	}
	var outBuf bytes.Buffer
	cmd := exec.CommandContext(ctx, ix.Node, bin, "--site", site, "--output-subdir", BundleDir, "--quiet")
	cmd.Dir = ix.ToolsDir
	cmd.Stdout, cmd.Stderr = &outBuf, &outBuf
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("%s %s in %s: %v\n%s", ix.Node, bin, ix.ToolsDir, err, strings.TrimSpace(outBuf.String()))
	}
	return Pages(siteDir)
}

// Pages reads the bundle manifest under siteDir and returns how many pages
// the index holds, or an error if there is no complete bundle.
func Pages(siteDir string) (int, error) {
	data, err := os.ReadFile(filepath.Join(siteDir, BundleDir, EntryFile))
	if err != nil {
		return 0, err
	}
	var entry struct {
		Languages map[string]struct {
			PageCount int `json:"page_count"`
		} `json:"languages"`
	}
	if err := json.Unmarshal(data, &entry); err != nil {
		return 0, fmt.Errorf("%s: %w", EntryFile, err)
	}
	n := 0
	for _, l := range entry.Languages {
		n += l.PageCount
	}
	return n, nil
}

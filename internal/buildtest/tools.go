package buildtest

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"hastcode.com/site/builder"
)

// RepoRoot is the repository root (the directory holding go.mod and
// package.json), found by walking up from the test's working directory.
var RepoRoot = sync.OnceValue(func() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("buildtest: go.mod not found above the working directory")
		}
		dir = parent
	}
})

func (t Tree) hasDiagram() bool {
	for _, data := range t {
		if bytes.Contains(data, []byte("```mermaid")) {
			return true
		}
	}
	return false
}

// requireDiagramRenderer skips the test when the Mermaid renderer cannot run
// here: no node on PATH, or `npm ci` not run in cfg.ToolsDir.
func requireDiagramRenderer(tb testing.TB, cfg builder.Config) {
	tb.Helper()
	if _, err := exec.LookPath("node"); err != nil {
		tb.Skip("skipping: fixture has a Mermaid diagram and `node` is not on PATH (install Node 22+, then `npm ci` at the repo root)")
	}
	pkg := filepath.Join(cfg.ToolsDir, "node_modules", "beautiful-mermaid", "package.json")
	if _, err := os.Stat(pkg); err != nil {
		tb.Skipf("skipping: fixture has a Mermaid diagram and %s is missing (run `npm ci` at the repo root)", pkg)
	}
}

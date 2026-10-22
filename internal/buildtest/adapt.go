package buildtest

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"
)

// Adapt runs the deploy adapter on outDir the way a deploy does (`make
// adapt TARGET=…`: the cmd/adapt command) and fails the test unless it
// exits 0.
func Adapt(tb testing.TB, outDir, target string) {
	tb.Helper()
	if code, output := RunAdapt(tb, outDir, target); code != 0 {
		tb.Fatalf("adapt -target %q exited %d:\n%s", target, code, output)
	}
}

// RunAdapt runs the deploy adapter command on outDir and returns its exit
// status and combined output, for tests of refusals.
func RunAdapt(tb testing.TB, outDir, target string) (exitCode int, output string) {
	tb.Helper()
	// Built rather than `go run`, which turns every failing status into 1.
	bin := filepath.Join(tb.TempDir(), "adapt")
	build := exec.Command("go", "build", "-o", bin, "./cmd/adapt")
	build.Dir = RepoRoot()
	if out, err := build.CombinedOutput(); err != nil {
		tb.Fatalf("build cmd/adapt: %v\n%s", err, out)
	}
	cmd := exec.Command(bin, "-target", target, "-out", outDir)
	var buf bytes.Buffer
	cmd.Stdout, cmd.Stderr = &buf, &buf
	err := cmd.Run()
	if exit := (*exec.ExitError)(nil); errors.As(err, &exit) {
		return exit.ExitCode(), buf.String()
	}
	if err != nil {
		tb.Fatalf("run adapt: %v\n%s", err, buf.String())
	}
	return 0, buf.String()
}

// OutputFiles lists every file under outDir, output-relative with "/"
// separators, sorted: what a deploy would upload.
func OutputFiles(tb testing.TB, outDir string) []string {
	tb.Helper()
	var files []string
	err := filepath.WalkDir(outDir, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(outDir, p)
		files = append(files, filepath.ToSlash(rel))
		return err
	})
	if err != nil {
		tb.Fatalf("walk %s: %v", outDir, err)
	}
	slices.Sort(files)
	return files
}

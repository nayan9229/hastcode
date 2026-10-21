package builder_test

import (
	"os"
	"strings"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

// The spec's weight budget: no page's HTML may exceed 50 KB uncompressed.
const htmlSizeLimit = 50 * 1024

func TestSizeGate_PageOver50KBFailsTheBuildNamingThePage(t *testing.T) {
	e := buildtest.ValidExplainer()
	padding := strings.Repeat("Leaders send heartbeats so followers do not start elections. ", 1000) // ~62 KB
	e.Body = strings.Replace(e.Body, "## Sources", "## Heartbeats keep followers quiet\n\n"+padding+"\n\n## Sources", 1)

	out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

	buildtest.RequireErrorRule(t, res, "", "html-size", explainerPath)
	buildtest.RequireErrorRule(t, res, "", "html-size", "50 KB")
	if buildtest.Exists(out, explainerPath) {
		t.Errorf("oversized page %s was written; want it left out of the output", explainerPath)
	}
}

func TestSizeGate_NormalPagePasses(t *testing.T) {
	out, res := buildtest.Build(t, buildtest.Config(), buildtest.ValidTree())

	buildtest.RequireNoProblemRule(t, res, "", "html-size")
	info, err := os.Stat(buildtest.StaticPath(out, explainerPath))
	if err != nil {
		t.Fatalf("normal page not written: %v", err)
	}
	if info.Size() > htmlSizeLimit {
		t.Errorf("fixture page is %d bytes; the normal-page fixture must stay under the limit", info.Size())
	}
}

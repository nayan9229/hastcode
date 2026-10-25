package builder_test

import (
	"strings"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

// The spec's weight budget: no page's HTML may exceed 50 KB uncompressed,
// not counting the shared inline stylesheet (ADR-0006).
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
	if size := buildtest.ContentSize(t, out, explainerPath); size > htmlSizeLimit {
		t.Errorf("fixture page is %d bytes without the stylesheet; the normal-page fixture must stay under the limit", size)
	}
}

// The inline stylesheet is the same on every page and belongs to the
// design, not the content: a page whose content is just under the limit
// passes even though stylesheet plus content is over it.
func TestSizeGate_InlineStylesheetDoesNotCount(t *testing.T) {
	base := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	content := buildtest.ContentSize(t, base, explainerPath)
	style := buildtest.FileSize(t, base, explainerPath) - content
	if style < 1000 {
		t.Fatalf("inline stylesheet is %d bytes; this test needs it present", style)
	}

	// Pad the content to 1 KB under the limit; the file is then over it.
	sentence := "Leaders send heartbeats so followers do not start elections. "
	need := int(htmlSizeLimit-1024-content) / len(sentence)
	e := buildtest.ValidExplainer()
	e.Body = strings.Replace(e.Body, "## Sources", "## Heartbeats keep followers quiet\n\n"+strings.Repeat(sentence, need)+"\n\n## Sources", 1)
	out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

	buildtest.RequireNoProblemRule(t, res, "", "html-size")
	if total := buildtest.FileSize(t, out, explainerPath); total <= htmlSizeLimit {
		t.Fatalf("padded page is %d bytes in total; want it over the limit so the test means something", total)
	}
	if got := buildtest.ContentSize(t, out, explainerPath); got > htmlSizeLimit {
		t.Errorf("padded page content is %d bytes; want it under the limit", got)
	}
}

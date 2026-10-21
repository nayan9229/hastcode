package builder_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildtest"
)

const explainerURL = "/distributed-systems/raft-leader-election/"
const explainerFile = "explainers/distributed-systems/raft-leader-election.md"

func mermaid(lines ...string) string {
	return "```mermaid\n" + strings.Join(lines, "\n") + "\n```\n"
}

// explainerWithDiagrams is ValidExplainer with a section holding extra
// Mermaid fences inserted before Sources (after the fence it already has).
func explainerWithDiagrams(fences ...string) buildtest.Explainer {
	e := buildtest.ValidExplainer()
	section := "## Log replication follows the election\n\nOnce elected, the leader appends entries.\n\n" + strings.Join(fences, "\n") + "\n"
	e.Body = strings.Replace(e.Body, "## Sources", section+"## Sources", 1)
	return e
}

func TestDiagram_EachMermaidFenceBecomesAFigureWithInlineSVGAndCaption(t *testing.T) {
	e := explainerWithDiagrams(mermaid(
		"%% title: Appending an entry",
		"sequenceDiagram",
		"  Client->>Leader: write x=1",
		"  Leader->>Follower: AppendEntries",
	))
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))
	body := buildtest.Page(t, out, explainerURL).Find(".e-content")

	figures := body.Find("figure")
	if figures.Length() != 2 {
		t.Fatalf("figures = %d, want 2", figures.Length())
	}
	if got := body.Find("pre code.language-mermaid").Length(); got != 0 {
		t.Errorf("%d Mermaid fences left as code blocks, want 0", got)
	}

	wantCaptions := []string{"A Raft election", "Appending an entry"}
	wantLabels := [][]string{
		{"Follower", "Candidate", "Leader"},
		{"Client", "Leader", "Follower", "write x=1", "AppendEntries"},
	}
	figures.Each(func(i int, fig *goquery.Selection) {
		if got := fig.ChildrenFiltered("svg").Length(); got != 1 {
			t.Errorf("figure %d: %d <svg> children, want 1", i, got)
		}
		if got := strings.TrimSpace(fig.Find("figcaption").Text()); got != wantCaptions[i] {
			t.Errorf("figure %d caption = %q, want %q", i, got, wantCaptions[i])
		}
		texts := map[string]bool{}
		fig.Find("svg text").Each(func(_ int, s *goquery.Selection) { texts[strings.TrimSpace(s.Text())] = true })
		for _, label := range wantLabels[i] {
			if !texts[label] {
				t.Errorf("figure %d: no <text> node %q in SVG; text nodes: %v", i, label, texts)
			}
		}
	})
}

// Colours must come from CSS (currentColor or var(--…)) so the page's light
// and dark themes both apply; a literal colour would be wrong in one of them.
func TestDiagram_SVGHasNoHardCodedColoursOrFontImports(t *testing.T) {
	e := explainerWithDiagrams(mermaid(
		"%% title: Raft server states",
		"stateDiagram-v2",
		"  [*] --> Follower",
		"  Follower --> Candidate: timeout",
		"  Candidate --> Leader: majority",
	))
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))
	doc := buildtest.Page(t, out, explainerURL)

	html, _ := doc.Html()
	for _, banned := range []string{"@import", "fonts.googleapis.com"} {
		if strings.Contains(html, banned) {
			t.Errorf("page contains %q: diagrams must not pull third-party fonts", banned)
		}
	}

	svgs := doc.Find("figure svg")
	if svgs.Length() != 2 {
		t.Fatalf("svgs = %d, want 2", svgs.Length())
	}
	if got := svgs.Find("style").Length(); got != 0 {
		t.Errorf("%d <style> elements inside diagrams; page-wide rules belong in site CSS", got)
	}
	svgs.Each(func(i int, svg *goquery.Selection) {
		if svg.AttrOr("viewBox", "") == "" || svg.AttrOr("width", "") == "" || svg.AttrOr("height", "") == "" {
			t.Errorf("svg %d lacks viewBox/width/height (layout shift): %v", i, svg.Nodes[0].Attr)
		}
	})
	svgs.Find("*").AddSelection(svgs).Each(func(_ int, el *goquery.Selection) {
		for _, a := range el.Nodes[0].Attr {
			switch a.Key {
			case "fill", "stroke", "color", "stop-color", "flood-color", "lighting-color":
				if !themeable(a.Val) {
					t.Errorf("<%s %s=%q>: hard-coded colour", goquery.NodeName(el), a.Key, a.Val)
				}
			case "style":
				for _, decl := range strings.Split(a.Val, ";") {
					_, v, _ := strings.Cut(decl, ":")
					if v = strings.TrimSpace(v); v != "" && !themeable(v) {
						t.Errorf("<%s style=%q>: hard-coded colour %q", goquery.NodeName(el), a.Val, v)
					}
				}
			}
		}
	})
}

func TestDiagram_MalformedOrUnsupportedFenceIsABuildErrorNamingThePage(t *testing.T) {
	for name, tc := range map[string]struct {
		fence string
		want  string
	}{
		"misspelt header":   {mermaid("%% title: Typo", "flowchat LR", "  A --> B"), "flowchat"},
		"gantt unsupported": {mermaid("%% title: Plan", "gantt", "  title Roadmap"), `"gantt" is not supported`},
		"pie unsupported":   {mermaid("%% title: Share", "pie", `  "Reads" : 80`), `"pie" is not supported`},
		"empty":             {mermaid("%% title: Nothing"), "empty"},
		"no nodes":          {mermaid("%% title: Header only", "flowchart LR"), "no nodes"},
		"unbalanced shape":  {mermaid("%% title: Bracket", "flowchart LR", "  A[[[Client --> B"), "unbalanced"},
		"dangling edge":     {mermaid("%% title: Dangling", "flowchart LR", "  A --> B", "  B -->"), "B -->"},
	} {
		t.Run(name, func(t *testing.T) {
			out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(explainerWithDiagrams(tc.fence)))

			buildtest.RequireError(t, res, explainerFile, tc.want)
			buildtest.RequireError(t, res, explainerFile, explainerURL)
			if buildtest.Exists(out, explainerURL) {
				t.Errorf("page %s was written despite the diagram error", explainerURL)
			}
		})
	}
}

// A rebuild with a warm cache must not need the renderer for unchanged
// diagrams; a changed diagram must reach it again.
func TestDiagram_UnchangedDiagramsAreServedFromTheCacheWithoutRerendering(t *testing.T) {
	cfg := buildtest.Config()
	cfg.CacheDir = t.TempDir()
	tree := buildtest.ValidTree()

	first := buildtest.Page(t, buildtest.MustBuild(t, cfg, tree), explainerURL).Find("figure.diagram")
	if first.Length() != 1 {
		t.Fatalf("first build: %d diagram figures, want 1", first.Length())
	}
	firstHTML, _ := goquery.OuterHtml(first)

	noRenderer := cfg
	noRenderer.Node = filepath.Join(t.TempDir(), "no-such-node")
	// Without Node there is no Pagefind either: a preview build may skip
	// the search index, so only the diagram renderer is missing.
	noRenderer.Mode = builder.Preview
	noRenderer.SkipSearchIndex = true

	second := buildtest.Page(t, buildtest.MustBuild(t, noRenderer, tree), explainerURL).Find("figure.diagram")
	if secondHTML, _ := goquery.OuterHtml(second); secondHTML != firstHTML {
		t.Errorf("cached figure differs from the first render:\nfirst:  %.300s\nsecond: %.300s", firstHTML, secondHTML)
	}

	changed := buildtest.ValidExplainer()
	changed.Body = strings.Replace(changed.Body, "Candidate --> Leader", "Candidate --> Leader --> Follower", 1)
	_, res := buildtest.Build(t, noRenderer, buildtest.NewTree().WithExplainer(changed))
	if !res.Failed() || !strings.Contains(res.String(), "no-such-node") {
		t.Errorf("changed diagram with no renderer: want a build error naming the renderer, got:\n%s", res)
	}
}

func themeable(v string) bool {
	return v == "none" || v == "currentColor" || v == "transparent" ||
		strings.HasPrefix(v, "var(--") || strings.HasPrefix(v, "url(#")
}

func TestDiagram_ElementIDsAreUniqueAcrossDiagramsOnAPage(t *testing.T) {
	e := explainerWithDiagrams(mermaid("%% title: Second", "flowchart TD", "  A --> B"))
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))
	doc := buildtest.Page(t, out, explainerURL)

	seen := map[string]bool{}
	doc.Find("[id]").Each(func(_ int, el *goquery.Selection) {
		id := el.AttrOr("id", "")
		if seen[id] {
			t.Errorf("duplicate id %q on page", id)
		}
		seen[id] = true
	})
	doc.Find("figure svg [marker-end]").Each(func(_ int, el *goquery.Selection) {
		ref := strings.TrimSuffix(strings.TrimPrefix(el.AttrOr("marker-end", ""), "url(#"), ")")
		if !seen[ref] {
			t.Errorf("marker-end references #%s, which is not on the page", ref)
		}
	})
}

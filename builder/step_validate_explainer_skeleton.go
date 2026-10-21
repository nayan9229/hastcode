package builder

import (
	"slices"

	"github.com/yuin/goldmark/ast"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/markdown"
)

func init() {
	registerStep(phaseValidate, 30, "validate explainer skeleton", validateExplainerSkeleton)
}

// Rules for the Explainer editorial skeleton (spec "Content model").
const (
	ruleTLDR           = "tldr"
	ruleSources        = "sources"
	ruleDiagram        = "diagram"
	ruleTradeOffsTable = "trade-offs-table"
)

const (
	headingTLDR      = "TL;DR"
	headingSources   = "Sources"
	headingTradeOffs = "Trade-offs"
)

// validateExplainerSkeleton enforces the structural editorial floor on
// every Explainer: TL;DR, Sources, a Mermaid diagram, and (Comparisons) a
// trade-offs table. Every broken rule is reported, not just the first.
func validateExplainerSkeleton(s *site) error {
	for _, e := range s.explainers {
		checkTLDR(s, e)
		checkSources(s, e)
		checkDiagram(s, e)
		checkTradeOffsTable(s, e)
	}
	return nil
}

// checkTLDR: the first section after the intro is "## TL;DR" and opens with
// a bullet list of 3–5 items.
func checkTLDR(s *site, e *content.Explainer) {
	_, sections := e.Body.Sections()
	i := slices.IndexFunc(sections, func(sec markdown.Section) bool { return sec.Heading == headingTLDR })
	switch {
	case i < 0:
		s.fail(e.File, ruleTLDR, "missing a ## TL;DR section immediately after the intro")
		return
	case i > 0:
		s.fail(e.File, ruleTLDR, "## TL;DR must come immediately after the intro, before ## %s", sections[0].Heading)
		return
	}
	blocks := sections[0].Blocks
	list, ok := firstBlock(blocks).(*ast.List)
	if !ok || list.IsOrdered() || len(blocks) != 1 {
		s.fail(e.File, ruleTLDR, "## TL;DR must be a bullet list of 3–5 bullets and nothing else")
		return
	}
	if n := list.ChildCount(); n < 3 || n > 5 {
		s.fail(e.File, ruleTLDR, "## TL;DR has %d bullets, want 3–5 bullets", n)
	}
}

// checkSources: the last section is "## Sources" and contains a link.
func checkSources(s *site, e *content.Explainer) {
	_, sections := e.Body.Sections()
	i := slices.IndexFunc(sections, func(sec markdown.Section) bool { return sec.Heading == headingSources })
	switch {
	case i < 0:
		s.fail(e.File, ruleSources, "missing a ## Sources section at the end")
	case i != len(sections)-1:
		s.fail(e.File, ruleSources, "## Sources must be the section at the end, but ## %s follows it", sections[len(sections)-1].Heading)
	case !markdown.Find(sections[i].Blocks, markdown.IsLink):
		s.fail(e.File, ruleSources, "## Sources must contain at least one link")
	}
}

// checkDiagram: the body contains at least one ```mermaid fence.
func checkDiagram(s *site, e *content.Explainer) {
	isMermaid := func(n ast.Node) bool { return e.Body.FenceLanguage(n) == "mermaid" }
	if !markdown.Find([]ast.Node{e.Body.Root}, isMermaid) {
		s.fail(e.File, ruleDiagram, "no architecture diagram: add at least one Mermaid diagram (a ```mermaid fence)")
	}
}

// checkTradeOffsTable: a Comparison has a Markdown table under "## Trade-offs".
func checkTradeOffsTable(s *site, e *content.Explainer) {
	if e.Kind != content.Comparison {
		return
	}
	_, sections := e.Body.Sections()
	for _, sec := range sections {
		if sec.Heading == headingTradeOffs && markdown.Find(sec.Blocks, markdown.IsTable) {
			return
		}
	}
	s.fail(e.File, ruleTradeOffsTable, "a Comparison needs a Markdown table under a ## Trade-offs heading")
}

func firstBlock(blocks []ast.Node) ast.Node {
	if len(blocks) == 0 {
		return nil
	}
	return blocks[0]
}

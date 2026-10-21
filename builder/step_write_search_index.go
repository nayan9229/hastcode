package builder

import (
	"context"
	"path/filepath"

	"hastcode.com/site/internal/pagefind"
)

func init() {
	// After every page is written and the link gate has removed broken
	// pages (render phase), so the index holds exactly what is deployed.
	registerStep(phaseWrite, 900, "index search", indexSearch)
}

// ruleSearchIndex: the Pagefind index is missing or empty. Production
// builds must ship it (spec "Search"); preview builds may skip it.
const ruleSearchIndex = "search-index"

// pagefindBundleURL is where the browser loads Pagefind from: the search
// page's lazy loader reads it, and the index step writes it.
const pagefindBundleURL = "/" + pagefind.BundleDir + "/"

// indexSearch runs Pagefind over the written site (Config.ToolsDir must
// hold node_modules/pagefind). Only regions marked data-pagefind-body (the
// Explainer and Brief articles) are indexed. A build that has already
// failed is not indexed: it is never deployed.
func indexSearch(s *site) error {
	if s.cfg.SkipSearchIndex {
		if s.cfg.Mode == Production {
			s.fail("", ruleSearchIndex, "production builds must build the search index; skipping it is for preview builds only")
		} else {
			s.warn("", ruleSearchIndex, "search index skipped: %s finds nothing in this build", searchURL)
		}
		return nil
	}
	if s.result.Failed() {
		return nil
	}
	ix := pagefind.Indexer{Node: s.cfg.Node, ToolsDir: s.cfg.ToolsDir}
	pages, err := ix.Index(context.Background(), filepath.Join(s.out.Dir, "static"))
	switch {
	case err != nil:
		s.fail("", ruleSearchIndex, "building the search index: %v", err)
	case pages == 0:
		s.fail("", ruleSearchIndex, "the search index holds no pages; Explainer and Brief pages must mark their article data-pagefind-body")
	}
	return nil
}

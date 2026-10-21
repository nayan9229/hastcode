package builder

import (
	"cmp"
	"fmt"
	"slices"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/vercel"
)

// The build is a pipeline of steps grouped into phases. Each step lives in
// its own step_*.go file and registers itself from init(), so adding a step
// (a gate, image processing, diagram rendering, a new page type) never
// means editing a shared function or list:
//
//	func init() { registerStep(phaseValidate, 20, "explainer skeleton gates", checkExplainerSkeleton) }
//
// Steps run in phase order, then by their order number (ties by name).
// Leave gaps in order numbers (10, 20, 30…) so later steps can slot between.
type phase int

const (
	phaseLoad      phase = iota // read content into s (pieces, authors)
	phaseValidate               // gates: record errors/warnings, mark pieces failed
	phaseTransform              // mutate parsed content (diagrams, hero variants)
	phaseRender                 // render pages into s.out (static/)
	phaseWrite                  // write manifests (config.json, sitemaps, feeds…)
)

type step struct {
	phase phase
	order int
	name  string
	run   func(*site) error
}

var steps []step

// registerStep adds a step to the pipeline. Call it from init().
func registerStep(p phase, order int, name string, run func(*site) error) {
	steps = append(steps, step{phase: p, order: order, name: name, run: run})
}

// site is the state threaded through every step.
type site struct {
	cfg        Config
	contentDir string
	out        *vercel.Output

	verticals  map[string]Vertical // by slug
	authors    map[string]*content.Author
	explainers []*content.Explainer
	briefs     []*content.Brief
	pages      []*content.Page  // static pages; see step_static_pages.go
	heroes     map[string]*hero // by content file; see step_transform_heroes.go
	// authorPhotos are encoded Author photos by Author slug; see
	// step_author_photos.go.
	authorPhotos map[string]*authorPhoto

	failed map[string]bool // content-relative files that produced an error
	result Result
}

func newSite(cfg Config, contentDir, outDir string) *site {
	s := &site{
		cfg:        cfg,
		contentDir: contentDir,
		out:        vercel.NewOutput(outDir),
		verticals:  map[string]Vertical{},
		authors:    map[string]*content.Author{},
		failed:     map[string]bool{},
		heroes:     map[string]*hero{},
	}
	for _, v := range cfg.Verticals {
		s.verticals[v.Slug] = v
	}
	return s
}

func (s *site) run() {
	ordered := slices.Clone(steps)
	slices.SortStableFunc(ordered, func(a, b step) int {
		return cmp.Or(cmp.Compare(a.phase, b.phase), cmp.Compare(a.order, b.order), cmp.Compare(a.name, b.name))
	})
	for _, st := range ordered {
		if err := st.run(s); err != nil {
			s.fail("", st.name, "%v", err)
			return
		}
	}
}

// fail records a build error for a content file and excludes that file's
// page from output. file is relative to the content directory ("" for
// site-wide problems).
func (s *site) fail(file, rule, format string, args ...any) {
	s.result.Errors = append(s.result.Errors, Problem{File: file, Rule: rule, Message: fmt.Sprintf(format, args...)})
	if file != "" {
		s.failed[file] = true
	}
}

// warn records a warning; warnings never fail the build.
func (s *site) warn(file, rule, format string, args ...any) {
	s.result.Warnings = append(s.result.Warnings, Problem{File: file, Rule: rule, Message: fmt.Sprintf(format, args...)})
}

// publishedExplainers are the Explainers that passed every gate so far, in
// load order. Render and write steps must use this, not s.explainers.
func (s *site) publishedExplainers() []*content.Explainer {
	var ok []*content.Explainer
	for _, e := range s.explainers {
		if !s.failed[e.File] {
			ok = append(ok, e)
		}
	}
	return ok
}

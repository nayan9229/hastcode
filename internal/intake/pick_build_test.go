package intake_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildtest"
	"hastcode.com/site/internal/intake"
)

// pickedTree picks one fixture Lead per kind into in's content tree and
// returns a Build-seam tree holding the four skeletons, the house byline
// they name and the Explainer the Brief links, with the skeletons' paths.
func pickedTree(t *testing.T) (buildtest.Tree, []string) {
	t.Helper()
	in := newPickIntake(t)
	linked := buildtest.ValidExplainer() // raft-leader-election
	withExplainer(t, in, linked.Vertical, linked.Slug)

	picks := []struct {
		title string
		opts  intake.PickOptions
	}{
		{regionOutage, intake.PickOptions{}},                     // case-study
		{"Why our B-tree pages are 8 KiB", intake.PickOptions{}}, // concept
		{raftLead, intake.PickOptions{Kind: "comparison"}},       // comparison
		{walLead, intake.PickOptions{Kind: "brief", Explainers: []string{linked.Slug}}},
	}
	tree := buildtest.NewTree().
		WithAuthor(buildtest.Author{Slug: "mira-vance", Kind: "organization", Name: "Mira Vance", Body: "HastCode's house byline.\n"}).
		WithExplainer(linked)
	var paths []string
	for _, p := range picks {
		s := readSkeleton(t, in, mustPick(t, in, p.title, p.opts))
		data, err := os.ReadFile(filepath.Join(in.ContentDir, filepath.FromSlash(s.Path)))
		if err != nil {
			t.Fatal(err)
		}
		tree.With(s.Path, data)
		paths = append(paths, s.Path)
	}
	return tree, paths
}

func TestPick_SkeletonOfEveryKindBuildsInPreviewWithOnlyTheTODOWarning(t *testing.T) {
	tree, skeletons := pickedTree(t)
	cfg := buildtest.Config()
	cfg.Mode = builder.Preview

	_, res := buildtest.Build(t, cfg, tree)

	if len(res.Errors) != 0 {
		t.Fatalf("a picked skeleton must build in preview, got errors:\n%s", res)
	}
	for _, file := range skeletons {
		buildtest.RequireWarningRule(t, res, file, "author-placeholder", "TODO(author)")
	}
	for _, w := range res.Warnings {
		if slices.Contains(skeletons, w.File) && w.Rule != "author-placeholder" {
			t.Errorf("skeleton warning other than the TODO(author) one: %s", w)
		}
	}
}

func TestPick_SkeletonsAreDraftsAbsentFromAProductionBuild(t *testing.T) {
	tree, skeletons := pickedTree(t)

	out, res := buildtest.Build(t, buildtest.Config(), tree)

	if len(res.Errors) != 0 {
		t.Fatalf("draft skeletons must not fail a production build:\n%s", res)
	}
	for _, file := range skeletons {
		slug := strings.TrimSuffix(filepath.Base(file), ".md")
		for _, p := range buildtest.OutputFiles(t, out) {
			if strings.Contains(p, slug) {
				t.Errorf("production output %s carries the draft %s", p, file)
			}
		}
	}
}

// A skeleton writes no prose: every body line is a line of its kind's
// template, except the Lead's line in Sources and a Brief's links to the
// Explainers it names.
func TestPick_SkeletonBodyIsTheTemplatesLinesAndTheLeadOnly(t *testing.T) {
	in := newPickIntake(t)
	withExplainer(t, in, "distributed-systems", "metastable-failures")
	withExplainer(t, in, "apis-networking", "http-connection-reuse")
	for _, tc := range []struct {
		title string
		opts  intake.PickOptions
		kind  string
		added []string // the lines a pick may add to the template's
	}{
		{regionOutage, intake.PickOptions{}, "case-study", []string{"- [" + regionOutage + "](https://engineering.example.com/posts/region-outage), Example Engineering"}},
		{"Why our B-tree pages are 8 KiB", intake.PickOptions{}, "concept", []string{"- [Why our B-tree pages are 8 KiB](https://db.example.org/blog/btree-pages), Example Database Blog"}},
		{raftLead, intake.PickOptions{Kind: "comparison"}, "comparison", []string{"- [" + raftLead + "](https://engineering.example.com/posts/raft-five-years), Example Engineering"}},
		{walLead, intake.PickOptions{Kind: "brief", Explainers: []string{"metastable-failures", "http-connection-reuse"}}, "brief", []string{
			"TODO(author): our own analysis, never the Source restated: the architectural cause, the lesson for a Reader's systems, and a link to the Explainer that teaches it, as [the idea](metastable-failures). Say \"we infer\" for anything the Source does not state.",
			"TODO(author): link [the idea](http-connection-reuse) where the analysis uses it.",
		}},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			tmpl, err := os.ReadFile(filepath.Join(in.Templates, tc.kind+"-template.md"))
			if err != nil {
				t.Fatal(err)
			}
			_, tmplBody, _ := strings.Cut(strings.TrimPrefix(string(tmpl), "---\n"), "\n---\n")
			allowed := append(strings.Split(tmplBody, "\n"), tc.added...)

			s := readSkeleton(t, in, mustPick(t, in, tc.title, tc.opts))

			for _, line := range strings.Split(s.Body, "\n") {
				if !slices.Contains(allowed, line) {
					t.Errorf("%s body line is not the template's: %q", s.Path, line)
				}
			}
		})
	}
}

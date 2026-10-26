package builder_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

// The Case Study, Concept, Comparison and Brief templates
// (docs/editorial/<kind>-template.md) are what `make lead` writes as a
// skeleton and what an Author copies to start by hand. These tests fill
// each one the way its writing guide says and build all four together, so
// a template can never drift from its kind's editorial gates.

// kindTemplate is one template, where its filled copy lives in the content
// tree, the values an Author would write for its frontmatter placeholders
// (by key) and for its prose tokens.
type kindTemplate struct {
	file        string
	contentPath string
	frontmatter map[string]string
	tokens      *strings.Replacer
	// links fills the Explainer slugs a Brief body links, as the Author
	// does with its `explainers` frontmatter, even when the prose is not.
	links *strings.Replacer
	rules []string // the kind's structural gates
}

var explainerTemplateRules = []string{"frontmatter", "tldr", "sources", "diagram", "trade-offs-table", "hero-missing"}

var kindTemplates = []kindTemplate{
	{
		file:        "docs/editorial/case-study-template.md",
		contentPath: "explainers/infrastructure-cloud/single-tenant-hosts.md",
		frontmatter: map[string]string{
			"title":       `"Single-tenant hosts: how Example Co gives one task its own VM"`,
			"description": `"How Example Co gives a few noisy workloads a dedicated instance under the same container API, and what that costs."`,
			"vertical":    `infrastructure-cloud`,
			"tags":        `["containers", "isolation"]`,
			"published":   `"2026-10-21T09:00:00Z"`,
			"reviewed":    `"2026-10-21T09:00:00Z"`,
		},
		tokens: strings.NewReplacer("{system}", "Single-tenant hosts", "{company}", "Example Co", "{title}", "Single-tenant hosts"),
		rules:  explainerTemplateRules,
	},
	{
		file:        "docs/editorial/concept-template.md",
		contentPath: "explainers/distributed-systems/retry-storms.md",
		frontmatter: map[string]string{
			"title":       `"Retry storms: why a brief blip can keep a service down"`,
			"description": `"A retry storm outlives its trigger because retries keep the service overloaded. How the loop works and what breaks it."`,
			"vertical":    `distributed-systems`,
			"tags":        `["retries", "overload"]`,
			"published":   `"2026-10-21T09:00:00Z"`,
			"reviewed":    `"2026-10-21T09:00:00Z"`,
		},
		tokens: strings.NewReplacer("{concept}", "a retry storm", "{title}", "Retry storms"),
		rules:  explainerTemplateRules,
	},
	{
		file:        "docs/editorial/comparison-template.md",
		contentPath: "explainers/apis-networking/http1-vs-http2.md",
		frontmatter: map[string]string{
			"title":       `"HTTP/1.1 vs HTTP/2: how each reuses a connection"`,
			"description": `"How HTTP/1.1 and HTTP/2 reuse connections, where each breaks, and which to use between services."`,
			"vertical":    `apis-networking`,
			"tags":        `["http", "http2"]`,
			"published":   `"2026-10-21T09:00:00Z"`,
			"reviewed":    `"2026-10-21T09:00:00Z"`,
		},
		tokens: strings.NewReplacer("{option A}", "HTTP/1.1", "{option B}", "HTTP/2", "{title}", "Connection reuse"),
		rules:  explainerTemplateRules,
	},
	{
		file:        "docs/editorial/brief-template.md",
		contentPath: "briefs/queue-fan-out.md",
		frontmatter: map[string]string{
			"title":       `"Example Co's queue backlog: a fan-out that flooded its own queue"`,
			"description": `"A fan-out of alert events saturated a shared queue. The fix was a separate lane with a fixed rate."`,
			"tags":        `["queues", "incidents"]`,
			"published":   `"2026-10-21T09:00:00Z"`,
			"source_url":  `https://status.example.com/incidents/1`,
			"explainers":  `[retry-storms]`,
		},
		tokens: strings.NewReplacer("{company}", "Example Co", "{event}", "queue backlog"),
		links:  strings.NewReplacer("{explainer}", "retry-storms"),
		rules:  []string{"frontmatter", "brief-length", "brief-sections", "brief-source", "brief-explainer-link", "broken-link", "hero-missing"},
	},
}

// fill reads the template and fills its frontmatter as `make lead` and the
// Author would before the first build: every placeholder value replaced,
// and the draft flag dropped as it is at Publish. When fillBody is set,
// every TODO(author) hint in the body becomes the prose it describes (the
// marker dropped, the {token}s named).
func (k kindTemplate) fill(t *testing.T, fillBody bool) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(buildtest.RepoRoot(), filepath.FromSlash(k.file)))
	if err != nil {
		t.Fatal(err)
	}
	rest, ok := strings.CutPrefix(string(data), "---\n")
	fm, body, ok2 := strings.Cut(rest, "\n---\n")
	if !ok || !ok2 {
		t.Fatalf("%s has no frontmatter block", k.file)
	}

	var lines []string
	filled := map[string]bool{}
	for _, line := range strings.Split(fm, "\n") {
		key, _, isKey := strings.Cut(line, ":")
		comment := strings.HasPrefix(line, "#")
		if isKey && !comment && key == "verified" {
			t.Errorf("%s sets verified: only the Author does", k.file)
		}
		switch {
		case line == "draft: true":
			continue
		case strings.Contains(line, "TODO(author)") && !comment:
			value, known := k.frontmatter[key]
			if !known {
				t.Fatalf("%s frontmatter key %q has a placeholder this test does not fill", k.file, key)
			}
			line = key + ": " + value
			filled[key] = true
		}
		lines = append(lines, line)
	}
	for key := range k.frontmatter {
		if !filled[key] {
			t.Errorf("%s frontmatter has no %s placeholder", k.file, key)
		}
	}
	for _, want := range []string{"\nauthor: mira-vance\n", "\ndraft: true\n"} {
		if !strings.Contains("\n"+fm+"\n", want) {
			t.Errorf("%s frontmatter lacks %q", k.file, strings.TrimSpace(want))
		}
	}
	if !regexp.MustCompile(`(?m)^# hero: `).MatchString(fm) || !regexp.MustCompile(`(?m)^# hero_alt: `).MatchString(fm) {
		t.Errorf("%s does not offer the optional hero and hero_alt as comments", k.file)
	}

	if k.links != nil {
		body = k.links.Replace(body)
	}
	if fillBody {
		body = k.tokens.Replace(strings.ReplaceAll(body, "TODO(author): ", ""))
		if strings.Contains(body, "TODO(author)") || strings.Contains(body, "{") {
			t.Fatalf("%s body has a placeholder not of the form \"TODO(author): hint\":\n%s", k.file, body)
		}
	}
	return []byte("---\n" + strings.Join(lines, "\n") + "\n---\n" + body)
}

// kindTemplateTree holds every filled template (the Brief links the
// Concept) and the house byline they name.
func kindTemplateTree(t *testing.T, fillBody bool) buildtest.Tree {
	tree := buildtest.NewTree().
		WithAuthor(buildtest.Author{Slug: "mira-vance", Kind: "organization", Name: "Mira Vance", Body: "HastCode's house byline.\n"})
	for _, k := range kindTemplates {
		tree.With(k.contentPath, k.fill(t, fillBody))
	}
	return tree
}

func TestKindTemplates_FilledWithProseBuildCleanUnderEveryGate(t *testing.T) {
	_, res := buildtest.Build(t, buildtest.Config(), kindTemplateTree(t, true))

	for _, k := range kindTemplates {
		for _, rule := range append(k.rules, "author-placeholder", "section-length") {
			buildtest.RequireNoProblemRule(t, res, k.contentPath, rule)
		}
	}
	if len(res.Errors) != 0 {
		t.Fatalf("the filled templates must build clean:\n%s", res)
	}
}

func TestKindTemplates_UnfilledHintsFailAProductionBuild(t *testing.T) {
	_, res := buildtest.Build(t, buildtest.Config(), kindTemplateTree(t, false))

	for _, k := range kindTemplates {
		buildtest.RequireErrorRule(t, res, k.contentPath, "author-placeholder", "TODO(author)")
		for _, rule := range k.rules {
			buildtest.RequireNoProblemRule(t, res, k.contentPath, rule)
		}
	}
}

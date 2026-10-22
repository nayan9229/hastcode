package builder_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

// The Repo Profile template (docs/editorial/repo-profile-template.md) is
// what an Author copies to start a Profile. These tests copy it into a
// fixture tree the way the writing guide says, fill its placeholders, and
// build it, so the template can never drift from the Repo Profile rules.

const templateFile = "docs/editorial/repo-profile-template.md"

// templateFrontmatter fills the template's frontmatter placeholders, by
// key, with the values an Author would write for apache/iggy.
var templateFrontmatter = map[string]string{
	"title":       `"Apache Iggy: a persistent message stream in Rust"`,
	"description": `"Iggy keeps Kafka's partitioned log and drops the JVM: one Rust binary, three transports, no clustering yet."`,
	"vertical":    `distributed-systems`,
	"repo":        `apache/iggy`,
	"who_for":     `"Teams that need Kafka-style streams without running a JVM cluster."`,
	"tags":        `["streaming", "rust"]`,
	"published":   `"2026-10-21T09:00:00Z"`,
	"reviewed":    `"2026-10-21T09:00:00Z"`,
}

// templatePlaceholders fills the tokens the template's prose hints use.
var templatePlaceholders = strings.NewReplacer(
	"https://github.com/{repo}", "https://github.com/apache/iggy",
	"{repo}", "Iggy",
	"{feature}", "QUIC",
	"{alternative}", "Kafka",
)

// copyTemplate reads the template and fills its frontmatter as an Author
// would before the first build: every placeholder value replaced, and the
// draft flag dropped as it is at Publish. When fillBody is set, every
// TODO(author) hint in the body is turned into the prose it describes
// (the marker dropped, the {repo}-style tokens named).
func copyTemplate(t *testing.T, fillBody bool) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(buildtest.RepoRoot(), filepath.FromSlash(templateFile)))
	if err != nil {
		t.Fatal(err)
	}
	src := string(data)
	rest, ok := strings.CutPrefix(src, "---\n")
	fm, body, ok2 := strings.Cut(rest, "\n---\n")
	if !ok || !ok2 {
		t.Fatalf("%s has no frontmatter block", templateFile)
	}

	var lines []string
	filled := map[string]bool{}
	keys := map[string]bool{}
	for _, line := range strings.Split(fm, "\n") {
		key, _, isKey := strings.Cut(line, ":")
		if isKey && !strings.HasPrefix(line, "#") {
			keys[key] = true
		}
		switch {
		case line == "draft: true":
			continue
		case strings.Contains(line, "TODO(author)") && !strings.HasPrefix(line, "#"):
			value, known := templateFrontmatter[key]
			if !known {
				t.Fatalf("template frontmatter key %q has a placeholder this test does not fill", key)
			}
			line = key + ": " + value
			filled[key] = true
		}
		lines = append(lines, line)
	}
	for key := range templateFrontmatter {
		if !filled[key] {
			t.Errorf("template frontmatter has no %s placeholder", key)
		}
	}
	for _, want := range []string{"kind", "repo", "who_for", "author", "draft"} {
		if !keys[want] {
			t.Errorf("template frontmatter lacks %s", want)
		}
	}
	if !strings.Contains(fm, "\nkind: repo-profile\n") {
		t.Errorf("template frontmatter does not set kind: repo-profile")
	}
	if !regexp.MustCompile(`(?m)^# hero: `).MatchString(fm) || !regexp.MustCompile(`(?m)^# hero_alt: `).MatchString(fm) {
		t.Errorf("template frontmatter does not offer the optional hero and hero_alt")
	}

	if fillBody {
		body = templatePlaceholders.Replace(strings.ReplaceAll(body, "TODO(author): ", ""))
		if strings.Contains(body, "TODO(author)") || strings.Contains(body, "{") {
			t.Fatalf("template body has a placeholder not of the form \"TODO(author): hint\":\n%s", body)
		}
	}
	return []byte("---\n" + strings.Join(lines, "\n") + "\n---\n" + body)
}

// templateTree holds the filled template as the apache/iggy Profile with
// its Fact Pack beside it, and the house byline it names.
func templateTree(md []byte) buildtest.Tree {
	fp := buildtest.ValidFactPack()
	return buildtest.NewTree().
		WithAuthor(buildtest.Author{Slug: "mira-vance", Kind: "organization", Name: "Mira Vance", Body: "HastCode's house byline.\n"}).
		With(repoProfileFile, md).
		With(strings.TrimSuffix(repoProfileFile, ".md")+".json", fp.JSON())
}

func TestRepoProfileTemplate_FilledWithProseBuildsCleanUnderEveryRepoProfileRule(t *testing.T) {
	out, res := buildtest.Build(t, buildtest.Config(), templateTree(copyTemplate(t, true)))

	rules := append([]string{"frontmatter", "tldr", "sources", "diagram", "fact-pack", "fact-pack-repo", "fact-sections", "hero-missing", "author-placeholder"}, repoProfileRules...)
	for _, rule := range rules {
		buildtest.RequireNoProblemRule(t, res, repoProfileFile, rule)
	}
	if len(res.Errors) != 0 {
		t.Fatalf("the filled template must build clean:\n%s", res)
	}

	doc := buildtest.Page(t, out, repoProfilePath)
	var questions []string
	faq := buildtest.JSONLDOfType(t, doc, "FAQPage")
	entities, _ := faq["mainEntity"].([]any)
	for _, e := range entities {
		q, _ := e.(map[string]any)
		name, _ := q["name"].(string)
		questions = append(questions, name)
	}
	if len(questions) != 5 || questions[0] != "Is Iggy production ready?" {
		t.Errorf("FAQPage questions = %q, want the template's five, starting with the production-ready search", questions)
	}
	if got := doc.Find("article blockquote").Text(); !strings.Contains(got, "Drafted with AI assistance") {
		t.Errorf("the page does not carry the AI-assistance note; blockquote = %q", got)
	}
	for _, sel := range []string{"section.signals", "section.changelog"} {
		if doc.Find(sel).Length() != 1 {
			t.Errorf("the page has no %s rendered from the Fact Pack", sel)
		}
	}
}

func TestRepoProfileTemplate_UnfilledHintsFailAProductionBuild(t *testing.T) {
	_, res := buildtest.Build(t, buildtest.Config(), templateTree(copyTemplate(t, false)))

	buildtest.RequireErrorRule(t, res, repoProfileFile, "author-placeholder", "TODO(author)")
	for _, rule := range repoProfileRules {
		buildtest.RequireNoProblemRule(t, res, repoProfileFile, rule)
	}
}

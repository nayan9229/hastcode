package builder_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

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

	rules := append([]string{"frontmatter", "tldr", "sources", "diagram", "fact-pack", "fact-pack-repo", "fact-sections", "hero-missing", "author-placeholder", "section-length"}, repoProfileRules...)
	rules = append(rules, proseRules...)
	for _, rule := range rules {
		buildtest.RequireNoProblemRule(t, res, repoProfileFile, rule)
	}
	if len(res.Errors) != 0 {
		t.Fatalf("the filled template must build clean:\n%s", res)
	}

	doc := buildtest.Page(t, out, repoProfilePath)
	var questions []string
	doc.Find("article .e-content h3").Each(func(_ int, h *goquery.Selection) {
		questions = append(questions, strings.TrimSpace(h.Text()))
	})
	if len(questions) != 5 || questions[0] != "Is Iggy production ready?" {
		t.Errorf("FAQ questions = %q, want the template's five, starting with the production-ready search", questions)
	}
	if got := doc.Find("article").Text(); strings.Contains(got, "AI assistance") {
		t.Errorf("the page carries AI-assistance wording; no AI note is published (ADR-0013):\n%s", got)
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

// The quickstart-ran line (docs/editorial/repo-profile-writing-guide.md,
// "In 60 seconds") is the Author's own report of running the quickstart.
// It is optional and Author-only, so the template offers it in its header
// comment, beside the other optional keys, never as a body line a Profile
// would have to fill or delete.

// quickstartRanComment is the template's header comment that offers the
// line: the run of frontmatter comment lines that mentions it, joined.
func quickstartRanComment(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(buildtest.RepoRoot(), filepath.FromSlash(templateFile)))
	if err != nil {
		t.Fatal(err)
	}
	rest, _ := strings.CutPrefix(string(data), "---\n")
	fm, _, _ := strings.Cut(rest, "\n---\n")
	var block []string
	for _, line := range strings.Split(fm+"\n", "\n") {
		comment, isComment := strings.CutPrefix(line, "#")
		switch {
		case isComment:
			block = append(block, strings.TrimSpace(comment))
		case strings.Contains(strings.Join(block, " "), "We ran the quickstart"):
			return strings.Join(block, " ")
		default:
			block = nil
		}
	}
	t.Fatalf("%s's header comment does not offer the optional \"We ran the quickstart\" line", templateFile)
	return ""
}

func TestRepoProfileTemplate_OffersTheQuickstartRanLineAsOptionalAndAuthorOnly(t *testing.T) {
	comment := quickstartRanComment(t)
	for _, want := range []string{
		"Optional",
		`after the quickstart block in "In 60 seconds"`,
		"when, on what, and whether it worked as written",
		"Never change the commands",
		"Only the Author writes it",
		"an agent never does",
	} {
		if !strings.Contains(comment, want) {
			t.Errorf("the quickstart-ran comment does not say %q:\n%s", want, comment)
		}
	}
}

// ranLine is the line as an Author writes it after a real run.
const ranLine = "We ran the quickstart on 2026-11-17 on macOS 26 (Apple M3), at Iggy 0.5.0: it worked as written."

// withRanLine puts ranLine after the quickstart block, where the guide
// says it goes: the first fence closed under "## In 60 seconds".
func withRanLine(t *testing.T, md []byte) []byte {
	t.Helper()
	head, rest, ok := strings.Cut(string(md), "## In 60 seconds\n")
	if !ok {
		t.Fatal("the template has no In 60 seconds section")
	}
	open := strings.Index(rest, "```")
	if open < 0 {
		t.Fatal("the In 60 seconds section has no fenced quickstart block")
	}
	closing := strings.Index(rest[open+3:], "\n```\n")
	if closing < 0 {
		t.Fatal("the quickstart block is never closed")
	}
	end := open + 3 + closing + len("\n```\n")
	return []byte(head + "## In 60 seconds\n" + rest[:end] + "\n" + ranLine + "\n" + rest[end:])
}

func TestRepoProfileTemplate_BuildsCleanWithAndWithoutTheQuickstartRanLine(t *testing.T) {
	for _, tc := range []struct {
		name    string
		md      []byte
		withRan bool
	}{
		{"skipped", copyTemplate(t, true), false},
		{"written by the Author", withRanLine(t, copyTemplate(t, true)), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, res := buildtest.Build(t, buildtest.Config(), templateTree(tc.md))
			for _, rule := range append([]string{"author-placeholder", "section-length"}, repoProfileRules...) {
				buildtest.RequireNoProblemRule(t, res, repoProfileFile, rule)
			}
			if len(res.Errors) != 0 {
				t.Fatalf("the Profile must build clean:\n%s", res)
			}

			page := buildtest.Page(t, out, repoProfilePath)
			article, _ := page.Find("article").Html()
			twin := readStatic(t, out, strings.TrimSuffix(repoProfilePath, "/")+".md")
			for where, text := range map[string]string{"page": article, "Markdown twin": twin} {
				if strings.Contains(text, "omitted") {
					t.Errorf("the %s carries an omitted-HTML marker", where)
				}
				if got := strings.Contains(text, "We ran the quickstart"); got != tc.withRan {
					t.Errorf("the %s shows the quickstart-ran line = %v, want %v", where, got, tc.withRan)
				}
			}
			if tc.withRan {
				after := page.Find("article pre").First().NextAllFiltered("p").First().Text()
				if after != ranLine {
					t.Errorf("the paragraph after the quickstart block = %q, want the line", after)
				}
			}
		})
	}
}

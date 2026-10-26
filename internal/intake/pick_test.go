package intake_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"hastcode.com/site/internal/intake"
)

// The skeleton headings of each kind, as its writing guide lists them under
// "The skeleton" (docs/editorial/<kind>-writing-guide.md).
var (
	caseStudyHeadings = []string{
		"TL;DR",
		"TODO(author): {system} at a glance, named for what the diagram shows",
		"TODO(author): the problem {company} had to solve",
		"TODO(author): the options {company} weighed",
		"TODO(author): the first part of the diagram and what it does",
		"TODO(author): the next part of the diagram and what it does",
		"TODO(author): what it costs",
		"TODO(author): the evidence it works, and what the Source leaves out",
		"TODO(author): what it makes possible",
		"Sources",
	}
	conceptHeadings = []string{
		"TL;DR",
		"TODO(author): what {concept} is",
		"TODO(author): how {concept} works, named for what the diagram shows",
		"TODO(author): where {concept} shows up in real systems",
		"TODO(author): why systems end up with {concept}",
		"TODO(author): the first practice that applies or defends against {concept}",
		"TODO(author): what {concept} is not",
		"TODO(author): what understanding {concept} makes possible",
		"Sources",
	}
	comparisonHeadings = []string{
		"TL;DR",
		"TODO(author): the options side by side, named for what the diagram shows",
		"TODO(author): {option A}: how it works",
		"TODO(author): {option B}: how it works",
		"TODO(author): the first dimension, answered for every option",
		"Trade-offs",
		"TODO(author): when each one wins",
		"Sources",
	}
	briefHeadings = []string{"What happened", "Why it matters"}
)

// The fixture Leads (testdata/feeds), by title.
const (
	regionOutage = "How we lost a region for 40 minutes" // Example Engineering, case-study, distributed-systems, 2026-10-24 14:30 UTC
	walLead      = "Inside the write-ahead log"          // Example Database Blog, concept, databases-storage, 2026-10-25 08:00 UTC
	raftLead     = "Raft in production: five years on"   // Example Engineering, 2026-10-19 09:00 UTC
)

// newPickIntake is newIntake over the example Watchlist with an empty temp
// content tree, the real kind templates, and one run done so the fixture
// Leads are open in the queue.
func newPickIntake(t *testing.T) *intake.Intake {
	t.Helper()
	in := newIntake(t, exampleWatchlist(t, newFeedServer(t)))
	in.ContentDir = filepath.Join(t.TempDir(), "content")
	if err := os.Mkdir(in.ContentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	in.Templates = filepath.Join("..", "..", "docs", "editorial")
	mustRun(t, in)
	return in
}

// leadID is the queue ID of the open Lead titled title.
func leadID(t *testing.T, in *intake.Intake, title string) string {
	t.Helper()
	for _, l := range readQueue(t, in) {
		if l.Title == title {
			return l.ID
		}
	}
	t.Fatalf("no open Lead %q in the queue", title)
	return ""
}

func mustPick(t *testing.T, in *intake.Intake, title string, opts intake.PickOptions) string {
	t.Helper()
	file, err := in.Pick(context.Background(), leadID(t, in, title), opts)
	if err != nil {
		t.Fatalf("Pick %q: %v", title, err)
	}
	return file
}

// skeleton is a picked skeleton as the Author opens it.
type skeleton struct {
	Path        string         // relative to the content tree, slash-separated
	Frontmatter map[string]any // the YAML frontmatter
	Keys        []string       // its keys in order
	Body        string
	Headings    []string // every "## " heading, in order, outside code fences
}

func readSkeleton(t *testing.T, in *intake.Intake, file string) skeleton {
	t.Helper()
	rel, err := filepath.Rel(in.ContentDir, file)
	if err != nil || strings.HasPrefix(rel, "..") {
		t.Fatalf("skeleton %s is not under the content tree %s", file, in.ContentDir)
	}
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	rest, ok := strings.CutPrefix(string(data), "---\n")
	fm, body, ok2 := strings.Cut(rest, "\n---\n")
	if !ok || !ok2 {
		t.Fatalf("%s has no frontmatter block:\n%s", rel, data)
	}
	s := skeleton{Path: filepath.ToSlash(rel), Body: body}
	if err := yaml.Unmarshal([]byte(fm), &s.Frontmatter); err != nil {
		t.Fatalf("%s frontmatter: %v", rel, err)
	}
	var node yaml.Node
	if err := yaml.Unmarshal([]byte(fm), &node); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(node.Content[0].Content); i += 2 {
		s.Keys = append(s.Keys, node.Content[0].Content[i].Value)
	}
	fence := false
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "```") {
			fence = !fence
		}
		if h, ok := strings.CutPrefix(line, "## "); ok && !fence {
			s.Headings = append(s.Headings, h)
		}
	}
	return s
}

// section is the body text under the heading h, up to the next heading.
func (s skeleton) section(t *testing.T, h string) string {
	t.Helper()
	_, after, ok := strings.Cut(s.Body, "\n## "+h+"\n")
	if !ok {
		t.Fatalf("%s has no ## %s", s.Path, h)
	}
	text, _, _ := strings.Cut(after, "\n## ")
	return strings.TrimSpace(text)
}

func TestPick_CaseStudySkeletonHasFrontmatterHeadingsAndTheLeadInSources(t *testing.T) {
	in := newPickIntake(t)

	s := readSkeleton(t, in, mustPick(t, in, regionOutage, intake.PickOptions{}))

	if want := "explainers/distributed-systems/how-we-lost-a-region-for-40-minutes.md"; s.Path != want {
		t.Errorf("skeleton at %s, want %s", s.Path, want)
	}
	for key, want := range map[string]any{
		"title":     regionOutage,
		"kind":      "case-study",
		"vertical":  "distributed-systems",
		"author":    "mira-vance",
		"draft":     true,
		"published": "2026-10-26T09:00:00Z", // provisional: the pick time
		"reviewed":  "2026-10-26T09:00:00Z",
	} {
		if got := s.Frontmatter[key]; got != want {
			t.Errorf("frontmatter %s = %#v, want %#v", key, got, want)
		}
	}
	for _, key := range []string{"verified", "hero", "hero_alt"} {
		if _, ok := s.Frontmatter[key]; ok {
			t.Errorf("frontmatter sets %s; a skeleton never does", key)
		}
	}
	if !slices.Equal(s.Headings, caseStudyHeadings) {
		t.Errorf("headings:\n%q\nwant the case-study template's:\n%q", s.Headings, caseStudyHeadings)
	}
	if got, want := s.section(t, "Sources"), "- ["+regionOutage+"](https://engineering.example.com/posts/region-outage), Example Engineering"; got != want {
		t.Errorf("## Sources = %q, want the Lead's link %q", got, want)
	}
}

func TestPick_ConceptSkeletonFollowsTheConceptTemplate(t *testing.T) {
	in := newPickIntake(t)

	s := readSkeleton(t, in, mustPick(t, in, walLead, intake.PickOptions{}))

	if want := "explainers/databases-storage/inside-the-write-ahead-log.md"; s.Path != want {
		t.Errorf("skeleton at %s, want %s", s.Path, want)
	}
	if s.Frontmatter["kind"] != "concept" || s.Frontmatter["vertical"] != "databases-storage" {
		t.Errorf("kind, vertical = %v, %v; want the Source's concept, databases-storage", s.Frontmatter["kind"], s.Frontmatter["vertical"])
	}
	if !slices.Equal(s.Headings, conceptHeadings) {
		t.Errorf("headings:\n%q\nwant the concept template's:\n%q", s.Headings, conceptHeadings)
	}
	if got, want := s.section(t, "Sources"), "- ["+walLead+"](https://db.example.org/blog/wal), Example Database Blog"; got != want {
		t.Errorf("## Sources = %q, want %q", got, want)
	}
}

func TestPick_KindAndVerticalOverridesReplaceTheSourcesSuggestion(t *testing.T) {
	in := newPickIntake(t)

	s := readSkeleton(t, in, mustPick(t, in, regionOutage, intake.PickOptions{Kind: "comparison", Vertical: "databases-storage"}))

	if want := "explainers/databases-storage/how-we-lost-a-region-for-40-minutes.md"; s.Path != want {
		t.Errorf("skeleton at %s, want %s", s.Path, want)
	}
	if s.Frontmatter["kind"] != "comparison" || s.Frontmatter["vertical"] != "databases-storage" {
		t.Errorf("kind, vertical = %v, %v; want the overrides comparison, databases-storage", s.Frontmatter["kind"], s.Frontmatter["vertical"])
	}
	if !slices.Equal(s.Headings, comparisonHeadings) {
		t.Errorf("headings:\n%q\nwant the comparison template's:\n%q", s.Headings, comparisonHeadings)
	}
}

// withExplainer puts an Explainer file at explainers/<vertical>/<slug>.md
// in in's content tree, for a Brief to link.
func withExplainer(t *testing.T, in *intake.Intake, vertical, slug string) {
	t.Helper()
	file := filepath.Join(in.ContentDir, "explainers", vertical, slug+".md")
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte("---\ntitle: \"An Explainer\"\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestPick_BriefSkeletonLinksTheLeadAndTheNamedExplainers(t *testing.T) {
	in := newPickIntake(t)
	withExplainer(t, in, "distributed-systems", "metastable-failures")
	withExplainer(t, in, "apis-networking", "http-connection-reuse")

	s := readSkeleton(t, in, mustPick(t, in, regionOutage, intake.PickOptions{
		Kind:       "brief",
		Explainers: []string{"metastable-failures", "http-connection-reuse"},
	}))

	if want := "briefs/how-we-lost-a-region-for-40-minutes.md"; s.Path != want {
		t.Errorf("skeleton at %s, want %s", s.Path, want)
	}
	wantKeys := []string{"title", "description", "tags", "published", "author", "source_url", "explainers", "draft"}
	if !slices.Equal(s.Keys, wantKeys) {
		t.Errorf("frontmatter keys %q, want the Brief template's %q (no kind, vertical, reviewed, verified or hero)", s.Keys, wantKeys)
	}
	for key, want := range map[string]any{
		"title":      regionOutage,
		"published":  "2026-10-26T09:00:00Z",
		"author":     "mira-vance",
		"source_url": "https://engineering.example.com/posts/region-outage",
		"draft":      true,
	} {
		if got := s.Frontmatter[key]; got != want {
			t.Errorf("frontmatter %s = %#v, want %#v", key, got, want)
		}
	}
	if got := s.Frontmatter["explainers"]; !slices.Equal(toStrings(got), []string{"metastable-failures", "http-connection-reuse"}) {
		t.Errorf("explainers = %#v, want exactly the named slugs", got)
	}
	if !slices.Equal(s.Headings, briefHeadings) {
		t.Errorf("headings %q, want %q", s.Headings, briefHeadings)
	}
	why := s.section(t, "Why it matters")
	for _, link := range []string{"](metastable-failures)", "](http-connection-reuse)"} {
		if !strings.Contains(why, link) {
			t.Errorf("## Why it matters does not link %s:\n%s", link, why)
		}
	}
}

func TestPick_MarksTheLeadPickedAndLaterQueuesLeaveItOut(t *testing.T) {
	in := newPickIntake(t)
	id := leadID(t, in, walLead)

	if _, err := in.Pick(context.Background(), id, intake.PickOptions{}); err != nil {
		t.Fatal(err)
	}

	if got := readState(t, in)[id].Decision; got != "picked" {
		t.Errorf("Lead decision %q after a pick, want picked", got)
	}
	mustRun(t, in)
	if got := titles(readQueue(t, in)); slices.Contains(got, walLead) {
		t.Errorf("the next run's queue %q still lists the picked Lead", got)
	}
	if _, err := in.Pick(context.Background(), id, intake.PickOptions{}); err == nil || !strings.Contains(err.Error(), "already picked") {
		t.Errorf("picking it again: err = %v, want a refusal naming it already picked", err)
	}
}

// requireRefused checks that picking the Lead titled title with opts fails
// with an error containing want, writes nothing under the content tree and
// leaves the Lead open.
func requireRefused(t *testing.T, in *intake.Intake, title string, opts intake.PickOptions, want string) {
	t.Helper()
	id := leadID(t, in, title)
	before := filesUnder(t, in.ContentDir)

	_, err := in.Pick(context.Background(), id, opts)

	if err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("Pick %q %+v: err = %v, want a refusal containing %q", title, opts, err, want)
	}
	if after := filesUnder(t, in.ContentDir); !slices.Equal(after, before) {
		t.Errorf("a refused pick changed the content tree: %q, was %q", after, before)
	}
	if got := readState(t, in)[id].Decision; got != "open" {
		t.Errorf("a refused pick left the Lead %s, want open", got)
	}
}

func TestPick_RefusesASlugAPieceAlreadyUses(t *testing.T) {
	in := newPickIntake(t)
	// Slugs are site-wide: the same slug under another Vertical counts.
	withExplainer(t, in, "languages-runtimes", "inside-the-write-ahead-log")

	requireRefused(t, in, walLead, intake.PickOptions{}, `slug "inside-the-write-ahead-log" is already used by explainers/languages-runtimes/inside-the-write-ahead-log.md`)
}

func TestPick_RefusesASlugABriefAlreadyUses(t *testing.T) {
	in := newPickIntake(t)
	file := filepath.Join(in.ContentDir, "briefs", "how-we-lost-a-region-for-40-minutes.md")
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte("the Author's Brief\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	requireRefused(t, in, regionOutage, intake.PickOptions{}, "already used by briefs/how-we-lost-a-region-for-40-minutes.md")
	if data, _ := os.ReadFile(file); string(data) != "the Author's Brief\n" {
		t.Errorf("the existing Brief was overwritten: %q", data)
	}
}

func TestPick_RefusesABriefFromALeadFourDaysOld(t *testing.T) {
	in := newPickIntake(t)
	withExplainer(t, in, "distributed-systems", "metastable-failures")
	// The Lead was published 2026-10-24 14:30 UTC.
	at(in, time.Date(2026, 10, 28, 14, 30, 0, 0, time.UTC))

	requireRefused(t, in, regionOutage, intake.PickOptions{Kind: "brief", Explainers: []string{"metastable-failures"}}, "a Brief needs a Lead at most 3 days")
}

func TestPick_TakesABriefFromALeadExactlyThreeDaysOld(t *testing.T) {
	in := newPickIntake(t)
	withExplainer(t, in, "distributed-systems", "metastable-failures")
	at(in, time.Date(2026, 10, 27, 14, 30, 0, 0, time.UTC))

	s := readSkeleton(t, in, mustPick(t, in, regionOutage, intake.PickOptions{Kind: "brief", Explainers: []string{"metastable-failures"}}))

	if s.Frontmatter["published"] != "2026-10-27T14:30:00Z" {
		t.Errorf("published = %v, want the pick time 2026-10-27T14:30:00Z", s.Frontmatter["published"])
	}
}

func TestPick_RefusesABriefThatNamesNoExplainer(t *testing.T) {
	in := newPickIntake(t)

	requireRefused(t, in, regionOutage, intake.PickOptions{Kind: "brief"}, "EXPLAINER=")
}

func TestPick_RefusesABriefNamingAnExplainerTheContentTreeLacks(t *testing.T) {
	in := newPickIntake(t)
	withExplainer(t, in, "distributed-systems", "metastable-failures")

	requireRefused(t, in, regionOutage, intake.PickOptions{Kind: "brief", Explainers: []string{"metastable-failures", "retry-storms"}}, `no Explainer "retry-storms"`)
}

func TestPick_RefusesAnExplainerSlugThatIsNotExact(t *testing.T) {
	in := newPickIntake(t)
	withExplainer(t, in, "distributed-systems", "metastable-failures")

	for _, slug := range []string{"meta*", "*", "metastable-failure?", "../distributed-systems/metastable-failures"} {
		requireRefused(t, in, regionOutage, intake.PickOptions{Kind: "brief", Explainers: []string{slug}}, fmt.Sprintf("EXPLAINER=%q is not a slug", slug))
	}
}

func TestPick_RefusesAnExplainerSlugThatOnlyPrefixesOne(t *testing.T) {
	in := newPickIntake(t)
	withExplainer(t, in, "distributed-systems", "metastable-failures")

	requireRefused(t, in, regionOutage, intake.PickOptions{Kind: "brief", Explainers: []string{"metastable"}}, `no Explainer "metastable"`)
}

func TestPick_RefusesExplainersForAnExplainerKind(t *testing.T) {
	in := newPickIntake(t)
	withExplainer(t, in, "distributed-systems", "metastable-failures")

	requireRefused(t, in, regionOutage, intake.PickOptions{Explainers: []string{"metastable-failures"}}, "only a Brief")
}

func TestPick_RefusesAnUnknownKindOrVertical(t *testing.T) {
	in := newPickIntake(t)

	requireRefused(t, in, regionOutage, intake.PickOptions{Kind: "repo-profile"}, `unknown kind "repo-profile"`)
	requireRefused(t, in, regionOutage, intake.PickOptions{Vertical: "frontend"}, `unknown Vertical "frontend"`)
}

func TestPick_RefusesAnUnknownLead(t *testing.T) {
	in := newPickIntake(t)

	if _, err := in.Pick(context.Background(), "zzzzzzzz", intake.PickOptions{}); err == nil || !strings.Contains(err.Error(), `no Lead "zzzzzzzz"`) {
		t.Errorf("err = %v, want no Lead", err)
	}
}

// toStrings is a parsed YAML sequence of strings (v as yaml.v3 decodes it
// into any) as a []string: "" for a non-string item, nil for a non-sequence.
func toStrings(v any) []string {
	items, _ := v.([]any)
	var out []string
	for _, item := range items {
		s, _ := item.(string)
		out = append(out, s)
	}
	return out
}

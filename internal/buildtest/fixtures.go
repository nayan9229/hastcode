// Package buildtest is the shared test support for the Build seam.
//
// Tests build a fixture content tree (Tree), run the builder into a temp
// output directory (Build / MustBuild), and assert on the output by parsing
// it (Page, JSONLD, VercelConfig). Tests never call internal packages; they
// only go through builder.Build and these helpers.
package buildtest

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"hastcode.com/site/builder"
)

// Tree is a fixture content tree: content-relative path -> file bytes.
type Tree map[string][]byte

// NewTree returns a tree holding the default Author only. Add pieces with
// the With* methods.
func NewTree() Tree {
	t := Tree{}
	t.WithAuthor(DefaultAuthor())
	return t
}

// ValidTree returns a tree with the default Author and one valid Explainer
// (ValidExplainer()).
func ValidTree() Tree {
	return NewTree().WithExplainer(ValidExplainer())
}

// With sets an arbitrary file in the tree.
func (t Tree) With(path string, data []byte) Tree {
	t[path] = data
	return t
}

// WithExplainer adds an Explainer at explainers/<vertical>/<slug>.md and,
// when e.Hero is set and HeroImage is non-nil, its hero raster beside it.
func (t Tree) WithExplainer(e Explainer) Tree {
	dir := filepath.Join("explainers", e.Vertical)
	t[filepath.Join(dir, e.Slug+".md")] = []byte(e.Markdown())
	if e.Hero != "" && e.HeroImage != nil {
		t[filepath.Join(dir, e.Hero)] = e.HeroImage
	}
	return t
}

// WithAuthor adds an Author at authors/<slug>.md.
func (t Tree) WithAuthor(a Author) Tree {
	t[filepath.Join("authors", a.Slug+".md")] = []byte(a.Markdown())
	return t
}

// Write materialises the tree in a fresh temp dir and returns its path.
func (t Tree) Write(tb testing.TB) string {
	tb.Helper()
	dir := tb.TempDir()
	for rel, data := range t {
		p := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			tb.Fatal(err)
		}
		if err := os.WriteFile(p, data, 0o644); err != nil {
			tb.Fatal(err)
		}
	}
	return dir
}

// Explainer is a fixture Explainer. Zero-value fields are omitted from the
// frontmatter, so tests can blank a field to make it invalid.
type Explainer struct {
	Vertical    string
	Slug        string
	Title       string
	Description string
	Kind        string
	Tags        []string
	Published   string
	Reviewed    string
	Author      string
	Hero        string
	HeroAlt     string
	HeroImage   []byte
	Related     []string
	Draft       bool
	Body        string
}

// ValidExplainer returns an Explainer that satisfies the whole editorial
// skeleton in the spec: intro, TL;DR (3 bullets), a Mermaid diagram,
// declarative sections, a Sources block, and a 1600x900 hero.
func ValidExplainer() Explainer {
	return Explainer{
		Vertical:    "distributed-systems",
		Slug:        "raft-leader-election",
		Title:       "How Raft elects a leader",
		Description: "Raft keeps a replicated log consistent by electing one leader per term with randomised timeouts.",
		Kind:        "concept",
		Tags:        []string{"consensus", "raft"},
		Published:   "2026-10-01T09:00:00+05:30",
		Reviewed:    "2026-10-15T18:30:00Z",
		Author:      DefaultAuthor().Slug,
		Hero:        "raft-leader-election.png",
		HeroAlt:     "Five Raft servers voting for a leader",
		HeroImage:   HeroPNG(1600, 900),
		Body:        ValidExplainerBody,
	}
}

// ValidExplainerBody is the Markdown body of ValidExplainer.
const ValidExplainerBody = `Raft splits consensus into leader election, log replication and safety.

## TL;DR

- Each term has at most one leader.
- Followers become candidates when their election timeout fires.
- A candidate needs votes from a majority of the cluster.

## Elections use randomised timeouts

Every follower waits a random **150–300 ms** before standing for election, so split votes are rare.

` + "```mermaid\n%% title: A Raft election\nflowchart LR\n  Follower --> Candidate --> Leader\n```" + `

## Sources

- [In Search of an Understandable Consensus Algorithm](https://raft.github.io/raft.pdf)
`

// ValidComparison returns a Comparison Explainer that satisfies the whole
// editorial skeleton, including a Markdown table under ## Trade-offs.
func ValidComparison() Explainer {
	return Explainer{
		Vertical:    "databases-storage",
		Slug:        "btree-vs-lsm-tree",
		Title:       "B-trees vs LSM-trees",
		Description: "B-trees favour reads and in-place updates; LSM-trees favour write throughput at the cost of compaction.",
		Kind:        "comparison",
		Tags:        []string{"storage-engines"},
		Published:   "2026-10-02T09:00:00Z",
		Reviewed:    "2026-10-16T09:00:00Z",
		Author:      DefaultAuthor().Slug,
		Hero:        "btree-vs-lsm-tree.png",
		HeroAlt:     "A B-tree beside an LSM-tree's levels",
		HeroImage:   HeroPNG(1600, 900),
		Body:        ValidComparisonBody,
	}
}

// ValidComparisonTradeOffs is the ## Trade-offs section of ValidComparisonBody.
const ValidComparisonTradeOffs = `## Trade-offs

| | B-tree | LSM-tree |
|---|---|---|
| Writes | Random, in place | Sequential, batched |
| Reads | One tree walk | Several levels plus Bloom filters |
`

// ValidComparisonBody is the Markdown body of ValidComparison.
const ValidComparisonBody = `Storage engines trade read cost against write cost.

## TL;DR

- B-trees update pages in place.
- LSM-trees buffer writes and merge sorted runs later.
- Pick by read/write ratio and tolerance for compaction.

## Writes take different paths

` + "```mermaid\n%% title: Write paths\nflowchart LR\n  Write --> Memtable --> SSTable\n```" + `

` + ValidComparisonTradeOffs + `
## Sources

- [The Log-Structured Merge-Tree](https://www.cs.umb.edu/~poneil/lsmtree.pdf)
`

// Markdown renders the fixture as a content file.
func (e Explainer) Markdown() string {
	var b strings.Builder
	b.WriteString("---\n")
	field(&b, "title", e.Title)
	field(&b, "description", e.Description)
	field(&b, "kind", e.Kind)
	field(&b, "vertical", e.Vertical)
	list(&b, "tags", e.Tags)
	field(&b, "published", e.Published)
	field(&b, "reviewed", e.Reviewed)
	field(&b, "author", e.Author)
	field(&b, "hero", e.Hero)
	field(&b, "hero_alt", e.HeroAlt)
	list(&b, "related", e.Related)
	if e.Draft {
		b.WriteString("draft: true\n")
	}
	b.WriteString("---\n\n")
	b.WriteString(e.Body)
	return b.String()
}

// Author is a fixture Author.
type Author struct {
	Slug string
	Name string
	Body string
}

// DefaultAuthor is the Author every default fixture references.
func DefaultAuthor() Author {
	return Author{Slug: "ada", Name: "Ada Lovelace", Body: "Writes about distributed systems.\n"}
}

// Markdown renders the fixture as a content file.
func (a Author) Markdown() string {
	var b strings.Builder
	b.WriteString("---\n")
	field(&b, "name", a.Name)
	b.WriteString("---\n\n")
	b.WriteString(a.Body)
	return b.String()
}

// HeroPNG returns a solid teal PNG of the given size, for hero fixtures.
func HeroPNG(w, h int) []byte {
	return HeroPNGFilled(w, h, color.RGBA{R: 0x0D, G: 0x94, B: 0x88, A: 0xFF})
}

// HeroPNGFilled returns a solid PNG of the given size and colour, for
// fixtures that need a hero with different content but the same size.
func HeroPNGFilled(w, h int, c color.Color) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// Config returns the builder configuration used by tests: the launch
// Verticals, a fixed site URL, production mode, the repo's Node tools (Build
// adds a per-build cache dir), and fast hero encoding with an image cache
// shared by every test (SharedImageCacheDir), so a fixture hero is encoded
// once, not once per test.
func Config() builder.Config {
	return builder.Config{
		SiteURL:  "https://hastcode.test",
		Mode:     builder.Production,
		ToolsDir: RepoRoot(),
		Images:   builder.ImageConfig{CacheDir: SharedImageCacheDir(), Fast: true},
		Verticals: []builder.Vertical{
			{Slug: "distributed-systems", Name: "Distributed Systems", Order: 1},
			{Slug: "databases-storage", Name: "Databases & Storage", Order: 2},
			{Slug: "apis-networking", Name: "APIs & Networking", Order: 3},
			{Slug: "infrastructure-cloud", Name: "Infrastructure & Cloud", Order: 4},
			{Slug: "ai-ml-systems", Name: "AI/ML Systems", Order: 5},
			{Slug: "languages-runtimes", Name: "Languages & Runtimes", Order: 6},
		},
	}
}

func field(b *strings.Builder, key, value string) {
	if value == "" {
		return
	}
	fmt.Fprintf(b, "%s: %q\n", key, value)
}

func list(b *strings.Builder, key string, values []string) {
	if len(values) == 0 {
		return
	}
	fmt.Fprintf(b, "%s:\n", key)
	for _, v := range values {
		fmt.Fprintf(b, "  - %q\n", v)
	}
}

// SharedImageCacheDir is the hero image cache Config uses: one directory
// under the OS temp dir, reused across tests and test runs. The cache is
// keyed on source bytes and encoder settings, so sharing it is safe. Tests
// about caching itself set Images.CacheDir to their own t.TempDir().
func SharedImageCacheDir() string {
	return filepath.Join(os.TempDir(), "hastcode-buildtest-image-cache")
}

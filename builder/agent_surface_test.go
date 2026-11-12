package builder_test

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildtest"
)

// The Agent Surface (ADR-0011): a Markdown twin beside every Piece page,
// llms.txt and sitemap.md, written for Agents rather than Readers.

// readStatic returns the text file written for urlPath (a twin, llms.txt,
// sitemap.md).
func readStatic(t *testing.T, out, urlPath string) string {
	t.Helper()
	data, err := os.ReadFile(buildtest.StaticPath(out, urlPath))
	if err != nil {
		t.Fatalf("nothing written at %s: %v", urlPath, err)
	}
	return string(data)
}

func TestAgentSurface_EveryPieceHasAMarkdownTwinWithAHeaderAndItsWrittenBody(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree().WithBrief(buildtest.ValidBrief()))

	twin := readStatic(t, out, "/distributed-systems/raft-leader-election.md")
	for _, want := range []string{
		"# How Raft elects a leader\n",
		"Canonical: https://hastcode.test/distributed-systems/raft-leader-election/\n",
		"Published: 2026-10-01\n",
		"Reviewed: 2026-10-15\n",
		"Author: Ada Lovelace\n",
		"© HastCode. Quote with a link to https://hastcode.test/distributed-systems/raft-leader-election/.\n",
		"## Elections use randomised timeouts\n",
		"```mermaid\n%% title: A Raft election\nflowchart LR\n",
		"## Sources\n",
	} {
		if !strings.Contains(twin, want) {
			t.Errorf("Explainer twin lacks %q", want)
		}
	}
	if strings.HasPrefix(twin, "---") || strings.Contains(twin, "\ntitle:") {
		t.Error("Explainer twin carries frontmatter; it should start with the generated header")
	}

	brief := readStatic(t, out, "/briefs/etcd-raft-rewrite.md")
	for _, want := range []string{
		"# etcd ships a rewritten Raft library\n",
		"Canonical: https://hastcode.test/briefs/etcd-raft-rewrite/\n",
		"Published: 2026-10-20\n",
		"Source: https://etcd.io/blog/2026/raft-rewrite/\n",
		"## What happened\n",
	} {
		if !strings.Contains(brief, want) {
			t.Errorf("Brief twin lacks %q", want)
		}
	}
}

func TestAgentSurface_ATwinFollowsItsPage_DraftsOnlyInPreviews(t *testing.T) {
	draft := buildtest.ValidExplainer()
	draft.Slug, draft.Draft = "draft-piece", true
	tree := buildtest.ValidTree().WithExplainer(draft)

	out := buildtest.MustBuild(t, buildtest.Config(), tree)
	if buildtest.Exists(out, "/distributed-systems/draft-piece.md") {
		t.Error("production wrote a twin for a draft")
	}
	cfg := buildtest.Config()
	cfg.Mode = builder.Preview
	out = buildtest.MustBuild(t, cfg, tree)
	if !buildtest.Exists(out, "/distributed-systems/draft-piece.md") {
		t.Error("preview did not write the draft's twin beside its page")
	}
}

func TestAgentSurface_TwinLinksAreAbsoluteSoAnAgentCanFollowThemAlone(t *testing.T) {
	linker := buildtest.ValidExplainer()
	linker.Slug, linker.Title = "raft-log-replication", "How Raft replicates the log"
	linker.Body = strings.Replace(buildtest.ValidExplainerBody,
		"## Sources",
		"Leader election is covered in [the election Concept](/distributed-systems/raft-leader-election/) and [its summary](/distributed-systems/raft-leader-election/#tldr); see also [the TL;DR above](#tldr) and [the paper](https://raft.github.io/raft.pdf).\n\n## Sources", 1)
	tree := buildtest.ValidTree().WithExplainer(linker).WithBrief(buildtest.ValidBrief())
	out := buildtest.MustBuild(t, buildtest.Config(), tree)

	twin := readStatic(t, out, "/distributed-systems/raft-log-replication.md")
	for _, want := range []string{
		"[the election Concept](https://hastcode.test/distributed-systems/raft-leader-election/)",
		"[its summary](https://hastcode.test/distributed-systems/raft-leader-election/#tldr)",
		"[the TL;DR above](#tldr)",
		"[the paper](https://raft.github.io/raft.pdf)",
	} {
		if !strings.Contains(twin, want) {
			t.Errorf("Explainer twin lacks %q", want)
		}
	}
	if brief := readStatic(t, out, "/briefs/etcd-raft-rewrite.md"); !strings.Contains(brief, "[how Raft elects a leader](https://hastcode.test/distributed-systems/raft-leader-election/)") {
		t.Errorf("Brief twin keeps the bare slug link; got:\n%s", brief)
	}
}

func TestAgentSurface_PiecePagesLinkTheirTwinInTheHeadAndOtherPagesDoNot(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree().WithBrief(buildtest.ValidBrief()))

	alternate := func(urlPath string) string {
		return buildtest.Page(t, out, urlPath).Find(`head link[rel="alternate"][type="text/markdown"]`).AttrOr("href", "")
	}
	if got, want := alternate("/distributed-systems/raft-leader-election/"), "https://hastcode.test/distributed-systems/raft-leader-election.md"; got != want {
		t.Errorf("Explainer page alternate = %q, want %q", got, want)
	}
	if got, want := alternate("/briefs/etcd-raft-rewrite/"), "https://hastcode.test/briefs/etcd-raft-rewrite.md"; got != want {
		t.Errorf("Brief page alternate = %q, want %q", got, want)
	}
	for _, urlPath := range []string{"/", "/distributed-systems/", "/briefs/", "/authors/ada/"} {
		if got := alternate(urlPath); got != "" {
			t.Errorf("%s links a Markdown twin (%q); only Pieces have one", urlPath, got)
		}
	}
}

func TestAgentSurface_EveryPageFootersAnAgentsRowWithTheSurfaceAndTheFeeds(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())

	for _, urlPath := range []string{"/", "/distributed-systems/raft-leader-election/", "/about/"} {
		var got []string
		buildtest.Page(t, out, urlPath).Find(`footer nav[aria-label="Agents"] a`).Each(func(_ int, a *goquery.Selection) {
			got = append(got, a.AttrOr("href", ""))
		})
		want := []string{"/llms.txt", "/sitemap.md", "/feed.xml", "/briefs/feed.xml"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s footer Agents row links = %q, want %q", urlPath, got, want)
		}
	}
}

func TestAgentSurface_LlmsTxtOrientsAnAgentWithoutListingPieces(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree().WithBrief(buildtest.ValidBrief()))

	llms := readStatic(t, out, "/llms.txt")
	for _, want := range []string{
		"# HastCode\n\n> How real systems are built.",
		"https://hastcode.test/briefs/x/ -> https://hastcode.test/briefs/x.md",
		"## Verticals\n",
		"- [Distributed Systems](https://hastcode.test/distributed-systems/): Consensus, replication and clocks.\n",
		"- [Languages & Runtimes](https://hastcode.test/languages-runtimes/): Compilers, garbage collectors and schedulers.\n",
		"- [Briefs](https://hastcode.test/briefs/)",
		"- [Every live Piece, with its Markdown twin](https://hastcode.test/sitemap.md)",
		"- [RSS, all Pieces](https://hastcode.test/feed.xml)",
		"- [RSS, Briefs](https://hastcode.test/briefs/feed.xml)",
		"- [Editorial Policy](https://hastcode.test/editorial-policy/)",
	} {
		if !strings.Contains(llms, want) {
			t.Errorf("llms.txt lacks %q", want)
		}
	}
	if strings.Contains(llms, "raft-leader-election") {
		t.Error("llms.txt lists a Piece; it stays short and leaves enumeration to sitemap.md")
	}
}

func TestAgentSurface_SitemapMdListsEveryLivePieceByVerticalThenBriefs(t *testing.T) {
	second := buildtest.ValidExplainer()
	second.Slug, second.Title, second.Vertical, second.Published = "write-ahead-logging", "Write-ahead logging", "databases-storage", "2026-10-05T12:00:00Z"
	draft := buildtest.ValidExplainer()
	draft.Slug, draft.Draft = "draft-piece", true
	tree := buildtest.ValidTree().WithExplainer(second).WithExplainer(draft).WithBrief(buildtest.ValidBrief())
	out := buildtest.MustBuild(t, buildtest.Config(), tree)

	sm := readStatic(t, out, "/sitemap.md")
	for _, want := range []string{
		"# HastCode: every live Piece\n",
		"## Distributed Systems\n\n- [How Raft elects a leader](https://hastcode.test/distributed-systems/raft-leader-election.md) · Concept · 2026-10-01\n",
		"## Databases & Storage\n\n- [Write-ahead logging](https://hastcode.test/databases-storage/write-ahead-logging.md) · Concept · 2026-10-05\n",
		"## Briefs\n\n- [etcd ships a rewritten Raft library](https://hastcode.test/briefs/etcd-raft-rewrite.md) · Brief · 2026-10-20\n",
	} {
		if !strings.Contains(sm, want) {
			t.Errorf("sitemap.md lacks %q", want)
		}
	}
	if strings.Contains(sm, "draft-piece") {
		t.Error("sitemap.md lists a draft in production")
	}
	if strings.Contains(sm, "## APIs & Networking") {
		t.Error("sitemap.md has a heading for a Vertical with no live Piece")
	}
}

func TestAgentSurface_TwinsStayOutOfSearch_NoindexOnBothTargetsAndNotInTheXMLSitemap(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree().WithBrief(buildtest.ValidBrief()))

	// Production (Cloudflare): the twin is noindex, its page is not.
	twin := "https://hastcode.com/briefs/etcd-raft-rewrite.md"
	if got := buildtest.CloudflareHeaders(t, out, twin).Get("X-Robots-Tag"); got != "noindex" {
		t.Errorf("%s X-Robots-Tag = %q, want noindex", twin, got)
	}
	if got := buildtest.CloudflareHeaders(t, out, "https://hastcode.com/briefs/etcd-raft-rewrite/").Get("X-Robots-Tag"); got != "" {
		t.Errorf("the Brief page itself got X-Robots-Tag %q", got)
	}
	for _, p := range []string{"/llms.txt", "/sitemap.md"} {
		if got := buildtest.CloudflareHeaders(t, out, "https://hastcode.com"+p).Get("X-Robots-Tag"); got != "noindex" {
			t.Errorf("%s X-Robots-Tag = %q, want noindex", p, got)
		}
	}
	// QA (Vercel): served, and noindex like every QA response.
	srv := buildtest.Serve(t, out)
	if res := srv.Get(t, "/briefs/etcd-raft-rewrite.md"); res.StatusCode != 200 || res.Header.Get("X-Robots-Tag") != "noindex" {
		t.Errorf("GET /briefs/etcd-raft-rewrite.md on QA = %d, X-Robots-Tag %q; want 200 and noindex", res.StatusCode, res.Header.Get("X-Robots-Tag"))
	}
	// Never in the XML sitemaps.
	for _, sm := range []string{"/sitemap.xml", "/news-sitemap.xml"} {
		if strings.Contains(readStatic(t, out, sm), ".md<") {
			t.Errorf("%s lists a Markdown twin", sm)
		}
	}
}

func TestAgentSurface_TwinLinkRewritingNeverTouchesCode(t *testing.T) {
	e := buildtest.ValidExplainer()
	e.Body = strings.Replace(buildtest.ValidExplainerBody, "## Sources",
		"Call `handlers[name](r)` once; in Go:\n\n```go\nout := table[key](req)\n```\n\nThen read [the paper](https://raft.github.io/raft.pdf) and [the TL;DR](/distributed-systems/raft-leader-election/#tldr).\n\n## Sources", 1)
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

	twin := readStatic(t, out, "/distributed-systems/raft-leader-election.md")
	for _, want := range []string{"`handlers[name](r)`", "out := table[key](req)", "[the TL;DR](https://hastcode.test/distributed-systems/raft-leader-election/#tldr)"} {
		if !strings.Contains(twin, want) {
			t.Errorf("twin lacks %q", want)
		}
	}
	if strings.Contains(twin, "raft-leader-election/r)") || strings.Contains(twin, "raft-leader-election/req)") {
		t.Error("twin rewrote a ](…) inside code as a link")
	}
}

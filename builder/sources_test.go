package builder_test

import (
	"strings"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

// Spec "Content model": `sources` is an optional frontmatter list (title,
// url) used in addition to inline links. It is merged into the Sources a
// Reader sees, once per URL.

func TestExplainerSources_FrontmatterSourcesJoinTheSourcesSectionOncePerURL(t *testing.T) {
	e := buildtest.ValidExplainer()
	e.Sources = []buildtest.Source{
		{Title: "Raft paper (extended)", URL: "https://raft.github.io/raft.pdf"}, // already in ## Sources
		{Title: "etcd Raft README", URL: "https://github.com/etcd-io/raft"},
		{Title: "etcd Raft README again", URL: "https://github.com/etcd-io/raft"},
	}
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))
	doc := buildtest.Page(t, out, explainerPath)

	sources := doc.Find("article .sources")
	if n := sources.Find(`a[href="https://github.com/etcd-io/raft"]`).Length(); n != 1 {
		t.Errorf("Sources links the frontmatter source %d times, want 1", n)
	}
	if got := sources.Find(`a[href="https://github.com/etcd-io/raft"]`).Text(); got != "etcd Raft README" {
		t.Errorf("frontmatter source link text = %q, want its title", got)
	}
	if n := sources.Find(`a[href="https://raft.github.io/raft.pdf"]`).Length(); n != 1 {
		t.Errorf("Sources links the paper %d times, want 1 (the body link; the frontmatter duplicate is dropped)", n)
	}

	feed := readFeed(t, out, "/feed.xml")
	if len(feed.Channel.Items) != 1 || !strings.Contains(feed.Channel.Items[0].Content, `href="https://github.com/etcd-io/raft"`) {
		t.Errorf("the full-text feed item lacks the frontmatter source")
	}
}

func TestBriefSources_FrontmatterSourcesAreListedWithTheSourceURLOncePerURL(t *testing.T) {
	b := buildtest.ValidBrief()
	b.Sources = []buildtest.Source{
		{Title: "etcd blog", URL: b.SourceURL}, // same as source_url
		{Title: "Release notes", URL: "https://github.com/etcd-io/raft/releases/tag/v4.0.0"},
	}
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree().WithBrief(b))
	doc := buildtest.Page(t, out, "/briefs/etcd-raft-rewrite/")

	line := doc.Find(".brief-source")
	if n := line.Find(`a[href="` + b.SourceURL + `"]`).Length(); n != 1 {
		t.Errorf("source line links source_url %d times, want 1", n)
	}
	rel := line.Find(`a[href="https://github.com/etcd-io/raft/releases/tag/v4.0.0"]`)
	if rel.Length() != 1 || rel.Text() != "Release notes" {
		t.Errorf("source line does not link the frontmatter source once by its title (found %d)", rel.Length())
	}
}

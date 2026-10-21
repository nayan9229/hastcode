package builder_test

import (
	"testing"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildtest"
)

// Lines marked TODO(author) in content are facts only the Author can
// supply. Production must never publish one; previews show them, warned.

func treeWithPlaceholderAbout() buildtest.Tree {
	tree := buildtest.ValidTree()
	about := buildtest.TransparencyPages()[0]
	about.Body = "TODO(author): who publishes HastCode.\n"
	return tree.WithPage(about)
}

func TestAuthorPlaceholder_InAPageFailsProductionNamingTheSourceFile(t *testing.T) {
	_, res := buildtest.Build(t, buildtest.Config(), treeWithPlaceholderAbout())

	buildtest.RequireErrorRule(t, res, "pages/about.md", "author-placeholder", "/about/")
}

func TestAuthorPlaceholder_InAnAuthorNameFailsProductionForArticlesAndFeeds(t *testing.T) {
	a := buildtest.DefaultAuthor()
	a.Name = "TODO(author): your full real name"
	_, res := buildtest.Build(t, buildtest.Config(), buildtest.ValidTree().WithAuthor(a))

	// The byline and Person JSON-LD carry the name onto the Explainer, and
	// dc:creator into the feeds; both point back at the Author file.
	buildtest.RequireErrorRule(t, res, "authors/ada.md", "author-placeholder", "/distributed-systems/raft-leader-election/")
	buildtest.RequireErrorRule(t, res, "authors/ada.md", "author-placeholder", "/feed.xml")
}

func TestAuthorPlaceholder_OnlyWarnsInPreview(t *testing.T) {
	cfg := buildtest.Config()
	cfg.Mode = builder.Preview
	out, res := buildtest.Build(t, cfg, treeWithPlaceholderAbout())

	if res.Failed() {
		t.Fatalf("preview build failed:\n%s", res)
	}
	buildtest.RequireWarningRule(t, res, "pages/about.md", "author-placeholder", "/about/")
	if !buildtest.Exists(out, "/about/") {
		t.Errorf("preview did not write /about/")
	}
}

func TestAuthorPlaceholder_CleanContentHasNone(t *testing.T) {
	_, res := buildtest.Build(t, buildtest.Config(), buildtest.ValidTree().WithBrief(buildtest.ValidBrief()))

	if res.Failed() {
		t.Fatalf("build failed:\n%s", res)
	}
	for _, p := range res.Warnings {
		if p.Rule == "author-placeholder" {
			t.Errorf("unexpected warning: %v", p)
		}
	}
}

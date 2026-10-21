package builder_test

import (
	"strings"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

// The broken internal link gate (spec user story 31): after every page is
// written, each internal href in each page must resolve to a page, an
// asset, a feed or sitemap, the subscribe function or the search page.

func TestBrokenLink_InAStaticPageFailsTheBuildAndDropsThePage(t *testing.T) {
	about := buildtest.TransparencyPages()[0]
	about.Body = "Read [our pipeline](/how-we-build/) first.\n"

	out, res := buildtest.Build(t, buildtest.Config(), buildtest.ValidTree().WithPage(about))

	buildtest.RequireErrorRule(t, res, "pages/about.md", "broken-link", "/how-we-build/")
	buildtest.RequireErrorRule(t, res, "pages/about.md", "broken-link", "/about/")
	if buildtest.Exists(out, "/about/") {
		t.Error("/about/ was written with a broken link")
	}
}

func TestBrokenLink_InAnExplainerBodyFailsTheBuild(t *testing.T) {
	cases := map[string]string{
		"missing Explainer":             "/databases-storage/raft-log-compaction/",
		"Explainer in another Vertical": "/databases-storage/raft-leader-election/",
		"missing hub page":              "/distributed-systems/page/2/",
		"absolute same-site URL":        "https://hastcode.test/nowhere/",
		"relative path":                 "../raft-log-compaction/",
	}
	for name, dest := range cases {
		t.Run(name, func(t *testing.T) {
			e := buildtest.ValidExplainer()
			e.Body = strings.Replace(e.Body, "so split votes are rare.", "so split votes are rare ([more]("+dest+")).", 1)

			out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(e))

			buildtest.RequireErrorRule(t, res, "explainers/distributed-systems/raft-leader-election.md", "broken-link", dest)
			if buildtest.Exists(out, raftURL) {
				t.Errorf("%s was written with a broken link", raftURL)
			}
		})
	}
}

func TestBrokenLink_SiteChromeLinkToAMissingPageFailsTheBuild(t *testing.T) {
	tree := buildtest.ValidTree().Without("pages/contact.md")

	_, res := buildtest.Build(t, buildtest.Config(), tree)

	// The footer links /contact/ from every page; name a page and the target.
	buildtest.RequireErrorRule(t, res, "", "broken-link", "/contact/")
	buildtest.RequireErrorRule(t, res, "", "broken-link", "page / links")
}

func TestBrokenLink_ResolvableInternalLinksPass(t *testing.T) {
	about := buildtest.TransparencyPages()[0]
	about.Body = strings.Join([]string{
		"- [an Explainer](/distributed-systems/raft-leader-election/#tldr)",
		"- [a hub](/databases-storage/)",
		"- [a kind page](/explainers/comparison/)",
		"- [the Author](/authors/ada/)",
		"- [contact, relative](../contact/)",
		"- [Briefs](/briefs/)",
		"- [the site feed](/feed.xml), [a Vertical feed](/distributed-systems/feed.xml), [Briefs feed](/briefs/feed.xml)",
		"- [sitemap](/sitemap.xml), [news sitemap](/news-sitemap.xml)",
		"- [search](/search/?q=raft)",
		"- [subscribe](/api/subscribe)",
		"- [the mark](/mark.svg)",
		"- [absolute home](https://hastcode.test/)",
		"- [elsewhere](https://raft.github.io/), [mail](mailto:a@example.com), [top](#main)",
	}, "\n") + "\n"

	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree().WithPage(about))

	if !buildtest.Exists(out, "/about/") {
		t.Error("/about/ was not written")
	}
}

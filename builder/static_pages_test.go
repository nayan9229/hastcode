package builder_test

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/internal/buildtest"
)

func treeWithTransparencyPages() buildtest.Tree {
	tree := buildtest.NewTree()
	for _, p := range buildtest.TransparencyPages() {
		tree.WithPage(p)
	}
	return tree
}

func TestStaticPages_PublishAtSlugPathWithTitleBodyAndBreadcrumb(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), treeWithTransparencyPages())

	for _, tc := range []struct{ url, title, body string }{
		{"/about/", "About Hast Code", "Hast Code explains how real systems are built."},
		{"/contact/", "Contact", "Email the Author."},
		{"/editorial-policy/", "Editorial Policy", "Every claim links a first-party source."},
	} {
		doc := buildtest.Page(t, out, tc.url)
		if got := strings.TrimSpace(doc.Find("main h1").Text()); got != tc.title {
			t.Errorf("%s h1 = %q, want %q", tc.url, got, tc.title)
		}
		if got := doc.Find("head title").Text(); got != tc.title+" | Hast Code" {
			t.Errorf("%s <title> = %q, want %q", tc.url, got, tc.title+" | Hast Code")
		}
		if got := strings.TrimSpace(doc.Find("main article p").First().Text()); got != tc.body {
			t.Errorf("%s body = %q, want %q", tc.url, got, tc.body)
		}
		if got := doc.Find(`head link[rel="canonical"]`).AttrOr("href", ""); got != "https://hastcode.test"+tc.url {
			t.Errorf("%s canonical = %q", tc.url, got)
		}

		var crumbs []string
		doc.Find(`nav[aria-label="Breadcrumb"] li`).Each(func(_ int, li *goquery.Selection) {
			crumbs = append(crumbs, strings.TrimSpace(li.Text())+" "+li.Find("a").AttrOr("href", ""))
		})
		if got, want := strings.Join(crumbs, "|"), "Home /|"+tc.title+" "; got != want {
			t.Errorf("%s breadcrumbs = %q, want %q", tc.url, got, want)
		}
		list := buildtest.JSONLDOfType(t, doc, "BreadcrumbList")
		items, _ := list["itemListElement"].([]any)
		if len(items) != 2 {
			t.Fatalf("%s BreadcrumbList = %#v, want Home and the page", tc.url, list)
		}
		if last, _ := items[1].(map[string]any); last["name"] != tc.title || last["item"] != "https://hastcode.test"+tc.url {
			t.Errorf("%s BreadcrumbList last item = %#v", tc.url, last)
		}
	}
}

func TestStaticPages_MissingTitleIsABuildError(t *testing.T) {
	p := buildtest.TransparencyPages()[0]
	p.Title = ""

	out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithPage(p))

	buildtest.RequireError(t, res, "pages/about.md", "title")
	if buildtest.Exists(out, "/about/") {
		t.Error("page was written without a title")
	}
}

// A static page lives at /<slug>/, so its slug must not take a path the
// site already uses (a Vertical hub, authors, briefs, search…).
func TestStaticPages_SlugThatTakesAnotherPagesPathIsABuildError(t *testing.T) {
	for _, slug := range []string{"databases-storage", "authors", "briefs", "search"} {
		p := buildtest.TransparencyPages()[0]
		p.Slug = slug

		_, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithPage(p))

		buildtest.RequireErrorRule(t, res, "pages/"+slug+".md", "page-path", "/"+slug+"/")
	}
}

package builder_test

import (
	"fmt"
	"image/color"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/internal/buildtest"
)

const adaURL = "/authors/ada/"

func TestAuthorPage_ShowsNameBioAndProfileLinks(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	doc := buildtest.Page(t, out, adaURL)

	card := doc.Find("main .h-card")
	if got := strings.TrimSpace(card.Find("h1.p-name").Text()); got != "Ada Lovelace" {
		t.Errorf("h1 = %q, want the Author's name", got)
	}
	if got := strings.TrimSpace(card.Find(".p-note").Text()); got != "Writes about distributed systems." {
		t.Errorf("bio = %q, want the Author file's body", got)
	}
	var links []string
	card.Find("a.u-url[rel~=me]").Each(func(_ int, a *goquery.Selection) { links = append(links, a.AttrOr("href", "")) })
	if got, want := strings.Join(links, " "), "https://github.com/ada-example https://www.linkedin.com/in/ada-example"; got != want {
		t.Errorf("profile links = %q, want %q", got, want)
	}
	if got := doc.Find("head title").Text(); got != "Ada Lovelace | HastCode" {
		t.Errorf("<title> = %q", got)
	}
}

// The Author page's Person is the same entity an Explainer's Article names
// as its author: identical url and sameAs.
func TestAuthorPage_PersonJSONLDMatchesTheExplainerAuthor(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	person := buildtest.JSONLDOfType(t, buildtest.Page(t, out, adaURL), "Person")
	if person == nil {
		t.Fatal("no Person JSON-LD block on the Author page")
	}
	article := buildtest.JSONLDOfType(t, buildtest.Page(t, out, raftURL), "Article")
	author, _ := article["author"].(map[string]any)

	if person["@context"] != "https://schema.org" || person["name"] != "Ada Lovelace" {
		t.Errorf("Person = %#v, want @context https://schema.org and name Ada Lovelace", person)
	}
	if person["url"] != "https://hastcode.test/authors/ada/" || person["url"] != author["url"] {
		t.Errorf("Person url = %#v, Article author url = %#v, want both https://hastcode.test/authors/ada/", person["url"], author["url"])
	}
	if got, want := fmt.Sprint(person["sameAs"]), fmt.Sprint(author["sameAs"]); got != want || got == "[]" || got == "<nil>" {
		t.Errorf("Person sameAs = %s, Article author sameAs = %s, want the same non-empty list", got, want)
	}
}

func TestAuthorPage_ListsTheAuthorsLatestPiecesNewestFirst(t *testing.T) {
	grace := buildtest.Author{Slug: "grace", Name: "Grace Hopper", Body: "Writes about compilers.\n"}
	byGrace := buildtest.ValidExplainer()
	byGrace.Slug = "raft-log-replication"
	byGrace.Title = "How Raft replicates its log"
	byGrace.Author = "grace"
	tree := buildtest.ValidTree().WithExplainer(buildtest.ValidComparison()).WithAuthor(grace).WithExplainer(byGrace)
	out := buildtest.MustBuild(t, buildtest.Config(), tree)

	var got []string
	buildtest.Page(t, out, adaURL).Find("main .latest li").Each(func(_ int, li *goquery.Selection) {
		a := li.Find("a")
		got = append(got, strings.TrimSpace(a.Text())+" "+a.AttrOr("href", ""))
	})
	want := []string{
		"B-trees vs LSM-trees /databases-storage/btree-vs-lsm-tree/", // published 2 Oct
		"How Raft elects a leader " + raftURL,                        // published 1 Oct
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("latest pieces = %q, want %q (Ada's only, newest first)", got, want)
	}
}

func TestExplainerByline_LinksTheAuthorPage(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	author := buildtest.Page(t, out, raftURL).Find("article.h-entry .byline a.p-author")

	if author.Length() != 1 {
		t.Fatalf("found %d linked p-author bylines, want 1", author.Length())
	}
	if got := author.AttrOr("href", ""); got != adaURL {
		t.Errorf("byline href = %q, want %s", got, adaURL)
	}
	if !buildtest.Exists(out, author.AttrOr("href", "")) {
		t.Error("byline links a page the build did not write")
	}
}

func authorWithPhoto() buildtest.Author {
	a := buildtest.DefaultAuthor()
	a.Photo = "ada.png"
	a.PhotoAlt = "Ada Lovelace, smiling"
	a.PhotoImage = buildtest.HeroPNG(900, 600)
	return a
}

func TestAuthorPage_ShowsASquarePhotoThatPersonJSONLDNamesAsItsImage(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree().WithAuthor(authorWithPhoto()))
	doc := buildtest.Page(t, out, adaURL)

	img := doc.Find("main .h-card img.u-photo")
	if img.Length() != 1 {
		t.Fatalf("found %d u-photo images, want 1", img.Length())
	}
	if got := img.AttrOr("alt", ""); got != "Ada Lovelace, smiling" {
		t.Errorf("photo alt = %q, want the Author's photo_alt", got)
	}
	src := img.AttrOr("src", "")
	photo := buildtest.DecodeImage(t, out, src)
	if photo.Width != photo.Height || photo.Width < 400 {
		t.Errorf("photo is %dx%d, want a square at least 400 px", photo.Width, photo.Height)
	}
	if img.AttrOr("width", "") == "" || img.AttrOr("height", "") == "" {
		t.Error("photo has no width/height attributes (layout shift)")
	}

	person := buildtest.JSONLDOfType(t, doc, "Person")
	if got := person["image"]; got != "https://hastcode.test"+src {
		t.Errorf("Person image = %#v, want https://hastcode.test%s", got, src)
	}
}

func TestAuthorPage_MissingPhotoFileIsABuildError(t *testing.T) {
	a := authorWithPhoto()
	a.PhotoImage = nil

	out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithAuthor(a))

	buildtest.RequireErrorRule(t, res, "authors/ada.md", "author-photo", "ada.png")
	if buildtest.Exists(out, adaURL) {
		t.Error("Author page was written with a missing photo")
	}
}

// WebP photos fail like WebP heroes (sniffed from the bytes, whatever the
// name), with an error that says what to do instead.
func TestAuthorPage_WebPPhotoIsABuildErrorAskingForPNGOrJPEG(t *testing.T) {
	for name, file := range map[string]string{"webp extension": "ada.webp", "png extension": "ada.png"} {
		t.Run(name, func(t *testing.T) {
			a := authorWithPhoto()
			a.Photo = file
			a.PhotoImage = buildtest.HeroWebP(900, 600, color.RGBA{R: 0x0D, G: 0x94, B: 0x88, A: 0xFF})

			out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithAuthor(a))

			buildtest.RequireErrorRule(t, res, "authors/ada.md", "author-photo", "WebP")
			buildtest.RequireErrorRule(t, res, "authors/ada.md", "author-photo", "PNG or JPEG")
			if buildtest.Exists(out, adaURL) {
				t.Error("Author page was written with a WebP photo")
			}
		})
	}
}

func TestAuthorPage_WithoutAPhotoWarns(t *testing.T) {
	_, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree())
	buildtest.RequireWarningRule(t, res, "authors/ada.md", "author-photo", "no photo")
}

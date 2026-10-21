package builder_test

import (
	"testing"

	"hastcode.com/site/internal/buildtest"
)

// organizationAuthor is the default fixture Author (same slug, so every
// default Explainer names it) published under a team byline.
func organizationAuthor() buildtest.Author {
	a := buildtest.DefaultAuthor()
	a.Kind = "organization"
	a.Name = "HastCode"
	return a
}

// An Author whose kind is organization is a team byline, not a person:
// the Explainer's Article names an Organization author and the Author page
// carries a top-level Organization, both with the Author page as url.
func TestOrganizationAuthor_IsAnOrganizationInJSONLD(t *testing.T) {
	tree := buildtest.ValidTree().WithAuthor(organizationAuthor())
	out := buildtest.MustBuild(t, buildtest.Config(), tree)

	article := buildtest.JSONLDOfType(t, buildtest.Page(t, out, raftURL), "Article")
	author, _ := article["author"].(map[string]any)
	if author["@type"] != "Organization" || author["name"] != "HastCode" || author["url"] != "https://hastcode.test/authors/ada/" {
		t.Errorf("Article author = %#v, want Organization HastCode with url https://hastcode.test/authors/ada/", author)
	}

	page := buildtest.Page(t, out, adaURL)
	if p := buildtest.JSONLDOfType(t, page, "Person"); p != nil {
		t.Errorf("Author page has a Person block %#v, want none for an organization", p)
	}
	org := buildtest.JSONLDOfType(t, page, "Organization")
	if org == nil || org["@context"] != "https://schema.org" || org["name"] != "HastCode" || org["url"] != "https://hastcode.test/authors/ada/" {
		t.Errorf("Author page Organization = %#v, want top-level HastCode with the Author page url", org)
	}
}

// A team byline has no headshot, so a missing photo is not worth a warning.
func TestOrganizationAuthor_WithoutAPhotoDoesNotWarn(t *testing.T) {
	_, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithAuthor(organizationAuthor()))
	buildtest.RequireNoProblemRule(t, res, "authors/ada.md", "author-photo")
}

func TestAuthor_UnknownKindIsABuildError(t *testing.T) {
	a := buildtest.DefaultAuthor()
	a.Kind = "robot"
	_, res := buildtest.Build(t, buildtest.Config(), buildtest.ValidTree().WithAuthor(a))
	buildtest.RequireErrorRule(t, res, "authors/ada.md", "author-kind", "robot")
}

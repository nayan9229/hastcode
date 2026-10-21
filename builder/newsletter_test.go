package builder_test

import (
	"strings"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

func TestNewsletter_FormEndsTheExplainerAndPostsToSubscribe(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	doc := buildtest.Page(t, out, "/distributed-systems/raft-leader-election/")

	forms := doc.Find(`form.newsletter-form`)
	if forms.Length() != 1 {
		t.Fatalf("found %d newsletter forms, want 1", forms.Length())
	}
	if forms.ParentsFiltered("article.h-entry").Length() == 0 {
		t.Errorf("newsletter form is not inside the Explainer article")
	}
	if forms.Closest(".e-content").Length() != 0 {
		t.Errorf("newsletter form is inside the article body, want it after the body")
	}
	if after := doc.Find("article.h-entry .e-content").NextAllFiltered(".newsletter").Length(); after != 1 {
		t.Errorf("newsletter section after the article body: found %d, want 1", after)
	}

	form := forms.First()
	if got := form.AttrOr("action", ""); got != "/api/subscribe" {
		t.Errorf("form action = %q, want /api/subscribe", got)
	}
	if got := strings.ToLower(form.AttrOr("method", "")); got != "post" {
		t.Errorf("form method = %q, want post", got)
	}

	email := form.Find(`input[name="email"]`)
	if email.AttrOr("type", "") != "email" || !email.Is("[required]") {
		t.Errorf("email input = type %q required %v, want type=email required", email.AttrOr("type", ""), email.Is("[required]"))
	}
	if id := email.AttrOr("id", ""); id == "" || form.Find(`label[for="`+id+`"]`).Length() != 1 {
		t.Errorf("email input has no associated <label>")
	}

	source := form.Find(`input[type="hidden"][name="source"]`)
	if got := source.AttrOr("value", ""); got != "/distributed-systems/raft-leader-election/" {
		t.Errorf("hidden source = %q, want the Explainer's path", got)
	}

	hp := form.Find(`input[name="website"]`)
	if hp.Length() != 1 {
		t.Fatalf("honeypot field: found %d, want 1", hp.Length())
	}
	if hp.AttrOr("tabindex", "") != "-1" || hp.AttrOr("autocomplete", "") != "off" {
		t.Errorf("honeypot tabindex/autocomplete = %q/%q, want -1/off", hp.AttrOr("tabindex", ""), hp.AttrOr("autocomplete", ""))
	}
	if hp.Closest(`[aria-hidden="true"]`).Length() == 0 {
		t.Errorf("honeypot is not inside an aria-hidden wrapper")
	}

	if form.Find(`button[type="submit"]`).Length() != 1 {
		t.Errorf("form has no submit button")
	}
	if form.Find(`[role="status"]`).Length() != 1 {
		t.Errorf("form has no status region for the scripted result")
	}
}

func TestNewsletter_EnhancementScriptIsSmall(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	doc := buildtest.Page(t, out, "/distributed-systems/raft-leader-election/")

	script := doc.Find(".newsletter script")
	if script.Length() != 1 {
		t.Fatalf("found %d newsletter scripts, want 1 inline script", script.Length())
	}
	if src := script.AttrOr("src", ""); src != "" {
		t.Errorf("newsletter script loads %q, want it inline", src)
	}
	if n := len(script.Text()); n > 1024 {
		t.Errorf("newsletter script is %d bytes, want ≤ 1 KB (JS budget)", n)
	}
}

func TestNewsletter_SubscribedPageExistsAndIsNoindex(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	doc := buildtest.Page(t, out, "/subscribed/")

	if got := doc.Find(`meta[name="robots"]`).AttrOr("content", ""); !strings.Contains(got, "noindex") {
		t.Errorf("meta robots = %q, want noindex", got)
	}
	if got := doc.Find("h1").Text(); !strings.Contains(strings.ToLower(got), "inbox") {
		t.Errorf("h1 = %q, want it to tell the Reader to check their inbox", got)
	}
	if doc.Find(`main a[href="/"]`).Length() == 0 {
		t.Errorf("no link back to the home page")
	}
}

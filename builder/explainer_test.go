package builder_test

import (
	"strings"
	"testing"
	"time"

	"hastcode.com/site/internal/buildtest"
)

func TestExplainer_IsPublishedAtVerticalSlugPath(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())

	doc := buildtest.Page(t, out, "/distributed-systems/raft-leader-election/")

	if got := doc.Find("html").AttrOr("lang", ""); got != "en" {
		t.Errorf("html lang = %q, want en", got)
	}
	if got := doc.Find("head title").Text(); !strings.Contains(got, "How Raft elects a leader") {
		t.Errorf("<title> = %q, want it to contain the Explainer title", got)
	}
	if got := doc.Find(`meta[name="description"]`).AttrOr("content", ""); got != "Raft keeps a replicated log consistent by electing one leader per term with randomised timeouts." {
		t.Errorf("meta description = %q", got)
	}
	if got := doc.Find(`meta[name="viewport"]`).AttrOr("content", ""); !strings.Contains(got, "width=device-width") {
		t.Errorf("meta viewport = %q", got)
	}

	article := doc.Find("article.h-entry")
	if got := strings.TrimSpace(article.Find("h1.p-name").Text()); got != "How Raft elects a leader" {
		t.Errorf("h1 = %q", got)
	}
	if got := strings.TrimSpace(article.Find(".p-category").Text()); got != "Distributed Systems" {
		t.Errorf("Vertical name = %q, want Distributed Systems", got)
	}
	if got := strings.TrimSpace(article.Find(".p-author").Text()); got != "Ada Lovelace" {
		t.Errorf("Author name = %q, want Ada Lovelace", got)
	}
}

func TestExplainer_ShowsPublishedAndLastReviewedDates(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	doc := buildtest.Page(t, out, "/distributed-systems/raft-leader-election/")

	for _, tc := range []struct {
		selector string
		want     string
		visible  string
	}{
		{"time.dt-published", "2026-10-01T09:00:00+05:30", "1 October 2026"},
		{"time.dt-updated", "2026-10-15T18:30:00Z", "15 October 2026"},
	} {
		el := doc.Find(tc.selector)
		if el.Length() != 1 {
			t.Fatalf("%s: found %d elements, want 1", tc.selector, el.Length())
		}
		got, err := time.Parse(time.RFC3339, el.AttrOr("datetime", ""))
		want, _ := time.Parse(time.RFC3339, tc.want)
		if err != nil || !got.Equal(want) {
			t.Errorf("%s datetime = %q, want instant %s", tc.selector, el.AttrOr("datetime", ""), tc.want)
		}
		if text := strings.TrimSpace(el.Text()); text != tc.visible {
			t.Errorf("%s text = %q, want %q", tc.selector, text, tc.visible)
		}
	}
}

func TestExplainer_BodyMarkdownIsRenderedToHTML(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	body := buildtest.Page(t, out, "/distributed-systems/raft-leader-election/").Find(".e-content")

	if got := body.Find("h2").First().Text(); got != "TL;DR" {
		t.Errorf("first h2 = %q, want TL;DR", got)
	}
	if got := body.Find("h2").First().Next().Find("li").Length(); got != 3 {
		t.Errorf("TL;DR bullets = %d, want 3", got)
	}
	if got := body.Find("strong").Text(); got != "150–300 ms" {
		t.Errorf("<strong> = %q", got)
	}
	if got := body.Find(`a[href="https://raft.github.io/raft.pdf"]`).Length(); got != 1 {
		t.Errorf("Sources link count = %d, want 1", got)
	}
}

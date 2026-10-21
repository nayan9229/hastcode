package builder_test

import (
	"testing"

	"hastcode.com/site/internal/buildtest"
)

func TestHomeJSONLD_WebSiteNamesHastCodeWithASearchAction(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	ws := buildtest.JSONLDOfType(t, buildtest.Page(t, out, "/"), "WebSite")
	if ws == nil {
		t.Fatal("home has no WebSite JSON-LD")
	}

	for key, want := range map[string]string{
		"@context":    "https://schema.org",
		"name":        "HastCode",
		"url":         siteURL + "/",
		"description": "How real systems are built.",
	} {
		if ws[key] != want {
			t.Errorf("WebSite %s = %v, want %q", key, ws[key], want)
		}
	}
	action, _ := ws["potentialAction"].(map[string]any)
	target, _ := action["target"].(map[string]any)
	if action["@type"] != "SearchAction" || target["urlTemplate"] != siteURL+"/search/?q={search_term_string}" || action["query-input"] != "required name=search_term_string" {
		t.Errorf("WebSite potentialAction = %v, want a SearchAction on /search/?q={search_term_string}", ws["potentialAction"])
	}
}

func TestHomeJSONLD_OrganizationHasNameTaglineLogoAndProfiles(t *testing.T) {
	cfg := buildtest.Config()
	cfg.SameAs = []string{"https://github.com/hastcode"}
	out := buildtest.MustBuild(t, cfg, buildtest.ValidTree())
	org := buildtest.JSONLDOfType(t, buildtest.Page(t, out, "/"), "Organization")
	if org == nil {
		t.Fatal("home has no Organization JSON-LD")
	}

	for key, want := range map[string]string{
		"@context":    "https://schema.org",
		"name":        "HastCode",
		"url":         siteURL + "/",
		"description": "How real systems are built.",
	} {
		if org[key] != want {
			t.Errorf("Organization %s = %v, want %q", key, org[key], want)
		}
	}
	logo, _ := org["logo"].(map[string]any)
	if logo["url"] != siteURL+"/mark.svg" {
		t.Errorf("Organization logo = %v, want the absolute /mark.svg URL", org["logo"])
	}
	if !buildtest.Exists(out, "/mark.svg") {
		t.Error("the logo file /mark.svg was not written")
	}
	if same, _ := org["sameAs"].([]any); len(same) != 1 || same[0] != "https://github.com/hastcode" {
		t.Errorf("Organization sameAs = %v, want the configured profiles", org["sameAs"])
	}
}

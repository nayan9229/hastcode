package builder_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildtest"
)

const indexNowKey = "3f6c1e0a9b2d4c7e8f1a2b3c4d5e6f70"

func indexNowConfig(t *testing.T) builder.Config {
	cfg := buildtest.Config()
	cfg.IndexNowKey = indexNowKey
	cfg.Now = buildClock
	return cfg
}

// indexNowBody is the changed-URL artefact: an IndexNow JSON POST body
// (indexnow.org/documentation) the post-deploy step sends as is.
type indexNowBody struct {
	Host        string   `json:"host"`
	Key         string   `json:"key"`
	KeyLocation string   `json:"keyLocation"`
	URLList     []string `json:"urlList"`
}

func readIndexNow(t *testing.T, out string) (indexNowBody, bool) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(out, "indexnow.json"))
	if os.IsNotExist(err) {
		return indexNowBody{}, false
	}
	if err != nil {
		t.Fatal(err)
	}
	var b indexNowBody
	if err := json.Unmarshal(data, &b); err != nil {
		t.Fatalf("indexnow.json does not decode: %v\n%s", err, data)
	}
	return b, true
}

func TestIndexNow_KeyFileIsServedFromTheSiteRoot(t *testing.T) {
	out := buildtest.MustBuild(t, indexNowConfig(t), buildtest.ValidTree())

	data, err := os.ReadFile(buildtest.StaticPath(out, "/"+indexNowKey+".txt"))
	if err != nil {
		t.Fatalf("no key file at /%s.txt: %v", indexNowKey, err)
	}
	if got := strings.TrimSpace(string(data)); got != indexNowKey {
		t.Errorf("key file holds %q, want the key", got)
	}
}

func TestIndexNow_ArtefactListsPagesNewOrChangedSinceThePreviousBuild(t *testing.T) {
	cfg := indexNowConfig(t)
	cfg.CacheDir = t.TempDir() // the build cache Vercel keeps between builds

	first := buildtest.MustBuild(t, cfg, buildtest.ValidTree())
	body, ok := readIndexNow(t, first)
	if !ok {
		t.Fatal("first build wrote no IndexNow artefact")
	}
	if !slices.Contains(body.URLList, siteURL+"/distributed-systems/raft-leader-election/") || !slices.Contains(body.URLList, siteURL+"/about/") {
		t.Errorf("first build: every page is new, got %v", body.URLList)
	}

	// Second build: one new Brief, one revised Explainer, the rest as before.
	revised := buildtest.ValidExplainer()
	revised.Body = strings.Replace(revised.Body, "split votes are rare", "split votes are rare and short", 1)
	second := buildtest.MustBuild(t, cfg, buildtest.NewTree().WithExplainer(revised).WithBrief(buildtest.ValidBrief()))
	body, ok = readIndexNow(t, second)
	if !ok {
		t.Fatal("second build wrote no IndexNow artefact")
	}
	for _, want := range []string{
		"/briefs/etcd-raft-rewrite/",                 // new
		"/distributed-systems/raft-leader-election/", // revised
		"/",        // lists the new Brief
		"/briefs/", // lists the new Brief
	} {
		if !slices.Contains(body.URLList, siteURL+want) {
			t.Errorf("second build's urlList lacks %s: %v", want, body.URLList)
		}
	}
	for _, unchanged := range []string{"/about/", "/contact/", "/databases-storage/"} {
		if slices.Contains(body.URLList, siteURL+unchanged) {
			t.Errorf("second build's urlList has unchanged %s", unchanged)
		}
	}
	if body.Host != "hastcode.test" || body.Key != indexNowKey || body.KeyLocation != siteURL+"/"+indexNowKey+".txt" {
		t.Errorf("host/key/keyLocation = %q / %q / %q", body.Host, body.Key, body.KeyLocation)
	}

	// Third build of the same tree: nothing to ping.
	third := buildtest.MustBuild(t, cfg, buildtest.NewTree().WithExplainer(revised).WithBrief(buildtest.ValidBrief()))
	if body, ok := readIndexNow(t, third); ok {
		t.Errorf("unchanged rebuild wrote an artefact listing %v", body.URLList)
	}
}

func TestIndexNow_ArtefactAndManifestStayOutOfThePublicOutput(t *testing.T) {
	cfg := indexNowConfig(t)
	cfg.CacheDir = t.TempDir()
	out := buildtest.MustBuild(t, cfg, buildtest.ValidTree())

	for _, p := range buildtest.URLPaths(t, out) {
		// Pagefind's bundle manifest is public by design (the search page
		// loads it); nothing else may be JSON.
		if strings.HasSuffix(p, ".json") && !strings.HasPrefix(p, "/pagefind/") {
			t.Errorf("public output serves %s", p)
		}
	}
}

func TestIndexNow_NoKeyMeansNoKeyFileAndNoArtefact(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())

	if _, ok := readIndexNow(t, out); ok {
		t.Error("wrote an IndexNow artefact without a key")
	}
	for _, p := range buildtest.URLPaths(t, out) {
		if strings.Count(p, "/") == 1 && strings.HasSuffix(p, ".txt") && p != "/robots.txt" {
			t.Errorf("wrote a key file %s without a key", p)
		}
	}
}

func TestIndexNow_MalformedKeyFailsTheBuild(t *testing.T) {
	cfg := indexNowConfig(t)
	cfg.IndexNowKey = "short"
	_, res := buildtest.Build(t, cfg, buildtest.ValidTree())

	buildtest.RequireErrorRule(t, res, "", "indexnow-key", "8 to 128")
}

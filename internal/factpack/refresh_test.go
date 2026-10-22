package factpack_test

import (
	"bytes"
	"context"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// twoProfiles writes two Repo Profiles with their Fact Packs fetched from
// the fixtures: apache/iggy and nats-io/nats-server.
func twoProfiles(t *testing.T, h *harness) (iggyJSON, natsJSON string) {
	t.Helper()
	h.gh.serveRepo("apache-iggy")
	h.gh.serveRepo("nats-io-nats-server")
	iggy := h.writeProfile("distributed-systems", "apache-iggy", "apache/iggy")
	nats := h.writeProfile("distributed-systems", "nats-server", "nats-io/nats-server")
	for _, p := range []struct{ repo, md string }{{"apache/iggy", iggy}, {"nats-io/nats-server", nats}} {
		if _, err := h.fetcher.Fetch(context.Background(), p.repo, p.md); err != nil {
			t.Fatal(err)
		}
	}
	h.log.Reset()
	h.fetcher.Now = func() time.Time { return fetchedAt.Add(7 * 24 * time.Hour) }
	return strings.TrimSuffix(iggy, ".md") + ".json", strings.TrimSuffix(nats, ".md") + ".json"
}

func TestRefresh_RewritesOnlyTheChangedFactPackAndSaysWhatChanged(t *testing.T) {
	h := newHarness(t)
	iggyJSON, natsJSON := twoProfiles(t, h)
	iggyBefore, natsBefore := readFile(t, iggyJSON), readFile(t, natsJSON)
	// Upstream, apache/iggy gained stars and promoted its release candidate.
	repo := string(readFile(t, filepath.Join("testdata", "github", "apache-iggy", "repo.json")))
	h.gh.answer("/repos/apache/iggy", answer{status: http.StatusOK, body: []byte(strings.Replace(repo, `"stargazers_count": 2468`, `"stargazers_count": 2600`, 1))})
	releases := string(readFile(t, filepath.Join("testdata", "github", "apache-iggy", "releases.json")))
	h.gh.answer(releasesPath, answer{status: http.StatusOK, body: []byte(strings.Replace(releases, `"prerelease": true`, `"prerelease": false`, 1))})

	report, err := h.fetcher.Refresh(context.Background(), h.contentDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(report.Changes) != 1 || report.Changes[0].Repo != "apache/iggy" {
		t.Fatalf("changes = %+v, want exactly apache/iggy", report.Changes)
	}
	if report.Profiles != 2 {
		t.Errorf("profiles checked = %d, want 2", report.Profiles)
	}
	if bytes.Equal(readFile(t, iggyJSON), iggyBefore) {
		t.Error("apache/iggy's Fact Pack was not rewritten")
	}
	fp := parseFactPack(t, iggyJSON)
	if fp.Stars != 2600 || fp.LatestRelease == nil || fp.LatestRelease.Tag != "server-0.6.0-rc1" || !fp.FetchedAt.Equal(fetchedAt.Add(7*24*time.Hour)) {
		t.Errorf("rewritten Fact Pack: stars %d, latest %+v, fetched_at %s", fp.Stars, fp.LatestRelease, fp.FetchedAt)
	}
	if !bytes.Equal(readFile(t, natsJSON), natsBefore) {
		t.Errorf("nats-io/nats-server's Fact Pack changed although nothing changed upstream:\n%s", firstDifference(natsBefore, readFile(t, natsJSON)))
	}

	lines := changeLines(h.log.String())
	if len(lines) != 1 {
		t.Fatalf("want one line for the changed repository, got %d:\n%s", len(lines), h.log)
	}
	for _, want := range []string{"apache/iggy", "stars 2468 -> 2600", "latest release server-0.5.0 -> server-0.6.0-rc1", "explainers/distributed-systems/apache-iggy.json"} {
		if !strings.Contains(lines[0], want) {
			t.Errorf("the change line lacks %q: %s", want, lines[0])
		}
	}
	if strings.Contains(h.log.String(), "nats-io/nats-server:") {
		t.Errorf("an unchanged repository must not be reported:\n%s", h.log)
	}
}

func TestRefresh_NothingChangedLeavesEveryFileAloneAndReportsNone(t *testing.T) {
	h := newHarness(t)
	iggyJSON, natsJSON := twoProfiles(t, h)
	iggyBefore, natsBefore := readFile(t, iggyJSON), readFile(t, natsJSON)

	report, err := h.fetcher.Refresh(context.Background(), h.contentDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(report.Changes) != 0 {
		t.Errorf("changes = %+v, want none", report.Changes)
	}
	if !bytes.Equal(readFile(t, iggyJSON), iggyBefore) || !bytes.Equal(readFile(t, natsJSON), natsBefore) {
		t.Error("an unchanged Fact Pack was rewritten (its fetched_at must not move)")
	}
	if lines := changeLines(h.log.String()); len(lines) != 0 {
		t.Errorf("no change lines expected, got:\n%s", strings.Join(lines, "\n"))
	}
}

func TestRefresh_FetchesAProfileWhoseFactPackIsMissing(t *testing.T) {
	h := newHarness(t)
	h.gh.serveRepo("apache-iggy")
	h.writeProfile("distributed-systems", "apache-iggy", "apache/iggy")

	report, err := h.fetcher.Refresh(context.Background(), h.contentDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(report.Changes) != 1 || report.Changes[0].Repo != "apache/iggy" {
		t.Fatalf("changes = %+v, want apache/iggy as new", report.Changes)
	}
	got := readFile(t, filepath.Join(h.contentDir, "explainers", "distributed-systems", "apache-iggy.json"))
	if !bytes.Equal(got, readFile(t, filepath.Join("testdata", "apache-iggy.json"))) {
		t.Error("the new Fact Pack is not the golden one")
	}
}

func TestRefresh_ContentTreeWithoutProfilesChangesNothing(t *testing.T) {
	h := newHarness(t)

	report, err := h.fetcher.Refresh(context.Background(), h.contentDir)
	if err != nil {
		t.Fatal(err)
	}
	if report.Profiles != 0 || len(report.Changes) != 0 {
		t.Errorf("report = %+v, want empty", report)
	}
}

// changeLines are the log lines naming a repository's changes
// ("owner/name: …").
func changeLines(log string) []string {
	var out []string
	for _, line := range strings.Split(strings.TrimSpace(log), "\n") {
		if repo, _, ok := strings.Cut(line, ": "); ok && strings.Count(repo, "/") == 1 && !strings.Contains(repo, " ") {
			out = append(out, line)
		}
	}
	return out
}

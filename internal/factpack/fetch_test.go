package factpack_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"hastcode.com/site/internal/content"
	"hastcode.com/site/internal/github"
)

// The fetch seam (spec "Testing Decisions"): the Fetcher runs against an
// in-process GitHub serving recorded JSON; the tests assert on the files it
// writes, what it prints, and the errors it stops with.

func TestFetch_WritesTheFactPackBesideTheMarkdown(t *testing.T) {
	h := newHarness(t)
	h.gh.serveRepo("apache-iggy")
	md := h.markdownPath("distributed-systems", "apache-iggy")

	if _, err := h.fetcher.Fetch(context.Background(), "apache/iggy", md); err != nil {
		t.Fatalf("fetch: %v", err)
	}

	got := readFile(t, filepath.Join(filepath.Dir(md), "apache-iggy.json"))
	want := readFile(t, filepath.Join("testdata", "apache-iggy.json"))
	if !bytes.Equal(got, want) {
		t.Fatalf("Fact Pack differs from testdata/apache-iggy.json at %s\n--- got ---\n%s", firstDifference(want, got), got)
	}
	if entries := strings.Count(h.log.String(), "\n"); entries == 0 {
		t.Errorf("expected the fetch to say what it wrote; log is empty")
	}
}

func TestFetch_SecondRunOnTheSameAnswersChangesOnlyFetchedAt(t *testing.T) {
	h := newHarness(t)
	h.gh.serveRepo("apache-iggy")
	md := h.markdownPath("distributed-systems", "apache-iggy")
	jsonFile := filepath.Join(filepath.Dir(md), "apache-iggy.json")

	if _, err := h.fetcher.Fetch(context.Background(), "apache/iggy", md); err != nil {
		t.Fatal(err)
	}
	first := readFile(t, jsonFile)
	h.fetcher.Now = func() time.Time { return fetchedAt.Add(7 * 24 * time.Hour) }
	if _, err := h.fetcher.Fetch(context.Background(), "apache/iggy", md); err != nil {
		t.Fatal(err)
	}
	second := readFile(t, jsonFile)

	if bytes.Equal(first, second) {
		t.Fatal("fetched_at did not move on the second run")
	}
	if !bytes.Equal(withoutFetchedAt(first), withoutFetchedAt(second)) {
		t.Fatalf("second run changed more than fetched_at: %s", firstDifference(withoutFetchedAt(first), withoutFetchedAt(second)))
	}
	if !strings.Contains(string(second), `"fetched_at": "2026-10-26T06:00:00Z"`) {
		t.Errorf("second run's fetched_at is not the new clock:\n%s", second)
	}
}

// withoutFetchedAt drops the fetched_at line.
func withoutFetchedAt(factPack []byte) []byte {
	var keep [][]byte
	for _, line := range bytes.Split(factPack, []byte("\n")) {
		if !bytes.Contains(line, []byte(`"fetched_at"`)) {
			keep = append(keep, line)
		}
	}
	return bytes.Join(keep, []byte("\n"))
}

func TestFetch_NOASSERTIONLicenseNeedsReviewAndIsNoted(t *testing.T) {
	h := newHarness(t)
	h.gh.serveRepo("nats-io-nats-server")
	md := h.markdownPath("distributed-systems", "nats-server")

	res, err := h.fetcher.Fetch(context.Background(), "nats-io/nats-server", md)
	if err != nil {
		t.Fatal(err)
	}

	fp := parseFactPack(t, filepath.Join(filepath.Dir(md), "nats-server.json"))
	if fp.License == nil || !fp.License.NeedsReview || fp.License.SPDXID != "NOASSERTION" {
		t.Errorf("license = %+v, want NOASSERTION with needs_review true", fp.License)
	}
	if fp.License != nil && fp.License.URL != "" {
		t.Errorf("an unclassified license has no URL, got %q", fp.License.URL)
	}
	if !strings.Contains(h.log.String(), "NOASSERTION") || !strings.Contains(h.log.String(), "LICENSE") {
		t.Errorf("the note should tell the Author to read the LICENSE file by hand; log:\n%s", h.log)
	}
	if len(res.Notes) == 0 {
		t.Error("Result.Notes is empty")
	}
}

func TestFetch_ClassifiedLicenseDoesNotNeedReview(t *testing.T) {
	h := newHarness(t)
	h.gh.serveRepo("apache-iggy")
	md := h.markdownPath("distributed-systems", "apache-iggy")

	if _, err := h.fetcher.Fetch(context.Background(), "apache/iggy", md); err != nil {
		t.Fatal(err)
	}

	fp := parseFactPack(t, filepath.Join(filepath.Dir(md), "apache-iggy.json"))
	if fp.License == nil || fp.License.NeedsReview || fp.License.SPDXID != "Apache-2.0" {
		t.Errorf("license = %+v, want Apache-2.0 with needs_review false", fp.License)
	}
	if strings.Contains(h.log.String(), "NOASSERTION") {
		t.Errorf("no NOASSERTION note expected; log:\n%s", h.log)
	}
}

func TestFetch_READMETextLandsInTheCacheDirectoryNotTheContentTree(t *testing.T) {
	h := newHarness(t)
	h.gh.serveRepo("apache-iggy")
	md := h.markdownPath("distributed-systems", "apache-iggy")

	res, err := h.fetcher.Fetch(context.Background(), "apache/iggy", md)
	if err != nil {
		t.Fatal(err)
	}

	want := filepath.Join(h.cacheDir, "readmes", "apache", "iggy", "README.md")
	if res.READMEFile != want {
		t.Errorf("README written to %s, want %s", res.READMEFile, want)
	}
	text := readFile(t, want)
	if !strings.HasPrefix(string(text), "# Iggy\n") || !strings.Contains(string(text), "docker run -p 8090:8090 apache/iggy:latest") {
		t.Errorf("README text was not decoded:\n%s", text)
	}
	for _, f := range treeFiles(t, h.contentDir) {
		if f != filepath.Join("explainers", "distributed-systems", "apache-iggy.json") {
			t.Errorf("unexpected file in the content tree: %s", f)
		}
	}
	fp := parseFactPack(t, filepath.Join(filepath.Dir(md), "apache-iggy.json"))
	if fp.README.URL != "https://github.com/apache/iggy/blob/main/README.md" || fp.README.SHA != "7c6b2a0f9e8d4c3b1a0f9e8d7c6b5a4f3e2d1c0b" {
		t.Errorf("readme ref = %+v", fp.README)
	}
	if strings.Contains(string(readFile(t, filepath.Join(filepath.Dir(md), "apache-iggy.json"))), "Quickstart") {
		t.Error("the Fact Pack carries README text; it must carry only the URL and SHA")
	}
}

func TestFetch_RepositoryWithoutREADMEOrDotGithub(t *testing.T) {
	h := newHarness(t)
	h.gh.serveRepo("nats-io-nats-server")
	md := h.markdownPath("distributed-systems", "nats-server")

	res, err := h.fetcher.Fetch(context.Background(), "nats-io/nats-server", md)
	if err != nil {
		t.Fatal(err)
	}

	if res.READMEFile != "" {
		t.Errorf("no README expected, got %s", res.READMEFile)
	}
	fp := parseFactPack(t, filepath.Join(filepath.Dir(md), "nats-server.json"))
	if fp.README != (content.READMERef{}) || fp.Files.README != "" {
		t.Errorf("readme = %+v files.readme = %q, want empty", fp.README, fp.Files.README)
	}
	if fp.Files.License != "LICENSE" || len(fp.Files.CI) != 1 || fp.Files.CI[0] != ".travis.yml" || len(fp.Files.Manifests) != 1 || fp.Files.Manifests[0] != "go.mod" {
		t.Errorf("files = %+v", fp.Files)
	}
	if fp.LatestRelease != nil || len(fp.Releases) != 0 {
		t.Errorf("no releases expected, got latest %+v releases %+v", fp.LatestRelease, fp.Releases)
	}
	if fp.Contributors.Total != 1 {
		t.Errorf("contributors.total = %d, want 1 (no Link header: the one page is the count)", fp.Contributors.Total)
	}
	// stargazers/count is not served for this repository: the repository's
	// stargazers_count stands in, with a note, and the endpoint is not a source.
	if fp.Stargazers.Count != 17820 {
		t.Errorf("stargazers.count = %d, want the repository's 17820", fp.Stargazers.Count)
	}
	for _, s := range fp.Sources {
		if strings.HasSuffix(s, "/stargazers/count") || strings.HasSuffix(s, "/readme") || strings.Contains(s, "/contents/.github") {
			t.Errorf("source %s answered 404 and must not be listed", s)
		}
	}
	if !strings.Contains(h.log.String(), "stargazers_count") {
		t.Errorf("expected a note about the stargazers fallback; log:\n%s", h.log)
	}
}

func TestFetch_SendsTheTokenAndAnIdentifyingUserAgentOnEveryRequest(t *testing.T) {
	h := newHarness(t)
	h.gh.serveRepo("apache-iggy")
	md := h.markdownPath("distributed-systems", "apache-iggy")

	if _, err := h.fetcher.Fetch(context.Background(), "apache/iggy", md); err != nil {
		t.Fatal(err)
	}

	reqs := h.gh.requestsTo("/")
	if len(reqs) == 0 {
		t.Fatal("no requests recorded")
	}
	for _, r := range reqs {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("%s: Authorization = %q", r.Path, got)
		}
		if ua := r.Header.Get("User-Agent"); !strings.Contains(ua, "hastcode") {
			t.Errorf("%s: User-Agent %q does not identify the site", r.Path, ua)
		}
		if got := r.Header.Get("Accept"); got != "application/vnd.github+json" {
			t.Errorf("%s: Accept = %q", r.Path, got)
		}
	}
}

func TestFetch_MakesAtMostTenCallsPerRepository(t *testing.T) {
	h := newHarness(t)
	h.gh.serveRepo("apache-iggy")
	md := h.markdownPath("distributed-systems", "apache-iggy")

	if _, err := h.fetcher.Fetch(context.Background(), "apache/iggy", md); err != nil {
		t.Fatal(err)
	}

	if n := len(h.gh.requestsTo("/")); n > 10 {
		t.Errorf("%d requests for one repository, the budget is 10", n)
	}
}

func TestFetch_RefusesWithoutAToken(t *testing.T) {
	_, err := github.NewClient(func(string) string { return "" })
	if err == nil {
		t.Fatal("expected a refusal")
	}
	if !strings.Contains(err.Error(), "GITHUB_TOKEN") {
		t.Errorf("the refusal must name the variable: %v", err)
	}
}

func TestFetch_RefusesToWriteBesideAProfileOfAnotherRepository(t *testing.T) {
	h := newHarness(t)
	h.gh.serveRepo("apache-iggy")
	md := h.writeProfile("distributed-systems", "nats-server", "nats-io/nats-server")

	_, err := h.fetcher.Fetch(context.Background(), "apache/iggy", md)
	if err == nil || !strings.Contains(err.Error(), "nats-io/nats-server") {
		t.Fatalf("expected a refusal naming the Profile's repository, got %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(filepath.Dir(md), "nats-server.json")); statErr == nil {
		t.Error("the Fact Pack was written anyway")
	}
}

func parseFactPack(t *testing.T, file string) *content.FactPack {
	t.Helper()
	fp, err := content.ParseFactPack(readFile(t, file))
	if err != nil {
		t.Fatalf("%s: %v", file, err)
	}
	return fp
}

// treeFiles lists the files under root, root-relative.
func treeFiles(t *testing.T, root string) []string {
	t.Helper()
	var files []string
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(root, p)
		files = append(files, rel)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

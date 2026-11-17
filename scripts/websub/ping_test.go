// Package websub tests the publish workflow's side of WebSub (README
// "Feeds, sitemaps, robots and IndexNow"): only the Production job pings,
// after the deploy, right after "Ping IndexNow", behind WEBSUB_ENABLED,
// and a failed ping warns without failing the job. ping.sh is Bash, so
// these tests run it with bash against a local stand-in hub: nothing
// leaves the machine. `make test` runs them with the Go tests, in GitHub
// Actions too (the runners have bash, curl and jq).
package websub

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"gopkg.in/yaml.v3"
)

type workflow struct {
	Jobs map[string]struct {
		Steps []workflowStep `yaml:"steps"`
	} `yaml:"jobs"`
}

type workflowStep struct {
	Name string `yaml:"name"`
	ID   string `yaml:"id"`
	Uses string `yaml:"uses"`
	If   string `yaml:"if"`
	Run  string `yaml:"run"`
}

func readWorkflow(t *testing.T, name string) (workflow, string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", name))
	if err != nil {
		t.Fatal(err)
	}
	var w workflow
	if err := yaml.Unmarshal(data, &w); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return w, string(data)
}

const stepName = "Ping WebSub hub"

func TestWorkflow_OnlyTheProductionJobPingsAfterTheDeployNextToIndexNow(t *testing.T) {
	w, raw := readWorkflow(t, "publish.yml")

	var names []string
	var ping workflowStep
	deploy, indexNow, hub, firstSave := -1, -1, -1, -1
	for i, s := range w.Jobs["cloudflare"].Steps {
		names = append(names, s.Name+s.Uses)
		switch {
		case s.ID == "deploy":
			deploy = i
		case s.Name == "Ping IndexNow":
			indexNow = i
		case s.Name == stepName:
			hub, ping = i, s
		case strings.HasPrefix(s.Uses, "actions/cache/save") && firstSave < 0:
			firstSave = i
		}
	}
	if hub < 0 || !(deploy < indexNow && indexNow+1 == hub && hub < firstSave) {
		t.Fatalf("cloudflare steps %v: want %q right after \"Ping IndexNow\", after the deploy and before the cache save", names, stepName)
	}
	if ping.If != "vars.WEBSUB_ENABLED == 'true'" {
		t.Errorf("%s runs if %q, want vars.WEBSUB_ENABLED == 'true'", stepName, ping.If)
	}
	if strings.TrimSpace(ping.Run) != `scripts/websub/ping.sh "$RUNNER_TEMP/websub.json"` {
		t.Errorf("%s runs %q, want scripts/websub/ping.sh on the set-aside artefact", stepName, ping.Run)
	}

	// The artefact is set aside before the adapter and upload, like
	// indexnow.json; QA and Previews only ever delete it.
	if !strings.Contains(raw, `mv .vercel/output/websub.json "$RUNNER_TEMP/websub.json"`) {
		t.Error("publish.yml does not set websub.json aside before the Cloudflare deploy")
	}
	requireNoPing(t, "publish.yml vercel", w.Jobs["vercel"].Steps)
	preview, _ := readWorkflow(t, "preview.yml")
	for name, j := range preview.Jobs {
		requireNoPing(t, "preview.yml "+name, j.Steps)
	}

	// The header comment lists the variable beside INDEXNOW_ENABLED.
	header := raw[:strings.Index(raw, "\nname:")]
	if !strings.Contains(header, "INDEXNOW_ENABLED") || !strings.Contains(header, "WEBSUB_ENABLED") {
		t.Error("publish.yml's header comment does not list WEBSUB_ENABLED next to INDEXNOW_ENABLED")
	}
}

// requireNoPing fails when a QA or Preview job could notify a hub. A job
// that deletes indexnow.json before its upload deletes websub.json too.
func requireNoPing(t *testing.T, job string, steps []workflowStep) {
	t.Helper()
	uploads, removed := false, false
	for _, s := range steps {
		text := s.Name + s.If + s.Run
		for _, ping := range []string{"WEBSUB_ENABLED", "websub/ping.sh", "hub.mode", "pubsubhubbub"} {
			if strings.Contains(text, ping) {
				t.Errorf("%s step %q can ping a WebSub hub (%s)", job, s.Name, ping)
			}
		}
		switch strings.TrimSpace(s.Run) {
		case "rm -f .vercel/output/indexnow.json":
			uploads = true
		case "rm -f .vercel/output/websub.json":
			removed = true
		}
	}
	if uploads && !removed {
		t.Errorf("%s does not delete .vercel/output/websub.json before its upload", job)
	}
}

// hubRequest is what the stand-in hub received.
type hubRequest struct {
	method, contentType string
	form                url.Values
}

// standInHub answers every request with status and records it.
func standInHub(t *testing.T, status int) (string, func() []hubRequest) {
	var mu sync.Mutex
	var got []hubRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		form, _ := url.ParseQuery(string(body))
		mu.Lock()
		got = append(got, hubRequest{method: r.Method, contentType: r.Header.Get("Content-Type"), form: form})
		mu.Unlock()
		w.WriteHeader(status)
	}))
	t.Cleanup(srv.Close)
	return srv.URL + "/", func() []hubRequest {
		mu.Lock()
		defer mu.Unlock()
		return append([]hubRequest(nil), got...)
	}
}

// runPing runs ping.sh the way the step does (GitHub's bash -e) on
// artefact, written to RUNNER_TEMP (none when ""), and returns its exit
// error, its output and what it appended to the run summary.
func runPing(t *testing.T, artefact string) (string, string, error) {
	t.Helper()
	for _, tool := range []string{"bash", "curl", "jq"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Fatalf("%s is not on PATH; ping.sh needs it (GitHub's runners have it)", tool)
		}
	}
	script, err := filepath.Abs("ping.sh")
	if err != nil {
		t.Fatal(err)
	}
	tmp := t.TempDir()
	notice := filepath.Join(tmp, "websub.json")
	if artefact != "" {
		if err := os.WriteFile(notice, []byte(artefact), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	summary := filepath.Join(tmp, "summary.md")
	cmd := exec.Command("bash", "--noprofile", "--norc", "-e", script, notice)
	cmd.Dir = tmp
	cmd.Env = append(os.Environ(), "RUNNER_TEMP="+tmp, "GITHUB_STEP_SUMMARY="+summary)
	out, err := cmd.CombinedOutput()
	s, _ := os.ReadFile(summary)
	return string(out), string(s), err
}

// notice is a websub.json artefact as the build writes it.
func notice(hub string, topics ...string) string {
	quoted := make([]string, len(topics))
	for i, u := range topics {
		quoted[i] = `"` + u + `"`
	}
	return `{"hub": "` + hub + `", "topics": [` + strings.Join(quoted, ", ") + "]}\n"
}

func TestPing_PostsEveryChangedFeedInOneFormRequest(t *testing.T) {
	hub, requests := standInHub(t, http.StatusNoContent)
	topics := []string{"https://hastcode.com/briefs/feed.xml", "https://hastcode.com/feed.xml"}

	out, summary, err := runPing(t, notice(hub, topics...))
	if err != nil {
		t.Fatalf("ping failed: %v\n%s", err, out)
	}
	got := requests()
	if len(got) != 1 {
		t.Fatalf("hub got %d requests, want 1:\n%s", len(got), out)
	}
	r := got[0]
	if r.method != http.MethodPost || r.contentType != "application/x-www-form-urlencoded" {
		t.Errorf("request = %s %s, want a form POST", r.method, r.contentType)
	}
	if r.form.Get("hub.mode") != "publish" || !reflect.DeepEqual(r.form["hub.url"], topics) {
		t.Errorf("form = %v, want hub.mode=publish and hub.url=%v", r.form, topics)
	}
	if !strings.Contains(out, "answered 204 for 2 feed(s)") || !strings.Contains(summary, "WebSub: "+hub+" answered 204 for 2 feed(s).") {
		t.Errorf("success is not reported; output:\n%s\nsummary:\n%s", out, summary)
	}
}

func TestPing_AFailedPingWarnsWithoutFailingTheJob(t *testing.T) {
	refusing, _ := standInHub(t, http.StatusServiceUnavailable)
	gone := httptest.NewServer(http.NotFoundHandler())
	gone.Close()
	for name, hub := range map[string]string{
		"hub answers 503": refusing,
		"hub unreachable": gone.URL + "/",
	} {
		t.Run(name, func(t *testing.T) {
			out, summary, err := runPing(t, notice(hub, "https://hastcode.com/feed.xml"))
			if err != nil {
				t.Fatalf("a failed ping failed the step (%v); it must only warn:\n%s", err, out)
			}
			if !strings.Contains(out, "::warning title=WebSub ping failed::") || !strings.Contains(summary, "WebSub ping failed: "+hub) {
				t.Errorf("no warning; output:\n%s\nsummary:\n%s", out, summary)
			}
		})
	}
}

func TestPing_NoArtefactMeansNoPing(t *testing.T) {
	out, summary, err := runPing(t, "")
	if err != nil {
		t.Fatalf("ping failed without an artefact: %v\n%s", err, out)
	}
	if !strings.Contains(out, "no feed changed") || summary != "" {
		t.Errorf("output does not say why nothing was pinged, or the summary is not empty:\n%s\nsummary:\n%s", out, summary)
	}
}

func TestSetup_TheWizardSetsTheSwitchOff(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "setup-deploy.sh"))
	if err != nil {
		t.Fatal(err)
	}
	script := string(data)
	if !strings.Contains(script, `set_var WEBSUB_ENABLED "false"`) {
		t.Error(`scripts/setup-deploy.sh does not set WEBSUB_ENABLED "false"`)
	}
	if strings.Contains(script, `set_var WEBSUB_ENABLED "true"`) {
		t.Error("scripts/setup-deploy.sh turns WEBSUB_ENABLED on; only the Author does, by hand")
	}
}

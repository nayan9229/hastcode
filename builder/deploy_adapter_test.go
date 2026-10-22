package builder_test

import (
	"os"
	"slices"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

// The deploy-adapter rule (ticket 05, ADR-0004): one build output, a file
// set per Target. Neither Target may serve the other's configuration. The
// adapter runs as deploys run it: the cmd/adapt command (`make adapt`).

// The files only one Target reads (README "One build, two Targets").
var (
	cloudflareOnly = []string{"static/_headers", "static/_redirects", "wrangler.json"}
	vercelOnly     = []string{"config.json"}
)

// without returns files minus drop, keeping the order.
func without(files, drop []string) []string {
	return slices.DeleteFunc(slices.Clone(files), func(f string) bool { return slices.Contains(drop, f) })
}

func TestDeployAdapter_EachTargetKeepsEverythingButTheOtherTargetsManifest(t *testing.T) {
	for target, foreign := range map[string][]string{"cloudflare": vercelOnly, "vercel": cloudflareOnly} {
		t.Run(target, func(t *testing.T) {
			out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree())
			built := buildtest.OutputFiles(t, out)
			for _, f := range append(slices.Clone(cloudflareOnly), vercelOnly...) {
				if !slices.Contains(built, f) {
					t.Fatalf("the build wrote no %s; files:\n%q", f, built)
				}
			}

			buildtest.Adapt(t, out, target)

			want := without(built, foreign)
			if got := buildtest.OutputFiles(t, out); !slices.Equal(got, want) {
				t.Errorf("after adapt -target %s the output holds\n%q\nwant everything but %q:\n%q", target, got, foreign, want)
			}
			// Adapting twice is harmless (a re-run deploy step).
			buildtest.Adapt(t, out, target)
			if got := buildtest.OutputFiles(t, out); !slices.Equal(got, want) {
				t.Errorf("a second adapt changed the output to\n%q", got)
			}
		})
	}
}

// publish.sh's short codes and the Publish tags name the same Targets.
func TestDeployAdapter_TargetMayBeNamedByShortCodeOrPublishTag(t *testing.T) {
	for spelling, foreign := range map[string][]string{
		"cf":                    vercelOnly,
		"VL":                    cloudflareOnly,
		"PUBCF-20261022-101500": vercelOnly,
		"PUBVL-20261022-101500": cloudflareOnly,
	} {
		t.Run(spelling, func(t *testing.T) {
			out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree())
			built := buildtest.OutputFiles(t, out)

			buildtest.Adapt(t, out, spelling)

			if got, want := buildtest.OutputFiles(t, out), without(built, foreign); !slices.Equal(got, want) {
				t.Errorf("after adapt -target %s the output holds\n%q\nwant\n%q", spelling, got, want)
			}
		})
	}
}

func TestDeployAdapter_UnknownTargetExits2AndChangesNothing(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree())
	built := buildtest.OutputFiles(t, out)
	robots, err := os.ReadFile(buildtest.StaticPath(out, "/robots.txt"))
	if err != nil {
		t.Fatal(err)
	}

	for _, target := range []string{"", "pages", "aws", "PUBXX-20261022-101500", "pubcf-20261022-101500"} {
		code, output := buildtest.RunAdapt(t, out, target)
		if code != 2 {
			t.Errorf("adapt -target %q exited %d, want 2:\n%s", target, code, output)
		}
	}

	if got := buildtest.OutputFiles(t, out); !slices.Equal(got, built) {
		t.Errorf("a refused adapt changed the output to\n%q\nwant\n%q", got, built)
	}
	if got, _ := os.ReadFile(buildtest.StaticPath(out, "/robots.txt")); string(got) != string(robots) {
		t.Errorf("a refused adapt rewrote robots.txt:\n%s", got)
	}
}

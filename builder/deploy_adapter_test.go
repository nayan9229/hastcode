package builder_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"hastcode.com/site/internal/buildtest"
	"hastcode.com/site/internal/deploy"
)

// The deploy-adapter rule (ticket 05): one build output, a file set per
// Target. Neither Target may serve the other's configuration.

var (
	cloudflareOnly = []string{"wrangler.json", "static/_headers", "static/_redirects"}
	vercelOnly     = []string{"config.json"}
)

func TestDeployAdapter_EachTargetsFileSetExcludesTheOthersManifest(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree())

	vercel, err := deploy.Files(out, deploy.Vercel)
	if err != nil {
		t.Fatal(err)
	}
	cloudflare, err := deploy.Files(out, deploy.Cloudflare)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range cloudflareOnly {
		if slices.Contains(vercel, f) {
			t.Errorf("Vercel file set contains Cloudflare's %s", f)
		}
		if !slices.Contains(cloudflare, f) {
			t.Errorf("Cloudflare file set lacks its own %s", f)
		}
	}
	for _, f := range vercelOnly {
		if slices.Contains(cloudflare, f) {
			t.Errorf("Cloudflare file set contains Vercel's %s", f)
		}
		if !slices.Contains(vercel, f) {
			t.Errorf("Vercel file set lacks its own %s", f)
		}
	}
	// The site itself is shared: both Targets serve the same pages.
	for _, f := range []string{"static/index.html", "static/404.html", "static/robots.txt"} {
		if !slices.Contains(vercel, f) || !slices.Contains(cloudflare, f) {
			t.Errorf("%s missing from a Target's file set (Vercel %t, Cloudflare %t)", f, slices.Contains(vercel, f), slices.Contains(cloudflare, f))
		}
	}
}

func TestDeployAdapter_ApplyLeavesExactlyTheTargetsFileSetOnDisk(t *testing.T) {
	for _, target := range []deploy.Target{deploy.Vercel, deploy.Cloudflare} {
		t.Run(string(target), func(t *testing.T) {
			out := buildtest.MustBuild(t, buildtest.Config(), buildtest.NewTree())
			want, err := deploy.Files(out, target)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := deploy.Apply(out, target); err != nil {
				t.Fatal(err)
			}
			var got []string
			err = filepath.WalkDir(out, func(p string, d os.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return err
				}
				rel, _ := filepath.Rel(out, p)
				got = append(got, filepath.ToSlash(rel))
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(got, want) {
				t.Errorf("after Apply(%s) the output holds\n%q\nwant its file set\n%q", target, got, want)
			}
			// Applying twice is harmless (a re-run deploy step).
			if _, err := deploy.Apply(out, target); err != nil {
				t.Errorf("second Apply: %v", err)
			}
		})
	}
}

func TestDeployAdapter_ParseTargetAcceptsOnlyTheTwoTargets(t *testing.T) {
	for in, want := range map[string]deploy.Target{"cloudflare": deploy.Cloudflare, "vercel": deploy.Vercel} {
		if got, err := deploy.ParseTarget(in); err != nil || got != want {
			t.Errorf("ParseTarget(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	for _, in := range []string{"", "pages", "Cloudflare", "cf"} {
		if _, err := deploy.ParseTarget(in); err == nil {
			t.Errorf("ParseTarget(%q) succeeded; want an error", in)
		}
	}
}

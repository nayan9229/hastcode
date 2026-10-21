package builder_test

import (
	"testing"
	"time"

	"hastcode.com/site/internal/buildtest"
)

func TestBuildOutput_ConfigJSONIsVersion3(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())

	cfg := buildtest.VercelConfig(t, out)
	if got := cfg["version"]; got != float64(3) {
		t.Errorf("config.json version = %v, want 3", got)
	}
}

func TestBuildOutput_FixtureTreeBuildsWellUnderASecond(t *testing.T) {
	// Hero encoding is paid once per source (image cache keyed on source
	// bytes); warm it so this measures the per-build cost every test pays.
	buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())

	start := time.Now()
	buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	elapsed := time.Since(start)

	t.Logf("fixture build took %s", elapsed)
	if elapsed > 500*time.Millisecond {
		t.Errorf("fixture build took %s, want well under a second (< 500ms)", elapsed)
	}
}

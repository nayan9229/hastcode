package builder_test

import (
	"os"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

const googleAdsTxt = "google.com, pub-8711471207758471, DIRECT, f08c47fec0942fa0"

// /ads.txt (IAB Authorized Digital Sellers) lists the ad systems allowed to
// sell the site's inventory, one record per line, verbatim from site.yaml.
func TestAdsTxt_WritesConfiguredRecordsAtTheRoot(t *testing.T) {
	cfg := buildtest.Config()
	cfg.AdsTxt = []string{googleAdsTxt, "example-exchange.com, 12345, RESELLER"}
	out := buildtest.MustBuild(t, cfg, buildtest.ValidTree())

	data, err := os.ReadFile(buildtest.StaticPath(out, "/ads.txt"))
	if err != nil {
		t.Fatalf("/ads.txt not written: %v", err)
	}
	want := googleAdsTxt + "\nexample-exchange.com, 12345, RESELLER\n"
	if string(data) != want {
		t.Errorf("/ads.txt =\n%s\nwant\n%s", data, want)
	}
}

func TestAdsTxt_OmittedWhenNoRecordsConfigured(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	if buildtest.Exists(out, "/ads.txt") {
		t.Error("/ads.txt written with no records configured")
	}
}

// A record needs the ad system's domain, the publisher's account id and
// the relationship (DIRECT or RESELLER); the certification authority id is
// optional. A malformed line would silently de-authorise a seller.
func TestAdsTxt_MalformedRecordFailsTheBuild(t *testing.T) {
	for _, bad := range []string{
		"google.com, pub-8711471207758471",          // no relationship
		"google.com, pub-8711471207758471, PARTNER", // unknown relationship
		", pub-8711471207758471, DIRECT",            // no domain
		"google.com, , DIRECT",                      // no account id
	} {
		cfg := buildtest.Config()
		cfg.AdsTxt = []string{bad}
		_, res := buildtest.Build(t, cfg, buildtest.ValidTree())
		buildtest.RequireErrorRule(t, res, "", "ads-txt", bad)
	}
}

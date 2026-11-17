package crawlcheck_test

import (
	"strings"
	"testing"
)

func TestEverythingPassingExitsZero(t *testing.T) {
	s := newStandIn(t)
	r, out := run(t, s.checker())
	if r.ExitCode() != 0 {
		t.Fatalf("exit code %d, want 0; output:\n%s", r.ExitCode(), out)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	last := lines[len(lines)-1]
	if !strings.Contains(last, "passed") {
		t.Errorf("last line %q is not the summary", last)
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "FAIL") {
			t.Errorf("unexpected failure: %s", line)
		}
	}
	for _, crawler := range []string{"Googlebot", "Bingbot", "OAI-SearchBot", "PerplexityBot", "Claude-SearchBot", "Applebot", "DuckDuckBot", "YandexBot", "Amzn-SearchBot"} {
		if !strings.Contains(out, "PASS "+crawler+" gets 200") {
			t.Errorf("no access line for %s; output:\n%s", crawler, out)
		}
	}
}

func TestCrawlerGetting403NamesItAndRunbookStep1(t *testing.T) {
	s := newStandIn(t)
	s.block("PerplexityBot", 403)
	r, out := run(t, s.checker())
	requireFailed(t, r, out, "Cloudflare is blocking PerplexityBot", "runbook step 1")
	if r.Failed() != 1 {
		t.Errorf("%d checks failed, want only PerplexityBot's; output:\n%s", r.Failed(), out)
	}
}

func TestOtherNon200IsReportedWithItsCode(t *testing.T) {
	s := newStandIn(t)
	s.block("DuckDuckBot", 503)
	r, out := run(t, s.checker())
	requireFailed(t, r, out, "DuckDuckBot", "503")
}

func TestHTTPAnswering200NamesAlwaysUseHTTPS(t *testing.T) {
	s := newStandIn(t)
	s.setHTTP(200, "")
	r, out := run(t, s.checker())
	requireFailed(t, r, out, s.http.URL, "200", "Always Use HTTPS", "runbook step 2")
}

func TestHTTPRedirectingElsewhereFails(t *testing.T) {
	s := newStandIn(t)
	s.setHTTP(301, "https://elsewhere.invalid/")
	r, out := run(t, s.checker())
	requireFailed(t, r, out, "https://elsewhere.invalid/", "Always Use HTTPS")
}

func TestMissingSitemapLineFails(t *testing.T) {
	s := newStandIn(t)
	robots := strings.Replace(cloudflareBlock+ourRobots(s.site.URL), "Sitemap: "+s.site.URL+"/news-sitemap.xml\n", "", 1)
	s.setRobots(robots)
	r, out := run(t, s.checker())
	requireFailed(t, r, out, "Sitemap: "+s.site.URL+"/news-sitemap.xml")
}

func TestMissingAllowGroupFails(t *testing.T) {
	s := newStandIn(t)
	s.setRobots(strings.Replace(cloudflareBlock+ourRobots(s.site.URL), "User-agent: Bingbot\nAllow: /\n\n", "", 1))
	r, out := run(t, s.checker())
	requireFailed(t, r, out, "Bingbot allow group")
}

func TestPrependedGroupDisallowingACheckedCrawlerFails(t *testing.T) {
	s := newStandIn(t)
	s.setRobots("User-agent: Claude-SearchBot\nDisallow: /\n\n" + cloudflareBlock + ourRobots(s.site.URL))
	r, out := run(t, s.checker())
	requireFailed(t, r, out, "robots.txt stops Claude-SearchBot", "Disallow: /", "line 2")
}

func TestPrependedStarDisallowBindsEveryUnnamedCrawler(t *testing.T) {
	// The * groups combine (RFC 9309): a Disallow: / in Cloudflare's *
	// group, against our Allow: / in ours, ties on length and Allow wins.
	// A longer Disallow does not tie: Disallow: /* beats Allow: /.
	s := newStandIn(t)
	s.setRobots(strings.Replace(cloudflareBlock, "Allow: /\n", "Disallow: /*\n", 1) + ourRobots(s.site.URL))
	r, out := run(t, s.checker())
	requireFailed(t, r, out, "robots.txt stops DuckDuckBot", "the * groups")
	if strings.Contains(out, "FAIL robots.txt stops Googlebot") {
		t.Errorf("Googlebot has its own group, so the * groups do not bind it; output:\n%s", out)
	}
}

func TestStarGroupsCombine(t *testing.T) {
	// Cloudflare's * group says Disallow: /, ours says Allow: /. Read as
	// one group the two rules tie and Allow wins, so nobody is shut out.
	s := newStandIn(t)
	s.setRobots(strings.Replace(cloudflareBlock, "Allow: /\n", "Disallow: /\n", 1) + ourRobots(s.site.URL))
	if r, out := run(t, s.checker()); r.ExitCode() != 0 {
		t.Fatalf("exit code %d, want 0; output:\n%s", r.ExitCode(), out)
	}
}

func TestApplebotFollowsGooglebotsGroup(t *testing.T) {
	s := newStandIn(t)
	robots := strings.Replace(ourRobots(s.site.URL), "User-agent: Googlebot\nAllow: /\n", "User-agent: Googlebot\nAllow: /\nDisallow: /$\n", 1)
	s.setRobots(cloudflareBlock + robots)
	r, out := run(t, s.checker())
	requireFailed(t, r, out, "robots.txt stops Applebot", "the Googlebot group")
	if strings.Contains(out, "Applebot-Extended") {
		t.Errorf("Applebot-Extended's group was read as Applebot's; output:\n%s", out)
	}
}

func TestApplebotNamedInItsOwnGroupIgnoresGooglebots(t *testing.T) {
	s := newStandIn(t)
	robots := strings.Replace(ourRobots(s.site.URL), "User-agent: Googlebot\nAllow: /\n", "User-agent: Googlebot\nDisallow: /$\n\nUser-agent: Applebot\nAllow: /\n", 1)
	s.setRobots(cloudflareBlock + robots)
	_, out := run(t, s.checker())
	if !strings.Contains(out, "PASS robots.txt lets Applebot crawl / (it follows the Applebot group)") {
		t.Errorf("Applebot did not follow its own group; output:\n%s", out)
	}
}

func TestHTTPRedirectBy308Passes(t *testing.T) {
	s := newStandIn(t)
	s.setHTTP(308, s.site.URL+"/")
	if r, out := run(t, s.checker()); r.ExitCode() != 0 {
		t.Fatalf("exit code %d, want 0; output:\n%s", r.ExitCode(), out)
	}
}

package builder_test

import (
	"os"
	"strings"
	"testing"

	"hastcode.com/site/internal/buildtest"
)

// robots.txt allows every crawler, names the search and answer-engine
// crawlers the site wants citations from, and points at both sitemaps.
// Training crawlers (GPTBot, Google-Extended) stay allowed: opting out is
// a separate, later decision (spec "Out of Scope").
func TestRobotsTxt_AllowsEveryCrawlerAndNamesBothSitemaps(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), buildtest.ValidTree())
	data, err := os.ReadFile(buildtest.StaticPath(out, "/robots.txt"))
	if err != nil {
		t.Fatal(err)
	}

	want := `User-agent: *
Allow: /

User-agent: Googlebot
Allow: /

User-agent: Bingbot
Allow: /

User-agent: OAI-SearchBot
Allow: /

User-agent: PerplexityBot
Allow: /

Sitemap: https://hastcode.test/sitemap.xml
Sitemap: https://hastcode.test/news-sitemap.xml
`
	if got := string(data); got != want {
		t.Errorf("robots.txt =\n%s\nwant\n%s", got, want)
	}
	for _, bot := range []string{"GPTBot", "Google-Extended"} {
		if strings.Contains(string(data), bot) {
			t.Errorf("robots.txt names %s; training crawlers are left to the * group", bot)
		}
	}
	if strings.Contains(string(data), "Disallow") {
		t.Error("robots.txt disallows something")
	}
}

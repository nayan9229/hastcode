package builder

import (
	"fmt"
	"strings"
)

func init() {
	registerStep(phaseWrite, 40, "write robots.txt", writeRobots)
}

// namedCrawlers get their own allow group in robots.txt (spec "Structured
// data and head"; docs/research/aeo-citation-signals.md): Google Search and
// AI Overviews, Bing and Copilot, ChatGPT search, Perplexity. A crawler
// that matches a named group ignores the * group, so each group must carry
// every rule; today that is only "Allow: /". GPTBot and Google-Extended
// (model training) are deliberately not named: they fall under *.
var namedCrawlers = []string{"Googlebot", "Bingbot", "OAI-SearchBot", "PerplexityBot"}

// writeRobots writes /robots.txt: allow all, the named crawler groups,
// and both sitemaps.
func writeRobots(s *site) error {
	var b strings.Builder
	for _, agent := range append([]string{"*"}, namedCrawlers...) {
		fmt.Fprintf(&b, "User-agent: %s\nAllow: /\n\n", agent)
	}
	fmt.Fprintf(&b, "Sitemap: %s\nSitemap: %s\n", s.absURL(sitemapURL), s.absURL(newsSitemapURL))
	return s.out.WriteStatic("/robots.txt", []byte(b.String()))
}

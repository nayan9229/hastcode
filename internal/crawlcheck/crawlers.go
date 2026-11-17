package crawlcheck

// Crawler is one search or answer crawler the check sends requests as.
type Crawler struct {
	// Name is how the output names it.
	Name string
	// UserAgent is the full user-agent string its owner publishes, sent
	// verbatim (placeholders such as Chrome/W.X.Y.Z included).
	UserAgent string
	// Tokens are the robots.txt product tokens it obeys, most specific
	// first: the first token some group names decides its rules, and with
	// none named it follows the combined * groups (RFC 9309).
	Tokens []string
	// Source is the owner's page the string and tokens come from.
	Source string
}

// Crawlers are the search and answer crawlers we want reaching the site
// (docs/research/indexing-submission-channels.md §0, §4, §6). Training
// crawlers (GPTBot, ClaudeBot, CCBot, Amazonbot) are left out on purpose:
// blocking them is the Author's choice. Strings checked against each
// Source on 2026-11-17 (Bingbot and Claude-SearchBot: see their notes); a
// version number in one may change, so update it here when its owner's
// page does.
var Crawlers = []Crawler{
	{
		Name:      "Googlebot",
		UserAgent: "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Googlebot/2.1; +http://www.google.com/bot.html) Chrome/W.X.Y.Z Safari/537.36",
		Tokens:    []string{"Googlebot"},
		Source:    "https://developers.google.com/search/docs/crawling-indexing/google-common-crawlers",
	},
	{
		// Bing's help page is client-rendered (the note's Open question
		// 4) and could not be read on 2026-11-17; this is the desktop
		// string Bing has published there. Check it in a browser.
		Name:      "Bingbot",
		UserAgent: "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm) Chrome/W.X.Y.Z Safari/537.36",
		Tokens:    []string{"Bingbot"},
		Source:    "https://www.bing.com/webmasters/help/which-crawlers-does-bing-use-8c184ec0",
	},
	{
		Name:      "OAI-SearchBot",
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36; compatible; OAI-SearchBot/1.4; +https://openai.com/searchbot",
		Tokens:    []string{"OAI-SearchBot"},
		Source:    "https://developers.openai.com/api/docs/bots",
	},
	{
		Name:      "PerplexityBot",
		UserAgent: "Mozilla/5.0 (compatible; PerplexityBot/1.0; +https://perplexity.ai/perplexitybot)",
		Tokens:    []string{"PerplexityBot"},
		Source:    "https://docs.perplexity.ai/guides/bots",
	},
	{
		// Anthropic's page names the token but publishes no full string;
		// this is the shape of its other crawlers' strings.
		Name:      "Claude-SearchBot",
		UserAgent: "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Claude-SearchBot/1.0; +Claude-SearchBot@anthropic.com)",
		Tokens:    []string{"Claude-SearchBot"},
		Source:    "https://support.claude.com/en/articles/8896518-does-anthropic-crawl-data-from-the-web-and-how-can-site-owners-block-the-crawler",
	},
	{
		// Not named in robots.txt, Applebot follows Googlebot's group
		// (Apple's page).
		Name:      "Applebot",
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Safari/605.1.15 (Applebot/0.1; +http://www.apple.com/go/applebot)",
		Tokens:    []string{"Applebot", "Googlebot"},
		Source:    "https://support.apple.com/en-us/119829",
	},
	{
		Name:      "DuckDuckBot",
		UserAgent: "DuckDuckBot/1.1; (+http://duckduckgo.com/duckduckbot.html)",
		Tokens:    []string{"DuckDuckBot"},
		Source:    "https://duckduckgo.com/duckduckgo-help-pages/results/duckduckbot",
	},
	{
		// "User-agent: Yandex" covers every Yandex robot.
		Name:      "YandexBot",
		UserAgent: "Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)",
		Tokens:    []string{"YandexBot", "Yandex"},
		Source:    "https://yandex.com/support/webmaster/en/robot-workings/check-yandex-robots",
	},
	{
		Name:      "Amzn-SearchBot",
		UserAgent: "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Amzn-SearchBot/0.1) Chrome/W.X.Y.Z Safari/537.36",
		Tokens:    []string{"Amzn-SearchBot"},
		Source:    "https://developer.amazon.com/amazonbot",
	},
}

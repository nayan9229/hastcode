# 03: make crawl-check

**What to build:** The Author runs one read-only command, `make crawl-check`, against live Production. In plain words it says whether search and answer engines can still reach the site, and it exits non-zero if any of them can't. It replaces the hand-typed `curl` loop in "Weekly check" (`docs/research/indexing-submission-channels.md`, runbook). Two things outside the build can shut a crawler out without anyone noticing. One is Cloudflare's managed robots.txt block, which Cloudflare prepends to our file. The other is an edge rule that returns 403 to a crawler. Both are live today: the managed block disallows 31 user agents, and the edge 403s PerplexityBot (note TL;DR items 1 and 2, §0 "robots.txt, line by line" and "Edge blocking"). Either could also undo ticket 01 later.

It checks:

- **robots.txt.** `/robots.txt` still holds our explicit allow groups (`Googlebot`, `Bingbot`, `OAI-SearchBot`, `PerplexityBot`) and both `Sitemap:` lines (`/sitemap.xml`, `/news-sitemap.xml`). No group in the file, including Cloudflare's prepended managed block, disallows `/` for any crawler on the list below. The check reads the groups the RFC 9309 way (§0): a crawler named in its own group follows only that group, and all the `User-agent: *` groups combine into one for everyone else. Applebot is the one exception. It is not named, and Apple says it then follows Googlebot's group (§0), so it is checked against that group.
- **Crawler access.** For each search or answer crawler, the command sends a GET for the home page and `/sitemap.xml` with that crawler's user-agent string. Every one should get `200`. The list: Googlebot, Bingbot, OAI-SearchBot, PerplexityBot, Claude-SearchBot, Applebot, DuckDuckBot, YandexBot and Amzn-SearchBot. These are the crawlers the note says we want reaching us (§0, §4, §6). All of them got 200 on 2026-11-16 except PerplexityBot. Each user-agent string is the one its owner publishes. The strings are kept in one list, each with a link to its owner's page (taken from the note's Sources), so they are easy to update. A `403` is reported as "Cloudflare is blocking <crawler>", with a pointer to runbook step 1. Any other non-200 answer is reported with its code. Training crawlers (GPTBot, ClaudeBot, CCBot, Amazonbot) are not checked, because blocking them is the Author's choice (ticket 01).
- **Redirects** (checked without following them). `http://hastcode.com/` answers `301` or `308` with location `https://hastcode.com/`, or the command says "turn on Always Use HTTPS (runbook step 2)". `https://www.hastcode.com/` still redirects to the apex (§0).

On the 2026-11-16 observations, it fails on PerplexityBot and on the http redirect until ticket 01 fixes them, and passes after that.

Limits, stated in its output and in the README: it sends crawler user-agent strings from the Author's own IP address. So it shows what Cloudflare does with a user agent, not what a verified crawler sees from its own IP addresses (§0 "Edge blocking"). It only reads. It never submits anything, never pings IndexNow or a WebSub hub, and never touches Search Console or Bing. It runs locally only, like `make intake`, and never in GitHub Actions: a weekly check from CI's IP addresses would test a different path through Cloudflare's edge. Pushes and tags also don't need a live-site check to pass (ADR-0004: Actions tests, builds Previews and publishes).

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] `make crawl-check` prints one plain line per check (pass or fail and why). It ends with a one-line summary, then exits 0 when everything passes and non-zero on any failure, including a network error or timeout.
- [ ] It reads the Production host from `site.yaml` `site_url`. It derives the `http://` and `www.` addresses from that host, and each of the three can be overridden. The tests use the overrides to run it against local stand-in servers, so `make test` never touches the network.
- [ ] Tests cover:
  - all checks passing → exit 0
  - one crawler getting 403 → non-zero, and the message names the crawler and runbook step 1
  - http answering 200 → non-zero, and the message names Always Use HTTPS
  - a missing `Sitemap:` line or allow group → non-zero
  - a prepended group that disallows `/` for a checked crawler → non-zero
  - an unreachable host → non-zero with a plain message
- [ ] It makes only GET or HEAD requests to the site under check, with a short timeout. It never sends a request to IndexNow, a WebSub hub or any search engine API (a test asserts that only the stand-in host was contacted).
- [ ] No workflow runs it; README says it is local only and why
- [ ] Run it once against live Production and paste the output here. The only failures should be the parts of ticket 01 not yet done: PerplexityBot and http if it hasn't started, none if it is finished.
- [ ] `docs/author-runbook.md`: today its "Getting found by search engines" paragraph (section 3) points to the note's hand-typed "Weekly check". It now gives a weekly check of its own: `make crawl-check` to copy and paste, what a failure means and which runbook step in the note fixes it, then the Search Console and Bing reports to look at. README documents the target beside "Feeds, sitemaps, robots and IndexNow", and the Makefile comment explains it like the other targets. Keep the paragraph's other lines (the get-indexed steps, the IndexNow line, any hub-ping lines from ticket 05) as they are; tickets 01 and 05 edit the same paragraph.
- [ ] `make test` passes; `make lint-workflows` passes

## Comments

**Live run, 2026-11-17** (`make crawl-check` from the Author's laptop, before ticket 01). The two failures are the parts of ticket 01 not yet done: PerplexityBot and the http redirect.

```
go run ./cmd/crawlcheck -config site.yaml 
crawl-check: read-only GETs to https://hastcode.com, http://hastcode.com and https://www.hastcode.com.
It sends crawler user agents from this machine's IP address, so it shows what Cloudflare does with each user agent, not what a verified crawler sees from its own addresses.
PASS robots.txt has our Googlebot allow group
PASS robots.txt has our Bingbot allow group
PASS robots.txt has our OAI-SearchBot allow group
PASS robots.txt has our PerplexityBot allow group
PASS robots.txt has Sitemap: https://hastcode.com/sitemap.xml
PASS robots.txt has Sitemap: https://hastcode.com/news-sitemap.xml
PASS robots.txt lets Googlebot crawl / (it follows the Googlebot group)
PASS robots.txt lets Bingbot crawl / (it follows the Bingbot group)
PASS robots.txt lets OAI-SearchBot crawl / (it follows the OAI-SearchBot group)
PASS robots.txt lets PerplexityBot crawl / (it follows the PerplexityBot group)
PASS robots.txt lets Claude-SearchBot crawl / (it follows the * groups)
PASS robots.txt lets Applebot crawl / (it follows the Googlebot group)
PASS robots.txt lets DuckDuckBot crawl / (it follows the * groups)
PASS robots.txt lets YandexBot crawl / (it follows the * groups)
PASS robots.txt lets Amzn-SearchBot crawl / (it follows the * groups)
PASS Googlebot gets 200 on / and /sitemap.xml
PASS Bingbot gets 200 on / and /sitemap.xml
PASS OAI-SearchBot gets 200 on / and /sitemap.xml
FAIL Cloudflare is blocking PerplexityBot: 403 on / and /sitemap.xml with its user agent. Fix: runbook step 1 in docs/research/indexing-submission-channels.md
PASS Claude-SearchBot gets 200 on / and /sitemap.xml
PASS Applebot gets 200 on / and /sitemap.xml
PASS DuckDuckBot gets 200 on / and /sitemap.xml
PASS YandexBot gets 200 on / and /sitemap.xml
PASS Amzn-SearchBot gets 200 on / and /sitemap.xml
FAIL http://hastcode.com/ answers 200 instead of redirecting to https://hastcode.com/: turn on Always Use HTTPS (runbook step 2 in docs/research/indexing-submission-channels.md)
PASS https://www.hastcode.com/ redirects (301) to https://hastcode.com/
crawl-check: 2 of 26 checks failed.
exit status 1
make: *** [crawl-check] Error 1
```

User-agent strings: Googlebot, OAI-SearchBot, PerplexityBot, Applebot, DuckDuckBot, YandexBot and Amzn-SearchBot were checked against their owners' pages on 2026-11-17. Bing's crawler page is client-rendered and could not be read, so Bingbot's is the desktop string Bing has published; Anthropic publishes only the `Claude-SearchBot` token, not a full string, so that one follows the shape of Anthropic's other crawlers' strings. Both are noted in `internal/crawlcheck/crawlers.go`.

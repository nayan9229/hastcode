# 01: Get-indexed dashboard steps

**What to build:** Nothing in code. The Author does the dashboard steps that only an account owner can do, so search and answer engines can reach Production and know what to crawl. Then the Author writes down what each check showed. The steps are in "Runbook: get indexed, in order", steps 1 to 4, in `docs/research/indexing-submission-channels.md`. Follow them there. This ticket does not copy them; it lists what to record. The evidence comes from the same note: TL;DR, §0 "What the site emits today", §1 (Search Console) and §2 (Bing). Request indexing has a quota and repeats don't help: `docs/research/first-traffic-new-domain.md` §1.1.

These facts were observed on 2026-11-16. With PerplexityBot's user agent, Cloudflare's edge answers 403, while a browser and Googlebot's user agent get 200 (§0 "Edge blocking"). `http://hastcode.com/` answers 200 instead of redirecting to HTTPS. A `google-site-verification` TXT record is live, but nobody has yet confirmed in Search Console who owns the property or whether it is still verified (§1.2).

IndexNow already runs on every Production Publish: `INDEXNOW_ENABLED` is `true`, and every ping since 2026-10-30 has got `200` (§0 "IndexNow wiring", runbook step 5). **Do not send IndexNow by hand** and do not use Bing's Submit URLs. Leave the Google Indexing API alone too (§1.5). The dashboard work changes no file. The only repo changes are the notes written into this file and, once those notes exist, the two doc edits in the last box.

**Blocked by:** None (can start immediately)

**Status:** ready-for-human

- [ ] **Why PerplexityBot is blocked.** In Cloudflare, find which setting returns the 403. It could be an AI Crawl Control Block on PerplexityBot, a WAF or custom rule, or the "Block AI bots" toggle (note, Open questions 1). Use AI Crawl Control → Security, and Security → Events filtered by user agent. Record which one it was, here.
- [ ] **PerplexityBot allowed** (runbook step 1). The step-1 `curl` with PerplexityBot's user agent now prints `200`. Record the date and the code.
- [ ] **Training crawlers decided on purpose** (runbook step 1.3). Record the choice for GPTBot, ClaudeBot, CCBot, Amazonbot, Google-Extended and Applebot-Extended (Blocked or Allowed). Also confirm that OAI-SearchBot, Claude-SearchBot, Perplexity-User and Amzn-SearchBot are still allowed.
- [ ] **Always Use HTTPS on** (runbook step 2). `curl -sI http://hastcode.com/` now shows a `301` or `308` with `location: https://hastcode.com/`. Record the status line.
- [ ] **Search Console Domain property** (runbook step 3.1). Record that the `hastcode.com` Domain property exists, that it is verified, and that the Author is Owner. The `google-site-verification` TXT record stays in Cloudflare DNS.
- [ ] **Both sitemaps submitted** (runbook step 3.2). Both `/sitemap.xml` and `/news-sitemap.xml` show **Success** under Indexing → Sitemaps. Record each one's status and its discovered-URL count. Resubmit later only if one says Couldn't fetch.
- [ ] **Request indexing, once per URL** (runbook step 3.3). Done for the home page and the newest 3 to 5 Pieces only. Record which URLs.
- [ ] **Baseline Page indexing counts** (runbook step 3.4). Record today's Indexed, "Discovered – currently not indexed" and "Crawled – currently not indexed" counts, to compare against weekly.
- [ ] **Bing Webmaster Tools import** (runbook step 4). The site has been imported from Search Console, and both sitemaps are listed under Sitemaps. Record the date. Within 48 h, check that IndexNow Insights shows recent submitted URLs from the CI pings, and record what it shows.
- [ ] Optional: if you find which service the unlabelled apex TXT `78ace81ef8ec61ef2998c8c214526a80` belongs to, record it (note, Open questions 3).
- [ ] If ticket 03 (`make crawl-check`) is done by then, it passes. Paste its output here.
- [ ] Once the PerplexityBot and HTTPS results are recorded, the README section "Feeds, sitemaps, robots and IndexNow" no longer says the edge blocks PerplexityBot. Today it says "PerplexityBot as of 2026-11-16". The README instead names whatever the edge still blocks, with its date. The "Getting found by search engines" paragraph in `docs/author-runbook.md` (section 3, "One-time setup") still matches what remains to do. An agent may make these two doc edits from the results recorded here, and nothing else. Change only the PerplexityBot and HTTPS wording. Keep any lines tickets 02, 03 or 05 added to the same README section or runbook paragraph.

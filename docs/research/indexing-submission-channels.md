# Getting hastcode.com indexed: how to submit it, and where else

Research for hastcode.com: a static Go-built publication served from Cloudflare Workers Static Assets, ad-monetised, live since 2026-10-25, 87 URLs in `/sitemap.xml` today. The Author asked: "How do I get the site indexed: how do I submit it, and where else can I submit it?"

This note builds on three earlier ones and does not repeat them:
- [`first-traffic-new-domain.md`](first-traffic-new-domain.md) (2026-10-30): sitemaps are hints, Request indexing has an unpublished daily limit, Bing import from Search Console, IndexNow basics, Cloudflare Crawler Hints, how long indexing takes.
- [`google-news-discover-requirements.md`](google-news-discover-requirements.md) (2026-10-20): Google News is automatic, Publisher Center takes no submissions, Bing PubHub and Apple News Publisher requirements.
- [`aeo-citation-signals.md`](aeo-citation-signals.md) (2026-10-20): the robots.txt token for each answer engine, and what gets a page cited.

Researched **2026-11-16**. Labels: **[P]** primary (the engine's or service's own docs, help centre, blog or spec, fetched that day) / **[P, snippet]** a first-party page I could not fetch (client-rendered, blocked or timed out), where I saw its text only as a search-engine snippet, so treat it as unverified / **[observed]** what a live request, DNS lookup, CI log or the repo showed today / **[inferred]** my deduction / **[opinion]**.

---

## TL;DR

- **Cloudflare blocks PerplexityBot at the edge, even though our robots.txt allows it.** A request with PerplexityBot's user agent gets `HTTP 403` and "Your request was blocked." from `server: cloudflare`. A plain browser gets 200 [observed]. Perplexity says PerplexityBot is what "surface[s] and link[s] websites in search results on Perplexity" [P]. Right now this is the one thing actively keeping us out of an answer engine. The fix is a dashboard toggle (runbook step 1).
- **Cloudflare adds its own block-list to the top of our robots.txt.** The live `/robots.txt` starts with a "Cloudflare Managed Content" block: a `Content-Signal: search=yes,ai-train=no,use=reference` line, then `Disallow: /` for 31 user agents. These include ClaudeBot, GPTBot, CCBot, Amazonbot, Google-Extended, Applebot-Extended and GoogleOther, and two search engines, **Baiduspider and PetalBot**. Our own allow groups follow it [observed]. Cloudflare says it "prepends" this when Bot Preference Sync is on [P]. No search crawler we care about is blocked by it. The edge test also returns 403 for GPTBot, ClaudeBot, CCBot and Amazonbot. Googlebot, Bingbot, OAI-SearchBot, Claude-SearchBot, Applebot, DuckDuckBot, YandexBot, SeznamBot, Yeti (Naver), MojeekBot, AhrefsBot (Yep) and archive.org_bot all get 200 [observed].
- **IndexNow now works. The earlier note's main blocker is fixed.** `INDEXNOW_ENABLED` has been `true` since 2026-10-29. The first ping after that sent 66 URLs and got `202` (key pending). Every ping since has got `200`, most recently 27 URLs on 2026-11-15 [observed]. One ping reaches all seven engines in IndexNow's registry: Bing, Yandex, Seznam, Naver, Yep, Amazon and the Internet Archive [P]. Signing up to those engines' own consoles adds reports, not reach.
- **Only two dashboards matter: Google Search Console and Bing Webmaster Tools.** DNS already has a `google-site-verification=` TXT record, so a Domain property has been set up at least once [observed]. Bing can import the property and its sitemaps from Search Console [P]. Yahoo, DuckDuckGo, Ecosia (by default), Copilot and ChatGPT search all draw on Bing's index, per those companies' own pages [P]. Startpage is said to show Google results, but I could not fetch its pages.
- **The Google Indexing API is not for us.** It "can only be used to crawl pages with either JobPosting or BroadcastEvent embedded in a VideoObject", and "All submissions through the Indexing API undergo rigorous spam detection" [P]. Using it for articles is outside its stated scope.
- **No AI answer engine has a submission form.** Each one only needs its crawler let in. Google AI Overviews need Googlebot ("no additional requirements"). ChatGPT search needs `OAI-SearchBot`. Perplexity needs `PerplexityBot`, which is blocked today. Claude needs `Claude-SearchBot`, and blocking `ClaudeBot` affects training only. Copilot uses Bing [P]. Brave has its own index, follows Googlebot's robots rules and has a "submit your page for re-fetching" form [P].
- **Three small fixes help every crawler.** (1) `http://hastcode.com/` answers 200 instead of redirecting to HTTPS [observed]. Cloudflare "Always Use HTTPS" fixes that on the Free plan [P]. (2) The favicon is only a `data:` URI and `/favicon.ico` is 404 [observed, ADR-0006]. Google requires that "Googlebot-Image must be able to crawl the favicon file" [P], and Google News now takes its logo from the favicon (prior note). (3) There is no `<link rel="alternate" type="application/rss+xml">` in `<head>` [observed], so feed readers and Naver/Flipboard-style feed intake cannot find the feed automatically.
- **Skip these:** Baidu (Chinese-language, ICP-filing limits, and our robots.txt blocks Baiduspider), Naver and Seznam consoles (Korean and Czech audiences, and IndexNow reaches them anyway), Kagi Small Web (it excludes sites with ads and LLM-generated content), Apple News (country-limited, see prior note), and any "submit to 500 search engines" service (§8).

---

## Channel table

"Own index?" means whether the engine crawls and ranks itself, or shows another engine's results. "Worth it?" is for an English tech site with ads, one Author, three weeks old [opinion].

| Channel | Own index? | How to submit | Verification | Quota | Worth it? |
|---|---|---|---|---|---|
| **Google Search** | Yes | Search Console → Sitemaps (both files); URL Inspection → Request indexing for a handful of key URLs; `Sitemap:` lines in robots.txt (already there) | Domain property: DNS TXT (a `google-site-verification` TXT is live) | Request indexing: "a quota", number unpublished; sitemap ≤ 50,000 URLs / 50 MB | **Yes**, confirm done |
| Google Indexing API | n/a | Only for JobPosting / BroadcastEvent pages | Service account | 200 publish requests/day default | **No**: outside its scope |
| **Bing** | Yes | BWT → Import from Google Search Console (brings sitemaps); Sitemaps; Submit URLs; IndexNow (live) | GSC import, Domain Connect, XML file, `msvalidate.01` meta tag, DNS TXT | Submit URLs up to 10,000/day, adaptive [P blog 2019; help page unreadable]; batch ≤ 500 per call | **Yes** |
| **IndexNow** | Protocol shared by 7 engines | `POST https://api.indexnow.org/indexnow` (CI does this on every Production Publish) | Key file `/465f05c90b9917457c1ad83b8df48051.txt` | ≤ 10,000 URLs per POST; no published rate; 429 = slow down | **Yes**, live |
| Yahoo, DuckDuckGo, Ecosia (default) | No: Bing (DDG "largely" Bing, plus its own DuckDuckBot) | Through Bing | — | — | Covered by Bing |
| Startpage | No: Google results [unverified, site blocked to this fetcher] | Through Google | — | — | Covered by Google |
| Qwant | Own crawler, with a Bing fallback; EUSP/Staan index with Ecosia | No submission found | — | — | No action |
| Yandex | Yes | Yandex Webmaster: Sitemap files, Reindex pages; already receives IndexNow | HTML file (recommended), meta tag, DNS TXT | Reindex: a daily limit, number unpublished; IndexNow: no request-count limit | **Later** (for its reports) |
| Naver | Yes | Search Advisor: sitemap + RSS; receives IndexNow | Meta tag or HTML file | — | No: Korean audience |
| Seznam | Yes | Seznam Webmaster; sitemap via robots.txt; receives IndexNow | `seznam-wmt` meta tag or file | add-URL API 500/day [P, snippet] | No: Czech audience |
| Baidu | Yes | Ziyuan: API push, sitemap, manual | File, meta tag, CNAME | API + manual: 100k/day shared [P, snippet]; sites without ICP filing: 1 sitemap/day [P, snippet, 2021] | **No** |
| Yep | Ahrefs' index (AhrefsBot) | Receives IndexNow | — | — | Covered already |
| Mojeek | Yes | No add-URL; does not use sitemaps; earned links only | — | — | No action |
| Brave Search | Yes; follows Googlebot's robots rules | `search.brave.com/submit-url` ("submit your page for re-fetching") | None | Not stated | **Later**, one-off |
| Apple (Spotlight, Siri, Safari) | Yes (Applebot) | None. Applebot follows Googlebot's rules when not named | — | — | No action (allowed) |
| Kagi | Teclis / TinyGem plus paid providers | Small Web list only (personal, single-author, no ads, no LLM content) | GitHub PR | — | No: not eligible |
| Google News / Discover | Google | None: automatic | — | — | No action (prior note) |
| Bing News PubHub | Bing | Inside BWT | BWT | — | **Later**, after a news track record |
| Apple News Publisher | Apple | News Publisher | Apple Account with 2FA | AU/CA/UK/US publishers | No (prior note) |
| Flipboard | Flipboard | Add RSS in the app / FlipBlogger@Flipboard.com (2017-era page) | — | Full feed, ≥ 30 items | Later |
| ChatGPT search | OAI-SearchBot (plus Bing) | None | — | robots change ~24 h | Yes: keep allowed (it is) |
| **Perplexity** | PerplexityBot | None | — | robots change ~24 h | **Yes: unblock it** |
| Claude | Claude-SearchBot | None | — | — | Yes: allowed (it is) |
| Google AI Overviews / AI Mode | Googlebot | None | — | — | Covered by Google |
| Copilot | Bing | Through Bing | — | — | Covered by Bing |
| Amazon (Alexa etc.) | Amzn-SearchBot; Amazonbot | Receives IndexNow | — | ~24 h | Covered (Amazonbot blocked: training) |
| Common Crawl | Yes (CCBot) | No form; reads robots.txt `Sitemap:` | — | — | Owner decision (blocked today) |
| Internet Archive | Yes | Save Page Now (one page at a time); receives IndexNow | — | — | Optional |
| WebSub (Google-run hub) | — | `atom:link rel="hub"` + a hub ping | — | — | Later (code change) |

---

## 0. What the site emits today [observed]

All of these are from `builder/` (step files named below) and from live GETs on 2026-11-16.

| URL | Status | What it is | Written by |
|---|---|---|---|
| `https://hastcode.com/robots.txt` | 200 | Cloudflare Managed block, then ours: `*`, `Googlebot`, `Bingbot`, `OAI-SearchBot`, `PerplexityBot` all `Allow: /`; two `Sitemap:` lines | `step_write_robots.go` + Cloudflare |
| `https://hastcode.com/sitemap.xml` | 200 | 87 `<loc>` with `lastmod` (14 Briefs, 63 Explainer/Vertical/kind pages, home, About, Contact, Editorial Policy, Privacy, Authors) | `step_write_sitemap.go` |
| `https://hastcode.com/news-sitemap.xml` | 200 | Briefs from the last 48 h (1 today) | `step_write_news_sitemap.go` |
| `https://hastcode.com/feed.xml`, `/briefs/feed.xml`, `/<vertical>/feed.xml` | 200 | Full-text RSS 2.0 with `atom:link rel="self"`; no `rel="hub"` | `step_write_feeds.go` |
| `https://hastcode.com/llms.txt`, `/sitemap.md`, `*.md` twins | 200 | Agent Surface, `X-Robots-Tag: noindex` (ADR-0011) | `step_write_llms_txt.go`, `step_write_markdown_sitemap.go` |
| `https://hastcode.com/465f05c90b9917457c1ad83b8df48051.txt` | 200 | IndexNow key file; body is the key | `step_write_indexnow.go`; key in `site.yaml` `indexnow_key` |
| `https://hastcode.com/ads.txt` | 200 | AdSense record | `site.yaml` `ads_txt` |
| `https://www.hastcode.com/` | 301 → apex | Fixed since the 2026-10-30 note | Cloudflare rule |
| `http://hastcode.com/` | **200** (no redirect) | Serves the page over HTTP; `rel="canonical"` points to https | Cloudflare setting |

**robots.txt, line by line.** There are two `User-agent: *` groups, Cloudflare's and ours. RFC 9309 says "the matching groups' rules MUST be combined into one group" [P], so the result is `Allow: /` plus Cloudflare's Content-Signal. A crawler named in its own group uses only that group. Applebot is not named but Googlebot is, so "the Apple robot will follow Googlebot instructions" [P, Apple]. Brave's crawler has no user agent of its own and will not crawl what Googlebot cannot [P, Brave]. Disallowed by the Cloudflare block: Amazonbot, Applebot-Extended, Bytespider, CCBot, ClaudeBot, Diffbot, Google-Extended, GPTBot, omgili, anthropic-ai, Claude-Web, cohere-ai, MistralAI-Training, meta-externalagent, GoogleOther, **Baiduspider**, **PetalBot**, AwarioSmartBot/AwarioRssBot, Google-CloudVertexBot, QualifiedBot, Cotoyogi, ICC-Crawler, atlassian-bot, FishBot, BorderxBot, NavuBot, SemrushBot-SWA, WARDBot, magpie-crawler, KimiBot, CitibotSiteCrawler.

What blocking each of those costs us in search, from the owners' own docs [P]:
- `ClaudeBot`: training only. Search visibility depends on `Claude-SearchBot` (Anthropic).
- `GPTBot`: training only. "Each setting is independent of the others" (OpenAI).
- `Google-Extended`: "does not impact a site's inclusion in Google Search" (prior AEO note).
- `GoogleOther`: its rules "don't affect any specific product". `Google-CloudVertexBot` "has no effect on Google Search" (Google common crawlers doc).
- `Applebot-Extended`: "Webpages that disallow Applebot-Extended can still be included in search results" (Apple).
- `Amazonbot`: "may also be used to train Amazon AI models". `Amzn-SearchBot`, which feeds Alexa search, is separate and allowed (Amazon, updated 2026-10-26).
- `CCBot`: Common Crawl. No search engine documents relying on it [inferred].
- `Baiduspider`, `PetalBot`: these are search crawlers (Baidu, Huawei Petal Search). Cloudflare's own verified-bots page lists "Baidubot" as a Search Engine Crawler [P]. Blocking them only matters if we ever want China or Huawei-device search [inferred].

**Edge blocking, separate from robots.txt.** I sent a GET to one article with each crawler's published user-agent string, from a normal (non-crawler) IP:
- **403**: GPTBot, ClaudeBot, CCBot, Amazonbot, **PerplexityBot**.
- **200**: everything else tested (Googlebot, Bingbot, OAI-SearchBot, ChatGPT-User, Claude-SearchBot, Claude-User, Perplexity-User, Applebot, DuckDuckBot, DuckAssistBot, YandexBot, SeznamBot, Yeti, MojeekBot, AhrefsBot, Amzn-SearchBot, Baiduspider, PetalBot, archive.org_bot, facebookexternalhit).

Cloudflare's AI Crawl Control "identifies crawlers by user agent string" on the free plan, and per-crawler Allow/Block lives at AI Crawl Control → Security → Crawlers [P]. So the PerplexityBot 403 is almost certainly an AI Crawl Control Block on PerplexityBot [inferred]. A real crawler from its own IPs may be treated differently, but the UA match makes that unlikely [inferred].

**IndexNow wiring.** The build writes `.vercel/output/indexnow.json` (`{host, key, keyLocation, urlList}`) listing pages whose `<main>` changed since the last Production build. `.github/workflows/publish.yml` step "Ping IndexNow" POSTs it to `https://api.indexnow.org/indexnow` when `vars.INDEXNOW_ENABLED == 'true'` [observed]. Results [observed, `gh run view`]:

| Run | Date | URLs | Response |
|---|---|---|---|
| 40943579139 | 2026-10-29 | 66 | 202 |
| 41054099392 | 2026-10-30 | 22 | 200 |
| 43691945311 | 2026-11-15 | 27 | 200 |

Previews and QA delete the artefact, so they never ping.

**Verification records.** Apex TXT records: `google-site-verification=DYIiThuL…`, an unlabelled 32-hex value `78ace81ef8ec61ef2998c8c214526a80`, and SPF. No `msvalidate.01`, `yandex-verification`, `naver-site-verification` or `seznam-wmt` meta tag is in the home page `<head>` [observed]. I do not know which service the unlabelled TXT value belongs to (Open questions).

---

## 1. Google Search Console, in more depth

### 1.1 Domain property vs URL-prefix property

- **Domain property** (`hastcode.com`): "aggregates data for all subdomains, protocols, and subpaths", so it covers http, https, www and any future `images.` or `qa.` subdomain. Its only verification method is DNS ("DNS record verification"). It is required "unless the property is on a Google product such as Blogger or Google Sites" [P, [answer/34592](https://support.google.com/webmasters/answer/34592)].
- **URL-prefix property** (`https://hastcode.com/`): only URLs starting with that exact prefix, protocol and host included. It has "Many possible methods" of verification [P, same].
- **Use the Domain property** [opinion]. It also shows `http://` and `www` traffic, which matters while `http://` is not redirected (§0).

### 1.2 DNS verification, on Cloudflare

- Google's steps: pick **TXT**, add the record at the domain host with Name blank or `@`, and paste the string as the value. If it fails, "wait a day or two". Manually added records "can take two or three days" [P, [answer/9008080](https://support.google.com/webmasters/answer/9008080)].
- **Keep the record forever.** "To stay verified, don't remove the DNS record from your provider, even after verification succeeds". Search Console re-checks periodically, and after a grace period "your permissions expire" [P, same].
- On Cloudflare: DNS → Records → **Add record** → choose a Type → fill the fields → **Save** [P, [Cloudflare create DNS records](https://developers.cloudflare.com/dns/manage-dns-records/how-to/create-dns-records/), 2026-05-06]. Google's page does not mention Cloudflare by name. If Search Console offers Cloudflare in its provider list, the automatic route "should then happen immediately" [P, answer/9008080].
- **Our state**: the TXT record is live [observed]. Whether the property is still verified, and who owns it, can only be seen inside Search Console (runbook step 3).

### 1.3 Sitemaps report

- Steps [P, [answer/7451001](https://support.google.com/webmasters/answer/7451001)]: you need owner permission, otherwise "list the sitemap in your robots.txt" (we do). Test the URL with a live URL inspection, paste it into "Add a new sitemap", then **Submit**.
- Statuses: **Success** means fetched and parsed, and all URLs are "queued for crawling". **Has errors** means it parsed with problems, but clean URLs are still queued. **Couldn't fetch** means a robots block, 404 or server error. After a failed fetch, Google "retries for a few days and then stops", so fix it and resubmit [P, same].
- Resubmitting a healthy sitemap is unnecessary: Google re-reads it on its own schedule [P, same]. The old `google.com/ping?sitemap=` endpoint is gone: "Pings to the old endpoint will return a 404" [P, [Search Central blog 2023-06-26](https://developers.google.com/search/blog/2023/06/sitemaps-lastmod-ping)]. Do not build a ping step.
- `lastmod` must "consistently match reality" or Google may stop trusting it [P, same]. Our builder ties it to a Piece's review/publication date (README "Feeds, sitemaps, robots and IndexNow") [observed]. That is the right design.
- Limits: 50 MB uncompressed, 50,000 URLs per sitemap, news sitemaps 1,000 URLs [P, answer/7451001]. We are far below them.

### 1.4 URL Inspection and Request indexing

Covered in [`first-traffic-new-domain.md`](first-traffic-new-domain.md) §1.1: an unpublished daily quota, owners and full users only, and "requesting a recrawl multiple times for the same URL won't get it crawled any faster" [P, [Ask Google to recrawl](https://developers.google.com/search/docs/crawling-indexing/ask-google-to-recrawl)]. Google's own split: URL Inspection is for "individual URLs", and a sitemap is for "large numbers of URLs" [P, same].

### 1.5 Indexing API: not for articles

- Scope, verbatim: "The Indexing API can only be used to crawl pages with either JobPosting or BroadcastEvent embedded in a VideoObject." [P, [quickstart](https://developers.google.com/search/apis/indexing-api/v3/quickstart); [using the API](https://developers.google.com/search/apis/indexing-api/v3/using-api)].
- "All submissions through the Indexing API undergo rigorous spam detection". Abuse, such as using several accounts to exceed quotas, "may lead to access being revoked" [P, quickstart].
- Default quota: 200 publish requests/day per project, 180 metadata reads/minute, 380 requests/minute overall [P, [quota page](https://developers.google.com/search/apis/indexing-api/v3/quota-pricing)].
- **For hastcode**: our pages are Article/NewsArticle, so they are out of scope by Google's own definition. The pages I fetched do not use the words "against policy". The accurate statement is "outside its documented scope, and every submission is spam-screened". Many SEO plugins advertise using it for any URL. Do not [opinion].

### 1.6 Page indexing report: statuses and what Google says to do

From [answer/7440203](https://support.google.com/webmasters/answer/7440203) [P]:

| Status | Meaning (Google) | What Google says to do | For us [inferred] |
|---|---|---|---|
| **Discovered – currently not indexed** | Found but not crawled yet. Google likely rescheduled because crawling then "would have overloaded the site". Last crawl date is blank | No specific action given | Normal for a 3-week-old site. Internal links and time help; Request indexing only for the few that matter |
| **Crawled – currently not indexed** | Crawled, then not indexed; "doesn't necessarily indicate an issue that needs attention" | For important pages, check them with URL Inspection; if you suspect a site-wide bug, use the report-an-indexing-issue form | If a strong Explainer sits here for weeks, the remedy is quality and links, not resubmitting |
| URL marked 'noindex' | Google saw a noindex | Remove it if you want the page indexed | Expected for `/search/`, `/subscribed/`, 404, `*.md` twins, `llms.txt`, `*.workers.dev` |
| Alternate page with proper canonical tag | Points to an indexed canonical | Nothing | Expected for `http://` variants until HTTPS redirect is on |
| Duplicate without user-selected canonical | Google picked another URL | Check Google's choice; set the canonical explicitly | Our pages all carry `rel="canonical"` [observed] |
| Page with redirect | Non-canonical, redirects | Inspect the target | Expected for `www` and the `/home` alias |

Validate fix: fix every instance first, click **Validate fix** once, and expect "up to about two weeks" [P, same].

### 1.7 Two Google details that touch us

- **HTTPS**: "Google prefers HTTPS pages over equivalent HTTP pages as canonical". Redirects are "A strong signal" [P, [Consolidate duplicate URLs](https://developers.google.com/search/docs/crawling-indexing/consolidate-duplicate-urls)]. Our `http://` 200 plus canonical is acceptable, but a 301 is cleaner [opinion].
- **Favicon**: Google needs "Googlebot-Image must be able to crawl the favicon file". Supported `rel` values are `icon`, `shortcut icon`, `apple-touch-icon` and `apple-touch-icon-precomposed`. It must be square, at least 8x8, with a recommended size "larger than 48x48px", and the URL "must be stable" [P, [Favicon in Search](https://developers.google.com/search/docs/appearance/favicon-in-search)]. The page does not mention `data:` URIs. Ours is only a `data:` URI (ADR-0006, chosen for Lighthouse speed), so there is no file to crawl [observed]. I infer Google may show a generic icon in results and in auto-generated Google News publication pages. This trade-off is the Author's decision under ADR-0006 (Open questions).

---

## 2. Bing Webmaster Tools

- **Bing's help pages are client-rendered.** Fetching `bing.com/webmasters/help/*` returns only UI strings (third time, after the 2026-10-30 and 2026-11-13 notes). The claims below come from Bing's blog and Microsoft Learn [P], plus labelled snippets.
- **Verification** [P, [Bing blog 2025-06-17](https://blogs.bing.com/webmaster/2025/6/Start-Using-Bing-Webmaster-Tools-to-Improve-Your-Site-Visibility/)]: "Domain Connect (recommended if available)", "Import your site from Google Search Console", "Upload an XML verification file", "Add a meta tag to your homepage", "Add a DNS TXT record manually". **Import is the cheapest for us**: it brings "verified sites and sitemaps" across (prior note).
- **Sitemaps**: "Bing processes sitemaps at least once every 24 hours" [P, same].
- **Submit URLs (tool and API)**: the API was launched at "up to 10,000 URLs per day". The quota is "determined based on the site verified age in Bing Webmaster tool, site impressions and other signals" [P, [Bing blog 2019-01](https://blogs.bing.com/webmaster/2019/1/bingbot-Series-Get-your-content-indexed-fast-by-now-submitting-up-to-10%2C000-URLs-per-day-to-Bing/)]. Batch calls: "The maximum supported batch size in this API is 500 URLs per request", and "the daily submission limit of up to 10,000 URLs still applies". Bing "now recommends IndexNow as the primary" real-time method [P, [Bing blog 2019-06, with later note](https://blogs.bing.com/webmaster/2019/6/bingbot-Series-Introducing-Batch-mode-for-Adaptive-URL-submission-API/)]. The current help page reportedly says 10,000 per domain per day, resets at midnight UTC, and the API "may be deprecated in the future" [P, snippet, [URL Submission help](https://www.bing.com/webmasters/help/URL-Submission-62f2860b)]. `GetUrlSubmissionQuota` returns your site's real figure [P, [Learn](https://learn.microsoft.com/en-us/dotnet/api/microsoft.bing.webmaster.api.interfaces.iwebmasterapi.geturlsubmissionquota?view=bing-webmaster-dotnet)].
- **API notice**: "Legacy SOAP and POX APIs will be retired on August 31, 2026. Migrate to our REST APIs" [P, [Bing Webmaster API](https://learn.microsoft.com/en-us/bingwebmaster/), updated 2026-08-07]. Whether the documented `…/api.svc/json/SubmitUrlbatch` endpoint counts as "REST" after that date is not stated (Open questions). Since IndexNow already covers this, we do not need the Bing API at all [inferred].
- **API key**: BWT → Settings → API Access → accept terms → **Generate API Key**. It is "generated for a user and not a site" [P, [Getting access](https://learn.microsoft.com/en-us/bingwebmaster/getting-access)].
- **Content Submission API** (`SubmitContent`): still documented. It takes a base64-encoded HTTP response, up to "10MB payload uncompressed", with a separate `GetContentSubmissionQuota` [P, [Learn, 2023-11-14](https://learn.microsoft.com/en-us/dotnet/api/microsoft.bing.webmaster.api.interfaces.iwebmasterapi.submitcontent?view=bing-webmaster-dotnet)]. It uses the POX/SOAP paths being retired, so treat it as legacy and do not build on it [inferred].
- **IndexNow inside BWT**: the "IndexNow Insights" report shows counts of URLs "submitted, crawled, and indexed", your "latest 1000 submitted URLs" with index status and first-indexed date, and an "Important URLs Missing" tab [P, [Bing blog 2024-03-03](https://blogs.bing.com/webmaster/March-2024/Optimize-your-Impact-with-IndexNow-Insights)]. This is where to check that our CI pings are being used.
- **PubHub**: the BWT page bundle still carries PubHub strings ("Add your URLs/RSS/Sitemaps", "Reviewed at", "Bing PubHub Content License Agreement") [observed in page source, 2026-11-16]. `pubhub.bing.com` did not answer [observed]. A search snippet mentioned a retirement banner that I could not see (Open questions). Requirements are in the [prior note](google-news-discover-requirements.md) §8.

---

## 3. IndexNow: exact format

From the spec [P, [indexnow.org/documentation](https://www.indexnow.org/documentation); [FAQ](https://www.indexnow.org/faq)]:

- **Single URL (GET)**: `https://<searchengine>/indexnow?url=<url-encoded>&key=<key>`. Add `&keyLocation=` only if the key file is not at the root.
- **Bulk (POST)**: `POST /indexnow` with `Content-Type: application/json; charset=utf-8`. Body fields: `host`, `key`, `keyLocation` (needed only for a non-root key), and `urlList`. Up to **10,000 URLs** per POST.
- **Key rules**: 8–128 characters of `a-z A-Z 0-9 -`. The file is named `{key}.txt`, UTF-8, and contains the key. Root is recommended. A key file in `/catalog/` only covers `/catalog/…` URLs. Ours, `465f05c90b9917457c1ad83b8df48051`, is 32 hex characters at the root and matches `builder/step_validate_indexnow_key.go` [observed].
- **Responses**:

| Code | Meaning |
|---|---|
| 200 | Received. Received only, not indexed |
| 202 | Received; key validation pending |
| 400 | Invalid format |
| 403 | Key not valid (missing file or wrong key) |
| 422 | URL not on host, or key/protocol mismatch |
| 429 | Too many requests, "potential spam" |

- **Etiquette**: "Submit only when content has changed". Wait at least 5 minutes before resubmitting a frequently updated URL, and after a 429 "wait at least 10 minutes". A submission "does not guarantee immediate indexing" [P, FAQ].
- **Sharing**: engines "agree to share submitted URLs automatically", so any one endpoint is enough [P, documentation].
- **Who receives it** ([searchengines.json](https://www.indexnow.org/searchengines.json)) [P]: bing, yandex, seznam, naver, yep, **internetarchive**, **amazonbot**. The FAQ prose lists six (no Internet Archive). Google is not a participant.
- **Yandex** adds that there is no limit on the number of IndexNow requests, but asks you to pace them [P, snippet, [Yandex IndexNow](https://yandex.com/support/webmaster/en/indexing-options/index-now)]. **Seznam** says the same: no URL cap through IndexNow, unlike its 500/day add-URL API [P, snippet, [Seznam IndexNow](https://o-seznam.cz/napoveda/vyhledavani/seznambot/protokol-indexnow/)].
- **Amazon and our robots.txt**: Amazon takes IndexNow pings, and Cloudflare blocks `Amazonbot`. `Amzn-SearchBot` is allowed, and Amazon says it "improves search experiences in Amazon products such as Alexa" [P, [Amazonbot](https://developer.amazon.com/amazonbot)]. Which of its bots follows up on an IndexNow ping is not documented (Open questions).

---

## 4. Other search engines

**Bing reuse, no action.** Bing says it powers "Yahoo, DuckDuckGo, Microsoft Start, and AI tools like ChatGPT and Inflection.ai" [P, Bing 2025-06-17]. DuckDuckGo: results come "largely" from Bing, plus its own DuckDuckBot, which "respects" robots.txt. Its help pages describe no submission [P, [sources](https://duckduckgo.com/duckduckgo-help-pages/results/sources), [DuckDuckBot](https://duckduckgo.com/duckduckgo-help-pages/results/duckduckbot)]. Ecosia: results "come from our partners Microsoft Bing, Google and EUSP". Bing is the default, Google appears only in some regions with consent, and EUSP "is currently only available in France". Site owners can only advertise, not submit [P, [Ecosia help, updated 2026-08-12](https://support.ecosia.org/article/579-search-results-providers)].

**Startpage**: I could not verify anything. Its domain is blocked to this research tool. It is widely described as showing Google results; check by hand before relying on that.

**Qwant**: Qwant says it runs its own crawler. Older posts say it uses Bing "to supplement search results on which we do not have sufficient relevance". It now runs on the EUSP "Staan" index co-owned with Ecosia [P, snippets of [betterweb.qwant.com](https://betterweb.qwant.com/en/2024/11/08/ecosia-and-qwant-join-forces-to-develop-european-search-index/) and the [2025 transparency report](https://about.qwant.com/legal/rapport-de-transparence-qwant-2025/)]. No webmaster submission was found.

**Yandex Webmaster.** It has its own index.
- Verify by **HTML file (recommended)**, `yandex-verification` meta tag or **DNS TXT**. www and non-www are separate sites. Keep the code in place or "your rights" can be reset. Verification can take up to 24 h [P, [Yandex rights](https://yandex.com/support/webmaster/en/service/rights)].
- Notify it of pages through Reindex pages, Sitemap files, Tags crawl (needs a Yandex Metrica tag) or IndexNow. Pages get indexed "between 3 and 7 days" after a reindex request [P, [new/changed pages](https://yandex.com/support/webmaster/en/robot-workings/new-changed)]. Reindex has a daily limit "for your primary domain and all subdomains", number unpublished. Processed data "should update within two weeks" [P, [Reindex](https://yandex.com/support/webmaster/en/robot-workings/site-reindex)].
- We already reach it by IndexNow. The console would add indexing and error reports for a mostly Russian/CIS audience. Later, if ever [opinion].

**Naver Search Advisor.** It has its own index; the robot is `Yeti`.
- Registration is "host 단위만 지원" (host-level only). Ownership is checked by **meta tag or HTML file upload** (no DNS method is listed). Submit an XML sitemap under 요청 → 사이트맵 제출 and a full-text RSS under 요청 → RSS 제출. Content shows in search "최대 1주일 이내" (within a week) after the robot visits [P, [Naver guide](https://searchadvisor.naver.com/guide/seo-basic-intro), fetched raw].
- Its audience is Korean, and it already receives IndexNow. No [opinion].

**Seznam Webmaster.** It has its own index; the robot is `SeznamBot`.
- Log in with a Seznam account → "přidat web" (add site). Verify with `<meta name="seznam-wmt" …>` in `<head>` (not injected by script) or a `seznam-wmt-<code>.txt` file at the root, and keep it, because Seznam re-verifies. HTTP and HTTPS are verified separately. Data appears within about 24 h [P, [Seznam Webmaster](https://o-seznam.cz/napoveda/vyhledavani/seznam-webmaster/)].
- The sitemap is read from robots.txt (absolute URL) [P, snippet]. Czech audience, and it already receives IndexNow. No.

**Baidu Ziyuan (百度搜索资源平台).** It has its own index. I could not fetch the pages (`ziyuan.baidu.com` reset the connection), so all of the following is [P, snippet]:
- Verification is by file, HTML tag or CNAME. Submission ("普通收录") is by API push, sitemap or manual. API plus manual share a 100,000/day quota, and the sitemap quota depends on the site.
- A May 2021 notice cut sites **without ICP filing (备案)** to one sitemap submission per day. ICP filing requires a mainland-China entity and hosting [inferred, not verified here].
- Our robots.txt already disallows Baiduspider through Cloudflare's block. Chinese-language search. **No.**

**Yep (Ahrefs).** "AhrefsBot remains the primary source of data for both Yep and Ahrefs". YepBot is "used only to make requests according to IndexNow protocol" [P, snippet, [YepBot](https://yep.com/yepbot/)]. AhrefsBot gets 200 from us [observed]. Covered by IndexNow.

**Mojeek.** It has its own index; the robot is `MojeekBot`, which obeys the first `MojeekBot` group or else `*` [P, [MojeekBot](https://www.mojeek.com/bot.html)]. It has no add-URL service because of "spammy submissions" and does not use sitemaps for discovery. It asks for "prominent and genuine links" [P, snippets of Mojeek's blog/forum]. No action.

**Brave Search.** It has its own index: "Every result comes from Brave's own independent index" [P, [Brave Search API](https://brave.com/search/api/)]. Pages come from its crawler plus the opt-in Web Discovery Project, which needs a URL to be visited by many people before it is sent [P, snippet].
- The crawler "does not advertise a differentiated user agent", and "if a domain or page is not crawlable by Googlebot, then Brave Search's bot will not crawl it either". Delisting is by `noindex`.
- Brave has a form: "You can submit your page for re-fetching" at `search.brave.com/submit-url` [P, [Brave crawler help](https://search.brave.com/help/brave-search-crawler)].
- Later: one-off, for the home page and two or three top Explainers [opinion].

**Apple (Spotlight, Siri, Safari).** Applebot data powers "search in Spotlight, Siri, and Safari". Applebot-Extended only controls training. If Applebot is not named, it follows Googlebot's rules. No submission process is described; the only contact is applebot@apple.com [P, [Apple, 2026-09-04](https://support.apple.com/en-us/119829)]. No action: Applebot gets 200 [observed].

**Kagi.** Its own indexes are "Teclis" (web) and "TinyGem" (news), plus "anonymized API calls to all major search result providers" [P, [Kagi sources](https://help.kagi.com/kagi/search-details/search-sources.html)]. The only submission path is the Small Web list: English, "Only personal blogs… (no multi-author blogs)", "must not contain any forms of advertisements", and no "auto-generated, LLM-generated" content [P, [kagisearch/smallweb](https://github.com/kagisearch/smallweb)]. hastcode has AdSense and AI-drafted Pieces, so it is **not eligible**.

**Petal Search (Huawei)**: not researched; Cloudflare's block disallows PetalBot (§0).

---

## 5. News surfaces (changes since the 2026-10-20 note)

- **Google News / Top stories / Discover**: still automatic, nothing to submit. See [prior note](google-news-discover-requirements.md) §1–§3. One new link: the auto-generated publication page takes its logo from the site favicon (prior note §1), which ties to our `data:`-only favicon (§1.7).
- **Bing News PubHub**: still inside BWT per the page bundle; status unclear (§2). Apply only once the Briefs have a steady track record. Bing's newsworthiness criterion excludes how-to content (prior note §8), so only Briefs qualify [inferred].
- **Microsoft Start**: Bing says its index powers it [P], but I found no first-party publisher-application page (Open questions).
- **Apple News Publisher**: no change. Publishers must be based in AU/CA/UK/US, English or Canadian French only (prior note §8).
- **Flipboard**: Flipboard's own blog says add the RSS feed in the app (search for the feed URL → Social tab). The feed must be "the full feed, not just titles or extracts", each post needs an image of at least 400 px on its shortest side, the blog needs "at least 30 items", and help is at FlipBlogger@Flipboard.com [P, [Flipboard blog, ~2017](https://about.flipboard.com/inside-flipboard/rss-feed-on-flipboard/)]. Our `/feed.xml` is full-text with more than 30 items [observed]. The page is old, so treat it as Later.

---

## 6. AI answer engines: how each one finds content

No vendor offers a submission form. Each says, in effect, "let our crawler in". The robots tokens below are from each vendor's own doc [P]. Citation behaviour is in [`aeo-citation-signals.md`](aeo-citation-signals.md).

| Engine | Finds pages through | Token(s) and what blocking does | Submission | hastcode today |
|---|---|---|---|---|
| Google AI Overviews / AI Mode | Google Search index. "no additional requirements to appear"; the page must be indexed and snippet-eligible ([AI features](https://developers.google.com/search/docs/appearance/ai-features)) | `Googlebot` controls it; `Google-Extended` "doesn't govern Search itself" | None | Allowed |
| ChatGPT search | `OAI-SearchBot`, "for search"; Bing per Bing's claim | `GPTBot` (training) is independent; changes take "~24 hours" ([OpenAI bots](https://developers.openai.com/api/docs/bots)) | None; allow the bot and OpenAI's IP ranges (`openai.com/searchbot.json`) | OAI-SearchBot 200; GPTBot blocked (training) |
| Perplexity | `PerplexityBot`, "not used to crawl content for AI foundation models"; `Perplexity-User` "generally ignores robots.txt" ([Perplexity bots](https://docs.perplexity.ai/guides/bots)) | Changes take up to 24 h | None | **PerplexityBot 403 at the edge** |
| Claude | `Claude-SearchBot` (search), `Claude-User` (user fetches), `ClaudeBot` (training). Blocking ClaudeBot excludes content from training only. From 2026-10-30 Anthropic also honours an IETF-aipref-style training signal; a `ClaudeBot` rule takes priority ([Anthropic](https://support.claude.com/en/articles/8896518-does-anthropic-crawl-data-from-the-web-and-how-can-site-owners-block-the-crawler)) | Disabling Claude-SearchBot "may reduce your site's visibility" | None | Claude-SearchBot 200; ClaudeBot blocked (training) |
| Copilot | Bing index (prior AEO note §4) | `Bingbot` | Through BWT | Allowed |
| Amazon / Alexa | `Amzn-SearchBot` (search), `Amazonbot` (may train), `Amzn-User` | Changes take about 24 h | IndexNow | Amzn-SearchBot 200; Amazonbot blocked |
| DuckDuckGo Duck.ai / DuckAssist | `DuckAssistBot` (details not fetched) | — | None | 200 |

**Our robots.txt Content-Signal line.** `search=yes,ai-train=no,use=reference` comes from Cloudflare. `ai-input` is unset, which Cloudflare's own header text says "neither grants nor restricts" [observed]. Which vendors read Content-Signal is not documented on the pages I fetched. Anthropic's page says it honours a training-preference signal but does not name it (Open questions).

---

## 7. Open and archival indexes

- **Common Crawl**: no request form. It archives "a randomly selected subset" of a site and uses the robots.txt `Sitemap:` lines. Its token is `CCBot` [P, [Common Crawl FAQ](https://commoncrawl.org/faq)]. Blocked today by Cloudflare's block and at the edge [observed]. Its FAQ does not name search engines as users. Keeping it blocked costs us nothing in search [inferred]. That is an Author decision tied to the training opt-out.
- **Internet Archive**: Save Page Now: "Put a URL into the form, press the button, and we save the page". It saves one page, does not save outlinks, and keeps no IP address [P, [IA help](https://help.archive.org/help/save-pages-in-the-wayback-machine/)]. The Internet Archive is also in IndexNow's registry [P], so our CI pings probably reach it already [inferred]. Optional: save the home page and About page once.
- **WebSub**: Google says "If you use Atom or RSS, you can use WebSub to broadcast your changes to search engines, including Google" [P, [Build a sitemap](https://developers.google.com/search/docs/crawling-indexing/sitemaps/build-sitemap)]. Google runs a public hub at `pubsubhubbub.appspot.com` ("run by Google as a service"; page dated ©2022). Publishers add `<atom:link rel="hub" href="https://pubsubhubbub.appspot.com/"/>` and POST `hub.mode=publish&hub.url=<feed>` after each update [P, [hub](https://pubsubhubbub.appspot.com/)]. Our feeds have `rel="self"` but no `rel="hub"` [observed]. It would need a builder change plus a post-deploy POST, like the IndexNow step. Later, not urgent: IndexNow does not reach Google, and WebSub is the only push channel Google documents for feeds [inferred].
- **RSS readers (Feedly, Inoreader, etc.)**: I found no first-party "submit your feed" path; they find feeds when users subscribe [not verified]. The practical step is the missing `<link rel="alternate" type="application/rss+xml" href="/feed.xml">` in `<head>` (§0), which is how readers and browsers find a feed from a page URL [inferred].

---

## 8. Excluded: spammy submission services and link directories

- **"Submit your site to 500+ search engines" services**: no. Every engine in this note that has an index takes submissions only through its own console, its form or IndexNow (§1–§4). These services' "engines" are mostly link directories or Bing/Google resellers [inferred]. Google lists "Low-quality directory or bookmark site links" and "Using automated programs or services to create links to your site" as link spam [P, [Spam policies](https://developers.google.com/search/docs/essentials/spam-policies), quoted in the prior notes].
- **Paid or "SEO" web directories and backlink packages**: no, for the same reason.
- **"Instant indexing" tools built on the Google Indexing API**: no. They run outside the API's scope (§1.5).
- **Pinging the old Google/Bing sitemap ping URLs**: pointless; they return 404 (§1.3).
- **Mass Request indexing or IndexNow resubmission of unchanged URLs**: no. Google: repeats "won't get it crawled any faster". IndexNow: submit "only when content has changed", and 429 means "potential spam" (§1.4, §3).

---

## Runbook: get indexed, in order

Plain steps. **[by hand]** means a dashboard or DNS step only the Author can do. Steps with commands only read; none of them changes the site. Do not commit anything from this runbook: it is all dashboard work or read-only checks.

### Step 1. Let PerplexityBot in (10 minutes) [by hand, Cloudflare]

1. Open https://dash.cloudflare.com → your account → **hastcode.com** → **AI Crawl Control** → **Security** tab.
2. In the **Crawlers** table find **PerplexityBot**. In **Actions** choose **Allow**.
3. While you are there, decide each of these on purpose. They are training crawlers, so blocking them does not hurt search:
   - GPTBot, ClaudeBot, CCBot, Amazonbot, Google-Extended, Applebot-Extended: keep Blocked if you do not want model training.
   - Leave OAI-SearchBot, Claude-SearchBot, Perplexity-User and Amzn-SearchBot allowed.
4. Wait a few minutes, then check from your laptop. Expect `200`:

```sh
curl -s -o /dev/null -w "%{http_code}\n" \
  -A "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; PerplexityBot/1.0; +https://perplexity.ai/perplexitybot)" \
  https://hastcode.com/
```

5. Optional: in Cloudflare, check the robots.txt / Bot Preference Sync setting that adds the "Cloudflare Managed Content" block. It disallows Baiduspider and PetalBot, which are search engines. That is fine unless you ever want Chinese or Huawei search.

### Step 2. Redirect http to https (2 minutes) [by hand, Cloudflare]

1. Cloudflare → hastcode.com → **SSL/TLS** → **Edge Certificates** → turn on **Always Use HTTPS**.
2. Check. Expect a 301 (or 308) line, then `location: https://hastcode.com/`:

```sh
curl -sI http://hastcode.com/ | head -3
```

### Step 3. Google Search Console (15 minutes) [by hand, Google]

1. Open https://search.google.com/search-console. In the property picker, look for **hastcode.com** (a Domain property, no `https://`).
   - If it is there and you are Owner: go to 3.2.
   - If not: **Add property** → **Domain** → `hastcode.com` → choose **TXT**. Copy the `google-site-verification=…` value. In Cloudflare: **DNS** → **Records** → **Add record** → Type **TXT**, Name `@`, Content = the value → **Save**. Back in Search Console click **Verify**; if it fails, retry in a day. Never delete this record.
   - Check the record is live: `dig +short TXT hastcode.com`
2. **Indexing → Sitemaps**. Add `https://hastcode.com/sitemap.xml` and press **Submit**, then the same for `https://hastcode.com/news-sitemap.xml`. Both rows should say **Success**. Do not resubmit them later unless one says Couldn't fetch.
3. **URL Inspection**: paste `https://hastcode.com/` → **Request indexing**. Do the same for the newest 3–5 Pieces only. Do this once per URL; repeats do not help.
4. **Indexing → Pages**: note the counts of **Indexed**, **Discovered – currently not indexed** and **Crawled – currently not indexed**. For a 3-week-old site, "Discovered" is normal. Recheck weekly, not daily.
5. **Do not** use the Indexing API or any plugin or tool that does "instant indexing" through it.

### Step 4. Bing Webmaster Tools (10 minutes) [by hand, Microsoft]

1. Open https://www.bing.com/webmasters → sign in → **My Sites** → **Import** → sign in with the Google account that owns the Search Console property → select hastcode.com → Import. This verifies the site and copies the sitemaps.
2. **Sitemaps**: confirm both sitemaps are listed. If not, add them by URL.
3. **IndexNow** (left menu, "IndexNow Insights"): check that recent URLs show as submitted. Our CI pings every Production Publish. Data can take up to 48 h after import.
4. Nothing else to submit. Leave **Submit URLs** alone; IndexNow already does that job.

### Step 5. IndexNow: already automatic. Check, don't re-send [read-only]

After each `scripts/publish.sh cf`, the GitHub run's **Ping IndexNow** step prints `200` or `202`. Check the latest run:

```sh
gh run list --workflow publish.yml --limit 2
gh run view <run-id> --log | grep -E "changed URL\(s\)$|Z (200|202|400|403|422|429)$"
```

If it ever prints `403` (key problem), check the key file, which must print `465f05c90b9917457c1ad83b8df48051`:

```sh
curl -s https://hastcode.com/465f05c90b9917457c1ad83b8df48051.txt
```

**Manual send: only if a Publish's ping failed or was skipped.** Do not re-send unchanged pages. Single URL:

```sh
curl -s -w "\n%{http_code}\n" \
  "https://api.indexnow.org/indexnow?url=https%3A%2F%2Fhastcode.com%2Fbriefs%2Fgithub-availability-october-2026%2F&key=465f05c90b9917457c1ad83b8df48051"
```

Several URLs (up to 10,000):

```sh
curl -s -w "\n%{http_code}\n" -X POST https://api.indexnow.org/indexnow \
  -H 'Content-Type: application/json; charset=utf-8' \
  --data-binary @- <<'JSON'
{
  "host": "hastcode.com",
  "key": "465f05c90b9917457c1ad83b8df48051",
  "keyLocation": "https://hastcode.com/465f05c90b9917457c1ad83b8df48051.txt",
  "urlList": [
    "https://hastcode.com/briefs/github-availability-october-2026/",
    "https://hastcode.com/apis-networking/fly-io-6pn-routing/"
  ]
}
JSON
```

What the codes mean: `200`/`202` good · `403` bad key · `422` a URL is not on hastcode.com · `429` stop and wait at least 10 minutes.

### Step 6. Optional, once, in a spare 15 minutes [by hand]

- **Brave**: open https://search.brave.com/submit-url and submit `https://hastcode.com/` and two top Explainers.
- **Internet Archive**: open https://web.archive.org/save, paste `https://hastcode.com/` → **Save**. Repeat for `/about/`.
- **Yandex Webmaster** (only if you want its reports): https://webmaster.yandex.com → **Add site** → `https://hastcode.com` → choose **DNS** → add the TXT record in Cloudflare (Name `@`) → **Verify** → **Indexing → Sitemap files** → add `https://hastcode.com/sitemap.xml`.

### Step 7. Later (after about 2 months of steady Briefs) [by hand]

- Bing News PubHub inside BWT, if it is still offered (prior note §8).
- Flipboard: add `https://hastcode.com/feed.xml` in the app, or email FlipBlogger@Flipboard.com.

### Do not do

- "Submit to 500 search engines", SEO directories, backlink packages (§8).
- Google Indexing API "instant indexing" tools (§1.5).
- Baidu, Naver, Seznam consoles, Kagi Small Web, Apple News: wrong market or not eligible (§4, §5).
- Re-requesting indexing or re-pinging IndexNow for pages that did not change.

### Weekly check (5 minutes, read-only)

```sh
# robots.txt still has our allow groups and both sitemap lines
curl -s https://hastcode.com/robots.txt | grep -E "^(User-agent: (Googlebot|Bingbot|OAI-SearchBot|PerplexityBot)|Sitemap:)"
# search and answer crawlers still get 200 (403 means Cloudflare is blocking one)
for b in "Googlebot/2.1" "bingbot/2.0" "OAI-SearchBot/1.3" "PerplexityBot/1.0" "Claude-SearchBot/1.0" "Applebot/0.1"; do
  printf "%-22s " "$b"; curl -s -o /dev/null -w "%{http_code}\n" -A "Mozilla/5.0 (compatible; $b)" https://hastcode.com/
done
```

Then look at Search Console **Pages** and **Performance**, and Bing **IndexNow Insights**.

---

## Code-side follow-ups (for `.scratch/` tickets; not done here) [opinion]

1. **Favicon file for crawlers**: keep the `data:` icon for browsers (ADR-0006), and add a crawlable square PNG or ICO of at least 48 px (e.g. `/favicon.png` or `rel="apple-touch-icon"`). Measure the Lighthouse cost first; it is an ADR-0006 amendment.
2. **Feed autodiscovery**: add `<link rel="alternate" type="application/rss+xml" title="HastCode" href="/feed.xml">` to `<head>` (and the Vertical/Briefs feed on those pages).
3. **WebSub (optional)**: `atom:link rel="hub"` in the feeds plus a post-deploy `hub.mode=publish` POST beside "Ping IndexNow".
4. **README / workflow comment**: README "Feeds, sitemaps, robots and IndexNow" describes our robots.txt as "allow all". Production's file now also carries Cloudflare's managed block. Note that so nobody is surprised.

---

## Open questions

1. **Is the PerplexityBot 403 an AI Crawl Control Block, a WAF/custom rule, or "Block AI bots"?** The response body is "Your request was blocked." with no `cf-mitigated` header. Check AI Crawl Control → Security, and Security → Events filtered by user agent.
2. **Does Google use a `data:` URI favicon at all?** Google's favicon page neither allows nor forbids it, but it requires a crawlable "favicon file". Check after indexing: search `site:hastcode.com` on mobile and see whether results show our mark or a generic globe.
3. **Unlabelled TXT `78ace81ef8ec61ef2998c8c214526a80`** on the apex: which service added it (Bing? another verification)? Check Cloudflare DNS comments, or whatever tool you used around launch.
4. **Bing help pages** (URL Submission quota, Add and Verify site) are client-rendered. The 10,000/day figure and the midnight-UTC reset come from a 2019 blog and a search snippet. Read the live page in a browser.
5. **Bing API after 2026-08-31**: is `ssl.bing.com/webmaster/api.svc/json/*` the "REST" API that survives the SOAP/POX retirement? Only matters if we ever use the Bing API, which we should not need.
6. **Bing PubHub**: a search snippet mentioned a retirement banner on the PubHub support page; the page bundle still has PubHub forms. Check inside BWT.
7. **Which Amazon bot follows up IndexNow pings?** Amazon's Amazonbot page does not mention IndexNow, and we block `Amazonbot`.
8. **Content-Signal readers**: which vendors honour Cloudflare's `Content-Signal` line? Anthropic says it honours an aipref-style training signal from 2026-10-30 but does not name it.
9. **Startpage, Qwant, Baidu, Yep, Seznam (parts)** could only be read through snippets or not at all (blocked or connection reset). Facts marked [P, snippet] need a browser check before anyone relies on them.
10. **Microsoft Start publisher programme**: no first-party application page was found.
11. **Whether the Google WebSub hub is still maintained** (page says ©2022).
12. **Crawler Hints on Workers Static Assets** (carried over from the prior note): still undocumented. Low priority, since our own IndexNow works.

---

## Sources (accessed 2026-11-16)

**Google**
- https://support.google.com/webmasters/answer/34592 (property types)
- https://support.google.com/webmasters/answer/9008080 (verification)
- https://support.google.com/webmasters/answer/7451001 (Sitemaps report)
- https://support.google.com/webmasters/answer/7440203 (Page indexing report)
- https://developers.google.com/search/docs/crawling-indexing/ask-google-to-recrawl
- https://developers.google.com/search/apis/indexing-api/v3/quickstart ; https://developers.google.com/search/apis/indexing-api/v3/using-api ; https://developers.google.com/search/apis/indexing-api/v3/quota-pricing
- https://developers.google.com/search/docs/crawling-indexing/sitemaps/build-sitemap
- https://developers.google.com/search/blog/2023/06/sitemaps-lastmod-ping
- https://developers.google.com/search/docs/appearance/favicon-in-search
- https://developers.google.com/search/docs/crawling-indexing/consolidate-duplicate-urls
- https://developers.google.com/search/docs/crawling-indexing/google-common-crawlers
- https://developers.google.com/search/docs/appearance/ai-features
- https://developers.google.com/search/docs/essentials/spam-policies (via prior notes)
- https://pubsubhubbub.appspot.com/

**Bing / Microsoft**
- https://blogs.bing.com/webmaster/2025/6/Start-Using-Bing-Webmaster-Tools-to-Improve-Your-Site-Visibility/
- https://blogs.bing.com/webmaster/2019/1/bingbot-Series-Get-your-content-indexed-fast-by-now-submitting-up-to-10%2C000-URLs-per-day-to-Bing/
- https://blogs.bing.com/webmaster/2019/6/bingbot-Series-Introducing-Batch-mode-for-Adaptive-URL-submission-API/
- https://blogs.bing.com/webmaster/March-2024/Optimize-your-Impact-with-IndexNow-Insights
- https://learn.microsoft.com/en-us/bingwebmaster/ ; https://learn.microsoft.com/en-us/bingwebmaster/getting-access
- https://learn.microsoft.com/en-us/dotnet/api/microsoft.bing.webmaster.api.interfaces.iwebmasterapi.submitcontent?view=bing-webmaster-dotnet
- https://learn.microsoft.com/en-us/dotnet/api/microsoft.bing.webmaster.api.interfaces.iwebmasterapi.geturlsubmissionquota?view=bing-webmaster-dotnet (search result only)
- Client-rendered, text not retrievable: https://www.bing.com/webmasters/help/URL-Submission-62f2860b ; https://www.bing.com/webmasters/help/add-and-verify-site-12184f8b ; https://www.bing.com/webmasters/help/bing-news-pubhub-support-e85a0280 (page bundle strings read)

**IndexNow**
- https://www.indexnow.org/documentation ; https://www.indexnow.org/faq ; https://www.indexnow.org/searchengines.json
- https://www.bing.com/indexnow/meta.json ; https://web-static.archive.org/indexnow/meta.json

**Other engines**
- Yandex: https://yandex.com/support/webmaster/en/service/rights ; https://yandex.com/support/webmaster/en/robot-workings/new-changed ; https://yandex.com/support/webmaster/en/robot-workings/site-reindex ; snippet: https://yandex.com/support/webmaster/en/indexing-options/index-now
- Naver: https://searchadvisor.naver.com/guide/seo-basic-intro (fetched raw with curl; WebFetch blocked)
- Seznam: https://o-seznam.cz/napoveda/vyhledavani/seznam-webmaster/ ; snippets: https://o-seznam.cz/napoveda/vyhledavani/seznambot/protokol-indexnow/ , https://o-seznam.cz/napoveda/vyhledavani/seznambot/sitemap-xml/
- Baidu (snippets only; connection reset): https://ziyuan.baidu.com/college/courseinfo?id=267&page=2 ; https://ziyuan.baidu.com/college/courseinfo?id=267&page=17 ; https://ziyuan.baidu.com/college/articleinfo?id=3134
- Yep: snippet https://yep.com/yepbot/
- Mojeek: https://www.mojeek.com/bot.html ; snippets: https://blog.mojeek.com/2015/03/how-to-submit-your-site-to-mojeek.html , https://community.mojeek.com/t/how-can-i-add-my-website-to-mojeek-search-results/1192
- Brave: https://search.brave.com/help/brave-search-crawler ; https://brave.com/search/api/
- Apple: https://support.apple.com/en-us/119829
- DuckDuckGo: https://duckduckgo.com/duckduckgo-help-pages/results/sources ; https://duckduckgo.com/duckduckgo-help-pages/results/duckduckbot
- Ecosia: https://support.ecosia.org/article/579-search-results-providers
- Qwant (snippets): https://betterweb.qwant.com/en/2024/11/08/ecosia-and-qwant-join-forces-to-develop-european-search-index/ ; https://about.qwant.com/legal/rapport-de-transparence-qwant-2025/
- Kagi: https://help.kagi.com/kagi/search-details/search-sources.html ; https://github.com/kagisearch/smallweb
- Startpage: not fetchable (domain blocked to this research tool)

**AI crawlers**
- https://developers.openai.com/api/docs/bots
- https://docs.perplexity.ai/guides/bots
- https://support.claude.com/en/articles/8896518-does-anthropic-crawl-data-from-the-web-and-how-can-site-owners-block-the-crawler
- https://developer.amazon.com/amazonbot

**Cloudflare**
- https://developers.cloudflare.com/bots/additional-configurations/managed-robots-txt/ (updated 2026-11-02)
- https://developers.cloudflare.com/ai-crawl-control/ ; https://developers.cloudflare.com/ai-crawl-control/features/manage-ai-crawlers/
- https://developers.cloudflare.com/bots/concepts/bot/verified-bots/ (updated 2026-11-12)
- https://developers.cloudflare.com/ssl/edge-certificates/additional-options/always-use-https/
- https://developers.cloudflare.com/dns/manage-dns-records/how-to/create-dns-records/

**Archives, feeds, news**
- https://commoncrawl.org/faq
- https://help.archive.org/help/save-pages-in-the-wayback-machine/
- https://about.flipboard.com/inside-flipboard/rss-feed-on-flipboard/
- https://www.rfc-editor.org/rfc/rfc9309.html (robots.txt group merging)

**Observed (2026-11-16)**
- `curl` of `https://hastcode.com/` and `{robots.txt, sitemap.xml, news-sitemap.xml, feed.xml, briefs/feed.xml, llms.txt, sitemap.md, 465f05c90b9917457c1ad83b8df48051.txt, ads.txt, favicon.ico, mark.svg}`; `http://hastcode.com/`; `https://www.hastcode.com/`.
- Crawler user-agent GETs (from a non-crawler IP) against `/` and `/distributed-systems/quepaxa-consensus/`.
- `dig +short TXT hastcode.com`.
- `gh variable list`; `gh run view` 40943579139, 41054099392, 43691945311.
- Repo: `builder/step_write_{robots,sitemap,news_sitemap,markdown_sitemap,indexnow,llms_txt,feeds,cloudflare_manifest,vercel_config}.go`, `builder/step_validate_indexnow_key.go`, `builder/urls.go`, `internal/views/layout_head.templ`, `.github/workflows/publish.yml`, `scripts/setup-deploy.sh`, `site.yaml`, README "Feeds, sitemaps, robots and IndexNow", ADR-0006.

# How does a new developer-content domain get its first real traffic, and how do we measure it honestly?

Research for hastcode.com: a static Go-built publication on Cloudflare Workers Static Assets, about 40 Pieces live since 2026-10-25 (21 of them since 2026-10-28 12:54 UTC), with GA4 showing roughly zero visits after about 1.5 days. Goals are Raptive eligibility (25k PV/month, 50% Tier-1; see `docs/research/ad-rpm-tech-audience.md` §4) and newsletter sponsorship. Google News and Discover inclusion is covered in [`google-news-discover-requirements.md`](google-news-discover-requirements.md), and AEO/answer-engine signals in [`aeo-citation-signals.md`](aeo-citation-signals.md). This note links to them and does not repeat them.

Researched **2026-10-30**. Labels: **[P]** primary (the platform's own docs, policy or spec page, fetched on that date) / **[observed]** what a live request, DNS lookup or our CI log returned / **[inferred]** my deduction / **[opinion]**; "[secondary, unverified]" marks the one non-primary claim mentioned. Pages I could not fetch are named where they come up and listed under Open questions.

---

## TL;DR

- **IndexNow has never been pinged.** The 2026-10-28 Production publish built a 48-URL IndexNow artefact, but `Ping IndexNow` was **skipped** (`INDEXNOW_ENABLED=false`). The run then saved the page manifest anyway, so no later build will list those URLs again [observed]. Bing and the other IndexNow engines have not been told about the launch.
- **Near-zero traffic after 1.5 days is what Google's docs predict.** "Crawling can take anywhere from a few days to a few weeks". "the most common reason that a site is not indexed is because it's just too new". "the vast majority of the new pages Google finds every day are through links" [P]. None of the Google pages fetched mentions a "sandbox".
- **Sitemaps and Request indexing are hints, not levers.** A sitemap "is merely a hint". Request indexing has an unpublished "daily limit", and repeating it "won't get it crawled any faster" [P].
- **Bing is one cheap step.** Import the site and its sitemaps from Search Console. Bing says it powers Yahoo, DuckDuckGo and "AI tools such as ChatGPT". IndexNow submissions are "shared across all IndexNow-enabled search engines" [P].
- **Each channel's own rules limit self-promotion.** HN: "It's ok to post your own stuff part of the time", and Show HN excludes "blog posts". Lobsters: new users cannot submit a domain nobody has submitted before or post self-promo. DEV's AI guidelines say AI-assisted posts should not "Promote any business … (including your own)" [P]. I could not fetch Reddit's rules (403).
- **GA4 undercounts EEA/UK/CH readers, and that includes Tier-1 UK.** Consent defaults to denied there [observed]. Until the property sends ≥1,000 denied events/day and has ≥1,000 consented users/day, "events triggered by users who decline consent aren't reported" [P]. GA4 will therefore understate the Tier-1 share Raptive checks [inferred].
- **Cloudflare zone analytics is the server-side cross-check.** It works from edge logs with no script, "cannot be blocked", and includes crawlers [P]. I found no first-party figure for how much ad blockers undercount.
- **www 522**: www is proxied by Cloudflare, but nothing behind it answers [observed]. Cloudflare's documented fix: a proxied `AAAA 100::` record plus a 301 redirect rule [P].

---

## 1. Indexing: how new pages get found, and how long it takes

### 1.1 Google

- **Discovery is mostly through links.** "Google primarily finds pages through links from other pages it already crawled". "the vast majority of the new pages Google finds every day are through links" [P, [SEO Starter Guide](https://developers.google.com/search/docs/fundamentals/seo-starter-guide), updated 2025-12-10]. The other route: "Still other pages are discovered when you submit a list of pages" (a sitemap) [P, [How Search Works](https://developers.google.com/search/docs/fundamentals/how-search-works), 2025-12-18].
- **No guarantees and no paid speed-up.** "Google doesn't guarantee that it will crawl, index, or serve your page"; "Google doesn't accept payment to crawl a site more frequently" [P, How Search Works]. "we cannot make predictions or guarantees about when or if your URLs will be crawled or indexed" [P, [Crawling and indexing FAQ](https://developers.google.com/search/help/crawling-index-faq), 2025-12-10].
- **Timing, in Google's words.** "Crawling can take anywhere from a few days to a few weeks" [P, [Ask Google to recrawl](https://developers.google.com/search/docs/crawling-indexing/ask-google-to-recrawl), 2025-12-10]. For a single inspected URL: "Indexing typically takes only a day or so, but can take much longer in some cases" and "Indexing can take up to a week or two" [P, [URL Inspection tool help](https://support.google.com/webmasters/answer/9012289)]. For changes in general: "Some changes might take effect in a few hours, others could take several months". Google suggests you "wait a few weeks to assess whether your work had beneficial effects" [P, SEO Starter Guide].
- **New sites.** "the most common reason that a site is not indexed is because it's just too new; be patient and ask Google to crawl and index it" [P, Crawling FAQ]. **"Sandbox"**: none of the Google pages fetched (including the [Site position FAQ](https://developers.google.com/search/help/site-position-in-search-faq)) uses the word, and a site-restricted search of Google's docs found nothing. SEO blogs attribute denials to Google staff [secondary, unverified].
- **Sitemaps.** Submit through the Search Console Sitemaps report or a robots.txt line, but "submitting a sitemap is merely a hint" [P, [Build a sitemap](https://developers.google.com/search/docs/crawling-indexing/sitemaps/build-sitemap), 2026-07-08]. `lastmod` is used "if it's consistently and verifiably" accurate [P, same]. A sitemap "does not guarantee indexing or increase your site's ranking" [P, Crawling FAQ]. hastcode.com lists `/sitemap.xml` (54 `<loc>`) and `/news-sitemap.xml` (3 `<loc>`) in robots.txt [observed].
- **Request indexing.** "There is a daily limit to how many index requests you can submit". Google publishes no number. "Submitting a request does not guarantee that the page will appear in the Google Index". You need to be an owner or full user of the property [P, URL Inspection help; Ask Google to recrawl]. "requesting a recrawl multiple times for the same URL won't get it crawled any faster" [P, Ask Google to recrawl].
- **Google News / Discover** have their own eligibility (a "consistent history" requirement, no submissions). See [`google-news-discover-requirements.md`](google-news-discover-requirements.md) §1 and §9. Expect nothing from them in week one [opinion, as in that doc].

### 1.2 Bing, IndexNow and the engines that use Bing's index

- **Bing Webmaster Guidelines** ([page](https://www.bing.com/webmasters/help/webmaster-guidelines-30fba23a)) are client-rendered, and the fetch returned UI strings only. What follows comes from Bing's own blog.
- **Import from Search Console**: "Navigate to My Sites page on Bing Webmaster Tools and click Import", then sign in with Google. The import connects "verified sites and sitemaps". "it might take up to 48 hours to get traffic data" [P, [Bing blog 2019-09-13, updated 2025-06](https://blogs.bing.com/webmaster/2019/9/Import-sites-from-Search-Console-to-Bing-Webmaster-Tools/)].
- **Sitemaps**: Bing "processes sitemaps at least once every 24 hours" [P, [Bing blog 2025-06-17](https://blogs.bing.com/webmaster/2025/6/Start-Using-Bing-Webmaster-Tools-to-Improve-Your-Site-Visibility/)].
- **Who uses Bing's index.** In the same post Bing says it powers "a wide range of AI and search experiences", including Yahoo, DuckDuckGo and "AI tools such as ChatGPT" [P, Bing's own claim]. DuckDuckGo says its traditional links come from Bing, "which we largely source from Bing", alongside its own DuckDuckBot [P, [DuckDuckGo help](https://duckduckgo.com/duckduckgo-help-pages/results/sources)]. Copilot citations are covered in [`aeo-citation-signals.md`](aeo-citation-signals.md) §4.
- **IndexNow** (spec: [documentation](https://www.indexnow.org/documentation), [FAQ](https://www.indexnow.org/faq)) [P]:
  - Participating engines include Amazon, Bing, Naver, Seznam.cz, Yandex and Yep, and a submission "will be shared across all IndexNow-enabled search engines". Google is not listed.
  - "Submitting a URL does not guarantee immediate indexing."
  - It is "not designed for submitting every URL on your site at once", except after a site-wide change such as a migration or redesign. Use it alongside an XML sitemap.
  - Each POST takes up to 10,000 URLs. Responses: 200 OK; 202 key validation pending; 403 key not valid; 422 URL not on the host; 429 rate-limited, "may be flagged as spam". Exact rate limits are not published.
- **Cloudflare Crawler Hints** "uses Cloudflare cache signals to tell search engines when your content has likely changed". It supports IndexNow, works on all plans including Free, and is switched on under Caching → Configuration. It keys off "cache-status MISS", it is "a global setting for your entire website", and its docs say nothing about Workers [P, [Crawler Hints](https://developers.cloudflare.com/cache/advanced-configuration/crawler-hints/), 2026-08-14]. The IndexNow FAQ confirms "Cloudflare offers native IndexNow integration" [P]. Whether Workers Static Assets responses produce the cache signals Crawler Hints needs is not documented [inferred → Open questions]. hastcode already has its own IndexNow implementation, with the key file live at `/465f05c90b9917457c1ad83b8df48051.txt` (HTTP 200) [observed], so Crawler Hints is not needed.
- **Our IndexNow ping has never run** [observed]. In `.github/workflows/publish.yml`, `Ping IndexNow` runs only when `vars.INDEXNOW_ENABLED == 'true'`, and the variable is `false` (set 2026-10-25). Run 40691800972 logged "48 changed URL(s)", skipped the ping, and saved the `hastcode-indexnow-*` manifest cache. The workflow's guard ("A failed ping fails the job before the cache is saved") does not cover a **skipped** ping, so those URLs now count as announced [inferred].

---

## 2. Distribution channels: their rules for posting your own articles

| Channel | What its own rules say [P] | What it means for hastcode [inferred/opinion] |
|---|---|---|
| **Hacker News** ([guidelines](https://news.ycombinator.com/newsguidelines.html), [FAQ](https://news.ycombinator.com/newsfaq.html), [Show HN](https://news.ycombinator.com/showhn.html)) | "Please don't use HN primarily for promotion." "It's ok to post your own stuff part of the time, but the primary use of the site should be for curiosity." "use the original title, unless it is misleading or linkbait; don't editorialize." "Don't solicit upvotes, comments, or submissions"; the FAQ says sites that break this get penalised or banned. "Please don't delete and repost." Reposts are OK if a story "has not had significant attention in the last year or so". Show HN is for "something you've made that other people can play with"; "Off topic: blog posts, sign-up pages, newsletters, lists, and other reading material." Also: "Please don't put generated text in HN posts." | Post Explainers as **regular** submissions, never Show HN. Use one account with real non-promotional history, and post the Piece's own title. The "generated text" rule is about what you write on HN (comments, post text), so the submitter writes their own comments [inferred]. |
| **Lobsters** ([about](https://lobste.rs/about)) | Invite-only: "Invitations are used as a mechanism for spam-control". "self-promo should be less than a quarter of one's stories and comments in the last year". Users count as new for their first 70 days, and new users cannot "submit links to domains we haven't seen submitted before" and are "not allowed to submit self-promo". The focus is "pretty narrowly on computing". | Not a channel for this week. hastcode.com is an unseen domain, so a new account cannot submit it, and you need an invite first. A long-standing member submitting a strong Piece on its merits is the only clean route [opinion]. |
| **Reddit** (r/programming, r/golang, r/ExperiencedDevs) | **Not fetched.** Reddit's spam help article returned 403, redditinc.com is blocked to this fetcher, and `/r/<sub>/about/rules.json` returned 403 for all three subreddits. No claim is made here. | Read each subreddit's sidebar rules by hand before posting. Note what they say about blog links, self-promotion ratios and AI content in Open questions. |
| **DEV (dev.to)** ([editor guide](https://dev.to/p/editor_guide), [Forem API](https://developers.forem.com/api/v1), [AI guidelines](https://dev.to/guidelines-for-ai-assisted-articles-on-dev), last updated 2024-04-08) | Front matter `canonical_url`: "link for the canonical version of the content". API: "If this post was originally published elsewhere, specify the canonical URL to maintain SEO integrity." AI-assisted posts should "Disclose the fact that they were generated or assisted by AI … using the tag #ABotWroteThis", be "checked for factual accuracy before publishing", and should not "Promote any business, program, or course (including your own)." | Cross-posting is technically supported. But our Pieces are AI-drafted, and DEV's AI rule against promoting "your own" may make cross-posts with a link back a problem. Ask DEV before doing it (Open questions). |
| **Hashnode** ([Hashnode CLI post, 2026-10-06](https://hashnode.com/blog/hashnode-cli); [Townhall republishing post](https://townhall.hashnode.com/why-you-should-republish-your-devblog-posts-and-how-to-do-it), undated, ~2019) | "Add a subtitle, a series, or a canonical URL with flags such as `--subtitle`, `--series-id`, and `--original-article-url`." The web editor lets you "add the original URL". | Supported. I found no Hashnode policy on AI content or self-promotion. |
| **Medium** ([canonical help](https://help.medium.com/hc/en-us/articles/360033930293-Set-a-canonical-link), [Medium Rules](https://medium.com/policy/medium-rules-30e5502c4eb4)) | **Both returned 403. Not verified.** | No claim made. |

**Google on syndication** [P, [Canonicalization troubleshooting](https://developers.google.com/search/docs/crawling-indexing/canonicalization-troubleshooting), 2026-08-21]: "The canonical link element is not recommended for those who want to avoid duplication by syndication partners … because the pages are often very different. The most effective solution is for partners to block indexing of your content." DEV, Hashnode and Medium offer canonical links but not noindex [inferred from the pages above]. Google may therefore rank a full cross-post above the original. A short excerpt with a link back avoids that risk [opinion].

**Keep promotion clean: Google's link-spam policy** [P, [Spam policies](https://developers.google.com/search/docs/essentials/spam-policies), 2026-08-28]: "Link spam is the practice of creating links to or from a site primarily for the purpose of manipulating search rankings". Its examples include buying or selling links, "Excessive link exchanges", automated link creation, low-quality directory or bookmark links, forum comments with optimised links, and guest posts or press releases with ranking-credit links. Paid or sponsored links are fine "as long as they are qualified with a" `rel="nofollow"` or `rel="sponsored"`. That matters for future newsletter sponsorship links [inferred]. Google's own promotion advice: "Social media promotion / Community engagement / Advertisement … Word of mouth". It also warns "you can overdo promoting your site and actually harm it: … search engines may perceive some of the practices as manipulation of search results" [P, SEO Starter Guide].

---

## 3. Measurement

### 3.1 What GA4 counts on hastcode.com

- **Our setup** [observed in the live HTML and `internal/views/analytics.go`]: gtag loads on every page. `gtag("consent","default",{analytics_storage:"denied",wait_for_update:500,region:[…EEA…,"GB","CH"]})` is followed by a global `analytics_storage:"granted"`. That is **advanced** consent mode: "Google tags load when a user opens the website" [P, [Consent mode overview](https://developers.google.com/tag-platform/security/concepts/consent-mode), 2026-07-30].
- **When storage is denied**, "cookieless pings are sent to Google Analytics" and "No Analytics cookies are set, accessed, or read". "Google Analytics will use cookieless pings for modeling" [P, [GA4 consent mode](https://support.google.com/analytics/answer/9976101)].
- **Modelling thresholds** [P, [Behavioral modeling](https://support.google.com/analytics/answer/11161109)]: "at least 1,000 events per day with analytics_storage='denied' for at least 7 days" **and** "at least 1,000 daily users sending events with analytics_storage='granted'" on 7 of the previous 28 days. Until then, "events triggered by users who decline consent aren't reported". We are far below both thresholds [inferred]. So EEA/UK/CH readers who do not accept the AdSense consent message are **absent from GA4**. Raptive counts the UK as Tier-1 (`ad-rpm-tech-audience.md` §4), so GA4 will understate the Tier-1 share [inferred].
- **Processing lag**: "Data processing can take 24-48 hours". Realtime takes "a few minutes" [P, [Data freshness](https://support.google.com/analytics/answer/11198161)]. Standard reports for 2026-10-28/29 may still be filling in [inferred].
- **Verifying the tag**: Realtime shows "user activity during the last 30 minutes", and one of its listed uses is to "Verify that the measurement code is working" [P, [Realtime](https://support.google.com/analytics/answer/9271392)]. For DebugView, use Tag Assistant or `debug_mode:true`. Setting it to `false` does not turn it off. Filter developer traffic so it "doesn't negatively impact the data in your reports" [P, [DebugView](https://support.google.com/analytics/answer/7201382)]. Neither page covers consent-denied pings, so test from a non-EEA connection [inferred].
- **Ad and tracker blockers**: I found **no first-party Google statement or figure** on GA4 loss to blockers. Cloudflare says of its own beacon: "Cloudflare is aware that the analytics beacon is blocked by these services" [P, [Web Analytics FAQ](https://developers.cloudflare.com/web-analytics/faq/)]. Any percentage you see quoted elsewhere is a secondary claim and is unverified.

### 3.2 Cloudflare as the server-side cross-check

- **Zone analytics** (free plan: Requests, Bandwidth, Unique Visitors, countries) is "based on Cloudflare's edge logs, with no need for third party scripts or trackers". "A page view is defined as a successful HTTP response with a content-type of HTML". The figures "include legitimate user requests as well as crawlers and threats". The page does not define "Unique Visitors" [P, [Zone analytics](https://developers.cloudflare.com/analytics/account-and-zone-analytics/zone-analytics/), 2026-10-02]. Treat it as an **upper bound**: blockers cannot hide from it, but bots inflate it [inferred].
- **Web Analytics (RUM)** is "free" and uses "a JavaScript snippet (also known as a beacon)" that "does not collect or use your visitors' personal data" [P, [About](https://developers.cloudflare.com/web-analytics/about/), 2026-04-17]. On proxied sites it is injected automatically unless the origin sends `no-transform` [P, [Get started](https://developers.cloudflare.com/web-analytics/get-started/)]. Our apex sends `cache-control: public, max-age=0, must-revalidate` [observed]. Unsampled data is kept for 7 days, then aggregated to about 10% [P, FAQ]. Ad blockers block it too, so it is a second client-side view, not an independent check [inferred].
- **Workers metrics** count "All incoming requests registered by a Worker" [P, [Workers metrics](https://developers.cloudflare.com/workers/observability/metrics-and-analytics/), 2026-07-01]. "Requests to static assets are free and unlimited" [P, [Static assets billing](https://developers.cloudflare.com/workers/static-assets/billing-and-limitations/)]. Neither page says whether asset-only requests are counted, so use zone analytics for page views [inferred].
- **Honest weekly view** [opinion]: report three numbers side by side and never add them together:
  1. Search Console clicks and impressions (what search actually sent).
  2. GA4 page views (consented, script-dependent; this is the number Raptive will see).
  3. Cloudflare zone HTML page views (includes bots; an upper bound).

---

## 4. www.hastcode.com returns 522

- **Observed**: `www.hastcode.com` resolves to Cloudflare anycast IPs (172.67.223.186, 104.21.86.187), the same as the apex, and answers `HTTP/1.1 522` in about 0.6 s. The Worker config attaches assets to the `hastcode` Worker with no www route [observed: `dig`, `curl`, `.vercel/output/wrangler.json`].
- **What 522 means**: "Error `522` occurs when Cloudflare times out contacting the origin web server". It is also returned when a target "is an A record pointing to a reserved address, such as `100::` or `192.0.2.0`" in an Origin Rules context [P, [Error 522](https://developers.cloudflare.com/support/troubleshooting/http-status-codes/cloudflare-5xx-errors/error-522/), 2026-07-31]. So www has a proxied DNS record, but no Worker or redirect handles it, and the origin it points at does not answer [inferred].
- **The documented fix** [P, [Workers Custom Domains](https://developers.cloudflare.com/workers/configuration/routing/custom-domains/), 2026-09-29]: Custom Domains need an exact hostname match, so a Worker on `example.com` does not receive `www.example.com`. To redirect, add "a proxied DNS record" with "an `A` record to `192.0.2.0` or an `AAAA` record to `100::`" and use a redirect rule. "You cannot create a Custom Domain on a hostname with an existing CNAME DNS record."
- **The redirect rule** [P, [Redirect www to root](https://developers.cloudflare.com/rules/url-forwarding/examples/redirect-www-to-root/), 2026-05-05]: a wildcard on Request URL `https://www.*`, target `https://${1}`, status 301, "Preserve query string". That pattern matches only HTTPS. Plain `http://www` relies on the zone's HTTPS upgrade [inferred].

---

## Open questions

1. **Reddit**: Reddit's spam/self-promotion help and the rules of r/programming, r/golang and r/ExperiencedDevs could not be fetched (403 / blocked). Read them by hand in a browser, and record what each says about blog links, self-promo ratio and AI content.
2. **DEV's AI rule**: does "should not … Promote any business … (including your own)" stop a canonical cross-post of an AI-drafted, human-verified Piece? Also check whether the 2024-04-08 guideline has been revised. Ask DEV before cross-posting.
3. **Medium**'s canonical-link help and Medium Rules returned 403. Not verified.
4. **Bing Webmaster Guidelines** text is client-rendered and was not retrievable. Its current wording on discovery, URL submission quotas and AI-generated content is unverified here.
5. **Crawler Hints on Workers Static Assets**: does a response served from Workers assets produce the "cache-status MISS" signal Crawler Hints relies on? Not documented. Low priority while our own IndexNow ping exists.
6. **Do cookieless (denied) pings show in GA4 Realtime or DebugView?** Not stated on either help page.
7. **Unique Visitors** in Cloudflare zone analytics is not defined on the page fetched.
8. **Workflow**: should a skipped IndexNow ping still save the manifest cache? As written, it silently drops changed URLs. This is a ticket for `.scratch/`.

---

## What to do this week (ordered) [opinion]

1. **Fix www.** Remove whatever proxied record now sits on `www` (§4 observed). Add a proxied `AAAA www 100::` and a 301 Single Redirect `https://www.*` → `https://${1}` with the query string preserved. Then confirm with `curl -sI https://www.hastcode.com/` (expect 301). *(Cloudflare Custom Domains + redirect example, §4.)*
2. **Turn IndexNow on and announce the launch once.** Set `INDEXNOW_ENABLED=true`. Then make one manual POST of the 54 sitemap URLs to `https://api.indexnow.org/indexnow` and record the status code (200/202 is good; 403 means a key problem). The 48 URLs from 2026-10-28 were never pinged and will not be re-listed. A first launch is the site-wide-change case the IndexNow FAQ allows. *(§1.2 observed; IndexNow FAQ and spec.)*
3. **Search Console:** submit `/sitemap.xml` and `/news-sitemap.xml` in the Sitemaps report. Use URL Inspection → Request indexing once each for the home page and your 5-10 strongest Pieces, and do not repeat it. *(Google: sitemap "merely a hint"; daily limit; repeats don't help, §1.1.)*
4. **Bing Webmaster Tools:** import hastcode.com from Search Console, which brings its sitemaps along. Expect traffic data within 48 h. *(Bing import post, §1.2; feeds DuckDuckGo/Yahoo per their own pages.)*
5. **Prove GA4 works before judging traffic.** Check Realtime and DebugView from a non-EEA connection and add a developer-traffic filter. Do not read the standard reports for 2026-10-28/29 until 48 h have passed. *(GA4 Realtime, DebugView, Data freshness, §3.1.)*
6. **Start the three-number weekly log** (Search Console / GA4 / Cloudflare zone HTML page views). Note in it that GA4 omits non-consenting EEA/UK/CH readers. *(Behavioral modeling thresholds; zone analytics, §3.)*
7. **One human distribution step.** Submit the single strongest Explainer to HN as a regular (not Show HN) submission, with its original title, from an account with real history. Write your own comments and ask no one to upvote. Hold Lobsters (you need an invite), Reddit (read the rules by hand first) and DEV (Open question 2). *(HN guidelines/FAQ/Show HN, Lobsters about, DEV AI guidelines, §2.)*
8. **Do not buy or swap links.** Mark any future sponsorship link `rel="sponsored"`. *(Google link-spam policy, §2.)*
9. **Judge in weeks, not days.** Re-check indexing coverage in Search Console around 2026-11-13. *(Google: "a few days to a few weeks"; "wait a few weeks to assess", §1.1.)*

---

## Sources (accessed 2026-10-30)

**Google Search**
- https://developers.google.com/search/docs/fundamentals/seo-starter-guide
- https://developers.google.com/search/docs/fundamentals/how-search-works
- https://developers.google.com/search/docs/crawling-indexing/ask-google-to-recrawl
- https://support.google.com/webmasters/answer/9012289 (URL Inspection tool)
- https://developers.google.com/search/help/crawling-index-faq
- https://developers.google.com/search/help/site-position-in-search-faq
- https://developers.google.com/search/docs/crawling-indexing/sitemaps/build-sitemap
- https://developers.google.com/search/docs/crawling-indexing/canonicalization-troubleshooting
- https://developers.google.com/search/docs/crawling-indexing/consolidate-duplicate-urls (checked; has no syndication guidance)
- https://developers.google.com/search/docs/essentials/spam-policies

**Bing, IndexNow, DuckDuckGo**
- https://blogs.bing.com/webmaster/2019/9/Import-sites-from-Search-Console-to-Bing-Webmaster-Tools/
- https://blogs.bing.com/webmaster/2025/6/Start-Using-Bing-Webmaster-Tools-to-Improve-Your-Site-Visibility/
- https://www.bing.com/indexnow
- https://www.indexnow.org/documentation ; https://www.indexnow.org/faq
- https://duckduckgo.com/duckduckgo-help-pages/results/sources
- https://www.bing.com/webmasters/help/webmaster-guidelines-30fba23a (client-rendered; text not retrievable)

**Distribution channels**
- https://news.ycombinator.com/newsguidelines.html ; https://news.ycombinator.com/newsfaq.html ; https://news.ycombinator.com/showhn.html
- https://lobste.rs/about
- https://dev.to/p/editor_guide ; https://developers.forem.com/api/v1 ; https://dev.to/guidelines-for-ai-assisted-articles-on-dev
- https://hashnode.com/blog/hashnode-cli ; https://townhall.hashnode.com/why-you-should-republish-your-devblog-posts-and-how-to-do-it
- Not fetched: https://help.medium.com/hc/en-us/articles/360033930293-Set-a-canonical-link (403), https://medium.com/policy/medium-rules-30e5502c4eb4 (403), https://support.reddithelp.com/hc/en-us/articles/360043504051-Spam (403), https://redditinc.com/policies/reddit-rules (blocked), `https://www.reddit.com/r/{programming,golang,ExperiencedDevs}/about/rules.json` (403)

**Google Analytics / tag platform**
- https://developers.google.com/tag-platform/security/concepts/consent-mode
- https://support.google.com/analytics/answer/9976101 (consent mode)
- https://support.google.com/analytics/answer/11161109 (behavioral modeling)
- https://support.google.com/analytics/answer/11198161 (data freshness)
- https://support.google.com/analytics/answer/9271392 (Realtime)
- https://support.google.com/analytics/answer/7201382 (DebugView)

**Cloudflare**
- https://developers.cloudflare.com/cache/advanced-configuration/crawler-hints/
- https://developers.cloudflare.com/web-analytics/about/ ; https://developers.cloudflare.com/web-analytics/get-started/ ; https://developers.cloudflare.com/web-analytics/faq/
- https://developers.cloudflare.com/analytics/account-and-zone-analytics/zone-analytics/
- https://developers.cloudflare.com/workers/observability/metrics-and-analytics/
- https://developers.cloudflare.com/workers/static-assets/billing-and-limitations/
- https://developers.cloudflare.com/support/troubleshooting/http-status-codes/cloudflare-5xx-errors/error-522/
- https://developers.cloudflare.com/workers/configuration/routing/custom-domains/
- https://developers.cloudflare.com/rules/url-forwarding/examples/redirect-www-to-root/

**Observed (2026-10-29/30)**: `curl` of https://hastcode.com/ (headers, consent snippet), `/robots.txt`, `/sitemap.xml`, `/news-sitemap.xml`, `/465f05c90b9917457c1ad83b8df48051.txt`, https://www.hastcode.com/ (522); `dig www.hastcode.com`; `gh run view 40691800972` (Publish PUBCF-20261028-125421: `Ping IndexNow: skipped`); `gh variable list` (`INDEXNOW_ENABLED false`); `.github/workflows/publish.yml`; `.vercel/output/wrangler.json`; `internal/views/analytics.go`.

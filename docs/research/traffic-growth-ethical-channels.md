# How can hastcode.com grow traffic by ethical means, with one Author and two Pieces a day?

Research for hastcode.com: a static Go-built developer publication (Briefs + Explainers: Case Studies, Concepts, Comparisons, Repo Profiles with fetched Fact Packs), about 67 Pieces live, two new Pieces a day, one solo Author under the disclosed house byline Mira Vance, AdSense live, GA4, RSS, a Buttondown newsletter, and since 2026-11-13 an Agent Surface (llms.txt, sitemap.md, Markdown twins). "Ethical" means inside Google's spam policies and each platform's own rules: no bought traffic, no vote rings, no sockpuppets, affiliation disclosed. Indexing, IndexNow and the three-number log are settled in [`first-traffic-new-domain.md`](first-traffic-new-domain.md); buying traffic is answered "no" in [`buying-traffic.md`](buying-traffic.md); answer engines in [`aeo-citation-signals.md`](aeo-citation-signals.md); News/Discover eligibility in [`google-news-discover-requirements.md`](google-news-discover-requirements.md); why US/UK readers matter in [`ad-rpm-tech-audience.md`](ad-rpm-tech-audience.md); the video persona in [`persona-inspiration-cleo-abram.md`](persona-inspiration-cleo-abram.md). This note starts where those stop.

Researched **2026-11-13**. Labels: **[P]** primary (the platform's own guideline, help, policy or API page, fetched that day; the page's date where shown) / **[P, snippet]** an official page whose text reached me only via search snippets because the page returned 403 or refused the connection / **[S]** secondary (blogs, vendor guides, mirrors) / **[observed]** what a live request returned / **[inferred]** my deduction / **[opinion]**. Quotation marks are verbatim.

---

## TL;DR

- **Hacker News is the one channel that is both allowed and large, and it has never seen hastcode.com** (0 submissions in the HN Algolia index [observed]). "It's ok to post your own stuff part of the time"; "Don't solicit upvotes, comments, or submissions"; Show HN excludes "blog posts" [P]. Budget: one self-submission a week from an account that mostly submits other people's work.
- **Lobsters is closed for 70 days**: hastcode.com has never been submitted (`/domains/hastcode.com` is 404 [observed]) and new users "can't ... submit links to domains we haven't seen submitted before" or "submit self-promo" [P]. Get an invite, participate, wait.
- **Reddit is unverified**: every Reddit page is blocked to the fetcher and the search tool. A mirror of r/golang's wiki says LLM-generated posts are forbidden [S]. Read each subreddit's rules by hand first.
- **Canonical republishing is supported on DEV, Hashnode and Medium, but our disclosed AI-assisted drafting is a real constraint**: DEV says AI-assisted posts should not "Promote any business, program, or course (including your own)" [P, 2024-04-08]; daily.dev rejects "AI-generated content" and "Content used primarily for self-promotion", and "Company blogs and personal blogs aren't eligible" as sources [P, 2026-11-07].
- **One developer newsletter publishes a route that welcomes your own work**: Changelog News, "Submitting your own work is also encouraged" [P]. Cooperpress has an editorial address; TLDR, Pointer, SWLW, ByteByteGo, Programming Digest and Hacker Newsletter publish none. Architecture Weekly takes PRs but "is not for advertising ourselves" [P]. Pragmatic Engineer closed guest pitches in Aug 2026 [P].
- **Earned links only**: Google lists buying links, "Excessive link exchanges", advertorials, "optimized anchor text in articles, guest posts, or press releases", forum signatures and "Low-quality directory or bookmark site links" as link spam [P]. Fact Packs are first-party data; data and good Profiles are the legitimate bait [opinion].
- **Two free Google levers**: the Preferred Sources button/deeplink (domain-level sites; Top Stories and, since 2026-05-27, AI Mode and AI Overviews) [P], and Discover "follow" via Search profiles, which need 10,000 followers on one platform [P, snippet], so not yet.
- **Owned channels are rule-light if consent is real**: Buttondown "requires double opt-in for all newsletters by default" [P]; CAN-SPAM penalties are "up to $53,088" per email with "no exception for business-to-business email" [P]; GDPR consent is "a clear affirmative act" [P]. A weekly RSS-to-email digest is a paid Buttondown feature [P].
- **Shorts carry no link**: "URLs placed in YouTube Shorts comments and Shorts descriptions are non-clickable" [P]. No video in the next 30 days [opinion].
- **Tag everything you can**: GA4 files HN, Lobsters and Mastodon as Referral by hostname; Email needs `utm_medium=email`; always set `utm_source`, `utm_medium`, `utm_campaign`, lowercase [P].
- **"Working" at this scale** [S 2009-2023; opinion]: one HN front page is 2k-50k visits in a day. By day 30: one such spike, Search Console clicks rising for three weeks, 50-100 confirmed subscribers.

---

## 1. Aggregators and communities

| Platform | Rule, verbatim [P unless marked] | Verdict [inferred/opinion] |
|---|---|---|
| **Hacker News** ([guidelines](https://news.ycombinator.com/newsguidelines.html), [FAQ](https://news.ycombinator.com/newsfaq.html), [Show HN](https://news.ycombinator.com/showhn.html); undated) | "Please don't use HN primarily for promotion." "It's ok to post your own stuff part of the time, but the primary use of the site should be for curiosity." "Don't solicit upvotes, comments, or submissions." FAQ: "We penalize or ban submissions, accounts, and sites that break this rule." Reposts: "a small number of reposts is ok" if no "significant attention in the last year"; "Please don't delete and repost." "Don't post generated text or AI-edited text." "please don't automate posting." "use the original title ... don't editorialize." Show HN is "something you've made that other people can play with"; off topic: "blog posts, sign-up pages, newsletters, lists, and other reading material." A moderator emailed a blogger over repeated self-submission; he now self-submits "a couple times per month"; on paid upvotes: "That's basically the equivalent of a capital offense on HN and we ban accounts and sites that do it" [S, Geerling 2022-02-10]. | **Primary channel.** One real account used mostly for others' stories; one Piece a week, Case Studies first, original title, never Show HN, hand-written comments, no asks. Lands as GA4 Referral (`news.ycombinator.com`). |
| **Lobsters** ([about](https://lobste.rs/about)) | "New users (first 70 days) are not allowed to submit self-promo." "New users can't send invites, submit links to domains we haven't seen submitted before". "self-promo should be less than a quarter of one's stories and comments in the last year." "The quickest way to receive an invitation is to talk to someone you recognize from the site." Domain page 404 [observed]. The "author" checkbox is behind login, unverified. | Invite from someone you know; 70 days of real participation; then one Piece with the author box, under a quarter of activity. Earlier = sockpuppet. |
| **Reddit** (sitewide rules, help centre, r/programming, r/golang, r/ExperiencedDevs, r/softwarearchitecture, r/devops) | **Not fetched**: reddit.com, old.reddit.com, redditinc.com, support.reddithelp.com blocked; search tool refuses reddit.com. Mirrors: r/golang's 2022 rule "Must Be Go Related" and wiki line that "using LLMs to generate blog posts or comments remains forbidden" [S]; a vendor guide claims r/programming self-promotion is Saturdays-only for projects [S, 2026]. | Read each rule set by hand, record it, post nowhere that bans LLM-assisted content; elsewhere rarely, disclosed, from an account with history. |
| **daily.dev** ([suggest a source](https://docs.daily.dev/docs/for-content-creators/suggest-new-source), 2026-11-07; [content guidelines](https://docs.daily.dev/docs/for-content-creators/content-guidelines), 2026-10-21) | "You need 250 reputation." "It must be a well-known publication or developer blogging platform." "Company blogs and personal blogs aren't eligible." "It needs a public RSS or Atom feed." "Reviews can take up to 30 days." Rejected: "AI-generated content, automated algorithms, or content with characteristics typical of AI-generated material"; "Content used primarily for self-promotion"; "Paid or sponsored content distributed outside our official programs." "Individual creators can start a Squad or post directly." | Source listing: not yet. A Squad is the sanctioned self-post route (low reach). Re-apply after months of verified output. |
| **Flipboard** ([RSS post](https://about.flipboard.com/inside-flipboard/rss-feed-on-flipboard/), undated) | "The feed needs to be the full feed, not just titles or extracts."; an image "at least 400px on its shortest side" per post; "at least 30 items." Publisher sign-up with feed review [S, search summary]. | Cheap; do it once the feed is confirmed full-text with images. Small, non-developer traffic. |
| **Feedly** ([help](https://docs.feedly.com/article/288-how-to-add-a-source), 2024-11-04) | "If you're a publisher and want your website available in Feedly, just make sure the RSS URL is public." | Nothing to submit; keep `<link rel="alternate">` and a descriptive feed title. |
| **DEV** ([editor guide](https://dev.to/p/editor_guide); [AI guidelines](https://dev.to/guidelines-for-ai-assisted-articles-on-dev), 2024-04-08) | `canonical_url`: "link for the canonical version of the content". AI-assisted posts should "Disclose the fact that they were generated or assisted by AI" (`#ABotWroteThis`) and not "Promote any business, program, or course (including your own)." | Ask DEV in writing first. If yes: excerpt + canonical, disclosed; never full text (syndicated copies can outrank the original, `first-traffic-new-domain.md` §2). |
| **Hashnode** ([CLI post](https://hashnode.com/blog/hashnode-cli), 2026-10-06) | "a canonical URL with flags such as `--subtitle`, `--series-id`, and `--original-article-url`." No AI/self-promo policy found. | Supported; same excerpt rule; low priority. |
| **Medium** (import and canonical help pages: 403) | [P, snippet] imported stories "will automatically apply a canonical URL referencing the original source material"; "you may only import content ... you have the right to". Medium Rules: 403. | Supported; same rule; not our Reader. |
| **Substack** (support: 403, snippets) | No support article on canonical tags found; cross-post is for other publications' posts. | Do not mirror Pieces: no canonical, and Buttondown exists. |
| **LinkedIn** ([publishing guidelines](https://www.linkedin.com/help/linkedin/answer/a519782); [Professional Community Policies](https://www.linkedin.com/legal/professional-community-policies)) | "You can republish something that you have published somewhere else, as long as it is original content that you own the rights to." "Advertisements and promotions for events, products, or services aren't appropriate content for article publishing." "Don't do things to artificially increase engagement with your content." "don't agree with others ahead of time to like or re-share each other's content." No canonical. | Native 150-word take plus link; a newsletter edition may carry the TL;DR, never the full Piece. No pods. Organic Social. |
| **Bluesky** ([guidelines](https://bsky.social/about/support/community-guidelines), 2025-09-19) | "Do not send spam or repeatedly post content in ways that disrupt normal conversations". "Do not artificially manipulate features or social signals to gain unearned reach". "Do not post undisclosed commercial content, sponsored material, or advertising." | Allowed: an account that says it is hastcode, with a sentence of why per Piece. |
| **Mastodon** (Fosstodon [rules API](https://fosstodon.org/api/v1/instance/rules); Hachyderm hung up twice) | "DO NOT post commercial promotions, or advertise"; "DO NOT use automated tools to post without also monitoring and/or interacting from your account." | Human-run account that replies; no feed-only bot. Ask the instance's admins first. Tag links (not in GA4's social list). |
| **X** (authenticity and copypasta policies: 403) | [P, snippet] not allowed: "Operating multiple accounts that interact with the same or substantially similar content or in order to inflate or manipulate the prominence of content"; link-only posting as most of an account's activity; duplicates "subject to review and enforcement under our platform manipulation and spam policy". | Allowed with commentary from one account; low yield for a text site. |

---

## 2. Developer newsletters

| Newsletter | Stated route [P] | Verdict |
|---|---|---|
| **Changelog News** ([submit](https://changelog.com/news/submit)) | "Submitting your own work is also encouraged." Not accepted: "How-to's and tutorials", "Commercial products/services", "Reader-hostile websites". "Do your best to convince us why something is newsworthy." Podcast ideas: separate form. | Submit one Case Study or Brief a week with a two-line why. |
| **Golang Weekly / Cooperpress** ([site](https://cooperpress.com/)) | No form; "Editorial" contact `editor@cooperpress.com`. | One email per Go-relevant Piece, subject = original title, never repeated. |
| **Architecture Weekly** ([README](https://github.com/oskardudycz/ArchitectureWeekly)) | "Feel free to add your links by sending PR with update to Upcoming.md"; "we should not be adding links to our posts - this repository is not for advertising ourselves but sharing the knowledge." | Never PR your own Pieces; PR the first-party Sources you cite [inferred]. |
| **Pragmatic Engineer** ([guest page](https://blog.pragmaticengineer.com/pragmatic-engineer-guest-article/)) | "Aug 2026: The Pragmatic Engineer is no longer accepting guest post pitches." | Closed. |
| **TLDR, Pointer, Software Lead Weekly, Programming Digest, ByteByteGo, Hacker Newsletter** | No submission route on any site [observed]; Hacker Newsletter is "hand-curated the best articles from Hacker News". | Curators find links on HN and Lobsters; §1 is the route. |

---

## 3. Search growth beyond indexing

**Helpful content and AI drafting** [P, [Creating helpful content](https://developers.google.com/search/docs/fundamentals/creating-helpful-content); [Using generative AI content](https://developers.google.com/search/docs/fundamentals/using-gen-ai-content), updated 2026-10-01]: "People-first content means content that's created primarily for people, and not to manipulate search engine rankings." AI disclosures "are useful for content where someone might think" how it was made. Mass-produced AI pages "may violate Google's spam policy on scaled content abuse"; it is "critical to manually factcheck and review all AI-generated content for accuracy". The Editorial Policy note and the `verified:` gate (ADR-0007) meet this; the risk is pace, not method [inferred].

**Link spam vs earned links** [P, [Spam policies](https://developers.google.com/search/docs/essentials/spam-policies)]: "Link spam is the practice of creating links to or from a site primarily for the purpose of manipulating search rankings." Examples: "Buying or selling links for ranking purposes", "Excessive link exchanges", "Using automated programs or services to create links to your site", "Advertorials or native advertising where payment is received for articles that include links that pass ranking credit", "links with optimized anchor text in articles, guest posts, or press releases distributed on other sites", "Low-quality directory or bookmark site links", "Forum comments with optimized links in the post or signature". Paid links are fine "as long as they are qualified with a `rel="nofollow"` or `rel="sponsored"`". Scraped content includes "Reproducing content feeds from other sites without providing some type of unique benefit to the user": a Brief that only re-feeds a Source [inferred].

**Internal linking and topical clusters** [P, [Link best practices](https://developers.google.com/search/docs/crawling-indexing/links-crawlable)]: "Every page you care about should have a link from at least one other page on your site." "Good anchor text is descriptive, reasonably concise, and relevant to the page that it's on." Google publishes no "topical authority" rule; the nearest is the [Discover core update](https://developers.google.com/search/blog/2026/02/discover-core-update) (2026-02-05), which (paraphrased by the fetch tool) favours in-depth, original, timely content from sites with expertise judged topic by topic. For us: every Brief links an Explainer, every Explainer its Related Explainers, Vertical pages are the clusters [inferred].

**Preferred sources** [P, [doc](https://developers.google.com/search/docs/appearance/preferred-sources); changelog 2026-05-27 "Preferred sources is available in AI Mode and AI Overviews", 2026-08-20 custom button]: "Only domain-level and subdomain-level sites are eligible". Add `https://news.google.com/swg/js/v1/publisher.js` + `<div google-add-preferred-source-btn></div>`, or the deeplink `https://www.google.com/preferences/source?q=hastcode.com`; AI surfaces require the site be "included in Search generative AI features in Search Console". Not a ranking claim; it changes what opted-in readers see. Cheapest unpulled lever [opinion].

**Discover follow and Search profiles** [P, [Google blog 2025-09-17](https://blog.google/products-and-platforms/products/search/discover-updates-september-2025/)]: "Now, you can "follow" publishers or creators right on Discover". The RSS-based Follow is gone: "The Follow feature is no longer shown in Google Discover" (changelog 2025-11-19). Search profiles [P, snippet]: Canada, India, US; 10,000 followers on YouTube, Instagram, X or TikTok. Revisit when any channel passes 10k.

**Bing and IndexNow** [P]: participants "Amazon", "Bing", "Naver", "Seznam.cz", "Yandex", "Yep"; submissions "shared across all IndexNow-enabled search engines"; Google absent ([FAQ](https://www.indexnow.org/faq)). Bing "processes sitemaps at least once every 24 hours" and powers "Yahoo, DuckDuckGo, Microsoft Start, and AI tools like ChatGPT and Inflection.ai" ([2025-06-17](https://blogs.bing.com/webmaster/2025/6/Start-Using-Bing-Webmaster-Tools-to-Improve-Your-Site-Visibility/)). Bing's Guidelines and URL-submission pages are client-rendered: no text, third attempt.

**Backlinks that are clearly allowed** (each needs someone else's editorial decision):
- *Original data.* A quarterly Piece from the Fact Packs with a downloadable table is cite-able first-party data [opinion].
- *Cited by the project you profile.* Neither Apache Iggy's nor FrankenPHP's README has a community-articles section [observed]; rare and per-repo. Tell maintainers a Profile exists; never PR your own link [opinion].
- *Awesome lists.* awesome-go takes packages, not articles, and wants descriptions "clear, concise, and non-promotional" [P]; awesome-scalability: "Make sure the list is useful before submitting", one PR per item [P]. Submit only where articles are in scope, disclosed.
- *Guest posts with editorial oversight* are allowed; anchor-text campaigns are not [P]. The one named venue with a form has closed.
- *Podcasts.* Changelog has a request form [P]; show-note links are earned [inferred].
- *Wikipedia is not a channel*: "you should avoid linking to a site that you own, maintain, or represent" [P, [External links](https://en.wikipedia.org/wiki/Wikipedia:External_links)].

---

## 4. Direct and owned channels

- **Email law** [P, [FTC CAN-SPAM guide](https://www.ftc.gov/business-guidance/resources/can-spam-act-compliance-guide-business), Aug 2023/Jan 2024]: "Don't use false or misleading header information." "Don't use deceptive subject lines." "Identify the message as an ad." "Tell recipients where you're located." "Honor opt-out requests promptly." "Monitor what others are doing on your behalf." Penalties "up to $53,088" per email; "no exception for business-to-business email." [P, [GDPR Recital 32](https://gdpr-info.eu/recitals/no-32/)]: "Consent should be given by a clear affirmative act." "Silence, pre-ticked boxes or inactivity should not therefore constitute consent." So: no pre-ticked boxes, a postal address in the footer, no imported contacts who never signed up [inferred].
- **Buttondown** [P]: "requires double opt-in for all newsletters by default" ([double opt-in](https://docs.buttondown.com/double-opt-in)). [Standards](https://docs.buttondown.com/standards) (2025-05-30): complaint rate "0.1% or lower", permanent failure "1% or lower", temporary failure "5% or lower"; imports sent "in small initial batches". [RSS-to-email](https://docs.buttondown.com/rss-to-email): "automatically send a new email every time you publish a new blog post"; weekly "on a day and time of your choosing"; digests via a template loop; paid feature. Terms say nothing about purchased lists; the law is the floor.
- **Weekly digest** [opinion]: Monday, 5-7 Pieces with one-line whys, drafted by RSS-to-email, sent by the Author, every link tagged `utm_source=buttondown&utm_medium=email&utm_campaign=weekly-YYYY-WW`. Growth inside the rules: a signup line at the end of every Piece, a `/newsletter/` page with no AdSense units (`buying-traffic.md`), public archive, no first-view pop-ups.
- **RSS** [observed]: `/feed.xml` is live (RSS 2.0). Confirm full-text with an image per item (Flipboard) and that Markdown twins are not in it.
- **Web push** [P, [web.dev](https://web.dev/articles/push-notifications-permissions-ux)]: "The worst thing you can do is to show the permission dialog to users as soon as they land on your site." "First show a dialog that your website controls, explaining the value". Opt-in from an ad-blocking developer audience will be tiny [inferred]; not in 30 days.

---

## 5. Video and audio

- **YouTube** [P, [monetization policies](https://support.google.com/youtube/answer/1311392), 2025-07-15]: not allowed: "Content downloaded or copied from another online source without any substantive modifications"; "Similar or repetitive content with low educational value, commentary, narratives, or minimal variation across videos" (policy renamed "inauthentic content"). Links [P, [Sharing links](https://support.google.com/youtube/answer/13748639)]: Shorts URLs "non-clickable"; "To add clickable external links, you must enable your channel's access to advanced features"; channel page "up to 14 links". A templated "read the Piece" Short per Piece is the inauthentic pattern and carries no link [inferred].
- **TikTok** [P, snippet; guidelines page refused twice]: bought engagement banned; multiple accounts allowed "for fan content or creative expression" but "not to deceive others or break the rules".
- **LinkedIn video**: same policies as §1.
- **Effort** [opinion]: one explainer video at the persona note's standard (pitch, research, visual outline, script) costs a Piece. The house byline is a disclosed non-person, so on-camera needs a named human. Not in 30 days.

---

## 6. Measurement

- **GA4 channels** [P, [Default channel group](https://support.google.com/analytics/answer/9756891)]: Direct = source "(direct)" and medium "(not set)"/"(none)"; Organic Search = search-site list or medium "organic"; Organic Social = Google's social-site list or medium in ("social", "social-network", "social-media", "sm", "social network", "social media"); Email = source or medium `email|e-mail|e_mail|e mail`; Referral = medium "referral", "app", "link"; Organic Video = video list or medium `^(.*video.*)$`. So HN, Lobsters, daily.dev, Flipboard and Mastodon instances are Referral by hostname; an untagged newsletter click is Direct or Referral, never Email [inferred].
- **UTM convention** [P, [Collect campaign data](https://support.google.com/analytics/answer/10917952)]: "should always use" `utm_source`, `utm_medium`, `utm_campaign`; case sensitive; one source per platform, one medium per channel. Proposal: `utm_source` = `buttondown`/`bluesky`/`mastodon`/`linkedin`/`x`; `utm_medium` = `email`/`social`; `utm_campaign` = Piece slug or `weekly-YYYY-WW`. Never tag HN or Lobsters links; the hostname attributes them and the communities read tracking as promotion [opinion]. AI Overviews/AI Mode clicks sit in Search Console's "Web" total (`aeo-citation-signals.md`).
- **Per channel**: HN/Lobsters/Reddit → Referral; DEV/Hashnode/Medium copies → Referral; newsletter → Email (tagged); Bluesky/LinkedIn/X → Organic Social; Mastodon → tagged `social`; Changelog/Cooperpress → Referral; Preferred Sources/Discover → Search Console, not GA4.
- **Benchmarks** [S 2009-2023]: #1 on HN for a day ≈ 49.5k visits (Willison 2023); #1 for six hours ≈ 31.8k pageviews (Burke); a typical front page 2k-10k visitors (Schmalbach). **Working by day 30** [opinion]: one front page or three stories above 20 points; Search Console clicks up three weeks running; one newsletter link; 50-100 confirmed subscribers; weekly three-number log kept. Raptive's 25k/month at 50% Tier-1 stays a 6-12-month target.

---

## 7. What not to do

| Practice | Forbidden by [P unless marked] |
|---|---|
| Asking for upvotes; upvoting from a second account; paid upvotes | HN "Don't solicit upvotes, comments, or submissions"; "We penalize or ban submissions, accounts, and sites"; "capital offense" [S quoting the moderator]. |
| Delete-and-repost; AI-written comments; automated submissions | "Please don't delete and repost." "Don't post generated text or AI-edited text." "please don't automate posting." |
| Submitting hastcode.com to Lobsters from a new or borrowed account | New users "can't ... submit links to domains we haven't seen submitted before"; "not allowed to submit self-promo". |
| Posting AI-drafted Pieces where LLM content is banned | r/golang wiki (mirror) [S]; daily.dev "AI-generated content ... characteristics typical of AI-generated material". |
| Engagement pods; several X accounts boosting one link | LinkedIn "don't agree with others ahead of time to like or re-share"; X "Operating multiple accounts that interact with the same or substantially similar content" [P, snippet]. |
| Undisclosed promotion on Bluesky; a feed-only bot on Fosstodon | "Do not post undisclosed commercial content"; "DO NOT use automated tools to post without also monitoring". |
| Full-text syndication without canonical; Substack or LinkedIn mirrors | Google: syndicated copies can outrank the original (`first-traffic-new-domain.md` §2); no canonical on either. |
| Undisclosed DEV cross-posts with a link back | DEV: disclose; not "Promote any business ... (including your own)". |
| Buying links, link swaps, keyword-anchor guest posts, forum signatures, directories, self-links to Architecture Weekly or Wikipedia | Google link-spam list (§3); "not for advertising ourselves"; "avoid linking to a site that you own, maintain, or represent". |
| Paid placements presented as editorial | Google "Advertorials or native advertising"; daily.dev "Paid or sponsored content distributed outside our official programs"; FTC disclosure (`buying-traffic.md` §5). |
| Briefs that re-feed Sources; many thin pages | "Reproducing content feeds from other sites without ... unique benefit"; scaled content abuse. |
| Imported contacts, pre-ticked boxes, no postal address | CAN-SPAM; GDPR Recital 32; Buttondown double opt-in. |
| Templated link-bait Shorts | YouTube inauthentic content; Shorts URLs "non-clickable". |
| Buying traffic of any kind | `buying-traffic.md`. |

---

## 30-day plan, one Author, two Pieces a day [opinion]

Distribution gets 45 minutes a day and never a Piece's verification time.

| Days | Action | Effort | GA4 channel |
|---|---|---|---|
| 1-2 | Preferred Sources button/deeplink in the article footer and `/about/`; confirm "Search generative AI features" in Search Console. | 1 h | Search Console, not GA4 |
| 1-3 | UTM convention in `docs/author-runbook.md`; tag newsletter links; `/newsletter/` page without AdSense units. | 2 h | Email (from day 8) |
| 1-3 | Confirm `/feed.xml` is full-text with images; submit to Flipboard; add `<link rel="alternate">` if missing. | 1 h | Referral |
| 3-7 | HN account: if it has no non-self history, spend week 1 submitting and commenting on others' work only. | 15 min/day | — |
| 7, 14, 21, 28 | One HN self-submission a week: strongest verified Case Study or Repo Profile, original title, hand-written comments, no asks; a flop waits 30 days. | 30 min | Referral |
| 7, 14, 21, 28 | Same Piece to Changelog News with a two-line why; Go-relevant Pieces once to `editor@cooperpress.com`. | 15 min | Referral / Email |
| 8, then Mondays | Weekly digest via RSS-to-email draft, Author sends, tagged. | 1 h, then 30 min/week | Email |
| 8-30 | One Bluesky and one LinkedIn post per Piece, 1-2 sentences of why, disclosed; reply by hand. Mastodon only after asking an instance's admins. | 15 min/day | Organic Social |
| 10 | Ask a Lobsters member you know for an invite; start the 70-day clock with others' work. | 1 h, then 10 min/day | — (Referral after day 70) |
| 10 | Email DEV about disclosed, verified, canonical cross-posts; act only on a written yes. | 30 min | Referral |
| 12 | Read five subreddits' rules by hand; record them; post nowhere that bans LLM-assisted content. | 1 h | Referral |
| 15-30 | Internal-link audit: every Brief → an Explainer; every Explainer → Related Explainers; Vertical pages link every Piece. | 3 h | Organic Search (later) |
| 20 | Tell two profiled repos' maintainers, in their channels, that a Profile exists; no link request. | 1 h | Referral (their choice) |
| 30 | Review: three-number log, Referral by hostname, Email clicks, subscribers, Search Console trend; decide on a quarterly Fact Pack data Piece. | 1 h | — |

Not in the plan: daily.dev source listing, Substack, Medium, video, push, paid anything.

---

## Not fetched / open

1. **Reddit**: all pages blocked; read subreddit rules by hand.
2. **Hachyderm** rules: connection hung up (`/about`, two API endpoints); Fosstodon stands in.
3. **Medium Rules** and canonical/import help: 403; snippets.
4. **X** authenticity and copypasta policies: 403; snippets.
5. **TikTok** guidelines: connection refused twice; snippets.
6. **Substack**: no official page on canonical tags either way.
7. **Lobsters** author checkbox: behind login.
8. **Bing Webmaster Guidelines** / URL-submission quota: client-rendered, no text.
9. **Flipboard** review criteria: undated blog post plus search summaries.
10. **Search profiles** eligibility: snippet of Google Search Help.
11. **DEV**: AI-assisted canonical cross-posts, open since 2026-10-30.
12. `/feed.xml` full-text/images and the Buttondown plan's RSS-to-email: not checked.

---

## Sources (accessed 2026-11-13)

**Aggregators and communities**
- https://news.ycombinator.com/newsguidelines.html ; https://news.ycombinator.com/newsfaq.html ; https://news.ycombinator.com/showhn.html
- https://hn.algolia.com/api/v1/search?query=%22hastcode.com%22&advancedSyntax=true&restrictSearchableAttributes=url (0 hits) [observed]
- https://www.jeffgeerling.com/blog/2022/i-almost-got-banned-hacker-news/ (2022-02-10) [S]
- https://lobste.rs/about ; https://lobste.rs/domains/hastcode.com (404) [observed] ; https://lobste.rs/stories/new (login wall)
- https://docs.daily.dev/docs/for-content-creators/suggest-new-source (2026-11-07) ; https://docs.daily.dev/docs/for-content-creators/content-guidelines (2026-10-21)
- https://about.flipboard.com/inside-flipboard/rss-feed-on-flipboard/ ; https://about.flipboard.com/publishers/
- https://docs.feedly.com/article/288-how-to-add-a-source (2024-11-04)
- https://dev.to/p/editor_guide ; https://dev.to/guidelines-for-ai-assisted-articles-on-dev (2024-04-08)
- https://hashnode.com/blog/hashnode-cli (2026-10-06)
- https://help.medium.com/hc/en-us/articles/214550207-Importing-a-post-to-Medium (403) ; https://help.medium.com/hc/en-us/articles/360033930293-Set-a-canonical-link (403) ; https://policy.medium.com/medium-rules-30e5502c4eb4 (403)
- https://support.substack.com/hc/en-us/articles/360037830351 ; https://support.substack.com/hc/en-us/articles/10522003894932 (403)
- https://www.linkedin.com/help/linkedin/answer/a519782/publishing-platform-guidelines ; https://www.linkedin.com/legal/professional-community-policies
- https://bsky.social/about/support/community-guidelines (2025-09-19)
- https://fosstodon.org/api/v1/instance/rules ; https://hachyderm.io/about (hung up)
- https://help.x.com/en/rules-and-policies/authenticity (403) ; https://help.x.com/en/rules-and-policies/copypasta-duplicate-content (403)
- Not fetched: https://support.reddithelp.com/hc/en-us/articles/204536055 (403); redditinc.com, old.reddit.com (blocked); r/golang 2022-04-14 mod post and wiki via redlib mirrors [S]; https://www.teract.ai/resources/reddit-marketing-developer-tools-2026 [S]

**Newsletters**
- https://changelog.com/news ; https://changelog.com/news/submit
- https://cooperpress.com/ ; https://golangweekly.com/
- https://github.com/oskardudycz/ArchitectureWeekly (README)
- https://blog.pragmaticengineer.com/pragmatic-engineer-guest-article/ (closed Aug 2026)
- https://tldr.tech/ ; https://www.pointer.io/ ; https://softwareleadweekly.com/ ; https://programmingdigest.net/ ; https://www.bonobopress.com/ ; https://hackernewsletter.com/ ; https://blog.bytebytego.com/about

**Google Search**
- https://developers.google.com/search/docs/fundamentals/creating-helpful-content ; https://developers.google.com/search/docs/fundamentals/using-gen-ai-content
- https://developers.google.com/search/docs/essentials/spam-policies ; https://developers.google.com/search/docs/crawling-indexing/links-crawlable
- https://developers.google.com/search/docs/appearance/preferred-sources ; https://developers.google.com/search/docs/appearance/search-profiles ; https://support.google.com/websearch/answer/16905608 (snippet)
- https://developers.google.com/search/docs/appearance/google-discover ; https://developers.google.com/search/blog/2026/02/discover-core-update (2026-02-05) ; https://blog.google/products-and-platforms/products/search/discover-updates-september-2025/ (2025-09-17)
- https://developers.google.com/search/updates (entries 2025-11-19, 2026-01-30, 2026-04-30, 2026-05-15, 2026-05-27, 2026-08-20, 2026-09-16, 2026-10-01, 2026-11-12)
- https://en.wikipedia.org/wiki/Wikipedia:External_links

**Bing / IndexNow**
- https://www.indexnow.org/faq ; https://blogs.bing.com/webmaster/2025/6/Start-Using-Bing-Webmaster-Tools-to-Improve-Your-Site-Visibility/ (2025-06-17) ; https://www.bing.com/webmasters/help/url-submission-62f2860b (no text)

**Backlink venues**
- https://github.com/avelino/awesome-go/blob/main/CONTRIBUTING.md ; https://github.com/binhnguyennus/awesome-scalability/blob/master/CONTRIBUTING.md
- https://github.com/apache/iggy ; https://github.com/dunglas/frankenphp (READMEs; no community-articles section) [observed]

**Owned channels**
- https://www.ftc.gov/business-guidance/resources/can-spam-act-compliance-guide-business (Aug 2023 / Jan 2024) ; https://gdpr-info.eu/recitals/no-32/
- https://docs.buttondown.com/double-opt-in ; https://docs.buttondown.com/standards (2025-05-30) ; https://docs.buttondown.com/rss-to-email ; https://buttondown.com/legal/terms
- https://web.dev/articles/push-notifications-permissions-ux
- https://hastcode.com/feed.xml ; https://hastcode.com/llms.txt (200) [observed]

**Video**
- https://support.google.com/youtube/answer/1311392 (2025-07-15) ; https://support.google.com/youtube/answer/13748639
- https://www.tiktok.com/community-guidelines/en/integrity-authenticity (refused; snippets via tiktok.com/safety and creator-academy pages)

**Measurement**
- https://support.google.com/analytics/answer/9756891 ; https://support.google.com/analytics/answer/10917952
- [S] https://simonwillison.net/2023/Feb/17/analytics/ ; https://kevin.burke.dev/kevin/hn-popular-post-stats/ ; https://www.vincentschmalbach.com/analyzing-a-year-of-hacker-news-traffic/

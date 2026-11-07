# Should hastcode.com buy traffic, and if so from where?

Research for hastcode.com: a static Go-built developer publication (news Briefs + architecture Explainers, ~69 Pieces live since 2026-10-28) on Cloudflare, with AdSense live, GA4, an Atom feed and a Buttondown newsletter; the Author is solo and India-based. Goals: Raptive eligibility (25k pageviews/month, 50% from US/UK/CA/AU/NZ, GA4 required; see [`ad-rpm-tech-audience.md`](ad-rpm-tech-audience.md) §4), then newsletter sponsorship, ultimately $10k/month. The tension: paid traffic costs per click, display ads pay per thousand views, and every ad network polices where traffic comes from. Free channels, indexing and HN/Lobsters/DEV rules are in [`first-traffic-new-domain.md`](first-traffic-new-domain.md); RPMs and network entry bars in [`ad-rpm-tech-audience.md`](ad-rpm-tech-audience.md); answer-engine signals in [`aeo-citation-signals.md`](aeo-citation-signals.md). This note does not repeat them.

Researched **2026-11-07**. Labels: **[P]** primary (the platform's own policy, help page or rate card, fetched that day, with the page's updated date when shown) / **[S]** secondary (trade press, agency benchmarks, blogs, with year) / **[observed]** something fetched or measured live / **[inferred]** my deduction / **[opinion]**. Policy wording is quoted verbatim; nothing in quotation marks is a paraphrase. "[P, snippet]" marks an official page whose text reached me only through a search-engine snippet because the page itself returned 401/403; those pages are listed under Open questions.

---

## TL;DR

- **Buying ads to an AdSense site is allowed; buying "traffic" is where the trouble starts.** AdSense: "AdSense publishers are welcome to use online advertising as a way to get new users to their sites", but "AdSense publishers are ultimately responsible for the traffic to their ads", and "we strongly urge you to use caution when partnering with third-party traffic services" [P]. Paid-to-click, paid-to-surf, autosurf and click-exchange programs are banned outright [P].
- **Raptive treats bought traffic as a policy violation, not a neutral source.** Its Content and Traffic Policies list "Traffic purchased with the intent of making more money from advertising" beside "Fraudulent traffic", "Suspicious traffic sources" and "Non-human traffic" [P, updated 2026-07-10]. Mediavine wants "Real visitors, not bots, incentivized clicks, or traffic from sources that put advertisers at risk" and reviews "Traffic sources" [P]. Neither says paid traffic counts toward its threshold.
- **Paid traffic does nothing for ranking.** Google: "Does buying ads improve my search ranking?" "No. Google's advertising programs are independent of our search results." Discover uses "many of the same signals and systems used by Search" [P].
- **The arbitrage does not pencil out, by two orders of magnitude.** At the kindest scenario in this note ($15 page RPM, no ad-blocking, 1.2 pages/session) a visitor is worth **$0.018**; with the developer ad-block rate of 72% it is **$0.005**. The cheapest real click (Meta ≈ $0.60, Google Display ≈ $0.63, Reddit $0.50-2.50) costs **30-150×** that. Buying the 25,000 pageviews Raptive asks for would cost roughly $12,500/month at $0.60/click, and would then fail Raptive's own traffic policy [S benchmarks; inferred].
- **The only "traffic" cheaper than its ad value is the kind that closes accounts.** SparkTraffic sells visits from "$0.08 per 1 000 credits" and says itself: "Paid traffic to a page with pay-per-click or pay-per-impression ads breaks the rules of publisher ad programs" [observed]. AdSense: "If we observe high levels of invalid traffic on your account, we may suspend or close the account"; deductions "are considered final and are not subject to appeal" [P]. Hard no.
- **Developer ad channels are priced for products, not publications.** EthicalAds: $3.10-8.00 CPM, "$1,000 minimum ad buy" [P, Q4 2026]. daily.dev: no rate card, quarterly plans. Carbon: quoted per campaign, "1M+ impressions a month". Stack Overflow, TLDR, ByteByteGo, Cooperpress, Pointer, Console: sales-led, prices not published [P]. Bonobo Press and Changelog publish rates, and they imply **$54-88 per 1,000 subscribers per issue** for 21-29k-subscriber developer newsletters [P].
- **Where money could eventually go: newsletter subscribers, not pageviews.** A verified subscriber costs about $2.50 via beehiiv paid recommendations or "$2-$7 per confirmed referral" via SparkLoop [P]; at developer-newsletter sponsorship rates that pays back in roughly a year, but only once a list is big enough to sell. hastcode has no rate card and (as of the last log) near-zero readers, so the honest answer today is **spend $0 on traffic**.

---

## 1. AdSense and Google Publisher Policies on traffic sources

**Program policies** [P, [AdSense Program policies](https://support.google.com/adsense/answer/48182), updated 2026-08-04]:
- "Google ads may not be placed on pages receiving traffic from certain sources." Sites may not "Use third-party services that generate clicks or impressions such as paid-to-click, paid-to-surf, autosurf and click-exchange programs", nor "Be promoted through unsolicited mass emails or unwanted advertisements on third-party websites", nor "Display Google ads, search boxes or search results as a result of the actions of software applications such as toolbars", nor "Be loaded by any software that can trigger pop-ups".
- The page never uses "incentivised". Its equivalent: publishers may not "Compensate users for viewing ads or performing searches, or promise compensation to a third party for such behavior", and "Except for rewarded inventory, publishers may not ask others to click or view their ads."
- "Any method that artificially generates clicks or impressions on your Google ads is strictly prohibited." Consequence: "we reserve the right to disable ad serving to your site and/or disable your AdSense account at any time", and "If your account is disabled, you will not be eligible for further participation in the AdSense program."

**Invalid traffic** [P, [Invalid traffic](https://support.google.com/adsense/answer/16737)]: "Invalid traffic includes any clicks or impressions that may artificially inflate an advertiser's costs" and "covers intentionally fraudulent traffic as well as accidental clicks". Examples: "Automated clicking tools or traffic sources, robots, or other deceptive software", "Publishers encouraging clicks on their ads", "Repeated ad clicks or impressions generated by one or more users". "We understand that a third party may generate invalid traffic on your ads without your knowledge or permission", but "ultimately it is your responsibility as the publisher to ensure that the traffic on your ads is valid." On detection: "If we observe high levels of invalid traffic on your account, we may suspend or close the account"; "if we are unable to verify the quality of your traffic, we may limit your ad serving"; "you may also see a difference between your estimated and finalized earnings."

**Deductions** [P]: "we can't share any specific details on where or how much invalid activity may have been detected"; deductions appear as "a line item deduction in your "Transactions" page for invalid activity" ([answer 6053295](https://support.google.com/adsense/answer/6053295)). "No, all deductions from earnings are considered final and are not subject to appeal"; "In most cases, deductions correspond to traffic within the past 60 to 90 days" ([Deductions FAQ](https://support.google.com/adsense/answer/2808531)).

**Buying advertising vs buying traffic.** Google draws the line itself:
- "AdSense publishers are welcome to use online advertising as a way to get new users to their sites." Pages the ads land on "are required to comply with the spirit of Google's Landing Page Quality Guidelines": "relevant and original content, transparency, and site navigability". "Ads shouldn't take users to pages that consist mostly of more ads, general search results" or "pages whose primary purpose is to redirect traffic to or mirror the content of a parent site" [P, [Use of online advertising](https://support.google.com/adsense/answer/1348727)].
- "You're welcome to promote your site in any manner that complies with our program policies. However, AdSense publishers are ultimately responsible for the traffic to their ads." "we've found that some of these services actually send artificial traffic to websites, despite their appearance" and "we strongly urge you to use caution when partnering with third-party traffic services." If a source looks suspicious, Google recommends "stopping or pausing your relationship" with it; and "If you are considering purchasing traffic to drive to your site, we recommend that you start by reading our traffic provider checklist" [P, [If you want to purchase traffic](https://support.google.com/adsense/answer/1348722)].
- The [checklist](https://support.google.com/adsense/answer/3332805) [P]: "If you don't recognize the pages where your ads will be placed, ask how those pages get their traffic"; "It's a warning sign if you are unable to get sample URLs where you can view your ads"; "If prices seem too good to be true" the provider may not be "appearing on high value pages"; and if traffic arrives "that you weren't expecting", "turn off that traffic source immediately."
- Ad Traffic Quality, [Resources for Publishers](https://www.google.com/ads/adtrafficquality/publishers) [P]: patterns that can violate policy are "Usage of automated bots that imitate user behavior", "Buying unrelated traffic from unknown networks", "Participating in pay-to-click programs", "Receiving traffic from unwanted email or software applications". "While traffic exchange services aren't strictly prohibited, some exchange services may lead to invalid clicks or impressions." The [account-closure list](https://support.google.com/adsense/answer/2660562) adds that automated traffic can be "received through purchased traffic" [P].

**Google Publisher Policies** [P, [overview](https://support.google.com/publisherpolicies/answer/10502938)] contain no "traffic sources" section; the relevant behavioural rule is inventory value: "We do not allow Google-served ads on screens: with more ads or other paid promotional material than publisher-content", where publisher-content "includes images, videos, games, article text and user-generated content created or managed by the publisher" but "does not include whitespace, header or footer content, or links to other content on the site" [P, [answer 11169917](https://support.google.com/publisherpolicies/answer/11169917)].

**Google Ads (advertiser side).** Google does not publish an "arbitrage" or "made for advertising" policy under those names; the [Abusing the ad network](https://support.google.com/adspolicy/answer/6020954) page says only that "Using the Google Network to gain an unfair traffic advantage over other participants in the auction is not allowed" [P]. The operative rule is [Destination requirements](https://support.google.com/adspolicy/answer/6368661): "ad destinations must offer unique value to users"; disallowed are "Destination content that is designed for the primary purpose of showing ads", "Destination content that is replicated from another source without adding value", and "Destinations that are solely designed to send users elsewhere" [P]. An Explainer with one or two ad units is not that; a thin page built to be a landing page would be [inferred].

---

## 2. Raptive, Mediavine and Journey on paid vs organic traffic

| Network | What its own pages say [P] | Does paid traffic count? [inferred] |
|---|---|---|
| **Raptive** ([eligibility](https://help.raptive.com/hc/en-us/articles/360032840891-Who-is-eligible-for-Raptive), updated 2026-09-26, via Zendesk JSON) | "Our minimum pageview requirement is 25,000 pageviews a month." 25k-99,999: "50% of traffic comes from the United States, the United Kingdom, Canada, New Zealand, and/or Australia"; "Google Analytics is correctly set up"; domain ≥ 6 months. The article says nothing about paid, social or bot traffic. | The threshold page is silent, so the policy page decides. |
| **Raptive** ([Content and Traffic Policies](https://help.raptive.com/hc/en-us/articles/29872333292443-Raptive-Content-and-Traffic-Policies), updated 2026-07-10) | "Examples of traffic that would violate our policies include but are not limited to: Fraudulent traffic, Suspicious traffic sources, Non-human traffic, Traffic purchased with the intent of making more money from advertising". No consequence is stated on the page. | **Disqualifying** when the purpose is ad revenue. Intent is the test and cannot be proven either way, so any paid campaign to an ad-bearing page is a risk to the application. |
| **Raptive** ([bots article](https://help.raptive.com/hc/en-us/articles/46089937895579), updated 2026-10-16; [applying blog](https://raptive.com/blog/applying-to-join-adthrive/), 2025-11-24) | "Invalid traffic is any non-human or non-intentional human behavior." "If high amounts of IVT are detected, it can get your domain flagged and greatly reduce advertiser spend." Reviewers check that pageviews come from "high-quality, trustworthy sources". Fix: ask the host to block bots; "Add the paid version of Cloudflare". | Bot-style "traffic" harms the domain even before an application. |
| **Mediavine** ([approval](https://help.mediavine.com/what-does-it-take-to-get-approved-by-mediavine); [requirements](https://www.mediavine.com/mediavine-requirements/)) | "Generate a minimum of $5,000 in annual ad revenue"; "Maintain clean, human, brand-safe traffic": "Real visitors, not bots, incentivized clicks, or traffic from sources that put advertisers at risk." Evaluated items include "Traffic countries of origin" and "Traffic sources". Sites are reviewed for "content quality, audience engagement, traffic integrity, and overall reader experience". | Not excluded by name; "sources that put advertisers at risk" is the catch-all. Paid social from a mainstream platform is probably tolerated, bought "visits" are not [opinion]. |
| **Journey by Mediavine** ([blog](https://www.mediavine.com/introducing-journey-new-ad-management/), 2024-03-14, updated 2026-09-29; requirements page) | "a minimum of 1,000 sessions from tier-1 countries, such as the US, Canada, UK, and Australia, within a 30-day period", elsewhere "a minimum threshold of 1,000 premium sessions per month", plus "our other content and traffic quality requirements". | Sessions must be Tier-1 and "premium"; the page does not say paid sessions are excluded, but the same "traffic standards" apply. |

I found no Mediavine help article specifically on paid, social or Pinterest traffic; its "known-bot" article covers GA4 filtering only, and a search snippet attributing "social traffic ... valuable" to Mediavine did not appear on the page it pointed to (Open questions). GA4 itself: "traffic from known bots and spiders is automatically excluded" and "you cannot disable known bot traffic exclusion or see how much known bot traffic was excluded" [P, [GA4 help](https://support.google.com/analytics/answer/9888366)], so unknown bots still inflate the GA4 number Raptive reads [inferred].

---

## 3. Google Search and Discover: does paid traffic move ranking?

- "Does buying ads improve my search ranking?" "No. Google's advertising programs are independent of our search results." "Participation in an advertising program doesn't positively or negatively affect inclusion or ranking in the Google search results." "Inclusion and ranking in Google Search results don't cost money; we don't accept payment to expedite inclusion or improve a site's ranking for particular keywords." [P, [Site position FAQ](https://developers.google.com/search/help/site-position-in-search-faq)]. The companion statement that "Google doesn't accept payment to crawl a site more frequently" is quoted in [`first-traffic-new-domain.md`](first-traffic-new-domain.md) §1.1.
- **Discover**: "Content is automatically eligible to appear in Discover if it is indexed by Google and meets Discover's content policies." "No special tags or structured data are required." "Discover makes use of many of the same signals and systems used by Search to determine what is helpful, people-first content." "Traffic from Discover is less predictable or dependable when compared to keyword-driven search visits." Nothing on the page mentions traffic volume or paid promotion as an input [P, [Discover](https://developers.google.com/search/docs/appearance/google-discover)].
- **Scaled content abuse** [P, [Spam policies](https://developers.google.com/search/docs/essentials/spam-policies)]: "Scaled content abuse is when many pages are generated for the primary purpose of manipulating search rankings and not helping users", including "Using generative AI tools or other similar tools to generate many pages without adding value for users". Buying traffic to AI-assisted Pieces neither triggers nor cures this: the policy is about why pages were made, not how readers arrived. The risk that does change is AdSense's landing-page rule (§1): a paid click must land on "substantial and unique content", so never on a thin or list page [inferred].

---

## 4. Unit economics: what a click costs vs what a visitor earns

### 4.1 Platform minimums and benchmark CPCs

| Platform | Minimum / pricing rule [P unless noted] | Benchmark cost per click for a tech or B2B audience [S] |
|---|---|---|
| Google Ads, Search | No minimum published on the bidding page; "Bidding amounts are ultimately up to you" [P, [bidding](https://business.google.com/us/ad-tools/bidding/)] | $5.42 all-industry; "Business Services" $5.87, "Education & Instruction" $4.81 (LocaliQ/WordStream 2026, Google+Microsoft search, updated 2026-06-01) [S] |
| Google Ads, Demand Gen / Display | Demand Gen runs "across YouTube (including Shorts), Discover, Gmail, Maps, and the Google Display Network"; budget guidance is a daily budget of at least 10× target CPA; Display campaigns migrate to Demand Gen from June 2026 [P, [About Demand Gen](https://support.google.com/google-ads/answer/13695777)] | Display ≈ $0.63 all-industry (vendor compilations of WordStream data, 2025-26); technology display $0.51 and search $3.80 are WordStream figures from a 2016 page [S]. Demand Gen CPC: not published. |
| Microsoft Ads | "You only pay when someone clicks your ad" [P, [new-advertiser FAQ](https://about.ads.microsoft.com/en/get-started/new-advertiser-faqs)]; "no minimum spend" wording seen only in a snippet [P, snippet] | Blended into the LocaliQ $5.42 figure above [S] |
| Reddit Ads | "at least $5.00/USD per day" (self-serve) [P, snippet; the cost page returned 401] | $0.50-2.50 across verticals; "Technology/SaaS subreddits" CPC ≈ $1.50-4.00, CPM $6-12 (Stackmatix, 2026-09-11, no dataset named) [S] |
| Meta (Facebook/Instagram) | "start with at least $5 for your budget and to choose a duration over six days" [P, [pricing](https://www.facebook.com/business/ads/pricing)]; USD 1/day for impression-billed sets [P, snippet; the minimum-budget page served an empty body] | Traffic campaigns $0.60 all-industry, "Business Services" $0.70 (LocaliQ 2026, updated 2026-09-23) [S] |
| X Ads | "there's no minimum campaign spend"; minimum bid "$.01 per desired action" [P, snippet; pages returned 403] | ≈ $0.74 average (agency roundups, 2026, no primary data) [S] |
| Quora Ads | "$5 minimum daily budget per campaign"; "no account-level minimums, no setup fees, and no required spend commitments" [P, [pricing](https://business.quora.com/pricing)]; CPC "Minimum bid is $0.01", CPM "$0.20" [P, [bidding](https://business.quora.com/resources/bidding-on-quora-ads)] | $1.50-4 for US B2B SaaS (vendor guides, 2026) [S] |
| LinkedIn | "The minimum daily budget amount required to start on LinkedIn is $10, for any ad format"; lifetime minimum $100 [P, [budget page](https://business.linkedin.com/marketing-solutions/success/best-practices/maximize-your-budget)] | $11.12 average CPC (Quora for Business citing Search Engine Land 2026) [S] |
| Taboola (Realize) | "The recommended daily budget per campaign is 10 times your CPA goal, or at least $50 USD"; "your daily cap could be up to double the daily cap set" [P, [budget](https://realize.com/help/en/articles/3878086-setting-your-budget-and-pacing), 2026-07-23]; minimum CPC $0.01, suggested start $0.20-0.70 [P, snippet] | Non-conversion campaigns: "100-200 clicks per day", e.g. "$0.50 CPC" → $50-100/day [P, [new campaign](https://realize.com/help/en/articles/11426249-setting-up-a-new-campaign)] |
| Outbrain | "The minimum daily budget you need to create a campaign is 20USD, 20EUR, or 20GBP per day" [P, [help](https://www.outbrain.com/help/advertisers/budget-duration/)]; raised: "The minimum daily budget for a new campaign will be increased from $20.00 to $25.00 USD" effective 2025-11-03 [P, [API group](https://groups.google.com/g/outbrain-amplifyapi/c/l2rPqTr11i4)] | Suggested starting CPC $0.65-0.85 [P, snippet; the setup page returned 404] |

### 4.2 The arbitrage table

Assumptions: page RPM scenarios from [`ad-rpm-tech-audience.md`](ad-rpm-tech-audience.md) ($3 India-heavy AdSense, $8 ≈ the one tech data point, $15 Tier-1 on a premium network); 1.2 pages per session; 72% of developer sessions serve no ad (Stack Overflow 2018, same doc §7). Revenue per visitor = RPM ÷ 1,000 × 1.2 × share of ad-served pageviews.

| Page RPM | Revenue per visitor, 72% blocked | Revenue per visitor, 0% blocked | Break-even CPC (72% / 0% blocked) |
|---|---|---|---|
| $3 | $0.0010 | $0.0036 | $0.0010 / $0.0036 |
| $8 | $0.0027 | $0.0096 | $0.0027 / $0.0096 |
| $15 | $0.0050 | $0.0180 | $0.0050 / $0.0180 |

| Channel (cheapest plausible click, from §4.1) | Cost per visitor | Cost ÷ value at $15 RPM, 72% blocked | Cost ÷ value at $15 RPM, 0% blocked |
|---|---|---|---|
| Taboola at its suggested $0.20 start | $0.20 | 40× | 11× |
| Reddit, broad, low end | $0.50 | 100× | 28× |
| Meta traffic campaign | $0.60 | 120× | 33× |
| Google Display / Demand Gen (display proxy) | $0.63 | 126× | 35× |
| X | $0.74 | 148× | 41× |
| Reddit, technology subreddits / Quora US B2B | $1.50 | 300× | 83× |
| EthicalAds at $3.10 CPM, assuming 0.2% CTR [inferred; CTR not published] | $1.55 | 310× | 86× |
| Google Search, Business Services | $5.87 | 1,174× | 326× |
| LinkedIn | $11.12 | 2,224× | 618× |
| SparkTraffic "credits" (§6) | $0.00008-0.00016 | 0.02-0.03× | 0.004-0.009× |

**Verdict [inferred]:** no legitimate channel gets within an order of magnitude of break-even, even at an RPM hastcode will not see for years and with every visitor seeing ads. The one line that "works" is bot traffic priced below a tenth of a cent, which is exactly what §1 and §2 prohibit. Reframed as a Raptive play: 25,000 pageviews ÷ 1.2 pages ≈ 20,800 clicks ≈ **$12,500/month at $0.60**, for a site whose ads would earn **$75-375** on those views, and whose application would then fail "Traffic purchased with the intent of making more money from advertising". Arbitrage is a dead end, not a question of finding a cheaper network.

---

## 5. Developer-specific paid channels, newsletter sponsorships and community rules

| Channel | Price / minimum [P unless noted] | Format | Audience claimed | Will it sell to a content site? |
|---|---|---|---|---|
| [EthicalAds](https://www.ethicalads.io/advertisers/pricing/) (Q4 2026) | CPM: "All developers" $3.10 run-of-network / $5.25 US-CA-UK-AU-NZ-IE; "Backend web development" $3.80 / $6.25; "DevOps" $4.10 / $6.50; "AI / machine learning" $4.60 / $7.50; "Custom niche targeting" $5.00 / $8.00. "$1,000 minimum ad buy"; 10% off at $3,000, 15% at $25,000. | Text + image, "approved by Acceptable Ads" | "over 35 million developer impressions each month" | No stated restriction; "campaigns are approved by the EthicalAds team". |
| [Carbon Ads](https://www.carbonads.net/) | "We quote each campaign based on timing, targeting and volume"; "Pay per thousand ad views (CPM)"; "we plan around 1M+ impressions a month, run for about three months". Native CPC product: "priced per click (CPC)", figure not published. | Native unit on "350+ hand-picked tech sites" | not published beyond site count | Sales-led (/contact); no eligibility stated. |
| [daily.dev Ads](https://business.daily.dev/pricing) | "We do not publish a rate card"; "one monthly plan in three sizes: Starter, Growth and Scale"; "A plan runs for a first quarter and then renews quarter by quarter". A $5,000 minimum appears only in daily.dev blog posts [S]. | In-feed native, Digest newsletter, homepage sponsorship ("2M+ to 60M+" impression floors) | "millions of developers" | Not a fit for "teams with no technical content yet" or "products developers do not use"; a publication is an unusual buyer. |
| [Stack Overflow Advertising](https://stackoverflow.co/advertising/) | Not published; "contact our team" | Banner, native, tag/site sponsorship, newsletter, podcast | "millions" monthly uniques | Sales-led; buyers listed are marketing, employer-brand and tech teams. |
| [TLDR](https://advertise.tldr.tech/sponsorship-opportunities) | Not published; custom proposal | Quick Links (≤30 words), Secondary (≤50), Primary (≤100, logo); max 3 advertisers per issue | TLDR Tech 1.6M, TLDR Dev 470k, TLDR DevOps 350k; "41%" open rate | "tech brands"; no stated rule against publishers. |
| ByteByteGo ([Substack post](https://blog.bytebytego.com/p/newsletter-sponsorships), 2022-05-24) | Not published; email hi@bytebytego.com; a marketplace listing shows "Contact for Pricing" [S] | newsletter placement | "Engineering: 94%"; ~1M+ subscribers per Substack milestone post [S] | Not stated. |
| Cooperpress / Golang Weekly ([advertise](https://cooperpress.com/advertise)) | Not published; media kit PDF (text not extractable here); "we maintain a waiting list"; "Net open rates 35%-60%" | sponsorships and "content development" | "over 450,000 developers and software engineers"; Golang Weekly ≈ 39,400 subscribers (Feb 2026, media-kit figure via search snippet) [S] | Not stated. |
| [Pointer](https://www.pointer.io/sponsorship) | Not published; form or email | newsletter placement | "55K" subscribers, "40%" open rate, "90%" "US based", "50%" manage teams | Not stated. |
| Bonobo Press ([media kit](https://www.bonobopress.com/media-kit/), updated 2026-10-19) | Programming Digest **$1,135/issue** (21,149 subs, 45.57% open); Leadership in Tech **$2,570** (29,158; secondary $1,799); C# Digest $1,345 (21,077); React Digest $1,375 (22,463; secondary $962) | "The format of the newsletter is text only": link, title <100 chars, description <400 | Europe 35-48%, US 30-35% | Ads cover "software tools and products, jobs, conferences, webinars, books, and courses"; no publisher exclusion. |
| [Console](https://console.dev/advertise) | Not published; email advertise@console.dev; advertisers "must meet our selection criteria" | two slots per Thursday edition: main ad; contents-only ad | "36% of our readers are in the US, 30% in Europe, 12% in the UK, and 22% in the rest of the world" | Criteria page not fetched; Console features devtools, so a publication is doubtful [inferred]. |
| [Changelog](https://changelog.com/sponsor/pricing) | Changelog News "$1,500 per week (no minimums)": "22,000+ subscribers with a 51.35% open rate", "20,000 to 30,000 listens per episode"; The Changelog "$3,000 per week", "Four (4) week minimum"; Fallthrough "$200 per episode" | podcast pre/mid-roll + newsletter | "100,000 developers weekly"; "650K monthly listens" | Email adam@changelog.com; no eligibility stated. |
| Hacker Newsletter, Software Lead Weekly | No sponsorship section on either site [observed]; SWLW "$2,749 for a 6-week" slot appears only in a third-party list [S] | — | "60,000+" / "31,020 readers" | Ask by email. |
| Hacker News | No ads. "Please don't post job ads as submissions to HN." YC-funded startups' job ads "appear on the front page, but are not stories: they have no vote arrows, points, or comments"; "startup launches get placed on the front page" for YC companies only [P, [FAQ](https://news.ycombinator.com/newsfaq.html)]. Asking for upvotes gets sites "penalized or banned". | — | — | Nothing to buy. |
| Lobsters | Invite-only; new users cannot submit unseen domains or self-promo; self-promo "less than a quarter" of activity [P, quoted in `first-traffic-new-domain.md` §2] | — | — | Nothing to buy. |
| Reddit communities | Reddit's rules pages and the r/programming, r/golang, r/ExperiencedDevs rule pages could not be fetched (403/blocked, same as on 2026-10-30). No claim made. | — | — | Read by hand (Open questions). |
| Bluesky, Mastodon | Bluesky: "advertising couldn't be our dominant business model" [P, [2023 plan](https://bsky.social/about/blog/7-05-2023-business-plan)]; no ad product exists, though the CEO has not ruled ads out [S]. Mastodon: "No algorithms or ads to waste your time." "We will never serve ads or push profiles for you to see." [P, [joinmastodon.org](https://joinmastodon.org/)] | — | — | Nothing to buy. |

FTC, one line: where "there's a connection between an endorser and the marketer that a significant minority of consumers wouldn't expect", "that connection should be disclosed clearly and conspicuously" [P, [FTC Endorsement Guides FAQ](https://www.ftc.gov/business-guidance/resources/ftcs-endorsement-guides-what-people-are-asking)]. Any paid newsletter mention of hastcode must say it is paid [inferred].

**Implied sponsorship CPM per issue** from the two published rate cards [P, arithmetic mine]:

| Newsletter | Price per issue | Subscribers | $ per 1,000 subscribers per issue |
|---|---|---|---|
| Programming Digest | $1,135 | 21,149 | $53.67 |
| React Digest | $1,375 | 22,463 | $61.21 |
| C# Digest | $1,345 | 21,077 | $63.81 |
| Changelog News | $1,500 | 22,000 | $68.18 |
| Leadership in Tech | $2,570 | 29,158 | $88.14 |
| Paved marketplace, "highly targeted" average | — | — | $23.01 [S, [Paved guide](https://www.paved.com/blog/the-ultimate-newsletter-sponsorship-guide), 2026-03-09, "most rebooked newsletters on Paved"] |
| Paved marketplace, "general interest" average | — | — | $6.53 [S, same] |

---

## 6. Cheap "website traffic" sellers: a hard no

- **What they sell** [observed, [SparkTraffic](https://www.sparktraffic.com/)]: credits, "from $0.08 per 1 000 credits" one-off or "$0.16 down to $0.07 per 1,000" monthly, one visit per credit or more with targeting; "Visits that start with a person seeing an ad are real people"; an "engaged" session is defined as over 10 seconds, two pageviews or a key event (the same three criteria as GA4's engaged session, which the service can therefore tune to). The site's own warning: "Paid traffic to a page with pay-per-click or pay-per-impression ads breaks the rules of publisher ad programs", and it recommends an ad-free landing page. Fiverr-style "10k visitors" gigs sell the same thing at the same price band [inferred from the price].
- **What it is under Google's definition**: "Automated clicking tools or traffic sources, robots, or other deceptive software" and "Buying unrelated traffic from unknown networks" (§1). The checklist's tell, "If prices seem too good to be true", applies: $0.08 per 1,000 visits is below any ad platform's minimum bid [P; inferred].
- **Fiverr's own rules** [P, snippet; the Community Standards pages returned 403]: artificial engagement "through bots, click farms, or other illegitimate techniques intended to increase followers, views, likes, fans, or traffic" is prohibited, and Gigs that "violate the Terms of Service of another website" are not allowed. So such gigs breach Fiverr's rules as well as Google's.
- **What GA4 would show** [P definitions; inferred outcome]: "Direct is the channel by which users arrive at your site/app via a saved link or by entering your URL", rule "Source exactly matches "(direct)"" [P, [Default channel group](https://support.google.com/analytics/answer/9756891)]; SparkTraffic says "Many ad formats do not pass a referrer" so visits land as direct unless UTM-tagged [observed]. A Raptive or Mediavine reviewer who sees a Direct spike with sub-10-second sessions, or oddly uniform "engaged" sessions, has what Raptive calls "Suspicious traffic sources" (§2).
- **Consequences**: AdSense "may suspend or close the account", deductions are "final and are not subject to appeal", and Raptive says high IVT "can get your domain flagged and greatly reduce advertiser spend" (§1, §2). Label: **never**.

---

## 7. Where paid spend could make sense, and the recommendation

**(a) Newsletter subscriber acquisition.** The only paid model with a revenue line on the other side.

| Input | Figure | Source |
|---|---|---|
| Cost per verified subscriber, beehiiv paid recommendations | you set the CPA; example "you will pay $2.50 per subscriber acquisition" and the recommender "will be paid $2"; "Available on paid beehiiv plans" | [P, [beehiiv help](https://www.beehiiv.com/support/article/14492963616279-growing-your-audience-with-beehiiv-boosts), 2026-10-20] |
| Cost per confirmed referral, SparkLoop Partner Network | "Generally, you'll get paid $2-$7 per confirmed referral" (publisher side; the buyer pays that plus fees) | [P, [SparkLoop help](https://help.sparkloop.app/en/articles/9768833-everything-you-need-to-know-about-paid-recs), 2024-08-21] |
| Sponsorship revenue per subscriber per issue, developer newsletters | $0.054-0.088 (Bonobo, Changelog rate cards, §5) | [P, arithmetic] |
| Sponsorship revenue per subscriber per issue, marketplace average | $0.023 niche / $0.0065 general | [S, Paved 2026] |
| Revenue per subscriber per year, one slot in 52 weekly issues | $2.79 at Programming Digest's rate; $4.58 at Leadership in Tech's; $1.20 at Paved's niche average | [inferred] |
| Payback on a $2.50 subscriber | ≈ 11 months at $2.79/yr; ≈ 25 months at $1.20/yr; never at the general-interest rate | [inferred] |

The catch: every rate above belongs to a list of 21-29k subscribers that sells every issue. A list of 500 has no rate card and no buyer, so a subscriber bought today earns $0 until the list is sellable, and churn eats the payback [inferred]. hastcode is on Buttondown; beehiiv paid recommendations require a beehiiv plan, and whether SparkLoop supports Buttondown was not checked (Open questions).

**(b) Seeding one launch Piece.** $50-100 on Reddit or X aimed at US/UK buys 25-150 visits at the §4.1 CPCs. That is nothing for pageviews and matters only if one of those readers links or reposts. It is allowed under AdSense's "welcome to use online advertising" as long as the page is "substantial and unique content", but it lands inside Raptive's "Traffic purchased" wording if ads are on the page [inferred]. If done at all: land on a page without AdSense units, tag with UTM so GA4 files it under Paid Social, and keep it below one such test a quarter [opinion].

**(c) Retargeting existing readers.** Needs an ads pixel. Google's EU user consent policy requires "end users' legally valid consent to" "the use of cookies or other local storage where legally required" and "the collection, sharing, and use of personal data for personalization of ads" on "Google products used on any site, app or other property that is under your control" [P, [policy](https://www.google.com/about/company/user-consent-policy/)]. hastcode's tag already defaults `analytics_storage` to denied in EEA/UK/CH and serves an AdSense consent message ([`first-traffic-new-domain.md`](first-traffic-new-domain.md) §3.1), so ad-storage consent would have to be added to that flow. Google remarketing segments need "a minimum of 100 active visitors or users within the last 30 days" to serve on Display, Search or YouTube [P, [Audience manager](https://support.google.com/google-ads/answer/7558048)]. With consented traffic near zero the list would not reach 100. Not viable now [inferred].

**(d) Google Ads Demand Gen, Tier-1 only.** Google recommends a daily budget of at least 10× the target CPA (§4.1) and publishes no CPC; the display proxy of ≈ $0.63 sits 35-126× above a visitor's ad value (§4.2). As a newsletter-signup campaign it collapses to case (a) with a worse CPA than paid recommendations, because every click still has to convert [inferred].

**Recommendation [opinion].** Spend **$0 on traffic for the site**, now and for as long as display ads are the revenue model: no network, developer-specific or otherwise, sells a click for less than roughly 30× what a hastcode pageview can earn, Raptive names purchased traffic as a policy violation, and the only cheap "traffic" is invalid by Google's definition. Spend on **newsletter growth only**, and only after three things exist: a published hastcode sponsorship price, a list large enough that someone has paid it at least once, and a measured open rate. When that day comes, the guardrails are: paid recommendations (beehiiv or SparkLoop-style, verified subscribers only), **CAC ceiling = one year of sponsorship revenue per subscriber at the rate actually charged** (≈ $1.50-2.50 at today's developer-newsletter rates), Tier-1 geo only, a cap of one month's realised sponsorship revenue per month, the signup page (no AdSense units) as the only landing page, UTM on everything, and a weekly check of GA4's Paid channel share and AdSense's invalid-traffic line item. Until then the money is better spent on the free channels already listed in [`first-traffic-new-domain.md`](first-traffic-new-domain.md).

---

## Open questions

1. **Reddit**: Reddit's rules, the three subreddits' rules, and Reddit Ads' cost and billing pages all returned 401/403 or were blocked. The "$5.00/USD per day" minimum rests on search snippets of the official page. Read by hand.
2. **X Ads**: pricing, campaigns-101 and bidding FAQ pages returned 403. "no minimum campaign spend" and "$.01 per desired action" are snippet-only.
3. **Meta**: the "Best practices for minimum budgets" page served an empty body in two locales; the USD 1/day impression minimum is snippet-only.
4. **Microsoft Ads**: no fetched page contains "no minimum spend"; only "You only pay when someone clicks your ad" was confirmed.
5. **Taboola/Realize**: the "minimum budget" article 404s after redirect; the $0.01 minimum CPC and $0.20-0.70 start are snippet-only. Outbrain's setup page (starting CPC $0.65-0.85) returned 404.
6. **WordStream** benchmark pages returned 403; the LocaliQ (same company) 2026 pages were used instead and carry no explicit data window beyond "2026".
7. **Cooperpress media kit**: the PDF is image/outline-only to the tools here; Golang Weekly's 39,398 subscribers (Feb 2026) comes from a search snippet of it. Rates are not published anywhere fetched.
8. **Fiverr Community Standards**: three URLs returned 403; wording is from snippets.
9. **Mediavine on social/paid traffic**: no help article found; the snippet "social traffic ... valuable" was not on the 2019 pageviews post it pointed to. Ask Mediavine support whether paid social sessions count toward Journey's 1,000.
10. **Raptive intent test**: the policy names "Traffic purchased with the intent of making more money from advertising" but not how intent is judged, nor what share of paid traffic is tolerated. Ask before any paid campaign.
11. **EthicalAds CTR** is not published, so its CPC-equivalent in §4.2 is an assumption.
12. **Console's "selection criteria"** page and Swapstack's pricing page were not fetched.
13. **SparkLoop + Buttondown**: whether SparkLoop paid recommendations work with Buttondown was not checked; beehiiv's require a beehiiv plan.
14. **Google Ads "arbitrage"**: Google publishes no policy under that name; the closest wording is "Destination content that is designed for the primary purpose of showing ads". Whether a two-ad-unit Explainer page ever trips it is unknown; the landing-page guidance says it should not.

---

## What to do (ordered) [opinion]

1. **Decide it: no paid traffic to hastcode.com pages.** Record it in an ADR so the question is not reopened every quarter. *(§1 AdSense responsibility and checklist; §2 Raptive "Traffic purchased"; §4.2 the 30-150× gap.)*
2. **Never buy "visitors".** No Fiverr gigs, SparkTraffic or traffic bots, including "to test analytics". *(§6; AdSense deductions "final"; Raptive IVT flagging.)*
3. **Keep the ad units out of any page you might ever advertise.** Add a no-ads template for the newsletter signup page now, so a future campaign has a compliant landing page. *(§1 landing-page guidance; §7b.)*
4. **Build the sponsorship side before the acquisition side.** Write a one-line rate card only when the list and open rate exist; price from the §5 table (≈ $55-90 per 1,000 subscribers per issue for developer newsletters, $23 marketplace niche average). *(§5, §7a.)*
5. **When a rate card exists, test paid recommendations, not ads.** $100-200 once, Tier-1 geo, CAC ceiling = one year of sponsorship revenue per subscriber, verified subscribers only; stop if cost per verified subscriber exceeds the ceiling. *(§7a; beehiiv/SparkLoop pricing.)*
6. **Instrument before any spend.** UTM on every paid link; watch GA4 Paid Social/Paid Search channel share and the AdSense Transactions page for an invalid-activity line. *(§1 "monitor your reports closely"; §6 GA4 Direct definition.)*
7. **Ask Raptive and Mediavine in writing** how paid-social sessions are treated before any campaign, and keep the answer with the ADR. *(Open questions 9-10.)*
8. **Spend the time, not the money, on the channels in `first-traffic-new-domain.md`**: indexing, Bing import, one HN submission, newsletter and feed. Paid cannot substitute for links, which are also what Discover and Search reward. *(§3.)*

---

## Sources (accessed 2026-11-07)

**Google AdSense / Publisher Policies / Ad Traffic Quality**
- https://support.google.com/adsense/answer/48182 (Program policies, updated 2026-08-04)
- https://support.google.com/adsense/answer/16737 (Invalid traffic)
- https://support.google.com/adsense/answer/6053295 (record of deductions) ; https://support.google.com/adsense/answer/2808531 (Deductions FAQ)
- https://support.google.com/adsense/answer/1348722 (If you want to purchase traffic) ; https://support.google.com/adsense/answer/3332805 (Traffic provider checklist) ; https://support.google.com/adsense/answer/1348727 (Use of online advertising) ; https://support.google.com/adsense/answer/10654953 (Leverage Google Analytics) ; https://support.google.com/adsense/answer/2660562 (violations that lead to closure) ; https://support.google.com/adsense/answer/9724 (eligibility)
- https://support.google.com/publisherpolicies/answer/10502938 ; https://support.google.com/publisherpolicies/answer/11169917
- https://www.google.com/ads/adtrafficquality/ ; https://www.google.com/ads/adtrafficquality/publishers
- https://support.google.com/admanager/answer/2796181 (Ad Manager: Purchase traffic to your site)

**Google Ads policy and help**
- https://support.google.com/adspolicy/answer/6368661 (Destination requirements) ; https://support.google.com/adspolicy/answer/6020954 (Abusing the ad network)
- https://support.google.com/google-ads/answer/13695777 (Demand Gen) ; https://support.google.com/google-ads/answer/7558048 (audience segment minimums)
- https://business.google.com/us/ad-tools/bidding/ ; https://www.google.com/about/company/user-consent-policy/

**Google Search / Analytics**
- https://developers.google.com/search/help/site-position-in-search-faq ; https://developers.google.com/search/docs/appearance/google-discover ; https://developers.google.com/search/docs/essentials/spam-policies
- https://support.google.com/analytics/answer/9756891 (Default channel group) ; https://support.google.com/analytics/answer/12195621 (Engagement rate) ; https://support.google.com/analytics/answer/9888366 (Known bot exclusion)

**Raptive / Mediavine**
- https://help.raptive.com/api/v2/help_center/en-us/articles/360032840891.json (eligibility, updated 2026-09-26) ; …/29872333292443.json (Content and Traffic Policies, 2026-07-10) ; …/46089937895579.json (bots and IVT, 2026-10-16) ; https://raptive.com/blog/applying-to-join-adthrive/ (2025-11-24)
- https://help.mediavine.com/what-does-it-take-to-get-approved-by-mediavine ; https://www.mediavine.com/mediavine-requirements/ ; https://www.mediavine.com/introducing-journey-new-ad-management/ (updated 2026-09-29) ; https://www.mediavine.com/blog/how-many-pageviews-do-i-need-for-mediavine/ (2019-10-14)

**Ad platforms (pricing / minimums)**
- Reddit: https://business.reddithelp.com/en/categories/billing-payment/how-much-do-reddit-ads-cost (401) ; https://business.reddithelp.com/articles/Knowledge/max-campaigns
- Quora: https://business.quora.com/pricing ; https://business.quora.com/resources/bidding-on-quora-ads ; https://business.quora.com/marketing-advertising-resources/best-b2b-advertising-benchmarks-2026 [S] ; https://business.quora.com/marketing-advertising-resources/linkedin-ads-cost-2026 [S]
- LinkedIn: https://business.linkedin.com/marketing-solutions/success/best-practices/maximize-your-budget
- Meta: https://www.facebook.com/business/ads/pricing ; https://www.facebook.com/business/help/203183363050448 (empty body)
- X: https://business.x.com/en/help/overview/ads-pricing (403) ; https://business.x.com/en/help/campaign-setup/campaigns-101 (403) ; https://business.x.com/en/help/troubleshooting/bidding-and-auctions-faqs (403)
- Microsoft: https://about.ads.microsoft.com/en/get-started/new-advertiser-faqs ; https://about.ads.microsoft.com/en/get-started/best-practices-price-cost-control
- Taboola/Realize: https://realize.com/help/en/articles/3878086-setting-your-budget-and-pacing ; https://realize.com/help/en/articles/11426249-setting-up-a-new-campaign ; https://realize.com/help/en/articles/8740185-what-is-the-minimum-budget-to-advertise-on-taboola (404)
- Outbrain: https://www.outbrain.com/help/advertisers/budget-duration/ ; https://groups.google.com/g/outbrain-amplifyapi/c/l2rPqTr11i4 (2025-10-17) ; https://www.outbrain.com/help/advertisers/first-campaign-setup/ (404)

**Benchmarks [S]**
- https://localiq.com/blog/search-advertising-benchmarks/ (updated 2026-06-01) ; https://localiq.com/blog/facebook-advertising-benchmarks/ (updated 2026-09-23) ; WordStream pages https://www.wordstream.com/blog/google-ads-benchmarks etc. (403)
- https://www.storegrowers.com/google-ads-benchmarks/ (2026-01-08; technology $3.80 search / $0.51 display from a 2016 WordStream page)
- https://stackmatix.com/blog/reddit-ads-cost-guide-2026 (2026-09-11) ; X CPC and Quora CPC ranges from agency roundups surfaced by search (no primary data)

**Developer channels and newsletters**
- https://www.ethicalads.io/advertisers/ ; https://www.ethicalads.io/advertisers/pricing/ (Q4 2026)
- https://www.carbonads.net/ ; https://www.carbonads.net/native-cpc
- https://business.daily.dev/ ; https://business.daily.dev/pricing (daily.dev/advertise: 404)
- https://stackoverflow.co/advertising/
- https://advertise.tldr.tech/sponsorship-opportunities (tldr.tech/advertise redirects there)
- https://blog.bytebytego.com/p/newsletter-sponsorships (2022-05-24; bytebytego.com/sponsorship: 404)
- https://cooperpress.com/advertise ; https://cooperpress.com/mediakit (PDF, text not extractable; cooperpress.com/sponsor/: 404)
- https://www.pointer.io/sponsorship (pointer.io/sponsor: 404)
- https://www.bonobopress.com/media-kit/ (updated 2026-10-19) ; https://www.bonobopress.com/
- https://console.dev/advertise (console.dev/advertise/: 404)
- https://changelog.com/sponsor ; https://changelog.com/sponsor/pricing
- https://hackernewsletter.com/ ; https://softwareleadweekly.com/
- https://news.ycombinator.com/newsfaq.html ; https://lobste.rs/about (quoted via `first-traffic-new-domain.md`)
- https://bsky.social/about/blog/7-05-2023-business-plan ; https://joinmastodon.org/
- https://www.ftc.gov/business-guidance/resources/ftcs-endorsement-guides-what-people-are-asking

**Newsletter growth and marketplaces**
- https://www.beehiiv.com/support/article/14492963616279-growing-your-audience-with-beehiiv-boosts (2026-10-20) ; https://www.beehiiv.com/ad-network ; https://www.beehiiv.com/features/boosts (beehiiv.com/boosts: 404)
- https://help.sparkloop.app/en/articles/9768833-everything-you-need-to-know-about-paid-recs (2024-08-21) ; https://sparkloop.app/pricing
- https://www.paved.com/blog/the-ultimate-newsletter-sponsorship-guide (2026-03-09) ; https://www.paved.com/publishers/smart-pricing (2024-12-20) ; https://www.swapstack.co/

**Traffic sellers and Fiverr**
- https://www.sparktraffic.com/ [observed]
- https://help.fiverr.com/hc/en-us/articles/32243087473553 ; …/32243195699985 ; …/32242973123985 (all 403; wording from search snippets)

**Not fetched (Reddit)**: https://support.reddithelp.com/hc/en-us/articles/204536055 (403), https://redditinc.com/policies/content-policy (blocked), https://old.reddit.com/r/programming/about/rules (blocked), https://business.reddithelp.com/en/categories/overview/how-purchase-reddit-ads (401), https://business.reddithelp.com/en/categories/billing-payment/billing-frequently-asked-questions (401)

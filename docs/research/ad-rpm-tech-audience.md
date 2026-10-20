# Display-ad RPMs and network entry requirements for a tech/developer audience

Research for hastcode.com (tech news + software-architecture explainers; US / Europe / India; goal: $10k/month from display ads).
All sources accessed 2026-10-20. Each number carries a URL and a label: **[P]** = primary (network-owned page, Google doc, or raw survey data) / **[S]** = secondary (publisher income report, blog, trade press). Anything I could not source is marked **unsourced**.

Terminology: *page RPM* = revenue per 1,000 pageviews. Mediavine and Ezoic report per **session** (Mediavine "RPM", Ezoic "EPMV"), Raptive and AdSense per **pageview**. Session RPM ÷ pages-per-session ≈ page RPM, so a $30 session RPM on a 1.5-pages/session site is ≈ $20 page RPM. Treat cross-network comparisons with that in mind.

## TL;DR

- **Entry bars moved a lot in 2025–26.** Raptive dropped to 25k pageviews/month (Oct 2025) but still needs 50% of traffic from US/UK/CA/NZ/AU; Mediavine replaced its 50k-session rule with "$5k/yr ad revenue" and Journey now starts at 1,000 Tier-1 sessions/30 days; Ezoic went the other way and now requires **250k monthly users** for new sites (Feb 2026) unless you get into its 20-sites/month Incubator. AdSense has no traffic minimum.
- **Geography dominates RPM.** The only named per-country datasets I found (Mediavine, travel niche) put India at **$4 vs $28 US** (Nov 2021) and **$0.71 vs $21.58 US** (Sep 2023–Feb 2024): India is roughly **3–15% of US**, UK/CA/AU ≈ 75–90% of US, continental EU ≈ 55–60% of US. Nobody publishes an official India/US multiplier.
- **Tech is a mid-to-low display niche, not a premium one.** The one named-publisher figure for a tech site on Mediavine is **$11.88 session RPM** (Feb 2022, vs $25–39 for home/garden, travel, fashion on the same account). Third-party 2026 aggregates claim $12–35 for "tech/SaaS" but publish no underlying data.
- **Developers block ads at ~72%.** Computed from the raw Stack Overflow 2018 survey data (n=76,055): 72.3% have an ad blocker installed; 76.7% in the US, 82.4% in Germany, 55.6% in India. General-population rates are ~30% (GWI Q2 2025). Expect a large share of pageviews to earn $0.
- **$10k/month needs roughly 0.5–1.3M monetisable pageviews/month** at plausible blended page RPMs of $8–20 — and 2–3M+ if India-heavy traffic and developer ad-blocking pull the blended RPM toward $3–5. Display alone is a weak primary model for this audience; see the math table.

## 1. Google AdSense

- **Traffic minimum:** none. Google's eligibility page lists only: own original content that "must be high-quality, original, and attract an audience", access to the site's HTML source, age 18+, and compliance with Program policies. [P] https://support.google.com/adsense/answer/9724?hl=en
- **Policies that matter for a news/explainer site:** invalid-click/traffic rules, no paid-to-click traffic, no ads on non-content pages, ads.txt authorisation, Google Publisher Policies (content, "low-value content" screens, Better Ads Standards). No numeric content or traffic threshold is stated anywhere in these pages. [P] https://support.google.com/adsense/answer/48182?hl=en and [P] https://support.google.com/publisherpolicies/answer/10502938?hl=en
- **EEA/UK consent:** since 16 Jan 2024 a Google-certified CMP integrated with IAB TCF is required to serve personalised ads to EEA/UK users (Switzerland since 31 Jul 2024). Relevant because "Europe" is a target geo. [P] https://support.google.com/adsense/answer/13554116?hl=en
- **Revenue share:** AdSense for content pays "80% of the revenue, after the advertiser platform takes its fee"; when the buyer is Google Ads the publisher keeps "about 68%". [P] https://support.google.com/adsense/answer/180195?hl=en
- **RPM:** Google publishes nothing. Named secondary data points: a food blogger reports ~$1 page RPM on AdSense vs $10 in month one on Mediavine (same 100k-pageview site) [S] https://www.productiveblogging.com/make-money-with-mediavine-ads/ (updated 19 Jan 2026). A Jan-2026 comparison guide (author interviews, not audited) gives AdSense $3–10, Ezoic $3–18 EPMV, Mediavine $15–40+ [S] https://eastondev.com/blog/en/posts/media/20260110-adsense-alternatives-comparison/. Developer-blog AdSense claims of "$6–15" come from an AdSense-approval vendor and are **unsourced** in any meaningful sense.

## 2. Ezoic

- **Minimum (new sites):** "sites are generally required to have 250,000+ monthly active users", verified via Google Analytics. Sites monetising with Ezoic before **19 Feb 2026** are grandfathered as long as they do not remove Ezoic for more than 7 days. Content must follow AdSense policies (AdSense approval itself is not required). [P] https://support.ezoic.com/kb/article/getting-started-ezoics-requirements
- **Incubator (sub-250k):** "250,000 monthly users" is the standard minimum; Ezoic "selects 20 sites or apps each month"; "There are no exceptions to the audience requirement or the selection process"; applications stay in the monthly review pool, and a site is auto-approved if it crosses the threshold later. [P] https://www.ezoic.com/incubator
- **Levels:** the current Levels/tiers page lists only audience-size tiers: eligibility "250K+ monthly users", Business ">250K monthly users" (dedicated account team, Google AdX access), VIP & Enterprise "1M+ monthly users". No per-level economics are published. [P] https://www.ezoic.com/levels. The original Levels launch post (13 Sep 2021) describes Levels 1–4 + VIP assigned on "revenue and visits", with a 90-day integration rule and a promise to "match revenue, plus 25%" if a site earns less with Ezoic; thresholds were never listed. [P] https://www.ezoic.com/blog/ezoic-levels-everything-you-need-to-know
- **Revenue share:** Ezoic's Terms (last updated 12 Jun 2026) pay publishers "no less than eighty-two percent (82%) … of the Net Revenue" from display ads. [P] https://www.ezoic.com/terms/. The older "ad-funded or 10% share" model is described only by third-party reviews now; Ezoic's current pricing page shows no percentages [P] https://www.ezoic.com/pricing — treat the 10% figure as historical.
- **Geo data from Ezoic:** Ezoic's Ad Revenue Index is a 1–100 relative index "across all verticals" with no public country/niche breakdown [P] https://adrevenueindex.ezoic.com/. Its July 2020 post published country index values (US 87, Germany 60, Canada 59, UK 47, France 42, India 38 on 22 Jul 2020) but states explicitly the values "are not comparable with one another in terms of money earned". [P] https://www.ezoic.com/blog/display-ad-rates-geolocation-2020. Ezoic community staff: "USA, UK, Canada and Australia traditionally have the highest EPMV"; India is a "low EPMV-generating country". [P] https://community.ezoic.com/d/16470-my-epmv-is-very-low-when-compared-to-sites-that-have-similar-trafic
- **RPM (EPMV) reports:** no Ezoic-published niche figures. Secondary: "$5 to $12+" EPMV (May 2024, no source given) [S] https://www.publisher-collective.com/blog/how-much-does-ezoic-pay; individual publishers reporting $11.32 (2023) and $15 falling to $5 [S] via https://adscollab.com/blog/adsense-vs-ezoic-2026-full-comparison-revenue-rpm-which-one-wins (vendor blog).

Practical read: for a greenfield site Ezoic is effectively closed until ~250k monthly users, or a 20-per-month lottery.

## 3. Mediavine (and Journey by Mediavine)

- **Mediavine (main):** "$5,000 in annual ad revenue" replaces the old 50,000-sessions rule; plus original audience-first content, "clean, human, brand-safe traffic", good standing with AdSense/AdX, and a reader experience that supports premium ads. Mediavine also reviews "traffic countries of origin, traffic sources, reader demographics". [P] https://www.mediavine.com/blog/new-mediavine-requirements/ and [P] https://help.mediavine.com/what-does-it-take-to-get-approved-by-mediavine. (The former 50k-sessions rule is confirmed as retired by [S] https://www.productiveblogging.com/make-money-with-mediavine-ads/, Jan 2026.)
- **Journey:** "a minimum of 1,000 sessions from tier-1 countries, such as the US, Canada, UK, and Australia, within a 30-day period" plus content/traffic-quality standards; sites are auto-evaluated for full Mediavine on reaching $5k annual ad revenue. [P] https://www.mediavine.com/introducing-journey-new-ad-management/ (published 14 Mar 2024, updated 29 Sep 2026). Secondary sources say the threshold fell from 10,000 to 1,000 sessions on 15 Jan 2026 [S] https://investors.club/?p=83583 and that Grow must run for 30 days before applying [S] https://bloggingguide.com/journey-by-mediavine-review/ (May 2024) — neither is on Mediavine's current page.
- **Revenue share (sites launched after 1 Jan 2026):** Official 75% (<$100k/yr), Select 80%, Signature 85%, Premiere 90% ($500k+). Journey is not listed on this page. [P] https://help.mediavine.com/revenue-share. Journey "70% share" is stated only by [S] https://bloggingguide.com/journey-by-mediavine-review/ (May 2024) — **unverified against Mediavine**.
- **Named RPM reports (per 1,000 sessions):**
  - Fat Stacks (Jon Dykstra, Mediavine, 18–24 Feb 2022): Tech **$11.88**, Entertainment $14.12, Parenting $23.40, Travel $25.53, Fashion $30.44, Cycling $33.53, Home/Garden $39.10. Author calls tech one of his "two worse niches". [S] https://fatstacksblog.com/?p=20869
  - World Travel Family (Mediavine, Nov 2021, USD): USA $28, Canada $25, Australia $24, UK $21, France $16, UAE $16, Singapore $15, **India $4**, Nepal $0.33. [S] https://worldtravelfamily.com/2021-blogging-income-report/
  - Productive Blogging (food blog): $10 RPM in month one; guidance that UK-majority audiences see $15–25 and US-majority $25–35 initially; "my average RPM from Italy is actually quite low – only $3.80". [S] https://www.productiveblogging.com/make-money-with-mediavine-ads/
  - Reddit r/juststart (Mediavine, travel/outdoors, Sep 2023–Feb 2024): US ≈ $21.58, **India ≈ $0.71** (lowest listed). Page could not be fetched directly; figures are from the search snippet of https://lr.us.psf.lt/r/juststart/comments/1b0kkow/rpm_cpm_by_country — [S], low confidence.

## 4. Raptive (formerly AdThrive)

- **Minimum:** "Our minimum pageview requirement is 25,000 pageviews a month."
  - 25,000–99,999 pageviews: "50% of traffic comes from the United States, the United Kingdom, Canada, New Zealand, and/or Australia"; site build makes ads easy to manage; long-form content on most pages; original content with "meaningful human involvement"; GA set up correctly; domain ≥ 6 months old.
  - 100,000+: "40% of your traffic comes from the United States, the United Kingdom, Canada, New Zealand, and/or Australia" + content/GA/domain-age rules.
  [P] https://help.raptive.com/hc/en-us/articles/360032840891-Who-is-eligible-for-Raptive (article updated 26 Sep 2026; the HTML is Cloudflare-gated, text retrieved via the Zendesk API endpoint https://help.raptive.com/api/v2/help_center/en-us/articles/360032840891.json). GA4 (not UA) is required: [P] https://help.raptive.com/hc/en-us/articles/6681661647515-Applying-to-Raptive-with-Google-Analytics
- **Change history:** cut from 100k to 25k on 16 Oct 2025; the 50k–100k "Rise" pilot was retired. [S] https://www.searchenginejournal.com/raptive-drops-traffic-requirement-by-75-to-25000-views/558780/ (20 Oct 2025) and [S] https://ppc.land/raptive-drops-pageview-requirement-to-25-000-monthly-visits/ (19 Oct 2025)
- **Revenue share:** "75% to creators, 25% to Raptive", NET 45 paid on the 15th. [P] https://raptive.com/tbex/. Raptive's own argument that share % is not the whole story: [P] https://raptive.com/revenue-share-doesnt-tell-the-whole-story/
- **RPM:** Raptive publishes lift claims only ("Average 2x RPM lift"; "+103%" for travel creators who switched in 2025, "Raptive network data and Comscore (Dec 2025)") [P] https://raptive.com/tbex/ — marketing, not benchmarks. No tech-niche Raptive RPM from a named publisher was found; third-party 2026 tables claiming $25–52 for tech are unaudited [S] https://earnifyhub.com/blog/blogging/mediavine-vs-raptive-vs-ezoic-rpm.
- **Important for hastcode:** continental Europe and India do **not** count toward the 50%/40% geo rule. A US/EU/India mix will likely fail eligibility unless US+UK alone clear 40–50%.

## 5. Other networks (brief)

| Network | Stated requirement | Source |
|---|---|---|
| Monumetric (Propel) | "10k – 80K pageviews per month"; WordPress or Blogger; ≥6 ad slots desktop+mobile; 300px sidebar; approved by major ad providers; "Minimum 50% traffic from US, UK, CAN & AUS"; "$99 one-time implementation fee" | [P] https://www.monumetric.com/propel-program/ |
| SHE Media Partner Network | "50k+ page views per month" (case-by-case below); "At least 30% US traffic"; site ≥ 90 days old; exclusivity on mobile+desktop banner inventory; paid 35 days after month end; international partners via PayPal | [P] https://www.shemedia.com/partner-network/faqs (older third-party reviews citing 20k pageviews / 80% US / 70% female are out of date: [S] https://imperfectidealist.com/she-media-review/) |
| Newor Media | Elevate program: "There are no minimum page views to qualify"; English-only; exclusivity; ≥2 ad units/page | [P] https://newormedia.com/node/285 |
| Newor Media (Display / AIVision) | "at least 10,000 monthly pageviews" (Display), "30,000 monthly pageviews or more" (AIVision); "At least 50% of traffic should come from the US, UK, Canada, or Australia" | [P] https://newormedia.com/blog/ad-monetization-onboarding-with-newor-media/ (11 Aug 2025) |
| Playwire | "a minimum of 500,000 page views a month" (websites) or 1,000 DAU (apps); IVT < 7%; HTTPS; ≥3 viewable ads per device; Net 60 | [P] https://www.playwire.com/faq and https://www.playwire.com/hubfs/llms.txt |
| Publift | 500,000 monthly pageviews or $2,000/month ad revenue | [P] https://www.publift.com/blog/the-best-cpm-ad-network-for-publishers |

## 6. Geo differences: US vs Western Europe vs India

No network publishes dollar RPM by country. Best available evidence, all secondary unless noted:

| Geo | Evidence | Approx. multiple of US |
|---|---|---|
| US | Mediavine travel: $28 (Nov 2021) [S] worldtravelfamily; ≈ $21.58 (Sep 2023–Feb 2024) [S] r/juststart | 1.0× |
| UK / CA / AU | $21 / $25 / $24 (Nov 2021) [S] worldtravelfamily; Mediavine lists US, CA, UK, AU as its "tier-1" examples [P] mediavine.com/introducing-journey… | ~0.75–0.9× |
| Continental W. Europe | France $16 (Nov 2021) [S] worldtravelfamily; Italy $3.80 on an English-language food blog [S] productiveblogging (undated comment); Ezoic 2020 index: Germany 60, France 42 vs US 87 but "not comparable… in terms of money earned" [P] ezoic.com | ~0.55–0.6× (France); Italy much lower on English content |
| India | **$4** (Nov 2021) [S] worldtravelfamily; **≈$0.71** (2023–24) [S] r/juststart; Ezoic: "rates for audiences in the US were more than double those for audiences in India" (citing Ezoic 2020 data) [S] publisher-collective.com; Indian SEO agency's AdSense benchmark: ₹40–₹830 RPM, CPC ₹4–₹33 vs ₹83–₹830 international (marketing content) [S] https://upgrowth.in/how-to-calculate-google-adsense-earnings-india/ | **~0.03–0.15×** on premium networks; ">2×" gap is the floor per Ezoic |
| Generic tiering claims | "3–10 times higher" for Tier 1 vs developing markets [S] https://www.aditude.com/blog/rpm-for-publishers (Jul 2025, **unsourced**); "US traffic is 3–5× higher than … India" [S] earnifyhub (Apr 2026, no data linked) | — |

Two structural points for an India-inclusive site: (a) Indian sessions also dilute eligibility for Raptive, Monumetric and Newor (Tier-1 % rules) and for Journey (1,000 Tier-1 sessions); (b) EEA/UK traffic needs a TCF-certified CMP to be personalised at all, and non-consenting EU users fall to limited ads (lower CPM) [P] https://support.google.com/adsense/answer/13554116?hl=en.

## 7. Developer-audience caveats

**Ad-blocker usage (primary data).** Computed from the raw public dataset of the Stack Overflow 2018 Developer Survey (fielded 8–28 Jan 2018, 98,855 rows; download https://info.stackoverflowsolutions.com/rs/719-EMH-566/images/stack-overflow-developer-survey-2018.zip; official results https://insights.stackoverflow.com/survey/2018). [P]

- Q: "Do you have ad-blocking software installed on any computers you use regularly?" (n = 76,055): **Yes 72.3%**, No 22.5%, Not sure 5.1%.
- By respondent country: United States 76.7% (n=17,354), United Kingdom 75.1% (n=5,151), Germany 82.4% (n=5,145), France 86.2% (n=2,012), Netherlands 81.9%, Canada 75.9%, Australia 74.5%, **India 55.6%** (n=9,161).
- Q: "In the past month, have you disabled your ad blocker for any reason…?" (n = 54,898): Yes 70.7%.
- Reasons for disabling (n = 37,745, multi-select): forced by site 58.3%; "wanted to support the website … by viewing their ads" 44.8%; display issues 38.2%; site asked 33.6%.
- Caveats: 2018 data (nothing newer from Stack Overflow found), self-selected Stack Overflow users, "installed on any computer" ≠ active on every pageview, and it says nothing about mobile. The widely quoted "72%" in marketing articles traces to this same survey [S] https://www.b2bmarketing.net/?p=8051.

**General-population comparison.** GWI via DataReportal (Q2 2025, ages 16–64): 29.5% globally, US 32.5%, UK 28.5%, Germany 31.5%; men 25–34 highest at 34.5%. India not listed. [S] https://backlinko.com/ad-blockers-users (updated 19 Mar 2026; the underlying DataReportal slides could not be text-extracted). Blockthrough's 2022 PageFair report measured a 21% average adblock rate across >10B pageviews on 9,453 sites (14 Jul 2022). [P-industry] https://blockthrough.com/blog/2022-pagefair-adblock-report/

Implication: a developer site's monetisable pageviews may be ~40–70% of GA pageviews, and GA itself undercounts blocker users (many blockers also block analytics), so the ratio is worse than the dashboard suggests. Networks quote RPM on **ad-served** sessions; your GA-based RPM will look lower.

**Tech-vertical CPMs.** No publisher-side CPM-by-vertical benchmark covering tech for US/UK/DE/IN was found (searches surfaced only advertiser-side social CPM reports). Named evidence that tech is not a premium display vertical: Fat Stacks' $11.88 tech session RPM vs $25–39 for lifestyle niches on the same Mediavine account (Feb 2022) [S]. Third-party 2026 aggregates (earnifyhub, "450+ blogs", no data published) put tech/SaaS at $12–35 with top-10% $48+, versus finance $15–40 [S] https://earnifyhub.com/blog/blogging/blog-display-ad-rpm-by-niche-2026 — low confidence. Nothing credible was found on "dark-mode" effects on CPM: **unsourced**.

## 8. Network comparison table

| Network | Minimum traffic | Geo requirement | Publisher rev share | Typical RPM evidence (per 1,000) | Source |
|---|---|---|---|---|---|
| Google AdSense | None | None (CMP needed for EEA/UK personalised ads) | 80% after ad-platform fee; ~68% via Google Ads | $1 page RPM (food blog) [S]; $3–10 (interview-based guide) [S]; tech-specific **unsourced** | [P] support.google.com/adsense/answer/9724; /180195; /13554116 |
| Ezoic | 250k monthly users (new sites, from 19 Feb 2026); Incubator: 20 sites/month | None stated | ≥82% of display Net Revenue (Terms, 12 Jun 2026) | $5–12+ EPMV (2024, unsourced blog) [S]; $5–15 individual reports [S] | [P] support.ezoic.com/kb/article/getting-started-ezoics-requirements; ezoic.com/incubator; ezoic.com/terms |
| Journey by Mediavine | 1,000 Tier-1 sessions / 30 days | Sessions must be from Tier-1 (US, CA, UK, AU…) | 70% [S, unverified] | None published | [P] mediavine.com/introducing-journey-new-ad-management/ |
| Mediavine | $5,000/yr ad revenue (no session minimum) | Reviews traffic countries; no % stated | 75% (<$100k/yr) → 90% | Tech $11.88 (2022) [S]; US $28 / UK $21 / FR $16 / IN $4 (2021, travel) [S] | [P] mediavine.com/blog/new-mediavine-requirements/; help.mediavine.com/revenue-share |
| Raptive | 25,000 pageviews/month | 50% US/UK/CA/NZ/AU (25k–99k); 40% (100k+); domain ≥6 months | 75% | Only lift claims ("2x", "+103% travel 2025") [P-marketing]; third-party $25–52 tech [S, unaudited] | [P] help.raptive.com article 360032840891; raptive.com/tbex/ |
| Monumetric (Propel) | 10k–80k pageviews | 50% US/UK/CA/AU | not stated | none found | [P] monumetric.com/propel-program/ |
| SHE Media | 50k+ pageviews (case-by-case below) | 30% US; site ≥90 days | not stated | none found | [P] shemedia.com/partner-network/faqs |
| Newor Media | Elevate: none; Display 10k; AIVision 30k | 50% US/UK/CA/AU (Display/AIVision) | not stated | none found | [P] newormedia.com/node/285; onboarding post |
| Playwire | 500k pageviews/month | "good traffic footprint" in US/UK/CA | rev share, % not stated | none found | [P] playwire.com/faq |

## 9. The $10k/month math

Formula: pageviews needed = $10,000 ÷ blended page RPM × 1,000. "Monetisable" = pageviews on which an ad actually renders (after ad-blocking).

| Scenario | Blended page RPM | Monetisable PV/month for $10k | If 50% of PV are ad-blocked, GA PV needed |
|---|---|---|---|
| AdSense, India-heavy mix | $3 | 3,333,000 | 6.7M |
| AdSense, US/EU-heavy tech | $6 | 1,667,000 | 3.3M |
| Mid-tier (Journey / Monumetric / Newor) | $10 | 1,000,000 | 2.0M |
| Mediavine/Raptive, tech niche, Tier-1 heavy | $15 | 667,000 | 1.3M |
| Mediavine/Raptive, strong US skew, good layout | $20 | 500,000 | 1.0M |
| Best plausible for tech | $25 | 400,000 | 0.8M |

Worked geo blend (premium network, per-geo page RPMs taken from the 2021 Mediavine report scaled by ~0.7 for pages-per-session, i.e. US $20, UK/CA/AU $16, W. Europe $11, India $2):

- Mix 40% US / 30% Europe / 30% India → 0.4×20 + 0.3×11 + 0.3×2 = **$11.9** blended page RPM → 840k monetisable PV/month.
- Mix 25% US / 25% Europe / 50% India → 5 + 2.75 + 1 = **$8.75** → 1.14M monetisable PV/month.
- Apply developer ad-blocking: at 72% blocked (SO 2018 ceiling) only 28% of PV serve ads, so the 40/30/30 mix needs 840k ÷ 0.28 ≈ **3.0M GA pageviews/month**; at a kinder 45% blocked (mobile and non-developer readers), ≈ **1.5M**.
- Using the Fat Stacks tech figure directly ($11.88 session RPM ≈ $8 page RPM at 1.5 pages/session) and 50% blocking: 10,000 ÷ 8 × 1,000 ÷ 0.5 = **2.5M GA pageviews/month**.

Reality check: the only tech-site data point in hand ($11.88) plus a 72% developer ad-block rate implies a display-only $10k/month target needs low-single-digit millions of monthly pageviews, which is top-1% territory for an independent dev blog. Display should be the floor, not the plan; sponsorships, newsletter, or affiliate (hosting/SaaS) carry the tech audience in every income report found.

## 10. Open questions

1. **Tech-niche RPM on Raptive/Mediavine today.** Only one named figure (2022). Need 2025–26 income reports from tech/dev publishers or a direct ask to Mediavine/Raptive sales ("average RPM for technology vertical, US traffic").
2. **Journey revenue share** (70% per a May-2024 review) is absent from Mediavine's current revenue-share page — confirm before modelling.
3. **Ezoic economics**: Terms say ≥82% to publisher; third-party reviews still describe "ad-funded or 10%". Which applies to Incubator sites?
4. **India ad-block rate for general users** (GWI) not found; developer rate (55.6%) is from 2018.
5. **Pages-per-session** assumption (1.5) drives all session→page RPM conversions; measure on the real site.
6. **Effect of EU consent rates** on realised European RPM under TCF — no public data found.
7. "Dark-mode" CPM effect — nothing found; likely a non-factor versus viewability and ad density.
8. Raptive counts GA4 pageviews; whether ad-blocked (and therefore GA-invisible) traffic hurts or helps eligibility is unclear.

## Sources (all accessed 2026-10-20)

Primary (network / Google / raw data)
- https://support.google.com/adsense/answer/9724?hl=en — AdSense eligibility
- https://support.google.com/adsense/answer/48182?hl=en — AdSense Program policies
- https://support.google.com/publisherpolicies/answer/10502938?hl=en — Google Publisher Policies
- https://support.google.com/adsense/answer/180195?hl=en — AdSense revenue share
- https://support.google.com/adsense/answer/13554116?hl=en — EEA/UK/CH CMP requirement
- https://support.ezoic.com/kb/article/getting-started-ezoics-requirements — Ezoic 250k requirement, 19 Feb 2026 grandfathering
- https://www.ezoic.com/incubator — Ezoic Incubator
- https://www.ezoic.com/levels — Ezoic tiers
- https://www.ezoic.com/pricing — Ezoic tiers/pricing page (no percentages)
- https://www.ezoic.com/blog/ezoic-levels-everything-you-need-to-know — Levels launch, 13 Sep 2021
- https://www.ezoic.com/terms/ — Ezoic Terms, last updated 12 Jun 2026 (82% floor)
- https://www.ezoic.com/blog/display-ad-rates-geolocation-2020 — Ezoic geo index, 24 Jul 2020
- https://adrevenueindex.ezoic.com/ — Ezoic Ad Revenue Index
- https://community.ezoic.com/d/16470-my-epmv-is-very-low-when-compared-to-sites-that-have-similar-trafic — Ezoic staff on country value
- https://www.mediavine.com/blog/new-mediavine-requirements/ — Mediavine + Journey requirements (2026)
- https://help.mediavine.com/what-does-it-take-to-get-approved-by-mediavine — Mediavine approval criteria
- https://help.mediavine.com/revenue-share — Mediavine revenue share (sites launched after 1 Jan 2026)
- https://www.mediavine.com/introducing-journey-new-ad-management/ — Journey (14 Mar 2024, updated 29 Sep 2026)
- https://help.raptive.com/hc/en-us/articles/360032840891-Who-is-eligible-for-Raptive — Raptive eligibility (updated 26 Sep 2026; read via https://help.raptive.com/api/v2/help_center/en-us/articles/360032840891.json)
- https://help.raptive.com/hc/en-us/articles/6681661647515-Applying-to-Raptive-with-Google-Analytics — GA4 requirement
- https://raptive.com/tbex/ — Raptive 75% share, NET 45, lift claims (Dec 2025 data)
- https://raptive.com/revenue-share-doesnt-tell-the-whole-story/ — Raptive on rev share
- https://www.monumetric.com/propel-program/ — Monumetric Propel
- https://www.shemedia.com/partner-network/faqs — SHE Media FAQ
- https://newormedia.com/node/285 — Newor Media Elevate FAQ
- https://newormedia.com/blog/ad-monetization-onboarding-with-newor-media/ — Newor Display/AIVision thresholds (11 Aug 2025)
- https://www.playwire.com/faq and https://www.playwire.com/hubfs/llms.txt — Playwire requirements
- https://www.publift.com/blog/the-best-cpm-ad-network-for-publishers — Publift threshold
- https://info.stackoverflowsolutions.com/rs/719-EMH-566/images/stack-overflow-developer-survey-2018.zip — Stack Overflow 2018 raw data (percentages computed locally; official results page https://insights.stackoverflow.com/survey/2018)
- https://blockthrough.com/blog/2022-pagefair-adblock-report/ — PageFair/Blockthrough 2022 (14 Jul 2022)

Secondary (publisher reports, trade press, blogs)
- https://worldtravelfamily.com/2021-blogging-income-report/ — Mediavine RPM by country, Nov 2021
- https://fatstacksblog.com/?p=20869 — Mediavine RPM by niche incl. Tech $11.88, Feb 2022 (published 5 Sep 2022)
- https://www.productiveblogging.com/make-money-with-mediavine-ads/ — AdSense vs Mediavine RPM, UK/US/Italy notes (updated 19 Jan 2026)
- https://lr.us.psf.lt/r/juststart/comments/1b0kkow/rpm_cpm_by_country — r/juststart Mediavine RPM by country (Feb 2024; snippet only)
- https://www.searchenginejournal.com/raptive-drops-traffic-requirement-by-75-to-25000-views/558780/ — Raptive 25k change (20 Oct 2025)
- https://ppc.land/raptive-drops-pageview-requirement-to-25-000-monthly-visits/ — Raptive 25k change (19 Oct 2025)
- https://investors.club/?p=83583 — Journey 10k→1k sessions (15 Jan 2026 claim)
- https://bloggingguide.com/journey-by-mediavine-review/ — Journey 70% share, Grow 30-day rule (8 May 2024)
- https://www.publisher-collective.com/blog/how-much-does-ezoic-pay — Ezoic EPMV $5–12+, US >2× India (7 May 2024)
- https://adscollab.com/blog/adsense-vs-ezoic-2026-full-comparison-revenue-rpm-which-one-wins — individual Ezoic EPMV anecdotes (vendor blog)
- https://eastondev.com/blog/en/posts/media/20260110-adsense-alternatives-comparison/ — AdSense/Ezoic/Mediavine ranges (10 Jan 2026, interview-based)
- https://earnifyhub.com/blog/blogging/blog-display-ad-rpm-by-niche-2026 — tech/SaaS $12–35 (2 Apr 2026; no data linked)
- https://earnifyhub.com/blog/blogging/mediavine-vs-raptive-vs-ezoic-rpm — network RPM table (unaudited)
- https://www.aditude.com/blog/rpm-for-publishers — Tier-1 "3–10×" claim (Jul 2025; unsourced)
- https://upgrowth.in/how-to-calculate-google-adsense-earnings-india/ — Indian AdSense RPM/CPC estimates (agency marketing)
- https://backlinko.com/ad-blockers-users — GWI/DataReportal ad-block rates by country (updated 19 Mar 2026)
- https://www.b2bmarketing.net/?p=8051 — "72% of developers use ad blockers" (cites SO 2018)
- https://imperfectidealist.com/she-media-review/ — older SHE Media criteria (2020, superseded)

# Google News, Google Discover, Bing News and Apple News: inclusion and technical requirements

Research for hastcode.com (tech news + architecture explainers; Go SSR; mobile-first; US / Europe / India).
Researched 2026-10-20. All primary-source quotes were pulled from the live pages on that date; "last updated" is the date printed on the page where one is shown. Secondary sources are marked **[secondary]**. Opinions are marked **[opinion]**.

---

## TL;DR

- **Google News does not take submissions any more.** Sites are "automatically considered" for Google News, Top stories and the News tab; the application process was retired in Dec 2019, Publisher Center stopped accepting new Google News publications on 2024-04-25, and since late March 2025 all publication pages are auto-generated and Publisher Center RSS/URL submissions are ignored. Publisher Center now only matters for News Showcase / Reader Revenue Manager. (Sources: Publisher Center answers 9607025, 14859160, 15898024; Search Central blog 2021-07-16.)
- **Eligibility is editorial, not technical:** follow the Google News content policies (10 content categories + ads/sponsorship, misleading previews, transparency: "Clear dates and bylines", author/publisher info, company info, contact information), and build "a consistent history of producing original news-related content" with E-E-A-T. Ranking factors Google names: relevance, prominence, authoritativeness, freshness, usability, location and language.
- **Discover is automatic too**: any indexed page that meets the Discover content policies is eligible; "No special tags or structured data are required." The one concrete technical lever is large images: >= 1200 px wide, > 300,000 px total, enabled by `max-image-preview:large` (or AMP). Google explicitly says Discover traffic is "less predictable or dependable" than Search and should be treated as supplemental.
- **Technical minimum set for the site:** news sitemap (<= 1,000 `news:news` entries, articles from the last 2 days only, `publication.name`, `publication.language`, `publication_date`, `title`), `NewsArticle` JSON-LD with `headline`, `datePublished`/`dateModified` (ISO 8601 with timezone), `author` (name + url), `image` (16x9 / 4x3 / 1x1, multiple sizes), `<meta name="robots" content="max-image-preview:large">`, `og:image`, HTTPS, good Core Web Vitals (LCP <= 2.5 s, INP < 200 ms, CLS < 0.1), and `hreflang` only if you actually ship region/language variants.
- **AMP has not been required for Top stories since the May 2021 page-experience rollout** ("Any page that meets the Google News content policies will be eligible", Search Central blog 2020-11-10). Google now "applies the same standard to all pages, regardless of the technology used to build the page."
- **AI-assisted writing is allowed** ("Appropriate use of AI or automation is not against our guidelines", 2023-02-08) but "scaled content abuse" (many pages generated primarily to manipulate rankings, "whether content is made by automation, humans, or a mix") is a spam policy since 2024-03-05. Google asks that automation be "self-evident to visitors through disclosures or in other ways" where readers would reasonably expect it, and says fabricated author profiles are "a form of deception". Apple News *requires* AI-generated articles to be marked as such.
- **Bing News (PubHub inside Bing Webmaster Tools) and Apple News Publisher are still submission/review-based.** Bing judges newsworthiness, originality, authority, readability. Apple News Publisher is for "professional publications of a journalistic nature", excludes personal blogs and business-promotion sites, is limited to publishers based in the US, UK, Canada and Australia, and supports English and Canadian French only, so it is likely out of reach for an India-based brand-new tech blog (see section 8).

---

## Technical checklist

| # | Requirement | Where it lives on hastcode.com | Source |
|---|-------------|--------------------------------|--------|
| 1 | Do **not** block `Googlebot-News` in robots.txt ("blocking Googlebot-News excludes your content from Google News") | `/robots.txt` | https://support.google.com/news/publisher-center/answer/9607025 |
| 2 | Visible byline, visible publish date (and updated date if edited) on every article | Article template | https://support.google.com/news/publisher-center/answer/6204050 (Transparency) |
| 3 | Author pages (bio, credentials, link target for bylines); no fabricated personas | `/authors/{slug}` | https://developers.google.com/search/docs/fundamentals/creating-helpful-content (last updated 2026-10-05) |
| 4 | About page naming the publisher and "the company or network behind the content"; Contact page with real contact info | `/about`, `/contact`, site footer | https://support.google.com/news/publisher-center/answer/6204050 |
| 5 | Sponsored content clearly labelled; ads must not "exceed your content" | Article template / ad slots | https://support.google.com/news/publisher-center/answer/6204050 (Ads & sponsored content) |
| 6 | News sitemap: `xmlns:news="http://www.google.com/schemas/sitemap-news/0.9"`, one `<news:news>` per URL, `<news:publication><news:name>` (must match the name shown on news.google.com) + `<news:language>` (ISO 639), `<news:publication_date>` (W3C datetime, original first-publish time), `<news:title>`; max 1,000 entries; only articles created in the last 2 days; drop old URLs or their `news:` block after 2 days; keep a standalone file for better Search Console tracking | `/news-sitemap.xml` generated by the Go server; reference it from robots.txt and submit in Search Console | https://developers.google.com/search/docs/crawling-indexing/sitemaps/news-sitemap (last updated 2025-12-10) |
| 7 | Regular XML sitemap for all content (explainers are not "news") | `/sitemap.xml` | https://developers.google.com/search/docs/crawling-indexing/sitemaps/overview (general sitemap guidelines referenced by the news-sitemap doc) |
| 8 | JSON-LD `NewsArticle` (news) / `Article` or `BlogPosting` (explainers): `headline`, `image[]`, `datePublished`, `dateModified` (ISO 8601 **with timezone**), `author` (`Person` with `name` + `url`), per author-markup best practices. No property is strictly required. | `<head>` of article template | https://developers.google.com/search/docs/appearance/structured-data/article (last updated 2026-09-08) |
| 9 | Article images: crawlable, indexable, Google-Images-supported format, multiple high-res variants in 16x9, 4x3 and 1x1 (min 50K pixels w*h) | Image pipeline / CDN | same as #8 |
| 10 | Discover-grade hero image: >= 1200 px wide, > 300,000 total pixels (e.g. 1280x720, 16x9); not a logo, not text-heavy; declared via `og:image` and/or schema.org `image` | Hero image + `og:image` meta | https://developers.google.com/search/docs/appearance/google-discover (last updated 2026-03-09) |
| 11 | `<meta name="robots" content="max-image-preview:large">` on article pages (needed to unlock the large Discover preview if you are not on AMP) | `<head>` sitewide | https://developers.google.com/search/docs/crawling-indexing/robots-meta-tag (last updated 2026-03-24) and Discover doc above |
| 12 | HTTPS everywhere; mobile-friendly rendering; good Core Web Vitals: LCP within 2.5 s, INP < 200 ms, CLS < 0.1 | Go SSR + CDN; no layout shift from ads/images (set width/height) | https://developers.google.com/search/docs/appearance/page-experience (2026-09-22); https://developers.google.com/search/docs/appearance/core-web-vitals (2025-12-10) |
| 13 | AMP: **not required** for any news surface | n/a (skip AMP) | https://developers.google.com/search/blog/2020/11/timing-for-page-experience (2020-11-10); https://developers.google.com/search/docs/crawling-indexing/amp/about-amp (2026-07-01) |
| 14 | One language per page; do not IP-redirect or IP-adapt content; if you later add regional/language variants, use `hreflang` (HTML `<link>`, HTTP header or sitemap `xhtml:link`) with return links on every variant plus `x-default` | Only if/when regional editions exist | https://developers.google.com/search/docs/specialty/international/localized-versions (2026-09-21); https://developers.google.com/search/docs/specialty/international/managing-multi-regional-sites (2025-12-10) |
| 15 | Site name and favicon are what Google News now uses as publication title and logo; set them via Search (WebSite structured data / favicon), not Publisher Center | `<link rel="icon">`, `WebSite` JSON-LD on home | https://support.google.com/news/publisher-center/answer/15898024 |
| 16 | Optional: "Add as preferred source" button (`news.google.com/swg/js/v1/publisher.js` + `<div google-add-preferred-source-btn>`) or deeplink `https://www.google.com/preferences/source?q=hastcode.com`; only domain/subdomain-level sites are eligible | Header / footer / end of article | https://developers.google.com/search/docs/appearance/preferred-sources (2026-09-18) |
| 17 | Verify the property in Google Search Console (Search performance filtered by "News" type = News tab; separate "Google News" report = news.google.com + apps; Discover report) and in Bing Webmaster Tools (prerequisite for PubHub) | DNS / HTML verification | https://support.google.com/news/publisher-center/answer/9607025; https://support.microsoft.com/en-us/bing/how-to-submit-your-news-website-to-bing |
| 18 | AI disclosure where readers would reasonably expect it; never give AI an author byline | Editorial policy page + per-article note | https://developers.google.com/search/blog/2023/02/google-search-and-ai-content; https://developers.google.com/search/docs/fundamentals/creating-helpful-content |
| 19 | Valid RSS 2.0 / Atom feed (UTF-8, HTTPS with valid cert, channel `title`/`language`/`link`, item `title`/`link`/`description`) - needed for Apple News RSS ingestion and useful generally | `/feed.xml` | https://support.apple.com/guide/news-publisher/apdc2c7520ff/icloud |

---

## 1. Google News: automatic inclusion vs Publisher Center

**Current stance (primary sources):**

- Publisher Center, "News content across Google" (https://support.google.com/news/publisher-center/answer/9607025, no last-updated date shown): "Google may find your site through our normal web crawl." "Publishers are automatically considered for 'Top stories' or the News tab of Search." "Having eligible content isn't a guarantee that content will necessarily rank well on Google News." Blocking `Googlebot-News` "excludes your content from Google News". The Google News performance report in Search Console is described as the best way to check whether eligible content is appearing.
- Search Central blog, "Answers to some common questions about appearing in Google News" (https://developers.google.com/search/blog/2021/07/google-news-top-questions, 2021-07-16, page carries an "information may be outdated" notice): "Your site is automatically considered for Google News and news surfaces in Search - no application required." The old application process was retired when Publisher Center launched in December 2019. Publisher Center is described as an optional tool for managing site details; "being approved to use the Publisher Center" or submitting content through it does not guarantee appearance.
- Publisher Center, "An update on publication pages" (https://support.google.com/news/publisher-center/answer/14859160, published 2024-04-25, updated 2024-12-18): publishers "can no longer set up Google News publications in Publisher Center" effective 2024-04-25. Having a publication in Publisher Center does not guarantee its content will appear across Google surfaces. Content from publishers that follows content policies is automatically eligible; ranking depends on "relevance, prominence, authoritativeness, freshness, location, and language."
- Publisher Center, "Google News transitions to automatically-generated publication pages" (https://support.google.com/news/publisher-center/answer/15898024, published 2025-02-10, updated 2025-03-20): "Google News is fully transitioning to automatically generated publication pages in late March 2025." "Google News will no longer use RSS feeds or web locations that were submitted in Publisher Center." Custom sections removed; the Google News tile removed from Publisher Center; logos now come from site favicons and titles from site names, set "through Google Search"; country restrictions via Publisher Center discontinued except for News Showcase panels; YouTube submissions discontinued. Still in Publisher Center: News Showcase country controls and logos, Reader Revenue Manager logos. "This doesn't change what content is eligible to appear in Google News and other news surfaces."
- Announcements index: https://support.google.com/news/publisher-center/announcements/10146168.

**Timeline, as documented by Google:**

| Date | Change | Source |
|------|--------|--------|
| Dec 2019 | New Publisher Center launched; the Google News application process retired; sites automatically considered | 2021-07-16 Search Central blog (Google's own retrospective). The original Dec 2019 blog.google post was not located in this research; coverage at ppc.land and searchenginejournal.com **[secondary]** matches. |
| 2021-01-12 | Google News performance report added to Search Console | Announcements page |
| 2024-04-25 | Publisher Center stops accepting new Google News publications | answer/14859160 |
| 2025-02-10 (announced), late March 2025 (completed) | All publication pages auto-generated; Publisher Center RSS / URL submissions, custom sections, country restrictions, custom logos dropped | answer/15898024 |

**What this means for hastcode.com:** there is nothing to submit. Do not plan work around Publisher Center unless you later join News Showcase or Reader Revenue Manager. The levers are content policy compliance, originality/consistency, transparency pages, and the technical items in the checklist.

**Google's three stated eligibility criteria** (2021-07-16 blog): sites that (a) have "high levels of expertise, authority, and trustworthiness", (b) have "a consistent history of producing original news-related content", and (c) comply with Google's News policies. The same post says these criteria apply to Google News, Top stories and the News tab and do not affect regular web results, and that "automated systems run regularly and will detect when a site meets the criteria."

## 2. Google News content policies and transparency requirements

Source: https://support.google.com/news/publisher-center/answer/6204050 (no last-updated date shown; references the site reputation abuse policy effective 2024-05-05, so it has been edited at least since then). Google Search's overall policies and spam policies apply in addition.

**Content categories that disqualify content from news surfaces:** Dangerous content; Deceptive practices; Harassing content; Hateful content; Manipulated media; Medical content; Sexually explicit content; Violent extremist content; Violent & gory content; Vulgar language & profanity.

**Feature-specific policies:**
- Ads & sponsored content: "Advertising and other paid promotional material on your pages shouldn't exceed your content." Sponsorship must be clearly disclosed; sponsored content must not be presented as independent editorial (links to the site reputation abuse spam policy).
- Misleading content: preview content (title, snippet, image) must not promise what the page does not deliver.
- Transparency (verbatim): "Visitors to your site want to trust and understand who writes and publishes the content they read. That's why news sources on Google should provide: Clear dates and bylines; Information about the authors, publication, and publisher; Information about the company or network behind the content; Contact information."

**Enforcement:** violating content may be removed from news surfaces; repeated or egregious violations may make a site ineligible. Exceptions for artistic, educational, historical, documentary, scientific purposes or substantial public benefit.

**Related Search guidance that backs the transparency items** (https://developers.google.com/search/docs/fundamentals/creating-helpful-content, last updated 2026-10-05): "Is it self-evident to your visitors who authored your content?" "Do bylines lead to further information about the author or authors involved[?]" "We strongly encourage adding accurate authorship information" where readers might expect it. "Fabricating creator profiles" with AI headshots, invented names or false credentials "is a form of deception." "E-E-A-T itself isn't a specific ranking factor"; "Of these aspects, trust is most important."

**Google News ranking factors** (https://support.google.com/news/publisher-center/answer/9606702, no date): "Relevance of content", "Prominence", "Authoritativeness", "Freshness", "Usability", "Location and language". "It doesn't accept payments to expedite or improve a site's search appearance or ranking." AdSense participation "has no effect on ranking."

## 3. Google Discover

Source: https://developers.google.com/search/docs/appearance/google-discover (last updated 2026-03-09). Content policies: https://support.google.com/websearch/answer/9982767 (no date shown).

- Eligibility: content is automatically eligible if it is indexed and meets Discover's content policies. "No special tags or structured data are required." "Note that being eligible to appear in Discover is not a guarantee of appearing."
- "As part of Google Search, Discover makes use of many of the same signals and systems used by Search"; Google points to its "creating helpful, reliable, people-first content" guidance (which is where E-E-A-T lives; the Discover doc itself does not use the acronym).
- Images: "At least 1200 px wide", "High resolution of more than 300,000 total pixels", 16x9 example (1280x720). Large previews are "Enabled by the max-image-preview:large setting, or by using AMP." Use schema.org markup or `og:image` to name the preferred image. Avoid site logos and text-heavy images.
- Content guidance: Discover favours content that is "timely for current interests, tells a story well, or provides unique insights." "Avoid clickbait and similar tactics to artificially inflate engagement" via misleading or exaggerated titles, snippets or images. "Avoid sensationalism tactics that manipulate appeal by catering to morbid curiosity, titillation, or outrage." Titles should "capture the essence of the content." Discover may exclude job applications, petitions, forms, code repositories and satire without context; SafeSearch and a "shocking or unexpected" filter apply.
- Traffic realism (Google's words): Discover traffic is "less predictable or dependable" than keyword-driven search and should be treated as supplemental; it fluctuates with changing interests, content type adjustments and Search updates.
- Where Discover surfaces (https://support.google.com/websearch/answer/2819496): in the Google app, "On your Android phone or tablet, in your browser on google.com", and on some devices by swiping right from the home screen. It is a mobile surface; Google's doc does not list a separate "mobile requirement" beyond a strong page experience, but the audience is mobile by construction, so the mobile-first SSR decision is well aligned.
- Discover content policies mirror the Google News list (same 10 categories, plus Ads & sponsored content, Misleading content, Transparency). Violations can trigger Discover manual actions in Search Console. The "Following" feed uses a slightly different policy list.

## 4. Technical requirements

### 4.1 News sitemap
Spec: https://developers.google.com/search/docs/crawling-indexing/sitemaps/news-sitemap (last updated 2025-12-10). Namespace `http://www.google.com/schemas/sitemap-news/0.9`.
- Tags: `<news:news>` (one per `<url>`), `<news:publication>` with `<news:name>` (must "exactly match the name as it appears on your articles on" news.google.com, excluding anything in parentheses) and `<news:language>` (ISO 639; `zh-cn` / `zh-tw` special cases), `<news:publication_date>` (W3C format; "Specify the original date and time when the article was first published on your site"), `<news:title>` (as on the site; no author/publication name/date in it).
- Limits: "up to 1,000 `news:news` tags" per sitemap, split otherwise. "Only include recent URLs for articles that were created in the last two days." After two days, remove the URL or its `news:` metadata. Add fresh articles to the existing file rather than creating a new file per update. Google News crawls news sitemaps "as often as the rest of the site". An empty file may raise an "Empty Sitemap" warning but causes no Search problems. Standalone news sitemap recommended for Search Console tracking.
- Note: the old Publisher Center "news sitemap" help URL (answer/9606710) now 301s to this Search Central page, confirming it is the canonical spec.

### 4.2 Article / NewsArticle structured data
Source: https://developers.google.com/search/docs/appearance/structured-data/article (last updated 2026-09-08).
- Types: `Article`, `NewsArticle`, `BlogPosting`. "There are no required properties; instead, add the properties that apply to your content."
- Recommended: `author` (Person/Organization; follow the author-markup best practices section), `author.name`, `author.url` (or `sameAs`), `dateModified` and `datePublished` (ISO 8601; "We recommend that you provide timezone information; otherwise, we will default to the timezone used by Googlebot"), `headline` ("long titles may be truncated on some devices"), `image`.
- Image rules (verbatim): "Image URLs must be crawlable and indexable." "Images must represent the marked up content." "Images must be in a file format that's supported by Google Images." "For best results, we recommend providing multiple high-resolution images (minimum of 50K pixels when multiplying width and height) with the following aspect ratios: 16x9, 4x3, and 1x1."
- Technical guidelines: for multi-part articles point `rel=canonical` at each page or a view-all page, never page 1; use paywall structured data if gating content. No AMP-specific requirements remain on this page.
- Google says structured data is not required to appear in Top stories (2021 Q&A does not mention it; coverage of the current Article doc at searchenginejournal.com **[secondary]** states markup "isn't needed to rank in Top Stories, though Article markup helps Google understand a page").

### 4.3 max-image-preview
Source: https://developers.google.com/search/docs/crawling-indexing/robots-meta-tag (last updated 2026-03-24). Values: `none`, `standard`, `large` ("A larger image preview, up to the width of the viewport, may be shown"). Applies "to all forms of search results (such as Google web search, Google Images, Discover, Assistant)". If you do not want large thumbnails for AMP or canonical article versions in Search or Discover, use `standard` or `none`. For hastcode.com: ship `max-image-preview:large` sitewide.

### 4.4 AMP status
- Search Central blog, "Timing for bringing page experience to Google Search" (https://developers.google.com/search/blog/2020/11/timing-for-page-experience, 2020-11-10): the Top stories change "will also roll out in May 2021"; "Any page that meets the Google News content policies will be eligible", and Google will "prioritize pages with great page experience, whether implemented using AMP or any other web technology, as we rank the results."
- About AMP (https://developers.google.com/search/docs/crawling-indexing/amp/about-amp, last updated 2026-07-01): Google "applies the same standard to all pages, regardless of the technology used to build the page"; "AMP on desktop doesn't get search-specific features". The page does not list any feature that requires AMP.
- Conclusion: AMP is not required for Top stories, the News tab, Google News or Discover. The only residual AMP mention relevant here is that AMP is one of two ways to unlock large Discover image previews; `max-image-preview:large` is the other.

### 4.5 Core Web Vitals / page experience
- Page experience doc (https://developers.google.com/search/docs/appearance/page-experience, last updated 2026-09-22): "Core Web Vitals are used by our ranking systems." "There is no single signal." Other page-experience aspects "don't directly help your website rank higher in search results" but good page experience "can contribute to success in Search". Self-assessment covers Core Web Vitals, HTTPS, mobile display, intrusive interstitials.
- Core Web Vitals doc (https://developers.google.com/search/docs/appearance/core-web-vitals, last updated 2025-12-10): LCP within 2.5 s, INP under 200 ms, CLS under 0.1. "We highly recommend site owners achieve good Core Web Vitals for success with Search."
- Top stories specifically: the current docs do not carve out a separate page-experience rule for Top stories; the 2020-11-10 blog post is the statement that page experience is used in ranking Top stories results. Treat CWV as a ranking input, not an eligibility gate.

## 5. Top stories carousel vs News tab vs Google News app

- Google defines the three surfaces in the 2021-07-16 blog: Google News app/site (iOS, Android, news.google.com with "For you" and Headlines); Top stories ("A carousel that often appears at the top of Search results for timely topics"); News tab ("In Search, it narrows results to news-related content"). Eligibility criteria (section 1) "apply to Google News, Top stories, and the News tab" and are determined automatically.
- Publisher Center answer/9607025: "Publishers are automatically considered for 'Top stories' or the News tab of Search"; Top stories content is "automatically selected for this feature."
- Measurement differs: the Search Console "Google News" performance report covers the apps and news.google.com; the News tab shows up in the Search performance report filtered by search type = News (launch coverage: https://searchengineland.com/google-news-performance-report-added-to-google-search-console-345332 **[secondary]**).
- Google's Search Central visual elements gallery (https://developers.google.com/search/docs/appearance/visual-elements-gallery, last updated 2026-02-04) does not define Top stories as a separate element, so there is no separate markup or opt-in doc for it.
- Preferred sources (https://developers.google.com/search/docs/appearance/preferred-sources, last updated 2026-09-18): users can pick your site; then "your content is more likely to appear in" Top stories with a "preferred" badge, and in AI Mode / AI Overviews for those users. Available "globally for the 'Top Stories' feature" in all Search languages. "Only domain-level and subdomain-level sites are eligible." The on-site button/SDK/deeplink is optional: "It's not required to do them in order to appear as a preferred source." Launch post: https://blog.google/products/search/preferred-sources/ (2025-08-12).
- **[secondary]** Search Engine Land (2019, https://searchengineland.com/google-all-sites-are-eligible-to-be-in-google-news-but-not-all-content-will-appear-in-google-news-350521, returned 403 on fetch; from search snippet only) quotes Google saying Google News and the News tab follow "2 different algorithms". Treat as unverified.

## 6. Geo / language: US, Europe, India

- Google News ranking includes "Location and language" as a factor (answer/9606702) but does not say how. Since late March 2025 country restrictions can no longer be set in Publisher Center (answer/15898024), so there is no Publisher Center region setting for a normal publisher to configure: Google decides geographic surfacing from the content itself.
- Language detection (https://developers.google.com/search/docs/specialty/international/managing-multi-regional-sites, last updated 2025-12-10): "Google uses the visible content of your page to determine its language." "We don't use any code-level language information such as lang attributes, or the URL." Keep one language per page. "Don't use IP analysis to adapt your content." "don't redirect based on what you think the user's language may be." Geotargeting "can improve your page rankings in the target country, but at the expense of results in other locales or languages." URL-structure options compared: ccTLD (single country), subdomain, subdirectory ("Low maintenance (same host)"), URL parameters ("Not recommended").
- hreflang (https://developers.google.com/search/docs/specialty/international/localized-versions, last updated 2026-09-21): three equivalent methods (HTML `<link rel="alternate" hreflang>`, HTTP `Link` header, sitemap `<xhtml:link>`); every version lists itself and all others; "If two pages don't both point to each other, the tags will be ignored"; ISO 639-1 language, optional ISO 3166-1 alpha-2 region; "You can't specify the country code by itself"; `UK`/`EU` are invalid (use `GB`); `x-default` "is used when no other language/region matches the user's browser setting". "Localized versions of a page are only considered duplicates if the main content of the page remains untranslated."
- News sitemap language: one `<news:language>` per article (ISO 639), so a single-language English site declares `en` everywhere.
- Practical reading **[opinion]**: a single English `.com` with no regional editions needs no hreflang at all and should not geotarget in Search Console. If you later ship e.g. `/in/` edition pages or Hindi/German translations, add hreflang pairs with return links and `x-default`; keep each page monolingual. For Google News the practical regional lever is content: Indian-angle stories (datelines, local sources, INR pricing) will be the thing that earns India surfacing, not metadata.

## 7. AI-generated content: Google's policies

- "Google Search's guidance about AI-generated content" (https://developers.google.com/search/blog/2023/02/google-search-and-ai-content, 2023-02-08): "Appropriate use of AI or automation is not against our guidelines." "Our focus on the quality of content, rather than how content is produced, is a useful guide." Content made "with the primary purpose of manipulating ranking in search results" is spam however produced. "Using AI doesn't give content any special gains." AI/automation disclosures "are useful for content where someone might think 'How was this created?'." Bylines should be used "when readers would reasonably expect it"; Google News publishers specifically should use bylines and author information; giving AI an author byline "is probably not the best way" to disclose.
- Spam policies (https://developers.google.com/search/docs/essentials/spam-policies, last updated 2026-08-28): "Scaled content abuse is when many pages are generated for the primary purpose of manipulating search rankings"; examples include "generative AI tools or other similar tools to generate many pages without adding value for users", scraping feeds/search results, stitching content, keyword-stuffed pages. Site reputation policy (third-party content on a host "mainly because of that host's already-established ranking signals"; outside the EEA a manual action may apply, inside the EEA such pages may be ranked separately). Expired domain abuse. The intro also lists "attempting to manipulate generative AI responses in Google Search" as spam.
- March 2024 announcement (https://developers.google.com/search/blog/2024/03/core-update-spam-policies, 2024-03-05): three new policies, "expired domain abuse, scaled content abuse, and site reputation abuse"; scaled content abuse applies whether content is "made by automation, humans, or a mix of both"; site reputation abuse "will take effect starting May 5, 2024". (The "40% less unhelpful content" figure is in the companion Keyword blog post, not this Search Central post; not verified here.)
- Helpful-content guidance (https://developers.google.com/search/docs/fundamentals/creating-helpful-content, 2026-10-05), "How" section: readers benefit from knowing "how a piece of content was produced", "That can include automated, AI-generated, and AI-assisted content." "Is the use of automation, including AI-generation, self-evident to visitors through disclosures or in other ways?" Consider explaining "how automation or AI-generation was used to create content"; add disclosures "when it would be reasonably expected." Using automation or AI primarily to manipulate rankings violates the spam policies.
- Apple News (https://support.apple.com/guide/news-publisher/ai-generated-content-apd8da39c6f2/1.0/icloud/1.0): stricter. Articles created "by or with the assistance of artificial intelligence (AI)" "must be marked as AI-generated in News Publisher or with the corresponding Apple News Format metadata property" and carry a labelled byline/co-byline. Using AI "as a tool in idea generation, headline generation, research, editing, or reporting" does not require disclosure when a journalist writes and verifies the piece. "Using AI to mislead readers is not allowed" and "could lead to your channel's suspension." "You're responsible for the accuracy of all content published to Apple News."

**Implications for an AI-assisted workflow [opinion]:** human-authored bylines with real author pages; an editorial/AI-use policy page; per-article disclosure when drafting was materially AI-assisted; cap volume to what humans actually review (a 50-articles-a-day pipeline on a 1-person site is exactly the "scaled content abuse" pattern); original reporting or analysis per piece rather than rewrites of other outlets' stories (which also fails Google's "original news-related content" criterion and Bing's "originality").

## 8. Bing News PubHub and Apple News Publisher (requirements only)

### Bing News PubHub
- Microsoft support, "How to submit your news website to Bing" (https://support.microsoft.com/en-us/bing/how-to-submit-your-news-website-to-bing, no date): (1) "Make sure that your site follows the Bing Webmaster Guidelines", (2) "Verify the ownership of your site using the Bing Webmaster Tools", (3) "Submit your website to the Bing News PubHub for evaluation". Links to Bing Publisher Guidelines (pubhub.bing.com/Home/Help) which "explain how we decide which sites to index as news in Bing". Note: `pubhub.bing.com` did not resolve on 2026-10-20; PubHub lives inside Bing Webmaster Tools' left navigation (https://searchengineland.com/bing-news-pubhub-moves-to-bing-webmaster-tools-380070 **[secondary]**).
- Bing Webmaster blog, "Reach millions of people with Bing News PubHub" (https://blogs.bing.com/webmaster/february-2022/Reach-millions-of-people-with-Bing-News-PubHub, 2022-02-08): criteria are Newsworthiness ("Reporting on timely events and topics that are interesting to the people using our services"; how-to, job postings, advice columns, product promotion, raw weather/stock data do not qualify), Originality ("unique facts or points of view"), Authority ("Submissions must identify sources, authors, and attribution of all content"), Readability (correct grammar/spelling, easy navigation, "Advertising should never interfere with the site experience"). Steps: verify in Bing Webmaster Tools, complete PubHub details, then access to the Bing Publisher Network. No review SLA and no RSS/sitemap requirement is stated. Original 2016 launch post with the same four criteria: https://blogs.bing.com/search/june-2016/Bing-helps-publishers-expand-their-reach.
- No official Bing page found dated 2025-2026; treat the criteria as current-but-unconfirmed.

### Apple News Publisher
- "Publishing on Apple News" (https://support.apple.com/guide/news-publisher/publishing-on-apple-news-apde42330c66/1.0/icloud/1.0): Apple News is "a distribution platform for professional publications of a journalistic nature." Not intended: publications that are not editorial/journalistic, personal blogs or portfolios, those primarily promoting a business/service/nonprofit, those primarily reposting others' content, those publishing factual inaccuracies or failing "widely accepted journalistic standards", and publications "Based outside countries where Apple News is available (Australia, Canada, the United Kingdom, and the United States)."
- Overview (https://support.apple.com/guide/news-publisher/overview-apd19eaffd5b/1.0/icloud/1.0): "News Publisher allows content in English and Canadian French only and uses Apple News Format (ANF)." Content can come from a CMS, uploaded ANF files, or be authored in News Publisher.
- RSS requirements (https://support.apple.com/guide/news-publisher/apdc2c7520ff/icloud): valid RSS 2.0 or Atom, UTF-8; channel `title`, `language` (with locale), `link`; item `title`, `link`, `description`; full-text feed recommended; HTTPS feeds need a valid certificate or "your content doesn't appear in News"; titles 15-60 chars, no HTML; no ads, share buttons, "read more" or footer links in descriptions; author and thumbnail are read from the article HTML (author meta, `og:image` / schema.org image); include `atom:link rel="self"`.
- CMS / API (https://support.apple.com/guide/news-publisher/use-your-cms-with-news-publisher-apd88c8447e6/1.0/icloud/1.0): "you need to be signed in to iCloud with an Apple Account that has two-factor authentication turned on"; obtain channel ID, API Key ID and secret from Settings > Apple News API; Apple News API reference on developer.apple.com.
- AI policy: see section 7.
- Channel review/approval timing and the exact submission flow are not described on the pages fetched (open question below).

## 9. Realism for a brand-new domain

**What Google actually says (primary):**
- Eligible sites are those with "a consistent history of producing original news-related content" (2021-07-16 blog). No numeric site age or article count is published anywhere in the docs fetched.
- "automated systems run regularly and will detect when a site meets the criteria" (same post).
- Discover: eligibility "is not a guarantee of appearing"; traffic is "less predictable or dependable" and "supplemental" (Discover doc, 2026-03-09).
- Google News: "Having eligible content isn't a guarantee that content will necessarily rank well" (answer/9607025). "We can't provide much feedback regarding ranking" (answer/9606702).
- Expired domain abuse is a spam policy (section 7): do not buy an aged domain to shortcut this.

**Secondary / historical:** Search Engine Roundtable (https://www.seroundtable.com/new-publishers-google-news-29975.html **[secondary]**) reports John Mueller acknowledging in Oct 2020 that the Google News team was backed up accepting new sites, and in 2019 that the bar had been made "a little bit more strict". Treat as dated colour, not current policy.

**[opinion]** Expect zero News-tab/Top-stories presence for the first weeks-to-months regardless of technical perfection; Google's "consistent history" criterion is time-gated by construction. Discover can arrive earlier than news surfaces because it keys off indexed pages plus interest signals, but it is volatile. The realistic early traffic for hastcode.com is classic organic search on the architecture explainers (evergreen, long-tail, not freshness-dependent), plus direct/social/newsletter. Build the technical checklist on day one anyway: it is cheap in Go SSR, it is a prerequisite rather than an accelerator, and retrofitting bylines, author pages and dates later is more expensive than shipping them from the first article. A credible path to news surfaces is a narrow beat (e.g. Go ecosystem releases, infra/cloud outages, database vendor news) covered consistently with original angles rather than broad "tech news".

## Open questions

1. Does Google still run any manual review for news surfaces? Google says inclusion is automated; the 2021 post warns it may be outdated and the Dec 2019 blog.google post was not retrievable. Nothing newer contradicts the automated model.
2. Is the `<news:name>` in the news sitemap now simply the site name Google derives for the auto-generated publication page? The sitemap doc still says "as it appears on your articles on news.google.com"; answer/15898024 says titles come from site names. Assumed consistent; confirm once the site is indexed.
3. Bing PubHub current state: `pubhub.bing.com` did not resolve; no 2025-2026 official page found. Verify inside Bing Webmaster Tools after verifying the property.
4. Apple News: eligibility is restricted to publishers based in AU/CA/UK/US. Whether a company registered in India can open a channel at all was not determinable from the help pages; the channel review process and timing are not documented on the pages fetched.
5. Google News content policies page and ranking page show no last-updated date; they may have changed since the research date.
6. Whether Google News treats "architecture explainers" as news at all: the 2021 post and Bing's criteria both exclude how-to content from the news definition; keep explainers out of the news sitemap and mark them `Article`/`BlogPosting`, not `NewsArticle` (**[opinion]** grounded in the Bing newsworthiness definition; Google does not publish a definition).
7. Discover's "Following" feed uses RSS. The Search Central "feeds"/"follow" doc URLs tried returned 404; find the current doc before relying on feed-specific behaviour.

## Sources (accessed 2026-10-20)

**Google: Publisher Center help**
- https://support.google.com/news/publisher-center/answer/9607025 - News content across Google (no date shown)
- https://support.google.com/news/publisher-center/answer/6204050 - Google News content policies (no date shown)
- https://support.google.com/news/publisher-center/answer/9606702 - Ranking within Google News (no date shown)
- https://support.google.com/news/publisher-center/answer/14859160 - An update on publication pages (2024-04-25, updated 2024-12-18)
- https://support.google.com/news/publisher-center/answer/15898024 - Google News transitions to automatically-generated publication pages (2025-02-10, updated 2025-03-20)
- https://support.google.com/news/publisher-center/announcements/10146168 - Announcements index
- https://support.google.com/news/publisher-center/answer/9606710 - legacy news sitemap help (301 -> Search Central news-sitemap doc)

**Google: Search Central docs**
- https://developers.google.com/search/docs/appearance/google-discover (last updated 2026-03-09)
- https://developers.google.com/search/docs/crawling-indexing/sitemaps/news-sitemap (2025-12-10)
- https://developers.google.com/search/docs/appearance/structured-data/article (2026-09-08)
- https://developers.google.com/search/docs/crawling-indexing/robots-meta-tag (2026-03-24)
- https://developers.google.com/search/docs/appearance/page-experience (2026-09-22)
- https://developers.google.com/search/docs/appearance/core-web-vitals (2025-12-10)
- https://developers.google.com/search/docs/crawling-indexing/amp/about-amp (2026-07-01)
- https://developers.google.com/search/docs/essentials/spam-policies (2026-08-28)
- https://developers.google.com/search/docs/fundamentals/creating-helpful-content (2026-10-05)
- https://developers.google.com/search/docs/specialty/international/localized-versions (2026-09-21)
- https://developers.google.com/search/docs/specialty/international/managing-multi-regional-sites (2025-12-10)
- https://developers.google.com/search/docs/appearance/preferred-sources (2026-09-18)
- https://developers.google.com/search/docs/appearance/visual-elements-gallery (2026-02-04)

**Google: Search Central blog**
- https://developers.google.com/search/blog/2021/07/google-news-top-questions (2021-07-16)
- https://developers.google.com/search/blog/2020/11/timing-for-page-experience (2020-11-10)
- https://developers.google.com/search/blog/2023/02/google-search-and-ai-content (2023-02-08)
- https://developers.google.com/search/blog/2024/03/core-update-spam-policies (2024-03-05)

**Google: other**
- https://support.google.com/websearch/answer/9982767 - Discover content policies (no date shown)
- https://support.google.com/websearch/answer/2819496 - Discover help (where it appears)
- https://blog.google/products/search/preferred-sources/ (2025-08-12, from search result; not fetched directly)

**Bing**
- https://support.microsoft.com/en-us/bing/how-to-submit-your-news-website-to-bing (no date)
- https://blogs.bing.com/webmaster/february-2022/Reach-millions-of-people-with-Bing-News-PubHub (2022-02-08)
- https://blogs.bing.com/search/june-2016/Bing-helps-publishers-expand-their-reach (2016-06, from search result)
- https://pubhub.bing.com/Home/Help - Bing Publisher Guidelines (DNS failure on access date)

**Apple**
- https://support.apple.com/guide/news-publisher/toc - News Publisher User Guide
- https://support.apple.com/guide/news-publisher/publishing-on-apple-news-apde42330c66/1.0/icloud/1.0
- https://support.apple.com/guide/news-publisher/overview-apd19eaffd5b/1.0/icloud/1.0
- https://support.apple.com/guide/news-publisher/apdc2c7520ff/icloud - RSS feed requirements
- https://support.apple.com/guide/news-publisher/use-your-cms-with-news-publisher-apd88c8447e6/1.0/icloud/1.0
- https://support.apple.com/guide/news-publisher/ai-generated-content-apd8da39c6f2/1.0/icloud/1.0

**Secondary (labelled)**
- https://searchengineland.com/google-news-performance-report-added-to-google-search-console-345332 (2021)
- https://searchengineland.com/bing-news-pubhub-moves-to-bing-webmaster-tools-380070
- https://searchengineland.com/google-all-sites-are-eligible-to-be-in-google-news-but-not-all-content-will-appear-in-google-news-350521 (403 on fetch; snippet only)
- https://www.seroundtable.com/new-publishers-google-news-29975.html (2020)
- https://ppc.land/google-launches-a-new-publisher-center (Dec 2019 Publisher Center launch coverage)
- https://www.searchenginejournal.com/google-updates-article-structured-data-guidelines/462334 (Top stories / markup not required)

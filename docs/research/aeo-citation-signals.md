# AEO citation signals: what answer engines cite and what demonstrably moves it

Research note for hastcode.com (Go SSR tech-news + architecture explainers, mobile-first).
Compiled 2026-10-20. All URLs accessed 2026-10-20 unless noted. This space changes monthly; every source carries its own publication / last-updated date so you can judge staleness.

Evidence labels used throughout:

- **[PRIMARY]** official documentation or statement from the engine operator (Google Search Central, OpenAI, Perplexity, Microsoft Bing, schema.org) or a peer-reviewed paper.
- **[STUDY]** correlational / observational study by a third party. Sample size and date given. Correlation is not causation.
- **[SECONDARY]** trade-press or blog reporting of a primary source that I could not fetch directly.
- **[OPINION]** practitioner opinion with no data.
- **[UNSOURCED]** a claim I could not tie to a URL; treat as unknown.

---

## TL;DR

1. **No engine documents a special "AEO" markup, file, or writing style.** Google says so explicitly and repeatedly (AI features doc, 2025-12-10; AI optimization guide, 2026-07-10): no llms.txt, no AI text files, no special schema.org type, no "chunking", no rewriting for AI. Google Search ignores llms.txt. No other engine has officially said it reads llms.txt. [PRIMARY]
2. **The gate is the same as classic search: be crawlable, indexed, and snippet-eligible.** Google AI Overviews/AI Mode are controlled by Googlebot robots.txt + snippet controls, not Google-Extended. ChatGPT search requires `OAI-SearchBot` allowed (not GPTBot). Perplexity requires `PerplexityBot` allowed. Copilot rides on Bingbot / the Bing index. [PRIMARY]
3. **Strongest correlational predictor of being cited is off-page authority** (referring domains, brand mentions on Reddit/Quora/review sites) — SE Ranking, 129k domains, Nov 2025 [STUDY]. Overlap between AI citations and Google top-10 varies wildly by engine and study: Perplexity ~91% domain overlap (Semrush, 5k keywords, Jul 2025); AI Overviews 38% (Ahrefs, 863k SERPs, Mar 2026) vs ~17% (BrightEdge, Feb 2026). [STUDY]
4. **Word count does not matter for AI Overviews** (Spearman 0.04, n=174k pages, Ahrefs Dec 2025); 53% of cited pages are under 1,000 words. For ChatGPT, SE Ranking found longer pages cited more, but attributes it to subtopic coverage. [STUDY]
5. **On-page traits with the best evidence:** quotations, statistics and cited sources (+27–41% visibility in the KDD 2024 GEO paper, GPT-3.5 setup; replicated directionally on Perplexity) [PRIMARY peer-reviewed, but lab setting]; clear headings with 120–180-word sections; freshness (updated within 3 months: 6.0 vs 3.6 avg ChatGPT citations; AI-cited pages 25.7% fresher than organic across 17M URLs). Question-style headings and FAQ schema correlated *negatively* in SE Ranking's ChatGPT data. Keyword stuffing hurt in the GEO paper. [STUDY]
6. **Structured data to emit:** `Article`/`NewsArticle`/`BlogPosting` (Google-supported subtypes) with `author` (Person + url/sameAs), `datePublished`, `dateModified` (ISO 8601 with timezone, matching visible dates), `headline`, `image`; `Organization` on the home page. **Skip FAQPage (removed from Google Search 2026-05-07) and HowTo (removed Sept 2023).** `Speakable` is a US-English Google Assistant beta with no documented AI-search effect. [PRIMARY]
7. **Server-side render everything that matters.** Google renders JS (evergreen Chromium) but says "not all bots can run JavaScript"; the only large measurement (Vercel/MERJ, Dec 2024, ~1.3B AI-crawler requests) found GPTBot, ClaudeBot and PerplexityBot fetch JS but do not execute it. Bing's Oct 2025 guidance (secondary) says its crawler does not click, expand accordions, or scroll. Your Go SSR choice is the right default. [PRIMARY + STUDY]

---

## What to actually do (checklist for an SSR Go site)

### robots.txt allowlist

Allow by default; the only "owner decision" tokens are the training-only ones. All UA tokens below are quoted from the operators' own docs.

| Token (robots.txt `User-agent`) | Operator | What it gates | Recommendation | Source |
|---|---|---|---|---|
| `Googlebot` | Google | Google Search index, snippets, **AI Overviews and AI Mode** ("robots.txt directives for Googlebot is the control for site owners to manage access to how their sites are crawled" — AI features doc) | **Allow** | https://developers.google.com/search/docs/appearance/ai-features (2025-12-10) |
| `Google-Extended` | Google | Training of future Gemini models + grounding in Gemini Apps / Vertex AI. "does not impact a site's inclusion in Google Search nor is it used as a ranking signal". Has **no separate UA string**; it is a control token only. | Owner decision. No official sentence says blocking it removes you from AI Overviews; the AI features doc points to Googlebot as the AI-features control and describes Google-Extended as for "some of Google's other systems". | https://developers.google.com/search/docs/crawling-indexing/google-common-crawlers (2026-10-15) |
| `OAI-SearchBot` | OpenAI | ChatGPT search results & citations. "Sites that are opted out of OAI-SearchBot will not be shown in ChatGPT search answers". ~24h to take effect. | **Allow** (required for citations) | https://developers.openai.com/api/docs/bots |
| `ChatGPT-User` | OpenAI | On-demand page fetch when a user asks. "robots.txt rules may not apply". "not used to determine whether content may appear in Search". | Allow (blocking is mostly ineffective anyway) | same |
| `GPTBot` | OpenAI | Training data for foundation models. Independent of search. | Owner decision (training only) | same |
| `PerplexityBot` | Perplexity | "designed to surface and link websites in search results on Perplexity"; "not used to crawl content for AI foundation models". ~24h to take effect. | **Allow** | https://docs.perplexity.ai/guides/bots |
| `Perplexity-User` | Perplexity | User-initiated fetch; "generally ignores robots.txt rules". | n/a | same |
| `Bingbot` | Microsoft | Bing index, which Copilot and "AI powered search engines like Bing Copilot" rely on (Bing blog 2025-07-31). ChatGPT Enterprise/Edu search also uses Bing (OpenAI help, via Semrush quote). | **Allow** | https://blogs.bing.com/webmaster/2025/7/Keeping-Content-Discoverable-with-Sitemaps-in-AI-Powered-Search/ |

Out of scope but worth knowing: Anthropic (`ClaudeBot`, `Claude-SearchBot`, `Claude-User`) and Apple (`Applebot`) have analogous split tokens; not researched here.

Minimal robots.txt for hastcode.com (allow everything; add training blocks only if the owner wants them):

```
User-agent: *
Allow: /

# Optional: opt out of model training without affecting search/answer citations
# User-agent: GPTBot
# Disallow: /
# User-agent: Google-Extended
# Disallow: /

Sitemap: https://hastcode.com/sitemap.xml
```

### Robots meta / snippet controls (do NOT restrict)

- Google: `nosnippet`, `data-nosnippet`, `max-snippet`, `noindex` all limit AI Overviews/AI Mode exposure ("More restrictive permissions will limit how your content is featured in our AI experiences" — Google blog 2025-05-21). Leave them off on article pages. [PRIMARY]
- Bing: `noarchive` "prevents content from being used in Copilot responses and grounding results"; `nocache` "limits Copilot to using only the URL, title, and snippet"; Bing supports `data-nosnippet` for Copilot since 2025-10-15. Leave `noarchive`/`nocache` off. [PRIMARY for data-nosnippet; SECONDARY (SEJ 2026-02-27) for noarchive/nocache wording because the guideline page is JS-rendered and could not be fetched]

### Indexing plumbing

- XML sitemap with accurate `lastmod` (Bing: "freshness signals directly influence how quickly updates are reflected in search results and AI generated answers", 2025-07-31). [PRIMARY]
- Implement IndexNow (`/indexnow` POST, key file at root, up to 10,000 URLs per request) and ping on publish/update. It feeds Bing (and therefore Copilot) and other participants. IndexNow docs themselves do not mention AI. [PRIMARY]
- Verify in Google Search Console (AI Overviews/AI Mode clicks appear under the "Web" search type, not separately) and Bing Webmaster Tools (AI Performance report, public preview since 2026-02-10: total citations, avg cited pages, grounding queries, page-level citations). [PRIMARY]

### Structured data to emit (JSON-LD in `<head>`, Google's recommended format)

| Type | Emit? | Notes | Source |
|---|---|---|---|
| `NewsArticle` (news posts) / `BlogPosting` or `Article` (explainers) | **Yes** | Google's Article doc supports exactly `Article`, `NewsArticle`, `BlogPosting`. Recommended: `headline`, `image` (16x9, 4x3, 1x1), `author` (Person/Organization with `name` + `url` or `sameAs`), `datePublished`, `dateModified` (ISO 8601 **with timezone**). | https://developers.google.com/search/docs/appearance/structured-data/article (2026-09-08) |
| `TechArticle` | Optional / uncertain | schema.org subtype of Article with `dependencies`, `proficiencyLevel` (Beginner/Expert). Not listed in Google's Article doc. Safe pattern: emit `BlogPosting` (or `Article`) as the Google-supported type; if you want `TechArticle` semantics, flag as open question (see below). | https://schema.org/TechArticle (v30.1, 2026-09-16) |
| `Person` for author (via `author`) | **Yes** | `author.url` "a link to a page that uniquely identifies the author" (or `sameAs`). Build `/authors/<slug>` pages. | Article doc above |
| `Organization` | **Yes, home page only** | name (same as site name), url, logo (>=112px, crawlable), sameAs, description. "We recommend placing this information on your home page". | https://developers.google.com/search/docs/appearance/structured-data/organization (2026-09-08) |
| `FAQPage` | **No** | Restricted to government/health sites 2023-09-14; "will no longer appear in Google Search starting May 7, 2026"; docs removed 2026-06-15. SE Ranking data: pages with FAQ schema averaged *fewer* ChatGPT citations (3.6 vs 4.2). | https://developers.google.com/search/updates |
| `HowTo` | **No** | Docs removed 2023-09-14: "this rich result is no longer shown in search results, on both desktop and mobile devices". | same |
| `Speakable` | Low priority | Beta; US English; Google Assistant news read-aloud. No mention of AI Overviews/AI Mode. Cheap to add via `cssSelector` on your TL;DR block, but no evidence it affects citation. | https://developers.google.com/search/docs/appearance/structured-data/speakable (2026-09-08) |
| Any "AI-specific" markup, llms.txt, Markdown mirrors | **No** (for Google) | "You don't need to create new machine readable files, AI text files, markup, or Markdown to appear in Google Search... Google Search ignores them." | https://developers.google.com/search/docs/fundamentals/ai-optimization-guide (2026-07-10) |

Rules Google states for structured data: "Making sure your structured data matches the visible text on the page"; don't overfocus on structured data (listed as a myth). Dates: visible, labelled ("Published", "Last updated"), one prominent date, structured values must match visible values, never future dates, never the date of the event described.

### Page structure rules (evidence-weighted)

1. **Open with a direct answer / summary paragraph** (2–4 sentences) under the H1, then the body. Evidence: Google AI features doc "keep important content in text form" [PRIMARY]; Growth Memo reports 44.2% of ChatGPT citations come from the first 30% of a page (18,012 citations) [SECONDARY, snippet only, not verified against the article]; GEO paper's "Fluency" and "Easy-to-understand" both lifted visibility [PRIMARY peer-reviewed].
2. **Headings every 120–180 words, plain descriptive headings.** SE Ranking (ChatGPT, 216k pages): 120–180-word sections averaged 4.6 citations vs 2.7 for <50-word sections; question-style headings underperformed plain headings 3.4 vs 4.3. [STUDY] Google: content "organized by paragraphs and sections, along with headings that provide a clear structure". [PRIMARY]
3. **Include concrete numbers, quotations from named sources, and outbound citations with links.** GEO paper: Quotation Addition +40.9% / Statistics Addition +30.6% / Cite Sources +27.5% position-adjusted word count vs baseline (GPT-3.5 engine); on live Perplexity, quotes ≈ +22% PAWC, statistics ≈ +37% subjective impression. Caveat: benefits accrue mainly to lower-ranked sources; top-ranked sources lost share. [PRIMARY peer-reviewed, lab]
4. **Don't pad.** Word count has ~zero correlation with AI Overview citation (Ahrefs, 0.04 Spearman, 174k pages); keyword stuffing was the one strategy that *reduced* visibility in the GEO paper. Aim for complete coverage of the subtopics a reader needs, not length.
5. **Keep it fresh and say so.** Update substantively, bump `dateModified` and the visible "Last updated" date together. SE Ranking: updated within 3 months = 6.0 avg ChatGPT citations vs 3.6 outdated. Ahrefs (17M URLs): AI-cited URLs are 25.7% newer by publish date, 13.1% by last-update date than organic results; AI Overviews specifically are *not* fresher than organic (+16 days), ChatGPT is the most recency-biased (−393 to −458 days). [STUDY]
6. **Lists and tables:** no primary-source statement from any engine; no large-sample study isolating them was found in this research. Treat "use lists/tables" as [OPINION] (widely repeated, not measured). Google's guide mentions paragraphs, sections and headings only; "does not discuss tables".
7. **Unique, non-commodity content.** Google's single strongest stated factor: "Don't just recycle what others on the internet have already said, or could easily be produced by a generative AI model." Firsthand benchmarks, real code, real incident write-ups. [PRIMARY]
8. **Author pages + Organization identity + off-site presence.** Off-page authority is the strongest correlate in SE Ranking's ChatGPT model; Google lists "Seeking inauthentic 'mentions'" as a myth, i.e. earn real ones. [STUDY + PRIMARY]
9. **No content behind interaction.** Server-render all text, nav, metadata. Don't hide body text in tabs/accordions that require JS to populate (Bing, secondary); `<details>` with content present in HTML is fine for Google (it renders), unverified for others.
10. **Multimodal support:** Google recommends "quality images and videos where relevant" with the key content still in text. Notable: YouTube URLs are 5.6% of all AI Overview citations (Ahrefs Mar 2026). [PRIMARY + STUDY]

---

## 1. Google AI Overviews and AI Mode

**Primary docs**

- *AI features and your website* — https://developers.google.com/search/docs/appearance/ai-features (last updated 2025-12-10) [PRIMARY]
  - "There are no additional requirements to appear in AI Overviews or AI Mode, nor other special optimizations necessary."
  - "You don't need to create new machine readable files, AI text files, or markup to appear in these features."
  - "There's also no special schema.org structured data that you need to add."
  - Eligibility: page must be indexed and eligible to show with a snippet; "Indexing and serving isn't guaranteed."
  - Controls: `nosnippet`, `data-nosnippet`, `max-snippet`, `noindex`; "robots.txt directives for Googlebot is the control for site owners to manage access to how their sites are crawled".
  - Google-Extended is described only as: "To limit AI training and grounding in some of Google's other systems, read more about Google-Extended."
  - Best practices listed: allow crawling (incl. CDN/WAF), internal links, page experience, important content in text, supporting images/video, Merchant Center/Business Profile currency, Search Console verification.
  - Measurement: AI-feature clicks are "included in the overall search traffic in Search Console" under the Web search type.
- *Optimizing your website for generative AI features on Google Search* — https://developers.google.com/search/docs/fundamentals/ai-optimization-guide (published 2026-05-15, last updated 2026-07-10) [PRIMARY]
  - "our generative AI features on Google Search are rooted in our core Search ranking and quality systems."
  - "optimizing for generative AI search is optimizing for the search experience, and thus still SEO."
  - "Structured data isn't required for generative AI search, and there's no special schema.org markup you need to add."
  - "You don't need to create new machine readable files, AI text files, markup, or Markdown to appear in Google Search" ... "Google Search ignores them."
  - "There's no requirement to break your content into tiny pieces for AI to better understand it." / "You don't need to write in a specific way just for generative AI search."
  - Mythbusting list: LLMS.txt and "special" markup; chunking; rewriting for AI; seeking inauthentic mentions; overfocusing on structured data.
  - Strongest stated factor: unique, non-commodity content with "unique expert or experienced takes that go beyond common knowledge".
  - Structure: content "organized by paragraphs and sections, along with headings that provide a clear structure".
  - Mentions RAG used to "improve the quality, accuracy, and freshness of AI responses"; no other freshness rule.
  - AI agents: early guidance; mentions Universal Commerce Protocol (UCP).
- Announcement blog: https://developers.google.com/search/blog/2026/05/a-new-resource-for-optimizing (May 2026). [PRIMARY]
- *Top ways to ensure your content performs well in Google's AI experiences on Search* — https://developers.google.com/search/blog/2025/05/succeeding-in-ai-search (2025-05-21, John Mueller) [PRIMARY]. Eight points: unique non-commodity content; page experience; technical requirements; preview controls ("More restrictive permissions will limit how your content is featured in our AI experiences"); structured data matches visible content; multimodal; measure full value of visits; evolve.
- Changelog — https://developers.google.com/search/updates [PRIMARY]: 2026-05-15 guide added; 2026-05-15 spam policies clarified to apply to generative AI responses; 2026-05-27 "Preferred sources" rolling out to AI Overviews and AI Mode; 2026-06-15 llms.txt note added.

**Google-Extended vs Googlebot** — https://developers.google.com/search/docs/crawling-indexing/google-common-crawlers (last updated 2026-10-15) [PRIMARY]

- "a standalone product token that web publishers can use to manage" whether content "may be used for training future generations of Gemini models" and for grounding ("providing content from the Google Search index to the model at prompt time") in Gemini Apps and Vertex AI.
- "Google-Extended does not impact a site's inclusion in Google Search nor is it used as a ranking signal in Google Search."
- "Google-Extended doesn't have a separate HTTP request user agent string." Crawling uses existing Google UAs; the token is a control.
- Googlebot robots.txt tokens: `Googlebot`, `Googlebot-Image`, `Googlebot-Video`, `Googlebot-News`.

**Does blocking Google-Extended affect AI Overviews?** No Google page I could find says the words "Google-Extended does not affect AI Overviews". What the docs do say: (a) AI Overviews/AI Mode are Search features gated by Googlebot robots.txt and snippet controls; (b) Google-Extended governs Gemini training/grounding in "some of Google's other systems" and has no effect on Search inclusion or ranking. The inference that Google-Extended does not gate AI Overviews is consistent with both docs but is an inference. Secondary sources (e.g. Marie Haynes, https://www.mariehaynes.com/should-you-use-google-extended-in-robots-txt/) state it outright. [PRIMARY for the two facts; the conclusion is inferred]

**Structured data Google says matters for AI features:** none specifically. The only instruction is that markup must match visible content. Rich-result-oriented types (Article etc.) remain relevant for classic Search eligibility, which is the prerequisite.

## 2. Perplexity

- Bot docs — https://docs.perplexity.ai/guides/bots (no date on page; `perplexity.ai/perplexitybot` 301-redirects here) [PRIMARY]
  - `PerplexityBot` UA: `Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; PerplexityBot/1.0; +https://perplexity.ai/perplexitybot)`. "designed to surface and link websites in search results on Perplexity." "is not used to crawl content for AI foundation models." Recommends allowing in robots.txt; IPs at https://www.perplexity.com/perplexitybot.json.
  - `Perplexity-User` UA: `...compatible; Perplexity-User/1.0; +https://perplexity.ai/perplexity-user`. "Since a user requested the fetch, this fetcher generally ignores robots.txt rules." IPs at https://www.perplexity.com/perplexity-user.json.
  - Robots.txt changes take "up to 24 hours". Recommends WAF rules match UA + published IPs.
  - **No criteria for source selection or ranking are published.**
- Help centre FAQ *How does Perplexity work?* — https://hub-prod.perplexity.ai/hub/faq/how-does-perplexity-work (returned 404 on direct fetch; search snippet quotes: "Every answer includes footnotes linking to the original sources.") [PRIMARY, snippet only]
- What index Perplexity uses: not officially documented in anything I found. [UNSOURCED]
- Correlational: Semrush (5,000 keywords, >150,000 citations, 2025-07-21) found Perplexity citations had "over 91% domain overlap" and 82% URL overlap with Google's top 10 — the tightest coupling of any engine studied. [STUDY] Ahrefs freshness (2025-07-28): Perplexity-cited pages 250 days newer by publish date than organic. [STUDY]
- GEO paper (KDD 2024) tested its methods on live Perplexity: Quotation Addition ≈ +22% position-adjusted word count; Statistics Addition ≈ +37% subjective impression. [PRIMARY peer-reviewed]

Practical reading: for Perplexity, ranking well in Google (and Bing) is the dominant lever; on-page, quotes and stats measurably helped in a controlled test.

## 3. OpenAI ChatGPT search

- Bot docs — https://developers.openai.com/api/docs/bots (platform.openai.com/docs/bots 301s here; no date on page) [PRIMARY]
  - `OAI-SearchBot` ("...compatible; OAI-SearchBot/1.4; +https://openai.com/searchbot"): "OAI-SearchBot is for search." "Sites that are opted out of OAI-SearchBot will not be shown in ChatGPT search answers, though can still appear as navigational links." Allow it + its IPs (https://openai.com/searchbot.json); "~24 hours from a site's robots.txt update".
  - `GPTBot` ("...GPTBot/1.4; +https://openai.com/gptbot"): training. "Disallowing GPTBot indicates a site's content should not be used in training generative AI foundation models." IPs https://openai.com/gptbot.json.
  - `ChatGPT-User` ("...ChatGPT-User/1.0; +https://openai.com/bot"): user actions / GPT Actions; "not used for crawling the web in an automatic fashion"; "Because these actions are initiated by a user, robots.txt rules may not apply"; "not used to determine whether content may appear in Search". IPs https://openai.com/chatgpt-user.json.
  - "Each setting is independent of the others." "If your site has allowed both bots, we may use the results from just one crawl for both use cases."
  - Nothing about Bing, ranking, or citation criteria.
- Bing dependency: OpenAI help centre pages (https://help.openai.com/en/articles/9237897-chatgpt-search and https://help.openai.com/en/articles/10093903-chatgpt-search-for-enterprise-and-edu) returned 403 to the fetcher. Via search snippets and Semrush's quotation (2025-08-06, https://www.semrush.com/blog/chatgpt-definitely-uses-google/): "ChatGPT searches based on prompts and may share disassociated search queries with the Bing search engine to return web results" and, for Enterprise/Edu, "the only third-party search provider we use for Enterprise and Edu workspaces is Bing". Launch post (2024-10-31, https://openai.com/index/introducing-chatgpt-search/, also 403): "leverages third-party search providers, as well as content provided directly by our partners". [PRIMARY text, obtained SECONDARY]
- Semrush's Aug 2025 experiment (fictional-term page visible only to Google was cited by ChatGPT Plus) suggests Google is also used for paid tiers; OpenAI has not confirmed. [STUDY, n≈1 experiment, replicated once by Backlinko]
- Correlational, ChatGPT-specific: SE Ranking *Ranking factors for ChatGPT* (https://seranking.com/blog/ranking-factors-for-chatgpt/, reported by SEJ 2025-11-26; 129,000 domains, 216,524 pages, 20 niches; predictive model; ChatGPT version and timeframe unspecified) [STUDY]:
  - Referring domains strongest single predictor (1.6–1.8 avg citations ≤2,500 RDs vs 8.4 at >350k; doubling past ~32k).
  - Length: <800 words 3.2 citations; >2,900 words 5.1.
  - Freshness: updated within 3 months 6.0 vs outdated 3.6.
  - FAQ schema 3.6 vs 4.2 without (authors attribute to simpler support pages). llms.txt "negligible impact".
  - Sections of 120–180 words 4.6 vs <50 words 2.7. Question headings 3.4 vs plain 4.3.
  - Brand mentions: Quora heavy presence 7.0 vs 1.7; Reddit >10M mentions 7.0 vs 1.8; multi-platform review listings 4.6–6.3 vs 1.8.
  - SE Ranking's own caveat: correlational, "directional", not proof.
- Semrush (2025-07-21): ChatGPT had the lowest overlap with Google top-10 of the four engines (no exact figure given). Ahrefs (2025-07-28): ChatGPT the most recency-biased engine (−393/−458 days vs organic). [STUDY]

## 4. Bing / Copilot

- Bing Webmaster Guidelines — https://www.bing.com/webmasters/help/webmaster-guidelines-30fba23a. The page is client-rendered and returned only its title to the fetcher. SEJ's report of the Feb 2026 rewrite (https://www.searchenginejournal.com/bing-adds-geo-to-official-guidelines-expands-ai-abuse-definitions/568442/, 2026-02-27) [SECONDARY]:
  - Guidelines now cover Copilot and grounding API results; "grounding results and citations" are eligibility outcomes.
  - "NOARCHIVE prevents content from being used in Copilot responses and grounding results." "NOCACHE limits Copilot to using only the URL, title, and snippet." Bing advises against NOCACHE if you want richer citations.
  - Names "Generative Engine Optimization (GEO)" and says "GEO doesn't guarantee citations". Added a "Prompt Injection and AI Manipulation" section.
  - Third-party summaries also quote: "Bing and Copilot search experiences rely on the same core crawling, indexing, and ranking foundation as traditional search." I could not verify this sentence on the primary page. [SECONDARY]
- Bing blog, *Keeping Content Discoverable with Sitemaps in AI Powered Search* — https://blogs.bing.com/webmaster/2025/7/Keeping-Content-Discoverable-with-Sitemaps-in-AI-Powered-Search/ (2025-07-31) [PRIMARY]: "AI powered search engines like Bing Copilot"; sitemaps "remain a foundational signal"; "freshness signals directly influence how quickly updates are reflected in search results and AI generated answers"; "no tool can guarantee when or how your content will appear in AI-generated results."
- Bing blog, *data-nosnippet support* — https://blogs.bing.com/webmaster/2025/10/Bing-Introduces-Support-for-the-data-nosnippet-HTML-Attribute/ (2025-10-15) [PRIMARY]: "precise control over what content appears in search results and AI-generated answers"; example comment "This section does not appear in Bing Search and Copilot answers."
- Bing blog, *IndexNow Drives Smarter and Faster Content Discovery* — https://blogs.bing.com/webmaster/2025/5/IndexNow-Drives-Smarter-and-Faster-Content-Discovery/ (2025-05-19) [PRIMARY]: only AI link is "empowering content owners to stay current in a search environment shaped by immediacy and AI".
- IndexNow protocol — https://www.indexnow.org/documentation [PRIMARY]: key file (8–128 chars) at root; GET with `url`+`key` or JSON POST to `/indexnow` with up to 10,000 URLs; submissions shared with all participating engines; no AI mention.
- AI Performance report — https://blogs.bing.com/webmaster/February-2026/Introducing-AI-Performance-in-Bing-Webmaster-Tools-Public-Preview (2026-02-10; page timed out on fetch; title/date confirmed from the blog index and SEJ https://www.searchenginejournal.com/bing-webmaster-tools-adds-ai-citation-performance-data/566874/) [PRIMARY date, SECONDARY content]: metrics Total Citations, Average Cited Pages, Grounding Queries, page-level citations across Copilot, Bing AI summaries, partner integrations; data is a sample. Topics Portfolio beta added 2026-10-13 (https://blogs.bing.com/webmaster/2026/10/introducing-topics-portfolio-beta-in-ai-performance/).
- "Bing's crawler does not click, expand accordions, switch tabs, or scroll" — attributed to an Oct 2025 Bing team post by third parties; I could not locate the primary post. [SECONDARY, unverified]

## 5. Structured data that matters

- Formats: "Google recommends using JSON-LD" — https://developers.google.com/search/docs/appearance/structured-data/intro-structured-data (2025-12-10). [PRIMARY]
- Article — https://developers.google.com/search/docs/appearance/structured-data/article (2026-09-08) [PRIMARY]: types `Article`, `NewsArticle`, `BlogPosting`; recommended `author` (+`name`, `url`/`sameAs`), `datePublished`, `dateModified` ("Add the dateModified property if you want to provide more accurate date information to Google"; timezone recommended, else Googlebot's timezone assumed), `headline`, `image`. `TechArticle` not mentioned.
- TechArticle — https://schema.org/TechArticle (v30.1, 2026-09-16) [PRIMARY]: "A technical article - Example: How-to (task) topics, step-by-step, procedural troubleshooting, specifications, etc."; adds `dependencies`, `proficiencyLevel`.
- FAQPage: 2023-09-14 restricted to "well-known, authoritative government and health websites"; 2026-05-08 changelog: "will no longer appear in Google Search starting May 7, 2026"; 2026-06-15 docs removed. https://developers.google.com/search/updates [PRIMARY]
- HowTo: 2023-09-14 "Removed the How-to structured data documentation, as this rich result is no longer shown in search results, on both desktop and mobile devices." Blog: https://developers.google.com/search/blog/2023/08/howto-faq-changes (Aug 2023). [PRIMARY]
- Speakable — https://developers.google.com/search/docs/appearance/structured-data/speakable (2026-09-08) [PRIMARY]: "This feature is in beta and subject to change"; US / English / Google Home; 2–3 sentences per section; no AI Overviews mention. schema.org property — https://schema.org/speakable: `SpeakableSpecification` with `cssSelector` or `xpath`, on `Article`/`WebPage`.
- Organization — https://developers.google.com/search/docs/appearance/structured-data/organization (2026-09-08) [PRIMARY]: no required properties; `name`, `url`, `logo`, `sameAs`, `description`, `address`, `contactPoint`, `foundingDate`...; "We recommend placing this information on your home page, or a single page that describes your organization."
- Relevance to AI citation specifically: Google says none is required and lists "overfocusing on structured data" as a myth. SE Ranking found FAQ schema correlated negatively for ChatGPT. No study found that Article/Person/Organization markup correlates with AI citation; treat them as classic-search hygiene that keeps you eligible. [PRIMARY + STUDY]

## 6. Freshness

- Google dates doc — https://developers.google.com/search/docs/appearance/publication-dates (2025-12-10) [PRIMARY]: visible, labelled date; `datePublished`/`dateModified` on a `CreativeWork` subtype; "Ensure that the date (and optional time and timezone) match between the equivalent user-visible and structured values"; "Don't specify future dates, or the date of the action described on the page"; "Minimize the presence of other dates on the page". Google estimates "when a page was published or significantly updated" — i.e. a cosmetic `dateModified` bump is not what it is looking for.
- Google AI guide: RAG is used to "improve the quality, accuracy, and freshness of AI responses" — the only official freshness statement tied to AI features. [PRIMARY]
- Bing: "freshness signals directly influence how quickly updates are reflected in search results and AI generated answers" (2025-07-31). [PRIMARY]
- Ahrefs, *Do AI assistants prefer to cite fresh content?* — https://ahrefs.com/blog/do-ai-assistants-prefer-to-cite-fresh-content/ (2025-07-28; 16.975M cited URLs across ChatGPT, Perplexity, Gemini, Copilot, AI Overviews vs organic; publish age = first seen by Ahrefs crawler, update age = Content Explorer last-published) [STUDY]: AI-cited URLs average 1,064 days old vs 1,432 organic (25.7% fresher by publish date); 909 vs 1,046 by update date (13.1%). Per engine vs organic (publish / update): AI Overviews +16 / +21 days (not fresher); Perplexity −250 / −53; Gemini −298 / −216; Copilot −360 / −182; ChatGPT references −393 / −182; ChatGPT citations −458 / −58. Average cited page is still ~2.9 years old.
- SE Ranking (ChatGPT, Nov 2025): updated within 3 months 6.0 avg citations vs 3.6 outdated. [STUDY]
- The oft-repeated "76% of ChatGPT's most-cited pages updated in the last 30 days" could not be tied to a primary report. [UNSOURCED]

## 7. Content structure: what correlates with citation

| Trait | Direction | Evidence | Label |
|---|---|---|---|
| Quotations from named sources | + (strongest in GEO paper) | GEO-bench 10k queries, GPT-3.5 engine: +40.9% PAWC, +28.0% subjective; live Perplexity ≈ +22% PAWC | PRIMARY peer-reviewed (lab) |
| Statistics / numbers | + | +30.6% PAWC, +22.8% subjective; Perplexity ≈ +37% subjective | PRIMARY peer-reviewed (lab) |
| Citing sources (links) | + | +27.5% PAWC; but −30.3% for already rank-1 sources, +115% for rank-5 | PRIMARY peer-reviewed (lab) |
| Fluency / readability | + | Fluency +28.0% PAWC; Easy-to-understand +14.0% | PRIMARY peer-reviewed (lab) |
| Authoritative tone | + modest | +10.4% PAWC, +18.7% subjective | PRIMARY peer-reviewed (lab) |
| Technical terms | + modest | +17.6% PAWC | PRIMARY peer-reviewed (lab) |
| Keyword stuffing | − | −8.3% PAWC (only negative method) | PRIMARY peer-reviewed (lab) |
| Word count | ~0 (AIO) / + (ChatGPT) | Ahrefs 174k pages: Spearman 0.04, avg 1,282 words, 53.4% <1,000, 16% >2,000, position 1 avg 1,270 vs positions 4–10 avg 1,690 (2025-12-03). SE Ranking: >2,900 words 5.1 vs <800 3.2 (attributed to subtopic coverage) | STUDY |
| Section length / headings | + for 120–180-word sections | SE Ranking 2025-11 | STUDY |
| Question-style headings | − | SE Ranking 3.4 vs 4.3 plain | STUDY (contradicts common [OPINION]) |
| Answer-first / early placement | + | Growth Memo: 44.2% of 18,012 ChatGPT citations from first 30% of page (snippet, unverified) | SECONDARY |
| Lists / tables | unknown | No primary statement; no isolating study found | OPINION |
| FAQ schema | − | SE Ranking 3.6 vs 4.2 | STUDY |
| Freshness | + | see section 6 | STUDY (+ PRIMARY statements) |
| Referring domains / brand mentions | + (strongest) | SE Ranking | STUDY |
| Unique, non-commodity content | + | Google's stated #1 factor | PRIMARY (statement, not measurement) |
| Ranking in Google top 10 | + varies by engine | Perplexity 91% domain overlap; AI Overviews 86% domain / 67% URL (Semrush Jul 2025) vs 37.9% URL in top 10 blocks (Ahrefs Mar 2026, 863k SERPs, 4M URLs) vs ~17% (BrightEdge Feb 2026 weekly panel, Feb 2025–Feb 2026); AI Mode ~54% domain / ~35% URL; ChatGPT lowest | STUDY (methodologies differ; don't average) |
| Cross-engine agreement | low | BrightEdge: pairwise source overlap 16–59% across five AI surfaces (via SEJ, unverified) | SECONDARY |

GEO paper details: Aggarwal, Murahari, Rajpurohit, Kalyan, Narasimhan, Deshpande, *GEO: Generative Engine Optimization*, KDD 2024, arXiv:2311.09735 (v1 2023-11-16, v3 2024-06-28), https://arxiv.org/abs/2311.09735, HTML https://arxiv.org/html/2311.09735v3. Benchmark 10K queries (8K/1K/1K) from MS MARCO, ORCAS-1, Natural Questions, AllSouls, LIMA, Davinci-Debate, Perplexity Discover, ELI5, GPT-4-generated. Important caveats: evaluation engine is a GPT-3.5-turbo pipeline from 2023, not today's Google/OpenAI products; percentages above are computed from the paper's Table 1 against a 19.3 baseline; the Perplexity test is a smaller appendix experiment.

Kevin Indig / Growth Memo: no primary, non-paywalled study on AI Overview page traits was retrievable. The research summary (https://www.growth-memo.com/p/2026-growth-memo-research-summary, 2026-06-24) is paywalled; search snippets attribute to it: Omnia dataset showing only 2.4% of cited URLs shared across ChatGPT/Perplexity/AI Overviews for the same prompt (3.7M citations), and an analysis of >500k AI Overviews finding most citations point to established authoritative sources. [SECONDARY, unverified]

## 8. llms.txt

- Proposal: https://llmstxt.org/ — Jeremy Howard (Answer.AI / fast.ai), first published 2024-09-03, modified 2026-08-10. Markdown file at `/llms.txt`: H1 name, blockquote summary, H2 link sections, optional "Optional" section; suggests `.md` mirrors and `rel="alternate" type="text/markdown"`. The page claims thousands of sites publish one and that OpenAI, Anthropic, Gemini publish one for their *own docs*; it makes **no** claim that any search engine or AI product reads it. [PRIMARY for the proposal]
- Google: "Google Search ignores them" (AI optimization guide, note added 2026-06-15; changelog: "these files aren't needed for Google Search (and won't negatively or positively impact your visibility or rankings)"). [PRIMARY]
- John Mueller on Reddit r/TechSEO (reported by SEJ 2025-04-17, https://www.searchenginejournal.com/google-llms-txt-comparable-keywords-meta-tag/544804/): "AFAIK none of the AI services have said they're using LLMs.TXT (and you can tell when you look at your server logs that they don't even check for it)... comparable to the keywords meta tag". [PRIMARY statement, SECONDARY reporting]
- OpenAI, Perplexity, Microsoft: no official statement found that their search products consume llms.txt. OpenAI's and Perplexity's bot docs do not mention it. [absence of evidence, strictly]
- SE Ranking (ChatGPT, Nov 2025): llms.txt "negligible impact". [STUDY]
- Chrome Lighthouse 13.3 reportedly added an "Agentic Browsing" audit that checks for llms.txt — aimed at browser agents, not search. [SECONDARY, unverified]

Verdict: harmless, cheap, zero documented citation benefit. If you ship one, do it for agent tooling, not AEO.

## 9. Technical: rendering and JavaScript

- Google — https://developers.google.com/search/docs/crawling-indexing/javascript/javascript-seo-basics [PRIMARY]: three phases crawl → render queue ("may stay on this queue for a few seconds, but it can take longer") → index; "evergreen version of Chromium"; "Google Search won't render JavaScript from blocked files or on blocked pages"; SSR/pre-rendering "still a great idea" because "not all bots can run JavaScript".
- OpenAI, Perplexity, Microsoft publish **no** statement on whether their search crawlers execute JavaScript. [absence]
- Vercel + MERJ, *The rise of the AI crawler* — https://vercel.com/blog/the-rise-of-the-ai-crawler (2024-12-17; one month of Vercel-network traffic: Googlebot 4.5B, GPTBot 569M, Claude 370M, AppleBot 314M, PerplexityBot 24.4M requests) [STUDY]: "none of the major AI crawlers currently render JavaScript"; ChatGPT crawlers fetched JS for 11.50% of requests and Claude 23.84% but "They can't read client-side rendered content"; Gemini "leverages Googlebot's infrastructure, enabling full JavaScript rendering"; AppleBot renders. Recommendation: "Prioritize server-side rendering for critical content." Caveat: 22 months old; crawler capabilities may have changed; no newer large-sample measurement found.
- Bing: guideline rewrite reportedly dropped its JavaScript-rendering section (SEJ 2026-02-27); third parties attribute "does not click, expand accordions, switch tabs, or scroll" to Bing (Oct 2025). [SECONDARY]

Implication for a Go SSR site: emit complete HTML (title, meta, canonical, JSON-LD, headings, body, dates, nav, pagination) with no JS dependency; keep `<details>`/tab content present in the initial HTML; avoid lazy-loading article text; serve fast (Google lists page experience among its AI-feature best practices).

---

## Open questions

1. **Does Google treat `TechArticle` as an `Article` for Article rich-result purposes?** Google's doc lists only `Article`, `NewsArticle`, `BlogPosting`. Not stated either way. Safe default: `BlogPosting`/`NewsArticle`; test `TechArticle` in the Rich Results Test before committing.
2. **Does blocking `Google-Extended` affect AI Overviews/AI Mode?** No explicit official sentence; inferred "no" from the Googlebot-is-the-control statement and the "no impact on inclusion or ranking" statement.
3. **Do OAI-SearchBot and PerplexityBot execute JavaScript today?** No official doc; only the Dec 2024 Vercel measurement (no).
4. **Does ChatGPT (consumer tiers) use Google results?** Semrush's Aug 2025 experiment says yes for Plus; OpenAI only names Bing (Enterprise/Edu) and "other search providers".
5. **Exact wording of Bing's Feb 2026 Webmaster Guidelines** (noarchive/nocache/GEO): page is client-rendered; verify in a browser.
6. **Do lists/tables help?** No measurement found. Could be tested with your own Bing AI Performance / Search Console data once live.
7. **Perplexity ranking inputs** remain entirely undocumented.
8. **How much does the GEO paper's 2023 GPT-3.5 lab result transfer to 2026 production engines?** Only the small Perplexity appendix test exists as a bridge.

---

## Sources (all accessed 2026-10-20)

Google (PRIMARY)
- https://developers.google.com/search/docs/appearance/ai-features — AI features and your website (updated 2025-12-10)
- https://developers.google.com/search/docs/fundamentals/ai-optimization-guide — Optimizing your website for generative AI features (published 2026-05-15, updated 2026-07-10)
- https://developers.google.com/search/blog/2026/05/a-new-resource-for-optimizing — announcement (May 2026)
- https://developers.google.com/search/blog/2025/05/succeeding-in-ai-search — Top ways to ensure your content performs well in Google's AI experiences (2025-05-21)
- https://developers.google.com/search/docs/crawling-indexing/google-common-crawlers — Google-Extended, Googlebot tokens (updated 2026-10-15)
- https://developers.google.com/search/docs/crawling-indexing/overview-google-crawlers — crawler overview (updated 2026-06-12)
- https://developers.google.com/search/updates — documentation changelog (entries 2026-05-08, 2026-05-15, 2026-05-27, 2026-06-15; 2023-09-14)
- https://developers.google.com/search/docs/appearance/structured-data/article — Article structured data (updated 2026-09-08)
- https://developers.google.com/search/docs/appearance/structured-data/organization — Organization (updated 2026-09-08)
- https://developers.google.com/search/docs/appearance/structured-data/speakable — Speakable (updated 2026-09-08)
- https://developers.google.com/search/docs/appearance/structured-data/intro-structured-data — formats, JSON-LD recommended (updated 2025-12-10)
- https://developers.google.com/search/docs/appearance/publication-dates — publication/update dates (updated 2025-12-10)
- https://developers.google.com/search/blog/2023/08/howto-faq-changes — Changes to HowTo and FAQ rich results (Aug 2023)
- https://developers.google.com/search/docs/crawling-indexing/javascript/javascript-seo-basics — JavaScript SEO basics

OpenAI (PRIMARY; help-centre pages returned 403 and are quoted via secondary sources)
- https://developers.openai.com/api/docs/bots — OAI-SearchBot, GPTBot, ChatGPT-User (platform.openai.com/docs/bots redirects here)
- https://help.openai.com/en/articles/9237897-chatgpt-search — ChatGPT search help (403 on fetch)
- https://help.openai.com/en/articles/10093903-chatgpt-search-for-enterprise-and-edu — Enterprise/Edu: Bing as only third-party provider (403 on fetch; quoted by Semrush)
- https://openai.com/index/introducing-chatgpt-search/ — launch post (2024-10-31; 403 on fetch)

Perplexity (PRIMARY)
- https://docs.perplexity.ai/guides/bots — PerplexityBot and Perplexity-User (undated)
- https://hub-prod.perplexity.ai/hub/faq/how-does-perplexity-work — FAQ (404 on fetch; snippet only)

Microsoft Bing (PRIMARY unless noted)
- https://www.bing.com/webmasters/help/webmaster-guidelines-30fba23a — Webmaster Guidelines (client-rendered; not readable by fetcher)
- https://blogs.bing.com/webmaster/2025/7/Keeping-Content-Discoverable-with-Sitemaps-in-AI-Powered-Search/ (2025-07-31)
- https://blogs.bing.com/webmaster/2025/10/Bing-Introduces-Support-for-the-data-nosnippet-HTML-Attribute/ (2025-10-15)
- https://blogs.bing.com/webmaster/2025/5/IndexNow-Drives-Smarter-and-Faster-Content-Discovery/ (2025-05-19)
- https://blogs.bing.com/webmaster/February-2026/Introducing-AI-Performance-in-Bing-Webmaster-Tools-Public-Preview (2026-02-10; timed out on fetch)
- https://blogs.bing.com/webmaster/2026/10/introducing-topics-portfolio-beta-in-ai-performance/ (2026-10-13)
- https://www.indexnow.org/documentation — IndexNow protocol
- https://www.searchenginejournal.com/bing-adds-geo-to-official-guidelines-expands-ai-abuse-definitions/568442/ — SECONDARY report of guideline rewrite (2026-02-27)
- https://www.searchenginejournal.com/bing-webmaster-tools-adds-ai-citation-performance-data/566874/ — SECONDARY (Feb 2026)

schema.org (PRIMARY)
- https://schema.org/TechArticle (v30.1, 2026-09-16)
- https://schema.org/speakable

Peer-reviewed (PRIMARY)
- https://arxiv.org/abs/2311.09735 — Aggarwal et al., GEO: Generative Engine Optimization, KDD 2024 (v3 2024-06-28); full text https://arxiv.org/html/2311.09735v3

llms.txt
- https://llmstxt.org/ — proposal, Jeremy Howard (2024-09-03, modified 2026-08-10)
- https://www.searchenginejournal.com/google-llms-txt-comparable-keywords-meta-tag/544804/ — Mueller Reddit comment (SEJ 2025-04-17) [SECONDARY reporting of PRIMARY statement]

Correlational studies [STUDY]
- https://ahrefs.com/blog/ai-overview-citations-top-10/ — Ahrefs, 863K SERPs / 4M AIO URLs, 37.9% in top 10 (2026-03-02)
- https://ahrefs.com/blog/short-vs-long-content-in-ai-overviews/ — Ahrefs, 560,346 AIOs / 174,048 pages, Spearman 0.04 (2025-12-03)
- https://ahrefs.com/blog/do-ai-assistants-prefer-to-cite-fresh-content/ — Ahrefs, 16.975M cited URLs, 25.7% / 13.1% fresher (2025-07-28)
- https://seranking.com/blog/ranking-factors-for-chatgpt/ — SE Ranking, 129,000 domains / 216,524 pages / 20 niches (Nov 2025); reported at https://www.searchenginejournal.com/new-data-top-factors-influencing-chatgpt-citations/ (2025-11-26)
- https://www.semrush.com/blog/ai-mode-comparison-study/ — Semrush, 5,000 keywords / >150,000 citations, overlap with Google top 10 (2025-07-21)
- https://www.semrush.com/blog/chatgpt-definitely-uses-google/ — Semrush, ChatGPT/Google experiment + OpenAI Bing quotes (2025-08-06)
- https://help.brightedge.com/resources/weekly-ai-search-insights/ai-overviews-one-year-presence-size-citing — BrightEdge, ~17% top-10 overlap, Feb 2025–Feb 2026 panel (2026-02-12)
- https://vercel.com/blog/the-rise-of-the-ai-crawler — Vercel/MERJ, AI crawlers and JavaScript (2024-12-17)
- https://www.growth-memo.com/p/2026-growth-memo-research-summary — Kevin Indig, paywalled (2026-06-24) [SECONDARY, snippets only]

Other secondary
- https://www.mariehaynes.com/should-you-use-google-extended-in-robots-txt/ — Google-Extended vs AI Overviews commentary [SECONDARY]

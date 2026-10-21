# Spec: hastcode.com v1 — static-first publication

Status: ready-for-agent
Created: 2026-10-21
Glossary: ../../GLOSSARY.md
ADRs: ../../docs/adr/0001-replace-genuin-community-with-static-publication.md, ../../docs/adr/0002-static-first-go-build-on-vercel.md
Research: ../../docs/research/ (ad-rpm-tech-audience, aeo-citation-signals, google-news-discover-requirements, vercel-go-ssr-viability)

## Problem Statement

The Author wants to publish software-architecture Explainers and Briefs at hastcode.com and earn $10k/month from display ads plus newsletter sponsorship within 12–18 months. Today the domain serves an unrelated Genuin community app. There is no publishing system: no way to turn a Markdown file into a fast, mobile-first page that Google, Discover, Google News, and answer engines (AI Overviews, ChatGPT search, Perplexity, Copilot) will index and cite, and no way to capture Readers into a newsletter list.

The Author produces ~5 Explainers and ~10 Briefs a week, AI-drafted and human-reviewed, and needs publishing to be `git push`. Readers are mid/senior backend engineers in the US, Europe, and India, mostly on phones, 72% with ad blockers. Every page must load fast from every region, carry the structured data and transparency signals the engines require, and never ship a page that violates the editorial skeleton.

## Solution

A Go program that is a **build tool**, not a server (ADR-0002). It reads Markdown with frontmatter from `content/`, renders pages through templ templates, and writes a complete static site into Vercel's Build Output API directory: HTML for every page type, JSON-LD, RSS feeds, sitemaps (standard + news), an IndexNow ping manifest, responsive image variants, a 404 page, and a 410 map for the retired community URLs. Vercel serves it from every edge. One Go serverless function handles newsletter signups by forwarding to Buttondown.

The build is also the **gatekeeper**: it fails on any page that breaks the editorial skeleton, exceeds the size budget, lacks a hero image or JSON-LD, has a broken internal link, or duplicates a slug. The Author therefore cannot publish a page that would hurt SEO or AEO.

The site ships with zero client-side framework, dark mode, Pagefind search, reserved (empty) ad slots, Plausible analytics, and the transparency pages Google News needs. Ads, consent management, sponsorship pages and contributors are deliberately not in v1.

## User Stories

### Reading

1. As a Reader on a phone in India, I want an Explainer to render in under 1.5 s on a slow connection, so that I read it instead of bouncing.
2. As a Reader, I want every page to work with JavaScript disabled, so that content and navigation never depend on scripts.
3. As a Reader, I want the site to follow my system dark/light preference and let me override it, so that I can read comfortably at night.
4. As a Reader, I want an Explainer to open with a TL;DR, so that I get the answer before deciding to read on.
5. As a Reader, I want short declarative sections with at least one architecture diagram, so that I can skim and still understand the system.
6. As a Reader, I want a Comparison to always include a trade-offs table, so that I can decide between alternatives at a glance.
7. As a Reader, I want a Sources block at the end of every Explainer linking first-party material, so that I can verify claims.
8. As a Reader, I want to see "Published" and "Last reviewed" dates on every Explainer, so that I know how current it is.
9. As a Reader, I want a Brief to say what happened, why it matters architecturally, and link to the relevant Explainers, so that news becomes learning.
10. As a Reader, I want "More in {Vertical}" and "Related" links on every Explainer, so that I can keep reading on the same subject.
11. As a Reader, I want a Vertical hub page listing its Explainers newest-first with pagination, so that I can browse a subject.
12. As a Reader, I want a Briefs index page, so that I can see what happened recently.
13. As a Reader, I want a home page showing the latest Explainers, latest Briefs, and the six Verticals, so that I can orient quickly.
14. As a Reader, I want to filter Explainers by kind (Case Study, Concept, Comparison), so that I can find the format I want.
15. As a Reader, I want site search that works without a backend, so that I can find an Explainer by keyword.
16. As a Reader, I want an Author page with a real bio, photo and links, so that I know who is accountable for what I read.
17. As a Reader, I want About, Contact and Editorial Policy pages, so that I understand how the site is made, including AI assistance.
18. As a Reader, I want to subscribe to the newsletter from any article without leaving the page flow, so that I can keep up without ads.
19. As a Reader, I want a clear confirmation or error after subscribing, so that I know whether it worked.
20. As a Reader, I want a helpful 404 page, so that a mistyped URL still leads somewhere useful.
21. As a Reader, I want RSS for the whole site, for each Vertical, and for Briefs only, so that I can follow in my reader.
22. As a Reader, I want diagrams to be crisp at any zoom and readable in dark mode, so that they are useful on a phone.
23. As a Reader with an ad blocker, I want no broken layout or empty gaps where ads would be, so that the page still looks intentional.

### Authoring and publishing

24. As the Author, I want to publish by committing Markdown and pushing to `main`, so that publishing takes no extra tooling.
25. As the Author, I want a preview deployment for every branch, so that I can review a draft at its real URL before merging.
26. As the Author, I want frontmatter to declare kind, Vertical, tags, dates, hero image, description and related slugs, so that pages are built from metadata not hand-written HTML.
27. As the Author, I want to write Mermaid in Markdown and get an inline SVG, so that I never hand-draw diagrams or manage image files for them.
28. As the Author, I want to drop one hero raster per piece and have the build produce all variants, so that I never resize images by hand.
29. As the Author, I want the build to fail with a clear message if an Explainer lacks a TL;DR, a Sources block, a hero image, or a diagram, so that I cannot publish below the editorial floor.
30. As the Author, I want the build to fail if a Comparison lacks a trade-offs table, so that the kind's promise holds.
31. As the Author, I want the build to fail on any broken internal link, so that Readers and crawlers never hit dead ends.
32. As the Author, I want the build to fail on duplicate slugs, so that URLs are unambiguous.
33. As the Author, I want the build to fail if any page's HTML exceeds 50 KB uncompressed, so that the weight budget holds.
34. As the Author, I want the build to fail if a hero image is narrower than 1200 px, so that every piece is Discover-eligible.
35. As the Author, I want the build to warn (not fail) on Explainers not Reviewed in 90 days, so that freshness decays visibly.
36. As the Author, I want a `draft: true` flag that excludes a piece from production output but includes it in previews, so that I can stage work.
37. As the Author, I want the build to run in seconds for hundreds of pages, so that publishing stays instant.
38. As the Author, I want a local command that builds and serves the output, so that I can check pages without deploying.
39. As the Author, I want image variants cached between builds, so that adding one piece does not re-encode every image.
40. As the Author, I want the Brief body limited to its skeleton (what happened, why it matters, Explainer links, source), so that Briefs stay short and the build enforces it.
41. As the Author, I want new Verticals added by editing one configuration list, so that adding Security later is a one-line change.

### Search engines and answer engines

42. As Googlebot, I want every page to be complete HTML with no content behind JavaScript, so that indexing needs no rendering.
43. As Googlebot, I want an `Article` JSON-LD block on every Explainer with headline, author Person, datePublished and dateModified (with timezone, matching visible dates), image in 16:9, 4:3 and 1:1, so that rich results and AI features can use it.
44. As Googlebot, I want a `NewsArticle` JSON-LD block on every Brief, so that Briefs are treated as news.
45. As Googlebot, I want `BreadcrumbList`, `WebSite` and `Organization` JSON-LD where appropriate, so that site structure is machine-readable.
46. As Googlebot, I want `max-image-preview:large` on every page and an `og:image` at least 1200 px wide, so that Discover can feature the page.
47. As Googlebot, I want a standard sitemap of every canonical URL and a separate news sitemap of Briefs from the last 48 hours (name, language, publication_date, title), so that both indexes stay current.
48. As Googlebot, I want canonical tags, no duplicate URLs with and without trailing slashes, and no `.html` suffixes, so that signals consolidate.
49. As Googlebot, I want every retired community URL to return 410, and `/home` to 301 to `/`, so that the old site is cleanly removed (ADR-0001).
50. As OAI-SearchBot, PerplexityBot and Bingbot, I want to be explicitly allowed in robots.txt, so that ChatGPT search, Perplexity and Copilot can cite the site.
51. As Bingbot, I want an IndexNow ping with the key file served from the site root for every new or changed URL, so that Bing and Copilot pick up changes quickly.
52. As an answer engine, I want sections of 120–180 words, direct statements, statistics and quotations with citations, so that I can extract and attribute an answer.
53. As any crawler, I want the hero image served with width/height attributes and modern formats with fallbacks, so that LCP is stable and CLS is zero.
54. As a crawler, I want `hreflang` absent (English only) and `lang="en"` set, so that language is unambiguous.

### Measurement and monetisation readiness

55. As the Author, I want cookieless analytics on every page, so that I see traffic by country and page without a consent banner.
56. As the Author, I want reserved, fixed-height ad slots in the Explainer and Brief layouts that render empty in v1, so that turning ads on later causes no layout shift.
57. As the Author, I want ad slots to be a single template component controlled by one flag, so that enabling AdSense is a configuration change.
58. As the Author, I want newsletter signups stored in Buttondown with a source tag of the page, so that I know which Explainers convert.
59. As the Author, I want the subscribe function to reject obviously invalid emails and rate-limit abuse, so that the list stays clean.
60. As the Author, I want the subscribe function to never expose the Buttondown API key to the browser, so that the key stays secret.
61. As the Author, I want Lighthouse budgets enforced on preview deployments (LCP ≤ 1.5 s, CLS 0, total JS ≤ 60 KB), so that regressions are caught before merge.

### Brand

62. As a Reader, I want a consistent wordmark, teal accent and the tagline "How real systems are built." across pages, OG images and feeds, so that the site is recognisable.
63. As the Author, I want a generated SVG favicon/mark and OG badge, so that no bitmap logo assets exist.

## Implementation Decisions

### Shape

- One Go module. Two executables: the **builder** (`content/` → `.vercel/output/`) and the **subscribe function** (Vercel Go serverless). Everything else is packages behind those two.
- The builder is the single deep module. Its interface is: a content directory, a configuration (site URL, Verticals, author, flags such as `ads_enabled`, `draft_mode`), and an output directory. It returns a list of errors (gate failures) and warnings. Internals (Markdown parsing, Mermaid rendering, templ rendering, image encoding, feed generation, sitemap generation, search indexing, redirect map) are not exported.
- Rendering with **templ** components. Component families: layout (head, header, footer, theme toggle), article (Explainer, Brief), listing (home, Vertical hub, Briefs index, kind filter), page (Author, About, Contact, Editorial Policy, 404), partials (hero figure, TL;DR, Sources, trade-offs table, related links, newsletter form, ad slot, breadcrumb, diagram figure).
- No client-side framework. Client JS allowed: Pagefind (lazy-loaded on search focus), theme toggle (inline, ≤1 KB), newsletter form progressive enhancement (works as a plain form POST without JS), Plausible script. Ads/CMP scripts exist behind a flag and are off.
- System font stack, no webfonts. CSS is a single inlined stylesheet per page (well under the budget) or a single external file with far-future cache; decide at implementation, measured against the 50 KB HTML gate.

### Content model

- Content lives under `content/explainers/<vertical>/<slug>.md` and `content/briefs/<slug>.md`. Hero raster sits beside the Markdown with the same basename. Author data in `content/authors/<slug>.md`. Static pages in `content/pages/`.
- Frontmatter (YAML): `title`, `description`, `kind` (case-study|concept|comparison; Explainers only), `vertical` (Explainers only), `tags`, `published` (RFC 3339 with timezone), `reviewed` (RFC 3339; Explainers), `author` slug, `hero` (filename), `hero_alt`, `related` (slugs), `draft` (bool), `sources` optional list (title, url) used in addition to inline links. Briefs add `source_url` and `explainers` (slugs linked in body).
- Explainer skeleton is enforced structurally, not by prose heuristics: a TL;DR is a `## TL;DR` section immediately after the intro with 3–5 bullets; a Sources block is a `## Sources` section at the end with ≥1 link; a diagram is ≥1 Mermaid fence; a Comparison needs ≥1 Markdown table under a `## Trade-offs` heading. Sections under H2 longer than 250 words raise a warning (not failure).
- Brief skeleton: body ≤ 250 words, exactly the sections "What happened", "Why it matters", and a source link; ≥1 link to an Explainer slug.
- Verticals are configured in one list (slug, name, description, order). Launch list: distributed-systems, databases-storage, apis-networking, infrastructure-cloud, ai-ml-systems, languages-runtimes.

### URLs

- Explainer `/{vertical}/{slug}/`, Vertical hub `/{vertical}/` (paginated `/{vertical}/page/{n}/`), kind filter `/explainers/{kind}/`, Brief `/briefs/{slug}/`, Briefs index `/briefs/` (paginated), Author `/authors/{slug}/`, pages `/about/`, `/contact/`, `/editorial-policy/`, search `/search/`. Trailing slash canonical; Vercel config set to `trailingSlash: true` with directory-index output. No dates, no `.html`.
- Feeds: `/feed.xml`, `/{vertical}/feed.xml`, `/briefs/feed.xml` (full text, RSS 2.0 with content:encoded). Sitemaps: `/sitemap.xml`, `/news-sitemap.xml`. `robots.txt` references both. IndexNow key file at root; a build artefact lists changed URLs for a post-deploy ping step.
- Retired URLs: `/video/*`, `/community/*`, `/group/*` → 410 via Vercel routes in the Build Output config; `/home` → 301 `/`. Any other unknown path → the static 404 page.

### Structured data and head

- Explainer: `Article` JSON-LD (headline, description, image [3 ratios], datePublished, dateModified, author → Person with url and sameAs, publisher → Organization, mainEntityOfPage, articleSection = Vertical name, keywords = tags). Brief: `NewsArticle` with the same core plus `dateline`-free. All pages: `BreadcrumbList`. Home: `WebSite` (+ search action pointing at `/search/?q=`) and `Organization` (name "HastCode", url, logo as data-URI-free SVG URL, sameAs).
- No `FAQPage`, `HowTo`, `Speakable`, or `TechArticle` (research: removed, deprecated, no evidence, and unlisted by Google respectively).
- Every page: `<meta name="robots" content="max-image-preview:large">`, canonical, Open Graph + Twitter card with the hero, `lang="en"`, theme-color for both schemes, preload for the hero on article pages.
- `robots.txt`: `User-agent: *` allow all; explicit allow stanzas for `OAI-SearchBot`, `PerplexityBot`, `Bingbot`, `Googlebot`; `GPTBot` and `Google-Extended` left allowed (training opt-out is a separate later decision). Sitemap lines for both sitemaps.

### Images and diagrams

- Hero: the build validates width ≥ 1200 px, then writes AVIF and WebP at widths 480/800/1200/1600 and the three crops (16:9, 4:3, 1:1) for schema/OG, with a JPEG fallback. Output filenames are content-hashed; a build cache keyed on source hash skips re-encoding.
- Mermaid fences render to inline SVG at build time (headless renderer invoked by the builder; the exact renderer is an implementation choice but must run in Vercel's build environment). SVG uses `currentColor` and CSS variables so it adapts to dark mode; text stays as text for crawlers. Each diagram gets a `<figure>` with a caption from the fence's title line.

### Subscribe function

- `POST /api/subscribe` with form fields `email`, `source` (page path), honeypot field. Validates syntax, rejects honeypot hits, applies a simple in-memory per-IP rate limit (best-effort; serverless), forwards to Buttondown's subscriber API with the source as a tag, returns 303 to `/subscribed/` on success for no-JS clients or JSON for `Accept: application/json`. Buttondown base URL and key come from environment variables; the base URL is injectable for tests.

### Search

- Pagefind runs over the output directory after the build as a post-step; the `/search/` page loads its UI lazily. The builder ensures the index exists and fails if the post-step did not run in production builds.

### Analytics, ads, consent

- Plausible script included via one component, domain from config. Ad slot component renders a fixed-height placeholder (`aside` with `aria-hidden`) in two positions (after TL;DR, after the body) when `ads_enabled=false`; when true it renders the AdSense tag. The CMP (consent) component exists as a stub behind the same flag and is out of scope to wire.

### Brand

- Wordmark `hastcode` in CSS (system monospace stack, lowercase, teal dot or bracket mark). SVG mark generated as a file at build and referenced as favicon and in Organization logo. Accent `#0D9488` light / `#2DD4BF` dark; text `#111827` / `#F3F4F6`; background `#FFFFFF` / `#0B0F19`. Tagline "How real systems are built." in header subtitle, OG default, feed description, Organization description. Prose and schema name "HastCode".

### Deployment

- Vercel Pro project, Git integration on `main`, branch previews. Build command runs the Go builder then Pagefind. Output: `.vercel/output/` (Build Output API v3) with `config.json` carrying routes for 410s, the 301, trailing-slash handling and cache headers (`Cache-Control: public, max-age=0, must-revalidate` on HTML with `s-maxage` long + SWR on the CDN; immutable far-future on hashed assets). Cloudflare DNS-only for apex and `www` (www → apex 301 at Vercel).
- Cutover to apex only after ≥30 Explainers across all six Verticals and ≥20 Briefs exist (ADR-0001 consequence); until then the project runs on its Vercel preview domain.

## Testing Decisions

- A good test exercises the two agreed seams only and asserts on externally observable output. Tests never import or call internal packages directly.
- **Build seam**: each test creates a temporary content tree from small fixtures (a valid Explainer per kind, a valid Brief, an Author, a Vertical config), runs the builder into a temporary output directory, and asserts by parsing the output: HTML with an HTML parser (presence and content of TL;DR, Sources, diagram `<figure>`, breadcrumb, canonical, meta robots, OG image, ad placeholder, dates), JSON-LD by decoding the script block and checking fields, feeds and sitemaps by XML parsing, the news sitemap's 48-hour window by fixture dates, image variants by decoding output files and checking dimensions, routes by reading the Build Output `config.json`, size gate by byte count. Gate tests use invalid fixtures (missing TL;DR, narrow hero, duplicate slug, broken link, Comparison without table, oversize page) and assert the builder returns the expected error and writes nothing for that page. `draft: true` excluded in production mode, included in preview mode.
- **Subscribe seam**: `httptest.NewServer` stands in for Buttondown; tests send requests to the handler and assert status, redirect target or JSON body, the request the fake Buttondown received (email, tag), rejection of invalid email and honeypot, and that rate limiting trips after N requests from one IP.
- Prior art: none (greenfield). The first build test establishes the fixture and assertion helpers that all later tests reuse; keep them in one test-support package.
- Lighthouse CI budgets run in GitHub Actions against the preview URL; not a Go test.
- No tests for templ components, the Markdown parser, image encoders, or Mermaid rendering in isolation.

## Out of Scope

- Live ads, AdSense account approval, consent management (CMP/TCF) wiring, sponsorship pages, "sponsored Explainer" format.
- Multiple Authors or contributor workflow, comments, reactions, any community feature (ADR-0001).
- AI drafting tooling (separate project). Editorial Review scheduling beyond the build warning.
- Localisation, hreflang, non-English content.
- Newsletter sending, templates, or digest automation (Buttondown's job).
- ISR/SSR for any page; dynamic rendering of any kind beyond the subscribe function.
- Migration of any old community content.
- Bing PubHub submission and Apple News (unavailable to an India-based publisher).
- Training-crawler opt-out decisions (`GPTBot`, `Google-Extended`, `ClaudeBot`).

## Further Notes

- Target metrics from research: HTML ≤ 50 KB, LCP ≤ 1.5 s on mobile, CLS 0, TTFB "good" ≤ 800 ms everywhere (static edge serving makes this the default).
- Editorial rules from research, to be enforced by the build where structural and by the Author where not: answer first; declarative headings (question-style headings correlated negatively with citation); sections 120–180 words; quote and cite primary sources; include real numbers; Review every 90 days and only bump `reviewed` on real change.
- First business milestone is Raptive eligibility (25k pageviews/month, ≥50% Tier-1 traffic, domain age ≥6 months from cutover), not a revenue figure.
- The 50 KB gate is uncompressed HTML including inlined CSS and JSON-LD; if inlined CSS makes long Explainers exceed it, move CSS to one external hashed file rather than raising the gate.

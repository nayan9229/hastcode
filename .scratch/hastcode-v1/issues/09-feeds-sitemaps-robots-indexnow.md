# 09: Feeds, sitemaps, robots and IndexNow

**What to build:** Everything crawlers and readers pull. Full-text RSS 2.0 (with `content:encoded`) at `/feed.xml` (all pieces), `/{vertical}/feed.xml`, and `/briefs/feed.xml`. A standard sitemap of every canonical URL with lastmod. A news sitemap containing only Briefs published in the last 48 hours with publication name, language, publication_date and title. `robots.txt` allowing all, with explicit stanzas for `Googlebot`, `Bingbot`, `OAI-SearchBot` and `PerplexityBot`, and both sitemap lines. An IndexNow key file at the site root, and a build artefact listing URLs changed since the previous build so a post-deploy step can ping. Home page gains `WebSite` (with search action pointing at `/search/?q=`) and `Organization` JSON-LD.

**Blocked by:** 06 (Briefs), 07 (Listing pages and internal linking)

**Status:** ready-for-agent

- [ ] Build-seam tests parse each feed as XML: item count, order, full HTML body, absolute links, site tagline as channel description
- [ ] Sitemap contains every page from the fixture tree exactly once, no drafts, trailing-slash URLs
- [ ] News sitemap contains only Briefs within 48 h of the builder's clock; an older Brief fixture is excluded; Explainers never included; at most 1,000 entries
- [ ] robots.txt content asserted exactly; `Google-Extended` and `GPTBot` not disallowed
- [ ] IndexNow key file present at root; changed-URL artefact lists new/modified pages between two builds of differing fixtures
- [ ] `WebSite` and `Organization` JSON-LD on home with name "Hast Code", tagline, logo URL and search action

# 06: Briefs

**What to build:** The Author writes a Brief (short, dated commentary on something that just happened) and it publishes at `/briefs/{slug}/` with a `NewsArticle` JSON-LD block, plus a paginated Briefs index at `/briefs/`. The build enforces the Brief skeleton: body ≤ 250 words, exactly the sections "What happened" and "Why it matters", a `source_url`, and at least one link to an existing Explainer slug (a link to a non-existent Explainer is a broken internal link and fails the build). Briefs reuse the hero, head and breadcrumb work from earlier tickets.

**Blocked by:** 02 (Explainer editorial gates), 05 (Explainer head and structured data)

**Status:** ready-for-agent

- [ ] Brief frontmatter parsed per spec (`source_url`, `explainers`, plus shared fields); kind/vertical absent
- [ ] Build-seam tests: valid Brief → page at expected path with `NewsArticle` JSON-LD (headline, dates with tz, author Person, publisher, image ratios); over-length body, missing section, missing Explainer link, link to unknown slug → specific build errors
- [ ] Briefs index lists Briefs newest-first, 20 per page, with `/briefs/page/{n}/` for later pages and rel prev/next links
- [ ] Breadcrumb Home → Briefs → Brief
- [ ] Explainer links in the Brief body rendered as internal links with trailing-slash URLs

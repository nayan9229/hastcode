# 07: Listing pages and internal linking

**What to build:** Readers can browse. The home page shows the latest Explainers, the latest Briefs, and the six Verticals with descriptions. Each Vertical has a hub at `/{vertical}/` listing its Explainers newest-first, 20 per page with `/{vertical}/page/{n}/`. Kind filter pages at `/explainers/{kind}/` list Case Studies, Concepts or Comparisons across Verticals. Every Explainer page ends with "More in {Vertical}" (3 latest others in the Vertical) and "Related" (up to 3 by shared tags, excluding itself, then `related` frontmatter slugs first). Verticals come from one configuration list; adding one needs no code change. Drafts never appear in listings in production.

**Blocked by:** 06 (Briefs)

**Status:** ready-for-agent

- [ ] Build-seam tests with a fixture tree of several Explainers across Verticals and kinds, plus Briefs: home, each hub, each kind page exist with the expected items in the expected order
- [ ] Pagination produces correct page count, rel prev/next, and no empty pages
- [ ] "More in Vertical" and "Related" sections present with correct candidates; an Explainer never relates to itself
- [ ] Adding a Vertical to the config in a test fixture produces its hub and nav entry without code changes
- [ ] Breadcrumbs on hubs and kind pages
- [ ] Each listing item shows title, description, kind badge (Explainers), date, and hero thumbnail (480 px variant)

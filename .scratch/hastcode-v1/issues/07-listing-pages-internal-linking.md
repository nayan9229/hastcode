# 07: Listing pages and internal linking

**What to build:** Readers can browse. The home page shows the latest Explainers, the latest Briefs, and the six Verticals with descriptions. Each Vertical has a hub at `/{vertical}/` listing its Explainers newest-first, 20 per page with `/{vertical}/page/{n}/`. Kind filter pages at `/explainers/{kind}/` list Case Studies, Concepts or Comparisons across Verticals. Every Explainer page ends with "More in {Vertical}" (3 latest others in the Vertical) and "Related" (up to 3 by shared tags, excluding itself, then `related` frontmatter slugs first). Verticals come from one configuration list; adding one needs no code change. Drafts never appear in listings in production.

**Blocked by:** 06 (Briefs)

**Status:** resolved

- [x] Build-seam tests with a fixture tree of several Explainers across Verticals and kinds, plus Briefs: home, each hub, each kind page exist with the expected items in the expected order
- [x] Pagination produces correct page count, rel prev/next, and no empty pages
- [x] "More in Vertical" and "Related" sections present with correct candidates; an Explainer never relates to itself
- [x] Adding a Vertical to the config in a test fixture produces its hub and nav entry without code changes
- [x] Breadcrumbs on hubs and kind pages
- [x] Each listing item shows title, description, kind badge (Explainers), date, and hero thumbnail (480 px variant)

## Resolution notes

- Tests: `builder/listings_test.go` (hubs, kind pages, home, pagination, breadcrumbs, item content, new Vertical), `builder/related_test.go` (More in / Related), `builder/links_test.go` (broken-link gate). Drafts in listings are covered by the existing `TestDraft_AbsentFromProductionOutput`, which scans every output file.
- Home shows the 6 newest Explainers and 5 newest Briefs, the kind pages, and every Vertical with its `description` (now set in `site.yaml`). Kind pages are paginated like hubs (`/explainers/{kind}/page/{n}/`); there is no `/explainers/` page, and the kind-page breadcrumb is Home > Case Studies.
- Related: `related` slugs first (in frontmatter order), then by number of shared tags, ties newest first; never the Explainer itself. An unknown `related` slug is a `broken-link` error. Both lists are omitted when empty and carry `data-pagefind-ignore="all"`.
- The Briefs index items now use the shared `PieceList` partial (adds the 480 px thumbnail); `views.BriefItem` became `views.ListingItem`.
- Added scope: the site-wide broken internal link gate (spec story 31), `builder/step_render_link_gate.go`; see README "Internal links". `buildtest.NewTree()` now includes the transparency pages, since the footer links them from every page.
- Added scope: the phone-width check owed by ticket 11 (see its resolution notes); one CSS fix for breadcrumbs.

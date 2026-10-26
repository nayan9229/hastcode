# 10: Build-generated hero title card for every kind

**What to build:** Case Studies, Concepts, Comparisons and Briefs get their hero from the build, the way Repo Profiles already do, so no Piece needs a hand-made image (Author decision, 2026-10-26). The build draws a branded card at least 1600×900 (title, kind label, Vertical name for Explainers or "Brief", the `hastcode` wordmark, brand palette) in pure Go with the fonts the OG badge already uses, and feeds it through the existing hero pipeline (variants, crops, preload, `og:image`, JSON-LD image). The `hero`/`hero_alt` frontmatter is no longer used for these kinds: follow whatever the Repo Profile kind does today (forbidden or ignored), consistently. Alt text is generated from the title. Remove the placeholder hero files and hero fields from the four drafts on main (`netflix-single-tenant-titus`, `metastable-failures`, `http-connection-reuse`, `github-dependabot-fan-out`). Update the templates/writing guides only if they exist (ticket 08 is not done; note it there instead).

Spec: `../spec.md`. Prior art: the Repo Profile hero card and the OG badge renderer.

**Blocked by:** None (Author decision; independent of 07–09)

**Status:** ready-for-agent

- [ ] Build-seam tests: a Case Study, Concept, Comparison and Brief without `hero` get a generated card; decoded output dimensions ≥1200 px wide; JSON-LD image and og:image point at its crops; alt text set; cards differ when titles differ and are byte-identical across rebuilds (content-hashed, cached)
- [ ] Behaviour for a `hero` field on these kinds matches the Repo Profile rule and is tested
- [ ] The four drafts on main build in preview with no hero files; production build unchanged
- [ ] README documents the card

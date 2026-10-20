# 04: Hero image pipeline

**What to build:** The Author drops one raster hero beside the Markdown and the build produces everything the page, Discover, Open Graph and structured data need. The build fails if the hero is missing or narrower than 1200 px. It writes AVIF and WebP at widths 480/800/1200/1600 plus a JPEG fallback, and three crops (16:9, 4:3, 1:1) at ≥1200 px wide for schema and OG. Output filenames are content-hashed under a static assets path. A build cache keyed on the source hash skips re-encoding on later builds. The Explainer page renders the hero in a `<picture>` with `srcset`, explicit width and height, `alt` from frontmatter, and a preload hint in the head.

**Blocked by:** 01 (Walking skeleton)

**Status:** ready-for-agent

- [ ] Build-seam tests: valid hero → expected variant files exist with correct dimensions (decode and measure); narrow hero → build error; missing hero → build error
- [ ] `<picture>` markup present with AVIF, WebP and JPEG sources, width/height attributes, and alt text
- [ ] Preload link for the hero in the head of Explainer pages
- [ ] Second build with unchanged source does not re-encode (verified via cache directory or timing/marker)
- [ ] Hashed filenames change when the source changes
- [ ] Image cache directory location is configurable so deployment can persist it between builds

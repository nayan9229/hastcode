# 01: Hero variants and quality retune

**What to build:** Readers get the smallest hero that still looks right on their screen. Every hero is encoded at widths 480, 800, 1000, 1200 and 1472 px (never upscaled past the source) as AVIF with 4:4:4 chroma at quality 60 below 1200 px and 55 at 1200 px and above, plus the stdlib JPEG fallback at quality 82. WebP is gone from the encoder, the view model, the `<picture>` markup and the tests. A new 1200×630 JPEG crop feeds `og:image` and the Twitter card; the 16:9, 4:3 and 1:1 1200-px JPEG crops stay for JSON-LD. When a hero's 1200-px AVIF exceeds 150 KB the build warns (`hero-weight`), naming the piece. The encoder cache key carries a settings version so this retune re-encodes v1's cached variants exactly once.

Spec: `../spec.md` ("Images package"). Research with measurements: `docs/research/image-optimisation-targets.md`. Author photos and the OG badge keep their encoding.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] Build-seam tests decode the output and assert the five widths (and fewer when the source is narrower), AVIF + JPEG only, no WebP files or `<source type="image/webp">` anywhere
- [ ] AVIF chroma is 4:4:4 and quality tiers apply at the 1200-px boundary (assert via the encoder options the build reports, or by measuring bytes against a Q60 encode of the same fixture)
- [ ] `og:image` and `twitter:image` point at a 1200×630 JPEG; JSON-LD `image` still lists the three ratio crops
- [ ] A fixture hero whose 1200-px AVIF exceeds 150 KB raises `hero-weight` as a warning, not an error; a normal fixture raises nothing
- [ ] Changing the settings version invalidates the cache (second build re-encodes once, third build is cached)
- [ ] `make test` green; `make build MODE=preview` on real content passes; README image section updated

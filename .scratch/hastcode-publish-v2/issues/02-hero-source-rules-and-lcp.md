# 02: Hero source rules and LCP attributes

**What to build:** The Author drops a hero straight from a phone or an export tool and it comes out right. A JPEG with an EXIF orientation tag is rotated upright before resizing. A source carrying an ICC profile that is not sRGB raises a `hero-colour` warning naming the piece and telling the Author to export as sRGB (nothing is converted). A WebP source is rejected with the existing `hero-format` error, because Go's WebP decoders return wrong colour ranges. On the page, the hero `<img>` carries `fetchpriority="high"` and no `loading` attribute, so the LCP image is fetched first; the preload hint stays.

Spec: `../spec.md` ("Images package"). Research: `docs/research/image-optimisation-targets.md` sections 5–6 and 8.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] Build-seam test: a JPEG fixture with EXIF orientation 6 produces variants whose decoded dimensions and pixel layout are upright
- [ ] A PNG/JPEG fixture with an embedded non-sRGB ICC profile raises `hero-colour` as a warning; an sRGB or profile-less fixture raises nothing
- [ ] A `.webp` hero (and a WebP file with a `.png` name) fails with `hero-format` naming the piece
- [ ] The hero `<img>` has `fetchpriority="high"` and no `loading` attribute on Explainer and Brief pages; the preload link is still present
- [ ] `make test` green; README documents the source rules (PNG/JPEG, sRGB, ≥1200 px)

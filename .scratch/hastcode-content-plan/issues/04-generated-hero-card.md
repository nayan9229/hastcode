# 04: Generated hero card

**What to build:** A Repo Profile needs no sourced image. When its frontmatter has no `hero`, the builder draws a 1600×900 card from the Fact Pack (wordmark, repository full name, description truncated to fit, language, stars, license) on the brand palette using the bundled Go font family already in the module graph, and feeds it into the existing hero pipeline as if it were a file, so the page gets the five widths, AVIF and JPEG, the 1200×630 Open Graph crop, the preload and the high-priority image like every other Explainer. Alt text is generated as "{repo}: {description}". A Profile that does supply `hero` and `hero_alt` uses them and skips the card. `hero`/`hero_alt` stay required for every other kind.

Spec: `../spec.md` ("Page composition and gates": generated hero, User Stories 25, 26). Prior art: the hero tests (variants, OG crop, preload, cache) and the fixture hero generators.

**Blocked by:** 01 (Repo Profile kind renders from a Fact Pack)

**Status:** ready-for-agent

- [ ] A Profile without `hero` builds clean and the page offers the same widths, formats, OG crop and preload the hero tests require
- [ ] The card contains the wordmark, repository name, language, stars and license (verified by rendering the card and checking the pixels against a golden image, or by asserting the drawn text through the renderer's test hook)
- [ ] Alt text is set on the `<img>` and in the Article image metadata
- [ ] A Profile with `hero` uses it; a Case Study without `hero` still fails as today
- [ ] The card changes when stars or the latest release change, so file names change and caches bust; an unchanged Fact Pack re-uses the cache
- [ ] Card renders deterministically (same bytes for the same Fact Pack) so builds stay reproducible

---
status: accepted
date: 2026-10-25
supersedes: ADR-0003
amended-by: ADR-0008 (Google Analytics 4 replaces Plausible; the rules for Previews stand)
---

# Inline the stylesheet, embed the favicon, keep analytics out of Previews

Lighthouse CI on the first Cloudflare Preview failed the spec's mobile LCP budget (1.5 s) on every page at 1.58–1.62 s. The LCP element was a text paragraph everywhere; server time (~0.16 s), fonts and main-thread work were not the cause. Measured locally with the same Lighthouse settings (mobile, simulated throttling), the time came from requests that start before LCP, each of which the simulation makes LCP wait for:

| Change, cumulative | Home LCP | Search LCP | Largest Repo Profile LCP |
|---|---|---|---|
| Before (external CSS, analytics tag, favicon file) | 1.38 s | 1.38 s | 1.68 s |
| Plausible omitted from preview builds | 1.05 s | 0.90 s | 1.20 s |
| + stylesheet inlined | 0.90 s | 0.90 s | 1.20 s |
| + favicon embedded as a data: URI | 0.78 s | 0.77 s | 1.05 s |

CI ran about 0.2 s slower than these local numbers on the same pages.

We decided three things:

1. **Inline the stylesheet once, as the single `<style>` in every page's `<head>`** (superseding ADR-0003), and make the 50 KB size gate measure each page **without that shared block**, because the gate bounds content, not design.
2. **Embed the favicon as a `data:image/svg+xml` URI.** `/mark.svg` is still written; it is the Organization logo.
3. **Preview builds (Previews, `make dev`) carry no analytics.** Their visits are not Readers', and Lighthouse CI, which only tests Previews, then measures the page's own code. Production keeps Plausible's standard deferred tag: real browsers do not wait on a deferred script for LCP, only the simulation does.

## Considered options

- **Raise the LCP budget to 2.0 s**: hides the cause rather than removing it.
- **Load Plausible after the `load` event via a small inline loader**: built and tested, then reverted. On a fast connection `load` fires before the measured LCP, so the simulation still counts the request; delaying further with timers would drop real pageviews.
- **Critical CSS inline, rest deferred**: risks layout shift against the zero-CLS budget.

## Consequences

- Every page carries ~14 KB uncompressed (~3 KB brotli) of CSS and ~1 KB of favicon; repeat views no longer get the stylesheet from cache. Content can still never exceed 50 KB (tests measure with `buildtest.ContentSize`).
- No `/assets/site.*.css` is written; `/assets/` still holds images and the OG badge, so its immutable rule stands.
- A Lighthouse run against Production (PageSpeed Insights) will report a slower lab LCP than CI because of the analytics tag; real-user LCP is the number that matters there.
- Any new request a page makes before first paint (a font, a tracking pixel, a preconnect-driven fetch) will show up as LCP in CI: that is now the known failure mode.

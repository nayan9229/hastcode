# 12: Size gate, ad-slot placeholders and analytics

**What to build:** Monetisation readiness without monetising. The build fails if any output HTML exceeds 50 KB uncompressed. Explainer and Brief layouts contain two fixed-height ad slot components (after the TL;DR, after the body) that render as empty, `aria-hidden` placeholders when `ads_enabled` is false and the AdSense tag when true; placeholders reserve the exact height so enabling ads causes no layout shift and ad-blocked Readers see no gap. The Plausible analytics script is included via one component with the domain from configuration, and omitted entirely when the configuration leaves it blank (tests and local builds).

**Blocked by:** 11 (Design system, dark mode and brand)

**Status:** resolved

- [x] Build-seam test: a fixture page padded past 50 KB fails the build naming the page; a normal page passes
- [x] With `ads_enabled=false`: two placeholder elements with fixed height present, no third-party script; with `ads_enabled=true`: AdSense script and two ad units present
  - Revised after code review (spec c3): with `ads_enabled=false` no ad-slot element is rendered and no space is reserved; with `ads_enabled=true` the two units sit in fixed-height boxes (Explainer: after the TL;DR and the body; Brief: after "What happened", `data-ad-position="after-what-happened"`, and the body).
- [x] Plausible script present only when a domain is configured; `defer`, no cookies
- [x] Total client JS per page stays within the 60 KB budget from the spec (asserted by summing script sizes in output)

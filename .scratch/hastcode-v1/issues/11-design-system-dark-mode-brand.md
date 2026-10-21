# 11: Design system, dark mode and brand

**What to build:** The site looks like hastcode. One stylesheet, system font stack, no webfonts, mobile-first layout with a 16 px gutter, readable measure for long Explainers, header with the `hastcode` wordmark (CSS, lowercase, with a teal mark) and tagline "How real systems are built.", nav for Verticals, Briefs and search, footer with About/Contact/Editorial Policy/feeds. Dark and light palettes from the spec via CSS variables, following `prefers-color-scheme` with a ≤1 KB inline toggle that persists the override. A generated SVG mark used as favicon and as the Organization logo. Diagrams and code blocks readable in both modes. Every page works with JavaScript disabled.

**Blocked by:** 01 (Walking skeleton)

**Status:** resolved

- [x] Build-seam tests: header contains wordmark text, tagline and nav links; footer contains transparency links and feed links; favicon link points at the SVG mark in output; theme-color metas for both schemes present
- [x] Palette values from the spec defined as CSS variables; dark variants under `prefers-color-scheme: dark` and under an explicit data attribute for the override
- [x] Inline toggle script ≤ 1 KB and the only script in the page apart from analytics/search (asserted by counting `<script>` elements)
- [x] CSS delivery decided and documented (inline vs single hashed file) with the 50 KB HTML gate in mind
- [ ] Manual check on a phone-width viewport: no horizontal scroll on home, hub, Explainer with table and diagram, Brief

## Resolution notes

- Tests: `builder/design_system_test.go`. CSS delivery: one hashed external file, see `docs/adr/0003-css-as-one-hashed-external-file.md` (stylesheet ~7 KB; theme script 676 bytes).
- Manual check (2026-10-21): built a local Explainer with a 4-column Comparison table, a long Go code block, long inline code and a long URL; served it and loaded it in headless Chrome inside 320, 360 and 375 px iframes, light and dark (`data-theme`), and at 1280 px. `documentElement.scrollWidth == clientWidth` at every phone width: no page-level horizontal scroll. Only the Verticals nav strip, the table and the `<pre>` overflow, each scrolling inside itself. The raw-HTML diagram stand-in was stripped by the Markdown renderer, so diagram layout depends on CSS rules only: `figure svg` capped at 100% width, `color: var(--fg)`.
- Not met yet: home, Vertical hub and Brief pages do not exist yet (tickets 06/07). Re-run the same phone-width check when they land. The Mermaid `<figure>` (ticket 03) also needs a visual check in both themes.

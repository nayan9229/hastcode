# 11: Design system, dark mode and brand

**What to build:** The site looks like hastcode. One stylesheet, system font stack, no webfonts, mobile-first layout with a 16 px gutter, readable measure for long Explainers, header with the `hastcode` wordmark (CSS, lowercase, with a teal mark) and tagline "How real systems are built.", nav for Verticals, Briefs and search, footer with About/Contact/Editorial Policy/feeds. Dark and light palettes from the spec via CSS variables, following `prefers-color-scheme` with a ≤1 KB inline toggle that persists the override. A generated SVG mark used as favicon and as the Organization logo. Diagrams and code blocks readable in both modes. Every page works with JavaScript disabled.

**Blocked by:** 01 (Walking skeleton)

**Status:** resolved

- [x] Build-seam tests: header contains wordmark text, tagline and nav links; footer contains transparency links and feed links; favicon link points at the SVG mark in output; theme-color metas for both schemes present
- [x] Palette values from the spec defined as CSS variables; dark variants under `prefers-color-scheme: dark` and under an explicit data attribute for the override
- [x] Inline toggle script ≤ 1 KB and the only script in the page apart from analytics/search (asserted by counting `<script>` elements)
- [x] CSS delivery decided and documented (inline vs single hashed file) with the 50 KB HTML gate in mind
- [x] Manual check on a phone-width viewport: no horizontal scroll on home, hub, Explainer with table and diagram, Brief

## Resolution notes

- Tests: `builder/design_system_test.go`. CSS delivery: one hashed external file, see `docs/adr/0003-css-as-one-hashed-external-file.md` (stylesheet ~7 KB; theme script 676 bytes).
- Manual check (2026-10-21): built a local Explainer with a 4-column Comparison table, a long Go code block, long inline code and a long URL; served it and loaded it in headless Chrome inside 320, 360 and 375 px iframes, light and dark (`data-theme`), and at 1280 px. `documentElement.scrollWidth == clientWidth` at every phone width: no page-level horizontal scroll. Only the Verticals nav strip, the table and the `<pre>` overflow, each scrolling inside itself. The raw-HTML diagram stand-in was stripped by the Markdown renderer, so diagram layout depends on CSS rules only: `figure svg` capped at 100% width, `color: var(--fg)`.
- Re-run (2026-10-21, ticket 07): built a fixture site (26 Explainers over 3 Verticals and all kinds, with real Mermaid diagrams and Comparison tables, deliberately long titles, an unbroken 60-character word in descriptions and a Brief title; 3 Briefs) with the real `site.yaml`, served it over HTTP and drove headless Chrome 155 through the DevTools protocol with mobile device emulation at 320, 360 and 375 px, `prefers-color-scheme` light and dark. Pages: home, Vertical hub, kind page, Briefs index, a Brief, a Concept Explainer (diagram), a Comparison (table + diagram), Author, About, 404. Check: `scrollWidth <= clientWidth` on `<html>`/`<body>`, listing elements past the right edge outside scroll containers. One failure found and fixed: a Brief breadcrumb with an unbreakable word pushed 320 px pages to 334 px (flex items do not shrink below their longest word); `.breadcrumbs li { min-width: 0; overflow-wrap: anywhere; }`. After the fix: no page-level horizontal scroll on any page at any width in either scheme; only the Verticals strip, tables and `<pre>` scroll inside themselves. The Mermaid figure's colours in dark mode were not inspected by eye beyond fitting the column.

# 11: Design system, dark mode and brand

**What to build:** The site looks like hastcode. One stylesheet, system font stack, no webfonts, mobile-first layout with a 16 px gutter, readable measure for long Explainers, header with the `hastcode` wordmark (CSS, lowercase, with a teal mark) and tagline "How real systems are built.", nav for Verticals, Briefs and search, footer with About/Contact/Editorial Policy/feeds. Dark and light palettes from the spec via CSS variables, following `prefers-color-scheme` with a ≤1 KB inline toggle that persists the override. A generated SVG mark used as favicon and as the Organization logo. Diagrams and code blocks readable in both modes. Every page works with JavaScript disabled.

**Blocked by:** 01 (Walking skeleton)

**Status:** ready-for-agent

- [ ] Build-seam tests: header contains wordmark text, tagline and nav links; footer contains transparency links and feed links; favicon link points at the SVG mark in output; theme-color metas for both schemes present
- [ ] Palette values from the spec defined as CSS variables; dark variants under `prefers-color-scheme: dark` and under an explicit data attribute for the override
- [ ] Inline toggle script ≤ 1 KB and the only script in the page apart from analytics/search (asserted by counting `<script>` elements)
- [ ] CSS delivery decided and documented (inline vs single hashed file) with the 50 KB HTML gate in mind
- [ ] Manual check on a phone-width viewport: no horizontal scroll on home, hub, Explainer with table and diagram, Brief

---
status: accepted
date: 2026-10-21
---

# CSS ships as one hashed external file, not inlined

The spec left CSS delivery open: one stylesheet either inlined into every page or served as a single external file with a far-future cache, decided against the 50 KB uncompressed-HTML gate (ticket 12). We decided on **one external file at `/assets/site.<sha256-12>.css`**, linked from every page's `<head>`. The builder writes it from `internal/views/site.css` plus the colour tokens generated from `internal/views/brand.go`; the file name is the content hash, so it can be served `immutable` and changes bust the cache automatically.

## Why

- **The 50 KB gate belongs to content.** Long Explainers carry inline Mermaid SVG (verbose: ~5–20 KB per diagram), JSON-LD and ad-slot markup. The stylesheet is ~7 KB minified-ish and would be repeated in every page, taking ~14% of the budget. The spec says to move CSS out rather than raise the gate if it ever bites; doing it from day one means no Explainer fails the gate because of styling.
- **Readers move between pages.** A Brief links into Explainers and Explainers link to each other ("More in Vertical", "Related"); every page after the first reuses the cached file.
- **First-view cost is small.** Same origin, same HTTP/2 connection as the HTML, served from Vercel's edge, ~2 KB compressed. It is render-blocking, which also means no flash of unstyled content.
- **No flash of the wrong theme.** The only inline code is the ≤1 KB theme script, which sits before the stylesheet link and sets `data-theme` before first paint.

## Considered options

- **Inline `<style>` per page**: one fewer request on a cold first view, but spends the HTML budget on bytes the browser already has from the previous page, and would need re-deciding as soon as a diagram-heavy Explainer nears 50 KB.
- **Inline critical CSS + async rest**: best of both for LCP, but needs a critical-CSS extraction step and two code paths for a site whose whole stylesheet is already small. Not worth it.

## Consequences

- Ticket 10 must give `/assets/*` `Cache-Control: public, max-age=31536000, immutable`.
- The stylesheet is not counted by the HTML size gate; keep it small anyway (no framework, system fonts, no webfonts, no `url()` resources).
- Colours live only in `brand.go` (palette structs). `site.css` uses the tokens and never hard-codes a colour.

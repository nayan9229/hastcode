---
status: accepted
date: 2026-10-20
---

# Static-first Go build on Vercel, not Go SSR

The brief said "Go SSR on Vercel". Research (see `docs/research/vercel-go-ssr-viability.md`) found Vercel's Go runtime is Beta, single-region (`iad1` by default), excluded from Fluid Compute, and has no streaming, so every uncached India/EU request would cross continents and threaten the 800 ms TTFB target. We decided the Go program is a **build tool**: it reads Markdown from `content/`, renders `templ` templates, and writes finished HTML, feeds, sitemaps, and image variants to `.vercel/output/` via the Build Output API v3, so every page is served from every edge with no function in the path. One Go serverless function remains, for the newsletter signup POST. Publishing is `git push` to `main`.

## Considered options

- **Go SSR function + CDN cache headers**: works, but cache is per-region and best-effort; first hits are slow exactly where our India/EU readers are.
- **Fly.io multi-region Go**: first-class Go and true multi-region, but no Mumbai region, more ops, and loses Vercel's Git-integration previews.
- **Cloudflare Workers**: Go only via Wasm; rejected.

## Consequences

- Page count is bounded by build time, not request latency. At ~15 pieces/week this is seconds for years.
- Anything dynamic (search, personalisation) must be client-side or go through the single function. Search is Pagefind, a static index.
- Vercel Hobby forbids ad-monetised sites; the project runs on Pro from the first deploy.
- Cloudflare stays DNS-only (grey cloud); proxying through Cloudflare would stack two CDNs and break Vercel's cache.

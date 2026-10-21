# 10: Routes: 410 map, redirects, trailing slash and cache headers

**What to build:** The Build Output `config.json` makes the old Genuin community disappear cleanly (ADR-0001) and makes the new site cache correctly. `/video/*`, `/community/*` and `/group/*` return 410 Gone. `/home` 301s to `/`. Requests without a trailing slash 308 to the slash form; `index.html` is never exposed in URLs. HTML responses carry `Cache-Control` that keeps browsers revalidating while the CDN serves long `s-maxage` with stale-while-revalidate; hashed assets are immutable far-future. Unknown paths fall through to the 404 page with status 404.

**Blocked by:** 08 (Author page, static pages and 404)

**Status:** resolved

- [x] Build-seam test reads `config.json` and asserts the 410 routes, the `/home` redirect, trailing-slash handling, the 404 fallback, and the header rules for HTML vs hashed assets
- [x] Route order is correct: 410s and redirects before the filesystem handler, 404 last
- [x] A local serving command (or the Vercel CLI in a documented check) demonstrates each status on a sample path; documented in the README
- [x] No page in output is reachable at two URLs (with and without slash) without a redirect between them

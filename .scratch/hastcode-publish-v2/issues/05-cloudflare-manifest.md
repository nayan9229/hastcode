# 05: Cloudflare manifest and deploy-adapter rule

**What to build:** The one build now also serves the Production Target. Beside the static tree and the Vercel manifest, the builder writes a Cloudflare Workers manifest: the Workers config declaring the static assets directory with trailing-slash forcing and the 404 page as not-found handling, `run_worker_first` restricted to the signup path, and `main` pointing at the Worker module at the agreed path `worker/index.js` (ticket 06 writes that file; it may not exist yet); a `_headers` file (immutable one-year cache for hashed assets, revalidate for HTML); and a `_redirects` file (`/home` → `/` 301). Retired community paths get no rule and fall through to 404 (ADR-0001 as amended). Define and document the deploy-adapter rule: which generated files belong to which Target, so the Vercel upload never ships Cloudflare's files and vice versa.

Spec: `../spec.md` ("One build, per-Target manifests"). Research: `docs/research/cloudflare-r2-and-deploy.md` sections 5–6 (Workers Static Assets config keys, `_headers`/`_redirects` support).

**Blocked by:** None (can start immediately)

**Status:** resolved

- [x] Build-seam tests read the Workers config and assert the assets directory, trailing-slash and 404 handling, `run_worker_first` containing only the signup path, and `main`
- [x] `_headers` asserts immutable caching for the hashed-assets path and revalidation for HTML; `_redirects` contains exactly the `/home` rule
- [x] No rule for `/video`, `/community`, `/group`; an unknown path resolves to the 404 page per the config
- [x] The adapter rule is implemented as a documented list (or function) both deploy steps will use, with a test that the Vercel set excludes Cloudflare files and the Cloudflare set excludes `config.json`
- [x] README "Deploying" documents the Cloudflare manifest and the adapter rule

## Comments

**Implementer (05):** Done on `v2/05-cf-manifest`.

- Config: `<out>/wrangler.json` (plain JSON), written by `builder/step_write_cloudflare_manifest.go` (write phase, order 105) via `internal/cloudflare`. Deploy with `wrangler deploy --config .vercel/output/wrangler.json`; relative paths resolve against the config file, so `main` is `../../worker/index.js` for the default `OUT` (anchored on `-tools`, the project root). No `routes` (no automatic domain attach); `workers_dev` and `preview_urls` true.
- `_headers`: `/assets/*` immutable; `/`, `/*/`, `/*.html` revalidate (non-overlapping, since Workers joins duplicate headers); `https://:worker.:account.workers.dev/*` gets `X-Robots-Tag: noindex`. `_redirects`: `/home / 301`, `/home/ / 301`.
- Adapter rule: `internal/deploy` (`Files`, `Apply`, `ParseTarget`), CLI `cmd/adapt`, `make adapt TARGET=cloudflare|vercel`.
- `buildtest.URLPaths` skips `_headers`/`_redirects` (not URLs on either Target); the production noindex test whitelists `static/_headers` (host-scoped rule).
- Verified by hand with wrangler 4.149.0: `wrangler deploy --dry-run --config .vercel/output/wrangler.json` bundles the Worker and reads the assets; `wrangler dev` on the same config: `/home` and `/home/` 301 `/`; `/video/abc`, `/video`, `/community/x`, `/group/` 404 with the 404 page; `/about` and `/about/index.html` 307 `/about/`; `/_headers`, `/_redirects` 404; hashed CSS `public, max-age=31536000, immutable`; HTML and feeds `public, max-age=0, must-revalidate`; `POST /api/subscribe` reaches the Worker (bad email 400).

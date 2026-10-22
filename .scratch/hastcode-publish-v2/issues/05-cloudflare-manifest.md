# 05: Cloudflare manifest and deploy-adapter rule

**What to build:** The one build now also serves the Production Target. Beside the static tree and the Vercel manifest, the builder writes a Cloudflare Workers manifest: the Workers config declaring the static assets directory with trailing-slash forcing and the 404 page as not-found handling, `run_worker_first` restricted to the signup path, and `main` pointing at the Worker module at the agreed path `worker/index.js` (ticket 06 writes that file; it may not exist yet); a `_headers` file (immutable one-year cache for hashed assets, revalidate for HTML); and a `_redirects` file (`/home` → `/` 301). Retired community paths get no rule and fall through to 404 (ADR-0001 as amended). Define and document the deploy-adapter rule: which generated files belong to which Target, so the Vercel upload never ships Cloudflare's files and vice versa.

Spec: `../spec.md` ("One build, per-Target manifests"). Research: `docs/research/cloudflare-r2-and-deploy.md` sections 5–6 (Workers Static Assets config keys, `_headers`/`_redirects` support).

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] Build-seam tests read the Workers config and assert the assets directory, trailing-slash and 404 handling, `run_worker_first` containing only the signup path, and `main`
- [ ] `_headers` asserts immutable caching for the hashed-assets path and revalidation for HTML; `_redirects` contains exactly the `/home` rule
- [ ] No rule for `/video`, `/community`, `/group`; an unknown path resolves to the 404 page per the config
- [ ] The adapter rule is implemented as a documented list (or function) both deploy steps will use, with a test that the Vercel set excludes Cloudflare files and the Cloudflare set excludes `config.json`
- [ ] README "Deploying" documents the Cloudflare manifest and the adapter rule

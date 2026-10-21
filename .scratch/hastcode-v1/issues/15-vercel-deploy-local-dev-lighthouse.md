# 15: Vercel deployment, local dev loop and Lighthouse CI

**What to build:** Pushing to `main` deploys the full site; every branch gets a preview. The Vercel project (Pro, since ads are commercial use) runs the Go builder then Pagefind as its build command, persists the image cache between builds, serves the Go subscribe function, and has the environment variables documented. A single local command builds and serves the output with the routes from ticket 10 so the Author can check pages before pushing. GitHub Actions runs Lighthouse CI against each preview URL with budgets (mobile LCP ≤ 1.5 s, CLS 0, total JS ≤ 60 KB) and fails the check on regression. The repo README covers the whole loop. Domain cutover (apex + `www` → apex, Cloudflare DNS-only) is a human step: produce a `/wizard` script for it rather than doing it; cutover itself waits for the content threshold in ADR-0001.

**Blocked by:** 10 (Routes), 12 (Size gate, ad slots, analytics), 13 (Site search), 14 (Newsletter subscribe)

**Status:** ready-for-human

- [x] Vercel configuration in repo: build command, output directory, Go function, Node/Pagefind availability, cache directory wiring
- [ ] First deploy to the project's Vercel preview domain succeeds and serves home, an Explainer, a Brief, feeds, sitemaps, robots, 410 on an old path, and a working subscribe POST against a test Buttondown key
- [x] Local build-and-serve command documented and working, honouring 410/301/trailing-slash rules
- [x] Lighthouse CI workflow with budgets runs on PRs and reports against the preview URL
- [x] Wizard script for domain cutover: Vercel domain add, Cloudflare DNS-only records, `www` redirect, verification; not executed
- [x] README: publish flow, draft/preview flow, env vars, how to add a Vertical, how to run tests

## Remaining human steps

Everything outward-facing is in `scripts/setup-deploy.sh` (run from the repo root; safe to re-run). See README "Deploying".

1. Merge `integration/hastcode-v1` into `main`.
2. Wizard stages 1-9: tools and logins; GitHub repo; Vercel team on Pro; project + `vercel link` (Node 24.x, Git **not** connected); Buttondown production + test newsletter keys (plan with tags); Vercel env var `BUTTONDOWN_API_KEY` (Production / Preview); Protection Bypass for Automation secret; Vercel token + GitHub secrets (`VERCEL_TOKEN`, `VERCEL_ORG_ID`, `VERCEL_PROJECT_ID`, `VERCEL_AUTOMATION_BYPASS_SECRET`, variable `INDEXNOW_ENABLED=false`); push `main` and open a draft PR.
3. Stage 10 checks the first preview: home, an Explainer, a Brief, feeds, sitemaps, robots, IndexNow key, 410 (with 404 body), `/home` 301, 308s, query string kept on redirects, 404, HTML cache headers with `Vercel-CDN-Cache-Control` stripped, `/assets` immutable, subscribe POST (JSON 200 and no-JS 303) against the test key. Then confirm in the Actions logs that the build ran and the cache was restored on the second run, and that Lighthouse CI is green on the PR. When all pass, tick the "First deploy" box. (An Explainer and a Brief can only be checked once one is published.)
4. Stages 11-13 (domain cutover) only after ADR-0001's threshold: ≥ 30 Explainers across all six Verticals and ≥ 20 Briefs.

## Implementation notes

- Deploy path: GitHub Actions running `vercel build` + `vercel deploy --prebuilt` (`.github/workflows/production.yml`, `preview.yml`); Vercel Git deployments off in `vercel.json`. Why: Vercel's build image has no Go on PATH for the build command, the IndexNow artefact must be POSTed from the build machine after the deploy, and the build cache is kept with `actions/cache`.
- A local `vercel build` (CLI 65.0.0, hand-written `.vercel/project.json`, no login, nothing deployed) found that `@vercel/go` renames the root module to `handler`, so `api/subscribe.go` could not import `internal/subscribe` ("use of internal package not allowed"). The package moved to `subscribe/`. After that, the merged Build Output had `functions/api/subscribe.func`, the builder's routes first and in order, and `indexnow.json` intact; `MODE=preview` reached `make build` through `vercel build`.
- Lighthouse CI run locally (lhci 0.15.1, same `lighthouserc.json`, against `make dev`) found CLS 0.0008 on every page: revealing the theme toggle widened the header nav, because the global `[hidden]{display:none !important}` collapsed it. Fixed in `site.css` (`.theme-toggle[hidden] { display: inline-block !important; visibility: hidden; }`) with a Build-seam test. After the fix: CLS 0, LCP about 750 ms, script about 1.9 KB on home, /search/, /briefs/ and /about/.

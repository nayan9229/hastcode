# 15: Vercel deployment, local dev loop and Lighthouse CI

**What to build:** Pushing to `main` deploys the full site; every branch gets a preview. The Vercel project (Pro, since ads are commercial use) runs the Go builder then Pagefind as its build command, persists the image cache between builds, serves the Go subscribe function, and has the environment variables documented. A single local command builds and serves the output with the routes from ticket 10 so the Author can check pages before pushing. GitHub Actions runs Lighthouse CI against each preview URL with budgets (mobile LCP ≤ 1.5 s, CLS 0, total JS ≤ 60 KB) and fails the check on regression. The repo README covers the whole loop. Domain cutover (apex + `www` → apex, Cloudflare DNS-only) is a human step: produce a `/wizard` script for it rather than doing it; cutover itself waits for the content threshold in ADR-0001.

**Blocked by:** 10 (Routes), 12 (Size gate, ad slots, analytics), 13 (Site search), 14 (Newsletter subscribe)

**Status:** ready-for-agent

- [ ] Vercel configuration in repo: build command, output directory, Go function, Node/Pagefind availability, cache directory wiring
- [ ] First deploy to the project's Vercel preview domain succeeds and serves home, an Explainer, a Brief, feeds, sitemaps, robots, 410 on an old path, and a working subscribe POST against a test Buttondown key
- [ ] Local build-and-serve command documented and working, honouring 410/301/trailing-slash rules
- [ ] Lighthouse CI workflow with budgets runs on PRs and reports against the preview URL
- [ ] Wizard script for domain cutover: Vercel domain add, Cloudflare DNS-only records, `www` redirect, verification; not executed
- [ ] README: publish flow, draft/preview flow, env vars, how to add a Vertical, how to run tests

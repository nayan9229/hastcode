# 09: Wizard, README and first-run checklist

**What to build:** The Author can set up and run the whole publish path without this conversation. `scripts/setup-deploy.sh` is rewritten with stages for: tools and logins; GitHub repository; Cloudflare account, Workers project, R2 bucket with `images.hastcode.com` custom domain, two API tokens (Workers edit; R2 object read/write scoped to the bucket) and Worker secrets; Buttondown keys; Vercel Hobby project, `qa.hastcode.com`, Vercel Authentication; GitHub secrets and the IndexNow gate variable; first `publish vl` with the QA checklist; first branch Preview with its checklist; the Production DNS flip (apex and `www` to the Worker) behind the ADR-0001 content gate; first `publish cf` with the live checklist drawn from the spec's "Further Notes". Vercel Pro, Vercel domain and bypass stages are removed. README is updated end to end (publish, QA, Preview, Image Store modes, prune, secrets, Worker). Call the Skill tool with `wizard` for the script; the stages above are pre-approved.

Spec: `../spec.md` ("Wizard and docs", "Further Notes"). ADR-0004.

**Blocked by:** 04 (Prune command), 08 (Workflows and the publish command)

**Status:** ready-for-human

- [ ] Wizard runs through every stage in dry conversation (no real accounts touched here), writes `.env` and GitHub secrets via the library helpers, and `confirm`s before the DNS flip and before `publish cf` — traced statically only: `bash -n` and shellcheck clean (only the library's known SC1087), library above the STAGES marker byte-identical to the skill template, every `set_secret`/`set_var` name matches a `secrets.*`/`vars.*` reference in the workflows, `confirm` gates before every push, deploy, Publish (`publish.sh vl`/`cf`), the DNS flip and `--delete`. Executing it with stubbed CLIs in a pty was refused by the session's permission classifier, so the first real run is also the first execution.
- [x] Each checklist item maps to a spec story or a research open question (each line in the script carries `(story N)`, `(research QN)` or `(Further Notes)`)
- [x] README has no remaining reference to Vercel Pro, Vercel production, push-to-main deploys or the Vercel preview workflow
- [x] First deploy criteria that only a human can complete are left unticked and the ticket ends `ready-for-human` with a "Remaining human steps" list

## Implementation notes

- Wizard: `scripts/setup-deploy.sh`, 19 stages, `--from N` resumes (values come back from `.env`). Stages: 1 tools and logins; 2 GitHub repo and Actions settings; 3 Cloudflare account, zone, workers.dev; 4 R2 bucket + `images.hastcode.com`; 5 Workers API token; 6 R2 API token; 7 Buttondown keys; 8 first Worker deploy (local preview-mode bootstrap) + Worker secret; 9 what survives a `wrangler deploy` (throwaway `deploy-check.hastcode.com`, a dashboard var, the secret); 10 Vercel Hobby project; 11 `qa.hastcode.com` + Vercel Authentication; 12 GitHub secrets + `INDEXNOW_ENABLED=false`; 13 push main; 14 first `publish vl` + QA checklist; 15 first branch Preview + checklist (`check-worker.sh` 13 cases locally, safe subset live); 16 ADR-0001 content gate + first `publish cf` while hastcode.com still points at the old host (so the flip never exposes the drafts bootstrap); 17 DNS flip; 18 go live (IndexNow on, `hastcode-cloudflare-*` caches deleted, `publish cf`, live checklist); 19 first `prune-images --delete` by hand against a dummy object + R2 billing checks.
- Vercel Authentication on Hobby: Vercel's Deployment Protection docs (page updated 2026-09-15, read 2026-10-22) say scope "All Deployments" is available on all plans and covers production domains; "Standard Protection" does not. The wizard asks for All Deployments and stage 14 checks an anonymous request is refused.
- Wrangler 4.149.0 `deploy --help`: without `--keep-vars`, Wrangler deletes all vars before setting the config's. The Worker therefore has no dashboard vars; stage 9 confirms it live, and that the secret survives.
- Buttondown keys: the Worker gets the test newsletter's key until the DNS flip (Previews share Worker secrets), then the real one; Vercel (QA) and `.env` get the test key.
- The CI token is saved in `.env` as `CF_WORKERS_API_TOKEN`, not `CLOUDFLARE_API_TOKEN`, so it never overrides `wrangler login` locally; it is set as the GitHub secret `CLOUDFLARE_API_TOKEN`.
- README: "First-time setup and first run" rewritten (stage table, what each stage writes, the live-only checks and their stages); stale v1 references fixed (Preview env key, wizard stage numbers, news-sitemap rebuild wording, Pagefind caching, build command). Makefile comment on `MODE` updated.

## Remaining human steps

Run `scripts/setup-deploy.sh` (from the repo root, on an up-to-date `main` that holds `integration/publish-v2`). Each item below is a wizard stage; failed or unconfirmed items are listed at the end of the run.

- [ ] Stages 1–13: accounts, tokens, R2 bucket + custom domain, Buttondown keys, first Worker deploy, Vercel Hobby project + `qa.hastcode.com` + Vercel Authentication (All Deployments), GitHub secrets and `INDEXNOW_ENABLED`, push `main`
- [ ] Stage 9: a dashboard custom domain survives `wrangler deploy` without routes (if not: add `routes` with `custom_domain: true` to the builder's `wrangler.json` before the flip); dashboard vars do not survive (expected); the Worker secret does
- [ ] Stage 14: first `publish vl`; R2 single PUT with checksums "when required" (research Q1); `images.hastcode.com` sends the object's `Cache-Control` and edge-caches it (research Q4); the lightweight tag resolves to the tagged commit (research Q7); Vercel Authentication refuses anonymous requests to `qa.hastcode.com`; QA noindex header, disallow-all robots, canonical on hastcode.com, signup to the test list
- [ ] Stage 15: first branch Preview green; Workers Preview URL reachable by Lighthouse CI; the 13 `check-worker.sh` cases pass; live Preview checks
- [ ] Stage 16: ADR-0001 content gate met (≥ 30 published Explainers across all six Verticals, ≥ 20 Briefs, no `TODO(author)`; today `content/` has none), then the first `publish cf` to the Worker
- [ ] Stage 17: Production DNS flip (apex and `www` to the Worker, `www` → apex Redirect Rule, Worker secret to the real newsletter key)
- [ ] Stage 18: `INDEXNOW_ENABLED=true`, Production caches cleared, `publish cf`, live checklist (IndexNow 200/202, custom domains still attached after the deploy, secret still set, `www` 301, Plausible hostnames, Search Console and Bing sitemaps)
- [ ] Stage 19: first `make prune-images PRUNE_FLAGS=--delete` by hand (DeleteObjects checksum accepted); R2 metrics: HEAD miss + PUT = one Class A (research Q2), DeleteObjects billing (research Q3)

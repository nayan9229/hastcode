# 08: Workflows and the publish command

**What to build:** Publishing becomes a tag. `scripts/publish.sh cf|vl` refuses on a dirty tree or off `main`, creates a lightweight tag `PUB<CF|VL>-<UTC yyyymmdd-HHMMSS>` on HEAD, pushes it, and after `vl` prints the QA URL and the reminder to tag `cf` with the same commit. Three workflows replace v1's two: **tests** on push to `main`; **preview** on push to any other branch (preview-mode build with `images.store: r2`, deploy to a Cloudflare preview URL, Lighthouse CI with v1's budgets); **publish** on tags matching `PUBCF-*` or `PUBVL-*`, which reads the Target from the prefix, fails unless the tagged commit is an ancestor of `main`, builds once in production mode, deploys with that Target's adapter (Cloudflare via the Wrangler action; Vercel via `vercel deploy --prebuilt --prod`), and for `PUBCF` only pings IndexNow when the changed-URL artefact exists and the gating repo variable is on. One concurrency group per Target, no cancel-in-progress. Image and diagram caches and the IndexNow manifest persist per Target. The Vercel preview workflow is deleted. All workflows pass actionlint; none can be run here.

Spec: `../spec.md` ("Publishing"). Research: `docs/research/cloudflare-r2-and-deploy.md` sections 5 and 7 (Wrangler action, token scopes, tag triggers, ancestor gate).

**Blocked by:** 03 (Image Store: R2 mode), 05 (Cloudflare manifest), 06 (Signup Worker), 07 (QA Target hardening)

**Status:** ready-for-agent

- [ ] `publish.sh` tested by hand: refuses dirty tree, refuses off-`main`, creates and pushes a correctly named tag (push step may be dry-run locally)
- [ ] `tests.yml` triggers only on `main` pushes; `preview.yml` only on non-`main` branch pushes; `publish.yml` only on the two tag patterns; secrets referenced by name only
- [ ] Publish workflow: ancestor-of-`main` gate with a clear failure; Target chosen from the tag; IndexNow step guarded by Target and variable; concurrency groups per Target
- [ ] Cloudflare deploy step uses the adapter's file set; Vercel deploy step uses its file set with the `robots.txt` overlay
- [ ] actionlint clean; README "Deploying" rewritten around tags, with the list of secrets and variables

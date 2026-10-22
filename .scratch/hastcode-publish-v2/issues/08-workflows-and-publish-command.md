# 08: Workflows and the publish command

**What to build:** Publishing becomes a tag. `scripts/publish.sh cf|vl` refuses on a dirty tree or off `main`, creates a lightweight tag `PUB<CF|VL>-<UTC yyyymmdd-HHMMSS>` on HEAD, pushes it, and after `vl` prints the QA URL and the reminder to tag `cf` with the same commit. Three workflows replace v1's two: **tests** on push to `main`; **preview** on push to any other branch (preview-mode build with `images.store: r2`, deploy to a Cloudflare preview URL, Lighthouse CI with v1's budgets); **publish** on tags matching `PUBCF-*` or `PUBVL-*`, which reads the Target from the prefix, fails unless the tagged commit is an ancestor of `main`, builds once in production mode, deploys with that Target's adapter (Cloudflare via the Wrangler action; Vercel via `vercel deploy --prebuilt --prod`), and for `PUBCF` only pings IndexNow when the changed-URL artefact exists and the gating repo variable is on. One concurrency group per Target, no cancel-in-progress. Image and diagram caches and the IndexNow manifest persist per Target. The Vercel preview workflow is deleted. All workflows pass actionlint; none can be run here.

Spec: `../spec.md` ("Publishing"). Research: `docs/research/cloudflare-r2-and-deploy.md` sections 5 and 7 (Wrangler action, token scopes, tag triggers, ancestor gate).

**Blocked by:** 03 (Image Store: R2 mode), 05 (Cloudflare manifest), 06 (Signup Worker), 07 (QA Target hardening)

**Status:** resolved

- [x] `publish.sh` tested by hand: refuses dirty tree, refuses off-`main`, creates and pushes a correctly named tag (push step may be dry-run locally)
- [x] `tests.yml` triggers only on `main` pushes; `preview.yml` only on non-`main` branch pushes; `publish.yml` only on the two tag patterns; secrets referenced by name only
- [x] Publish workflow: ancestor-of-`main` gate with a clear failure; Target chosen from the tag; IndexNow step guarded by Target and variable; concurrency groups per Target
- [x] Cloudflare deploy step uses the adapter's file set; Vercel deploy step uses its file set with the `robots.txt` overlay
- [x] actionlint clean; README "Deploying" rewritten around tags, with the list of secrets and variables

## Resolution

- `scripts/publish.sh cf|vl [--dry-run]` and `scripts/publish-gate.sh TAG [MAIN_REF]` (name check, peel, `merge-base --is-ancestor`, prints `tag`/`target`/`commit` for `$GITHUB_OUTPUT`). Hand-tested in a throwaway clone with a local bare `origin`: dirty (modified and untracked), off-`main`, detached, unpushed `HEAD`, bad args all refused; `--dry-run` leaves no tag; real `vl` and `cf` pushed lightweight `PUBVL-…`/`PUBCF-…` tags; the gate accepted them (plus an annotated tag and an older ancestor) and refused a branch-only commit, a malformed name, a lowercase name, a missing tag and a missing main ref.
- `tests.yml` (push `main`), `preview.yml` (push to non-`main` branches; `wrangler versions upload`, Version Preview URL from the action's `deployment-url` output with a fallback to Wrangler's "Version Preview URL:" line; Lighthouse CI unchanged), `publish.yml` (tags `PUBCF-*`/`PUBVL-*`); `production.yml` deleted. actionlint 1.7.12 clean. `wrangler versions upload … --message … --dry-run` and `wrangler deploy … --message … --dry-run` pass against a real build after `make adapt TARGET=cloudflare`.
- **Deliberate deviation:** `publish.yml` also has a `workflow_dispatch` trigger (input `tag`). GitHub scopes caches by ref ("Workflow runs also cannot restore caches created for different tag names"), so a tag-triggered run could never restore the previous Publish's cache (images, diagrams, the IndexNow manifest). The tag run is the gate and hands off with `gh workflow run publish.yml --ref main -f tag=…`; the dispatched run (on `main`, so caches persist in `main`'s scope, keys per Target) re-runs the gate and deploys. Still three workflows; the dispatch doubles as "re-deploy this tag" by hand.
- `make test` runs in each publish job before the build (a commit inside a multi-commit push to `main` is never tested on its own).

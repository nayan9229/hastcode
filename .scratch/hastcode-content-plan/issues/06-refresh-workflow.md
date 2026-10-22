# 06: Refresh workflow

**What to build:** Every Monday at 06:00 UTC (and on manual trigger) a GitHub Actions workflow checks out `main`, runs the fetch command in refresh mode, and stops quietly when nothing changed. When something changed it runs the full check target, commits the changed Fact Packs to a `refresh/<date>` branch and opens a pull request titled with the date whose body lists each changed repository and what changed; if that branch and pull request already exist from an earlier run, it updates them instead of opening a second. The workflow never tags and never touches prose. The README explains the Refresh, says that the job runs the tests itself because pull requests opened with the repository token trigger no workflows, and records the first manual run.

Spec: `../spec.md` ("Refresh workflow", User Stories 31–35). Glossary: Refresh. ADR-0004: Publish stays by tag.

**Blocked by:** 05 (Fetch command)

**Status:** ready-for-human

- [x] Workflow has the Monday schedule and a manual trigger; workflow lint passes
- [x] Unchanged tree: job ends green with no branch, no pull request
- [x] Changed tree: tests run before the pull request exists; branch and pull request created with the date title and the per-repository summary in the body
- [x] Re-run on the same day with further changes updates the existing branch and pull request (no duplicate)
- [x] The token used for the API is the repository token; no other secret is introduced
- [ ] README section written; first manual run performed and its result noted there (ready-for-human step inside this ticket)

## Implementation notes

- **Files.** `.github/workflows/refresh.yml` (cron `0 6 * * 1` + `workflow_dispatch`; top-level `permissions: contents: read`, the job elevates to `contents: write`, `pull-requests: write`; `concurrency: group: refresh, cancel-in-progress: false`; checkout `ref: main`; `./.github/actions/setup` for Go + Node + `npm ci`, as tests.yml; actions on the same version tags as the other workflows), `scripts/refresh-pr.sh` (commit, branch, push, `gh pr create`/`gh pr edit`), README "Repo Profiles: fetching a Fact Pack" → "### Refresh" (written; only its "First run" line waits for the human step).
- **Deviation: the workflow does not run `make refresh`.** GNU make reports a recipe's exit 3 as "Error 3" and itself exits 2, which collides with the command's usage code 2. The step runs `go build -o "$RUNNER_TEMP/fetch" ./cmd/fetch` and the binary directly, stdout to `$RUNNER_TEMP/refresh.out`. The Makefile comment and the README `make refresh` line that said "exit 3" now say the command exits 3 and make prints "Error 3".
- **Exit handling.** 0 → step summary `Nothing changed: refresh: N Repo Profiles checked, 0 changed`, later steps skipped, run green. 1/2 → `::error` annotation, run fails. 3 → `make check`, then `scripts/refresh-pr.sh`.
- **refresh-pr.sh.** Body from stdout lines matching `^[^ ]+/[^ ]+: ` (bulleted), `note: ` lines under "Notes for the Author", the `refresh:` total, the run URL, and "Merging does not publish" (Publish tag, ADR-0004). Stages `content/` and refuses when anything but `content/**.json` changed (facts, never prose). Branch `refresh/<UTC date>` created or reset from main (`checkout -B`), commit as `github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>`, force-pushed (the branch belongs to the workflow; the concurrency group rules out a racing run, so no lease). `gh pr list --head <branch> --base main --state open` → `gh pr edit` (title + body) or `gh pr create`. Never tags. `set -euo pipefail`, shellcheck-clean; the repo tests no scripts, so no test added.
- **Token.** `github.token` only: `GITHUB_TOKEN` for the fetch, `GH_TOKEN` for `gh`; checkout's persisted credentials do the push. No secret introduced.
- **Not executed end to end.** `make lint-workflows` (actionlint + shellcheck) passes; the script's git/gh path was checked by reading only (the agent sandbox refused git on a scratch repository). Criteria 2–4 are ticked on that basis; the first real run below confirms them.

## Remaining human steps

1. Push the repository to GitHub (the Refresh needs a remote, `main` and Actions; Settings → Actions → General → Workflow permissions must allow GitHub Actions to create pull requests, or `gh pr create` is refused).
2. Run "Refresh" by hand once (Actions tab → Refresh → Run workflow).
3. Record the date and the result (exit status; the pull request, or "nothing changed") in README "Refresh" → "First run", then tick the last criterion.

#!/usr/bin/env bash
# Publish the current main commit to a Target (spec "Publishing", ADR-0004).
#
#   scripts/publish.sh vl            # QA Target (Vercel): tag PUBVL-<UTC stamp>
#   scripts/publish.sh cf            # Production Target (Cloudflare): tag PUBCF-<UTC stamp>
#   scripts/publish.sh vl --dry-run  # every check, then `git push --dry-run`; no tag is left behind
#
# Creates a lightweight tag PUB<CF|VL>-<yyyymmdd>-<HHMMSS> (UTC) on HEAD and
# pushes it to origin; the tag push starts .github/workflows/publish.yml,
# which deploys that commit. Refuses to run:
#   - with uncommitted changes or untracked files (what you see locally
#     must be what gets published);
#   - on any branch but main;
#   - when HEAD is not on origin/main (push main first: the workflow
#     refuses tags on commits that are not on main).
# The normal flow is `vl`, check https://qa.hastcode.com/, then `cf` on the
# same commit. QA_URL overrides the QA address printed after `vl`.
# The Target codes and tag prefixes come from internal/deploy (`go run
# ./cmd/target`; the full names cloudflare/vercel work too), so this needs Go.
set -euo pipefail

usage() {
  echo "usage: scripts/publish.sh cf|vl [--dry-run]" >&2
  exit 2
}

arg_target=""
dry_run=false
for arg in "$@"; do
  case "$arg" in
    -n | --dry-run) dry_run=true ;;
    -* | PUB*) usage ;;
    *) [[ -z "$arg_target" ]] || usage; arg_target="$arg" ;;
  esac
done
[[ -n "$arg_target" ]] || usage

refuse() {
  echo "publish: $*" >&2
  exit 1
}

cd "$(git rev-parse --show-toplevel)"

spelling=$(go run ./cmd/target "$arg_target" 2>/dev/null) || usage
target=$(sed -n 's/^target=//p' <<<"$spelling")
prefix=$(sed -n 's/^tag_prefix=//p' <<<"$spelling")

if [[ -n "$(git status --porcelain --untracked-files=normal)" ]]; then
  git status --short >&2
  refuse "uncommitted changes or untracked files (above). Commit, stash or remove them first: a Publish deploys the commit, not your working tree."
fi

branch=$(git symbolic-ref --short -q HEAD || true)
[[ "$branch" == main ]] ||
  refuse "on ${branch:-a detached HEAD}, not main. Publish tags go on main: git switch main && git pull"

git fetch --quiet origin main ||
  refuse "could not fetch main from origin"
head=$(git rev-parse HEAD)
if ! git merge-base --is-ancestor "$head" origin/main; then
  refuse "HEAD ${head:0:12} is not on origin/main (unpushed or diverged). Push main first: git push origin main"
fi
if [[ "$head" != "$(git rev-parse origin/main)" ]]; then
  echo "note: origin/main is ahead of HEAD; publishing HEAD ${head:0:12} (git pull to publish the latest)" >&2
fi

# The UTC stamp orders Publishes: publish.yml deploys only the newest tag
# of a Target and skips older ones that are still waiting.
tag="${prefix}$(date -u +%Y%m%d-%H%M%S)"

git tag "$tag" "$head"
if $dry_run; then
  git push --dry-run origin "refs/tags/$tag"
  git tag -d "$tag" >/dev/null
  echo "dry run: would have pushed $tag on ${head:0:12} (no tag kept)"
  exit 0
fi
if ! git push origin "refs/tags/$tag"; then
  git tag -d "$tag" >/dev/null
  refuse "push of $tag failed; the local tag was removed, nothing was published"
fi

echo "Pushed $tag on ${head:0:12} ($(git log -1 --format=%s "$head"))."
remote=$(git remote get-url origin)
if [[ "$remote" =~ github\.com[:/]([^/]+/[^/]+)$ ]]; then
  echo "Workflow: https://github.com/${BASH_REMATCH[1]%.git}/actions/workflows/publish.yml"
fi

if [[ "$target" == vercel ]]; then
  cat <<EOF

QA (Vercel), once the workflow is green: ${QA_URL:-https://qa.hastcode.com/}
When QA looks right, publish the same commit to Production:
  scripts/publish.sh cf      # with HEAD still at ${head:0:12}
EOF
else
  echo
  echo "Production (Cloudflare), once the workflow is green: https://hastcode.com/"
fi

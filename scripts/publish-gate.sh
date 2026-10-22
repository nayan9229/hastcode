#!/usr/bin/env bash
# The publish gate (spec "Publishing", ADR-0004): decides whether a Publish
# tag may deploy, and where to.
#
#   scripts/publish-gate.sh TAG [MAIN_REF]     # MAIN_REF default origin/main
#
# Refuses (exit 1, reason on stderr) unless:
#   - TAG is named PUBCF-<yyyymmdd>-<HHMMSS> or PUBVL-<yyyymmdd>-<HHMMSS>
#     (what scripts/publish.sh creates; the strict pattern also makes the
#     name safe to pass to later steps);
#   - refs/tags/TAG exists here (fetch tags first);
#   - the commit it points at (annotated tags are peeled) is MAIN_REF or an
#     ancestor of it, so an unmerged branch can never be published.
# On success prints, for $GITHUB_OUTPUT:
#   tag=<TAG>
#   target=cloudflare|vercel     (PUBCF -> cloudflare, PUBVL -> vercel)
#   commit=<full sha>
#
# Used by .github/workflows/publish.yml; runs anywhere with bash + git.
set -euo pipefail

fail() {
  if [[ "${GITHUB_ACTIONS:-}" == true ]]; then
    echo "::error title=Publish refused::$*" >&2
  else
    echo "publish refused: $*" >&2
  fi
  exit 1
}

tag="${1:-}"
main_ref="${2:-origin/main}"
[[ -n "$tag" ]] || fail "usage: publish-gate.sh TAG [MAIN_REF]"

if [[ ! "$tag" =~ ^PUB(CF|VL)-[0-9]{8}-[0-9]{6}$ ]]; then
  fail "tag '$tag' is not PUBCF-<yyyymmdd>-<HHMMSS> or PUBVL-<yyyymmdd>-<HHMMSS>; create Publish tags with scripts/publish.sh"
fi
case "${BASH_REMATCH[1]}" in
  CF) target=cloudflare ;;
  VL) target=vercel ;;
esac

commit=$(git rev-parse --verify --quiet "refs/tags/$tag^{commit}") ||
  fail "tag $tag does not exist in this clone (fetch tags, e.g. checkout with fetch-depth: 0)"
main=$(git rev-parse --verify --quiet "$main_ref^{commit}") ||
  fail "$main_ref does not exist in this clone (fetch main, e.g. checkout with fetch-depth: 0)"

if ! git merge-base --is-ancestor "$commit" "$main"; then
  fail "tag $tag points at ${commit:0:12}, which is not on main ($main_ref is ${main:0:12}). Only commits merged to main can be published: merge first, then run scripts/publish.sh on main. Delete this tag with: git push origin :refs/tags/$tag"
fi

echo "tag=$tag"
echo "target=$target"
echo "commit=$commit"

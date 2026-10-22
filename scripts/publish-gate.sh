#!/usr/bin/env bash
# The publish gate (spec "Publishing", ADR-0004): decides whether a Publish
# tag may deploy, and where to.
#
#   scripts/publish-gate.sh TAG [MAIN_REF]     # MAIN_REF default origin/main
#
# Refuses (exit 1, reason on stderr) unless:
#   - TAG is named PUB<code>-<yyyymmdd>-<HHMMSS> for a Target's code
#     (PUBCF-… or PUBVL-…; what scripts/publish.sh creates; the strict
#     pattern also makes the name safe to pass to later steps);
#   - refs/tags/TAG exists here (fetch tags first);
#   - the commit it points at (annotated tags are peeled) is MAIN_REF or an
#     ancestor of it, so an unmerged branch can never be published.
# On success prints, for $GITHUB_OUTPUT:
#   tag=<TAG>
#   target=cloudflare|vercel     (from the prefix, via cmd/target)
#   commit=<full sha>
#   superseded_by=<newer tag>|   (empty unless a newer tag exists, below)
#
# Newest tag wins: tag names embed their UTC creation time, so a later name
# is a later Publish. superseded_by names the newest tag for the same Target
# that is later than TAG and passes the same checks (its commit on
# MAIN_REF); a deploy for TAG then must not run, or an older commit would
# replace a newer one. Run it on a fresh fetch of every tag.
#
# Target spellings come from internal/deploy (`go run ./cmd/target`), so
# this needs Go. Used by .github/workflows/publish.yml; runs anywhere with
# bash, git and go.
set -euo pipefail
export LC_ALL=C # byte-wise name comparisons below

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

stamp='[0-9]{8}-[0-9]{6}'
if [[ ! "$tag" =~ ^PUB[A-Z]{2}-$stamp$ ]] ||
  ! spelling=$(cd "$(dirname "$0")/.." && go run ./cmd/target "$tag" 2>/dev/null); then
  fail "tag '$tag' is not PUBCF-<yyyymmdd>-<HHMMSS> or PUBVL-<yyyymmdd>-<HHMMSS>; create Publish tags with scripts/publish.sh"
fi
target=$(sed -n 's/^target=//p' <<<"$spelling")
prefix=$(sed -n 's/^tag_prefix=//p' <<<"$spelling")

commit=$(git rev-parse --verify --quiet "refs/tags/$tag^{commit}") ||
  fail "tag $tag does not exist in this clone (fetch tags, e.g. checkout with fetch-depth: 0)"
main=$(git rev-parse --verify --quiet "$main_ref^{commit}") ||
  fail "$main_ref does not exist in this clone (fetch main, e.g. checkout with fetch-depth: 0)"

if ! git merge-base --is-ancestor "$commit" "$main"; then
  fail "tag $tag points at ${commit:0:12}, which is not on main ($main_ref is ${main:0:12}). Only commits merged to main can be published: merge first, then run scripts/publish.sh on main. Delete this tag with: git push origin :refs/tags/$tag"
fi

# Ascending by name = by time; the last match is the newest.
superseded_by=""
while IFS= read -r other; do
  [[ "$other" =~ ^$prefix$stamp$ && "$other" > "$tag" ]] || continue
  other_commit=$(git rev-parse --verify --quiet "refs/tags/$other^{commit}") || continue
  if git merge-base --is-ancestor "$other_commit" "$main"; then
    superseded_by="$other"
  fi
done < <(git tag --list "$prefix*" | LC_ALL=C sort)

echo "tag=$tag"
echo "target=$target"
echo "commit=$commit"
echo "superseded_by=$superseded_by"

#!/usr/bin/env bash
# make release (ADR-0012): the Author's release step, after verifying every
# draft in the review queue.
#
#   make release                       # release every draft
#   make release FILES="content/briefs/a.md content/explainers/v/b.md"
#   make release DRY=1                 # print the plan, change nothing
#
# Per Piece: removes draft: true, sets verified to today (in the Piece's own
# offset), sets reviewed to now on an Explainer, keeps published as the
# draft carries it (a backfilled Gap Day stays that day). Refuses before
# touching anything when a Publication Day would exceed two Pieces. Then
# runs the production build (the gates) and commits. Nothing is pushed or
# tagged: next is git push origin main, then scripts/publish.sh cf.
# An agent never runs this (CLAUDE.md): verified: is the Author's record.
set -euo pipefail
cd "$(dirname "$0")/.."
FILES="${FILES:-}" DRY="${DRY:-}"

if [ -n "$DRY" ]; then
  go run ./cmd/release -content content -files "$FILES" -dry
  exit 0
fi
if ! { git diff --quiet && git diff --cached --quiet; }; then
  echo "make release commits its release: commit or stash your tracked changes first" >&2; exit 2
fi
out=$(go run ./cmd/release -content content -files "$FILES")
printf '%s\n' "$out"
paths=$(printf '%s\n' "$out" | awk -F'\t' '$1=="released"{print "content/" $2}')
titles=$(printf '%s\n' "$out" | awk -F'\t' '$1=="released"{t=(t==""?"":t", ")$5} END{print t}')
[ -n "$paths" ] || { echo "nothing released" >&2; exit 1; }

echo "Production build (the gates):"
if ! make -s build; then
  echo "build failed: the drafts above are released in the working tree; fix the error and run make build, or git checkout -- content to undo" >&2
  exit 1
fi
day=$(date -u +%Y-%m-%d)
printf '%s\n' "$paths" | xargs git add --
git commit -q -m "Release $day: $titles, verified by the Author" -m "Released by make release (ADR-0012); Publication Days kept as drafted."
echo "Committed on $(git branch --show-current). Next: git push origin main && scripts/publish.sh cf"

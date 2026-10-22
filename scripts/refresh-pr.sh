#!/usr/bin/env bash
# The Refresh pull request (spec "Refresh workflow"; refresh.yml runs it
# after the fetch command exited 3 and `make check` passed).
#
#   scripts/refresh-pr.sh REFRESH_STDOUT
#
# REFRESH_STDOUT is the fetch command's stdout: one `owner/name: …` line
# per changed repository, `note: …` lines, then the `refresh: …` total.
# Commits the changed Fact Packs, on top of the checked-out main, to
# refresh/<UTC date> as the github-actions bot (the branch is created, or
# reset onto main and force-updated: it belongs to the workflow), pushes it
# and opens the pull request, or updates the open one for that branch
# (title and body), so a re-run on the same day never makes a second.
# Refuses when anything but a .json file under content/ changed: a
# Refresh touches facts, never prose. Never tags (ADR-0004).
# Needs GH_TOKEN (the repository token); RUN_URL, when set, is linked.
set -euo pipefail

out=${1:?usage: scripts/refresh-pr.sh REFRESH_STDOUT}
date=$(date -u +%Y-%m-%d)
branch="refresh/$date"
title="Refresh: Fact Packs $date"

changes=$(grep -E '^[^ ]+/[^ ]+: ' "$out" || true)
notes=$(grep -E '^note: ' "$out" || true)
total=$(grep -E '^refresh: ' "$out" || true)
if [ -z "$changes" ]; then
  echo "refresh-pr: no changed-repository lines in $out" >&2
  exit 1
fi

git add -A -- content
other=$(git diff --cached --name-only | grep -v '^content/.*\.json$' || true)
if [ -n "$other" ]; then
  printf 'refresh-pr: refusing, files other than Fact Packs changed:\n%s\n' "$other" >&2
  exit 1
fi
if git diff --cached --quiet; then
  echo "refresh-pr: the fetch command reported changes but no Fact Pack differs" >&2
  exit 1
fi

git checkout -q -B "$branch"
git -c user.name='github-actions[bot]' \
  -c user.email='41898282+github-actions[bot]@users.noreply.github.com' \
  commit -q -m "$title" -m "$changes"
git push -q --force origin "HEAD:refs/heads/$branch"

bullets() { while IFS= read -r line; do printf -- "- %s\n" "$line"; done; }

body=$(mktemp)
{
  echo "Fact Packs re-fetched on $date. Changed repositories:"
  echo
  bullets <<<"$changes"
  if [ -n "$notes" ]; then
    echo
    echo "Notes for the Author:"
    echo
    bullets <<<"$notes"
  fi
  echo
  echo "${total:+$total. }\`make check\` passed in ${RUN_URL:-the Refresh run} before this pull request was opened (pull requests opened with the repository token trigger no workflows)."
  echo
  echo "**Merging does not publish.** Readers see these facts after the next Publish tag (\`scripts/publish.sh vl\`, then \`cf\`; ADR-0004)."
} >"$body"

pr=$(gh pr list --head "$branch" --base main --state open --json number --jq '.[0].number // empty')
if [ -n "$pr" ]; then
  gh pr edit "$pr" --title "$title" --body-file "$body"
  action="Updated pull request #$pr"
else
  gh pr create --base main --head "$branch" --title "$title" --body-file "$body"
  action="Opened a pull request"
fi
echo "$action ($branch): $(wc -l <<<"$changes" | tr -d ' ') changed Fact Pack(s)." >>"${GITHUB_STEP_SUMMARY:-/dev/null}"

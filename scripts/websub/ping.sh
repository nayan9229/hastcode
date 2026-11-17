#!/usr/bin/env bash
# scripts/websub/ping.sh <websub.json>: tell the WebSub hub which feeds
# changed. Run by the publish workflow's "Ping WebSub hub" step, in the
# Production job only, after the deploy is live, when the repository
# variable WEBSUB_ENABLED is "true" (README "Feeds, sitemaps, robots and
# IndexNow"). Never run it by hand: it sends a notice to a third party.
#
# <websub.json> is the build's artefact {hub, topics}: the hub every feed
# names in atom:link rel="hub" and the feeds (their rel="self" URLs) whose
# bytes changed since the last Production build. No file means no feed
# changed: nothing is sent. Otherwise one form POST, as the hub's own page
# (https://pubsubhubbub.appspot.com/, read 2026-11-17) asks:
# hub.mode=publish and one hub.url per changed feed. The hub answers 204.
#
# A failed ping only warns (a GitHub ::warning and a line in the run
# summary, $GITHUB_STEP_SUMMARY) and exits 0: the Publish is already live,
# and the job must go on to save the build cache.
set -uo pipefail

notice=${1:?usage: scripts/websub/ping.sh <websub.json>}
summary=${GITHUB_STEP_SUMMARY:-/dev/null}

if [ ! -f "$notice" ]; then
  echo "no artefact: no feed changed, nothing to ping"
  exit 0
fi

hub=$(jq -r .hub "$notice")
n=$(jq '.topics | length' "$notice")
args=(--data-urlencode hub.mode=publish)
while IFS= read -r topic; do
  echo "changed feed: $topic"
  args+=(--data-urlencode "hub.url=$topic")
done < <(jq -r '.topics[]' "$notice")

code=$(curl -sS -o /dev/null -w '%{http_code}' --max-time 30 -X POST "$hub" "${args[@]}") || code="no answer (curl ${code:-failed})"
echo "$hub answered $code for $n feed(s)"

case "$code" in
  2??)
    echo "WebSub: $hub answered $code for $n feed(s)." >>"$summary"
    ;;
  *)
    echo "::warning title=WebSub ping failed::$hub answered $code for $n feed(s). The Publish itself is fine."
    echo "WebSub ping failed: $hub answered $code for $n feed(s). The Publish itself is fine; the hub hears about these feeds at the next Publish that changes them." >>"$summary"
    ;;
esac
exit 0

#!/usr/bin/env bash
# Prints the URLs Lighthouse CI audits on a deployment, one per line:
# home, /search/, the Briefs index, one Vertical hub, one Explainer and one Brief (the last
# three only once the build wrote them), taken from the build's own
# sitemap so the list follows the content.
#
#   scripts/lighthouse-urls.sh https://abc12345-hastcode.example.workers.dev [.vercel/output]
#
# Used by .github/workflows/preview.yml (Cloudflare Version Preview URL); runs anywhere with bash + grep/sed.
set -euo pipefail

base="${1:?usage: lighthouse-urls.sh BASE_URL [BUILD_OUTPUT_DIR]}"
base="${base%/}"
out="${2:-.vercel/output}"
sitemap="$out/static/sitemap.xml"
[[ -f "$sitemap" ]] || { echo "no $sitemap: build first" >&2; exit 1; }

# Root-relative paths of every <loc>.
paths=$(sed -n 's#.*<loc>https\{0,1\}://[^/]*\(/[^<]*\)</loc>.*#\1#p' "$sitemap")

# Explainers live at /<vertical>/<slug>/; other two-segment paths are
# Briefs, Author pages, kind pages or pagination.
explainer=$(grep -E '^/[a-z0-9-]+/[a-z0-9-]+/$' <<<"$paths" |
  grep -vE '^/(briefs|authors|explainers)/|/page/$' | head -n1 || true)
brief=$(grep -E '^/briefs/[a-z0-9-]+/$' <<<"$paths" | grep -v '^/briefs/page/' | head -n1 || true)

echo "$base/"
echo "$base/search/"
echo "$base/briefs/"
if [[ -n "$explainer" ]]; then
  echo "$base/${explainer#/}" | sed 's#^\(.*://[^/]*/[^/]*/\).*#\1#'   # its Vertical hub
  echo "$base$explainer"
fi
if [[ -n "$brief" ]]; then
  echo "$base$brief"
fi

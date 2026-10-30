#!/usr/bin/env bash
# Preview-build the site into <out dir> and report on one Piece (make daily,
# ADR-0009): every build error or warning naming the file, then the page's
# HTML size without the inline stylesheet, as the html-size gate counts it.
#   scripts/daily/check-page.sh content/explainers/<vertical>/<slug>.md <out dir>
set -uo pipefail
file="$1" out="$2"
slug="$(basename "$file" .md)"
case "$file" in
  content/briefs/*) url="briefs/$slug" ;;
  content/explainers/*) url="$(basename "$(dirname "$file")")/$slug" ;;
  *) echo "not a Piece: $file" >&2; exit 2 ;;
esac
rm -rf "$out"
go run ./cmd/builder -content content -config site.yaml -out "$out" -mode preview -skip-search-index 2>&1 | grep -F "$slug.md" || echo "no build errors or warnings name $slug.md"
page="$out/static/$url/index.html"
[ -f "$page" ] || { echo "page not written: $page"; exit 1; }
python3 -c "import re,sys;s=open(sys.argv[1]).read();print('page bytes', len(re.sub(r'<style>.*?</style>','',s,count=1,flags=re.S).encode()), '(at most 49500)')" "$page"

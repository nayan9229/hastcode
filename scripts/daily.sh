#!/usr/bin/env bash
# make daily (ADR-0009): pick and draft today's Pieces; the Author verifies.
#
#   make daily                    # Intake, choose, pick, draft, fact-check, commit
#   make daily DRY=1              # Intake and choose only: print the picks, write nothing
#   make daily EXPLAINERS=2 BRIEFS=2 PARALLEL=4 BACKLOG=4
#
# 1. make intake refreshes the queue.
# 2. A headless Claude Code session (claude -p) chooses: every fresh Lead that
#    can carry a Brief with a live Explainer teaching its lesson (BRIEFS, default
#    2, or all), and up to EXPLAINERS Leads that can carry a Case Study, Concept
#    or Comparison (default 2, or 0 while BACKLOG (4) Explainer drafts wait).
#    It explains every rejection.
# 3. Each pick goes through make lead, one at a time (the intake state is not
#    safe for parallel writers), and is renamed to the chosen slug.
# 4. PARALLEL (default 4) pipelines: one session drafts each Piece by its
#    kind's writing guide, a second fact-checks it against its Sources. A
#    draft the Source cannot carry is removed (its Lead stays picked).
# 5. The drafts and the Author's check list
#    (.scratch/daily/<date>/review-queue.md) are committed on the current
#    branch. Nothing is pushed or tagged; every draft keeps draft: true and
#    no verified: (CLAUDE.md, ADR-0007).
# Run logs and each session's JSON are kept under .intake/daily/<run>/.
set -euo pipefail
cd "$(dirname "$0")/.."

CLAUDE="${CLAUDE:-claude}"
WORKER_TOOLS="Read Edit Write Glob Grep WebFetch WebSearch Bash(scripts/daily/check-page.sh:*)"
SELECT_TOOLS="Read Glob Grep WebFetch WebSearch"
SIZE_LIMIT='SIZE LIMIT: the check prints "page bytes"; keep it at most 49,500. Production pages carry about 1.4 KB more (ads, analytics) against a hard limit of 51,200. Two Mermaid diagrams plus 2,800 words usually overflows: prefer one diagram and about 2,400 words.'

render() { # render <template> KEY=value...: substitute {{KEY}}
  python3 - "$@" <<'PY'
import sys, pathlib
text = pathlib.Path(sys.argv[1]).read_text()
for kv in sys.argv[2:]:
    k, v = kv.split("=", 1)
    text = text.replace("{{" + k + "}}", v)
print(text)
PY
}

claude_json() { # claude_json <out.json> <schema> <tools> <prompt>
  "$CLAUDE" -p "$4" --output-format json --json-schema "$2" \
    --permission-mode dontAsk --allowedTools "$3" >"$1" 2>"$1.err" || true
  python3 -c 'import json,sys
try: d=json.load(open(sys.argv[1]))
except Exception: sys.exit(1)
sys.exit(0 if isinstance(d.get("structured_output"),dict) else 1)' "$1"
}

DRAFT_SCHEMA='{"type":"object","properties":{"status":{"type":"string","enum":["drafted","not-viable"]},"title":{"type":"string"},"reason":{"type":"string"},"words":{"type":"number"},"page_bytes":{"type":"number"},"author_checks":{"type":"array","items":{"type":"string"}},"build":{"type":"string"}},"required":["status","title","reason","build"]}'
CHECK_SCHEMA='{"type":"object","properties":{"claims_checked":{"type":"number"},"fixed":{"type":"array","items":{"type":"string"}},"removed":{"type":"array","items":{"type":"string"}},"unconfirmed":{"type":"array","items":{"type":"string"}},"page_bytes":{"type":"number"},"build":{"type":"string"}},"required":["claims_checked","fixed","removed","unconfirmed","build"]}'
SELECT_SCHEMA='{"type":"object","properties":{"picks":{"type":"array","items":{"type":"object","properties":{"id":{"type":"string"},"kind":{"type":"string","enum":["brief","case-study","concept","comparison"]},"vertical":{"type":"string"},"slug":{"type":"string"},"explainers":{"type":"array","items":{"type":"string"}},"angle":{"type":"string"}},"required":["id","kind","slug","angle"]}},"skipped":{"type":"array","items":{"type":"object","properties":{"id":{"type":"string"},"reason":{"type":"string"}},"required":["id","reason"]}}},"required":["picks","skipped"]}'

# __piece <run dir> <n>: draft then fact-check pick n (called in parallel).
if [ "${1:-}" = "__piece" ]; then
  run="$2" n="$3"
  eval "$(python3 -c 'import json,sys,shlex
p=json.load(open(sys.argv[1]))[int(sys.argv[2])]
for k in ("path","kind","slug","angle"): print(f"{k}={shlex.quote(p[k])}")
print("explainers="+shlex.quote(",".join(p.get("explainers") or [])))' "$run/picked.json" "$n")"
  case "$kind" in
    brief) examples="content/briefs/github-dependabot-fan-out.md and content/briefs/boat-bedrock-fitness-coach.md"
           limits="BRIEF LIMITS: 250 words at most; sections exactly What happened and Why it matters; link every Explainer in explainers: ($explainers) in Why it matters, and no other Explainer." ;;
    case-study) examples="content/explainers/infrastructure-cloud/netflix-single-tenant-titus.md and content/explainers/ai-ml-systems/meta-thin-orchestrator.md"; limits="$SIZE_LIMIT" ;;
    concept) examples="content/explainers/distributed-systems/metastable-failures.md and content/explainers/databases-storage/apache-iceberg-table-format.md"; limits="$SIZE_LIMIT" ;;
    comparison) examples="content/explainers/apis-networking/http-connection-reuse.md and content/explainers/databases-storage/cascade-vs-application-vs-soft-deletes.md"; limits="$SIZE_LIMIT" ;;
  esac
  out="$run/build-$n"
  if ! claude_json "$run/$slug.draft.json" "$DRAFT_SCHEMA" "$WORKER_TOOLS" \
      "$(render scripts/daily/draft.md KIND="$kind" PATH="$path" ANGLE="$angle" EXAMPLES="$examples" LIMITS="$limits" OUT="$out")"; then
    echo "  ✗ $slug: the drafting session failed (log: $run/$slug.draft.json.err); skeleton kept"; exit 0
  fi
  status=$(python3 -c 'import json,sys;print(json.load(open(sys.argv[1]))["structured_output"]["status"])' "$run/$slug.draft.json")
  if [ "$status" != drafted ]; then
    rm -f "$path"; echo "  – $slug: not viable, skeleton removed"; exit 0
  fi
  notes=$(python3 -c 'import json,sys;print(json.dumps(json.load(open(sys.argv[1]))["structured_output"].get("author_checks",[])))' "$run/$slug.draft.json")
  claude_json "$run/$slug.check.json" "$CHECK_SCHEMA" "$WORKER_TOOLS" \
    "$(render scripts/daily/check.md KIND="$kind" PATH="$path" LIMITS="$limits" OUT="$out-check" NOTES="$notes")" \
    || echo "  ! $slug: the fact-check session failed; draft kept unchecked"
  echo "  ✓ $slug"
  exit 0
fi

# The Publish cap is two Pieces a day (ADR-0010, enforced by the build), so
# drafting keeps pace with it: up to 2 Briefs (they expire after 72 hours)
# and, while fewer than BACKLOG (4) Explainer drafts wait, up to 2 Explainers.
backlog=$( { grep -l '^draft: true' content/explainers/*/*.md 2>/dev/null || true; } | wc -l | tr -d ' ')
BRIEFS="${BRIEFS:-2}"; PARALLEL="${PARALLEL:-4}"; DRY="${DRY:-}"
if [ -z "${EXPLAINERS:-}" ]; then
  if [ "$backlog" -ge "${BACKLOG:-4}" ]; then EXPLAINERS=0; else EXPLAINERS=2; fi
fi
echo "Explainer drafts waiting: $backlog (Explainers today: $EXPLAINERS, Briefs: up to $BRIEFS)"
command -v "$CLAUDE" >/dev/null || { echo "make daily needs Claude Code (claude) on PATH" >&2; exit 2; }
if [ -z "$DRY" ] && ! { git diff --quiet && git diff --cached --quiet; }; then
  echo "make daily commits its drafts: commit or stash your tracked changes first" >&2; exit 2
fi
now="$(date -u +%Y-%m-%dT%H:%M:%SZ)" day="${now%%T*}"
run=".intake/daily/$(date -u +%Y%m%d-%H%M%S)"
mkdir -p "$run"

echo "1/5 Intake"
make -s intake 2>&1 | grep -v '(✓)' | tail -3
python3 scripts/daily/candidates.py .intake/state.json .intake/queue.md content "$now" "$run"

echo "2/5 Choosing (Briefs: $BRIEFS, Explainers: up to $EXPLAINERS)"
[ "$BRIEFS" = all ] && briefs_rule="pick every Lead that qualifies" || briefs_rule="pick at most $BRIEFS"
claude_json "$run/select.json" "$SELECT_SCHEMA" "$SELECT_TOOLS" \
  "$(render scripts/daily/select.md RUN="$run" BRIEFS="$briefs_rule" EXPLAINERS="$EXPLAINERS")" \
  || { echo "the choosing session failed: $run/select.json.err" >&2; exit 1; }
python3 - "$run" <<'PY'
import json, sys
d = json.load(open(sys.argv[1] + "/select.json"))["structured_output"]
for p in d["picks"]:
    print(f"  {p['kind']:11} {p['id']}  {p['slug']}" + (f"  → {', '.join(p.get('explainers') or [])}" if p['kind']=='brief' else f"  ({p.get('vertical','')})"))
print(f"  {len(d['picks'])} picked, {len(d['skipped'])} rejected with reasons (select.json)")
PY
[ -n "$DRY" ] && { echo "DRY=1: nothing picked or written. Choices: $run/select.json"; exit 0; }

echo "3/5 Picking"
python3 - "$run" <<'PY' >"$run/pick.sh"
import json, shlex, sys
live = {x["slug"] for x in json.load(open(sys.argv[1] + "/explainers.json")) if x["type"] == "explainer" and x["live"]}
for p in json.load(open(sys.argv[1] + "/select.json"))["structured_output"]["picks"]:
    if p["kind"] == "brief" and not (p.get("explainers") and set(p["explainers"]) <= live):
        print(f"  ✗ {p['id']}: a Brief must link live Explainers, got {p.get('explainers')}", file=sys.stderr); continue
    args = [f"ID={p['id']}", f"KIND={p['kind']}"]
    if p["kind"] == "brief": args.append("EXPLAINER=" + ",".join(p.get("explainers") or []))
    elif p.get("vertical"): args.append(f"VERTICAL={p['vertical']}")
    print(" ".join(shlex.quote(a) for a in [p["id"], p["slug"]] + args))
PY
: >"$run/picked.jsonl"
while read -r line; do
  eval "set -- $line"; id="$1" slug="$2"; shift 2
  if ! written=$(make -s lead "$@" 2>&1 | sed -n 's/^wrote //p' | tail -1) || [ -z "$written" ]; then
    echo "  ✗ $id: make lead refused (run make lead $* to see why)"; continue
  fi
  target="$(dirname "$written")/$slug.md"
  if [[ "$slug" =~ ^[a-z0-9]+(-[a-z0-9]+)*$ ]] && [ "$written" != "$target" ] && ! ls content/explainers/*/"$slug".md content/briefs/"$slug".md >/dev/null 2>&1; then
    mv "$written" "$target"; written="$target"
  fi
  python3 - "$run" "$id" "$written" <<'PY'
import json, pathlib, sys
run, id_, path = sys.argv[1:4]
p = next(x for x in json.load(open(run + "/select.json"))["structured_output"]["picks"] if x["id"] == id_)
p.update(path=path, slug=pathlib.Path(path).stem)
open(run + "/picked.jsonl", "a").write(json.dumps(p) + "\n")
PY
  echo "  $written"
done <"$run/pick.sh"
python3 -c 'import json,sys;r=sys.argv[1];json.dump([json.loads(l) for l in open(r+"/picked.jsonl")],open(r+"/picked.json","w"),indent=1)' "$run"
count=$(python3 -c 'import json,sys;print(len(json.load(open(sys.argv[1]))))' "$run/picked.json")
[ "$count" -gt 0 ] || { echo "nothing picked today"; exit 0; }

echo "4/5 Drafting and fact-checking $count Pieces, $PARALLEL at a time"
seq 0 $((count - 1)) | xargs -P "$PARALLEL" -I{} "$0" __piece "$run" {}

echo "5/5 Check list and commit"
python3 scripts/daily/report.py "$run" "$day"
files=$(python3 -c 'import json,sys,os;print("\n".join(p["path"] for p in json.load(open(sys.argv[1])) if os.path.exists(p["path"])))' "$run/picked.json")
if [ -n "$files" ]; then
  printf '%s\n' "$files" | xargs git add -- ".scratch/daily/$day"
  git commit -q -m "Daily $day: drafts from make daily (draft: true, unverified)" -m "Check list: .scratch/daily/$day/review-queue.md (ADR-0009)."
  echo "Committed $(printf '%s\n' "$files" | wc -l | tr -d ' ') drafts on $(git branch --show-current). Verify, set verified:, remove draft: true, then Publish."
fi

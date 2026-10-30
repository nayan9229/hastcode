#!/usr/bin/env python3
"""Write the Author's check list for one make daily run (ADR-0009).

usage: report.py <run dir> <yyyy-mm-dd>
Reads picked.json and each <slug>.draft.json / <slug>.check.json (claude -p
JSON output) from the run dir; writes .scratch/daily/<date>/review-queue.md
(appending when several runs share a day) and prints the run summary.
"""
import json, os, pathlib, sys

run, day = pathlib.Path(sys.argv[1]), sys.argv[2]
picked = json.load(open(run / "picked.json"))
select = json.load(open(run / "select.json"))

def load(name):
    try:
        d = json.load(open(run / name))
        return d.get("structured_output") or {}, d.get("total_cost_usd") or 0
    except Exception:
        return {}, 0

cost = select.get("total_cost_usd") or 0
lines, drafted, dropped = [], 0, []
for p in picked:
    d, c1 = load(f"{p['slug']}.draft.json")
    c, c2 = load(f"{p['slug']}.check.json")
    cost += c1 + c2
    if not os.path.exists(p["path"]):
        dropped.append(f"- `{p['id']}` {p['slug']}: {d.get('reason') or 'not drafted'}")
        continue
    if d.get("status") != "drafted":
        dropped.append(f"- `{p['id']}` {p['slug']}: drafting failed, skeleton kept at `{p['path']}`")
        continue
    drafted += 1
    lines += [f"### {d.get('title')}", "",
              f"`{p['path']}` · {p['kind']} · {d.get('words', '?')} words · fact-checker: "
              + (f"{c.get('claims_checked')} claims checked, {len(c.get('fixed', []))} fixed, {len(c.get('removed', []))} removed" if c else "did not run"), ""]
    if p["kind"] == "brief":
        lines += ["**Publish within 72 hours of the Lead's date**, with or after its Explainer: " + ", ".join(p.get("explainers") or []), ""]
    if d.get("author_checks"):
        lines += ["**Drafter asks you to check:**"] + [f"- [ ] {x}" for x in d["author_checks"]] + [""]
    if c.get("unconfirmed"):
        lines += ["**Fact-checker could not confirm:**"] + [f"- [ ] {x}" for x in c["unconfirmed"]] + [""]
    if c.get("fixed") or c.get("removed"):
        lines += ["<details><summary>What the fact-checker changed</summary>", ""] \
            + [f"- Fixed: {x}" for x in c.get("fixed", [])] + [f"- Removed: {x}" for x in c.get("removed", [])] + ["", "</details>", ""]

skipped = select.get("structured_output", {}).get("skipped", [])
out = pathlib.Path(".scratch/daily") / day
out.mkdir(parents=True, exist_ok=True)
q = out / "review-queue.md"
head = [] if q.exists() else [f"# Daily review queue, {day}", "",
        "Drafted by `make daily` (ADR-0009). Every Piece is `draft: true` and unverified. Verify every claim against its Source, then set `verified:` and remove `draft: true` yourself (ADR-0007).", ""]
tail = ["### Not drafted", ""] + (dropped or ["- none"]) + ["", "<details><summary>Leads the editor rejected</summary>", ""] \
    + [f"- `{s['id']}`: {s['reason']}" for s in skipped] + ["", "</details>", ""]
with open(q, "a") as f:
    f.write("\n".join(head + [f"## Run {run.name}", ""] + lines + tail) + "\n")
print(f"  {drafted} drafted, {len(dropped)} not drafted, {len(skipped)} Leads rejected; check list {q}")
print(f"  Claude usage this run: ${cost:.2f}")

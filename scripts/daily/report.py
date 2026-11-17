#!/usr/bin/env python3
"""Write the Author's check list for one make daily run (ADR-0009).

usage: report.py <run dir> <yyyy-mm-dd>
Reads picked.json and each <slug>.draft.json / <slug>.check.json (claude -p
JSON output) from the run dir; writes .scratch/daily/<date>/review-queue.md
(appending when several runs share a day) and prints the run summary.
"""
import json, os, pathlib, re, sys

run, day = pathlib.Path(sys.argv[1]), sys.argv[2]

# The self-edit checklist (docs/editorial/self-edit-checklist.md): one
# unticked box per "- [ ]" line under its "## The checklist" groups.
CHECKLIST = pathlib.Path(__file__).resolve().parents[2] / "docs/editorial/self-edit-checklist.md"
CHECKLIST_LINK = "../../../docs/editorial/self-edit-checklist.md"  # from .scratch/daily/<date>/
NOTE = re.compile(r"\s*\*\(([^*]*)\)\*\s*$")

def checklist():
    """[(group, line, [notes])] from the checklist doc ([] when it is missing)."""
    items, group, inside = [], None, False
    text = CHECKLIST.read_text() if CHECKLIST.exists() else ""
    for raw in text.splitlines():
        if raw.startswith("## "):
            inside = raw == "## The checklist"
        elif inside and raw.startswith("### "):
            group = raw[4:].strip()
        elif inside and group and raw.startswith("- [ ] "):
            text, notes = raw[6:].rstrip(), []
            while (m := NOTE.search(text)):
                notes.insert(0, m.group(1))
                text = text[:m.start()]
            items.append((group, text, notes))
    return items

KIND_NAMES = {"brief": "Brief", "case-study": "Case Study", "concept": "Concept",
              "comparison": "Comparison", "repo-profile": "Repo Profile"}

def self_edit(kind):
    """The Author's unticked self-edit boxes for one drafted Piece: a Brief
    drops the "(not in a Brief)" lines, every other kind the "(Brief only)"
    ones, and a "(<Kind>: …)" note is shown only on that kind's entry."""
    name, boxes = KIND_NAMES.get(kind, kind), []
    for group, text, notes in CHECKLIST_ITEMS:
        if ("not in a Brief" in notes and kind == "brief") or ("Brief only" in notes and kind != "brief"):
            continue
        own = [n.split(": ", 1)[1] for n in notes if n.startswith(name + ": ")]
        boxes.append(f"- [ ] {group}: {text}" + "".join(f" ({n})" for n in own))
    if not boxes:  # never drop the step: one box pointing at the doc
        boxes = [f"- [ ] Self-edit the draft with `docs/editorial/self-edit-checklist.md` (not found when this list was written)"]
    return [f"**Self-edit before you verify** ([checklist]({CHECKLIST_LINK})), ticked by you while you read the draft:"] + boxes + [""]

CHECKLIST_ITEMS = checklist()
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
    if p.get("fills"):
        lines += [f"**Fills Gap Day {p['fills']}** (ADR-0012): `published` is that day; `make release` keeps it.", ""]
    if d.get("author_checks"):
        lines += ["**Drafter asks you to check:**"] + [f"- [ ] {x}" for x in d["author_checks"]] + [""]
    if c.get("unconfirmed"):
        lines += ["**Fact-checker could not confirm:**"] + [f"- [ ] {x}" for x in c["unconfirmed"]] + [""]
    lines += self_edit(p["kind"])
    if c.get("fixed") or c.get("removed"):
        lines += ["<details><summary>What the fact-checker changed</summary>", ""] \
            + [f"- Fixed: {x}" for x in c.get("fixed", [])] + [f"- Removed: {x}" for x in c.get("removed", [])] + ["", "</details>", ""]

skipped = select.get("structured_output", {}).get("skipped", [])
out = pathlib.Path(".scratch/daily") / day
out.mkdir(parents=True, exist_ok=True)
q = out / "review-queue.md"
head = [] if q.exists() else [f"# Daily review queue, {day}", "",
        "Drafted by `make daily` (ADR-0009). Every Piece is `draft: true` and unverified. "
        f"First tick each draft's self-edit boxes while you read it, before you verify ([the checklist]({CHECKLIST_LINK})): "
        "they check the writing, not the facts, and the self-edit does not replace verification. "
        "Then verify every claim against its Source, then `make release` (ADR-0012): it sets `verified:`, removes `draft: true` and keeps each Publication Day. "
        "Ticked self-edit boxes never stand in for `verified:`.", ""]
tail = ["### Not drafted", ""] + (dropped or ["- none"]) + ["", "<details><summary>Leads the editor rejected</summary>", ""] \
    + [f"- `{s['id']}`: {s['reason']}" for s in skipped] + ["", "</details>", ""]
with open(q, "a") as f:
    f.write("\n".join(head + [f"## Run {run.name}", ""] + lines + tail) + "\n")
print(f"  {drafted} drafted, {len(dropped)} not drafted, {len(skipped)} Leads rejected; check list {q}")
print(f"  Claude usage this run: ${cost:.2f}")

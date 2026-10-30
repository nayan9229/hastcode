#!/usr/bin/env python3
"""Write the inputs `make daily`'s selector reads (ADR-0009).

usage: candidates.py <state.json> <queue.md> <content dir> <now RFC3339> <out dir>

Writes <out>/leads.json (open, uncovered Leads, newest first, each with its
age in hours and whether it can still be a Brief: at most 72 hours old by
its date, published else first seen, the rule make lead applies) and
<out>/explainers.json (every Explainer and Brief slug with title, kind,
Vertical and whether it is live; a Brief may link only a live Explainer).
"""
import datetime, json, pathlib, re, sys

state, queue, content, now, out = sys.argv[1:6]
now = datetime.datetime.fromisoformat(now.replace("Z", "+00:00"))
out = pathlib.Path(out)
out.mkdir(parents=True, exist_ok=True)

covered = set(re.findall(r"^- `([a-z0-9]+)`.*· covered", pathlib.Path(queue).read_text(), re.M))
leads = []
for lead in json.load(open(state))["leads"].values():
    if lead.get("decision") != "open" or lead["id"] in covered:
        continue
    date = lead.get("published") or lead.get("first_seen")
    age = (now - datetime.datetime.fromisoformat(date.replace("Z", "+00:00"))).total_seconds() / 3600
    leads.append({k: lead.get(k) for k in ("id", "title", "link", "source", "kind", "vertical")}
                 | {"date": date, "age_hours": round(age, 1), "brief_ok": age <= 72})
leads.sort(key=lambda l: l["age_hours"])
json.dump(leads, open(out / "leads.json", "w"), indent=1)

def front(path):
    text = path.read_text()
    fm = text.split("\n---\n", 1)[0]
    get = lambda k: (re.search(rf"^{k}:\s*\"?(.*?)\"?\s*$", fm, re.M) or [None, ""])[1]
    return {"title": get("title"), "kind": get("kind"), "live": not re.search(r"^draft:\s*true", fm, re.M)}

pieces = []
for p in sorted(pathlib.Path(content, "explainers").glob("*/*.md")):
    pieces.append({"slug": p.stem, "vertical": p.parent.name, "type": "explainer"} | front(p))
for p in sorted(pathlib.Path(content, "briefs").glob("*.md")):
    pieces.append({"slug": p.stem, "type": "brief"} | front(p))
json.dump(pieces, open(out / "explainers.json", "w"), indent=1)
print(f"{len(leads)} open Leads ({sum(l['brief_ok'] for l in leads)} fresh enough for a Brief); "
      f"{sum(p['type']=='explainer' and p['live'] for p in pieces)} live Explainers")

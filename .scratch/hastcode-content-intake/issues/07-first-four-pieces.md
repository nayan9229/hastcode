# 07: First four pieces, by hand

**What to build:** The new kinds are proven on real material before any template is fixed (ADR-0007, format before pipeline). A one-off Intake run proposes three Leads each for a Case Study, a Concept and a Comparison, with a line on why each fits. The Author picks one per kind. Each is drafted in a Claude session with the house-byline voice guide and the editorial rules (first-party Sources only, every claim linked, short credited quotes, diagrams redrawn in Mermaid, code only under its licence), as `draft: true`. The Author verifies each and sets `verified:` and `published` themselves. Then a Brief from a fresh Lead (at most 3 days old) that links one of the three.

Spec: `../spec.md` ("Format before pipeline"). ADR-0007. Voice: `docs/editorial/house-byline-voice.md`.

**Blocked by:** 01 (verified gate), 06 (real Watchlist)

**Status:** ready-for-human

- [ ] Three Lead proposals per kind recorded here; the Author's picks recorded
- [ ] One Case Study, one Concept, one Comparison drafted, verified by the Author (`verified:` set by the Author only), passing every build gate
- [ ] One Brief from a Lead ≤ 3 days old, linking one of them, verified by the Author
- [ ] Notes on what worked and what did not, per kind, for ticket 08

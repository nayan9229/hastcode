# 07: First four pieces, by hand

**What to build:** The new kinds are proven on real material before any template is fixed (ADR-0007, format before pipeline). A one-off Intake run proposes three Leads each for a Case Study, a Concept and a Comparison, with a line on why each fits. The Author picks one per kind. Each is drafted in a Claude session with the house-byline voice guide and the editorial rules (first-party Sources only, every claim linked, short credited quotes, diagrams redrawn in Mermaid, code only under its licence), as `draft: true`. The Author verifies each and sets `verified:` and `published` themselves. Then a Brief from a fresh Lead (at most 3 days old) that links one of the three.

Spec: `../spec.md` ("Format before pipeline"). ADR-0007. Voice: `docs/editorial/house-byline-voice.md`.

**Blocked by:** 01 (verified gate), 06 (real Watchlist)

**Status:** ready-for-human

- [ ] Three Lead proposals per kind recorded here; the Author's picks recorded
- [ ] One Case Study, one Concept, one Comparison drafted, verified by the Author (`verified:` set by the Author only), passing every build gate
- [ ] One Brief from a Lead ≤ 3 days old, linking one of them, verified by the Author
- [ ] Notes on what worked and what did not, per kind, for ticket 08

## Proposals (2026-10-26, from the first real Intake run, 189 Leads)

| Kind | Lead | Source | Vertical |
|---|---|---|---|
| Case Study | **`pk4mvo2l` Single-tenant Titus: VM-level isolation on a multi-tenant container platform** | Netflix TechBlog | Infrastructure & Cloud |
| Case Study | `g7w5f2tx` How D1 scales to millions of databases | Cloudflare blog | Databases & Storage |
| Case Study | `rdsbakwp` 2026 World Cup: reliability at scale | Meta Production Engineering | Distributed Systems |
| Concept | **`2asmegc4` Metastability and Lost Writes** (metastable failures) | Marc Brooker | Distributed Systems |
| Concept | `zkmxgotq` Push-Based KV Transfer for Disaggregated Inference | vLLM blog | AI/ML Systems |
| Concept | `5ldu3gbw` JEP 543: Structured Concurrency | Inside Java | Languages & Runtimes |
| Comparison | **`dwrxqhav` Better HTTP/1 connection reuse → HTTP/1.1 vs HTTP/2 vs HTTP/3 connection reuse** | Go blog | APIs & Networking |
| Comparison | `oqyj2tvr` Shard design → hash vs range vs directory sharding | PlanetScale | Databases & Storage |
| Comparison | `xtsbvhnf` Elastic Expert Parallelism → tensor vs pipeline vs expert parallelism | vLLM blog | AI/ML Systems |

**Author's picks (2026-10-26):** Titus (Case Study), metastable failures (Concept), HTTP connection reuse (Comparison). Bold above.

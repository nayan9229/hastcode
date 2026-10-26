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

## Draft: Case Study, Netflix single-tenant Titus (`infrastructure-cloud/netflix-single-tenant-titus`)

~2,630 words, draft, placeholder hero. Lead read via Netflix's RSS (the blog returns 403 to direct fetches); three earlier Netflix posts described as the new post describes them. Labelled inferences: per-cell autoscaler/scheduler; single-tenant tag as Kubernetes metadata; the 1-vCPU overhead arithmetic is our hypothetical.

Check before verifying:
- Netflix wrote "Expected Departure Time"; Cilium's docs say "Earliest Departure Time" (piece says "the approach Cilium's Bandwidth Manager calls EDT").
- "No benchmark scores published" (the post's measurement sentence reads garbled in the feed).
- Priority-class bullet "never evicted to make room for another pod" is a gloss of "exempts the task from preemption".
- The three microVM costs, incl. "nested virtualization was not ready at the time".
- Shared fleet packs denser: Netflix's own report, no numbers.
- Headline numbers: "millions of containers a week across tens of thousands of EC2 instances"; r7i.2xlarge vs m7a.8xlarge warm-pool example.

## Draft: Comparison, HTTP connection reuse (`apis-networking/http-connection-reuse`)

~2,450 words, draft, placeholder hero. Sources: Go blog post (CC BY 4.0), RFC 9112 §9, RFC 9113, RFC 9114, RFC 9000, Go `net/http` docs (go1.27.3). Labelled inferences: HTTP/3 from Go needs a non-stdlib client; idle cap default 2 causes re-dials above 2 concurrent; IdleConnTimeout below the server's makes close races rarer; coalescing matters mainly with a shared wildcard cert behind one proxy; the whole "When each wins" section.

Check before verifying:
1. Intro's 710 ms / 150 ms figures and the "on my machine" caveat.
2. RFC 9112 §9 says nothing about cross-host connection sharing (the Coalescing section and table's "Not described" rely on it).
3. GOAWAY + RFC 9113 §8.7 back the "even a POST" retry claim.
4. DefaultTransport vs custom Transport HTTP/2 defaults (custom Transport with Protocols nil is usually HTTP/1 only).
5. disable_active_migration wording: piece says "a server can forbid", RFC 9000 §9 says "the peer".
6. Every quote verbatim at its anchor.
7. Replace the hero.

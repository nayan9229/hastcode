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

## Draft: Concept, metastable failures (`distributed-systems/metastable-failures`)

~2,500 words, draft, placeholder hero. Opens on the 2015 DynamoDB retry storm. 13 Sources: Brooker's Lead post and four earlier posts, HotOS'21 (Bronson et al.), OSDI'22 (Huang et al.) and its reproduction repo, AWS 2015 DynamoDB summary, AWS "Exponential Backoff And Jitter", two Builders' Library articles, the Raft paper, Stability Sim. Labelled inferences (both in "Where lost writes fit"): how a member lying about fsync() can lose a committed Raft entry; why Brooker pairs lost writes with metastability.

Check before verifying:
1. OSDI counts: abstract 22 failures vs body 21; draft says "22 from 11 organisations", "retries in more than half", "4 of 15 AWS outages".
2. "Fix to break" passage keeps OSDI's wording; no claim the SimpleDB fix caused the DynamoDB event.
3. Jitter result "call count fell by more than half (100 clients)": source baseline ambiguous.
4. Token-bucket paraphrase "roughly 0.1 retries well above 10% failures".
5. HotOS worked example: 280 → 560 QPS; recovery below 150 QPS or retries below 20 QPS.
6. Builders' Library URLs now redirect to builder.aws.com.
7. The lost-writes inference.
8. Replace the hero.

## Draft: Brief, GitHub Dependabot fan-out (`briefs/github-dependabot-fan-out`)

Lead `wnzcpeld` (GitHub status, incident first posted 2026-10-23 14:58 UTC, resolved 18:51 UTC). Links the metastable-failures Concept. 173 words, draft, placeholder hero. One labelled inference: capacity alone did not clear the backlog quickly (the report does not say so directly). Freshness: the 3-day rule ends 2026-10-26 between 14:58 and 18:51 UTC depending on which post counts; drafted 17:30 UTC.

Check before verifying:
1. Both quotes verbatim against the incident report.
2. "Updates from 15:39 to 16:43 UTC describe more processing capacity rolling out".
3. Whether "the durable fix was not more consumers but a separate lane with a fixed rate" reads as fair to the report.
4. The metastable link is a lesson by analogy: the report never says "metastable".
5. Replace the hero.

## Status 2026-10-26

All four un-drafted at the Author's request without `verified:` (published 2026-10-26, outside the gate's date rule). Remaining for this ticket: the Author checks each list above and adds `verified:`. Tickets 08 and 09 proceed from the four pieces as written.

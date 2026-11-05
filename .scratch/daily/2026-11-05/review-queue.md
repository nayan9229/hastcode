# Daily review queue, 2026-11-05

Drafted by `make daily` (ADR-0009). Every Piece is `draft: true` and unverified. Verify every claim against its Source, then set `verified:` and remove `draft: true` yourself (ADR-0007).

## Run 20261105-125621

### Fly.io's flaky Petsem: SQLite has one writer, so keep writes bounded

`content/briefs/fly-io-petsem-sqlite-lock-contention.md` · brief · 244 words · fact-checker: 17 claims checked, 3 fixed, 0 removed

**Publish within 72 hours of the Lead's date**, with or after its Explainer: cloudflare-d1

**Drafter asks you to check:**
- [ ] Date: page shows 'October 15 ... (09:16UTC)' with no year; the infra-log index also gives no year (entry sits above posts dated 2025). Queue dates the post 2026-11-04. Brief states '15 October, 09:16 UTC (the page gives no year)'. Confirm the year and the 3-day window before publishing.
- [ ] Quote (the only one): 'tipped the database into sustained lock contention', verbatim from https://fly.io/infra-log/petsem-was-flaky/, credited to Fly.io.
- [ ] Numbers: 'a little over half of write requests started failing'; replicas and reads from the primary 'were not affected'. Brief says 'a little over half of write requests failed, while replicas and reads from the primary stayed healthy'.
- [ ] Source says 'runs on a SQLite primary in IAD'; Brief says 'one SQLite primary in IAD'. 'One' is our reading of 'a SQLite primary'.
- [ ] Source motive: 'exploring a change away from SQLite, which may or may not happen, but which requires more record tracking'. Brief: 'A deployed change, groundwork for a possible move off SQLite, added record tracking.'
- [ ] Fix: Source says 'built and released a new version of Petsem without the tracking code, which stopped the intermittent errors'. Brief: 'A new build without the tracking code stopped the errors.'
- [ ] IAD = 'Ashburn, Virginia (US)' per the Fly.io Regions table at https://docs.fly.io/reference/regions (fly.io/docs/reference/regions/ 301-redirects there).
- [ ] SQLite one-writer claim: https://www.sqlite.org/whentouse.html, 'High Concurrency' bullet: 'SQLite supports an unlimited number of simultaneous readers, but it will only allow one writer at any instant in time ... Writers queue up.'
- [ ] WAL claim: https://www.sqlite.org/wal.html section 2.2 Concurrency: 'writers and readers can run at the same time. However, since there is only one WAL file, there can only be one writer at a time.'
- [ ] Labelled inference: Petsem runs in WAL mode (from healthy reads beside failing writes); the Source does not name the journal mode. Also implicit: the contended lock is SQLite's write lock; Source says 'lock contention' after 'extra write volume'.
- [ ] Explainer link: /databases-storage/cloudflare-d1/ for 'scales by adding databases'; the Case Study cites D1 docs calling each database 'inherently single-threaded' and says an application grows by adding databases, not by growing one.
- [ ] Frontmatter: sources: added with the two SQLite pages; published left at pick time 2026-11-05T13:01:30Z (provisional); draft: true kept; no verified:. Title 68 chars; description 146 chars; tags [sqlite, locking, incidents, fly-io].

**Fact-checker could not confirm:**
- [ ] Year of the '15 October' incident heading: page and index give none; inferred 2026 from the feed's 2026-11-04 entry date and the entry sitting first in feed above a 2026-10-23 entry
- [ ] 3-day window: Lead queue date 2026-11-04 (feed published date) is within 72h of 2026-11-05, but the incident the post describes is dated 15 October, about three weeks earlier; Author to decide whether the post date or the incident date governs

<details><summary>What the fact-checker changed</summary>

- Fixed: Lead sentence: replaced 'entry dated 15 October, 09:16 UTC (the page gives no year)' with 'entry of 4 November 2026' linked to https://fly.io/infra-log/feed.xml (Atom <published> 2026-11-04T00:00:00+00:00), 'headed 15 October, 09:16 UTC' linked to the post; page heading reads 'October 15: Petsem was flaky for a bit (09:16UTC)' with no year anywhere on page or index
- Fixed: 'one SQLite primary' -> 'a SQLite primary' in body and description; Source: 'Our secrets store, Petsem, runs on a SQLite primary in IAD.'
- Fixed: Unlinked sentence 'A deployed change ... added record tracking.' now links the Source post (Source: 'requires more record tracking in the database' / 'Some time after being deployed, the extra write volume from this change')

</details>

### vLLM's Kimi K3.1 Pro day 0: merge shape and cache layout set the cost

`content/briefs/vllm-kimi-k31-pro-day-0-support.md` · brief · 241 words · fact-checker: 17 claims checked, 4 fixed, 1 removed

**Publish within 72 hours of the Lead's date**, with or after its Explainer: decode-context-parallelism, disaggregated-inference

**Drafter asks you to check:**
- [ ] Date and byline: post dated 'November 4, 2026', byline 'Inferact and the vLLM Team'; no time of day given, so the Brief says only 'On 4 November'. Lead picked 2026-11-05, within the 3-day window.
- [ ] Sparse attention claim: 'keeps the top 2,048 tokens per query' against the post's 'scores every earlier token for each query with 32 lightweight heads. Then it keeps the top 2,048 tokens.' (section 'Sparse Attention in K3.1 Pro', anchor #sparse-attention-in-k31-pro).
- [ ] Generic merge: post says 'The generic merge sends every GPU's full local top-2048, as FP32 score and token id pairs, to every other GPU.' Brief writes 'top-2,048' with a comma.
- [ ] Packed merge: post says 'Each GPU picks only its best 1.5 × 2048 / W candidates straight from the BF16 scores and packs each with its score into 32 bits, about 10× fewer bytes.' Brief paraphrases as 'sends fewer candidates, each id packed with its score into 32 bits, about 10× fewer bytes' (anchor #an-exact-top-k-merge-with-10-less-traffic).
- [ ] Speedup: Figure 5 caption 'One DSA attention module at DCP8 on 8× GB300 with 5 query tokens per request. The exact one-round merge is 1.27–1.51× faster than the generic merge.' TL;DR rounds this to '1.3–1.5× faster at DCP8'. Brief uses 1.27–1.51× and 'eight GB300s'; confirm the figure is the whole attention module, as the TL;DR states, and note the 5-query-token, 64-request setup is omitted for length.
- [ ] KV cache layout: post says 'vLLM packs the cache by block instead of by layer.' Brief: 'packed by block, not by layer' (anchor #a-kv-cache-for-three-kinds-of-state).
- [ ] The one quote: 'one large RDMA read instead of dozens of small per-layer reads' is verbatim from the 'Disaggregated serving' bullet: 'Each block is one contiguous range of memory, so a P/D transfer moves it with one large RDMA read instead of dozens of small per-layer reads.' Credited as 'the post says'.
- [ ] 'every attention layer of every decode step must reconcile what the GPUs hold' links the 'Long context' section, which says 'the global top-2048 has to be merged across GPUs in every layer before attention starts'. Confirm anchor #long-context-decode-context-parallelism-with-a-sparse-indexer resolves (headings' ids inferred from the page's table of contents, not from raw HTML).
- [ ] Inference labelled: 'We infer a rule the post does not state: design the per-layer merge and the cache layout together, before adding GPUs.' Also check 'the shape of that exchange sets the cost, not the GPU count' reads as our analysis; the post's own framing is 'DCP trades per-step latency for KV capacity'.
- [ ] 'All figures are vendor-reported': every number in the Brief is from the vLLM post's own measurements by Inferact and the vLLM team; no third-party figure used.
- [ ] All four anchors (#sparse-attention-in-k31-pro, #an-exact-top-k-merge-with-10-less-traffic, #a-kv-cache-for-three-kinds-of-state, #long-context-decode-context-parallelism-with-a-sparse-indexer) are inferred from the post's table of contents; click each on the live page.
- [ ] Word count was done by hand (Bash denied for a count script): 238 body words, 241 with the > and ## markers; the build's brief-length check did not fire. Description is 137 characters, title 69.

**Fact-checker could not confirm:**
- [ ] Anchor targets: all four fragments appear as hrefs in the page's table of contents, but DOM id attributes were not inspected (curl denied); clicking on the live page still advised
- [ ] 'every attention layer of every decode step': post says 'merged across GPUs in every layer before attention starts' and 'DCP trades per-step latency for KV capacity'; 'every decode step' is a light extrapolation
- [ ] 1.27–1.51× attributed to 'the attention module': caption says 'The exact one-round merge is 1.27–1.51× faster than the generic merge' measured on 'One DSA attention module'; TL;DR says 'makes the attention module 1.3–1.5× faster at DCP8'; arithmetic (merge 258 µs of 671 µs module) supports module-level, but the post never states it explicitly
- [ ] Figure 5 setup '5 query tokens per request' omitted from the Brief for length

<details><summary>What the fact-checker changed</summary>

- Fixed: description: 'packs the DCP top-k merge into 32 bits' -> 'packs each DCP top-k candidate into 32 bits' (post packs each candidate with its score, not the merge)
- Fixed: 'as FP32 pairs' -> 'as FP32 score-and-id pairs' (post: 'as FP32 score and token id pairs')
- Fixed: 'a prefill-to-decode transfer is "one large RDMA read…"' -> 'transfer uses' (post: 'moves it with one large RDMA read')
- Fixed: 'the shape of that exchange sets the cost, not the GPU count' -> 'the shape of that exchange, we infer, sets more of the cost than the GPU count' (labelled inference; post frames DCP as trading per-step latency for KV capacity)
- Removed: 'Both changes follow from one choice.' (unsupported: the P/D transfer cost comes from disaggregated serving, not from sharding the cache by token)

</details>

### Not drafted

- none

<details><summary>Leads the editor rejected</summary>

- `5jswvc4b`: Fetched. Google Cloud europe-west2 power incident now has a full report with root cause (PLC logic misconfiguration after a utility cable fault), but the page's first post was 20 Oct 2026 08:35 UTC, so it is outside the 72-hour Brief window counted from first post; Explainer slots are 0 today. Strong Case Study candidate for a later day.
- `qwggc3mr`: Fetched. GitHub Actions incident: no cause stated, root cause analysis promised later.
- `s6vjv67p`: Fetched. GitHub Actions and Pages incident: no cause stated.
- `3ittfw7g`: Fetched. GitHub Pull Requests, Git Operations and Issues incident: no cause stated.
- `eanovtqi`: Fetched. Copilot Mission Control: only partial, unlinked explanations (upstream database throttling, then a specific client version); no live Explainer teaches that lesson cleanly.
- `intnmdgl`: Fetched. GitHub Git Operations incident: mitigation rolled out, no cause stated.
- `coq3fsuq`: Fetched. GitHub API, Git and Copilot incident: no cause stated.
- `liiwq2tm`: Fetched. GPT-5.4 Nano errors: cause is only 'an issue with an upstream model provider'; window also closes 2026-11-05T17:46Z.
- `m2nyvca7`: Fetched. Go 1.27's encoding/json/v2 is a real design change, but no live Explainer teaches its lesson (a new major version with v1 reimplemented on v2 for gradual migration). Possible Concept later.
- `pzqcdif2`: Fetched. Kinference is a ~280-word product announcement by Buoyant with vendor benchmarks and no mechanism beyond feature names.
- `bclvs75y`: Fetched. rustup 1.29.3 is a two-item bug-fix point release.
- `2d3e7tr3`: Fetched. APNIC guest commentary on Bangladesh outages and dependency mapping; not a first-party account of an event.
- `22ir6vsm`: Meta F3 Lightning is a Case Study candidate, not a Brief (no event); Explainer slots are 0 today.
- `rbq2cyng`: Duplicate of 22ir6vsm; same reason.
- `p4vynmq3`: Cloudflare DNS targeted zone builds: engineering change, but 70.9h old at intake with the window closing 2026-11-05T14:00Z, and no live Explainer teaches incremental builds.
- `gcdsrrow`: AWS large-transaction binlog contribution to MySQL: no live Explainer on replication or binlogs, and the window closes 2026-11-05T16:57Z.
- `vsfmm6jz`: Same Google Cloud incident page as 5jswvc4b, older queue entry; out of window.
- `cpeb3uo6`: Podcast with no readable text.
- `grtrs2bn`: Call for proposals.
- `mvpdxzbi`: Google Research post on deliberation methodology; not an event at a system.
- `oomzdelu`: Cloudflare product feature announcement.
- `ek72wcz6`: Solar PV design guide; outside the Verticals' lesson set, not an event.
- `tmbblbn7`: Devoxx talk write-up for Project Detroit; not a release or event with a design change.
- `xy35vrlu`: AWS Kiro dashboard tutorial; marketing.
- `36fas6c4`: Cloudflare Zero Wire hub launch; announcement.
- `7k3wqvk5`: Meetup recap; community post.
- `mgjuzw4t`: curl post on scrutinising security assessments; process opinion, not an event.
- `pqb2tzwf`: IBM Vault Dedicated product announcement.
- `vfpsj2yq`: APNIC opinion restating the case for DNSSEC; no event.
- `6btegobi`: DuckDB quant desktop post; product positioning, no incident or design change.
- `6cyiigrd`: PlanetScale latency explainer post; a concept, not an event.
- `qvst6zr3`: pgbadger point release.
- `vj2epacc`: SIG announcement.
- `3qhfgrlb`: AWS DR guidance series; no named system or event.
- `iidnbskk`: AWS DR guidance series; no named system or event.
- `v63c76og`: AWS DR guidance series; no named system or event.
- `t2oik657`: AWS DRS agent lifecycle how-to; no event.

</details>


# Daily review queue, 2026-11-16

Drafted by `make daily` (ADR-0009). Every Piece is `draft: true` and unverified. Verify every claim against its Source, then `make release` (ADR-0012): it sets `verified:`, removes `draft: true` and keeps each Publication Day.

## Run 20261116-093749

### Cyber Chaos Monkey: how Shopify lets AI agents attack production

`content/explainers/infrastructure-cloud/shopify-cyber-chaos-monkey.md` · case-study · 2756 words · fact-checker: 62 claims checked, 12 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Lead was read through WebFetch, which would not return the whole text word for word. The facts were collected over four targeted passes, so check every quote against the live page: "a rising tide that lifts all boats.", "Attackers target our production systems, so our testing has to as well.", "Our egress proxy is effectively an allowlist on top of nginx,", "instantly kills the offensive agent.", "a bet on asymmetry", and the TPI definition "the idea that a second party must review the work of the acting party".
- [ ] The editor's angle says the agent holds only placeholder credentials. The Source says the agent is seeded with ephemeral Shopify credentials (e.g., a session cookie and CSRF token) that it can read but not modify; placeholders are only for knowledge services and model providers. The draft follows the Source. Check this wording.
- [ ] The diagram routes the executor's requests to Shopify production through the egress proxy. We infer this from 'egress proxy pods are the only ones capable of exiting the VPC', and the text says so. Check it against Shopify's sandbox-network figure.
- [ ] The 'recent incidents where AI agents in security testing reached systems they were never meant to touch' sentence came back truncated. Check its wording, and check that the post names no specific incident.
- [ ] The list of what the agent does in the sandbox ('install packages, build fuzzers, make live requests to Shopify and read production logs') comes from a summary of the post's six steps. Check it is the Source's wording.
- [ ] 'a report sent back to Discovery and upwards': check the word 'upwards' is in the Source.
- [ ] 'Keeping everything pod-local (but isolated) reduces the risk of network partitions' came back truncated. Check that the partitions are between the agents, which is our gloss.
- [ ] The Netflix Simian Army post returned 403 and was not read. The 2011 date and 'randomly shut down production instances' are credited to Shopify's post, not Netflix's.
- [ ] July harness figures ($50 to $300 per full scan, over 300 findings, over $400,000 in equivalent bounty payouts, fix-author agents creating a branch per finding) come from a WebFetch summary of Zack Deveau's 29 July 2026 post. Check them, and check the date and author.
- [ ] The Kubernetes sidecar facts (init containers with restartPolicy: Always, stable since v1.33, a sidecar does not prevent a Job completing) came from the docs page with anchors #pod-sidecar-containers and #jobs-with-sidecar-containers.
- [ ] The options section is built from three alternatives the post argues against implicitly. The post has no explicit options comparison, so check this framing.
- [ ] The related links /infrastructure-cloud/openshell/, /ai-ml-systems/meta-thin-orchestrator/ and /infrastructure-cloud/shopify-overexposed/ are set in related: and inline. Check that the Overexposed comparison ('would rather report inconclusive than guess') is fair.

**Fact-checker could not confirm:**
- [ ] Netflix Simian Army post (netflixtechblog.com/the-netflix-simian-army-16e57fbab116) returned 403; not read. 2011 and 'randomly shut down production instances' remain credited to Shopify's post; search suggests the original is techblog.netflix.com/2011/07/netflix-simian-army.html
- [ ] Shopify's sandbox-network figure checked only through its alt text, not the image itself

<details><summary>What the fact-checker changed</summary>

- Fixed: Opening: 'live credentials and everything it knows about its own systems' corrected to ephemeral live credentials plus the Source's list (source code, threat intelligence, past reports, infrastructure records); 'one tool call' -> 'a tool call'
- Fixed: TL;DR: 'cost of finding vulnerabilities towards zero for everyone' -> Source's 'cost of vulnerability discovery towards zero' across the industry
- Fixed: TL;DR: placeholder credentials scoped to internal services (Shopify session credentials are real, ephemeral)
- Fixed: TL;DR: 'fails closed' is not the Source's term; marked '(our term)' and link now carries the Source's wording (judge times out, crashes or fails to adjudicate)
- Fixed: Diagram note: our Step 1-6 sections distinguished from the post's own six steps
- Fixed: Diagram note: executor-through-egress path now stated as 'we infer', citing the VPC sentence and the sandbox-network figure (alt text: pod routes requests through credential and egress proxies)
- Fixed: 300% bug bounty increase now linked to the CCM post rather than only the bug bounty page
- Fixed: Options section: post has no options comparison; reframed as 'we infer three alternatives'
- Fixed: Options bullet: source code/threat intelligence/past reports/infrastructure records claim now linked; 'outside researchers do not have' marked as inference
- Fixed: Network partitions: Source says between orchestrator, executor and judge, not 'between the agents'
- Fixed: Netflix: 'Shopify borrows it from' -> 'Shopify's homage to' (Source's framing)
- Fixed: Overexposed comparison: CCM's 'preference for stopping over improvising' marked as our reading

</details>

### Durable execution: how workflows survive crashes by replaying a log

`content/explainers/distributed-systems/durable-execution.md` · concept · 2450 words · fact-checker: 74 claims checked, 8 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Word count is an estimate (about 2,450 words in the body). Bash was denied for wc, so it was not measured. No section-length warning came up, so every H2 is under 250 words.
- [ ] Lead quotes were read through WebFetch summaries because Bash/curl was denied. Check each quote against the live page word for word: the definition ('a programming model in which the runtime keeps a record of the result of each step that interacts with the outside world'), 'no new quotes, no second purchase order, no second request to the manager.', 'A persistent log records what the agent saw, not what actually happened.', 'Replay runs the new code against a log the old code wrote.' (from the 'What can go wrong' table, cited at #what-can-go-wrong), and 'plus others that arise from non-determinism'.
- [ ] Lead anchors come from the page's table of contents (#running-example, #agentic-loop, #what-can-go-wrong, #durable-execution, #making-the-loop-durable, #entity-and-workflow, #waiting, #retries, #durable-is-not-correct, #unseen-outcomes, #world-moves, #change, #memory, #untrusted-text, #conclusion). Confirm they resolve. The Write-Ahead Log sidebar is cited at #durable-execution; check that it sits in that section.
- [ ] Wake-up timings (14 ms on Temporal and about 0.1 s on Cloudflare in memory; about 2 s on Cloudflare and under 0.1 s on Temporal when rebuilt; no CPU while waiting) come from the Lead's footnote 4, cited at #waiting. Confirm the footnote is reachable there.
- [ ] The retry policy is shown as 'six attempts, starting two seconds apart and doubling'. The Lead's code is RetryPolicy(maximum_attempts=6, initial_interval=2s, backoff_coefficient=2). Confirm the paraphrase is fair.
- [ ] Temporal limits (51,200 events or 50 MB, warning at 10,240 events or 10 MB) are cited to /workflow-execution/limits. The termination wording on /workflow-execution/event#event-history-limits gives only the event count, not 50 MB.
- [ ] Cloudflare anchors: #name-steps-deterministically is confirmed. #do-not-rely-on-state-outside-of-a-step and #ensure-apibinding-calls-are-idempotent were inferred from the headings; check they resolve.
- [ ] The Restate URL docs.restate.dev/concepts/durable_execution served a 'Key Concepts' page. Confirm it is the canonical address and holds both quotes ('without changing how you write business logic' and 'replays the journal, skipping completed steps and resuming from exactly where it left off.').
- [ ] The DBOS anchors #how-workflow-recovery-works and #upgrading-workflow-code were inferred from the headings.
- [ ] Hook framing: the Lead says the agent asked the manager for approval and 'while it waited, the process crashed', and stayed down a day. The draft walks this as the approval wait in the diagram. That reading is ours; the Lead names no specific wait step.
- [ ] The paragraph that walks the first run lists the steps before the wait (quotes, draft PO, ask manager). It is reconstructed from the Lead's tool list and its restart sentence; check it matches the Lead's timeline.
- [ ] 'We infer' is used once, in the determinism section (the log answers 'what did step N return?').
- [ ] Two bylines added to Sources after checking them: Parallel Change is by Danilo Sato (13 May 2014), not Fowler, and Idempotent Receiver is by Unmesh Joshi. Published Interface is credited to Martin Fowler without checking.
- [ ] Frontmatter now sets related: [metastable-failures, agentic-development] (uncommented). Confirm you want it.
- [ ] Inngest is named only as one of the Lead's listed implementations; no first-party Inngest Source was added.

**Fact-checker could not confirm:**
- [ ] Quotes over 125 characters (the definition, the full Restate replay sentence) were checked through WebFetch, which only quotes short pieces. Each piece matched word for word, but I couldn't compare the raw HTML because Bash/curl was denied
- [ ] Cloudflare anchors #do-not-rely-on-state-outside-of-a-step, #ensure-apibinding-calls-are-idempotent and #use-conditional-logic-carefully are inferred from the heading names (Starlight slug pattern). Only #name-steps-deterministically is confirmed on the page
- [ ] Temporal anchors #worker-versioning and #workflow-patching come from the page's table of contents. The patching heading reads 'Versioning with Patching'
- [ ] Restate URL docs.restate.dev/concepts/durable_execution serves a page titled 'Key Concepts'. No redirect was reported, but I couldn't confirm it is the canonical address
- [ ] The closing line 'wait a week, survive a day-long crash, and still place exactly one order' joins three separate runs in the Lead. Read it as our summary, not one reported run
- [ ] Hook: the Lead does say the crash happened during the manager-approval wait ('While it waited, the process crashed and stayed down for a day'), so the diagram's reading is supported
- [ ] related: [metastable-failures, agentic-development]: both Pieces exist under content/explainers and the inline links match their paths. Keeping them is the Author's choice

<details><summary>What the fact-checker changed</summary>

- Fixed: TL;DR: 'returns it from the log rather than running again' anchor moved from #making-the-loop-durable to #durable-execution, where the Lead says it
- Fixed: Replay walk-through: 'runs again from the top' anchor moved to #durable-execution (the 'restarts from the top' sentence sits there)
- Fixed: Practice three quote 'Replay runs the new code against a log the old code wrote.' is in the 'Essential and accidental complexity' sidebar, not the 'What can go wrong' table; anchor changed to #EssentialAndAccidentalComplexity
- Fixed: Cloudflare: Math.random()/Date.now() warning is under 'Use conditional logic carefully', not the state section; added a link to #use-conditional-logic-carefully and used the page's spelling 'behavior'
- Fixed: Waiting: 'spend most of their time waiting' changed to 'much', as the Lead says
- Fixed: Idempotent Receiver: Joshi's pattern ignores duplicates using a client ID plus request number, not an 'idempotency key'. Split the sentence and linked the idempotency-key wording to the Lead (#retries)
- Fixed: Prompt injection: link text 'does not solve it directly' changed to 'does not help with it per se', the Lead's wording
- Fixed: Sources: Restate page title is 'Key Concepts', not 'Durable execution'; added dates for Idempotent Receiver (23 November 2023) and Published Interface (26 December 2003, by Martin Fowler, now checked)

</details>

### G1 garbage collector: regions, mixed collections and a soft pause goal

`content/explainers/languages-runtimes/g1-region-based-garbage-collection.md` · concept · 2300 words · fact-checker: 78 claims checked, 16 fixed, 4 removed

**Fills Gap Day 2026-11-13** (ADR-0012): `published` is that day; `make release` keeps it.

**Drafter asks you to check:**
- [ ] WebFetch returned only summaries with partial quotes, never exact page text: check every quotation against the live dev.java and Oracle JDK 27 pages.
- [ ] The Oracle #GUID-... anchors came from each page's in-page link lists, not from the HTML ids: confirm each one jumps to the intended heading.
- [ ] Humongous threshold: dev.java says 'larger than half' a region, Oracle says 'larger or equal the size of half a region'. The draft uses 'larger than half'.
- [ ] Oracle's G1 page spells two flags two ways (MaxGCPauseMillis vs MaxGCPauseTimeMillis, GCPauseIntervalMillis vs GCPauseTimeInterval). The draft uses MaxGCPauseMillis and GCPauseIntervalMillis.
- [ ] The -Xlog:gc+phases=debug mention is now unlinked; confirm which Oracle section states it.
- [ ] openjdk.org returned 403, so JEP 523 is linked only as cited by dev.java and its content was not read.
- [ ] The dev.java 'single CPU and less than 2 GB' wording is a partial quote; check its context.
- [ ] The Lead page shows no date: 'November 2026' in Sources comes from the URL (2026/11/10).
- [ ] We infer: G1HeapWastePercent is what ends the Space-Reclamation phase (about 410 MB on the example heap); why SATB suits G1; why regions beat whole-old-generation collection; that a lower pause goal means fewer old regions per mixed pause.
- [ ] Worked-example numbers are our own arithmetic on documented defaults: 8 GB heap, 4 MB regions, 2 MB humongous threshold, 400 candidates / 8 = 50 regions (200 MB), 8,192 cards per region.
- [ ] The AWS HBase post does not say whether MaxGCPauseMillis was changed; the hook says so. Its exact post title was taken from the URL slug and fetch summary; confirm it.
- [ ] The 10% refinement default (G1RSetUpdatingPauseTimePercent) is linked to the tuning guide's Tunable Defaults table anchor; confirm it is listed there.
- [ ] In the hook, the middle quote 'pause of approximately 3021ms' carries no link of its own; it comes from the same AWS post as the linked quotes before and after it.
- [ ] Oracle's docs do not describe G1's write barrier or card table, so the draft leaves them out.

**Fact-checker could not confirm:**
- [ ] JEP 523 content: openjdk.org returned 403; claim rests on dev.java's statement that JEP 523 made G1 the default in all environments
- [ ] Oracle #GUID anchors: confirmed against each page's own in-page Topics links; raw HTML ids not inspected (curl was blocked)
- [ ] Quotes were checked through WebFetch's sentence-level verbatim excerpts, not raw page text
- [ ] Inside Java 'November 2026' date comes only from the URL; the page shows no date
- [ ] dev.java's 'better' is in double quotes in the original; the draft nests it as 'better' inside a quotation
- [ ] Hook quote 'pause of approximately 3021ms' still has no link of its own (it is verbatim, from the same AWS post linked on either side)

<details><summary>What the fact-checker changed</summary>

- Fixed: Hook: 'counted what each RegionServer's JvmPauseMonitor reported... it logged' -> YCSB benchmark, RegionServer reports unusually long pauses via JvmPauseMonitor; 'they saw' (AWS post: 'During YCSB benchmarking, G1 GC exhibited 75 such pauses')
- Fixed: TL;DR default claim: anchor #a-brief-history-of-g1 -> #introduction; 'across all system configurations' -> 'in all environments' (dev.java wording)
- Fixed: Basic Concepts glosses (generational/incremental/STW/...) prefixed 'In our words:'; Oracle lists the terms but defines only evacuating
- Fixed: Humongous threshold: now 'half a region or more' linked to Oracle Humongous Objects ('larger or equal the size of half a region'), noting dev.java says 'larger than half'; 'skips eden and survivor' linked to dev.java; matches the worked example's '2 MB or more'
- Fixed: IHOP 45%: linked to Ergonomic Defaults table (not stated in the IHOP section)
- Fixed: Full GC trigger: 'When collections cannot free enough memory' -> 'If the application runs out of memory' (Oracle cycle section)
- Fixed: G1HeapWastePercent: Oracle defines it as allowed unreclaimed space applied when choosing candidates; now sourced to the defaults table, '410 MB on a full 8 GB heap', end-of-phase link kept as 'We infer'
- Fixed: Remembered set: 'those cards become heap roots' -> 'the remembered sets determine the heap roots' (Oracle); 'dropping duplicates' -> 'dropping many duplicates'
- Fixed: 'G1 keeps them only where needed' -> 'G1 builds them mostly lazily', linked to Remembered Set section
- Fixed: Cards / 512 bytes: linked to Remembered Set section
- Fixed: Refinement spill-over: linked to Tuning for Throughput
- Fixed: 1792 MB / two-core selection: anchor #a-brief-history-of-g1 -> #using-g1
- Fixed: 'built for multiple processors and tens of GB': now linked; 'working with a single CPU' -> 'performing well with a single CPU' (dev.java context)
- Fixed: Latency direction: linked to General Recommendations ('modify the pause time target')
- Fixed: -Xms=-Xmx + AlwaysPreTouch and Automatic Heap Sizing: linked to dev.java #using-g1
- Fixed: -Xlog:gc+phases=debug: linked to Oracle Garbage Collection Process section
- Removed: 'Every fix moves a cost rather than removing it...' (unsourced, for size)
- Removed: dev.java tutorial clause in the closing paragraph (duplicated a Source, for size)
- Removed: 'Headroom is also what SATB's conservatism and humongous leftovers quietly consume.' (unsourced, for size)
- Removed: 'which is what the headroom advice below supplies' (for size)

</details>

### Postgres migration: dump/restore vs logical vs trigger replication

`content/explainers/databases-storage/dump-restore-vs-logical-vs-trigger-replication.md` · comparison · 2500 words · fact-checker: 68 claims checked, 11 fixed, 0 removed

**Fills Gap Day 2026-11-14** (ADR-0012): `published` is that day; `make release` keeps it.

**Drafter asks you to check:**
- [ ] Lead copy-speed figures: 50-100 GB/h with indexes in place; about 60 to nearly 300 GB/h after deferral; rebuild about 2 hours; about 2.5x faster overall
- [ ] 1 TB arithmetic (about 17 h at 60 GB/h, about 3.3 h at 300 GB/h) is labelled inference and uses migrator COPY rates, not pg_dump rates
- [ ] Logical replication cutover steps (catch up, set sequences, repoint) are inference; the Lead uses 'seconds-scale cutover' for the ideal outcome
- [ ] Write overhead from triggers is inferred from 'every insert, update and delete is recorded in per-table delta and track tables'; no Source measures it
- [ ] Heroku feature table: Managed Logical Replication 'Not Supported' on Classic Postgres; the fetch summary could not see the column headers
- [ ] Lead anchors (#copy-restarting-from-zero, #silent-network-drops, #the-delta_check-stall, #why-leaving-heroku-is-unique, #the-heroku-migrator, #the-migration-path-step-by-step, #before-you-migrate-the-discovery-tool-and-liftoff)
- [ ] Exact error string 'index row size 3328 exceeds btree version 4 maximum 2704' and the 'about 30 hours' detail in Silent network drops
- [ ] Every Lead quote: the Lead was read only through a fetch summarizer, not raw text
- [ ] Inference that the trigger path has no fallback like replica identity FULL ('The Lead describes no fallback')
- [ ] pg_dump parallel-dump lock note: a queued exclusive lock request can make the dump fail

**Fact-checker could not confirm:**
- [ ] All quotes verified via WebFetch summarizer (exact-phrase yes/no plus surrounding sentence), not raw page text; Bash/curl was denied
- [ ] 'Dump and restore needs a connection that can read every table' has no direct Source link (basic, unlinked)
- [ ] Downtime section 'keeps writes stopped from the moment pg_dump starts...' is an unlabelled inference following the dump section's 'we infer'

<details><summary>What the fact-checker changed</summary>

- Fixed: Quote 'the only change-capture mechanism available on a standard Heroku plan' relinked to #the-heroku-migrator (was #why-leaving-heroku-is-unique); 'tends to be more fragile' given its own #why-leaving-heroku-is-unique link
- Fixed: Delta_check stall: 'changes piled up into the millions' belongs to Silent network drops; replaced with Lead's 'backlog of unreplicated changes on the source kept growing'
- Fixed: Silent network drops: 30-hour Bucardo death and stalled worker were different migrations; added 'on another migration'
- Fixed: 'clean, seconds-scale cutover' describes most online logical-replication migrations, not an ideal; reworded and anchored to #why-leaving-heroku-is-unique
- Fixed: Replica identity FULL: manual says subscriber search 'can be very inefficient' without suitable index; fixed quote; table 'slow fallback' -> 'fallback'
- Fixed: 'Bucardo installs triggers' -> 'the migrator installs a bucardo_delta trigger on every source table' per Lead
- Fixed: Unlinked 'significantly higher performance over the serial mode' quote linked to populate.html#POPULATE-PG-DUMP (14.4.9); '--jobs' -> '-j'
- Fixed: Liftoff assessment claim relinked to Lead #before-you-migrate-the-discovery-tool-and-liftoff (Liftoff page does not describe it)
- Fixed: TL;DR: Heroku documents logical replication (not wal_level) as Advanced-only; reworded
- Fixed: max_wal_senders rule: added 'plus any physical replicas' per manual
- Fixed: Intro: 'next run emptied every table and started over' -> 'would empty ... start over' (Lead describes behaviour)

</details>

### Not drafted

- none

<details><summary>Leads the editor rejected</summary>

- `p27jcgwv`: Brief: GitHub PR merge-timeout status page states no cause, RCA promised later.
- `b2t2dl4z`: Brief: Copilot degradation attributed only to an upstream model provider, no mechanism.
- `hxhbrezm`: Brief: Jira email failures attributed only to a Microsoft 365 degradation, no mechanism.
- `emdwjqhx`: Brief: AWS Health Dashboard returned no text; no stated cause for the US-EAST-1 connectivity event.
- `dyqtqhih`: Brief: duplicate AWS US-EAST-1 status update; no stated cause.
- `ngptcql4`: Brief: duplicate AWS US-EAST-1 status update; no stated cause.
- `67rd2uwf`: Brief: duplicate AWS US-EAST-1 status update; no stated cause.
- `juslopgo`: Brief: duplicate AWS US-EAST-1 status update; no stated cause.
- `e5jzkuy5`: Brief: duplicate AWS US-EAST-1 status update; no stated cause.
- `hxstaai6`: Kitesurf post is mainly an open-source announcement with a short overview; no live Explainer teaches its lesson, too thin for a Case Study.
- `bbrn7vi5`: Newscast show notes with no transcript; ZGC adaptive heap sizing already covered by zgc-adaptive-heap-sizing.
- `s2q2xzby`: Conference topic roundup (DNS OARC 47), not an event at a named system.
- `l5laiihp`: Product release announcement (Datasentinel), no design change.
- `2mbhwtqc`: Product announcement (Cloudflare Spend Management).
- `mvap46l3`: Troubleshooting how-to for an Aurora error, not an incident; no live Explainer teaches it.
- `ucr5ttmb`: RDS Proxy pinning how-to, not an event; databases-storage already has today's pick.
- `u6yodzgr`: RDS Proxy monitoring how-to, not an event.
- `idkm25dl`: Product feature announcement (Antigravity /teamwork).
- `ayik5ghc`: AWS customer story; no live Explainer for a Brief; ai-ml Vertical already thick.
- `frjwk2ei`: Security-model opinion piece, not an event with a mechanism.
- `gfaluotq`: Strong Case Study (Zillow on Step Functions Express) but overlaps today's durable-execution pick in the same Vertical; hold for a later day.
- `uphw4t5j`: Genuine four-way ElastiCache isolation Comparison, but databases-storage already has today's pick; hold for a later day.

</details>


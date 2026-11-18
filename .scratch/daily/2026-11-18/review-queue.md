# Daily review queue, 2026-11-18

Drafted by `make daily` (ADR-0009). Every Piece is `draft: true` and unverified. First tick each draft's self-edit boxes while you read it, before you verify ([the checklist](../../../docs/editorial/self-edit-checklist.md)): they check the writing, not the facts, and the self-edit does not replace verification. Then verify every claim against its Source, then `make release` (ADR-0012): it sets `verified:`, removes `draft: true` and keeps each Publication Day. Ticked self-edit boxes never stand in for `verified:`.

## Run 20261118-062855

### Fly.io's connection exhaustion: many pools, no shared limit

`content/briefs/fly-io-database-connection-exhaustion.md` · brief · 222 words (recounted by the build's `brief-length` counter after the note was removed, 2026-11-18) · fact-checker: 19 claims checked, 7 fixed, 0 removed

**Publish within 72 hours of the Lead's date**, with or after its Explainer: cloudflare-hyperdrive-mysql-caching

**Drafter asks you to check:**
- [ ] DATE: the Lead is dated 2026-11-16 in the queue, but the post describes an incident on 'October 27' at 13:12 UTC with no year, no end time and no publication date on the page. The Brief says 'On 27 October at 13:12 UTC' with no year. Confirm the year is 2026 and that the 3-day rule holds counted from the post's first publication, not the incident.
- [ ] Source was read only through WebFetch's summariser, not the raw page. Check every Fly.io fact against the page itself: https://fly.io/infra-log/database-connection-exhaustion/
- [ ] Fact: 'the RDS database behind Fly.io's main web app ran out of connection slots' (Source wording reported as 'Our main web app's RDS database ran out of connection slots').
- [ ] Fact: 'Some dashboard and API requests returned 500 and 503 errors' (Source: requests 'intermittently failed with 500s and 503s').
- [ ] Fact: 'existing Machines and apps were unaffected'.
- [ ] Fact: 'Many Rails processes each kept their own connection pool' (Source: runs 'a lot of Rails processes to absorb request spikes', each with its own pool).
- [ ] Number: background worker machines held 'around 2,000 connections, most of the total' (Source: 'around 2,000 connections, making up the large majority of open connections').
- [ ] Number: moving container registry authentication to a read replica 'freed a little over 10% of the primary's connections' (Source: those services 'had been holding a little over 10% of connections on the primary'; the move was 'enough to stop the immediate failures', which the Brief leaves out).
- [ ] Fact: 'Fly.io plans to set safe limits on its workers' (Source: plans to work out safe limits for job worker machines, threads and pool sizes, enforced with queue throttling).
- [ ] Fact: 'Fly.io names both as options' for PgBouncer and RDS Proxy (Source: 'looking into putting a connection pooler such as PgBouncer or RDS Proxy in front of the database'; they are examples, not a shortlist).
- [ ] INFERENCE (labelled 'We infer'): the database's own cap was the only limit all callers shared. The Source does not say this; it rests on the planned limits and the pooler not yet being in place.
- [ ] ANALYSIS, not labelled as inference: 'That cap refuses whoever asks next, a user's request or a periodic job.' The Source says new connections were refused and dashboard/API requests failed; the generalisation is ours. Add 'we infer' if you want it marked.
- [ ] ANALYSIS: 'Total demand is processes times pool size, so adding processes for a spike uses more of the primary's slots' is our arithmetic, not the Source's.
- [ ] PgBouncer claim: 'In transaction mode, PgBouncer assigns a database connection for one transaction only' against https://www.pgbouncer.org/features.html ('A server connection is assigned to a client only during a transaction').
- [ ] RDS Proxy claim: 'RDS Proxy queues or throttles connections it cannot serve at once' against https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html (read in full: 'RDS Proxy queues or throttles application connections that can't be served immediately from the connection pool'; it also rejects connections past the limits you set).
- [ ] Hyperdrive claim: 'a SET ends with its transaction' against https://developers.cloudflare.com/hyperdrive/concepts/how-hyperdrive-works/#pooling-mode ('supports SET statements for the duration of a transaction or a query'; the Brief drops 'or a query').
- [ ] The Source never names the database engine. PgBouncer is Postgres-only and the linked Explainer covers Hyperdrive's MySQL support; the Brief claims neither engine for Fly.io. Check you are happy with the link framed as 'shows the pattern'.
- [ ] Left out for length: the recurrence about a day later (spike of 503s), staggering periodic jobs off the top of the hour, and deploy overlap raising connection counts. The term 'backpressure' appears only as a tag, not in the body.
- [ ] No quote is used. The sources: frontmatter is left commented out, as in the published examples; the PgBouncer, AWS and Cloudflare pages are linked inline only.

**Fact-checker could not confirm:**
- [ ] Incident year: the post, its .md twin and the index show no year. 2026 rests only on the feed's <published> 2026-11-16. The Brief states no incident year.
- [ ] 3-day rule: counted from first publication (feed 2026-11-16T00:00Z, day granularity), published 2026-11-18T06:32:27Z is 54.5 h, inside 72 h. It breaks if make daily moves the Brief to a Gap Day on or after 2026-11-19T00:00Z. Author to check before release.
- [ ] Raw page not read: curl was denied in this session, so Fly.io facts were confirmed through WebFetch with targeted verbatim-quote prompts (quotes capped near 125 characters), on both the HTML page and its .md twin. Each fact below matched a returned quote; full paragraphs were not seen in one piece.
- [ ] 'That cap refuses whoever asks next, a user's request or a periodic job' stays unlabelled analysis. Source supports the user half ('When new connections were refused, requests to the dashboard and API intermittently failed'; errors surface 'even if what caused us to reach the limit is something non-critical'). It does not say a periodic job was refused.
- [ ] 'Peak demand is processes times pool size, so adding processes for a spike uses more of the primary's slots' is our arithmetic. Source only says it runs 'a lot of Rails processes to absorb request spikes' and 'each one keeps its own pool'.
- [ ] Title 'many pools, no shared limit' sits loosely against the Source, which says the database itself is the one limit. Left unchanged as an editorial call for the Author.
- [ ] Engine: the Source never names Postgres or MySQL. PgBouncer is Postgres-only and the linked Explainer covers Hyperdrive's MySQL support. The Brief claims no engine for Fly.io; 'shows the pattern' framing left as is.
- [ ] 'around 2,000 connections, most of the total': Source says 'around 2,000 connections' across all background worker machines, 'making up the large majority of open connections', with no total and no time given for the count.
- [ ] No end time or duration is on the page; the Brief gives none. Left out by the drafter and still out: the move to a read replica 'was enough to stop the immediate failures', the 503 spike 'about a day later', staggered periodic jobs, fewer worker machines, deploy overlap.

**Self-edit before you verify** ([checklist](../../../docs/editorial/self-edit-checklist.md)), ticked by you while you read the draft:
- [ ] Point: I can say the point in one sentence: ______________________________
- [ ] Point: The piece covers one idea; anything else is a Related link.
- [ ] Top of the page: Title: key term in the first two words, plain words, no pun, under ~65 chars.
- [ ] Top of the page: Description: adds the "so what", has a verb, doesn't repeat the title, <=150 chars.
- [ ] Top of the page: Nothing comes before "What happened".
- [ ] Structure: Read only the headings: do they tell the story in order?
- [ ] Structure: No heading says "Background", "Overview", "Introduction" or "Conclusion".
- [ ] Structure: No question headings outside a Repo Profile FAQ.
- [ ] Structure: Each paragraph has one topic, stated in its first sentence.
- [ ] Structure: The whole Brief is 250 words at most, headings included (`brief-length`); cut, don't pad.
- [ ] Sentences: Average about 20 words or fewer; every sentence over ~30 words is split or earns its length.
- [ ] Sentences: Active voice with a named actor ("the client retried", not "requests were retried").
- [ ] Sentences: Hardest parts use the shortest sentences.
- [ ] Words: Every technical term is defined at first use, and kept where it is the right term.
- [ ] Words: One name per thing: the Source's term every time; a synonym appears once, in brackets, only if Readers search for it.
- [ ] Words: No hype or filler word from voice rule 5's banned list (`house-byline-voice.md`), and no "there is/are" openings.
- [ ] Words: At most one analogy, and it says where it stops being true.
- [ ] Words: No idioms or pop-culture references a non-native reader would miss.
- [ ] Teaching: A concrete example comes before the abstract statement.
- [ ] Teaching: Every "what" has a "why" or a cost.
- [ ] Teaching: The common confusion is named and corrected ("what it is not").
- [ ] Teaching: Code: introduced by a sentence, realistic, under its licence; README commands verbatim.
- [ ] Evidence: Every claim has its link where the claim is made.
- [ ] Evidence: Every number has a unit, a date or version, what it is compared with, and its caveat.
- [ ] Evidence: Quotes are exact, short, credited, and followed by our own analysis.
- [ ] Evidence: Each inference says "we infer" in its own sentence.
- [ ] Our contribution: Name what this Piece adds that no Source says: ____________________ ("Why it matters" is the contribution: our own analysis, never the Source restated.)
- [ ] Our contribution: Delete every paragraph that only restates a Source.
- [ ] Ending: The last section tells the Reader what they can now do, measure or decide.
- [ ] Ending: No new idea in the last paragraph; no recap of what they just read.
- [ ] Last pass: Read it aloud once. Fix every place you stumbled or got bored.
- [ ] Last pass: Cut big (sections, paragraphs), then small (words). Then cut 10% more.

<details><summary>What the fact-checker changed</summary>

- Fixed: Mislabelled inference: 'We infer the database's own cap was the only limit all callers shared' was not an inference. Source states it ('Currently, the only thing that enforces any limit is the database itself'). Now '[Fly.io says](source) only the database itself enforces any limit', linked; 'We infer' dropped.
- Fixed: Date: 'On 27 October at 13:12 UTC ... its infra log reports' replaced with 'Fly.io's [infra log entry of 16 November 2026](feed.xml), headed [27 October, 13:12 UTC](post), says ...'. Page heading reads 'October 27: ... (13:12UTC)' with no year; Atom feed gives <published> 2026-11-16T00:00:00+00:00 (updated 2026-11-17T14:21:47). Same pattern as the Author-verified Petsem Brief.
- Fixed: Hyperdrive: 'a SET ends with its transaction' -> 'a SET lasts only for its transaction or query' (Cloudflare: 'supports SET statements for the duration of a transaction or a query').
- Fixed: PgBouncer/RDS Proxy: 'Fly.io names both as options' -> 'Fly.io, looking into a pooler, names both' (Source: 'looking into putting a connection pooler such as PgBouncer or RDS Proxy in front of the database'; examples, not a shortlist).
- Fixed: Worker limits: 'plans to set safe limits on its workers' -> 'plans safe limits for its workers' (Source: 'plan to work out safe limits for job worker machines, threads, and pool sizes'); sentence moved after 'That cap refuses...' so 'That cap' follows the database-limit sentence it refers to.
- Fixed: Arithmetic: 'Total demand is processes times pool size' -> 'Peak demand ...' (processes x pool size is the ceiling, not the standing total).
- Fixed: Description: 'dashboard and API calls failed' -> 'some dashboard and API calls failed' (Source: 'intermittently failed'); 145 characters.

</details>

### Deterministic simulation testing: replay any failure from a seed

`content/explainers/distributed-systems/deterministic-simulation-testing.md` · concept · 2300 words · fact-checker: 96 claims checked, 15 fixed, 1 removed

**Drafter asks you to check:**
- [ ] All Sources were read through WebFetch, which returns a small model's extraction, not the raw page. Every quotation in the draft needs a character-for-character check against its page; several were returned as clauses of longer sentences.
- [ ] Hook: FoundationDB engineering page wording for 'exclusively in simulation for the first 18 months' (I paraphrase it as 'its first 18 months'); confirm the #simulation anchor on both engineering.html and testing.html.
- [ ] Hook: 'roughly one trillion CPU-hours of simulation' is from testing.html; engineering.html says 'a trillion CPU-hours of simulated stress testing'. The docs are version 8.0.1 and the estimate is undated.
- [ ] Wilson post ('Is something bugging you?', 13 Feb 2024): the 2010 start, the quoted fragment 'we first wrote a fully-deterministic event-based network simulation', 'driven by the same random number generator', 'one or two bugs reported by a customer', Antithesis started in 2018, and the fragment 'write a hypervisor which emulates a deterministic computer'.
- [ ] Antithesis docs: the definition quote, 'FoundationDB and Amazon Web Services around 2010', 'generally impractical for systems already in production', 'run many times with different seeds'. The claim that exploring the state space efficiently is as hard as making the system deterministic is a paraphrase I did not see verbatim. Anchors (#what-is-..., #how-do-i-implement-...) come from the page's table of contents as reported, not from the HTML.
- [ ] TigerBeetle Safety docs: '24/7 on 1024 cores' and 'at 1000x speed', and that both sit under #software-reliability.
- [ ] TigerBeetle VOPR page (linked at main, unpinned): stubs clock/network/disk, may corrupt reads and writes to the 'disk', replay from seed plus Git commit hash, credits FoundationDB and Antithesis.
- [ ] TigerBeetle liveness post (matklad, 6 Jul 2023): quote 'does not check all liveness properties'; my paraphrase 'partitions always healed in that mode, so the team added a second one'; my gloss of liveness as 'the cluster keeps making progress'.
- [ ] Lead numbers: 10-tick timeout, 20-tick delay, majority of two, fault probability 0.08, request probability 0.4, heal every 40 ticks, seed 123 for 500 ticks, TigerBeetle 10 ms tick, etcd 100 ms tick with 10-tick election timeout, 40 ms vs 120 ms.
- [ ] Lead: the six rule names and their enforcement (ArchUnit suite for the first five, runtime thread check in the future for the sixth); the quoted fragment 'must follow six strict coding rules'; all #anchors (taken from the article's table of contents).
- [ ] Lead sidebar: FoundationDB simulates at byte-stream level, TigerBeetle's packet simulator at message level, Tickloom framing code runs only on the real network; 'the message-level kind is simpler' is my paraphrase.
- [ ] etcd-io/raft README (linked at main, unpinned): both quotes, and the paraphrase that time is an abstract tick advanced by Node.Tick().
- [ ] Jepsen linearizability page: both quotes.
- [ ] Antithesis hypervisor post (Alex Pshenichkin, 20 Mar 2024): 'runs on just one physical CPU core' and the #cpu-parallelism anchor.
- [ ] ArchUnit described as 'a library for checking Java architecture rules': my wording, archunit.org was not fetched.
- [ ] Three 'we infer' lines are our own: the resemblance to durable execution, 'DST is not a proof', and the application-level trade (no real kernel, disk or network stack). Also our own: 'tick-based timeouts alone are not DST' and the 'where is the seam' advice in the last section.
- [ ] Diagram: step order (inject or heal faults, advance clocks, tick network, clients, servers, maybe issue a client request) is from the Lead's Cluster.tick() and SimulationRunner description. The Lead does not say where the simulated disk is ticked, so the diagram leaves the disk out.
- [ ] Closing: 'clone Tickloom and change the seed in its fault-injecting test' was not run by us; sim.tigerbeetle.com was not opened.
- [ ] Internal links to check on the Preview: /distributed-systems/etcd/, /distributed-systems/durable-execution/, /infrastructure-cloud/aws-resilience-hub-testing/.
- [ ] The FoundationDB SIGMOD 2021 paper could not be read (PDF tooling missing, ACM returned 403), so it is not cited; its stated simulation limits would strengthen 'Practice two' if the Author wants them.
- [ ] Description is 136 characters by my hand count; title 64.

**Fact-checker could not confirm:**
- [ ] All checks ran through WebFetch, which caps quotations at 125 characters and returns a small model's reading. Raw-HTML download (curl) was denied in this session, so every quote was confirmed as EXACT by that model, not by a byte-level comparison I ran myself. The Author should spot-check the longer quotes on the pages.
- [ ] Anchors (#MakingProgressExplicit, #what-is-deterministic-simulation-testing, #software-reliability, #cpu-parallelism, #simulation, #correctness-and-performance and the rest) were confirmed from each page's table of contents as WebFetch reported it, not from the HTML id attributes. Check them by clicking on the Preview.
- [ ] sim.tigerbeetle.com could not be read (the page is JavaScript only; WebFetch saw a title and one image). 'Inject faults by hand' rests on TigerBeetle's vopr.md: 'you can inject different faults yourself'.
- [ ] Tickloom was not cloned or run. The README confirms seeded simulation tests (SimulationRunnerTest, seeds 42L, 123L, 111_111L); whether the Lead's faultInjectedHistoryShouldBeCheckable test exists in the repository under that name was not checked.
- [ ] The Lead's Cluster.tick() does not show where the simulated disk is ticked. The diagram leaves the disk out, which matches the Lead, but the sentence 'Disk operations wait for ticks too' rests on SimulatedFileIO.tick() alone.
- [ ] 'Wilson co-founded Antithesis in 2018' combines the byline ('Co-founder and CEO') with 'So in 2018, we started Antithesis'. The post never states it in one sentence.
- [ ] Wilson's post says 'When we sat down in 2010'. Whether Wilson himself was at FoundationDB in 2010 is not stated in the post; the draft attributes the account to him and says 'the FoundationDB team'.
- [ ] 'Roughly one trillion CPU-hours of simulation' is from the FoundationDB 8.0.1 docs and is undated; 'has since run' is the draft's framing. engineering.html words it as 'the equivalent of a trillion CPU-hours of simulated stress testing'.
- [ ] tigerbeetle vopr.md and the etcd-io/raft README are linked at main, unpinned; their wording can change after today.
- [ ] Our own glosses, not from a Source and not marked 'we infer': gray failure as 'a component is degraded but not dead' (the Lead says 'components rarely fail cleanly'), 'the new log line appears in the same run', and the 'where is the seam' advice in the last section.
- [ ] The durable-execution resemblance and the description of AWS Resilience Hub rest on HastCode's own Pieces, which I did not re-verify. All three internal links resolve to non-draft Pieces under content/.
- [ ] The FoundationDB SIGMOD 2021 paper remains unread and uncited (drafter's note).

**Self-edit before you verify** ([checklist](../../../docs/editorial/self-edit-checklist.md)), ticked by you while you read the draft:
- [ ] Point: I can say the point in one sentence: ______________________________
- [ ] Point: The piece covers one idea; anything else is a Related link.
- [ ] Top of the page: Title: key term in the first two words, plain words, no pun, under ~65 chars.
- [ ] Top of the page: Description: adds the "so what", has a verb, doesn't repeat the title, <=150 chars.
- [ ] Top of the page: Hook: one real, linked moment (system, number, failure), not a definition or history.
- [ ] Top of the page: Nut graf: the hook's last sentence says what this piece gives the Reader.
- [ ] Top of the page: TL;DR: the answer in plain words; a reader who stops here is not misled.
- [ ] Structure: Read only the headings: do they tell the story in order?
- [ ] Structure: No heading says "Background", "Overview", "Introduction" or "Conclusion".
- [ ] Structure: No question headings outside a Repo Profile FAQ.
- [ ] Structure: Each paragraph has one topic, stated in its first sentence.
- [ ] Structure: Each H2 is under 250 words; cut, don't pad.
- [ ] Sentences: Average about 20 words or fewer; every sentence over ~30 words is split or earns its length.
- [ ] Sentences: Active voice with a named actor ("the client retried", not "requests were retried").
- [ ] Sentences: Hardest parts use the shortest sentences.
- [ ] Words: Every technical term is defined at first use, and kept where it is the right term.
- [ ] Words: One name per thing: the Source's term every time; a synonym appears once, in brackets, only if Readers search for it.
- [ ] Words: No hype or filler word from voice rule 5's banned list (`house-byline-voice.md`), and no "there is/are" openings.
- [ ] Words: At most one analogy, and it says where it stops being true.
- [ ] Words: No idioms or pop-culture references a non-native reader would miss.
- [ ] Teaching: A concrete example comes before the abstract statement.
- [ ] Teaching: At least one worked example with real numbers (Concept, Case Study, Comparison).
- [ ] Teaching: Every "what" has a "why" or a cost.
- [ ] Teaching: The common confusion is named and corrected ("what it is not").
- [ ] Teaching: The text walks the diagram box by box; the diagram is never left to speak alone.
- [ ] Teaching: Code: introduced by a sentence, realistic, under its licence; README commands verbatim.
- [ ] Evidence: Every claim has its link where the claim is made.
- [ ] Evidence: Every number has a unit, a date or version, what it is compared with, and its caveat.
- [ ] Evidence: Quotes are exact, short, credited, and followed by our own analysis.
- [ ] Evidence: Each inference says "we infer" in its own sentence.
- [ ] Our contribution: Name what this Piece adds that no Source says: ____________________
- [ ] Our contribution: Delete every paragraph that only restates a Source.
- [ ] Ending: The last section tells the Reader what they can now do, measure or decide.
- [ ] Ending: No new idea in the last paragraph; no recap of what they just read.
- [ ] Last pass: Read it aloud once. Fix every place you stumbled or got bored.
- [ ] Last pass: Cut big (sections, paragraphs), then small (words). Then cut 10% more.

<details><summary>What the fact-checker changed</summary>

- Fixed: Hook: 'co-founder Will Wilson' read as FoundationDB co-founder; the post's byline says only 'Co-founder and CEO' and calls FoundationDB 'our last company'. Now 'Antithesis co-founder Will Wilson'.
- Fixed: Hook: 'A failing run could be run again, and it failed the same way' was unlinked. Linked to Wilson's post ('always perfectly reproducible', 'run it over and over again with the same random seed') and reworded to 'run again with the same seed'.
- Fixed: Worked example: 'Byzantium's acknowledgement is due 20 ticks out' was imprecise. The Lead says it leaves 'a few ticks later, and arrives 20 ticks after that'. Now 'spends 20 ticks in transit'.
- Fixed: Wrong anchor: 'could go either way' pointed at #TheNondeterminismProblem; the Run 1 / Run 2 timeout passage sits in #SourcesOfNondeterminism. Anchor corrected.
- Fixed: Disk section: 'TigerBeetle's simulator goes a step further' implied Joshi's simulated disk does not corrupt data, which the Lead does not state. Now 'faults the disk as well'.
- Fixed: Checker section: 'Following Jepsen's format' changed to 'Jepsen's approach'. The Lead says 'The approach used by Jepsen is to record client operations and their results in EDN format' and does not say Tickloom uses that format.
- Fixed: Jepsen quote 2 was presented as a standalone sentence with a capital 'If'. On the page it is a lowercase clause after 'e.g.,'. Now lowercase, introduced with a colon, full stop outside the quotation.
- Fixed: ArchUnit quote ended with a full stop inside the quotation marks, but the Lead's sentence continues after an em dash ('at runtime—it cannot distinguish…'). Full stop moved outside the quotation.
- Fixed: Practice one: 'The cost falls on dependencies… has to stay behind the seam' was stated as the Lead's claim; the Lead never mentions third-party libraries. Now marked 'We infer', with the link kept on 'behind the seam' (#Trade-offs).
- Fixed: Practice two: 'Joshi argues this rarely hurts. Many distributed systems are I/O bound' overstated the Lead ('Most of the distributed systems are often IO bound'; single-threaded processing 'actually works better'). Now 'Joshi argues this suits I/O-bound work, which most distributed systems are.'
- Fixed: Liveness gloss 'the cluster keeps making progress' replaced with the post's own definition ('the cluster remains available if the operating environment is sufficiently healthy'): 'a healthy enough cluster stays available'.
- Fixed: Antithesis docs: link text 'run many times with different seeds' corrected to 'run multiple times' (page: 'running tests multiple times with different seeds').
- Fixed: Antithesis docs: 'exploring the state space efficiently is as hard as making the system deterministic' is not what the page says (it calls it 'a complex undertaking as well', beside the 'main technical challenge'). Now 'is hard in its own right'.
- Fixed: What it is not: 'tick-based timeouts alone are not DST' is our own classification; now reads 'We also infer that…'.
- Fixed: Closing: 'change the seed in its fault-injecting test' narrowed to 'its simulation tests'. The Tickloom README confirms seeded SimulationRunner tests; it does not say the seeded runner injects faults.
- Removed: 'and recovery time' from the chaos-engineering sentence: the Lead says chaos engineering monitors infrastructure health and availability, and does not mention recovery time.

</details>

### kube-proxy iptables vs ipvs vs nftables: which backend to run

`content/explainers/infrastructure-cloud/kube-proxy-iptables-vs-ipvs-vs-nftables.md` · comparison · 2350 words · fact-checker: 78 claims checked, 8 fixed, 1 removed

**Drafter asks you to check:**
- [ ] Sources were read through WebFetch (a small model relays the page), not raw HTML: check every quoted phrase character for character, especially the Lead quotes ('use 100% of one CPU all the time just to process iptables updates', 'do not share their IPVS scheduler state', 'generally don't work better than random selection', 'still not as good as the `ipvs` backend', 'even the largest numbers in this graph are quite tiny').
- [ ] Removal timeline: the Lead says 'plan to remove it in Kubernetes 1.40'; the docs and KEP-5495 say off by default in 1.40 behind the KubeProxyIPVS feature gate and code deleted in 1.43. The piece states both; confirm the docs wording ('disabled by default from Kubernetes v1.40 ... fully removed in Kubernetes v1.43') and the KEP stages.
- [ ] Latency comparisons come from the 2025 post, not the Lead: nftables p50 about equal to iptables p01 at 5,000 and 10,000 Services; nftables p99 beats iptables p01 'by a few microseconds' at 30,000 Services. No absolute figures are given in text by any Source, so none are in the piece; KEP-5495 has no latency figures at all (the editor's angle expected some).
- [ ] KEP-3866 claims were relayed as fragments: 'acquire a lock and download the entire ruleset from the kernel' then re-upload (check the full sentence and what operation it describes), 'do not provide enough functionality to fully implement Kubernetes services', each nftables user has its own table, 'multiple commands ... in a single atomic transaction'. Check the GitHub heading anchors resolve (#the-iptables-kernel-subsystem-has-unfixable-performance-problems, #the-ipvs-mode-of-kube-proxy-will-not-save-us, #tables, #communicating-with-the-kernel-nftables-subsystem).
- [ ] KEP-5495 anchors and wording: #motivation (SIG Network 'lacks maintainers who are familar with the ipvs backend code'), #proposal ('will be out of LTS by the end of 2026'), #graduation-criteria.
- [ ] 'Eleven schedulers' is our count of the docs' list (rr, wrr, lc, wlc, lblc, lblcr, sh, dh, sed, nq, mh); recount on the page.
- [ ] 'Since Kubernetes v1.28' for iptables partial updates is the docs' wording; the 2025 post says the work began in 1.26. The piece uses v1.28 only.
- [ ] Our own glosses, not from a Source: EndpointSlice ('the object that lists a Service's endpoints'), ipset ('kernel sets of addresses those rules can match'), destination NAT ('rewrites the packet's destination address'), conntrack ('the kernel's connection tracking'), LTS ('long-term support'), 'IPVS is a kernel load-balancing subsystem', verdict map ('a table from a key to an action').
- [ ] Inferences labelled 'we infer': iptables as the behaviour the others are measured against; a least-connection scheduler counting only its own node's connections; the whole 'When to run' section (including treating a high sync_proxy_rules_duration_seconds as the signal to move to nftables).
- [ ] Diagram: chain and map names (KUBE-SERVICES, services, service-ips) come from the 2025 post's sample rulesets; the 10.0.0.1:1234 Service is the docs' example. The ipvs and nftables last boxes say 'Traffic redirected to the endpoint' because the docs state destination NAT explicitly only for iptables mode.
- [ ] Docs anchors used: #proxy-modes, #proxy-mode-iptables, #packet-processing-iptables, #minsyncperiod, #minimize-iptables-restore, #proxy-mode-ipvs, #proxy-mode-nftables, #migrating-from-iptables-mode-to-nftables. Lead anchors: #why-was-the-ipvs-backend-written, #why-is-the-ipvs-backend-now-being-deprecated, #what-about-ipvs-schedulers, #switching-to-nftables, #switching-to-iptables.
- [ ] 'Default in 1.37' and 'alpha in Kubernetes 1.37' (localhost NodePorts) are tied to the docs version live on 2026-11-18; recheck if the docs have moved to a newer release.
- [ ] Not included because only the February 2025 post states it and it may be dated: nft command-line tool 1.0.0 or later on the host. Add if the Author confirms it still applies.
- [ ] Table cell 'Not measured in the Sources' (nftables, many changes at once) rests on the 2025 post's 'Unfortunately I don't have cool graphs for this part' and on the Lead giving no control-plane measurement.
- [ ] related: [cilium] added to frontmatter (live Repo Profile on replacing kube-proxy); remove if the Author prefers no Related link.

**Fact-checker could not confirm:**
- [ ] All Sources were read through WebFetch (a small model relays the page); Bash/curl for raw HTML was denied. Every quoted phrase was returned verbatim by the relay, but straight-vs-curly apostrophes cannot be confirmed (e.g. 'generally don't work better than random selection', 'Unfortunately I don't have cool graphs for this part'). Author should eyeball those quotes on the live pages.
- [ ] GitHub KEP anchors (#the-iptables-kernel-subsystem-has-unfixable-performance-problems, #the-ipvs-mode-of-kube-proxy-will-not-save-us, #tables, #communicating-with-the-kernel-nftables-subsystem, #motivation, #proposal, #graduation-criteria) were confirmed from the exact heading text in the raw READMEs, not by loading the rendered anchors.
- [ ] Docs anchors (#proxy-mode-iptables, #packet-processing-iptables, #minimize-iptables-restore, #proxy-mode-ipvs, #proxy-mode-nftables) were confirmed from kubernetes/website main Markdown source, not the rendered v1.37 HTML; #proxy-modes, #minsyncperiod and #migrating-from-iptables-mode-to-nftables are auto-generated ids inferred from heading text. 2025 post anchors (#why-nftables-part-1-data-plane-latency, #why-nftables-part-2-control-plane-latency, #why-not-nftables, #future-plans, #learn-more) are likewise inferred from heading text. Lead anchors were all seen on the page, including #switching-to-iptables.
- [ ] 'Default in 1.37' and 'alpha in Kubernetes 1.37' confirmed against the docs live on 2026-11-18 (page text names Kubernetes 1.37; feature state 'Alpha since Kubernetes v1.37; disabled by default'); will date when the docs move to a newer release.
- [ ] TL;DR 'ipvs one call per change' is our reading of the Lead's 'doesn't have a way of doing batch updates'; no Source uses the words 'one call'. Left as is.
- [ ] Diagram boxes 'pick an endpoint at random' (iptables, nftables) omit session affinity, which the docs name as the alternative for iptables; the prose states it. Left as is.
- [ ] The post calls p50 'the average' (the piece follows the post); p50 is conventionally the median. Left as the Source words it.
- [ ] nft command-line tool 1.0.0 or later: confirmed present in the February 2025 post ('earlier than 1.0.0'), still left out per the drafter's note; Author's call.
- [ ] KEP-5495 contains no latency figures (confirmed), so the editor's expected angle is not met from that Source; the piece gives none.
- [ ] related: [cilium] resolves to content/explainers/apis-networking/cilium.md (repo-profile, no draft flag); keeping it is the Author's call.

**Self-edit before you verify** ([checklist](../../../docs/editorial/self-edit-checklist.md)), ticked by you while you read the draft:
- [ ] Point: I can say the point in one sentence: ______________________________
- [ ] Point: The piece covers one idea; anything else is a Related link.
- [ ] Top of the page: Title: key term in the first two words, plain words, no pun, under ~65 chars.
- [ ] Top of the page: Description: adds the "so what", has a verb, doesn't repeat the title, <=150 chars.
- [ ] Top of the page: Hook: one real, linked moment (system, number, failure), not a definition or history.
- [ ] Top of the page: Nut graf: the hook's last sentence says what this piece gives the Reader. (the opening's closing question, the one the piece answers, does this job.)
- [ ] Top of the page: TL;DR: the answer in plain words; a reader who stops here is not misled.
- [ ] Structure: Read only the headings: do they tell the story in order?
- [ ] Structure: No heading says "Background", "Overview", "Introduction" or "Conclusion".
- [ ] Structure: No question headings outside a Repo Profile FAQ.
- [ ] Structure: Each paragraph has one topic, stated in its first sentence.
- [ ] Structure: Each H2 is under 250 words; cut, don't pad.
- [ ] Sentences: Average about 20 words or fewer; every sentence over ~30 words is split or earns its length.
- [ ] Sentences: Active voice with a named actor ("the client retried", not "requests were retried").
- [ ] Sentences: Hardest parts use the shortest sentences.
- [ ] Words: Every technical term is defined at first use, and kept where it is the right term.
- [ ] Words: One name per thing: the Source's term every time; a synonym appears once, in brackets, only if Readers search for it.
- [ ] Words: No hype or filler word from voice rule 5's banned list (`house-byline-voice.md`), and no "there is/are" openings.
- [ ] Words: At most one analogy, and it says where it stops being true.
- [ ] Words: No idioms or pop-culture references a non-native reader would miss.
- [ ] Teaching: A concrete example comes before the abstract statement.
- [ ] Teaching: At least one worked example with real numbers (Concept, Case Study, Comparison).
- [ ] Teaching: Every "what" has a "why" or a cost.
- [ ] Teaching: The common confusion is named and corrected ("what it is not").
- [ ] Teaching: The text walks the diagram box by box; the diagram is never left to speak alone.
- [ ] Teaching: Code: introduced by a sentence, realistic, under its licence; README commands verbatim.
- [ ] Evidence: Every claim has its link where the claim is made.
- [ ] Evidence: Every number has a unit, a date or version, what it is compared with, and its caveat.
- [ ] Evidence: Quotes are exact, short, credited, and followed by our own analysis.
- [ ] Evidence: Each inference says "we infer" in its own sentence.
- [ ] Our contribution: Name what this Piece adds that no Source says: ____________________
- [ ] Our contribution: Delete every paragraph that only restates a Source.
- [ ] Ending: The last section tells the Reader what they can now do, measure or decide.
- [ ] Ending: No new idea in the last paragraph; no recap of what they just read.
- [ ] Last pass: Read it aloud once. Fix every place you stumbled or got bored.
- [ ] Last pass: Cut big (sections, paragraphs), then small (words). Then cut 10% more.

<details><summary>What the fact-checker changed</summary>

- Fixed: Intro: 'Every node in a cluster runs kube-proxy' now carries the docs' caveat ('Unless you have replaced it, ...'); docs say 'unless you have deployed your own alternative component'.
- Fixed: TL;DR bullet 5: 'Kubernetes recommends nftables on Kubernetes 1.36 or later' narrowed to 'to `ipvs` users'; the Lead's recommendation is addressed to ipvs users on 1.36+ with kernel 5.13+.
- Fixed: iptables section: KEP-3866 lock/download sentence now states the operation it describes ('To add a single rule, KEP-3866 says, the iptables binary must acquire a lock and download the entire ruleset...'), matching the KEP's full sentence.
- Fixed: nftables section: verdict map action 'is a jump to' changed to 'sends the packet to' the Service's chain; the 2025 post's sample uses `goto`, not jump.
- Fixed: nftables section: 'arrived as alpha in Kubernetes 1.29' changed to 'arrived in Kubernetes 1.29'; the linked Lead says 'introduced', never 'alpha'.
- Fixed: Lookup section: 'The documentation agrees: the packet gain...' changed to 'The documentation says the packet gain over `iptables`...'; the docs sentence compares nftables with iptables mode, not with ipvs as the preceding Lead quote does.
- Fixed: Service features section: 'Everything beyond basic load balancing falls back to iptables' softened to 'Various features beyond basic load balancing fall back'; the Lead says 'various features of Kubernetes Services beyond basic "load balancing"'.
- Fixed: Service features section: the quote 'was never able to implement all of the edge cases of Kubernetes Service functionality correctly' now links the Lead (it had no link in a paragraph whose only link was a different Source).
- Removed: iptables section: 'Every component on the node shares it, too.' The linked KEP-3866 section says nothing about other components sharing iptables; the supported version (no global lock contention between components, 2025 post) is already stated and linked in the control-plane section.

</details>

### Campaign Orchestrator: how Zillow sends 100M notifications a day

`content/explainers/distributed-systems/zillow-notification-workflows.md` · case-study · 2650 words · fact-checker: 88 claims checked, 11 fixed, 0 removed

**Fills Gap Day 2026-11-17** (ADR-0012): `published` is that day; `make release` keeps it.

**Drafter asks you to check:**
- [ ] Reading method: WebFetch returns model-extracted text, not raw HTML, and Bash curl was denied. Every quote and number came through that extraction. Check each quotation mark against the page as published, especially: "about 90 percent cheaper", "would have added another infrastructure component to operate", "required custom lifecycle management, worker pools, queue integration, and monitoring", "at roughly a quarter of the estimated build effort", the two table verdicts, "Some of these controls are specific to how we order and throttle campaign sends, so they live outside Step Functions.", "Campaigns don't share one generic workflow that picks modules at runtime.", "stateless and schema-driven", "deliver, retry, or drop each message individually", "At our target throughput, service quotas are part of the architecture.", "fell along the same lines, from roughly 8 hours per week to 2".
- [ ] Tail of the post (from 'Batch mode optimization' onward: 1.2 s, 12 s, 1.4 s, 10 to 20 percent, 256 KiB, 64 MB, quotas, per-campaign queues and circuit breakers, TestState in CI, 9 vs 32 development weeks, 8 to 2 hours, author bios) was read from the AWS Architecture Blog RSS feed because the page fetch truncated. Confirm against the page itself.
- [ ] Author affiliations: Yuzhong Li 'founding member and tech lead of Zillow Group's Product Comms System'; Nishita Kheur 'Senior Technical Account Manager at AWS'. Only one of the two authors is Zillow's. Publication date 11 November 2026.
- [ ] Hook softened: the post says each notification 'runs as its own short-lived workflow execution', but also describes batching many messages into one execution. The draft says Zillow 'runs them as short-lived executions' and lists the tension as a gap. Check you agree with that reading.
- [ ] Diagram edge that is partly ours: results returning from the workflow to the router (dotted), and the router handing off to the delivery service. The post says the router uses per-message results to deliver, retry or drop, but its Figure 2 sentence only says 'the delivery service sends the result'. Compare with Figure 2.
- [ ] Diagram label 'one queue per campaign' and 'one Express Workflow per campaign' rest on 'Each campaign still runs its own state machine and its own Amazon SQS queue'. Confirm wording.
- [ ] '8 hours per week to 2': the post does not make clear whether 8 is the estimate for an in-house engine or the old system's measured load. The draft quotes it and calls both effort figures comparisons against an estimated in-house engine. Check that framing.
- [ ] The post names both CloudFormation ('the Step Functions definition and CloudFormation template that gets deployed') and Terraform ('writes the result into the Terraform we deploy'). The draft says only 'deployment template' and names neither. Decide whether to mention the inconsistency.
- [ ] Our arithmetic, labelled as ours: batch saving between 88% and 90% (requests 10 to 1; billed time 12 s to 1.4 s); 100 million a day is about 1,160 executions a second; times four HTTP calls is about 4,600 a second; 9/32 = 28% and 2/8 = 25%.
- [ ] AWS quota figures cited from the Step Functions service quotas page: HTTP Task token bucket 300 with refill 300 per second; Express StartExecution bucket 6,000 with refill 6,000 per second; HTTP Task duration 60 seconds (hard quota); 256 KiB maximum input or output. The post calls the first 'HTTP Invoke concurrency' while the docs express it as a rate; the draft uses both terms. Check anchors #service-limits-http-task, #service-limits-api-action-throttling-general, #service-limits-task-executions resolve.
- [ ] Express pricing wording ('by the number of requests and by duration', memory in 64-MB chunks) is from aws.amazon.com/step-functions/pricing/. Confirm.
- [ ] Express table rows from 'Choosing workflow type in Step Functions': one year vs five minutes; exactly-once vs at-least-once (asynchronous) and at-most-once (synchronous); billed by state transitions vs executions, duration and memory; workflow type immutable; Express execution history not captured by Step Functions. Anchor #express-at-least-once-execution.
- [ ] Inferences marked 'we infer' or 'our reading', to confirm you accept each: extracting modules first kept the engine's job narrow; Zillow needed none of Temporal's durable features; CloudWatch logging quota follows from Express not keeping history; the idempotency check guards against the same message reaching the router twice; per-campaign workflows fix steps at build time and limit a bad config to one campaign; duplicate sends are stopped by the router's check or the delivery service.
- [ ] Gaps stated as absent from the post: no latency, error or duplicate rates, no dollar cost, no production batch size, no async vs sync start mode, no detail of the idempotency key, no incident counts for the old design. The extraction reported 'not mentioned' for these; confirm on the page.
- [ ] Gloss of 'hydrate' ('load the details a message needs, such as a home's price and photos') extends the post's definition of property hydration to hydration in general. Circuit breaker and idempotent definitions are ours.
- [ ] Module list shortened: 'common data' is given without its contents (subscriptions, saved homes, eligibility flags, propensity scores). Criteria list omits 'priced to fit volume', which appeared in one extraction but not in the quoted sentence.
- [ ] Internal link /distributed-systems/durable-execution/ used twice and added as related: [durable-execution]. Confirm that Explainer is live and that the sentence 'Temporal is a durable execution runtime, built so a workflow can wait for days and survive a crash' matches what it teaches.
- [ ] Title uses '100M' to stay at 64 characters; the Source says 'more than 100 million'. Description is 141 characters.

**Fact-checker could not confirm:**
- [ ] Tail of the post (from 'Batch mode optimization' onward: 1.2 s, 12 s, 1.4 s, 10 to 20 percent, 'about 90 percent cheaper', 256 KiB, 64 MB, quotas paragraph, per-campaign queues and circuit breakers, TestState in CI, 9 vs 32 weeks, 8 to 2 hours, product velocity, ownership, author bios) was confirmed only against the AWS Architecture Blog RSS feed full text; the page fetch truncates at 'assuming each step supports' at any offset. Author should check these on the live page.
- [ ] Figure 2 image could not be viewed; the diagram was checked against the post's text only. Results returning to the router and the router handing off to the delivery service are inferred and now labelled so.
- [ ] 'In practice a batch of 10 takes about 1.4 seconds': the post says 'a batch runs only 10 to 20 percent longer than a single message, about 1.4 seconds instead of 12'; batch size 10 is implied by the next sentence, not stated. Left as is.
- [ ] Gap claims (no error or duplicate rates, no dollar cost, no production batch size, no sync vs async start mode, no idempotency key detail, no incident counts) rest on model extraction reporting 'not found' for both the page and the feed tail.
- [ ] Whether '8 hours per week' is the old system's measured load or the in-house engine estimate is not stated in the post; the draft now says so.
- [ ] The post names both CloudFormation ('the Step Functions definition and CloudFormation template that gets deployed') and Terraform ('writes the result into the Terraform we deploy'); the draft says 'deployment template' and names neither. Inconsistency left unmentioned.
- [ ] All quotes were checked through WebFetch model extraction, not raw HTML; one typographic difference seen (curly apostrophe in 'don't').

**Self-edit before you verify** ([checklist](../../../docs/editorial/self-edit-checklist.md)), ticked by you while you read the draft:
- [ ] Point: I can say the point in one sentence: ______________________________
- [ ] Point: The piece covers one idea; anything else is a Related link.
- [ ] Top of the page: Title: key term in the first two words, plain words, no pun, under ~65 chars.
- [ ] Top of the page: Description: adds the "so what", has a verb, doesn't repeat the title, <=150 chars.
- [ ] Top of the page: Hook: one real, linked moment (system, number, failure), not a definition or history.
- [ ] Top of the page: Nut graf: the hook's last sentence says what this piece gives the Reader.
- [ ] Top of the page: TL;DR: the answer in plain words; a reader who stops here is not misled.
- [ ] Structure: Read only the headings: do they tell the story in order?
- [ ] Structure: No heading says "Background", "Overview", "Introduction" or "Conclusion".
- [ ] Structure: No question headings outside a Repo Profile FAQ.
- [ ] Structure: Each paragraph has one topic, stated in its first sentence.
- [ ] Structure: Each H2 is under 250 words; cut, don't pad.
- [ ] Sentences: Average about 20 words or fewer; every sentence over ~30 words is split or earns its length.
- [ ] Sentences: Active voice with a named actor ("the client retried", not "requests were retried").
- [ ] Sentences: Hardest parts use the shortest sentences.
- [ ] Words: Every technical term is defined at first use, and kept where it is the right term.
- [ ] Words: One name per thing: the Source's term every time; a synonym appears once, in brackets, only if Readers search for it.
- [ ] Words: No hype or filler word from voice rule 5's banned list (`house-byline-voice.md`), and no "there is/are" openings.
- [ ] Words: At most one analogy, and it says where it stops being true.
- [ ] Words: No idioms or pop-culture references a non-native reader would miss.
- [ ] Teaching: A concrete example comes before the abstract statement.
- [ ] Teaching: At least one worked example with real numbers (Concept, Case Study, Comparison).
- [ ] Teaching: Every "what" has a "why" or a cost.
- [ ] Teaching: The common confusion is named and corrected ("what it is not").
- [ ] Teaching: The text walks the diagram box by box; the diagram is never left to speak alone.
- [ ] Teaching: Code: introduced by a sentence, realistic, under its licence; README commands verbatim.
- [ ] Evidence: Every claim has its link where the claim is made.
- [ ] Evidence: Every number has a unit, a date or version, what it is compared with, and its caveat.
- [ ] Evidence: Quotes are exact, short, credited, and followed by our own analysis.
- [ ] Evidence: Each inference says "we infer" in its own sentence.
- [ ] Our contribution: Name what this Piece adds that no Source says: ____________________
- [ ] Our contribution: Delete every paragraph that only restates a Source.
- [ ] Ending: The last section tells the Reader what they can now do, measure or decide.
- [ ] Ending: No new idea in the last paragraph; no recap of what they just read.
- [ ] Last pass: Read it aloud once. Fix every place you stumbled or got bored.
- [ ] Last pass: Cut big (sections, paragraphs), then small (words). Then cut 10% more.

<details><summary>What the fact-checker changed</summary>

- Fixed: Diagram node 'Ingress service queues a campaign event' -> 'queues a message' (post: 'The ingress service queues it')
- Fixed: Diagram prose: router-to-delivery hand-off now marked as inferred alongside the dotted results edge ('we infer two edges... The post states neither path')
- Fixed: Old-design gap sentence: dropped 'engineer-hours lost' (post gives an '8 hours per week' maintenance figure that may be the old system's load); now 'no count of incidents caused by the old design'
- Fixed: Temporal/Flowable bullet: 'Neither is a managed AWS service' now attributed to the post and matched to its wording ('Neither is available as a managed AWS service, the post says')
- Fixed: A/B example: 'one config change and no new service' (not in post) -> 'a config change, with no Step Functions definition edited directly' (post: 'without directly touching Step Functions or CloudFormation definition')
- Fixed: Rollout: 'the rollout went in phases' -> 'was planned in phases' (post: 'We planned the rollout in phases'); 'Each phase let Zillow validate' -> 'That let Zillow validate'
- Fixed: Maintenance quote 'fell along the same lines, from roughly 8 hours per week to 2' was unlinked; now linked to the post
- Fixed: Gaps list: 'no latency figures' -> 'no delivery latency' (post gives 1.2 s and 1.4 s run times in its example)
- Fixed: Gaps list: effort bullet narrowed to the build comparison, plus 'The post does not say what the 8 hours a week measured'
- Fixed: AWS docs phrase 'could potentially run more than once' put in quotation marks and credited to AWS
- Fixed: Unmarked verbatim source text reworded into our own words: closing practices list, batch rule of thumb, 'more than 20 email and push campaigns... more than 10 microservices' sentence, Figure 2 flow sentence

</details>

### Agent Host Protocol: how GitHub syncs one session across devices

`content/explainers/apis-networking/github-agent-host-protocol.md` · case-study · 2570 words · fact-checker: 78 claims checked, 13 fixed, 3 removed

**Fills Gap Day 2026-11-17** (ADR-0012): `published` is that day; `make release` keeps it.

**Drafter asks you to check:**
- [ ] How the Lead was read: direct download was blocked, so the post was read through WebFetch, which paraphrases and caps quotes at about 125 characters. Check every quoted phrase against the published page; each is a fragment the tool returned in quotation marks, not a full sentence I saw in context.
- [ ] Reconnect behaviour differs from the editor's angle. The post says a client "sends all queued actions" and waits for the server to acknowledge them. The AHP guide (reconciliation#reconnection) says the client "clears pendingActions" after reconnect. The draft reports both and says neither Source settles whether an unacknowledged action is sent again. Confirm the post's full sentence (it also mentions a "loading" state) and decide if that framing is fair.
- [ ] Line counts (about 369,000 TypeScript, about 305,000 Rust, about 415,000 for Copilot Chat, all non-test) come from the chart's alt text as paraphrased by the tool, not from body prose. The body gives only "over 400,000 lines of TypeScript, not counting tests". Check the chart.
- [ ] The session timeline ("Fix the flaky login test": 9:02 AM Copilot app, 9:40 phone approval, 11:15 VS Code diff review, 11:32 cloud agent lands the PR) comes from a figure's alt text, paraphrased. Check wording and times, or cut the paragraph.
- [ ] The N:M arithmetic (4 clients x 2 hosts = 8 integrations vs 4 + 2 = 6 implementations) is ours and labelled so. The four clients and two places an agent runs are my reading of the post's names; it never counts them. Its heading says "five clients and two hosts" but the body does not list five, so the draft avoids that count.
- [ ] ACP appears as a fourth design "set aside", but only the AHP docs (guide/ahp-and-acp) discuss it; the post does not mention ACP. The bullet says AHP sits above ACP, not instead of it. Check the heading does not overstate this.
- [ ] The LSP/DAP bullet ("A protocol that fixes agent behaviour") is my gloss on the post's sentence beginning "Unlike the domains that produced the Language Server Protocol and Debug Adapter Protocol". Check that the post contrasts them the way the draft says.
- [ ] Specification version: the docs site offers 1.0.0, 1.1.0 and a 1.2.0 draft. The homepage labels the Specification v1.1.0, while the overview page I read may be the 1.2.0 draft. The draft says only "as published on 18 November 2026". Confirm which version the linked pages show and whether the examples (serverSeq 10/11; lastSeenServerSeq 42, replay 43-44, snapshot fromSeq 50) match.
- [ ] Required channels: the post says root, session and chat must be implemented. The AHP doctrine page says a minimal host needs root and session channels. The draft follows the post.
- [ ] Spec anchors were reported by the fetch tool, not clicked: subscriptions.html#channel-uris and #channels, overview.html#channels-are-the-routing-key, #base-protocol and #extensions, transport.html#common-transports, actions.html#reducers, #actions and #action-envelope, reconciliation.html#client-side-state, #example, #reconciliation-algorithm, #why-rebasing-is-simple and #reconnection, lifecycle.html#reconnection and #connection-handshake, doctrine.html#anti-goals and #ahp-keeps-escape-hatches-explicit, ahp-and-acp.html#the-mutex-analogy, terminals.html#data-flow, what-is-ahp.html#design-requirements.
- [ ] Three labelled inferences to judge: (1) channels let a phone follow one chat without loading every terminal; (2) a replayed envelope with the client's own origin tells it the action landed before the drop; (3) every confirmation depends on the host, so a cut-off client cannot confirm. Also "We infer that one implementation on each side had hidden both" gaps, and "In our reading" for the deciding argument.
- [ ] Wording to check against the post: an update to VS Code "interrupted" the agent (the tool's paraphrase was "ended or disrupted"); "after the first couple of weeks" (post: "after the first couple weeks"); the team started in early 2026 with the Copilot SDK as the foundation; GitHub's host uses metadata "to link to project boards".
- [ ] Authors and date in Sources (Connor Peet and Rob Lourens, 17 November 2026) are as the tool reported. The draft says "One of the post's authors keeps a macropad"; the post says "I" without saying which author.
- [ ] Two internal links: /languages-runtimes/agentic-development/ (definition of harness, itself quoting Martin Fowler's site) and /distributed-systems/durable-execution/ (replaying a log). Confirm both are live and that the comparison with durable execution is fair.
- [ ] `published`/`reviewed` are left at the skeleton's 2026-11-17T06:32:29Z; frontmatter `related` is left commented out.

**Fact-checker could not confirm:**
- [ ] Spec version of linked pages: unversioned doc URLs show a selector listing '1.2.0 (draft)', 1.1.0, 1.0.0; homepage button says 'Specification v1.1.0'. WebFetch could not show which is active (two fetches said 1.2.0 draft selected, three could not tell). Draft still says only 'as published on 18 November 2026'. Author should open page and name version if wanted. Examples do match linked pages: seq 10/11 in guide; lastSeenServerSeq 42, replay 43–44, snapshot fromSeq 50 in lifecycle.
- [ ] All Source text read through WebFetch only (direct download denied in this session; tool refuses full verbatim dumps, returns fragments under ~120 chars). Every quoted phrase was confirmed YES by targeted query and surrounding sentences were returned in fragments, but nobody in this pipeline has seen the raw page. Post uses curly apostrophes (It’s, wasn’t, isn’t); draft uses straight.
- [ ] Line counts (about 369,000 TS, about 305,000 Rust, about 415,000 Copilot Chat, non-test) exist only in chart alt text, confirmed in fragments; body prose gives only 'over 400,000 lines of TypeScript, not counting tests'. Draft's opening sentence links these to the post without saying 'chart'; later paragraph does say 'on the post's chart'. Left as is.
- [ ] Session timeline (9:02 AM Copilot app, 9:40 AM phone approval, 11:15 AM VS Code diff, 11:32 AM cloud agent lands PR; title "Fix the flaky login test") is figure alt text only; fragments match draft. Kept.
- [ ] Reconnect framing kept as fair: post sentence sits under 'Reconnection and replay in poor network conditions are easy: a client sends all queued actions, presents the right "loading" state, and waits for the server to reconcile and acknowledge the actions before confirming them'; guide says client 'clears `pendingActions`'. Draft omits the 'loading' state; Author may want it.
- [ ] `_meta` link anchor #connection-handshake exists, but the `_meta` sentence sits under sub-heading 'Initialize Response (Server → Client)' inside it. Anchors reported by tool, not clicked. All other draft anchors appear in tool-reported heading lists.
- [ ] Copilot app link https://github.com/features/ai/github-app is the URL the post itself links; page title returned was 'GitHub Copilot · Your AI coding agent', not an app-specific title.
- [ ] Macropad: post says 'I have a macropad sitting on my desk' without saying which of two authors; draft's 'One of the post's authors' kept.
- [ ] Three labelled inferences (phone follows one chat; replayed own-origin envelope means action landed; cut-off client cannot confirm) plus 'one implementation on each side had hidden both' and 'In our reading' are HastCode's own and stay labelled; no Source states them. Durable-execution comparison is our analogy, unsourced by design.
- [ ] Internal links: /languages-runtimes/agentic-development/ and /distributed-systems/durable-execution/ both exist under content/, carry verified: and no draft flag; harness definition wording matches the Piece. Not checked on live site.
- [ ] published/reviewed still skeleton 2026-11-17T06:32:29Z; related still commented out; draft: true kept; no verified: added.

**Self-edit before you verify** ([checklist](../../../docs/editorial/self-edit-checklist.md)), ticked by you while you read the draft:
- [ ] Point: I can say the point in one sentence: ______________________________
- [ ] Point: The piece covers one idea; anything else is a Related link.
- [ ] Top of the page: Title: key term in the first two words, plain words, no pun, under ~65 chars.
- [ ] Top of the page: Description: adds the "so what", has a verb, doesn't repeat the title, <=150 chars.
- [ ] Top of the page: Hook: one real, linked moment (system, number, failure), not a definition or history.
- [ ] Top of the page: Nut graf: the hook's last sentence says what this piece gives the Reader.
- [ ] Top of the page: TL;DR: the answer in plain words; a reader who stops here is not misled.
- [ ] Structure: Read only the headings: do they tell the story in order?
- [ ] Structure: No heading says "Background", "Overview", "Introduction" or "Conclusion".
- [ ] Structure: No question headings outside a Repo Profile FAQ.
- [ ] Structure: Each paragraph has one topic, stated in its first sentence.
- [ ] Structure: Each H2 is under 250 words; cut, don't pad.
- [ ] Sentences: Average about 20 words or fewer; every sentence over ~30 words is split or earns its length.
- [ ] Sentences: Active voice with a named actor ("the client retried", not "requests were retried").
- [ ] Sentences: Hardest parts use the shortest sentences.
- [ ] Words: Every technical term is defined at first use, and kept where it is the right term.
- [ ] Words: One name per thing: the Source's term every time; a synonym appears once, in brackets, only if Readers search for it.
- [ ] Words: No hype or filler word from voice rule 5's banned list (`house-byline-voice.md`), and no "there is/are" openings.
- [ ] Words: At most one analogy, and it says where it stops being true.
- [ ] Words: No idioms or pop-culture references a non-native reader would miss.
- [ ] Teaching: A concrete example comes before the abstract statement.
- [ ] Teaching: At least one worked example with real numbers (Concept, Case Study, Comparison).
- [ ] Teaching: Every "what" has a "why" or a cost.
- [ ] Teaching: The common confusion is named and corrected ("what it is not").
- [ ] Teaching: The text walks the diagram box by box; the diagram is never left to speak alone.
- [ ] Teaching: Code: introduced by a sentence, realistic, under its licence; README commands verbatim.
- [ ] Evidence: Every claim has its link where the claim is made.
- [ ] Evidence: Every number has a unit, a date or version, what it is compared with, and its caveat.
- [ ] Evidence: Quotes are exact, short, credited, and followed by our own analysis.
- [ ] Evidence: Each inference says "we infer" in its own sentence.
- [ ] Our contribution: Name what this Piece adds that no Source says: ____________________
- [ ] Our contribution: Delete every paragraph that only restates a Source.
- [ ] Ending: The last section tells the Reader what they can now do, measure or decide.
- [ ] Ending: No new idea in the last paragraph; no recap of what they just read.
- [ ] Last pass: Read it aloud once. Fix every place you stumbled or got bored.
- [ ] Last pass: Cut big (sections, paragraphs), then small (words). Then cut 10% more.

<details><summary>What the fact-checker changed</summary>

- Fixed: VS Code update: draft said update "interrupted" agent; post says "Update VS Code, and every agent stops as the application restarts." Reworded to match.
- Fixed: Remote machine: draft said connecting to remote also ended agent; post only says you get "an entirely new extension host commit-locked against the version of the VS Code workbench you have". Reworded; dropped 'So did'.
- Fixed: Options section: post lists no rejected designs (WebFetch confirmed it never says team considered and rejected reusing Copilot Chat). Heading changed from 'The designs GitHub set aside before writing a protocol' to 'The alternatives AHP rules out, and why'; intro now says 'The post lists no rejected designs. We read three alternatives out of it; the AHP docs supply a fourth.' Intro paragraph phrase 'the designs GitHub set aside' changed to 'the alternatives it rules out'.
- Fixed: N:M bullet: removed drafter's own count of 4 clients x 2 places (8 vs 6), which contradicted post heading 'Shipping five clients and two hosts'. Now uses heading's count (5 x 2 = 10 vs 5 + 2 = 7, still labelled our arithmetic) and notes body never lists the five. N:M quote now tied to what post's 'This' refers to: coordinating sessions on own devices and in cloud sandboxes.
- Fixed: Required channels: 'AHP requires three channels' changed to 'The post requires three channels' — only the post says root/session/chat are mandatory; doctrine page says minimal host needs 'root and session channels'; spec overview/subscriptions pages name no mandatory set.
- Fixed: Protocol-notification link: 'does not replay protocol notifications such as `sessionAdded`' was linked to specification/lifecycle.html#reconnection, which names no notification. Re-linked to guide/reconciliation.html#reconnection, which says 'Protocol notifications (like `sessionAdded`/`sessionRemoved`) are not replayed'.
- Fixed: Reconnect paragraph: 'Neither Source says whether a client should send an unacknowledged action again' overstated (post says client 'sends all queued actions' on reconnect). Changed to 'The specification's lifecycle page does not say…', which the page confirms.
- Fixed: Evidence section: draft called both findings 'gaps'. Post calls PTY-less shell output 'a gap in the protocol' and session-URI structure 'bad assumptions'. Reworded to 'a real gap and a bad assumption'.
- Fixed: 'client libraries in Rust, TypeScript, Kotlin, Swift, .NET and Go' → 'libraries' (post: 'Agent Host Protocol libraries are available in…'; says clients and servers can be built).
- Fixed: TL;DR: 'open protocol' → 'open standard' (post's wording).
- Fixed: Terminology note: post uses both 'agent host' and 'server'; sentence now says so.
- Fixed: 'Its only numbers are lines of code' → 'Its only measurements' (post also has times and client/host counts).
- Fixed: Trimmed options section to stay under 250-word section-length warning after edits.
- Removed: Drafter's 4-client / 2-place enumeration (VS Code, Copilot app, iOS, Android; user's device, cloud sandbox) and 8-vs-6 arithmetic — post never counts them that way.
- Removed: Claim that connecting to a remote machine ended the agent ('So did…') — post does not say this.
- Removed: Claim that the post 'names three' designs GitHub set aside.

</details>

### Not drafted

- none

<details><summary>Leads the editor rejected</summary>

- `qj3x47yl`: GitHub status page for the Pull Requests incident gives merge-timeout rates but no cause; root cause analysis is promised later.
- `frx4gkkq`: Discord status page states no cause or mechanism, only that a fix was implemented.
- `emdwjqhx`: AWS Health Dashboard item links the generic status page and states no cause.
- `l4p7nqgs`: Rust parallel frontend on nightly is a real design change with benchmarks, but no live Explainer teaches its lesson, and the post does not describe the mechanism (query system, locking) needed for an Explainer.
- `k5g42mqo`: vLLM v0.34.0 is a long release list; its design change (KV cache sized from measured free memory) is not the lesson of any live Explainer, only the same topic.
- `7qtr7wjs`: Cloudflare's post on Quad9 leaving France is policy advocacy; no live Explainer teaches its lesson and the technical content is introductory.
- `542g6wfu`: Netflix 'Engineering AI Discover' returned HTTP 403 on both fetches, so its mechanism could not be judged; reconsider on a later run.
- `gifu2zhy`: Valhalla value objects is a talk video landing page with a 47-word abstract and no readable text.
- `5sgodn7v`: Rakbank observability has mechanism and numbers, but authors' Rakbank affiliation is unstated, it reads as an AWS reference architecture, and Distributed Systems already takes two picks today.
- `3u3twadt`: Flickr's Redis to Valkey migration has good mechanism, but the cutover lesson overlaps the Topstep Brief and planned-database-switchover, and Databases & Storage is not thin; keep for a later run.
- `ucr5ttmb`: RDS Proxy pinning is a candidate Concept (connection multiplexing and session pinning), but the fetched text was truncated after the first of three conditions and Databases & Storage is not thin; keep for a later run.
- `eq4gd4tl`: Apache Cloudberry 2.2.0 is a release announcement with a kernel upgrade; no live Explainer teaches a lesson from it.
- `2ucuqqlb`: Daniel Stenberg's post on curl rewrites in Rust is opinion, not an event at a system or a taught idea.
- `s2q2xzby`: APNIC's DNS OARC 47 post is a conference roundup of several topics, not one event or one idea.
- `pxs4bbg7`: RDS read replicas behind a Network Load Balancer is an AWS how-to, not a named system or a genuine choice between alternatives.
- `uphw4t5j`: ElastiCache multi-tenancy isolation models is a possible Comparison, but Databases & Storage already has four Comparisons; not fetched in depth.
- `ofa5rk22`: Meta Intelligent Delivery is a first-party Case Study candidate, but AI/ML Systems is the fullest Vertical with six Case Studies, three of them Meta's; not fetched.

</details>


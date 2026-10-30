# Daily review queue, 2026-10-30

Drafted by `make daily` (ADR-0009). Every Piece is `draft: true` and unverified. Verify every claim against its Source, then set `verified:` and remove `draft: true` yourself (ADR-0007).

## Run 20261030-095803

### PrivacyLib: how Meta checks privacy policy at every data access

`content/explainers/distributed-systems/meta-privacylib.md` · case-study · 2600 words · fact-checker: 68 claims checked, 8 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Source read only through WebFetch, in four passes (fetching the raw HTML was blocked); check every quote against the live page
- [ ] Options section: the post has no formal comparison; four designs it moved away from (per-system code, remote config fetch, network-only check, fleet-wide config push) are framed as options, and the section says so
- [ ] 'Propagation enabled' on the Policy Zones module is shown but not explained in the post; the gloss is from the 2025 Policy Zones post (information flow control, labels carried downstream)
- [ ] ServiceRouter called 'Meta's global service mesh' from the OSDI 2023 paper abstract, not from the PrivacyLib post
- [ ] 'we infer' sentences: the old cost multiplies across requirements and systems; Request entry fixes context once; why build-time compilation works; the token lets the network layer be strict and nudges owners onto typed clients; simulation and the runtime control answer the two incidents
- [ ] Agents section: 'a check at the data does not need to predict the caller' is our analysis
- [ ] Source titles for the lineage and data-understanding posts come from the post's link text, and for the Messenger and Hack posts from their URLs; purpose-limitation, Investing in Privacy, Policy Zones and ServiceRouter titles were confirmed
- [ ] Quote 'low per-check latency' and 'just the starting point' wording should be checked against the live page
- [ ] Code identifiers (rollingOut(), enabled(), logger.input.write.zones, MetadataTraceLoggerConfig, MESSAGING_TRAFFIC_METADATA) appear only as names; no code is quoted
- [ ] Word count of about 2,600 is an estimate; the word-count command was blocked

**Fact-checker could not confirm:**
- [ ] The full post could not be read verbatim (WebFetch only summarises and paraphrases). Every quote was confirmed through targeted verbatim checks, but the Author should spot-check the quotes against the live page
- [ ] ServiceRouter is cited to the OSDI 2023 page, which confirms 'Meta's global service mesh'. The PrivacyLib post itself links to atscaleconference.com/servicerouter-hyperscale-service-mesh-at-meta/ instead
- [ ] The Bloom-filter 'fixed rate' and the 'fine-grained sampling rates adjust' wording were confirmed only through a paraphrased fetch answer
- [ ] Inference 'the replay simulation and the runtime dial exist for exactly those two incidents' is marked 'we infer'. The post may draw this link itself; not checked

<details><summary>What the fact-checker changed</summary>

- Fixed: Lede: 'the common path never makes a network call' -> 'avoids a network call', now linked (the post says most checks are cached and the common path avoids a network call; token-less RPCs still fall back to a network check)
- Fixed: PAI paragraph: the post says the 'Can we prove...' question underlies many privacy requirements, not PAI itself; reworded, PAI link moved to the earlier-posts sentence
- Fixed: Pre-PrivacyLib link text: 'integrated privacy enforcement separately into the clients of each ... system' -> 'implementing the same enforcement logic separately in each system' (Source wording)
- Fixed: Typed client vs raw packet sentence was unmarked analysis; now 'We infer why:'
- Fixed: Messaging metadata example: 'such as delivery timing' -> 'such as when messages are sent or delivered' (Source wording)
- Fixed: Annotation change 'is a two-line edit' -> 'can be a two-line edit' (the post's example, not a rule)
- Fixed: AgentLib 'built on' -> 'based on' PrivacyLib's design principles (Source wording)
- Fixed: Next steps: 'the plan is to evaluate agent behaviour' -> 'the harder problems are evaluating internal agent behaviour while it happens and reasoning about...' (Source frames them as harder problems)

</details>

### Logfwdr and Buftee: how Cloudflare keeps analytics on linear hardware

`content/explainers/distributed-systems/cloudflare-analytics-pipeline.md` · case-study · 2350 words · fact-checker: 68 claims checked, 12 fixed, 3 removed

**Drafter asks you to check:**
- [ ] The Lead was read only through WebFetch summaries, section by section, because curl was denied. Check every number and paraphrase against the published page. The quotes were confirmed word for word in a separate fetch: 'unsustainable against our exponential traffic growth', 'instant, independent query failover', 'attenuates storage of the "most seen" keys', 'knocking out roughly a year's worth of compressed traffic growth', 'not compromising the quality of data delivered', 'deferring two and a half years of hardware expansion', 'without backpressure'. 'one and only schema' and 'to eliminate bias' were not part of that check.
- [ ] The Source contradicts itself and the draft says so: load balancing's net saving is 36 GB/s in its own section but ~38 GB/s in Impact, and Ready-Analytics fell 28→16 GB/s in a chart caption but ~26→13 GB/s in Impact (described as 43% in one place and 50% in the other).
- [ ] Ready-Analytics' sampled rate of ~8 GB/s comes from text the fetch reported as a typo ('~8GB/'). Confirm it.
- [ ] Lead links have no heading anchors because the fetch showed no id attributes. Add #load-balancing and the like if the page has them.
- [ ] Inferences are labelled: 'non-regional' as the data outside Customer Metadata Boundaries; double sampling reusing _sample_interval (from the Jan 2025 post); why aggregation repeats; Buftee turning replication into an async job; each region needing capacity for all non-regional traffic; 'replicate at the narrowest point'; the deciding argument in the options section.
- [ ] Our own arithmetic, marked as ours: the relative error on a slice shrinks about as 1/sqrt(rows kept). Check it is fair and clearly labelled.
- [ ] The order in the diagram (sample, then aggregate, inside Logfwdr) is our reading. The Lead describes the two in separate sections and does not order them inside Logfwdr.
- [ ] Earlier-post titles and dates come from WebFetch. 'How we make sense of too much data' is dated 27 Jan 2025 from JSON-LD and no author is shown. The incident post is dated 26 Nov 2024 and the ClickHouse post 14 May 2026. The docs page title 'Customer Metadata Boundary' was not fetched.
- [ ] Buftee facts ('over one million buffers', per-buffer isolation) come from the Nov 2024 incident post, which describes Buftee in the context of Cloudflare Logs/Logpush jobs. Check the generalisation is fair.
- [ ] No related HastCode Explainer genuinely teaches a step, so the related: comment was left as the skeleton had it.

**Fact-checker could not confirm:**
- [ ] Heading anchors on blog.cloudflare.com pages: WebFetch shows no id attributes, so no #anchors were added (curl denied, raw HTML not inspected)
- [ ] The Lead was checked only through WebFetch quote extraction, not raw HTML. Quotes came back verbatim, but every claim should get one last read on the live page
- [ ] Chart internals (axes of the rows-vs-estimate chart, Ready-Analytics 'ghost' curves) not verifiable from text; draft relies only on the text and captions
- [ ] Our labelled 1/sqrt(rows) error heuristic is a standard sampling-variance approximation, not from any Source

<details><summary>What the fact-checker changed</summary>

- Fixed: 2022 split: 'small product-specific messages' -> 'product-specific schemas'; 'every consumer' -> 'every internal consumer' (Source: Cloudflare-internal consumers ingest everything, product-specific schemas)
- Fixed: 2024 compression: 'little impact on CPUs' -> 'no noticeable difference in CPU usage'; link text 'a daily issue' -> 'occurring daily' (Source says congestion 'on a daily basis')
- Fixed: Traffic doubling: now says it was manageable with a handful of streams, and traffic doubled as disaster resilience expanded to almost all analytics streams (Source wording)
- Fixed: Aggregation vs sampling: 'leaves the choice per dataset to producer teams' -> 'offers producer teams aggregation as another option besides sampling'
- Fixed: Sampling exponents: 'falls short of x^0.25' -> 'the post says x^0.25 was not observed'; 100x example labelled as our arithmetic and corrected from ~10x to ~9x (100^0.485 ≈ 9.3)
- Fixed: Ready-Analytics described as the ClickHouse post does (internal teams stream data into a single shared table), not a 'product'; unconfirmed #the-setup... anchor dropped
- Fixed: Buftee: unconfirmed #buftee anchor dropped; isolation claim narrowed to what the incident post says (each Logpush job gets its own buffer, so one customer's failed job does not affect another's)
- Fixed: Late replication: unsupported 'in the region that received it' and 'local inserter' removed; now quotes the post: a remote inserter makes the second ClickHouse copy after Buftee
- Fixed: Diagram intro now says two placements are our reading: sampling before aggregation in Logfwdr, and Buftee plus the first inserter in the receiving region
- Fixed: 'replica costs a sixth' labelled as our arithmetic
- Fixed: Aggregation fit: 'works best for datasets that collapse easily... whose sums and counts must be exact' -> 'is powerful for datasets that can collapse easily, such as usage metrics'
- Fixed: Sources: unconfirmed byline 'Cloudflare Data team' removed (no byline on the page); Jan 2025 post retitled to the full 'Over 700 million events/second: How we make sense of too much data'
- Removed: 'Cloudflare plans to keep moving aggregation earlier' (not in Source; replaced with 'expects the savings to keep mounting')
- Removed: Duplicate 'Cloudflare leaves the choice to producer teams, dataset by dataset' in the costs section (choice per dataset not stated)
- Removed: 'A local inserter writes this region's ClickHouse' (term not in Source)

</details>

### Expression volatility: what a sharded SQL router can push to shards

`content/explainers/databases-storage/volatile-expressions-sharded-sql.md` · concept · 2429 words · fact-checker: 62 claims checked, 18 fixed, 2 removed

**Drafter asks you to check:**
- [ ] The Lead was read through WebFetch summaries because curl was denied. The quotes from it were checked word for word in a second WebFetch pass, but italics and bold were dropped. Check (a)–(j) against the live page: 'will produce the same output each time…', 'produces the same output within a statement', 'may produce a different output on every single invocation', 'returns a new, unique value each time it's evaluated', 'we cannot generate the sequence values in advance and then draw from a cache', 'if we can relay the entire query to the shard…', 'it is not guaranteed that Neki and Postgres agree…', 'Each Neki shard is its own complete database system…' (cut short), 'the physical host for all sequences', 'the natural clock drift and the different arrival times…'.
- [ ] Diagram: the post does not say whether the router batches nextval requests to the authoritative shard. The diagram shows one request per row, and the text says 'we draw one per value'. Confirm, or reword.
- [ ] Diagram/hook: the SELECT … FOR UPDATE is shown going to 'Shards' (all of them). The alt text says shards return rows, but the post does not state the fan-out explicitly.
- [ ] Table 1 and Table 2 values (now() under UTC and Pacific/Auckland; to_char TMDay giving Wednesday/Onsdag/星期三) come from the images' alt text. Check them against the rendered tables.
- [ ] Hook: 'Swap nextval for lower(name) and none of that machinery is needed' rests on Query 6, where immutable expressions run shard-locally. The post does not show an UPDATE with lower(name).
- [ ] Citus #240: the dates (opened as pg_shard#108 on 5 May 2015, closed via PR #472 on 28 Apr 2016) and the maintainer quotes (jasonmp85, anarazel) come from the issue page. What #472 actually changed was not confirmed, so the text says only that the issue closed through it.
- [ ] WarehousePG is EDB's docs for a Greenplum-derived engine. Check that it is acceptable as a first-party Source for the 'where it shows up' section.
- [ ] Vitess: '1,000 suggested' for the cache column, block loss on a tablet crash, and QPS limited by RPC latency are cited to the v24.0 docs anchors. Check the anchors resolve.
- [ ] Inferences marked 'we infer': an unlabeled (VOLATILE) function stays on the router; a mislabeled IMMUTABLE function could compute under each shard's own settings; the general rule that state outside the arguments must be pinned to one place; per-shard nextval would risk duplicates.
- [ ] Rydberg's closing exercise (how to handle now() in the UPDATE) is paraphrased from the summary. Check the wording.

**Fact-checker could not confirm:**
- [ ] Exact italics and bold in the Neki post quotes. curl was denied, so WebFetch confirmed the words but not the emphasis. The Citus quote keeps _all_ italic, which matches the issue page.
- [ ] Heading anchors were inferred from section headings, not from rendered id attributes: the Neki post (#putting-it-all-together, #router-only-expressions, #enter-volatility, #table-expressions and others), the Neki architecture post (#router, #shards) and Vitess (#initializing-a-sequence, #mysql-auto_increment-feature). The headings themselves were confirmed.
- [ ] Table 1 and Table 2 values were confirmed only from the images' alt text (UTC 2026-09-09 00:00:00+00, Pacific/Auckland 12:00:00+12; en_US Wednesday, sv_SE Onsdag, zh_CN 星期三), not from the rendered images.
- [ ] Diagram fan-out: the SELECT … FOR UPDATE going to all shards rests on the post's Query 9 caption, 'shards 1, 2 and 3 return rows'. The text does not state it outright.
- [ ] Citus PR #472: the text claims only that the issue was closed through it on 28 Apr 2016. What the PR changed was not checked.
- [ ] WarehousePG: the docs are hosted by EDB and imported from the warehouse-pg/whpg-docs repository. We treat them as the product's own documentation. Whether that counts as a first-party Source is the Author's call.

<details><summary>What the fact-checker changed</summary>

- Fixed: TL;DR: the 'now() evaluated once on the router' link now points to #putting-it-all-together (Query 6). The router-only section says only that it is evaluated in the router.
- Fixed: 'How Postgres classifies expressions': 'from its definition and inputs' replaced with the post's own wording, 'one of the phases of query planning'.
- Fixed: Diagram: 'Locked rows, streamed back' changed to 'Ids of the locked rows'. The post never says 'streamed'.
- Fixed: Trace: 'Each shard would call its own nextval(), and values must be unique' replaced with the post's reason, 'centralized governance to avoid duplicate values across shards' (#enter-volatility). 'Number of rows not known in advance' changed to 'Neki does not know how many rows each shard will return'.
- Fixed: Trace: 'streams their keys back' changed to say the rows stay locked until the per-row updates arrive and the shards return their ids. Linked to the post.
- Fixed: Trace: 'we draw one per value' changed to 'the diagram draws one per value'. The post does not say whether requests are batched.
- Fixed: Citus: 'a maintainer warned' changed to 'Jason Petersen warned', and 'another proposed' to 'Andres Freund proposed'. The issue page does not state anyone's maintainer role.
- Fixed: Network cost: 'a major cost' changed to the post's 'a very large cost'.
- Fixed: Hook: 'none of that machinery is needed' for lower(name) is now marked 'we infer'. The post shows no UPDATE that uses lower(name).
- Fixed: General case: 'router computes every column of every row' for SELECT age + 1, lower(name), now() is now marked 'we infer'. The post's example there is SELECT age, name.
- Fixed: 'Postgres guarantees them by having one clock, one session and one sequence' is now marked 'we infer'.
- Fixed: now() quote 'the wall clock time recorded at the start of the transaction': anchor corrected from #the-general-case-of-neki-query-planning to #router-only-expressions.
- Fixed: Removed 'Each shard starts its own transaction', which is not in the post. 'Evaluated once' is now linked to Query 6 (#putting-it-all-together).
- Fixed: Citus snapshot condition: 'stable functions' changed to 'STABLE functions' to match the comment's wording.
- Fixed: Table expressions: 'cannot be deduplicated even when inputs repeat' replaced with the post's 'must be evaluated afresh on each row', linked.
- Fixed: Vitess: 'limited mainly by RPC latency' changed to 'only', matching the docs' 'QPS only limited by the RPC latency'.
- Fixed: 'What volatility is not': 'to_char() could run on a shard whose lc_time matches' is now marked 'we infer'.
- Fixed: Closing exercise reworded to match the post: how a stable expression such as now() would be handled for 'a variation of' that UPDATE (Query 7), not 'that same UPDATE'.
- Removed: 'Postgres does not check it:' (Practice one). The volatility docs do not say this.
- Removed: 'Each shard starts its own transaction' (Practice two). The Neki post does not say this.

</details>

### RDS for PostgreSQL snapshot restore: pg_prewarm vs scans vs waiting

`content/explainers/databases-storage/snapshot-restore-hydration.md` · comparison · 2600 words · fact-checker: 70 claims checked, 9 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Lead facts came through WebFetch summaries, not the raw page (Bash curl was denied). Check every Lead quote word for word: "10x to 15x slower", "typically lasts hours depending on storage type", "isn't available for RDS managed Amazon EBS volumes", "VACUUM isn't reliable for hydration", "depends entirely on the dead-tuple ratio when the snapshot was taken", "Passive hydration only", "expected and intentional", "higher throughput and reduced latency", "full database heap hydration in one operation", "Minimal", "Targeted", "Partial".
- [ ] Lead test numbers: db.m5.4xlarge, gp3 at 20,000 IOPS and 600 MiBps, 15 parallel sessions, about 5 h vs 23% at 5 h and over 24 h estimated, 44 s vs 571 s, 123 s vs 1,531 s, ReadLatency 2–3 ms vs up to 40 ms, CLI sample at 65 percent. Also check the author names (Sandeep Modappady, Arunkumar Pattabiraman, Mahesh Thakur) and the date 14 Oct 2026.
- [ ] 'while the other instance was still hydrating' during the post-prewarm queries is our reading of the graph times (14:45–15:35) and the 23% figure. Confirm the Lead supports it.
- [ ] Inferences marked 'we infer': COUNT(*) may use an index-only scan and leave the heap cold; COUNT(*) never reads TOAST; pg_dump never reads index blocks; waiting also loads catalogs and non-main forks; pg_prewarm block counts work as a per-table receipt. Decide whether each is safe as written.
- [ ] The Cost section says no option is billed per read. That is inferred from the RDS storage doc's 'charged for the provisioned resources whether or not you use them' and 'pay for only the additional provisioned storage performance'. Also confirm that no RDS page offers the EBS Provisioned Rate for Volume Initialization.
- [ ] The 12x to 13x slowdown in the Cost section restates the Lead's 13x and 12x speedups the other way round. Check the wording.
- [ ] PostgreSQL section numbers in Sources: F.30 pg_prewarm, 66.2.1 TOAST, 11.9 Index-Only Scans (PG18). Confirm 11.9.
- [ ] In the Trade-offs table, the Waiting 'AWS speed rating: Not rated' and VACUUM 'Least I/O' cells are our summaries, not Lead ratings.
- [ ] No related HastCode Explainer genuinely teaches a step here, so no related links were added.

<details><summary>What the fact-checker changed</summary>

- Fixed: pg_dump: Lead's table puts 'All heap' in both heap and TOAST columns (draft said 'all TOAST')
- Fixed: pg_dump reads every TOAST value: marked 'we infer' (TL;DR, pg_dump section, Coverage)
- Fixed: 'while the other instance was still hydrating': now an explicit inference from 23% at 5 h
- Fixed: prefetch speed: Lead rates it 'fastest for hydration' and says it queues parallel I/O; removed draft's causal 'because'
- Fixed: 'sustained high ReadLatency' -> 'sustained elevated' (Lead wording)
- Fixed: 'Passive hydration only' is from the Lead's notes on the gp3 no-prewarming row, not an instance label: reworded
- Fixed: '12x to 13x slower afterwards' -> 'took 12x to 13x as long' (Lead's speedups restated)
- Fixed: Waiting inference: 'blocks no SQL-level method reads, every fork' -> 'blocks SQL-level methods skip unless you name them, such as non-main forks'
- Fixed: VACUUM ANALYZE: Lead recommends it for statistics refresh in production DR; attribution narrowed

</details>

### Bundled vs platform trust stores: who revokes a bad CA, how fast

`content/explainers/languages-runtimes/bundled-vs-platform-trust-stores.md` · comparison · 2571 words · fact-checker: 64 claims checked, 7 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Go blog anchors (#emudhra-incident, #incident-timeline, #lessons-learned-and-future-improvements, #upgrading-and-verifying-your-deployments) are guessed from the Go blog's id pattern and were not seen in the HTML. Bash curl was denied.
- [ ] Cloudflare blog anchors (#timeline, #our-response, #looking-forward) are guessed from the heading text and were not seen.
- [ ] Go blog quotes were read through WebFetch excerpts, so check the exact wording: 'trivially weak', 'recovery of the private keys in microseconds', 'present in most trusted root programs', 'effectively blocking their usage in Firefox', 'does not have a uniform built-in mechanism to handle revocation checking', 'which similarly removed trust', 'will automatically reject', 'just under four and a quarter hours', and the upki description 'provide reliable, privacy-preserving and efficient certificate revocation'.
- [ ] Go timeline: alert at 11:45 PDT, Mozilla confirmed at 12:30, patch at 14:00, release started at 15:00, announced at 19:13 on 19 October. 'First-ever URGENT' is from the blog's intro.
- [ ] Go release announcement CVE-2026-111943 and the quote 'will now filter out chains which contain these intermediates'. The date 20 October in Sources is the timestamp golang-announce displays, in an unknown time zone.
- [ ] OpenJDK page: CVE written as 'CVE 2026-114820' (the page's format), the four aliases including comodoaaaca [jdk] (the page gives no reason for it), the quote 'take one of the following actions immediately', and the page's last update of 24 October 2026. The claim that removing roots distrusts everything under them is labelled 'we infer'.
- [ ] Debian bug 1069251: a later comment saying the bug blocks removing an eMudhra intermediate came from a WebFetch summary. Check it on the bug page.
- [ ] The Mozilla wiki CRLite status may be stale (the page says early 2019). The draft says only that CRLite uses the same delivery mechanism as OneCRL.
- [ ] Cloudflare author Bas Westerbaan comes from the page's author section, not a byline. The publish date of 20 October comes from the page's structured data.
- [ ] The inferences that 'a root store cannot exclude specific intermediates' and that Java CRL checking is opt-in are marked as inference or worded from OpenJDK's 'if you have not already enabled'.
- [ ] Not used: Apple and Microsoft block dates, the eMudhra revocation date and the fact that six intermediates were TLS and two code-signing. Only third parties or unquoted summaries gave them.

**Fact-checker could not confirm:**
- [ ] Mozilla wiki anchors #OneCRL and #CRLite: the headings were confirmed but not their ids. They were kept because MediaWiki builds ids from the heading text
- [ ] pkg.go.dev anchors #pkg-overview, #Certificate.Verify and #SystemCertPool follow pkgsite's standard pattern but their ids were not seen in the HTML
- [ ] Where the WebFetch output was a paraphrase rather than an exact excerpt, quotes were checked through its yes/no answer with a quoted sentence, not by reading the raw HTML (Bash curl was denied). 'In order of preference' for the OpenJDK actions is based on the page's labels 'Preferred / Otherwise / Or else'

<details><summary>What the fact-checker changed</summary>

- Fixed: Removed guessed Go blog anchors (#emudhra-incident, #incident-timeline, #lessons-learned-and-future-improvements, #upgrading-and-verifying-your-deployments): the headings exist but their id attributes could not be seen, so links now point at the bare post URL
- Fixed: Removed guessed Cloudflare anchors (#timeline, #our-response, #looking-forward) for the same reason
- Fixed: Removed unconfirmed go.dev security policy anchors (#urgent, #private) in the body and Sources
- Fixed: Go section: 'So no list could be pushed to it' was our own reasoning, now labelled 'We infer'
- Fixed: Debian bug 1069251: the title has since been changed, so it now reads 'originally titled'. The paraphrase of the comment now matches what Stian Mathiassen wrote on 22 October 2026: the bug keeps Java's keystore from updating when ca-certificates removes the eMudhra intermediates. It no longer says the bug 'blocks removing' them
- Fixed: CRLite: the draft said Mozilla 'extends' the list to leaf certificates, but the wiki only describes CRLite as an expected, not-yet-complete mechanism (as of 2019). It now says Mozilla's page describes CRLite as its leaf-certificate mechanism and keeps the verbatim delivery quote
- Fixed: upki and Linux parity: the plan is now worded as the Go blog puts it ('bring parity to revocation support on Linux') and linked to the blog, since the upki quote comes from the blog

</details>

### Adaptive heap sizing: how ZGC trades memory for GC CPU on its own

`content/explainers/languages-runtimes/zgc-adaptive-heap-sizing.md` · concept · 2450 words · fact-checker: 58 claims checked, 7 fixed, 2 removed

**Drafter asks you to check:**
- [ ] Anchors on openjdk.org/jeps/546 (#Summary, #Goals, #Motivation, #Configuring-the-maximum-heap-size-is-difficult, #Configuring-the-JVM-to-be-a-good-neighbor-is-difficult, #Adaptive-heap-sizing, #Description, #Tuning-resource-usage, #A-good-steward-of-CPU-and-memory, #A-good-neighbor, #Alternatives, #Risks-and-Assumptions) were built from the JEP 545 anchor pattern. Only #Case-study-in-adaptive-heap-sizing was seen in the page. Check that each one resolves.
- [ ] All JEP 546 quotes were read through a fetch tool that caps quotes at 125 characters. Check each quote word for word against the page: 'The Java Virtual Machine presents the illusion of an infinite heap.', both Goals quotes, 'GC CPU usage tends to vary inversely with heap size.', 'smoothed out over time in order to avoid disruptive swings', 'If the maximum heap size is too small, ZGC will use more CPU time', 'neither necessary nor desirable', 'as if a soft maximum heap size was set', 'The relationship between intensity and actual GC CPU usage is intentionally unspecified', 'ZGC prioritizes being a good neighbor over achieving its target intensity', 'ZGC will use more memory than it did in past releases', 'finds a reasonable heap size that yields good performance'.
- [ ] The 'Four JVMs in 16GB' section maps each case-study phase onto a branch of the diagram. That mapping is our reading, and so is the claim that the contraction in minutes 12 to 24 was the good-neighbour bias. Both are marked 'we infer' / 'our reading'. The JEP gives no GC CPU numbers for any phase.
- [ ] The hook says 'Every heap size in that chart was the garbage collector's choice'. The JEP does not say whether -Xmx was set in the case study. Soften this if the Author wants it stricter.
- [ ] The paper 'Heap Size Adjustment with CPU Control' could not be read: Semantic Scholar and ACM returned 403, and the PDF would not parse. It is cited only as 'credited by the JEP' for smoothing. A search snippet gave the authors (Tavakolisomeh, Shimchenko, Österlund, Bruno, Ferreira, Wrigstad), MPLR '23 and DOI 10.1145/3617651.3622988. None of that is in the text. Add it if verified.
- [ ] The JEP 545 claims come through the fetch tool: the initial heap went from 1/64 of physical memory to 2MB, and aggressive expansion during warmup. Check them on the JEP 545 page. The JEP 546 link for warmup expansion points to openjdk.org/jeps/8329758#Improved-warmup-with-aggressive-heap-expansion, which failed to fetch. JEP 545 is cited instead.
- [ ] Per Liden post: the claims are that JDK 13 introduced SoftMaxHeapSize, the quote 'soft limit on how large the Java heap can grow', ZUncommitDelay defaults to 300s, and memory is uncommitted after a spike. The name on the page is spelt 'Liden', without the accent.
- [ ] Inside Java 2023 (Shimchenko): the claims are guest author and PhD student at Uppsala, an Oracle/Uppsala/KTH project, targets of 5/10/15/20%, higher targets using less memory and more energy, and execution time generally unaffected. 'We infer this is the ancestor of ZGCIntensity' is our link.
- [ ] The Brooker post's date comes from its URL (2024/08/14); the page itself shows no date. The Inside Java lead's date of 19 October 2026 also comes from its URL.
- [ ] 'Integrated in JDK 28' is based on the JEP's Status: Integrated and Release: 28. The JEP does not say whether ZAdaptiveHeapSizing is a product or experimental flag, and the draft does not claim either.
- [ ] The JEP does not say how ZGC measures host memory inside a container with its own limit. The draft says so and leaves it open.
- [ ] The internal links /languages-runtimes/jvm-warmup-and-aot-compilation/ and /distributed-systems/metastable-failures/ both point to verified live Explainers.

**Fact-checker could not confirm:**
- [ ] Paper 'Heap Size Adjustment with CPU Control': its content not read; only confirmed that JEP 546 credits it for smoothing and links the same Semantic Scholar URL
- [ ] Inside Java lead date 19 October 2026 and Brooker date August 2024 come from URLs only; neither page shows a date
- [ ] JEP 546/545 fragment anchors (#Goals, #A-good-neighbor, #Risks-and-Assumptions etc.) built from confirmed heading text, not individually resolved in a browser

<details><summary>What the fact-checker changed</summary>

- Fixed: Misquote: 'as if a soft maximum heap size was set' corrected to JEP 546's verbatim 'as though a soft maximum heap size had been set'
- Fixed: Link text 'just as difficult' (paraphrase not in JEP) changed to 'no easier', JEP's wording
- Fixed: 'a fourth giving memory back when a hungry neighbour arrived' corrected to 'the first' (it was the initial SPECjbb2015 JVM)
- Fixed: Shimchenko results: 'higher targets used less memory and more energy' softened to 'could use more energy' per Source
- Fixed: JEP 545 warmup sentence: 'which JEP 546's loop relies on' (unsupported) replaced with what JEP 546 says: ZGC expands aggressively during warmup
- Fixed: 'We infer' added to unsourced advice on who should raise/lower ZGCIntensity
- Fixed: Missing blank line before '## Why picking a heap size was so hard' heading
- Removed: Hook: 'Every heap size in that chart was the garbage collector's choice' (JEP does not show case-study command lines or say -Xmx was unset)
- Removed: TL;DR and loop section: 'measures its CPU use after each cycle' / 'When a collection cycle ends, ZGC knows what it just spent' (JEP never says measured after each cycle); reworded to 'compares its GC CPU use with a target intensity'

</details>

### Origin Range Requests: how Cloudflare serves ranges from cache chunks

`content/explainers/apis-networking/cloudflare-origin-range-requests.md` · case-study · 2450 words · fact-checker: 64 claims checked, 10 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Inference (Step 1): the blog never gives Pingora's chunk size. We infer it is 1 MiB from the range-requests docs line 'Cloudflare aligns origin range requests to 1 MiB cache boundaries'. The '2 MiB, two chunks' figure is our own arithmetic built on that inference.
- [ ] The Lead was read through WebFetch summaries, not the raw page (Bash curl was denied). Check every paraphrase against the published post, especially the four old-design problems and 'popular files needed a separate origin fetch for every range' (which we left out).
- [ ] Quotes were confirmed verbatim by a targeted fetch but deserve a recheck: 'cut responses short', 'No value works equally well', 'a request that wouldn't fill a floppy disk', 'the extra work is bounded and predictable', 'petabytes less data'.
- [ ] Publication date 23 October 2026 comes from the page's JSON-LD datePublished, not from visible text.
- [ ] Inference (Step 2): before the follow-up change, a partial hit refetched bytes the cache already held, and that is the traffic the 6% measures.
- [ ] Inference (Step 3): when an in-flight download covers only part of a request, the coordinator streams that part and fetches the rest.
- [ ] Inference (Step 4): the aligned fetch runs at each tier-cache hop. The post never places the coordinator in a layer or says it runs at every tier.
- [ ] Inference (evidence): the improvement on cache hits comes from removing slice blocking. The post gives no reason.
- [ ] nginx blog post (https://blog.nginx.org/blog/smart-efficient-byte-range-caching-nginx) timed out twice, so nginx claims cite the nginx slice-module docs instead. The memory warning is paraphrased from the #slice directive text.
- [ ] The Concurrent Streaming Acceleration post's JSON-LD date looked odd (2019 published, 2026 modified), so no date is given in Sources.
- [ ] Docs anchors were inferred from headings (#origin-range-requests, #origin-request-behavior, #partial-response-requirements, #origin-responses-that-do-not-honor-range, #conditional-requests, #multiple-ranges, #caching-limitations, Cache Reserve #limits). Confirm they resolve.
- [ ] The RFC 9110 section anchors (13.1.5 If-Range, 14.4 Content-Range, 15.3.7 206, 15.5.17 416) are from our knowledge of the RFC, not fetched.
- [ ] Internal links: /distributed-systems/hdr-histograms-tail-latency/ and /apis-networking/http-connection-reuse/ (both verified Explainers).

**Fact-checker could not confirm:**
- [ ] scripts/daily/check-page.sh was NOT run: Bash is denied in this session. Page bytes are unknown and must be checked before review. The size limit is 49,500. The draft has 1 diagram and roughly 2,400 words, and the fixes added about 40 words.
- [ ] Cloudflare docs anchors (#origin-range-requests, #origin-request-behavior, #partial-response-requirements, #origin-responses-that-do-not-honor-range, #conditional-requests, #multiple-ranges, #caching-limitations, cache-reserve #limits): each matches a heading on its page, but WebFetch shows no element ids, so we could not confirm the ids themselves.
- [ ] The 23 October 2026 date comes only from the JSON-LD datePublished. The page shows no visible date.
- [ ] RFC 9110 sections: 15.3.7 (206) and 15.5.17 (416) were confirmed from their headings. 13.1.5 (If-Range) and 14.4 (Content-Range) were confirmed only from cross-references in the RFC text.
- [ ] Quotes were checked through WebFetch summaries, which quoted these phrases back verbatim: 'cut responses short', 'No value works equally well', 'a request that wouldn't fill a floppy disk', 'the extra work is bounded and predictable', 'petabytes less data'. Raw HTML was not inspected.

<details><summary>What the fact-checker changed</summary>

- Fixed: Old-design problems: draft said 'four problems' and dropped one. The post lists five, so the text now says five and adds 'Popular files needed a separate origin fetch for every range.'
- Fixed: Concurrent streaming: changed 'a file still being fetched' to 'still being written', the wording used in the Concurrent Streaming Acceleration post and the Origin Range Requests post.
- Fixed: Range sizes: draft said customers who moved had configured 'slices from tens of KB to tens of GB'. The post says 'Configured range sizes on the old system ran from tens of kilobytes to tens of gigabytes', so the text now says that.
- Fixed: Small slices: changed 'each holding memory while the pieces are combined' to 'while it waited', as the post says.
- Fixed: Large slices: changed 'for clients who had already left' to 'after the viewer had already left', as the post says.
- Fixed: Tuning: changed 'let each customer choose' to 'set it per zone'. The post says only 'configured per zone' and never says customers chose the value.
- Fixed: 'No value works equally well' paraphrase: added 'a large media file' and 'every origin and cache topology' to match the full sentence.
- Fixed: Changed 'should not need deep knowledge of the cache' to 'should not need to understand Cloudflare's capabilities in depth ... to have files served efficiently', matching the post's wording.
- Fixed: Changed 'parsing the Range header is the easy part' to 'is straightforward', the post's word.
- Fixed: Mermaid diagram: the 'no' branch claimed that any failed check returns the full representation. The post states that outcome only for an If-Range mismatch and gives no outcome for range, size or ETag mismatches. Edge labels are now 'valid' and 'If-Range mismatch'.

</details>

### Root KSK rollover: how validating resolvers learn a new trust anchor

`content/explainers/apis-networking/dnssec-root-key-rollover.md` · concept · 2591 words · fact-checker: 95 claims checked, 7 fixed, 1 removed

**Drafter asks you to check:**
- [ ] Lead was read through WebFetch summaries and extracted quotes, not raw HTML (Bash curl was denied at first). Check against the page: India 55%→59%, Brazil 54%→50%, global 39.7%→40.6%, 1,700,275,794 of 4,308,517,174, +45.8M net, India +45,423,360, Brazil −7,780,266, 13,538 networks, 35 up >50%, 58 down ≥50%, 50 networks added KSK-2024 between 7 and 13 Oct, and the quotes 'not so fortunate', 'Surprisingly, none of these networks was stranded' and 'The DNS is a highly distributed but loosely coupled system.'
- [ ] Hook says APNIC 'reached a major Indian ISP' on 9 Oct whose resolvers did not trust the incoming key. Confirm the Lead's wording supports 'did not yet trust'. The 9 Oct post lists Reliance Jio (AS55836) at 2.50%, but the draft deliberately does not name it as the contacted ISP.
- [ ] The 'Explanation for the 35 networks' paraphrases the Lead as trusting the incoming KSK 'but not the incumbent'. The extract only gave the clause up to 'trusted the incoming KSK,'. Check the rest of the sentence.
- [ ] The 48-hour cache window comes from the Lead (the shaded 48-hour period) and the 2018 post (the '48-hour TTL'). No source gives 172800 s, so the draft does not use it.
- [ ] 9 Oct post: 'around 90% of reporting resolvers' is Verisign's figure as cited by Huston, and 'around 70% of tests' supporting RFC 8509 is APNIC's. The 6M and 1.8M daily tests are also from this post. The May post gives very different earlier figures (under 20%, 17%→12% in April 2026). Those are left out.
- [ ] Verisign: the 'above 95% after about 500 days' for both the 2017 and 2024 rolls, 3.5% still retaining KSK-2010, the trust anchor file names (bind.keys, root.key for Unbound and PowerDNS Recursor, root.keys for Knot), and the 'confirm automatic updates enabled / writable key storage' advice. The quote 'out of an abundance of caution' is lower-cased from 'Out of…'.
- [ ] IANA dates: KSK-2024 generated in April 2024 (2024-04-26), published in the root zone 2025-01-11, RFC 5011 trust from 2025-02-10, KSK-2017 signing 2018-10-11 to 2026-10-11, revocation 2027-01-11, removal 2027-03-22. These were paraphrased by the fetch tool, not quoted. Note the May post gives the revocation as 'early 2027'.
- [ ] 2018 post: impact 0.2–0.3% of end users against the 0.5% rollback threshold, the halt in 2017 due to trust anchor data, and the quote 'is unable to resolve any name, whether or not the name is DNSSEC-signed'. Note the fetch flagged an internal table vs text inconsistency in that post, which the draft does not cite.
- [ ] RFC 5011 section anchors used: #section-2.3, 2.4.1, 2.4.2, 3, 4, 4.2, 6, 8.2. Confirm the anchors resolve on rfc-editor.org. The same for RFC 8509: #section-2.1, 2.2, 3.
- [ ] 'We infer' claims to check: India's gain being almost the whole global net gain; resolvers holding only the new key failing for 20 months and then working; unwritable key storage reverting to the old file at restart; operators who have never done a roll; why resolver-weighted and user-weighted shares differ; the Sentinel signal guiding the 9 Oct outreach; at least daily RFC 5011 refresh given a 48-hour TTL; a second authoritative provider not helping. The last links the live Explainer /apis-networking/dns-provider-redundancy/.
- [ ] Cloudflare dnstest.dev/ksk-2024 was described from a WebFetch of the page. Confirm it is still up after the roll and the wording 'If all three sentinel names resolve…'.
- [ ] Word count from wc is 2,591 including frontmatter, URLs and Sources. The prose is about 2,250 words. The check reported no section-length warnings.

**Fact-checker could not confirm:**
- [ ] RFC 5011 and RFC 8509 #section-* anchors: WebFetch saw working TOC links, not the raw HTML id attributes
- [ ] Verisign '3.5% of sources still retain KSK-2010': the sentence says only 'About 3.5% still retain KSK-2010'; 'sources' is taken from the preceding sentence
- [ ] Brazil's loss being the largest and India's gain the largest rest on each being the top row of its table in the 21 Oct post, not on a sentence saying so
- [ ] The 5 May 2026 APNIC post is listed in Sources but not cited in the text (title, author and date confirmed)

<details><summary>What the fact-checker changed</summary>

- Fixed: Marked inference: 'By 14:00 UTC on 13 October, we infer, every resolver had fetched the new set' (Source gives only the 48-hour cached-lifetime window)
- Fixed: Replaced unsourced 'The root servers kept answering throughout' with 'Huston reports no major DNS outage' linked to the 21 Oct APNIC post
- Fixed: Marked inference: 'We infer that a stale anchor costs nothing on a normal day'
- Fixed: Marked inference: 'We infer from KSK-2010's 3.5% tail that some will not'
- Fixed: Linked 'KSK-2017 took over on 11 October 2018' to IANA (signing start 2018-10-11)
- Fixed: Linked 'three names' to RFC 8509 #section-4.3 (the bogus / not-ta current / is-ta new procedure is in 4.3, not 3)
- Fixed: RFC 5011 states: 'and finally Revoked' -> 'and Revoked' (Missing and Removed are also states; Removed follows Revoked)
- Removed: 'adding that he did not expect it' after the 'Surprisingly, none of these networks was stranded' quote (not stated beyond 'Surprisingly')

</details>

### Overexposed: how Shopify finds apps that answer without a login

`content/explainers/infrastructure-cloud/shopify-overexposed.md` · case-study · 2671 words · fact-checker: 62 claims checked, 5 fixed, 0 removed

**Drafter asks you to check:**
- [ ] The Source refused verbatim fetch: WebFetch only summarised it, and curl was denied. Every quote was confirmed as verbatim by targeted WebFetch questions, not read off the page directly. Check each quoted phrase against the live page: "a clue, not always an answer", "uncertainty does not mean protected", "not an exposure", "tells you nothing about authentication", "not enough on their own", "a test can pass but check the wrong thing", "the tests are wrong", "runs on a schedule", "One search gets the resources, and a second search gets their access policies without needing cluster credentials.", and the caption "The repository shows intent. The cluster shows what is enforced."
- [ ] Two fetches gave different wordings for the summary line. One gave "what the cluster enforces", the other "what the cluster actually runs". The draft uses "what a team wanted, not what the cluster actually runs" (TL;DR bullet 2 and the closing section). Confirm the exact wording.
- [ ] Status-code table: the two fetches disagreed on the 5xx row. One said "any 5xx", the other "other 5xx". The draft writes 5xx → inconclusive (server error). Also, Cloud Run 200 is shown as "not in the table" because one fetch said the table lists 200 for Ingress only. Verify both against the EXACT_OUTCOMES snippet.
- [ ] The no-status → unreachable and other → inconclusive rows are shown for both surfaces. The Source may list them once, as generic rows.
- [ ] "four steps per run": the post numbers its steps 1–4 (find candidates, probe, classify and apply policy, store). Confirm the count.
- [ ] Gist claims to check: GET fallback only on 405 or 501; DNS answers must be global unicast; allowlist matches lower-cased host and case-sensitive path, applied after classification; inventory "reads configuration only and does not check IAM or prove that an endpoint is reachable"; and the quote "illustrates the core design, not a complete production scanner". The gist states no licence, so no code is quoted.
- [ ] The "Keep reading config" bullet in the options section is our inference, labelled as such. The post names only the config-reading approach it tried.
- [ ] Inferences labelled "we infer": that pinning the checked IP closes a re-resolution gap (the post does not name a DNS-rebinding attack); why allowlisted endpoints are still probed (with a hypothetical status-page example); that coverage beyond two surfaces needs new queries and rules; and that the network restrictions are infrastructure someone must own.
- [ ] Links to the Kubernetes docs are for definitions only: Ingress, and admission controllers intercepting API requests and possibly modifying the resource. The Wikipedia mutation-testing link the post uses was dropped because it is not first-party.
- [ ] Internal link /ai-ml-systems/meta-jittesting/ (verified 2026-10-29) sits in the Tests-can-lie cost paragraph. Confirm it fits the piece.
- [ ] The post links IAP and the URL-map docs. The draft cites these claims to the Shopify post, not to Google's pages.

**Fact-checker could not confirm:**
- [ ] Every quote was confirmed through targeted WebFetch questions (a summarising model), not read off the raw HTML, because Bash and curl were denied. Each answer agreed with the draft's wording, including 'what a team wanted, not what the cluster actually runs'
- [ ] The Google Cloud ingress-doc anchor #settings was inferred from in-page links that point to it. The element ID itself was not seen

<details><summary>What the fact-checker changed</summary>

- Fixed: Four steps per run: replaced our paraphrase with the post's own step names (ask the cloud what is running, probe like a stranger, interpret the result with policy, record a verdict and evidence)
- Fixed: Gist claim 'reads configuration only and does not check IAM or prove that an endpoint is reachable' is not in the gist. Replaced it with the gist's verbatim lines: 'reads no IAM and trusts no auth annotation' and inventory 'is not proof of public reachability'
- Fixed: 'It can cover the whole organisation from one place' was unsupported. Replaced it with the post's snippet, which takes an organisation, folder or project as its scope
- Fixed: Status table intro: the no-status, 5xx and anything-else rows come from the post's classify_status fallbacks, not from EXACT_OUTCOMES. The intro now says so and notes that they apply to both surfaces. 5xx is 'any 5xx' (500 <= status < 600); Cloud Run 200 is correctly shown as not in the table
- Fixed: 403 reason: 'reasons that have nothing to do with access control' changed to the Source's 'for other reasons'

</details>

### Workers AI KV cache exposure: how a rejected prompt read stale pages

`content/explainers/ai-ml-systems/workers-ai-kv-cache-exposure.md` · case-study · 2300 words · fact-checker: 62 claims checked, 16 fixed, 0 removed

**Drafter asks you to check:**
- [ ] The Lead's text came through WebFetch summaries, not the raw page (Bash was denied, so no curl). Check every claim against the live post, especially the short quotes: "HTTP 200 instead of HTTP 400", "We were aware of this mismatch before launch", "still held what earlier requests left in them", "We found no evidence of deliberate exploitation".
- [ ] The Cloudflare post's headings have no anchor ids, so every link to the Lead goes to the bare URL.
- [ ] The 'defences' section is our framing. The Lead does not list alternatives it rejected; the six options come from its fixes, its operator advice and its open work. 'The deciding argument is depth' is labelled as inference only in its last sentence. Check the framing is fair.
- [ ] Version claims: WebFetch summaries of the Lead were inconsistent about where #36988 shipped (v0.5.20, or merged 6 Sept and unreleased). The draft avoids naming a 'fixed in' version and says only 'v0.5.22 or later carries both'. Confirm.
- [ ] #38935: merged 16 Sept 2026, author ShangmingCai. The '300-second timeout to an immediate 400' figure is from the PR's 1P1D test with both workers rejecting. The PR says it does not cover the case where only prefill rejects; the draft leaves that caveat out.
- [ ] Issue #30233: opened 6 July 2026 by kevin-14 on GLM-5.2-FP8, closed by a bot for inactivity on 5 Sept 2026. The draft says 'three days before Cloudflare's report arrived' (report received 8 Sept, 23:11 UTC). Check the close date on GitHub.
- [ ] The claim that responses could contain 'other users' prompts and partial generations' came from a WebFetch summary. Check the Lead's exact wording.
- [ ] The page-tag prototype: 'tested on GLM 5.2 in mid-August, but disabled it shortly afterwards due to reliability issues'. Confirm 'it' means the prototype, not the model (WebFetch said the prototype).
- [ ] The SGLang post's affected range for the RDMA race: v0.4.x through v0.5.21 (v0.5.22 for MORI); the draft simplifies this to 'v0.4.x through v0.5.21'. The 60 to 120 s timeout advice is for v0.5.19 and v0.5.20 set by hand.
- [ ] Inferences labelled 'we infer': why decode reserves pages early (RDMA needs room ready), how stale K/V becomes another customer's text, that the larger pool cost something elsewhere, and that held pages shrink room under abort-heavy load.
- [ ] Topaz and Meerkat links were taken from the Lead's own links but not opened. Confirm the URLs resolve.
- [ ] The hook says 'From 27 August to 9 September 2026'. The Lead dates the exposure from launch on 27 Aug until the model went private on 9 Sept at 14:00 UTC.

**Fact-checker could not confirm:**
- [ ] Exact quotes confirmed only through WebFetch extractions (Bash curl denied), not raw HTML: 'HTTP 200 instead of HTTP 400', 'We were aware of this mismatch before launch', 'still held what earlier requests left in them', 'We found no evidence of deliberate exploitation', 'a small number of customers'. All matched consistently across several fetches.
- [ ] Cloudflare's PR-to-release table puts #36988 in v0.5.20, but the post's prose also says it was 'merged September 6 … but not yet released'. The draft's 'v0.5.22 or later carries both' matches the post's own recommendation, so it was kept.
- [ ] LMSYS 'two linked problems' is our paraphrase in the link text (the post links them through 'an admission rejection … aborts request A mid-transfer'); 'spends most of its length on the race' is based on WebFetch's estimate of section lengths.

<details><summary>What the fact-checker changed</summary>

- Fixed: Step 2: '#38935 fixed it' → 'closed the admission gap'; 'went through' → 'could complete a bootstrap handshake, reserve a metadata buffer and run a forward pass' (PR calls itself a waste-and-latency fix, not a correctness fix)
- Fixed: RDMA race range: 'v0.4.x through v0.5.21, silently' → 'from the v0.4.x line until v0.5.21 (v0.5.22 for MORI), and each occurrence was silent' (LMSYS wording)
- Fixed: 60–120 s timeout advice scoped to v0.5.19 or v0.5.20, not 'older versions'
- Fixed: LMSYS 'no CPU signal' paraphrase → bypasses receiver's CPU, decode gets 'no signal for each write' it could wait on
- Fixed: Not zeroing: 'common practice for performance' / 'needed for performance' → Cloudflare's 'common among inference engines' and 'necessary for the best performance', linked
- Fixed: Zeroing: 'stale reads return nothing' → PR #46435's 'zeros instead of stale data'; noted it is an opt-in draft
- Fixed: Timeline: 25-minute reproduction attributed to Product Security, two independent accounts, as the post says
- Fixed: #38935 'follows Cloudflare's draft patch' now linked to the Cloudflare post (the PR page does not say it); added 'with both workers rejecting' to the 300 s → 400 test
- Fixed: Page tags: 'a read … fails that request only' (stated as prototype behaviour) → the post's stated aim; 'attention backends' → 'attention backends and model architectures'; checksum comparison linked and phrased as the post's 'check every read'
- Fixed: Formal verification: 'track where production SGLang diverges from it' → 'track and fix any divergence from the safety properties it wants'
- Fixed: 'Different model' for the 1% figure → 'tested on GLM-5.2, not GLM-5.3', linked
- Fixed: 'Other users' prompts and partial generations' now attributed to stale pages and linked
- Fixed: Quote punctuation: period moved outside "still held what earlier requests left in them"
- Fixed: Defences list: missing links added for proxy advice, backstop, integrity checks; backstop wording matched to the post ('the day the pool-size check is wrong')
- Fixed: 'The deciding argument is depth' labelled as our reading
- Fixed: Sources: real titles for #46435, Topaz and Meerkat posts

</details>

### Not drafted

- none

<details><summary>Leads the editor rejected</summary>

- `covxn3zy`: Brief: GitHub status page, Copilot agent-sessions errors; no cause given, RCA only promised.
- `b6h3yjql`: Brief: status page blames an unnamed upstream model provider; no mechanism, no live Explainer teaches it.
- `vsfmm6jz`: Brief: GCP europe-west2-b power/PLC incident first posted 2026-10-20 (same page as 55sjwtvu), outside 72h; no live Explainer teaches power or PLC failover.
- `55sjwtvu`: Same GCP incident as vsfmm6jz; too old for a Brief, and a facility power event is outside the Verticals' Explainer scope.
- `grxwqfpc`: Brief/Explainer: inside.java teaser with no changes listed; no live GC Explainer to link.
- `rt7w4cwg`: Brief: Node 26.13.0 only marks the line LTS; no code changes.
- `6qhqbqua`: Brief: Node 24.22.2 point release, V8 backports and one stream fix, no stated security mechanism.
- `fmyk4h44`: Brief: quarterly roundup of 49 disruptions (shutdowns, power, quakes); no single system mechanism a live Explainer teaches.
- `mm72sfrp`: Brief: AWS Health Dashboard root URL, no cause.
- `hf74t2u5`: Brief: DDoS threat report is statistics, not an event at a named system.
- `uzpknzyx`: Brief: AWS customer story (BMW on Graviton); no live Explainer teaches its lesson.
- `2sekcmoe`: Brief: RFC publication, not an event at a system; no live Explainer.
- `7hjwr6fi`: Brief: Kubernetes patch releases (also hlu4lise, jpvyzfyk, y2d3zdxr), bug fixes only.
- `baainlg7`: Brief: PostgreSQL news product announcements and CfP (also i7tzhiny, rudhkots, uynhsn6g).
- `zj552qbr`: Community/working-group post, not an event or architecture.
- `tjxmyhel`: Measurement research, not an event; weak as an Explainer.
- `4pkal76v`: Explainer: GitHub star-history endpoint is solid, but the source carries ~1,000–1,500 words, not 2,000.
- `rdsbakwp`: Explainer: Meta World Cup post is a ~1,000-word overview of many systems with no internals.
- `lvup6ddk`: Fly.io London BGP outage is a ~150-word infra-log note; too thin.
- `msztzyde`: DuckDB JSON-to-typed-columns is a SQL how-to on three rows; no architectural idea big enough.
- `w7dafiay`: DynamoDB Bulk Executor is mostly a tool walkthrough and the fetch was truncated; the Vertical slot went to stronger picks.
- `3lcwadch`: Netflix Pensive AI: netflixtechblog/medium returns 403, couldn't confirm depth.
- `ftnyh3pn`: Netflix ads event pipeline Part II: 403 on fetch, couldn't confirm depth.
- `3uywsrjr`: Cloudflare eMudhra post used as a supporting Source for bundled-vs-platform-trust-stores (lead xpqwuoig).
- `h3yyvfmm`: OpenJDK eMudhra teaser; its openjdk.org page is a supporting Source for bundled-vs-platform-trust-stores.
- `pk4mvo2l`: Already covered by live netflix-single-tenant-titus.
- `fdlpxa74`: PCP/DCP already covered by live decode-context-parallelism.
- `jtyoqdlc`: Blue/green MariaDB upgrade: customer story, largely covered by planned-database-switchover.
- `dm3iiss5`: NYT Oracle migration: customer-success story with little mechanism.

</details>


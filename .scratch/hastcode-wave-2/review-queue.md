# Wave 2 review queue

Drafted 2026-10-29 by agents: 11 Explainers, every one `draft: true`, none verified. Each was fact-checked by a second agent against its Sources. Verify every claim yourself before removing `draft: true` and adding `verified:` (ADR-0007). Publish at most two a day (ADR-0007 as amended); after the 21 published on 28 October, the plan is no new Explainers before about 8 November.

Renamed after drafting: the vLLM Mistral Large 4 Lead became the Concept `decode-context-parallelism`; the PyTorch 2.15 Lead became `fully-sharded-data-parallelism`.

## Meta Andromeda: how Meta scores every eligible ad on one GPU server

`content/explainers/ai-ml-systems/meta-andromeda-ads-retrieval.md` · 2720 words · fact-checker: 88 claims checked, 16 fixed, 0 removed

**Drafter asks you to check:**
- [ ] The Lead was read as published: I downloaded its HTML and extracted the text, because WebFetch would only summarise it. Every number comes from that text. The Lead's headings have no anchors, so all links point to the bare post URL.
- [ ] 5.3x from 2.4 billion works out to about 12.7 billion ads/s. The post only says 'over 10 billion', so this is consistent but worth a look.
- [ ] The 2.4% conversion lift ('Since 2025, Meta Andromeda has delivered...'): the post does not say which generation it covers. The draft says so.
- [ ] The post says both 'on the same GPU' and 'on a single GPU server'. The draft flags that we cannot tell whether 1 billion ads/s is one GPU's figure or one server's.
- [ ] The diagram places the user experts outside the fused kernel. The post does not say where they are computed, and the draft says so.
- [ ] Inferences are labelled in their sections: why retrieval quality matters ('an ad that retrieval drops never reaches the ranking model'), the reason for several user towers, that 2024 system work ran as separate kernels through HBM, that 'index' means the targeting index, that exact targeting contrasts with approximate scoring, the hardware-binding cost, the capacity cost (each server holding the index and embeddings), and the shift from hierarchical-index pruning to scoring every eligible ad.
- [ ] The 1 KB intermediate per ad at 10B ads/s = 10 TB/s example is labelled as our numbers, not Meta's.
- [ ] MTIA: the 2024 post names MTIA beside Grace Hopper and plans further integration with it. The draft says the 2026 post mentions only NVIDIA hardware. Check that wording.
- [ ] The 'options weighed' section is built from the three earlier designs the post says it replaced, because the post lists no rejected alternatives. It adds the 2024 post's comparison with two-tower networks and approximate nearest neighbor search. The section says this framing is ours.
- [ ] The gloss of index matching as 'which ads this user is eligible to see at all' goes beyond the Source's words ('enforcing advertisers' targeting rules').
- [ ] Description is 135 characters and the title 67. related: [expert-parallelism] is now set in the frontmatter. draft: true is kept and verified: is not set.

**Fact-checker could not confirm:**
- [ ] The 2026 post has no heading anchors, so every link points at the bare post URL. The Author may want more specific anchors, but the page offers none.
- [ ] The 1 KB per ad at 10B ads/s = 10 TB/s example is our own arithmetic, labelled as ours. The arithmetic is correct, but no Source states it.
- [ ] The diagram's placement of the user experts (outside the kernel) and the ad tower (inside) is our choice. The caption says the post does not say where either is computed.
- [ ] The 5.3x from 2.4B is consistent with the post's 'over 10 billion' (5.3 x 2.4 is about 12.7B). The figures are reported as the post gives them.

<details><summary>What the fact-checker changed</summary>

- Fixed: Opening: replaced the unsupported 'the biggest win came from moving data less' with what the post supports: the pipeline moves far less data through memory and across the network (fused kernel, about 100x less network traffic).
- Fixed: TL;DR: 'the richer model scores as fast as a plain dot product' changed to 'the optimized gating kernel runs as fast as a simple dot-product scoring kernel'. The post says this about the optimized gating kernel.
- Fixed: TL;DR: linked the bullet about the request reaching one GPU server, which had no link.
- Fixed: Diagram: 'Gather the ad's embeddings' changed to 'Gather embeddings', because the post does not say whose embeddings. The caption now also says we placed the ad tower inside the kernel and that the post does not say where it is computed.
- Fixed: GH200 claim: the link moved from the 2024 post, which says only 'NVIDIA Grace Hopper Superchip', to the 2026 post, which names the GH200.
- Fixed: 'Three bottlenecks, none of them arithmetic' changed to two bottlenecks (memory traffic, network) plus a third pressure, cost, from another section of the post. Linked the unlinked 'tens of millions of ads per second per GPU card'.
- Fixed: 'Designs Meta moved past': the intro now says the framing is ours. 'Each old design was a boundary... cost bandwidth' was wrong for the single user tower; it now says, as our reading, that the last two were boundaries.
- Fixed: Index matching: the gloss 'which ads this user is eligible to see at all' is now labelled 'We read that as...'. 'Meta is plain that this was the hard part' changed to the post's 'not straightforward'.
- Fixed: Ad tower: the gloss 'depends on where and when it would be shown' replaced with the post's 'context-aware representations'.
- Fixed: End-to-end: 'shared layouts' corrected to 'co-designing data layouts', and 'latency' to 'end-to-end latency'.
- Fixed: Grace Blackwell: 'already doing it' (redoing kernel work) changed to 'already extending its co-design work', as the post words it.
- Fixed: 'The table above is not yet what serves ads' softened: the production fused kernel already runs at over 10B ads/s, so now 'not everything in the table serves ads yet, and the post does not say which parts do'.
- Fixed: 2024 results: '+6% recall and +8% ads quality on selected segments' now has a comma, because the 'selected segments' qualifier applies only to +8%.
- Fixed: The quote 'Near ranking-level performance' was not verbatim. Corrected to the post's 'near-ranking-level performance'.
- Fixed: MTIA: 'planned further integration with MTIA' changed to what the 2024 post says: integrating Andromeda with MTIA would push retrieval further. 'The 2026 post mentions only NVIDIA hardware' changed to 'does not name MTIA, and the only hardware it names is NVIDIA's'. The 2026 post does refer to 'custom AI hardware' when it describes 2024.
- Fixed: Closing: the unlabelled inference about model and systems teams moving together is now marked 'we infer'.

</details>

## Lens as data: how Smarsh automates AWS Agentic AI Lens reviews

`content/explainers/ai-ml-systems/aws-agentic-ai-lens-reviews.md` · 2440 words · fact-checker: 72 claims checked, 9 fixed, 0 removed

**Drafter asks you to check:**
- [ ] The paragraph on a fourth option is ours. It rests on our own count on 30 October 2026: aws-samples/sample-well-architected-custom-lens agentic-ai-lens.json held 40 questions and 146 best practices (186 choices minus 40 'None of these'). The appendix lists 41 questions and 150 practices, and the file lacks AGENTOPS07. Re-count before publishing, because the file may change. We do not know whether the file existed when Smarsh started, and the draft says so.
- [ ] The 90 High / 59 Medium / 1 Low split is in the post. We also re-counted it from Appendix A and named the one Low item, AGENTPERF06-BP03 'Optimize meta-tool utilization and tool chaining'. Check both.
- [ ] The JSON block is our placeholder version of the contract's shape, not the post's example, so the post's text is not copied. Confirm the four field names (status, evidence, gap, remediation).
- [ ] The hook says the script 'refuses any best practice the model makes up'. That glosses the post's 'unknown best-practice IDs or invalid statuses abort the build ... hallucinated items cannot silently ship'. Check that the wording holds.
- [ ] The second bullet under the post's options ('No effort, but no review either, and Smarsh's gate does not let an agent ship without one') is our gloss, combining 'waiting for console support' with the rule that every agentic solution must pass a review before it reaches a customer.
- [ ] 'When Claude calls it, it also validates row counts and flags blank statuses': the post puts this in its workflow step, not in its description of build_report.py, so it is unclear which component does it.
- [ ] Step 1 numbers come from Anthropic's Agent Skills overview (platform.claude.com): about 100 tokens per skill for metadata, and 'only their output enters context' for scripts. docs.claude.com redirects there, so check the anchors #level-1-metadata-always-loaded and #level-3-resources-and-code-loaded-as-needed.
- [ ] The lens publication date of 10 June 2026 and 'initial release' come from the lens's document-revisions page.
- [ ] Every 'we infer' sentence is our analysis. These include: consistency is the main cost of manual reviews; the reason the baseline produced short checklists; what read-only access is for; how one fact bears on several pillars; what live COUNTIF formulas are for; and that the dataset needs re-extracting when the lens is revised.
- [ ] 'What it costs' says the trade-offs are 'mostly by implication' in the post. Its four costs are our framing, built on sourced facts: the dependency-scanning miss, the questions that went to AWS, and 'humans on judgment'.
- [ ] The post's conclusion quote is credited to 'the post', because it is co-authored by Smarsh and AWS. The executive-summary step and the AWS agenda items are paraphrased from the post.

**Fact-checker could not confirm:**
- [ ] The custom-lens JSON count is only as of today: re-counted on 2026-10-30 at 40 questions and 186 choices, of which 40 are 'None of these', leaving 146 best practices, with no AGENTOPS07. The Author should re-count before publishing because the file may change. Whether it existed when Smarsh started is unknown, and the draft says so.

<details><summary>What the fact-checker changed</summary>

- Fixed: Hook: 'a Smarsh architect turned the lens into a JSON file' changed to 'Smarsh turned the lens...'. The post says Smarsh built the skill and does not say who extracted the data, which the draft itself admits in Step 2. Also changed 'refuses any best practice the model makes up' to 'rejects any best-practice ID the model makes up', because the validator checks IDs and statuses.
- Fixed: TL;DR: 'without the skill, the same model passed 8 of 16' changed to 'baseline runs without the skill passed 8 of 16'. The post never says the baseline used the same model.
- Fixed: Diagram intro: 'The model works only inside the assess box' changed to 'The model's judgement happens only inside the assess box'. Claude also gathers the evidence and calls build_report.py.
- Fixed: Smarsh customers: 'many of them financial institutions' changed to the post's 'from small businesses to large multinational financial institutions'. The post gives no share.
- Fixed: Options bullet 2: replaced the unsourced 'No effort' gloss with 'No review until AWS ships it, and Smarsh's gate does not let an agentic solution reach a customer without one'. The gate is sourced: every agentic solution must pass a review before it reaches an end customer.
- Fixed: Mermaid: the executive summary moved out of build_report's output into its own node, 'Skill returns an executive summary of top high-risk gaps', to match the post's workflow (step 6 comes after the script).
- Fixed: Step 5: the sentence giving the script row-count validation now says the post puts this in the workflow step and does not say whether the script or the model does it.
- Fixed: Step 6: 'evidence the repository had overlooked' was copied near-verbatim from the post without quote marks. It now reads 'evidence the review had missed'.
- Fixed: Evidence section: 'Most expectations check completeness and form' now reads 'Of the expectations the post names, most check...', so the claim rests on the post's own list.

</details>

## Immutable data in Java: state the JVM can trust never to change

`content/explainers/languages-runtimes/java-immutable-data.md` · 2440 words · fact-checker: 78 claims checked, 8 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Release status as of 30 Oct 2026: JEP 401 and JEP 539 show Status Integrated, Release 28, and both are on the JDK 28 page's targeted list (updated 22 Oct). JEP 547 and JEP 468 are still Candidate. The draft says 547 is 'proposed for a fourth round in JDK 28' and calls 468 'still a candidate'.
- [ ] Two links point at JEP 401's anchor 'When-can-flattening-and-scalarization-be-done?', written as %3F because a bare '?' would start a query string. Check the browser lands on that section. It supports the claims that flattening typically cannot apply to Object or T, and that the fallback is most common during warmup.
- [ ] The hook quotes JEP 401's '48 bits' and the jshell results d1 == d3 false (before) and true (with --enable-preview), taken from the Motivation and Programming-with-value-objects sections.
- [ ] 'Since JDK 15 and 16' for hidden classes and records being barred from final-field mutation comes from JEP 500's Motivation. 'Static finals cannot be updated by reflection' paraphrases JEP 547's Risks section.
- [ ] The value-based history (JDK 8 for Optional and java.time, JDK 16 for the primitive wrappers) comes from Parlog's heads-up, not from a JEP.
- [ ] The fixes for record constructors ('remove the statement, make it static-like, or assign every field then call super()') paraphrase the heads-up. 'Static-like' stands in for its 'very similar to a static context'.
- [ ] Practice three ends on a 'we infer': value-object scalarization can often keep derived copies out of the heap. No Source states this for 'with' expressions.
- [ ] Author credits use the JEP pages' Author or Owner lines. JEP 395, 401 and 539 list only an Owner (Gavin Bierman for 395, Dan Smith for 401 and 539), and JEP 269's Owner is Stuart Marks.
- [ ] All Java snippets are short syntax fragments we wrote (for example 'r with { denom = 0; }' and 'dates[0] = dates[0].plusYears(30)'), modelled on the JEP examples. No JEP code block is reproduced.

**Fact-checker could not confirm:**
- [ ] 'Number and Record as abstract value classes': JEP 401 describes Number as an abstract value class and lists Record among the 30 value classes. It never calls Record abstract. java.lang.Record is an abstract class in the JDK API, so we kept the wording, but the cited JEP does not state it.
- [ ] The '%3F' links to JEP 401's 'When-can-flattening-and-scalarization-be-done?' anchor: the page's HTML has id="When-can-flattening-and-scalarization-be-done?" and the text there supports both claims (Object/T and warmup). Browsers decode the percent-encoding of a fragment when they look for a matching id, so this should land on the section, but nobody has clicked it in a browser to check.

<details><summary>What the fact-checker changed</summary>

- Fixed: JEP 395 quote was not verbatim: the Summary says 'transparent carriers for immutable data' (plural). The sentence now reads 'JEP 395 defines records as classes that act as "transparent carriers for immutable data"', and the next sentence starts 'A record's component fields'.
- Fixed: Flattening atomicity: JEP 401 gives a hard limit, 'no more than 64 bits', not 'about 64 bits'. Changed to 'limits mutable fields that store them to no more than 64 bits'.
- Fixed: 'Constant folding across them is "not only possible but routine"': JEP 547's Risks section says this about static final fields only, not about record or hidden-class fields. Changed 'them' to 'static finals'.
- Fixed: The sun.misc.Unsafe sentence had no Source. JEP 500's Motivation says JDK 24 began removing Unsafe methods that can mutate final fields, so the sentence is now linked to that and the memory-access part goes to our off-heap Explainer: 'removing methods of sun.misc.Unsafe that can also mutate final fields; our off-heap memory Explainer follows its memory-access methods out.'
- Fixed: The claim that non-record value classes need writeReplace/readResolve to serialize had no link. It now links to JEP 401's #serialization anchor.
- Fixed: Escape analysis versus value objects: JEP 401 says scalarization of value objects is 'more predictable and far-reaching', not that it makes it 'predictable'. Changed to 'more predictable'.
- Fixed: JEP 269 says the behaviour of a Set whose element is mutated 'could become unspecified'. Changed 'becomes unspecified' to 'can become unspecified'.
- Fixed: Typo: '; Existing classes' changed to '. Existing classes'.

</details>

## DuckLake: a lakehouse format that keeps its metadata in SQL tables

`content/explainers/databases-storage/ducklake-lakehouse-format.md` · 2456 words · fact-checker: 96 claims checked, 10 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Hook figures are Matt Martin's own reports in the Lead (#why-ducklake): Iceberg catalog operations of seconds resolving in under a millisecond in DuckLake 'in his experience'. They are anecdote, not a measurement; the draft frames them that way.
- [ ] Manifesto INSERT walk-through (#ducklake): the draft gives only the values 42 and 43, snapshot 2 and 'inserted_into_table:1'. It leaves out the numeric columns of the v0.1 INSERT rows (2 rows, 279 bytes, 164 footer) because the v0.1 column order may differ from spec 1.1.
- [ ] Manifesto date 'May 2025' is linked to the DuckLake 1.0 post, which says 'In May 2025, we published the DuckLake manifesto'. The news index dates the original post 2025-05-27. The manifesto page itself shows no date.
- [ ] Benchmark table (10,000 rows, 1,000 batches; 1,148.77 s vs 10.88 s, 83.06 s vs 0.09 s, 52.83 s vs 0.28 s) and the 401-file pyiceberg count (100 Parquet, 200 Avro, 101 JSON) come from the DuckDB team's own inlining post. It is a vendor benchmark of its own format, and the draft says so. The post's TL;DR gives 926x for aggregation, while its Iceberg table gives 923x; the draft cites neither multiplier, only 105x for inserts.
- [ ] 'Thousands of transactions per second' for PostgreSQL and 'as many table changes as the catalog database can commit' are the manifesto's claims (#speed), not measured; the draft notes there is no independent measurement.
- [ ] The '1.5 million downloads' figure: the 1.1 post says DuckLake had more than 1.5 million downloads 'last week (October 12 to 18, 2026)'. The draft attributes it to the extension for that week; check the wording.
- [ ] The 'we infer' sentences: DuckLake's primary-key conflict check as the same optimistic scheme as Iceberg's pointer swap; and the catalog's primary keys guarding metadata, not user rows.
- [ ] Code snippet (INSTALL ducklake; ATTACH 'ducklake:my_ducklake.ducklake' ...) is from the DuckLake docs, repo duckdb/ducklake-web, MIT licensed, 'Copyright 2018-2025 Stichting DuckDB Foundation'. The notice is credited in Sources.
- [ ] Spec version is 1.1 (specification/introduction). DuckLake 1.1 requires DuckDB v2.0.0 or later (duckdb/introduction#installation). Both are current as of 2026-10-30 and may change.
- [ ] Deletion vectors in Puffin files are marked experimental in the DuckLake 1.0 post (#deletion-vectors); the draft says so.
- [ ] Docs URLs are cited without .html and return HTTP 200 directly; anchors were taken from the page headings (e.g. conflict_resolution#data, #compaction; data_inlining#insertion-inlining).

**Fact-checker could not confirm:**
- [ ] 'Mark Raasveldt and Hannes Mühleisen of DuckDB': the manifesto names both authors but does not state 'of DuckDB'. Left as is.
- [ ] 'Swap the path for a PostgreSQL connection string and an S3 DATA_PATH, and the same tables serve a team': this is a general reading of the docs and has no specific linked passage. Left as is.
- [ ] The SQL snippet joins two separate docs code blocks: 'INSTALL ducklake;' from Installation and the ATTACH line from Creating a New Database. Both lines are verbatim from the docs, and the MIT notice (Copyright 2018-2025 Stichting DuckDB Foundation) is confirmed in the duckdb/ducklake-web LICENSE.

<details><summary>What the fact-checker changed</summary>

- Fixed: Lead: the sentence 'a query that used to plan in seconds can slow down significantly because thousands of metadata files need to be resolved first' copied the Source's words without quote marks. It is now paraphrased.
- Fixed: TL;DR: the claim that a read sends one query that prunes by schema, partition and statistics had no link. It now links to the manifesto's #speed section, which says this.
- Fixed: Visibility rule (begin_snapshot/end_snapshot filter): added a link to specification/queries#select, where the condition is given.
- Fixed: Deletes: 'Parquet file of row positions' became 'deleted row indexes', since the spec says 'row index'. The 'at most one delete file per data file in a snapshot' rule now links to queries#select, where it appears; it was not on #delete.
- Fixed: Train-stations example: 'AT (VERSION => 1) returns the old name' became 'reads it as of snapshot 1, before the rename'. The docs show that query but not its output.
- Fixed: Expire snapshots: added a link to the expire_snapshots page for the ducklake_expire_snapshots options.
- Fixed: Cleanup: 'generally safe after a few days' became 'generally safe for files scheduled more than a few days ago', to match the docs.
- Fixed: Tiered compaction: 'larger ones rarely' became 'larger ones less often'. The docs run tier 1 occasionally and tier 2 rarely.
- Fixed: Iceberg compatibility: reworded to 'the data and positional delete files DuckLake writes', to match the manifesto's 'data and (positional) deletion files that DuckLake writes'.
- Fixed: Downloads: 'the extension had more than 1.5 million downloads from 12 to 18 October 2026' became 'DuckLake had ... in the week of 12 to 18 October 2026'. The 1.1 post says DuckLake and does not name the extension.

</details>

## Resilience Hub tests: how AWS turns a service model into FIS faults

`content/explainers/infrastructure-cloud/aws-resilience-hub-testing.md` · 2330 words · fact-checker: 96 claims checked, 9 fixed, 3 removed

**Drafter asks you to check:**
- [ ] The Sources disagree on whether a payment gateway is soft or hard. The Lead's worked example treats Payment Gateway as soft and its run passes, while the guide's classification page lists 'payment gateway' as a typical hard dependency. The draft reports this as a disagreement in the evidence section; check the wording is fair.
- [ ] The Sources also disagree on reports. The Lead says 'Each test run auto-generates a report', while the runs-and-report page says no report is generated without an S3 destination. The draft states both.
- [ ] 'We infer' that 'selected based on your service's architecture and resilience policy' mostly means pre-filling parameters and applying the policy's thresholds. Neither the Lead nor the guide explains how selection works; this inference is ours.
- [ ] 'We infer' that a stale hand-built experiment can keep passing because its target has moved and nothing is hit. This illustrates the Lead's 'drift out of date' point and is not stated by AWS.
- [ ] 'We infer' that ECS task replacement lets the fleet heal itself out of the experiment. Built from the docs' note that packet loss does not re-apply to replacement tasks.
- [ ] The $3 example (one action, 30 minutes, at $0.10 per action-minute) is labelled as our arithmetic. The draft also says we have not found whether a skipped action is billed. Check the pricing page has not changed.
- [ ] Dependency validation limits come from the Lead, not the template docs: up to 10 dependencies, 1 minute to 12 hours, default 30 minutes. Multi-Region isolation's 3-hour default and the recovery tests' 'RTO plus 30 minutes' default come from the configuring page.
- [ ] The AZ recovery action list is condensed from the guide's table. Confirm the condensed wording is still accurate: stops instances, blocks launches and Auto Scaling capacity, cuts subnet traffic, RDS failover if the writer is in the AZ, ElastiCache node termination, EBS I/O pause on half the volumes, ARC zonal autoshift if enabled.
- [ ] Most AWS docs pages were read through curl and HTML-to-text extraction, not as rendered, so recheck the quotes against the live pages. The quotes are: 'Do not use it as a guarantee or as evidence. If you need evidence, run a resilience test.' and 'no cascading failures, no data corruption, no hung connections'.
- [ ] The diagram is ours. The Lead's two figures (test-run timelines) are linked, not reproduced. Check the Mermaid renders on the Preview.
- [ ] Inline link to /distributed-systems/metastable-failures/ (related: [metastable-failures]). It is tied to the decision framework's advice to add retry or fallback; confirm the pairing is apt.

**Fact-checker could not confirm:**
- [ ] No Source states whether a skipped action is billed. The draft already says we have not found this.
- [ ] The absence claim (no figures on customer adoption or failure rates) can only be checked by noting that the post and guide pages read contain none.
- [ ] The 'cascading failures' quote and the decision-framework wording were checked against the raw HTML of the blog. WebFetch's rendered view cut off before that section. The quote matches the HTML word for word.
- [ ] Mermaid was checked only in the built Preview HTML: it is a figure.diagram with inline SVG containing the new labels. It was not looked at in a browser.

<details><summary>What the fact-checker changed</summary>

- Fixed: Diagram node 'one per AZ, Region or dependency fault' was wrong because there are two Region templates. It now reads 'AZ, Region and dependency faults'.
- Fixed: Diagram result node now reads 'report if an S3 destination is set'. The runs-and-report and configuring pages say no report is made without an S3 destination.
- Fixed: 'AWS does not let either stand as proof' became 'AWS does not let the assessment stand as proof'. The Source says this only about the next-gen assessment, not about v1's checks.
- Fixed: 'Still the tool for a fault no template covers' is no Source's claim. It is now marked 'which, we infer, stays the tool...'.
- Fixed: Decision framework: 'means a misclassification...' / 'means the service is more resilient...' overstated the post's 'might be misclassified'. Both now say 'may be misclassified'.
- Fixed: 'That is the framework's point...' is our reading, so it now starts 'We read that as the framework's point'.
- Fixed: The report disagreement is now described fairly. The post's opening says each run auto-generates a report, but its own configuration walkthrough later agrees with the guide that a report needs an S3 bucket.
- Fixed: Added Source links to two claims that had none: 'Multi-Region RTO the policy must define' (Multi-Region recovery page) and 'add one by DNS name' (dependency validation page).
- Fixed: Trimmed four bits of filler (diagram reading guide, 'Every input...' lead-in, evidence opener, two lead-in sentences) to get the page from 49,777 to 49,440 bytes.
- Removed: 'Recommended tests are no different.' (filler, cut for size)
- Removed: 'AWS FIS already gave teams the raw parts.' (filler, cut for size)
- Removed: The diagram walkthrough sentence listing the four stages, replaced by 'It reads top to bottom; Steps 1 to 5 below walk it in order.'

</details>

## Revinate's partner API: how API Gateway reaches EKS with no ingress

`content/explainers/apis-networking/revinate-api-gateway.md` · 2416 words · fact-checker: 74 claims checked, 7 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Viability: the post has no versioning scheme, no throttle/quota values and no traffic, latency or cost numbers (it calls its results 'qualitative'). Decide whether this meets the brief's bar before verifying the rest.
- [ ] The post's TargetGroupBinding option lists both 'shared NLB across multiple Target Groups' (pro) and 'different applications will need individual NLB' (con). The draft quotes both, and our reading (one NLB carrying several target groups per application) is labelled as ours.
- [ ] Not stated in the post, and flagged as such in the draft: where the WAF web ACL is attached, the API endpoint type (edge-optimized or Regional behind its own CloudFront), the VPC link version (V1 legacy, NLB only, or V2), whether the authorizer caches results, and the API key source (HEADER or AUTHORIZER). The authorizer-source preference is labelled 'we infer'.
- [ ] Inferences labelled 'we infer': Terraform owns the NLB and target groups while TargetGroupBinding manages only the pods; path routing lives in API Gateway because the NLB forwards by port; usage-plan attribution can measure but not hard-cap spend; the portal's 'Try it' calls go through the same Lambda authorizer.
- [ ] Account default throttle of 10,000 RPS with a 5,000 burst is from the API Gateway quotas page (most Regions; some default to 2,500/1,250). Confirm it is still current.
- [ ] 'Weeks to days' onboarding is Revinate's figure; 'weeks to minutes' is AWS's Nov 2025 portal launch note. The draft contrasts the two, so check both quotes.
- [ ] Source disclosure: three of the four authors are AWS staff; Jakarry Williams (Revinate, Engineering Integration team lead) is the fourth. Taken from the post's author bios.
- [ ] Agentgateway is linked as an example of a self-managed gateway in the options section. It is an LLM/MCP-focused gateway, so check the link earns its place.
- [ ] In the last section, the line saying inbound APIs turn usage plans from a 'read gate' into a 'write gate' where a hard cap matters more is our own analysis.

**Fact-checker could not confirm:**
- [ ] Viability, for the Author to decide: the post gives no traffic, latency, cost, throttle or quota numbers and no versioning scheme. It calls its results 'qualitative', and its only figures are weeks-to-days onboarding and 'sub-second latency at scale'. The draft says this openly in its evidence section, but whether that meets the Case Study bar is the Author's call.
- [ ] Opening line: 'the design question it wrote up in most detail' is our reading. The EKS connection section is the post's most detailed comparison, and the post itself calls it 'A critical design decision'.
- [ ] Agentgateway link (/apis-networking/agentgateway/, content/explainers/apis-networking/agentgateway.md, not a draft): kept. HastCode's own page describes it as a general-purpose HTTP and gRPC data plane as well as an LLM/MCP gateway, so it is a fair example of a self-managed open-source proxy. The Author may still prefer a more conventional API gateway as the example.

<details><summary>What the fact-checker changed</summary>

- Fixed: TL;DR: 'its one measured result' became 'its one before-and-after figure'. The post calls its results qualitative, and weeks-to-days is not a measurement.
- Fixed: Path section: 'with AWS WAF in front' became 'protected by AWS WAF', because the post does not put WAF in front of CloudFront. 'single entry point' became the post's own phrase, 'unified entry point'.
- Fixed: Options section: 'plugs into the same account' became 'integrates natively with the AWS services it already runs'. The post never mentions an account; it talks about native integration with EKS and Lambda.
- Fixed: Step 3: the account throttle default of 10,000 RPS with a 5,000 burst now says 'in most Regions', with '(2,500 and 1,250 in some)'. This is still current per the quotas page, which lists 13 Regions at 2,500/1,250.
- Fixed: Step 3: the X-API-Key header claim had no link. It now links to the 'Choose an API key source' doc page, which is also added to Sources.
- Fixed: Last section: the unlabelled 'read gate into a write gate' analysis now starts with 'Our view:'.
- Fixed: Last section: 'lets a small team ship' became 'lets a team ship'. The post never says how big the team is.

</details>

## Planned switchover: stop the primary, pause writes or buffer queries

`content/explainers/databases-storage/planned-database-switchover.md` · 2570 words · fact-checker: 78 claims checked, 11 fixed, 1 removed

**Drafter asks you to check:**
- [ ] MySQL's two pages disagree on the timeout for group_replication_set_as_primary(). 14.18.1.1 says 0 (immediately) to 3600 s. 20.5.1.1 says 1 to 3600 s, with 0 meaning no timeout. The draft says only 'up to 3,600 seconds', with no default. Pick the reading to keep.
- [ ] dev.mysql.com returned 403 / Oracle 'Technical Difficulties' to curl, so the two MySQL pages were read only through WebFetch's model summary. Check every MySQL quote word for word against the live pages: 'a read-only secondary', 'the read-write primary', disconnect at expiry, commit-phase transactions allowed to complete, and explicit BEGIN being blocked.
- [ ] We link the claim that Patroni stops Postgres during a switchover to the Lead only. Patroni's own docs say candidates join the leader race 'after the leader stepped down' and do not use the word shutdown.
- [ ] Everything the Stop-the-primary option shows the client (failed connections, writes failing until rerouted, a retry wave) is marked 'we infer'. The table cells for it also say '(we infer)'.
- [ ] Patroni dynamic_configuration.html is linked at page level (no anchors) for ttl/leader lock, primary_stop_timeout and use_pg_rewind. The draft says primary_stop_timeout sends SIGKILL only with synchronous mode on.
- [ ] The 'Distributed Configuration Store' wording and etcd as a DCS come from Patroni's YAML config page (#etcd).
- [ ] The pg_promote() default (wait up to 60 s) is from functions-admin, Recovery Control. The pg_rewind requirements (checksums or wal_log_hints) are from Patroni's use_pg_rewind entry.
- [ ] The Vitess default values are buffer-size 1000, buffer-window 10s and buffer-max-failover-duration 20s, with buffering off by default. The guidance that client read/write timeouts should exceed buffer_max_failover_duration comes from the flags section. The doc's 'Lost connection' section compares against buffer_window instead.
- [ ] The draft says Vitess and the MySQL change-primary page do not describe a rollback. That rests on reading the whole Reparenting page and the whole 20.5.1.1 page.
- [ ] The draft says reads served by replicas are unaffected, quoting Neki. For Vitess it says only that VTGate buffers requests to the primary, from 'Types of queries that can be buffered'.
- [ ] The closing sentence of 'When to use each' generalises the Lead's examples (pod replacement, kernel upgrade, resizes) to all three options. It sits inside the inference section.
- [ ] The page is at 49,161 bytes, 339 under the 49,500 cap. Any additions need matching trims.

**Fact-checker could not confirm:**
- [ ] MySQL quotes were checked only through WebFetch's model output. The auto-mode classifier blocked a direct curl of dev.mysql.com as bot-check circumvention. WebFetch gave back 14.18.1.1's sentence 'The current primary becomes a read-only secondary, and the specified group member becomes the read-write primary.' It also gave back the timeout, commit-phase, BEGIN and DDL sentences, and these support the draft. 20.5.1.1 has the same sentence but spells it 'read/write primary'; the draft's quotes link to 14.18.1.1, so they match. The Author should still check the MySQL wording word for word in a browser.
- [ ] Timeout reading: kept 'up to 3,600 seconds' with no default, as 14.18.1.1 says. The two pages disagree on what 0 means (14.18.1.1: immediately; 20.5.1.1: no timeout, with a 1-3600 range), so the draft deliberately leaves the lower bound out.
- [ ] 'A leader lock held in a Distributed Configuration Store such as etcd' is pieced together from two Patroni pages. dynamic_configuration expands DCS and describes the leader lock, and yaml_configuration has an Etcd section and mentions the leader key in DCS. No single sentence says 'etcd is a DCS'.
- [ ] 'Patroni stops Postgres during a switchover' rests only on the Lead. Patroni's docs say only 'after the leader stepped down'.
- [ ] Patroni REST API sample output shows 4.0.0, while the docs header is 'Patroni 4.1.5 documentation'. The Sources entry says 4.1.5, which matches the header.

<details><summary>What the fact-checker changed</summary>

- Fixed: Patroni leader-race quote ('all the eligible nodes ... after the leader stepped down') comes from the REST API /switchover section, not patronictl. Moved it so it follows the /switchover link, and the patronictl sentence now comes after it.
- Fixed: Neki quote now matches the Lead word for word: 'took longer to complete, but did not fail because the primary changed'. The cut-down version left out the condition.
- Fixed: Neki sync claim reworded to the Lead's own claim: the target is 'the only synchronous standby still streaming'. The draft said 'the only standby that can acknowledge a commit'.
- Fixed: TL;DR pause bullet: 'clients must find the new primary themselves' becomes 'MySQL's page does not say how clients find the new primary'. The page leaves this out, so the draft's line was unsupported.
- Fixed: TL;DR: added '(we infer)' after 'a failed connection', since the client side of stop-the-primary is inferred.
- Fixed: Diagram intro: the stop-the-primary client behaviour is now marked 'we infer'. The pause-in-place blocking now applies only when a timeout is set, because MySQL says new transactions can start when no timeout is set.
- Fixed: Mermaid: the Stop subgraph title now says 'we infer the client side'. P2 is limited to 'With a timeout'. P4 now reads 'Client reaches the new primary (routing not described)'.
- Fixed: Pause-in-place app-sees line: 'Then the client reconnects' becomes 'How the client then reconnects, the page does not describe'.
- Fixed: Trade-offs table, buffer 'A write sees': 'A pause, then success' becomes 'A pause; errors if the buffer fills or times out'. Vitess says buffering does not guarantee zero errors, and the Lead says queries fail when the buffer is full or the deadline passes.
- Fixed: The etcd profile link now says the profile covers 'its lock and election APIs', which is what that profile describes.
- Fixed: Closing inference: 'pod moves' becomes 'pod replacements', matching the Lead's example.
- Removed: 'replacing the usual election' in the group_replication_set_as_primary() sentence. Neither MySQL page says this.

</details>

## Harness engineering: steering coding agents with guides and sensors

`content/explainers/languages-runtimes/agentic-development.md` · 2481 words · fact-checker: 78 claims checked, 5 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Slug: the file is still agentic-development.md, but the Concept is harness engineering. Consider renaming it to harness-engineering.md before publishing (I could not rename it under this task's edit-only-this-file rule).
- [ ] OpenAI's 'Harness engineering' post (openai.com/index/harness-engineering/, cited by Böckeler) returned 403 and could not be read, so it is left out entirely. If the Author can read it, it may be the origin of the term and worth adding to 'What a harness is'.
- [ ] Hook mixes two accounts that may be different runs: the Lead says 'a week', '80% done' and that it was run twice with different teams; 'An Accidental Blackboard' (2 Sep 2026) says 10 engineers in Barcelona and 'built one in four days'. Check that both describe the same run.
- [ ] Name and title differ: the Lead calls him 'Giles Alexander, our Regional CTO for Europe'; his own byline is 'Giles Edwards-Alexander', CTO for Europe, Middle East and India. The draft uses the byline name in Sources only.
- [ ] Naresh Jain founded Specmatic and built his harness with it. The draft discloses this; decide whether that is enough given he is a vendor.
- [ ] Jain's table: Sol's tokens, time and cost (1,076,371; 6 min 51 s; ~$0.406) come from the table near the top, which has no anchor, so that row links the bare URL. The Luna figures and the 3.9x token ratio are from #TestingACost-optimizedModel.
- [ ] Stripe: 'over a thousand PRs merged each week', 'devbox in 10 seconds', 'lints in under five seconds', 'over three million tests', 'at most two rounds of CI' and the quote 'often one, at most two, CI runs...' all come from the stripe.dev post. The page has no section anchors.
- [ ] DORA quotes come from the Google Cloud blog announcement (Nathen Harvey, dated 2025-09-23 in its metadata, shown as 24 September on the page; Sources says 'September 2025'), not the full report PDF.
- [ ] The CI section's sentence beginning 'We infer the mechanism' (CI as a shared sensor whose red mainline stops pointing at one change) is our own reconstruction. Laycock only lists possible causes.
- [ ] The definition section's 'not a sensor in this sense' boundary is marked 'we infer'. Böckeler does not say it explicitly.
- [ ] The GitHub figure 'pushes grew 4.9x year over year' was checked on the GitHub blog post by Brian Celenza (6 Oct 2026, updated 26 Oct 2026).
- [ ] Charity Majors's 'Only prod is prod' (quoted in the Lead via an X post) is deliberately left out, because the X post could not be checked first-hand.
- [ ] related: is set to [meta-thin-orchestrator, github-agent-scale-git, pgrust]. All three are live and linked inline.

**Fact-checker could not confirm:**
- [ ] Laycock says the experiment ran twice with different teams. The opening still reads as one run, mixing the Barcelona write-up (10 engineers, four days) with Laycock's account (80% done, two-thirds of time at the whiteboard). We could not tell which run each figure belongs to.
- [ ] OpenAI's 'Harness engineering' post, cited by Böckeler, was not read and is not cited. It may be where the term comes from.
- [ ] Naresh Jain is a vendor: he founded Specmatic and built the harness with it. The page discloses this, but the Author should decide whether that is enough.
- [ ] Inferences marked 'we infer' have no direct Source, as intended: 'not a sensor in this sense', CI as a shared sensor with finite throughput, and the Meta move being the same as Stripe's.
- [ ] The slug is still agentic-development, but the Concept is harness engineering. Rename the file to harness-engineering.md before publishing, if the Author agrees.

<details><summary>What the fact-checker changed</summary>

- Fixed: Opening: the two accounts of the build conflicted. The page said 'They built one in four days', but Laycock says they got 80% done. It now credits both: 'The write-up says they built one in four days; Thoughtworks CTO Rachel Laycock, who links it, says they got 80% done.' Laycock's 'here' link points to An Accidental Blackboard, so both describe the same experiment. 'Told Thoughtworks CTO Rachel Laycock' became 'told her'.
- Fixed: Steering loop: cut 'not another comment' from 'the fix is a new guide or sensor'. Böckeler says only that repeated issues should lead to better controls.
- Fixed: Stripe: three sentences had no links. The devbox that spins up in 10 seconds, the agent rules applied by subdirectory, and the more than three million tests now each link the stripe.dev post. That page has no anchors.
- Fixed: Jain: the sentence 'He built the harness with Specmatic' had no link. It now links #KeepingGuidesAndSensorsInSync.
- Fixed: Defence one: 'where it is hardest to build' (legacy code) had no link. The words come from Böckeler's Ambient affordances sidebar, so it now links #AmbientAffordances.

</details>

## Decode context parallelism: splitting the KV cache by token

`content/explainers/ai-ml-systems/decode-context-parallelism.md` · 2507 words · fact-checker: 68 claims checked, 14 fixed, 0 removed

**Drafter asks you to check:**
- [ ] The diagram's order of steps (gather query heads, attend over each rank's tokens, partial output plus LSE, combine) is our inference from the DCP-only operations listed in the Lead and from PR #61285. The text says 'we infer'.
- [ ] LSE-weighted combine: the draft says each partial output is weighted by its LSE, based on PR #61285's 'fp32 LSE-weighted reference' and Helix's 'lightweight communication step'. Neither Source spells out the formula.
- [ ] 'Each rank stores about 7,500 of 60,000 tokens' is our own arithmetic (60,000 / 8, assuming an even interleaved split), marked 'we compute'.
- [ ] 2.38M to 18.88M (FP8) and 3.60M to 28.76M (NVFP4): the Lead gives the hardware (8x GB300, gpu-memory-utilization 0.85) only in the NVFP4 chart caption. The draft does not state hardware for the FP8 figures.
- [ ] TTFT 6.89 s to 4.73 s, TPOT 17.82 to 12.50 ms and c64 139.1 to 204.0 tok/s are p50 values on TP8 DCP8 FP8 at 60k/1k, and they include the decode fusions as well as the collectives port. Check the wording.
- [ ] DCP8 vs TP8 per-token slowdown: the 11.1 vs 8.7 ms and 11.4 vs 9.1 ms figures come from 8x GB300 NVFP4 with no speculative decoding. The 16-28% range is the Lead's own figure.
- [ ] In the hook, the agentic P/D test is quoted as the Lead gives it: TP4 DCP4 decode node at 0.7 s TTFT with 32 sessions, a plain TP4 node out of KV at 15 s. The 'same four GPUs' framing is ours.
- [ ] Helix Parallelism (arXiv:2507.07120): the abstract does not state affiliation, so the draft does not name NVIDIA. Add it if the paper states it.
- [ ] GLM-5.3 GSM8K figures (0.918 TP8, 0.914 TP8+DCP8 a2a) are from PR #61285 on 2x4 GB300. The draft says 'eight GB300 GPUs'.
- [ ] 'PCP plus DCP still requires ag_rs' comes from PR #61285, merged 13 October 2026. It may have changed since.
- [ ] The interleaving granularity (per token or per block) is not stated. The draft gives only the docs' wording and the 16-token physical block that represents 64 tokens with DCP=4.
- [ ] The slug stays vllm-day-0-model-support because the task allowed editing only this file. The Author may want to rename it to decode-context-parallelism.

**Fact-checker could not confirm:**
- [ ] The blog does not state that Mistral Large 4 was released on 29 October 2026. 'The day Mistral Large 4 came out' rests on the post's 'day-0 support' framing
- [ ] The LSE definition (log-sum-exp of the attention scores) is standard background and is not in any Source
- [ ] Helix affiliation: the arXiv abstract page lists the authors without an affiliation, so the draft still does not name NVIDIA
- [ ] The builder printed a TODO(author) warning that names pytorch-2-15.md together with this page's URL. This file contains no TODO text, so the warning appears to belong to the other draft
- [ ] Quotes, numbers and anchors were checked against raw page text fetched with curl, because WebFetch would not return verbatim text. All anchors used (#highlights, #usage, #model-architecture, #long-context-decode-and-prefill-context-parallelism, #dspark-speculative-decoding, #disaggregated-serving-with-nixl, #hisparse-host-memory-sparse-kv-offload, #whats-next, #decode-context-parallel, #prefill-context-parallel, #context-parallel-deployment) exist

<details><summary>What the fact-checker changed</summary>

- Fixed: Hook: 'The same four GPUs ... Nothing about the hardware changed' overstated the blog, which names two different decode nodes. Now reads 'A decode node with the same four-GPU count, configured as TP4 DCP4 ...', and the cause is marked 'we infer'
- Fixed: TL;DR: 'one exchange per layer merges' is true only of the a2a backend (ag_rs uses an all-gather plus a reduce-scatter). Now 'an exchange in every layer merges'
- Fixed: (0, -inf): PR #61285 says this applies to rows with no local top-k candidates on the FlashMLA sparse path, not to 'a rank with nothing to contribute'. Reworded to match
- Fixed: LSE-weighted combine: no Source states the formula. Now marked 'we infer', with the PR's fp32 LSE-weighted reference as support
- Fixed: Helix: 'describes the same step' was our inference. Now 'which we infer does a comparable merge', and the link text tracks the abstract ('to preserve exact attention behavior')
- Fixed: DeepSeek-R1 case study: the docs say -dcp 8 'reduces' the duplication. Changed 'removes it' to 'reduces that'
- Fixed: GLM-5.3: PR #61285 enables the a2a combine for FlashMLA sparse and never calls it 'faster'. Reworded, and the model and precision named (GLM-5.3 FP8 on eight GB300, matching 2x4 GB300)
- Fixed: Duplication factor: the docs' case studies do not include Mistral Large 4. Its factor of 8 is now marked 'we compute'
- Fixed: 'each rank only sees candidates from its own tokens' marked 'we infer'
- Fixed: TTFT 6.89 to 4.73 s and TPOT 17.82 to 12.50 ms are p50 values. Now labelled 'median', with 60k/1k added to the TTFT sentence
- Fixed: Added a link to the sentence with the 11.1 vs 8.7 ms and 11.4 vs 9.1 ms TPOT figures, which had none
- Fixed: '128 concurrent sequences' changed to max_num_seqs 128, the blog's wording
- Fixed: The PCP 'amortise the prefill computation and control TTFT' claim comes from the vLLM docs, not the blog. Added a link to the docs
- Fixed: 'PCP plus DCP still requires ag_rs' is now dated to the PR ('as of that pull request ... still required')

</details>

## Catching JiTTests: how Meta writes tests meant to fail on a risky diff

`content/explainers/ai-ml-systems/meta-jittesting.md` · 2645 words · fact-checker: 112 claims checked, 8 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Paper internal inconsistency: it reports 8 strong catches from 41 reach-outs (19.5%), but one sentence in section 5 says 'we were able to catch 7 issues from 41 reach outs ... 19.5%'. 8/41 = 19.5%, so the draft uses 8 (as the abstract does).
- [ ] Paper Table 6 date range reads '11th November 202 to 11th December 2025' (typo in the source); the draft says 'in the month to 11 December 2025' for the 7,987 passing tests across 2,517 diffs.
- [ ] 'Roughly 920 engineers were contacted' is our own arithmetic (589 / 0.64), labelled as ours; Meta does not state the outreach total.
- [ ] The 60% true-positive rate is what engaged engineers 'reported' (post text; the post's impact figure labels it 'Precision (True Positives)'). The draft says the post does not explain how it is counted. Check that this wording is fair.
- [ ] LLM-as-Judge model names (Llama3.3-70B, Gemini 3 Pro, Claude Sonnet 4) come from the paper and describe the paper-era deployment; the post only says 'several state-of-the-art LLMs'. The draft says 'at the time'.
- [ ] Stage 1 detail 'overnight (west coast USA)' and 'all diffs flagged as high risk' come from the paper; 'similar to Diff Risk Score' and 'spare capacity' come from both sources.
- [ ] The mutant 'builds and passes existing tests' check and the 'test builds and passes on parent' check are read from Meta's Figure 1 (an image), not from the post's text.
- [ ] Statement that the initial deployment was human-in-the-loop (paper section 6) while production is human-on-the-loop (post) is presented as an evolution; the post does not explicitly say the mode changed.
- [ ] Inferences labelled in text: targeter is what makes it affordable; most dodgy-diff failures are intended behaviour changes; authors hear the next day at the earliest; RubFake rules need upkeep; the dotted 'humans on the loop' edges in the diagram are ours; the comparison to the thin orchestrator is ours.
- [ ] Table of baseline weak-catch rates (0.2%, 2.0%/3.6%, 0.8%/1.6%, 2.5%/4.0%, 6.4%/7.9%) transcribed from the paper's Table 1 (HTML version); the TestGen-LLM row used is the 'need not add coverage' variant.
- [ ] Paper's '98% precision' for identifying false positives is reported with the paper's own caveat (unbalanced population); the draft treats the 70% review-load saving as the headline result, as the paper does.
- [ ] RubFake's 'scores a likely true-positive pattern higher when it is also quick to dismiss' is from paper section 5.1, which says 'in the rule based approach'. The draft attributes it to the rule-based approach and does not name RubFake in that sentence.

**Fact-checker could not confirm:**
- [ ] 'mid-September to the end of October 2025': paper S5 gives no year. The 2025 comes from the paper's statement that production started in September 2025, so the year is inferred.
- [ ] Preview build warning not caused by this file: the page's related-cards strip shows the untracked skeleton pytorch-2-15.md, which still has 'TODO(author)' summary text. The builder reports this against /ai-ml-systems/meta-jittesting/. It needs fixing in pytorch-2-15.md (outside this task).

<details><summary>What the fact-checker changed</summary>

- Fixed: TL;DR: 'they rarely read the test code at all' changed to 'they rarely need to read the test code' (post: 'engineers rarely need to read test code')
- Fixed: Hardening tests: 'Most test generation, Meta's included, produces' changed to 'Most work on unit test generation produces' (paper S1 wording; 'Meta's included' was not stated)
- Fixed: Oracle paragraph: 'A catching test usually fails because behaviour changed' (unsourced) replaced with the paper's point that catching tests cannot rely on the implicit oracle, so whether a failure reveals a bug depends on the general oracle
- Fixed: 'Meta's Table 1' changed to 'The post's Table 1', so it is not confused with the paper's Table 1 (comment/failing-test contrast checked against the post's Table 1 image)
- Fixed: Mutant build-and-pass check now links to Meta's post (Figure 1 image confirmed: 'Mutant build and pass existing tests?')
- Fixed: Coupling hypothesis: 'real faults tend to be coupled to simple mutants' changed to the paper's hedged wording, 'reasonable to suppose real faults may be coupled to such mutants' ('simple' is not in the Source)
- Fixed: 'The authors call the sample "admittedly limited"' changed to 'call that number of bugs' (paper S5: 'admittedly limited number of bugs')
- Fixed: 60% gap sentence: replaced 'or whether it covers engineers who did not engage' (the post does say the engaged engineers reported it) with: the post does not say how the 60% is counted, its impact figure labels it precision, and the number comes only from engineers who engaged

</details>

## Fully sharded data parallelism: training models no GPU can hold

`content/explainers/ai-ml-systems/fully-sharded-data-parallelism.md` · 2300 words · fact-checker: 78 claims checked, 3 fixed, 0 removed

**Drafter asks you to check:**
- [ ] ZeRO numbers (16 bytes per parameter = 2+2+12; GPT-2 1.5B needs 3 GB of weights and at least 24 GB in total; 7.5B on 64 GPUs: 120 / 31.4 / 16.6 / 1.88 GB) are cited to ar5iv.labs.arxiv.org/html/1910.02054. The S5.T1 anchor was reported by the fetch tool. Section 3.1's anchor is not confirmed, so those links use the bare URL. ar5iv was slow to load with curl.
- [ ] The PVLDB paper's GPT-175B result (more than 173 and 186 TFLOPS per GPU, about 55% and 60% of peak) is in section 5.4 (#S5.SS4). Backward prefetch's 18% and DDP's OOM above 2.28B are in section 5.2. The 512 A100 80GB GPUs are in section 5.1. The 33M-element all-gather threshold and the 1.5x communication over DDP are in section 3.2.1. All were checked against arxiv.org/html/2304.11277v2.
- [ ] The 'Meta AI authors' affiliation for the FSDP paper comes from the arXiv HTML. No affiliation is given for the ZeRO authors.
- [ ] The hook's '2x reduce-scatter bytes' and the before/after bf16/fp32 description are quoted from the body of PR #194434, written by the PR author, not from the release blog.
- [ ] Figures for the 1.75B KDA-style model (all-gather copy-in 4.2 to 4.9% slower, forward plus backward about 53 ms) and 'keeping only one AG per fully_shard' come from the body of PR #196626.
- [ ] The release date of 28 October 2026 comes from the GitHub release's published_at timestamp (2026-10-28T17:09:49Z). The blog post itself was fetched without a visible date.
- [ ] The '3k vs 14k lines' figure comes from the RFC's progress tracker as relayed by the fetch summary. Re-check the exact wording in issue #114299.
- [ ] Three 'we infer' sentences: (1) the peak for a per-layer-wrapped transformer is the shard plus one or two full layers; (2) an int reshard_after_forward keeps the backward gather within fewer ranks; (3) a hidden collective costs memory rather than time. A fourth, the 'FSDP moves weights to the data' boundary with tensor and expert parallelism, is inferred from the user contract.
- [ ] The release blog is cited without section anchors because its heading IDs were not confirmed.

<details><summary>What the fact-checker changed</summary>

- Fixed: reshard_after_forward as an int: the draft marked 'keeps the backward gather within fewer ranks' as our inference, but the 2.15 fully_shard docs state it directly ('allows the all-gather in backward to be over a smaller world size at the cost of higher memory usage than setting to True'). The sentence now reports the sourced fact: 'so the backward all-gather runs over a smaller world size, at the cost of more memory than `True`.'
- Fixed: Practice one: 'Each prefetched layer is another full layer in memory.' had no link and goes further than the docs and tutorial, which say only 'more reserved memory' or 'increased memory'. It is now labelled as inference: 'Each prefetched layer, we infer, is another full layer in memory.'
- Fixed: PR #196626: 'reuses one byte buffer' changed to 'reuses the existing uint8 communication buffer', which is the PR body's wording ('invokes the existing uint8 communication buffer').

</details>

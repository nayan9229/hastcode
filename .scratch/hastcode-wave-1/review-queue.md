# Wave 1 review queue

Drafted 2026-10-27 by agents: 21 Pieces, every one `draft: true`, none verified. Each was fact-checked by a second agent against its Sources. Verify every claim yourself before removing `draft: true` and adding `verified:` (ADR-0007). Publish at most two a day (ADR-0007 as amended).

**Order:** the two Briefs first (news goes stale), each together with or after the Explainer it links: the JDK 28 Brief needs `jvm-warmup-and-aot-compilation` live, the OCP Brief needs `meta-prometheus-supercluster` live. A published Piece that links a draft fails the production build (broken-link).

**Not drafted:** Netflix Conductor (`pcwu4qit`): the post returns 403 and has left Netflix's feed; skeleton removed. npm outage Brief (`zt2nrtbh`): the status page gives no cause; skeleton removed. Both Leads stay marked picked in the local intake state, so they will not come back to the queue; to retry Conductor, copy the case-study template by hand.

## OpenJDK's JEP 544: cached compiled code, with the JIT as fallback

`content/briefs/jdk-28-aot-code-compilation.md` · 243 words · fact-checker: 17 claims checked, 3 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Freshness: the Lead's date comes from the Inside Java feed (published 2026-10-25T00:00:00Z, date only, author John Rose). The 72-hour window runs to about 2026-10-28T00:00Z, so verify and publish fast or the Brief goes stale. The Source gives no time of day, so 'What happened' states only the date, not a UTC time.
- [ ] The JDK 28 GA date of 9 March 2027 comes from the 'Schedule (as proposed)' on openjdk.org/projects/jdk/28/ and is worded 'as proposed'. Check it hasn't changed.
- [ ] The JEP 544 header now shows Status 'Integrated', Release 28 (updated 2026/10/14). The Brief says only 'targeted to JDK 28', as the Lead does.
- [ ] Numbers: 'around 65% to 80%' with AOT code and '50% to 70%' without come from JEP 544 #Performance. That run used five framework benchmarks on a two-core Linux/x64 system with a JDK 28 early-access build. The Brief says 'two-core machine' and leaves out Linux/x64 and the EA build.
- [ ] 'with no new options' paraphrases JEP 544 #Description ('No additional options or settings are required to generate or use AOT code'; the workflow is unchanged).
- [ ] 'only one set of hot spots' paraphrases the JEP's static-compilation section ('can be optimized for only one set of hot spots'). The link is there, but it has no quotation marks. The Brief contains no direct quotes.
- [ ] Fallback claim: a CPU feature or GC mismatch makes HotSpot warn, skip the AOT code and still run correctly (JEP 544 #Consistency-of-training-and-production-runs). The 'if the workload shifts, it recompiles' clause is not linked; it comes from the JEP's Summary and Shifting-compilation-work sections.
- [ ] Labelled inference: 'a cache of expensive work is safe to ship when a miss falls back to the slow path, not to an error' is our generalisation. The JEP does not state it.
- [ ] AOTMode=required 'exits when a constraint is violated' comes from JEP 544 #Observing-AOT-cache-usage-in-production ('report an error and exit if the AOT cache violates any constraint').
- [ ] Anchor links use the JEP's heading ids as they appear in the page HTML, including '#What-about-static-compilation?'. Check that they resolve.
- [ ] I uncommented `sources` in the frontmatter and filled it with JEP 544 and the JDK 28 project page. Title and tags were replaced. The AI note, author, published and draft: true are unchanged, and no verified: was added. The Concept the Brief links was still being drafted in parallel. Check that its interpreter-to-C2 ladder section survives in the final version, because the Brief points readers to it.

**Fact-checker could not confirm:**
- [ ] Announcement date 25 October 2026: the Inside Java page shows no date, but the URL path (2026/10/25) and the Inside Java Atom feed both give it (published 2026-10-25T00:00:00Z, author John Rose). Today is 2026-10-27, so the Brief's 3-day freshness window ends around 2026-10-28.
- [ ] JDK 28 General Availability on 9 March 2027 is confirmed from 'Schedule (as proposed)' on openjdk.org/projects/jdk/28/ (page last updated 2026/10/22). It is a proposed date and could still change.
- [ ] The linked Concept jvm-warmup-and-aot-compilation exists and has the section 'The ladder from interpreter to C2'. It was being drafted in parallel and is still draft: true, so the Brief's link depends on that Concept being published first.

<details><summary>What the fact-checker changed</summary>

- Fixed: Unlinked claim 'the JIT ... redoes profiling and compiling on every start' now links to JEP 544 #Shifting-compilation-work-to-training-runs. That section says the cache spares HotSpot 'having to recompile it at the start of each production run', and #Motivation describes profiling in every run.
- Fixed: Unlinked claim 'if the workload shifts, it recompiles' now has 'recompiles' linked to JEP 544 #Shifting-compilation-work-to-training-runs ('if the application's workload changes then HotSpot can, as usual, deoptimize and reoptimize').
- Fixed: Fallback consequence read 'only slower to warm up'. It now reads 'only slower to start and warm up', to match JEP 544 #Consistency-of-training-and-production-runs: 'the application may start up and warm up more slowly, but its execution will still be correct'.

</details>

## Meta's OCP 2026 networking: three fabrics for three kinds of traffic

`content/briefs/meta-ocp-2026-ai-networking.md` · 231 words · fact-checker: 16 claims checked, 3 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Date: the Source shows only 'Posted on October 26, 2026' with no time, so the Brief says 'On 26 October' with no UTC time. Check the post time if a UTC time is wanted. The Lead was picked 27 Oct, so it is fresh, but check that it is still under 72 hours old when published.
- [ ] Scale-out claim: 'connects each GPU to eight network planes and flattens a three-tier Clos to two' paraphrases the Source's highlight '3-tier → 2-tier Clos'.
- [ ] Scale-up claim: Tahan is described as 'an Ethernet switch for its MTIA 300 system'. The Source calls it the 'Ethernet scale-up switch for MTIA 300 based on Broadcom Tomahawk5'.
- [ ] Why-it-matters basis: the claim that 'latency- and bandwidth-sensitive operations concentrate within a model replica, while traffic between replicas tolerates more latency' comes from the Prometheus post (19 Oct), not the talks post. Confirm the wording against it.
- [ ] Labelled inference: 'densest in the rack, flatter across a building, longest between sites' is our own characterisation of the three tiers. The sentence that follows says 'We infer' that the traffic split is why there are three fabrics. Check that the hedge covers the characterisation too, or tighten it.
- [ ] Paraphrase, not a quote: 'collective operations run only as fast as the slowest node' paraphrases the MetaRoCE talk blurb ('all-reduce and all-to-all performance gated by the slowest node across hundreds of thousands of accelerators'). The Brief has no direct quotes.
- [ ] The closing advice ('oversubscribe the cold ones') is HastCode's own analysis. It echoes the Case Study's BAG oversubscription point, which the talks post does not make.
- [ ] Frontmatter: I removed 'TODO(author)' from the skeleton's YAML comment so the placeholder scan cannot flag this file. I also added an optional `sources:` entry for the Prometheus post. I kept draft: true and did not add verified:.

**Fact-checker could not confirm:**
- [ ] Publication time: the talks post shows only 'October 26, 2026' with no time, so 'On 26 October' is correct, but the brief guide's UTC time is not available. The 72-hour freshness window should be checked at publish.
- [ ] 'For scale-up inside the rack': the post calls Tahan 'the Ethernet scale-up switch for MTIA 300' and MTIA 300 a rack-scale system with 72 accelerators per rack, but never says 'in-rack' outright. This is a reasonable reading.
- [ ] Description: 'each built for the traffic it carries' is a mild interpretive gloss. The post lists the three domains but does not explain why they are separate. Left as is; the Author may want to soften it.
- [ ] Closing advice ('oversubscribe the cold ones') is HastCode's own analysis. It matches the Prometheus Case Study's BAG oversubscription point (~4.5:1, from Meta's Feb 2026 BAG post), not either source of this Brief.

<details><summary>What the fact-checker changed</summary>

- Fixed: Why it matters opened with an unhedged cause, 'Meta runs three networks, not one, because training traffic is not uniform', and the 'We infer' hedge came three sentences later. Neither source gives that reason. It now reads 'We infer that Meta runs three networks, not one, because training traffic is not uniform; the talks post does not say so.' The trailing 'We infer that this split is the reason...' sentence is removed as redundant.
- Fixed: The tier characterisation ('densest in the rack, flatter across a building, longest between sites') read as fact. It is our own reading, so it now opens with 'On our reading, that lets each tier trade bandwidth against distance differently'.
- Fixed: The MetaRoCE claim had no inline link. 'MetaRoCE talk' now links to the talks post, which says 'all-reduce and all-to-all performance gated by the slowest node'. 'adds why it pays' became 'shows the stakes', so the talk is no longer credited with justifying the three-fabric split.

</details>

## TimeSeries cold tier: how Netflix serves reads from Cassandra backups

`content/explainers/databases-storage/netflix-timeseries-tiered-storage.md` · 2800 words · fact-checker: 92 claims checked, 8 fixed, 0 removed

**Drafter asks you to check:**
- [ ] The Lead returned 403 to both WebFetch and curl. It was read in full through the Netflix TechBlog RSS feed (content:encoded, pubDate 14 Oct 2026, five named authors). Check it against the post as published, as the guide asks.
- [ ] '80% to 7%' disk usage comes only from a figure caption ('Storage Reduction During Rollout: from 80% disk usage to 7%'). It is used in the hook and the evidence section.
- [ ] We put the streaming tier's 'up to 1,250 reads per second globally' beside the Cassandra tier's 'thousands of reads per second in any given region'. The two were measured at different times, and the comparison is ours.
- [ ] Inferences labelled 'we infer': netflix_views is Netflix's own layer that aggregates node-local virtual tables; the cold cluster must match topology because backups are keyed by node token; CompressionInfo turns a position into an S3 byte range; SSTable overlap costs an S3 round trip per file; backups are now production data, so retention rules matter; a late write to a frozen slice has nowhere to go. That last point is in the cost section, which also carries an unlabelled 'two costs the post does not discuss'.
- [ ] The text says the control plane switches a table to the backup Memtable. The post shows only a CREATE TABLE with memtable = 'backupmemtable:timestamp=...' and mentions 'DDL change protection'. It never says ALTER TABLE, and the draft avoids claiming it.
- [ ] Our own analysis, not Netflix's words: 'In our reading' the sub-second need rules out a batch data lake; the Memtable hook as 'the design's neatest move'; 'every cold byte was written twice'; turning off read repair and speculative retry as a general rule to revisit safety mechanisms that assumed local-disk latency.
- [ ] Short quotes from Netflix: 'is expensive and adds moving pieces that widen the team's on-call scope'; 'a more radical question: Can we read directly from existing Cassandra backups in S3?'; 'turning a cost-saving feature into a reliability risk'. Also short quotes from the Cassandra docs (SSTables definition, Data.db, virtual tables).
- [ ] Parquet's recommended row group size of 512 MB to 1 GB comes from the Parquet configurations page. Netflix links that page but does not state the figure itself.
- [ ] Prose runs to about 2,800 words, the top of the guide's range. The largest section is 228 words, under the 250 limit. related: was left commented out.

**Fact-checker could not confirm:**
- [ ] The Netflix post was read in full through the Netflix TechBlog RSS feed (content:encoded, pubDate Wed 14 Oct 2026, five named authors in the byline), not the live page, which refuses direct fetches. The Author should check the claims against the post as published.
- [ ] Figure contents were checked only for two images (the 'Additional Infrastructure and Operational Cost' diagram has no dollar figures, which supports the 'no money' gap; the lifecycle figure matches the Mermaid order). The '80% to 7%' and '~30% p90' figures rest on the caption and body text.
- [ ] 'The latency numbers come from dashboards in the post, not a published benchmark' describes the post's charts. It is a fair description but was not checked against every image.

<details><summary>What the fact-checker changed</summary>

- Fixed: Hook: 'on the clusters it reached, disk usage fell from 80% to 7%' went further than the Source, which has only one rollout chart caption. It now says 'a chart of the rollout shows disk usage falling from 80% to 7%'.
- Fixed: Compaction job: 'On the largest datasets the job took about 12 hours' now says 'On some of the largest datasets the job typically took about 12 hours', matching the Source's 'On some of our largest datasets, this job typically takes about 12 hours'. A link was added.
- Fixed: Cost section: 'a sudden spike of cold reads' now says 'a sudden spike in read requests'. The Source does not say the spike is of cold reads.
- Fixed: Step 4: the draft said the netflix_views keyspace 'returns one row per host'. Only the cluster_view query returns one row per host; complete_backups has no host column. It now reads 'a netflix_views keyspace whose cluster_view returns one row per host'.
- Fixed: Cost section: 'Netflix states them as choices: ... cold reads are slower' claimed Netflix framed slower cold reads as a choice, which it does not. It now says Netflix accepts some over-fetch on small reads from larger compression chunks (linked), says cold reads stay slower, and marks 'fine for an audit tool' as 'in our reading'.
- Fixed: Evidence bullet: the cold tier's 'thousands per region' next to the streaming tier's '1,250 globally' is now marked as our comparison, not Netflix's. The drafter's note that the two were 'measured at different times' is not in the Source: the post gives the 1,250 figure 'at the time of writing', so the claim was not added.
- Fixed: Unmarked inferences now carry labels: 'in our reading the contest is about operations more than speed'; ROCS 'so, we infer, each key lives in one file'; 'we infer a slice that is not ready simply waits'; 'In our reading, the last two protect a live, changing table...'.
- Fixed: Added inline Source links where a number or quote had none: the 150 ms / 500 ms latency sentence, the 'is expensive and adds moving pieces' quote, the 'a more radical question' quote, the complete_backups view, and ReadOnlyCompactionStrategy in the TL;DR.

</details>

## Sprites Block Device: how Fly.io runs ext4 on a bucket

`content/explainers/infrastructure-cloud/fly-io-sbd.md` · 2884 words · fact-checker: 96 claims checked, 7 fixed, 0 removed

**Drafter asks you to check:**
- [ ] WebFetch would not return the Lead word for word, so its text was read from a curl of the published HTML. Check the quotes against the live page: "so each one could be right while the others were wrong", "stops and says so", "because a clean restart is the only way back to a disk you can trust", "orders of magnitude more reliable", "increasingly clever ways to keep three ledgers in agreement", "spaghetti that comes with a repair kit is gonna leave a weird taste in your mouth" and "have left the building".
- [ ] The diagram puts SBD's local cache on the server's NVMe. The Lead says only 'its own local cache'. The NVMe tier comes from the current Lifecycle docs, which do not name SBD.
- [ ] In the diagram, the checkpoint catalog is a separate box fed by sbd serve. The Lead does not say where the catalog is stored.
- [ ] Three inferences are labelled 'we infer': that fsync becomes an NBD flush and the flush waits for the bucket; that SBD must order its uploads so the map never points at a missing block; and that running JuiceFS without writeback would close this crash window but leave three records. The Lead discusses none of these.
- [ ] The options list has four entries. Two come from the Lead (keep JuiceFS and repair it; the block device). Two come from the earlier design post (attached NVMe, as Fly Volumes do; the LSVD feature on Fly Machines). The no-writeback option is ours. The Lead never sets these up as alternatives Fly.io weighed.
- [ ] The draft keeps 'created before 10/6' exactly as the Lead writes it, with no date conversion. Decide whether to state it as 6 October 2026.
- [ ] The Lead presents the JuiceFS failure as 'a story which may or may not be based on lived experience'. The draft keeps that hedge in the quote and does not say the failure happened in production.
- [ ] The Lead gives no publication date in its text. The date 22 October 2026 and the author Daniel Botha come from the page's metadata (time datetime="2026-10-22" and the byline). The design post's date, 14 January 2026, and its author, Thomas Ptacek, come from its own page.
- [ ] Saying Sprites keep their storage stack in the VM's root namespace relies on the January design post, which described the JuiceFS stack. That sbd serve runs there too is implied by 'a daemon inside the Sprite', not stated.
- [ ] Total is about 2,880 words counting the TL;DR and Sources (about 2,810 without them), slightly over the guide's 2,800 target. The options section is the longest at about 241 words.

**Fact-checker could not confirm:**
- [ ] I checked every quote word for word against a curl of the live Lead (fly.io/blog/sprites-block-device/) and of the design post. All seven quotes the drafter flagged match exactly. All anchors exist on both posts and the docs pages. The Lead's date (2026-10-22) and author (Daniel Botha) are confirmed. So are the design post's date (2026-01-14) and author (Thomas Ptacek).
- [ ] Diagram: the Lead says only 'its own local cache'. Labelling it 'Local cache on the server's NVMe' relies on the current Lifecycle docs ('A hot NVMe cache'), which do not name SBD. I kept it because the docs describe every current Sprite. The Author should decide.
- [ ] 'Sprites created before 10/6' is kept exactly as the Lead writes it. The Lead does not give the year.
- [ ] That `sbd serve` runs in the VM's root namespace is implied, not stated. The Lead says 'a daemon inside the Sprite'. The January design post says the storage stack 'lives' in the root namespace, but that post described the JuiceFS stack. The prose says only that it 'follows a pattern', which I left as it is.
- [ ] Preview build: the only error reported against fly-io-sbd.md is a broken-link to /databases-storage/apache-iceberg-table-format/. Nothing in this file causes it. The site generates that link (author/vertical/related listing), and the Iceberg page drops out of the build because it fails its own html-size gate (52,334 bytes > 51,200). Other pages show the same error, and it clears once the Iceberg draft is trimmed. There are no other errors or warnings for this file, and the options H2 is 239 words, under the 250 section-length limit.

<details><summary>What the fact-checker changed</summary>

- Fixed: Options intro said 'The post and Fly.io's earlier design post name four routes', but neither post names the no-writeback route. It now says 'Fly.io's two posts point to three routes, and we add a fourth.'
- Fixed: Step 3 said the SQLite map was 'made durable by Litestream on its own schedule'. No Source mentions Litestream's timing, and the Lead puts the gap down to writeback mode. It now says 'made durable separately from the chunks' (Lead: 'Each one was made durable separately').
- Fixed: Step 3 said 'a missing or wrong chunk surfaced as an ext4 checksum failure'. The Lead describes only a missing chunk (a 404), so 'or wrong' was removed.
- Fixed: The fail-stop section said that after a restart 'only services the runtime owns come back'. The Lifecycle docs say Services are processes the runtime owns and that it restarts them. It now reads 'only services, the processes the runtime owns, are started again.'
- Fixed: The migration section said Fly.io 'recommends the slower one'. The docs say in-place is 'simpler, but carries more risk' and do not call either route slower. It now says 'recommends the more cautious one'.
- Fixed: The rolling-back sentence said the docs 'warn against upgrading a second time'. The docs say 'Don't run the upgrade again after rolling back'. It now reads 'warn against running the upgrade again after a rollback'.
- Fixed: Added the missing inline links in the TL;DR: the `sbd serve` daemon now links the Lead's #just-a-disk, fail-stop links #dumber-checkpoints-are-better-checkpoints, and the owner-started migration links the docs at sprite-maintenance#storage-upgrades.

</details>

## Meta's thin orchestrator: how a 1,300-line router steers 14 AI agents

`content/explainers/ai-ml-systems/meta-thin-orchestrator.md` · 2555 words · fact-checker: 74 claims checked, 11 fixed, 1 removed

**Drafter asks you to check:**
- [ ] WebFetch refused to give the Lead verbatim, so I downloaded the HTML with curl and pulled out the text myself. Please check every quote against the published page: 'performs zero domain work itself', 'It worked in demos. In production, it started hallucinating.', 'This wasn't a prompt engineering problem. It was an architectural one.', 'constrain[s] the LLM's decision space', 'It's a router.', 'No agent sees the full picture', 'templated work, not greenfield engineering', 'every update requires redeployment', 'Zero new agent code', 'with human review always gating any changes', and the closing line about 'the one that does the least work'.
- [ ] Meta calls the Mount Sinai study a 'preprint'. The page it links is a published npj Health Systems Brief Communication (vol. 3, article 23, online 9 March 2026), so the draft says 'study' and cites the paper directly. The 65.3%, 16.6% and 65-fold numbers come from the paper's abstract and are about clinical tasks, not Meta's pipeline.
- [ ] Three inferences are ours and are marked 'we infer' in the text. (1) Each template phase is one pass through the seven-step protocol, which is how the diagram's loop draws it. (2) The instruction file read in protocol step 1 is what supplies the subagent's pre-gathering checklist. (3) The April 2026 'tribal knowledge' post is linked to this system. That one rests only on two shared authors (Plawan Rath, Ashwini Verma) and the shared 'config-as-code' pipeline wording. The October post does not link April's, and April says four repositories where October says three.
- [ ] The post gives two counts of data sources: 'population patterns for 10+ data sources' for the first single agent, and the population agent covers 'each of 10 data sources'. The draft uses 'more than ten' and 'ten' to match each.
- [ ] Who wrote the 70 to 80% / 10 to 20% numbers and what they describe: Meta calls it 'an alternative single-agent approach' that it evaluated. The draft presents it as a design Meta tried. Check that this does not go beyond what the Source says.
- [ ] The steps are numbered two ways. The draft's 'Step 1-5' H2s are our walk through the diagram, separate from Meta's seven protocol steps, which the text calls 'protocol step N'. Check that the reader cannot mix them up.
- [ ] The cost section's claims about upkeep (knowledge files spread mistakes as fast as fixes) and about requests that fit no template going to a person are our analysis, labelled 'we infer'. The 'roughly a third still need significant edits' figure is our arithmetic from Meta's 'about two-thirds'.
- [ ] The closing paragraph compares the design to a workflow engine that sequences idempotent workers. That is our framing, not Meta's, and it is not linked. Keep it or cut it.

**Fact-checker could not confirm:**
- [ ] The quotes were checked against the article text pulled from the published HTML with curl, not through WebFetch. Every listed quote matched verbatim, including 'performs zero domain work itself', 'It worked in demos. In production, it started hallucinating.', 'This wasn't a prompt engineering problem. It was an architectural one.', 'constrain[s] the LLM's decision space', 'It's a router.', 'No agent sees the full picture', 'templated work, not greenfield engineering', 'every update requires redeployment', 'Zero new agent code', 'with human review always gating any changes', 'the most important component in our AI agent system is the one that does the least work.'. The Author should still look at the live page once.
- [ ] The link between the April 'tribal knowledge' post and this system is our inference, and it is marked as one. It rests only on the shared authors (Plawan Rath, Ashwini Verma) and the shared 'config-as-code' wording. The October post does not link the April one.
- [ ] 'About a third still need significant human edits' (TL;DR and cost section) is our arithmetic from Meta's 'about two-thirds land without significant human edits'. The cost section says so, but the TL;DR bullet does not.
- [ ] Mount Sinai: Meta calls it a 'preprint', but the page Meta links is the published npj Health Systems Brief Communication (vol. 3, art. 23, online 2026-03-09). The 65.3%, 16.6% and 65-fold figures match its abstract. The draft's 'study' wording stands.

<details><summary>What the fact-checker changed</summary>

- Fixed: Diagram labels 'Step 6' and 'Step 7' could be confused with the article's own Step 1-5 H2s. They now say 'Protocol step 6' and 'Protocol step 7'. The body text uses 'Protocol step 6/7' and 'protocol steps 1 to 5' in the same places. A sentence under the diagram now says the Step 1-5 sections are our walk through the diagram, not Meta's seven protocol steps.
- Fixed: H2 'The designs Meta tried before going thin' / 'three designs it moved away from' said Meta tried the heavyweight orchestrator, which the post does not say. It now reads 'The designs Meta weighed before going thin' / 'three designs it argues against'. The prompt-per-task bullet now says 'Meta evaluated this alternative', which matches the post's 'We also evaluated an alternative single-agent approach'.
- Fixed: 'A router that follows a fixed protocol has no reasoning step to get wrong, so any hallucination is confined to one specialist's step' went further than the post, which says it 'leaves little room for hallucination'. It now follows the post: little room for hallucination, Meta argues, because routing and checks don't rely on a model, and the risk moves into isolated contexts where failures are contained.
- Fixed: 'Intent leads to a workflow template' was stated as fact, but the post never says intent picks the template. It is now marked: 'We infer that intent selects a workflow template; the post does not say how one is chosen.'
- Fixed: 'The gate means a known warning or a required reviewer stops the run' was not supported (the post only says the gate checks for active warnings or required reviewers). The paragraph is now marked 'We infer...' and the gate sentence says 'looks for active warnings or required reviewers before work begins'.
- Fixed: 'delegates "with full context", meaning everything steps 1 to 5 gathered' was an unmarked gloss. It now reads 'which we read as everything protocol steps 1 to 5 gathered'.
- Fixed: 'in Meta's words the 15th looks structurally identical' presented a paraphrase as if it were a quote. It now reads 'Meta says the 15th looks...'.
- Fixed: The April post comparison said April's pipeline 'spans four' repositories. April's opening says four, but it later says 'across three repositories'. The text now says April's 'opens by counting four (it later says three)'.
- Fixed: Unmarked analysis is now labelled as ours: 'We infer that each mistake is plausible on its own...', 'We infer that a bad step costs one step', and 'We infer that the router can only run workflows someone has written a template for, and that a request that fits no template...'.
- Fixed: 'Meta's planned answer is automated feedback loops' overstated it, because the post says Meta is 'exploring' them. It now reads 'Meta is exploring automated feedback loops...'.
- Fixed: The closing workflow-engine comparison was our framing stated as fact, and it included an unsourced 'idempotent workers'. It now reads 'In our view, none of this is new...: a workflow engine that sequences workers...', and Meta's part is now that its post 'argues it holds' instead of 'showing' it.
- Removed: 'idempotent' (describing the workers in the closing workflow-engine comparison): no Source says Meta's subagents are idempotent.

</details>

## Hyperdrive: how Cloudflare pools and caches MySQL for Workers

`content/explainers/apis-networking/cloudflare-hyperdrive-mysql-caching.md` · 2456 words · fact-checker: 96 claims checked, 7 fixed, 0 removed

**Drafter asks you to check:**
- [ ] The Lead has no section anchors, so every claim from it links to the bare post URL. I read it from the raw HTML, not the feed (WebFetch would only summarise it).
- [ ] The 'options weighed' section uses alternatives from the April 2025 architecture post (regional pooler, serverless drivers, session or statement pooling, regional pools), not from the Lead. The Lead names no alternatives to the canonical-encoding fix. The two alternatives in Step 5 (capability set in the cache key, unshared caches) are labelled as our analysis.
- [ ] There are three 'we infer' sentences: after the fix, a Client serving mysql2 rewrites cached bytes on every hit; the shared-cache bug mostly hit Workers running near their database; mysql2 misreads the untranslated result (implied, not stated; we only say the post reports symptoms). Check that each is worded as inference.
- [ ] Most latency numbers are Postgres and from 2025, not MySQL: the Chicago-London benchmark of 1,200 / 500 / 320 ms, 125 ms per round trip, the 625 ms handshake, Client handshake p90 4 ms and p50 2 ms, and cache hits cut from 20-70 ms to about 4 ms. The piece says the latency evidence is older and from Postgres.
- [ ] 'One in 100,000' is not used; the hook says 'about 0.001% of the time', as the Lead does. 99.99% is the Lead's 'success rate for Client connections', which the piece also calls the beta's production availability target, as the Lead does.
- [ ] The MySQL docs (dev.mysql.com) refused direct fetches. The capability flag, text-resultset, COM_RESET_CONNECTION and utf8 claims were read through WebFetch summaries. Check the utf8mb3 '3-byte' wording, and that the text-resultset page says EOF follows the column definitions only when CLIENT_DEPRECATE_EOF is not set.
- [ ] 'latin1 on some servers' comes from the Lead's example ('a MySQL server instance that defaults to latin1'). Check that this gloss is acceptable.
- [ ] The frontmatter comment line was reworded from 'replace every TODO(author) line' to 'no author notes are left', because the builder flagged the literal TODO(author) string in the YAML comment as author-placeholder.
- [ ] Added related: [cloudflare-d1, http-connection-reuse] and an internal link to /databases-storage/cloudflare-d1/. Confirm both slugs are the intended Related links.
- [ ] The 'over 300 data centre locations' figure comes from the Hyperdrive docs page (it describes Cloudflare's network).

**Fact-checker could not confirm:**
- [ ] MySQL docs (dev.mysql.com) were read only through WebFetch summaries. A direct curl was denied by the permission classifier. The summaries confirm: CLIENT_DEPRECATE_EOF exists ('use OK_Packet instead'); the text resultset has an EOF after the column definitions only when CLIENT_DEPRECATE_EOF is not set, and ends with an OK packet when it is set; COM_RESET_CONNECTION does not re-authenticate or close the connection; utf8 is a deprecated alias (utf8mb3 by default, unless INTERPRET_UTF8_AS_UTF8MB4 is set); and the utf8mb3 section is titled '3-Byte UTF-8 Unicode Encoding'. The Lead also states the 1-3 bytes figure.
- [ ] 'D1's control plane leans on its cache to keep heartbeats off Postgres' is supported only by the internal page /databases-storage/cloudflare-d1/, which cites the D1 control-plane post. That post was not re-fetched here.
- [ ] Related slugs cloudflare-d1 and http-connection-reuse both exist under content/explainers/. Whether they are the intended Related links is an editorial choice for the Author.

<details><summary>What the fact-checker changed</summary>

- Fixed: Hook: 'The network was fine and so was the database' -> 'No network issue explained it'. The Lead only says no network issue accounted for the bug, and says nothing that clears the database.
- Fixed: TL;DR: 'cache common read queries in its CDN cache' -> 'cache common queries in its global CDN cache', the Lead's wording, since the sentence links to the Lead.
- Fixed: Options section: 'The second half is what later exposed the bug' is our own analysis, so it is now marked 'we infer'.
- Fixed: Step 1: 'accepts whatever capabilities the Worker's driver offers' -> 'adapts to the capabilities the Worker's driver expects'. The Lead says Hyperdrive adjusts to the capabilities the driver expects.
- Fixed: Step 4: 'latin1 on some servers' -> 'such as latin1'. The Lead gives latin1 only as an example of a server default.
- Fixed: Costs: 'the reset fix adds a command after every checkout' -> 'after every connection reset'. The Lead says the utf8mb4 command runs after each COM_RESET_CONNECTION and does not say a reset happens on every checkout.
- Fixed: Evidence: 'Workers Placement ... Workers that make several queries per request' -> 'Placement ... Workers that make multiple sequential queries per request', the docs' name and wording.

</details>

## Spokes: how GitHub reshapes its Git quorum for agent-scale load

`content/explainers/infrastructure-cloud/github-agent-scale-git.md` · 2950 words · fact-checker: 96 claims checked, 7 fixed, 1 removed

**Drafter asks you to check:**
- [ ] Our own arithmetic, labelled as ours in the text: 4 round trips x 25ms is about 100ms per write, or about 10 writes/s under one repository lock. The 4-round-trip count is from the 2017 Stretching Spokes post. The Lead says only 'multiple'. Check whether the illustration is acceptable or should be cut.
- [ ] Our own arithmetic: 0.00001% of 'about a billion repositories' is on the order of a hundred repositories still outside Azure. GitHub gives no count.
- [ ] Rollout status: the draft says the Lead does not state how much of the per-reference quorum or the virtual replicas serves production traffic. It calls them a change GitHub is making and a first step measured in internal tests. Confirm this reading of the Lead's tense ('we're changing', 'the first step is', 'in internal testing').
- [ ] Dates: the earlier post's figures are read as 2026. 'In September alone' (7.38B commits) and 'August' (about 1B requests for the busiest repository) come from a post dated 6 Oct 2026, and the draft says 'in September' and 'in August' without a year.
- [ ] The migration-routing sentence ends 'the old region stops receiving requests for that repository'. The Source says 'entirely', inside a per-repository process. Check that this gloss is fair.
- [ ] Inferences labelled 'we infer': the sweepers drive replica moves; the per-reference quorum builds on Git's per-reference transaction locks; 'for now' compaction on the virtual write replica is a stopgap, given the earlier post's separate maintenance workers; the default or trunk branch stays serialized at the slowest replica's pace; background-reconciling read replicas can briefly serve an older reference value; the GHE non-voting geo-replica is a precedent for reads outside the quorum.
- [ ] Bylines and dates are taken from page metadata: Lead by Eleanor White and Victoria Dye, 21 Oct 2026; earlier post by Brian Celenza, 6 Oct 2026; Stretching Spokes by Michael Haggerty, 13 Oct 2017.
- [ ] The 35x figure is from the earlier post's internal benchmarks of the end-state architecture. The draft says explicitly that it is not a bridge result.
- [ ] Frontmatter: related is set to [metastable-failures, etcd]. Both are live. One word in the skeleton's header comment was changed from 'TODO(author) line' to 'placeholder line' to clear the author-placeholder warning. draft: true is kept, and no verified: key was added.
- [ ] Quotes are short and credited: 'converges on a single point in our architecture'; 'adding replicas to absorb read load makes writes slower'; 'enough time to migrate safely, but not enough to do it slowly'; 'move the data first, then evolve the architecture once it arrived'; 'don't serialize what doesn't need to be'; 'won't take it to the next billion and beyond'; 'without anyone agreeing on them'. Check each against the page as published.

**Fact-checker could not confirm:**
- [ ] The rollout reading (line 52 and line 134: the posts do not say how much of the per-reference quorum or the virtual replicas carry production traffic) is fair. The Source says 'we're changing', 'The first step is' and 'In internal testing'. However, it also says in the present tense that these changes 'get the last remaining repositories into Azure'. The Author may want to soften or keep the reading.
- [ ] The illustration of 4 round trips x 25ms = about 100ms, about 10 writes/s (line 66) is labelled as our arithmetic. The figure of four round trips is confirmed verbatim in Stretching Spokes #reducing-round-trips, but it describes the 2017 protocol. The Bridge post says only 'multiple'. The Author should decide whether to keep it.
- [ ] The figure of 0.00001% of about a billion, on the order of a hundred (line 134), is labelled as our arithmetic. The arithmetic is correct, and 'a billion repositories' is now linked to the earlier post. GitHub gives no count.
- [ ] The gloss 'the old region stops receiving requests for that repository' (line 80) is a fair reading. The Source says 'stops receiving requests entirely' inside a per-repository process ('Every repository goes through this process').
- [ ] Line 118, 'The quorum also still waits for every healthy member up to a timeout, and the virtual write replica is now one of them', combines two sourced statements. One is the wait behaviour of the current quorum; the other is that the virtual replica is a quorum member whose durability gates the acknowledgement. The Source does not say this explicitly of the bridge quorum.
- [ ] The definition of a reference as 'a named pointer to a commit, such as a branch or a tag' (line 76) was confirmed only as reachable (HTTP 200) on Pro Git 'Git References'. Its wording is a paraphrase of the chapter, not checked line by line. The same applies to the git-update-ref docs, the Pro Git 'Git Objects' page and the git/git compare link. The compare link is the one Stretching Spokes itself links.

<details><summary>What the fact-checker changed</summary>

- Fixed: Line 126: said full replicas, the write node, the read nodes and the object store were 'all under one quorum'. Virtual read replicas do not join the quorum. It now says the full replicas, the ephemeral write node and the object store are inside the quorum and the read nodes are outside it.
- Fixed: Line 36: said virtual read replicas keep 'no durable state'. It now uses the Source's word: 'no persistent state' (Bridge post intro).
- Fixed: Mermaid diagram: it routed fetches only to the Azure replicas and the virtual read replicas. The Source says Spokes routes reads to a replica in the same region and zone, which includes the full replicas in GitHub's data centers, so a 'Spokes -->|fetch| DC' edge was added.
- Fixed: Line 74: 'the Lead leaves the end state...' used our internal term for the Source. It now reads 'the new post', with a link to #h-what-comes-next.
- Fixed: Line 100: an inference was marked 'we read "for now" as a stopgap'. It now says 'we infer that "for now" marks a stopgap', per the rule that inferences say 'we infer'.
- Fixed: Line 108: the cache-miss failure behaviour was presented as a property of the bridge. It comes from the earlier post's description of the new architecture, and the sentence now says so.
- Fixed: Unlinked claims and quotes now link inline at the most specific anchor. These cover the TL;DR price and read-replica bullets; the 'converges on a single point' quote and the CI fan-out sentence; the repository-wide lock; 'says only "multiple"'; falling further behind; the how-we-got-here quotes, the prototyping sentence and the Tier 0 claim; the 'don't serialize' quote; the Aug 2025 Central US, two-of-five and half-of-traffic figures; per-reference waiting; the latency of the virtual write replica; compaction 'for now'; reads and writes competing for resources and the memory example; read replica warm-up; the quorum membership of the virtual write replica; in-line reconciliation; the 'next billion' quote; 'a billion repositories' (earlier post); the last repositories moving and the old region stopping; the Universe final post; the 'without anyone agreeing' quote; the small consensus group.
- Removed: Line 142: 'raised write throughput without new hardware'. No Source mentions hardware, so 'without new hardware' was cut.

</details>

## Prometheus: how Meta runs a gigawatt AI cluster as one machine

`content/explainers/ai-ml-systems/meta-prometheus-supercluster.md` · 2745 words · fact-checker: 96 claims checked, 13 fixed, 1 removed

**Drafter asks you to check:**
- [ ] Spare-GPU arithmetic in 'The first cost' is ours. 32,768/64 = 512 racks, so 523 means 11 reserve racks (about 2%). The 4,888 GPUs outside the active slice are 13% of the 37,656 provisioned but about 15% on top of the 32,768 active. So Meta's 'about 13% above the active footprint' only works as a share of the total, and the draft says so. Check this wording is fair to Meta.
- [ ] Interruption arithmetic in 'The evidence' is ours. It assumes '8K GPUs' means 8,192 and that interruptions scale linearly with GPU count. From that: the peak works out to roughly one interruption per 200 GPU hours, 54k MTBI to about 3.6 a day per 8K GPUs, and a 32,768-GPU job to one interruption about every 100 minutes (about every 3 hours at the 100k target).
- [ ] 'We infer' that one data-parallel replica sits inside one building fabric. The post says only that the network hierarchy matches the replica hierarchy.
- [ ] 'We infer' the time-scale ordering of the power layers (capacitors and rack batteries fastest, then site batteries and Dimmer, then the utility). The post lists the parts but gives no time scales.
- [ ] 'Light in fibre travels roughly 5 microseconds per kilometre' is our figure, labelled as such, not Meta's.
- [ ] The PAFT (arxiv 2602.00277) and Dynamo Dimmer (arxiv 2605.24461) papers were linked but not read: the arxiv fetch was slow. Every claim about them comes from Meta's post, and Sources calls them 'paper linked from Meta's post'. Confirm the authors are Meta's before treating them as first-party.
- [ ] The Lead was read from the published HTML, not a feed. WebFetch would not reproduce it, so the page was downloaded and its text extracted. The headings have no anchors, so every link goes to the post URL.
- [ ] '16-48 Pbps per region pair' comes from the February 2026 BAG post, which gives it as an example. The October post says only 'multiple Pbps'.
- [ ] Catalina's 'about 60 kW a rack' is our arithmetic (120 kW across the two-rack pair). The step from there to Clemente's 135 kW in one rack is our framing.
- [ ] Wording of the options section: the post is not a decision memo. The three forks (Ariel vs stock GB200, DSF vs NSF, capping GPU power vs hard stops) are framed by us from Meta's stated reasons. Only the Ariel quote ('expected that imbalance to cost performance') is verbatim.
- [ ] Two quotes to check word for word: 'falls out of the 64-of-72 slice geometry rather than being an absolute reliability tax' and 'a primary engineering discipline'.
- [ ] Rewrote the skeleton's frontmatter comment, which contained the literal TODO(author) marker. The author note under the first H2 was removed because the task said to replace every TODO; this check list replaces it.
- [ ] related: is left commented out. No live Explainer was a clear fit; the Brief that links this piece will be the inbound link.

**Fact-checker could not confirm:**
- [ ] arXiv 2602.00277 (PAFT) and 2605.24461 (Dimmer) list no affiliations, so they could not be confirmed as first-party Meta papers. Ahmed Sharif, a co-author of Meta's post, is an author of the PAFT paper, and Chunqiang Tang is on both. Every claim the draft makes about them is taken from Meta's post, which confirms it.
- [ ] The Dimmer paper's abstract never mentions 'Dynamo Dimmer' or GPU power capping. The link mirrors where Meta's own post points, so it was kept.
- [ ] The Prometheus post had to be downloaded and its text extracted with curl, as the drafter did. Every number, quote and name was checked against that text, and it matched apart from the fixes listed.
- [ ] Our own arithmetic and inferences are labelled in the text, and all of it checks out: 523x72=37,656, 4,888 is 13.0% of the total and 14.9% of the active, 11/512 is about 2%, 8,192x24 is about 196,600 GPU hours, 54k MTBI gives about 3.6 a day, and 32,768 GPUs gives about 99 minutes. They are kept for the Author to judge. Meta's 'about 13% above the active footprint' only works as a share of the total, and the draft says so fairly.
- [ ] 'A slower step costs less than a restart' and 'which fits a design that sends little traffic between replicas' are our analysis, not Meta's words. Neither is marked 'we infer'. Both were left as editorial reasoning.

<details><summary>What the fact-checker changed</summary>

- Fixed: TL;DR: 'reliability work was hardware maturing, not software' overstated it, because the post also credits firmware and software fixes. Now: 'how much of the reliability work was maturing first-generation hardware'.
- Fixed: Ariel fork: 'kept the spec uniform as the region grew into hardware where 800G was standard'. The post says Ariel 'gave us more uniformity of spec as the region expanded over time and 800G scale-out became available'. The wording now follows that.
- Fixed: Step 2: the 'second dedicated data network' link pointed at the March 2024 GenAI infrastructure post, which never describes a back-end network. It now points to the October 2026 Prometheus post, which does. The 2024 post is no longer cited, so it was dropped from Sources.
- Fixed: Step 2, NSF: the switch details (51.2T, Minipack3 with Dynamic Load Balancing, Minipack3N with Adaptive Routing) come from the Prometheus post. The OCP 2025 post does not mention Dynamic Load Balancing. A link to the Prometheus post was added on those details.
- Fixed: Step 3: 'Building fabrics use shallow-buffer switches' was too broad. The BAG post says this of the shallow-buffer NSF switches, and DSF is not shallow-buffer. Now: 'NSF building fabrics use shallow-buffer switches'.
- Fixed: Step 4: 'monitors NVLink switch health at millisecond resolution'. The post says 'fabric health' under the NVLink Switch and Fabric Health heading. Now: 'NVLink switch and fabric health'.
- Fixed: Step 5: 'Meta's answer is layered by time scale' stated as fact the time ordering that the draft itself later marks as our inference. Now: 'Meta's answer comes in layers.'
- Fixed: Step 5: 'Beneath it ... facility-level battery storage absorbs more, alongside close coordination' is not what the post says. It says Dimmer is 'reinforced by' BBUs, CBUs, BESS and 'deep coordination with utility providers'. The sentence was rewritten to match.
- Fixed: Evidence: 'no breakdown of interruptions by cause, no cost or energy figures, no comparison with H100 clusters' was inaccurate. The post does give qualitative top causes, per-rack kW figures and an H100 compute comparison. Now: 'no numeric breakdown of interruptions by cause, no cost or total energy figures, no reliability comparison with Meta's earlier H100 clusters'.
- Fixed: Sources: the Grand Teton post's real title is 'OCP Summit 2022: Open hardware for AI infrastructure' (18 October 2022). It was listed as 'OCP Summit 2022: Grand Teton'.
- Fixed: Sources: the PAFT entry now uses the paper's real title, 'Training LLMs with Fault Tolerant HSDP on 100,000 GPUs' (Omkar Salpekar et al., January 2026), and notes that Meta's post links it for PAFT.
- Fixed: Sources: the Dynamo Dimmer entry now uses the paper's real title, 'Provisioning to Runtime Optimization of a 100 MW-Scale AI Cluster' (Ehsan K. Ardestani et al., May 2026), and notes that Meta's post links it for Dynamo Dimmer.
- Fixed: Sources: the DSF post's authors (Ron He and Ankur Singh) were added. The closed-loop cooling title now uses its real capitalisation, credited to Meta Newsroom, 27 August 2026.
- Removed: Step 5: 'facility-level battery storage absorbs more' was removed. The post only lists BESS among the layers that reinforce Dimmer.

</details>

## tsrs: how a line-by-line Rust port of TypeScript cuts memory

`content/explainers/languages-runtimes/cloudflare-tsrs.md` · 2313 words · fact-checker: 118 claims checked, 11 fixed, 1 removed

**Drafter asks you to check:**
- [ ] Who built it: the Brief calls tsrs 'Cloudflare's port', but Max Schwenk at Owner.com built it, and Cloudflare, VoidZero and Oxc took it over. The title and TL;DR say so, and the filename slug still says cloudflare-tsrs.
- [ ] Default mode drift: the Lead (20 Oct) says tsrs ships an adaptive configuration by default. Release 0.20.0 (26 Oct) made TypeScript's checker settings the default and put the adaptive settings behind --fast. The draft follows 0.20.0.
- [ ] The Lead's benchmark charts are images with no alt text, so no number in the draft comes from them. The 8x and ~70% figures are from the Lead's own text.
- [ ] The repo moves daily, so every GitHub link is pinned to commit 4a406ed (read 2026-10-27). bench/RESULTS.md was regenerated that day; its ranges (1.94x-4.19x, 1.90x-6.56x with --fast, 1.85x-2.60x memory, 19 projects) come from that version. Owner's post says 17 projects.
- [ ] First-night numbers mix two Sources. Owner's post gives 78% byte-identical, then 99.97% (all but four of 13,462). STATUS.md's first run, 30 Sep, counts differently (10,479 pass of 15,197 variants). The draft cites Owner's percentages and STATUS only for the date and the first-run memory (20.5 GB tsrs vs 18.0 GB tsgo, one checker).
- [ ] Owner's post has interactive figure cards whose numbers came through the scrape garbled (19.53 GB, -21% and others). They were left out. Every Owner number used comes from its prose or its Olympus table.
- [ ] Inferences labelled 'we infer': memory rather than speed was the binding constraint for agents; the arena/handle model is why the port could go file by file; the gap between the Lead's 8x and the results page is mostly the faster 7.1 baseline. One point the Sources do not make is unlabelled: the 'tsrs takes both points literally' paragraph links Cavanaugh's 'forego garbage collection' remark to tsrs's leak arenas. Rewrite it as an inference or confirm it.
- [ ] Quote wording to confirm: Cavanaugh's 'a structurally similar codebase', 'might be a better fit for a ground-up rewrite' and 'can effectively forego garbage collection entirely'; PORTING.md's 'mechanical and faithful', 'improvements are bugs' and 'by decree, so the compiler does not police sharing'; README's 'appear, disappear, move or have different message text'.
- [ ] Owner's post says the repo bar requires a 1% instruction, 2% wall or 5% memory cut. The current AGENTS.md has dropped fixed thresholds. The draft leaves this out.
- [ ] RUST.md says the nightly 200,000-mutant fuzz workflow is now disabled; Owner's post describes it as running. The draft cites only 'the first 320,000 mutants' from RUST.md and does not call it nightly.
- [ ] Is building.owner.com first-party? Max Schwenk ported tsrs, so the draft treats his company's engineering blog as first-party. The post shows no explicit date ('Where things stand, October 19'), so Sources gives none.
- [ ] To pass author-placeholder, the skeleton's frontmatter comment was reworded to drop 'TODO(author)'. It now says the check list is in the pull request. Nothing else in the frontmatter changed beyond title, description and tags. draft: true kept, no verified:.

**Fact-checker could not confirm:**
- [ ] Ryan Cavanaugh as author of microsoft/typescript-go discussion #411: the page is 'Why Go?', dated 2025-03-05 and mentions RyanCavanaugh 13 times. The scrape did not isolate the author field. All three Cavanaugh quotes are confirmed verbatim.
- [ ] 'We infer the gap between these and Cloudflare's 8x is mostly the faster 7.1 baseline' is labelled as an inference but is weak. RESULTS.md uses default checker counts, while Cloudflare compared at equal checker counts. That difference may matter as much as the baseline. The Author may want to reword or drop it.
- [ ] 'Rust's ownership rules fight that shape' (Step 2) is general background with no Source. It was left as framing.
- [ ] 'That layout is what made an AI-driven port tractable' is a causal framing only loosely supported. Owner says signatures were written first so agents could work in parallel worktrees. It was left as is.
- [ ] 'The unsafe arena code runs under Miri' rests on Owner's post and on RUST.md (focused arena checks in miri.yml). RUST.md lists 'Miri on all of tsrs_core' as not yet tried, so the claim holds only for the arena code. The wording already says arena code.
- [ ] The STATUS.md first run counts its private monorepo as 38k files (37,942). Owner gives Olympus as about 37,000 files. They are probably the same codebase, but no Source says so. The draft keeps them separate.

<details><summary>What the fact-checker changed</summary>

- Fixed: Hejlsberg quote: the draft quoted '"10x faster"' as his words. The phrase appears only in the post's headline, so the sentence now quotes the headline verbatim: "A 10x Faster TypeScript".
- Fixed: 'tsrs takes both points literally' paragraph (the drafter flagged it): no Source links Cavanaugh's remark to tsrs's arenas. It now reads 'We infer that tsrs takes both points literally' and 'in the spirit of Cavanaugh's remark', replacing 'the way Cavanaugh describes'.
- Fixed: Schwenk quote: the draft said Go's memory was "millions of small objects...", dropping the word 'mostly' that the Owner post has. It also called this 'the deciding argument', which the post never says. It now says Go's memory was mostly "millions of small objects in a garbage-collected heap", attributed as Schwenk's own reasoning 'as he tells it', and that Rust 'would let him' choose size and lifetime.
- Fixed: PORTING.md rules for agents: the draft said they 'begin' with "Read the Go function, write the Rust function." They actually begin with 'The Go source is the specification.' Changed 'begin' to 'say'.
- Fixed: No lifetime parameters: PORTING.md names one exception (issue #646, DefaultResolver<'h>). The text now reads 'Apart from one documented exception, there are no lifetime parameters'.
- Fixed: Fast-mode scheduling: the README says fast mode groups nearby files on one checker. Module grouping starts only at 16 checkers (8 for programs of 24,000 or more files), and idle checkers 'take unstarted work'. The 'from the busiest' detail came from Owner's post, not the README. The sentence now follows the README. The Olympus 11%/16% figure is now tied to 'grouping files by the modules they import', which is what Owner describes.
- Fixed: WebAssembly: 'On 32-bit WebAssembly, P<T> is a plain pointer' is not in the linked wasm-build note. Removed it and kept 'Every layout assertion has a 32-bit twin', which the note supports.
- Fixed: 'This is the part Cloudflare cares most about' overstated the post. It now says 'the part Cloudflare's post dwells on'.
- Fixed: Oxlint: 'Cloudflare plans type-aware lint rules in Oxlint next' overstated the post, which says it is 'excited to experiment with deeper tsrs integrations ... for example for type-aware lint rules in Oxlint'. It now says Cloudflare wants to experiment with tsrs in type-aware lint rules for Oxlint.
- Fixed: The reader-facing body used the internal term 'the Lead's 8x'. It now says 'Cloudflare's 8x'.
- Fixed: Sources: the Owner post is dated 10/19/2026 on the page, so the Sources entry now gives '19 October 2026'. The drafter had said the post showed no date.
- Removed: 'On 32-bit WebAssembly, P<T> is a plain pointer' (Step 5): the linked wasm-build note does not say this. PORTING.md says only that P is a plain pointer 'on targets without mmap', and nothing ties that to wasm32.

</details>

## JVM warmup: why Java starts slow and how AOT caches shorten it

`content/explainers/languages-runtimes/jvm-warmup-and-aot-compilation.md` · 2414 words · fact-checker: 96 claims checked, 6 fixed, 1 removed

**Drafter asks you to check:**
- [ ] Hook: in JEP 545's scenario 1 the roughly 80 s warmup comes from ZGC committing memory and taking page faults, not from the JIT. The opening says this, and says JEP 544 is aimed at 'the other part of warmup: compiling the code'. Check that the framing does not suggest the 80 s was JIT time.
- [ ] The Inside Java post shows no date on the page. '16 October 2026' is taken from the URL path (/2026/10/16/).
- [ ] Two links use the anchor '#What-about-static-compilation?', with a literal question mark, because that is the id on the JEP 544 page. Check that it scrolls to the right place on the Preview.
- [ ] Bylines in Sources are the JEP header's Owner field for JEPs 545, 544, 514, 516, 295 and 410. JEP 483 and JEP 515 list Authors, and those are used instead. Decide whether an Owner is fine to credit as the author.
- [ ] JEP 544's performance numbers (around 50–70% without AOT code and 65–80% with it, on five framework benchmarks on two cores; javac first iteration about 30%, then about 75% total) are read from the prose under its charts. The benchmark apps are not named because they appear only in the chart images.
- [ ] JEP 545 uses the flag -XX:+ZAdaptiveHeapSizing, which comes from JEP 546 (adaptive heap sizing). That JEP is not cited here.
- [ ] Lines that say 'we infer': (1) a load-balancer ramp means a new instance carries less than its share until it is warm; (2) Leyden's AOT code differs from jaotc because it comes from HotSpot's C1 and C2, guided by a profile, and sits in the same cache as the classes.
- [ ] JEP 516 is summarised as 'before JDK 26, ZGC users could not enable ahead-of-time class loading and linking'. Check this against JEP 516's Motivation wording.
- [ ] The TL;DR says the cache is tied to one JDK release (JEP 483's consistency rules), while the CPU-feature and GC constraints apply only to AOT code (JEP 544). JEP 483 allows the garbage collector to differ for the class cache.
- [ ] The closing points to the JDK 28 early-access page, which notes that a recently targeted JEP may not be in the current build (build 21, 2026/10/22). The piece does not claim JEP 544 is in a published EA build.
- [ ] Before the description field, the frontmatter still says '# published and reviewed are the pick time, provisional' and keeps the commented hero and related lines. No related slugs were added, and verified: was not set.

**Fact-checker could not confirm:**
- [ ] The date '16 October 2026' for the Inside Java post: the page shows no date, and the date comes only from the URL path /2026/10/16/. Also, JEP 545 now shows Status 'Integrated', not 'Targeted'. The sentence only says the post announced the targeting, so it is still accurate.
- [ ] Bylines: JEPs 545, 544, 514, 516, 295 and 410 credit the header's Owner field. JEP 483 credits its Authors (Ioi Lam, Dan Heidinga, John Rose) and JEP 515 its Author field (Igor Veresov, John Rose). The Author must decide whether crediting an Owner as the author is acceptable.
- [ ] The '#What-about-static-compilation?' anchor matches the id on the JEP 544 page exactly. Whether it scrolls to the right place on the Preview was not tested in a browser.
- [ ] JEP 544's percentages (about 50-70% and 65-80%; javac about 30% and then about 75% in total) match the prose under its charts. The chart images themselves were not checked.
- [ ] The closing line 'What could go right is a new instance that serves its first ramp at close to the speed of the ones already running' is our own forward-looking framing. No single Source states it. It is worded as a possibility, and it is kept.
- [ ] 'What does not count: a slow first response caused only by scanning JARs and loading classes is a startup cost, not warmup' follows from the startup and warmup definitions in JEPs 544, 545 and 483. No Source says it in those words.

<details><summary>What the fact-checker changed</summary>

- Fixed: Opening: the line said JEP 545's future work names ahead-of-time code compilation as 'a second change for JDK 28', but that section does not say JDK 28. The sentence now also links JEP 544, whose header says Release 28 and Status Integrated.
- Fixed: Opening: the line did not say what caused the 80 s warmup, so a reader could take it for JIT time. It now says the P99 spikes came 'while ZGC grew the heap towards 20GB', which JEP 545 Scenario 1 supports. This answers the drafter's framing concern.
- Fixed: TL;DR and 'Small containers': these said the JIT 'competes' with the application for CPU on two cores. JEP 544 Performance says it 'is likely to compete', so both now say 'is likely to compete'.
- Fixed: JEP 515 HelloStreamWarmup: the 19% gain was put down to C2 starting 'without waiting for the interpreter and C1'. The JEP does not say that. It now says 'because the JIT ran earlier and with more accuracy', following JEP 515's Description.
- Fixed: 'It is not static compilation': the line said initial experiments 'produced' overly large cache files from near-total reliance on AOT code. JEP 544 Future Work says experiments 'suggest that minimizing interpreter use results in overly large AOT cache files', and the line now says that.
- Fixed: Observability: the line said that without -XX:AOTMode=required HotSpot 'quietly' falls back to the JIT. JEP 544 says HotSpot issues a warning, so the line now says it falls back 'with only a warning'.
- Removed: The line 'No closed world is assumed, and reflection keeps working' had no link and was not said about JEP 544. It is replaced by a short sourced quote: applications keep 'the agility, portability, and compatibility of dynamic compilation' (JEP 544, Shifting compilation work to training runs).

</details>

## Off-heap memory in Java: data the garbage collector does not manage

`content/explainers/languages-runtimes/java-off-heap-memory.md` · 2630 words · fact-checker: 72 claims checked, 7 fixed, 1 removed

**Drafter asks you to check:**
- [ ] Benchmark figures (15.862 vs 1.052 ns/op at 5 bytes on an Apple M4, speedups of 6.8x to 18.9x, parity above 64 bytes) are from the Lead. JDK-8385697's earlier experiment shows different numbers (15.734 vs 2.335 ns, about 6x). The draft cites only the Lead's figures.
- [ ] The 256-bytes-per-thread retention bound and the 'three native operations' framing come from CSR JDK-8390770, not the Lead. Check that the CSR's final text still says this.
- [ ] The Lead's 'more than 99.99% of confined arenas ... used less than 64 bytes' came from one instrumented test-suite run, and the Lead itself says it is not representative. The hook keeps that qualifier.
- [ ] The 'we infer' statements to check: an automatic arena behaves much like a direct ByteBuffer's memory, and JDK 28's pooled arena is the slicing-allocator idea built into the JDK.
- [ ] JEP 471's removal timeline ('Phase 3 ... in JDK 26 or later') is a plan stated in the JEP. Check whether JDK 26 or 27 already changed the default of --sun-misc-unsafe-memory-access, since the draft does not claim the current default.
- [ ] Both Java snippets are our own code, not copied. The errno example is ours and describes the call writing into the segment. Real errno capture in FFM goes through Linker.Option.captureCallState, so check that the comment does not mislead.
- [ ] The JDK 28 Arena and java.lang.foreign links point at early-access draft docs (28-ea+21). Re-point them to GA docs when JDK 28 ships.
- [ ] Lucene is cited from the 10.0.0 MMapDirectory javadoc. The Netty anchor #troubleshooting-buffer-leaks is a GitHub-wiki-style anchor, so confirm it resolves on netty.io.
- [ ] related: [structured-concurrency] was uncommented and added, because the Lead's section on virtual threads borrowing carrier pools links the two. Drop it if the link is too thin.
- [ ] The H2 'One confined arena, from open to close' runs about 283 raw words including its code and Mermaid fences. The build did not warn about section-length, so the build likely does not count code, but check it.

**Fact-checker could not confirm:**
- [ ] The 15.862 and 1.052 ns/op figures for 5 bytes on the Apple M4 match the Lead's listing, but PR #31365's own summary gives about 16 to 19.8 vs 1.1 to 1.2 ns/op for the M4. Nothing changed, because the draft attributes the numbers to the Lead.
- [ ] JEP 471 Phase 3 (exceptions by default 'in JDK 26 or later') is still a plan. The JDK 26 and JDK 27 JEP lists contain no Unsafe JEP, so 'plans' looks current. A non-JEP change to the default in 26 or 27 was not ruled out.
- [ ] The JDK 28 Arena, Linker.Option and java.lang.foreign links point at early-access draft docs (28-ea+21). Re-point them to GA docs when JDK 28 ships.
- [ ] The #wiki-h2-9 anchor on the Netty wiki is a generated heading id and could change if the wiki is regenerated.
- [ ] 'Since JDK 22 the standard way' and 'print a run-time warning on first use since JDK 24' are light extrapolations of the JEPs (final in 22; warning by default in 24). They are reasonable but not worded that way in the Sources.
- [ ] The related: [structured-concurrency] slug exists under content/. Whether the link is too thin is an editorial call left to the Author.

<details><summary>What the fact-checker changed</summary>

- Fixed: Hook: the draft said the instrumented test-suite run happened 'during the development of JDK 28'. The Lead says only 'one instrumented test-suite run', so that phrase is removed.
- Fixed: Errno example: the draft introduced it only as 'ours', with a plain JAVA_INT segment passed to the call. The intro now says the example is ours and simplified, and that real code links the handle with Linker.Option.captureCallState, which saves errno into a segment passed to the call. This is linked to the JDK 28 EA Linker.Option docs, anchor verified.
- Fixed: Netty: the draft said 'pools its buffers, direct ones included'. The Netty wiki never mentions direct buffers, so it now says 'pools its buffers'.
- Fixed: Netty leak link: the anchor #troubleshooting-buffer-leaks has no matching id in the netty.io HTML. The page's ids are wiki-h2-9 and user-content-troubleshooting-buffer-leaks, so the link now uses #wiki-h2-9.
- Fixed: JEP 471 timeline: 'remove the methods in releases after that' now reads 'in releases following JDK 26', which is the JEP's wording.
- Fixed: PR #31365: the draft stated as fact that the PR holds the benchmark configuration and raw results, but the PR page shows no JMH settings or raw runs. It now says the Lead puts its configuration and raw results in the PR, which is what the Lead says.
- Fixed: Closing section: the draft suggested running with --sun-misc-unsafe-memory-access=warn to 'see which libraries still reach for Unsafe'. Per JEP 471, warn prints at most one warning in total, so it now says =debug, which prints a warning and a stack trace on every use. This matches Practice four.
- Removed: 'of Oracle' after Per-Ake Minborg. The Lead's byline gives no affiliation; only the site footer names Oracle Corporation. The Sources entry 'Inside Java, Oracle' stays, because the footer supports it.

</details>

## Open table formats: how Iceberg turns files into a table

`content/explainers/databases-storage/apache-iceberg-table-format.md` · 2723 words · fact-checker: 112 claims checked, 6 fixed, 1 removed

**Drafter asks you to check:**
- [ ] Title and frontmatter: I retitled the piece to the concept ('Open table formats: ...'), but the file slug stays apache-iceberg-table-format, as the task required. Decide whether the slug should change to match.
- [ ] The phrase 'open table format' comes only from Iceberg's home page tagline. The meaning of 'open' (a published spec many engines can implement) is marked 'we take that to be'. No Source defines the term outright.
- [ ] The claim that the DuckDB file name 00003-<uuid>.metadata.json follows the spec's metastore-table naming scheme is our connection. The spec describes that scheme for metastore/database pointers, and the Lead uses a REST catalog fixture. Check that the wording ('the naming the spec gives for tables whose pointer lives in a metastore or database') does not overstate it.
- [ ] Sequence-number example: we apply the spec's scope rules to the DuckDB table (the delete at seq 3 applies to the 2 Oct data file at seq 1). The manifest sequence numbers 1/2/3 come from the Lead's manifest-list output. The rule comes from the spec's Scan Planning section.
- [ ] 'each of those files spans roughly 1 to 9,000': the Lead's actual bounds are 3-9000, 1-8998 and 2-8999. The 4 Oct file is 9,001-12,000.
- [ ] Lead table is format-version 2 and uses a position delete file. The spec says position delete files must not be added to v3 tables and are deprecated in v3. Check this reads consistently.
- [ ] Iceberg v4 claims (root manifest replaces manifest list, adaptive metadata tree, 'under active development and has not been formally adopted') come from the live spec page, which may change before publication.
- [ ] Quotes to check word for word: spec 'a large, slow-changing collection of files in a distributed file system or key-value store as a table'; 'This table format tracks individual data files in a table instead of directories.'; 'reused across snapshots'; 'only to data files that are older and in the same partition'; 'have no requirements and can always be applied'; 'unsafe in object stores and local file systems'; 'is not standardized by this spec'. Docs: 'a 10x performance improvement'; 'an index over a table's data'. DuckDB: 'a convenient hint rather than a guarantee'; 'That move is the commit.'; 'Writes are cheap, and reads pay a small cost, which grows with the number of delete files.'; the hook quote about every change writing new files. Hudi: 'a log of all actions performed on the table at different instants'.
- [ ] 'We infer' sentences: the shared idea across Iceberg, Delta and Hudi (tree vs log); why the spec is exact about sequence numbers (sharing across engines); the 'format describes state, engines do the work' boundary.
- [ ] Delta Lake is cited from PROTOCOL.md on the master branch, which may change. Consider pinning a commit or tag. The Hudi docs URL redirected to /docs/timeline/ (cited at that address), and the page is for version 1.2.1.

**Fact-checker could not confirm:**
- [ ] Delta Lake is cited from PROTOCOL.md on the master branch. Its Overview, Delta Log Entries and Checkpoints text matched today (2026-10-27), but it can change before publication. Pinning a tag or commit is the Author's call.
- [ ] Iceberg v4 statements ('under active development and has not been formally adopted', single root manifest per snapshot, adaptive metadata tree for efficient small commits) match the live spec today, but that spec page changes with v4 work.
- [ ] 'Its stated aims match Iceberg's' (Delta's serializable writes and snapshot isolation compared with Iceberg's serializable isolation) is the draft's own comparison of two goal lists. Each half is sourced, but 'match' is a judgment call.
- [ ] 'each of those files spans roughly 1 to 9,000': the Source bounds are 3-9000, 1-8998 and 2-8999, so 'roughly' is accurate. It is left as is.

<details><summary>What the fact-checker changed</summary>

- Fixed: Partition-skipping example: it listed the three manifests as covering '1 to 3 October, 2 October, and 4 October' and then said a 4 October query rules out 'the first and third'. In that order, that would skip the 4 October manifest. DuckDB counts in write order (seq 1 = 1-3 Oct, seq 2 = 4 Oct, seq 3 = 2 Oct). The text now reads 'the three manifests, in the order they were written, cover 1 to 3 October, 4 October, and 2 October'.
- Fixed: Retry rules: 'deletes and compactions first check that the files they remove are still in the table' is now 'deletes of specific files and compactions ...'. The spec says only deletes of specific files must verify; deletes based on an expression can always be applied.
- Fixed: Commit step: 'moves the atomic step to a catalog's check-and-put' is now 'describes the atomic step instead as a check-and-put on a pointer in a metastore or database', with a link to spec #metastore-tables. That is the spec's wording; it does not say 'catalog'.
- Fixed: Hive-style folders: 'the design Iceberg was built to replace' is now 'the Hive-style design whose correctness problems Iceberg was designed to solve'. The Reliability doc says Iceberg 'was designed to solve correctness problems that affect Hive tables running in S3', not that it replaces them.
- Fixed: First Mermaid diagram: the dotted edge from the position delete file to the 1-3 Oct data files now reads 'skip these row positions in the 2 Oct file'. The DuckDB delete file only names the 2 October data file.
- Fixed: Same section: dropped the lead-in 'What does not count:' so the section is back under the 250-word guide. The build had warned at 252 words after the Hive-style fix.
- Removed: 'and that compaction work is lost' (Practice two). The Reliability doc says only that the compaction 'must fail' if a concurrent commit deleted file_a.avro or file_b.avro. It says nothing about the work being lost.

</details>

## HDR histograms: measuring tail latency with bounded error

`content/explainers/distributed-systems/hdr-histograms-tail-latency.md` · 2660 words · fact-checker: 78 claims checked, 5 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Hook numbers come from wrk2 README 'Example 2'. Uncorrected: median 3.02 ms, p99 6.04 ms, mean 6.588 ms. Corrected: median 8.61 ms, p99 1.27 s, mean 63.66 ms. Both histograms come from one --u_latency run. The README's own '200x' is loose: 1.27 s / 6.04 ms is about 210x.
- [ ] The worked example (5,000 µs -> bucket 2, sub-bucket 1,250, slot 5,000-5,003, 0.08%; worst case about 0.098% at 4,096) is our own arithmetic from AbstractHistogram.java. The source comment there says bucket 1 covers '2048..4097', which looks like a typo; the draft says 2,048 to 4,095.
- [ ] The claim that the reported percentile is 'never below the true value, and at most one step above' is our reading of getValueAtPercentile returning highestEquivalentValue, and the draft labels it as such. Please check it against the code.
- [ ] The line-anchored permalinks into AbstractHistogram.java at commit de84b0a (L318-L333, L359-L372, L725-L748, L1432-L1454, L2420-L2437) were taken from master on 2026-10-27. Check that the ranges still frame the right code.
- [ ] Recording cost: the HdrHistogram README says 3-6 ns on 'circa 2012' Intel CPUs, but the github.io project page says 'circa 2014'. The draft cites the README's 2012.
- [ ] Prometheus figures (443 ms vs about 320 ms; summary range 270-330 ms) are from the current 'Histograms and summaries' page. Older versions of that page gave different numbers (442.5 ms), so recheck before publishing.
- [ ] The Lead's title says 'exact percentiles' while its body and the README say about 0.1% bounded error. The 'what it is not' section points this out; check the tone is fair to the project.
- [ ] Inferences marked 'we infer': the mean described almost no request in the paused run; the closed loop is the force behind coordinated omission; correction must happen before a histogram reaches pg_hdr_histogram, because its function list has no corrected-recording call.
- [ ] The pg_hdr_histogram figures (42 B / 12 KB blob sizes, about 110 µs per merge, default range 0..1e9 at 3 significant figures, PostgreSQL 14-18, HdrHistogram_c 0.11.10) are the README author's own dev-machine measurements and claims. The 1.0.0 GitHub release is dated 17 October 2026, and the PostgreSQL news post 26 October.
- [ ] The SQL snippet in Practice two is our own, written against the README's documented functions, so no licence notice is needed. It assumes hdr_add works as a GROUP BY aggregate (the announcement calls it one) and that a latency_minutes table exists, which the snippet does not create.
- [ ] related: [metastable-failures] was uncommented in the frontmatter, and the 'Follow ... replace every TODO(author) line' frontmatter comment was removed, because it tripped author-placeholder. draft: true is kept and verified: was not added.

**Fact-checker could not confirm:**
- [ ] GitHub shows the v1.0.0 release date as '17 Oct' with no year. The draft's '17 October 2026' fits, but the page does not confirm the year.
- [ ] The draft follows the README's 'circa 2012' for the 3-6 ns recording cost. The github.io project page says 'circa 2014'. Both are first-party sources, so the Author should decide which to cite.
- [ ] Bucket 1's range is given as 2,048 to 4,095. The AbstractHistogram.java comment (L359-L372) says '2048..4097', which looks like a typo in the source, so the draft's figure is our correction. The other permalinked ranges at de84b0a still frame the right code (sub-bucket sizing, add(), getValueAtPercentile, getBucketIndex/getSubBucketIndex).
- [ ] 'Never below the true value, and at most one step above' is our own reading of getValueAtPercentile returning highestEquivalentValue of the slot that reaches the ceil'd target count. The code supports it, and the draft labels it as our reading.
- [ ] The worked example (5,000 >> 2 = 1,250, slot 5,000-5,003, 0.08%; worst case 4/4,096 = about 0.098%) is our own arithmetic. I checked it against the source.
- [ ] The wrk2 paragraph names 'Michael Barker's C port', but the wrk2 README it links says 'Mike Barker'. The pg_hdr_histogram README uses 'Michael Barker', so I kept it.
- [ ] The HdrHistogram README's 1 µs to 1 ms / 1 ms to 1 s / 1 s to 1,000 s tiers come from its first example paragraph. The phrase 'always accurate to 3 decimal points' closes a second paragraph about the same histogram whose tiers differ slightly (~2 ms, ~2 s, 2,000 s). The pairing is fair but not one passage.
- [ ] 'SubBuckets = 2048' is quoted with the wrk2 output's column padding collapsed.
- [ ] The SQL snippet assumes hdr_add works as a GROUP BY aggregate. The announcement says 'Aggregates hdr_add and hdr_merge', which supports this. It also assumes a latency_minutes table the snippet does not create.

<details><summary>What the fact-checker changed</summary>

- Fixed: Hook: 'an Apache server on his own machine' is not in the wrk2 README, which shows only 127.0.0.1. It now says 'a local Apache server'.
- Fixed: Hook: 'Measured the way most load generators do' overstated the README's 'similar to the model found in many current load generators'. It now says 'many'.
- Fixed: What it is not: 'the same sentence gives ~0.1% relative precision' was wrong. The '~0.1%' line is in the announcement's body, not in its title. It now says 'its body gives', and 'The honest claim' is softened to 'The precise claim' to be fairer to the project.
- Fixed: What it makes possible: wrk2's spectrum was said to print 'out to 99.999% and beyond'. The README's printed rows stop near 99.998% and then jump to 100%. It now says 'past 99.99%, up to the maximum'.
- Fixed: Practice one: 'concurrent writers need one of the concurrent variants or a Recorder' left out the docs' first option. It now says 'need external synchronization, one of the concurrent variants, or a Recorder'.

</details>

## ASPA: how BGP routers check the AS path for route leaks

`content/explainers/apis-networking/aspa-bgp-route-leaks.md` · 2700 words · fact-checker: 96 claims checked, 8 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Worked example: APNIC's July post does not state the origin AS. We assume AS1 is the origin so we can apply the draft's ramp arithmetic (up-ramp 1 < N=2 at AS3; 1+1 < N=3 at AS4). Check the ramp counts against draft -29 sections 5.4 to 5.6, and check against APNIC's Figure 1 that AS1 originates.
- [ ] The 'remove AS1's ASPA, the result becomes Unknown' paragraph is our inference from section 5.4's min/max bounds, and is flagged 'we infer'.
- [ ] The closing applies ASPA to the 2019 leak: DQE's ASPA would not list Allegheny, so Verizon would see the path as Invalid. This is our inference (flagged), since ASPA did not exist in 2019. The same goes for the OTC counterfactual about Allegheny's routers.
- [ ] 'In 2019 ROV would have helped only because the prefixes were split' rests on Cloudflare's statement that its prefixes are signed with max length /20. 'A leak of the /20 itself would have been ROV Valid' is flagged as inference.
- [ ] BIRD: verification draft -29 says ASPA lives in a side branch, with release expected after RTR v2. APNIC's July post says there is 'good support' in BIRD. The piece follows the draft. Check whether a BIRD release now ships ASPA.
- [ ] Counts: 538 ASPAs under APNIC's trust anchor (APNIC, 22 Oct 2026) and roughly 2,400 globally as of July 2026 (profile draft -29, section 8.1). Two different scopes and dates; make sure both read correctly.
- [ ] APNIC says 'Since November 2025, the RIPE NCC, ARIN, and APNIC have deployed support'. We paraphrase it as deployed since November 2025, and the wording is ambiguous about who deployed when.
- [ ] The 'one hop between apexes = the single peering a route may cross' reading is labelled as ours, not the draft's wording.
- [ ] The draft URLs are pinned to revision -29 (verification draft dated 16 Oct 2026, profile draft 29 Jul 2026). If newer revisions or RFC numbers appear before publishing, re-anchor them.
- [ ] The Sources entry for the aspa-examples PDF credits K. Sriram, based on the repo path (ksriram25). The draft's reference itself gives no author.

**Fact-checker could not confirm:**
- [ ] The preview build still reports one error against this file that is not caused by it: broken-link to /apis-networking/cloudflare-hyperdrive-mysql-caching/. That sibling draft fails html-size (52,071 bytes, limit 51,200), so it is not written. Every page in the vertical links to it through generated navigation (agentgateway, cilium, dns-provider-redundancy and others report the same error). The fix belongs to the hyperdrive draft. Apart from that, this file has no errors or warnings.
- [ ] Worked example: APNIC's Figure 1 text does not say outright that AS1 originates the route. Figure 3 of the same post does say 'AS1 legitimately originates 192.0.2.0/24', so the assumption holds. I recomputed the ramp arithmetic against draft -29 sections 5.4 to 5.6: at AS3 max_up_ramp is 1, less than N=2, so Invalid. At AS4 max_up 1 plus max_down 1 is less than N=3, so Invalid. Without AS1's ASPA, at AS3 max_up is 2 and min_up is 1, so Unknown. All three results are correct.
- [ ] 'Crosses at most one peering' in the TL;DR is our reading of 'apexes no more than one hop apart' (section 5.2). The body labels it as ours; the TL;DR does not.
- [ ] The BGPsec line 'because a leaked update really did travel that path' is our explanation. Draft section 8.2 says only that BGPsec does not detect route leaks (valley-free violations) and that the two are complementary.
- [ ] The aspa-examples PDF credit (K. Sriram) is confirmed by the PDF's own metadata (Author: Kotikalapudi Sriram; contact ksriram@nist.gov), not by the draft's reference entry.
- [ ] Checked against the datatracker on 2026-10-27: both drafts are still at -29, with no RFC numbers yet. Re-check before publishing.

<details><summary>What the fact-checker changed</summary>

- Fixed: RFC 7908 section 3 calls the hairpin 'a common form' of leak, not the most common, and the 2019 leak split the prefix, so it was not a pure Type 1 'full prefix' leak. The text said 'The 2019 incident was the most common form, a hairpin turn'. It now says 'The 2019 incident took the shape of a common form, a hairpin turn'.
- Fixed: 'Peers exchange only their own customers' routes' had no link. It now links RFC 7908 section 3.2, which says lateral peers 'should only propagate to each other their respective customer prefixes'.
- Fixed: BIRD: the draft -29 implementation list (side branch mq-aspa, release after RTR v2) is out of date. BIRD 2.16 (December 2024) shipped ASPA support, and APNIC's July post says BIRD has 'good support'. The text still reports what the draft lists, now as 'a BIRD side branch', and adds: 'The BIRD entry lags: APNIC's July post already counts BIRD among routers with "good support"', linked to APNIC.
- Fixed: The registry sentence was reworded to follow APNIC's own wording and order: 'Harrison writes that since November 2025 the RIPE NCC, ARIN and APNIC [have deployed support]'.
- Fixed: The text said MyAPNIC suggests providers from BGP observations. APNIC says only that it gives 'recommendations based on BGP observations'. It now reads 'with recommendations based on BGP observations and by checking each change...'.
- Fixed: Appendix B Property 1 also requires that AS B receives the leaked route on a customer or lateral-peer interface. The text said 'leaked towards the other'. It now says 'leaked, and then received by the other from a customer or peer'.
- Fixed: The closing inference said validation 'at Verizon or any network further along' would mark the 2019 path Invalid. Under section 5.6, a Verizon customer running the downstream check would not necessarily get Invalid, because its down-ramp can stretch through ASes that publish no ASPA. It now says 'at Verizon, or failing that at any Verizon peer that validated' (upstream check: DQE's ASPA without Allegheny makes the up-ramp stop short). It is still flagged as inference.
- Fixed: The skeleton's front-matter comment contained the literal text 'TODO(author)'. Because of it, the builder blamed this file for the author-placeholder warnings coming from other drafts' pages (netflix-conductor, briefs and others). The comment now says 'author placeholder line', so this file has no placeholder warnings. The comment is not rendered.

</details>

## Disaggregated inference: why prefill and decode run on separate GPUs

`content/explainers/ai-ml-systems/disaggregated-inference.md` · 2600 words · fact-checker: 96 claims checked, 5 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Hook and worked example: DistServe Figure 1 and the text beside it give 1.6 rps colocated, 5.6 rps prefill-only, 10 rps decode-only, and an 'ideal' 2+1 GPU split at 10 rps, which is 3.3 per GPU and 2.1x. These are the paper's ideal estimate, not a measured deployment. Check the wording 'would, ideally'.
- [ ] The 1.13 GB / 11.3 GB/s / 90 Gbps sizing uses OPT-66B, while the hook uses a 13B model. The text says the sizing uses a larger model; check that this reads clearly.
- [ ] TL;DR '7.4 times more requests' is DistServe's 'up to 7.4x more requests or 12.6x tighter SLO'. Only the first figure is used.
- [ ] Splitwise is cited as arXiv v2, not as ISCA 2024, because I could not source the venue. Its figures are: 60-70% of time at 20 or fewer tokens (conversation); BLOOM-176B 1,500-token prompt equal to 6 output tokens; 8 ms / 5 ms non-overlapped transfer; 3% vs 0.8% of E2E; +16.5% vs +64% on the second token; 27 prompt + 3 token machines at 70 RPS (Splitwise-HH, coding, 2-minute trace).
- [ ] Splitwise says the first token comes from the prompt machine, so the copy shows up as second-token latency. The vLLM Lead says the copy is critical to TTFT. The piece presents both. Check that this does not read as a contradiction.
- [ ] Lead numbers: TTFT about 1,744 ms pull vs about 540 ms push at 8 QPS and 16k input; TPOT about 10.4 ms vs near 7 ms at 4-8 QPS and 16k. The setup is Qwen3-32B, 1P1D, P5en with 8x H200 and EFA, 128 output tokens, prefix caching disabled on decode.
- [ ] 'We infer' claims: (1) pull's inline descriptor work on the decode worker reintroduces, in miniature, the interference disaggregation removes; (2) KV cache tiering and disaggregation are separable ideas. The basis for (2) is that vLLM lists an OffloadingConnector beside its transfer connectors.
- [ ] Mooncake is cited as arXiv v5, dated 8 Oct 2026 on the abs page. Check that the 75% and 525% figures are unchanged in that version. Only the 75% figure is used.
- [ ] DeepSeek: EP32 over 4 nodes for prefill and EP144 over 18 nodes for decode; 73.7k input tok/s (including cache hits) and 14.8k output tok/s per H800 node; 56.3% of 608B input tokens hit the on-disk cache; the window starts 27 Feb 2025 12:00 UTC+8. The Sources entry title for the DeepSeek page is reconstructed, so check it against the page heading.
- [ ] The Dynamo v0.7.1 design doc is an old version (the page says the latest is v1.6.0). The conditional-disaggregation page is the dev version and is marked experimental. Defaults are 2,048 tokens and 0.7, and the example uses 4,096 and 0.6. The claim that the router needs KV events from backends came from a search-result summary of that page, not from my direct fetch, so check it on the page.
- [ ] Sources-list titles for the vLLM docs pages ('Disaggregated Prefilling (experimental)', 'NixlConnector Usage Guide') are reconstructed, not copied, so check them.
- [ ] The 'One request, two pools, one copy' section counts at about 271 words including the Mermaid fence. The build gave no section-length warning, but it is close to the 250-word limit if the fence is not excluded.

**Fact-checker could not confirm:**
- [ ] Mooncake: the article links the arXiv abs page, which serves the latest version. I checked the v5 PDF (arXiv:2407.00079v5, 8 Oct 2026). It still has the 75% (real workloads) and 525% figures, Conductor, early rejection and the CPU/DRAM/SSD cache.
- [ ] Splitwise venue: cited as arXiv v2 (20 May 2024), which matches the PDF header. I did not check the ISCA 2024 venue, and the article does not claim it.
- [ ] Two short sentences are the article's own reasoning with no direct Source statement. They are not marked 'we infer', and I left them in because they follow from the definitions. One is in 'Prefill builds the cache': 'Adding identical replicas does not count either, because every replica still mixes both phases.' The other is in 'What disaggregated inference is not': 'Those split one model across GPUs within a phase; disaggregation splits the phases.' The Author may want to mark or cut them.
- [ ] The Dynamo v0.7.1 design doc is an old version (the page says the latest is v1.6.0). The quoted sentence is verbatim there. The conditional-disaggregation page is the dev version and is marked 'Experimental.' I confirmed its defaults (2048, 0.7), its example (4096, 0.6) and the KV-events requirement on the page itself.

<details><summary>What the fact-checker changed</summary>

- Fixed: Practice three: the pull connector 'issues an RDMA read'. The vLLM blog says it issues a NIXL READ, so the text now says 'issues a NIXL read'.
- Fixed: Placement paragraph: 'With prefill and decode kept on the same machines'. DistServe section 6.3 says the prefill and decoding instances keep the same pipeline stage on one machine. It now says 'With the prefill and decode segments of each pipeline stage kept on one machine'.
- Fixed: Placement paragraph: 'over 95% of requests waited under 30 ms'. The paper's figure measures transmission time, not waiting. It now says 'spent under 30 ms in transmission'.
- Fixed: Continuous batching: 'which is why serving systems adopted it' was an unmarked causal inference. DistServe only says that current serving systems use it, so it now says 'and current serving systems use it'.
- Fixed: KV copy paragraph: the InfiniBand 800 Gbps / NVLink 600 GB/s sentence had no link. It is now linked to the DistServe paper, section 3.3, which says both.

</details>

## Expert parallelism: why each GPU holds only some of the experts

`content/explainers/ai-ml-systems/expert-parallelism.md` · 2640 words · fact-checker: 96 claims checked, 4 fixed, 0 removed

**Drafter asks you to check:**
- [ ] The sources disagree on DeepSeek's decode layout. The V3 report (sec. 3.4.2) gives EP320 on 40 nodes with one expert per GPU. The February 2025 Day 6 overview gives EP144 on 18 nodes with 2 routed experts and 1 shared expert per GPU. Both are cited and labelled (paper vs production). Check that the wording is fair.
- [ ] I read GShard through ar5iv.labs.arxiv.org because the PDF would not extract here. The quotes cited to arxiv.org/pdf/2006.16668 are: 'infeasible to be replicated...', 'the only viable strategy is to shard experts into many devices', 'leaving other experts untrained, slowing down the training', 'an overflowed token', 'tied the number of devices...for simplicity' and 'The most significant part of gating is communication'. Also check '2,048 experts per MoE layer' (MoE(2048E, 36L)) and '2,048 TPU v3 cores in 4 days' against the PDF.
- [ ] Three links are our own inference or arithmetic. (1) The paragraph opening 'We infer' in Practice four says EP makes elasticity harder than replica scaling. (2) The paragraph opening 'We infer' in the closing section says elastic EP lets operators follow demand in steps smaller than a replica. (3) 'its dispatch reaches nine GPUs' follows from 9 experts per token and one expert per GPU. Check that the wording holds.
- [ ] The hook's line that four new GPUs 'disturb the sixteen already serving' is our framing of the Lead's point that changing the EP communicator invalidates retained workers' CUDA graphs and changes the expert mapping.
- [ ] Check the DeepSeek-V3 report anchors (#S2.SS1.SSS2, #S3.SS2.SSS1, #S3.SS2.SSS2, #S3.SS4.SSS1, #S3.SS4.SSS2, #S3.SS5.SSS1, #S4.SS2) against arxiv.org/html/2412.19437v2. Check that v2 is the version to cite.
- [ ] The vLLM docs URL ends in .html and redirects to /serving/expert_parallel_deployment/ (trailing slash). The piece cites the resolved address. The 2.4 GB redundant-expert figure and the 32-redundant-experts recommendation come from 'latest' docs, which can change.
- [ ] Node-limited routing (at most 4 nodes, NVLink 160 GB/s vs IB 50 GB/s, up to 13 experts) is stated for training (sec. 2.1.2 and 3.2.2). Prefill reuses the same all-to-all method. The text says 'In training'. Check that this stays clear.
- [ ] In the description, 'How it works, what it costs.' is a second clause after the sentence. It is 147 characters.

**Fact-checker could not confirm:**
- [ ] The arXiv HTML deep anchors for the DeepSeek-V3 report (#S3.SS4.SSS1, #S3.SS4.SSS2, #S4.SS2 and others) could not be checked against the page markup. The fetch tool showed only #S3.SS4, #S2.SS1.SSS2.Px*, #S3.SS2.SSS1, #S3.SS2.SSS2 and #S3.SS5.SSS1 / #S4.SS2.SSS0.Px1. The content at each section was confirmed, and v2 exists and serves the cited text.
- [ ] The vLLM docs anchors #example-command_1 (32 redundant experts) and #memory-footprint-overhead (2.4 GB) were confirmed by section content but not by literal anchor id. These figures come from the 'latest' docs, which can change.
- [ ] The DeepEP dispatch and combine quotes are verbatim, but they come from the README's example-code docstrings, not its prose. They are short and the repository is MIT-licensed.
- [ ] Hook: the October post gives the 100 s up and 56 s down pauses as the interruption that remained after the May (first-version) Elastic EP. It does not label them a 'first version' measurement in so many words. They match the post's 'Blocking' baseline (99.62 s / 56.25 s).
- [ ] The 'joining workers ... warm up while existing workers keep serving' line follows the post's design description. The timing table shows that warm-up and graph capture run in the background only in the NIXL EP configuration.
- [ ] The vLLM 'Instead of sharding each expert ... selected experts' definition quote is verbatim at about 147 characters: one sentence, kept as the definition. The Author may want to trim it.
- [ ] The two 'We infer' paragraphs (elasticity harder than replica scaling; following demand in steps smaller than a replica) are our own analysis and are marked as such. The hook's 'disturb the sixteen already serving' framing is supported by the Lead's account of CUDA-graph invalidation on existing workers.

<details><summary>What the fact-checker changed</summary>

- Fixed: GShard attention: 'Attention, much smaller, was replicated on every device' added a size claim the paper does not make. It now reads 'Attention, by contrast, had its weights replicated to every device', which is the paper's wording (attention is split along the batch dimension and its weights replicated to all devices).
- Fixed: GShard model: 'The largest model had 2,048 experts per MoE layer and 600 billion parameters' was wrong because the paper also ran a 1-trillion-weight MoE(2048E, 60L) scalability run. It now reads 'Its best-quality model had ...'. The PDF says 'The 600B parameters model that achieved the best translation quality was trained with 2048 TPU v3 cores for 4 days'.
- Fixed: DeepSeek-V3 decode: 'so its dispatch reaches nine GPUs' is our own arithmetic and was not marked. It now reads 'so, we infer, its dispatch reaches nine GPUs'.
- Fixed: DeepSeek-V3 prefill: 'Tokens cross nodes over InfiniBand to the GPU with the same in-node index' brought a detail from the training section (3.2.2) into the prefill description (3.4.1). It now reads 'Tokens cross nodes over InfiniBand first, then hop over NVLink to the GPU that owns their expert', which matches 3.4.1.

</details>

## Hash, range and directory sharding: which queries stay on one shard

`content/explainers/distributed-systems/hash-vs-range-vs-directory-sharding.md` · 2550 words · fact-checker: 72 claims checked, 8 fixed, 0 removed

**Drafter asks you to check:**
- [ ] None of the Sources uses the word 'directory'. We use it for Vitess lookup Vindexes, Vitess numeric_static_map and Neki's GSI. Check that the label is fair, and note that Neki's GSI points to a shard key, not to a shard.
- [ ] Placing Vitess's 'numeric' identity Vindex under range sharding, with contiguous placement and cost 0 in the table, is our inference. The Vitess docs have no 'range' Vindex.
- [ ] Vitess claims come from the v24.0 docs, the version the unversioned URL redirects to. Check the cost figures (identity 0, functional 1, lookup unique 10, non-unique 20) and the 'too slow for resharding' quote.
- [ ] Consistency section: the Lead says Neki writes the GSI lookup row 'in the same request' and that cross-shard transactions have no atomic commit. We infer that a partial failure is the case to test. The post doesn't say what happens.
- [ ] Resharding: the claim that splitting a hash key range moves only that range's rows, compared with modulo moving nearly every row, is our inference. Neki's Reshard copies the declared tables into a new shard group either way.
- [ ] The hook's 'four shard requests instead of eight' comes from the lifecycle post's recent-orders query after orders are re-sharded by customer_id. Check that our wording keeps that context.
- [ ] MongoDB facts come from the 9.0 manual: range is the default, hashed keys are ideal for monotonic fields, sharding on continent allows at most 7 chunks, and reshardCollection exists since 5.0. We don't claim reshardCollection runs online.
- [ ] To clear the author-placeholder warning, the frontmatter comment from make lead was reworded, because its literal 'TODO(author)' text set it off. No other skeleton lines (author, published, reviewed, draft: true, AI note) were changed, and verified: is not set.

**Fact-checker could not confirm:**
- [ ] The 'directory' label is the drafter's own term. None of the Sources uses it. It fits Vitess lookup Vindexes and Neki's GSI, though the GSI maps columns to a shard key, not to a shard. The TL;DR line 'stores the key-to-shard mapping in a table' simplifies this. The Author should decide whether that is acceptable.
- [ ] Placing Vitess 'numeric' (an identity Vindex, cost 0) under range sharding is marked as our inference. Vitess has no range Vindex.
- [ ] The Mermaid directory path follows customer_id 42 through a lookup row. That is an illustration, not a sourced example. It matches the Vitess lookup description (value, then keyspace ID, then shard).
- [ ] Several claims are inferences and are marked as such: splitting a hash range moves only that range's rows, a GSI partial failure is the case to test, directory balancing is the operator's job, and a merchant 100-200 query touches only eu-1.

<details><summary>What the fact-checker changed</summary>

- Fixed: Hook: 'a customer's orders were scattered across all four shards' now reads 'could be spread across all four shards'. The lifecycle post says they 'can be spread'. The 'four shard requests instead of eight' figure is verbatim, and its context (the recent-orders join after orders are re-sharded by customer_id) holds.
- Fixed: Diagram intro: 'whose own figure shows the same paths' now says the figure 'draws the hash and range paths beside a modulo one'. PlanetScale's figure shows xxhash, modulo and range. It has no directory path.
- Fixed: Hash expressions: the left(order_id, 4) example now says it only applies where the store ID is the order ID's first four characters, which is the condition the post states.
- Fixed: Range section: 'PlanetScale calls range and modulo less common...' had no link. It now links to the post's #beyond-xxhash section.
- Fixed: Directory section: numeric_static_map is no longer called a 'static variant' of the lookup Vindex. It is now 'a static map', and the text says Vitess classes it as a functional Vindex at cost 1, not a lookup, and that we file it under directory because its mapping is stored. This addresses the drafter's note about the label.
- Fixed: Which-queries section: the link text 'close values share a chunk' now reads 'close values are likely to share a chunk', which is what MongoDB's Ranged Sharding page says.
- Fixed: Lookup cost: 'the overhead is amortized ... because the alternative is a query sent to every shard' stated a cause Vitess does not give. It now follows the Vitess docs: the price may be worth paying to keep high-QPS reads from going to every shard, and the overhead is amortized as the number of shards grows.
- Fixed: When-to-use: 'MongoDB's own fix for both is a hash' was wrong for low-cardinality keys, since hashing does not add distinct values. It now says MongoDB's fix for monotonic keys is a hash and for low-cardinality keys a key with higher cardinality.

</details>

## vLLM vs DeepInfer: general-purpose or model-specific MoE serving

`content/explainers/ai-ml-systems/vllm-vs-deepinfer.md` · 2430 words · fact-checker: 96 claims checked, 9 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Decode 'within 6%' is the vLLM team's measurement on GB200, set against values they read off DeepInfer's README charts. DeepInfer's own documents make no comparison with vLLM; the draft says so, and that rests on reading the README and Q&A in full.
- [ ] Source conflict: the vLLM post says DeepInfer's prefill chart does not state its hardware. DeepInfer's README text under the charts says '4 SM100 GPUs for prefill and 64 for decode, each within one scale-up domain' (that text is in the repo's only README commit, 20 Oct). The draft reports both. The Trade-offs cells 'Published curve, 64 SM100 GPUs' and '4 SM100 GPUs' rely on the README.
- [ ] 'nextn=3 (a draft length of three tokens)' is our reading. It comes from DeepInfer's README, where DEEPINFER_USE_NEXTN_DECODE sets the maximum draft length, and from run_sweep.sh. Please confirm the gloss.
- [ ] The ~$38 cost comparison: vLLM's $38.96 per GPU-hour at DeepSeek's off-peak prices (a benchmark on one rack) against DeepSeek's US$38 per GPU-hour revenue at valley pricing (a production cluster). We say not to rank the engines on these two figures. The vLLM post does not say how it converts throughput to dollars.
- [ ] The paragraph on DEP32 against EP64+ is marked 'we infer'. It attributes the width difference to the hardware each team describes (an NVL72 rack, 64–128-GPU supernodes), not to the engines.
- [ ] The router in the vLLM subgraph of the diagram ('llm-d, Dynamo or your own') comes from vLLM's disaggregated-serving guide and docs. The Lead itself does not name the router used for the rack run.
- [ ] We call DeepSeek's moe_lb (slowest rank over mean) and vLLM's balancedness (mean over max) 'the same ratio upside down', and say we do not compare the values. Check the gloss of 'slowest rank' as max load.
- [ ] vLLM batch invariance: the tested-models list names DeepSeek-V3, V3-0324, R1 and V3.1 and no V4 model. The draft infers that DeepInfer is the only one of the two that documents bit-identical outputs for V4.1-Flash today.
- [ ] The Sources we read do not say whether vLLM's Huawei Ascend plugin runs V4.1-Flash. The draft states that as not described.
- [ ] Anchors on GitHub markdown (README #performance-benchmarks, #acknowledgments; Q&A #inference-system-overview, #expert-parallel-load-balancing, #data-parallel-load-balancing, #kernel-implementation) are GitHub's auto-generated slugs. Check that they resolve.
- [ ] The quotes on vLLM routing ('We are actively working on routing solutions.', 'fine for development but not for production') come from docs and blog pages that change over time. Re-check them before publishing.

**Fact-checker could not confirm:**
- [ ] The vLLM docs and blog quotes ('We are actively working on routing solutions.', 'fine for development but not for production', 'experimental and subject to change', 'may impact performance') match the pages as of 27 Oct 2026. They sit on the docs' /en/latest pages, which change over time, so re-check them before publishing.
- [ ] The diagram's vLLM router node ('llm-d, Dynamo or your own') is backed by the disaggregation guide and the vLLM docs. The Lead does not say which router its rack run used.
- [ ] The ~$38 cost figures set a benchmark against a production cluster's revenue. Both figures match their Sources, but the vLLM post does not say how it converts throughput into dollars.
- [ ] The README text 'Measured on 4 SM100 GPUs for prefill and 64 for decode' is in the README today. Whether that is the repo's only README commit (20 Oct, from the drafter's notes) was not checked.
- [ ] The DEPn footnote is linked at #footnote-label. That anchor exists, but the footnote itself is #user-content-fn-terms, a more specific anchor the Author may prefer.
- [ ] The GitHub heading anchors (#deepinfer, #prerequisites, #performance-benchmarks, #acknowledgments, #about; Q&A #inference-system-overview, #expert-parallel-load-balancing, #data-parallel-load-balancing, #kernel-implementation) were checked against the heading text by the GitHub slug rule, not loaded in a browser.

<details><summary>What the fact-checker changed</summary>

- Fixed: Opening, the nextn=3 gloss: it said 'a draft length of three tokens'. It now says 'in DeepInfer's settings, a maximum draft length of three tokens' and links the README. The README says DEEPINFER_USE_NEXTN_DECODE controls the maximum draft length, and run_sweep.sh maps nextn to that variable.
- Fixed: Mermaid node D6: it said 'Not released: Entrypoint, load balancers, KV cache store'. It now says 'Not released: load balancers, KV cache system'. The Q&A says only the KV cache management and load-balancing components are not yet open-sourced; no Source says the Entrypoint is unreleased, and the repo does ship deepinfer/entrypoints/llm.py.
- Fixed: Online serving: it said 'The vLLM team's disaggregation guide'. The post is by Martin Hickey (IBM Research), so it now says 'A disaggregation guide on the vLLM blog, by Martin Hickey of IBM Research'.
- Fixed: Online serving: it said 'DeepSeek says those components are tightly coupled', where 'those' took in disaggregation too. The Q&A names only KV cache management and load balancing, so it now says 'its KV cache management and load balancing components'.
- Fixed: Reproducible outputs and When to pick: 'bit-identical' is now 'identical' in both places. DeepInfer's README says 'identical', not bit-identical.
- Fixed: Sources, the 7 Oct post: the credit was 'vLLM Team'. It is now 'Inferact and the vLLM Team', the post's byline.
- Fixed: Sources, the NIXL push post: the credit was 'vLLM Team'. It is now 'Sunita Nadampalli (AWS) and Nicolò Lucchesi (Mistral)', with the date 14 October 2026 added.
- Fixed: Sources, the disaggregation guide: the credit was 'vLLM Team'. It is now 'Martin Hickey (IBM Research)', with the date 29 September 2026 added.
- Fixed: Added inline links where claims had none: the GB300 NVL72 rack, DEP2 prefill and DEP32 decode sentence (Lead, #prefill-dep2-works-best, #the-scaling-of-wide-ep); 'the setups differ' and the 'we do not read the gap...' quote (#compared-to-deepinfer); 'does not state its hardware' (#compared-to-deepinfer); vLLM's hardware list (vLLM README #about); 'the Lead' for the NVIDIA-only numbers; 'tops out at DEP32' (#the-scaling-of-wide-ep); 'EP64 or above' (Q&A #inference-system-overview).

</details>

## Kueue job admission: why a Kubernetes Job waits with no Pods

`content/explainers/infrastructure-cloud/kueue-job-admission.md` · 2650 words · fact-checker: 78 claims checked, 9 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Hook and walk: the claim that the Lead's Job waits because 3 CPUs and 3 Pods exceed the queue's 2 CPUs and 2 Pod slots comes from the Lead's Steps 3 and 5 and its 'What we learned' section. The quoted message 'insufficient quota for cpu in flavor default-flavor' is from Kueue's Troubleshooting Jobs page, which also shows 3 CPU and 600Mi, not from the Lead. The draft says 'the docs show a message of the form', so check that this framing is fair.
- [ ] Job submission: the 2022 Introducing Kueue post uses a kueue.x-k8s.io/queue-name annotation. The current Run A Kubernetes Job page says label, and the draft follows the current docs.
- [ ] Kueue's Cluster Queue page is inconsistent: its Cohort section says '.spec.cohort', while its examples and the Cohort page use 'cohortName'. The draft names neither field, but check it if you add YAML.
- [ ] Unsuspending a Job: 'a single update to the Job' that sets suspend to false and adds the flavor's labels to nodeSelector comes from the 2022 post's spot example. The current ResourceFlavor page says Kueue adds the flavor's labels and tolerations once the Workload is admitted. Check that pairing the two sources is accurate for current versions.
- [ ] Preemption cost: SIGTERM, the grace period and 'do not count towards completions' come from the Kubernetes Job docs on suspending a Job. The claim that preemption suspends the Job and deletes its Pods comes from Kueue's Troubleshooting Jobs ('Is my Job suspended?'). The draft chains the two sources; check the chain holds.
- [ ] Two 'we infer' sentences need your judgement. Practice one: quota above node capacity gives admitted-but-Pending Pods, and quota below it leaves nodes idle. 'Not gang scheduling': all-or-nothing applies to quota, not to nodes, drawn from the 4 Running / 1 Pending example in Troubleshooting Jobs.
- [ ] Numbers to check against their pages: 12 + 9 = 21 CPUs; borrowingLimit 9 + 1 = 10; lendingLimit 1 giving 10; Cohort fair-sharing weights 0.75 and 0.25 giving 75% and 25%; the 16-GPU ProvisioningRequest example (32 GPUs of quota, four nodes with four GPUs each); MultiKueue beta since v0.9; the visibility API since v0.6.0.
- [ ] The TL;DR bullet on how common admission is has no sourced number, because no first-party Source gives an adoption figure. 'Where it shows up' uses the integrations list, ProvisioningRequest and MultiKueue as variants instead of named production adopters, since the project's ADOPTERS.md returned 404. Add a first-party adopter account if you want one.
- [ ] The quote 'considering queue policy, quota, ResourceFlavors, admission checks, and other waiting Workloads' is from the Lead, anchored at #first-separate-admission-from-pod-scheduling. 'can work the control plane to death' and 'the job-controller instantly creates the pods...' are from the 2022 post at #why-kueue. Check that the wording is exact.

**Fact-checker could not confirm:**
- [ ] No HastCode internal links appear in the body and the front matter has no related: list, so there was nothing to check under content/.
- [ ] Kueue's Cohort field name (.spec.cohort vs cohortName) is inconsistent across its docs. The draft names neither field and has no YAML, so nothing was changed.
- [ ] 'which is the Lead's case' (line 81): the Troubleshooting Jobs message ('insufficient quota for cpu in flavor default-flavor', 3 CPU / 600Mi) matches the Lead's numbers and its 'insufficient CPU quota' description. The Lead does not print the exact message text, and 'a message of the form' is a fair framing.
- [ ] The two remaining 'we infer' passages (quota above or below node capacity; all-or-nothing on quota, not on nodes, from the 4 Running / 1 Pending example) are reasonable readings of the cited pages. They need the Author's judgement.

<details><summary>What the fact-checker changed</summary>

- Fixed: Introducing Kueue quote context: the draft said 'With many tenants sharing finite capacity' the default 'can work the control plane to death'. The 2022 post says 'At scale', and the draft now says 'At scale, the authors write'.
- Fixed: Alternative batch schedulers: the draft said they 'often replaced kube-scheduler or the job controller'. The 2022 post lists replacing them as one of 'one or more' problems, so it now says they 'could replace' them.
- Fixed: Admission as a 'single update to the Job' (suspend false, flavor labels added to nodeSelector): this comes from the 2022 post's spot example. It is now framed as 'In the 2022 post's example', and a sentence was added with a link to the current ResourceFlavor page, which says Kueue also adds the flavor's tolerations once the Workload is admitted.
- Fixed: 'It influences placement only by setting the flavor's labels...': the docs do not support 'only' (Topology-Aware Scheduling also affects placement). It now reads 'It influences placement by adding the flavor's labels in nodeSelector and adding the flavor's tolerations'.
- Fixed: Quota reservation: the draft said 'The docs sometimes call it workload scheduling, and warn...'. The glossary says it 'is sometimes referred to as workload scheduling', and the draft now says 'It is sometimes called workload scheduling, the docs note, but it is not the same thing as Pod scheduling'.
- Fixed: The line about a Pending Pod after admission (node capacity, affinity, taints, volumes, image pull, assigned to the scheduling branch) had no link. It now links to the Lead's #if-the-workload-is-admitted-switch-views.
- Fixed: 'two Pod slots' = reserved `pods` resource: the Lead never names the pods resource, so the sentence is now marked 'we infer'.
- Fixed: 'a rule that halves every investigation' was an overstatement. It now reads 'splits every investigation in two'.
- Fixed: Reproduction step 'Then raise the quota to three and watch the Job resume' is not in the Lead. It is replaced by the Lead's sourced line that Kueue keeps the Job suspended 'until quota becomes available or the queue configuration changes', linked to #what-we-learned-without-leaving-headlamp, followed by 'so we infer that raising both quotas to three would let it start'.

</details>

## ON DELETE CASCADE vs application deletes vs soft deletes

`content/explainers/databases-storage/cascade-vs-application-vs-soft-deletes.md` · 2585 words · fact-checker: 74 claims checked, 7 fixed, 0 removed

**Drafter asks you to check:**
- [ ] The hook: the Lead tells the outage as a second-person story ('Then one day, you get paged'), not as a dated incident. The draft calls it a 'scenario' Tom Pang walks through. Confirm that framing, and that 'hundreds of gigabytes of WAL written by a single transaction' is quoted exactly.
- [ ] The Lead's slug is i-dont-like-delete-cascade but its title is 'I don't trust ON DELETE CASCADE'. Section anchors (#cascade-plus-migration-equals-outage and others) came from the page's own heading links.
- [ ] MySQL facts were read through WebFetch summaries, because dev.mysql.com blocked direct download. These are: NO ACTION equals RESTRICT and is rejected immediately; cascaded actions do not activate triggers; an index on the referencing columns is created automatically; metadata locks extend to related tables (the #foreign-key-locking anchor). Check each quote and anchor on the live 8.4 page.
- [ ] 'We infer' claims to check hardest: a soft-deleted parent still accepts new child rows through its foreign key; a partial unique index WHERE deleted_at IS NULL lets a name be reused (the docs give the partial unique index feature, not this use); application deletes inherit the same logical-replication traps; readers can see a half-deleted tree between batches; a migration waits at most one batch.
- [ ] The Locks section says the cascade holds ROW EXCLUSIVE on every table involved. That rests on the Lead. The PG docs only say DELETE takes ROW EXCLUSIVE on its target table and that locks are held until the transaction ends.
- [ ] The 10,000-rows-per-batch example is from the June 2026 PlanetScale post. The 'hundreds to low thousands' figure is from the Lead. They are different scales and the draft keeps them apart. Also, the June post says DELETE does not touch index data, while the Lead says it does. The draft avoids both claims.
- [ ] Trade-offs cell 'Undo: None once a batch commits' for application deletes rests on an inference sentence in the guarantees section. No source states it directly.
- [ ] PostgreSQL docs are linked at /docs/current/ (18 at fetch time). The Sources list labels them 'PostgreSQL 18'.

**Fact-checker could not confirm:**
- [ ] MySQL 8.4 page (dev.mysql.com): curl got an Oracle 'Technical Difficulties' page, so these were checked only through two WebFetch summaries. The summaries confirmed: CASCADE 'the matching rows in the child table'; the depth-first search on the foreign key's index; 'Cascaded foreign key actions do not activate triggers'; NO ACTION equals RESTRICT and is 'immediately rejected'; the default is NO ACTION; an index on the referencing columns is created automatically (#foreign-key-restrictions); metadata locks are extended to related tables to stop conflicting DML and DDL (#foreign-key-locking). The anchor ids were reported by the summarizer. Wording and anchors should be checked once on the live page.
- [ ] Locks section: 'ROW EXCLUSIVE on every table involved' rests on the Lead ('it holds RowExclusiveLock on every table involved'). The PG docs only say DELETE takes it on the target table. The draft labels this correctly as the Lead's claim.
- [ ] Inferences are marked 'we infer' and no Source states them: a soft-deleted parent still accepts new children; a partial unique index on deleted_at IS NULL allows a name to be reused; application deletes keep both replication traps; readers see partial state between batches; a migration waits at most one batch. The Trade-offs cell 'Undo: None once a batch commits' rests on the same inference.
- [ ] 'with no moment when a reader sees half the tree' (guarantees section) follows from the atomicity the Lead states. No Source says it directly.

<details><summary>What the fact-checker changed</summary>

- Fixed: Hook said 'a page-at-night scenario', but the Lead never mentions night. It now reads 'In a scenario PlanetScale's Tom Pang walks through ("Then one day, you get paged"),'. The scenario framing is confirmed, and the phrase 'hundreds of gigabytes of WAL written by a single transaction' matches the Lead word for word.
- Fixed: RESTRICT 'cannot be deferred' linked only to the CREATE TABLE RESTRICT entry, which doesn't say it. The claim now reads 'does not allow the check to be deferred' and links ddl-constraints#DDL-CONSTRAINTS-FK, where 5.5.5 says 'RESTRICT does not allow the check to be deferred'.
- Fixed: MySQL 'There is no deferral' went further than the 8.4 page, which describes deferred checks only for NDB and says InnoDB rejects 'immediately'. It now reads 'The check is immediate, so children always go first.'
- Fixed: 'Most forms of ALTER TABLE need ACCESS EXCLUSIVE' did not match PG 13.3.1, which says 'Many forms of ALTER INDEX and ALTER TABLE also acquire a lock at this level.' It now reads 'Many forms of ALTER TABLE take ACCESS EXCLUSIVE'.
- Fixed: The partial unique index was introduced as 'one tool helps on the read side', but a unique index constrains writes. It now reads 'One tool helps with uniqueness'. The quote 'enforces uniqueness among the rows that satisfy the index predicate' is verbatim (it is the sentence just before #INDEXES-PARTIAL-EX3).
- Fixed: The link for indexing only live rows (Finding the children, soft deletes) pointed at Example 11.3, the unique-index example. It now points at #INDEXES-PARTIAL-EX2 ('Setting up a Partial Index to Exclude Uninteresting Values'), and the Sources line reads 'Examples 11.2 and 11.3'.
- Fixed: The logical-replication opener had no link. 'streams row changes from a publisher to subscribers' now links the Lead's #cascades-and-logical-replication, which gives that definition.

</details>

## Single, secondary or multi-primary DNS: surviving a provider outage

`content/explainers/apis-networking/dns-provider-redundancy.md` · 2735 words · fact-checker: 74 claims checked, 2 fixed, 0 removed

**Drafter asks you to check:**
- [ ] The Fastly quotes came through WebFetch, which would not give me the full page. Check each one word for word against the live post: "Every web request starts with a DNS lookup.", "addresses the failure that actually happens, depending on one authoritative provider", "This is not a failover.", "losing one costs you capacity rather than resolution", "Fastly cannot create or edit records" and "Not every provider does, and that is the one thing worth confirming before you plan this out." Also check the date (20 October 2026) and the byline (Matt Rinehart).
- [ ] The GitHub post quotes also came through WebFetch: "One of the best options available is to split authority..." and "can either be complex to set up or a completely manual process". Check them, along with the dyn/route53 example, the default dry run and the --doit flag. The 2019 update to the post may have changed wording.
- [ ] Single-provider section: RFC 2182 (1997) is about servers, not managed DNS providers. Applying its "(such as a service provider)" routing wording to a single DNS provider is marked "we infer". Confirm that framing is fair.
- [ ] RFC 2182 calls multiple primaries "not advisable", while RFC 8901 (2020) treats multi-primary as a supported model. Our reading that the difference is tooling is labelled "we infer". The "23 years later" figure is 1997 to 2020.
- [ ] Outage section, multi-primary: "no EXPIRE clock tied to the other" and "keep editing through the healthy provider" are our inference from RFC 8901's model. No Source states them.
- [ ] RFC 1996 section 3.6 calls a 60-second interval and at most 5 retransmissions "reasonable defaults". They are suggestions, not requirements; the text says "suggests".
- [ ] Cloudflare's "zone transfers are only available to customers on an Enterprise plan" is one vendor's plan rule as of today. Re-check it before publishing.
- [ ] The RFC 9199 TTL advice (at least one hour, ideally 4, 8 or 24 hours) and its caching claims are RFC 9199 reporting Moura et al., so the text says "passes on" or "cites measurements". Confirm that attribution reads right.
- [ ] Where the RFCs say "master" and "slave", the text uses the modern words primary and secondary. RFC 2182 defines them as synonyms ("Sometimes known as a Master server").

**Fact-checker could not confirm:**
- [ ] Fastly quotes, date and byline: WebFetch confirmed 'Every web request starts with a DNS lookup.', 'addresses the failure that actually happens, depending on one authoritative provider', 'This is not a failover.', 'losing one costs you capacity rather than resolution', 'Fastly cannot create or edit records' and 'Not every provider does, and that is the one thing worth confirming before you plan this out.', plus 20 October 2026 and Matt Rinehart. WebFetch goes through a summarising model, so the Author should still check them word for word on the live page.
- [ ] GitHub post quotes: 'One of the best options available...' is confirmed (the draft lowercases 'One'). 'complex to set up or a completely manual process' is confirmed as a fragment about keeping records identical across providers. dyn and route53, the default dry run, --doit and YAML in Git are all confirmed. These came through WebFetch too. The post says 'Updated March 7, 2019' but does not say what changed.
- [ ] Cloudflare 'Zone transfers are only available to customers on an Enterprise plan' is confirmed today. It is a vendor plan rule and could change before publishing.
- [ ] Inferred claims not stated by any Source, all correctly marked 'we infer': the single-provider failure domain, 'the difference is tooling', no EXPIRE clock under multi-primary, sync-tool risks and credentials, and the 'Which setup to choose' section.
- [ ] Preview build still shows one error for this page: broken-link to /apis-networking/cloudflare-hyperdrive-mysql-caching/. This file's body does not link there. The link comes from site navigation. The Hyperdrive page itself fails html-size (52071 bytes, limit 51200), so it is not emitted, and every apis-networking page and the author page report the same broken link. It goes away once the Hyperdrive draft is trimmed. That file is outside this task's scope.

<details><summary>What the fact-checker changed</summary>

- Fixed: Intro: 'Fastly shipped secondary DNS before primary DNS' implied Fastly has shipped primary DNS. The post only says 'Primary DNS is the bigger build, so we started with the piece that makes an existing setup more resilient on day one.' The line now reads 'Fastly started with secondary DNS rather than the bigger primary DNS build'.
- Fixed: Front-matter skeleton comment: the line 'replace every TODO(author) line' made the builder blame this file for TODO(author) warnings on listing pages, which really come from other unfinished skeletons (netflix-conductor, npm brief and others). It now reads 'no author notes are left', the same wording the finished Hyperdrive draft uses. No content changed.

</details>

# Production trims after verification (2026-10-28)

Six Pieces were over the 50 KB page gate in production builds (ad slots and analytics add about 1.4 KB over a preview build). Agents cut prose only after the Author verified them: no claim added or reworded, frontmatter untouched. Reread these cuts if you want; every remaining claim is one you verified.

## explainers/infrastructure-cloud/fly-io-sbd.md

Fits: True, 50156 bytes.

- Diagram 1 (SBD at a glance): shortened node and edge labels only, with the same nodes and edges. User 'User code and agents'; Ext4 'ext4, mounted like any disk'; NBD 'Kernel NBD block device'; Serve 'sbd serve daemon in the Sprite'; Cache 'Local NVMe cache'; edge to Bucket 'dirty blocks; flush on demand'; Bucket 'Bucket: blocks and their metadata'; Bucket->Serve edge 'reads checked against writes'; Catalog 'Checkpoint catalog: versioned snapshots'; restore edge 'restore: unmount, rewind, mount'
- Diagram 2 (three ledgers): shortened labels only. Ext4 'ext4 in a 100 GB image on a loop device'; Meta 'SQLite metadata, via Litestream'; Ckpt 'checkpoints.db: image per checkpoint'
- At a glance: removed 'It follows one write from user code down to object storage, then the checkpoint path beside it.'
- At a glance: removed 'First, why object storage needs a layer like this at all, and why the first one failed.' (the 'Top to bottom' order sentence stays)
- Why a Sprite's disk lives in a bucket: removed the clause ', so that getting migration right took Fly.io three years'
- Three ledgers: removed 'The cracks showed as usage grew.'
- Crash window: removed the whole paragraph 'The lesson: a checkpoint that shares storage with what it protects is fast, but guards against your mistakes, not against the storage layer losing a block.' (the 'Know what a snapshot shares' bullet makes the same point)
- Options, NVMe bullet: removed the clause ', which is the opposite of a Sprite that wakes wherever there is room'
- Options: removed 'Every fix to the old stack added another thing to reconcile.'
- Step 1: removed 'The change of layer is the whole idea.' and 'JuiceFS worked at the level of files and chunks, and Fly.io used it to hold a single file, a disk image, so ext4 still did the real filesystem work on top.' (this repeats the three-ledgers section)
- Step 2: removed 'That is why a pause moves no data: on wake, the same paths hold the same bytes.'
- Step 3: removed 'Compare the three records of the old stack.'
- Step 3: removed the trailing '; the post does not say how' from the inference sentence (the paragraph already says the post does not describe the layout)
- Step 4: removed 'After the upgrade, checkpoints mark a point in your disk's history instead of copying the disk' and its link to the storage-upgrades docs. That Source is still cited in the TL;DR and other sections
- What it costs (fail-stop): removed the clause ', and only services, the processes the runtime owns, are started again'
- What it costs (fail-stop): removed 'A daemon crash that a smarter layer might have retried becomes a visible interruption for whatever was running.'
- What it costs (migration): removed 'The upgrade API refuses to change storage without an explicit `allow_storage_migration` flag: "Routine upgrades never move you."'
- What it costs (migration), rollback: removed the clause 'and the docs warn against running the upgrade again after a rollback, which reuses the first attempt's disk copy'
- Evidence: removed 'On that basis SBD is now the Sprites filesystem from here on, and Sprites created before 10/6 can be migrated.'
- Evidence: removed 'The gap is large.'
- What it makes possible: removed the clause ', with nothing outside the Sprite coordinating it'
- Closing paragraph: removed 'The cost is restarts instead of quiet retries, and a migration that drops old checkpoints.'
- Not changed: frontmatter (title, description, published, verified: 2026-10-28, reviewed), every H2 (TL;DR and Sources included) and the Sources list. Each Source is still cited, so none was dropped. No links or claims were added. The page went from 52096 to 50156 bytes, and the production build prints no error or warning naming this page

## explainers/apis-networking/cloudflare-hyperdrive-mysql-caching.md

Fits: True, 50112 bytes.

- Problem section: removed the clause ', even for MySQL, which does not use an OS process per connection' from the per-connection resource cost sentence.
- Options section: removed the inference sentence 'The second half, we infer, is what later exposed the bug.'
- Step 3: removed the analogy ', like a SYN and a payload bundled together'.
- Step 3: removed the paragraph 'There are now two handshakes per logical connection, Worker to Client and Endpoint to origin, and Cloudflare argues that is the point: the same five-plus round trips still happen, but over the smallest possible distances.' (it repeats the deciding argument in the options section; its link is still cited elsewhere).
- Step 4: removed the parenthetical '(itself a fix: plain `utf8` is a deprecated alias for the 3-byte `utf8mb3`, which broke emoji)' and its link.
- Sources: dropped 'MySQL 9.7 Reference Manual: The utf8 Character Set (Deprecated alias for utf8mb3)', which nothing cites any more.
- Costs section: removed 'The charset fix also adds a command after every connection reset, which the post does not measure either.'
- Evidence section: cut ', rare enough that raising the trace rate did not catch an end-to-end trace; the bug was found by reading telemetry, the MySQL protocol documentation and Hyperdrive's code' (the sentence now ends at 'for one high-volume customer.').
- Evidence section: removed the inference sentence 'We infer it was mostly Workers running near their database, which is what Placement encourages for Workers that make multiple sequential queries per request.' and its link.
- What it makes possible: removed the closing sentence of the lesson paragraph, 'An assumption that held for one protocol, here "the bytes I sent are the bytes I can store", is worth writing down before you port it.'
- Flowchart: shortened the Worker label from 'Worker with a MySQL driver, such as mysql2' to 'Worker with a MySQL driver'.
- Flowchart: shortened the Client label from 'Hyperdrive Client: acts as the MySQL server, negotiates the Worker's capabilities, translates results' to 'Hyperdrive Client: poses as the MySQL server, translates results'.
- Flowchart: shortened the Endpoint label from 'Hyperdrive Endpoint: acts as a MySQL client with Hyperdrive's fixed capabilities, holds the pool' to 'Hyperdrive Endpoint: MySQL client with fixed capabilities, holds the pool'.
- Flowchart: shortened both cache labels (CacheA, CacheB) from 'CDN cache of this data centre' to 'Local CDN cache'.
- Flowchart: shortened the Step 3 edge label from 'Step 3: on a miss, over a pre-warmed TCP connection' to 'Step 3: on a miss, pre-warmed TCP'. No nodes or edges changed in either diagram.
- Not changed: the frontmatter (title, description, published, verified: 2026-10-28) and every heading the case-study guide requires. The production build now exits 0 and prints no html-size, broken-link or other error or warning for this page. The page is 50112 bytes without the inline stylesheet, against the 51200 limit.

## explainers/databases-storage/apache-iceberg-table-format.md

Fits: True, 50259 bytes.

- The page now fits: 51,840 bytes before, 50,259 after, against the 51,200 limit. That is 941 bytes of margin, close to the ~50,200 target. The production build no longer reports html-size for this page, and the full build log has no other error or warning naming it (0 warnings, broken-link lines included). Frontmatter, every H2 and every Sources entry are unchanged; every listed Source is still cited.
- 'A table format decides which files are the table': cut the paragraph 'A folder of Parquet files partitioned by directory name is not a table format: the file listing is the state, the Hive-style design whose [correctness problems Iceberg was designed to solve](reliability).' It repeats the Hive/reliability point, which stays in 'Why tables became metadata over files'.
- 'The tree, from catalog pointer to Parquet file': cut 'Each snapshot carries a summary. Snapshot 1 added 9,000 records in three data files; snapshot 2 took the total to 12,000 in four; snapshot 3 added 1,000 position deletes and one delete file, and left the total at 12,000.' It repeats the intro's counts.
- 'How a reader skips most of the tree': cut the closing paragraph 'The Iceberg docs report that in some cases this bounds filtering is "a 10x performance improvement"' (performance#data-filtering). The section's other link to that anchor stays.
- 'A commit is one pointer swap': cut 'Walk the diagram. Both writers load version 2 and write their files. Writer A swaps first.' The paragraph now starts at 'Writer B's commit names version 2...'.
- 'Three formats, one idea': cut the Hudi sentence 'Each action moves from requested to inflight to completed.'
- 'Why tables became metadata over files': cut the paragraph 'The second pressure is sharing. Many engines on one table only work if the rules for which files are live are written down, not held in one engine's memory.'
- 'Why tables became metadata over files': cut the paragraph 'The same features carry the costs. Immutability means every change adds files. History means old files stay until someone removes them.'
- 'What a table format is not': cut the closing sentence 'How catalogs implement the swap deserves its own Explainer.'
- Tree flowchart, labels only (nodes and edges unchanged): C 'Catalog: pointer to current metadata file' -> 'Catalog: current pointer'; M 'Metadata file 00003 (JSON): schema, partition spec, 3 snapshots' -> 'Metadata 00003 (JSON): 3 snapshots'; S 'Current snapshot, sequence 3' -> 'Current snapshot, seq 3'; L 'Manifest list (Avro): one per snapshot' -> 'Manifest list (Avro)'; D1 'Data manifest, seq 1: 3 files, 9,000 rows' -> 'Data manifest, seq 1: 3 files'; D2 'Data manifest, seq 2: 1 file, 3,000 rows' -> 'Data manifest, seq 2: 1 file'; X 'Delete manifest, seq 3: 1 file, 1,000 rows' -> 'Delete manifest, seq 3: 1 file'; P1 'Parquet data files for 1, 2 and 3 Oct' -> 'Parquet, 1-3 Oct'; P2 'Parquet data file for 4 Oct' -> 'Parquet, 4 Oct'; PD 'Position delete file, 2 Oct' -> 'Position deletes, 2 Oct'; dotted edge label 'skip these row positions in the 2 Oct file' -> 'skip rows'.
- Commit sequence diagram, message text only (participants and messages unchanged): 'Load table, current is v2' -> 'Load v2' (both writers); 'Write data files, manifests, manifest list, metadata' and 'Write its own files and metadata' -> 'Write files and metadata'; 'Swap pointer from v2 to v3' -> 'Swap v2 to v3' (both); 'Committed, v3 is current' -> 'Committed v3'; 'Rejected, v2 is no longer current' -> 'Rejected, v2 not current'; 'Check its assumptions still hold on v3' -> 'Recheck on v3'; 'Swap pointer from v3 to v4' -> 'Swap v3 to v4'.
- The file's frontmatter diff (verified: 2026-10-28 added, draft: true removed) is the Author's own uncommitted change. It is unchanged here.

## explainers/apis-networking/aspa-bgp-route-leaks.md

Fits: True, 50066 bytes.

- Mermaid labels shortened in 'The check' diagram. Nodes and edges are unchanged. R 'Route arrives on an eBGP session' -> 'Route arrives on eBGP'; C 'Pre-checks: is the neighbor the first AS on the path, with no AS_SET?' -> 'Neighbor is first AS, no AS_SET?'; P 'Compress the path: drop repeated AS numbers' -> 'Compress: drop repeated ASes'; Q 'What role does the neighbor play?' -> 'Neighbor's role?'; U 'Upstream check: the whole path must be one up-ramp' -> 'Upstream: all up-ramp'; D 'Downstream check: an up-ramp, then a down-ramp' -> 'Downstream: up-ramp, then down-ramp'; A 'For each hop, look up the lower AS's ASPA: is the next AS its provider?' -> 'Each hop: is the next AS in the lower AS's ASPA?'; edge labels 'Even the longest possible ramps cannot cover the path' -> 'Ramps cannot cover path', 'The ramps might cover it, but some AS has no ASPA' -> 'Might cover; an AS has no ASPA', 'Attested ramps cover the whole path' -> 'Attested ramps cover path'; X 'Ineligible for route selection, kept in Adj-RIB-In' -> 'Ineligible, kept in Adj-RIB-In'; Y 'Treated the same as Valid' -> 'Treated as Valid'.
- Removed from 'A route leak is a route that went the wrong way' (hairpin paragraph): 'The RFC explains why such leaks spread: the second provider prefers a customer's announcement over a peer's announcement of the same prefix.' The same point is kept in 'Why leaks get past origin validation'.
- Removed from the same paragraph (a secondary example): 'Its examples include a Telekom Malaysia leak of about 179,000 prefixes that Level3 accepted and propagated.' (RFC 7908 §3.1 link). RFC 7908 is still cited elsewhere.
- Removed the whole last paragraph of 'Practice three: add OTC to stop your own leaks': 'The draft also names what ASPA replaces. **Peerlock** has large carriers protect each other with AS path filters built by a laborious manual process (verification draft §8.3). ASPA publishes the same knowledge in the RPKI, for any network to use.'
- Removed from 'Where ASPA runs today' (The routers): 'Most contributors tested against a shared set of worked examples.' The Sriram examples link stays in 'What understanding ASPA makes possible', so its Sources entry is still cited.
- Removed from 'Where ASPA runs today' (The registries), because the TL;DR already says it: 'Across the global RPKI, the profile draft counts roughly 2,400 objects as of July 2026.' (profile draft §8.1 link). The TL;DR bullet with the same figure and link is kept.
- Removed from 'Why leaks get past origin validation': 'The rest of the path is not consulted.' The APNIC sentence right after it says the same thing.
- Removed from 'Why leaks get past origin validation': 'The origin is still the rightful one.' This repeated the sentence before it.
- Removed from 'What ASPA is not' (forged-origin paragraph): 'Hijacks deserve their own Explainer.'
- Not changed: frontmatter (verified: 2026-10-28, title, description, published), all required Concept headings, and the Sources list (every entry is still cited). Size measured: 51389 before, 50066 bytes after, without the inline stylesheet. The real-repo production build currently fails on an unrelated page (infrastructure-cloud/fly-io-sbd is over html-size, plus the broken links that follow from it), so no page files get written. The exact 50066 comes from a production build of a scratch copy of content/ in which only that fly-io-sbd copy was shortened so the build could finish. That copy is at /private/tmp/claude-502/-Users-nayanmevada-Desktop-Genuin-demos-go-articles/91601c52-ad00-47ee-96c4-01da44086d97/scratchpad/ccopy, and the repo file was not touched. In the real-repo production build, no error or warning names aspa-bgp-route-leaks. The Author verified this Piece before these cuts, so they may want to reread the diff. Only sentences were removed and labels shortened; no claim was added or reworded.

## explainers/infrastructure-cloud/github-agent-scale-git.md

Fits: True, 50048 bytes.

- Why a 25ms round trip caps the busiest repositories: cut the sentence 'A repository already near its ceiling would fall further behind with every push' and its link. It repeated the queue point in the sentence before it.
- Step 1: cut the sentence 'The first copy of a repository landed in Azure's Central US region in August 2025' and its link. This was a secondary migration detail.
- Step 1: cut the whole second paragraph, 'Routing reads to a nearby copy is not new: in 2017 Spokes already sent reads to the closest replica that was in sync', and its Stretching Spokes link. This was a secondary precedent. Stretching Spokes is still cited elsewhere.
- Step 2: cut 'Each write still pays the cross-region round trips, but writes to different references can be in flight at the same time.' It repeated the previous sentence and the costs section.
- Step 3: cut the trailing clause ', just as it does for the cross-region hop' after the per-reference serialization link.
- Step 3: cut the inference sentence 'GitHub's end-state plan moves maintenance to separate workers that run directly against durable storage, so we infer that "for now" marks a stopgap' and its link to #h-minimize-coordination. 'For now it also runs compaction' stays.
- Step 4: cut 'Object storage is its source of truth, and it reconciles its local state against it in the background.' It repeated the reconciliation point made in the reads-that-lag section.
- What it costs, every write still pays the round trip: cut the inference sentence 'We infer that, for that agent, the bridge is about not getting slower rather than getting faster.' It restated the sentence before it.
- What it costs, reads that can lag: cut the trailing clause ', fetching only the reference state and data that changed since the replica last caught up' after the 'move reconciliation in-line' link.
- The evidence: cut the clause ', having absorbed GitHub's growth for most of the migration while its own data centers held steady'. The sentence now ends at the 'about half of Git traffic' link.
- The evidence: cut 'Nor do they say how far the per-reference quorum and virtual replicas have rolled out to production.' The diagram section already says this.
- What it makes possible: cut the sentence 'It promises to explain why a small consensus group remains somewhere in the stack' and its link.
- No other changes. Frontmatter (title, description, published, verified, reviewed), every heading and the Mermaid diagram are byte-identical to before. I tried shortening one Mermaid label, but it made the SVG bigger, so I reverted it. All six Sources are still cited in the body, so none were removed.
- Result: 51937 bytes before, 50048 bytes after, both without the inline stylesheet. The final production build passed with 0 errors and 0 warnings, and nothing in the log names this file. The first build after the cuts also showed a broken-link error on this page. It came from a link to fly-io-sbd, which was being dropped for its own size at the time. It cleared once that page was trimmed. A copy of the original file is at /private/tmp/claude-502/-Users-nayanmevada-Desktop-Genuin-demos-go-articles/91601c52-ad00-47ee-96c4-01da44086d97/scratchpad/github-agent-scale-git.orig.md

## explainers/infrastructure-cloud/kueue-job-admission.md

Fits: True, 49916 bytes.

- How final_bytes was measured: it is an estimate, not a direct reading. A failed production build writes no pages, and this build fails on other pages' broken-link errors, so the production index.html never exists to measure. I calibrated a preview build against production using the original text in a scratch copy: preview gave 50621 bytes against production's reported 51988, an offset of +1367. Final preview size is 48549, so production is about 49916. The production build no longer reports html-size or any other error or warning for this file, which shows the page is under 51200.
- 'Five objects carry the decision': removed the closing paragraph 'Most of Kueue's operations centre on the ClusterQueue; the ResourceFlavor and LocalQueue are mainly organisational wrappers.' and its link to the 2022 post's #resource-model.
- 'Admission decides whether a Job starts...': removed 'It is sometimes called workload scheduling, the docs note, but it is not the same thing as Pod scheduling.' because it repeats the section's point.
- 'The admission path': removed 'The current docs add that Kueue also adds the flavor's tolerations once the Workload is admitted.' and its resource_flavor link. The tolerations point is still made in 'What Kueue admission is not'.
- 'The admission path': removed 'From there the job controller and kube-scheduler work as they always do.' because node K of the diagram already says it.
- Mermaid labels shortened, with the same nodes and edges: J 'Job sent to a LocalQueue, kept suspended' became 'Job in a LocalQueue, suspended'. W 'Workload: sum of all Pod requests' became 'Workload: sum of Pod requests'. F 'Fits ClusterQueue quota, or can borrow from the Cohort?' became 'Fits quota, or can borrow from Cohort?'. K 'Job controller creates Pods; kube-scheduler places them' became 'Pods created; kube-scheduler places them'.
- 'Where admission shows up': removed '; the 2022 post names the two things such an API must offer, a suspend field and mutable scheduling directives' and its #overview link. The 'Why plain Kubernetes needed a queue' section repeats it ('added two knobs to the Job API').
- 'Why plain Kubernetes needed a queue': removed 'With spare nodes, that is the fastest path to a running Job.'
- 'Why plain Kubernetes needed a queue': removed 'The Kueue authors' summary: quotas are enforced on resource creation, and there is no queueing of requests.' and its #why-kueue link, which is still cited elsewhere. It repeated the ResourceQuota 403 sentence before it and the TL;DR.
- 'Practice two': removed the secondary example ': in the docs' example, weights of 0.75 and 0.25 steer two organisations towards 75% and 25% of common resources', so the sentence now ends at the fair-sharing weights link.
- 'Practice three': removed ', and its docs include a proof that two Workloads cannot keep preempting each other under Fair Sharing' and its fair_sharing link.
- 'What understanding admission makes possible': removed 'The alpha Kueue plugin for Headlamp shows the same objects as linked, read-only pages.' and its plugin link. The plugin's sample manifests are still cited, and the Headlamp Sources entry is unchanged.
- Sources: the Preemption entry no longer lists 'Fair Sharing'. Its only citation was cut above. Every other Sources entry is unchanged.
- Not touched: frontmatter (title, description, published, verified, reviewed are exactly as the Author left them), all H2 headings, TL;DR, the Mermaid nodes and edges, and every other claim and link.

# Daily review queue, 2026-11-07

Drafted by `make daily` (ADR-0009). Every Piece is `draft: true` and unverified. Verify every claim against its Source, then set `verified:` and remove `draft: true` yourself (ADR-0007).

## Run 20261107-060325

### Fly.io's fleet reboot: a VM is only as strong as the host's parser

`content/briefs/fly-io-ipv4-checksum-offload-cve.md` · brief · 244 words · fact-checker: 26 claims checked, 4 fixed, 0 removed

**Publish within 72 hours of the Lead's date**, with or after its Explainer: netflix-single-tenant-titus, openshell

**Drafter asks you to check:**
- [ ] Post date: page shows no date or byline (only illustration credit Annie Ruygt). Fly blog feed.xml gives published 2026-11-06T00:00:00+00:00, updated 2026-11-06T09:46:59+00:00. Brief says 'On 6 November'. Lead date 2026-11-06, pick 2026-11-07T06:09:09Z: inside 3-day window.
- [ ] 'every host patched and rebooted in 10 hours': post says 'By that evening, we had patched and rebooted our entire fleet' and 'The full rollout took 10 hours'.
- [ ] CVE-2026-99999 'published a day earlier': cve.org record datePublished 2026-10-15T12:19:22Z, CNA Linux; post says 'announced CVE-2026-99999 on the morning of the 15th'. Event day is 16 October.
- [ ] 'a wrong IPv4 check on checksum offload': post: 'The offload code has checks meant to reject bad metadata. For IPv4, they were wrong.' CVE title: 'ipv4: validate checksum_start before completing checksum'.
- [ ] Metadata sentence: post: 'the kernel stores metadata saying where the checksum should read from and write to' and 'The virtual network card accepts checksum metadata and passes it to the host.'
- [ ] 'rewrites two bytes' of a header 'already routed and filtered': post: 'The attacker now controls two bytes of an IPv4 header the host has already trusted' and 'lands after the host has already applied routing and firewall rules'.
- [ ] 'a firewall bypass, or 32 bytes of host memory, a neighbouring VM's packets included': post: 'The first is a firewall bypass', 'captures 32 bytes of memory as options on the outgoing packet', 'host packets, kernel data structures, and packets belonging to the VM next door'.
- [ ] '15 to 20 minutes' to reboot: post says 'around 15 to 20 minutes'. Waves: 'no MPG cluster had both nodes in the same wave'; 'our largest region, iad, needed 17 waves'. Brief paraphrases as 'so no Managed Postgres cluster lost both nodes'.
- [ ] Firecracker claim links docs.fly.io/reference/architecture: 'Application code runs in Firecracker microVMs.' The postmortem itself does not name Firecracker.
- [ ] Only quote: 'Assume every boundary leaks a little' is a closing recommendation heading in the post, verbatim. Credited to Fly.io, linked.
- [ ] 'a VM still reaches the host through device paths': post: 'A guest should not be able to influence the host outside of a few narrowly defined paths.' Post links the virtio 1.3 spec for the device path.
- [ ] Analysis, not in Source: 'hostile tenants add the host's device emulation to the attack surface'; 'a reboot patch is a planned outage: stage waves so no pair's two halves fall together'. Post does not use the term HA pairs; the generalisation is ours from the MPG waves.
- [ ] Netflix link: Explainer says Netflix 'needed resource and performance isolation, not protection from an attacker' and that microVMs/gVisor 'are built to contain hostile tenants'. Analogy is ours; Source does not mention Netflix.
- [ ] OpenShell link: Explainer says credentials are swapped for placeholders and 'if the supervisor disconnects, the sandbox freezes the agent' (supervisor is 'the trusted side'). Ties to Fly's advice 'If you can keep credentials away from agents, even better'. Source does not mention OpenShell.
- [ ] lore.kernel.org announce (linked from the post) returned 403 to us; cve.org record used instead. Confirm the lore link if you prefer it in sources.
- [ ] Description 147 chars and title 66 chars by hand count; ad-hoc Python was denied. Body ~244 words by hand count; builder passed brief-length (earlier draft reported 285, builder counts about 2 fewer than hand count).
- [ ] sources: frontmatter added (CVE record, Fly architecture docs); build accepted it without a frontmatter warning.

**Fact-checker could not confirm:**
- [ ] Post date 2026-11-06: the page shows no date or byline; only Fly's feed.xml gives published 2026-11-06T00:00:00+00:00
- [ ] cve.org record page is JavaScript-rendered and returned no text; CVE-2026-99999 (datePublished 2026-10-15T12:19:22Z, CNA Linux, title 'ipv4: validate checksum_start before completing checksum') confirmed via the CVE Services API JSON instead
- [ ] 'a 32-byte leak finds placeholders' is HastCode analysis of OpenShell's design, not stated by any Source

<details><summary>What the fact-checker changed</summary>

- Fixed: Checksum metadata description: 'where to compute a checksum' -> 'where a checksum reads and writes' (post: 'where the checksum should read from and write to')
- Fixed: Reboot time: '15 to 20 minutes' -> 'around 15 to 20 minutes' (post says 'around')
- Fixed: Managed Postgres waves: 'lost both nodes' -> 'rebooted both nodes at once' (post: 'no MPG cluster had both nodes in the same wave')
- Fixed: Linked previously unlinked claim 'trusted guest metadata after routing and filtering' to the Fly.io post; dropped leading 'And' to stay under 250 words

</details>

### GitHub's Project Rubicon: only the pointers need consensus

`content/briefs/github-ref-cas-git-storage.md` · brief · 246 words · fact-checker: 25 claims checked, 2 fixed, 0 removed

**Publish within 72 hours of the Lead's date**, with or after its Explainer: github-agent-scale-git

**Drafter asks you to check:**
- [ ] Date: post is dated November 5, 2026 with no time; Brief says 'On 5 November' (no UTC time available).
- [ ] Availability: Source final paragraph says 'available in a limited private preview now, with a broader rollout early next year'; Brief says 'now in limited private preview'.
- [ ] Push path: Source ('Introducing Project Rubicon') says a virtual replica 'indexes the incoming pack and uploads it to Blob Storage' then 'commits one transaction to Ref-CAS that lists the pack in the manifest and updates the references'. Brief drops 'Azure' before 'Blob Storage' for length.
- [ ] Ref-CAS contents: 'References and the list of packs that make up the repository (the pack manifest) live in a separate replicated log service' and 'Ref-CAS never holds or validates object data'.
- [ ] CAS and majority: 'Each reference update in Ref-CAS is a compare-and-swap committed by a quorum' and 'becomes visible once a majority of the log's replicas have recorded it'.
- [ ] LSN: 'Each commit to Ref-CAS gets a log sequence number (LSN), and the repository's version is the LSN of its latest commit'.
- [ ] Compaction/GC: 'a fleet of compaction nodes works in the background, to merge each repository's packs and remove unreachable objects. That work never runs on a virtual replica.' GC 'using dedicated workers'. Brief says 'run on a background fleet, never on a serving replica' covering both; confirm GC is covered by 'that work'.
- [ ] Numbers: lede says 'support 5,000 concurrent agents working on a repository, sending 7,200 pushes per second to a repo' from internal testing. Brief: 'In internal tests one repository took 7,200 pushes a second from 5,000 concurrent agents'.
- [ ] Quote: 'Only reference updates need agreement.' sits in the introduction before the first H2; credited to the byline 'Karthik Raman & Brian Celenza'. Only quote in the Brief.
- [ ] Pro Git link (https://git-scm.com/book/en/v2/Git-Internals-Git-Objects) supports 'Git names objects by content'; also listed in sources: frontmatter (build accepted the format).
- [ ] Inference (labelled 'we infer'): a slow Blob upload delays only that push, not its neighbours. Based on 'A virtual replica uploads a push's pack to Blob Storage before it commits the reference transaction' and 'A transaction locks only the references it touches'. Source does not state it and gives no latency figures; WebFetch found no sentence on Blob latency or slow Blob.
- [ ] 'where each vote still covered a full copy' characterises the bridge from our Explainer github-agent-scale-git (its sourced claim: replicas 'each hold a full copy of the repository'). Check it is a fair one-clause summary.
- [ ] 'Rubicon half answers it': the Explainer's open question was 'what happens to writes when Blob Storage is slow or that node fails'; the Source answers only that the upload sits on the write path before the commit.
- [ ] 'the latest log sequence number is the repository's version' and 'a vote only settles which commit a reference names' are our paraphrases; confirm against 'the repository's version is the LSN of its latest commit' and Source's 'A reference only ever points at objects that are already durable'.
- [ ] Word count: build counted 255 before trim; passes 250 now (approx 246 by the build's counter). Title 58 chars, description 144 chars.
- [ ] Published is the pick time 2026-11-07T06:09:09Z; Lead dated 2026-11-05, so the 3-day window closes 2026-11-08.

**Fact-checker could not confirm:**
- [ ] 'it gives no latency figures': confirmed for all article text and every image alt text; two figures (writes/read side-by-side visualization, pack-count chart) had no alt text in the fetched content, so their pixels could not be inspected. Nothing suggests they carry latency values.

<details><summary>What the fact-checker changed</summary>

- Fixed: What happened: 'uploads it to Blob Storage' -> 'uploads it to Azure Blob Storage' (Source's product name)
- Fixed: What happened: 'Compaction and garbage collection run on a background fleet, never on a serving replica' -> 'run on dedicated background nodes, never on a virtual replica' (Source: compaction fleet 'works in the background' and removes unreachable objects, 'That work never runs on a virtual replica'; GC uses 'dedicated workers'; 'serving replica' is not Source's term)

</details>

### Not drafted

- none

<details><summary>Leads the editor rejected</summary>

- `5jswvc4b`: Fetched: Google Cloud europe-west2 power loss (PLC undervoltage sensors wired upstream of an isolated breaker, generators never started, regional spillover saturated zones a and c) has a full first-party RCA, but the status page is updated in place and its first post was 2026-10-20, outside the 72-hour Brief window per the Brief guide; strong Case Study candidate when the Explainer quota allows.
- `s6vjv67p`: Fetched: GitHub's Actions/Pages/Webhooks incident states its cause (a config change to an event-routing service depended on feature flags not set as intended, the rollback left the code in place so a redeploy reintroduced it), but no live Explainer teaches flag defaults or config-versus-code rollback coupling; window also closes 16:05 UTC today counting from the first post.
- `4vfvzwgz`: Fetched: vLLM's migration to the LibTorch Stable ABI (v0.33.0 wheels, seven dependencies migrated, boxed StableIValue crossing a C shim) is a real design change with mechanism, but no live Explainer teaches ABI or FFI boundaries; Concept candidate later, Explainer quota is 0 today.
- `hlbeocrj`: Fetched: GitHub multi-service degradation of 2026-11-05/06 gives no cause; 'a detailed root cause analysis will be shared'.
- `4uzvwb5x`: Fetched: GitHub Actions incident t89g2q1800fr gives no cause, mitigation unnamed.
- `qwggc3mr`: Fetched: GitHub Actions incident h02n9572mx9y only points to the combined summary xr3myzj7ms51 (Lead s6vjv67p), no cause of its own.
- `3ittfw7g`: Fetched: GitHub Pull Requests/Git/Issues incident gives no cause or mechanism.
- `eanovtqi`: Fetched: Copilot Mission Control errors cite 'upstream database throttling' then 'a specific client version' with no mechanism connecting them, and no live Explainer fits.
- `5d62xswp`: Fetched: AWS Health Dashboard link is the dashboard root, renders by JavaScript only; no incident text, region or cause retrievable for 'Elevated Launch Failures'.
- `w5bonxa4`: Same AWS Health Dashboard root link as 5d62xswp; no cause retrievable.
- `kuvz2ptn`: Same AWS Health Dashboard root link as 5d62xswp; no cause retrievable.
- `ywbe3ike`: Same AWS Health Dashboard root link as 5d62xswp; no cause retrievable.
- `bityigxe`: Fetched: Cloudflare's post-quantum Privacy Pass work (MOLE, VOLE-in-the-Head blind signatures) is a research proposal with no deployment date, not a shipped change, and no live Explainer covers post-quantum credentials.
- `l4y5jthh`: Fetched: Trusted Agent Delegation is a coalition and protocol proposal; the shipped parts are product features (analytics category, access policies, early-access abuse protection) with no architectural mechanism, and no live Explainer teaches the lesson.
- `5hblcz6e`: Fetched: zero-mismatch async RL is a training recipe with results, partly in draft PRs; no live Explainer teaches RL training-inference alignment.
- `5crb3xfq`: OrderServ multi-tenant ECS is an AWS customer cost story, not an event; Explainer quota is 0 today.
- `r5ttoips`: PCAP-LM is a tool description, not an event at a named system.
- `zqewpxbh`: RFC 10073 is a standards publication, not an event at a named system; the only trust-store Explainer (bundled-vs-platform-trust-stores) is not live.
- `eybwjwjx`: DynamoDB bulk put items is a feature announcement; no live Explainer on DynamoDB or bulk loading.
- `tjfjy4gk`: TUI threat modeling with Bedrock is a customer-success story without system mechanism.
- `iwvczomr`: Google Research hypertension detection is outside the six Verticals and not an event.
- `o7dg4a4z`: Shopify Reflex is an architecture post, not an event; Case Study candidate when the Explainer quota allows.
- `eep3iy4k`: Martin Fowler's Data Product Platform Navigator is a framework essay, not an event and not a single named mechanism.
- `2d3e7tr3`: APNIC 'The Internet didn't break' is opinion.
- `4p5ahled`: DuckDB in DataGrip is an engineering post, not an event; possible Case Study later, Explainer quota is 0 today.
- `aa4wodac`: ReinDB is a product announcement.
- `aq7pxhuh`: PlanetScale's unplanned failover design post is not an event; Comparison candidate later, Explainer quota is 0 today.
- `cpeb3uo6`: Podcast with no readable text.
- `leqteksw`: Autobase 2.12 is a point release.
- `spzwfw4i`: myDBR 7.7 is a point release.
- `grtrs2bn`: Call for proposals.
- `mvpdxzbi`: Google Research 'Agile Deliberation' is a research post, not an event at a named system.
- `oomzdelu`: Cloudflare personal-agent visibility is a product announcement without architectural mechanism.
- `5mia7quk`: Fly.io Managed Postgres etcd space-quota incident would have matched the etcd Explainer, but it is 78 hours old (brief_ok false).

</details>


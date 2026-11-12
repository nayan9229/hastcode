# Daily review queue, 2026-11-12

Drafted by `make daily` (ADR-0009). Every Piece is `draft: true` and unverified. Verify every claim against its Source, then set `verified:` and remove `draft: true` yourself (ADR-0007).

## Run 20261112-091131

### vLLM's encoder disaggregation: a third GPU pool must earn its keep

`content/briefs/vllm-encoder-disaggregation-epd.md` · brief · 248 words · fact-checker: 21 claims checked, 3 fixed, 0 removed

**Publish within 72 hours of the Lead's date**, with or after its Explainer: disaggregated-inference

**Drafter asks you to check:**
- [ ] Date and byline: post dated November 11, 2026, byline 'Inferact, InferMatrix, and IBM' (post header).
- [ ] 'now in vLLM v0.32.0': post says 'The core vLLM capabilities described below, including MiMo and XPU support, are included in vLLM v0.32.0.' FLAG: GitHub release page for v0.32.0 (shows 18 Oct) did not visibly mention EPD/EC terms in the text fetched (page truncated at 100k chars); the GB200 benchmark ran 'main at 5314c58281 plus the DeepSeek metadata-only extension', not the tagged release. Decide whether 'now in' is fair.
- [ ] 'run no language model': post says 'An encoder-only instance runs no language model, maintains no language-model KV cache, and produces no sampled tokens.'
- [ ] EC connector backends: post names 'NIXL + CPU EC' (NIXL transfers between registered CPU mmap regions) and 'Mooncake EC' (Mooncake TransferEngine, TCP and RDMA). Post never expands 'EC'; docs page expands as 'encoder-cache (EC)'. Brief says embeddings 'reach prefill'; in E+PD mode they reach the combined PD instance.
- [ ] Metadata-only: post says a metadata-only request 'replaces the downstream media item with a reference containing the structural information' and 'The consumer reconstructs placeholders from that metadata and receives the actual embeddings through the EC connector.' Embeddings still travel over EC even in metadata-only mode; Brief's 'carries just the layout' refers to the downstream request, not the embeddings.
- [ ] GB200 setup: 'PD baseline used four GPUs for P and four for D'; 'EPD added four single-GPU encoder instances while keeping the same P/D configuration: TP=1, DP=4'; 12 GPUs total for EPD; DeepSeek V4.1 Flash; 512 MP-DocVQA validation questions, 4-6 document pages each; requests averaged 4,905 input and 5.2 output tokens; three runs per concurrency; EPD used Mooncake GPU RDMA for embeddings, both used NIXL for KV.
- [ ] Throughput definition: 'input plus output tokens, divided by elapsed time and all deployed GPUs: eight for PD, twelve for EPD, including all four encoder GPUs.' Brief: 'throughput per GPU', 'all 12 counted'.
- [ ] Numbers: 'At concurrency 64-128, EPD delivered 6.6-9.5% higher input-plus-output throughput per GPU than the eight-GPU PD baseline' and 'P90 TTFT was 44-47% lower.' 'At concurrency 16, however, EPD had lower per-GPU throughput; at 32 it was approximately break-even.' Brief renders break-even as 'roughly tied'. Post gives no absolute numbers in text (plot CSV linked).
- [ ] 'independently of decode': post says 'The demand on this stage varies with media size, modality, and request mix, and can change independently of prefill and decode demand.'
- [ ] 'small beside the language model' links vLLM docs 'Disaggregated Encoder' (https://docs.vllm.ai/en/latest/features/disagg_encoder/): 'Vision encoders are lightweight, while language models are orders of magnitude larger.' Docs speak of vision encoders; Brief says 'Encoders'.
- [ ] The one quote, verbatim from the DeepSeek V4.1 Flash section: 'The benefit depends on load and resource allocation, rather than separation alone.' Credited as 'The authors'.
- [ ] Labelled inference: idle encoder-only GPUs explain the per-GPU loss at concurrency 16. Post does not state a cause for the loss.
- [ ] Our framing, not in post: 'EPD adds a third pool and a second copy' (post: embeddings via EC connector, KV via NIXL, so two transfers); 'one GPU cannot be sized for both'.
- [ ] Description 144 chars, title 66 chars, body 248 words by manual count (build brief-length gate passed). Lead age: post 11 Nov, picked 12 Nov 09:17 UTC, inside 3-day window.
- [ ] sources frontmatter added with the docs page; optional, remove if unwanted.

**Fact-checker could not confirm:**
- [ ] "now in vLLM v0.32.0": supported only by the post's sentence "The core vLLM capabilities described below, including MiMo and XPU support, are included in vLLM v0.32.0"; GB200 run itself used "main at 5314c58281 plus the DeepSeek metadata-only extension"; release notes not checked
- [ ] "Embeddings reach prefill": post says embeddings reach "consumers"; true for E/P/D (the GB200 setup), in E+PD the consumer is the combined PD instance
- [ ] "independently of decode": post says "independently of prefill and decode demand" (brief states a subset)
- [ ] "so one GPU cannot be sized for both": HastCode analysis, not stated by Source, not labelled "we infer"

<details><summary>What the fact-checker changed</summary>

- Fixed: "Encoders are small beside the language model" -> "Vision encoders are small beside the language model" (docs Source: "Vision encoders are lightweight, while language models are orders of magnitude larger")
- Fixed: Removed "just" from "carries just the layout the encoder computed": embeddings still travel over the EC connector in metadata-only mode, so "just" overstated
- Fixed: Linked "a second copy" to the post ("EC connectors deliver encoder embeddings; KV connectors transfer the state between prefill and decode"); was an unlinked claim

</details>

### Shopify's 12-week native rebuild: the harness set the pace

`content/briefs/shopify-mobile-app-native-rebuild.md` · brief · 234 words · fact-checker: 24 claims checked, 3 fixed, 0 removed

**Publish within 72 hours of the Lead's date**, with or after its Explainer: agentic-development

**Drafter asks you to check:**
- [ ] Date and byline: post shows 'Published on Nov 10, 2026', author Remi Vanhee; no time of day given, Brief says 'On 10 November' with no time zone. Publish by 13 Nov 2026 to stay inside 72 hours.
- [ ] Intro facts: 'Thirty engineers just rebuilt the Shopify app from React Native to Swift and Kotlin in 12 weeks'; 'more than 300 screens'; 'Agents wrote most of the code'; later 'About 30 engineers worked on the rewrite, organized in pods of two'. Brief says 'about 30 engineers', '300-plus-screen', 'agents writing most of the code'.
- [ ] 'Before any screen moved, a small group built one architecture for both platforms': section 'How we prepared the app for agents' (small group built the surrounding system before migrating screens; one architecture shared by iOS and Android).
- [ ] ViewModel parity: main post, section 'Consistency over platform-specific architecture': 'An iOS ViewModel for the orders screen looks like the Android ViewModel, whether you're reading the Swift or Kotlin version.' Linked to the Reflex post (Thiago Magalhaes, 5 Nov 2026, shopify.engineering/reflex), which says Swift and Kotlin ViewModels share the same structure. Confirm Reflex post states this.
- [ ] CI divergence check: exact sentence 'CI fails when iOS and Android diverge in ways the architecture doesn't allow.' (section 'Consistency over platform-specific architecture'). Brief paraphrases as 'CI that fails when iOS and Android diverge in ways the architecture forbids'.
- [ ] CLI speed: 'a CLI that can load any screen in milliseconds' (section 'How we prepared the app for agents'; previous loop 'several minutes per iteration'). Reflex post says test cases run 'in a couple hundred milliseconds without a device'. Brief says 'loads any screen in milliseconds'.
- [ ] Agent review: 'Every PR gets an adversarial review against our architecture guidelines before a human looks at it.' (section 'Code reviews became a bottleneck'). Brief: 'an adversarial agent review before a human read it'.
- [ ] Small checkpoints: post says a PR is a small checkpoint not a full screen, the diff size a human can read is a constraint given to the agent, and the team prefers small checkpoints over large agent-generated PRs (sections 'Keeping agents on the rails' and 'The existing app was the spec'). Brief: 'pull requests stayed small checkpoints a human can read' and 'not from bigger agent tasks'.
- [ ] The one quote, exact wording including semicolon: 'A document explains why the rule exists; CI decides whether the rule was followed.' (section 'Keeping agents on the rails'), credited as 'Vanhee's rule'.
- [ ] Bottlenecks: section 'What are the new bottlenecks?' with H3s 'Laptops became a bottleneck', 'Code reviews became a bottleneck', 'Merging became a bottleneck', 'Coordination became a bottleneck'. Brief: 'Once code was cheap, laptops, review and merging became the queue'.
- [ ] Inference flagged: the guides/sensors mapping (ViewModel shape as guide; divergence checks, CLI, agent review as sensors) is ours; the post does not use those words. Confirm the post never says 'guide', 'sensor' or 'harness'.
- [ ] Editor's angle named 'lint rules'; the post lists lint rules among guardrails but the Brief cites only 'CI that fails' for divergence. Reflex post describes shared tests ('a PR can't merge if it only passes on one') rather than lint rules.
- [ ] Frontmatter: 'sources:' entry for the Reflex post added (build raised no frontmatter warning). 'related:' left commented out with meta-thin-orchestrator as a suggestion. Unverified numbers from the post (launch time, binary size, memory, crash-free rate) deliberately left out for length.
- [ ] Word count 234 is a hand count (note 26, headings 5, What happened 91, Why it matters 112); the build's brief-length gate passed. Why it matters is longer than What happened as the guide requires.

<details><summary>What the fact-checker changed</summary>

- Fixed: ViewModel claim: 'reads the same in Swift and Kotlin' -> 'reads almost the same' (Reflex post: 'the Kotlin side has the same structure and reads almost the same')
- Fixed: 'built before the agents started' -> 'built before the screens moved' (post says prep preceded screen migration; does not say agents were absent from prep)
- Fixed: 'Once code was cheap, laptops, review and merging became the queue' -> 'Once writing code was no longer the slowest part, laptops, review, merging and coordination became the bottlenecks' (post: 'Once writing code was no longer the slowest part of the system, the bottleneck moved'; four bottleneck H3s; 'cheap' absent)

</details>

### Not drafted

- none

<details><summary>Leads the editor rejected</summary>

- `st3piyl3`: Node.js 26.13.1 security release (fetched): real event with HTTP/2 stream-memory and TLS SNI mechanisms, and http-connection-reuse teaches where HTTP/2 multiplexing breaks; viable Brief but lost to the cap of 2 today. Fallback if a pick fails.
- `oemexus6`: Node.js 24.22.3: same security release as st3piyl3, duplicate event.
- `zuyldv2d`: Node.js 22.23.5: same security release as st3piyl3, duplicate event.
- `ehaht7qe`: Slack incident (fetched): cause given only as a likely trigger (a database change leading to storage capacity exhaustion), not confirmed, and no live Explainer teaches database storage exhaustion.
- `7fycx6wy`: GitHub multi-product incident (fetched): status page gives no cause, only 'detailed root cause analysis will be shared'.
- `hpnjx3fo`: GitHub Webhooks incident (fetched): four-minute incident, no cause stated.
- `xcjrs63o`: GitHub Git Operations and Pull Requests incident (fetched): no cause stated.
- `j6yor3hg`: GitHub Copilot model degradation (fetched): attributed only to 'an upstream model provider', no mechanism.
- `jqqos7pq`: GitHub Actions and Checks API incident (fetched): no cause stated.
- `qdtwws4n`: GitHub Webhooks and Actions incident (fetched): mitigation applied, no cause stated.
- `a5e543zb`: Discord homepage incident (fetched): one-line 'networking issue routing traffic', no mechanism.
- `pffam7ns`: AWS Health Dashboard root link (fetched, redirects to health.aws.amazon.com which returns no content): no incident text, no cause.
- `iv6ebz7q`: Same AWS status root link as pffam7ns, no cause available.
- `pwhjuphl`: Same AWS status root link as pffam7ns, no cause available.
- `t22geqx7`: Same AWS status root link as pffam7ns, no cause available.
- `rm5zerhw`: Same AWS status root link as pffam7ns, no cause available.
- `qitwivep`: Same AWS status root link as pffam7ns, no cause available.
- `ikcbyiq5`: Same AWS status root link as pffam7ns, no cause available.
- `fg647axo`: vLLM H3-Image day-0 (fetched): release with mechanism (Ulysses SP4, AdaLN caching, tiled VAE) but no live Explainer teaches sequence parallelism for diffusion, and a second vLLM Brief on the same day repeats last week's Kimi day-0 Brief.
- `66tamqmo`: PlanetScale Heroku migration (fetched): general description of the Heroku Migrator tool with anonymous customers, not a named migration event; strong future Comparison candidate (trigger-based vs logical replication vs snapshot) next to planned-database-switchover.
- `3u3twadt`: Flickr to ElastiCache for Valkey (fetched): real migration with mechanism, but no live Explainer teaches its lesson (decouple client from infrastructure change, verify replication before cutover); the Topstep Valkey Brief already covers Redis-to-Valkey this week.
- `5yr6p7kc`: Cloudflare QuePaxa / Meerkat (fetched): experimental system, not a release or incident, so no Brief event; strong future Concept candidate (leaderless consensus) beside the etcd Repo Profile.
- `3f744wo7`: Withdrawal of primitive patterns preview (fetched): real design change for JDK 28 but the page gives no reasons and no live Explainer teaches the lesson.
- `tuddqnt4`: Aurora MySQL 8.4 OOM (fetched): how-to on a known behaviour with a reproduced test, not an event.
- `ubit4guw`: Fastly static publishing v8 (fetched): tool major release, no live Explainer on edge static publishing.
- `q62iodb5`: DuckDB Ducks in a Column (fetched): tutorial on columnar features, not an event.
- `u7aui2bd`: DuckDB local analytics (fetched): interview essay, not an event.
- `gfaluotq`: Zillow 100M notifications on Step Functions: engineering post, no event for a Brief; Case Study candidate but Explainer cap is 0 today.
- `luif2rbp`: Go turns 17: anniversary announcement.
- `uqm5rpbi`: Python 3.16.0 alpha 2: pre-release announcement.
- `poot6ukk`: Rust GSoC results: community post.
- `azacbbvq`: Vite Ecosystem Fund: funding announcement.
- `g2x5llnm`: APNIC submarine cable podcast: no readable text.
- `drobnrha`: APNIC survey council response: community governance, not an event.

</details>


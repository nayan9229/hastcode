# Daily review queue, 2026-11-26

Drafted by `make daily` (ADR-0009). Every Piece is `draft: true` and unverified. First tick each draft's self-edit boxes while you read it, before you verify ([the checklist](../../../docs/editorial/self-edit-checklist.md)): they check the writing, not the facts, and the self-edit does not replace verification. Then verify every claim against its Source, then `make release` (ADR-0012): it sets `verified:`, removes `draft: true` and keeps each Publication Day. Ticked self-edit boxes never stand in for `verified:`.

## Run 20261126-073053

### GitHub's Git LFS overload: retries hid one region's shortfall

`content/briefs/github-lfs-unauthenticated-traffic-overload.md` · brief · 241 words · fact-checker: 25 claims checked, 2 fixed, 0 removed

**Publish within 72 hours of the Lead's date**, with or after its Explainer: metastable-failures

**Drafter asks you to check:**
- [ ] Window: Source says 'On November 25, 2026, between 21:30 UTC and 00:28 UTC'. Brief says '21:30 UTC on 25 November to 00:28 UTC the next day'. The 26 November date for 00:28 comes from the update timestamps (mitigation update posted Nov 26, 00:28 UTC), not from the summary sentence.
- [ ] Numbers: 'an average of 2% of Git LFS requests failed, peaking at 4%'. Brief: 'On average 2% of Git LFS requests failed, 4% at the peak'.
- [ ] Retry masking: Source says 'Many more requests initially failed but succeeded when Git LFS clients automatically retried them, which caused slower uploads and downloads for affected users'. Brief's 'The 2% leaves out every request a retry saved' is our reading of 'Many more'; the Source does not say outright that the 2% is counted after retries.
- [ ] Quote (the only one): "unable to respond to requests in time". Source: 'leaving them unable to respond to requests in time'. Check it is exact.
- [ ] Cause wording: Brief says 'heavy automated, unauthenticated traffic' for the Source's 'a large volume of automated, unauthenticated traffic', and 'one region's Git LFS servers' for 'the Git LFS servers in one of our regions'. The region is not named in the Source.
- [ ] Mitigation: Source says 'moving most Git LFS traffic from the affected region to other regions with more capacity'. Brief: 'moved most Git LFS traffic to regions with more capacity'.
- [ ] Analysis, not in the Source: 'Other regions had more capacity, yet it did not help until GitHub moved traffic to it'. The Source states the capacity and the move; the 'did not help until' framing is ours and carries no 'we infer'.
- [ ] Inference 1 (labelled): 'We infer each retry also added load to servers already past their limit.' The Source does not say retries added load.
- [ ] Inference 2 (labelled): 'We infer the servers admitted more work than they could finish.' It rests on 'The report describes late responses and mentions no refusals': the Source says requests 'failed' and servers were 'unable to respond ... in time'; it is silent on whether any request was refused early.
- [ ] 'A timed-out request wastes the server's work, then returns as a retry: the loop behind metastable failures' is grounded only in the Explainer (its defence-three section, citing Yanacek's Builders' Library article), not linked to an outside Source. The Brief does not claim this incident was metastable.
- [ ] Follow-up: Source says 'We are working to limit the impact of high-volume unauthenticated traffic on Git LFS, and to add capacity in the affected region'. Brief rewords it as 'working to limit what heavy unauthenticated traffic can do to Git LFS and to add capacity in that region'.
- [ ] Omitted from the Brief though in the Source: 'Other Git operations were not affected', the 'small number of customers' wording, and the 23:06 UTC first status post (the 21:30 start is from the summary only).
- [ ] Freshness: Lead's first post is 25 Nov 2026 23:06 UTC; the 3-day window closes 28 Nov 23:06 UTC.
- [ ] Description is 143 characters by hand count; title 61.

**Fact-checker could not confirm:**
- [ ] Title 'retries hid one region's shortfall' and the topic sentence 'Retries hid the overload's size' are our reading, not the Source's words. Source supports that many more requests failed first and succeeded on retry, and that the affected region lacked capacity (traffic moved to 'other regions with more capacity'; GitHub will 'add capacity in the affected region'). Left as is; Author to decide whether the title is too firm.
- [ ] '00:28 UTC the next day': the summary says 'On November 25, 2026, between 21:30 UTC and 00:28 UTC' with no end date. 26 November comes from the update timestamps (mitigation update Nov 26, 00:28 UTC; resolved Nov 26, 00:43 UTC) and 'Error rates returned to normal by 00:28 UTC'. Consistent, but not in the summary sentence.
- [ ] 'A timed-out request wastes the server's work, then returns as a retry: the loop behind metastable failures' and 'admit only what you can finish before the client gives up' are supported only by the linked Explainer (content/explainers/distributed-systems/metastable-failures.md, 'Defence three', citing Yanacek's Builders' Library article), not by an outside link in the Brief. The Brief does not claim this incident was metastable.
- [ ] 'mentions no refusals': confirmed only as absence. The Resolved summary does not contain rejected, refused, rate limit, throttle, shed, queue or timeout; closest wording is 'unable to respond to requests in time'.
- [ ] Source could be read only through a fetch agent that returns short excerpts (125 characters at most). Every sentence the Brief relies on was quoted exactly in excerpts, but I did not see the summary as one continuous text.
- [ ] Omitted from the Brief though in the Source: 'Other Git operations were not affected', the 'small number of customers' wording (23:10 and 23:45 UTC updates), the 23:06 UTC first post. Not errors.
- [ ] Freshness: Lead's first post is 25 Nov 2026 23:06 UTC, so the 3-day window closes 28 Nov 2026 23:06 UTC. Not enforced by the build.

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

- Fixed: Why it matters, para 1: 'yet it did not help until GitHub moved traffic to it' read more into the Source than it states (Source gives only the mitigation: 'moving most Git LFS traffic from the affected region to other regions with more capacity'). Now labelled: 'yet we infer it did not help until GitHub moved traffic to it.'
- Fixed: Why it matters, para 2: 'The 2% leaves out every request a retry saved' was stated as fact and linked to the Source, but the Source says only 'an average of 2% of Git LFS requests failed, peaking at 4%' and 'Many more requests initially failed but succeeded when Git LFS clients automatically retried them'; it never says the 2% is counted after retries. Now 'We infer the 2% leaves out every request a retry saved' (link kept). Following sentence changed to 'We also infer each retry added load...' to avoid a repeated opener; still labelled.

</details>

### vLLM's ERNIE 5.5 Flash plugin: model code outside the core

`content/briefs/vllm-ernie-55-flash-plugin-day-0.md` · brief · 243 words · fact-checker: 18 claims checked, 3 fixed, 0 removed

**Publish within 72 hours of the Lead's date**, with or after its Explainer: vllm-vs-deepinfer, expert-parallelism

**Drafter asks you to check:**
- [ ] Word count 243 is my hand count, headings included; the build's brief-length rule passed, but I did not get a count from the build.
- [ ] Quote, exact words: the plugin "mainly contains the ERNIE model code and a set of kernel optimizations" (plugin section of the vLLM post). Only quote in the Brief.
- [ ] Date and byline: post dated 24 November 2026, by "Baidu ERNIE Team and vLLM Team". Post gives no time of day or time zone, so the Brief gives the date only.
- [ ] "release-day support": post title says "Day-0", intro says vLLM supports the model "on release day through the vLLM-ERNIE plugin". Post does not state the model's own release date.
- [ ] MoE figure: architecture table row reads "512 routed experts (top-10) + 1 shared expert". Brief says the router "picks the top 10 of 512 routed experts"; shared expert left out.
- [ ] Install claim: "a standalone package that you install into the same environment as the official vllm==0.32.0 wheel"; Brief says "installed beside".
- [ ] Entry point: post says vLLM finds the plugin through the vllm.general_plugins entrypoint. Brief links the vLLM plugin docs page https://docs.vllm.ai/en/latest/design/plugin_system.html (loaded 200 via fetch agent); confirm the link resolves and names that group.
- [ ] Number: "The fused kernels run 1.55× to 6.53× faster than naive implementations." Brief adds "measured per kernel": the per-fusion bullets say "per kernel"; the range sentence itself does not. Gated attention's 1.55×–1.75× is prefill only; Brief does not say so.
- [ ] "each merging several operations into one" is my gloss of fusion, from the post's bullets (one GEMM for all router logits, routing chain in one kernel, single-pass fusion). Check it is fair for all four fusions.
- [ ] "mostly come from vLLM itself": Source sentence lists attention kernels, the KV cache manager, speculative decoding, parallelism and the server. Brief names only server, KV cache manager and parallelism.
- [ ] "The plugin has added features without waiting for a vLLM release": from the "Since the Base release, the plugin has added…" sentence. The fetch agent paraphrased its ending ("none of which had to wait for a vLLM release"); check the exact wording supports this.
- [ ] "applies hooks to shared vLLM code": Source says "applies only a few hooks to shared vLLM code, and those hooks take effect only for the ERNIE architectures". Brief drops "only a few" and the ERNIE-only scope; check that is not unfair in a sentence about cost.
- [ ] "runs on one vLLM wheel, 0.32.0" links the vLLM-ERNIE README: "It runs on the released vLLM 0.32.0 wheel, unmodified, on NVIDIA GPUs only." Description also says "unmodified vLLM 0.32.0".
- [ ] "ERNIE's code lives in Baidu's plugin": rests on the repo being github.com/baidu/vLLM-ERNIE and the post's "mainly contains the ERNIE model code".
- [ ] Inference 1 (labelled): a general-purpose engine keeps pace because per-model speed work moves outside the core. The post never frames it this way and never mentions DeepInfer; the link to vllm-vs-deepinfer is ours.
- [ ] Inference 2 (labelled): the fused kernels are among the plugin's "kernel optimizations". The post does not say where the fused kernels live; its acknowledgements say fused-MoE and attention components reuse community implementations, which may cut against this. The title avoids the claim ("model code outside the core").
- [ ] Inference 3 (labelled): each vLLM upgrade needs a plugin retest. Ours, from the hooks plus the single pinned wheel.
- [ ] Serve command: post shows --data-parallel-size 4 and --enable-expert-parallel "on a single 4-GPU node"; Highlights say "Data-parallel attention with expert parallelism". "each GPU holds whole experts" is from our expert-parallelism Explainer, not the post.
- [ ] Unread: the post's three images (hero.png, flashindexer.png, vllm-ernie-plugin.png) could not be read. If the plugin diagram shows where kernels live, it may settle inference 2.
- [ ] Left out on purpose: 280B/14B parameters, 1M context, FlashIndexer 1.4×, DSpark, quantization, hardware. The fetch found no hardware or end-to-end figure in the post's text; the Brief makes no claim about that absence, but the closing advice ("benchmark whole requests, not kernels") leans on it.

**Fact-checker could not confirm:**
- [ ] Inference 2 (fused kernels are among the plugin's kernel optimizations) stays an inference, labelled "we infer" with "the post does not say so". Post gives no location. Repo file listing shows mhc_kernels.py, nvidia/indexer_fused_topk.py, nvidia/moe.py, nvidia/attention.py and patches.py, pure Python with no native build, but no source file named for VHA, MultiMax or the gate; file contents not read. Acknowledgements say many fused-MoE and attention components reuse community implementations.
- [ ] The 1.55x lower bound comes from the gated-attention bullet, which reads "1.55x to 1.75x faster per kernel in prefill". Brief says "measured per kernel" without the prefill qualifier. Left as is: the Source's own range sentence ("The fused kernels run 1.55x to 6.53x faster than naive implementations.") carries no qualifier either. Author may want to add "in prefill for one".
- [ ] Fetch tool capped quotes at 125 characters, so long sentences were confirmed by exact fragment, not whole. Confirmed whole and exact: the Brief's one quote sentence, the hooks sentence, the "mostly come from vLLM itself" sentence, the range sentence, the first sentence of the post. Confirmed by exact fragment only: "none of which had to wait for a vLLM release", "vLLM supports the model on release day through the vLLM-ERNIE plugin", "a standalone package that you install into the same environment as the official vllm==0.32.0 wheel", "loads through its `vllm.general_plugins` entrypoint", "on a single 4-GPU node".
- [ ] Post's three images (hero.png, flashindexer.png, vllm-ernie-plugin.png) not read; all have empty alt text and no caption. The plugin diagram could settle inference 2.
- [ ] Post does not state the model's own release date. "release-day support" rests on the title ("Day-0 Support") and the intro fragment "on release day".
- [ ] "runs on one vLLM wheel, 0.32.0": README says "It runs on the released vLLM 0.32.0 wheel, unmodified, on NVIDIA GPUs only." and states no other version; repo pyproject.toml pins dependencies = ["vllm==0.32.0"]. Supported, but the word "one" comes from the pin, not from a README sentence. Left unchanged.
- [ ] Title says "vLLM's ERNIE 5.5 Flash plugin": the plugin is Baidu's repo (baidu/vLLM-ERNIE) and the model is "ERNIE 5.5 Flash Preview". Readable as "plugin for vLLM"; left unchanged for the Author to judge.
- [ ] Word count not printed by the build. Build raised no brief-length error; hand count after edits is about 245 (drafter's 243, plus 2 net).
- [ ] Closing advice ("benchmark whole requests, not kernels") leans on the post giving no end-to-end throughput or latency figure and no benchmark hardware; fetch of the post text found none, images unread.

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

- Fixed: "whose router picks the top 10 of 512 routed experts" -> "that selects the top 10 of 512 routed experts": post describes a Multi-Router Latent MoE (several routers, scores summed, then top-10), so a single "router" was wrong. Table row "512 routed experts (top-10) + 1 shared expert" confirmed.
- Fixed: "applies hooks to shared vLLM code" -> "applies a few hooks to shared vLLM code": Source sentence is "The plugin applies only a few hooks to shared vLLM code, and those hooks take effect only for the ERNIE architectures." Dropping "only a few" overstated the cost.
- Fixed: "ERNIE's code lives in Baidu's plugin" -> "ERNIE's model code lives in Baidu's plugin": matches Source ("mainly contains the ERNIE model code") and the Brief's own title.

</details>

### RPKI trust anchor tiebreaking: why an old TA certificate loses

`content/explainers/apis-networking/rpki-trust-anchor-tiebreaking.md` · concept · 2400 words · fact-checker: 78 claims checked, 9 fixed, 2 removed

**Drafter asks you to check:**
- [ ] Every quote was relayed by a fetch agent as a short fragment, not read by me in the raw page. Check each against its Source, especially RFC 10097 Section 2 step wording: "a current, validly self-signed RPKI CA certificate", "the locally cached TA", "If the notBefore of the newly retrieved TA is less recent, use the locally cached TA.", "to use the TA that is currently published".
- [ ] Claim of absence: RFC 8630 Section 3 (four steps, "an RP SHOULD:") says nothing about what to do with a certificate the RP already holds. The agent reported no cached-copy fallback in RFC 8630; confirm.
- [ ] Claim of absence: RFC 10097 steps 1 and 2 do not say what an RP with an empty cache does when the fetch or validation fails; only step 3 has a no-cache rule. Also that the RFC sets no maximum validity period.
- [ ] Hook numbers from Snijders' 24 June 2024 post: ARIN old certificate 1 Oct 2015 to 1 Oct 2025; AfriNIC old 22 Dec 2016 to 20 Dec 2026; "likely to end up rejecting 100% of subordinate products" (AfriNIC) and about 50% (ARIN). I quoted a fragment and paraphrased the ARIN half; check the full sentence's wording and that 100% belongs to AfriNIC.
- [ ] Certificates published in June 2024 per the same post: AfriNIC 30 Mar 2020 to 28 Mar 2030; ARIN from 30 Jan 2024; LACNIC 5 Mar 2024 to 5 Mar 2124; validity periods ARIN about 2 years, APNIC 5, AfriNIC 10, LACNIC and RIPE NCC 100. These are June 2024 values, stated as such; the piece does not say what any registry publishes today.
- [ ] Worked example inferences (marked "we infer"): that the 2016 AfriNIC certificate is validly self-signed and carries the same key as the TAL. The post does not state either.
- [ ] rpki-client 9.4 dates: release notes write "February 2nd, 2026" and "March 3rd 2027"; the piece paraphrases as 2 February 2026 and 3 March 2027. Check 15-year warning wording ("expiry date more than 15 years into the future") and the 3-year rejection.
- [ ] Not sourced and left out: whether LACNIC and the RIPE NCC actually reissued before 2 February 2026, and whether the 9.4 schedule still holds in current rpki-client. A Reader in November 2026 may ask; add only with a first-party notice.
- [ ] rpki-client 9.1 release date: the release notes print none; "June 2024" comes from the draft -06 reference list. The RFC dropped that appendix, so the draft is the only Source for the two named implementations.
- [ ] rpki-prover PR 218: merged 3 September 2024, title and commit name "Implement some tiebreaker scheme", and the "some version of" comment. The PR cites the earlier draft-spaghetti name. I call the commenter "the project's author" (GitHub user lolepezy).
- [ ] Mail archive URLs for the Bruijnzeels (zeleBN8M5i3P8lCUvtf9zvzgDQs) and Njuki (T3RxVOL5yx1Z-GbfhlYMM_IBYRY) replies came from an agent's report of the thread, not from a URL I supplied. Click both. Quotes: "We'll schedule to reissue", "blocking issue", and Harrison's "APNIC will reissue its TA certificate with a shorter validity period".
- [ ] Section anchors assumed, not all confirmed: rfc8630.html #section-2.1, #section-2.2, #section-2.3, #section-3, #section-4, #section-5; rfc6487.html #section-4.6 and #section-4.6.1; rfc6480.html #section-1 and #section-2.4. RFC 10097 #section-1/2/3 were confirmed. Check that the stable-key, same-URI, no-CRL and multiple-URI sentences sit in RFC 8630 Section 2.3.
- [ ] RFC 9691 is cited only for "support planned key rollovers" (abstract fragment). The -P option claim and the "TA validity period exceeds 15 years." output come from Snijders' 19 December 2024 reply.
- [ ] Spelling follows the Sources: "AfriNIC" (the registry now styles itself AFRINIC) and "the RIPE NCC". Change if house style differs.
- [ ] Diagram: check on the Preview that the 14-edge flowchart renders and reads at page width.

**Fact-checker could not confirm:**
- [ ] Raw page text never read directly: curl was denied and the fetch agent could only relay fragments under 125 characters. Each quote was confirmed as an exact-string match by the agent (RFC 10097 Section 2 steps relayed fragment by fragment), not by my own read of the raw page. Author should click through the RFC 10097 Section 2 quotes once.
- [ ] "arrive at deterministic outcomes" (RFC 10097 Section 3): one fetch reported exact match, a second fetch paraphrased the sentence as "reach deterministic outcomes". Kept; worth one click.
- [ ] HTML anchor ids (rfc8630 #section-2.1/2.2/2.3/3/4/5, rfc6487 #section-4.6/4.6.1, rfc6480 #section-1/2.4, rfc10097 #section-1/2/3): the pages' tables of contents link to them, raw id attributes not inspected.
- [ ] Simulated 2 Feb 2026 run in the 19 Dec 2024 mail reportedly prints "TA cert rejected: validity period exceeds 3 years." (15 years would be expected on that date). Piece quotes only "TA cert rejected", which holds either way.
- [ ] Not sourced, left out as drafted: whether LACNIC and the RIPE NCC reissued before 2 February 2026, and whether the 9.4 schedule still holds in current rpki-client.
- [ ] Worked-example inferences (2016 AfriNIC certificate validly self-signed, same key as TAL): the June 2024 mail says nothing on keys or self-signing. Both stay marked "we infer".
- [ ] "any two certificates from one TA can be ranked" is the piece's own gloss ("In our words") on "preorder"; a preorder is not total by definition, though this scheme's is.
- [ ] Diagram (14 edges) matches the six steps as sourced; build passes, but render at page width not viewed on the Preview.
- [ ] House style for "AfriNIC" / "the RIPE NCC" left as the Sources spell them.

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

- Fixed: RFC 8630 claim of absence was too strong: Section 3 does mention "the locally cached copy of the retrieved TA" (refetch before it expires). Now reads: it mentioned a cached copy only to say when to fetch again, not which certificate to keep.
- Fixed: `-P` claim was wrong: "TA validity period exceeds 15 years." came from the run with NO -P flag (19 Dec 2024 mail). The -P 1769990401 run simulates 2 Feb 2026 and prints "TA cert rejected: ..." for the RIPE NCC and LACNIC. Sentence now says Snijders used -P to simulate 2 February 2026 and the RIPE NCC output said "TA cert rejected".
- Fixed: Gate sentence quoted "the locally cached TA" for all three gates; RFC 10097 step 1 says "the locally cached copy of the TA referenced by the TAL previously retrieved" (only steps 2-5 use the quoted phrase). Quote marks dropped: "the RP uses its locally cached copy".
- Fixed: rpki-prover: "the project's author" changed to "the repository's owner" (GitHub marks lolepezy Owner; identity with Mikhail Puzanov not sourced). Added that the PR cites the draft under its earlier draft-spaghetti name (datatracker: draft-spaghetti-sidrops-rpki-ta-tiebreaker "Replaced by" draft-ietf-sidrops-rpki-ta-tiebreaker).
- Fixed: rpki-client 9.4 paragraph: "rpki-client enforces it" and present-tense rejections were not sourced for current rpki-client. Now tied to the 9.4 notes: "adopted it", "Its notes say it warns ... will reject such certificates after 2 February 2026 ... From 3 March 2027 it will reject".
- Fixed: Worked example anchored to June 2024 ("Suppose that in June 2024", "current then"): the 2016 AfriNIC certificate expires 20 December 2026, so an undated "is current" goes stale in under four weeks.
- Fixed: "five Regional Internet Registries" link moved from the June 2024 mail (lists five TA certs, never says RIRs or five) to the 18 Dec 2024 roadmap mail ("RIRs have to reissue shorter-lived TA certificates", names all five).
- Fixed: Unlinked claim "never in the tie" (key mismatch) now links RFC 10097 Section 2 (step 3: "If the keys do not match, use the locally cached TA.").
- Fixed: Unlinked claim of absence "sets no maximum" now links RFC 10097.
- Removed: Quote "TA validity period exceeds 15 years." as -P output (misattributed; replaced by short fragment "TA cert rejected" from the simulated run).
- Removed: Quote marks around "the locally cached TA" as wording for all three gates (not step 1's wording).

</details>

### Traffic mix: why key skew slows a lock table and speeds a cache

`content/explainers/distributed-systems/workload-skew-in-load-tests.md` · concept · 2200 words · fact-checker: 87 claims checked, 4 fixed, 2 removed

**Drafter asks you to check:**
- [ ] Microsoft report host: the PDF is Microsoft's own Post Incident Report (EX1550543/MO1550821, dated 20 November 2026) but is served from a third-party mirror, 365fanatic.com; the official copy sits in the M365 admin center behind sign-in. Decide whether a mirror of a first-party document meets the first-party rule. I read the PDF itself (15 pages); every hook fact and quote is from it, with #page= anchors.
- [ ] Report URL: Brooker's own link (…/2026/11/postincidentreport.pdf) now serves an HTML summary page on that site, not the PDF. The piece cites the PDF at …/postincidentreport-mo1550821.pdf. Confirm it resolves and that the #page= anchors land on pages 2, 5, 6, 8 and 15.
- [ ] Preliminary vs final wording: Brooker (18 Nov) quotes the preliminary report ("…because [of] the traffic mix and SKUs used in our labs"). The final report words it differently; the piece quotes only the final report ("the traffic conditions were different from what was observed during the incident" and the finding "Our testing did not adequately represent the combination of hardware and traffic mix…") and says Brooker read the preliminary version. I could not read the preliminary report.
- [ ] Worked example numbers are our arithmetic on Brooker's parameters, not his output (the post prints no throughput figures): lock table about 16,000/s uniform vs about 1,000/s on one key (16-fold); cache about 16,800/s at a 5% hit rate vs 1.6 million/s all hits (about 100-fold). Labelled as our arithmetic / "we infer". Check against the live simulators; note his general claim is "multiple orders of magnitude, in both directions" while our lock-table bound is only 16 times.
- [ ] Lead quotes came through a fetch tool that returned extracted text, not raw HTML: verify each quoted sentence character for character on the live page, especially "can vary by multiple orders of magnitude, in both directions" (the editor's angle said "either direction"; the post says "both directions"), "a coefficient of 0 is uniform", and the OCC / single cost/latency mode sentence ("mode" is italic in the post).
- [ ] Lead date 18 November 2026 comes from the URL and the blog index, not from text on the page. The Lead URL is kept as http:// as make lead wrote it; other Brooker posts are cited at https://.
- [ ] pgbench documentation: the tail of the page was read only through the fetch tool's summary, not verbatim. Verify: `random` is uniform ("equal probability"); `random_zipfian` parameter range [1.001, 1000]; the 2.5 → 5.66 example; "you'll mostly be measuring update contention" when clients exceed scale factor; every transaction updates one row of pgbench_branches, which has one row per scale-factor unit; "at least a few minutes" run advice; and that the #PGBENCH-BUILTIN-FUNCTIONS anchor exists. /docs/current/ pointed at PostgreSQL 19 on the day.
- [ ] Brooker 2023 hot-keys post: the Zipf formula is paraphrased as "key number k is drawn with probability proportional to 1/k^s" from a LaTeX formula; and the claim "errors climbed once the system had scaled to about 10 nodes" for s = 0.7, N = 1,000. Check both against the post, including what a node and an error mean in his simulation.
- [ ] Twitter cache study: the paper's abstract says 153 clusters; the released traces (cache-trace README) cover 54 clusters over one week in March 2020. Both numbers are used, each for its own claim. Only the USENIX abstract page and the README were read, not the paper body.
- [ ] YCSB wiki: workload A "a mix of 50/50 reads and writes", workload B 95/5, and `requestdistribution` default `uniform` with `zipfian` allowed. The wiki pages were last edited in 2020; confirm they still match the current code if that matters.
- [ ] Aurora DSQL user guide anchors (#dsql-transaction-conflicts, #dsql-perf-guidelines) and quotes ("evaluates conflicts at commit time", SQLSTATE 40001, "minimize high contention on single keys"): verify on the live page.
- [ ] Inferences marked "we infer" for the Author to accept or cut: the 16-fold range; why hot keys stay in an LRU cache; a uniform test errs twice; one read path means one read cost in DSQL; a captured trace ages; the skew sweep as a practice (no Source prescribes it); and the parallel between Microsoft's gateway lock contention and Brooker's lock table (the report does not say which properties of lab traffic differed, and never mentions key skew).
- [ ] Brooker is given no job title or employer in the piece, because none of the pages I read states one. Add it if you want it and can source it.
- [ ] Diagram: one Mermaid flowchart drawn from the two models' stated parameters; the post's own figures are interactive simulators, linked beside it. Check it renders on the Preview.

**Fact-checker could not confirm:**
- [ ] Microsoft report host: the PDF resolves (HTTP 200, 15 pages, report date November 20, 2026, Microsoft copyright notice) and every cited fact, quote and #page= anchor (2, 5, 6, 8, 15) matches it. But it is served from 365fanatic.com, a third-party mirror; I could not compare it with the official copy behind M365 admin center sign-in. Author decides whether a mirror of a first-party document meets the first-party rule. Brooker's own post links the same site.
- [ ] 'more than 80% of North America's gateway servers': report p8 says "more than 80% of servers were unresponsive" in the sentence after one about the North American server fleet. North America is the natural reading, not stated in that sentence.
- [ ] pgbench current page (PostgreSQL 19): the fetch tool returns only the first ~40,000 characters verbatim, so the random / random_zipfian text ([1.001, 1000], the 2.5 -> 5.66 example, 'equal probability') and Good Practices were read through the tool's summary of /docs/19 and /docs/18 (identical in both, and supporting the piece). The quoted sentence "you'll mostly be measuring update contention", "at least a few minutes" and the one-branch-row-per-scale-unit table were confirmed verbatim on the short /docs/9.0 page. The #PGBENCH-BUILTIN-FUNCTIONS anchor is a link target on the current page. Author: confirm the random_zipfian lines on the live page.
- [ ] Brooker simulator defaults: '5% of keys', 0.01 ms hit and 1 ms miss appear in the post's diagram and text, but the input default values and the chart axes are script-rendered and were not readable. Check on the live page.
- [ ] Worked-example numbers (16,000/s vs 1,000/s, 16-fold; 16,800/s vs 1.6 million/s, about 100-fold) are our arithmetic, labelled as such. The arithmetic is right (the cache ratio is 95, rounded to 'about 100'). The lock-table diagram says Poisson arrivals into a FIFO queue "never empty at saturation", consistent with a saturated-throughput reading, but the post prints no figures to compare. Brooker's claim is 'multiple orders of magnitude'; our lock-table bound is 16 times. Author to accept or cut.
- [ ] Dates of all Brooker posts (18 Nov 2026, 7 Feb 2023, 27 Aug 2021, 20 May 2025, 18 Oct 2023, 4 and 5 Dec 2024) come from URLs and the blog's 'Similar Posts' lists; no post page prints its date. The post says "last week's M365 outage", consistent with 18 November.
- [ ] Quotes on brooker.co.za were confirmed as exact fragments; the fetch agent capped quotes at 125 characters, so the longest sentences were confirmed in two fragments, not as one string. All quoted strings in the piece are within what was confirmed. The post has "its the one thing" (no apostrophe); the piece quotes from "the one thing", so it is unaffected.
- [ ] Inferences marked 'we infer' are left for the Author: 16-fold range; why hot keys stay cached; a uniform test errs twice; one read path means one read cost in DSQL; a captured trace ages; the skew sweep as a practice; the parallel between Microsoft's lock contention and the lock table (the report never mentions key skew). 'So even published workload studies can lag the mix in production' (Twitter paragraph) is our reading of the abstract and carries no 'we infer'.
- [ ] Brooker's employer is sourceable if wanted: his blog sidebar says "I am an engineer at Amazon Web Services (AWS) in Seattle". The piece gives no title; left as is.
- [ ] Mermaid diagram matches the sourced parameters (one lock per key, LRU cache, 5% of keys, 0.01 ms hits, 1 ms misses) and the build reports no diagram-render error; not viewed on a Preview.

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

- Fixed: Lock-table paragraph: 'Brooker calls this a reasonable mental model of a relational database' did not match the post, which says "This isn't an unreasonable mental model of a relational database". Now: Brooker says it is not "an unreasonable mental model of a relational database" (verbatim fragment, quoted and linked).
- Fixed: Microsoft 365 gateway paragraph: 'the servers stayed healthy for about 16 hours' changed to 'all services stayed healthy'. Report p6 says "all services continued to remain healthy for the next ~16 hours"; its p9 timeline says isolated server failures began in that window.
- Fixed: Why-wrong-mix section: 'Microsoft's lab ran older hardware than its fleet' changed to 'older hardware than its newest servers'. Report p6: the fleet holds both newer and older SKUs; the lab used older SKUs with fewer cores.
- Fixed: Why-wrong-mix section: 'Tool defaults are uniform' softened to 'often uniform', and the YCSB sentence now adds that the packaged workload A sets requestdistribution to zipfian (linked to workloads/workloada in the YCSB repo, also added to Sources). Code confirms: REQUEST_DISTRIBUTION_PROPERTY_DEFAULT = "uniform", but workloada and workloadb both set requestdistribution=zipfian.
- Removed: 'in a single step' from 'one failover changed the gateway's traffic in a single step': not in the report, and its p9 timeline says load on the front-end servers "increased gradually" after the failover.
- Removed: 'from uniform upward' and 'Zipf exponent' from the Practice two sentence on Brooker's curve (now 'throughput against skew'): the charts are script-rendered and the fetch tool could not read their axes or sweep range.

</details>

### vLLM-GR: how vLLM runs beam search inside the engine

`content/explainers/ai-ml-systems/vllm-gr-generative-recommendation.md` · case-study · 2700 words · fact-checker: 96 claims checked, 15 fixed, 5 removed

**Fills Gap Day 2026-11-25** (ADR-0012): `published` is that day; `make release` keeps it.

**Drafter asks you to check:**
- [ ] Word count is an estimate (~2,700 prose words; wc gives 3,251 for the whole file including frontmatter, Mermaid, tables and Sources). Page bytes pass with only 149 bytes of margin.
- [ ] All quotes were relayed by a fetch agent, not read by me on the page. Check each verbatim: "beam selection records the score and parent of each surviving path."; "The next decode uses those parent indices to read the right KV state."; "At beam width 256, vLLM-GR recorded the lowest P99 among the backends shown at batch sizes 4 and 8 on all six accelerators."; "summarize separate benchmark runs with different timing definitions, sample counts, and backend settings."; "large-scale production deployments" (blog). "All beams attend to the same cached prompt."; "The frontend submits the request to EngineCore once."; "an independent extension of vLLM"; "Ranking a supplied candidate set is planned." (README). "does not forward 128 Beam model rows for that first selection."; "intentionally prune the search"; "a fixed generated-token horizon"; "TP, PP, DP, DCP, and PCP sizes equal to 1." (offline example README). "writes current KV, selects successors and applies the parent KV mapping once before the next Forward." (beam_state.md). "Each stage restricts the next SID token to legal continuations of the prefix." (design constraint_tables.md). "ordinary requests are held back"; "Selected tokens stay on the device." (async_scheduling.md). "Running a list as one session does not guarantee identical beam rankings to separate B=1 calls." (onerec_quickstart.md).
- [ ] Offline example README: the fetch agent read its tail (after ~char 39,900) less reliably. The "fixed generated-token horizon" sentence, "Variable-length early stopping is not supported" and "intentionally prune the search" (#npu-per-stage-beam-widths-from-the-cli) come from that part; confirm wording and anchors. The "fixed generated-token horizon" link has no anchor because I could not confirm the section.
- [ ] Editor's angle says the roadmap 'admits' three limits. The post has no limitations section: distributed KV cache, variable-length generation and MoE are roadmap items ('explore'). The draft says so, and sources fixed length and parallelism=1 to the example README's restrictions instead. Confirm that framing is acceptable.
- [ ] Gloss: 'one model on one device' and the description's 'no model sharding' go beyond the Source's words ("TP, PP, DP, DCP, and PCP sizes equal to 1"; roadmap 'Evaluate model sharding'). Marked 'we read the list as' in the body. Also: I call these 'vLLM's parallelism settings' and expand only DCP (via our Explainer); the README does not expand the abbreviations.
- [ ] Upstream vLLM beam search: described from vllm/entrypoints/generate/beam_search/offline.py on main (frontend loop, max_tokens=1, logprobs=2*beam_width, one prompt per live beam). A main-branch link can drift; confirm the code still reads that way, and that 'frontend'/'entrypoint' is fair. The '256 prompts per step, 512 log probabilities each' figure is ours and labelled so.
- [ ] Our arithmetic, labelled as ours, to re-check: 4 x 128 = 512 and 4 x 64 = 256 decode rows; 256 x 1,024 = 262,144 vs 1,024 + 2 x 256 = 1,536 KV token positions (suffix of at most two per beam is our inference from two decode stages); SGLang P99 about 3.5x vLLM-GR's (1,280.1 / 369.4) and SID-GR about 16% higher (428.1 / 369.4) on H20, beam 256, batch 8.
- [ ] Step 2 table: 'Prefill | 4 prompts' is my reading of "one final prompt position per request" for B=4; the README states only the 512 and 256 decode totals and the 128/64/32 widths.
- [ ] Benchmark: A100 at beam width 256 (274.5 ms vs 451.5 ms) left out because the post's sentence does not name the batch size. H20 batch 4 is a tie ('about 191 ms' for both), which sits oddly with 'lowest P99 at batch sizes 4 and 8 on all six'; the table shows it without comment. 'Catalog constraints off' is stated for the RTX PRO 6000 run; I did not claim it for the other NVIDIA GPUs.
- [ ] Hook: '128 candidate item IDs' is from the post's OneRec-1.7B example sentence ("from a viewer's history to 128 candidate SIDs"); 'three tokens long' is from docs/configuration/constraint_tables.md ("three-token semantic ID"). Confirm both, and that 'announced on 24 November 2026' matches the post date.
- [ ] Inferences labelled 'we infer' to judge: hardware coverage was part of the reason to extend vLLM; admission rejects up front to avoid wasted stages; fixed length is one reason the host can prepare stages early; an instance is best run as a dedicated recommender; SGLang's 2W pool means the gap is not like for like.
- [ ] Discrepancies in the Sources: (1) README says vLLM's prefix cache already skips prefill for matching prompt blocks across requests, while the roadmap lists 'Reuse KV blocks for an unchanged prefix across requests' as future work; the draft cites the README claim. (2) The blog's framework figure lists HSTU and OneTrans as model examples but the roadmap lists them as future; the draft says only OneRec-1.7B is served, per the README. (3) OneRec model card metadata shows '2B params' against the 1.7B name; the draft avoids a parameter count.
- [ ] 'itemic tokens' is the OneRec model card's term ("Itemic Tokens derived from hierarchical vector quantization"); the card does not itself say 'semantic ID'. The SID definition sentence ('an item's ID written as a few model tokens') is our wording.
- [ ] Constraint Table details from the design note: 'trie', 'flat offset and value arrays', K cap 256 vs 1,024 unconstrained CUDA, online replacement on CUDA vs restart on Ascend. I did not say logits are masked: the docs do not say how selection applies the trie.
- [ ] Contributor names shortened: the post prints 'TikTok-Global E-commerce' and 'ByteDance-Commerce & Advertising'; the draft says 'TikTok, ByteDance'. Byline is 'vLLM-GR Team' with no individual authors.
- [ ] Anchors (#about, #why-vllm-gr, #inside-the-architecture, #key-features, #performance-evaluation, #future-roadmap, #getting-started, #acknowledgements on the blog; GitHub README and doc anchors) were reported by the fetch agent; click each. EngineCore glossed as 'where vLLM's scheduler runs' from the post's figure legend.
- [ ] Mermaid diagram is ours; check it renders on the Preview and that the Step 1-6 numbering in the boxes (our walk order, with the Constraint Table as Step 5 and the host's early preparation as Step 6) reads right.

**Fact-checker could not confirm:**
- [ ] Verbatim checks went through a fetch agent capped at 125 characters per quote; I could not fetch pages directly (curl denied). Every quote in the draft was returned as an exact match, but the Author should eyeball punctuation on the live pages (the blog uses curly apostrophes).
- [ ] Anchors: blog slugs come from the page's table-of-contents links; GitHub anchors are derived from heading text by the slug rule. All headings exist; none was clicked in a browser. New anchors added by me: design constraint_tables.md#artifact-layout and example README #per-stage-beam-widths-through-the-python-api.
- [ ] "EngineCore, where vLLM's scheduler runs" (body and diagram 'Host: EngineCore'): the post body says only "EngineCore owns each beam session from admission to completion" and "The scheduler dispatches prefill and decode stages to a worker"; the figure legend says "extensions to EngineCore scheduling". Left as a gloss.
- [ ] Diagram places 'Step 6: the host prepares the next stage' inside the EngineCore box; the async note says the CUDA worker prepares the next input layout while EngineCore issues the stage early. Close enough for 'host', but the Author may want to judge. Mermaid rendering not checked on a Preview.
- [ ] "NVIDIA SID-GR was faster at both widths" at batch size 1 on RTX PRO 6000: confirmed by the post's sentence under the RTX PRO 6000 heading; the beam-width-128 millisecond values are only in a figure the fetch could not read. Same for the 'lowest P99 at batch sizes 4 and 8 on all six accelerators' claim: confirmed as the post's sentence, per-device figures not readable.
- [ ] H20 batch 4 is a tie ("both measured about 191 ms") inside a 'lowest P99' claim; table shows it without comment, as the drafter noted.
- [ ] Description's "no model sharding" and the body's "one model on one device" are glosses on "TP, PP, DP, DCP, and PCP sizes equal to 1." plus the roadmap's "Evaluate model sharding"; body marks it 'we read the list as'. 'vLLM's parallelism settings' is not expanded by the README.
- [ ] Heading "every search runs a fixed number of tokens": the README's fixed-horizon sentence is scoped to the per-stage-width path, and Current restrictions says unconstrained generation supports 1 through gr.max_decode_steps tokens (a per-request length). The general claim rests on the roadmap listing variable-length generation as future work.
- [ ] "A session can carry more than one viewer, and the count is its batch size": the post does not state current batching; support is the example README ("Logical B=1..8 with shared generation parameters") and the quickstart ("Running a list as one session").
- [ ] "ranks and packs it on the GPU and copies one compact payload to the host" is from lines of a text flow diagram under Public entrypoint, not a prose sentence.
- [ ] Source discrepancy left as drafted: README says vLLM's prefix cache skips prefill for matching prompt blocks across requests, while the roadmap lists reusing KV blocks for the unchanged prefix across requests as future work.
- [ ] Our arithmetic re-checked and kept labelled as ours: 4 x 128 = 512, 4 x 64 = 256; 256 x 1,024 = 262,144 vs 1,024 + 2 x 256 = 1,536; 1,280.1 / 369.4 = 3.47 (about 3.5x); 428.1 / 369.4 = 1.159 (about 16%). The suffix-of-two-per-beam assumption is an inference and ignores paged block allocation.
- [ ] Upstream offline.py is linked on main and can drift; it read as described on 2026-11-26.
- [ ] Word count is still about 2,600 prose words (3,150 by wc for the whole file) against the preferred 2,400; page bytes are within the limit with 188 bytes of margin.

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

- Fixed: Wrong anchor: "intentionally prune the search" sits under #per-stage-beam-widths-through-the-python-api, not #npu-per-stage-beam-widths-from-the-cli; link corrected.
- Fixed: "a fixed generated-token horizon": anchor added (#per-stage-beam-widths-through-the-python-api) and the claim scoped to the README's per-stage-width path ("Variable-length early stopping is not supported by this dynamic-width implementation").
- Fixed: Wrong anchor: "EngineCore can issue the next stage before the worker acknowledges the previous one" and "Selected tokens stay on the device." are in the async_scheduling.md introduction, not #cuda-dispatch-timeline; anchor dropped there and moved to the layout/forward split sentence, which is what that section says.
- Fixed: "generates three [tokens]" was stated for the whole benchmark; the post says it only for the RTX PRO 6000 run. Now: overview charts hold input at 1,024 tokens; the one run described in full generates three.
- Fixed: P99 gloss "the time within which 99% of batches finished" applied to every figure; whole-batch timing is defined only for the RTX PRO 6000 run (figures use "different timing definitions"). Reworded to timed samples, and "whole-batch P99" added to the RTX PRO 6000 description.
- Fixed: SGLang candidate pool (2W vs K equal to W) is stated only in the RTX PRO 6000 paragraph; draft applied it generally and asserted it "favours vLLM-GR". Now scoped to that run, kept as a labelled inference, with a note that the post does not say which pool the H20 run used.
- Fixed: Upstream vLLM loop: no function named generate is called per step (code calls _render_and_run_requests); "one generate call per step" changed to resubmitting prompts / a new set of requests at every step (TL;DR and body). max_tokens=1, logprobs=2*beam_width and one prompt per live beam confirmed in offline.py on main.
- Fixed: Constraint Table storage: "flat offset and value arrays" is not in the design note (it says a single-Trie safetensors file contains six tensors; offsets address concatenated arrays). Reworded to "stored as tensors in a safetensors file" and linked to #artifact-layout; "ahead of time" dropped.
- Fixed: Online table replacement: Source says "CUDA asynchronous Beam serving supports online table replacement", not CUDA in general; qualifier added.
- Fixed: README example: per-stage widths 128/64/32 are for a B=1, 2 or 4 batch; only the 512 and 256 decode-row totals are stated for B=4. Sentence split so the widths are not tied to B=4, and the 'Prefill | 4 prompts' row is now labelled as our reading.
- Fixed: "Admission is all or nothing" / "Refusing at admission": the Source sentence does not say "at admission". Reworded to "A batch runs whole or not at all" and "Refusing up front" (still a labelled inference).
- Fixed: "The team publishes no timeline of the overlap": the note has a schematic swimlane figure; reworded to "no measured timeline".
- Fixed: Unlinked claims and quotes now linked: TL;DR bullet 3 (beam session, #inside-the-architecture); quote "does not forward 128 Beam model rows for that first selection."; quote "summarize separate benchmark runs..."; per-parent K cap 256 vs 1,024 (example README, #per-stage-beam-widths-through-the-python-api); benchmark setup sentence above the table (#performance-evaluation).
- Fixed: Closing analysis ("The design travels beyond recommendation") stated as fact; now begins "We infer".
- Fixed: Two sentence-length build warnings introduced by the edits were split; build is clean.
- Removed: Unsupported gloss on the second cost: "because each trip to the frontend is time the device spends waiting" (post says only "duplicate work and per-step coordination with the host").
- Removed: Inference paragraph that hardware coverage was part of the reason to extend vLLM (with its README link): redundant with the sourced sentence on NVIDIA GPUs and Ascend NPUs; cut for page size.
- Removed: "The post does not measure the upstream loop, so we cannot say how slow it was." (duplicate of the later 'no comparison with upstream vLLM's loop'; cut for size).
- Removed: Filler/opinion sentences cut for size: "Whichever interface it uses, the engine sees one unit of work."; "Two settings shape the search."; "This step is what makes a wide search affordable."; "Here is the common confusion."; "Measure P99 with your catalog constraint on."; "a separate repository"; "then turn to cost and evidence".
- Removed: Duplicate #performance-evaluation link on "the RTX PRO 6000 run" (the quote directly above now carries the link to the same section).

</details>

### Single-server inference: why a few GPU nodes need no RDMA fabric

`content/explainers/infrastructure-cloud/small-gpu-cluster-inference-network.md` · concept · 2400 words · fact-checker: 118 claims checked, 26 fixed, 0 removed

**Fills Gap Day 2026-11-25** (ADR-0012): `published` is that day; `make release` keeps it.

**Drafter asks you to check:**
- [ ] Word count ~2,400 is an estimate: no wc run was possible in this session (only the check script was permitted).
- [ ] The Lead was read through a fetch agent that capped quotes at ~125 characters. Check every quoted Lead phrase character for character against the page, especially: "4 servers × 8 GPUs", "The network problem is not how the GPUs talk to each other.", "the model fits in one server", the 300 bytes / 50 messages / 120 kbit/s sentence, "tens to a few hundred streams at once and no more", "a few hundred GBs", "no protocol or storage overhead", "so the adapters can read and write GPU memory directly", "both place connections and see nothing of the requests inside them", "can cost a hundred times the next", "receives text in clumps", "Set all of them to the longest silence you accept.", "by segment or by Access Control List (ACL)", the two-traffic-classes quote (capital 'One' after the colon is as reported), the shaper/policer sentence, "CS1 or Lower-Effort Per-Hop Behavior (LE PHB)", "could still be 3 to 5 years away, if not longer".
- [ ] "Single-server inference" is OUR name; the Lead names no pattern (the piece says so in the TL;DR and first section). Confirm you accept the coined term in the title, or retitle around the Lead's "a few GPU nodes" / "small GPU cluster".
- [ ] Author byline and bio: "Sabina Anja", "Chief Technologist at Broadcom" as reported by the fetch agent; confirm on the page. No pronoun is used for the author.
- [ ] Hook sentence "Reference designs for large AI clusters add a second network for GPU-to-GPU traffic, with 400G ports tied to it" is our compression of the Lead's reference-design, backend-fabric and 400G passages; check it does not overstate.
- [ ] GPU memory arithmetic is ours, not the Lead's: the Lead does not name the baseline GPU. We took the 48 GB and 80 GB sizes it mentions only as a comparison (8 x 48 = 384 GB, 8 x 80 = 640 GB; minus 141 GB = 243 to 499 GB; about 24 streams at 10 GB, about 499 at 1 GB). Cut if you think this reads as the Lead's claim.
- [ ] Our arithmetic to verify: 120 Mbit/s is under 0.5% of 25 GbE (0.48%); 141 GB = 1,128 gigabits; one of four servers loading = a quarter of capacity.
- [ ] Lead gives the 10 GbE load time two ways: "113 seconds" (Model loading section) and "about two minutes" (hardware section). The piece uses 113 seconds.
- [ ] Lead says both that the model file is spread across the GPUs in a chassis and that each server holds a complete copy; the piece keeps both (complete copy per server, split across its eight GPUs).
- [ ] NVLink ("a direct GPU-to-GPU link") and PCIe ("the bus that connects cards inside a server"), VRF ("a separate routing table"), MTU ("the largest packet a link carries") and jumbo frames ("a larger MTU") are our one-line glosses, not sourced definitions.
- [ ] Tensor parallelism description (a slice of every layer, exchange at every layer for every token) is from the Lead; vLLM's page does not define it in words (it links Megatron-LM).
- [ ] vLLM's L40S / no-NVLINK sentence sits inside the "Edge case: uneven GPU splits" note, not a general recommendation; the piece says vLLM's docs "point to it when" GPUs lack NVLINK. Check that wording against the page and the Lead's stronger reading.
- [ ] vLLM docs links have no section anchors: the rendered site exposed no heading ids. Pages cited are the 'latest' developer preview (parallelism_scaling, usage/security, features/disagg_prefill, features/automatic_prefix_caching); quotes were checked as short fragments only.
- [ ] vLLM security quote: the page says "All communications between nodes in a multi-node vLLM deployment are insecure by default"; we quote only "are insecure by default". The Lead's wording (unencrypted, can be exploited to run code) matches the parallelism page, not quoted here.
- [ ] "We infer" sentences to confirm: (1) ECMP + HTTP/2 means every request on a connection lands on one server (from RFC 2991 s3/s4 and RFC 9113 s5/s9.1); (2) the cached conversation work is vLLM's prefix cache; (3) least connections and consistent-hash affinity conflict and the Lead does not rank them; (4) long idle timers keep dead-stream state longer; (5) smaller models leave more GPU memory for streams.
- [ ] RFC 9113 s9.1 says "HTTP/2 connections are persistent."; linked as 'are persistent'. RFC 2991 s4 paraphrase 'hashing the packet header fields that identify a flow' is from the Hash-Threshold / Modulo-N text.
- [ ] RFC 8622: "scavenge otherwise-unused resources only" cited to section 1 (also in the abstract); recommended codepoint '000001' in section 6; CS1 replacement in section 10; starvation quote in section 4. 'Prefer LE where your switches support it' is our advice from the RFC's SHOULD-transition sentence.
- [ ] RFC 2475 link has no anchor (abstract is unnumbered); our paraphrase 'classifies and marks packets for a forwarding behaviour at each hop' rests on the abstract sentence.
- [ ] nginx: proxy_buffering default on; proxy_read_timeout default 60s; hash ... consistent remaps 'only a few keys'; periodic health checks 'available as part of our commercial subscription' (in the #upstream example text); passive max_fails/fail_timeout under #server. Confirm the commercial-subscription sentence still reads that way.
- [ ] llm-d page (https://llm-d.ai/docs/operations/lifecycle/readiness-probes): anchors #health-endpoint, #v1models-endpoint-openai-compatible, #recommended-probe-configuration as reported; '30 minutes' is the YAML comment '30s * 60 = 30min max startup time'. The Lead mentions llm-d without a link.
- [ ] The Lead lists five proxy settings; the piece puts balancing in Practice one, three in Practice two, and health checks in Practice five. The Lead's ESI-LAG / vPC / LACP detail and its firewalling / Internet-access notes are left out.
- [ ] Diagram is redrawn from the Lead's text; the Lead's Figure 4 is mentioned in prose and linked via the article URL (no image link, no embed). Lead figures are captioned 'Image generated with AI' on the source page; none is reproduced.
- [ ] Internal links added: /apis-networking/http-connection-reuse/, /ai-ml-systems/kv-cache-offloading/, /ai-ml-systems/disaggregated-inference/ (all live, not draft). The commented-out related: line was left as the skeleton had it.
- [ ] Sources list gives RFC 9113 without a date and RFCs without author names (not verified in this session); dates for RFC 2991 (Nov 2000), RFC 2475 (Dec 1998), RFC 8622 (June 2019) came from the fetched pages.

**Fact-checker could not confirm:**
- [ ] Lead quotes were read through a fetch agent that caps quotes at about 125 characters and converts the page to Markdown. Every quoted Lead phrase came back as matching, but straight vs curly apostrophes and exact punctuation could not be seen in raw HTML (curl was denied). None of the quotes kept in the piece contains an apostrophe. The Author should eyeball them on the page.
- [ ] Title keeps the coined term 'Single-server inference' (the Lead names no pattern; the TL;DR and first section say it is our name). Author's call whether to retitle around the Lead's 'a few GPU nodes'. Also 'need no RDMA fabric' is firmer than the Lead's 'you might not need it', though the Lead does say 'this design carries no RDMA traffic'.
- [ ] nginx anchor #example and llm-d anchor #problem-statement: the fetch agents reported the sections and, for llm-d, the anchor id; the nginx #example id is the standard one for 'Example Configuration' on nginx.org but was not seen in raw HTML. Click both.
- [ ] vLLM docs links have no section anchors (only H1 ids were visible) and point at the 'latest' developer preview, which can change.
- [ ] GPU-memory arithmetic (384/640 GB, 243 to 499 GB, about 24 to 499 streams) is ours: the Lead names no GPU for its example and mentions 48 GB and 80 GB only as a comparison. Arithmetic checks out and is labelled 'by our arithmetic'; the Author may still prefer to cut it.
- [ ] Our one-line glosses for NVLink, PCIe, VRF, MTU and jumbo frames are not sourced definitions.
- [ ] 'A loading server must not receive requests' is our lead-in to the Lead's readiness quote; the Lead does not state it as a rule.
- [ ] llm-d '30 minutes' rests on the YAML comment in the Recommended Probe Configuration ('30s * 60 = 30min max startup time'); the same page's disaggregated decode pod example allows 60 minutes.
- [ ] The Lead gives the 10 GbE load time as both '113 seconds' and 'about two minutes'; the piece uses 113 seconds.
- [ ] Lead figures are captioned 'Image generated with AI'; Figure 4's image itself could not be inspected (caption confirms spine-and-leaf). Our Mermaid diagram matches the Lead's text.
- [ ] Frontmatter says picked 2026-11-26 but published/reviewed are 2026-11-25T07:37:22Z; assumed to be a Gap Day assigned by make daily, left untouched.
- [ ] RFC author names are not in the Sources list (confirmed on the pages: RFC 2991 Thaler and Hopps; RFC 9113 Thomson and Benfield, eds.); RFC 2991 spells it 'Equal-Cost Multipath', the piece follows the Lead's 'Equal-Cost Multi-Path'.

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

- Fixed: Hook: 'Reference designs for large AI clusters add a second network ... 400G ports tied to it' was not what the Lead says (it says large reference designs 'simply don't make sense' at this scale; the second network belongs to 'a training cluster'). Now: a training cluster adds a second, dedicated network for traffic between GPUs in different servers, and the 400G ports go with it.
- Fixed: Hook: near-verbatim unquoted Lead sentence on the 45-second load reworded; 'install or a version change' corrected to 'at install and on a version change'; 'one 141 GB model file per server' corrected to 'a 141 GB model load per server'.
- Fixed: TL;DR: backend fabric described as what 'large clusters build'; now 'a training cluster has', RDMA marked 'ideally' as in the Lead.
- Fixed: TL;DR: 'the backend fabric returns' is not on the page; now 'GPUs in different servers exchange data, the traffic a backend fabric is built for'.
- Fixed: Section 1: the clause 'so no GPU exchanges data with a GPU in another chassis' was verbatim Lead text outside the quotation marks; now inside the quote.
- Fixed: Traffic walk-through: 'sends a prompt over HTTPS and the server streams the answer back' corrected to the Lead's order (server responds over HTTPS; a chatbot's answer is streamed).
- Fixed: GPU-memory arithmetic: '243 to 499 GB remain: the Lead's "a few hundred GBs"' equated our numbers with the Lead's; now 'in line with the Lead's'.
- Fixed: Use cases: 'The Lead aims at deployments where...' softened to what the page says (this scale 'can be the right fit' where users are internal and the priority is latency or privacy rather than training throughput).
- Fixed: Staff chatbot example: restored the Lead's condition on the 8B to 32B model ('if it answers from an internal knowledge base'); sentence split for the 30-word limit.
- Fixed: H2 'Why small clusters end up with a fabric they do not use' asserted something the Lead does not; now 'Why a fabric built too early carries little traffic'.
- Fixed: Backend fabric section: 'Large AI cluster reference designs include one' (not on page) replaced with 'A training cluster has one'; 'buying the fabric at the start looks prudent' (the Lead says the opposite: building early is the risk) replaced with 'the question is whether to build the fabric now'; costs now framed as costs 'of building it too early'.
- Fixed: Lossless configuration: the Lead applies the phrase to PFC and ECN only, not queue management; sentence corrected.
- Fixed: QoS check: 'marking by DSCP or CoS' corrected to 'DSCP or CoS support' (the checks list asks for support; marking comes later).
- Fixed: Practice one: near-verbatim unquoted Lead clause ('a model starts late, goes quiet and then produces...') reworded; the 'Requests differ in cost...' quote now credited to the Lead.
- Fixed: Practice two: 'receives text in clumps' sat in a sentence linked to nginx; now credited to the Lead.
- Fixed: Practice two: affinity sentence was the Lead's sentence verbatim without quotation marks; paraphrased and credited.
- Fixed: Practice three: 'VXLAN overlay segment' corrected to the Lead's 'VXLAN Network Identifier (VNI)'.
- Fixed: Practice three: unsourced 'Isolation alone does not stop a model transfer from delaying API traffic' replaced with the Lead's point (without QoS, responses wait behind the transfer in the default queue).
- Fixed: Practice four: verbatim unquoted clause 'the spines see every packet as default class' reworded; shaper advice restored the Lead's qualifier ('where a shaper is available'); scavenger sentence now reflects the Lead's 'such as CS1 or ... (LE PHB)'.
- Fixed: Practice five: llm-d /health quote sits under 'Problem Statement', not '/health Endpoint'; anchor changed to #problem-statement.
- Fixed: Practice five: nginx commercial-subscription sentence is in the Example Configuration section, not under the upstream directive; anchor changed to #example.
- Fixed: Multi-node paragraph: vLLM says 'Efficient tensor parallelism requires fast internode communication, preferably through high-speed network adapters such as InfiniBand'; subject and wording corrected.
- Fixed: 'That isolated, fast network is the backend fabric' was an unflagged inference (the Lead's closing section names only InfiniBand-class adapters); now 'We infer...'.
- Fixed: Disaggregated prefilling: the Lead warns only of an unnamed 'experimental mode that passes cached state from one instance to another'; the match to this vLLM feature is now marked 'We infer'.
- Fixed: Closing: '3 to 5 years' quote refers to deployments scaling, not to 'that day' of the two trigger events; lead-in corrected.
- Fixed: Sources: RFC 9113 date added (June 2022).

</details>

### Not drafted

- none

<details><summary>Leads the editor rejected</summary>

- `25ckh34z`: GitHub multi-product degradation: page says only 'high database load', no trigger or mechanism; RCA promised but not published.
- `qcemeclm`: Copilot MAI Code 2 degradation: cause is only 'an issue with an upstream model provider', no mechanism.
- `eqcnxwjq`: Copilot MAI Code 2 availability: same vague upstream-provider wording, no mechanism.
- `5fyie67p`: Code Completions disruption: status page states no cause.
- `njpbtd3i`: Jira oauth-server degradation: page says a cause was identified but never states it.
- `ab27ett6`: Valkey 9.2.1: bug-fix point release, no CVE, one line of mechanism per fix.
- `5egzzcj7`: Kubernetes v1.37.3: patch release page with no change list; changelog not readable.
- `jdtwbsrh`: Kubernetes v1.36.7: patch release, same as v1.37.3.
- `l7pe3zb4`: Kubernetes v1.35.11: patch release, same as v1.37.3.
- `43ozefh5`: JDK 28 jlink security-properties plugin: small tooling feature heads-up, not an event; bundled-vs-platform-trust-stores is adjacent but does not teach baked-in runtime security config.
- `qqh4pqip`: Kubernetes v1.38 sneak peek: preview list of planned changes; the one item with mechanism (cgroup driver from CRI) has no live Explainer, and nothing on kube-proxy or Kueue.
- `ygg444bk`: Cloudflare Search for Agents: product announcement, request flow only, no internals or numbers; no live Explainer teaches it.
- `vtkcrm6t`: Terraform tfpolicy/SCIM: GA feature roundup, high-level, no mechanism.
- `pqdchiqi`: Distributed E-Maj 0.9.1: test-only add-on announcement, one paragraph of mechanism.
- `orkyy5ha`: vLLM Singapore meetup: event recap of five talks, not one idea or system.
- `r5xybey4`: ExecuTorch-TensorRT: rich mechanism, but no live Explainer for a Brief and AI/ML already has a pick today; good Case Study candidate for a later day.
- `6omdpdxb`: Inkling 2 Mini looped model on vLLM: rich mechanism (per-pass KV caches), held back to keep AI/ML to one Explainer today; strong Concept candidate later.
- `suwoz4in`: FlyDSL MonoKernels: rich mechanism, held back for Vertical balance (AI/ML is among the thickest).
- `lregvcvd`: DuckDB Star Trek queries: playful LLM-plus-SQL experiment, no architectural idea.
- `2uz33vjs`: SproutMe on DuckDB: customer story (DuckDB for CI, BigQuery in production), thin mechanism, authors say results do not generalise.
- `7dxjfjvi`: Fly.io Holding the Door for Agents: product announcement spanning several features (signup spec, 402 payments, MCP server, software factory), marketing tone, no single system in depth.
- `47zffks3`: Fly.io Petsem: already covered by the Brief fly-io-petsem-sqlite-lock-contention.
- `bucmcx76`: APNIC supply-chain risk essay: opinion, no mechanism.
- `45dzl3k6`: Cloudflare nginx2root: rich CVE mechanism, but security-advisory material and a second APIs & Networking pick today; RFC 10097 preferred.
- `a7gcdyix`: Meta technique-propagation agents: flow is clear but agent internals absent; AI/ML already has four Meta Case Studies.
- `b7cpii7i`: Meta Standard Model Templates: mostly framing around outcome numbers, little mechanism.
- `rn3af3oh`: Google Duplex in Gemini: historical overview and product tie-in, no numbers for the current stack.
- `gk6ifqc4`: PlanetScale Insights Agent: product announcement; good before/after query numbers but the agent's internals are thin.
- `35nxujtp`: Shopify secure mobile apps: concrete review pipeline but no detection numbers, and it would be a third Shopify security Case Study in Infrastructure & Cloud.
- `hyqudmlz`: Evolving the Java Language: 230-word abstract of a keynote video, no readable text.
- `ors6i6yf`: JEP 547 Lazy Constants: idea already covered by java-immutable-data.
- `lzoq2a6f`: United Airlines CDC to Iceberg: Source body could not be fetched (AWS page returned navigation only; feed gave a summary only), so mechanism could not be judged.
- `mmvdcjtg`: Aderant Iceberg pipeline: Source body not fetchable, summary only; cannot judge mechanism.
- `4tdw4ppe`: Aurora PostgreSQL excessive logging: Source body could not be fetched; cannot judge mechanism. Worth a retry on a later run.
- `xkqx2xhz`: Which AWS Database should you use: Source body could not be fetched; likely a product catalogue.
- `sdr5okwe`: Country Road Group integration on AWS: Source body could not be fetched.
- `low6vsvu`: ElastiCache distributed KV caching for LLMs: Source not fetchable, and the idea overlaps kv-cache-offloading.
- `25uryifl`: ElastiCache Serverless turns three: Source not fetchable; anniversary post.
- `mzvcbq2f`: PostgreSQL 19 released: release feature list, outside the Brief window.
- `hwl36z5p`: APNIC 2027 Activity Plan: organisational announcement.
- `ybd66hrc`: DNS topics at IETF 127: meeting roundup, not one event or idea.
- `iev2h2ju`: APNIC Honeynet podcast: podcast, no readable text.
- `dz3qoweq`: PyTorch Day China schedule: event announcement.
- `jxmllhov`: Google verifiable social reasoning: ML research, not an architecture event or pattern for the Reader.
- `2cmd5fpc`: Google transformer residual connections: model-architecture research, outside the backend Reader's systems scope.

</details>


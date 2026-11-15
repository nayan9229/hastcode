# Daily review queue, 2026-11-15

Drafted by `make daily` (ADR-0009). Every Piece is `draft: true` and unverified. Verify every claim against its Source, then set `verified:` and remove `draft: true` yourself (ADR-0007).

## Run 20261115-070549

### GitHub's October incidents: clients that never backed off

`content/briefs/github-availability-october-2026.md` · brief · 231 words · fact-checker: 24 claims checked, 8 fixed, 1 removed

**Publish within 72 hours of the Lead's date**, with or after its Explainer: metastable-failures

**Drafter asks you to check:**
- [ ] Publication date: the post is dated 12 November 2026 (by Jakub Oleksy). Per the 3-day rule that window closes 15 Nov; check the Lead is still fresh at publish time.
- [ ] Times: the incidents started 6 Oct 19:57 UTC, 19 Oct 15:49 UTC and 21 Oct 14:52 UTC. The Source's headers and duration lines disagree (Oct 19 header 1h7m vs ~1h28m; Oct 21 header 4h30m vs ~6h46m), so the Brief states start times only.
- [ ] Oct 6: 'a burst of very large queries' 'used all of the database's available connections'; 'The importer did not reduce its load when the database became overloaded' (paraphrased in the Brief).
- [ ] Oct 19: the Actions service lost connections to 'some cache servers' in one cloud region; 'Every server running that service then retried those connections at the same time', saturating the link between GitHub's datacenter and that region. The Brief says 'every instance retried at once' and 'shared datacenter link'.
- [ ] Oct 21: servers were near their memory limits; clients opened new connections immediately; overloaded servers were taken out of service and their traffic moved to others, which also overloaded.
- [ ] Fixes: 'We have limited how quickly this actions service can open new connections to its cache servers'; 'We added randomized delays to reconnects after a failed request'; 'We are automating importer backoff' (the Brief says 'still being automated').
- [ ] The only quote: "a broader reconnect storm", from the intro paragraph on 'Connection limits and other client improvements'. Confirm the exact wording.
- [ ] Inference, labelled: the link to metastable failures is ours; the report does not use the term.
- [ ] Our framing 'three partial faults that clients made worse' and 'load that never backed off spread the fault' is analysis. Check it is fair to Oct 6, which was a query burst rather than retries.
- [ ] Description is about 141 characters; tags are retries, backoff, incidents, github.

**Fact-checker could not confirm:**
- [ ] Fetched the Source through WebFetch's summarising model (Bash curl was denied), so the facts are confirmed only through its extracted fragments, not raw HTML. The quote 'a broader reconnect storm' came back verbatim in the full sentence, as did the author and date (Jakub Oleksy, 12 November 2026), the UTC start times and the fix sentences.
- [ ] Freshness: the Lead was posted 12 Nov 2026 and today is 15 Nov, the last day of the 3-day window. If it is published after 15 Nov the Brief fails the Lead-age rule.

<details><summary>What the fact-checker changed</summary>

- Fixed: Incident count: Source reports eight incidents in October, not 'three partial faults'. Now reads 'eight incidents, three involving load that did not back off'.
- Fixed: Oct 19: Source says the service 'lost connections to some of its cache servers', not that it lost the servers. Now 'lost connections to some cache servers'; 'every instance' changed to 'all its servers' (Source: 'Every server running that service').
- Fixed: Oct 19 link: 'shared datacenter link' changed to 'shared network link' (Source: 'the shared network link').
- Fixed: Oct 21: Source does not name clients. It says 'When requests to these servers failed, new connections were opened immediately'. Reworded to match.
- Fixed: Why it matters, Oct 21: 'healthy ones' is unsupported (Source: 'other servers... Those servers then became overloaded too'). Now 'shifted traffic overloaded the servers that received it'.
- Fixed: Analysis framing: 'load that never backed off spread the fault' was unlabelled and overstated Oct 6, which was a query burst, not a spread. Now 'In each case, we infer, load that did not back off made the fault worse'.
- Fixed: Description: 'never backed off spread three GitHub faults' changed to 'did not back off worsened three GitHub faults'.
- Fixed: Length: trimmed to keep the Brief within 250 words (an estimated 259 before the trim).
- Removed: 'between GitHub's data center and that region' (accurate but cut for the word limit)

</details>

### QuePaxa: consensus where the leader is only a fast path

`content/explainers/distributed-systems/quepaxa-consensus.md` · concept · 2804 words · fact-checker: 60 claims checked, 11 fixed, 0 removed

**Drafter asks you to check:**
- [ ] The QuePaxa paper PDF (bford.info/pub/os/quepaxa/quepaxa.pdf) could not be read (fetch timed out and no PDF text extraction was available). Paper claims are cited from its abstract page, Cloudflare's summary, the UCL project page and EPFL news. Check the phase descriptions and the '50% chance of deciding' per round against the paper itself.
- [ ] LAN throughput disagrees between sources: bford.info's abstract says 600k cmd/sec and the UCL project page says 584k. The draft cites only the WAN figures (250k cmd/sec, median under 380ms), which both agree on.
- [ ] The quotes from the Cloudflare post were pulled by a summarising fetch tool, not copied by hand. Confirm the exact wording of: 'The leader is completely optional.', 'At no point did the system halt.', 'a five-millisecond hedge timeout only increases the number of active proposers.', 'not been battle-tested as extensively', and the 'upwards of 95%' and 'around 65%' figures.
- [ ] The leader-failure chart: the draft uses the caption's 'a minute after failing'. The chart's own marks (15:00:32 failure, 15:01:25 replacement) give about 53 seconds.
- [ ] The Cloudflare post says a decision that misses phase 0 runs at 'at most a third of the usual rate', but the measured aggregate rate was about 65%. The draft cites only the 65% and the three round trips. The closing sentence's 'about a third of its decisions' is our rounding of the 35% drop.
- [ ] The 2020 incident: the 6h33m impact and the 80x dashboard slowdown came mostly from the knock-on database failover, not from etcd itself, which was read-only only until the switch recovered six minutes later. The hook now says this. Cloudflare's 2026 post links the incident to 'a consensus-powered system' degrading. Check that the framing is fair.
- [ ] The FLP description (1985, 'no deterministic protocol can guarantee consensus in an asynchronous system where even one process may crash') and its Sources entry (authors, JACM 1985) come from general knowledge. The linked Princeton PDF could not be read.
- [ ] 'Cloudflare's implementation is in Rust' rests on the post's mention of a planned write-up on 'Meerkat and QuePaxa's Rust implementation'.
- [ ] No heading anchors on blog.cloudflare.com were verified, so the Cloudflare links use the bare post URL. Add section anchors if the page offers them.
- [ ] The gloss of a multi-armed bandit ('a learner that keeps trying options and favours what performs best') and the line 'our view is that a well-understood Raft ... is still the safer choice' are our own words. Keep them or cut them.
- [ ] The EPFL news article's date (23 October 2023) and author (Tanya Petersen) come from the fetch summary; confirm both on the page.

**Fact-checker could not confirm:**
- [ ] FLP Journal of the ACM 1985 (vol. 32 no. 2): from search results only; the PDFs are scanned images, the ACM page returned 403 and the Northeastern copy 404
- [ ] FLP 'deterministic': dropped from the body because the abstract does not say it and the paper body could not be read
- [ ] QuePaxa paper PDF not read: the 50% per round figure and phase details rest on Cloudflare's account of the paper

<details><summary>What the fact-checker changed</summary>

- Fixed: Hook: 'opened on that incident' -> 'cited that incident' (Cloudflare post cites it; opening position not confirmed)
- Fixed: TL;DR: removed causal 'so' between not battle-tested and formal verification
- Fixed: 'not Byzantine faults' linked to GitHub (repo never mentions Byzantine) -> 'crash fault tolerant; it does not claim to tolerate Byzantine faults'
- Fixed: EPFL 'without interfering with each other' not in source -> 'a second leader can help in the same round without interfering with the first one'
- Fixed: EPFL 'stalls the system' not in source -> 'waiting a long time' + verbatim 'where the real disaster can happen.'
- Fixed: 'leader-failure test ran on servers spanning nearly all continents' -> 'Cloudflare's own evaluations of Meerkat used servers spanning nearly all continents'
- Fixed: 'passed SOSP's artifact evaluation' -> 'went through' (EPFL wording)
- Fixed: 'dial between latency and extra proposers' (not in post) -> 'trades extra bandwidth for how soon other proposers help'
- Fixed: FLP description reworded to the abstract's own wording; link switched from unreadable Princeton PDF to MIT DSpace record
- Fixed: FLP Sources entry -> MIT-LCS-TR-282, September 1982 (later in the Journal of the ACM)
- Fixed: 'FLP rules out a deterministic bound' -> 'we infer that FLP is why no fixed bound is promised'

</details>

### Fly.io 6PN: how lazy routing replaced a fleet-wide phonebook

`content/explainers/apis-networking/fly-io-6pn-routing.md` · case-study · 2400 words · fact-checker: 72 claims checked, 16 fixed, 3 removed

**Drafter asks you to check:**
- [ ] Anchors: the Fly.io post's heading ids could not be seen through the fetch. All deep links use the usual slug form (#the-address-is-the-route, #machines-learn-to-move, #a-phonebook-in-elvish, #the-network-already-asks, #ask-write-it-down-move-on, #what-changed, #safer-too, #what-it-costs, #changing-the-engine-at-cruising-altitude, #rerouting-on-purpose). Check that each one resolves. The same goes for docs.fly.io #6pn-addresses-in-detail.
- [ ] 100 to 200: one read of the post put this in the intro as 'records in total', another as 'routes on most hosts' (the 'What changed' section). The draft uses 'on most hosts' throughout. Check the intro wording.
- [ ] Stale-route forwarding: one read placed it in 'What changed', another in 'Rerouting, on purpose'. The draft links it to #what-changed.
- [ ] Date: the post shows no date. '13 November 2026' in Sources comes from the intake feed (leads.json). Confirm it or drop it.
- [ ] Byline: the post names no writer, only the illustrator Annie Ruygt, so no author is credited in Sources.
- [ ] drift-agent 'configures host peers itself': this wording came from the fetch summary. Check the exact phrase.
- [ ] Inference, labelled in the draft: the netlink event is VXLAN's l3miss flag (from the ip-link man page). The post does not name the flag.
- [ ] Inference, labelled: the lookup is local because Corrosion gossips a SQLite copy to the host. Corrosion's own post also mentions 'regionalization', a two-level design, so the copy may not be fleet-wide.
- [ ] Inference, labelled: UDP senders without TCP retransmission depend on their own retries. The post discusses only TCP.
- [ ] Inference, labelled: the feature flag made the rollout fall back to the old path by default.
- [ ] Quote 'Mistakes had, indeed, been made.' and 'what felt like a thousand state machines': check the exact wording.
- [ ] The quote '...the route is asymmetric' is lower-cased mid-sentence. The original reads 'The route is asymmetric for a while.'
- [ ] The last-check-against-cap quote ('within the same order of magnitude as the cap') and 'Every host had a different phonebook.': check them verbatim.
- [ ] 'Linux's default retransmission timeout is one second' is attributed to the post. The RFC 6298 link backs the general TCP value, not Linux's own default.
- [ ] The kernel setting Fly.io changed to accept asymmetric routes is deliberately left unnamed, because the post does not name it.

**Fact-checker could not confirm:**
- [ ] docs.fly.io #6pn-addresses-in-detail anchor: heading exists but id not visible; slug pattern implied by the page's #private-network-vpn link
- [ ] Quotes were checked through WebFetch's summarising model one question at a time, not against the raw page text

<details><summary>What the fact-checker changed</summary>

- Fixed: TL;DR and phonebook section: BPF map / Corrosion subscription link moved from #a-phonebook-in-elvish to #machines-learn-to-move
- Fixed: drift-agent 'configures host peers itself' -> 'handles host peers itself', linked to #ask-write-it-down-move-on; 'A miss for a 6PN address' -> 'Everything else'
- Fixed: s6pnforwarderd used to preload every migrated Machine (not every Machine)
- Fixed: dnat6pn 'broke some of Fly.io's own apps' -> 'even broke Fly.io's own apps'
- Fixed: 'at the last check' -> 'when Fly.io last looked'
- Fixed: Corrosion sync broke 'without saying so' -> problems Fly.io might not have heard about
- Fixed: VXLAN UDP-port / can't-shout claim now linked to the post #the-network-already-asks; added ARP
- Fixed: Netlink event fires for 'any unknown address VXLAN sees' (dropped unsourced 'not only for its VXLAN peers')
- Fixed: First-packet wait: replaced unsourced 'sender does not know a route has appeared' with post's 'with TCP the retry is up to the sender's kernel'
- Fixed: Asymmetric quote now verbatim: 'The route is asymmetric for a while.'
- Fixed: Removed unsourced 'during the trial' after yyz Flycast quote
- Fixed: 'gone from every host' -> 'old swapping code is gone'
- Fixed: 'phonebook could not do without more eBPF' -> 'two things Fly.io names'
- Fixed: 'a few hundred entries' -> '100 to 200 routes on most hosts'
- Fixed: Gaps paragraph: post page shows no date; blog feed dates it 13 November 2026
- Fixed: Sources: date attributed to the blog feed
- Removed: 'during the trial' (yyz Flycast trouble timing, unsourced)
- Removed: 'not only for its VXLAN peers' (unsourced)
- Removed: 'the sender does not know a route has appeared; it finds out when it retransmits' (unsourced)

</details>

### Not drafted

- none

<details><summary>Leads the editor rejected</summary>

- `p27jcgwv`: GitHub PR incident status page: says only 'mitigating action'; root cause deferred to a later RCA.
- `gl7yhd7jx00v`: GitHub Webhooks API status page: no root cause or mechanism stated.
- `oryjkch6`: Same as gl7yhd7jx00v: GitHub Webhooks API incident with no cause given.
- `ykawcojl`: GitHub email delivery: blamed on an unnamed upstream email provider, no mechanism.
- `b2t2dl4z`: Copilot GPT degradation: upstream model provider issue, no cause given; no live Explainer teaches it.
- `hxhbrezm`: Jira incoming email: caused by a Microsoft 365 degradation; status page gives no mechanism.
- `bbrn7vi5`: Inside Java Newscast: video with no transcript; adaptive heap sizing is already live as zgc-adaptive-heap-sizing.
- `hxstaai6`: Kitesurf: open-source announcement with only a high-level architecture overview and no numbers; no live Explainer teaches it.
- `l5laiihp`: Datasentinel release: vendor product announcement.
- `2mbhwtqc`: Cloudflare Spend Management: product announcement.
- `ui6vizld`: Node.js alpha release: release list, no design change.
- `mvap46l3`: Aurora MySQL error 63955: troubleshooting how-to, not an event, and no Explainer teaches it.
- `wtamray2`: Durable execution (Fowler): strong Concept source, but QuePaxa already fills today's distributed-systems slot. Good backlog pick.
- `cmgchcqh`: Shopify Cyber Chaos Monkey: an AI security program, not a Brief event, and not chosen over the Fly.io and QuePaxa Explainers today.
- `ucr5ttmb`: RDS Proxy multiplexing Part 1: AWS how-to on a well-covered Vertical (databases-storage); possible later Concept on connection pinning.
- `u6yodzgr`: RDS Proxy multiplexing Part 2: diagnose and monitor how-to; same reasons as Part 1.
- `d3l7vjmh`: Fastly election night: traffic observations, no system mechanism, and no Explainer teaches it.

</details>


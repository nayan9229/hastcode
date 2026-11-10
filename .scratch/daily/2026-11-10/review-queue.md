# Daily review queue, 2026-11-10

Drafted by `make daily` (ADR-0009). Every Piece is `draft: true` and unverified. Verify every claim against its Source, then set `verified:` and remove `draft: true` yourself (ADR-0007).

## Run 20261110-063041

### GitHub's anonymous-traffic surge: shed the cheapest class first

`content/briefs/github-anonymous-traffic-surge.md` · brief · 239 words · fact-checker: 22 claims checked, 2 fixed, 0 removed

**Publish within 72 hours of the Lead's date**, with or after its Explainer: metastable-failures

**Drafter asks you to check:**
- [ ] 05:59 UTC 9 Nov: Investigating 'impacted performance for some GitHub services'
- [ ] 06:30 UTC: 'debugging a spike in anonymous traffic'; logged-out users 'may see errors'; CI jobs and release downloads 'may fail with a 504 Gateway Timeout status if they make anonymous calls'
- [ ] 06:51 UTC: Monitoring, 'The degradation has been mitigated'
- [ ] 07:43 UTC: Investigating, 'repeated increase in 504 errors for anonymous and unauthenticated traffic, including users who are not logged in and release downloads'
- [ ] 08:16 UTC quote (the Brief's only quote): 'Logged-in users and authenticated requests are unaffected.'
- [ ] 08:47 UTC: 'Error rates have declined again, and we're preparing targeted traffic controls while monitoring for a recurrence.'
- [ ] 09:09 UTC: Resolved; 'prepared additional protections'; RCA to be shared later
- [ ] 8 Nov incident pyynvrm1vskw: 10:36-12:32 UTC, component API Requests, 'degraded performance primarily for unauthenticated users', 'identified a problematic component and are working to redeploy it'
- [ ] Inference: anonymous and authenticated classes 'already kept apart' (page does not say how)
- [ ] Inference: targeted traffic controls act as load shedding
- [ ] Inference framing: metastable failures 'offers one reading' of the recurrence; not asserting GitHub's incident was metastable
- [ ] Claim: neither status page connects the 8 and 9 Nov incidents
- [ ] Brief's 'under an hour' = 06:51 mitigation to 07:43 recurrence (52 min)
- [ ] Quotes gathered via WebFetch summariser, not raw HTML (Bash denied); verify wording on the page

**Fact-checker could not confirm:**
- [ ] Page wording checked via WebFetch summariser (two independent passes, consistent), not raw HTML; Author should eyeball the 08:16 quote 'Logged-in users and authenticated requests are unaffected.' on the live page before setting verified:

<details><summary>What the fact-checker changed</summary>

- Fixed: 8 Nov sentence: 'until GitHub redeployed a component it called problematic' -> 'and GitHub said it was working to redeploy a problematic component' (page says 'working to redeploy', never says redeploy completed or caused mitigation)
- Fixed: 08:47 sentence: 'GitHub prepared targeted traffic controls' -> 'GitHub was preparing targeted traffic controls' (page says 'we're preparing targeted traffic controls')

</details>

### Topstep's Redis to Valkey cutover: the DNS TTL is the write window

`content/briefs/topstep-redis-to-valkey-migration.md` · brief · 240 words · fact-checker: 28 claims checked, 6 fixed, 0 removed

**Publish within 72 hours of the Lead's date**, with or after its Explainer: planned-database-switchover

**Drafter asks you to check:**
- [ ] Date: feed pubDate 'Mon, 09 Nov 2026 16:59:06 +0000'; Brief says 'On 9 November'.
- [ ] 'about 53 million keys': Source says 'roughly 53 million keys'.
- [ ] Source system: single self-managed Redis OSS 7.4 on one EC2 r5.4xlarge, no replica, no automatic failover; target Multi-AZ ElastiCache for Valkey 7.2.6.
- [ ] Quote 'falls outside the Redis OSS ≤ 7.2 range that the managed migration supports' is a fragment of the sentence beginning 'their post-fork Redis OSS 7.4 source ...'; check the '≤' glyph as the post renders it.
- [ ] AWS Migration-Prepare doc (docs.aws.amazon.com/AmazonElastiCache/latest/dg/Migration-Prepare.html) says you can migrate 'from any supported Valkey or Redis OSS version'; Brief attributes the 7.2 limit to Topstep ('Topstep says'), not to AWS.
- [ ] RedisShake 'in PSYNC mode' ran as a Kubernetes deployment on EKS 'until replication lag converged to near-zero'; RedisShake sync-reader doc says it 'simulates a Slave connecting to the Master node'.
- [ ] Route 53: private hosted zone A record (EC2 IP) changed to CNAME at the ElastiCache primary endpoint, 30-second TTL; 'services picked up the new target within a few seconds'.
- [ ] 'old node stayed up as a warm standby': Source: 'The source Redis OSS instance stayed running as a warm standby during a soak period.' RSS feed text says kept at least 24–48 hours; not used in Brief.
- [ ] Omitted for length: pre-cutover key-count parity check and application smoke tests against green; VPC peering between accounts; TLS allowlist via tls.checkServerIdentity; latency results (read P99 ~1.25 ms, write P99 ~0.5 ms).
- [ ] Inference 1 (labelled 'we infer'): some clients could still write to the old primary for up to 30 seconds; Source does not discuss writes during the DNS switch.
- [ ] Inference 2 (our analysis, unlabelled): 'stragglers follow only while replication runs, then they are lost'; Source does not say whether RedisShake kept running through or after cutover. Consider adding 'we infer'.
- [ ] Route 53 TTL doc: 'The amount of time, in seconds, that you want DNS recursive resolvers to cache information about this record.'
- [ ] Comparison link: 'calls that window unavoidable' rests on the Comparison's opening ('You cannot remove that window') quoting PlanetScale's 'unavoidable cost'.
- [ ] 'hostname in front of the cache from day one' paraphrases Source lesson 'Place DNS in front of your cache endpoint from the start'; 'warm rollback target left running' paraphrases lesson 2 (validate green, keep blue as rollback target).
- [ ] sources: frontmatter added with RedisShake and Route 53 docs; build accepted the key.
- [ ] Word count 240 is a hand count (Bash denied for wc); build's brief-length check passed.

**Fact-checker could not confirm:**
- [ ] Raw HTML rendering of the '≤' glyph in the quoted sentence: curl was denied; WebFetch rendered it as '≤' consistently on two fetches, kept as is
- [ ] Word count 241 is a hand count using the build's rule (whitespace split, headings and AI note included); build's brief-length check passed

<details><summary>What the fact-checker changed</summary>

- Fixed: Attribution: 'Topstep says its source' -> 'The post says Topstep's source' (sentence is in the post authors' third-person voice)
- Fixed: Added 'we infer' to 'stragglers follow only while replication runs, then they are lost' (post does not say whether RedisShake kept running after cutover)
- Fixed: 'a hostname in front of the cache from day one' -> 'a hostname already in front of the cache' (post says the A record existed before migration; lesson says 'from the start', not day one)
- Fixed: Description: 'moved about 53 million keys ... with one Route 53 change' -> 'cut about 53 million keys over ... with one Route 53 change' (RedisShake moved the data; the Route 53 change was the cutover); 147 chars
- Fixed: Added Source link on 'until lag was near zero' (claim previously linked only to the RedisShake doc)
- Fixed: Added Source link on 'Two choices made the cutover a record edit' (previously unlinked paraphrase of the post's lessons)

</details>

### Not drafted

- none

<details><summary>Leads the editor rejected</summary>

- `j6yor3hg`: Fetched: Copilot model degradation attributed only to 'an issue with an upstream model provider'; no mechanism, RCA pending.
- `jqqos7pq`: Fetched: Actions and Checks API incident gives no cause; RCA pending.
- `qdtwws4n`: Fetched: Webhooks and Actions delays, mitigation applied, no cause stated.
- `a5e543zb`: Fetched: single Discord update citing a 'networking issue' with no cause or mechanism.
- `bonucn3e`: Fetched: 8 Nov unauthenticated-user timeouts fixed by redeploying 'a problematic component'; no cause named. Folded into the 9 Nov Brief as context rather than picked alone.
- `jknpqkbf`: Fetched: npm publish issue, three lines, no cause; also only ~8 hours left in the 72-hour window.
- `auxsxxjo`: Fetched: DeepL Cloud WAN post is a repeatable drain procedure (AS_PATH prepends) with no dated event, and no live Explainer teaches closest-exit routing. Case Study candidate for a later day; Explainer cap is 0 today.
- `ldt6gffh`: Fetched: Valhalla Ticker is a tutorial on JDK 28 preview value classes, not an event at a named system.
- `rh5e2yi2`: Netflix synopsis fact-checking post describes a system, not an event; a Case Study candidate for a later day, Explainer cap is 0 today.
- `umoswrxq`: AWS customer story about agentic fleet analytics; no event and no live Explainer teaches it.
- `jkzd7bu7`: Google Research product post (Language Explorer); not an architectural event.
- `azacbbvq`: Funding announcement (Vite Ecosystem Fund).
- `ro3i2ch6`: APNIC opinion piece on post-quantum DNSSEC; no event. Possible Concept source later.
- `2d4mdp4p`: Fastly marketing post on green websites.
- `nafhp7cv`: Dasha 1.10 point release of a PostgreSQL monitoring tool; feature list, no design change at a named system.
- `ol3gy5ix`: Devoxx talk recording; no readable first-party text and not an event.

</details>


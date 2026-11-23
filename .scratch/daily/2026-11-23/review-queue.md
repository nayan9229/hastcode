# Daily review queue, 2026-11-23

Drafted by `make daily` (ADR-0009). Every Piece is `draft: true` and unverified. First tick each draft's self-edit boxes while you read it, before you verify ([the checklist](../../../docs/editorial/self-edit-checklist.md)): they check the writing, not the facts, and the self-edit does not replace verification. Then verify every claim against its Source, then `make release` (ADR-0012): it sets `verified:`, removes `draft: true` and keeps each Publication Day. Ticked self-edit boxes never stand in for `verified:`.

## Run 20261123-092320

### Zero Operator Access: how Bee keeps operators out of user data

`content/explainers/infrastructure-cloud/bee-zero-operator-access.md` · case-study · 2400 words · fact-checker: 130 claims checked, 16 fixed, 6 removed

**Drafter asks you to check:**
- [ ] AWS post read via its RSS feed (content:encoded) through a fetch tool, not the rendered page: verify every quote against the live page; curly apostrophes are typed straight.
- [ ] AWS post byline (Gopinath Jagadesan, Steve Korshakov, Ethan Sutin, Arvind Raghu) and date 20 November 2026 come from the feed's About-the-authors section and pubDate.
- [ ] Whitepaper date: the amazon.science landing page shows only 2026; the PDF filename suggests 09-23-26. Sources list says 2026.
- [ ] AWS post contradicts itself on Lambda: 'never see plaintext' (network section) vs plaintext with user consent (runtime step 4). The draft cites both and follows the whitepaper (p8, p22): consent-gated, consent max TTL 7 days.
- [ ] KMS-held key option: neither Bee source discusses it. The AWS doc route (#attestation-kms) and the 'deciding argument' are labelled as our inference.
- [ ] PCR4/PCR7/PCR12 meanings, 'signed by the Nitro Hypervisor' (nitrotpm-attestation-document-content) and the attestable-AMI definition phrase came from a fetch summary, not a direct read: check the wording.
- [ ] Rekor 'Auditors can monitor the log for consistency' and anchor #auditing-the-public-instance, and SCP anchor #not-restricted-by-scp for the management account, came from a fetch summary: check.
- [ ] Hook's 'a public log that Amazon does not operate' rests on whitepaper p11 (public-good Rekor operated as part of Sigstore/OpenSSF rather than by Amazon). The transparency service itself is a separate Amazon team.
- [ ] Trail of Bits 2026 assessment cited to whitepaper p24 only; the report itself was not read.
- [ ] Bedrock Mantle claim cited to the AWS post only; the Mantle post did not load.
- [ ] Egress: the AWS post lists only VPC endpoints to RDS, ElastiCache and Lambda; whitepaper p8-9 adds DNS, S3 artifact access and peered-VPC traffic via security groups. The draft gives both.
- [ ] Inference servers: the AWS post says they 'never receive keys'; the whitepaper says no access to raw customer keys. The diagram says 'no keys'.
- [ ] TL;DR 'up to 7 days after the user's last request' paraphrases whitepaper p17 (each client request resets an absolute expiry, max 7 days); the AWS post says expiry after 7 days without user activity.
- [ ] The 18:00-Friday key-expiry worked example is ours, labelled as such.
- [ ] Private CA account: whitepaper p22 says it is accessible only by the Bee team; the draft says 'a tightly controlled account'.
- [ ] Tags: the company tag is 'bee', not 'aws', following the netflix/meta/shopify convention.

**Fact-checker could not confirm:**
- [ ] Whitepaper #page=N anchors mapped from inline footer page tokens in an r.jina.ai text extraction, not a PDF viewer (25 printed pages = 25 PDF pages, all cited pages consistent); spot-check in a viewer
- [ ] AWS post read only via RSS feed content (rendered page returned navigation only); byline and 20 Nov 2026 date from feed
- [ ] Trail of Bits 2026 assessment cited to whitepaper p24 only; the report itself not read
- [ ] Inference paragraph's 'A key policy governs every KMS key' is an unlinked general KMS fact
- [ ] Quotes use straight apostrophes where AWS post uses curly ones

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

- Fixed: TL;DR: 'process plaintext on servers' was linked to p1, which does not say it; now links p4 ('must process plaintext'), with compute claim kept on p1 ('more compute than a personal device can provide')
- Fixed: TL;DR: PCR14 holds one hash over every container image and model file, not per-file hashes
- Fixed: Step 1: PCR meanings corrected to AWS doc: PCR4 hashes of binaries the UEFI environment runs, PCR7 a hash of the Secure Boot policy, PCR12 a hash of the boot command line
- Fixed: Step 4: X.509 extension carries the 'Sigstore bundle', not 'Rekor bundle' (whitepaper p12)
- Fixed: Options: 'AWS never holds Bee's key' -> 'the user's key'
- Fixed: Step 2: 'The manifest also replaces EC2 user data' -> Bee replaces EC2's built-in user data loading; 'Its manifest loader rejects legacy bash userdata' (p7)
- Fixed: Step 3: 'Each entry is the PCR map as canonical JSON' -> entry's payload is canonical JSON holding the PCR map; split 31-word sentence
- Fixed: Step 3: added p11 link for 'an S3 bucket replicated to each region' (was unlinked)
- Fixed: Second cost: 'compromised phones' -> 'compromised client devices' (p3)
- Fixed: First cost: 'A new container image, model file or security agent' -> 'Changing a container image, adding a model file or modifying the security agent' (p23)
- Fixed: First cost: 'Diagnostics run only through' -> 'must run through' (p23)
- Fixed: First cost: 'Review enforces that, not encryption' -> 'Review and deployment discipline enforce that, not encryption alone' (p23)
- Fixed: Options: 'discard session state after each request' -> 'once a request completes' (p5)
- Fixed: TL;DR: 'not AWS' -> 'not AWS's operators' (post: AWS infrastructure operators)
- Fixed: Second cost: Client Secret glossed as 'the secret held on the user's device', not 'the root secret' (Account Secret/M is the key-tree root)
- Fixed: Evidence: Lambda consent wording -> 'only when an active cryptographic consent record authorizes it' (p8)
- Removed: 'Bee set out to apply that rule to itself: an engineer with production access should find nothing to read.' (unsourced framing; size)
- Removed: 'The image also leaves an operator nothing to log in to.' (redundant; size)
- Removed: 'A measured instance with an open network path could still leak, so Bee closes the paths first.' (unsourced framing; size)
- Removed: 'Here is the lifetime with numbers.' (size)
- Removed: TL;DR p5 link 'in the background while no session is open' (size; body keeps p5 link)
- Removed: Diagram intro step-grouping paragraph shortened (size)

</details>

### Devin on macOS: how Cognition runs agent VMs on EC2 Mac

`content/explainers/infrastructure-cloud/cognition-devin-macos-vms.md` · case-study · 2300 words · fact-checker: 95 claims checked, 12 fixed, 0 removed

**Drafter asks you to check:**
- [ ] Lead text was read through the AWS Architecture Blog RSS feed (content:encoded) because the page body did not load for the fetcher. Check every quote against the live page: "the two systems raced to configure the same filter.", "a limit macOS enforces.", "each dialog would block Devin until a person approved it.", and the separation-of-lifetimes paraphrase.
- [ ] The AWS heading anchors used (#why-an-ai-coding-agent-needs-its-own-mac, #why-cognition-chose-aws, #solution-overview, #keep-session-state-with-disk-snapshots, #apply-network-access-rules-with-a-user-space-ethernet-gateway, #prepare-guest-images-before-the-first-boot, #let-devin-operate-native-apps-through-the-accessibility-tree, #results, #plan-your-own-macos-fleet-on-ec2-mac) came from the feed HTML ids. Confirm they exist on the live page. The devin.ai anchors (#why-devin-needs-a-mac, #disk-snapshots, #enter-ethernet, #bootable-is-not-ready) came from that page's table of contents.
- [ ] Lead date 20 November 2026 is the feed pubDate (01:10 UTC); the date shown on the page was not seen. Seven authors are listed under About the authors, but the feed creator is Adam Hesch only. The draft says three work at Cognition (Hesch, Olut, Almloff) and four at AWS (Almela, Quang, Krall, Mitine, the Worldwide GTM Lead for EC2 Mac).
- [ ] Intro facts linked to the Lead with no anchor: nearly 1,000 organisations within six weeks 'according to Cognition'; a fresh VM in under a minute across six AWS Regions; '15–20 times faster than setting up a Mac host themselves', by Cognition's estimate. The draft writes 'organisations' (British) outside quotes.
- [ ] Customer quotes: Evan Almloff, 'mac-m4pro.metal and mac-m4.metal fleet spans six regions' and 'in under a minute'. Sumit Oberoi (Wispr), 'cut down code review time by 70%' and 'improved our velocity by nearly 2x'. These are quoted as clauses; check that the clause boundaries do not change their meaning.
- [ ] The two-VMs-per-Mac licence claim is cited to the AWS post. apple.com/legal/sla did not render, so the macOS SLA clause itself was not checked.
- [ ] 'Available since macOS 14' for VZNetworkBlockDeviceStorageDeviceAttachment comes from Apple's documentation JSON (introducedAt 14.0). The file-handle gloss 'transmits raw network packets and frames over a datagram socket' paraphrases Apple's abstract.
- [ ] Glosses that go beyond the Source's words: vhost-user 'lets a user-space process on the same host serve a VM's device queues' (QEMU says it implements 'the control plane needed to establish virtqueue sharing with a user space process on the same host'). EC2 Dedicated Host is glossed as 'a physical server allocated to one AWS account'; the Dedicated Hosts page itself was not fetched. The Virtualization framework root page was not fetched.
- [ ] 'The unit of billing is the host' and 'only as On-Demand Instances' come from the EC2 Mac user guide (no anchor found). The 24-hour minimum and one-instance-per-host claims link #mac-instance-considerations.
- [ ] Cognition quotes: 'can be unreliable and expensive.' (UI automation, #bootable-is-not-ready); 'Stopping both processes together could interrupt the very work needed to preserve the session.' (#disk-snapshots); 'Each datagram carries one frame,' and 'The trade-off is that we now own the link-layer handling and validation that managed networking would otherwise provide.' (#enter-ethernet).
- [ ] The AWS post says sessions free and reclaim their VM 'multiple times' on average (an earlier fetch said 'several'). The draft uses 'multiple'. The 90% figure has no time window, and the draft says so.
- [ ] Labelled inferences: reuse as the deciding argument; the prepared image keeps start under a minute; the source-IP check keeps sessions apart; a host allocated for a short burst costs a full day; release-on-idle is the one density lever left; images need rechecking when macOS or Xcode changes. The worked example (48 VM-hours, three 20-minute bursts across eight hours) and the advice to measure idle share are ours and labelled as ours.
- [ ] 'Neither post says why bridged networking lost' is stated as a gap; there is no speculation about the reason.
- [ ] RFC 826 author and date (David C. Plummer, November 1982) were added from memory, not fetched.
- [ ] The feed has garbles not carried into the draft: '(SLA) (SLA)', '(Amazon S3) (Amazon S3)', a repeated sentence 'Each Mac VM has its own gateway and its own session policy enforcement', and an unclear 'was not supported until a couple months ago' (not used).
- [ ] Frontmatter: related: [fly-io-sbd] is now uncommented (the slug exists and is live). The comment was rewritten so no TODO(author) marker remains. draft: true is kept, and verified: was not added.
- [ ] Build: re-run check-page.sh once bee-zero-operator-access fits the size limit. The final page bytes have not been measured. The scratch build directory .intake/daily/20261123-092320/build-cognition-probe can be deleted.

**Fact-checker could not confirm:**
- [ ] The body of the AWS post's HTML page never loaded; only navigation came through. Every AWS claim and quote was checked against the post's content:encoded in the AWS Architecture Blog RSS feed. All four AWS quotes are verbatim there ('the two systems raced to configure the same filter.', 'a limit macOS enforces.', 'each dialog would block Devin until a person approved it.', and the Almloff and Oberoi clauses). They should still be checked on the live page.
- [ ] The AWS heading anchors come from the feed HTML ids, which match the ones the draft uses. They were not confirmed on the live page.
- [ ] The date shown on the AWS page was not seen. The Sources entry's 20 November 2026 is the feed pubDate (01:10 UTC), and the page may show 19 November in US time.
- [ ] The devin.ai anchors (#why-devin-needs-a-mac, #disk-snapshots, #enter-ethernet, #bootable-is-not-ready) appear as links in the page's table of contents. The id attributes on the headings themselves were not seen.
- [ ] The QEMU #introduction anchor is linked from the page's table of contents, but the heading's own id may be #id4.
- [ ] The macOS Software License Agreement itself (apple.com/legal/sla) was not fetched. The 'two additional copies' claim is cited to, and supported by, the AWS post.
- [ ] Absence claims ('Neither post names the databases', 'Neither post gives a price', no host count, utilisation, upload time or failure figures, no time window for the 90%) rest on full reads of the feed item and of the devin.ai page. A negative cannot be confirmed more firmly than that.

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

- Fixed: TL;DR: the snapshot sentence read as if the upload to S3 runs over NBD. The AWS post says NBD carries only the guest's disk I/O and names no transport for the S3 upload. Reworded to 'A snapshot process, reached over NBD, saves the disk to Amazon S3'.
- Fixed: TL;DR: 'An Auto Scaling group supplies hosts' now adds 'where none are preallocated'. The AWS post limits the Auto Scaling group to Regions without preallocated capacity.
- Fixed: Dedicated Host gloss: 'a physical server allocated to one AWS account' became 'a physical server dedicated to one customer's use'. The AWS Dedicated Hosts page says 'fully dedicated for your use', and the AWS post says hosts can be shared across accounts.
- Fixed: Problem section: 'customers fell back on their own CI or Mac fleet' became 'To build and test, it depended on the customer's own CI system or Mac fleet.' The AWS post says Devin depended on the customer's CI or Mac fleet. The sentence was also split to clear a sentence-length warning.
- Fixed: Options section: 'name an alternative at four points, and each lost for a stated reason' became 'give a reason for turning down four alternatives'. Both posts also name bridged networking and give no reason for passing it over.
- Fixed: Racks bullet: 'EC2 Mac hosts come through the EC2 API, run in an Amazon VPC' now says 'instances'. The AWS post says Mac instances run in a VPC and boot from EBS.
- Fixed: vhost-user gloss: 'serve a VM's device queues' became 'sharing a VM's virtqueues with a user-space process on the same host'. This keeps to QEMU's own wording.
- Fixed: Step 2: 'so system setup does not compete with Devin' was attached to turning off background daemons. Cognition's post states that reason for the boot pass ('We don't want Devin to compete with system initialization for resources') and says the daemons are disabled because they 'aren't needed'. Moved the reason and its link to the boot pass.
- Fixed: Step 2: 'warms up the iOS Simulator' became 'warms the Simulator', Cognition's wording.
- Fixed: Step 3: the 90% figure is now attributed ('Cognition reports that about 90%'), as the AWS post does.
- Fixed: Costs: 'Cognition edits macOS internal configuration databases directly' was unlinked. It now links the AWS post's #prepare-guest-images-before-the-first-boot section.
- Fixed: Closing section: 'Customers no longer need their own CI system or Macs' became 'It no longer needs a customer's CI system or Macs to build and test Apple apps'. The AWS post makes Devin the subject, not customers.

</details>

### Data contracts: an API spec for a dataset, tested in CI

`content/explainers/distributed-systems/data-contracts.md` · concept · 2300 words · fact-checker: 110 claims checked, 12 fixed, 4 removed

**Fills Gap Day 2026-11-20** (ADR-0012): `published` is that day; `make release` keeps it.

**Drafter asks you to check:**
- [ ] Every Lead anchor comes from the page's table-of-contents links, not the DOM. Confirm the odd ones resolve: #DataContractsAreThex201cproductDetailsPagex201dInADataMarketplace, #RolesampSecurity, #TheProblemWithx201cselectx201d, #Data-firstVs.Contract-first, #DataContractsAndTheDataMeshPrinciples, #TestingDatasetsVsTestingCode.
- [ ] The Lead was read only through a fetch agent, because direct curl was denied. Every quote was confirmed character for character against that agent's copy of the page, but not by me in a browser.
- [ ] Freshness vs latency: the Lead's example uses the SLA property `freshness` (24 hours on last_updated_at). The ODCS SLA page lists `latency` as "preferred to freshness", and the draft says both. Check the anchor #valid-values-for-sla-properties.
- [ ] 'Over 1 million downloads per month' is the Lead's figure. The CLI's own badge showed 1.4M/month on 2026-11-23. The draft discloses that Harrer co-founded Entropy Data, the copyright holder on the CLI's MIT licence, and calls it 'its makers' count'. Check that the framing is fair.
- [ ] PayPal's 2023 open-sourcing date comes from the Lead, not PayPal's README. The README only lists a third-party 2023-05-01 article.
- [ ] 'ODCS current version v3.2.0' comes from the GitHub README. The changelog lists v3.2.0 as APPROVED with no date, and v3.1.0 as 2025-12-08. The ODCS 3.1 (2025) support for the Data Contract Specification is from the Lead.
- [ ] Bitol as 'a Graduated project of LF AI & Data under the Apache 2.0 licence' is from bitol.io. The Lead says 'Bitol project at the Linux Foundation'.
- [ ] The removal of `password` (v3.0.0, 'credentials should not be stored in the data contract') links the changelog page without an anchor.
- [ ] `datacontract breaking` exiting with status 1 is per ERROR severity on cli.datacontract.com/compare-contract-versions#severity-levels. The rule to fail a pull request when a breaking change ships as a minor or patch version is the authors' (#AutomateAllTheThings), not a CLI default. The diagram edge 'breaking change, version not major' follows the authors' rule.
- [ ] OAAS RFC 0056 shows no status or date, and its Decision section is the template placeholder. The draft says 'The RFC records no decision yet'. The fields (purpose, dataContractId, status requested/approved/rejected, endTimestamp) come from its Fields table.
- [ ] These definitions are our own wording, not quotes: data mesh (paraphrasing the Lead's sidebar), data product, output port ('the interface a consumer reads from'), schema registry, and forward compatibility ('consumers on the old schema can read the new data'). The Lead says only that a minor version 'has to be forward compatible' and should 'ideally' be fully compatible.
- [ ] These lines are marked 'we infer': what does not count (a wiki page, an Iceberg table's own schema), that the provider feeds both output ports during a major migration, and that the split lets grants come and go without a new contract version.
- [ ] The CLI's test targets (Snowflake, Databricks, BigQuery, Kafka, Postgres, S3, and local Parquet/JSON/CSV/Delta) come from the CLI homepage FAQ, linked as #frequently-asked-questions; that anchor is unverified. The CLI's own counts of supported sources disagree (15+ vs 18+), so the draft gives no number.
- [ ] The Radar facts are Vol 35, November 2026, ring Adopt. The quote 'a sensible default for teams building data products' is from the blip text.
- [ ] The `datacontract import sql --source ddl.sql --dialect postgres --output datacontract.yaml` example is copied from the CLI commands index (MIT). Check that it is still verbatim.
- [ ] The Bertrand Meyer PDF title is unverified. The Sources line calls it 'the chapter the Lead links'.
- [ ] Alex is the Lead's illustrative persona ('a data engineer at a multinational insurance company'). The draft calls Alex 'the authors' example engineer'.
- [ ] The word count (~2,300) is an estimate. A shell word count was not allowed in this session.

**Fact-checker could not confirm:**
- [ ] Bertrand Meyer PDF title and author: the URL serves a PDF (HTTP 200) but its text could not be read; pdftoppm is not installed
- [ ] Quotes were checked by fetch agents against converted text, not raw HTML or a browser (curl denied); long sentences, such as the Alex persona sentence and the forward-compatibility sentence, were confirmed by fragment only (125-character cap)
- [ ] ODCS changelog: no v3.0.0 anchor found; the link stays without an anchor
- [ ] SLA anchor #valid-values-for-sla-properties: confirmed from the page's table-of-contents link; the raw id attribute was not seen
- [ ] 'Over 1 million downloads per month' is the Lead's figure only, not counted independently; the 'by its makers' count' framing is kept (Harrer is co-founder and CEO of Entropy Data GmbH, the copyright holder on the CLI's licence)
- [ ] Word count not measured

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

- Fixed: Schema: customer_name tag is the string `pii:true`, not `pii: true` (also fixed in 'What data contracts make possible')
- Fixed: Quality: row rule is an SQL count with `mustBeGreaterThan: 100000`, not `rowCount`; status rule is on `customer_status` (sentence split for the 30-word limit)
- Fixed: Origin: ODCS 3.1 (2025) brought the two formats together and Data Contract CLI gained native ODCS 3.1 support; the Data Contract Specification is 'another format', then retired (the draft said ODCS 3.1 added native DCS support)
- Fixed: Bitol: 'moved to Bitol' is now 'became ODCS, the founding standard of the Bitol project', linked to the Lead's #TheOriginOfTheOpenDataContractStandard
- Fixed: ODPS: 'builds on ODCS' and 'output port described by a data contract' now link to the Lead's #DataProductsBuildUponDataContracts (the ODPS README says neither)
- Fixed: Roles: replaced 'such as a Snowflake role with read access' (not on the spec's Roles page) with the Lead's 'an IAM role on the data platform, such as a Snowflake role' (#RolesampSecurity)
- Fixed: Incentives: unsupported 'payoff lands with the teams that read' is now the Lead's quote 'friction with little payoff' (#TheCost)
- Fixed: Data-first: import formats are now 'such as' (the Lead's list is not complete)
- Fixed: Contract-first: added 'among other tools' (the Lead also lists a YAML editor and an Excel template)
- Fixed: Access agreement: 'three parts' is now 'three main parts', per the Lead
- Fixed: Author bios sentence now links the Lead
- Fixed: Sources: 'TolerantReader' is now 'Tolerant Reader'; 'Semantic Versioning' is now 'Semantic Versioning 2.0.0'
- Removed: 'while the payoff lands with the teams that read' (no Source)
- Removed: 'such as a Snowflake role with read access' attributed to the ODCS Roles page (not there)
- Removed: '`rowCount`' metric name (not in the Lead's example)
- Removed: 'ODCS 3.1 added native support for the Data Contract Specification' (the Lead says the CLI gained native ODCS 3.1 support)

</details>

### RDS read replica load balancing: Aurora, DNS, proxy or NLB

`content/explainers/databases-storage/rds-read-replica-load-balancing.md` · comparison · 2400 words · fact-checker: 98 claims checked, 9 fixed, 2 removed

**Fills Gap Day 2026-11-20** (ADR-0012): `published` is that day; `make release` keeps it.

**Drafter asks you to check:**
- [ ] Lead text was read from the AWS Database Blog RDS category RSS feed (https://aws.amazon.com/blogs/database/category/database/amazon-rds/feed/), not the live page, which returned only navigation. Check on the live page: the authors (Vijay Chakraborty, Debapriyo Jogi, Arun Kumar N C), the date of 17 November 2026, and every heading anchor used (#cost-estimate, #aurora-reader-endpoint, #route-53-weighted-routing, #connection-pooling-with-proxysql-and-pgbouncer, #solution-overview, #how-it-works, #trigger-mechanism, #key-design-decisions, #hybrid-and-on-premises-connectivity, #considerations-and-limitations).
- [ ] Lead quotes were relayed with a 125-character cap, so check each exactly: "7 TB × 0.667 ≈ 4.7 TB", "If you run Aurora, you can use its built-in reader endpoints", "straightforward to set up", "no awareness of the caller's Availability Zone", "prone to client-side DNS caching that causes uneven traffic distribution", "A single proxy instance is pinned to one Availability Zone", "the stable entry point for read replica traffic", "Registration is immediate, while deregistration follows a cautious approach.", "receive load balancer IP addresses within their own Availability Zone", "On-premises clients don't natively resolve", "significant operational overhead", "no additional DNS infrastructure, no Resolver endpoints, and no private hosted zone configuration are needed", "Note that the Network Load Balancer balances connections, not queries", "you must deploy, patch, and scale the proxy tier yourself".
- [ ] $47/month is our arithmetic: 4.7 TB (decimal) × the RDS pricing page's "$0.01 per GB in/out" charged on the EC2 side. The EC2 on-demand pricing page showed no rate when fetched. 'Aurora's pricing matches' rests on the Aurora page saying EC2 Regional Data Transfer charges apply, with no figure given.
- [ ] Private hosted zone gloss ('names resolve only inside the VPCs you attach to it') is our own definition; I fetched no page that states it.
- [ ] Route 53 health checks: the draft cites the newer private endpoint health check (source in your VPC, 30 s or 10 s, failure threshold 3 by default). The older 'Configuring failover in a private hosted zone' page still says checkers sit outside the VPC; I cut that from the draft. Confirm the private check applies to weighted records pointing at an RDS replica, and that the 'additional charge' for 10-second fast checks (stated on health-checks-creating-values) covers private endpoint checks.
- [ ] Detection times of about 90 s (Route 53: 30 s × 3, plus one TTL) and about 60 s / 10 s (NLB: 30 s or 5 s × unhealthy threshold 2) are our arithmetic. The NLB docs call the interval approximate and say checks use consensus.
- [ ] 'By default it closes connections to an unhealthy target' comes from target_health_state.unhealthy.connection_termination.enabled default true, on the Target groups page (#target-group-attributes); that sentence has no inline link of its own. The 300 s deregistration delay is from the same section. TCP connections are not terminated at the end of the delay by default; the draft does not claim they are.
- [ ] ProxySQL defaults come from the monitor-variables page: ping 8,000 ms, replication lag 10,000 ms. That page contradicts itself on monitor_connect_interval (not used). The Aurora page gives check_interval_ms as 1,000 ms by default. Two ProxySQL pages describe lag_num_checks differently (not used). 'supports the PostgreSQL protocol' is from the PostgreSQL configuration page, which states no version.
- [ ] PgBouncer: the anchors #dns_max_ttl and #so_reuseport are unverified, because the fetch did not expose ids. '1.17.0, March 2022' is from the changelog (2022-03-23). 'Configuration and feature pages describe no read/write routing' is an absence claim. The inference that successive transactions in transaction mode can reach different replicas is labelled 'we infer'.
- [ ] Absence claims rest only on the Aurora Reader, Overview, Cluster and Custom endpoint pages: no AZ preference, no timing for a new replica, nothing on lag, nothing on on-premises clients. 'The RDS read replica guide describes no such shared endpoint' rests on USER_ReadRepl.html only; its 'Differences' sub-page was not read.
- [ ] NLCU worked example on the Lead's cost assumptions (500 new connections/s, 120 s average duration, 10 GB/hour) is our arithmetic. Active connections are 500 × 120 ≈ 60,000. The figure of 1 GB/hour per NLCU applies to IP targets. Check how the pricing page's 'per minute' charge fits with 'the dimension with the highest usage for the hour'.
- [ ] 'Up to 15 Aurora Replicas' is from the Aurora endpoint overview page. The NLB flow-hash quote is from the NLB introduction page; that link now has no anchor, to save bytes.
- [ ] Inference that a new replica's DNS name must be added to the Lambda function's configuration first rests on the Lead's DNS resolution section (names come from environment variables / the RDSReplicaDNSNames parameter).
- [ ] Internal link /apis-networking/cloudflare-hyperdrive-mysql-caching/: check that the Hyperdrive case study supports 'shows why a shared pool keeps database connections warm'.
- [ ] Title is 58 characters, with 'Aurora' standing for the reader endpoint. The description is 143 characters. Several inline links were folded into nearby links or credited 'the Lead says' to stay under the byte limit. Check that every claim still traces to a link in its section.

**Fact-checker could not confirm:**
- [ ] Live Lead page returns only navigation; all Lead checks (authors, 17 Nov 2026, 10 heading anchors, 14 quotes) used the RSS content:encoded text
- [ ] Two Lead quotes use straight apostrophes where the source has curly ones (caller's, don't); left to match house style
- [ ] PgBouncer anchors #dns_max_ttl and #so_reuseport: headings exist, ids not seen
- [ ] Route 53 private endpoint health check with weighted records: no page says it works (now hedged); the fast-check 'additional charge' is on the general Request interval value, which the private-check section points to
- [ ] $47/month (4,700 GB x $0.01) and Aurora sharing the EC2 rate are our arithmetic/inference
- [ ] Section-level linking remains for bytes: Route 53 TTL/health-check sentences (values page linked in the same section, not the same sentence); Lead quotes in the AZ affinity section (Lead linked once in-section); Lead claims in the 'Connections or queries' and 'What you operate' sections, with no Lead link of their own; ProxySQL '>0' and 'until it catches up' are on the mysql-tables page while the link points at the monitor variables page

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

- Fixed: Read replicas: 'copies of a primary ... each reached on its own endpoint' -> 'read-only copies of a DB instance, which applications connect to just as to any DB instance' (USER_ReadRepl.html wording)
- Fixed: Aurora pricing: 'matches' -> 'applies the same EC2 charges' (page names EC2 Regional Data Transfer charges, no rate)
- Fixed: Proxy tier: 'adds a hop' -> 'breaks AZ affinity' (Lead: 'This prevents AZ-affinity traffic flow')
- Fixed: ProxySQL Aurora: 'disables replicas over max_lag_ms' -> 'at or over' (source: 'greater or equal than max_lag_ms')
- Fixed: Route 53 private endpoint health check: ~90 s inference made conditional; the page shows the check only with failover routing, not weighted records
- Fixed: Private NLB name resolving through public DNS to private, static IPs: credited to the Lead (NLB docs page does not say it)
- Fixed: PgBouncer 'single-threaded and uses one CPU core per instance': put in quotation marks (verbatim from the so_reuseport entry)
- Fixed: Added inline links: PgBouncer config page x2 (host-list quotes), RDS read replica guide ('You must manually create read replicas'), NLB target groups #target-group-attributes (connection termination default)
- Fixed: Sources: added 'solution overview' to the Lead's sections used
- Removed: Private hosted zone gloss 'whose names resolve only inside the VPCs you attach to it' (unsourced; AWS also allows inbound endpoints)
- Removed: Filler cut for bytes: 'and the sections below source each step', 'For those, another option has to pick.', 'They work at different levels.', 'per the RDS pricing page' (now 'as of November 2026')

</details>

### Not drafted

- none

<details><summary>Leads the editor rejected</summary>

- `bucmcx76`: Brief window, but an opinion essay on network operators' supply-chain risk: no event at a named system.
- `inrxbn7h`: Community post profiling an elected SIG chair.
- `hyqudmlz`: Language-process perspective piece, no release or design change; no one architectural idea to teach.
- `y77kclvi`: Fetched: a benchmark write-up, not an event, and no live Explainer teaches KV offload, so not a Brief. A strong Concept (kv-cache-offloading), held only because ai-ml-systems is the thickest Vertical: pick it next run.
- `6cxhrm44`: Call for papers and sponsors.
- `c4lkzgaf`: Newsletter roundup; value objects and lazy constants already covered by java-immutable-data.
- `rn3af3oh`: Product feature announcement (Duplex in Gemini); no live Explainer teaches its lesson and no clear mechanism for an Explainer.
- `wgppn7d3`: Fetched: real Ruby heap overflow (ß case-mapping expansion, reachable pre-auth via Rack MethodOverride), but the Source doesn't say what the WAF rules match and no live Explainer teaches the lesson (virtual patching or input-expansion overflows).
- `zzorbmvg`: Fetched: PHP session-handler torn write leading to object injection; Cloudflare is not PHP's first-party Source, and no live Explainer teaches torn writes or virtual patching.
- `ldyxsokv`: Fragments: a link roundup and commentary.
- `l4p7nqgs`: Fetched: rustc parallel frontend on nightly; benchmark tables and flags only, with no internal mechanism for a Case Study or Concept; past the Brief window.
- `gifu2zhy`: Valhalla value objects already covered by live java-immutable-data.
- `ors6i6yf`: JEP 547 lazy constants already covered by live java-immutable-data.
- `kmex2pgr`: Fetched: good workload-skew idea (hot keys hurt lock tables, help caches), but short, and its numbers sit only in interactive figures; weaker than today's picks, so revisit.
- `suwoz4in`: Fetched: substantive persistent-kernel (MonoKernel) decode post, but ai-ml-systems is the thickest Vertical, so revisit.
- `plse63y4`: Fetched: HashiCorp product positioning; it names static vs dynamic credentials with no mechanism, TTLs or dimensions for a Comparison.
- `wquefmqb`: Fetched via feed: an operational how-to on vacuumlo and orphan cleanup; no comparison against bytea in depth or S3, so not enough for a Comparison or Concept.
- `542g6wfu`: Netflix TechBlog returned HTTP 403, so the mechanism could not be confirmed.
- `lzoq2a6f`: Article body not retrieved; CDC-to-Iceberg overlaps the live apache-iceberg-table-format, and today's slots went to verified Sources.
- `47zffks3`: Petsem already covered by live Brief fly-io-petsem-sqlite-lock-contention.
- `45dzl3k6`: Past the Brief window; a list of CVEs with WAF coverage, not one system or idea.
- `lanmblhy`: Past the Brief window; WAF rule announcement for a CVE, no live Explainer on virtual patching.
- `okz2zwx3`: Vendor virtual-patch announcement; past the Brief window, no Explainer-grade mechanism.
- `ttqmurga`: Resilience testing on AWS overlaps the live aws-resilience-hub-testing.
- `mzvcbq2f`: Major release notes list; past the Brief window.
- `w3xcwxy5`: Opinion essay.
- `2ucuqqlb`: Opinion and commentary on Rust rewrites of curl.

</details>


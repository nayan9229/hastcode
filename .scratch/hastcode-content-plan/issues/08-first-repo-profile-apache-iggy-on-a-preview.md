# 08: First Repo Profile: apache/iggy on a Preview

**What to build:** The whole path once, for real. Run the fetch command for apache/iggy, start from the template, draft the opinion sections from the Fact Pack and README (Architecture with a Mermaid diagram of the server's shape, Use it if / Don't use it if, three to five Alternatives, five to eight FAQ questions phrased as searches, Verdict), pass the production build, and push a branch so a Preview is deployed. The agent drafts; the Author verifies every claim against its source, edits, and decides whether it merges. The remaining four of the first five (xerj-org/xerj, NVIDIA/OpenShell, deepseek-ai/DeepInfer, php/frankenphp) follow the same path as editorial routine, not as tickets.

Spec: `../spec.md` ("Further Notes": first five). Editorial: `docs/editorial/repo-profile-selection.md`, the writing guide from 07, `content/pages/editorial-policy.md`.

**Blocked by:** 03 (Fact sections, Sources merge and structured data), 04 (Generated hero card), 05 (Fetch command), 07 (Editorial template, writing guide, roundup proof and docs)

**Status:** ready-for-human

- [x] Fact Pack fetched with the command; `needs_review` false (Apache-2.0)
- [x] Draft follows the template; every factual claim in prose cites a Fact Pack field, the README or a release page; quickstart commands are verbatim from the README (architecture and clustering claims cite iggy's own docs pages; the four Alternatives lines link each project's own site, as the guide allows, but we did not read those sites: see "Author attention")
- [x] Production build green: skeleton rule, section-length guide, 50 KB gate, link gate (verified on a scratch copy with `draft: true` and the deliberate `TODO(author)` line both removed; on the branch itself the draft is excluded from production output; see notes)
- [ ] Preview deployed and its URL recorded in this ticket. Not possible: the repository has no Git remote. Local preview build instead: `make build MODE=preview OUT=<dir>` renders `<dir>/static/distributed-systems/apache-iggy/index.html` (checked at `/private/tmp/claude-502/-Users-nayanmevada-Desktop-Genuin-demos-go-articles/03a3f308-7c88-42a4-9a9e-413a27ae22bd/scratchpad/preview/static/distributed-systems/apache-iggy/index.html`, a session scratch directory)
- [ ] Author verification: every claim checked, draft flag removed only by the Author (ready-for-human step; the agent stops at the Preview)
- [ ] Day-28 checkpoint note scheduled (date written here) using the template from 07. Day 28: 28 days after the Production Publish date (not yet published)

## Implementation notes

Branch `cp/08-first-profile-iggy`. Profile: `content/explainers/distributed-systems/apache-iggy.md` (`draft: true`); its Fact Pack beside it, `apache-iggy.json`, written by `make fetch` and not edited.

**Fact Pack** (fetched 2026-10-22T14:07:04Z): 5,160 stars, Apache-2.0 (`needs_review: false`), Rust, 154 contributors, latest release `server-0.9.0` (2026-09-18), default branch `master`, last push 2026-10-22.

**Draft:** about 1,170 words (about 1,055 outside code fences); every H2 under 250 words (FAQ was trimmed after a first `section-length` warning). Seven FAQ questions, four Alternatives, one Mermaid flowchart of the write path (client, listener on shard 0, `IggyNamespace` router, shard 0 metadata and `metadata/journal.wal`, owning shard, partition primary, VSR replicas, `.log`/`.index` segment), using only components the architecture docs name. In 60 seconds is the README's Docker section, both commands verbatim. The body's first line is the deliberate `TODO(author): verify every claim against its source, then remove draft: true and this line.`; no other `TODO(author)` remains. `published` and `reviewed` hold 2026-10-22T09:00:00Z as a placeholder: the Author sets the real publish time.

**Builds:**

- `make build MODE=preview`: 3 warnings, all `author-placeholder` for that one line (the page, `/distributed-systems/feed.xml`, `/feed.xml`). The page renders with the header card (`repo-profile-card`), Signals after Architecture, Changelog after Alternatives, the generated hero (`/assets/heroes/apache-iggy-card-*`) and the diagram; JSON-LD has Article, SoftwareSourceCode and FAQPage (7 Questions); HTML 35,422 bytes.
- `make build MODE=production` on this branch: green, 0 warnings, because a draft is left out of production output (the page is not written). The draft flag is the gate that keeps it from shipping.
- Production build of a scratch copy with only `draft: true` removed: fails, `build failed: 3 error(s)`, all `author-placeholder`: `/distributed-systems/apache-iggy/`, `/distributed-systems/feed.xml` and `/feed.xml` still contain "TODO(author)". Removing the flag without the line cannot ship.
- Production build of a scratch copy with both removed: green, 0 warnings, page 35,323 bytes (skeleton, section-length, 50 KB and link gates pass).
- `make check` passes.

**Sources consulted:** the README (cache copy from the fetch, blob `07990be`); https://iggy.apache.org/docs/introduction/architecture ; https://iggy.apache.org/docs/clustering/vsr ; https://github.com/apache/iggy/releases/tag/server-0.9.0 (a list of merged pull requests that links the release post); https://iggy.apache.org/blogs/2026/09/21/release-0.9.0 ; the Hacker News item from the shortlist (https://news.ycombinator.com/item?id=50181389), read only for its link.

**Left out on purpose:**

- Every throughput and latency number ("millions of messages per second", "over 5000 MB/s", microsecond p99): first-party claims that depend on hardware and that we cannot verify. The Profile points to the project's benchmarking platform instead.
- Adopters and production users: no first-party source we read names any.
- Apache Incubator or top-level-project status: not stated in the sources read. The license FAQ says only "Apache License 2.0, developed in the Apache Software Foundation's `apache` GitHub organisation".
- The VSR clustering part 3 post the HN item links (https://iggy.apache.org/blogs/2026/10/27/vsr-clustering-part-3/): the page states a publication date of October 27, 2026, after today and after the HN submission (2026-10-21). Not cited; the Author should check it.
- Why VSR rather than Raft: the only first-party reason found ("determinism") sits in that part 3 post.
- Connectors, MCP server, Web UI, Helm chart, encryption, S3 archiving, edge releases: true per the README but outside the one idea (the log and how it replicates).
- Stars, contributor counts and release dates in prose: the card, Signals and Changelog render them from the Fact Pack.

**Author attention:** the four Alternatives lines describe Kafka, Redpanda, NATS JetStream and Pulsar from general knowledge, each linked to its project's own site; check each against that site. "Linux" in Use it if and Verdict rests on the README's `io_uring` and "modern Linux systems" lines. The Verdict's opening ("serious, carefully documented") is opinion.

## Remaining human steps

1. Verify every claim, command and the diagram against its linked source; rewrite what is wrong.
2. When a remote exists, push the branch, check the Preview (header card, Signals, Changelog, hero card) and record its URL here.
3. Set `published` and `reviewed` to the real publish time, delete the `TODO(author)` line and `draft: true`, merge, then tag (`scripts/publish.sh vl`, check QA, `scripts/publish.sh cf`).
4. Write the day-28 date here (Production Publish date plus 28 days) and, on that day, fill in a copy of `docs/editorial/checkpoint-note-template.md`.

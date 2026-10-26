# 06: Draft the real Source Watchlist

**What to build:** The Author has a reviewed starting Watchlist. Draft it from the research tables: 6 Sources per Vertical, chosen for first-party signal and steady cadence, plus the 11 incident feeds (kind case-study), each with a one-line reason. Prefer feeds that sent validators; avoid the stale or blocked ones the research lists. Run one real `make intake` against it, fix any Source that fails, and record the run summary (Leads per Vertical) in the ticket.

Spec: `../spec.md` ("Source Watchlist"). Research: `docs/research/content-source-intake.md` (tables per Vertical, incident feeds, terms notes).

**Blocked by:** 02 (Intake tracer), 03 (polite fetching), 04 (GitHub Sources)

**Status:** resolved

- [x] Watchlist has 36 Vertical Sources (6 each) plus the incident feeds, every one validated; 40 of 47 fetched successfully in one real run (the 7 repository Sources failed only because `GITHUB_TOKEN` was not set in the agent's environment; see below)
- [x] Each Source carries a reason; any swap from the research's top picks is explained (in `sources.yaml` comments per Vertical, and below)
- [x] The run summary is pasted into this ticket
- [x] Ends `ready-for-human`: the Author reviews the list before it is committed; the agent does not commit it unreviewed (committed on branch `intake/06-watchlist` only; not merged until the Author reviews)

## What was done (2026-10-26)

`sources.yaml` at the repo root: 47 Sources, 6 per Vertical plus 11 incident feeds, 40 feeds and 7 GitHub repositories (etcd, Valkey, Envoy, Kubernetes, vLLM, PyTorch, Node.js). Every feed URL was re-fetched with the Intake User-Agent before the run (all 200 and RSS/Atom).

| Vertical | Sources (kind) |
|---|---|
| distributed-systems | AWS Architecture Blog (case-study), Netflix TechBlog (case-study), Meta Engineering, Production Engineering (case-study), Marc Brooker (concept), Martin Fowler (concept), etcd releases (brief, repo) |
| databases-storage | PostgreSQL news (brief), AWS Database Blog (case-study), DuckDB blog (concept), PlanetScale blog (comparison), TigerBeetle blog (concept), Valkey releases (brief, repo) |
| apis-networking | Cloudflare blog (case-study), RFC Editor (brief), APNIC blog (concept), Fastly blog (case-study), Daniel Stenberg (curl) (concept), Envoy releases (brief, repo) |
| infrastructure-cloud | Kubernetes blog (concept), Kubernetes releases (brief, repo), GitHub Engineering (case-study), Shopify Engineering (case-study), HashiCorp blog (comparison), Fly.io blog (case-study) |
| ai-ml-systems | vLLM blog (concept), vLLM releases (brief, repo), PyTorch blog (concept), PyTorch releases (brief, repo), Google Research blog (concept), Meta Engineering, ML Applications (case-study) |
| languages-runtimes | Go blog (brief), Rust blog (brief), Python Insider (brief), Python PEPs (concept), Inside Java (concept), Node.js releases (brief, repo) |
| incident feeds (case-study) | AWS Health Dashboard, Azure status, Google Cloud Service Health, GitHub status, Fly.io Infra Log (infrastructure-cloud); Cloudflare blog, Post Mortem tag (apis-networking); GitHub availability reports, Slack status, Discord status, Jira Software status (distributed-systems); npm status (languages-runtimes) |

### Swaps from the research's picks (and why)

- **Kafka and CockroachDB releases dropped as repository Sources**: `apache/kafka` and `cockroachdb/cockroach` publish tags but no GitHub Releases (their Releases pages say "There aren't any releases here"); the research's `releases.atom` counts were tags. `GET /releases` would list nothing. etcd stays for Distributed Systems; **Valkey** replaces CockroachDB for Databases.
- **Meta Engineering main feed → its Production Engineering category** (Distributed Systems): 5 of the main feed's 9 items were ML posts that the AI/ML Source (Meta ML Applications category) already lists, so the same post would appear twice under two Verticals. The category feed is not in the research tables; verified 2026-10-26 (200, ETag + Last-Modified, 5 items in 90 days).
- **Atlassian status → Jira Software status** (incident): `status.atlassian.com/history.rss` is stale (newest 2025-04-30, as the research noted); Atlassian's live product-level Statuspage `jira-software.status.atlassian.com/history.rss` (14 items in 90 days, ETag) replaces it.
- **Cloudflare status → npm status** (incident): on the first run `cloudflarestatus.com/history.atom` held only scheduled maintenance (25 of 25 items, dated up to a week ahead), which pushes real incidents out of the feed. Cloudflare outages still arrive through the Post Mortem tag Source. `status.npmjs.org/history.rss` (11 in 90 days, ETag; not in the research) covers the registry every Node.js install depends on. Revert if you would rather keep the research's eleven exactly.
- **Left out by policy** (per-Vertical comments in `sources.yaml` list each): firehoses (AWS What's New, Google Cloud blog, Azure updates, NVIDIA, Hugging Face, arXiv); aggregators (Planet PostgreSQL); feeds with no validators where a validator-sending equivalent exists (Node.js blog → Node.js releases, Redis, Tailscale, Stripe, Dropbox, Airbnb); low cadence (All Things Distributed, IETF blog, Istio, Cilium); and Go tags (not the REST API; Go blog is Go's release signal). Kept despite no validators: Netflix TechBlog, vLLM blog, Go blog (no equivalent first-party source).

### Run summary (one real `make intake`, 2026-10-26 16:02 UTC, fresh `.intake/`)

Command output (failure lines shortened):

```
1488 new Leads, 1488 open: .intake/queue.md
failed etcd releases: GITHUB_TOKEN is not set: …
failed Valkey releases: GITHUB_TOKEN is not set: …
failed Envoy releases: GITHUB_TOKEN is not set: …
failed Kubernetes releases: GITHUB_TOKEN is not set: …
failed vLLM releases: GITHUB_TOKEN is not set: …
failed PyTorch releases: GITHUB_TOKEN is not set: …
failed Node.js releases: GITHUB_TOKEN is not set: …
silent Cloudflare blog, Post Mortem tag: newest item 2026-05-01 (177 days ago)
```

- **Fetched:** 40 of 40 feed Sources (Azure status returned a valid feed with 0 items: it carries only live incidents).
- **Failed:** the 7 repository Sources, all for the same reason only: no `GITHUB_TOKEN` in the agent's environment or `.env`. Not a Source problem: run `GITHUB_TOKEN=… make intake` once to confirm them.
- **Silent:** Cloudflare blog, Post Mortem tag (expected: post mortems are months apart).

Leads per Vertical and kind. **Totals include back catalogue**: Intake has no age cutoff yet, so a first run lists every item each feed carries (Shopify's feed alone holds 436, Marc Brooker's 166). The 14-day column counts Leads whose published date is within 14 days of the run (from `.intake/state.json`), which is what a daily run will look like. Incident-feed Leads are included in their Vertical and also shown on their own.

| Vertical | Leads (includes back catalogue) | by kind | of which incident feeds | published ≤14 days | by kind | of which incident feeds |
|---|---|---|---|---|---|---|
| distributed-systems | 327 | case-study 131, concept 196 | 92 | 30 | case-study 24, concept 6 | 5 |
| databases-storage | 97 | case-study 20, concept 42, comparison 25, brief 10 | 0 | 30 | case-study 11, concept 4, comparison 5, brief 10 | 0 |
| apis-networking | 120 | case-study 65, concept 40, brief 15 | 20 | 46 | case-study 28, concept 17, brief 1 | 0 |
| infrastructure-cloud | 660 | case-study 590, concept 50, comparison 20 | 104 | 24 | case-study 23, concept 1 | 14 |
| ai-ml-systems | 169 | case-study 9, concept 160 | 0 | 39 | case-study 8, concept 31 | 0 |
| languages-runtimes | 115 | case-study 25, concept 20, brief 70 | 25 | 24 | case-study 4, concept 11, brief 9 | 4 |
| **total** | **1488** | case-study 840, concept 508, comparison 45, brief 95 | 241 | **193** | case-study 98, concept 70, comparison 5, brief 20 | 23 |

Briefs are thin in the 14-day view (20) because all 7 repository Sources, which are Brief Sources, failed for the token; with the token they add each project's recent releases.

### Observations for the code owners (not fixed here; `internal/intake` untouched)

- **No age cutoff** (known; fix in progress on another branch): first run = 1488 Leads, ~87% back catalogue.
- **Future-dated items**: on the first run, Cloudflare status's scheduled-maintenance entries (dated up to 2026-11-02) rendered as `2026-11-02 (today)` and sorted to the top of their Vertical. The queue's age text treats a future date as "today"; a status feed with maintenance windows will do this again.
- **Tracking parameters in links**: Medium-hosted feeds (Netflix TechBlog) link with `?source=rss----…`; covered detection (05) should normalise query strings away or Netflix Leads will never match a cited URL.

## For the Author's review

Check before merging `intake/06-watchlist`:

1. The 47 Sources and their one-line reasons in `sources.yaml` (kind and Vertical are copied onto every Lead, so a wrong default misfiles every Lead from that Source).
2. The four swaps above, in particular the two not in the research (Meta Production Engineering category, npm status) and dropping Cloudflare status.
3. Run `GITHUB_TOKEN=… make intake` once to confirm the 7 repository Sources (expect each to list its recent full releases; prereleases are skipped).
4. Terms: AWS and Cloudflare reserve all rights (link and quote briefly, never reuse diagrams); Medium-hosted Netflix terms were not read (research Open question 3).

## Author review

Approved as drafted by the Author on 2026-10-26 (47 Sources).

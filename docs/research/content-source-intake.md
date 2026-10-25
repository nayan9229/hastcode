# First-party content sources and the rules for using them (Case Studies, Concepts, Comparisons, Briefs)

Research for hastcode.com (software-architecture Explainers and Briefs under the disclosed House Byline Mira Vance; AI-drafted, human-verified; first-party sources only, every claim linked: `content/pages/editorial-policy.md`; curated, never exhaustive: ADR-0005). Question: which first-party sources can feed a steady, curated intake for the six Verticals in `site.yaml`, and what rules (Google Search / Google News policy, copyright, HTTP and GitHub API etiquette) govern using them?

Researched **2026-10-26**. Every feed URL below was fetched on that date. Labels: **[P]** primary (publisher's own page, statute, spec, official guidance) / **[observed]** what a live request returned / **[inferred]** my deduction / **[opinion]**. Not legal advice: section 4 states what the primary texts say; we are not lawyers.

---

## TL;DR

- **There is enough first-party supply.** I requested about 130 candidate feed URLs and verified 77 live first-party feeds in the six Vertical tables (13 Distributed Systems, 12 Databases & Storage, 13 APIs & Networking, 14 Infrastructure & Cloud, 10 AI/ML Systems, 15 Languages & Runtimes; 83 with the lower-priority extras listed under the tables), plus 11 incident/status feeds. Their items dated in the last 90 days add up to at least ~1,700 (a lower bound: many feeds only show their latest 10-20 items) [observed]. Supply is not the constraint; the Author's verification time is (Editorial Policy, ADR-0005).
- **The best Case Study sources are engineering blogs and post-incident reviews**: AWS Architecture/Database blogs, Netflix, Meta, Dropbox, Airbnb, Shopify, GitHub, Cloudflare, Fastly, plus AWS Post-Event Summaries ("will remain available for a minimum of 5 years"), Azure PIRs ("each retained for 5 years"), Cloudflare's Post Mortem tag feed, GitHub's monthly availability reports and Fly.io's Infra Log. **Concept** sources: Marc Brooker, Martin Fowler, All Things Distributed, Inside Rust/Inside Java, PEPs, RFCs, arXiv, Google/Microsoft Research. **Comparison** inputs: Jepsen analyses (no feed), MLCommons, vendor posts read against each other. **Briefs**: project blogs and release feeds (Kubernetes, Go, Rust, PostgreSQL, Python, Node, Envoy, vLLM, PyTorch) and the RFC Editor feed.
- **Not fetchable or no feed**: ACM Queue and CACM (403), MySQL blog (403), Uber Engineering (406), Jepsen, LinkedIn, CockroachDB blog and Anthropic engineering (no feed found); USENIX's RSS is stale since 2025-03 and MongoDB's since 2026-06.
- **Google's policies bite hardest on Briefs.** The spam policies name "Scraping feeds … to generate many pages", "Stitching or combining content from different web pages without adding value" and "Reproducing content feeds from other sites without providing some type of unique benefit to the user". The helpful-content guide lists "Are you mainly summarizing what others have to say without adding much value?" as a warning sign and says "**Attribution or giving credit to other sources doesn't replace the need for original effort.**" Google ranks "original reporting, ahead of those who merely cite it". So a Brief survives only on its "why it matters" analysis, picked by the Author, never generated from a feed.
- **Copyright: no regime gives a safe word count.** US Copyright Office: "there is no formula … specific number of words, lines, pages". EU InfoSoc Art. 5(3)(d) allows quotation "for purposes such as criticism or review" only with "the source, including the author's name" and "to the extent required by the specific purpose". India s.52(1)(a) allows fair dealing for criticism, review and "reporting of current events" but **excludes computer programmes**. The EU press publishers' right (DSM Art. 15) exempts "hyperlinking" and "individual words or very short extracts", and does not cover blogs outside a news publisher's editorial control (Recital 56).
- **Diagrams are the sharpest limit**: they are protected works in the US ("diagrams, models, and technical drawings", §101) and India ("a drawing (including a diagram, map, chart or plan)… whether or not any such work possesses artistic quality", s.2(c)). AWS's and Cloudflare's site terms reserve all rights. Redraw architectures from the facts (§102(b): no copyright in "any idea, procedure, process, system, method of operation") and link the original. Reuse wholesale only under a licence: Go (CC BY 4.0), Kubernetes docs (CC BY 4.0), LF/CNCF (CC BY 3.0). Google's SRE book is CC BY-NC-ND, which an ad-funded site cannot use to reproduce it.
- **Intake mechanics**: support RSS 2.0 and Atom (71 and 35 of the fetched feeds; **zero JSON Feeds** seen). Poll daily, sequentially, with `If-None-Match` / `If-Modified-Since` from stored validators (RFC 9110 §13.1); verified 304s from Kubernetes, Cloudflare and GitHub. The Go blog feed sends no validators. For GitHub releases use the REST API with the existing token: 5,000 requests/h, and an authenticated 304 "does not count against your primary rate limit".

---

## 1. Source inventory per Vertical

**How these were checked [observed].** Every feed URL was requested on 2026-10-26 with `curl -sL --compressed` and parsed as XML; "fetched" means HTTP 200 and a parseable RSS or Atom document. **Cadence** is what the feed itself showed, not a publisher claim: "N in 90 d" counts items dated in the 90 days to 2026-10-26; where every item in the feed falls inside a short window the feed is truncated, so the cadence is given as "N items in D days". **Cond.** shows whether the response carried an `ETag` (E) and/or `Last-Modified` (L), the validators conditional GET needs (section 5). **Fits** names the Explainer kind or Brief the source most naturally feeds **[opinion]**. **Terms** quotes the site's own licence or terms page where I fetched one; "not checked" means I did not find or fetch a statement, and the safe assumption is all rights reserved **[inferred]**.

Three cross-cutting facts frame every table:

- **GitHub `releases.atom` / `tags.atom` feeds are github.com web endpoints, not the REST API** [observed: `https://github.com/{owner}/{repo}/releases.atom` returns `application/atom+xml` with a weak `ETag` and answered `If-None-Match` with 304]. ADR-0005 chose the REST API as "the sanctioned path" for the Fact Pack fetcher; whether polling these Atom feeds counts as "scraping" under GitHub's Acceptable Use Policies is an open question (see Open questions). The REST alternative is in section 5.3.
- **Vendor engineering blogs are first-party for their own systems only.** A Case Study of Netflix's system cites Netflix's blog; the same blog is not a source for claims about anyone else's system (Editorial Policy, "Sources").
- **Aggregators are discovery aids, never citations.** Planet PostgreSQL and arXiv listings point at first-party posts and papers; cite the post or paper itself.

### 1.1 Distributed Systems

| Source | Site | Feed URL | Fetched | Cadence (observed) | Cond. | Fits | Terms |
|---|---|---|---|---|---|---|---|
| AWS Architecture Blog | https://aws.amazon.com/blogs/architecture/ | https://aws.amazon.com/blogs/architecture/feed/ | yes (RSS) | 20 items in 34 d (~4/week) | L | Case Study, Concept | AWS Site Terms: "limited license to access and make personal use"; excludes "any resale or commercial use … any derivative use" (https://aws.amazon.com/terms/). Link and quote briefly; do not reproduce diagrams. |
| Netflix TechBlog | https://netflixtechblog.com/ | https://netflixtechblog.com/feed | yes (RSS, Medium) | 10 items in 8 d (feed shows last 10) | none | Case Study | Medium-hosted; Medium's ToS page returned 403, not checked |
| Meta Engineering | https://engineering.fb.com/ | https://engineering.fb.com/feed/ | yes (RSS) | 9 items in 6 d (truncated) | E, L | Case Study | not checked |
| Marc Brooker (AWS distinguished engineer, personal blog) | https://brooker.co.za/blog/ | https://brooker.co.za/blog/rss.xml | yes (RSS, full archive 166 items since 2012) | 4 in 90 d | E, L | Concept | not checked |
| All Things Distributed (Werner Vogels, Amazon CTO) | https://www.allthingsdistributed.com/ | https://www.allthingsdistributed.com/index.xml | yes (RSS) | 3 in 90 d | E, L | Concept, Case Study | not checked |
| Martin Fowler | https://martinfowler.com/ | https://martinfowler.com/feed.atom | yes (Atom) | 30 in 90 d (includes article instalments and guest articles) | E, L | Concept, Comparison | not checked |
| Dropbox Tech | https://dropbox.tech/ | https://dropbox.tech/feed | yes (RSS) | 5 in 90 d | none | Case Study | not checked |
| Airbnb Engineering | https://medium.com/airbnb-engineering | https://medium.com/feed/airbnb-engineering | yes (RSS, Medium) | 9 in 90 d | none | Case Study | as Netflix |
| Microsoft Research blog | https://www.microsoft.com/en-us/research/blog/ | https://www.microsoft.com/en-us/research/feed/ | yes (RSS) | 10 items in 23 d | E, L | Concept (paper-backed) | not checked |
| arXiv cs.DC (Distributed, Parallel, and Cluster Computing) | https://arxiv.org/list/cs.DC/recent | https://rss.arxiv.org/rss/cs.DC | yes (RSS), but **0 items at access**: the channel declares `<skipDays>` Saturday and Sunday and was last built Sun 25 Oct | daily on weekdays; arXiv: "Feeds are updated daily at midnight Eastern Standard Time" (https://info.arxiv.org/help/rss.html) | E | Concept (papers) | Each paper carries the licence its authors chose; "A Creative Commons CC0 1.0 Universal Public Domain Dedication will apply to all metadata" (https://info.arxiv.org/help/license/index.html) |
| Release feeds: etcd, Apache Kafka, Temporal | https://github.com/etcd-io/etcd , https://github.com/apache/kafka , https://github.com/temporalio/temporal | `https://github.com/{etcd-io/etcd,apache/kafka,temporalio/temporal}/releases.atom` | yes (Atom) ×3 | etcd 10 on one day (patch releases on several branches); Kafka 10 in 28 d; Temporal 10 in 14 d | E | Brief | Release notes are repository content; GitHub ToS D.5 licenses other users only to "use, display, perform and reproduce (by forking)" through the Service (https://docs.github.com/en/site-policy/github-terms/github-terms-of-service) |

**Checked but not usable as a feed:**

- **Jepsen analyses** (https://jepsen.io/analyses): page fetched (200), no feed autodiscovery link; `/blog.xml`, `/analyses.xml`, `/feed.xml` all 404. Its author's personal blog https://aphyr.com/posts.atom fetched (Atom, 12 items, newest 2026-06-12, 0 in 90 d). Watch the analyses page by hand; each analysis is first-party test evidence about the database it tests, prime Case Study and Comparison material.
- **ACM Queue** (https://queue.acm.org/): `https://queue.acm.org/rss/feeds/queuecontent.xml` and the home page returned **403** to both a bot and a browser User-Agent; `https://cacm.acm.org/feed/` also 403. **Not fetched.**
- **USENIX** (OSDI, NSDI, FAST, ATC, SREcon): `https://www.usenix.org/rss.xml` fetched but **stale** (newest item 2025-03-26; 0 in 90 d). The site says "We banished paywalls in 2008, providing conference proceedings free of charge under our open access policy … recorded presentations of our events are publicly available after each event" (https://www.usenix.org/about). Use conference programme pages per event; no live feed found. No reuse licence was stated on the pages fetched.
- **Uber Engineering**: `https://www.uber.com/en-US/blog/engineering/rss/` redirects to `/US/en/blog/engineering/rss/` and answers **406**, also for the HTML page. **Not fetched.**
- **LinkedIn Engineering**: page fetched, no autodiscovery link; guessed `rss.xml` 404. **No feed found.**
- **Slack Engineering** (https://slack.engineering/feed/): fetched, newest 2026-07-14, 0 in 90 d. Low cadence; keep for incident write-ups (section 2).

### 1.2 Databases & Storage

| Source | Site | Feed URL | Fetched | Cadence (observed) | Cond. | Fits | Terms |
|---|---|---|---|---|---|---|---|
| PostgreSQL news | https://www.postgresql.org/about/newsarchive/ | https://www.postgresql.org/news.rss | yes (RSS) | 10 items in 4 d (truncated; includes community and extension announcements) | L | Brief (releases, security) | Copyright "The PostgreSQL Global Development Group"; the PostgreSQL Licence grants permission "to use, copy, modify, and distribute this software and its documentation for any purpose" with the notice retained (https://www.postgresql.org/about/licence/). It names software and documentation, not news posts. |
| Planet PostgreSQL (aggregator of contributors' blogs) | https://planet.postgresql.org/ | https://planet.postgresql.org/rss20.xml | yes (RSS) | 30 items in 5 d | L | discovery for Concept / Case Study; cite the underlying post | per-author |
| AWS Database Blog | https://aws.amazon.com/blogs/database/ | https://aws.amazon.com/blogs/database/feed/ | yes (RSS) | 20 items in 20 d | L | Case Study (Aurora, DynamoDB internals), Comparison | AWS Site Terms (see 1.1) |
| DuckDB blog | https://duckdb.org/news/ | https://duckdb.org/feed.xml | yes (Atom) | 10 items in 23 d | E | Concept, Case Study | not checked |
| PlanetScale blog | https://planetscale.com/blog | https://planetscale.com/blog/feed.atom | yes (Atom) | 25 items in 60 d | E | Concept, Comparison | not checked |
| ClickHouse blog | https://clickhouse.com/blog | https://clickhouse.com/rss.xml | yes (RSS, 920-item archive) | 102 in 90 d (~8/week, much of it product marketing) | none | Concept, Case Study (filter hard) | not checked |
| Redis blog | https://redis.io/blog/ | https://redis.io/blog/feed/ | yes (RSS) | 37 in 90 d | none | Brief, Concept | not checked |
| TigerBeetle blog | https://tigerbeetle.com/blog/ | https://tigerbeetle.com/blog/atom.xml | yes (Atom) | 2 in 90 d (32 since 2021) | E, L | Concept (testing, storage-fault models) | not checked |
| Release feeds: CockroachDB, DuckDB, Valkey | https://github.com/cockroachdb/cockroach , https://github.com/duckdb/duckdb , https://github.com/valkey-io/valkey | `…/releases.atom` | yes (Atom) ×3 | CockroachDB 10 in 49 d; DuckDB 2 in 90 d; Valkey 5 in 90 d | E | Brief | GitHub ToS D.5 (see 1.1) |
| SQLite check-in timeline | https://sqlite.org/src/timeline | https://sqlite.org/src/timeline.rss | yes (RSS) | 20 check-ins in ~1 d: commit-level, too noisy for intake; release notes live at https://sqlite.org/changes.html (no feed found) | none | Brief (by hand) | not checked |

**Checked but not usable:** **MongoDB blog** `https://www.mongodb.com/blog/rss` redirects to `/company/blog/rss`, fetched, but newest item is 2026-06-25 and 0 in 90 d: the feed looks abandoned while the blog continues [inferred]; `/company/blog/engineering` 404. **CockroachDB blog**: no autodiscovery link; `/blog/index.xml`, `/blog/rss.xml`, `/blog/feed/` all 404. **MySQL blog** (`https://blogs.oracle.com/mysql/rss` and the HTML page): **403**, not fetched.

### 1.3 APIs & Networking

| Source | Site | Feed URL | Fetched | Cadence (observed) | Cond. | Fits | Terms |
|---|---|---|---|---|---|---|---|
| Cloudflare blog | https://blog.cloudflare.com/ | https://blog.cloudflare.com/rss/ | yes (RSS) | 20 items in 11 d (~12/week) | L (304 on `If-Modified-Since`) | Case Study, Concept, Brief | Website Terms: "Except as expressly authorized by Cloudflare you may not make use of the Materials. Cloudflare reserves all rights to the Materials not granted expressly in these Terms." (https://www.cloudflare.com/website-terms/) Quote under fair use / quotation only. |
| RFC Editor: new RFCs | https://www.rfc-editor.org/ | https://www.rfc-editor.org/rfcrss.xml | yes (RSS) | last 15 RFCs, 2026-08-28 → 2026-10-14 | E | Brief (new RFC), Concept | IETF Trust Legal Provisions 5, §3.c: licence "to copy, publish, display and distribute IETF Contributions and IETF Documents **in full and without modification**"; code components separately (§4) (https://trustee.ietf.org/documents/trust-legal-provisions/tlp-5/). Excerpts fall back on fair use / quotation. |
| IETF blog | https://www.ietf.org/blog/ | https://www.ietf.org/blog/feed/ | yes (RSS, archive since 2012) | 8 in 90 d | L | Brief, Concept | as RFCs only if an IETF Document; otherwise not checked |
| W3C news | https://www.w3.org/news/ | https://www.w3.org/news/feed/ | yes (RSS) | 23 in 90 d | E | Brief (specs reaching CR/REC) | not checked |
| APNIC blog | https://blog.apnic.net/ | https://blog.apnic.net/feed/ | yes (RSS) | 30 items in 36 d | E, L | Concept (BGP, DNS, routing measurement) | not checked |
| Fastly blog | https://www.fastly.com/blog | https://www.fastly.com/blog_rss.xml (from autodiscovery; `/blog/rss.xml` 404) | yes (RSS) | 25 items in 48 d | E | Case Study, Concept | not checked |
| Tailscale blog | https://tailscale.com/blog | https://tailscale.com/blog/index.xml | yes (RSS) | 15 items in 68 d | none | Concept (NAT traversal, WireGuard) | not checked |
| curl / Daniel Stenberg | https://daniel.haxx.se/blog/ | https://daniel.haxx.se/blog/feed/ | yes (RSS) | 10 items in 26 d | E, L | Brief, Concept (HTTP/3, protocol edge cases) | not checked |
| Stripe blog | https://stripe.com/blog | https://stripe.com/blog/feed.rss | yes (RSS) | 10 items in 30 d (mostly product) | none | Case Study (API design, idempotency) when engineering | not checked |
| Release feeds: Envoy, nginx | https://github.com/envoyproxy/envoy , https://nginx.org/ | https://github.com/envoyproxy/envoy/releases.atom ; https://nginx.org/index.rss | yes ×2 | Envoy 10 in 16 d; nginx 7 in 90 d | E / E, L | Brief | GitHub ToS D.5 / not checked |
| Istio, Cilium blogs (service mesh, eBPF networking) | https://istio.io/latest/blog/ , https://cilium.io/blog/ | https://istio.io/latest/blog/feed.xml ; https://cilium.io/blog/rss.xml | yes ×2 | Istio 4 in 90 d; Cilium 3 in 90 d | E | Concept, Comparison | not checked |

**Checked but weak:** gRPC blog `https://grpc.io/blog/index.xml` fetched, newest 2026-06-02 (0 in 90 d). Chromium blog `https://blog.chromium.org/feeds/posts/default` fetched, newest 2026-03-31. IETF Datatracker `https://datatracker.ietf.org/feed/rfcs/` 404 (guess).

### 1.4 Infrastructure & Cloud

| Source | Site | Feed URL | Fetched | Cadence (observed) | Cond. | Fits | Terms |
|---|---|---|---|---|---|---|---|
| Kubernetes blog | https://kubernetes.io/blog/ | https://kubernetes.io/feed.xml | yes (RSS) | 25 in 90 d | E (304 on `If-None-Match`) | Brief (releases, deprecations), Concept | Footer: "© 2026 The Kubernetes Authors \| Documentation Distributed under CC BY 4.0" and "For website terms of use, trademark policy and other project policies please see https://lfprojects.org/policies/" (https://kubernetes.io/blog/; whether CC BY 4.0 covers blog posts, not only documentation, is not stated there) |
| Kubernetes contributor blog | https://www.kubernetes.dev/blog/ | https://www.kubernetes.dev/index.xml (autodiscovery; `kubernetes.dev/feed.xml` 404) | yes (RSS) | 10 in 90 d | E | Concept (design of SIG work) | as above |
| Kubernetes official CVE feed | https://kubernetes.io/docs/reference/issues-security/official-cve-feed/ | https://kubernetes.io/docs/reference/issues-security/official-cve-feed/feed.xml | yes (RSS, 94 items since 2017) | 3 in 90 d | E | Brief (security) | documentation, so CC BY 4.0 (as above) |
| Kubernetes releases | https://github.com/kubernetes/kubernetes | https://github.com/kubernetes/kubernetes/releases.atom | yes (Atom) | 10 items in 48 d | E (304 observed) | Brief | GitHub ToS D.5 |
| CNCF blog | https://www.cncf.io/blog/ | https://www.cncf.io/feed/ | yes (RSS) | 10 items in 13 d | E, L | discovery; Brief for graduations; member posts are vendor content | Footer "Terms of Use" links to the Linux Foundation terms: "Except as otherwise provided, Content on this site … is licensed under a Creative Commons Attribution 3.0 License" (https://www.linuxfoundation.org/legal/terms) |
| AWS What's New | https://aws.amazon.com/new/ | https://aws.amazon.com/about-aws/whats-new/recent/feed/ | yes (RSS) | 100 items in 9 d (~75/week): a firehose to filter, not read | L | Brief | AWS Site Terms |
| Google Cloud blog | https://cloud.google.com/blog | https://cloudblog.withgoogle.com/rss/ | yes (RSS) | 20 items in 3 d (firehose) | none | Case Study (infrastructure posts), Brief | not checked |
| Azure updates | https://azure.microsoft.com/en-us/updates/ | https://www.microsoft.com/releasecommunications/api/v2/azure/rss | yes (RSS, 200 items) | 195 in 90 d | none | Brief (filter) | not checked. The old `https://azure.microsoft.com/en-us/updates/feed/` now redirects to an HTML page. |
| GitHub Engineering | https://github.blog/engineering/ | https://github.blog/engineering/feed/ | yes (RSS) | 10 in 77 d | E, L | Case Study | not checked |
| Shopify Engineering | https://shopify.engineering/ | https://shopify.engineering/blog.atom | yes (Atom, archive since 2010) | 13 in 90 d | none | Case Study | not checked |
| HashiCorp blog | https://www.hashicorp.com/blog | https://www.hashicorp.com/blog/feed.xml | yes (Atom served as `application/rss+xml`) | 18 in 90 d | E, L | Brief, Comparison | not checked |
| Fly.io blog | https://fly.io/blog/ | https://fly.io/blog/feed.xml | yes (Atom) | 3 in 90 d | E, L | Case Study, Concept | not checked |
| Release feeds: containerd, Linux kernel | https://github.com/containerd/containerd , https://www.kernel.org/ | https://github.com/containerd/containerd/releases.atom ; https://www.kernel.org/feeds/kdist.xml | yes ×2 | containerd 10 in 27 d; kernel 9 items in 8 d | E / L | Brief | GitHub ToS D.5 / not checked |

Also fetched (vendor-heavy, lower priority): Grafana blog `https://grafana.com/blog/index.xml` (10 items in 6 d), Docker blog `https://www.docker.com/feed/` (10 items in 28 d).

### 1.5 AI/ML Systems

| Source | Site | Feed URL | Fetched | Cadence (observed) | Cond. | Fits | Terms |
|---|---|---|---|---|---|---|---|
| vLLM blog | https://vllm.ai/blog | https://vllm.ai/blog/rss.xml (autodiscovery; the old `blog.vllm.ai/feed.xml` redirects and 404s; Atom at `https://vllm.ai/blog/atom.xml`) | yes (RSS) | 46 in 90 d | none | Concept (paged attention, scheduling), Case Study | not checked |
| vLLM releases | https://github.com/vllm-project/vllm | https://github.com/vllm-project/vllm/releases.atom | yes (Atom) | 10 items in 12 d | E | Brief | GitHub ToS D.5 |
| PyTorch blog | https://pytorch.org/blog/ | https://pytorch.org/feed/ | yes (RSS) | 10 items in 3 d (truncated) | E, L | Concept, Brief | not checked |
| PyTorch releases | https://github.com/pytorch/pytorch | https://github.com/pytorch/pytorch/releases.atom | yes (Atom) | 10 on one day (includes nightly/RC tags) | E | Brief | GitHub ToS D.5 |
| Google Research blog | https://research.google/blog/ | https://research.google/blog/rss/ | yes (RSS) | 29 in 90 d | L | Concept (paper-backed), Case Study | not checked |
| Meta Engineering, ML Applications category | https://engineering.fb.com/category/ml-applications/ | https://engineering.fb.com/category/ml-applications/feed/ | yes (RSS) | 9 in 74 d | E, L | Case Study (training/serving infra) | not checked |
| NVIDIA Technical Blog | https://developer.nvidia.com/blog/ | https://developer.nvidia.com/blog/feed | yes (Atom) | 100 items in 63 d (~11/week) | E, L | Concept (inference serving, GPU memory) | not checked |
| Hugging Face blog | https://huggingface.co/blog | https://huggingface.co/blog/feed.xml | yes (RSS) | 60 in 90 d | E | Concept, Brief | not checked |
| MLCommons | https://mlcommons.org/ | https://mlcommons.org/feed/ | yes (RSS) | 12 items in 30 d | E, L | Brief (MLPerf results), Comparison | not checked |
| arXiv cs.LG | https://arxiv.org/list/cs.LG/recent | https://rss.arxiv.org/rss/cs.LG | yes (RSS), 0 items at access (weekend `skipDays`) | daily on weekdays | E | Concept (papers; very high volume) | per-paper licence; metadata CC0 |

Also fetched: OpenAI news `https://openai.com/news/rss.xml` (192 in 90 d, mostly product/company news; Brief only, filter hard); AWS Machine Learning blog `https://aws.amazon.com/blogs/machine-learning/feed/` (20 items in 3 d, firehose); Ray releases `https://github.com/ray-project/ray/releases.atom` (6 in 90 d); llama.cpp releases `https://github.com/ggml-org/llama.cpp/releases.atom` (10 on one day: a build per merge, too noisy for intake). **Anthropic engineering** (https://www.anthropic.com/engineering): page fetched, no autodiscovery link; `/engineering/rss.xml` redirects to `/research/rss.xml` → 404; `/rss.xml` 404. **No feed found.**

### 1.6 Languages & Runtimes

| Source | Site | Feed URL | Fetched | Cadence (observed) | Cond. | Fits | Terms |
|---|---|---|---|---|---|---|---|
| Go blog | https://go.dev/blog/ | https://go.dev/blog/feed.atom | yes (Atom) | 9 in 90 d | **none** (no `ETag`, `Last-Modified` or `Cache-Control` on the feed) | Concept, Brief (releases) | "Except as noted, the contents of this site are licensed under the Creative Commons Attribution 4.0 License, and code is licensed under a BSD license." (https://go.dev/copyright) |
| Go releases (tags; Go does not publish GitHub Releases) | https://go.dev/doc/devel/release (no feed) | https://github.com/golang/go/tags.atom | yes (Atom) | 10 items in 67 d | E | Brief | GitHub ToS D.5 |
| Rust blog | https://blog.rust-lang.org/ | https://blog.rust-lang.org/feed.xml | yes (Atom) | 10 items in 49 d | E, L | Brief (releases), Concept | "The Rust Programming Language and all other official projects, including this website, are generally dual-licensed: Apache License, Version 2.0 / MIT license" (https://www.rust-lang.org/policies/licenses) |
| Inside Rust | https://blog.rust-lang.org/inside-rust/ | https://blog.rust-lang.org/inside-rust/feed.xml | yes (Atom) | 10 items in 50 d | E, L | Concept (compiler/runtime design) | as Rust blog |
| Python Insider | https://blog.python.org/ | https://blog.python.org/feeds/posts/default | yes (RSS, served as `application/octet-stream`) | 26 in 90 d | E, L | Brief | not checked |
| Python PEPs | https://peps.python.org/ | https://peps.python.org/peps.rss | yes (RSS) | 10 items in 71 d | E, L | Concept, Brief (accepted PEPs, e.g. free-threading) | not checked |
| CPython downloads | https://www.python.org/downloads/ | https://www.python.org/downloads/feed.rss | yes (RSS) | 10 items in 20 d | L | Brief | not checked |
| Inside Java (Oracle Java team) | https://inside.java/ | https://inside.java/feed.xml | yes (Atom) | 10 items in 11 d | E, L | Concept (GC, Loom, Valhalla), Brief (JDK releases) | not checked |
| Node.js blog | https://nodejs.org/en/blog | https://nodejs.org/en/feed/blog.xml | yes (RSS, full archive) | 24 in 90 d | none | Brief (releases, security) | not checked |
| .NET blog | https://devblogs.microsoft.com/dotnet/ | https://devblogs.microsoft.com/dotnet/feed/ | yes (RSS) | 10 items in 23 d | E, L | Concept (runtime/GC), Brief | not checked |
| LLVM blog | https://blog.llvm.org/ | https://blog.llvm.org/index.xml | yes (RSS) | 6 in 90 d | E, L | Concept (compilers) | not checked |
| Kotlin, Swift blogs | https://blog.jetbrains.com/kotlin/ , https://www.swift.org/blog/ | https://blog.jetbrains.com/kotlin/feed/ ; https://www.swift.org/atom.xml | yes ×2 | Kotlin 12 in 56 d; Swift 6 in 90 d | E, L | Brief | not checked |
| Release feeds: Rust, Node.js | https://github.com/rust-lang/rust , https://github.com/nodejs/node | `…/releases.atom` | yes ×2 | Rust 5 in 90 d; Node 10 items in 34 d | E | Brief | GitHub ToS D.5 |

**Checked but weak:** V8 blog `https://v8.dev/blog.atom` fetched, newest 2025-08-04. TypeScript blog `https://devblogs.microsoft.com/typescript/feed/` 1 in 90 d. Zig news `https://ziglang.org/news/index.xml` 1 in 90 d.

### 1.7 Count

| Vertical | Feeds fetched and live (≥1 item in 90 d, or a valid weekday-only feed) | Checked and not fetched / no feed / stale |
|---|---|---|
| Distributed Systems | 13 (11 rows; the release row is 3 feeds) | Jepsen (no feed), ACM Queue (403), USENIX (stale), Uber (406), LinkedIn (no feed), Slack (0 in 90 d), aphyr (0 in 90 d) |
| Databases & Storage | 12 (10 rows; release row is 3 feeds) | MongoDB (stale), CockroachDB blog (no feed), MySQL (403) |
| APIs & Networking | 13 (11 rows; Envoy+nginx and Istio+Cilium rows are 2 each) | gRPC (stale), Chromium (stale) |
| Infrastructure & Cloud | 14 (13 rows) + Grafana, Docker | old Azure feed URL (now HTML) |
| AI/ML Systems | 10 + OpenAI, AWS ML, Ray, llama.cpp | Anthropic engineering (no feed) |
| Languages & Runtimes | 15 (13 rows; Kotlin+Swift and Rust+Node release rows are 2 each) | V8 (stale), TypeScript and Zig (low cadence) |

**[opinion] Volume vs ADR-0005.** These ~80 live feeds showed at least ~1,700 items dated in the last 90 days (a lower bound, since truncated feeds hide older items), dominated by AWS What's New, Azure updates, OpenAI, ClickHouse and NVIDIA. Polled raw they would bury the Author; intake should be a shortlist the Author picks from, the same rule ADR-0005 sets for trending repositories ("a shortlist aid for the Author, never a page generator").

## 2. Public postmortems and incident reports

A postmortem is Case Study material because it is the operator's own account of how the system behaved under failure. Google's SRE book defines it: "A postmortem is a written record of an incident, its impact, the actions taken to mitigate or resolve it, the root cause(s), and the follow-up actions to prevent the incident from recurring." (https://sre.google/sre-book/postmortem-culture/; the book is "Licensed under CC BY-NC-ND 4.0", a non-commercial, no-derivatives licence, so an ad-funded site cannot rely on it to reproduce chapters or figures; quote briefly under section 4 instead.)

Two kinds of source exist and they are not interchangeable **[opinion]**: **status-page incident histories** (timestamps and short updates written during the event: good for a dated Brief, too thin for a Case Study) and **post-incident reviews** (written after closure with contributing factors: Case Study material).

### 2.1 First-party collections

| Publisher | Where | Feed | Fetched | What the publisher commits to (quoted) |
|---|---|---|---|---|
| AWS Post-Event Summaries | https://aws.amazon.com/premiumsupport/technology/pes/ | none found | page yes; lists summaries back to 2011 (e.g. "Summary of the Amazon DynamoDB Service Disruption in the Northern Virginia (US-EAST-1) Region October 19, 2025") | "When an issue has broad and significant customer impact that results in the failure of a significant percentage of control plane API calls, impacts a significant percentage of a service's infrastructure, resources or APIs or is the result of total power failure or significant network failure, AWS is committed to providing a public Post-Event Summary (PES) following the closure of the issue. Post-Event Summaries will remain available for a minimum of 5 years and will provide the issue's scope of impact, factors that contributed to the issue and the actions taken to address the risks identified." AWS Site Terms apply (section 1.1). |
| AWS Health Dashboard | https://health.aws.amazon.com/health/status | https://status.aws.amazon.com/rss/all.rss | yes (RSS, 38 items, 2 in 90 d) | Same page: "the Service History for the preceding 12 months." |
| Azure Post Incident Reviews | https://azure.status.microsoft/en-us/status/history/ | https://azure.status.microsoft/en-us/status/feed/ (current status only) | page yes; feed yes but **0 items at access** (it carries live incidents, not history) | "This page contains Post Incident Reviews (PIRs) of previous service issues, each retained for 5 years. From November 20, 2019, this included PIRs for all issues about which we communicated publicly. From June 1, 2022, this includes PIRs for broad issues as described in our documentation." |
| Google Cloud Service Health | https://status.cloud.google.com/summary | https://status.cloud.google.com/en/feed.atom | yes (Atom, 1 item at access, dated 2026-10-22) | Incident pages on the same host; no retention commitment found on the summary page. |
| Cloudflare blog, "Post Mortem" tag | https://blog.cloudflare.com/tag/post-mortem/ | https://blog.cloudflare.com/tag/post-mortem/rss/ (linked on the page as "Subscribe to Post Mortem RSS feed") | yes (RSS, 20 items 2023-10-20 → 2026-05-01) | Long-form, engineer-written outage write-ups (e.g. "Post mortem on the Cloudflare Control Plane and Analytics Outage", 2023-11-04). Cloudflare Website Terms apply. |
| Cloudflare status | https://www.cloudflarestatus.com/ | https://www.cloudflarestatus.com/history.atom | yes (Atom, includes scheduled maintenance dated ahead) | status updates only |
| GitHub availability reports | https://github.blog/tag/github-availability-report/ | https://github.blog/tag/github-availability-report/feed/ | yes (RSS; monthly, newest "GitHub availability report: September 2026", 2026-10-14) | Monthly summary of incidents with causes, e.g. "In September, we experienced two incidents that resulted in degraded performance across GitHub services." |
| GitHub status | https://www.githubstatus.com/history | https://www.githubstatus.com/history.rss | yes (RSS, 25 in 90 d) | status updates only |
| Fly.io Infra Log | https://fly.io/infra-log/ | https://fly.io/infra-log/feed.xml | yes (Atom, 27 in 90 d) | "This log aims for 100% fidelity to internal incidents, and is a superset both of our status page events and of customer-impacting [incidents]" |
| Slack, Discord, Atlassian status (Statuspage-hosted) | https://slack-status.com/ , https://discordstatus.com/ , https://status.atlassian.com/ | https://slack-status.com/feed/rss ; https://discordstatus.com/history.rss ; https://status.atlassian.com/history.rss | yes ×3 (Slack 9 in 90 d; Discord 18 in 90 d; Atlassian stale, newest 2025-04-30) | status updates only |

**Statuspage pattern [observed + P].** Pages hosted on Atlassian Statuspage (githubstatus.com, discordstatus.com, cloudflarestatus.com, status.atlassian.com) all served `/history.rss` or `/history.atom`. Atlassian documents RSS and Atom as subscription types ("RSS subscribers are not counted against your subscriber total. Atom: Atom subscribers are not counted against your subscriber total.", https://support.atlassian.com/statuspage/docs/how-are-subscribers-counted/) but no Atlassian page found states the `/history.rss` URL pattern; treat it as observed, not guaranteed.

**Engineering-blog postmortems.** Several blogs in section 1 publish occasional incident write-ups (Slack Engineering, Discord, GitHub Engineering, Dropbox Tech, Meta Engineering). They arrive through the normal blog feeds; no separate tag feed was checked for them.

**Third-party collections are indexes, not sources.** Community lists of postmortems (for example the GitHub repository `danluu/post-mortems`) are useful to find old write-ups but are not first-party; cite only the company's own document they link to. Not fetched for this research.

**[opinion] Fit with the Editorial Policy.** A Case Study built on a PES/PIR must label inference ("Reconstruction is labelled"): these documents give contributing factors, not architecture diagrams. The safest Case Study shape is one incident, quoted sparingly, explained against the company's own architecture posts and docs. A same-day outage Brief from a status page is the riskiest Brief shape under section 3 (thin, time-sensitive, everyone else is writing it): only publish one when there is a first-party architecture link that explains *why*.

## 3. Google's policies on summarised, rewritten and aggregated content

All quotes below were copied from the live pages on 2026-10-26. `docs/research/google-news-discover-requirements.md` (2026-10-20) covers News/Discover eligibility, bylines and AI disclosure; this section adds only what bears on **pages built from other publishers' material**.

### 3.1 Spam policies (Search Central, last updated 2026-08-28)

Source: https://developers.google.com/search/docs/essentials/spam-policies

**Scaled content abuse** (verbatim):

> "Scaled content abuse is when many pages are generated for the primary purpose of manipulating search rankings and not helping users. This abusive practice is typically focused on creating large amounts of unoriginal content that provides little to no value to users, no matter how it's created. Examples of scaled content abuse include, but aren't limited to:
> - Using generative AI tools or other similar tools to generate many pages without adding value for users
> - Scraping feeds, search results, or other content to generate many pages (including through automated transformations like synonymizing, translating, or other obfuscation techniques), where little value is provided to users
> - Stitching or combining content from different web pages without adding value
> - Creating multiple sites with the intent of hiding the scaled nature of the content
> - Creating many pages where the content makes little or no sense to a reader but contains search keywords"

**Scraping** (verbatim):

> "Scraping refers to the practice of taking content from other sites, often through automated means, and hosting it with the purpose of manipulating search rankings. Examples of abusive scraping include:
> - Republishing content from other sites without adding any original content or value, or even citing the original source
> - Copying content from other sites, modify it only slightly (for example, by substituting synonyms or using automated techniques), and republish it
> - Reproducing content feeds from other sites without providing some type of unique benefit to the user
> - Creating sites dedicated to embedding or compiling content, such as videos, images, or other media from other sites, without substantial added value to the user"

Note the wording: "**or even** citing the original source". Citing is listed as something abusive scrapers fail to do, not as a cure. The helpful-content page (3.2) says so outright.

**Site reputation abuse** ("site reputation policy"): "applies where third-party content is published on a host site mainly because of that host's already-established ranking signals, which it has earned primarily from its first-party content." "Third-party content is content that's created by an entity that's separate from the established host site. Examples of separate entities include users of that site, freelancers, white-label services, and content created by people not employed directly by the host site." Listed as **not** inconsistent: "Wire service or press release service sites", "News publications that have syndicated news content from other news publications", "Columns, opinion pieces, articles, and other work of an editorial nature". **[inferred]** This policy targets hosting someone else's content for ranking; it does not bite on HastCode's own editorial Briefs. It would matter only if HastCode later published sponsored or partner-written pieces (Editorial Policy, "Independence").

### 3.2 Helpful, people-first content (Search Central, last updated 2026-10-05)

Source: https://developers.google.com/search/docs/fundamentals/creating-helpful-content

"Content and quality questions" (verbatim, the ones that bear on Briefs):

> "Does the content provide original information, reporting, research, or analysis?"
> "Does the content provide insightful analysis or interesting information that is beyond the obvious?"
> "If the content draws on other sources, does it avoid simply copying or rewriting those sources, and instead provide substantial additional value and originality?"
> "Does the content provide substantial value when compared to other pages in search results?"
> "Is the content mass-produced by or outsourced to a large number of creators, or spread across a large network of sites, so that individual pages or sites don't get as much attention or care?"

Warning signs ("Answering yes to some or all of the questions below is a warning sign that you should reevaluate how you're creating content"):

> "Are you using extensive automation to produce content on many topics?"
> "Are you mainly summarizing what others have to say without adding much value?"
> "Are you writing about things simply because they seem trending and not because you'd write about them otherwise for your existing audience?"
> "Does your content leave readers feeling like they need to search again to get better information from other sources?"
> "Are you changing the date of pages to make them seem fresh when the content has not substantially changed?"

What Search Quality raters are trained to assess (verbatim):

> "Effort: The extent to which human work went into creating the content or the systems powering it. For example, writing original analysis, translating a poem manually, or building a custom interactive map for a game represents high effort. Conversely, automatically generating pages from feeds or using generative AI to produce large amounts of text without manual oversight or curation represents little to no effort. **Attribution or giving credit to other sources doesn't replace the need for original effort.**"
> "Originality: The extent to which the content offers unique, original information or perspectives that aren't already available on other websites."

Authorship (already satisfied by the House Byline disclosure, quoted for completeness): "Fabricating creator profiles (such as by using AI-generated headshots, made-up names, or false credentials to make content appear as if it was written by human experts) is a form of deception." **[inferred]** Mira Vance is a *disclosed* house name with no invented credentials or photo (Editorial Policy); the risk would be any drift from that disclosure, for example an Author page that reads as a real individual.

AI disclosure: "Is the use of automation, including AI-generation, self-evident to visitors through disclosures or in other ways?" "Consider adding these when it would be reasonably expected." The Search Central AI post (2023-02-08, https://developers.google.com/search/blog/2023/02/google-search-and-ai-content): "Appropriate use of AI or automation is not against our guidelines. This means that it is not used to generate content primarily to manipulate search rankings, which is against our spam policies." and "Other systems work to elevate original news reporting."

### 3.3 Ranking systems that favour the original over the summary

- Search Central, "A guide to Google Search ranking systems" (last updated 2025-12-10, https://developers.google.com/search/docs/appearance/ranking-systems-guide), **Original content systems**: "We have systems to help ensure we are showing original content prominently in search results, including original reporting, ahead of those who merely cite it." The same page, **Removal-based demotion systems / Legal removals**: "When we receive a significant volume of valid copyright removal requests involving a given site, we are able to use that to demote other content from the site in our results." (Copyright compliance is therefore also a ranking matter; see section 4.)
- Search Central blog, "Understanding news topic authority" (2023-05, https://developers.google.com/search/blog/2023/05/understanding-news-topic-authority): signals include "Influence and original reporting: Our system looks at how original reporting (for example the publisher that first broke a story) is cited by other publishers to understand how a publication is influential and authoritative on a topic", and "Source reputation: … a publication's history doing original reporting or their journalistic awards are strong evidence of positive reputation for news websites." "Publishers with a specific expertise are more likely to surface when their content is most helpful to understand a news story about a specific topic or field."

### 3.4 Google News policies

- "Google News policies" (https://support.google.com/news/publisher-center/answer/6204050, no date shown): to be eligible, "content should follow best practices for article pages and must not violate Google Search's overall policies, Google Search's spam policies, or these policies for Search features". So scaled content abuse and scraping (3.1) are News-eligibility rules too. Misleading content: "We don't allow preview content that misleads users to engage with it by promising details which aren't reflected in the underlying content." Transparency: "Clear dates and bylines; Information about the authors, publication, and publisher; Information about the company or network behind the content; Contact information." Exceptions: "We may make exceptions to these policies based on artistic, educational, historical, documentary, or scientific considerations, or where there are other substantial benefits to the public." "These policies only apply to news surfaces."
- "Ranking within Google News" (https://support.google.com/news/publisher-center/answer/9606702, no date shown): "Google News aims to promote original journalism and expose users to diverse perspectives." Factors: "Relevance of content; Prominence; Authoritativeness; Freshness; Usability; Location and language".
- The eligibility criterion "a consistent history of producing original news-related content" (Search Central blog 2021-07-16) is quoted in `google-news-discover-requirements.md`.
- **No Google page fetched states a numeric limit** (words of quotation, percentage of summary, minimum original words). The tests are qualitative: "original", "substantial additional value", "beyond the obvious".

### 3.5 What this means for Briefs [opinion]

The Brief is the format most exposed to these policies, because by definition it reacts to someone else's announcement ("Curated from public sources; never original reporting", `GLOSSARY.md`). Read together, the texts above make five constraints:

1. **A Brief that restates the announcement is the "mainly summarizing" warning sign**, and citing the source does not fix it ("Attribution or giving credit to other sources doesn't replace the need for original effort"). The Brief's value has to be the architectural "why it matters" and the link into an Explainer, which is exactly the Glossary's definition: keep that section mandatory and longer than the "what happened" section.
2. **Never generate Briefs from feeds.** "Automatically generating pages from feeds … without manual oversight or curation represents little to no effort" and "Reproducing content feeds from other sites without providing some type of unique benefit" are named. Feeds feed a shortlist; the Author picks (ADR-0005's rule for trending repos, applied to Briefs).
3. **Cover only what the existing audience would get anyway** ("writing about things simply because they seem trending" is a warning sign): stay inside the six Verticals.
4. **Original reporting is what news ranking rewards, and HastCode does not do it** (Glossary). Expect Briefs to rank below the first-party announcement they discuss ("ahead of those who merely cite it"). Their job is to route Readers to Explainers and to be the newsletter, not to win news queries.
5. **Volume is a signal in itself** ("many pages", "lots of content on many different topics", "extensive automation"). The Editorial Policy's "We publish only as many pieces as the Author can verify" is the right cap; ADR-0005 already applies it to Profiles.

## 4. Copyright and quoting

**We are not lawyers.** This section states what the primary texts say. Case law (which decides how the factors are actually weighed: for example the CJEU on "quotation", US courts on "transformative" use, Indian courts on "fair dealing") was not reviewed and is where the real limits are set. Questions for a lawyer are in Open questions.

### 4.1 What copyright does not cover (useful for architecture writing)

- **US, 17 U.S.C. §102(b)**: "In no case does copyright protection for an original work of authorship extend to any idea, procedure, process, system, method of operation, concept, principle, or discovery, regardless of the form in which it is described, explained, illustrated, or embodied in such work." (https://www.copyright.gov/title17/92chap1.html)
- **India, Copyright Office handbook**: "Unlike the case with patents, copyright protects the expressions and not the ideas. There is no copyright in an idea." "Copyright does not ordinarily protect titles by themselves or names, short word combinations, slogans, short phrases, methods, plots or factual information. Copyright does not protect ideas or concepts." (https://copyright.gov.in/documents/handbook.html)

**[inferred]** How Netflix shards a datastore, what a PES says failed, what a release changes, are facts and methods that HastCode can explain in its own words and its own diagrams. The company's *prose* and *drawings* are protected expression.

### 4.2 United States: fair use, 17 U.S.C. §107

Statute (verbatim, https://www.copyright.gov/title17/92chap1.html):

> "Notwithstanding the provisions of sections 106 and 106A, the fair use of a copyrighted work, including such use by reproduction in copies or phonorecords or by any other means specified by that section, for purposes such as criticism, comment, news reporting, teaching (including multiple copies for classroom use), scholarship, or research, is not an infringement of copyright. In determining whether the use made of a work in any particular case is a fair use the factors to be considered shall include—
> (1) the purpose and character of the use, including whether such use is of a commercial nature or is for nonprofit educational purposes;
> (2) the nature of the copyrighted work;
> (3) the amount and substantiality of the portion used in relation to the copyrighted work as a whole; and
> (4) the effect of the use upon the potential market for or value of the copyrighted work.
> The fact that a work is unpublished shall not itself bar a finding of fair use if such finding is made upon consideration of all the above factors."

US Copyright Office explanation (https://www.copyright.gov/fair-use/):

- Factor 1: "'transformative' uses are more likely to be considered fair. Transformative uses are those that add something new, with a further purpose or different character, and do not substitute for the original use of the work." Commercial use is weighed, not disqualifying: "This does not mean, however, that all nonprofit education and noncommercial uses are fair and all commercial uses are not fair".
- Factor 2: "using a more creative or imaginative work (such as a novel, movie, or song) is less likely to support a claim of a fair use than using a factual work (such as a technical article or news item)."
- Factor 3: "If the use includes a large portion of the copyrighted work, fair use is less likely to be found; if the use employs only a small amount of copyrighted material, fair use is more likely. … in other contexts, using even a small amount of a copyrighted work was determined not to be fair because the selection was an important part—or the 'heart'—of the work."
- Factor 4: courts consider "whether the use is hurting the current market for the original work (for example, by displacing sales of the original) and/or whether the use could cause substantial harm if it were to become widespread."
- **No safe word count**: "there is no formula to ensure that a predetermined percentage or amount of a work—or specific number of words, lines, pages, copies—may be used without permission." The FAQ: "There are no legal rules permitting the use of a specific number of words, a certain number of musical notes, or percentage of a work." "However, in cases of doubt, the Copyright Office recommends that permission be obtained." (https://www.copyright.gov/help/faq/faq-fairuse.html)

Diagrams are protected works: §101 "'Pictorial, graphic, and sculptural works' include … maps, globes, charts, diagrams, models, and technical drawings, including architectural plans." A diagram is often a whole work in itself, so reproducing it uses 100% of that work (factor 3) **[inferred]**.

### 4.3 European Union

**Directive 2001/29/EC (InfoSoc), Art. 5(3)(d)**, an exception Member States *may* provide (verbatim, via the EU Publications Office, http://publications.europa.eu/resource/celex/32001L0029; EUR-Lex itself answered 202 with an empty body to scripted requests):

> "(d) quotations for purposes such as criticism or review, provided that they relate to a work or other subject-matter which has already been lawfully made available to the public, that, unless this turns out to be impossible, the source, including the author's name, is indicated, and that their use is in accordance with fair practice, and to the extent required by the specific purpose;"

Art. 5(3)(c) also allows "use of works or other subject-matter in connection with the reporting of current events, to the extent justified by the informatory purpose and as long as the source, including the author's name, is indicated, unless this turns out to be impossible". All exceptions are subject to the three-step test, Art. 5(5): "shall only be applied in certain special cases which do not conflict with a normal exploitation of the work or other subject-matter and do not unreasonably prejudice the legitimate interests of the rightholder."

So in the EU, unlike the US: **attribution is a condition of the exception** (source and author's name), the quotation must serve a purpose such as criticism or review, and its length is bounded by "the extent required by the specific purpose". Each Member State implements these in its own law; national texts were not reviewed.

**Directive (EU) 2019/790 (DSM), Art. 15, press publishers' right** (http://publications.europa.eu/resource/cellar/214471fe-786e-11e9-9f05-01aa75ed71a1.0006.03/DOC_1):

> "1. Member States shall provide publishers of press publications established in a Member State with the rights provided for in Article 2 and Article 3(2) of Directive 2001/29/EC for the online use of their press publications by information society service providers.
> The rights provided for in the first subparagraph shall not apply to private or non-commercial uses of press publications by individual users.
> The protection granted under the first subparagraph shall not apply to acts of hyperlinking.
> The rights provided for in the first subparagraph shall not apply in respect of the use of individual words or very short extracts of a press publication."

Para. 3: InfoSoc Articles 5 to 8 (so the quotation exception) "shall apply mutatis mutandis". Para. 4: the right expires "two years after the press publication is published", counted "from 1 January of the year following", and does not apply to press publications first published before 6 June 2019.

Who is a "press publication" (Art. 2(4)): "a collection composed mainly of literary works of a journalistic nature … which: (a) constitutes an individual item within a periodical or regularly updated publication under a single title, such as a newspaper or a general or special interest magazine; (b) has the purpose of providing the general public with information related to news or other topics; and (c) is published in any media under the initiative, editorial responsibility and control of a service provider." Recital 56: "Periodical publications published for scientific or academic purposes, such as scientific journals, should not be covered … Neither should that protection apply to websites, such as blogs, that provide information as part of an activity that is not carried out under the initiative, editorial responsibility and control of a service provider, such as a news publisher." Recital 58: "it is important that the exclusion of very short extracts be interpreted in such a way as not to affect the effectiveness of the rights provided for in this Directive."

**[inferred]** For HastCode's intake the Art. 15 right mostly does not bite, because the sources in section 1 are first-party company blogs, project sites, standards and papers, not news publishers. It *would* apply if a Brief quoted a news outlet (which the Editorial Policy already forbids as a source) or an engineering magazine that qualifies as a press publication (ACM Queue is a candidate: a periodical, but arguably "for scientific or academic purposes"; unresolved). Hyperlinking is expressly outside the right; "very short extracts" has no number in the Directive.

### 4.4 India: fair dealing, Copyright Act 1957 s.52

Section 52(1)(a) (verbatim, Copyright Office of India, https://copyright.gov.in/Copyright_Act_1957/chapter_xi.html):

> "Certain acts not to be infringement of copyright. — (1) The following acts shall not constitute an infringement of copyright, namely,— (a) a fair dealing with any work, not being a computer programme, for the purposes of— (i) private or personal use, including research; (ii) criticism or review, whether of that work or of any other work; (iii) the reporting of current events and current affairs, including the reporting of a lecture delivered in public."

Notes from the primary texts:

- "**not being a computer programme**": the fair-dealing clause excludes software. Quoting code from a repository relies on the repository's licence, not on s.52(1)(a) **[inferred from the text]**. (The licence route is usually open: most sources here publish code under BSD/MIT/Apache, e.g. Go's "code is licensed under a BSD license".)
- s.52(1)(a) as quoted states **no attribution condition**, unlike EU Art. 5(3)(d). The Act does not define "fair dealing".
- s.2(c): "'artistic work' means,— (i) a painting, a sculpture, a drawing (including a diagram, map, chart or plan), an engraving or a photograph, whether or not any such work possesses artistic quality" (https://copyright.gov.in/Copyright_Act_1957/chapter_i.html). Diagrams are protected regardless of artistic merit.
- Handbook: "Copyright in a work is considered as infringed only if a substantial part is made use of unauthorizedly. What is 'substantial' varies from case to case. More often than not, it is a matter of quality rather than quantity. For example, if a lyricist copy a very catching phrase from another lyricist's song, there is likely to be infringement even if that phrase is very short." (https://copyright.gov.in/documents/handbook.html)

### 4.5 Practical limits for HastCode [opinion, built on the texts above]

None of the three regimes gives a safe word count. The defensible house rules are those that satisfy all three texts at once:

| Rule | Why (primary-text hook) |
|---|---|
| Quote a sentence or two at a time, only where the exact words matter (a definition, a commitment, a number's wording), each quote attributed by source and author/organisation and linked | US factor 3 and the "heart" warning; EU Art. 5(3)(d) "source, including the author's name" and "to the extent required by the specific purpose"; EU DSM "very short extracts" |
| Every quote is followed by HastCode's own analysis; a quote never carries a section by itself | US factor 1 "transformative … add something new"; EU "for purposes such as criticism or review"; India s.52(1)(a)(ii); Google "Attribution … doesn't replace the need for original effort" |
| Never reproduce a third-party diagram. Redraw the architecture in Mermaid/SVG from the facts (§102(b): systems and methods are not protected) and link the original figure | Diagrams are protected works in the US (§101) and India (s.2(c)); a whole diagram is 100% of the work; AWS and Cloudflare terms reserve all rights |
| Reuse wholesale only under a licence that allows it, and honour its conditions: Go, Kubernetes docs (CC BY 4.0), Linux Foundation/CNCF site content (CC BY 3.0), arXiv papers per their stated licence. CC BY-NC-ND (Google's SRE book) does **not** permit an ad-funded site's reuse [inferred: "NC" = non-commercial] | Licences quoted in sections 1 and 2 |
| RFCs: quoting a passage still relies on quotation/fair use; TLP §3.c grants reproduction only "in full and without modification" (code components separately under §4) | TLP 5 |
| Briefs: headline facts in our own words + one short quote at most + the link; never the publisher's lede or summary paragraph | US factor 4 (substitution); Google scraping policy ("Republishing content from other sites without adding any original content or value") |
| Code samples from a repository: under the repository's licence with its notice, not under fair dealing | India s.52(1)(a) excludes computer programmes |
| Keep a takedown path (Contact page) and act fast on a valid notice | Google demotes sites with "a significant volume of valid copyright removal requests" (ranking-systems guide) |

## 5. Intake mechanics

### 5.1 Feed formats to support

**What the sources actually serve [observed].** Of the feeds that returned a parseable document on 2026-10-26, **71 were RSS 2.0 and 35 were Atom; none was JSON Feed**, and no fetched feed referenced `jsonfeed`/`application/feed+json`. Several served the wrong media type (Python Insider as `application/octet-stream`, HashiCorp's Atom as `application/rss+xml`, many as `text/xml` or `application/xml`), so detect the format from the root element, not from `Content-Type`.

| Format | Spec | Fields the intake should key on |
|---|---|---|
| RSS 2.0 | https://www.rssboard.org/rss-specification | Item identity: "guid stands for globally unique identifier. It's a string that uniquely identifies the item. When present, an aggregator may choose to use this string to determine if an item is new." `isPermaLink` "is optional, its default value is true". Fall back to `link` when `guid` is absent. Dates in `pubDate` (RFC 822). Channel hints: "ttl stands for time to live. It's a number of minutes that indicates how long a channel can be cached before refreshing from the source"; `skipHours` / `skipDays`: hours/days when aggregators "may not read the channel". [observed: only 3 fetched feeds carried `<ttl>` (5, 60 and 1440 minutes); 3 carried `skipHours`/`skipDays`, including arXiv's Saturday/Sunday.] |
| Atom 1.0 | RFC 4287, https://www.rfc-editor.org/rfc/rfc4287 | "atom:feed elements MUST contain exactly one atom:id element." "atom:entry elements MUST contain exactly one atom:id element." "atom:entry elements MUST contain exactly one atom:updated element." Key on `id`; treat a changed `updated` as an edit, not a new item. |
| JSON Feed 1.1 | https://www.jsonfeed.org/version/1.1/ | `version` and `title` required at top level; item "id (required, string) is unique for that item for that feed over time. If an item is ever updated, the id should be unchanged." `date_published` is RFC 3339. Served as `application/feed+json`. **[opinion]** Cheap to support, but no source in this inventory needs it; defer. |
| RDF/RSS 1.0 | not encountered | none of the fetched feeds used it |

### 5.2 Polling etiquette (HTTP)

From RFC 9110 (HTTP Semantics, https://www.rfc-editor.org/rfc/rfc9110):

- **If-None-Match** (§13.1.2): makes the request conditional on the server "having a selected representation with an entity tag that does not match any of those listed in the field value." "A recipient MUST use the weak comparison function when comparing entity tags for If-None-Match … since weak entity tags can be used for cache validation". (GitHub's `releases.atom` and REST responses send weak `W/"…"` tags [observed].)
- **If-Modified-Since** (§13.1.3): "Transfer of the selected representation's data is avoided if that data has not changed." "A recipient MUST ignore If-Modified-Since if the request contains an If-None-Match header field". "When used for cache updates, a cache will typically use the value of the cached message's Last-Modified header field to generate the field value of If-Modified-Since. This behavior is most interoperable for cases where clocks are poorly synchronized or when the server has chosen to only honor exact timestamp matches". **So store the server's `Last-Modified` string verbatim and send it back; do not synthesise a date from the local clock.**
- **304 Not Modified** (§15.4.5): "there is no need for the server to transfer a representation of the target resource because the request indicates that the client … already has a valid representation". The 304 MUST carry "Content-Location, Date, ETag, and Vary" and "Cache-Control and Expires" if a 200 would have: update the stored ETag from it.
- **Retry-After** (§10.2.3): "Servers send the 'Retry-After' header field to indicate how long the user agent ought to wait before making a follow-up request." Value is an HTTP-date or delay-seconds.
- **429 Too Many Requests** (RFC 6585 §4, https://www.rfc-editor.org/rfc/rfc6585): "indicates that the user has sent too many requests in a given amount of time ('rate limiting')" and "MAY include a Retry-After header".

**Observed support** (column "Cond." in section 1): most feeds send `ETag` and/or `Last-Modified`; I verified a real **304** from `kubernetes.io/feed.xml` (If-None-Match), `blog.cloudflare.com/rss/` (If-Modified-Since) and `github.com/kubernetes/kubernetes/releases.atom` (If-None-Match). **The Go blog feed sends neither validator nor `Cache-Control`**, and neither do Netflix/Medium, Dropbox, ClickHouse, Redis, Stripe, Tailscale, Google Cloud, Azure updates, Shopify, vLLM and Node.js: for those, every poll is a full download, so poll them less often.

**[opinion] Polling rules for the intake command** (it would run locally, like `make refresh`, never in GitHub Actions per CLAUDE.md / ADR-0005):

1. One request per feed per run, sequential (no fan-out), identifying `User-Agent` with a contact URL (the `/contact/` page).
2. Persist per feed: final URL after redirects, `ETag`, `Last-Modified` (verbatim), last status, last-seen item ids. Send `If-None-Match` and `If-Modified-Since` together (servers that understand both use the former per §13.1.3).
3. Honour `<ttl>`, `skipHours`/`skipDays`, `Retry-After`, 429 and 503; on repeated failure back off exponentially and report the feed as broken rather than retrying.
4. Run at most daily; the Author's weekly rhythm (ADR-0005: Refresh Monday, roundup Brief Friday) needs no more. Daily is far gentler than any `ttl` observed except arXiv's.
5. Output a dated shortlist (title, source, link, Vertical guess) for the Author. Never write Markdown pieces from it (section 3.5).
6. Update stored feed URLs on permanent redirects: four sources had moved (MongoDB, vLLM, Azure updates, Uber), and Azure's old feed URL now lands on HTML with a 200, which only a format check catches.

### 5.3 GitHub REST API for releases

The repo already fetches GitHub data locally (`cmd/fetch`, `internal/github`, `internal/factpack`; ADR-0005). The same token and client can watch releases.

**Endpoints** (https://docs.github.com/en/rest/releases/releases):

- `GET /repos/{owner}/{repo}/releases`: "Information about published releases are available to everyone. Only users with push access will receive listings for draft releases." `per_page` "max 100", default 30.
- `GET /repos/{owner}/{repo}/releases/latest`: "View the latest published full release for the repository. The latest release is the most recent non-prerelease, non-draft release, sorted by the created_at attribute. The created_at attribute is the date of the commit used for the release, and not the date when the release was drafted or published." **[inferred]** So `latest` can lag a release published today from an older commit; for a Brief, list releases and compare `published_at`.
- Projects that publish tags, not Releases: Go is the main case here. `GET /repos/golang/go/tags?per_page=3` returned `weekly.2012-03-27`, `weekly.2012-03-22`, `weekly.2012-03-13` [observed]: the tags endpoint is not newest-first, so it is unusable for "what's new" on golang/go. Use the Go blog feed or `https://go.dev/dl/?mode=json` (returned `go1.27.3` and `go1.26.10`, both `stable: true` [observed]; I did not find a page documenting this endpoint).

**Rate limits** (https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api):

- "The primary rate limit for unauthenticated requests is 60 requests per hour." Authenticated (personal access token): "your personal rate limit of 5,000 requests per hour."
- Secondary limits: "No more than 100 concurrent requests are allowed"; "No more than 900 points per minute are allowed for REST API endpoints"; "No more than 90 seconds of CPU time per 60 seconds of real time is allowed."
- On a limit: "If you exceed your primary rate limit, you will receive a 403 or 429 response, and the x-ratelimit-remaining header will be 0. You should not retry your request until after the time specified by the x-ratelimit-reset header." For secondary limits: honour `retry-after`; "Otherwise, wait for at least one minute before retrying. If your request continues to fail due to a secondary rate limit, wait for an exponentially increasing amount of time between retries, and throw an error after a specific number of retries."

**Conditional requests are free only when authenticated** (https://docs.github.com/en/rest/using-the-rest-api/best-practices-for-using-the-rest-api): "Making a conditional request does not count against your primary rate limit if a 304 response is returned and the request was made while correctly authorized with an Authorization header." [observed, unauthenticated: a 200 took `x-ratelimit-remaining` from 60 to 59 and the follow-up 304 took it to 58, consistent with the doc.] The same page: "Poll only as often as you need to, on a fixed schedule. If a response includes an x-poll-interval header, wait at least that many seconds before you poll the same endpoint again." and "You should subscribe to webhook events instead of polling the API for data" (webhooks need a receiving server, which ADR-0002/0004 rule out, so polling it is) **[inferred]**.

**[inferred] Budget.** ~25 watched repositories × 1 request per daily run = 25 requests/day against 5,000/hour; with ETags most will be 304s that cost nothing. The REST response also carries `cache-control: public, max-age=60` [observed], so polling the same repository more than once a minute is pointless.

**[opinion] REST vs `releases.atom`.** ADR-0005 already treats the REST API as "the sanctioned path" and cites GitHub's Acceptable Use Policies on scraping. `releases.atom` is convenient (no token, one URL per repo, newest first, real 304s) but it is a github.com web endpoint. Using the REST endpoint inside the existing fetch command keeps intake under the same token, rules and code as Fact Packs; feed polling stays for non-GitHub sources.

## Open questions

1. **Is polling `github.com/{owner}/{repo}/releases.atom` "scraping" under GitHub's Acceptable Use Policies?** ADR-0005 quotes the AUP definition ("Scraping does not refer to the collection of information through our API"). The Atom feeds are public, linked by GitHub and answer conditional GETs, but they are not the API. Ask GitHub or use REST (section 5.3).
2. **Legal review of the house quoting rules (4.5)** by someone qualified in US, EU (at least one Member State's implementation of Art. 5(3)(d) and Art. 15) and Indian law, since the audience spans all three. In particular: does the EU quotation exception cover a commercial, ad-funded explainer; is an engineering magazine such as ACM Queue a "press publication" or a periodical "for scientific or academic purposes"; how Indian courts read "fair dealing" for commercial commentary.
3. **Per-source licence audit.** Only Go, Kubernetes, CNCF/LF, Rust, PostgreSQL, AWS, Cloudflare, IETF, arXiv, USENIX, GitHub and the SRE book were checked. Medium's ToS returned 403 (Netflix, Airbnb, Pinterest are on Medium). Before a Case Study leans on a source, read its terms.
4. **ACM Queue / CACM, MySQL and Uber** blocked scripted requests (403/406). Check by hand in a browser whether they publish feeds a polite fetcher may use, or drop them from automated intake.
5. **Feedless high-value sources** (Jepsen analyses, AWS Builders' Library, USENIX proceedings, Anthropic engineering, CockroachDB blog, SQLite release notes) have no feed. Decide whether the Author checks them by hand on a schedule, or whether a page-change check is acceptable under each site's terms.
6. **Does the Glossary need a word for the shortlist?** The intake output (dated candidates the Author picks from) is a new artefact. ADR-0005 calls trending discovery "a shortlist aid". If intake is built, name it in `GLOSSARY.md` (domain-modeling) and decide where it lives (`.scratch/`, local only?).
7. **`go.dev/dl/?mode=json`** returned Go's current stable releases but I found no page documenting it; confirm before depending on it.
8. **Statuspage `/history.rss`/`/history.atom`** URLs worked on four status pages, but Atlassian's docs found here only list RSS/Atom as subscriber types, not the URL pattern.

## Sources (accessed 2026-10-26)

**Feeds.** Every feed URL in the section 1 and 2 tables was fetched on 2026-10-26; status, format, item count and validators are recorded there and are not repeated here.

**Repository context**
- `GLOSSARY.md`, `content/pages/editorial-policy.md`, `docs/adr/0005-repo-profiles-as-an-explainer-kind-with-a-fetched-fact-pack.md`, `site.yaml`, `docs/research/google-news-discover-requirements.md`

**Google**
- https://developers.google.com/search/docs/essentials/spam-policies (last updated 2026-08-28): scaled content abuse, scraping, site reputation policy
- https://developers.google.com/search/docs/fundamentals/creating-helpful-content (last updated 2026-10-05)
- https://developers.google.com/search/docs/appearance/ranking-systems-guide (last updated 2025-12-10): original content systems, legal removals
- https://developers.google.com/search/blog/2023/02/google-search-and-ai-content (2023-02-08)
- https://developers.google.com/search/blog/2023/05/understanding-news-topic-authority (2023-05)
- https://support.google.com/news/publisher-center/answer/6204050 (Google News policies; no date shown)
- https://support.google.com/news/publisher-center/answer/9606702 (Ranking within Google News; no date shown)

**Copyright: statutes and official guidance**
- US Title 17, chapter 1 (§§101, 102, 107), US Copyright Office: https://www.copyright.gov/title17/92chap1.html (uscode.house.gov returned "Page Not Found" while "operating in a temporary format")
- US Copyright Office Fair Use Index: https://www.copyright.gov/fair-use/
- US Copyright Office FAQ, fair use: https://www.copyright.gov/help/faq/faq-fairuse.html
- Directive 2001/29/EC (Art. 5(3)(c), 5(3)(d), 5(5)), EU Publications Office: http://publications.europa.eu/resource/celex/32001L0029 (EUR-Lex https://eur-lex.europa.eu/legal-content/EN/TXT/HTML/?uri=CELEX:32001L0029 answered 202 with no body to scripted requests)
- Directive (EU) 2019/790 (Art. 2(4), Art. 15, Recitals 56 and 58), EU Publications Office: http://publications.europa.eu/resource/cellar/214471fe-786e-11e9-9f05-01aa75ed71a1.0006.03/DOC_1 (CELEX 32019L0790, OJ L 130, 17.5.2019)
- India, Copyright Act 1957, Chapter XI (s.52): https://copyright.gov.in/Copyright_Act_1957/chapter_xi.html ; Chapter I (s.2(c)): https://copyright.gov.in/Copyright_Act_1957/chapter_i.html
- India, Copyright Office handbook: https://copyright.gov.in/documents/handbook.html

**Site terms and licences**
- https://go.dev/copyright
- https://kubernetes.io/blog/ (footer)
- https://www.linuxfoundation.org/legal/terms (linked as "Terms of Use" from https://www.cncf.io/blog/)
- https://www.rust-lang.org/policies/licenses
- https://www.postgresql.org/about/licence/
- https://aws.amazon.com/terms/
- https://www.cloudflare.com/website-terms/
- https://trustee.ietf.org/license-info/ and https://trustee.ietf.org/documents/trust-legal-provisions/tlp-5/
- https://info.arxiv.org/help/license/index.html ; https://info.arxiv.org/help/rss.html
- https://www.usenix.org/about
- https://docs.github.com/en/site-policy/github-terms/github-terms-of-service (D.5)
- https://sre.google/sre-book/postmortem-culture/ (definition and CC BY-NC-ND 4.0 notice)
- https://policy.medium.com/medium-terms-of-service-9db0094a1e0f (403, not read)

**Postmortem and status pages**
- https://aws.amazon.com/premiumsupport/technology/pes/
- https://health.aws.amazon.com/health/status
- https://azure.status.microsoft/en-us/status/history/
- https://status.cloud.google.com/summary
- https://blog.cloudflare.com/tag/post-mortem/
- https://github.blog/tag/github-availability-report/
- https://www.githubstatus.com/history
- https://fly.io/infra-log/
- https://support.atlassian.com/statuspage/docs/how-are-subscribers-counted/

**Feed and HTTP specifications**
- RSS 2.0: https://www.rssboard.org/rss-specification
- Atom: RFC 4287, https://www.rfc-editor.org/rfc/rfc4287
- JSON Feed 1.1: https://www.jsonfeed.org/version/1.1/
- HTTP Semantics: RFC 9110, https://www.rfc-editor.org/rfc/rfc9110 (§10.2.3, §13.1.2, §13.1.3, §15.4.5)
- 429 Too Many Requests: RFC 6585, https://www.rfc-editor.org/rfc/rfc6585 (§4)

**GitHub REST API**
- https://docs.github.com/en/rest/releases/releases
- https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api
- https://docs.github.com/en/rest/using-the-rest-api/best-practices-for-using-the-rest-api
- Observed: `https://api.github.com/repos/kubernetes/kubernetes/releases?per_page=1`, `https://api.github.com/repos/golang/go/tags?per_page=3`, `https://go.dev/dl/?mode=json`

**Pages checked for feeds (no feed found or blocked)**
- https://jepsen.io/analyses ; https://www.linkedin.com/blog/engineering ; https://www.cockroachlabs.com/blog/ ; https://www.anthropic.com/engineering ; https://aws.amazon.com/builders-library/ ; https://sre.google/books/ ; https://queue.acm.org/ (403) ; https://cacm.acm.org/feed/ (403) ; https://blogs.oracle.com/mysql/ (403) ; https://www.uber.com/en-US/blog/engineering/ (406)

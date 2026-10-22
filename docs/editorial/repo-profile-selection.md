# Choosing repositories for Repo Profiles

A Repo Profile (see `GLOSSARY.md`) evaluates one open-source repository for the Reader: a mid-to-senior backend engineer deciding whether to adopt it. We profile a handful a week, chosen by an Author; we never profile everything that trends. Decided 2026-10-22 (ADR-0005).

## Floors

A repository qualifies only if every line holds. If one is in doubt, skip it; there is always another.

1. **OSI-approved license.** GitHub reports the SPDX id in the Fact Pack. `NOASSERTION` means GitHub could not classify the LICENSE file: read it by hand (Bun, for example, is MIT plus bundled notices).
2. **A commit on the default branch in the last 90 days.**
3. **Fits exactly one Vertical.** If it straddles two, the Author places it and the Profile says why. If it fits none (agent plugin directories, "awesome" lists, how-to repositories, courses), it is not a Profile candidate however hot it is.
4. **At least 1,000 stars, or clearly rising**: star growth in the last 30 days that is a large share of its total (the shortlist records the evidence and its URL).
5. **Not already profiled.** A new major release of a profiled repository is a Refresh and, if the architecture changed, a Review; not a second Profile.
6. **Not primarily a UI kit, mobile library, game or consumer app.** The Reader runs backend systems.

## Cadence and mix

- Three to five Profiles a week, alongside five Explainers of other kinds and five Briefs.
- Spread across Verticals over a month rather than chasing one hot area.
- Mix rising repositories (traffic now) with evergreen anchors (search intent that lasts: "X vs Y", "is X production ready", "X alternatives"). Roughly three rising to two anchors.
- Prefer repositories whose users are platform and backend teams. Our ad-network target (Raptive) needs half the site's readers in the US, UK, Canada, Australia and New Zealand; GitHub's own population is about a quarter Tier-1, so audience shape is part of the choice.

## Where candidates come from

- The weekly shortlist: GitHub search (`stars:>1000 pushed:>DATE`, `created:>DATE sort=stars`), per-language and per-topic trending pages read by a person, Hacker News front page and Show HN. The all-languages trending front page is mostly agent-tooling noise; query per language and per topic.
- Never scrape github.com programmatically, never pool tokens: the REST API under one token is the sanctioned path (GitHub Acceptable Use Policies; `docs/research/github-trending-geography-and-data-access.md`).
- Reader submissions and maintainer requests, once a contact route exists.

## First five (2026-10-22)

| Repository | Vertical |
|---|---|
| apache/iggy | distributed-systems |
| xerj-org/xerj | databases-storage |
| NVIDIA/OpenShell | infrastructure-cloud |
| deepseek-ai/DeepInfer | ai-ml-systems |
| php/frankenphp | languages-runtimes |

Alternate: cilium/cilium (apis-networking). Reserves: malisper/pgrust, agentgateway/agentgateway, pranshuparmar/witr. Shortlist with evidence: `.scratch/hastcode-content-plan/shortlist-2026-10-22.md`.

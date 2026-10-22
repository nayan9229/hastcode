# go-articles

Content site: tech news + architecture explainers, SEO/AEO-optimised, ad-monetised. Greenfield.

## Agent skills

### Issue tracker

Local markdown: issues live as files under `.scratch/<feature>/`. See `docs/agents/issue-tracker.md`.

### Triage labels

Default five-role vocabulary (`needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`). See `docs/agents/triage-labels.md`.

### Domain docs

Single-context: `GLOSSARY.md` + `docs/adr/` at repo root. See `docs/agents/domain.md`.

## Repo Profiles

A Fact Pack is fetched, never hand-written: `make fetch REPO=owner/name MD=content/explainers/<vertical>/<slug>.md` writes it beside the Markdown (needs `GITHUB_TOKEN`; README "Repo Profiles: fetching a Fact Pack"). `make refresh` re-fetches every one and exits 3 when any changed. Nothing publishes without the Author: a Refresh is a pull request a human merges, and Publish stays a tag (ADR-0004, ADR-0005).

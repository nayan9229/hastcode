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

A Fact Pack is fetched, never hand-written: `make fetch REPO=owner/name MD=content/explainers/<vertical>/<slug>.md` writes it beside the Markdown (needs `GITHUB_TOKEN`; README "Repo Profiles: fetching a Fact Pack"). `make refresh` re-fetches every one; the fetch binary exits 3 when any changed (make reports it as `Error 3`). A Refresh is manual: the Author runs it locally, reviews the diff, commits and pushes; GitHub Actions never fetches from GitHub (it tests, builds Previews and publishes by tag). Nothing publishes without the Author's approval: an agent never removes `draft: true` from a Profile, merges a Profile pull request, or tags. An agent never sets `verified:` on any Piece either: it records that the Author checked every claim against its Source (ADR-0007). The Author does that after verifying every claim; Publish stays a tag (ADR-0004, ADR-0005).

Writing one: start from `docs/editorial/repo-profile-template.md`, follow `docs/editorial/repo-profile-writing-guide.md` (sections, grounding in the Fact Pack / README / release pages, quickstart verbatim, title patterns), choose with `docs/editorial/repo-profile-selection.md`, voice in `docs/editorial/house-byline-voice.md`. Day-28 numbers go in a copy of `docs/editorial/checkpoint-note-template.md`.

## Writing a Case Study, Concept, Comparison or Brief

Start from the kind's template, `docs/editorial/{case-study,concept,comparison,brief}-template.md` (the skeleton `make lead` writes), and follow its writing guide, `docs/editorial/<kind>-writing-guide.md` (first-party Sources, every claim linked, short credited quotes, diagrams redrawn in Mermaid, code only under its licence, sections, titles, what the build checks); voice in `docs/editorial/house-byline-voice.md`. A Brief also needs a Lead at most 3 days old, 250 words at most, and a "Why it matters" that is our own analysis.

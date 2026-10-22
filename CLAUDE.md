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

A Fact Pack is fetched, never hand-written: `make fetch REPO=owner/name MD=content/explainers/<vertical>/<slug>.md` writes it beside the Markdown (needs `GITHUB_TOKEN`; README "Repo Profiles: fetching a Fact Pack"). `make refresh` re-fetches every one and exits 3 when any changed. Nothing publishes without the Author's approval, for Profiles and Refresh pull requests alike: an agent never removes `draft: true` from a Profile, merges a Profile or Refresh pull request, or tags. The Author does that after verifying every claim; a Refresh is a pull request a human merges, and Publish stays a tag (ADR-0004, ADR-0005).

Writing one: start from `docs/editorial/repo-profile-template.md`, follow `docs/editorial/repo-profile-writing-guide.md` (sections, grounding in the Fact Pack / README / release pages, quickstart verbatim, title patterns), choose with `docs/editorial/repo-profile-selection.md`, voice in `docs/editorial/house-byline-voice.md`. Day-28 numbers go in a copy of `docs/editorial/checkpoint-note-template.md`.

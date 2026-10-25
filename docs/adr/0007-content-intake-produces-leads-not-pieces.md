---
status: accepted
date: 2026-10-26
---

# Content intake produces Leads, not Pieces

HastCode needs a steady supply of Case Studies, Concepts, Comparisons and Briefs to reach the launch gate (ADR-0001: ≥30 Explainers across all six Verticals, ≥20 Briefs) and then sustain 5 Explainers and 5 Briefs a week. Research (`docs/research/content-source-intake.md`) found supply is not the limit: 77 verified first-party feeds produced ~1,700 items in 90 days. The limits are Google's spam policies, which name "scraping feeds … to generate many pages" and say "attribution … doesn't replace the need for original effort", copyright on diagrams, and the Author's verification time. We decided:

- **Intake lists Leads; it never writes a Piece.** A local `make intake` checks the committed **Source Watchlist** (`sources.yaml`: 6 Sources per Vertical plus 11 incident feeds to start, each with a default Vertical and kind) with conditional GETs, and GitHub releases through the REST API with the existing token. It lists new Leads with the Source's suggested kind and Vertical; no AI classification. The Lead queue and seen-state are local and gitignored.
- **The Author picks; drafting is a separate step.** `make lead ID=…` writes a draft skeleton (frontmatter, Source link, the kind's empty headings, no prose); a Claude session drafts into it with the house-byline voice guide and the kind's writing guide. At most two Leads a day are picked; undecided Leads dismiss after 7 days.
- **Briefs come from Leads at most 3 days old**, and stand on their "why it matters" analysis, never on the Source's summary.
- **Diagrams are always redrawn** in Mermaid from the facts, with the original linked; a Source's images are never copied.
- **Nothing ships unverified.** Production builds require `verified: <date>` on every non-draft Piece published after 2026-10-26. Only the Author sets it, after checking every claim against its Source; an agent never does (CLAUDE.md). The ten earlier Repo Profiles are outside the date rule until the review queue is cleared.
- **Format before pipeline.** One Case Study, one Concept, one Comparison, then one Brief linking one of them, are written by hand from real Leads before any Intake code is built; each kind's template and writing guide is drawn from its piece.

## Considered options

- **Auto-draft Pieces from feed items**: the shape Google's scaled-content and feed-scraping policies name, and a verification load no solo Author can carry.
- **Scheduled Intake in GitHub Actions**: rejected to keep Actions deploy-only (tests, Previews, Publish by tag), as for the Repo Profile Refresh.
- **Follow all 77 verified feeds**: ~19 Leads a day, a backlog by the end of the first week.
- **Polling `github.com/…/releases.atom`**: a web endpoint that may count as scraping under GitHub's Acceptable Use Policies; the REST API is the sanctioned path.

## Consequences

- The `verified:` gate is new build behaviour: a production build fails on a new non-draft Piece without it.
- Brief volume is bounded by Explainer volume, since every Brief must link an Explainer.
- The Watchlist is a reviewed, versioned file; adding a Source is a commit with a reason.

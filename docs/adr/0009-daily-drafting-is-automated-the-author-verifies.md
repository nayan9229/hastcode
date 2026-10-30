---
status: accepted
date: 2026-10-30
amends: ADR-0007
---

# Daily drafting is automated; the Author verifies

Picking Leads and drafting by hand took the Author's mornings, and Briefs expire 72 hours after their Lead. The Author decided that everything up to verification runs from one command, and that the daily Publish cap is lifted until the site reaches 500 Pieces. We decided:

- **`make daily` picks and drafts; it never verifies or publishes.** It runs Intake, then a headless Claude Code session judges the fresh Leads: a Brief needs a real event with an architectural cause and a live Explainer that teaches the lesson; an Explainer needs a real system, one named idea or a genuine choice, not marketing, opinion or a release list. Every pick goes through `make lead`, then one session drafts it by the kind's writing guide and a second fact-checks it against its Sources. Drafts stay `draft: true` and are committed locally with the Author's check list; nothing is pushed or tagged.
- **The Author's part is unchanged in substance:** verify every claim, set `verified:`, remove `draft: true`, Publish by tag (ADR-0004, ADR-0007, CLAUDE.md).
- **No daily cap until 500 Pieces** (amending ADR-0007's two a day). The run has a per-run Explainer limit for cost, not pacing.

## Considered options

- **Agent shortlists, the Author picks:** kept the Author's judgement on every pick, but costs a step every morning and loses Briefs to the 72-hour window.
- **Keep the two-a-day Publish cap:** the safer path against Google's scaled-content policy, which applies from the first page, not after 500. The Author chose speed to the launch gate and accepts that risk until 500 Pieces.
- **Run on GitHub Actions on a schedule:** CLAUDE.md keeps Actions from fetching Sources, and it would need its own Claude credentials and budget.

## Consequences

- Each run costs Claude usage roughly in proportion to the Pieces drafted; the run prints its total.
- The Author's verification is now the only human check between a feed item and a live page; its quality bounds the site's.
- Leads a draft rejects stay marked picked in the local intake state.

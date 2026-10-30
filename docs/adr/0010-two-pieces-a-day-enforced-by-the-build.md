---
status: accepted
date: 2026-10-30
amends: ADR-0009
---

# Two Pieces a day, enforced by the build

ADR-0009 lifted the daily Publish cap until 500 Pieces. After 32 Explainers went live on 28 and 29 October, the Author restored it: Google's scaled-content policy applies from the first page, and a new domain that grows in bursts of AI-drafted pages is the pattern it names. We decided:

- **At most two Pieces a day**, Explainers and Briefs together, counted by `published` day. A production build fails (rule `publish-pace`) when a day from 2026-10-30 on has more; a Preview only warns. The days before 2026-10-30 are exempt, so nothing live comes down.
- **Drafting keeps pace.** `make daily` drafts up to two Briefs (they expire after 72 hours) and up to two Explainers, none while four or more Explainer drafts are waiting.
- **The backlog waits as drafts.** The build has no scheduling: a non-draft Piece is live on the next Publish whatever its date. On release day the Author lifts `draft: true` and sets `published` and `verified` to that day (the verified gate refuses a `verified` earlier than `published`).

The rest of ADR-0009 stands: picking and drafting are automated, and the Author verifies.

## Considered options

- **The cap as policy only, as on 27 October:** it was broken within a day.
- **Count in `scripts/publish.sh`** (Pieces newly live since the last tag): it misses a Publish from another machine, and a build rule is testable at the Build seam like every other gate.

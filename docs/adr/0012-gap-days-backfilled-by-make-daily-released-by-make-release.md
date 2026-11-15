---
status: accepted
date: 2026-11-15
amends: ADR-0009, ADR-0010
---

# Gap Days are backfilled by `make daily`; `make release` is the one release command

The two-a-day cap (ADR-0010) is counted per Publication Day, and the build has no scheduling. When the Author skips `make daily` for a day, that day stays empty for good, and the site's growth shows the gap. We decided:

- **A Gap Day is a Publication Day since 2026-10-30 with fewer than two Pieces, live or drafted for it.** `make daily` looks back at most three days (yesterday and the two before) and fills each Gap Day it finds, oldest first, with Explainers only: a draft's `published` is set to that day at the pick's clock time, and the review queue says which day it fills. Older gaps stay empty.
- **Never a Brief.** A Brief is news: its Lead must be under 72 hours old and the news sitemap lists 48 hours. A Brief dated two days back and first live today is stale and its date misleads. Today's Briefs are drafted as before.
- **At most six drafts a run**: today's picks first (up to two Briefs and two Explainers), then Gap Days, two Explainers each, until six. Gap Days left unfilled wait for the next run while inside the window. The four-waiting rule stands: while four or more Explainer drafts wait, no Explainers are drafted, Gap Days included.
- **`make release` releases drafts.** The Author runs it after verifying every draft in the review queue. With no `FILES=` it takes every draft; it first counts Pieces per Publication Day, live plus released, and refuses if any day would exceed two, naming the files. Then, per Piece, it removes `draft: true`, sets `verified` to today in the Piece's own offset, sets `reviewed` to now on Explainers, keeps `published` as the draft carries it, runs the production build, and commits. An agent never runs it (CLAUDE.md). The morning is `make daily`, verify, `make release`, `git push`, `scripts/publish.sh cf`.
- **The dating trade-off is accepted and bounded.** A backfilled Explainer shows the Gap Day as its date though it first went live later; `verified` and `reviewed` carry the real day. Explainers are evergreen, the lag is at most three days, and the Author chose this over leaving days empty, having done it by hand twice (2026-11-12).

## Considered options

- **Leave gaps empty**: honest dates, but the site shows every missed morning and the cap can never be caught up.
- **Raise the cap on catch-up days**: ADR-0010's whole point is steady growth; four live Pieces dated one day is the burst it forbids.
- **Backfill Briefs too**: rejected, above.
- **Unbounded look-back**: a week away would draft fourteen Explainers nobody can verify.

## Consequences

- `published` on a draft is no longer "provisional, set it when you publish": it is the Publication Day `make daily` assigned, and `make release` keeps it.
- Release commits read "Release <day>: <titles>, verified by the Author", as before, but are made by `make release`.

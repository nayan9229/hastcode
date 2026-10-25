# 06: Draft the real Source Watchlist

**What to build:** The Author has a reviewed starting Watchlist. Draft it from the research tables: 6 Sources per Vertical, chosen for first-party signal and steady cadence, plus the 11 incident feeds (kind case-study), each with a one-line reason. Prefer feeds that sent validators; avoid the stale or blocked ones the research lists. Run one real `make intake` against it, fix any Source that fails, and record the run summary (Leads per Vertical) in the ticket.

Spec: `../spec.md` ("Source Watchlist"). Research: `docs/research/content-source-intake.md` (tables per Vertical, incident feeds, terms notes).

**Blocked by:** 02 (Intake tracer), 03 (polite fetching), 04 (GitHub Sources)

**Status:** ready-for-agent

- [ ] Watchlist has 36 Vertical Sources (6 each) plus the incident feeds, every one validated and fetched successfully in one real run
- [ ] Each Source carries a reason; any swap from the research's top picks is explained
- [ ] The run summary is pasted into this ticket
- [ ] Ends `ready-for-human`: the Author reviews the list before it is committed; the agent does not commit it unreviewed

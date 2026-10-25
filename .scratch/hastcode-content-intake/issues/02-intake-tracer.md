# 02: Intake tracer: Watchlist to Leads

**What to build:** The Author runs `make intake` and sees what is new. A committed Source Watchlist (YAML at the repo root) lists Sources, each with name, site, a feed URL (RSS or Atom) or a GitHub repository, Vertical, kind and reason; it is validated at the start of every run (unknown Vertical or kind, both or neither of feed and repository, duplicate names fail the run). Intake fetches each feed Source, lists each item not seen before as a Lead (title, Source, date, link, the Source's kind and Vertical, age) in a local Markdown queue grouped by Vertical, newest first, and records seen-state so a second run lists nothing new. Queue and state live in a gitignored local directory; Intake writes nothing under the content tree. GitHub repository Sources are accepted by validation but fetched in ticket 04 (skipped with a note here).

Spec: `../spec.md` ("Source Watchlist", "Intake"). Research: `docs/research/content-source-intake.md` section 5. Prior art: the Fact Pack fetcher and its in-process fake GitHub tests. Ship a small example Watchlist used by tests; the real one is ticket 06.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] New Intake package with run (and stubs to extend for dismiss/pick) behind a thin command and `make intake`
- [ ] Intake-seam tests with `httptest` servers serving RSS 2.0 and Atom fixtures, temp state dir and fixed clock: new Leads listed with the Source's kind and Vertical; second run lists nothing new; Watchlist validation errors; nothing written under the content tree
- [ ] Lead IDs are stable across runs and short enough to type
- [ ] The local state directory is gitignored; README gains an Intake section (commands, where state lives)

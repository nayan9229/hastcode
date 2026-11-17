# 01: Self-edit checklist in the morning flow

**What to build:** Every draft gets a prose self-edit before the Author checks its facts, and the Author sees that checklist every morning without going looking for it. Today the house-voice drafting pipeline reads "draft in this voice → verify every claim → edit → publish", and nothing says what "edit" means. The Author's verification checks facts, and nothing checks the prose (`docs/research/article-writing-craft.md` §5.2 item 3). The research note's one-page checklist (§4) fills that gap, and §5.3 item 3 gives the pipeline wording. Breakdown approved 2026-11-17.

- The §4 checklist becomes an editorial doc beside the writing guides (`docs/editorial/self-edit-checklist.md`). It keeps the §4 groups (Point, Top of the page, Structure, Sentences, Words, Teaching, Evidence, Our contribution, Ending, Last pass) and each line's wording. It says which lines do not apply to a Brief: a Brief has no hook, nut graf, TL;DR, diagram or worked example, and only the AI-assistance note comes before "What happened". Its 250-word cap (`brief-length`) replaces the per-H2 line, and "Why it matters" is its contribution. It also says which lines work differently in a Repo Profile. Its FAQ `###` questions are the one place question headings are allowed. Its quickstart stays verbatim. A number the Fact Pack carries is pointed at in Signals, not retyped (the Repo Profile guide's "Do not retype what the Fact Pack carries"). The nut-graf line applies to Concepts and Case Studies (ticket 03); in a Comparison, the opening's closing question does that job. It says plainly that it checks the writing, not the facts, and that the Author's verification against Sources is unchanged (ADR-0007, ADR-0009).
- The drafting pipeline in `docs/editorial/house-byline-voice.md` gains the step before verification: "→ draft in this voice → self-edit with the checklist → verify every claim → edit →" (§5.3 item 3), linking the new doc.
- `make daily`'s drafter already reads the house-voice guide. Its instructions also name the checklist, so each draft is self-edited before the fact-check session runs.
- The review queue that `make daily` writes (`.scratch/daily/<date>/review-queue.md`) carries the checklist for every drafted Piece. Each draft gets its own unticked self-edit boxes next to "Drafter asks you to check" and "Fact-checker could not confirm". That can be the full list, or one box per checklist group plus a link to the doc; the entry for a Brief shows only the lines that apply to a Brief. Pieces listed under "Not drafted" get no boxes.
- The boxes are the Author's. `make daily` writes them unticked, and no agent ticks them for the Author. Ticking them is not verification and never stands in for `verified:` (which only `make release`, run by the Author, writes; CLAUDE.md, ADR-0012).

**Blocked by:** None (can start immediately)

**Status:** resolved

- [x] `docs/editorial/self-edit-checklist.md` exists, carries every §4 line, marks the lines that don't apply to Briefs (hook, nut graf, TL;DR, diagram, worked example, per-H2 length) and the lines that work differently in Repo Profiles, and links back to the research note
- [x] `house-byline-voice.md`'s drafting pipeline has the self-edit step before "verify every claim" and links the checklist. Each of the five writing guides' "The flow" lists a self-edit step after "Draft" and before "The Author verifies", linking the checklist (step numbers move up by one)
- [x] `make daily`'s drafting instructions name the checklist
- [x] Test: a review queue written from a sample run (one Explainer and one Brief drafted, one Lead not drafted) has unticked self-edit boxes, or a box per group plus the checklist link, under each drafted Piece; the Brief's entry leaves out the lines that don't apply to Briefs; the not-drafted entry has none. `make test` runs this test. Today nothing tests the review queue, and `make test` runs only the Go tests. The review queue's writer is not Go, so the test needs a way in. It also runs in GitHub Actions, which run `make test`, so it uses only fixture files: it must not fetch from GitHub or call Claude
- [x] The review queue's opening paragraph says the self-edit comes before verification and does not replace it
- [x] Runbook §4 Step 2 ("verify") lists the self-edit boxes as a fourth thing the queue shows each draft (after the file, "Drafter asks you to check" and "Fact-checker could not confirm"). It says in plain words to tick them while reading the draft, before checking claims against the Source. Runbook §6 Option B's "Draft … following …" prompt also names the checklist
- [x] README's daily-rhythm paragraph says the Author's check list now carries self-edit boxes for each draft, and links the checklist doc
- [x] No agent ticks a self-edit box, removes `draft: true`, sets `verified:` or runs `make release`
- [x] `make test` passes

## Comments

**Resolved, 2026-11-17.** Merged into `integration/indexing-writing-craft` as 7fb0cdde1f4d9a2148068b1ae4aa32fe576eb769; code-review fixes in 131d137. `make test` passes. Nothing left over. The self-edit boxes stay the Author's: no agent ticks them.

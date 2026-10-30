# 01: make daily picks and drafts; the Author verifies

**What to build:** one command each morning that does everything up to verification (ADR-0009, grilled 2026-10-30). Intake; a headless Claude Code session chooses Briefs (fresh Lead, real event, a live Explainer teaching the lesson) and up to N Explainers (real system, one idea, a genuine choice), explaining rejections; `make lead` per pick, serially, renamed to the chosen slug; parallel draft + fact-check sessions per Piece; a check list; a local commit. Never verifies, pushes or tags.

**Blocked by:** None

**Status:** resolved

- [x] `make daily`, `DRY=1`, `BRIEFS=`, `EXPLAINERS=`, `PARALLEL=`
- [x] A Brief is refused unless every Explainer it links is live
- [x] Not-viable drafts removed; failed sessions reported, not fatal
- [x] Check list in `.scratch/daily/<date>/review-queue.md`; run cost printed
- [x] Sessions run with `--permission-mode dontAsk` and an allow-list: read, edit, web fetch, and `scripts/daily/check-page.sh` only
- [x] README "The daily rhythm"; ADR-0009; ADR-0007 marked amended

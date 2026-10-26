# 01: The verified gate

**What to build:** Nothing new ships unchecked. Frontmatter on every Piece gains an optional `verified` date. A production build fails (rule `unverified`, naming the file) on any non-draft Piece whose `published` is after 2026-10-26 when `verified` is missing, earlier than `published`, or after the build clock. Preview builds report the same condition as a warning. Drafts and Pieces published on or before 2026-10-26 (the ten Repo Profiles) are outside the rule.

Spec: `../spec.md` ("Verification gate"). ADR-0007. CLAUDE.md already says an agent never sets `verified:`; tests use fixtures only.

**Blocked by:** None (can start immediately)

**Status:** resolved

- [x] Build-seam tests: new non-draft Piece without `verified` fails production naming the file; valid date passes; date before `published` fails; future date fails; draft unaffected; Piece published on or before 2026-10-26 unaffected; preview build warns instead of failing
- [x] Works for every kind of Explainer and for Briefs
- [x] Real `content/` still builds in production (the ten Profiles predate the rule)
- [x] README documents the field and the rule

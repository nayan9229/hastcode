# 06: Refresh workflow

**What to build:** Every Monday at 06:00 UTC (and on manual trigger) a GitHub Actions workflow checks out `main`, runs the fetch command in refresh mode, and stops quietly when nothing changed. When something changed it runs the full check target, commits the changed Fact Packs to a `refresh/<date>` branch and opens a pull request titled with the date whose body lists each changed repository and what changed; if that branch and pull request already exist from an earlier run, it updates them instead of opening a second. The workflow never tags and never touches prose. The README explains the Refresh, says that the job runs the tests itself because pull requests opened with the repository token trigger no workflows, and records the first manual run.

Spec: `../spec.md` ("Refresh workflow", User Stories 31–35). Glossary: Refresh. ADR-0004: Publish stays by tag.

**Blocked by:** 05 (Fetch command)

**Status:** ready-for-agent

- [ ] Workflow has the Monday schedule and a manual trigger; workflow lint passes
- [ ] Unchanged tree: job ends green with no branch, no pull request
- [ ] Changed tree: tests run before the pull request exists; branch and pull request created with the date title and the per-repository summary in the body
- [ ] Re-run on the same day with further changes updates the existing branch and pull request (no duplicate)
- [ ] The token used for the API is the repository token; no other secret is introduced
- [ ] README section written; first manual run performed and its result noted there (ready-for-human step inside this ticket)

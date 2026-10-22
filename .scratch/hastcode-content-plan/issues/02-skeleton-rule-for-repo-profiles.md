# 02: Skeleton rule for Repo Profiles

**What to build:** The build enforces the Repo Profile shape so every Profile reads the same. The Markdown body must carry exactly these H2 headings in this order: TL;DR, In 60 seconds, Architecture, Use it if, Don't use it if, Alternatives, FAQ, Verdict, Sources. A missing, extra or reordered heading fails the build with a named rule and a message that says which heading was expected where. A hand-written "Signals" or "Changelog" heading fails with a message saying those sections render from the Fact Pack. The FAQ section must contain five to eight H3 questions. The existing TL;DR and Sources rules keep applying.

Spec: `../spec.md` ("Page composition and gates", User Stories 3, 10, 22). Prior art: the Brief skeleton gate and the Explainer TL;DR/Sources rules.

**Blocked by:** 01 (Repo Profile kind renders from a Fact Pack)

**Status:** resolved

- [x] A Profile with the headings in order builds clean; each of missing, extra and reordered fails with the rule name and the expected heading
- [x] "Signals" or "Changelog" in the body fails with the "rendered from the Fact Pack" message
- [x] FAQ with four or nine H3s fails; five and eight pass; an H2 under FAQ is an extra heading
- [x] The rule applies only to `repo-profile`; Case Study, Concept and Comparison bodies are untouched by it
- [x] The valid Repo Profile fixture from 01 passes the rule unchanged

## Implementation notes

Built test-first at the builder seam: eight tests in `builder/repo_profile_skeleton_test.go` (one per criterion, plus the TL;DR/Sources carry-over and the all-kinds guard). `make check` passes. `buildtest.ValidRepoProfileBody` needed no change: it already followed the spec's order with five FAQ questions.

- **Step.** `builder/step_validate_repo_profile_skeleton.go`, registered as validate **32** "validate repo profile skeleton", between the Explainer skeleton (30) and the Brief skeleton (35). It loops `s.explainers` and returns early for every kind but `content.RepoProfile`, so Case Study, Concept and Comparison bodies are untouched. The Explainer rules (`tldr`, `sources`, `diagram`) still run on a Profile at 30; a Profile missing TL;DR or Sources is reported under both that rule and `repo-profile-sections`.
- **Headings (the one spelling).** Constants `headingIn60Seconds` "In 60 seconds", `headingArchitecture` "Architecture", `headingUseItIf` "Use it if", `headingDontUseItIf` "Don't use it if" (straight apostrophe; matching is exact, goldmark has no typographer), `headingAlternatives` "Alternatives", `headingFAQ` "FAQ", `headingVerdict` "Verdict", `headingSignals` "Signals", `headingChangelog` "Changelog", beside the existing `headingTLDR` and `headingSources`. `repoProfileHeadings` is the ordered body: TL;DR, In 60 seconds, Architecture, Use it if, Don't use it if, Alternatives, FAQ, Verdict, Sources. `repoProfileFactSections` maps each fact section to the written section it follows on the page (Signals -> Architecture, Changelog -> Alternatives); ticket 03's section splits should read it rather than respell the names. If ticket 03 declared constants of the same names on its branch, keep these and drop its copies at merge.
- **Rules and messages.**
  - `repo-profile-sections`: `missing ## Verdict (section 8, after ## FAQ); a Repo Profile has exactly these H2 sections in this order: TL;DR, In 60 seconds, …, Sources` (one per missing heading; section 1 has no "after"); `unexpected ## Roadmap (section 7); …` (one per heading outside the shape, a second copy of a heading, or an H2 written under FAQ, which starts a new section); when the set is right but the order is not, `## Don't use it if is section 4, where ## Use it if belongs; …` (the first heading out of place). The order check never fires when a missing or unexpected heading was reported (same shape as `brief-sections`).
  - `repo-profile-fact-section`: `## Signals is rendered from the Fact Pack (the build places it after ## Architecture): delete the hand-written section` (and Changelog / Alternatives). A fact section is excluded from the order check, so it gets this one error, never an "unexpected" one too.
  - `repo-profile-faq`: `## FAQ has 4 questions, want 5–8 (one ### heading per question, phrased the way a Reader searches)`; counts `*ast.Heading` level 3 among the FAQ section's blocks; skipped when there is no FAQ section (the sections rule reports that). Constants `minFAQQuestions` 5, `maxFAQQuestions` 8.
- **What ticket 07's template must match.** Exactly the nine H2s above, in that order, nothing else at H2; a 3–5-bullet TL;DR; a Mermaid fence (under Architecture, the stable diagram slot); five to eight `###` questions under FAQ, each followed by its answer; a Sources list with at least one link; no Signals or Changelog heading. Body H1s are ignored (as `Sections()` ignores them), so the template should not use one.
- **Docs.** README: the Conventions step list (Repo Profile skeleton 32), a Codebase map row for the step, and a "Repo Profiles" paragraph under Conventions beside "Briefs".

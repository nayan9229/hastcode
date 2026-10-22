# 02: Skeleton rule for Repo Profiles

**What to build:** The build enforces the Repo Profile shape so every Profile reads the same. The Markdown body must carry exactly these H2 headings in this order: TL;DR, In 60 seconds, Architecture, Use it if, Don't use it if, Alternatives, FAQ, Verdict, Sources. A missing, extra or reordered heading fails the build with a named rule and a message that says which heading was expected where. A hand-written "Signals" or "Changelog" heading fails with a message saying those sections render from the Fact Pack. The FAQ section must contain five to eight H3 questions. The existing TL;DR and Sources rules keep applying.

Spec: `../spec.md` ("Page composition and gates", User Stories 3, 10, 22). Prior art: the Brief skeleton gate and the Explainer TL;DR/Sources rules.

**Blocked by:** 01 (Repo Profile kind renders from a Fact Pack)

**Status:** ready-for-agent

- [ ] A Profile with the headings in order builds clean; each of missing, extra and reordered fails with the rule name and the expected heading
- [ ] "Signals" or "Changelog" in the body fails with the "rendered from the Fact Pack" message
- [ ] FAQ with four or nine H3s fails; five and eight pass; an H2 under FAQ is an extra heading
- [ ] The rule applies only to `repo-profile`; Case Study, Concept and Comparison bodies are untouched by it
- [ ] The valid Repo Profile fixture from 01 passes the rule unchanged

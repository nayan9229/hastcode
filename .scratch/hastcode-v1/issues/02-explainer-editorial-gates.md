# 02: Explainer editorial gates, draft mode and Review warning

**What to build:** The Author cannot publish an Explainer below the editorial floor. The build fails with a clear, page-attributed message when an Explainer lacks a `## TL;DR` section (3–5 bullets immediately after the intro), lacks a `## Sources` section at the end with at least one link, contains no Mermaid fence, is a Comparison without a Markdown table under `## Trade-offs`, or shares a slug with another piece. The build warns (not fails) when any H2 section exceeds 250 words or when `reviewed` is more than 90 days old. `draft: true` pieces are excluded from production output and included in preview builds. Nothing is written for a failed page.

**Blocked by:** 01 (Walking skeleton)

**Status:** resolved

- [x] Each gate has an invalid fixture and a Build-seam test asserting the specific error and that no output page is written for it
- [x] Errors name the content file and the rule broken; multiple failures are reported together, not one at a time
- [x] Warnings are reported separately from errors and do not fail the build
- [x] Duplicate slug across Verticals and across Explainer/Brief namespaces both detected
- [x] Draft fixture absent from production output, present in preview output, and absent from any listing in production
- [x] Review warning triggers at >90 days and not at 89 days (fixture dates relative to a clock the builder accepts)

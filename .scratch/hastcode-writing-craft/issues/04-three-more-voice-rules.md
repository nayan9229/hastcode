# 04: Three more voice rules

**What to build:** The house voice gains three rules that drafters and the Author can hold a Piece to. They are text only: no build check. Breakdown approved 2026-11-17.

- **One name per thing** (`docs/research/article-writing-craft.md` §5.2 item 7, wording §5.3 item 6): "One name per thing. Pick the Source's term and use it every time; mention a synonym once, in brackets, if Readers search for it." Microsoft's global writing tips say the same. The Case Study guide already learned it for one case ("Terms drift between Sources… use the Source's own term").
- **At most one analogy** (§5.2 item 6, wording §5.3 item 8): "Analogies: at most one per piece, and say where it stops being true." The note marks this rule as [opinion], based on Evans, Google's developer style guide and Orwell (§7). The rule text says it is a house choice, not a research finding.
- **Numbers carry context** (§5.2 item 12, wording §5.3 item 10, added to voice rule 7 "Show the receipts"): "Every number carries its unit, its date or version, and what it is compared with." The Comparison guide already asks for "a measurement quoted with its caveat", and BBC 3.4.8 asks the same.
- All three go into `docs/editorial/house-byline-voice.md`. Ticket 02 edits the same rule list (rule 5), which is why this one waits for it.
- **Repo Profiles keep their exception.** The Repo Profile guide says not to retype numbers the Fact Pack carries ("point at it" or "say what it means"). The number rule must not push a drafter to copy stars or release counts into prose. The rule, or the Repo Profile guide, says that pointing at Signals meets it.

**Blocked by:** 02 (Prose warnings in the build), which edits the same voice-rule list

**Status:** resolved

- [x] `house-byline-voice.md` has the one-name rule and the analogy rule as new numbered rules after rule 8, and rule 7 has the numbers sentence. Wording as §5.3 items 6, 8 and 10, edited only to fit the list
- [x] Rule 5 as ticket 02 left it is unchanged, and the list's numbering stays consistent
- [x] The analogy rule is labelled a house choice
- [x] The Repo Profile guide's "Do not retype what the Fact Pack carries" and the number rule agree (one line in either doc says how)
- [x] If the self-edit checklist (ticket 01) exists, its "One name per thing", analogy and "Every number has a unit…" lines say the same as the new rules
- [x] No build rule, test or content change. README and the runbook are unchanged, since nothing the Author runs behaves differently

## Comments

**Resolved, 2026-11-17.** Merged into `integration/indexing-writing-craft` as ed30c045a681f87ea79918f77cc083d0ffbbe401; code-review fixes in 131d137. `make test` passes. Nothing left over.

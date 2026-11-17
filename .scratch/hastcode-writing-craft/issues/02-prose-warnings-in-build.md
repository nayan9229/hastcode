# 02: Prose warnings in the build

**What to build:** The build points at long sentences, filler words and question headings the same way it already points at long sections: as warnings that never fail a build. The house-voice rules turn into numbers the build can check. Voice rule 5 in `docs/editorial/house-byline-voice.md` says "Short sentences" with no number, and its banned list has three words. Across our Explainers, about one sentence in eight runs past 35 words, and Briefs run as long on average (`docs/research/article-writing-craft.md` §5.2 items 2, 8 and 9; draft wording in §5.3 items 2, 7 and 9). Breakdown approved 2026-11-17.

- **Voice rule 5** reads (§5.3 item 2): "Short sentences: about 20 words on average. Split anything over 30 unless it is a deliberate list. Use the shortest sentences at the hardest points." Its banned list grows (§5.3 item 7) to: revolutionary, game-changing, magic, simply, just (as in "just add"), easy, obviously, "it is important to note", "please note".
- **Three warnings**, in the validate phase next to the `section-length` warning (validate 40), with the same shape: a warning naming the content file, the rule and where the problem is, and never an error. Production and preview builds both report them, so `make release`, `make dev` and `make daily`'s page check all show them.
  - `sentence-length`: a sentence over 30 words. The message gives the H2 it sits under, the sentence's first few words and its word count, so the Author can find it.
  - `banned-word`: one of the banned words or phrases, named in the message with its H2. "Just" warns only in the filler sense ("just add", "just run"), never in "just as", "just under", "just before" or "just in time".
  - `question-heading`: an H2 that ends in "?", in any Piece. A Repo Profile's FAQ questions are `###` under `## FAQ` and never warn (the guide requires them, §5.2 item 9).
- **What is never checked:** frontmatter, fenced code (Mermaid included), inline code, tables, blockquotes (the AI-assistance note among them), link URLs and the Fact Packs. Signals and Changelog are rendered from the Fact Pack after validation, so they are never seen. A banned word inside quotation marks is a Source's words and never warns. Words inside quotation marks still count toward a sentence's length, because the Reader reads them. Confirmed by the Author on 2026-11-17: quoted words count toward length but never trigger `banned-word`, and `question-heading` checks H2s only. A sentence-length warning on a sentence that quotes a Source is fixed by splitting our own words around the quote, never by rewording the quote. Headings are checked only by `question-heading`. Static pages (About, Contact, Editorial Policy, Privacy) are not Pieces and are not checked.
- **Scope:** every Explainer kind and every Brief.
- **`make daily` picks them up.** The drafter is already told to fix every build warning naming its file, so new drafts arrive with fewer long sentences. Its instructions say that a sentence that is a deliberate list may stay long, and that a Source's quoted words are never reworded to clear a warning.
- **Expect many warnings.** Rough baseline measured on 2026-11-17 over the 74 Pieces in `content/` (13 Briefs, 61 Explainers, no drafts): 7,266 sentences, mean 18.4 words, median 17. 912 sentences are over 30 words (13%), in 69 of the 74 Pieces. Banned words outside quotes: "simply" 10, "easy" 7, "just" 25 uses, most of them not the filler sense ("just as", "just under"). None of revolutionary, game-changing, magic, obviously, "it is important to note" or "please note". Question H2s: 0; all 64 question headings are Repo Profile FAQ `###`. This came from a throwaway script with its own sentence splitter, and the build's counts will differ. The research note's own count (mean 21.7) used a different splitter. That's fine: they are warnings, not failures, and nothing published needs to change for this ticket.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] Voice rule 5 in `house-byline-voice.md` carries the numbers and the extended banned list, worded as §5.3 items 2 and 7
- [ ] Build-seam tests, one per rule: a 31-word sentence warns `sentence-length` and a 30-word one does not; each banned word or phrase warns `banned-word`; "just as" and "just under" do not, and neither does a banned word inside quotation marks; an H2 ending in "?" warns `question-heading` in a Concept and in a Brief; Repo Profile FAQ `###` questions do not; nothing inside a code fence, a table, a blockquote or a link URL warns. Each case also checks that the build still succeeds (a warning, never an error)
- [ ] The filled kind templates and the Repo Profile template still build clean: their tests also check that none of the three new rules fires
- [ ] A preview build over today's `content/` succeeds. Its count per rule is written in this ticket (a Comments note) next to the baseline above, with a line on why they differ
- [ ] README: the pipeline-step order list names the new validate step(s) beside "section-length warning 40". README has no paragraph on the section-length warning today, so add a short one (beside the Briefs and Repo Profiles rule text) saying what each of the three rules checks and what it skips
- [ ] The Concept, Case Study, Comparison and Repo Profile guides' "What the build checks" list the three rules under Warnings, next to `section-length`. The Brief guide has no Warnings list (`section-length` does not check Briefs), so it gains one naming the three rules
- [ ] A build with roughly 900 `sentence-length` warnings still shows its errors clearly. The builder's last line already counts errors and warnings apart ("build failed: N error(s), M warning(s)"). Check that, with one error added to today's `content/`, the Author can find that error without scrolling through every warning. If not, say in the runbook how to list only the errors
- [ ] The runbook says in plain words what the three warnings mean, that they never block `make release` or a Publish, and what to do (split the sentence, cut the word, restate the heading), next to the "If `make build` fails" table
- [ ] Each writing guide's "The skeleton" gains one line (§5.3 item 9, confirmed by the Author on 2026-11-17): "Only Repo Profile FAQ headings may be questions; every H2 is a statement."
- [ ] `make test` passes

## Comments

**Implementer (W02), 2026-11-17:** Built on `impl/w02` (`builder/step_validate_prose.go`, validate 41; prose extraction `internal/markdown/prose.go`; tests `builder/prose_warnings_test.go`).

Preview build over today's `content/` (`go run ./cmd/builder -mode preview -skip-search-index`, 74 Pieces): succeeds, 0 errors, 937 warnings.

| Rule | Build | Baseline above |
|---|---|---|
| `sentence-length` | 918 warnings in 69 Pieces | 912 sentences over 30 words in 69 Pieces |
| `banned-word` | 18 in 17 Pieces: "simply" 10, "easy" 7, "just check" 1 | "simply" 10, "easy" 7, "just" 25 uses (most not filler) |
| `question-heading` | 0 | 0 question H2s |

Why they differ: the build splits sentences its own way (at `.`, `!` or `?` followed by a space and a capital, digit, quote, bracket or code; never after "e.g.", "i.e.", "vs." or an initial), counts an inline code span as one word and skips tokens with no letter or digit (a lone dash). The throwaway script split differently, so six more sentences cross 30 words here. "just" warns only before an instruction verb (a fixed list: add, run, use, check, …), so only one of the 25 uses is filler ("Just check first", github-agent-scale-git). One "easy" is in "easy-to-use", which warns because a hyphen is a word boundary.

Errors among the warnings: with one error added to a copy of `content/` (a bad `source_url`), the error is line 1 of 939 lines of output, so it scrolls off the top; the last line does count it ("build failed: 1 error(s), 936 warning(s)"). The runbook now gives `make build 2>&1 | grep -E '^(error|build failed)'` to list only the errors.

Brief: a question H2 in a Brief warns `question-heading`, but the Brief skeleton's `brief-sections` already fails any H2 other than its two literal ones, so that build fails on `brief-sections` (never on `question-heading`). The test checks exactly that.

A block that still holds a `TODO(author)` hint is skipped: it is a skeleton's placeholder, which `author-placeholder` reports (the picked-skeleton test in `internal/intake` expects only that warning).

# 03: Guide wording: word-count ceiling and "what you'll get" sentence

**What to build:** Two wording changes to the Explainer guides, so that drafts stop padding toward a length and every Concept and Case Study says early what the Reader will get. Breakdown approved 2026-11-17.

- **A ceiling, not a target.** The Concept, Case Study and Comparison writing guides each say "Aim for 2,000 to 2,800 words". Google lists "writing to a particular word count" as a sign of search-engine-first content, and its raters mark down "filler" (`docs/research/article-writing-craft.md` §5.2 item 1). `docs/research/aeo-citation-signals.md` found word count has about zero correlation with AI Overview citation. In all three guides the sentence becomes the §5.3 item 1 wording: "Most land between 2,000 and 2,800 words. This is a ceiling to check, not a target: stop when the point is made. Google has no preferred word count, and padding is the 'filler' its raters mark down." The "each H2 under 250 words (`section-length` …)" part of that line stays.
- **The "what you'll get" sentence (nut graf).** The Concept and Case Study openings end on a turn, and no sentence says what the piece will give the Reader before the AI-assistance note and TL;DR (§5.2 item 4). Both guides' "Opening (before TL;DR)" gain the §5.3 item 4 sentence: "End the hook with one sentence that says what this piece will show the Reader (the 'nut graf'): 'This is how that loop forms, and the three defences that break it.'" The Comparison guide already ends its opening "on the question the piece answers", so it does not change.
- **Templates show it.** The hook placeholder in the Concept and Case Study templates asks for that closing sentence. `make lead` writes its skeleton from these templates, so every new Concept or Case Study skeleton shows it too.
- **No build change.** The build checks no Explainer word count today, as a target or a floor. The only word limits are the `section-length` warning (250 words per H2) and the Brief's 250-word cap (`brief-length`), and both are ceilings. `make daily` has one place that reads as a floor: its selection step asks the chooser to confirm that a Source "can carry 2,000+ words". That is reworded to judge whether the Source has enough mechanism (components, data flow, numbers) for the kind. The page-size advice given to the drafter ("about 2,400 words" with one diagram) stays, because it is a ceiling set by the page-size limit.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] The Concept, Case Study and Comparison guides carry the ceiling wording; no editorial doc still says "Aim for" a word count
- [ ] The Concept and Case Study guides' "Opening (before TL;DR)" ends with the nut-graf sentence, with its example
- [ ] The Concept and Case Study templates' hook placeholder asks for the closing "what this piece will show the Reader" sentence. A `make lead` skeleton of each kind shows it (checked by a test or by running `make lead` on a real Lead and discarding the skeleton)
- [ ] The kind-template tests still pass: filled templates build clean under every gate, and unfilled hints still fail a production build with `author-placeholder`
- [ ] `make daily`'s selection step no longer asks for "2,000+ words"; the drafter's page-size advice is unchanged
- [ ] Runbook and README checked: neither presents 2,000 to 2,800 words as a target (today neither does; leave them alone unless a line does)
- [ ] `make test` passes

# Self-edit checklist

Use this after the draft is done and **before** the claim-by-claim check against Sources. It checks the writing, not the facts. The Author's verification is unchanged: every claim, number, quote and diagram is still checked against its Source, and only the Author sets `verified:` (ADR-0007, ADR-0009). Ticking a line here is not verification and never stands in for it.

The lines come from the one-page checklist in [`docs/research/article-writing-craft.md` §4](../research/article-writing-craft.md#4-one-page-self-edit-checklist), word for word; that section says where each line comes from. Write in [the house byline voice](house-byline-voice.md) and follow your kind's writing guide; this list does not replace either.

## Where it is used

- **The drafter** (a `make daily` session, or you or an agent drafting by hand) self-edits with this list before handing the draft on for verification. It is the "self-edit with the checklist" step in [the drafting pipeline](house-byline-voice.md#drafting-pipeline) and in each writing guide's "The flow".
- **The Author** gets the same boxes, unticked, under every drafted Piece in the morning's review queue (`.scratch/daily/<date>/review-queue.md`, written by `make daily`). Tick them while you read the draft, before you check its claims against the Source. They are your boxes: no agent ticks them for you.

`make daily` reads the boxes from [the checklist](#the-checklist) below, one per `- [ ]` line. Keep each line on one line, and keep the notes in the trailing `*(…)*` form, so the review queue stays right.

## By kind

Each line applies to every kind unless it carries a note:

- *(not in a Brief)*: the line does not apply to a Brief, and a Brief's entry in the review queue leaves it out.
- *(Brief only)*: the line replaces one of those for a Brief, and only a Brief's entry shows it.
- *(Brief: …)*, *(Comparison: …)*, *(Repo Profile: …)*: the line works differently in that kind; the note says how.

**Brief.** A Brief has no hook, nut graf, TL;DR, diagram or worked example: only the AI-assistance note comes before "What happened" (`brief-sections`). Its 250-word cap (`brief-length`) replaces the per-H2 line, and "Why it matters" is its contribution: our own analysis, never the Source restated ([the Brief writing guide](brief-writing-guide.md)).

**Concept and Case Study.** Every line applies. The nut-graf line is the sentence that ends the hook and says what this piece will give the Reader.

**Comparison.** Every line applies. The opening's closing question, the one the piece answers, does the nut graf's job ([the Comparison writing guide](comparison-writing-guide.md)).

**Repo Profile.** Its FAQ `###` questions are the one place question headings are allowed. Its quickstart in "In 60 seconds" stays verbatim. A number the Fact Pack carries is pointed at in Signals, not retyped ("Do not retype what the Fact Pack carries", [the Repo Profile writing guide](repo-profile-writing-guide.md#ground-every-claim)). `make daily` does not draft Repo Profiles, so their notes never reach the review queue.

## The checklist

Copy it into the ticket or pull request when you draft by hand: `SELF-EDIT: {slug}`.

### Point

- [ ] I can say the point in one sentence: ______________________________
- [ ] The piece covers one idea; anything else is a Related link.

### Top of the page

- [ ] Title: key term in the first two words, plain words, no pun, under ~65 chars.
- [ ] Description: adds the "so what", has a verb, doesn't repeat the title, <=150 chars.
- [ ] Hook: one real, linked moment (system, number, failure), not a definition or history. *(not in a Brief)* *(Repo Profile: one sentence, a sourced fact about the repository; a number the Fact Pack carries is pointed at in Signals, not retyped.)*
- [ ] Nut graf: the hook's last sentence says what this piece gives the Reader. *(not in a Brief)* *(Comparison: the opening's closing question, the one the piece answers, does this job.)* *(Repo Profile: the hook is one sentence; the TL;DR that follows says what the Reader gets.)*
- [ ] TL;DR: the answer in plain words; a reader who stops here is not misled. *(not in a Brief)* *(Repo Profile: its three bullets: what it is, the design choice that sets it apart, who should adopt it today.)*
- [ ] Only the AI-assistance note comes before "What happened". *(Brief only)*

### Structure

- [ ] Read only the headings: do they tell the story in order?
- [ ] No heading says "Background", "Overview", "Introduction" or "Conclusion".
- [ ] No question headings outside a Repo Profile FAQ. *(Repo Profile: its FAQ `###` questions, phrased the way the Reader types them, are the one place question headings are allowed.)*
- [ ] Each paragraph has one topic, stated in its first sentence.
- [ ] Each H2 is under 250 words; cut, don't pad. *(not in a Brief)*
- [ ] The whole Brief is 250 words at most, headings and the AI-assistance note included (`brief-length`); cut, don't pad. *(Brief only)*

### Sentences

- [ ] Average about 20 words or fewer; every sentence over ~30 words is split or earns its length.
- [ ] Active voice with a named actor ("the client retried", not "requests were retried").
- [ ] Hardest parts use the shortest sentences.

### Words

- [ ] Every technical term is defined at first use, and kept where it is the right term.
- [ ] One name per thing: the same term every time (no synonyms for the same component).
- [ ] No hype or filler: revolutionary, game-changing, magic, simply, just, easy, obviously, "it is important to note", "please note", "there is/are" openings.
- [ ] At most one analogy, and it says where it stops being true.
- [ ] No idioms or pop-culture references a non-native reader would miss.

### Teaching

- [ ] A concrete example comes before the abstract statement.
- [ ] At least one worked example with real numbers (Concept, Case Study, Comparison). *(not in a Brief)*
- [ ] Every "what" has a "why" or a cost.
- [ ] The common confusion is named and corrected ("what it is not").
- [ ] The text walks the diagram box by box; the diagram is never left to speak alone. *(not in a Brief)*
- [ ] Code: introduced by a sentence, realistic, under its licence; README commands verbatim. *(Repo Profile: the quickstart in "In 60 seconds" is the README's commands verbatim: same commands, same order, same flags; never improved, shortened or merged.)*

### Evidence

The writing side; the Author's verification checks the facts.

- [ ] Every claim has its link where the claim is made.
- [ ] Every number has a unit, a date or version, and its caveat. *(Repo Profile: a number the Fact Pack carries (stars, license, release dates, contributors) is pointed at in Signals, not retyped.)*
- [ ] Quotes are exact, short, credited, and followed by our own analysis.
- [ ] Each inference says "we infer" in its own sentence.

### Our contribution

- [ ] Name what this Piece adds that no Source says: ____________________ *(Brief: "Why it matters" is the contribution: our own analysis, never the Source restated.)*
- [ ] Delete every paragraph that only restates a Source.

### Ending

- [ ] The last section tells the Reader what they can now do, measure or decide.
- [ ] No new idea in the last paragraph; no recap of what they just read.

### Last pass

- [ ] Read it aloud once. Fix every place you stumbled or got bored.
- [ ] Cut big (sections, paragraphs), then small (words). Then cut 10% more.

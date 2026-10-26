# Writing a Case Study

A Case Study (`GLOSSARY.md`) reconstructs how one real, named system is built, from its makers' public material, for the Reader: a mid-to-senior backend engineer. It is the flagship kind. This guide and [the template](case-study-template.md) are drawn from the first one, [Single-tenant Titus: how Netflix gives one container its own VM](../../content/explainers/infrastructure-cloud/netflix-single-tenant-titus.md) (ticket 07). Write in [the house byline voice](house-byline-voice.md) and hold every sentence to the [Editorial Policy](../../content/pages/editorial-policy.md).

## The flow

1. **Pick a Lead.** Run `make intake`, pick from the queue (at most two Leads a day), then `make lead ID=<id>` (add `KIND=case-study` or `VERTICAL=…` to override the Source's suggestion). The skeleton is this template, with the Lead's title, kind and Vertical filled in and the Lead's link in Sources. To start by hand, copy [the template](case-study-template.md) to `content/explainers/<vertical>/<slug>.md` on a branch. The slug names the system in lower case (`netflix-single-tenant-titus`). Keep `draft: true`.
2. **Source notes.** Read the Lead in full, then the earlier first-party posts it builds on. Note each fact with its link and anchor before writing a sentence.
3. **Diagram first.** Redraw the system's path in Mermaid from the text. Its boxes, top to bottom, become the H2s that walk it.
4. **Draft** into the skeleton. AI assistance is fine; it is a draft, not a source.
5. **Write the check list.** In the ticket or pull request, list every claim you are unsure of: a word the Source spells differently, a number read from a garbled feed, a gloss that goes beyond the Source's words. The Author checks these first.
6. **Preview.** `make dev`, then push the branch for a Preview. Check the diagram renders and the generated title card reads right.
7. **The Author verifies** every claim, number, quote and diagram against its Source, rewrites what is wrong, and sets `verified:`. Nothing publishes until the Author has done this.
8. **Publish.** The Author removes `draft: true`, merges and tags (`scripts/publish.sh vl`, check QA, `scripts/publish.sh cf`).

## Ground every claim

- **First-party Sources only.** The company's engineering blog, documentation, talks by its engineers, papers by the people who built it, specifications, source code and release notes. No other outlets' summaries, aggregators or "people say". A fact only a third party reports stays out.
- **Every claim is linked**, inline where it is made: each fact, number and quote links the first-party page that supports it, at the most specific anchor the page offers. Sources, last, lists every page cited.
- **Say "we infer"** in the sentence that reconstructs what the Sources leave open. A number of our own is labelled as ours ("a hypothetical overhead of one vCPU (our number, not Netflix's)").
- **Quote short and credited.** A sentence or two at a time, only where the exact words matter (a definition, a commitment, the wording of a number), in quotation marks, credited to the person or organisation and linked. Quote only as long as the point needs, follow each quote with our own analysis, and never let a quote carry a section or reproduce a Source's lede or summary. No law gives a safe word count (`docs/research/content-source-intake.md` §4.5).
- **Diagrams are always redrawn** in Mermaid from the facts, never copied, traced or embedded, and the original figure is linked beside ours ("The diagram is ours, redrawn from the text of Netflix's post; Netflix's own figures are in the original"). A Source's images are never reproduced (ADR-0007): diagrams are protected works in the US and India.
- **Code only under its repository's licence.** Quote code only when the repository's licence allows reuse, and keep the notice it requires (name the licence beside the snippet or in Sources). Fair dealing does not cover code in India, and a non-commercial licence (CC BY-NC-ND) does not cover an ad-funded site. Otherwise describe the code and link it.
- **Credit is not a substitute for analysis.** Google's guidance: "Attribution or giving credit to other sources doesn't replace the need for original effort." The Case Study earns its place by explaining why, and what it costs, beyond what the Source says.

## Notes, hero and `verified`

- **TODO(author) lines.** A note for the Author is a `TODO(author): …` line, placed as the first line under the first H2 after the TL;DR. The TL;DR may hold only its bullets (the `tldr` gate fails anything else). A production build fails while any `TODO(author)` is left.
- **Hero.** Leave `hero` and `hero_alt` out. The build draws the hero card from the title, the kind and the Vertical, with the title as alt text; check it on the Preview, do not make one. An image HastCode has the right to use (never a Source's), beside the Markdown and named in `hero` with `hero_alt`, replaces the card.
- **`verified:`** records that the Author checked every claim against its Source (ADR-0007). Only the Author sets it, after doing that check; an agent never writes it (CLAUDE.md). A production build fails on a non-draft Piece published after 2026-10-26 without it (`unverified`).

## The skeleton

These are the template's H2s, in order, and exactly what `make lead` writes for a Case Study. Only `TL;DR` and `Sources` are literal; rename every other heading to a plain sentence that states its point ("Why sharing a kernel ran out of road"), and add one H2 per part of the diagram.

```text
## TL;DR
## TODO(author): {system} at a glance, named for what the diagram shows
## TODO(author): the problem {company} had to solve
## TODO(author): the options {company} weighed
## TODO(author): the first part of the diagram and what it does
## TODO(author): the next part of the diagram and what it does
## TODO(author): what it costs
## TODO(author): the evidence it works, and what the Source leaves out
## TODO(author): what it makes possible
## Sources
```

Frontmatter: `title`, `description`, `kind: case-study`, `vertical`, `tags`, `published`, `reviewed`, `author: mira-vance`, `draft: true`; `hero`, `hero_alt`, `related` and `sources` optional and commented out. No `verified`.

## Section by section

- **Opening (before TL;DR).** One or two hook sentences: a sourced fact about the system that makes the Reader look twice (Titus: "launches millions of containers a week", then the opposite move). Then the AI-assistance note as a blockquote, exactly as the template has it.
- **TL;DR.** Three to five bullets, answer first: what it is, why it was built, how it works, what it costs.
- **At a glance.** The diagram, a sentence saying it is ours and linking the original, and the order the sections walk it.
- **The problem.** What the old design did and who paid for it, in the Source's evidence.
- **The options weighed.** Every alternative the Source names and why it lost, then the deciding argument. This is where a Case Study beats a summary.
- **One H2 per part of the diagram**, in order. Titus used "Step 1: federation decides who qualifies", "Step 2", "Step 3".
- **What it costs.** The trade-offs, plainly. Split into one H2 per cost when there are several.
- **The evidence.** What the company measured, and the honest gap: say when the Source publishes no numbers.
- **What it makes possible.** Who uses it, what is next, and the lesson for the Reader's own systems.
- **Sources.** Last, the Lead first, with at least one link.

Aim for 2,000 to 2,800 words, each H2 under 250 words (`section-length` warns above that; answer engines quote short sections best).

## Titles, descriptions and search intent

| Search | Where it goes |
|---|---|
| "how {company} {does X}", "{system} architecture" | Title `{system}: how {company} {does X}`, e.g. "Single-tenant Titus: how Netflix gives one container its own VM" |
| "{system} vs {alternative}" | The options section; a full Comparison when the pair earns one, linked in `related` |
| "why {company} {chose X}" | The problem and options sections, headed in the Reader's words |

The title stays under about 70 characters. The description is one sentence of 150 characters or fewer naming the system and the trade-off: the first Case Study's ran to 149.

## What the build checks

`tldr` (TL;DR first after the intro, 3–5 bullets and nothing else), `diagram` (at least one Mermaid fence), `diagram-render` (a fence that does not render), `sources` (Sources last, with a link), `author-placeholder` (a `TODO(author)` left in a production build), `unverified` (no valid `verified:` on a non-draft Piece published after 2026-10-26; a warning in preview), `frontmatter` (unknown keys, bad dates, `repo`/`who_for` on a non-Profile). Warnings: `section-length` (an H2 over 250 words), `review-overdue`.

## What the first one taught us

From the Titus draft and its check list (ticket 07):

- **Worked:** opening on a scale fact and then the reversal; a redrawn diagram whose boxes became the "Step 1/2/3" sections; a full section on the alternatives Netflix rejected; labelling three inferences in their own sentences; saying plainly that the post publishes no benchmark scores.
- **Did not:** the Lead's site refused direct fetches, so it was read through the feed, where one sentence arrived garbled. Read the Source as published, not only the feed, before drafting a claim on it. Terms drift between Sources ("Expected" vs "Earliest Departure Time"): use the Source's own term and link the other's. A gloss ("never evicted") went beyond the Source's words ("exempts the task from preemption"): keep to the Source's wording when the difference matters.

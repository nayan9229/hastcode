# Writing a Concept

A Concept (`GLOSSARY.md`) defines and teaches one named architectural idea or pattern for the Reader: a mid-to-senior backend engineer. One idea per piece; a second idea is a second Explainer and a Related link. This guide and [the template](concept-template.md) are drawn from the first one, [Metastable failures: why systems stay down after the trigger is gone](../../content/explainers/distributed-systems/metastable-failures.md) (ticket 07). Write in [the house byline voice](house-byline-voice.md) and hold every sentence to the [Editorial Policy](../../content/pages/editorial-policy.md).

## The flow

1. **Pick a Lead.** Run `make intake`, pick from the queue (`make daily` picks and drafts, ADR-0009; at most two Pieces are Published a day, enforced by the build, ADR-0010), then `make lead ID=<id>` (add `KIND=concept` or `VERTICAL=…` to override the Source's suggestion). The skeleton is this template, with the Lead's title, kind and Vertical filled in, the Lead's link in Sources, and `published` and `reviewed` set to the pick time: provisional, set both when you publish. Retitle it: the Lead is a post, the Concept is the idea (the Lead "Metastability and Lost Writes" became "Metastable failures"). To start by hand, copy [the template](concept-template.md) to `content/explainers/<vertical>/<slug>.md` on a branch. The slug names the idea (`metastable-failures`). Keep `draft: true`.
2. **Source notes.** Find who named the idea (the paper, the RFC, the original post) and the first-party accounts of it in real systems. A Lead is often the newest post in a long thread: follow it back.
3. **Diagram first.** Draw the mechanism in Mermaid: the loop, the states, the flow. It is the "how it works" section.
4. **Draft** into the skeleton. AI assistance is fine; it is a draft, not a source.
5. **Write the check list.** In the ticket or pull request, list every claim you are unsure of: counts that differ between a paper's abstract and its body, a paraphrase of a formula, a URL that now redirects. The Author checks these first.
6. **Preview.** `make dev`, then push the branch for a Preview. Check the diagram renders and the generated title card reads right.
7. **The Author verifies** every claim, number, quote and diagram against its Source, rewrites what is wrong, and sets `verified:`. Nothing publishes until the Author has done this.
8. **Publish.** The Author removes `draft: true`, merges and tags (`scripts/publish.sh vl`, check QA, `scripts/publish.sh cf`).

## Ground every claim

- **First-party Sources only.** Papers by the people who named the idea, specifications, engineering blogs and postmortems of the systems it happened to, documentation and source code. No other outlets' summaries, aggregators or "people say". A fact only a third party reports stays out.
- **Every claim is linked**, inline where it is made: each fact, number and quote links the first-party page that supports it, at the most specific anchor the page offers. Sources, last, lists every page cited.
- **Say "we infer"** in the sentence that reconstructs what the Sources leave open, such as a mechanism no Source spells out or a link between two ideas the author does not explain.
- **Quote short and credited.** A sentence or two at a time, only where the exact words matter (a definition, a commitment, the wording of a number), in quotation marks, credited to the person or organisation and linked. Quote only as long as the point needs, follow each quote with our own analysis, and never let a quote carry a section or reproduce a Source's lede or summary. No law gives a safe word count (`docs/research/content-source-intake.md` §4.5).
- **Diagrams are always redrawn** in Mermaid from the facts, never copied, traced or embedded, and the original figure is linked beside ours when the idea came with one. A Source's images are never reproduced (ADR-0007): diagrams are protected works in the US and India.
- **Code only under its repository's licence.** Quote code only when the repository's licence allows reuse, and keep the notice it requires (name the licence beside the snippet or in Sources). Fair dealing does not cover code in India, and a non-commercial licence (CC BY-NC-ND, such as Google's SRE book) does not cover an ad-funded site. Otherwise describe the code and link it.
- **Credit is not a substitute for analysis.** Google's guidance: "Attribution or giving credit to other sources doesn't replace the need for original effort." A Concept earns its place by connecting Sources no single one connects, and by the worked example.

## Notes, hero and `verified`

- **TODO(author) lines.** A note for the Author is a `TODO(author): …` line, placed as the first line under the first H2 after the TL;DR. The TL;DR may hold only its bullets (the `tldr` gate fails anything else). A production build fails while any `TODO(author)` is left.
- **Hero.** Leave `hero` and `hero_alt` out. The build draws the hero card from the title, the kind and the Vertical, with the title as alt text; check it on the Preview, do not make one. An image HastCode has the right to use (never a Source's), beside the Markdown and named in `hero` with `hero_alt`, replaces the card.
- **`verified:`** records that the Author checked every claim against its Source (ADR-0007). Only the Author sets it, after doing that check; an agent never writes it (CLAUDE.md). A production build fails on a non-draft Piece published after 2026-10-26 without it (`unverified`).

## The skeleton

These are the template's H2s, in order, and exactly what `make lead` writes for a Concept. Only `TL;DR` and `Sources` are literal; rename every other heading to a plain sentence that states its point ("The loop that keeps it down"), and add one H2 per practice.

```text
## TL;DR
## TODO(author): what {concept} is
## TODO(author): how {concept} works, named for what the diagram shows
## TODO(author): where {concept} shows up in real systems
## TODO(author): why systems end up with {concept}
## TODO(author): the first practice that applies or defends against {concept}
## TODO(author): what {concept} is not
## TODO(author): what understanding {concept} makes possible
## Sources
```

Frontmatter: `title`, `description`, `kind: concept`, `vertical`, `tags`, `published`, `reviewed`, `author: mira-vance`, `draft: true`; `hero`, `hero_alt`, `related` and `sources` optional and commented out. No `verified`.

## Section by section

- **Opening (before TL;DR).** A real, sourced moment the idea happened to a named system, told concretely: the 2015 DynamoDB disruption, minute by minute, ending on "What held DynamoDB down for hours was its own retry traffic." No definition yet. Then the AI-assistance note as a blockquote, exactly as the template has it.
- **TL;DR.** Three to five bullets, answer first: the definition, the mechanism, how common it is, what to do and its cost.
- **What it is.** The definition from the people who named it, quoted short and credited, then in our words. Say what does not count.
- **How it works.** The diagram, walked in order with one worked example that puts numbers on each box (the HotOS paper's 280 → 560 QPS).
- **Where it shows up.** Named systems and variants, each from its first-party account.
- **Why systems end up with it.** The forces that lead there, often features that help on a normal day.
- **One H2 per practice** that applies the idea or defends against it, each with who recommends it and what it costs (the Concept had three defences and a recovery section).
- **What it is not.** The nearest idea and the boundary; a neighbouring idea that deserves its own piece gets one sentence and a future Explainer.
- **What it makes possible.** What the Reader can now do or measure, and where to try it.
- **Sources.** Last, the Lead first, with at least one link.

Aim for 2,000 to 2,800 words, each H2 under 250 words (`section-length` warns above that; answer engines quote short sections best).

## Titles, descriptions and search intent

| Search | Where it goes |
|---|---|
| "what is {concept}", "{concept}" | Title `{concept}: {the one-line answer}`, e.g. "Metastable failures: why systems stay down after the trigger is gone" |
| "how to prevent {concept}", "{concept} mitigation" | The practice sections, headed in the Reader's words ("Defence one: give retries a budget") |
| "{concept} vs {other idea}" | The "what it is not" section; a Comparison when the pair earns one |

The title stays under about 70 characters. The description is one sentence of 150 characters or fewer: what the idea is and why it matters. The first Concept's ran to 162; keep under 150.

## What the build checks

`tldr` (TL;DR first after the intro, 3–5 bullets and nothing else), `diagram` (at least one Mermaid fence), `diagram-render` (a fence that does not render), `sources` (Sources last, with a link), `author-placeholder` (a `TODO(author)` left in a production build), `unverified` (no valid `verified:` on a non-draft Piece published after 2026-10-26; a warning in preview), `frontmatter` (unknown keys, bad dates, `repo`/`who_for` on a non-Profile). Warnings: `section-length` (an H2 over 250 words), `review-overdue`.

## What the first one taught us

From the metastable-failures draft and its check list (ticket 07):

- **Worked:** an outage as the hook instead of a definition; the paper's own worked example walked through the diagram; following the Lead back through four earlier posts by the same author and the two papers they rest on; giving every defence its cost; a closing that points at runnable reproductions.
- **Did not:** the Lead paired two ideas (metastability and lost writes). One idea per piece: the second got one short section and is left for its own Explainer. Papers disagree with themselves (an abstract's 22 failures vs the body's 21): quote the figure you cite and note the other in the check list. A Source's URL that now redirects is cited at the address it resolves to, after checking it still says the same.

# Writing a Comparison

A Comparison (`GLOSSARY.md`) sets two or more alternatives side by side to help the Reader, a mid-to-senior backend engineer, choose. It answers the same questions for every option, in the same order, and ends on a table. This guide and [the template](comparison-template.md) are drawn from the first one, [Connection reuse in HTTP/1.1, HTTP/2 and HTTP/3: how each works and where each breaks](../../content/explainers/apis-networking/http-connection-reuse.md) (ticket 07). Write in [the house byline voice](house-byline-voice.md) and hold every sentence to the [Editorial Policy](../../content/pages/editorial-policy.md).

## The flow

1. **Pick a Lead.** Run `make intake`, pick from the queue (at most two Leads a day), then `make lead ID=<id>` (add `KIND=comparison` or `VERTICAL=…` to override the Source's suggestion). The skeleton is this template, with the Lead's title, kind and Vertical filled in and the Lead's link in Sources. Retitle it: the Lead is one option's news, the Comparison is the choice (the Lead "Better HTTP/1 connection reuse" became HTTP/1.1 vs HTTP/2 vs HTTP/3). To start by hand, copy [the template](comparison-template.md) to `content/explainers/<vertical>/<slug>.md` on a branch. The slug names the choice (`http-connection-reuse`). Keep `draft: true`.
2. **Pick the questions.** Choose the options (two to four) and the dimensions you will answer for every one of them. The dimensions become H2s and the rows of the Trade-offs table.
3. **Source notes.** Each option's own specification or documentation, at the section level. Note what each specification does not say as carefully as what it does: "not described" is a cell too.
4. **Diagram first.** One Mermaid diagram showing the same request or workload through each option, one subgraph per option.
5. **Draft** into the skeleton. AI assistance is fine; it is a draft, not a source.
6. **Write the check list.** In the ticket or pull request, list every claim you are unsure of: a default that depends on how a library is configured, a specification's "peer" glossed as "server", a measurement quoted with its caveat. The Author checks these first.
7. **Preview.** `make dev`, then push the branch for a Preview. Check the diagram and table render and the generated title card reads right.
8. **The Author verifies** every claim, number, quote, diagram and table cell against its Source, rewrites what is wrong, and sets `verified:`. Nothing publishes until the Author has done this.
9. **Publish.** The Author removes `draft: true`, merges and tags (`scripts/publish.sh vl`, check QA, `scripts/publish.sh cf`).

## Ground every claim

- **First-party Sources only.** Each option's specification or RFC, official documentation, source code and release notes, and its makers' own posts and measurements. No benchmarks from other outlets, aggregators or "people say". A fact only a third party reports stays out.
- **Every claim is linked**, inline where it is made: each fact, number and quote links the first-party page that supports it, at the most specific anchor the page offers (an RFC section, a documented field). Sources, last, lists every page cited, with the sections used.
- **Say "we infer"** in the sentence that reconstructs what the Sources leave open. A recommendation no Source makes is our inference and says so; a section of them says so once at the top.
- **Quote short and credited.** A sentence or two at a time, only where the exact words matter (a MUST in a specification, a definition, the wording of a number), in quotation marks, credited and linked. Quote only as long as the point needs, follow each quote with our own analysis, and never let a quote carry a section or reproduce a Source's lede or summary. No law gives a safe word count (`docs/research/content-source-intake.md` §4.5). RFCs may be reproduced only in full and unmodified, so quoting them relies on quotation like any other text.
- **Diagrams are always redrawn** in Mermaid from the facts, never copied, traced or embedded, and the original figure is linked beside ours when a Source has one. A Source's images are never reproduced (ADR-0007): diagrams are protected works in the US and India.
- **Code only under its repository's licence.** Quote code only when the repository's licence allows reuse, and keep the notice it requires (the Go blog is CC BY 4.0, credited in Sources). Fair dealing does not cover code in India, and a non-commercial licence (CC BY-NC-ND) does not cover an ad-funded site. Otherwise describe the code and link it.
- **Credit is not a substitute for analysis.** Google's guidance: "Attribution or giving credit to other sources doesn't replace the need for original effort." A Comparison earns its place by answering the same questions for every option, which no single option's documentation does.

## Notes, hero and `verified`

- **TODO(author) lines.** A note for the Author is a `TODO(author): …` line, placed as the first line under the first H2 after the TL;DR. The TL;DR may hold only its bullets (the `tldr` gate fails anything else). A production build fails while any `TODO(author)` is left.
- **Hero.** Leave `hero` and `hero_alt` out. The build draws the hero card from the title, the kind and the Vertical, with the title as alt text; check it on the Preview, do not make one. An image HastCode has the right to use (never a Source's), beside the Markdown and named in `hero` with `hero_alt`, replaces the card.
- **`verified:`** records that the Author checked every claim against its Source (ADR-0007). Only the Author sets it, after doing that check; an agent never writes it (CLAUDE.md). A production build fails on a non-draft Piece published after 2026-10-26 without it (`unverified`).

## The skeleton

These are the template's H2s, in order, and exactly what `make lead` writes for a Comparison. `TL;DR`, `Trade-offs` and `Sources` are literal (the build looks for the table under `Trade-offs`); rename every other heading to a plain sentence that states its point ("Head-of-line blocking moves down a layer each time"), and add one H2 per option and per dimension.

```text
## TL;DR
## TODO(author): the options side by side, named for what the diagram shows
## TODO(author): {option A}: how it works
## TODO(author): {option B}: how it works
## TODO(author): the first dimension, answered for every option
## Trade-offs
## TODO(author): when each one wins
## Sources
```

Frontmatter: `title`, `description`, `kind: comparison`, `vertical`, `tags`, `published`, `reviewed`, `author: mira-vance`, `draft: true`; `hero`, `hero_alt`, `related` and `sources` optional and commented out. No `verified`.

## Section by section

- **Opening (before TL;DR).** A sourced fact that makes the choice matter (the Go team's 710 ms vs 150 ms, with its "one machine" caveat), ending on the question the piece answers for every option. Then the AI-assistance note as a blockquote, exactly as the template has it.
- **TL;DR.** Three to five bullets: one per option, then which to pick when, labelled "we infer" unless a Source says it.
- **Side by side.** What the options share, and the diagram: the same workload through each.
- **One H2 per option**: how it works, from its own specification or documentation.
- **One H2 per dimension**, answering it for every option in the same order (the first Comparison had four: head-of-line blocking, limits, idle connections, coalescing).
- **Trade-offs.** A Markdown table: one row per dimension, one column per option, short cells. Open with "Each cell summarises a section above, where its source is linked." so the cells need no links of their own; every cell must trace to a linked sentence above.
- **When each wins.** One paragraph per option: where it is the right choice and what to watch. Our inference, and the section says so.
- **Sources.** Last, the Lead first, with at least one link.

Aim for 2,000 to 2,800 words, each H2 under 250 words (`section-length` warns above that; answer engines quote short sections best).

## Titles, descriptions and search intent

| Search | Where it goes |
|---|---|
| "{A} vs {B}" | Title `{A} vs {B}: {what the choice turns on}` for two options |
| "{A} vs {B} vs {C}", "{topic} in {A}, {B} and {C}" | Title `{topic} in {A}, {B} and {C}: {the question}`, e.g. the first Comparison |
| "when to use {A}", "{A} or {B} for {use}" | The "when each wins" section, headed in the Reader's words |

The title stays under about 70 characters; the first Comparison's 85 is too long for a search result. The description is one sentence of 150 characters or fewer naming the options and the question; the first one's ran to 164.

## What the build checks

`trade-offs-table` (a Markdown table under a `## Trade-offs` heading), `tldr` (TL;DR first after the intro, 3–5 bullets and nothing else), `diagram` (at least one Mermaid fence), `diagram-render` (a fence that does not render), `sources` (Sources last, with a link), `author-placeholder` (a `TODO(author)` left in a production build), `unverified` (no valid `verified:` on a non-draft Piece published after 2026-10-26; a warning in preview), `frontmatter` (unknown keys, bad dates, `repo`/`who_for` on a non-Profile). Warnings: `section-length` (an H2 over 250 words), `review-overdue`.

## What the first one taught us

From the HTTP connection-reuse draft and its check list (ticket 07):

- **Worked:** turning one option's news (a Go release) into the choice between three; one diagram with a subgraph per option; dimension sections that answer every option in the same order, so the table wrote itself; marking "When each wins" as inference in its first sentence; anchoring every specification claim to its RFC section.
- **Did not:** library behaviour depends on configuration (Go's `DefaultTransport` vs a hand-built `Transport` differ on HTTP/2): state which configuration a claim is about. A specification's wording ("the peer") was narrowed in a gloss ("a server"): keep the specification's word. "Not described" in a table cell rests on reading a whole section: cite the section the absence is read from.

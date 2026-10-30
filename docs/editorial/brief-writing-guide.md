# Writing a Brief

A Brief (`GLOSSARY.md`) is a short, dated commentary on something that just happened in tech, with the architectural "why it matters", linking into Explainers. It is curated from public Sources, never original reporting, and it is the kind most exposed to Google's policies on summarised content: a Brief that restates its Source is the "mainly summarizing what others have to say" warning sign, and citing the Source does not fix that (`docs/research/content-source-intake.md` §3.5). Its value is our analysis and the link into an Explainer. This guide and [the template](brief-template.md) are drawn from the first one, [GitHub's Dependabot backlog: a queue flooded by its own fan-out](../../content/briefs/github-dependabot-fan-out.md) (ticket 07). Write in [the house byline voice](house-byline-voice.md) and hold every sentence to the [Editorial Policy](../../content/pages/editorial-policy.md).

## The rules a Brief adds

- **A fresh Lead.** The Lead is at most 3 days (72 hours) old when the Brief is published, counted from the Lead's date in the queue. For a page updated in place (an incident report, a status page), count from its first post, not its last update. `make lead KIND=brief` refuses an older Lead; a Lead that has aged out can still become an Explainer.
- **At most 250 words**, the body's headings and the AI-assistance note included (`brief-length`).
- **Exactly two sections**, `## What happened` then `## Why it matters` (`brief-sections`). The only text before them is the AI-assistance note.
- **"Why it matters" is HastCode's own analysis**, never the Source restated or reworded: the architectural cause, the lesson for the Reader's own systems, and the link into an Explainer. It is at least as long as "What happened".
- **`explainers` lists exactly the slugs the body links**, at least one, no more (`brief-explainer-link`). Link an Explainer by its bare slug (`[metastable failures](metastable-failures)`) or its path (`/distributed-systems/metastable-failures/`). Every Brief needs an Explainer to point at, so Brief volume follows Explainer volume (ADR-0007).
- **`source_url` is the Lead's link**: the first-party page the Brief is about (`brief-source`). A Brief has no Sources section; every other page it cites is linked inline, and may also go in the optional `sources` frontmatter.
- **Inside the six Verticals.** Cover what the Reader would want anyway; a Lead that is only trending is a warning sign, not a Brief.

## The flow

1. **Pick a Lead** from that morning's `make intake` queue (`make daily` picks and drafts, ADR-0009; at most two Pieces are Published a day, enforced by the build, ADR-0010), then, with the Explainer it teaches in hand (step 2), `make lead ID=<id> KIND=brief EXPLAINER=<slug>` (several: `EXPLAINER=a,b`). The skeleton is this template, with the Lead's title in `title`, its link in `source_url`, the named Explainers in `explainers` and linked in "Why it matters", and `published` set to the pick time: provisional, set it when you publish. To start by hand, copy [the template](brief-template.md) to `content/briefs/<slug>.md` on a branch. The slug names the event (`github-dependabot-fan-out`). Keep `draft: true`.
2. **Find the Explainer** the event teaches. If none exists yet, the Lead is not a Brief today: `make lead` refuses a Brief without `EXPLAINER=`, or naming an Explainer that does not exist.
3. **Draft** into the skeleton. AI assistance is fine; it is a draft, not a source.
4. **Write the check list.** In the ticket or pull request, list every quote, time and inference for the Author to check against the Source.
5. **Preview.** `make dev`, then push the branch for a Preview. Check the generated title card reads right.
6. **The Author verifies** every claim and quote against the Source, rewrites what is wrong, and sets `verified:`, inside the 3-day window. Nothing publishes until the Author has done this.
7. **Publish.** The Author removes `draft: true`, merges and tags (`scripts/publish.sh vl`, check QA, `scripts/publish.sh cf`).

## Ground every claim

- **First-party Sources only.** The company's own announcement, incident report, status page, release notes or engineering post. Never another site's coverage of it.
- **Every claim is linked**, inline where it is made, to the first-party page that supports it.
- **Say "we infer"** in the sentence that reads more into the Source than it states ("We infer from the timeline that capacity alone did not clear the backlog quickly; the report does not say so directly.").
- **Quote short and credited.** One short quote at most, only where the exact words matter, in quotation marks, credited and linked, followed by our own analysis. State the headline facts in our own words; never reproduce the Source's lede or summary paragraph (research §4.5: a Brief that substitutes for its Source fails both copyright's market test and Google's scraping policy).
- **Diagrams are always redrawn** in Mermaid from the facts, never copied, and the original linked. A Brief rarely has room for one; link the Explainer's instead. A Source's images are never reproduced (ADR-0007).
- **Code only under its repository's licence**, with the notice it requires. Fair dealing does not cover code in India. In a Brief, link the code instead.
- **Credit is not a substitute for analysis.** Google's guidance: "Attribution or giving credit to other sources doesn't replace the need for original effort."

## Notes, hero and `verified`

- **TODO(author) lines.** A note for the Author is a `TODO(author): …` line, placed as the first line under `## What happened`. A production build fails while any `TODO(author)` is left.
- **Hero.** Leave `hero` and `hero_alt` out. The build draws the hero card from the title and the label "Brief", with the title as alt text; check it on the Preview, do not make one. An image HastCode has the right to use (never a Source's), beside the Markdown and named in `hero` with `hero_alt`, replaces the card.
- **`verified:`** records that the Author checked every claim against its Source (ADR-0007). Only the Author sets it, after doing that check; an agent never writes it (CLAUDE.md). A production build fails on a non-draft Piece published after 2026-10-26 without it (`unverified`).

## The skeleton

These are the template's H2s, in order, and exactly what `make lead KIND=brief` writes. Both are literal: the build accepts no other sections.

```text
## What happened
## Why it matters
```

Frontmatter: `title`, `description`, `tags`, `published`, `author: mira-vance`, `source_url`, `explainers`, `draft: true`; `hero`, `hero_alt`, `related` and `sources` optional and commented out. No `kind`, `vertical` or `reviewed` (a Brief rejects them), and no `verified`.

## Section by section

- **Opening.** The AI-assistance note as a blockquote, exactly as the template has it, and nothing else.
- **What happened.** The facts in our own words: who, what, when (UTC), linked. Two to four sentences. One short quote at most.
- **Why it matters.** Our analysis: the architectural cause in plain words, the lesson for the Reader, and the Explainer link that teaches it. Close on what a Reader can do differently.

## Titles and descriptions

Title: `{company}'s {event}: {the lesson}`, e.g. "GitHub's Dependabot backlog: a queue flooded by its own fan-out". It names the system and the architectural point, not only the news; under about 70 characters. The description is one sentence of 150 characters or fewer: what happened and the lesson. The first Brief's ran to 166.

## What the build checks

`brief-length` (over 250 words), `brief-sections` (sections other than exactly What happened, Why it matters), `brief-source` (missing or relative `source_url`), `brief-explainer-link` (no Explainer link, or `explainers` not exactly the linked slugs), `broken-link` (a link to an Explainer that does not exist), `author-placeholder` (a `TODO(author)` left in a production build), `unverified` (no valid `verified:` on a non-draft Piece published after 2026-10-26; a warning in preview), `frontmatter` (unknown keys such as `kind`, bad dates). The 3-day rule is enforced when the Lead is picked, not by the build: check the date before publishing.

The Friday roundup is an ordinary Brief; see [the Repo Profile writing guide](repo-profile-writing-guide.md#the-friday-roundup).

## What the first one taught us

From the Dependabot draft and its check list (ticket 07):

- **Worked:** a Lead from an incident feed, just under 3 days old at drafting; "What happened" in UTC times with the Source's own words only where they carried the point; a "why it matters" that named the general lesson (one producer's burst delays everyone on a shared queue; the fix was a separate, rate-limited lane) and linked the Concept that teaches it; one labelled inference; the link to the Concept framed as a lesson by analogy, because the report never uses the word.
- **Did not:** the 3-day window was ambiguous for an incident report updated in place (first post 14:58 UTC, resolution 18:51 UTC, Brief published at 18:00 UTC three days later): count from the first post from now on. It used two quotes in "What happened"; one is the limit. It left out the AI-assistance note, which the Editorial Policy requires on every piece drafted with AI assistance; the template now carries it, inside the 250 words. Its description ran past 150 characters.

---
status: proposed
date: 2026-11-18
supersedes: the "AI-assistance note on every Piece" practice (Editorial Policy, house voice guide, the five writing guides and kind templates; never recorded in an ADR)
---

# No AI note is published

Until now every Piece carried the Editorial Policy's note as a blockquote ("Drafted with AI assistance, researched, verified and edited by the HastCode editorial team. We publish only as many pieces as the Author can verify this way."): after the hook in an Explainer, first on the page in a Brief. The Editorial Policy had an "AI assistance" section that quoted it, and the Author page and About page said drafts may be written with AI assistance. Writing-craft ticket 06 asked the Author where the note should sit. On 2026-11-18 the Author decided, in their words:

> AI should not the part of content and not able to see by any end user. any ai note should not publish

We decided:

- **No Reader sees an AI-assistance note or AI-assistance wording anywhere on the site.** That covers every HTML page, the Markdown twins (ADR-0011), `llms.txt`, `sitemap.md`, the feeds, the sitemaps and the structured data. The note is removed from all 79 tracked Pieces (74 live, 5 drafts from `make daily` on 2026-11-18), from the five kind templates (the skeletons `make lead` writes), from `make daily`'s drafter instructions, and from the Editorial Policy, the About page and the Author page.
- **The build enforces it** with the error rule `ai-note` (`builder/step_validate_ai_note.go`, validate 13). It fails any Explainer, Brief, static page, Author file or site.yaml Vertical whose Reader-visible text (title, description, tags, hero alt, `who_for`, the whole body, an Author's name, photo alt and bio, a Vertical's name and description) says AI made the Piece or the site's Pieces: the note's own wording anywhere ("drafted with/by/using" AI, Claude, an LLM and the like; an "AI-assistance note"); a sentence that opens with how it was made ("Generated with Claude.") or is only an "AI-generated" label; or a sentence about this Piece or the site's own work ("this piece", "our pieces", HastCode, the editorial team, Mira Vance) that says AI made it ("written with the help of AI", "AI-generated", "with AI assistance"). Any case, across line breaks, with any Unicode space or dash. AI as a topic builds as before ("AI-assisted development", "Copilot's AI assistance in code review", "Was pgrust written by AI?"). Source titles are other publishers' titles and are not read. Frontmatter YAML comments are never published and are not read. A note in other words than these can still slip past the rule; the Author's verification before every Release is the last check.
- **Scope: the site, not the repository.** "Any end user" is read as a Reader of the site: every surface listed above. The source repository is public on GitHub today (`nayan9229/hastcode`), and anyone who finds it can read ADR-0009, `scripts/daily/`, the `.scratch/` review queues and this ADR, and the removed note in every Piece's git history. Making the repository private, or leaving it public, is the Author's decision; it is not taken here.
- **Honesty guard.** The site never claims or implies that the Pieces are written only by humans or without tools. Sentences that mentioned AI were removed, not replaced with other claims. Where a page said the team writes the Pieces, the claim was narrowed to what is true without the note: "written and edited by" became "edited by" (About, Contact), "written and checked by" became "checked by" (About's description), and the voice the team "writes in" became the voice it "publishes under" (Editorial Policy). What stays is true of the process: first-party Sources, every claim linked, claim-by-claim verification by the Author before publishing, the publishing pace that verification sets, the 90-day Review and Corrections. Mira Vance stays a disclosed house byline, not a real person (GLOSSARY "House Byline").
- **The drafting process is unchanged.** `make daily` still drafts with Claude (ADR-0009) and the Author still verifies every claim and sets `verified:` (ADR-0007). This decision is about what Readers see, not about how Pieces are made.

Not decided here: ticket 06's part (b), whether the Editorial Policy and the Author page say plainly who verifies behind the house byline. It stays open with the Author; this ADR leaves the "Author" wording on those pages as it was.

## Considered options

- **Keep the note where it is**: the research note's evidence for a disclosure (Google calls disclosures "useful for content where someone might think 'How was this created?'"; BBC 3.4.29), and the attention it takes from the opening (NN/g: the first two paragraphs get the most attention), `docs/research/article-writing-craft.md` §5.1 and §5.2 item 10. Rejected by the Author.
- **Move the note lower** (after the TL;DR, or to the foot of the page): the option ticket 06 asked about. Rejected by the Author: no note at all.
- **Keep the disclosure on the Editorial Policy only**: rejected; the Author's words cover every page.

## Consequences

- This supersedes: the house voice guide's line "Pieces drafted with AI assistance carry the disclosure in the Editorial Policy"; each writing guide's Opening rule (the note "exactly as the template has it"); the Brief guide's rules that the note is the only text before "What happened" and counts toward the 250-word cap; the self-edit checklist's two Brief lines that named the note; the Editorial Policy's "AI assistance" section; the row "AI-assistance note on every Piece" in the research note §5.1 and its open question §5.2 item 10 (now closed); the planning user stories that asked for the note (`.scratch/hastcode-v1/spec.md` story 17, `.scratch/hastcode-content-plan/spec.md` story 13).
- A Brief's 250 words are now all its own: the note's 26 words no longer count toward `brief-length`.
- The research note's evidence for a disclosure still stands; the Author chose against it. A later change of mind is a new ADR, and `ai-note` would have to go with it.
- EU AI Act Article 50(4) requires a deployer of an AI system that generates text published to inform the public on matters of public interest to disclose that the text is AI-generated, unless the text has undergone human review or editorial control and a natural or legal person holds editorial responsibility for its publication. The Author's claim-by-claim verification before every Release is that review and control. This is not legal advice; whether the exception applies is the Author's call.

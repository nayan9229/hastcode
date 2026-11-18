# 06: Decide AI-note placement and byline wording

**What to build:** Two decisions only the Author can make, recorded either way. If either one changes anything, it goes into an ADR (next free number: 0013) through `/grill-with-docs` or the domain-modeling skill. The follow-up edits to templates, guides, tests, pages and content then become their own tickets. If neither changes, this ticket closes with a one-line note saying so. The research note gives no recommendation on either question and lists both as open (`docs/research/article-writing-craft.md` §5.2 items 10 and 13, §7). Breakdown approved 2026-11-17.

**(a) Where the AI-assistance note sits.**

- *Today:* every Piece carries the Editorial Policy's note as a blockquote ("Drafted with AI assistance, researched, verified and edited by the HastCode editorial team. We publish only as many pieces as the Author can verify this way."). In an Explainer of any kind it comes right after the hook, before the TL;DR. Each kind's writing guide calls for this and each template carries it. In a Brief it is the very first thing on the page: the Brief guide allows only the note before "What happened", and the note counts toward the 250-word cap (`brief-length`). The Repo Profile template test checks that a Profile page carries the note. The guides also record why the Brief template carries it: a past Brief left it out.
- *Evidence (§5.2 item 10):* NN/g says the first two paragraphs get the most attention. Google wants the disclosure, as useful for content where someone might wonder "How was this created?", and says nothing about where on the page it goes. BBC 3.4.29 is cited for disclosure in §5.1. Whether moving it lower would be just as honest is left to the Author and the Editorial Policy.
- *What a change would touch:* all four kind templates and the Brief template, the five writing guides, the tests that look for the note, `make daily`'s drafter (which keeps the note as given), and the 74 published Pieces, if they are moved too.

**(b) Whether the Editorial Policy and the Author page say plainly who verifies the work behind the house byline.**

- *Today:* Mira Vance is a disclosed house byline, not a real individual (GLOSSARY "House Byline", `docs/editorial/house-byline-voice.md`). That stays as it is. The Editorial Policy's "Who is accountable" says the HastCode editorial team is accountable for every piece. Its "AI assistance" section says "the Author researches every claim against its first-party source … and decides whether the piece is published". The Author page says pieces are "researched, verified and edited by the HastCode editorial team, who are accountable for them". Its structured data names Mira Vance as an Organization (the team), never a Person. Neither page names a person. ADR-0009 says the Author's verification is "the only human check between a feed item and a live page", and ADR-0007 says only the Author sets `verified:`. No ADR records the house byline or the note's placement: both live in GLOSSARY, the house voice guide, the writing guides and the Editorial Policy.
- *A wording wrinkle to know about:* GLOSSARY defines Author as "the named person, or the team publishing under one byline", and says "at launch the only Author is the House Byline". The Policy's AI section, the AI note ("as many pieces as the Author can verify") and CLAUDE.md all use "the Author" for whoever verifies and releases. The Policy uses both senses on one page: its opening says "the Author named on each piece is accountable", which is the byline, and its AI section's "the Author researches every claim" is the person who verifies. A change to (b) may need the glossary entry looked at too.
- *Evidence (§5.2 item 13):* Google asks, "Is it self-evident to your visitors who authored your content?". The Quality Rater Guidelines accept an "alias or username" when "there is good reason for anonymity" (§2.5.2, §5.5). The note marks its own suggestion as [inferred]: that the Policy and the Author page say clearly who stands behind the verification.

**Blocked by:** None (can start immediately)

**Status:** ready-for-human

- [ ] The Author decides (a): the note stays where it is, or it moves, and if so where, per kind (Explainers, Briefs)
- [ ] The Author decides (b): the Policy and Author page stay as they are, or they say who verifies (and in what words: a role, a team, or a named person)
- [ ] Any change is recorded as an ADR, made with `/grill-with-docs` or the domain-modeling skill, including GLOSSARY if "Author" or "House Byline" changes meaning
- [ ] Each follow-up change (templates, guides, tests, Editorial Policy, Author page, existing Pieces) is a new ticket, not part of this one
- [ ] If nothing changes, a short note under a `## Comments` heading at the bottom of this ticket says so, with the date and one line per question (a and b)
- [ ] No agent makes either decision or edits the Editorial Policy, Author page, templates or Pieces for this ticket

## Comments

**2026-11-18, (a) decided by the Author:** "AI should not the part of content and not able to see by any end user. any ai note should not publish". So the note does not move: it goes, everywhere a Reader could see it. Recorded as ADR-0014 (`docs/adr/0014-no-ai-note-is-published.md`). The Author asked for the follow-up edits in the same request, so they were made on branch `impl/noai-1118` rather than as separate tickets: the note removed from the 79 tracked Pieces, 74 live and the 5 drafts of Daily 2026-11-18 (body line only, no frontmatter touched), the five kind templates, `make daily`'s drafter prompt and the test fixtures; the Editorial Policy, About, Contact and Author page rewritten without AI wording (sentences removed, never replaced with a claim that the Pieces are written without tools); a new build error rule `ai-note`; the voice guide, writing guides, self-edit checklist, runbook and README updated.

**2026-11-18, (b) not decided.** Whether the Editorial Policy and the Author page say plainly who verifies behind the house byline is still open with the Author. Nothing about it changed: the pages keep their "Author" and "HastCode editorial team" wording, and GLOSSARY "Author" and "House Byline" are unchanged.

**2026-11-18, review follow-up: for the Author to approve before the next `make release`.** This ticket asked that no agent edit the Editorial Policy, Author page, templates or Pieces, and that each follow-up be its own ticket. On the Author's one-line request an agent did all of it here. So ADR-0014 is now `status: proposed`, and the (a) and ADR boxes stay unticked, until the Author has read and approved:

- the ADR-0014 text (`docs/adr/0014-no-ai-note-is-published.md`), including its new "Scope: the site, not the repository" bullet;
- the rewritten public wording: About and Contact "written and edited by" became "edited by"; About's description "written and checked by" became "checked by"; the Editorial Policy's voice the team "writes in" became the voice it "publishes under", and its "AI assistance" section became "Verification"; the Author page lost its AI sentence;
- phrases left as they were, which a reviewer flagged as leaning toward "the team writes the Pieces" now the note is gone: About "no say over what we write", "companies it writes about"; Editorial Policy, Independence, "no say over what we write or how". Reword them or accept them;
- the source repository: `nayan9229/hastcode` is public, and its docs (ADR-0009, `scripts/daily/`, `.scratch/` review queues, ADR-0014) and every Piece's git history show the AI drafting and the old note. Whether to make it private is the Author's call; ADR-0014 treats the repository as out of scope until then.

The `ai-note` rule was widened the same day: it now fails any wording that says AI made this Piece or our Pieces ("written with the help of AI", "this piece is AI-generated", "Generated with Claude.", Unicode dashes and spaces included), lets "AI assistance" as a topic through, no longer reads Source titles, and also reads site.yaml's Vertical names and descriptions.

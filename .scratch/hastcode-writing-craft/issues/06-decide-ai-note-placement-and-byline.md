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

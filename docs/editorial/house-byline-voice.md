# House byline: Mira Vance

Mira Vance is HastCode's **house byline**: the single, disclosed editorial
voice the team publishes under. Mira is not a real individual and must never
be presented as one. The Author page says so, and the Editorial Policy
explains it. Inspiration: `docs/research/persona-inspiration-cleo-abram.md`.

Use this guide when drafting or editing any Explainer or Brief.

## Who Mira is (as a voice)

A curious, optimistic systems explainer who believes the best way to
understand software is to see how real systems are actually built, and who
gets visibly excited about what good architecture makes possible, without
ever overselling it.

## Voice rules

1. **Lead with a concrete hook.** Open on one real, sourced fact about a real
   system (a number, an outage, a design choice), not a definition.
2. **Answer first.** The TL;DR states the answer in plain words before any
   detail.
3. **One idea per piece.** If a second idea appears, it is a separate
   Explainer and a Related link.
4. **Diagram first.** Choose the key diagram before writing. The H2 sections
   walk through it in order.
5. **Plain language.** Define every term at first use. Short sentences:
   about 20 words on average. Split anything over 30 unless it is a
   deliberate list. Use the shortest sentences at the hardest points. No
   hype or filler words: revolutionary, game-changing, magic (a technical
   term such as "magic number" is fine), simply, just (as an instruction:
   "just add", "just run"), easy, obviously, "it is important to note",
   "please note". The build warns on both (`sentence-length`,
   `banned-word`); a Source's words inside quotation marks are never
   reworded to clear one. This is the list's one home: other docs point
   here, and `builder/step_validate_prose.go` holds the build's copy.
6. **Optimistic, honest.** Name the trade-offs plainly, then close on what the
   design makes possible: "what could go right" without hiding the cost.
7. **Show the receipts.** Every claim, number and quote links its
   first-party source. Say "we infer" when reconstructing. Every number
   carries its unit, its date or version, and what it is compared with. In a
   Repo Profile, pointing at Signals meets this for a number the Fact Pack
   carries (stars, license, release dates, contributors): Signals shows it
   with its date, so never retype it into prose.
8. **Talk to the Reader.** Second person is fine ("you'd expect..."); first
   person singular is not, because Mira is a team voice. Use "we".
9. **One name per thing.** Pick the Source's term and use it every time;
   mention a synonym once, in brackets, if Readers search for it.
10. **At most one analogy.** Analogies: at most one per piece, and say where
    it stops being true. This is a house choice, not a research finding.

## Drafting pipeline

Pitch (one sentence: the idea and why it matters) → source notes (first-party
links only) → diagram outline (Mermaid) → draft in this voice → self-edit
with [the checklist](self-edit-checklist.md) → verify every claim → edit →
publish under `author: mira-vance`.

The self-edit checks the writing, not the facts. It comes before the
Author's verification against Sources and never replaces it.

No Piece carries an AI note, and nothing a Reader sees says how a Piece was drafted (ADR-0013). The build's `ai-note` rule fails the words "AI assistance" or "drafted with AI" on any Piece, static page or Author page. Writing about AI as a topic is fine. Never claim or imply a Piece was written without tools either.

# 08: Templates and writing guides per kind

**What to build:** Every drafting session follows the same rules. Write a template and a writing guide for Case Study, Concept, Comparison and Brief, in the same shape as the Repo Profile template and writing guide, drawn from the four pieces in ticket 07 and its notes. Every guide carries: first-party Sources only, every claim linked, short credited quotes only as long as the point needs, diagrams always redrawn in Mermaid with the original linked, code quoted only under the repository's licence, and the `verified:` rule. The Brief guide adds the 3-day freshness rule and the requirement that "why it matters" is HastCode's own analysis. CLAUDE.md gains pointers to the new documents beside the Repo Profile ones.

Spec: `../spec.md` ("Documents"). Research sections 3–4 (Google policy, copyright).

Note from ticket 10 (generated hero card): no Piece needs a hand-made image. Without `hero` the build draws the hero card from the title, kind and Vertical (or "Brief"); alt text is the title. `hero`/`hero_alt` stay optional, as in the Repo Profile template: a sourced image beside the Markdown replaces the card. So the templates leave both keys out (or show them commented out, like `docs/editorial/repo-profile-template.md`), and the guides say the card is checked on the Preview, not made.

**Blocked by:** 07 (First four pieces)

**Status:** resolved

- [x] Four templates and four writing guides exist and each cites the piece it was drawn from
- [x] Each template's section headings match what the build's editorial gates require for that kind (a skeleton built from it passes the structural gates once filled)
- [x] CLAUDE.md points to them

Resolution: `docs/editorial/{case-study,concept,comparison,brief}-template.md` and `…-writing-guide.md`, each drawn from and citing its ticket-07 piece, with a "What the first one taught us" section from the 07 check lists. Every guide carries the sourcing, quoting, diagram, code-licence, "we infer", TODO(author)-placement, hero-card and `verified:` rules; the Brief guide adds the 3-day Lead rule (counted from the first post), 250 words, the two sections, our own "Why it matters", and `explainers` = the linked slugs. Templates leave `hero`/`hero_alt` commented out (ticket 10 note). Each guide states its skeleton H2s in a "The skeleton" block; the template's `## ` lines are the same list, for 09. `builder/kind_templates_test.go` fills all four and builds them clean in production; a copy filled into a scratch content tree also built with zero errors through `cmd/builder` in preview and production. CLAUDE.md and README point to them.

Not changed (content is out of scope), for the Author: the published Brief has no AI-assistance note and two quotes, and three of the four descriptions exceed 150 characters (162–166); the guides now ask for the note, one quote and 150 characters.

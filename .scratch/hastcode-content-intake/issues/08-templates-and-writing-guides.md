# 08: Templates and writing guides per kind

**What to build:** Every drafting session follows the same rules. Write a template and a writing guide for Case Study, Concept, Comparison and Brief, in the same shape as the Repo Profile template and writing guide, drawn from the four pieces in ticket 07 and its notes. Every guide carries: first-party Sources only, every claim linked, short credited quotes only as long as the point needs, diagrams always redrawn in Mermaid with the original linked, code quoted only under the repository's licence, and the `verified:` rule. The Brief guide adds the 3-day freshness rule and the requirement that "why it matters" is HastCode's own analysis. CLAUDE.md gains pointers to the new documents beside the Repo Profile ones.

Spec: `../spec.md` ("Documents"). Research sections 3–4 (Google policy, copyright).

Note from ticket 10 (generated hero card): no Piece needs a hand-made image. Without `hero` the build draws the hero card from the title, kind and Vertical (or "Brief"); alt text is the title. `hero`/`hero_alt` stay optional, as in the Repo Profile template: a sourced image beside the Markdown replaces the card. So the templates leave both keys out (or show them commented out, like `docs/editorial/repo-profile-template.md`), and the guides say the card is checked on the Preview, not made.

**Blocked by:** 07 (First four pieces)

**Status:** ready-for-agent

- [ ] Four templates and four writing guides exist and each cites the piece it was drawn from
- [ ] Each template's section headings match what the build's editorial gates require for that kind (a skeleton built from it passes the structural gates once filled)
- [ ] CLAUDE.md points to them

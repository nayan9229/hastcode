---
status: accepted
date: 2026-11-12
---

# An Agent Surface: Markdown twins, llms.txt and sitemap.md, for agent tooling, not for AEO

Research (`docs/research/aeo-citation-signals.md` §8) found that Google ignores llms.txt, no answer engine has said it reads one, and measured citation impact is negligible. We still publish an **Agent Surface** (GLOSSARY): a Markdown twin beside every live Piece (`/briefs/<slug>.md`, `/<vertical>/<slug>.md`), `/llms.txt` (short orientation per llmstxt.org) and `/sitemap.md` (every live Piece), linked from each Piece's head (`rel="alternate" type="text/markdown"`) and a footer row. The reason is agent tooling, not rankings: a coding agent or LLM tool a Reader points at a page gets the Piece's own Markdown at a fraction of the HTML's tokens and quotes it accurately. It is cheap because every file is static and written by the builder, so it fits ADR-0002 and ADR-0004.

- **A twin is the Piece's written Markdown**, frontmatter removed, under a generated header (title, canonical URL, dates, author, "© HastCode. Quote with a link to <canonical>"). Mermaid diagrams stay as code; internal links become absolute. A Repo Profile's Fact Pack sections are in its HTML page only.
- **Twins follow their page**: written wherever the HTML page is written (Previews include drafts, Production never), so there is no second draft gate.
- **Twins stay out of search**: not in `sitemap.xml`, the news sitemap or IndexNow, and served with `X-Robots-Tag: noindex` on both Targets, so Google never indexes a `.md` as a duplicate of its page.
- **Bot policy is unchanged**: `robots.txt` already allows every crawler, training bots included; the twins do not change who may fetch.

## Considered options

- **A JSON discovery document at `/api`, an MCP endpoint, a search API** (what aihero.dev ships): hastcode has no API or product catalogue, `/api/*` is reserved for the signup function, and an MCP server needs a Worker, against ADR-0002. Rejected.
- **Markdown by `Accept: text/markdown` negotiation**: not possible on a static Target. Rejected.
- **Nothing**: cheapest, but every agent read then pays for the full HTML page. Rejected as a small, static cost.

## Consequences

- Twin URLs are public contracts once agents link them; retiring one is a 410, like any retired path.
- A twin carries no ads; agent reads earn nothing, as they never did.
- The head link and footer row add 280 bytes of chrome to every page. Three live pages sat within 150 bytes of the 50 KB HTML gate, so the gate moved to 52 KB (2026-11-13): the Author chose headroom over trimming verified pages. The daily drafting guidance moved with it (51,500 bytes in a Preview).

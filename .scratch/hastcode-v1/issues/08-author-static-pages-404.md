# 08: Author page, static pages and 404

**What to build:** The transparency surface Google News requires. An Author page at `/authors/{slug}/` from an Author Markdown file (name, bio, photo, links) with `Person` JSON-LD, listing that Author's latest pieces. Static pages at `/about/`, `/contact/` and `/editorial-policy/` built from Markdown under the pages content directory, the last carrying the sourcing rule, AI-assistance disclosure wording and corrections policy. A static 404 page at the path Vercel's fallback route will point at. Every Explainer and Brief links its Author's page from the byline.

**Blocked by:** 05 (Explainer head and structured data)

**Status:** resolved

- [x] Build-seam tests: Author fixture → page with name, bio, photo, links and `Person` JSON-LD whose `url`/`sameAs` match the Explainer JSON-LD author reference
- [x] Static page fixtures → pages at the expected paths with title, body and breadcrumb
- [x] 404 page exists in output, is noindex, and links home, Verticals and search
- [x] Byline on article pages links to the Author page
- [x] A piece referencing an unknown Author slug is a build error (Explainers: `TestFrontmatter_UnknownAuthorIsABuildError`; Briefs must reuse the rule in ticket 06)

## Notes (implementation)

- Author page lists Explainers now via `s.authorPieces(slug)` in `builder/step_render_authors.go`; Briefs (06) add one loop there with `Type: "Brief"`.
- Briefs use `@Byline(views.Link{...})` (`partial_byline.templ`) so the byline links the Author page.
- 404 is written at `/404.html`; ticket 10 adds the Vercel fallback route to it.
- Real content: `content/pages/{about,contact,editorial-policy}.md` and `content/authors/nayan.md` carry `TODO(author):` markers the Author must fill before cutover.

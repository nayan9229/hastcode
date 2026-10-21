# 13: Site search with Pagefind

**What to build:** A Reader opens `/search/`, types a keyword and gets Explainers and Briefs without any backend. Pagefind indexes the output directory as a post-build step; the search page lazy-loads Pagefind's UI on first focus, works as a plain page without it, and respects dark mode. Production builds fail if the index is missing; preview and local builds may skip it with a flag. Indexed content excludes nav, footer, ad slots and the newsletter form.

**Blocked by:** 07 (Listing pages and internal linking), 11 (Design system, dark mode and brand)

**Status:** resolved

- [x] Build-seam test with the post-step enabled: Pagefind index directory exists in output and the search page references it; with the post-step skipped in production mode the build fails
- [x] Search page exists, is noindex, has a visible input and no-JS explanation
- [x] Article templates mark the body as the indexable region and exclude chrome
- [x] Documented install/run of Pagefind for the Vercel build command

**Resolution notes:** Pagefind runs inside the builder as the last write step (`step_write_search_index.go`, via `Config.Node` + `ToolsDir`), so the Vercel build command is just `make build` after `npm ci` (README "Search"). `-skip-search-index` is preview-only; production fails with rule `search-index`. Index restricted to `article[data-pagefind-body]` on Explainers and Briefs; diagrams, Sources, Brief source line, newsletter, ads and reading lists ignored. `/search/` lazy-loads the Pagefind modular UI on first focus; checked end-to-end in headless Chrome (light, dark, 375 px).

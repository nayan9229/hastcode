# 01: Walking skeleton: one Explainer becomes one static page

**What to build:** The Author drops a single valid Explainer (Markdown with frontmatter) into the content tree, runs the builder with a content directory, a configuration (site URL, Vertical list, Author, flags) and an output directory, and gets a Vercel Build Output API v3 tree containing that Explainer at `/{vertical}/{slug}/index.html` plus a minimal `config.json`. The page is complete HTML rendered through a templ layout and Explainer template: title, description, Vertical name, Author name, body Markdown rendered to HTML, visible published and last-reviewed dates. This is the spine every later ticket hangs on; it also establishes the test harness at the Build seam.

Spec: `../spec.md`. Respect ADR-0002 (Go is a build tool; no server).

**Blocked by:** None (can start immediately)

**Status:** resolved

- [x] Go module with a `builder` command taking content dir, output dir, config, and a production/preview mode flag; exits non-zero and prints errors when the build fails
- [x] Frontmatter parsed per spec (title, description, kind, vertical, tags, published, reviewed, author, hero, hero_alt, related, draft); unknown Vertical is a build error
- [x] Output written under the Build Output API layout (`static/` directory-index HTML, `config.json` with version 3)
- [x] templ is wired with a `go generate` step; layout + Explainer page components exist with the minimal head (title, meta description, `lang="en"`, viewport)
- [x] A test-support package provides: write fixture content tree to a temp dir, run the build, load and parse an output HTML file, decode JSON-LD blocks, read `config.json`
- [x] First Build-seam test: valid Explainer fixture → page exists at expected path, title/description/dates/body present when parsed
- [x] Build of the fixture tree completes in well under a second

## Comments

Implemented on `ticket/01-walking-skeleton`. Notes for later tickets:

- The Author's name comes from `content/authors/<slug>.md` (spec content model), not from `site.yaml`; the site config holds `site_url` and `verticals`. An unknown Author slug is already a build error (ticket 08 only needs to add its page).
- Frontmatter is decoded strictly (unknown keys fail); `published`/`reviewed` must be RFC 3339 with offset; `kind` must be case-study|concept|comparison. `draft` is parsed but not yet acted on (ticket 02).
- Pipeline, templ and test-support conventions are in README.md "Codebase map" / "Conventions".

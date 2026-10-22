# hastcode

Static-first Go publication for hastcode.com ("HastCode"). The Go program is a
**build tool**, not a server (ADR-0002): it reads Markdown from `content/`,
renders templ components, and writes a Vercel Build Output API v3 tree to
`.vercel/output/`. Domain terms (Explainer, Brief, Vertical, Author, Review,
Reader) are defined in `GLOSSARY.md`; use them in code and test names.

## Quick start

templ is a Go tool dependency (`go.mod` `tool` directive), so no global install
is needed. Generated `*_templ.go` files are **gitignored**: generate before
compiling, vetting or testing anything. Mermaid diagrams need Node 22+ and the
pinned npm packages (see "Mermaid diagrams" below).

```sh
npm ci          # Node build tools from package.json (Mermaid renderer, Pagefind)
make generate   # go generate ./...  (runs `go tool templ generate`)
make test       # generate + go test ./...
make vet        # generate + go vet ./...
make build      # generate + go run ./cmd/builder (content/ -> .vercel/output/)
make dev        # build (preview mode) + serve site and subscribe on http://localhost:3000
```

## Working on the site

**Run the tests.** `npm ci` once, then `make test` (templ generate + `go
test ./...`, about a minute; `go test -short ./...` skips the slow
1,000-entry news-sitemap case). Build-seam tests skip, with the reason,
when Node or the npm packages are missing, so run `npm ci` first or most
of them will not run. `make vet` before committing. CI runs `make test`
on every push to a branch other than `main` (and on pull requests from
forks) and before every production deploy.

**Check pages locally.** `make dev` builds in **preview mode** (drafts
included) and serves `.vercel/output` plus the subscribe function on
http://localhost:3000 (`PORT=3001 make dev` for another port) behind the
built `config.json` routes: `/video/*` 410, `/home` 301, `/x` -> `/x/`
308, missing paths 404 with the 404 page. `DEV_MODE=production make dev`
shows exactly what production will publish. Signups go to whatever
`BUTTONDOWN_API_KEY` says: the shell first, else `.env` (the setup wizard
writes the *test* newsletter's key there), else none (every POST 502s).
To try the form without Buttondown, point it at any fake that answers
`POST /v1/subscribers` with 201:
`BUTTONDOWN_API_KEY=x BUTTONDOWN_BASE_URL=http://127.0.0.1:8089 make dev`.

**Draft and preview flow.** Write the piece with `draft: true` in its
frontmatter on a branch. `make dev` shows it locally; push the branch
(a pull request is optional), and the Preview workflow deploys a preview build
(drafts included, behind Vercel Authentication) and posts its URL in the
run summary, then runs Lighthouse CI against it. Production builds never
contain drafts.

**Publish flow.** Remove `draft: true`, merge the pull request to `main`.
The Production workflow runs `make test`, builds in production mode,
deploys with `vercel deploy --prebuilt --prod`, then pings IndexNow with
the pages that changed (after the domain cutover). Feeds, sitemaps and the
news sitemap update in the same build. A build error (failed editorial
gate, broken link, oversized page, a page without JSON-LD, a leftover
`TODO(author)`) fails the workflow and nothing deploys. The shipped
`content/` still has `TODO(author)` placeholders, so production builds
fail until the Author fills them (see "Content layout"); `make dev` and
previews work meanwhile.

**Add a Vertical.** Add one entry to `site.yaml` `verticals` (`slug`,
`name`, `description`, `order`), then put its Explainers in
`content/explainers/<slug>/`. The build gives it a hub at `/<slug>/`, a
header nav link, a home page entry and `/<slug>/feed.xml`. No code change.
The slug is a top-level URL forever: pick it once.

**Environment variables.** Nothing in the build reads secrets. The
subscribe function reads two (Vercel project settings, see "Subscribe
function"): `BUTTONDOWN_API_KEY` (required, Sensitive; production list in
Production, test list in Preview) and `BUTTONDOWN_BASE_URL` (optional,
default `https://api.buttondown.com`). CI uses GitHub secrets
`VERCEL_TOKEN`, `VERCEL_ORG_ID`, `VERCEL_PROJECT_ID`,
`VERCEL_AUTOMATION_BYPASS_SECRET` and the repository variable
`INDEXNOW_ENABLED`. `make build` itself also reads `MODE`, `FLAGS`, `OUT`
from the environment (the preview job sets `MODE=preview`).

**Deploy.** See "Deploying" below; first-time setup is
`scripts/setup-deploy.sh`.

The Makefile exports `CGO_ENABLED=0` and `GOFLAGS=-tags=nodynamic`: the AVIF
encoder (gen2brain/avif) then runs as WASM in-process and never loads a
system libavif, so image output is byte-identical everywhere. Use the same settings when calling `go` directly.

Builder flags: `-content` (default `content`), `-config` (default `site.yaml`),
`-out` (default `.vercel/output`), `-mode production|preview`, `-tools` (dir
with `package.json`/`node_modules`, default `.`), `-cache` (build cache, default
`<tools>/node_modules/.cache/hastcode`), `-image-cache` (overrides
`images.cache_dir`), `-skip-search-index` (preview only: no Pagefind run; a
production build with it fails). It prints every
error and warning and exits non-zero if the build failed.

Content layout (spec "Content model"): `content/explainers/<vertical>/<slug>.md`,
`content/authors/<slug>.md` (frontmatter `kind` = `person` (default) or
`organization` for a team byline, `name`, `same_as` list, optional `photo`
raster beside it + `photo_alt`; body = bio), `content/pages/<slug>.md`
(static pages at `/<slug>/`: `title`, `description`, Markdown body). Briefs come
in `content/briefs/<slug>.md`. Optional frontmatter `sources` (list of
`title` + `url`) adds to the inline links: an Explainer's are appended to
its `## Sources` list (page and feed), a Brief's are listed after its
`source_url`; a URL already linked there is not repeated
(`step_transform_sources.go`). Lines marked `TODO(author):` in `content/` are facts only the
Author can supply. A **production** build fails while any written page or
feed still contains `TODO(author)` (rule `author-placeholder`, naming the
page and the content file to fix; `step_write_author_placeholders.go`), so
the Production workflow stays red until the Author page and the About,
Contact and Editorial Policy pages are filled in. Preview builds (`make
dev`, `make build MODE=preview`, the Preview workflow) only warn, so the
drafts can be read in place.
Add a Vertical by adding one entry to `site.yaml`.

Heroes: the Explainer names a raster beside its Markdown in `hero` (and
`hero_alt`). The build fails if it is missing, does not decode, or is narrower
than 1200 px. It writes AVIF and JPEG at 480/800/1000/1200/1472 px wide
(never upscaled past the source; `heroWidths`), plus 1200 px JPEG crops at
16:9, 4:3 and 1:1 (JSON-LD `image`) and 1200x630 (`og:image`,
`twitter:image`), to `/assets/heroes/<name>-<w>x<h>.<hash>.<ext>` (hash of
the file bytes, so immutable). No WebP: AVIF covers nearly every browser
and the JPEG covers the rest. Encoder settings
(`docs/research/image-optimisation-targets.md`): AVIF 4:4:4 chroma, speed
8, quality 60 below 1200 px and 55 at 1200 px and above (`heroQuality`;
the 1200 px file is the mobile LCP image); stdlib JPEG quality 82. When a
hero's 1200 px AVIF is over 150 KB the build warns (rule `hero-weight`,
naming the piece): re-export a lighter hero. Encoded variants are cached in
`images.cache_dir` (default `node_modules/.cache/hastcode/images`, which CI
persists between deploys), keyed on the source bytes, the encoder settings
and `settingsVersion` (`internal/images`; bump it when encoding changes in
a way the settings do not capture, and every cached variant re-encodes
once); an unchanged hero is never re-encoded. `images.fast: true` raises
AVIF speed to 10 and drops JPEG to quality 75 for quicker encodes (tests);
AVIF quality tiers and chroma stay the same.

Hero source rules: drop the file straight from a phone or an export tool;
the build reads its bytes, not its name.

- **PNG or JPEG only.** Anything else, including a WebP (even one renamed
  `.png`), fails with `hero-format`: Go's WebP decoders return the wrong
  colour range, so a WebP source would ship washed out.
- **sRGB.** Profiles are never converted and the outputs are untagged
  (read as sRGB), so a source carrying a non-sRGB ICC profile (Display P3
  from an iPhone, Adobe RGB) would look desaturated. The build warns with
  `hero-colour`, naming the piece and the profile: re-export it as sRGB.
  Untagged, sRGB-tagged and grey sources are fine.
- **At least 1200 px wide** once upright (`hero-too-narrow`).
- **EXIF orientation is honoured**: a sideways-stored phone photo is
  rotated upright before resizing; no need to pre-rotate.

On the page the hero `<img>` carries `fetchpriority="high"` and no
`loading` attribute (it is the LCP image), and `<head>` keeps one AVIF
preload with the same `imagesrcset`/`imagesizes`.

## Codebase map

| Path | What lives there |
|------|------------------|
| `cmd/builder/` | CLI only: flags, `LoadConfig`, print `Result`, exit code. No logic. |
| `builder/` | The single deep module. Public API: `Config`, `Vertical`, `Mode`, `LoadConfig`, `Build(cfg, contentDir, outDir) Result`, `Result`/`Problem`. Everything else unexported. |
| `builder/pipeline.go` | Phases, `registerStep`, the `site` state struct, `fail`/`warn`, the "skip failed files" filter `passed` behind `publishedExplainers()`/`publishedBriefs()`/`publishedPages()`. |
| `builder/step_*.go` | One file per pipeline step (load, validate, transform, render, write). |
| `builder/urls.go` | The URL scheme: one func per page type (`explainerURL`, …), every fixed page, feed and sitemap path, the asset directories, and the retired paths (`retiredSections`, `homeAliasURL`). |
| `builder/render.go` | Output helpers: `writePage(urlPath, templComponent)` incl. the 50 KB HTML size gate (rule `html-size`), `writeXML`, `isWritten`, rule `markdown`. The output reset is `step_render_reset.go`. |
| `builder/load.go` | `eachMarkdown` and rule `frontmatter`, shared by the `step_load_*.go` steps. |
| `builder/*_test.go` | Build-seam tests (package `builder_test`), one file per feature area. |
| `internal/content/` | Parse content files (frontmatter + body) into `Explainer`, `Brief`, `Author`, `Page`. No config or output knowledge. |
| `internal/markdown/` | goldmark `Parse` -> `*Document` (AST + source) and `Render(doc)`. Steps inspect/rewrite the AST between the two. `doc.Diagrams()` finds ```` ```mermaid ```` fences (title from a leading `%% title:` line); `doc.ReplaceDiagram(d, html)` swaps one for builder HTML. |
| `internal/mermaid/` | Mermaid -> SVG: one batch Node process (`render.mjs`, embedded) per build, type/syntax checks, post-processing, on-disk cache by source hash, `PrefixIDs`. |
| `internal/views/` | templ components + view-model structs (`models.go`). |
| `internal/images/` | Resize/crop/encode (AVIF, JPEG; per-variant quality) through a worker pool with an on-disk cache. Knows nothing of heroes or URLs. |
| `internal/views/brand.go`, `og_badge.go` | Brand constants (`PublicationName`, `Tagline`, `MarkPath`), the light/dark palettes, the generated SVG mark, the generated 1200x630 PNG OG badge (`OGBadge()`, Go fonts, `/assets/og-badge.<hash>.png`), `StaticAssets()`. |
| `internal/images/source.go` | Source rules, pure Go: `Inspect(path) Info` sniffs PNG/JPEG from the bytes (`ErrWebP` for WebP), reads the EXIF orientation (upright `Width`/`Height`) and the ICC profile description (`Profile`, `SRGB`); the encoder applies the orientation right after decoding. |
| `internal/buildtest/hero_sources.go` | Hand-built hero source fixtures: `HeroWebP`, `HeroJPEG`, `WithEXIFOrientation`, `ICCProfile` + `WithJPEGICCProfile`/`WithPNGICCProfile`, `Quadrants`, and `DecodePixels` for pixel checks. |
| `internal/views/site.css`, `stylesheet.go`, `theme.js` | The one stylesheet (hashed file, ADR-0003) and the inline theme-toggle script. |
| `builder/step_render_authors.go`, `step_{validate,transform,render}_author_photos.go` | Author pages (`/authors/<slug>/`: h-card, photo, profiles, latest pieces via `s.authorPieces`, top-level `Person` or `Organization` by the Author's `kind`) and the square photo (validate, encode, write). |
| `builder/step_{load_static_pages,validate_static_page_paths,render_static_pages}.go`, `step_render_not_found.go` | Static pages from `content/pages/` (slug may not take a Vertical or reserved path, rule `page-path`) and the noindex `/404.html`. |
| `builder/step_render_home.go`, `step_render_verticals.go`, `step_render_kinds.go` | Listings: home `/` (6 latest Explainers, 5 latest Briefs, every Vertical with its `description`, kind links), Vertical hubs `/{vertical}/` and kind pages `/explainers/{kind}/` (both paginated, `s.listingPage`). |
| `builder/listing.go`, `builder/related.go` | Shared listing helpers: `newestFirst` (via `newestExplainers`/`newestBriefs`), `s.explainerItems`/`s.briefItems` (-> `views.ListingItem` with the 480 px hero thumbnail), `explainerKinds` (badge and kind-page names); `s.readingLists(e)` ("More in {Vertical}" and "Related" at the end of an Explainer). |
| `builder/step_render_link_gate.go` | The site-wide broken internal link gate (rule `broken-link`), last in the render phase. |
| `builder/step_render_search.go`, `step_write_search_index.go` | The noindex `/search/` page and the Pagefind post-step (rule `search-index`) that indexes the written output into `/pagefind/`. |
| `internal/pagefind/` | Runs the Pagefind CLI (`node node_modules/pagefind/lib/runner/bin.cjs --site …`) and reads the bundle's page count. |
| `internal/views/page_search.templ`, `search.js` | Search page and its inline lazy loader (Pagefind modular UI on first focus). |
| `builder/chrome.go` | `s.page(urlPath, title, description)`: every page's `views.Page` incl. canonical, default Open Graph, header/footer nav, ads flag and Plausible domain. |
| `builder/head.go` | Head + JSON-LD shared by page types: `s.articlePage(articleHead{Type, URL, Headline, …, Breadcrumbs})` (Article/NewsArticle + BreadcrumbList + og:type article with the 1200x630 hero crop + visible breadcrumb), and the schema builders `s.person(author)`, `s.organization()`, `s.breadcrumbs(links)`, `s.absURL(path)`. |
| `internal/schema/` | schema.org JSON-LD node types (`Article` with `TypeArticle`/`TypeNewsArticle`, `Person`, `Organization`, `BreadcrumbList`, `WebSite` + `NewSearchAction`). Plain data; `views.StructuredData` marshals each `Page.StructuredData` block into `<head>`. |
| `internal/vercel/` | Build Output API writer (`static/` files, `config.json` with `Route`s) and `NewRouter`, the local router that applies a built `config.json`'s routes. |
| `internal/buildtest/` | Test support for the Build seam (fixtures, run build, parse output). |
| `subscribe/` | The subscribe function: `New(Config) http.Handler`, `ConfigFromEnv()`, form contract constants (`Path`, `SubscribedPath`, `HoneypotField`). Buttondown client in `buttondown.go`, per-IP limiter in `ratelimit.go`. Subscribe-seam tests in `subscribe_test.go`. |
| `api/subscribe.go` | Vercel Go function entry (package `handler`, `func Handler`). Wiring only. No other `.go` files here: each becomes a function. |
| `cmd/subscribe/` | Local run of the subscribe function plus the built site, behind the built `config.json` routes. |
| `builder/step_write_vercel_config.go` | The routing table (`routes()`): cache headers, 410s, `/home` 301, trailing-slash 308s, 404 fallback. |
| `builder/step_write_page_index.go`, `step_write_page_manifest.go` | Reads every written HTML page back (canonical, noindex, `<time>`s, hash of `<main>`) and diffs it against the previous build's page manifest in the build cache; the last write step saves the new manifest when the build passed. |
| `builder/step_write_feeds.go`, `step_write_sitemap.go`, `step_write_news_sitemap.go`, `step_write_robots.go`, `step_validate_indexnow_key.go`, `step_write_indexnow.go` | Full-text RSS feeds, `/sitemap.xml` + `/news-sitemap.xml`, `/robots.txt`, the IndexNow key file and changed-URL artefact. See "Feeds, sitemaps, robots and IndexNow". |
| `site.yaml` | Site config: `site_url`, `verticals`, `plausible_domain` (blank = no analytics), `ads_enabled`, `adsense_client` (required when ads are on), `indexnow_key` (public, stable; blank = IndexNow off), `same_as` (home Organization profiles). |
| `package.json` | Pinned Node build tools (`beautiful-mermaid`, `pagefind`). Not a JS app. |

Planned homes for parallel tickets: editorial gates -> `builder/step_validate_*.go`;
hero images -> `builder/step_validate_heroes.go`, `step_transform_heroes.go`, `step_render_hero_assets.go` +
`internal/images/`; Mermaid -> `builder/step_transform_diagrams.go` +
`internal/mermaid/`; design system -> `internal/views/layout_*.templ` + CSS.

## Conventions

**Adding a pipeline step.** Create `builder/step_<phase>_<what>.go` with an
`init()` that calls `registerStep(phase, order, name, fn)`: one step per
file, so a feature spanning phases (heroes: validate, encode, write) is one
file per phase. Do not edit a central list. A step returning an error stops
the build with rule `build-step`. Phases run in order: `phaseLoad` -> `phaseValidate` ->
`phaseTransform` -> `phaseRender` -> `phaseWrite`; within a phase, by `order`
(use gaps: 10, 20, 30…). Existing orders: load authors 10, load explainers 20,
load static pages 25, load briefs 25, exclude drafts (production only; Explainers
and Briefs) 30; validate ads config 5, references 10, Brief references 12,
unique slugs 20, static page paths 25, Explainer skeleton 30, Brief skeleton
35, section-length warning 40, Review-overdue warning 50, validate heroes 50
(Explainers and Briefs), author photos 55; transform render diagrams 10,
rewrite Brief Explainer links 20, merge Explainer frontmatter sources 30,
encode heroes 50, encode author photos 55;
render reset 0, write site assets 10, write hero assets 15, write author
photos 16, render explainers 20, render briefs 30, render author pages 40,
render briefs index 40, render static pages 45, render vertical hubs 50,
render kind pages 52, render home 54, render subscribed page 60,
render search page 62, render 404 page 65, check internal links 1000
(keep it last: it reads every written page); validate IndexNow key 6; write
phase: index written pages 5 (reads the output back; later write steps use
`s.index`), feeds 20, sitemap 30, news sitemap 35, robots.txt 40, IndexNow
key + changed URLs 50, author placeholders 60, config.json 100, index search 800 (Pagefind over the
final output), save page manifest 900 (keep it last).

**Internal links.** After every page is written, `checkInternalLinks`
reads each HTML file in the output and resolves every `href` that stays on
the site (root-relative, page-relative, or absolute on `site_url`'s host;
query and fragment ignored). It must name a written file or a directory
index (`/a/` -> `/a/index.html`; `/a` without the slash is broken), or a
path the build does not write itself (`linkTargetsNotWritten`:
`/api/subscribe`, `/feed.xml`, `/{vertical}/feed.xml`,
`/briefs/feed.xml`, `/sitemap.xml`, `/news-sitemap.xml`; feeds and
sitemaps are written in the write phase, after the gate, so their entries
stay). Otherwise it is a `broken-link`
error naming the page and the target, attributed to the page's content
file ("" for home and listings), and the page is removed from the output.
Header and footer links are checked on every page, so the About, Contact
and Editorial Policy pages must exist (`buildtest.NewTree()` includes them).
Brief Explainer links and `related` slugs are checked earlier, at the
source, under the same rule.

**Listings.** Every listing shows `views.ListingItem`s through
`@PieceList(class, items, headingLevel)` (`partial_piece_list.templ`):
480 px hero thumbnail (AVIF/JPEG, lazy, size reserved), kind badge
(Explainers), title, description, date, as h-entry. Listings use
`s.publishedExplainers()`/`s.publishedBriefs()`, so drafts never appear in
production. Verticals come only from `site.yaml` (`slug`, `name`,
`description`, `order`): a new entry gets a hub, a header nav link and a
home entry.

**Briefs.** `content/briefs/<slug>.md` -> `/briefs/<slug>/` (NewsArticle via
`s.articlePage`). Skeleton (`step_validate_brief_skeleton.go`): ≤ 250 words incl.
headings (`brief-length`), H2s exactly "What happened" then "Why it matters"
(`brief-sections`; an intro is allowed), absolute http(s) `source_url`
(`brief-source`), ≥ 1 body link to an Explainer and `explainers` frontmatter
= the linked slugs (`brief-explainer-link`). Body Explainer links are a bare
slug (`[text](raft-leader-election)`) or `/{vertical}/{slug}` (slash and
`#fragment` optional); an unknown target is `broken-link`, a known one is
rewritten to the canonical trailing-slash URL (`brief_links.go`). The first
ad slot follows "What happened".

**Paginated listings.** `paginate(items, listingPageSize, baseURL)` in
`builder/pagination.go` returns `[]listingPage[T]{Items, URL, Pagination}`
(items already in display order; page 1 at `baseURL`, page n at
`baseURL + "page/n/"`; always ≥ 1 page, never an empty later page). Pass
`lp.Pagination` to the view and render `@PaginationNav(p)`
(`partial_pagination.templ`: `<a rel="prev">`/`<a rel="next">`, nothing on
a single page). Example: `step_render_briefs_index.go`.
Date rules use `s.cfg.Now` (the build clock; `Config.Now`, zero = wall clock),
never `time.Now()`. Gates inspect the body via `markdown.Document.Sections()`
(intro + H2 sections), `Words`, `Text`, and `markdown.Find` with `IsLink`/`IsTable`.
`fn(s *site) error` returns an error only for fatal I/O (aborts the build).
Content problems go through `s.fail(file, rule, fmt, args…)` (error; marks that
file failed so its page is not written) or `s.warn(...)` (never fails). `file`
is content-relative with `/` separators (`e.File`); `rule` is a short kebab-case
constant declared next to the step. Render/write steps iterate
`s.publishedExplainers()`, never `s.explainers`, so failed pages are skipped.

**Adding a templ component.** File name = component family prefix:
`layout*.templ` (shell, head, header, footer), `article*.templ` (Explainer,
Brief pages), `listing*.templ`, `page*.templ`, and **one file per partial**
`partial_<name>.templ`. Components take view-model structs from
`internal/views/models.go`, never `content.*` types; the builder maps content to
view models (e.g. `explainerPage` in `step_render_explainers.go`). Article markup
uses microformats2 h-entry classes (`h-entry`, `p-name`, `p-category`,
`p-author`, `dt-published`, `dt-updated`, `e-content`); tests select on them, so
keep them. Run `make generate` after editing `.templ` files.

**Adding a page type.** Add a URL func in `builder/urls.go`, a view model in
`models.go`, a templ file, and a `step_render_<type>.go` that calls
`s.writePage(url, component)`. Build the page's `views.Page` with
`s.page(url, title, description)` (never a bare `views.Page{}`) so the
canonical, header and footer nav are filled, and wrap the component in
`@Layout(page)`. Article pages (Explainer, Brief) use
`s.articlePage(articleHead{...})` instead, which calls `s.page` and adds
the hero preload, Open Graph article tags, the breadcrumb trail and the
Article/NewsArticle + BreadcrumbList JSON-LD. Other pages set
`page.Breadcrumbs` (Home first, the page last; rendered above the content
by `Layout`) and append `s.breadcrumbs(page.Breadcrumbs)` plus any
top-level block (`Person`/`Organization` with `Context: schema.Context`)
to `page.StructuredData`.

**Head and structured data.** Every page's `<head>` carries `lang="en"`,
robots `max-image-preview:large` (or `noindex` when `Page.NoIndex`; only
`/subscribed/`, `/404.html` and `/search/` are noindex in production; the 404 has no
canonical), an absolute trailing-slash
canonical, Open Graph + Twitter card and both theme-colors. Every page has
an `og:image` ≥ 1200 px (repeated as `twitter:image`) and `twitter:card`
`summary_large_image`: an article its 1200x630 hero crop, every other page
the 1200x630 OG badge (`s.page` sets
it; drawn at build time from the wordmark, tagline and light palette). JSON-LD URLs
are absolute (`s.absURL`), dates are RFC 3339 with the frontmatter's
offset (the same instants as the visible `<time datetime>`), and an
Explainer's `dateModified` is its `reviewed` date, never the build time.
Author `url` is `/authors/<slug>/`; `sameAs` comes from the Author file's
`same_as` list. `s.author` emits `@type` `Person`, or `Organization` when the
Author's `kind` is `organization` (unknown kinds fail with `author-kind`;
organizations get no `author-photo` warning). Article bylines use `@Byline(views.Link{Name, URL: authorURL(slug)})`
(`partial_byline.templ`), which links the Author page. No `FAQPage`, `HowTo`, `Speakable` or `TechArticle`.
`writePage` fails an indexable page with no `application/ld+json` block
(rule `jsonld-missing`; the page is not written), so a new page type needs
at least its BreadcrumbList.

**Design system.** One stylesheet, `internal/views/site.css`, served as
`/assets/site.<hash>.css` (ADR-0003); no inline `<style>`, no webfonts, no
framework, no `url()`. Colours exist only as tokens generated from
`brand.go`: `--bg`, `--fg`, `--accent` (spec values), `--link` (AA-safe link
ink), `--muted`, `--border`, `--surface`; dark values apply under
`prefers-color-scheme: dark` and under `<html data-theme="dark">` (the
toggle's override; `data-theme="light"` forces light). Use the tokens, never
hex values, in CSS and in generated SVG (Mermaid: `currentColor` resolves to
`--fg` inside `figure`). Layout: `.wrap` is the page column (16 px gutter on
phones, 72rem max); `main > article` gets the reading measure (68ch). Element
styles already cover `figure`/`figcaption`, `img`/`svg` (never wider than the
column), `table` (scrolls inside itself), `pre`/`code`, `blockquote`, forms
and `button`. Generic hooks for partials: `.box` (bordered surface panel,
e.g. TL;DR, newsletter), `.muted`, `.inline-list` (breadcrumbs, tags,
links), `.visually-hidden`. Add component-specific rules to `site.css` under
a comment naming the partial, one declaration block per line. The inline
scripts are the theme toggle (≤ 1 KB, in `<head>`), the newsletter form's
enhancement (≤ 1 KB, inside `.newsletter`) and the search page's loader
(≤ 1.5 KB, `search.js`; Pagefind itself loads only on interaction, so it is
not part of any page's initial JS); analytics and Pagefind are the
only other scripts allowed (`TestThemeToggle_IsTheOnlyScriptAndUnder1KB`,
`TestNewsletter_EnhancementScriptIsSmall`), plus the AdSense library and its
per-unit push only when `ads_enabled`. Total client JS per page stays
≤ 60 KB **uncompressed, analytics included**: the authority is the
Build-seam test `TestJSBudget_EveryPageStaysWithin60KB`, which sums every
script a page loads initially (inline code, same-origin files, and an
upper bound for each third-party script: Plausible 4 KB, measured 2,841
bytes). Give any new third-party script a size there. Lighthouse CI's
script check is only a secondary guard (it measures compressed transfer
and must block Plausible on previews).

**Weight, ads, analytics.** Every page written through `s.writePage` is
checked against the 50 KB uncompressed HTML gate (the stylesheet is external
and not counted, ADR-0003); an oversized page is a site-wide `html-size`
error naming the URL and is not written. Ads come from one component,
`@AdSlot(page.Ads, views.AdAfterTLDR|views.AdAfterWhatHappened|views.AdAfterBody)`
(`partial_ad_slot.templ`), controlled by `ads_enabled`. With ads off it
renders **nothing and reserves no space** (v1 pages carry no blank boxes;
turning ads on later is a deliberate layout change, reviewed on a
preview). With ads on it is an AdSense unit in an `aside.ad-slot` whose
height is fixed by `--ad-slot-height`, so the unit filling in causes no
layout shift; it carries `data-pagefind-ignore="all"`. Article layouts
render the body in two halves (`Document.SplitAfterSection` +
`markdown.RenderBlocks`: after "TL;DR" for an Explainer, after "What
happened" for a Brief) with the first slot between them and the second
after `.e-content`. Plausible comes
from `layout_analytics.templ` in every `<head>` when `plausible_domain` is set.
Everything must work without JS.

**Tests.** Only two seams are tested (spec "Testing Decisions"): the Build seam
and the Subscribe HTTP seam. Build-seam tests live in `builder/` as package
`builder_test`, import only `builder` and `internal/buildtest`, and assert on
parsed output. Helpers in `internal/buildtest`:

- Fixtures: `ValidTree()`, `NewTree()` (default Author + `TransparencyPages()`, which the footer links), `Tree.WithExplainer(e)`, `Tree.Without(path)`, `Tree.WithPage(p)`,
  `Tree.WithAuthor(a)`, `Tree.With(path, bytes)`; `ValidExplainer()` (meets the
  full editorial skeleton: TL;DR, Mermaid fence, Sources, 1600x900 PNG hero),
  `ValidComparison()` (adds a `## Trade-offs` table), `ValidBrief()` +
  `Tree.WithBrief(b)` (meets the Brief skeleton; links `ValidExplainer()` by
  slug, so use it with `ValidTree()`), `DefaultAuthor()`,
  `HeroPNG(w, h)`, `HeroPNGFilled(w, h, colour)`, `Config()` (launch Verticals,
  production; set `cfg.Now` to pin the clock; fast hero encoding into the
  shared `SharedImageCacheDir()` so each fixture hero is encoded once; set
  `cfg.Images.CacheDir = t.TempDir()` to test caching itself).
  Mutate a copy of `ValidExplainer()` to make an invalid fixture.
- Run: `Build(t, cfg, tree) (outDir, Result)`, `MustBuild(t, cfg, tree) outDir`.
  `Config()` points `ToolsDir` at the repo root; `Build` gives each build a
  fresh temp `CacheDir` unless the test sets one, and **skips** the test (with
  the reason) when the tree has a Mermaid fence but `node` or
  `node_modules/beautiful-mermaid` is missing. `ValidExplainer()` has a fence,
  so most build tests need `npm ci` to run rather than skip.
- Serve: `Serve(t, outDir)` runs the output through the local router (the
  built `config.json` routes, `static/`, a stub subscribe function counting
  `FunctionCalls`); `srv.Get(t, path)` / `srv.Do(t, method, path, body)`
  never follow redirects. `URLPaths(t, outDir)` lists every output URL.
- Assert: `Page(t, outDir, "/vertical/slug/")` (goquery document),
  `Exists(outDir, urlPath)`, `StaticPath(outDir, urlPath)`, `JSONLD(t, doc)`,
  `JSONLDOfType(t, doc, "Article")`, `VercelConfig(t, outDir)`,
  `DecodeImage(t, outDir, urlPath)` / `Images(t, outDir)` (decode output images
  and measure them), `Srcset(t, attr)`,
  `RequireError(t, res, file, substring)`, `RequireWarning(...)`,
  `RequireErrorRule(t, res, file, rule, substring)`, `RequireWarningRule(...)`,
  `RequireNoProblemRule(t, res, file, rule)`, `FilesContaining(t, outDir, s)`
  (prove a piece appears in no output file).

Subscribe-seam tests live in `subscribe/subscribe_test.go` (package
`subscribe_test`): they build the handler with `subscribe.New(Config{APIKey,
BaseURL: fake.URL})`, send form posts through `ServeHTTP`, and assert status,
`Location`/JSON body and what the `httptest` fake Buttondown received.

## Subscribe function

The newsletter form (`internal/views/partial_newsletter.templ`, at the end of
every Explainer) is a plain `POST /api/subscribe` with `email`, `source` (the
page path, sent to Buttondown as a tag) and a honeypot `website`. Without JS
the function answers `303 /subscribed/` (a static noindex page the builder
writes); the form's inline script (< 1 KB) sends `Accept: application/json`
and shows `{"ok", "message"}` in place. Errors: 400 invalid email, 429 per-IP
rate limit (5 per 10 min per warm instance, best-effort), 502 when Buttondown
fails or the key is missing, 405 for any method but POST. A honeypot hit or
an already-subscribed email gets the normal success answer and (honeypot) no
upstream call.

Two implementations serve this one contract (ADR-0004):

| Target | Code | Verified by |
|--------|------|-------------|
| QA (Vercel) and `make dev` | Go handler in `subscribe/`, deployed as `api/subscribe.go` | `make test` (`subscribe/subscribe_test.go`) |
| Production (Cloudflare) | Worker module `worker/index.js` | `scripts/check-worker.sh` (by hand; tests stay Go-only) |

A change to one is a change to both: the check script curls the same cases
the Go tests assert.

### Go function (QA Target)

Environment (Vercel project settings, Production and Preview):

| Variable | Required | Notes |
|----------|----------|-------|
| `BUTTONDOWN_API_KEY` | yes, secret | Buttondown API key. Use a separate test newsletter's key for Preview. Never sent to the browser. |
| `BUTTONDOWN_BASE_URL` | no | Default `https://api.buttondown.com`. Point at a fake for local runs. |

The Buttondown plan must include tags (otherwise Buttondown answers 403
`feature_disabled` and every signup is a 502).

Local run (no Vercel CLI needed): `make dev` (see "Working on the site"),
or by hand `make build` then `BUTTONDOWN_API_KEY=... go run ./cmd/subscribe`
(`-out` picks the Build Output dir; `-out ""` for the function only).
`vercel dev` is not used: `cmd/subscribe` runs the same handler behind the
same routes without a Vercel login.

Deploying: Vercel's zero-config detection turns `api/subscribe.go` into the
function `api/subscribe` (Go runtime, built from the root `go.mod`).
`vercel build` runs our build command first (`make build`, which writes
`.vercel/output`), keeps that tree, then `@vercel/go` adds
`functions/api/subscribe.func` and appends its routes to `config.json`
after ours. Checked with a local `vercel build` (CLI 65.0.0, a hand-written
`.vercel/project.json`, no login): the merged `config.json` keeps every
builder route in order and adds `{"src":"^/api(/.*)?$","status":404}`
after `{"handle":"filesystem"}` plus `error`/`miss` phases. That check
found one trap: `@vercel/go` renames the root module to `handler` while
compiling, so `api/subscribe.go` cannot import anything under
`internal/` ("use of internal package not allowed"); the function's
package therefore lives at `subscribe/`. The builder's routes keep `/api/*`
out of the trailing-slash redirects (a pass-through route before them) and
put the 404 catch-all after `{"handle":"filesystem"}`, so a POST is never
redirected to `/api/subscribe/` or answered by the 404. A POST on a real
deployment is still to be seen (setup wizard, stage 10).

### Worker (Production Target)

`worker/index.js` exports `fetch`: `/api/subscribe` runs the same contract as
the Go handler; every other request goes to the static assets binding
(`env.ASSETS`). The builder's Workers config sets `main` to this file and
`run_worker_first` to the signup path only, so pages and assets are served
without running any code. Client IP for the rate limit and Buttondown's
`ip_address` is `CF-Connecting-IP` (set by Cloudflare's edge). The rate
limit counters live in the isolate's memory: each Cloudflare location and
each isolate keeps its own, and an eviction resets them, so the limit is
best-effort abuse damping, as on Vercel. Email syntax matches the Go check
for everything a Reader types (it rejects quoted local parts, which Go's
`net/mail` accepts).

Worker secret and variable (same names as the Go function):

| Name | Kind | Notes |
|------|------|-------|
| `BUTTONDOWN_API_KEY` | secret, required | `wrangler secret put BUTTONDOWN_API_KEY` (or the dashboard). Sent upstream only; never in a response or a log line. Missing: every signup is a 502. |
| `BUTTONDOWN_BASE_URL` | var, optional | Default `https://api.buttondown.com`. The check script points it at a stub. |

Check it locally (no Cloudflare login, no network):

```sh
npm ci                         # pinned wrangler (package.json)
scripts/check-worker.sh        # assets from .vercel/output/static; builds if missing
scripts/check-worker.sh DIR    # or any built static tree
```

The script starts a stub Buttondown (logs each request's email,
`Authorization: Token` header and tags), runs the Worker under `wrangler dev`
with a throwaway Workers config (same assets settings as the builder's), and
prints PASS/FAIL for: form 303 to `/subscribed/` with the upstream request
checked; JSON 200 `{"ok":true}`; invalid emails 400 (JSON and form, link back
to the source page); honeypot success with no upstream call; rate limit after
5 from one IP, another IP unaffected; upstream 5xx and unreachable upstream
502 without the key or upstream body; already-subscribed success;
Buttondown-rejected email 400; GET 405 with `Allow: POST`; missing key 502
with no upstream call; non-signup paths served from assets (200, trailing-slash
307, 404 page); the key absent from every response. Exit status is the number
of failed cases. `STUB_PORT`/`WORKER_PORT` (default 18790/18791, plus 18792
for the missing-key instance) move the ports.

## Feeds, sitemaps, robots and IndexNow

All written in the write phase, after the link gate, from what the build
actually wrote (a page the gate removed is in no feed or sitemap).

- **Feeds** (RSS 2.0, full text in `content:encoded`, root-relative links
  made absolute): `/feed.xml` (Explainers and Briefs), `/{vertical}/feed.xml`,
  `/briefs/feed.xml`; newest 50 each, channel description = tagline.
  Feed readers load no site CSS, so diagrams in feeds have their CSS
  variables resolved to the light palette (`views.FeedSafeDiagrams`);
  pages keep the themable variables.
- **`/sitemap.xml`**: the canonical URL of every written page that is not
  noindex, once. `lastmod`: an Explainer's Review date (or publication, if
  later), a Brief's publication, else the newest `<time datetime>` on the
  page (listings, Author pages), else when the page's `<main>` last changed
  between builds (static pages).
- **`/news-sitemap.xml`**: Briefs published in the 48 h up to the build
  clock (`Config.Now`), newest first, at most 1,000, with `news:name`
  "HastCode", `news:language` `en`, `news:publication_date` (RFC 3339) and
  `news:title` (Google's news sitemap doc, last updated 2025-12-10).
  Explainers never appear. Vercel only rebuilds on deploy, so a Brief drops
  out of the file at the first build after it is 48 h old, not on the dot.
- **`/robots.txt`**: allow all, explicit allow groups for `Googlebot`,
  `Bingbot`, `OAI-SearchBot`, `PerplexityBot`, both sitemap lines. Training
  crawlers (`GPTBot`, `Google-Extended`) fall under `*` (allowed).
- **IndexNow**: `indexnow_key` in `site.yaml` (public by design; keep it
  stable) is served as `/<key>.txt`. The build diffs every indexable page's
  `<main>` hash against the previous successful build's manifest
  (`<cache>/pages/<mode>.json`, default `node_modules/.cache/hastcode/`,
  which CI persists; never in `static/`) and, when anything is new,
  changed or gone, writes `.vercel/output/indexnow.json`: a ready IndexNow
  POST body `{host, key, keyLocation, urlList}` (≤ 10,000 URLs). It sits
  beside `static/`, so it is not served. No artefact means nothing to ping.
  The build itself makes no network call. A failed build does not update the
  manifest, so its changes are listed again next time; a build that passes
  but is never deployed (or whose ping fails) does, so those changes are
  not re-listed.

Post-deploy ping: the Production workflow moves the artefact out of
`.vercel/output` before `vercel deploy --prebuilt --prod` (so it is never
uploaded) and, once the deploy is live and the repository variable
`INDEXNOW_ENABLED` is `true`, runs the equivalent of

```sh
curl -fsS -X POST https://api.indexnow.org/indexnow \
  -H 'Content-Type: application/json; charset=utf-8' --data-binary @indexnow.json
```

200/202 are success; 403 means the key file is not live yet; 422 means a
URL is not on `host`. Only production builds of `site_url` ping, so
`INDEXNOW_ENABLED` stays `false` until hastcode.com points at Vercel (the
setup wizard flips it at cutover and clears the Actions cache, so the
first build after cutover lists every page). A failed ping fails the job
before the build cache is saved, so the next build lists those URLs again.

## Routes and caching

`builder/step_write_vercel_config.go` writes every route into
`.vercel/output/config.json` (Build Output API v3 `routes`; patterns stay
RE2-compatible, no lookaround, so the local router can apply the same file).
In order:

1. Cache headers (`continue`): HTML (paths ending in `/` or `.html`) gets
   `Cache-Control: public, max-age=0, must-revalidate` for browsers and
   `Vercel-CDN-Cache-Control: s-maxage=31536000, stale-while-revalidate=86400`
   for Vercel's CDN (stripped before the browser; the CDN cache is scoped to
   the deployment, so a new deploy is never served stale). `/assets/*` is
   `public, max-age=31536000, immutable`: put only content-hashed files
   there. Everything else (feeds, `/mark.svg`, `/pagefind/*`, whose names are
   stable) keeps Vercel's default `public, max-age=0, must-revalidate`.
2. `/video/*`, `/community/*`, `/group/*` -> 410 with the 404 page as body
   (ADR-0001); `/home` -> 301 `/`.
3. `/api/*` and `/.well-known/*` pass through untouched.
4. One URL per page: `/x` -> 308 `/x/`, `/x/index.html` -> 308 `/x/`,
   `/file.ext/` -> 308 `/file.ext` (paths whose last segment has a dot are
   files and keep no slash).
5. `{"handle":"filesystem"}` (static files, the subscribe function), then
   any miss -> `/404.html` with status 404.

Check locally: `make dev` (or `make build`, `go run ./cmd/subscribe`), then

```sh
for p in /video/abc /home /about /about/ /about/index.html /mark.svg/ /nope/; do
  curl -s -o /dev/null -w "%{http_code} $p %{redirect_url}\n" "http://localhost:3000$p"
done   # 410, 301 /, 308 /about/, 200, 308 /about/, 308 /mark.svg, 404
```

The local router (`vercel.NewRouter`) refuses a `config.json` using route
features it does not implement (other `handle` phases, `has`/`missing`, a
non-RE2 regex), so the check cannot silently diverge from the file. It is not
Vercel itself: the setup wizard (stage 10) repeats the check against the
first preview. Note that `vercel build`'s merged `config.json` (with
`@vercel/go`'s `error`/`miss` phases) is not servable by this router: run
`make dev`, which serves the builder's own output.

## Search

`/search/` (noindex) searches Explainers and Briefs in the browser with
[Pagefind](https://pagefind.app) `1.5.2` (pinned exactly in `package.json`;
`npm ci` installs the platform binary as an optional dependency). There is
no backend.

- **Index (post-step inside the builder).** The last pipeline step
  (`step_write_search_index.go`, write phase 800, after the link gate has
  removed broken pages; before the page manifest is saved) runs `node node_modules/pagefind/lib/runner/bin.cjs
  --site <out>/static --output-subdir pagefind --quiet` in `-tools`, then
  reads `static/pagefind/pagefind-entry.json`. A missing or empty index is a
  `search-index` error, so the build fails. `-skip-search-index` (Config
  `SkipSearchIndex`) skips the run: an error in production, a warning in
  preview. A build that already failed is not indexed.
- **What is indexed.** Explainer and Brief templates put
  `data-pagefind-body` on `<article class="h-entry">`; once any page has it,
  pages without it (listings, static pages, `/search/`, 404) are not
  indexed. Inside the article, ad slots, the newsletter box and the reading
  lists carry `data-pagefind-ignore="all"`; diagram figures, the Explainer's
  Sources section (`div.sources`) and the Brief's source line carry
  `data-pagefind-ignore`. Header, footer and breadcrumbs are outside the
  article. (A tree with no published Explainer or Brief has no
  `data-pagefind-body` at all, so Pagefind falls back to indexing every
  page's `<body>`; this ends with the first Explainer.)
- **Search page.** Plain HTML first: a `role=search` form with a labelled
  `input type=search`, a `<noscript>` note that search needs JavaScript,
  and links to every Vertical hub and the Briefs. The inline loader fetches
  `/pagefind/pagefind-modular-ui.js` on the first focus of the input, a
  submit, or a `?q=` in the URL; the modular UI then loads `pagefind.js` and
  the index chunks on demand. Results are styled in `site.css` with the
  colour tokens (light, dark, toggle); Pagefind's CSS is never loaded.
- **Caching.** `/pagefind/*` keeps Vercel's default revalidating cache: its
  entry files have stable names (see "Routes and caching").
- **Tests** (`builder/search_test.go`) read the index through
  `buildtest.SearchIndex(t, out)` (URL -> indexed text, decoded from the
  gzipped fragments). `buildtest.Build` skips a test when `node` or
  `node_modules/pagefind` is missing, unless the config skips indexing.

**Build command.** `npm ci && make build`: CI runs `npm ci` itself and
`vercel build` then runs `make build` (vercel.json; see "Deploying"). No
separate `npx pagefind` step: the builder runs Pagefind and fails the build
if the index is missing. Preview deployments build with `MODE=preview`,
which still indexes; `FLAGS=-skip-search-index` (preview mode only) is for
quick local builds.

## Mermaid diagrams

A ```` ```mermaid ```` fence in an Explainer becomes `<figure class="diagram">`
with inline SVG and a `<figcaption>` from the fence's first line
`%% title: <caption>` (no title line, no caption). Step:
`builder/step_transform_diagrams.go`; renderer: `internal/mermaid`.

- **Renderer:** [`beautiful-mermaid`](https://github.com/lukilabs/beautiful-mermaid)
  `1.1.3` (pure JS + elkjs, no browser), pinned exactly in `package.json` /
  `package-lock.json`. The builder runs `node --input-type=module --eval
  <render.mjs>` once per build in `-tools` with every uncached fence on stdin.
  Chosen over `mmdc` (needs Chrome + `nss` on Vercel, hex-only theming) and
  Rust renderers (need glibc 2.35; Vercel's AL2023 has 2.34).
- **Supported types:** `flowchart`/`graph`, `sequenceDiagram`,
  `stateDiagram(-v2)`, `classDiagram`, `erDiagram`. Anything else (`gantt`,
  `pie`, `xychart-beta`, …), an empty fence, a flowchart line with unbalanced
  brackets or an edge with no target, or SVG with no nodes is a build error
  (`diagram-render` rule) attributed to the Explainer and naming its URL; that page
  is not written. The renderer itself is lenient, so these checks run in Go.
- **Output:** colours are `currentColor` and CSS variables only; the renderer's
  `<style>` (Google Fonts `@import`, page-global rules) is stripped and its
  rules live once in `internal/views/site.css` (the `.diagram` block), in the
  hashed site stylesheet (ADR-0003). Each SVG maps its `--bg`/`--fg`/`--accent` to
  `--diagram-bg`/`currentColor`/`--diagram-accent`, which `.diagram` sets from
  the palette tokens, so light, dark and the theme toggle all apply. Feeds
  carry a copy with those variables resolved to light-palette colours,
  computed from the same `.diagram` rules (`internal/views/feed_diagrams.go`). Element ids are prefixed per diagram (`d1-`, `d2-`…). Labels stay
  `<text>`. `width`/`height`/`viewBox` are kept (no layout shift).
- **Cache:** rendered SVG is stored at `<cache>/mermaid/<sha256>.svg`, keyed on
  the fence source, the installed renderer version and a post-processing
  version (`cacheVersion` in `internal/mermaid`; bump it when changing render
  options or post-processing). The default cache dir is
  `node_modules/.cache/hastcode/`, which CI persists between deploys with
  `actions/cache` (restored after `npm ci`, which deletes `node_modules/`).
  Fully cached builds never start Node.
- **CI build:** `npm ci` runs before `make build` from the repo root, with
  Node 24 (Node 22+ works).

## Deploying

Production is `git push` to `main` (ADR-0002); every other branch gets a
preview on each push. Both run in GitHub Actions with the Vercel CLI, not in Vercel's
Git integration, which `vercel.json` turns off (`git.deploymentEnabled:
false`):

- `vercel.json`: `framework: null` ("Other"), `installCommand: ""` (CI
  installs; an empty string skips Vercel's install step), `buildCommand:
  "make build"`. No `outputDirectory`: the build writes Build Output API v3
  to `.vercel/output` directly, which `vercel build` detects and keeps.
  `api/subscribe.go` is picked up by zero-config detection and compiled by
  `@vercel/go` with the Go version from `go.mod`.
- **Why CI and not Vercel's own builds.** (1) The build command needs Go on
  `PATH`; Vercel's build image lists Node, Python and Ruby for builds and
  provides Go only to the `@vercel/go` function builder. (2) The IndexNow
  artefact exists only on the build machine and must be POSTed after the
  deploy is live (ticket 09). (3) The build cache (heroes, diagrams, page
  manifest) persists in `actions/cache` under our control.
- `.github/workflows/production.yml` (push to `main`, or run by hand):
  `make test`, then `vercel pull` -> `vercel build --prod` -> `vercel deploy
  --prebuilt --prod` -> IndexNow ping (when `INDEXNOW_ENABLED`) -> save the
  build cache. One run at a time, in order.
- `.github/workflows/preview.yml` (push to any branch but `main`; a
  same-repo pull request is covered by its pushes, so it is not built
  twice; a fork's pull request only runs `make test`): `make test`; `vercel build` with `MODE=preview` -> `vercel deploy --prebuilt`
  (URL in the run summary); then **Lighthouse CI** (`lighthouserc.json`)
  against home, `/search/`, `/briefs/`, a Vertical hub, an Explainer and a
  Brief (`scripts/lighthouse-urls.sh` picks them from the build's sitemap).
  Budgets, mobile emulation with simulated throttling: LCP ≤ 1500 ms
  (median of 3 runs), CLS 0 (worst run), first-party script transfer
  (compressed) ≤ 60 KB as a secondary guard; the spec's total-JS budget
  (uncompressed, analytics included) is enforced by `make test`. A miss
  fails the check; reports are the `lighthouse-reports` artefact (never
  public storage: previews hold drafts). The preview is behind Vercel
  Authentication, so Lighthouse sends `x-vercel-protection-bypass`; it
  blocks Plausible and ad hosts so no third party sees that header and runs
  are not counted as pageviews (third-party JS sizes are covered by
  `TestJSBudget_EveryPageStaysWithin60KB` instead).
- Shared setup: `.github/actions/setup` (Go from `go.mod`, Node 24, `npm
  ci`, `vercel@65.0.0`).

**First-time setup** is a script that walks you through every account
step and writes the values where they belong:

```sh
scripts/setup-deploy.sh
```

Stages: tools and logins; GitHub repo (private, remote `origin`); Vercel
team on Pro; Vercel project + `vercel link` + Node 24.x; Buttondown
production and test newsletter keys; Vercel env vars
(`BUTTONDOWN_API_KEY`, Sensitive, Production and Preview); Protection
Bypass for Automation secret; Vercel token and GitHub secrets
(`VERCEL_TOKEN`, `VERCEL_ORG_ID`, `VERCEL_PROJECT_ID`,
`VERCEL_AUTOMATION_BYPASS_SECRET`, variable `INDEXNOW_ENABLED=false`);
first push and a draft PR; an automated check of the first preview
(every route, cache header and subscribe POST that could not be checked
without a real deploy); then the **domain cutover**, behind a
confirmation that ADR-0001's threshold is met (≥ 30 Explainers across all
six Verticals and ≥ 20 Briefs): `vercel domains add`, Cloudflare records
DNS only, `www` -> apex 301 at Vercel, live checks, `INDEXNOW_ENABLED=true`.
It is safe to re-run; values are kept in `.env` (gitignored).

Local equivalent of a CI build, for debugging: `npm ci`, `vercel pull`,
`vercel build` (needs the project linked). For page checks use `make dev`.

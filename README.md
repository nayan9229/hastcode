# hastcode

Static-first Go publication for hastcode.com ("Hast Code"). The Go program is a
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
npm ci          # Node build tools from package.json (Mermaid renderer)
make generate   # go generate ./...  (runs `go tool templ generate`)
make test       # generate + go test ./...
make vet        # generate + go vet ./...
make build      # generate + go run ./cmd/builder (content/ -> .vercel/output/)
```

The Makefile exports `CGO_ENABLED=0` and `GOFLAGS=-tags=nodynamic`: the AVIF
and WebP encoders (gen2brain/avif, gen2brain/webp) then run as WASM in-process
and never load a system libavif/libwebp, so image output is byte-identical
everywhere. Use the same settings when calling `go` directly.

Builder flags: `-content` (default `content`), `-config` (default `site.yaml`),
`-out` (default `.vercel/output`), `-mode production|preview`, `-tools` (dir
with `package.json`/`node_modules`, default `.`), `-cache` (build cache, default
`<tools>/node_modules/.cache/hastcode`), `-image-cache` (overrides
`images.cache_dir`). It prints every
error and warning and exits non-zero if the build failed.

Content layout (spec "Content model"): `content/explainers/<vertical>/<slug>.md`,
`content/authors/<slug>.md`; Briefs and static pages come in later tickets.
Add a Vertical by adding one entry to `site.yaml`.

Heroes: the Explainer names a raster beside its Markdown in `hero` (and
`hero_alt`). The build fails if it is missing, does not decode, or is narrower
than 1200 px. It writes AVIF, WebP and JPEG at 480/800/1200/1600 px wide (never
upscaled) plus 1200 px JPEG crops at 16:9, 4:3 and 1:1, to
`/assets/heroes/<name>-<w>x<h>.<hash>.<ext>` (hash of the file bytes, so
immutable). Encoded variants are cached in `images.cache_dir` (default
`node_modules/.cache/hastcode/images`, which Vercel persists between builds),
keyed on the source bytes and encoder settings; an unchanged hero is never
re-encoded. `images.fast: true` trades file size for encode speed (tests).

## Codebase map

| Path | What lives there |
|------|------------------|
| `cmd/builder/` | CLI only: flags, `LoadConfig`, print `Result`, exit code. No logic. |
| `builder/` | The single deep module. Public API: `Config`, `Vertical`, `Mode`, `LoadConfig`, `Build(cfg, contentDir, outDir) Result`, `Result`/`Problem`. Everything else unexported. |
| `builder/pipeline.go` | Phases, `registerStep`, the `site` state struct, `fail`/`warn`, `publishedExplainers()`. |
| `builder/step_*.go` | One file per pipeline step (load, validate, transform, render, write). |
| `builder/urls.go` | The URL scheme: one func per page type (`explainerURL`, …). |
| `builder/render.go` | `writePage(urlPath, templComponent)` incl. the 50 KB HTML size gate (rule `html-size`); output reset step. |
| `builder/*_test.go` | Build-seam tests (package `builder_test`), one file per feature area. |
| `internal/content/` | Parse content files (frontmatter + body) into `Explainer`, `Author` (later `Brief`, `Page`). No config or output knowledge. |
| `internal/markdown/` | goldmark `Parse` -> `*Document` (AST + source) and `Render(doc)`. Steps inspect/rewrite the AST between the two. `doc.Diagrams()` finds ```` ```mermaid ```` fences (title from a leading `%% title:` line); `doc.ReplaceDiagram(d, html)` swaps one for builder HTML. |
| `internal/mermaid/` | Mermaid -> SVG: one batch Node process (`render.mjs`, embedded) per build, type/syntax checks, post-processing, on-disk cache by source hash, `PrefixIDs`. |
| `internal/views/` | templ components + view-model structs (`models.go`). |
| `internal/images/` | Resize/crop/encode (AVIF, WebP, JPEG) through a worker pool with an on-disk cache. Knows nothing of heroes or URLs. |
| `internal/views/brand.go` | Brand constants (`PublicationName`, `Tagline`, `MarkPath`), the light/dark palettes, the generated SVG mark, `StaticAssets()`. |
| `internal/views/site.css`, `stylesheet.go`, `theme.js` | The one stylesheet (hashed file, ADR-0003) and the inline theme-toggle script. |
| `builder/chrome.go` | `s.page(title, description)`: every page's `views.Page` incl. header/footer nav, ads flag and Plausible domain. |
| `internal/vercel/` | Build Output API writer: `static/` files and `config.json`. |
| `internal/buildtest/` | Test support for the Build seam (fixtures, run build, parse output). |
| `internal/subscribe/` | The subscribe function: `New(Config) http.Handler`, `ConfigFromEnv()`, form contract constants (`Path`, `SubscribedPath`, `HoneypotField`). Buttondown client in `buttondown.go`, per-IP limiter in `ratelimit.go`. Subscribe-seam tests in `subscribe_test.go`. |
| `api/subscribe.go` | Vercel Go function entry (package `handler`, `func Handler`). Wiring only. No other `.go` files here: each becomes a function. |
| `cmd/subscribe/` | Local run of the subscribe function plus the built site. |
| `site.yaml` | Site config: `site_url`, `verticals`, `plausible_domain` (blank = no analytics), `ads_enabled`, `adsense_client` (required when ads are on). |
| `package.json` | Pinned Node build tools (`beautiful-mermaid`). Not a JS app. |

Planned homes for parallel tickets: editorial gates -> `builder/step_validate_*.go`;
hero images -> `builder/step_transform_heroes.go` +
`internal/images/`; Mermaid -> `builder/step_transform_diagrams.go` +
`internal/mermaid/`; design system -> `internal/views/layout_*.templ` + CSS.
Briefs: end the Brief template with `@NewsletterForm(views.Newsletter{Source: briefURL(b)})`
like `article_explainer.templ` does.

## Conventions

**Adding a pipeline step.** Create `builder/step_<phase>_<what>.go` with an
`init()` that calls `registerStep(phase, order, name, fn)`. Do not edit a
central list. Phases run in order: `phaseLoad` -> `phaseValidate` ->
`phaseTransform` -> `phaseRender` -> `phaseWrite`; within a phase, by `order`
(use gaps: 10, 20, 30…). Existing orders: load authors 10, load explainers 20,
exclude drafts (production only) 30; validate ads config 5, references 10,
unique slugs 20, Explainer skeleton 30, section-length warning 40,
Review-overdue warning 50, validate heroes 50; transform render diagrams 10,
encode heroes 50; render reset 0, write site assets 10, write hero assets 15,
render explainers 20, render subscribed page 60; write config.json 100.
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
`s.page(title, description)` (never a bare `views.Page{}`) so the header and
footer nav are filled, and wrap the component in `@Layout(page)`.

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
scripts are the theme toggle (≤ 1 KB, in `<head>`) and the newsletter form's
enhancement (≤ 1 KB, inside `.newsletter`); analytics and Pagefind are the
only other scripts allowed (`TestThemeToggle_IsTheOnlyScriptAndUnder1KB`,
`TestNewsletter_EnhancementScriptIsSmall`), plus the AdSense library and its
per-unit push only when `ads_enabled`. Total client JS per page stays
≤ 60 KB (`TestJSBudget_EveryPageStaysWithin60KB`; give any new third-party
script a size there).

**Weight, ads, analytics.** Every page written through `s.writePage` is
checked against the 50 KB uncompressed HTML gate (the stylesheet is external
and not counted, ADR-0003); an oversized page is a site-wide `html-size`
error naming the URL and is not written. Ads come from one component,
`@AdSlot(page.Ads, views.AdAfterTLDR|views.AdAfterBody)`
(`partial_ad_slot.templ`): an empty `aria-hidden` placeholder with ads off,
an AdSense unit with ads on, in a box whose height is fixed by
`--ad-slot-height` so enabling ads causes no layout shift; it carries
`data-pagefind-ignore="all"`. Article layouts render the body in two halves
(`Document.SplitAfterSection("TL;DR")` + `markdown.RenderBlocks`) with the
first slot between them and the second after `.e-content`. Plausible comes
from `layout_analytics.templ` in every `<head>` when `plausible_domain` is set.
Everything must work without JS.

**Tests.** Only two seams are tested (spec "Testing Decisions"): the Build seam
and the Subscribe HTTP seam. Build-seam tests live in `builder/` as package
`builder_test`, import only `builder` and `internal/buildtest`, and assert on
parsed output. Helpers in `internal/buildtest`:

- Fixtures: `ValidTree()`, `NewTree()` (default Author only), `Tree.WithExplainer(e)`,
  `Tree.WithAuthor(a)`, `Tree.With(path, bytes)`; `ValidExplainer()` (meets the
  full editorial skeleton: TL;DR, Mermaid fence, Sources, 1600x900 PNG hero),
  `ValidComparison()` (adds a `## Trade-offs` table), `DefaultAuthor()`,
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
- Assert: `Page(t, outDir, "/vertical/slug/")` (goquery document),
  `Exists(outDir, urlPath)`, `StaticPath(outDir, urlPath)`, `JSONLD(t, doc)`,
  `JSONLDOfType(t, doc, "Article")`, `VercelConfig(t, outDir)`,
  `DecodeImage(t, outDir, urlPath)` / `Images(t, outDir)` (decode output images
  and measure them), `Srcset(t, attr)`,
  `RequireError(t, res, file, substring)`, `RequireWarning(...)`,
  `RequireErrorRule(t, res, file, rule, substring)`, `RequireWarningRule(...)`,
  `RequireNoProblemRule(t, res, file, rule)`, `FilesContaining(t, outDir, s)`
  (prove a piece appears in no output file).

Subscribe-seam tests live in `internal/subscribe/subscribe_test.go` (package
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
fails or the key is missing. A honeypot hit or an already-subscribed email
gets the normal success answer and (honeypot) no upstream call.

Environment (Vercel project settings, Production and Preview):

| Variable | Required | Notes |
|----------|----------|-------|
| `BUTTONDOWN_API_KEY` | yes, secret | Buttondown API key. Use a separate test newsletter's key for Preview. Never sent to the browser. |
| `BUTTONDOWN_BASE_URL` | no | Default `https://api.buttondown.com`. Point at a fake for local runs. |

The Buttondown plan must include tags (otherwise Buttondown answers 403
`feature_disabled` and every signup is a 502).

Local run (no Vercel CLI needed): `make build`, then
`BUTTONDOWN_API_KEY=... go run ./cmd/subscribe` serves the built site and the
function on http://localhost:3000 (`-static ""` for the function only).
`vercel dev` should also serve `api/subscribe.go` through the Go runtime, but
that has not been tried yet (ticket 15).

Deploying: Vercel's zero-config detection turns `api/subscribe.go` into the
function `api/subscribe` (Go runtime, built from the root `go.mod`). On a
project whose build command writes `.vercel/output` (this one), `vc build`
runs the frontend builder first, keeps that Build Output tree, then lets
`@vercel/go` add `functions/api/subscribe.func` and merges its routes into
`config.json` (vercel/vercel `internals/builder-orchestration`
`sort-builders.ts`, `write-build-result.ts`). Not yet checked on a real
deploy: ticket 15's first preview verifies a working POST. Route rules the
builder adds to `config.json` (ticket 10) must not shadow `/api/subscribe`:
keep catch-alls (404, trailing-slash redirect) after `{"handle":"filesystem"}`
or exclude `/api/`, so a POST is never redirected to `/api/subscribe/`.

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
  rules live once in `internal/views/diagram.css`, appended to the hashed site
  stylesheet (ADR-0003). Each SVG maps its `--bg`/`--fg`/`--accent` to
  `--diagram-bg`/`currentColor`/`--diagram-accent`, which `.diagram` sets from
  the palette tokens, so light, dark and the theme toggle all apply. Element ids are prefixed per diagram (`d1-`, `d2-`…). Labels stay
  `<text>`. `width`/`height`/`viewBox` are kept (no layout shift).
- **Cache:** rendered SVG is stored at `<cache>/mermaid/<sha256>.svg`, keyed on
  the fence source, the installed renderer version and a post-processing
  version (`cacheVersion` in `internal/mermaid`; bump it when changing render
  options or post-processing). The default cache dir is under
  `node_modules/.cache/` because Vercel persists only `node_modules/**`
  between builds. Fully cached builds never start Node.
- **Vercel build:** the install step must run `npm ci` (Vercel does this by
  default when `package.json` is present) before the build command runs the
  Go builder from the repo root, e.g. `npm ci && make build` if the install
  step is overridden. Node 22+ is in the default build image.

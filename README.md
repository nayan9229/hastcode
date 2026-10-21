# hastcode

Static-first Go publication for hastcode.com ("Hast Code"). The Go program is a
**build tool**, not a server (ADR-0002): it reads Markdown from `content/`,
renders templ components, and writes a Vercel Build Output API v3 tree to
`.vercel/output/`. Domain terms (Explainer, Brief, Vertical, Author, Review,
Reader) are defined in `GLOSSARY.md`; use them in code and test names.

## Quick start

templ is a Go tool dependency (`go.mod` `tool` directive), so no global install
is needed. Generated `*_templ.go` files are **gitignored**: generate before
compiling, vetting or testing anything.

```sh
make generate   # go generate ./...  (runs `go tool templ generate`)
make test       # generate + go test ./...
make vet        # generate + go vet ./...
make build      # generate + go run ./cmd/builder (content/ -> .vercel/output/)
```

Builder flags: `-content` (default `content`), `-config` (default `site.yaml`),
`-out` (default `.vercel/output`), `-mode production|preview`. It prints every
error and warning and exits non-zero if the build failed.

Content layout (spec "Content model"): `content/explainers/<vertical>/<slug>.md`,
`content/authors/<slug>.md`; Briefs and static pages come in later tickets.
Add a Vertical by adding one entry to `site.yaml`.

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
| `internal/markdown/` | goldmark `Parse` -> `*Document` (AST + source) and `Render(doc)`. Steps inspect/rewrite the AST between the two. |
| `internal/views/` | templ components + view-model structs (`models.go`). |
| `internal/views/brand.go` | Brand constants (`PublicationName`, `Tagline`, `MarkPath`), the light/dark palettes, the generated SVG mark, `StaticAssets()`. |
| `internal/views/site.css`, `stylesheet.go`, `theme.js` | The one stylesheet (hashed file, ADR-0003) and the inline theme-toggle script. |
| `builder/chrome.go` | `s.page(title, description)`: every page's `views.Page` incl. header/footer nav, ads flag and Plausible domain. |
| `internal/vercel/` | Build Output API writer: `static/` files and `config.json`. |
| `internal/buildtest/` | Test support for the Build seam (fixtures, run build, parse output). |
| `site.yaml` | Site config: `site_url`, `verticals`, `plausible_domain` (blank = no analytics), `ads_enabled`, `adsense_client` (required when ads are on). |

Planned homes for parallel tickets: editorial gates -> `builder/step_validate_*.go`;
Mermaid -> `builder/step_transform_diagrams.go` (+ helper package under
`internal/`); hero images -> `builder/step_transform_heroes.go` +
`internal/images/`; design system -> `internal/views/layout_*.templ` + CSS;
newsletter -> handler in `internal/subscribe/`, Vercel entry in `api/`, form in
`internal/views/partial_newsletter.templ`.

## Conventions

**Adding a pipeline step.** Create `builder/step_<phase>_<what>.go` with an
`init()` that calls `registerStep(phase, order, name, fn)`. Do not edit a
central list. Phases run in order: `phaseLoad` -> `phaseValidate` ->
`phaseTransform` -> `phaseRender` -> `phaseWrite`; within a phase, by `order`
(use gaps: 10, 20, 30…). Existing orders: load authors 10, load explainers 20,
exclude drafts (production only) 30; validate ads config 5, references 10, unique slugs 20,
Explainer skeleton 30, section-length warning 40, Review-overdue warning 50;
render reset 0, write site assets 10, render explainers 20; write config.json 100.
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
a comment naming the partial, one declaration block per line. The only
inline script is the theme toggle (≤ 1 KB); analytics and Pagefind are the
only other scripts allowed (`TestThemeToggle_IsTheOnlyScriptAndUnder1KB`),
plus the AdSense library and its per-unit push only when `ads_enabled`.
Total client JS per page stays ≤ 60 KB (`TestJSBudget_EveryPageStaysWithin60KB`;
give any new third-party script a size there).

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
  `HeroPNG(w, h)`, `Config()` (launch Verticals, production; set `cfg.Now` to
  pin the clock).
  Mutate a copy of `ValidExplainer()` to make an invalid fixture.
- Run: `Build(t, cfg, tree) (outDir, Result)`, `MustBuild(t, cfg, tree) outDir`.
- Assert: `Page(t, outDir, "/vertical/slug/")` (goquery document),
  `Exists(outDir, urlPath)`, `StaticPath(outDir, urlPath)`, `JSONLD(t, doc)`,
  `JSONLDOfType(t, doc, "Article")`, `VercelConfig(t, outDir)`,
  `RequireError(t, res, file, substring)`, `RequireWarning(...)`,
  `RequireErrorRule(t, res, file, rule, substring)`, `RequireWarningRule(...)`,
  `RequireNoProblemRule(t, res, file, rule)`, `FilesContaining(t, outDir, s)`
  (prove a piece appears in no output file).

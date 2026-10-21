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

Builder flags: `-content` (default `content`), `-config` (default `site.yaml`),
`-out` (default `.vercel/output`), `-mode production|preview`, `-tools` (dir
with `package.json`/`node_modules`, default `.`), `-cache` (build cache, default
`<tools>/node_modules/.cache/hastcode`). It prints every
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
| `builder/render.go` | `writePage(urlPath, templComponent)`; output reset step. |
| `builder/*_test.go` | Build-seam tests (package `builder_test`), one file per feature area. |
| `internal/content/` | Parse content files (frontmatter + body) into `Explainer`, `Author` (later `Brief`, `Page`). No config or output knowledge. |
| `internal/markdown/` | goldmark `Parse` -> `*Document` (AST + source) and `Render(doc)`. Steps inspect/rewrite the AST between the two. `doc.Diagrams()` finds ```` ```mermaid ```` fences (title from a leading `%% title:` line); `doc.ReplaceDiagram(d, html)` swaps one for builder HTML. |
| `internal/mermaid/` | Mermaid -> SVG: one batch Node process (`render.mjs`, embedded) per build, type/syntax checks, post-processing, on-disk cache by source hash, `PrefixIDs`. |
| `internal/views/` | templ components + view-model structs (`models.go`). |
| `internal/vercel/` | Build Output API writer: `static/` files and `config.json`. |
| `internal/buildtest/` | Test support for the Build seam (fixtures, run build, parse output). |
| `site.yaml` | Site config: `site_url`, `verticals`. |
| `package.json` | Pinned Node build tools (`beautiful-mermaid`). Not a JS app. |

Planned homes for parallel tickets: editorial gates -> `builder/step_validate_*.go`;
hero images -> `builder/step_transform_heroes.go` +
`internal/images/`; design system -> `internal/views/layout_*.templ` + CSS;
newsletter -> handler in `internal/subscribe/`, Vercel entry in `api/`, form in
`internal/views/partial_newsletter.templ`.

## Conventions

**Adding a pipeline step.** Create `builder/step_<phase>_<what>.go` with an
`init()` that calls `registerStep(phase, order, name, fn)`. Do not edit a
central list. Phases run in order: `phaseLoad` -> `phaseValidate` ->
`phaseTransform` -> `phaseRender` -> `phaseWrite`; within a phase, by `order`
(use gaps: 10, 20, 30…). Existing orders: load authors 10, load explainers 20;
validate references 10; transform render diagrams 10; render reset 0, render
explainers 20; write config.json 100.
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
`s.writePage(url, component)`.

**Tests.** Only two seams are tested (spec "Testing Decisions"): the Build seam
and the Subscribe HTTP seam. Build-seam tests live in `builder/` as package
`builder_test`, import only `builder` and `internal/buildtest`, and assert on
parsed output. Helpers in `internal/buildtest`:

- Fixtures: `ValidTree()`, `NewTree()` (default Author only), `Tree.WithExplainer(e)`,
  `Tree.WithAuthor(a)`, `Tree.With(path, bytes)`; `ValidExplainer()` (meets the
  full editorial skeleton: TL;DR, Mermaid fence, Sources, 1600x900 PNG hero),
  `DefaultAuthor()`, `HeroPNG(w, h)`, `Config()` (launch Verticals, production).
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
  `RequireError(t, res, file, substring)`, `RequireWarning(...)`.

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
  (`diagram` rule) attributed to the Explainer and naming its URL; that page
  is not written. The renderer itself is lenient, so these checks run in Go.
- **Output:** colours are `currentColor` and CSS variables only; the renderer's
  `<style>` (Google Fonts `@import`, page-global rules) is stripped and its
  rules live once in `internal/views/diagram.css` (inlined by the
  `DiagramStyles` partial in `Head`; the site stylesheet may absorb it). Theme
  hooks: `--diagram-bg`, `--diagram-accent` and `color` on `.diagram`, falling
  back to `--color-bg` / `--color-accent` / `--color-text`, then the brand
  palette. Element ids are prefixed per diagram (`d1-`, `d2-`…). Labels stay
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

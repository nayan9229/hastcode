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
| `builder/render.go` | `writePage(urlPath, templComponent)`; output reset step. |
| `builder/*_test.go` | Build-seam tests (package `builder_test`), one file per feature area. |
| `internal/content/` | Parse content files (frontmatter + body) into `Explainer`, `Author` (later `Brief`, `Page`). No config or output knowledge. |
| `internal/markdown/` | goldmark `Parse` -> `*Document` (AST + source) and `Render(doc)`. Steps inspect/rewrite the AST between the two. |
| `internal/views/` | templ components + view-model structs (`models.go`). |
| `internal/vercel/` | Build Output API writer: `static/` files and `config.json`. |
| `internal/buildtest/` | Test support for the Build seam (fixtures, run build, parse output). |
| `internal/subscribe/` | The subscribe function: `New(Config) http.Handler`, `ConfigFromEnv()`, form contract constants (`Path`, `SubscribedPath`, `HoneypotField`). Buttondown client in `buttondown.go`, per-IP limiter in `ratelimit.go`. Subscribe-seam tests in `subscribe_test.go`. |
| `api/subscribe.go` | Vercel Go function entry (package `handler`, `func Handler`). Wiring only. No other `.go` files here: each becomes a function. |
| `cmd/subscribe/` | Local run of the subscribe function plus the built site. |
| `site.yaml` | Site config: `site_url`, `verticals`. |

Planned homes for parallel tickets: editorial gates -> `builder/step_validate_*.go`;
Mermaid -> `builder/step_transform_diagrams.go` (+ helper package under
`internal/`); hero images -> `builder/step_transform_heroes.go` +
`internal/images/`; design system -> `internal/views/layout_*.templ` + CSS.
Briefs: end the Brief template with `@NewsletterForm(views.Newsletter{Source: briefURL(b)})`
like `article_explainer.templ` does.

## Conventions

**Adding a pipeline step.** Create `builder/step_<phase>_<what>.go` with an
`init()` that calls `registerStep(phase, order, name, fn)`. Do not edit a
central list. Phases run in order: `phaseLoad` -> `phaseValidate` ->
`phaseTransform` -> `phaseRender` -> `phaseWrite`; within a phase, by `order`
(use gaps: 10, 20, 30…). Existing orders: load authors 10, load explainers 20;
validate references 10; render reset 0, render explainers 20, render subscribed
page 60; write config.json 100.
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
- Assert: `Page(t, outDir, "/vertical/slug/")` (goquery document),
  `Exists(outDir, urlPath)`, `StaticPath(outDir, urlPath)`, `JSONLD(t, doc)`,
  `JSONLDOfType(t, doc, "Article")`, `VercelConfig(t, outDir)`,
  `RequireError(t, res, file, substring)`, `RequireWarning(...)`.

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

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
make lint-workflows  # pinned actionlint over .github/workflows/
```

## Working on the site

**Run the tests.** `npm ci` once, then `make test` (templ generate + `go
test ./...`, about a minute; `go test -short ./...` skips the slow
1,000-entry news-sitemap case). Build-seam tests skip, with the reason,
when Node or the npm packages are missing, so run `npm ci` first or most
of them will not run. `make vet` before committing. CI runs `make test`
on every push (to `main`: that is all it runs) and before every Publish.

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
(a pull request is optional), and the Preview workflow uploads a preview build
(drafts included) to a Cloudflare preview URL, posts it in the run
summary, then runs Lighthouse CI against it. Production builds never
contain drafts.

**Publish flow.** Remove `draft: true`, merge the pull request to `main`
(that only runs the tests). Then, on an up-to-date, clean `main`:
`scripts/publish.sh vl` deploys that commit to the QA Target
(https://qa.hastcode.com/); check it; `scripts/publish.sh cf` deploys the
same commit to Production (hastcode.com) and pings IndexNow with the pages
that changed. Feeds, sitemaps and the news sitemap update in the same
build. A build error (failed editorial gate, broken link, oversized page,
a page without JSON-LD, a leftover `TODO(author)`, missing R2
credentials) fails the workflow and nothing deploys. See "Deploying". The shipped
`content/` still has `TODO(author)` placeholders, so production builds
fail until the Author fills them (see "Content layout"); `make dev` and
previews work meanwhile.

**Add a Vertical.** Add one entry to `site.yaml` `verticals` (`slug`,
`name`, `description`, `order`), then put its Explainers in
`content/explainers/<slug>/`. The build gives it a hub at `/<slug>/`, a
header nav link, a home page entry and `/<slug>/feed.xml`. No code change.
The slug is a top-level URL forever: pick it once.

**Environment variables.** Nothing in the build reads secrets. The
subscribe function reads two (see "Subscribe function"):
`BUTTONDOWN_API_KEY` (required, secret) and `BUTTONDOWN_BASE_URL`
(optional, default `https://api.buttondown.com`); the Worker reads the same
names. Which newsletter each key points at: the Worker (Production) gets
the real newsletter's key at the DNS flip and the test newsletter's key
before it; the Vercel project (QA) and `.env` (`make dev`) get the test
newsletter's key.
An `images.store: r2` build reads the R2 credentials (see "Image Store").
The fetch command reads `GITHUB_TOKEN` (see "Repo Profiles: fetching a
Fact Pack").
CI's GitHub secrets and variables are listed under "Deploying". `make build`
itself also reads `MODE`, `FLAGS`, `OUT` from the environment (the QA
publish sets `MODE=production` for `vercel build`).

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
every Publish stays red until the Author page and the About,
Contact and Editorial Policy pages are filled in. Preview builds (`make
dev`, `make build MODE=preview`, the Preview workflow) only warn, so the
drafts can be read in place.
Add a Vertical by adding one entry to `site.yaml`.

Heroes: the Explainer names a raster beside its Markdown in `hero` (and
`hero_alt`). The build fails if it is missing, does not decode, or is narrower
than 1200 px. The one exception is a Repo Profile, which may omit `hero`:
the build then draws a 1600x900 card from its Fact Pack (wordmark,
repository, description, language, stars, license, latest release; light
palette, Go fonts; alt text "{repo}: {description}") and feeds it through
the same pipeline, so a changed star count or release renames every
variant and an unchanged Fact Pack hits the cache
(`step_load_hero_cards.go`, `internal/views/hero_card.go`). It writes AVIF and JPEG at 480/800/1000/1200/1472 px wide
(never upscaled past the source; `heroWidths`), plus 1200 px JPEG crops at
16:9, 4:3 and 1:1 (JSON-LD `image`) and 1200x630 (`og:image`,
`twitter:image`), to `/assets/heroes/<name>-<w>x<h>.<hash>.<ext>` (hash of
the file bytes, so immutable). No WebP: AVIF covers nearly every browser
and the JPEG covers the rest. Encoder settings
(`docs/research/image-optimisation-targets.md`): AVIF 4:4:4 chroma, speed
8, quality 60 below 1200 px and 55 at 1200 px and above (`heroQuality`;
the 1200 px file is the mobile LCP image); stdlib JPEG quality 82. When a
hero's 1200 px AVIF is over 150 KB (150,000 bytes, decimal: `heroWeightBudget`) the build warns (rule `hero-weight`,
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

Author photos follow the same format rule: PNG or JPEG, sniffed from the
bytes; a WebP photo fails with `author-photo`.

On the page the hero `<img>` carries `fetchpriority="high"` and no
`loading` attribute (it is the LCP image), and `<head>` keeps one AVIF
preload with the same `imagesrcset`/`imagesizes`.

## Codebase map

| Path | What lives there |
|------|------------------|
| `cmd/builder/` | CLI only: flags, `LoadConfig`, print `Result`, exit code. No logic. |
| `builder/` | The single deep module. Public API: `Config`, `Vertical`, `Mode` (`Production`, `Preview`, `ParseMode`), `ImageConfig` + `ImageStoreMode` (`StoreLocal`, `StoreR2`) + `DefaultImageCacheDir`, `LoadConfig`, `Build(cfg, contentDir, outDir) Result`, `Result`/`Problem`, and `Prune`/`PruneOptions`/`PruneReport` (the prune command). Everything else unexported. |
| `internal/buildflags/` | The flags `cmd/builder` and `cmd/prune` share (`-content`, `-config`, `-tools`, `-cache`, `-image-cache`) and `LoadConfig()`, which applies them. |
| `builder/pipeline.go` | Phases, `registerStep`, the `site` state struct, `fail`/`warn`, the "skip failed files" filter `passed` behind `publishedExplainers()`/`publishedBriefs()`/`publishedPages()`. |
| `builder/step_*.go` | One file per pipeline step (load, validate, transform, render, write). |
| `builder/urls.go` | The URL scheme: one func per page type (`explainerURL`, …), every fixed page, feed and sitemap path, the asset directories, and the retired paths (`retiredSections`, `homeAliasURL`). |
| `builder/render.go` | Output helpers: `writePage(urlPath, templComponent)` incl. the 50 KB HTML size gate (rule `html-size`), `writeXML`, `isWritten`, rule `markdown`. The output reset is `step_render_reset.go`. |
| `builder/load.go` | `eachMarkdown` and rule `frontmatter`, shared by the `step_load_*.go` steps. |
| `builder/*_test.go` | Build-seam tests (package `builder_test`), one file per feature area. |
| `internal/content/` | Parse content files (frontmatter + body) into `Explainer`, `Brief`, `Author`, `Page`, and a Repo Profile's `FactPack` (`fact_pack.go`: the Fact Pack v1 shape with its JSON names, `ParseFactPack`, and `FactPackFile(markdownFile)`, the naming rule: the same name with `.json`, beside the Markdown), and `Repo` (`repo.go`: owner/name, `ParseRepo`, the one check the frontmatter `repo` key and the fetch command's `-repo` share). No config or output knowledge. |
| `builder/repo_profile.go`, `step_load_fact_packs.go`, `step_validate_fact_packs.go`, `step_transform_fact_sections.go` | Repo Profiles (ADR-0005): each Profile's Fact Pack is loaded beside its Markdown (rule `fact-pack`: missing or unparsable), checked against its `repo` (`fact-pack-repo`) and the build clock (`fact-pack-stale`, a warning when it has been unchanged upstream for over 30 days: the repository may be dormant), and mapped to the header card (`repoProfileCard` -> `views.RepoProfileCard`, `partial_repo_profile_card.templ`), the Signals and Changelog sections inserted into the body at transform 35 (`repoProfileSignals`, `repoProfileChangelog` -> `views.RepoProfileSignals`, `views.RepoProfileChangelog`; `partial_repo_profile_signals.templ`, `partial_repo_profile_changelog.templ`; rule `fact-sections`) and the `SoftwareSourceCode` JSON-LD (`softwareSourceCode`). Tests: `builder/repo_profile_test.go`, `builder/repo_profile_facts_test.go`. |
| `builder/step_validate_repo_profile_skeleton.go` | The Repo Profile skeleton: the one place the nine written headings are spelled (`repoProfileHeadings`) and where the fact sections are placed (`repoProfileFactSections`: Signals after Architecture, Changelog after Alternatives). Rules `repo-profile-sections` (a missing, extra or reordered H2, naming the heading and its position), `repo-profile-fact-section` (a hand-written Signals or Changelog) and `repo-profile-faq` (five to eight `###` questions). See "Repo Profiles" under Conventions. Tests: `builder/repo_profile_skeleton_test.go`. |
| `internal/factpack/`, `cmd/fetch/` | The fetch command (ADR-0005): `Fetcher{GitHub, CacheDir, Now, Log}` with `Fetch(ctx, repo, markdownFile)` (ten REST calls at most, `assemble.go`; file flags from the trees, `files.go`; the Fact Pack written beside the Markdown in a stable layout, `write.go`) and `Refresh(ctx, contentDir)` (`refresh.go`; what changed, fetched-at aside, `diff.go`), whose `RefreshReport.ExitStatus()` is the command's exit status (`ExitChanged` 3 when a Fact Pack was rewritten, `ExitNothingChanged` 0). `cmd/fetch` is flags and exit codes only. Tests: `internal/factpack/*_test.go` against the in-process GitHub of `fake_test.go`. See "Repo Profiles: fetching a Fact Pack". |
| `internal/github/` | The GitHub REST client the fetch command uses: `NewClient(getenv)` (token from `GITHUB_TOKEN`, refused when unset), `Get(ctx, path, v)` (Bearer token, identifying User-Agent, `Retry-After` and primary rate-limit headers honoured with a bounded backoff, 202 and 5xx retried a bounded number of times, `RateLimitError` when the wait exceeds the bound), `LastPage(header)`. |
| `internal/markdown/` | goldmark `Parse` -> `*Document` (AST + source) and `Render(doc)`. Steps inspect/rewrite the AST between the two. `doc.Diagrams()` finds ```` ```mermaid ```` fences (title from a leading `%% title:` line); `doc.ReplaceDiagram(d, html)` swaps one for builder HTML; `doc.InsertHTMLAfterSection(heading, html)` places builder HTML at the end of an H2 section (a Repo Profile's fact sections). |
| `internal/mermaid/` | Mermaid -> SVG: one batch Node process (`render.mjs`, embedded) per build, type/syntax checks, post-processing, on-disk cache by source hash, `PrefixIDs`. |
| `internal/views/` | templ components + view-model structs (`models.go`). |
| `internal/images/` | Resize/crop/encode (AVIF, JPEG; per-variant quality) through a worker pool with an on-disk cache. Knows nothing of heroes or URLs. |
| `builder/step_load_hero_cards.go`, `internal/views/hero_card.go` | The generated hero card: a Repo Profile without `hero` gets a 1600x900 PNG drawn from its Fact Pack (`views.HeroCardPNG(views.HeroCard)`; `views.HeroCardText` is the renderer's text seam, see "Tests"), written under the build cache (`hero-cards/`) and registered as the Profile's hero source before validate heroes, so the hero pipeline encodes it like a sourced file (variants named `<slug>-card-<w>x<h>.<hash>.<ext>`). Tests: `builder/hero_card_test.go`, `internal/views/hero_card_test.go`. |
| `builder/image_store.go`, `step_render_upload_images.go` | The Image Store seam: `s.imageURL(sitePath)` (the URL pages use), `s.storeImage(sitePath, file, data, type)` (write in local mode, queue in r2 mode), `openImageStore` (config + credentials, checked before any step), and the upload step. See "Image Store". |
| `internal/imagestore/` | The R2 client (aws-sdk-go-v2): env var names, `FromEnv`, `New(creds, bucket)`, `Upload` (HEAD then PUT on a bounded pool), `List(prefix)` (ListObjectsV2 pages), `Delete(keys)` (DeleteObjects, 1,000 keys per request). Shared with the prune command. |
| `builder/prune.go`, `cmd/prune/` | The prune command: `Prune(ctx, cfg, contentDir, PruneOptions{Delete, Log}) (PruneReport, error)` builds the tree locally (Production + Preview) for `ImageKeys`, lists the bucket, reports and optionally deletes the rest; `make prune-images` runs it as a dry run. Tests: `builder/prune_test.go`. See "Image Store". |
| `internal/buildtest/image_store.go` | `NewImageStore(t)`: the fake Image Store (httptest S3: HEAD, PUT, ListObjectsV2, DeleteObjects). |
| `internal/views/brand.go`, `og_badge.go` | Brand constants (`PublicationName`, `Tagline`, `MarkPath`), the light/dark palettes, the generated SVG mark, the generated 1200x630 PNG OG badge (`OGBadge()`, Go fonts, `/assets/og-badge.<hash>.png`), `StaticAssets()`. |
| `internal/images/source.go` | Source rules, pure Go: `Inspect(path) Info` sniffs PNG/JPEG from the bytes (`ErrWebP` for WebP), reads the EXIF orientation (upright `Width`/`Height`) and the ICC profile description (`Profile`, `SRGB`); the encoder applies the orientation right after decoding. |
| `internal/buildtest/hero_sources.go` | Hand-built hero source fixtures: `HeroWebP`, `HeroJPEG`, `WithEXIFOrientation`, `ICCProfile` + `WithJPEGICCProfile`/`WithPNGICCProfile`, `Quadrants`, and `DecodePixels` for pixel checks. |
| `internal/views/site.css`, `stylesheet.go`, `theme.js` | The one stylesheet (inlined in every page's `<head>`, ADR-0006) and the inline theme-toggle script. |
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
| `internal/schema/` | schema.org JSON-LD node types (`Article` with `TypeArticle`/`TypeNewsArticle`, `Person`, `Organization`, `BreadcrumbList`, `WebSite` + `NewSearchAction`; `SoftwareSourceCode` for Repo Profiles only). Plain data; `views.StructuredData` marshals each `Page.StructuredData` block into `<head>`. |
| `internal/vercel/` | Build Output API writer (`static/` files, `config.json` with `Route`s) and `NewRouter`, the local router that applies a built `config.json`'s routes. |
| `internal/buildtest/` | Test support for the Build seam (fixtures, run build, parse output). |
| `subscribe/` | The subscribe function: `New(Config) http.Handler`, `ConfigFromEnv()`, form contract constants (`Path`, `SubscribedPath`, `HoneypotField`). Buttondown client in `buttondown.go`, per-IP limiter in `ratelimit.go`. Subscribe-seam tests in `subscribe_test.go`. |
| `api/subscribe.go` | Vercel Go function entry (package `handler`, `func Handler`). Wiring only. No other `.go` files here: each becomes a function. |
| `cmd/subscribe/` | Local run of the subscribe function plus the built site, behind the built `config.json` routes. |
| `builder/step_write_vercel_config.go` | The routing table (`routes()`): cache headers, 410s, `/home` 301, trailing-slash 308s, 404 fallback. |
| `builder/step_write_cloudflare_manifest.go`, `internal/cloudflare/` | The Production Target's manifest: `<out>/wrangler.json` (Workers Static Assets config, `main` -> `worker/index.js`) plus `static/_headers` and `static/_redirects`. See "Deploying". |
| `internal/deploy/`, `cmd/adapt/` | The deploy-adapter rule: `Apply(out, target)` (remove the other Target's files in place, write the Target's overlays); `make adapt TARGET=cloudflare\|vercel` runs it. Also the one map of Target spellings: `ParseTarget` reads a name (`cloudflare`), short code (`cf`) or Publish tag (`PUBCF-…`); `Target.Code()`, `Target.TagPrefix()`. |
| `cmd/target/` | `go run ./cmd/target <name\|code\|tag>` prints `target=`, `code=`, `tag_prefix=`: how `scripts/publish.sh` and `scripts/publish-gate.sh` spell a Target. |
| `builder/step_write_page_index.go`, `step_write_page_manifest.go` | Reads every written HTML page back (canonical, noindex, `<time>`s, hash of `<main>`) and diffs it against the previous build's page manifest in the build cache; the last write step saves the new manifest when the build passed. |
| `builder/step_write_feeds.go`, `step_write_sitemap.go`, `step_write_news_sitemap.go`, `step_write_robots.go`, `step_validate_indexnow_key.go`, `step_write_indexnow.go` | Full-text RSS feeds, `/sitemap.xml` + `/news-sitemap.xml`, `/robots.txt`, the IndexNow key file and changed-URL artefact. See "Feeds, sitemaps, robots and IndexNow". |
| `site.yaml` | Site config: `site_url`, `verticals`, `images` (`store`, `base_url`, `bucket`, `cache_dir`; see "Image Store"), `plausible_domain` (blank = no analytics), `ads_enabled`, `adsense_client` (required when ads are on), `ads_txt` (records written verbatim to `/ads.txt`; malformed lines fail with `ads-txt`; empty = no file), `indexnow_key` (public, stable; blank = IndexNow off), `same_as` (home Organization profiles). |
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
load fact packs 22 (Repo Profiles), load static pages 25, load briefs 25,
exclude drafts (production only; Explainers
and Briefs) 30, draw hero cards 35 (Repo Profiles without `hero`: the
card is drawn, cached and registered as the hero source); validate ads config 5, references 10, Brief references 12,
Fact Packs 15, unique slugs 20, static page paths 25, Explainer skeleton 30, Repo Profile
skeleton 32, Brief skeleton 35, section-length warning 40, Review-overdue warning 50, validate heroes 50
(Explainers and Briefs), author photos 55; transform render diagrams 10,
rewrite Brief Explainer links 20, merge Explainer frontmatter sources 30
(and a Repo Profile's Fact Pack sources), insert Repo Profile fact sections 35,
encode heroes 50, encode author photos 55;
render reset 0, write site assets 10 (and the OG badge), write hero assets 15,
write author photos 16, upload images 17 (r2 mode), render explainers 20, render briefs 30, render author pages 40,
render briefs index 40, render static pages 45, render vertical hubs 50,
render kind pages 52, render home 54, render subscribed page 60,
render search page 62, render 404 page 65, check internal links 1000
(keep it last: it reads every written page); validate IndexNow key 6, validate ads.txt 6; write
phase: index written pages 5 (reads the output back; later write steps use
`s.index`), feeds 20, sitemap 30, news sitemap 35, robots.txt 40, ads.txt 41, IndexNow
key + changed URLs 50, author placeholders 60, config.json 100, Cloudflare manifest 105, index search 800 (Pagefind over the
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

**Repo Profiles.** A `kind: repo-profile` Explainer keeps every Explainer
rule (TL;DR first, Sources last with a link, a Mermaid fence) and adds the
Repo Profile skeleton (`step_validate_repo_profile_skeleton.go`): the
body's H2s are exactly TL;DR, In 60 seconds, Architecture, Use it if,
Don't use it if, Alternatives, FAQ, Verdict, Sources, in that order
(`repo-profile-sections`; the message names the heading and its position,
or the first heading out of place); no hand-written Signals or Changelog,
which the build renders from the Fact Pack after Architecture and
Alternatives (`repo-profile-fact-section`); and five to eight `###`
questions under FAQ (`repo-profile-faq`; an H2 under FAQ is an extra
section). Heading text matches exactly (a straight apostrophe in "Don't").
The constants in that file are the only spelling of the headings; the
fact-section placement reuses them, and the editorial template
(`docs/editorial/repo-profile-template.md`) is built against them by
`builder/repo_profile_template_test.go`.

The fact sections (`step_transform_fact_sections.go`, transform 35) are
builder HTML placed into the body at the end of the written section each
follows (`doc.InsertHTMLAfterSection`, the way a rendered diagram replaces
its fence), so the page, the full-text feeds and the search index all carry
them. Signals (`partial_repo_profile_signals.templ`: `section.signals`, one
`div[data-signal=…]` row per fact: last commit, releases in the twelve
months before the fetch ("10 or more" when the ten listed all fall in it),
latest release, contributors, the top contributor's share of the listed
(top ten) contributions, open issues, security policy, license, stars,
forks, watchers, documentation files and directories) and Changelog
(`partial_repo_profile_changelog.templ`: `section.changelog`, the Fact
Pack's releases newest first, each tag linking its release page, a
pre-release marked in text and by `data-prerelease`). A fact the Fact Pack
lacks reads "not found" (`views.FactNotFound`); nothing is scored, graded
or invented. The Fact Pack's `sources` join the Sources section after the
frontmatter sources, once per URL (`step_transform_sources.go`). The
page's JSON-LD adds `SoftwareSourceCode` (full name, repository URL, and
the language, license URL and description the Fact Pack carries) beside
the Article graph (`repo_profile.go`: `softwareSourceCode`). The FAQ is
visible prose only: no `FAQPage` (dropped 2026-10-22, Google retired FAQ
results). Rule `fact-sections` fails a published
Profile lacking the section a fact section follows: a safety net the
skeleton rule makes unreachable. Both components paint with the brand
tokens only, so the dark theme applies (`repo_profile_facts_test.go`
greps their rules for colour literals).

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
(`partial_byline.templ`), which links the Author page. No `HowTo`, `Speakable`, `TechArticle` or `FAQPage`
anywhere; `SoftwareSourceCode` only on a Repo Profile (see
"Repo Profiles"), never on another kind (`TestExplainerJSONLD_NoUnsupportedTypes`).
`writePage` fails an indexable page with no `application/ld+json` block
(rule `jsonld-missing`; the page is not written), so a new page type needs
at least its BreadcrumbList.

**Design system.** One stylesheet, `internal/views/site.css`, inlined once
as the single `<style>` in every page's `<head>` (ADR-0006, superseding
ADR-0003: it saves the render-blocking round trip that held mobile LCP
above 1.5 s); no other `<style>`, no webfonts, no framework, no `url()`. Colours exist only as tokens generated from
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
checked against the 50 KB uncompressed HTML gate, not counting the shared
inline stylesheet (ADR-0006; tests measure with `buildtest.ContentSize`); an oversized page is a site-wide `html-size`
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
from `layout_analytics.templ` in every `<head>` when `plausible_domain` is set,
in production builds only: preview builds (Previews, `make dev`) carry no
analytics, so their visits are never counted and Lighthouse CI measures the
page's own code (ADR-0006).
Everything must work without JS.

**Tests.** Four seams are tested, and only these (spec "Testing
Decisions"): the Build seam, the Subscribe HTTP seam, the Fetch seam
(`internal/factpack` against an in-process GitHub) and the hero card's
text seam (`views.HeroCardText`, below). Build-seam tests live in `builder/` as package
`builder_test`, import only `builder` and `internal/buildtest`, and assert on
parsed output. Helpers in `internal/buildtest`:

- Fixtures: `ValidTree()`, `NewTree()` (default Author + `TransparencyPages()`, which the footer links), `Tree.WithExplainer(e)`, `Tree.Without(path)`, `Tree.WithPage(p)`,
  `Tree.WithAuthor(a)`, `Tree.With(path, bytes)`; `ValidExplainer()` (meets the
  full editorial skeleton: TL;DR, Mermaid fence, Sources, 1600x900 PNG hero),
  `ValidComparison()` (adds a `## Trade-offs` table), `ValidRepoProfile()` +
  `ValidFactPack()` + `Tree.WithRepoProfile(e)` (a Repo Profile of apache/iggy
  with the spec's headings, `repo`, `who_for` and its Fact Pack written
  beside it, and no `hero`, so the build draws its card; `FactPack.JSON()` spells the keys by hand, independently of
  the parser, and `e.FactPack = nil` leaves the JSON out), `ValidBrief()` +
  `Tree.WithBrief(b)` (meets the Brief skeleton; links `ValidExplainer()` by
  slug, so use it with `ValidTree()`), `RoundupProfiles()` +
  `RoundupBrief()` (the Friday roundup: five Profiles with Fact Packs and
  the Brief linking them), `DefaultAuthor()`,
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
- Deploy adapter: `Adapt(t, outDir, target)` runs the real `cmd/adapt` (what
  `make adapt` runs) and fails on a non-zero exit; `RunAdapt(t, outDir,
  target) (exitCode, output)` for refusals; `OutputFiles(t, outDir)` lists
  every output file (output-relative, sorted).
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
- Image Store: `store := NewImageStore(t)` starts the fake (one bucket,
  path-style, in memory); `store.Configure(cfg)` switches a config to r2
  against it (`ConfigureEnv(cfg, env)` with a variable deleted from
  `store.Env()` tests missing credentials); `store.Seed(key, StoredObject{…})`
  pre-loads an object; `store.Fail(method, key, status)` makes that request
  fail; inspect with `store.Requests(method)` / `RequestKeys(method)`
  (each `StoreRequest` has `Op`, `Key`, `Header`, `Body`), `store.Object(key)`,
  `store.Keys()`, and `ResetRequests()` between builds.

Subscribe-seam tests live in `subscribe/subscribe_test.go` (package
`subscribe_test`): they build the handler with `subscribe.New(Config{APIKey,
BaseURL: fake.URL})`, send form posts through `ServeHTTP`, and assert status,
`Location`/JSON body and what the `httptest` fake Buttondown received.

The hero card's text seam is `internal/views/hero_card_test.go` (package
`views_test`): which words the card draws (wordmark, repository,
description wrapped and cut with an ellipsis, the facts line without
absent facts or runes the font lacks) cannot be read back from pixels at
the Build seam, so `views.HeroCardText(views.HeroCard)` returns the lines
the renderer lays out, the same layout `HeroCardPNG` draws. Everything
else about the card (pipeline, alt text, 1600x900, crop-safe layout read
from the card the build writes under its cache, same bytes for the same
Fact Pack) is asserted at the Build seam in `builder/hero_card_test.go`.

Fetch-seam tests live in `internal/factpack/*_test.go` (package
`factpack_test`): `newHarness(t)` starts the fake GitHub (`fake_test.go`),
a `Fetcher` against it with a fixed clock and the log captured, and a
temporary content tree and cache directory; `h.gh.serveRepo("apache-iggy")`
loads `testdata/github/apache-iggy/*.json` (one recorded answer per
endpoint; a missing file is a 404, a `.headers` file adds headers);
`h.gh.answer(path, …)` / `answerSequence` override one endpoint (403s,
202s); `h.gh.requestsTo(prefix)` and `h.gh.sleptFor()` show what the
client did. The golden Fact Pack is `testdata/apache-iggy.json`. The Refresh tests also assert the exit status a Refresh maps to (3 when anything changed, 0 when nothing did) through `RefreshReport.ExitStatus()`; `cmd/fetch` only returns it.

## Repo Profiles: writing one

Start from `docs/editorial/repo-profile-template.md` (frontmatter, the
nine headings in order, five FAQ questions, the AI-assistance note) and
follow `docs/editorial/repo-profile-writing-guide.md`: copy, fetch the
Fact Pack, draft against the Fact Pack, README and release pages, check
the Preview, have the Author verify every claim, then remove `draft: true`
and tag. Nothing publishes without the Author. Which repositories qualify:
`docs/editorial/repo-profile-selection.md`. On each Profile's day 28, fill
in `docs/editorial/checkpoint-note-template.md` (the ADR-0005 metrics and
pass rule). The Friday roundup is an ordinary Brief linking the week's
Profiles (`builder/roundup_brief_test.go` proves the Brief gates take it).

## Repo Profiles: fetching a Fact Pack

A Repo Profile (ADR-0005) renders its fact sections from a Fact Pack: the
JSON beside its Markdown, named after it (`apache-iggy.md` ->
`apache-iggy.json`). The fetch command writes it; nothing writes it by
hand, and a Refresh rewrites it.

```sh
export GITHUB_TOKEN=github_pat_…   # one token; the command refuses without it
make fetch REPO=apache/iggy MD=content/explainers/distributed-systems/apache-iggy.md
make refresh                       # every Profile; the command exits 3 when anything changed (make: "Error 3")
# same as: go run ./cmd/fetch -repo apache/iggy -markdown content/explainers/distributed-systems/apache-iggy.md
#          go run ./cmd/fetch -refresh -content content   (go run turns exit 3 into 1)
```

- **Token.** `GITHUB_TOKEN`: the variable GitHub Actions provides to every
  job (1,000 requests an hour per repository). Locally, a fine-grained
  personal access token with read access to public repositories and no
  other permission. Unset, the command refuses with a message naming the
  variable (exit 2): unauthenticated GitHub allows 60 requests an hour.
  One token only; GitHub's terms forbid pooling tokens or accounts, and
  the command never scrapes github.com.
- **What it writes.** The Fact Pack beside the Markdown (`-markdown`;
  the file need not exist yet, and when it does it must be a Repo Profile
  of that repository), with keys in a fixed order, two-space indentation
  and every list present, so a Refresh diff shows values and never shape.
  The README text goes to `<cache>/readmes/<owner>/<name>/<README file>`
  for drafting (`-cache`, default `node_modules/.cache/hastcode`, the
  builder's build cache), never into the content tree: the Fact Pack
  carries only the README's URL and blob SHA. It prints each file written
  and every note (a `NOASSERTION` license, a fallback).
- **Calls.** At most ten GET requests per repository, serially: the
  repository; `languages`;
  `releases?per_page=10` (`latest_release` is the newest non-prerelease
  among them, drafts skipped); `contributors?per_page=10` (the top ten)
  and `contributors?per_page=1`, whose `Link` header's last page is the
  total (contributors with a GitHub account); `git/trees/<default
  branch>` (root files) and, when the root has `.github`,
  `contents/.github?ref=<branch>` (its URL is stable, a tree SHA is not);
  `readme` (page URL, blob SHA and the text, base64 in the same answer);
  `stargazers/count`. A `readme` 404 means no README; a
  `stargazers/count` 404, or an answer without a `count`, falls back to
  the repository's `stargazers_count` with a note.
- **Sources.** Each call that answered cites, in the Fact Pack's
  `sources`, the github.com page a Reader can open to check it, never the
  API endpoint, each page once and in call order: the repository page
  (repository, `languages`, the tree, `stargazers/count`), the releases
  page, the contributors graph page (`/graphs/contributors`), the
  `.github` tree page (`/tree/<branch>/.github`) and the README's blob
  page. The builder lists them in Sources after the frontmatter sources,
  once per URL, so an Author citing the same page in frontmatter or in
  the body's Sources section is not repeated. `watchers` is GitHub's `subscribers_count`
  (`watchers_count` is stars). A license URL is its spdx.org page;
  `NOASSERTION` (GitHub could not classify it) sets `needs_review: true`
  and prints a note telling the Author to read the LICENSE file by hand;
  it is never an error. File flags come from the root and `.github`
  listings (`internal/factpack/files.go` holds the CI and manifest names).
- **Rate limits.** Every request carries `User-Agent: hastcode-fetch/1
  (+https://hastcode.com; Repo Profile Fact Packs)`. A 403 or 429 with
  `Retry-After` waits that long; with `X-RateLimit-Remaining: 0` it waits
  until `X-RateLimit-Reset`; a secondary limit without either waits a
  minute. Any wait over two minutes, or a fourth retry, stops the command
  with `GitHub rate limit reached: GET … ; stopping instead of waiting.`
  (exit 1): run it again after the reset. A 202 (GitHub still computing
  the answer) is retried three times after 2, 4 and 8 seconds, a 5xx
  three times; both then fail with exit 1.
- **Refresh.** `-refresh` walks `content/explainers/` for every
  `kind: repo-profile`, re-fetches each Fact Pack and rewrites only those
  whose facts changed, `fetched_at` aside: an unchanged file keeps its
  bytes and its `fetched_at` (so a dormant repository's `fact-pack-stale`
  warning means "nothing has changed in 30 days"). A Profile without a
  Fact Pack gets one. It prints one line per changed repository saying
  what changed, `apache/iggy: stars 2468 -> 2600, latest release
  server-0.5.0 -> server-0.6.0 (explainers/distributed-systems/apache-iggy.json)`,
  then `refresh: N Repo Profiles checked, M changed`. Prose is never
  touched. A fetch that fails stops the Refresh (exit 1); the Fact Packs
  rewritten before it stand.
- **Exit codes.** 0 done (refresh: nothing changed); 1 a fetch or write
  failed (rate limit reached, network, I/O); 2 usage (bad flags,
  `GITHUB_TOKEN` unset); 3 refresh rewrote at least one Fact Pack.
  `make refresh` builds the binary first because `go run` reports every
  non-zero status as 1.

### Refresh

A Refresh is manual (Author decision, 2026-10-22; it replaced a scheduled
GitHub Actions workflow). From the repository root, with a `GITHUB_TOKEN`
in the environment (`gh auth token` works):

```sh
GITHUB_TOKEN=$(gh auth token) make refresh
```

It re-fetches every Repo Profile's Fact Pack under `content/`, rewrites
only the changed ones (an unchanged Fact Pack keeps its `fetched_at`) and
prints one line per changed repository plus any `note:` lines for the
Author (a `NOASSERTION` license, for instance). The fetch binary exits 3
when anything changed and 0 when nothing did; make reports 3 as `Error 3`,
which is the "something changed" signal, not a failure. Then review the
diff (facts only: a Refresh never touches prose; re-read the Profile's
prose against the new facts when a release or license changed), run
`make check`, commit and push. Readers see the new facts after the next
Publish tag (`scripts/publish.sh vl`, then `cf`; ADR-0004). GitHub Actions
only tests, builds Previews and publishes to the Targets; nothing in CI
fetches from GitHub.

Rate limit: a personal token allows 5,000 requests an hour and a Refresh
makes about 9 per Profile, so hundreds of Profiles fit in one run. Cadence:
weekly, Monday, so the Changelog stays fresh; the 90-day Review re-reads
the prose.

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

Environment (Vercel project settings, environment `Production`: every QA
deploy is `vercel deploy --prod`, and Vercel serves nothing but QA):

| Variable | Required | Notes |
|----------|----------|-------|
| `BUTTONDOWN_API_KEY` | yes, secret | Buttondown API key: the test newsletter's (the setup wizard sets it). Never sent to the browser. |
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
The setup wizard no longer sets Vercel up (its Vercel stages were removed 2026-10-25 at the Author's request); configure the QA deployment by hand if the QA Target is kept.

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

## Image Store

Every encoded image (hero variants and crops, Author photos, the OG badge)
goes through one seam (`builder/image_store.go`, ADR-0004). `site.yaml`:

```yaml
images:
  store: local          # local (default) | r2
  base_url: https://images.hastcode.com/   # public URL of the bucket (r2)
  bucket: hastcode-images                  # R2 bucket (r2)
  cache_dir: node_modules/.cache/hastcode/images
```

`images.store` in `site.yaml` is the one switch for every build, local and
CI; the workflows pass no flag. **It is `local` today** (ADR-0004 as amended
2026-10-25: R2 needs a payment method on the Cloudflare account, and local
images fit the free plan's 20,000-file cap for roughly 1,300 pieces at about
15 files each). To move to R2: enable R2, run the wizard's stages 4, 6 and
19, add the two R2 secrets, set `store: r2`, Publish.

- **`local`** (current; `make dev`, tests, CI): images are written under
  `/assets/` in the output, ship inside each deploy, and pages use site
  paths. Needs no credentials.
- **`r2`** (`store: r2`, or `-image-store r2` / `make build
  FLAGS=-image-store=r2` for a one-off): nothing image-shaped is written to the
  output; every image URL in HTML, JSON-LD, Open Graph and feeds is
  absolute on `base_url` (`s.absURL` passes absolute URLs through, so
  nothing is double-prefixed). For each key the build sends `HEAD` and,
  only on 404, `PUT` with `Content-Type` (`image/avif`, `image/jpeg`,
  `image/png`) and `Cache-Control: public, max-age=31536000, immutable`, on
  8 workers. A non-2xx answer fails the build naming the key.
- **Credentials come from the environment only**: `CLOUDFLARE_ACCOUNT_ID`,
  `R2_ACCESS_KEY_ID` (the R2 API token's id), `R2_SECRET_ACCESS_KEY`
  (SHA-256 of the token value); token permission Object Read & Write on the
  bucket. Optional `R2_ENDPOINT` overrides
  `https://<account id>.r2.cloudflarestorage.com` (a jurisdiction endpoint,
  or a test fake). An r2 build missing any of them, or `base_url`/`bucket`,
  fails before loading or encoding anything (rule `image-store`, naming the
  variable).
- **Client** (`internal/imagestore`): aws-sdk-go-v2 built from explicit
  options (nothing read from `~/.aws` or `AWS_*`), region `auto`,
  `BaseEndpoint`, path-style requests, request checksum calculation and
  response checksum validation `WhenRequired` (R2 does not document the
  SDK's default CRC32 on single PUTs; research §1.3).
- **Key layout**: the v1 asset path without `/assets/`:
  `heroes/<source name>-<w>x<h>.<hash10>.{avif,jpg}`,
  `authors/<photo name>-400x400.<hash10>.jpg`, `og-badge.<hash12>.png`. The
  hash is of the encoded bytes, so an existing key never needs re-uploading
  and changing `base_url` means a full rebuild and republish.
- **Keys a build references**: `Result.ImageKeys` (sorted, unique; the same
  in both modes), which the prune command diffs against the bucket.

### Pruning orphans

Builds only ever add objects: a re-exported hero, a removed piece or a
retuned encoder leaves the old keys in the bucket. The prune command
lists them. **It is manual and a dry run by default**; nothing runs it
automatically.

```sh
make prune-images                          # dry run: list orphans, delete nothing
make prune-images PRUNE_FLAGS=--delete     # delete the listed orphans
# same as: go run ./cmd/prune -content content -config site.yaml [--delete]
```

- It reads the same `site.yaml` (`images.bucket`, `images.base_url`) and the
  same environment as an r2 build (`CLOUDFLARE_ACCOUNT_ID`,
  `R2_ACCESS_KEY_ID`, `R2_SECRET_ACCESS_KEY`, optional `R2_ENDPOINT`),
  whatever `images.store` says. A missing variable fails with the
  builder's message (rule `image-store`) before anything is built or
  requested.
- It builds the content tree twice with the local store into a temporary
  directory, Production and Preview, and keeps the union of their
  `ImageKeys` (nothing is uploaded, no search index is built). It then
  lists the whole bucket (`ListObjectsV2`, 1,000 keys a page) and prints
  every other object with its size, then the count and total bytes.
- If the Production build fails (beyond the skipped search index) it
  refuses and prints the errors: a failed piece drops its images from the
  key set and they would look orphaned. Fix the content first.
- `--delete` removes exactly the listed orphans with `DeleteObjects`, at
  most 1,000 keys per request.
- **Preview images may be pruned.** A draft's images are kept only while
  the draft is in the content tree you prune from (normally `main`). Images
  a Preview uploaded for a draft that lives only on another branch, or for
  an earlier version of a hero, are orphans and `--delete` removes them; the
  next Preview of that branch uploads them again.
- Pages already published elsewhere (an old deployment, a cached feed or
  link preview) may still point at an orphan; deleting it breaks that
  image. Prune after the new Production deploy is live, not before.
- Not yet verified live: whether R2 accepts the checksum header the SDK
  sends on `DeleteObjects` (a checksum-required operation, research
  §1.3). Run the first `--delete` by hand and check the result: the setup
  wizard's stage 16 does it against a dummy object (`prune-check.txt`).
- Credentials for a local run: `set -a; . ./.env; set +a; make
  prune-images` (the wizard saves `CLOUDFLARE_ACCOUNT_ID`,
  `R2_ACCESS_KEY_ID` and `R2_SECRET_ACCESS_KEY` there).

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
  Explainers never appear. The file only changes when a Production Publish
  rebuilds the site, so a Brief drops out at the first Publish after it is
  48 h old, not on the dot.
- **`/robots.txt`**: allow all, explicit allow groups for `Googlebot`,
  `Bingbot`, `OAI-SearchBot`, `PerplexityBot`, both sitemap lines. Training
  crawlers (`GPTBot`, `Google-Extended`) fall under `*` (allowed). That is
  Production's file; on the QA Target (Vercel) the deploy adapter
  overlays `User-agent: *` / `Disallow: /` (see "Deploying").
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

Post-deploy ping: only a Production Publish (`PUBCF-*`) pings. The
Cloudflare publish job moves the artefact out of `.vercel/output` before
`wrangler deploy` and, once the deploy is live and the repository variable
`INDEXNOW_ENABLED` is `true`, runs the equivalent of

```sh
curl -fsS -X POST https://api.indexnow.org/indexnow \
  -H 'Content-Type: application/json; charset=utf-8' --data-binary @indexnow.json
```

200/202 are success; 403 means the key file is not live yet; 422 means a
URL is not on `host`. Only Production Publishes ping, so
`INDEXNOW_ENABLED` stays `false` until hastcode.com is served by Cloudflare
(the setup wizard flips it at cutover and deletes the `hastcode-indexnow-*`
Actions caches, which hold only the page manifest, so the first Publish
after cutover lists every page while encoded images and diagrams stay
cached). The manifest the artefact is diffed against is Production's own
(only the Cloudflare job restores and saves it), so a QA Publish never hides a change from the next Production
ping. A failed ping fails the job before the build cache is saved, so the
next Production Publish lists those URLs again.

## Routes and caching

This section is the Vercel (QA Target) routing; the Production Target's
equivalent is the Cloudflare manifest (`wrangler.json`, `_headers`,
`_redirects`; see "Deploying").

`builder/step_write_vercel_config.go` writes every route into
`.vercel/output/config.json` (Build Output API v3 `routes`; patterns stay
RE2-compatible, no lookaround, so the local router can apply the same file).
In order:

0. `X-Robots-Tag: noindex` on every response, whatever its status
   (`continue`): Vercel is only ever the QA Target (ADR-0004), so nothing
   it serves may be indexed. `make dev` sends it too. Canonical, Open Graph
   and JSON-LD URLs still point at `site_url` (`https://hastcode.com`).
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
Vercel itself: QA sits behind Vercel Authentication, so it is checked in a
signed-in browser by hand (the wizard's Vercel stages were removed
2026-10-25); the Production equivalents are curled automatically against
the first Preview, the workers.dev URL and hastcode.com (stages 12, 13, 15). Note that `vercel build`'s merged `config.json` (with
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
- **Caching.** `/pagefind/*` keeps the host's default revalidating cache
  (Workers and Vercel alike): its entry files have stable names (see
  "Routes and caching" and "Deploying").
- **Tests** (`builder/search_test.go`) read the index through
  `buildtest.SearchIndex(t, out)` (URL -> indexed text, decoded from the
  gzipped fragments). `buildtest.Build` skips a test when `node` or
  `node_modules/pagefind` is missing, unless the config skips indexing.

**Build command.** `npm ci && make build`: CI runs `npm ci` itself, then
`make build` (Production Publishes and Previews) or `vercel build`, which
runs `make build` (QA Publishes; vercel.json; see "Deploying"). No
separate `npx pagefind` step: the builder runs Pagefind and fails the build
if the index is missing. Previews build with `MODE=preview`,
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
  inline site stylesheet (ADR-0006). Each SVG maps its `--bg`/`--fg`/`--accent` to
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
  `actions/cache`, per Target (restored after `npm ci`, which deletes
  `node_modules/`; see "Deploying").
  Fully cached builds never start Node.
- **CI build:** `npm ci` runs before `make build` from the repo root, with
  Node 24 (Node 22+ works).

## Deploying

### One build, two Targets (ADR-0004)

One `make build` writes one output directory (`$(OUT)`, default
`.vercel/output`) that holds the static tree and a manifest per Target:

| File (relative to `$(OUT)`) | Target | What it is |
|---|---|---|
| `static/**` (except the two below) | both | The site. |
| `config.json` | Vercel (QA) | Build Output API routes ("Routes and caching"), incl. `X-Robots-Tag: noindex` on everything. |
| `static/robots.txt` | both, overlaid for Vercel | Production's real file; `make adapt TARGET=vercel` overwrites it with `User-agent: *` / `Disallow: /`. |
| `wrangler.json` | Cloudflare (Production) | Workers Static Assets config. |
| `static/_headers` | Cloudflare | Cache policy; parsed by Workers, never served. |
| `static/_redirects` | Cloudflare | `/home` and `/home/` -> `/` 301; parsed, never served. |

**The Cloudflare manifest.** `wrangler.json` is plain JSON (Wrangler reads
`wrangler.json`, `.jsonc` or `.toml`; the builder marshals JSON). Wrangler
only discovers a config in the working directory and its parents and
resolves relative paths against the config file, so deploy from the repo
root with `wrangler deploy --config .vercel/output/wrangler.json`
(`wrangler dev` and `wrangler versions upload` take the same flag). Keys:
`name: "hastcode"`; `main` is `worker/index.js` written relative to the
output (`../../worker/index.js` for the default `OUT`; the builder anchors
it on `-tools`, the project root); `compatibility_date` (a constant in the
step file; never later than the pinned Wrangler's runtime);
`workers_dev: true` and `preview_urls: true` (smoke-test and Preview URLs);
no `routes`, so a deploy never attaches hastcode.com by itself (the DNS flip
is a wizard step); and `assets`:

- `directory: "static"`, `binding: "ASSETS"` (the Worker's `env.ASSETS`).
- `html_handling: "force-trailing-slash"`: `/x` and `/x/index.html` -> 307
  `/x/` (Vercel answers 308; both are fine for GET).
- `not_found_handling: "404-page"`: every unknown path gets `/404.html` with
  status 404, including the retired `/video/…`, `/community/…`, `/group/…`,
  which have no rule anywhere (ADR-0001 as amended: 404, not 410, on
  Production).
- `run_worker_first: ["/api/subscribe"]`: the signup Worker is the only code
  that runs; every other request is a free static-asset hit.

`_headers`: `/assets/*` `public, max-age=31536000, immutable`; `/`, `/*/`
and `/*.html` `public, max-age=0, must-revalidate` (Workers joins a header
two rules set, so the patterns never overlap); every other file keeps
Workers' default (the same revalidate value); `*.workers.dev` hostnames
(the `workers_dev` URL and Previews) get `X-Robots-Tag: noindex`. Headers
and redirects never apply to the Worker's own responses.

Check it locally against real Wrangler (no Cloudflare login needed):
`make build MODE=preview FLAGS=-skip-search-index`, then
`npx wrangler deploy --dry-run --config .vercel/output/wrangler.json`
(validates the config and bundles the Worker), or `npx wrangler dev
--config .vercel/output/wrangler.json` and curl: `/home` 301 `/`,
`/video/x` 404 page, `/about` 307 `/about/`, `/assets/*.css` immutable,
`/_headers` 404. `scripts/check-worker.sh` covers the signup path.

**The deploy-adapter rule** (`internal/deploy`): each Target's file set is
every file in the output except those only the other Target reads (table
above), with the Target's overlays written in place (today one: Vercel's
disallow-all `robots.txt`). Neither host may serve the other's configuration. Every deploy step
runs the adapter once, after the build and before the upload; the output
then belongs to that Target:

```sh
make build                      # or `vercel build` (runs make build)
make adapt TARGET=cloudflare    # removes config.json
npx wrangler deploy --config .vercel/output/wrangler.json

make adapt TARGET=vercel        # removes wrangler.json, static/_headers, static/_redirects;
                                # overlays a disallow-all static/robots.txt
vercel deploy --prebuilt
```

Tag prefix to Target: `PUBCF-*` -> `cloudflare`, `PUBVL-*` -> `vercel`. A
new per-Target file or overlay (like the QA `robots.txt`) goes in
`internal/deploy`, never in a workflow.

### Publishing: a tag per Target

Pushing `main` runs the tests and nothing else. Going live is a tag on a
`main` commit, created and pushed by one command:

```sh
git switch main && git pull     # clean tree, HEAD on origin/main
scripts/publish.sh vl           # QA Target: tag PUBVL-<yyyymmdd>-<HHMMSS> (UTC)
# check https://qa.hastcode.com/ (Vercel Authentication), then the same commit:
scripts/publish.sh cf           # Production Target: tag PUBCF-<yyyymmdd>-<HHMMSS>
scripts/publish.sh vl --dry-run # every check, then `git push --dry-run`; no tag kept
```

`publish.sh` refuses with uncommitted changes or untracked files, on any
branch but `main`, and when `HEAD` is not on `origin/main` (push first).
It creates a lightweight tag on `HEAD`, pushes it, prints the workflow
link, and after `vl` prints the QA URL and the reminder to run `cf` while
`HEAD` is still the same commit. Never type a tag by hand: the workflow
only accepts the exact `PUB(CF|VL)-<8 digits>-<6 digits>` shape.

`.github/workflows/publish.yml` (tags `PUBCF-*`, `PUBVL-*`) runs in two
hops:

1. **Gate** (the tag-push run): `scripts/publish-gate.sh` reads the Target
   from the prefix (`PUBCF` -> `cloudflare`, `PUBVL` -> `vercel`, through
   `go run ./cmd/target`: the spellings live in `internal/deploy`), peels
   the tag to its commit and fails, with the reason in the run's
   annotations, unless that commit is `origin/main` or an ancestor of it
   (`git merge-base --is-ancestor`, checkout with `fetch-depth: 0`). It
   then starts the same workflow on `main` with the tag as input
   (`gh workflow run publish.yml --ref main -f tag=…`).
2. **Deploy** (the "Publish PUB…" run on `main`): waits its turn (below),
   runs the gate again, and stops there, green, with "Skipped <tag>:
   superseded by <newer tag>" in the run summary when a newer tag for the
   same Target is on `main`. Otherwise the
   Target's job checks out the gated commit and runs `make test`, one
   production-mode build (images per `images.store`), the Target's adapter
   and the upload:
   - **Cloudflare** (`production` environment): `make build MODE=production`
     -> set `indexnow.json` aside -> `make adapt
     TARGET=cloudflare` -> `cloudflare/wrangler-action@v4` running `wrangler
     deploy --config .vercel/output/wrangler.json --message <tag>` (Wrangler
     4.149.0, as in `package.json`) -> IndexNow ping when
     `INDEXNOW_ENABLED` is `true` -> save the cache.
   - **Vercel** (`qa` environment): `vercel pull --environment=production`
     -> `vercel build --prod` (runs `make build` with `MODE=production`;
     `@vercel/go` adds the subscribe function) ->
     delete `indexnow.json` -> `make adapt TARGET=vercel` (QA's disallow-all
     `robots.txt`) -> `vercel deploy --prebuilt --prod` -> save the cache.
     QA never pings IndexNow.

Why two hops: the build cache (`node_modules/.cache/hastcode`: encoded
images, rendered diagrams, the IndexNow page manifest) must carry from one
Publish to the next, and GitHub scopes caches by ref ("Workflow runs also
cannot restore caches created for different tag names"). A run on `main`
saves into `main`'s scope, which every later Publish and every branch can
restore. Build-cache keys are per Target (`hastcode-cloudflare-*`,
`hastcode-vercel-*`; path list `images/` + `mermaid/` under
`node_modules/.cache/hastcode`, identical in `preview.yml`, since the path
list is part of a cache's version), and Production's IndexNow page manifest
(`pages/`) is a cache of its own (`hastcode-indexnow-*`, Cloudflare job
only), so Production's IndexNow diff is always against the previous
Production build and can be reset without re-encoding every image. Running "Publish" by hand from the Actions tab
(ref `main`, input `tag`) re-deploys an existing tag through the same gate,
if it is still its Target's newest.

**Newest tag wins.** Deploy runs have one concurrency group per Target
(`publish-cloudflare`, `publish-vercel`), never cancelled, so Production
deploys run one at a time and QA never waits on Production or the reverse.
The group has `queue: max`: every waiting run stays queued (GitHub's
default keeps one and silently cancels the others). Gates finish in any
order, so runs can queue out of tag order; that is why each deploy run,
once it starts, checks for a newer tag (tag names embed the UTC time
`publish.sh` created them) whose commit is on `main`, and skips itself
with the note above if there is one. So an older Publish never deploys
over a newer one, the newest always deploys, and every skipped tag says
so in its own run. A later Publish of an older commit is still a newer
Publish: to roll back, put the old commit on `main` (revert) and publish
again, or tag it with `publish.sh` from a `main` checkout at that commit.
`make lint-workflows` ignores actionlint's one complaint here (1.7.12
predates `queue`). A failed build or test deploys nothing to either
Target.

### Previews and Lighthouse

`.github/workflows/preview.yml` (push to any branch but `main`): `make
test`, and in parallel a preview-mode build (drafts included, images to
R2, cache restored from the branch's own or else the last Production
Publish's) -> `make adapt TARGET=cloudflare` -> `wrangler versions upload
--config .vercel/output/wrangler.json`. A version upload is not a deploy:
hastcode.com keeps serving the last `wrangler deploy`, and the upload gets
its own Version Preview URL (`<version prefix>-hastcode.<subdomain>.workers.dev`;
`preview_urls: true` in `wrangler.json`; `*.workers.dev` hosts send
`X-Robots-Tag: noindex`), posted in the run summary. Branch pushes never
reach Vercel: it is only the QA Target, deployed by Publish tags.

Then **Lighthouse CI** (`lighthouserc.json`) runs against that URL: home,
`/search/`, `/briefs/`, a Vertical hub, an Explainer and a Brief
(`scripts/lighthouse-urls.sh` picks them from the build's sitemap).
Budgets, mobile emulation with simulated throttling: LCP ≤ 1500 ms
(median of 3 runs), CLS 0 (worst run), first-party script transfer
(compressed) ≤ 60 KB as a secondary guard; the spec's total-JS budget
(uncompressed, analytics included) is enforced by `make test`. A miss
fails the check; reports are the `lighthouse-reports` artefact (never
public storage: previews hold drafts). Preview builds carry no analytics
(ADR-0006); Lighthouse also blocks Plausible and ad hosts as a backstop.
Found on the first Preview (ADR-0006): Lighthouse's simulation makes LCP
wait for every request that starts before it, including a deferred
analytics script and the favicon, so neither may be requested during first
paint.

`.github/workflows/tests.yml` (push to `main`): `make test` and `make
lint-workflows` (pinned actionlint, plus shellcheck on `run:` scripts),
nothing else. Run `make lint-workflows` after editing a workflow.

### Secrets, variables and settings

GitHub repository secrets (workflows reference names only):

| Secret | Used by | Value |
|---|---|---|
| `CLOUDFLARE_API_TOKEN` | publish (Cloudflare), preview | Account API token from the "Edit Cloudflare Workers" template, scoped to this account (and zone). |
| `CLOUDFLARE_ACCOUNT_ID` | publish (both), preview | Cloudflare account id; also the R2 endpoint's account. |
| `R2_ACCESS_KEY_ID` | publish (both), preview; only when `images.store` is `r2` | R2 API token's access key id (Object Read & Write on `hastcode-images`). |
| `R2_SECRET_ACCESS_KEY` | publish (both), preview; only when `images.store` is `r2` | That token's secret access key. |
| `VERCEL_TOKEN` | publish (Vercel) | Vercel account token. |
| `VERCEL_ORG_ID` | publish (Vercel) | From `.vercel/project.json` after `vercel link`. |
| `VERCEL_PROJECT_ID` | publish (Vercel) | Same file. |

Repository variable `INDEXNOW_ENABLED`: `false` until hastcode.com is
served by Cloudflare, then `true`. GitHub environments `production` and
`qa` are created on first use (protection rules optional). The publish
workflow's hand-off job needs `actions: write`, which it requests itself.

Outside GitHub: the Cloudflare Worker `hastcode` (first `wrangler deploy`
creates it; version uploads need it to exist), its workers.dev subdomain
and Preview URLs on, its secret `BUTTONDOWN_API_KEY` (no dashboard vars:
`wrangler deploy` deletes them), hastcode.com and www.hastcode.com
attached as its custom domains at cutover plus a `www` -> apex Redirect
Rule; the R2 bucket `hastcode-images` with `images.hastcode.com`
connected; the Vercel project (Hobby, Vercel Authentication with scope
All Deployments, `qa.hastcode.com` assigned to its `Production`
environment, env `BUTTONDOWN_API_KEY` with the test newsletter's key). `vercel.json` keeps Vercel's Git deployments off:
`framework: null`, `installCommand: ""` (CI installs), `buildCommand:
"make build"`; `api/subscribe.go` is picked up by zero-config detection.

Why CI and not the hosts' own builds: the build needs Go on `PATH` (Vercel
gives Go only to `@vercel/go`), the IndexNow artefact exists only on the
build machine and is POSTed after the deploy is live, and the build cache
persists in `actions/cache` under our control. Shared setup:
`.github/actions/setup` (Go from `go.mod`, Node 24, `npm ci`, optionally
`vercel@65.0.0`).

### First-time setup and first run

`scripts/setup-deploy.sh` walks you through every account step in order,
writes each value where it belongs, and runs the first-run checklists. It
opens each dashboard page, says what to click and copy, asks before
anything irreversible (a push, a deploy, a Publish, the DNS flip, a
prune), and ends with everything still to do. Re-running is safe: values
are kept in `.env` (gitignored, mode 0600). To resume, pass the stage
number: values from earlier stages then come from `.env`.

```sh
scripts/setup-deploy.sh             # all 19 stages
scripts/setup-deploy.sh --from 13   # e.g. once the content threshold is met
```

| # | Stage | Writes |
|---|---|---|
| 1 | Tools and logins: git, gh, go, node, make, curl, jq, `npm ci` (pinned Wrangler), `vercel@65.0.0`; `gh`, `vercel` and `wrangler` logins | |
| 2 | GitHub repo (private, remote `origin`); Actions must allow `cloudflare/wrangler-action`; offers to delete a leftover v1 Lighthouse secret | `.env` `GITHUB_REPO` |
| 3 | Cloudflare account that holds the hastcode.com zone; workers.dev subdomain | `.env` `CLOUDFLARE_ACCOUNT_ID` |
| 4 | R2 bucket `hastcode-images` with custom domain `images.hastcode.com` (r2.dev off, no CORS). Skipped while `images.store` is `local` | |
| 5 | API token from the "Edit Cloudflare Workers" template (account + hastcode.com zone) | `.env` `CF_WORKERS_API_TOKEN` (not `CLOUDFLARE_API_TOKEN`, which would override your `wrangler login`) |
| 6 | R2 API token: Object Read & Write on `hastcode-images` only. Skipped while `images.store` is `local` | `.env` `R2_ACCESS_KEY_ID`, `R2_SECRET_ACCESS_KEY` |
| 7 | Buttondown: the real newsletter's key and a test newsletter's key | `.env` `BUTTONDOWN_PROD_API_KEY`, `BUTTONDOWN_TEST_API_KEY`, `BUTTONDOWN_API_KEY` (= test, for `make dev`) |
| 8 | First `wrangler deploy` of a local preview-mode build, so the Worker `hastcode` exists before any Preview; Worker secret `BUTTONDOWN_API_KEY` = test key | `.env` `WORKER_URL` |
| 9 | What survives a `wrangler deploy` (no routes, no vars in the config): a dashboard custom domain on a throwaway `deploy-check.hastcode.com`, a dashboard variable, the secret | |
| 10 | Vercel token; every GitHub secret in the table above; `INDEXNOW_ENABLED=false` | `.env` `VERCEL_TOKEN`; GitHub secrets and variable |
| 11 | Push `main` (tests only) | |
| 12 | First branch Preview (`deploy-check`, empty commit) and its checklist; `scripts/check-worker.sh` | `.env` `PREVIEW_URL` |
| 13 | ADR-0001 content gate (≥ 30 published Explainers across all six Verticals, ≥ 20 Briefs, no `TODO(author)`), then the first `scripts/publish.sh cf`, while hastcode.com still points at the old host | |
| 14 | Production DNS flip: Worker secret to the real newsletter's key; delete the old @ and `www` records; attach hastcode.com and www.hastcode.com to the Worker as custom domains; Redirect Rule `www` -> apex | |
| 15 | Go live: `INDEXNOW_ENABLED=true`, delete the `hastcode-indexnow-*` caches (the page manifest only), `scripts/publish.sh cf`, the live checklist; Plausible hostnames, Search Console and Bing sitemaps | GitHub variable |
| 16 | First `make prune-images PRUNE_FLAGS=--delete` by hand (against a dummy `prune-check.txt`) and the R2 billing checks. Skipped while `images.store` is `local` | |

What only a live run can verify, and the stage that checks it (each line
in the script names its user story or research question):

- R2 accepts single PUTs with checksums set to "when required" (research
  Q1): stage 13, the first Publish build.
- Custom-domain responses carry the object's `Cache-Control` and are
  cached at the edge (`cf-cache-status: HIT`) (research Q4): stages 13, 15.
- A HEAD miss then a PUT costs one Class A operation; `DeleteObjects` is
  free like `DeleteObject` (research Q2, Q3): stage 16, from R2 metrics.
- R2 accepts the SDK's checksum header on `DeleteObjects`: stage 16.
- The lightweight Publish tag resolves to the tagged commit (research Q7):
  stages 13, 15.
- Workers Preview URLs are reachable by Lighthouse CI: stage 12.
- A dashboard custom domain survives `wrangler deploy` without `routes`
  (stage 9 on a throwaway host, stage 15 on hastcode.com); if it does not,
  add `routes` with `custom_domain: true` to the builder's `wrangler.json`
  before the flip. Dashboard variables do not survive a deploy (Wrangler
  deletes them without `--keep-vars`), so the Worker has none; its secret
  does.
- The Worker's 13 signup cases (`scripts/check-worker.sh`, locally) and
  the safe subset against the live Preview: stage 12.

Local equivalent of a CI build, for debugging: `npm ci`, `make build`,
`make adapt TARGET=cloudflare`, `npx wrangler deploy --dry-run --config
.vercel/output/wrangler.json` (or `vercel pull` + `vercel build` with the
project linked). For page checks use `make dev`.

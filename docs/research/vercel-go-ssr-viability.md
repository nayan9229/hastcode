# Go SSR on Vercel for hastcode.com: viability and constraints

Research date: 2026-10-20. All Vercel pages were read on that date; each page's own `last_updated` stamp is noted where Vercel publishes one. Vercel changes limits often, so re-check the cited URLs before committing to a number.

Scope: Go backend, server-rendered HTML, deployed on Vercel, for an SEO-first content site (tech news + architecture explainers) serving US / Europe / India at 1–3M pageviews/month.

---

## TL;DR

- **Verdict: viable-with-caveats.** Go on Vercel works today and is listed as an *official* runtime, but the Go runtime page itself says it is "available in Beta on all plans", and Go is **not** in the published Fluid Compute runtime list (Node.js, Python, Edge, Bun, Rust). Treat Go as a second-tier runtime relative to Node.js: no documented streaming, no large-function support, no extended max duration, and no published cold-start figures. ([Go runtime](https://vercel.com/docs/functions/runtimes/go), [Fluid compute](https://vercel.com/docs/fluid-compute), [Runtimes](https://vercel.com/docs/functions/runtimes))
- **Compute runs in one region by default (`iad1`, Washington DC) and Hobby is locked to a single region.** India and EU visitors reach the function over a transcontinental hop on every cache miss. Pro allows multi-region (docs disagree: "5 regions" on the region page vs "Up to 3" on the Fluid page). ([Regions config](https://vercel.com/docs/functions/configuring-functions/region), [Fluid compute defaults](https://vercel.com/docs/fluid-compute#default-settings-by-plan))
- **The fix is caching, not compute placement.** Return `Cache-Control: max-age=0, s-maxage=86400, stale-while-revalidate=...` (or `Vercel-CDN-Cache-Control`) from the Go handler and the HTML is cached by Vercel's CDN. But the CDN cache is **segmented per region** and best-effort, so each of the 19 regions warms independently and rarely-hit pages get evicted. ([CDN cache](https://vercel.com/docs/caching/cdn-cache), [Cache-Control headers](https://vercel.com/docs/caching/cache-control-headers))
- **Static HTML beats SSR for this content type.** A Go build step can emit `.vercel/output/static/**/*.html` via the Build Output API v3; static files are "automatically cached on Vercel's global network" and need no function at all. Use SSR (or Build Output API Prerender Functions = ISR) only for the long tail and dynamic bits. ([Build Output API primitives](https://vercel.com/docs/build-output-api/primitives), [ISR](https://vercel.com/docs/incremental-static-regeneration))
- **Hobby is not an option for this site.** Vercel's Terms: "You shall only use the Services under a Hobby plan for your personal or non-commercial use", and the fair-use page lists "the inclusion of advertisements, including ... Google AdSense" as commercial. Hobby also includes only 1M CDN requests/month, which 1M pageviews will exceed on assets alone. ([Terms §4](https://vercel.com/legal/terms), [Fair use](https://vercel.com/docs/limits/fair-use-guidelines))
- **Pro cost at 1–3M pageviews/month is dominated by CDN requests, not Go compute.** Pro is $20/mo with a $20 credit and the included Flat Rate CDN tier is 1M requests + 1TB; the next tiers are $20/mo (10M requests) and $100/mo (50M requests). With HTML cached at the edge, Go invocations stay in the low hundreds of thousands. ([Pro plan](https://vercel.com/docs/plans/pro-plan), [Flat Rate CDN](https://vercel.com/docs/pricing/flat-rate-cdn))
- **Core Web Vitals:** Google's TTFB "good" threshold is ≤0.8s at the 75th percentile. Edge cache HITs from a nearby region clear this easily; a MISS from Mumbai to `iad1` with a cold Go start is the case that threatens it, so the architecture must keep the HIT ratio high for hot pages and accept MISS latency only on the long tail. ([web.dev TTFB](https://web.dev/articles/ttfb))

---

## Limits table (Vercel, as read 2026-10-20)

| Limit | Hobby | Pro | Source |
|---|---|---|---|
| Function max duration | 300s default and max | 300s default, 800s max (1800s beta, Node/Bun/Python only) | [Limits](https://vercel.com/docs/functions/limitations) (last_updated 2026-08-24) |
| Function memory / CPU | 2 GB / 1 vCPU (default and max) | 2 GB / 1 vCPU default, 4 GB / 2 vCPU max | [Limits](https://vercel.com/docs/functions/limitations#memory-size-limits) |
| Bundle size (uncompressed) | 250 MB | 250 MB (5 GB "large functions" beta is Node/Bun/Python only) | [Limits](https://vercel.com/docs/functions/limitations#bundle-size-limits) |
| Request/response body | 4.5 MB | 4.5 MB | [Limits](https://vercel.com/docs/functions/limitations#request-body-size) |
| Concurrency | auto-scales to 30,000 | auto-scales to 30,000 (Enterprise 100,000+) | [Limits](https://vercel.com/docs/functions/limitations) |
| Function regions | Single region (default `iad1`) | "5 regions" per region page; "Up to 3" per Fluid page (inconsistent) | [Region config](https://vercel.com/docs/functions/configuring-functions/region#limits), [Fluid defaults](https://vercel.com/docs/fluid-compute#default-settings-by-plan) |
| Function failover regions | – | – (Enterprise only) | [Region config](https://vercel.com/docs/functions/configuring-functions/region#automatic-failover) |
| Functions per deployment (non-framework `api/`) | 12 | unlimited | [Runtimes](https://vercel.com/docs/functions/runtimes#functions-created-per-deployment) |
| Function archiving | archived after 2 weeks idle (prod); un-archive cold start "at least 1 second longer" | same | [Runtimes](https://vercel.com/docs/functions/runtimes#archiving) |
| File descriptors | 1,024 shared | 1,024 shared | [Limits](https://vercel.com/docs/functions/limitations#file-descriptors) |
| Writable disk | `/tmp` up to 500 MB | same | [Runtimes](https://vercel.com/docs/functions/runtimes#file-system-support) |
| Env var total size (Go runtime) | 64 KB | 64 KB | [Limits page](https://vercel.com/docs/limits) (last_updated 2026-09-16) |
| CDN max cacheable function response | 10 MB (20 MB streaming) | same | [CDN cache limits](https://vercel.com/docs/caching/cdn-cache#limits) |
| CDN max cache time | 1 year (`s-maxage`, `max-age`, `stale-while-revalidate`) | same | [CDN cache limits](https://vercel.com/docs/caching/cdn-cache#limits) |
| Fast Data Transfer (bandwidth) included | 100 GB/mo | 1 TB/mo via included Flat Rate CDN tier | [Hobby plan](https://vercel.com/docs/plans/hobby), [Flat Rate CDN](https://vercel.com/docs/pricing/flat-rate-cdn) |
| CDN Requests (formerly Edge Requests) included | 1,000,000/mo | 1,000,000/mo via included Flat Rate CDN tier; $20/mo tier = 10M, $100/mo = 50M, $300/mo = 150M | [Hobby plan](https://vercel.com/docs/plans/hobby), [Flat Rate CDN](https://vercel.com/docs/pricing/flat-rate-cdn) |
| Function invocations included | 1,000,000/mo | none included; $0.60 per 1M on demand (offset by $20 credit) | [Fluid pricing](https://vercel.com/docs/functions/usage-and-pricing) |
| Active CPU included | 4 CPU-hrs/mo | none included; $0.128/hr (iad1) to $0.221/hr (gru1) | [Fluid pricing](https://vercel.com/docs/functions/usage-and-pricing#regional-pricing) |
| Provisioned memory included | 360 GB-hrs/mo | none included; $0.0106/GB-hr (iad1) to $0.0183 | [Fluid pricing](https://vercel.com/docs/functions/usage-and-pricing#regional-pricing) |
| Fast Origin Transfer included | 10 GB/mo | usage-based, $0.06–$0.43/GB by region | [Hobby plan](https://vercel.com/docs/plans/hobby), [Regional pricing](https://vercel.com/docs/pricing/regional-pricing) |
| Runtime log retention | 1 hour | 1 day | [Limits page](https://vercel.com/docs/limits) |
| Deployments per day | 100 | 6,000 | [Limits page](https://vercel.com/docs/limits) |
| Build time per deployment | 45 min | 45 min | [Limits page](https://vercel.com/docs/limits) |
| Commercial use | **Prohibited** | Allowed | [Terms §4](https://vercel.com/legal/terms), [Fair use](https://vercel.com/docs/limits/fair-use-guidelines#commercial-usage) |

Note on the public pricing page: vercel.com/pricing (read 2026-10-20) presents Pro as including "10M Edge Requests" and "1 TB" transfer, while the docs' Flat Rate CDN page says the included Pro tier is "1M CDN requests and 1 TB". The docs are more specific and more recently stamped (2026-09-14); treat 1M as the safe assumption and confirm in the billing dashboard.

---

## Recommended architecture sketch (Go on Vercel)

Principle: **render at build time wherever possible; render at request time only behind an edge cache; never serve uncached HTML to a crawler or a far-away visitor.**

### 1. Routing

Two supported Go models, both on the same page ([Go runtime](https://vercel.com/docs/functions/runtimes/go), last_updated 2026-08-11):

- **Go Framework Preset (recommended for one Go binary serving all routes).** Vercel "detects a root `go.mod` file and one of these entrypoints: `main.go`, `cmd/api/main.go`, or `cmd/server/main.go`. Your server must listen on the `PORT` environment variable." Requires `"framework": "go"` in `vercel.json`. Works with "standard `net/http` servers and frameworks such as `chi` or `gin`". The whole server is one function, so there is one cold start to manage and no 12-function Hobby cap to worry about.
- **`api/*.go` file functions.** "Each file that exports an `http.HandlerFunc` becomes a separate Vercel Function." To route everything to one handler use a catch-all rewrite, e.g. `{"rewrites":[{"source":"/(.*)","destination":"/api/index"}]}` (rewrite syntax from [vercel.json rewrites](https://vercel.com/docs/project-configuration/vercel-json#rewrites); note "precedence is given to the filesystem prior to rewrites being applied", so static files in the output directory win over the rewrite, which is exactly what we want).

Minimal `vercel.json` for the framework preset:

```json
{
  "$schema": "https://openapi.vercel.sh/vercel.json",
  "framework": "go",
  "regions": ["iad1"],
  "build": { "env": { "GO_BUILD_FLAGS": "-ldflags '-s -w'" } }
}
```

Region choice: `iad1` is the default and the cheapest compute region ($0.128/hr Active CPU). If the content DB lives in Europe, `fra1`/`dub1`/`lhr1` are available; `bom1` (Mumbai) exists as a compute region too. On Pro you can list several regions and "Vercel selects the region geographically closest to the request" ([vercel.json functionFailoverRegions section](https://vercel.com/docs/project-configuration/vercel-json#functionfailoverregions)). For a read-mostly content site, multi-region compute matters far less than cache hit ratio.

### 2. Caching headers (set in the Go handler)

Vercel's default for function responses is `cache-control: public, max-age=0, must-revalidate`, i.e. no CDN caching ([Cache-Control headers](https://vercel.com/docs/caching/cache-control-headers)). To cache, the response must contain `s-maxage=N`, optionally with `stale-while-revalidate=Z` and `stale-if-error=Z` ([CDN cache: using Vercel Functions](https://vercel.com/docs/caching/cdn-cache#using-vercel-functions)).

Vercel's own recommendation for "Semi-static (product pages, blogs, marketing)" is `max-age=120, s-maxage=86400`, and for "Server-rendered, same for all visitors" `max-age=0, s-maxage=86400` ([recommended settings](https://vercel.com/docs/caching/cache-control-headers#recommended-settings)).

Suggested policy for hastcode.com:

```go
// Article page: cache a day at the CDN, serve stale for a week while revalidating,
// serve stale for a day if the origin errors. Browser: 60s.
w.Header().Set("Cache-Control", "public, max-age=60")
w.Header().Set("Vercel-CDN-Cache-Control", "public, s-maxage=86400, stale-while-revalidate=604800, stale-if-error=86400")
```

Why `Vercel-CDN-Cache-Control`: it "has top priority", "controls caching behavior only within Vercel's Cache", and is consumed by the proxy and "not returned to the browser or forwarded to other CDNs" ([Behavior](https://vercel.com/docs/caching/cache-control-headers#behavior)). That keeps browser caching short (so a republish shows up quickly) while the edge holds the page for a day. If you ever put Cloudflare in front, add `CDN-Cache-Control` for it.

Cacheability rules the handler must respect ([cacheable response criteria](https://vercel.com/docs/caching/cdn-cache#cacheable-response-criteria)): GET/HEAD only; status 200/404/410/301/302/307/308; body ≤10 MB; **no `Set-Cookie`**; no `private`/`no-cache`/`no-store`; no `Vary: *`; and no `Vary: Cookie` ("Vercel's CDN doesn't cache a response whose `Vary` names ... `Cookie`"). So: no session cookies on article pages, no per-user rendering, analytics via client-side JS only.

Invalidation on publish: `stale-while-revalidate` means stale pages are served while the function regenerates in the background; for an immediate refresh use the CDN purge API/dashboard ([Purging CDN cache](https://vercel.com/docs/caching/cdn-cache/purge)) or a `Pragma: no-cache` request to force synchronous revalidation (`x-vercel-cache: REVALIDATED`) ([Cache-Control headers](https://vercel.com/docs/caching/cache-control-headers#stale-while-revalidate)). Verify with the `x-vercel-cache` response header: `HIT`, `MISS`, `STALE`, `PRERENDER`, `REVALIDATED`, `BYPASS` ([Response headers](https://vercel.com/docs/headers/response-headers#x-vercel-cache), [Cache status and reasons](https://vercel.com/docs/caching/cache-status)).

Important constraint: "When you specify `Cache-Control` headers in a function, responses will be cached in the region the function was requested from" and "Vercel's CDN Cache is segmented by region" ([CDN cache](https://vercel.com/docs/caching/cdn-cache#cache-control-options), [Limits](https://vercel.com/docs/caching/cdn-cache#limits)). Also "cache times are best-effort and not guaranteed ... If your asset is rarely requested (e.g. once a day), it may be evicted from the regional cache." A page's first visitor from each of the 19 regions pays a MISS.

### 3. Build

Preferred: a two-output build.

1. **Static pass (Go SSG).** A Go program renders every known article/listing/tag page to HTML at build time and writes it to `.vercel/output/static/...` (plus `sitemap.xml`, `robots.txt`, `feed.xml`, hashed CSS). Static files "are served with the Vercel Edge CDN", made available at `/` with names unmodified ([Build Output API: static files](https://vercel.com/docs/build-output-api/primitives#static-files)), and "automatically cached on Vercel's global network for the lifetime of the deployment after the first request" ([CDN cache: static files](https://vercel.com/docs/caching/cdn-cache#static-files-caching)). No function, no region, no cold start, no invocation cost. Use `overrides` in `config.json` to serve `blog/foo.html` at `/blog/foo` (clean URLs) ([configuration: overrides](https://vercel.com/docs/build-output-api/configuration#overrides), [features: cleanUrls](https://vercel.com/docs/build-output-api/features#cleanurls)).
2. **Dynamic pass (Go SSR behind edge cache).** The same Go binary (framework preset) handles whatever is not pre-rendered: search, preview/draft pages, very long tail, 404s. It sets the headers above.

Simplest variant if you don't want to hand-write Build Output API files: use the Go framework preset only and set `Vercel-CDN-Cache-Control` on everything. It works, but every page's first hit per region is a MISS to `iad1`, and a new deployment resets the cache ("Vercel scopes cached responses to the deployment that produced them" — [Cache status: Cold](https://vercel.com/docs/caching/cache-status#cold)).

Middle variant: Build Output API **Prerender Functions** (= ISR for any framework) with a build-time `fallback` HTML. See section 4.

---

## Alternatives comparison

| Platform | Go support | Regions near US/EU/India | HTML edge caching | Entry price | Primary source |
|---|---|---|---|---|---|
| **Vercel** (functions) | Official runtime, page says Beta; not in Fluid runtime list; `net/http` server or `api/*.go` | 19 compute regions incl. iad1, fra1, lhr1, dub1, cdg1, arn1, bom1; functions single-region by default | Yes via `s-maxage`/`Vercel-CDN-Cache-Control`; per-region cache; 126 PoPs | Pro $20/mo (+ usage) | [Runtimes](https://vercel.com/docs/functions/runtimes), [Regions](https://vercel.com/docs/regions) |
| **Cloudflare Workers** | Not first-class (JS/TS/Python/Rust are); Go only via WebAssembly; WASI "experimental", no threads, larger binaries start slower | Global, no region choice needed | HTML not cached by default; Cache Rules needed; static assets free | Free 100k req/day, 10ms CPU; Paid $5/mo, 10M req | [Languages](https://developers.cloudflare.com/workers/languages/), [Wasm](https://developers.cloudflare.com/workers/runtime-apis/webassembly/), [Pricing](https://developers.cloudflare.com/workers/platform/pricing/) |
| **Fly.io** | First-class: generated Dockerfile, Go binary in a Firecracker VM | iad/ord/sjc/ewr/dfw/lax; ams/fra/lhr/cdg/arn; sin/nrt/syd; **no Mumbai in the docs region table** | Not a CDN; you add Cloudflare or cache in-app | No free tier; shared-cpu-1x 256MB $2.19/mo; egress $0.02/GB NA/EU, $0.12/GB India | [Go guide](https://docs.fly.io/languages-and-frameworks/golang), [Regions](https://docs.fly.io/reference/regions), [Pricing](https://docs.fly.io/about/pricing) |
| **Railway** | Nixpacks/Railpack auto-build from `go.mod`; `$PORT` | us-west2, us-east4, europe-west4 (Amsterdam), asia-southeast1 (Singapore); multi-region replicas route "to the nearest region" | Not a CDN | Hobby $5/mo ($5 usage incl.), Pro $20/mo; vCPU $20/mo, RAM $10/GB/mo, egress $0.05/GB | [Plans](https://docs.railway.com/reference/pricing/plans), [Regions](https://docs.railway.com/reference/deployment-regions), [Scaling](https://docs.railway.com/reference/scaling) |
| **Render** | First-class: `go build -tags netgo -ldflags '-s -w' -o app`, `./app`, PORT 10000 | Oregon, Ohio, Virginia, Frankfurt, Singapore; one region per service, cannot move | "Edge caching for static assets"; dynamic responses not detailed | Free tier spins down after 15 min idle, ~1 min spin-up, 750 hrs/mo; paid instance prices not captured | [Go quickstart](https://render.com/docs/deploy-go-nethttp), [Regions](https://render.com/docs/regions), [Free](https://render.com/docs/free) |
| **VPS + Cloudflare CDN** | Anything; single Go binary | Wherever the VPS is; Cloudflare's network in front | Cloudflare "does not cache HTML or JSON by default"; Cache Rules + origin `s-maxage`; note Cloudflare says "Do not use s-maxage with stale-while-revalidate" | VPS cost + Cloudflare Free | [Default cache behavior](https://developers.cloudflare.com/cache/concepts/default-cache-behavior/), [Cache-Control](https://developers.cloudflare.com/cache/concepts/cache-control/) |

---

## 1. Vercel Go runtime: how it works, versions, support status

Status. The [Runtimes](https://vercel.com/docs/functions/runtimes) page (last_updated 2026-08-12) lists Go under **"Official runtimes"** ("The Go runtime takes in a Go program that defines a singular HTTP handler and outputs it as a Vercel Function"). Community runtimes listed there are Bash, Deno and PHP only. However the [Go runtime page](https://vercel.com/docs/functions/runtimes/go) (last_updated 2026-08-11) opens with "The Go runtime is available in Beta on all plans". So: official, but Beta.

Two deployment models (both on the Go runtime page):

- **Go Framework Preset**: root `go.mod` + `main.go` | `cmd/api/main.go` | `cmd/server/main.go`; "Your server must listen on the `PORT` environment variable"; "Running a Go server requires the `framework` preset to be set to `go`". The builder source confirms: when `framework` is `go` (or in a Services project) the builder calls `buildStandaloneServer` instead of the per-handler wrapper ([packages/go/src/index.ts](https://raw.githubusercontent.com/vercel/vercel/main/packages/go/src/index.ts)). The [@vercel/go CHANGELOG](https://github.com/vercel/vercel/blob/main/packages/go/CHANGELOG.md) records this as "An experimental Go runtime framework preset ... for standalone HTTP servers" added in 3.4.0; the package is at 10.0.1 as of 2026-10-20.
- **`api/*.go` functions**: "Each file that exports an `http.HandlerFunc` becomes a separate Vercel Function"; "The function name can be any exported Go identifier as long as it matches the `http.HandlerFunc` signature". The builder renames your function, generates a `main__vc__go__.go` wrapper, compiles with `GOOS=linux GOARCH=amd64`, and ships a `bootstrap` binary as the Lambda handler using a "provided" runtime obtained from `getProvidedRuntime()` ([index.ts](https://raw.githubusercontent.com/vercel/vercel/main/packages/go/src/index.ts)). Vercel's bundle-size doc confirms the AWS Lambda substrate: limits "are enforced by AWS" ([Limits](https://vercel.com/docs/functions/limitations#bundle-size-limits)).

Go versions. "Vercel reads the Go version from the `go` directive in `go.mod`. If a `toolchain` directive is also present, Vercel uses that version instead. `go.mod` must be at the project root. If `go.mod` does not declare a version, Vercel uses the latest supported Go version." ([Go runtime](https://vercel.com/docs/functions/runtimes/go#go-version)). The builder's `versionMap` (read 2026-10-20 from [packages/go/src/go-helpers.ts](https://raw.githubusercontent.com/vercel/vercel/main/packages/go/src/go-helpers.ts)) maps 1.13 → 1.27, newest-first, with 1.27 → 1.27.1 and a minimum of Go 1.13; a `go 1.N` newer than the map resolves to `1.N.0`. Toolchains are downloaded and cached on first use.

Dependencies and build. `go.mod` + committed `go.sum`; vendoring supported ("If your project contains a `vendor` directory with a `modules.txt` file, Vercel passes `-mod=vendor`"); `GO_BUILD_FLAGS` env var (default `-ldflags '-s -w'`); `buildCommand` override runs "with the selected Go toolchain on `PATH`, and `GOOS=linux` already set", and the binary is located in the project root or `bin/`, else via `$VERCEL_OUTPUT_FILE`. Private modules via `GIT_CREDENTIALS`. ([Go runtime](https://vercel.com/docs/functions/runtimes/go))

Container option. Vercel also runs OCI images from a `Dockerfile.vercel` ("Container Runtime is available in Beta on all plans"), auto-adding "a rewrite rule to route all traffic to the container image", listening on port 80 or `$PORT`, with "The same limits and Active CPU pricing model of Vercel Functions"; instances "not receiving any traffic for 5 minutes in production ... will automatically scale down" ([Container runtime](https://vercel.com/docs/functions/runtimes/container), last_updated 2026-07-07). This is a fallback if the Go preset proves limiting, at the cost of slower cold starts (container pull).

## 2. Function limits for Go, cold starts, regions, Fluid

Numbers are in the table above. Go-specific caveats:

- **Fluid compute.** "Fluid compute is available for the following runtimes: Node.js, Python, Edge, Bun, Rust" — Go is absent ([Fluid compute](https://vercel.com/docs/fluid-compute#available-runtime-support), last_updated 2026-08-24). "Optimized concurrency ... is available when using Node.js or Python runtimes." Bytecode caching is Node-only. So a Go function should be assumed to run one request per instance (classic Lambda model) with no documented pre-warming. Counter-signal: the [Services pricing page](https://vercel.com/docs/services/pricing) says "Fluid Compute applies, so concurrent requests can share an instance" for every service, and the Go page says Go servers can be deployed as a Service. This is unresolved; see Open questions.
- **Streaming.** "The Node.js runtime supports streaming by default. Streaming is also supported when using the Python runtime." ([Runtimes: streaming](https://vercel.com/docs/functions/runtimes#streaming)). Go is not mentioned anywhere on the [Streaming](https://vercel.com/docs/functions/streaming-functions) page. Assume buffered responses for Go.
- **Max duration.** Go gets the plan's base limits (300s Hobby; 300s default / 800s max Pro). The 1800s extended beta and `"max"` are for Node/Bun/Python ([Limits](https://vercel.com/docs/functions/limitations#max-duration)). Irrelevant for page rendering.
- **Cold starts.** Vercel publishes no Go cold-start figures. What is documented: archived functions (2 weeks idle in prod) add "at least 1 second"; Fluid's cold-start mitigations (instance reuse, bytecode caching, pre-warmed instances) are described for Node ([What is compute](https://vercel.com/docs/fundamentals/what-is-compute#cold-starts), [cold-start KB](https://vercel.com/kb/guide/improve-function-cold-start-performance-on-vercel)). A statically linked Go binary with `-s -w` is small, so the Lambda init itself should be short, but this is inference, not a Vercel statement.
- **Region.** "By default, Vercel Functions execute in Washington, D.C., USA (`iad1`) for all new projects"; Hobby: "Single region"; Pro: "5 regions"; Enterprise: all; "Deploying to more regions than your plan allows causes the deployment to fail before the build step" ([Region config](https://vercel.com/docs/functions/configuring-functions/region)). The Fluid page's table says Pro "Up to 3" multi-region functions; the region page is the more specific document. Automatic cross-region failover is Enterprise-only; AZ-level failover is default.
- **Hobby function cap.** "For Hobby, this approach [one file = one function] is limited to 12 Vercel Functions per deployment" ([Runtimes](https://vercel.com/docs/functions/runtimes#functions-created-per-deployment)). The framework preset avoids this because the whole server is one function.

## 3. Caching SSR HTML at the edge

Mechanics (all from [CDN cache](https://vercel.com/docs/caching/cdn-cache), last_updated 2026-09-14, and [Cache-Control headers](https://vercel.com/docs/caching/cache-control-headers), last_updated 2026-09-14):

- Trigger: any of `s-maxage=N`, `s-maxage=N, stale-while-revalidate=Z`, `s-maxage=N, stale-while-revalidate=Z, stale-if-error=Z`. `proxy-revalidate` unsupported.
- `s-maxage`: "sets the number of seconds a response is considered 'fresh' by the CDN. After this period ends, Vercel's CDN will serve the 'stale' response from the edge until the response is asynchronously revalidated". Range 1 s to 31,536,000 s.
- `stale-while-revalidate`: "serve content from the Vercel CDN cache while simultaneously updating the cache in the background". The first request is synchronous; later ones are served from cache and revalidated asynchronously when stale. Vercel consumes the directive and strips it before the browser sees it.
- `stale-if-error`: CDN "serves a stale response when an error is encountered instead of returning the error to the client".
- Precedence: `Vercel-CDN-Cache-Control` > `CDN-Cache-Control` > `Cache-Control`. If only `Cache-Control` is set, Vercel strips `s-maxage` before the browser; the browser receives `cache-control: public, max-age=0, must-revalidate` ([Response headers](https://vercel.com/docs/headers/response-headers#cache-control)).
- Keys: `Accept` and `Accept-Encoding` are already in the cache key; add `Vary` only for small-cardinality headers such as `X-Vercel-IP-Country` or `Accept-Language`. `Vary: Cookie` makes the response uncacheable.
- Scope: per region and best-effort ("segmented by region"; rarely requested assets "may be evicted from the regional cache"); scoped to the deployment ("after a new deployment" is a listed cold-miss cause).

Global reach. Vercel runs "126 PoPs" that "terminate TCP and route requests over a private network to the nearest Vercel region" and "19 compute-capable regions"; "popular content is more likely to be available in each region's cache" ([Regions](https://vercel.com/docs/regions)). An Indian visitor is served from the `bom1` region's cache on a HIT; on a MISS the request crosses to `iad1` (or whichever region you pin).

ISR route for a non-Next framework. Vercel contrasts Cache-Control caching with ISR: "With `Cache-Control` headers alone, Vercel doesn't know a path is cacheable until it receives the response, so these features aren't available" — the features being request collapsing, durable storage that survives deployments, "Globally consistent purging ... within 300ms", instant rollbacks ([ISR](https://vercel.com/docs/incremental-static-regeneration#caching-on-vercel), last_updated 2026-08-28). ISR is available to any framework through "Build Output API: Define Prerender Functions". ISR durable storage lives in one region (your default function region) with CDN copies in every region; "CDN cache reads and writes are free", ISR reads/writes bill at $0.40–0.64 / $4.00–6.40 per million 8 KB units ([ISR pricing](https://vercel.com/docs/incremental-static-regeneration/limits-and-pricing)).

## 4. Static vs SSR vs ISR via the Build Output API v3

The [Build Output API](https://vercel.com/docs/build-output-api) is "a file-system-based specification for a directory structure that can produce a Vercel deployment"; "If you are not using a framework ... you can create the `.vercel/output` directory and populate it according to this specification yourself." Examples live in [vercel/examples/build-output-api](https://github.com/vercel/examples/tree/main/build-output-api).

- **Static**: `.vercel/output/static/` — files served at `/`, unmodified, by the Edge CDN; `config.json` `overrides` to rename `blog.html` → `/blog`; `_vercel/immutable/` + `immutable.json` for content-hashed assets shared across deployments ([Primitives: static files](https://vercel.com/docs/build-output-api/primitives#static-files)).
- **Functions**: `.vercel/output/functions/<path>.func/.vc-config.json` with `runtime`, `handler`, `memory`, `maxDuration`, `regions`, `supportsResponseStreaming`, `architecture` (`x86_64` default). A `.func` directory "may be a symlink to another `.func` directory" so many paths can share one Go binary ([Primitives: functions](https://vercel.com/docs/build-output-api/primitives#functions)).
- **Prerender functions (ISR)**: "A Prerender asset is a Vercel Function that will be cached by the Vercel CDN in the same way as a static file. This concept is also known as Incremental Static Regeneration." Config `<name>.prerender-config.json`: `expiration` (seconds, or `false` = never), `group`, `bypassToken`, `fallback` (build-time HTML served "while there is not yet a cached version"; the function is then "invoked 'out-of-band' to re-generate"), `allowQuery`, `passQuery`, `initialHeaders`, `initialStatus` ([Primitives: prerender](https://vercel.com/docs/build-output-api/primitives#prerender-functions)). On-demand revalidation: `GET`/`HEAD` with header `x-prerender-revalidate: <bypassToken>` ([Features: on-demand ISR](https://vercel.com/docs/build-output-api/features#on-demand-incremental-static-regeneration-isr)).
- **Routes**: `config.json` `{"version":3,"routes":[...]}` using the low-level `src`/`dest`/`headers`/`continue`/`handle` syntax; `@vercel/routing-utils` `getTransformedRoutes()` converts `rewrites`/`headers`/`cleanUrls` into it ([Configuration](https://vercel.com/docs/build-output-api/configuration), [Features](https://vercel.com/docs/build-output-api/features)).
- **Known limitation**: "when building locally, your build tools will compile native dependencies targeting your machine's architecture" — build the Go binary with `GOOS=linux GOARCH=amd64` ([Build Output API](https://vercel.com/docs/build-output-api#known-limitations)).

Caveat for Go: the public docs show only Node.js and Edge `.vc-config.json` examples; the `runtime` string the Go builder emits is computed by `getProvidedRuntime()` in `@vercel/build-utils` and is not documented on the Build Output API pages. Emitting Go `.func` directories by hand is therefore undocumented territory (see Open questions). The reliable combination is: hand-emit **static** files via the Build Output API (fully documented) and let the **Go framework preset** build the function.

Practical decision for hastcode.com:

| Page type | Render | Why |
|---|---|---|
| Articles, listings, tags, home, sitemap, feeds | Static at build (Go SSG → `.vercel/output/static`) | Served from every region, no invocation, no cold start, survives nothing to survive |
| Newly published article before next build | SSR via Go function + `Vercel-CDN-Cache-Control` | Immediately live; first hit per region is a MISS |
| Search, previews, 404 | SSR via Go function | Dynamic by nature |
| Optional: long tail with on-demand revalidate | Prerender Functions (ISR) with `fallback` | Durable cache, request collapsing, 300 ms global purge |

## 5. Pricing at 1–3M pageviews/month

Hobby (free): 100 GB Fast Data Transfer, 1M CDN requests, 1M invocations, 4 CPU-hrs, 360 GB-hrs, 10 GB Fast Origin Transfer; on exceeding "you will have to wait until 30 days have passed before you can use the feature again" ([Hobby plan](https://vercel.com/docs/plans/hobby), last_updated 2026-09-14). "Hobby teams are restricted to non-commercial personal use only ... Examples ... The inclusion of advertisements, including ... Google AdSense" ([Fair use: commercial usage](https://vercel.com/docs/limits/fair-use-guidelines#commercial-usage)); Terms of Service §4 "Hobby Plan": "You shall only use the Services under a Hobby plan for your personal or non-commercial use" ([Terms](https://vercel.com/legal/terms)). A tech-news site with ads or sponsorships is commercial.

Pro: "$20/month Pro platform fee — 1 deploying team seat included — $20/month in usage credit"; "Pro includes the lowest capacity tier at no extra cost, with a capacity of 1 million CDN requests and 1 TB of data transfer each month" ([Pro plan](https://vercel.com/docs/plans/pro-plan), last_updated 2026-09-15). Flat Rate CDN tiers: included 1M req/1 TB; $20/mo 10M/50 TB; $100/mo 50M/50 TB; $300/mo 150M/50 TB; above that, on-demand ([Flat Rate CDN](https://vercel.com/docs/pricing/flat-rate-cdn)). On-demand: Fast Data Transfer $0.15–0.35/GB, CDN Requests $2.00–3.20/M, Fast Origin Transfer $0.06–0.43/GB, invocations $0.60/M, Active CPU $0.128–0.221/hr, memory $0.0106–0.0183/GB-hr ([Regional pricing](https://vercel.com/docs/pricing/regional-pricing), [Fluid pricing](https://vercel.com/docs/functions/usage-and-pricing)).

Rough sizing (assumptions marked; not Vercel numbers):

- CDN requests: every HTML, CSS, image, font and script fetch counts ("Static assets and functions all incur CDN Requests", [CDN pricing](https://vercel.com/docs/manage-cdn-usage#cdn-requests)). Assume ~8 requests/pageview after aggressive bundling → 8–24M requests/month → the $20/mo (10M) or $100/mo (50M) Flat Rate tier. This is the main line item.
- Bandwidth: assume ~250 KB/pageview → 250–750 GB/month, inside the 1 TB included tier.
- Invocations: with static pre-rendering and a 95%+ CDN HIT ratio on SSR pages, assume 50–200k Go invocations/month → under $0.15 on demand.
- Compute: assume 20 ms Active CPU × 200k = ~1.1 CPU-hr ≈ $0.15; memory negligible at 2 GB × sub-second.

Estimated total: $20 (Pro) + $20–100 (Flat Rate CDN tier) ≈ **$40–120/month**, almost all of it CDN. Compute is a rounding error if caching is done right; without caching, each pageview becomes an invocation plus Fast Origin Transfer, and India/EU TTFB suffers.

## 6. Alternatives (one paragraph each)

**Cloudflare Workers / Pages.** Workers' first-class languages are "JavaScript, TypeScript, Python Workers, and Rust"; Go is only listed among languages that can compile to WebAssembly ([Languages](https://developers.cloudflare.com/workers/languages/)). The Wasm page warns that "WASI support is experimental", threads are unavailable, Wasm Workers "are usually larger than equivalent JavaScript Workers" and "larger Workers may take longer to start" ([WebAssembly](https://developers.cloudflare.com/workers/runtime-apis/webassembly/)). Running `net/http`-style Go on Workers means TinyGo or a `GOOS=js` build plus a JS shim; Cloudflare publishes no Go guide, so this is community-supported only. Pricing: Free 100k requests/day with 10 ms CPU; Paid $5/mo minimum with 10M requests and 30M CPU-ms included, $0.30 per extra million; "Requests to static assets are free and unlimited" ([Pricing](https://developers.cloudflare.com/workers/platform/pricing/), [Static assets](https://developers.cloudflare.com/workers/static-assets/)). Best fit: Go SSG at build time → static assets on Workers/Pages (truly global, free egress), with any dynamic bits in JS/TS rather than Go.

**Fly.io.** First-class Go: `fly launch` detects a Go app, generates a Dockerfile (Debian base with CGO, Alpine optional), writes `fly.toml` with `primary_region`, and deploys to Fly Machines; default VM `shared-cpu-1x` 1 GB ([Go guide](https://docs.fly.io/languages-and-frameworks/golang)). Multi-region is explicit: `fly scale count 3 --region ams,iad` places Machines per region; a shared IPv4 and "unlimited Anycast IPv6 addresses for global load balancing" come free ([Scale count](https://docs.fly.io/launch/scale-count), [Pricing](https://docs.fly.io/about/pricing)). Regions include iad/ord/sjc/ewr/dfw/lax and ams/fra/lhr/cdg/arn plus sin/nrt/syd; **Mumbai (`bom`) is not in the docs' region table** as of 2026-10-20 ([Regions](https://docs.fly.io/reference/regions)). No free tier (trial of 2 hours of Machine runtime or 7 days); `shared-cpu-1x` 256 MB is $2.19/mo in iad; egress $0.02/GB NA/EU but **$0.12/GB to Africa and India** ([Pricing](https://docs.fly.io/about/pricing)). Fly is not a CDN, so you would still put Cloudflare in front for India and static assets.

**Railway / Render.** Railway auto-builds Go from `go.mod` with Nixpacks/Railpack, serves on `$PORT`, and offers regions us-west2 (California), us-east4 (Virginia), europe-west4 (Amsterdam) and asia-southeast1 (Singapore); multi-region replicas are supported and "public traffic goes to the nearest region" ([Regions](https://docs.railway.com/reference/deployment-regions), [Scaling](https://docs.railway.com/reference/scaling)); Hobby $5/mo with $5 usage, Pro $20/mo with $20; vCPU $20/mo, RAM $10/GB/mo, egress $0.05/GB ([Plans](https://docs.railway.com/reference/pricing/plans)). Render's Go quickstart is `go build -tags netgo -ldflags '-s -w' -o app` / `./app` binding `0.0.0.0:10000` ([Go quickstart](https://render.com/docs/deploy-go-nethttp), [Web services](https://render.com/docs/web-services)); regions are Oregon, Ohio, Virginia, Frankfurt, Singapore, one region per service and "Render doesn't currently support changing the region for an existing service" ([Regions](https://render.com/docs/regions)); the free tier "spins down after 15 minutes" and spin-up "takes about one minute", with 750 free instance hours/month and no edge caching ([Free](https://render.com/docs/free)). Neither has an India region; both are single-origin PaaS and need a CDN in front for a global content site.

**Plain VPS + Go binary + Cloudflare CDN.** The classic choice for a Go content site: one `net/http` binary, systemd, Cloudflare proxying the domain. Cloudflare's CDN "does not cache HTML or JSON by default"; you create a Cache Rule to cache everything, then control TTL with origin headers ([Default cache behavior](https://developers.cloudflare.com/cache/concepts/default-cache-behavior/)). Cloudflare honours `s-maxage`, `stale-while-revalidate` and `stale-if-error`, and `Cloudflare-CDN-Cache-Control` targets Cloudflare alone, but note "Do not use s-maxage with stale-while-revalidate" on Cloudflare (SWR is disabled by `s-maxage`), and `stale-if-error` only fires on origin 5xx ([Cache-Control](https://developers.cloudflare.com/cache/concepts/cache-control/)). Cheapest and most controllable; you own ops, TLS renewals (Cloudflare handles the edge), and single-origin latency on misses.

## 7. Go SSR templating: html/template vs templ vs gomponents

- **`html/template`** (stdlib). "Implements data-driven templates for generating HTML output safe against code injection"; same API as `text/template`; contextual auto-escaping decided at parse time (HTML, attribute, URL, CSS, JS contexts); `Execute(w io.Writer, data)`; "A template may be executed safely in parallel" ([pkg.go.dev/html/template](https://pkg.go.dev/html/template)). Zero dependencies and no codegen. Render into a `bytes.Buffer` first so you can set status, `Content-Length`, and an `ETag` before writing — that is what makes the response cleanly cacheable.
- **templ** (a-h/templ). Components in `.templ` files are "compiled into performant Go code" by `templ generate`; `templ.Handler(component, opts...)` is a ready `http.Handler`, components expose `Render(ctx, w io.Writer)` ([templ.guide](https://templ.guide/), [HTTP server](https://templ.guide/server-side-rendering/creating-an-http-server-with-templ/)). Default behaviour: "`templ.Handler` renders the template into a buffer before writing it to the response" so error status codes still work; `templ.WithStreaming()` plus `@templ.Flush()` streams, but then "response codes and headers can no longer be changed" ([Streaming](https://templ.guide/server-side-rendering/streaming/)). Static generation is a one-liner: render to an `os.File` ([Static HTML](https://templ.guide/static-rendering/generating-static-html-files-with-templ/)). Best type-safety/ergonomics; adds a codegen step to the Vercel build (`go generate` or a committed `_templ.go`).
- **gomponents** (maragudk). "Pure Go" HTML components, no code generation, no dependencies; "renders directly to an io.Writer"; `Text`/`Textf` escape, `Raw`/`Rawf` do not; README benchmark: a realistic full page ≈ 63,560 ns and 2,858 allocations on an M4 ([gomponents](https://github.com/maragudk/gomponents), [gomponents.com](https://www.gomponents.com/)). Simplest build pipeline; markup-as-Go is a taste question.

Caching and streaming interplay on Vercel: because Go streaming is undocumented on Vercel and a cached edge response must have a final status and headers, **prefer buffered rendering** (`html/template` into a buffer, templ's default handler, or gomponents into a buffer). That also lets the Go handler compute `ETag` and `Content-Length`, which Vercel recommends to cut Fast Origin Transfer ("Ensure your Function supports `If-Modified-Since` or `ETag`", [CDN pricing](https://vercel.com/docs/manage-cdn-usage#functions)). Streaming only helps TTFB on cache misses, and those should be rare.

## 8. Core Web Vitals implications

Thresholds (Google): TTFB "Good ... 0.8 seconds or less", "poor values are greater than 1.8 seconds", measured at the 75th percentile ([web.dev TTFB](https://web.dev/articles/ttfb)). TTFB is not itself a Core Web Vital; the CWVs are LCP ≤ 2.5 s, INP ≤ 200 ms, CLS ≤ 0.1 at p75, and TTFB is listed as a supporting metric "for diagnosing issues with LCP" ([web.dev vitals](https://web.dev/articles/vitals)). Since LCP ≥ TTFB + HTML parse + hero asset, an 800 ms TTFB leaves little room for a 2.5 s LCP on mobile in India.

How Vercel's topology maps onto that for a US/EU/India audience (latency figures below are reasoning, not Vercel-published numbers):

| Request path | Who | Expected TTFB contributors | Verdict |
|---|---|---|---|
| Static file or `x-vercel-cache: HIT` in the visitor's region (e.g. `bom1`, `fra1`) | Any | PoP TCP/TLS termination + regional cache read | Well inside 0.8 s |
| `STALE` (SWR window) | Any | Same as HIT; revalidation happens in the background | Inside 0.8 s |
| `MISS` from India/EU to a warm Go instance in `iad1` | India/EU | PoP → private network to `iad1` (intercontinental RTT) + Go render + DB round trips near `iad1` + response back | Typically still OK if render is fast and data is co-located with the function, but it is the p75 risk if MISSes are >25% of a region's traffic |
| `MISS` + cold start (new deployment, idle instance, or archived function) | India/EU | Above + Lambda init (Vercel: archived functions add "at least 1 second") | Likely blows the 0.8 s "good" budget for that request |
| `REVALIDATED` (synchronous, after a purge) | Any | Full origin round trip | Rare, operator-triggered |

Design consequences:

1. Pre-render everything you can at build time so most URLs never depend on `iad1`.
2. For SSR pages, use long `s-maxage` + very long `stale-while-revalidate` so that once a region has a copy it keeps serving (`STALE` is as fast as `HIT`), and expect a cold MISS per region per deployment.
3. Keep data sources in the function region ("Functions should be executed in the same region as your database", [Regions](https://vercel.com/docs/regions#compute-defaults)); a Go render that fans out to a database in another continent doubles the miss penalty.
4. Measure with real-user data segmented by country (Vercel Speed Insights or your own RUM) and by `x-vercel-cache` status; crawlers from Googlebot are mostly US-based and will usually see `iad1`-local latencies, so a US-only lab test will overstate the site's global health.
5. Keep the HTML small and fonts/CSS inline or hashed-immutable so a good TTFB turns into a good LCP on Indian mobile networks.

---

## Open questions

1. **Does Fluid compute apply to Go?** The Fluid page omits Go from the supported runtimes; the Services pricing page says Fluid applies to every service; the Limits page describes all limits "with Fluid compute". Needs a Vercel support answer or an empirical test (observe instance reuse and Active CPU vs duration billing on a Go deployment).
2. **Streaming for Go functions.** Not documented. Test whether `http.Flusher` works through the Go runtime and whether `supportsResponseStreaming` is set by the builder.
3. **Go runtime identifier for hand-written `.vc-config.json`.** The builder gets it from `getProvidedRuntime()` in `@vercel/build-utils`; not published on the Build Output API pages. Needed only if you want Prerender Functions (ISR) backed by Go without going through the framework preset. Alternative: inspect `.vercel/output` after `vercel build` on the Go preset and reuse what it emits.
4. **Pro multi-region count:** 5 (region page) vs 3 (Fluid page). Confirm in dashboard.
5. **Pro included CDN requests:** docs say 1M (Flat Rate CDN included tier); vercel.com/pricing marketing page shows 10M. Confirm on the billing page before sizing.
6. **Go cold-start latency on Vercel** has no published figure; measure p50/p95 of `MISS` responses from `bom1`/`fra1` PoPs after a fresh deployment.
7. **Does the Go framework preset get instance pre-warming on Pro production deployments?** The cold-start KB attributes pre-warming to Fluid/Node; unclear for Go.
8. **Cache persistence across deployments** for `Cache-Control`-cached pages is per deployment; confirm whether Prerender Functions' "fallback" plus durable ISR storage is worth the extra build complexity versus simply redeploying the static pass.

---

## Sources (all accessed 2026-10-20)

Vercel docs
- https://vercel.com/docs/functions/runtimes/go (last_updated 2026-08-11)
- https://vercel.com/docs/functions/runtimes (2026-08-12)
- https://vercel.com/docs/functions/configuring-functions/runtime (2026-08-12)
- https://vercel.com/docs/functions/runtimes/container (2026-07-07)
- https://vercel.com/docs/functions/limitations (2026-08-24)
- https://vercel.com/docs/limits (2026-09-16)
- https://vercel.com/docs/fluid-compute (2026-08-24)
- https://vercel.com/docs/fundamentals/what-is-compute (2026-08-11)
- https://vercel.com/docs/functions/configuring-functions/region (2026-08-11)
- https://vercel.com/docs/regions (2026-08-11)
- https://vercel.com/docs/functions/streaming-functions (2026-09-01)
- https://vercel.com/docs/functions/usage-and-pricing (2026-06-16)
- https://vercel.com/docs/caching/cdn-cache (2026-09-14)
- https://vercel.com/docs/caching/cache-control-headers (2026-09-14)
- https://vercel.com/docs/caching/cache-status (2026-09-15)
- https://vercel.com/docs/headers/response-headers (2026-08-11)
- https://vercel.com/docs/incremental-static-regeneration (2026-08-28)
- https://vercel.com/docs/incremental-static-regeneration/limits-and-pricing (2026-08-11)
- https://vercel.com/docs/build-output-api (2026-08-11)
- https://vercel.com/docs/build-output-api/primitives (2025-03-04)
- https://vercel.com/docs/build-output-api/configuration (2026-07-27)
- https://vercel.com/docs/build-output-api/features (2025-03-04)
- https://vercel.com/docs/project-configuration/vercel-json (2026-08-14)
- https://vercel.com/docs/services (2026-08-10)
- https://vercel.com/docs/services/pricing (2026-08-31)
- https://vercel.com/docs/plans/hobby (2026-09-14)
- https://vercel.com/docs/plans/pro-plan (2026-09-15)
- https://vercel.com/docs/pricing/flat-rate-cdn (2026-09-14)
- https://vercel.com/docs/pricing/regional-pricing (2026-09-14)
- https://vercel.com/docs/manage-cdn-usage (2026-09-14)
- https://vercel.com/docs/limits/fair-use-guidelines (2026-09-14)
- https://vercel.com/legal/terms (Section 4, Hobby Plan)
- https://vercel.com/pricing (marketing page; numbers differ from docs, see note)
- https://vercel.com/kb/guide/improve-function-cold-start-performance-on-vercel
- https://vercel.com/kb/guide/deploy-go-using-docker-vercel

Vercel Go runtime source (vercel/vercel, `main` branch)
- https://github.com/vercel/vercel/tree/main/packages/go
- https://github.com/vercel/vercel/blob/main/packages/go/CHANGELOG.md (latest 10.0.1; 3.4.0 added the framework preset; 3.4.4 added Go 1.26 + 1.25)
- https://raw.githubusercontent.com/vercel/vercel/main/packages/go/src/index.ts (bootstrap handler, `getProvidedRuntime()`, `GOOS=linux GOARCH=amd64`, `buildStandaloneServer`)
- https://raw.githubusercontent.com/vercel/vercel/main/packages/go/src/go-helpers.ts (`versionMap` 1.13–1.27, min 1.13)
- https://github.com/vercel/examples/tree/main/build-output-api

Cloudflare
- https://developers.cloudflare.com/workers/languages/
- https://developers.cloudflare.com/workers/runtime-apis/webassembly/
- https://developers.cloudflare.com/workers/platform/limits/
- https://developers.cloudflare.com/workers/platform/pricing/
- https://developers.cloudflare.com/workers/static-assets/
- https://developers.cloudflare.com/cache/concepts/default-cache-behavior/
- https://developers.cloudflare.com/cache/concepts/cache-control/

Fly.io
- https://docs.fly.io/languages-and-frameworks/golang
- https://docs.fly.io/reference/regions
- https://docs.fly.io/launch/scale-count
- https://docs.fly.io/about/pricing

Railway
- https://docs.railway.com/reference/pricing/plans
- https://docs.railway.com/reference/deployment-regions
- https://docs.railway.com/reference/scaling
- https://docs.railway.com/guides/gin

Render
- https://render.com/docs/deploy-go-nethttp
- https://render.com/docs/web-services
- https://render.com/docs/regions
- https://render.com/docs/free
- https://render.com/pricing (did not render; instance prices not captured)

Google / web.dev
- https://web.dev/articles/ttfb
- https://web.dev/articles/vitals

Go templating
- https://pkg.go.dev/html/template
- https://templ.guide/
- https://templ.guide/server-side-rendering/creating-an-http-server-with-templ/
- https://templ.guide/server-side-rendering/streaming/
- https://templ.guide/static-rendering/generating-static-html-files-with-templ/
- https://www.gomponents.com/
- https://github.com/maragudk/gomponents

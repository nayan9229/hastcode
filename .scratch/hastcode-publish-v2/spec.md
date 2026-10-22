# Spec: hastcode.com publish v2 — Cloudflare Production, Vercel QA, Publish by tag, images in R2

Status: ready-for-agent
Created: 2026-10-22
Glossary: ../../GLOSSARY.md (Publishing section: Target, Publish, Preview, Image Store)
ADRs: ../../docs/adr/0004-cloudflare-production-vercel-qa-publish-by-tag-images-in-r2.md (governs), ../../docs/adr/0001-replace-genuin-community-with-static-publication.md (amended: 404 on Production), ../../docs/adr/0002-static-first-go-build-on-vercel.md (build half still stands)
Research: ../../docs/research/cloudflare-r2-and-deploy.md, ../../docs/research/image-optimisation-targets.md
Predecessor: ../hastcode-v1/spec.md (v1, all tickets resolved; ticket 15 was ready-for-human)

## Problem Statement

v1 deploys to Vercel Pro from a push to `main`, with every image variant inside the site output. The Author now wants hastcode.com served from Cloudflare (free, unmetered, already holds the DNS), Vercel kept only as a place to check a release before it goes live, and the moment of going live under explicit control rather than "whatever was last pushed". Images inside the output also collide with Cloudflare's 20,000-file cap (about 1,200 Explainers), and the current encoder settings were chosen quickly in v1 rather than from measurement.

So today: no way to deploy to Cloudflare; no QA step; publishing is implicit; images bloat every deploy; hero bytes on the LCP path are larger than they need to be.

## Solution

Publishing becomes a deliberate act: tag a commit on `main` with `PUBVL-…` to put it on the QA Target (Vercel), check it, then tag the same commit `PUBCF-…` to put it on the Production Target (Cloudflare), which is the only host hastcode.com ever points at. Pushing `main` runs tests only. Branch pushes still produce Previews, now on Cloudflare preview URLs.

The build stays one Go program producing one static tree, now with a manifest per Target. Every raster image is uploaded at build time to an R2 bucket served at `images.hastcode.com` and referenced by absolute URL, so each Explainer ships as one file. The images package is retuned from the research: new width set, AVIF 4:4:4 with quality tiers, WebP dropped, a dedicated Open Graph crop, EXIF orientation honoured, sRGB enforced by warning, and `fetchpriority="high"` on the hero.

The newsletter signup is ported to a Cloudflare Worker, the only code that runs on Production. The Go function stays for the QA Target. A setup wizard and README carry the human steps.

## User Stories

### Publishing

1. As the Author, I want pushing `main` to run tests and nothing else, so that merging never changes what Readers see.
2. As the Author, I want to run one command (`publish cf` or `publish vl`) that tags the current `main` commit and pushes the tag, so that publishing is explicit and I never type a tag by hand.
3. As the Author, I want a `PUBVL-*` tag to deploy that commit to the QA Target in production mode, so that QA shows exactly the bytes Production would ship.
4. As the Author, I want a `PUBCF-*` tag to deploy that commit to the Production Target, so that going live is a single, auditable act.
5. As the Author, I want both tags on one commit to run independently, so that QA and Production never wait on each other.
6. As the Author, I want a tag that does not point at a commit on `main` to be refused with a clear message, so that I cannot publish an unmerged branch by accident.
7. As the Author, I want the publish command to print the QA URL after `vl` and remind me to tag `cf` with the same commit, so that the normal flow is obvious.
8. As the Author, I want one Production Publish at a time, in tag order, so that the IndexNow manifest diff is always against the previous live build.
9. As the Author, I want a failed build (editorial gate, broken link, size gate, missing image credentials) to deploy nothing to either Target, so that Readers never see a half-built site.
10. As the Author, I want IndexNow pinged only after a Production Publish is live, so that Bing is never told about QA URLs.
11. As the Author, I want a branch push to deploy a Preview (drafts included) to a Cloudflare preview URL and run Lighthouse CI against it, so that I review drafts where they will really be served.
12. As the Author, I want the Vercel preview workflow removed, so that there is one Preview path.

### QA Target

13. As the Author, I want the QA Target reachable at `qa.hastcode.com`, so that I check releases on a stable address.
14. As the Author, I want every QA deployment to send `X-Robots-Tag: noindex` and serve a `robots.txt` that disallows everything, so that Google never indexes a duplicate of the site.
15. As the Author, I want QA deployments protected by Vercel Authentication, so that only I can open them.
16. As the Author, I want canonical and Open Graph URLs on QA to still point at `https://hastcode.com/...`, so that nothing QA-specific leaks into page metadata.
17. As Googlebot, I want the QA host to tell me not to index it in both headers and robots.txt, so that I never rank it.

### Production Target

18. As a Reader, I want every page served from Cloudflare's edge as a static asset, so that nothing waits on a function.
19. As a Reader, I want `/about` to reach `/about/` and `/x/index.html` to reach `/x/`, so that one page has one address.
20. As a Reader, I want an unknown path to show the 404 page with status 404, so that mistakes still lead somewhere useful.
21. As a Reader, I want old community URLs (`/video/…`, `/community/…`, `/group/…`) to return 404, so that dead links fail honestly (ADR-0001 as amended).
22. As a Reader, I want `/home` to redirect to `/`, so that the old entry point still works.
23. As a Reader, I want hashed assets cached for a year and HTML revalidated, so that pages stay fresh and assets fast.
24. As a Reader, I want the newsletter form to work on Production exactly as on QA: plain POST without JavaScript, confirmation page, errors shown, so that signing up never depends on scripts.
25. As the Author, I want the Worker to run only for the signup path, so that no other request executes code.
26. As the Author, I want the Worker to validate email syntax, drop honeypot hits silently, rate-limit per IP best-effort, tag the signup with the page it came from, and never expose the Buttondown key, so that Production signups are as safe as QA's.
27. As the Author, I want the Worker's Buttondown base URL injectable, so that it can be tested against a fake.

### Image Store

28. As the Author, I want every raster (hero variants, crops, author photos, OG badge) uploaded to R2 and referenced by absolute URL on `images.hastcode.com`, so that each Explainer ships as one file and Cloudflare's file cap stops mattering.
29. As the Author, I want uploads to skip objects that already exist, so that a rebuild uploads only new images and costs nothing.
30. As the Author, I want each object uploaded with the right `Content-Type` and an immutable one-year `Cache-Control`, so that the edge and browsers cache it forever.
31. As the Author, I want `make dev` and tests to need no R2 credentials, so that local work stays offline.
32. As the Author, I want a production or QA build configured for R2 but missing credentials to fail immediately with a clear message, so that it never silently ships local image paths.
33. As the Author, I want Previews to upload to the same bucket, so that a Preview shows real image URLs.
34. As the Author, I want a prune command that lists R2 objects no build references, dry-run by default and deleting only with an explicit flag, so that orphans are cleaned deliberately.
35. As Googlebot and social scrapers, I want JSON-LD `image`, `og:image` and feed images to be absolute `images.hastcode.com` URLs, so that rich results and link previews resolve.
36. As a Reader with an RSS reader, I want feed images to load from the Image Store, so that feeds look right outside the site.

### Image optimisation

37. As a Reader on a 2× or 3× phone, I want hero variants at 480, 800, 1000, 1200 and 1472 px, so that my device downloads the closest size, never an upsized one.
38. As a Reader on slow 4G, I want the 1200-px hero to be as small as visually acceptable (AVIF quality 55 at 1200 px and above, 60 below), so that the LCP image arrives fast.
39. As a Reader, I want AVIF with 4:4:4 chroma, so that illustration edges and text in heroes do not smear.
40. As a Reader on a browser without AVIF, I want a JPEG fallback, so that the hero still shows.
41. As the Author, I want WebP variants gone, so that each hero produces fewer files for 0.9 points of browser coverage.
42. As the Author, I want a build warning naming the piece when its 1200-px AVIF exceeds 150 KB, so that I re-export a heavy hero.
43. As the Author, I want a phone photo with EXIF rotation to come out upright, so that I never pre-rotate by hand.
44. As the Author, I want a warning when a hero carries a non-sRGB colour profile, so that I know why it looks desaturated and export it as sRGB.
45. As the Author, I want WebP sources rejected with a clear error, so that Go's wrong-range WebP decoding never corrupts a hero.
46. As a Reader, I want the hero marked `fetchpriority="high"` and never lazy-loaded, so that the LCP image is fetched first.
47. As a social scraper, I want a dedicated 1200×630 JPEG for `og:image`, so that link previews are not cropped oddly.
48. As Googlebot, I want the 16:9, 4:3 and 1:1 1200-px JPEG crops kept in JSON-LD, so that rich results have every ratio.
49. As the Author, I want the encoder cache keyed on the new settings, so that retuning re-encodes once and never serves stale variants.

### Operations

50. As the Author, I want a wizard that walks me through the Cloudflare account steps (Workers project, R2 bucket and custom domain, API tokens, GitHub secrets), the Vercel QA setup (Hobby project, `qa.hastcode.com`, Vercel Authentication), and the Production DNS flip, so that I never re-explain setup to an agent.
51. As the Author, I want the README to describe the publish flow, the QA flow, the Preview flow, Image Store modes, the prune command and every secret, so that the project is operable without this conversation.
52. As the Author, I want the first real Production Publish to come with a checklist of what could only be verified live, so that nothing is assumed.

## Implementation Decisions

### Publishing

- Three workflows: **tests** on push to `main`; **preview** on push to any other branch (preview-mode build, deploy to a Cloudflare preview URL, Lighthouse CI with v1's budgets); **publish** on tag push matching `PUBCF-*` or `PUBVL-*`. The publish workflow reads the Target from the tag prefix, fails unless the tagged commit is an ancestor of `main`, builds once in production mode, deploys with the Target's adapter, and for `PUBCF` only, pings IndexNow when the changed-URL artefact exists and the repo variable that gates pinging is on. Concurrency groups: one per Target, no cancel-in-progress.
- Tags are lightweight; the publish command names them `PUB<CF|VL>-<UTC yyyymmdd-HHMMSS>` and refuses to run with uncommitted changes or off `main`.
- Image and diagram caches, and the IndexNow page manifest, persist between runs via the Actions cache keyed per Target, as in v1.

### One build, per-Target manifests

- The builder keeps writing the static tree. It also writes a Vercel manifest (today's routes config, unchanged except as noted under QA) and a Cloudflare manifest: a Workers config declaring the static assets directory with trailing-slash forcing and the 404 page as not-found handling, `run_worker_first` restricted to the signup path, plus a `_headers` file (immutable for hashed assets, HTML revalidation) and a `_redirects` file (`/home` → `/`). Retired community paths need no rule; they fall through to 404.
- Manifests that must live inside the static directory for Cloudflare are excluded from the Vercel upload, and vice versa, by the deploy adapters, so neither Target serves the other's configuration.

### QA Target

- QA builds are production mode. The Vercel manifest adds `X-Robots-Tag: noindex` on every response and the deploy adapter swaps `robots.txt` for a disallow-all file. `site_url` stays `https://hastcode.com`, so canonical, Open Graph and JSON-LD URLs are unchanged.
- Vercel project on Hobby with Vercel Authentication on; the Lighthouse/bypass secret from v1 is no longer needed because Lighthouse runs against Cloudflare Previews.

### Worker

- One Worker module exporting `fetch`. It serves the signup path (`POST`, same contract as the Go handler: form or JSON by `Accept`, honeypot, best-effort per-IP limit, Buttondown subscriber call with the source tag, 303 to the confirmation page or JSON) and forwards everything else to the assets binding. Buttondown key and base URL come from Worker secrets/vars. Rate limiting is in-isolate memory, best-effort, documented as such.
- The Go subscribe handler is untouched and still deployed to QA.

### Image Store

- Config gains an image store section: `store: local | r2` (default `local`), `base_url` (public hostname), bucket name; credentials (account id, access key id, secret) come from environment variables only. In `r2` mode: for every encoded variant, HEAD the key; PUT only on miss, with `Content-Type` and `Cache-Control: public, max-age=31536000, immutable`; nothing image-shaped is written to the static tree; every image URL in view models is absolute on `base_url`. In `local` mode behaviour is v1's. The URL helper that prefixes `site_url` must leave already-absolute URLs alone.
- Object keys mirror v1's asset paths without the `/assets/` prefix (`heroes/…`, `authors/…`, `og-badge.<hash>.png`).
- S3-compatible client per the research: Cloudflare's documented SDK with region `auto`, custom endpoint, and request/response checksum calculation set to "when required" (R2 rejects the SDK's default CRC32 on single PUTs). Uploads run on a bounded worker pool.
- Missing credentials in `r2` mode is a fatal configuration error before any encoding starts.
- Prune: a separate command that lists the bucket, diffs against the keys the current build would produce, prints orphans, and deletes only with `--delete`. It reads the same config and environment.

### Images package

- Widths: 480, 800, 1000, 1200, 1472 (never upscaled). Formats: AVIF + JPEG; WebP removed from encoding, view models, `<picture>` and tests.
- AVIF: 4:4:4 chroma, speed 8, quality 60 below 1200 px and 55 at 1200 px and above. JPEG: stdlib, quality 82, used for the fallback `<img>` and all crops.
- Crops, all 1200 px wide JPEG: 16:9, 4:3, 1:1 for JSON-LD; a new 1200×630 for `og:image` and Twitter.
- Decode: PNG and JPEG sources only (WebP is a `hero-format` error naming the piece). EXIF orientation is read and applied before resizing. A source carrying an ICC profile that is not sRGB raises a warning (`hero-colour`) naming the piece; nothing is converted.
- A new warning (`hero-weight`) fires when the 1200-px AVIF exceeds 150 KB.
- The encoder cache key includes the settings version so the retune invalidates v1's cache once.
- Hero `<img>` carries `fetchpriority="high"` and no `loading` attribute; the preload hint stays. The `sizes` value is unchanged.
- Author photos and the OG badge go through the same store; their encoding is unchanged.

### Wizard and docs

- `setup-deploy.sh` is rewritten around the new Targets: tools; GitHub repo; Cloudflare account, Workers project, R2 bucket + `images.hastcode.com`, two API tokens (Workers edit; R2 object read/write scoped to the bucket); Buttondown keys; Cloudflare Worker secrets; Vercel Hobby project + `qa.hastcode.com` + Vercel Authentication; GitHub secrets and the IndexNow gate variable; first `publish vl` and the QA checklist; first Preview and its checklist; Production DNS flip (apex and `www` to the Worker) behind the ADR-0001 content gate; first `publish cf` and the live checklist. Vercel-Pro and Vercel-domain stages are removed.

## Testing Decisions

- A good test exercises an agreed seam and asserts on observable output: parsed HTML, decoded images, received HTTP requests, files on disk, process exit status. Nothing imports internals.
- **Build seam** (existing, in the builder's test package with `buildtest`): add a fake Image Store, an `httptest` server speaking the handful of S3 operations the builder and prune command use (HEAD, PUT, LIST, DELETE), injected as the R2 endpoint. Tests cover: `r2` mode uploads exactly the expected keys with the expected headers and writes no image files; a second build HEADs and PUTs nothing; `local` mode is unchanged; missing credentials fail before encoding; absolute image URLs in HTML, JSON-LD, Open Graph and feeds; the URL helper does not double-prefix; the Cloudflare manifest files and their contents; the Vercel manifest's `noindex` header rule in QA; new widths, AVIF-only variants, 4:4:4, quality tiers, the 1200×630 crop, decoded and measured; `hero-weight`, `hero-colour`, `hero-format` (WebP) and EXIF-orientation behaviour with fixtures; cache invalidation on settings change; `fetchpriority` on the hero; prune lists and deletes correctly against the fake. Prior art: v1's hero, head, feeds and routes tests.
- **Subscribe HTTP seam, Go** (existing): unchanged.
- **Worker HTTP seam, JavaScript** (new): Vitest with Cloudflare's Workers test pool; a stubbed global `fetch` stands in for Buttondown. One test per Go case: form 303, JSON 200, bad email 400, honeypot swallowed, rate limit trips, upstream failure, method not allowed, key never in a response, non-signup paths forwarded to assets. Run by `make test` alongside Go tests; CI runs both.
- Workflows are linted with actionlint; the publish script is exercised by hand on the first run; real R2, Vercel and Cloudflare behaviour is covered by the wizard's checklists.

## Out of Scope

- Cloudflare Images transformations, Cloudflare Pages, Vercel Pro, any automated DNS flip.
- 410 responses on Production (ADR-0001 amended), CMP/consent, ads going live, sponsorship.
- Colour-profile conversion (warning only), JPEG XL, WebP sources, LQIP/blur placeholders.
- Pruning during the build, multi-bucket layouts, per-image encoder overrides.
- Any change to content, editorial gates, search, feeds beyond image URLs, or the design system.

## Further Notes

- First-run checklist (from the research's open questions): confirm R2 accepts single PUTs with checksums set to when-required; confirm custom-domain responses carry the object `Cache-Control`; confirm a `HEAD` miss then `PUT` costs one Class A operation; confirm `github.sha` resolves the tagged commit for lightweight tags; confirm Workers preview URLs are reachable by Lighthouse CI; confirm Vercel Authentication on Hobby blocks anonymous access to `qa.hastcode.com`.
- Cost expectation: Cloudflare free plan and R2 free tier cover the 18-month traffic target; Vercel Hobby for QA is a terms grey area accepted in ADR-0004.
- Measured in research: AVIF Q60 is 46–58% of JPEG Q82 bytes; 4:4:4 adds ~1%; speed 6 would cost 3–4× encode time for ≤3% bytes.

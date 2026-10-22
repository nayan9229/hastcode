# 03: Image Store: R2 mode

**What to build:** In CI the build uploads every raster (hero variants, crops, author photos, OG badge) to R2 and the site references them by absolute URL on `images.hastcode.com`, so each Explainer ships as one file. Config gains an image store section: `store: local | r2` (default `local`, today's behaviour), the public `base_url`, and the bucket name; the account id, access key id and secret come from environment variables only. In `r2` mode the build HEADs each content-hashed key and PUTs only on a miss, with the right `Content-Type` and `Cache-Control: public, max-age=31536000, immutable`, on a bounded worker pool; it writes no image files into the static tree; every image URL in HTML, JSON-LD, Open Graph and feeds is absolute on `base_url`; the helper that prefixes `site_url` leaves absolute URLs alone. Object keys mirror v1's asset paths without the `/assets/` prefix. A build in `r2` mode without credentials fails before any encoding with a clear message. The client follows the research: Cloudflare's documented SDK, region `auto`, custom endpoint, request and response checksum calculation set to "when required".

Spec: `../spec.md` ("Image Store"). Research: `docs/research/cloudflare-r2-and-deploy.md` sections 1–3. Test support gets a fake Image Store: an `httptest` server speaking S3 HEAD, PUT, LIST and DELETE for a handful of keys, injected as the endpoint, so the builder (this ticket) and the prune command (ticket 04) test against it.

**Blocked by:** 01 (Hero variants and quality retune)

**Status:** resolved

- [x] Build-seam tests with the fake store: `r2` mode PUTs exactly the expected keys with the expected `Content-Type` and `Cache-Control`; a second build HEADs every key and PUTs nothing; no file under the assets path exists in the output
- [x] HTML `<picture>`/`<img>`, JSON-LD `image`, `og:image`, author `Person`/`Organization` image, the OG badge and feed images all use absolute `base_url` URLs; nothing is double-prefixed with `site_url`
- [x] `local` mode output is byte-identical to before this ticket (existing tests untouched and green)
- [x] `r2` mode with any credential missing fails before encoding, naming the missing variable; `make dev` needs nothing
- [x] Upload failures (non-2xx) fail the build naming the key
- [x] README documents store modes, config keys, environment variables and the key layout

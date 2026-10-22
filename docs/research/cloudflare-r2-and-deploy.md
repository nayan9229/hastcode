# Cloudflare R2 for build-time images, and a tag-driven Cloudflare deploy target

Research date: 2026-10-22. Every Cloudflare, AWS, GitHub, Go and W3C page cited below was read on that date (Cloudflare pages via their `index.md` "View as Markdown" variants where exact wording mattered). Cloudflare moves fast; re-check the cited URL before committing to a number. Claims that are inference rather than something a source states are marked **[inferred]**; claims I could not source are marked **[unsourced]**.

Scope: the Go static builder that runs in GitHub Actions and today writes content-hashed image variants (`/assets/heroes/<name>-<w>x<h>.<hash10>.<ext>`, AVIF/WebP/JPEG) into the Vercel Build Output and deploys with `vercel deploy --prebuilt`. Goal: (1) upload those variants to R2 at build time and reference them by absolute URL from HTML, JSON-LD, Open Graph and RSS; (2) add a second deploy target on Cloudflare (Pages or Workers Static Assets), selected by pushing a git tag.

---

## TL;DR

- **Go → R2: use `aws-sdk-go-v2` exactly as Cloudflare's own example does** — `config.WithRegion("auto")`, static credentials, `o.BaseEndpoint = https://<ACCOUNT_ID>.r2.cloudflarestorage.com`. Access Key ID = the R2 API token's `id`; Secret = SHA-256 of the token value. Cloudflare publishes no minio-go example and minio-go's README does not mention R2. ([aws-sdk-go example](https://developers.cloudflare.com/r2/examples/aws/aws-sdk-go/), [R2 tokens](https://developers.cloudflare.com/r2/api/tokens/), [S3 API: region](https://developers.cloudflare.com/r2/api/s3/api/))
- **Turn off the SDK's default checksums.** Since `service/s3` v1.73.0 (2025-01-15) aws-sdk-go-v2 adds a CRC32 checksum to every Put by default; R2's compatibility table lists CRC32 as `COMPOSITE`-only (multipart) and lists no checksum headers at all under `PutObject`. AWS itself says to set `WHEN_REQUIRED` for third-party S3 implementations. Set `aws.RequestChecksumCalculationWhenRequired` and `aws.ResponseChecksumValidationWhenRequired`. ([AWS SDK announcement](https://github.com/aws/aws-sdk-go-v2/discussions/2960), [AWS data-integrity reference](https://docs.aws.amazon.com/sdkref/latest/guide/feature-dataintegrity.html), [R2 checksum types](https://developers.cloudflare.com/r2/api/s3/api/#checksum-types))
- **Skipping re-uploads is safe and cheap.** R2 implements `If-None-Match` on `PutObject`, "evaluated atomically against the committed bucket state"; a second writer gets `412 PreconditionFailed`. `HeadObject` is also supported and is a Class B op (10 M/month free) versus Class A for PUT (1 M/month free). With content-hashed names either approach is correct. ([R2 S3 API](https://developers.cloudflare.com/r2/api/s3/api/), [R2 pricing](https://developers.cloudflare.com/r2/pricing/))
- **`Cache-Control` and `Content-Type` are set per object on PUT** (both are listed as supported system metadata). Images this size never need multipart (R2 single PUT limit is 4.995 GiB; `manager.Uploader` only switches to multipart above its 5 MiB default part size, and that package is now marked deprecated in favour of `feature/s3/transfermanager`). ([R2 S3 API](https://developers.cloudflare.com/r2/api/s3/api/), [R2 limits](https://developers.cloudflare.com/r2/platform/limits/), [s3/manager](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/s3/manager))
- **Serve from a custom domain, not `r2.dev`.** Cloudflare: `r2.dev` "is intended for non-production traffic", "is rate-limited and should only be used for development purposes" (variable limit in the "hundreds of requests/second", then 429), and "Avoid creating a CNAME record pointing to the r2.dev subdomain. This is an unsupported access path". A custom domain must be a zone in the same account and unlocks Cloudflare Cache, WAF, Access and Bot Management. ([Public buckets](https://developers.cloudflare.com/r2/buckets/public-buckets/), [R2 limits](https://developers.cloudflare.com/r2/platform/limits/))
- **The CDN will cache the images and honour `Cache-Control`.** AVIF, WEBP, JPG, JPEG, PNG and WOFF2 are all in the default cached-extension list (default edge TTL 2 h for 200s); Origin Cache Control is on by default for Free/Pro/Business and "strictly respect[s] Cache-Control directives received from the origin"; `immutable` is browser-only. The one-hop statement "R2 custom-domain responses carry the object's `Cache-Control` header" is **[inferred]** from R2's `httpMetadata` docs plus S3 semantics; verify with `curl -I` once the domain is live. ([Default cache behavior](https://developers.cloudflare.com/cache/concepts/default-cache-behavior/), [Cache-Control](https://developers.cloudflare.com/cache/concepts/cache-control/), [R2 Workers API `httpMetadata`](https://developers.cloudflare.com/r2/api/workers/workers-api-reference/))
- **CORS: none for `<img>`, required for fonts.** HTML makes image requests in `no-cors` mode unless `crossorigin` is set; CSS Fonts 4 fetches `@font-face` sources with CORS mode `cors`, so a font on `images.hastcode.com` needs a bucket CORS policy allowing `https://hastcode.com`. ([HTML: CORS settings attributes](https://html.spec.whatwg.org/multipage/urls-and-fetching.html#cors-settings-attributes), [CSS Fonts 4 §4.8.2](https://www.w3.org/TR/css-fonts-4/#font-fetching-requirements), [R2 CORS](https://developers.cloudflare.com/r2/buckets/cors/))
- **Image transformations are on the Free plan** (5,000 unique transformations/month, then error 9422; paid: $0.50 per 1,000 after the first 5,000) and can source from the same zone, so `/cdn-cgi/image/...` on `images.hastcode.com` is available as a fallback — but the builder already emits sized variants, so this is optional. ([Images pricing](https://developers.cloudflare.com/images/pricing/), [Transform via URL](https://developers.cloudflare.com/images/transform-images/transform-via-url/), [Sources](https://developers.cloudflare.com/images/transform-images/sources/))
- **R2 price: $0.015/GB-month, Class A $4.50/M, Class B $0.36/M, egress $0; free tier 10 GB-month, 1 M Class A, 10 M Class B per month**, usage rounded *up* to the next unit. `DeleteObject` is free; `ListObjects` is Class A. A hero-image library of a few GB and a few thousand objects stays inside the free tier. ([R2 pricing](https://developers.cloudflare.com/r2/pricing/))
- **Deploy target: Workers Static Assets, not Pages.** Cloudflare's Pages landing page now says, verbatim: "Are you sure you want to use Pages? Workers supports most Pages use cases and offers a broader feature set. It is Cloudflare's primary platform for building applications. Start new projects with Workers." Workers natively supports `_headers`/`_redirects`, 20,000 files (Free) / 100,000 (Paid), free unlimited static-asset requests, and `run_worker_first` lets a tiny Worker emit `410` for `/video/*`, `/community/*`, `/group/*` while everything else is served without billing. ([Pages landing](https://developers.cloudflare.com/pages/), [Workers limits](https://developers.cloudflare.com/workers/platform/limits/), [Billing and limitations](https://developers.cloudflare.com/workers/static-assets/billing-and-limitations/), [Routing: Worker script](https://developers.cloudflare.com/workers/static-assets/routing/worker-script/))
- **`_redirects` cannot emit 410 on either platform.** Both pages list only "Redirects (301, 302, 303, 307, 308)" and proxying (200) as supported and mark "Rewrites (other status codes)" ❌. ([Workers `_redirects`](https://developers.cloudflare.com/workers/static-assets/redirects/), [Pages `_redirects`](https://developers.cloudflare.com/pages/configuration/redirects/))
- **GitHub Actions tag trigger:** `on: push: tags: ['PUBCF-*']`. "If you define only `tags`/`tags-ignore` or only `branches`/`branches-ignore`, the workflow won't run for events affecting the undefined Git ref" — so a tag push does not fire a `branches:`-only workflow and vice versa. Read the tag as `github.ref_name` (`github.ref` is `refs/tags/<tag>`, `github.ref_type` is `tag`). Gate on "tag points at a commit on `main`" with `fetch-depth: 0` + `git merge-base --is-ancestor`. Use a `concurrency.group` per target. ([Workflow syntax](https://docs.github.com/en/actions/writing-workflows/workflow-syntax-for-github-actions), [Default variables](https://docs.github.com/en/actions/reference/workflows-and-actions/variables), [git merge-base](https://git-scm.com/docs/git-merge-base), [Concurrency](https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/control-the-concurrency-of-workflows-and-jobs))

---

## Decision table

| Question | Answer | Source |
|---|---|---|
| Go SDK for R2 | `github.com/aws/aws-sdk-go-v2` (`config`, `credentials`, `service/s3`) with `BaseEndpoint` | [aws-sdk-go example](https://developers.cloudflare.com/r2/examples/aws/aws-sdk-go/) |
| Endpoint | `https://<ACCOUNT_ID>.r2.cloudflarestorage.com` (EU: `.eu.`, US: `.us.`, FedRAMP: `.fedramp.` — only for jurisdiction buckets) | [R2 tokens](https://developers.cloudflare.com/r2/api/tokens/) |
| Region | `auto` (`""` and `us-east-1` alias to it) | [R2 S3 API](https://developers.cloudflare.com/r2/api/s3/api/) |
| Credentials | R2 API token → Access Key ID = token `id`, Secret = SHA-256(token value); shown once in the dashboard | [R2 tokens](https://developers.cloudflare.com/r2/api/tokens/) |
| Token permission | **Object Read & Write**, scoped to the one bucket (S3 API only; Admin needed only to create buckets) | [R2 tokens](https://developers.cloudflare.com/r2/api/tokens/) |
| Skip existing | `HeadObject` (Class B) or `PutObject` + `If-None-Match: *` (412 if exists, atomic) | [R2 S3 API](https://developers.cloudflare.com/r2/api/s3/api/) |
| Per-object headers on PUT | `Content-Type`, `Cache-Control`, `Content-Encoding`, `Content-Disposition`, `Content-Language`, `Expires` | [R2 S3 API](https://developers.cloudflare.com/r2/api/s3/api/) |
| Multipart | min part 5 MiB (last part exempt), max part 5 GiB, 10,000 parts, all parts equal size; single PUT ≤ 4.995 GiB | [Multipart](https://developers.cloudflare.com/r2/objects/multipart-objects/), [Limits](https://developers.cloudflare.com/r2/platform/limits/) |
| Rate limits | 1 write/s per key (429 above); bucket-management ops 50/s; no documented per-bucket object-op limit | [R2 limits](https://developers.cloudflare.com/r2/platform/limits/) |
| Public serving | Custom domain on the bucket (zone in same account); up to 100 per bucket | [Public buckets](https://developers.cloudflare.com/r2/buckets/public-buckets/), [Limits](https://developers.cloudflare.com/r2/platform/limits/) |
| Deploy target | Workers Static Assets via `wrangler deploy` + `[assets]` | [Pages landing](https://developers.cloudflare.com/pages/), [Wrangler config](https://developers.cloudflare.com/workers/wrangler/configuration/) |
| 410 Gone | Worker with `run_worker_first = ["/video/*", "/community/*", "/group/*"]` | [Binding](https://developers.cloudflare.com/workers/static-assets/binding/) |
| CI action | `cloudflare/wrangler-action@v4` (defaults to Wrangler v4) | [wrangler-action](https://github.com/cloudflare/wrangler-action) |
| API token for CI | "Edit Cloudflare Workers" template (Workers); Account → Cloudflare Pages → Edit (Pages) | [GitHub Actions guide](https://developers.cloudflare.com/workers/ci-cd/external-cicd/github-actions/), [Pages direct upload + CI](https://developers.cloudflare.com/pages/how-to/use-direct-upload-with-continuous-integration/) |

---

## 1. R2 from Go

### 1.1 Endpoint, region, auth

- Endpoint: "`https://<ACCOUNT_ID>.r2.cloudflarestorage.com`". Jurisdictional buckets must use their matching endpoint (`<ACCOUNT_ID>.eu.r2.cloudflarestorage.com`, `.us.`, `.fedramp.`); a plain bucket uses the default. ([R2 tokens](https://developers.cloudflare.com/r2/api/tokens/))
- Region: "When using the S3 API, the region for an R2 bucket is `auto`. For compatibility with tools that do not allow you to specify a region, an empty value and `us-east-1` will alias to the `auto` region." ([R2 S3 API](https://developers.cloudflare.com/r2/api/s3/api/))
- Credentials come from an **R2 API token** (R2 → Manage API tokens), account-owned or user-owned. "Access Key ID: The `id` of the API token." "Secret Access Key: The SHA-256 hash of the API token `value`." The dashboard shows both once: "You will not be able to access your Secret Access Key again after this step." Permission levels: Admin Read & Write / Admin Read only (bucket + object ops) and **Object Read & Write** / Object Read only, which "Allows the ability to read, write, and list objects in specific buckets" and work only with the S3 API. Bucket scoping applies to the Object-level permissions. Account tokens "remain valid until manually revoked or until they reach a preset expiration date." ([R2 tokens](https://developers.cloudflare.com/r2/api/tokens/))
- For CI: one **Object Read & Write** token scoped to the images bucket, stored as two GitHub secrets (`R2_ACCESS_KEY_ID`, `R2_SECRET_ACCESS_KEY`) plus `CLOUDFLARE_ACCOUNT_ID`. **[inferred]** from the permission table above.

### 1.2 Which SDK Cloudflare documents

Cloudflare's R2 "S3 SDKs" examples include a Go page that uses **aws-sdk-go v2** (`github.com/aws/aws-sdk-go-v2/{aws,config,credentials,service/s3}`). The documented client construction is:

```go
cfg, err := config.LoadDefaultConfig(context.TODO(),
    config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
    config.WithRegion("auto"), // "Required by SDK but not used by R2" (Cloudflare's comment)
)
client := s3.NewFromConfig(cfg, func(o *s3.Options) {
    o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
})
```

([aws-sdk-go example](https://developers.cloudflare.com/r2/examples/aws/aws-sdk-go/); the example goes on to `ListObjectsV2`, `ListBuckets` and a presigned PUT.) The R2 examples index lists Rclone, "S3 SDKs", Terraform and SSE-C; there is no minio-go page ([R2 examples](https://developers.cloudflare.com/r2/examples/)), and the minio-go README does not mention R2 ([minio-go](https://github.com/minio/minio-go)). Use aws-sdk-go-v2.

### 1.3 Checksums: the one thing Cloudflare's example omits

- AWS: "clients default to enabling an additional checksum on all Put calls and enabling validation on Get calls" from `service/s3` v1.73.0 (2025-01-15). The default algorithm for the Go SDK is **CRC32**. The announcement warns this "may result in compatibility issues with 3rd-party S3 emulators" and recommends `WHEN_REQUIRED` for both `request_checksum_calculation` and `response_checksum_validation`. ([aws-sdk-go-v2 discussion #2960](https://github.com/aws/aws-sdk-go-v2/discussions/2960), [AWS SDK reference: data integrity](https://docs.aws.amazon.com/sdkref/latest/guide/feature-dataintegrity.html))
- R2's checksum table: CRC64NVME supports `FULL_OBJECT` only; CRC32, CRC32C, SHA1, SHA256 support `COMPOSITE` only. The `PutObject` row of the compatibility table lists conditional headers and system metadata but **no checksum headers at all**, so whether R2 accepts or ignores the SDK's default `x-amz-checksum-crc32` on a single-part PUT is not documented. ([R2 S3 API](https://developers.cloudflare.com/r2/api/s3/api/)) Cloudflare's Go example sets nothing here ([aws-sdk-go example](https://developers.cloudflare.com/r2/examples/aws/aws-sdk-go/)).
- Safe setting (both helpers exist in `config/load_options.go`; the `s3.Options` struct also carries both fields; env vars `AWS_REQUEST_CHECKSUM_CALCULATION` / `AWS_RESPONSE_CHECKSUM_VALIDATION` are read by `config/env_config.go`):

```go
cfg, err := config.LoadDefaultConfig(ctx,
    config.WithRegion("auto"),
    config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(id, secret, "")),
    config.WithRequestChecksumCalculation(aws.RequestChecksumCalculationWhenRequired),
    config.WithResponseChecksumValidation(aws.ResponseChecksumValidationWhenRequired),
)
```

([config/load_options.go](https://github.com/aws/aws-sdk-go-v2/blob/main/config/load_options.go), [aws/checksum.go](https://github.com/aws/aws-sdk-go-v2/blob/main/aws/checksum.go), [config/env_config.go](https://github.com/aws/aws-sdk-go-v2/blob/main/config/env_config.go)). Note `DeleteObjects` is a checksum-*required* operation on S3 ("The Content-MD5 request header is required for all Multi-Object Delete requests"); with `WhenRequired` the SDK still sends a checksum for it. ([S3 DeleteObjects](https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObjects.html), [discussion #2960](https://github.com/aws/aws-sdk-go-v2/discussions/2960))

### 1.4 Skipping objects that already exist

Two documented options, both correct given content-hashed keys:

1. **`HeadObject` first.** R2 implements `HeadObject` with `If-Match`/`If-None-Match`. Class B operation ($0.36/M, 10 M free). In aws-sdk-go-v2 a 404 surfaces as the error type `*types.NotFound` (`service/s3/types/errors.go`), so `errors.As(err, &nf)` → upload. ([R2 S3 API](https://developers.cloudflare.com/r2/api/s3/api/), [R2 pricing](https://developers.cloudflare.com/r2/pricing/), [s3/types/errors.go](https://github.com/aws/aws-sdk-go-v2/blob/main/service/s3/types/errors.go))
2. **Conditional PUT.** "The conditional headers on `PutObject` (such as `If-Match` and `If-None-Match`) are evaluated atomically against the committed bucket state, as part of the same write that commits the object. For example, when multiple clients concurrently issue `PutObject` with `If-None-Match: *` for a previously absent key, at most one write succeeds and the others fail with `412 PreconditionFailed`." The SDK field is `PutObjectInput.IfNoneMatch` ("Expects the '*' (asterisk) character"; a 409 `ConditionalRequestConflict` means retry). ([R2 S3 API](https://developers.cloudflare.com/r2/api/s3/api/), [R2 S3 extensions](https://developers.cloudflare.com/r2/api/s3/extensions/), [api_op_PutObject.go](https://github.com/aws/aws-sdk-go-v2/blob/main/service/s3/api_op_PutObject.go))
3. **One `ListObjectsV2` sweep** of the prefix before uploading (1,000 keys per call, Class A) builds an in-memory set and is cheaper than N `HeadObject`s once N > ~1,000 objects per build. **[inferred]** from the pricing classes. ([S3 ListObjectsV2](https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjectsV2.html), [R2 pricing](https://developers.cloudflare.com/r2/pricing/))

Whether a conditional PUT that fails with 412 is billed as a Class A operation is not stated anywhere I read — see Open questions. For a site with thousands of images per build, HEAD-then-PUT (option 1) is the simplest and costs nothing at free-tier scale.

### 1.5 Per-object headers, sizes, multipart, concurrency

- `PutObject` system metadata R2 stores: `Content-Type`, `Cache-Control`, `Content-Disposition`, `Content-Encoding`, `Content-Language`, `Expires`, `Content-MD5`; storage class `STANDARD` / `STANDARD_IA`. Not supported on PUT: `x-amz-tagging`, ACL headers, object lock, SSE-KMS. ([R2 S3 API](https://developers.cloudflare.com/r2/api/s3/api/)) SDK fields: `CacheControl`, `ContentType`, `ContentEncoding` on `PutObjectInput`. ([api_op_PutObject.go](https://github.com/aws/aws-sdk-go-v2/blob/main/service/s3/api_op_PutObject.go))
- Size limits: object 5 TiB; single-part upload "5 GiB" (footnote: 4.995 GiB); max 10,000 parts; key ≤ 1,024 bytes; metadata ≤ 8,192 bytes. ([R2 limits](https://developers.cloudflare.com/r2/platform/limits/))
- Multipart rules: "Minimum part size: 5 MiB (except for the last part)", "Maximum part size: 5 GiB", "Maximum number of parts: 10,000", "All parts except the last must be the same size." Cloudflare says you "do not typically need to implement multipart logic yourself when using the S3 API." ([Multipart objects](https://developers.cloudflare.com/r2/objects/multipart-objects/))
- `feature/s3/manager.Uploader`: `DefaultUploadPartSize` = `MinUploadPartSize` = 5 MiB, `DefaultUploadConcurrency` = 5, `MaxUploadParts` = 10,000; "Uploads that are smaller than the Uploader's PartSize will be uploaded using the PutObject API operation." The package is marked deprecated, "superceded by feature/s3/transfermanager". Image variants are far below 5 MiB, so call `PutObject` directly with a `*bytes.Reader` body (seekable, so the SDK can set `Content-Length`). ([s3/manager](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/s3/manager))
- Rate limits: "Concurrent writes to the same object name (key) at a higher rate [than 1 per second] return HTTP 429"; bucket-management ops 50/s per bucket; the 1,200 req/5 min limit is for the Cloudflare REST API, not the S3 endpoint. No per-bucket object-operation ceiling is published. ([R2 limits](https://developers.cloudflare.com/r2/platform/limits/)) A bounded worker pool (8–16 goroutines) with retry on 429/5xx is reasonable; the number is **[unsourced]** engineering judgement.

### 1.6 Sketch

```go
type uploader struct{ c *s3.Client; bucket, prefix string }

func (u uploader) put(ctx context.Context, key, ctype string, body []byte) error {
    _, err := u.c.HeadObject(ctx, &s3.HeadObjectInput{Bucket: &u.bucket, Key: &key})
    if err == nil { return nil }                // already there: hashed key => identical bytes
    var nf *types.NotFound
    if !errors.As(err, &nf) { return err }
    _, err = u.c.PutObject(ctx, &s3.PutObjectInput{
        Bucket: &u.bucket, Key: &key, Body: bytes.NewReader(body),
        ContentType:  aws.String(ctype),                                   // image/avif, image/webp, image/jpeg
        CacheControl: aws.String("public, max-age=31536000, immutable"),
        IfNoneMatch:  aws.String("*"),                                     // belt and braces; 412 == someone else won
    })
    var ae smithy.APIError
    if errors.As(err, &ae) && ae.ErrorCode() == "PreconditionFailed" { return nil }
    return err
}
```

The `errors.As`/`smithy.APIError` idiom is standard aws-sdk-go-v2 usage; the specific error-code string for 412 is **[unsourced]** (R2 docs say `412 PreconditionFailed`; the SDK surfaces the XML `Code`). Verify with a two-run test.

---

## 2. Serving R2 publicly

### 2.1 Custom domain vs `r2.dev`

- Custom domain: "The domain being used must have been added as a zone in the same account as the R2 bucket" (partial/CNAME setup is allowed for domains not fully on Cloudflare). Bucket → Settings → Custom Domains → Add; "The status takes a few minutes to change from Initializing to Active." Enterprise zones must release any zone hold first. 100 custom domains per bucket. ([Public buckets](https://developers.cloudflare.com/r2/buckets/public-buckets/), [R2 limits](https://developers.cloudflare.com/r2/platform/limits/))
- `r2.dev`: "This endpoint is intended for non-production traffic." "Public access through `r2.dev` subdomains is rate-limited and should only be used for development purposes." "Avoid creating a CNAME record pointing to the `r2.dev` subdomain. This is an **unsupported access path**, and we cannot guarantee consistent reliability or performance." The limits page: a "variable rate limit" in the "hundreds of requests/second", excess gets `429 Too Many Requests` and bandwidth may be throttled. "To enable access management, cache, and bot management features, you must set up a custom domain." ([Public buckets](https://developers.cloudflare.com/r2/buckets/public-buckets/), [R2 limits](https://developers.cloudflare.com/r2/platform/limits/))
- Limitation: "Currently, public buckets do not let you list the bucket contents at the root of your (sub) domain." ([Public buckets](https://developers.cloudflare.com/r2/buckets/public-buckets/))
- Egress: "does not incur data transfer (egress) charges and is free" whether via Workers, S3 API or r2.dev; a custom domain in front is Cloudflare CDN, which is not metered either on these plans. ([R2 pricing](https://developers.cloudflare.com/r2/pricing/))

So: `images.hastcode.com` (a record in the existing `hastcode.com` zone) as the bucket's custom domain; all absolute URLs in HTML/JSON-LD/OG/RSS point there.

### 2.2 CDN caching of R2 objects and `Cache-Control`

- "Domain access through a custom domain allows you to use Cloudflare Cache to accelerate access to your R2 bucket." "Configure your cache to use Smart Tiered Cache to have a single upper-tier data center next to your R2 bucket." "By default, only certain file types are cached. To cache all files in your bucket, you must set a Cache Everything page rule." ([Public buckets](https://developers.cloudflare.com/r2/buckets/public-buckets/))
- The default cached-extension list includes AVIF, WEBP, JPG, JPEG, PNG and WOFF2; "Cloudflare only caches based on file extension and not by MIME type"; default edge TTL for 200/206/301 is 120 minutes, 302/303 20 minutes, 404/410 3 minutes; HTML and JSON are not cached by default. So image variants are cached with no extra rule; a Cache Rule is only needed if you ever put HTML/JSON in the bucket. ([Default cache behavior](https://developers.cloudflare.com/cache/concepts/default-cache-behavior/))
- Origin Cache Control: Cloudflare "should strictly respect `Cache-Control` directives received from the origin server" (on by default for Free, Pro and Business). `max-age` sets the edge lifetime; `s-maxage` overrides it for shared caches; Cloudflare "respects whichever value is higher: the Browser Cache TTL in Cloudflare or the `max-age` header"; `no-store`/`private` are never cached; "`immutable` ... has no effect on public caches like Cloudflare, but does change browser behavior"; values must be integers. ([Cache-Control](https://developers.cloudflare.com/cache/concepts/cache-control/))
- R2 keeps `cacheControl`, `contentType`, `contentEncoding`, etc. as `httpMetadata` that "match the HTTP metadata passed when the object was created" and `writeHttpMetadata` "applies their corresponding HTTP headers". That the public-bucket/custom-domain path emits the object's `Cache-Control` as a response header so the CDN sees it is **[inferred]** — it is the ordinary S3 behaviour and the reason the metadata exists, but no single page states it for custom domains. Confirm with `curl -I https://images.hastcode.com/<key>` and check `cf-cache-status` on the second hit. ([R2 Workers API](https://developers.cloudflare.com/r2/api/workers/workers-api-reference/))
- Recommendation: `Cache-Control: public, max-age=31536000, immutable` on every hashed object (content-hashed names make year-long caching safe), Smart Tiered Cache on, no Cache Everything rule needed for images.

### 2.3 CORS

- `<img>`: HTML's "create a potential-CORS request" sets "mode be "no-cors" if corsAttributeState is No CORS", and the `crossorigin` attribute's "missing value default is the No CORS state". No CORS headers needed for plain `<img>`, `<picture>`/`<source>`, OG images or RSS enclosures. ([HTML standard](https://html.spec.whatwg.org/multipage/urls-and-fetching.html#cors-settings-attributes))
- Fonts: CSS Fonts 4 fetches `@font-face` URLs with "destination "font", CORS mode "cors"", and notes "fonts will typically not be loaded cross-origin unless authors specifically take steps to permit cross-origin loads. Sites can explicitly allow cross-site loading of font data using the Access-Control-Allow-Origin HTTP header." ([CSS Fonts 4 §4.8.2](https://www.w3.org/TR/css-fonts-4/#font-fetching-requirements))
- R2: set a bucket CORS policy (dashboard → bucket → Settings → CORS Policy, or `npx wrangler r2 bucket cors set <BUCKET> --file cors.json`). Custom domains "automatically return CORS response headers" when a policy exists; the request must carry an `Origin` header; "CORS rule propagation can, in rare cases, take up to 30 seconds"; purge the hostname's cache after changing a policy because cached assets may lack the headers. `AllowedOrigins` entries are `scheme://host[:port]` with at most one `*`. Example policy for fonts served to the site: `[{"AllowedOrigins":["https://hastcode.com"],"AllowedMethods":["GET"]}]`. ([R2 CORS](https://developers.cloudflare.com/r2/buckets/cors/))

### 2.4 Image transformations (`/cdn-cgi/image/`)

- Plan: "The Free plan includes access to the transformations feature, which lets you optimize images stored outside of Images." Free: "up to 5,000 unique transformations each month for free"; beyond that new transformations return error `9422` and "You will not be charged for exceeding the limits in the Free plan." Paid: first 5,000 included, then "$0.50 / 1,000 unique transformations / month". A unique transformation is a source image + parameter combination, counted once per calendar month; "Cloudflare excludes `format` and `quality` when determining whether a transformation is unique." ([Images pricing](https://developers.cloudflare.com/images/pricing/))
- Enable per zone: Images → Transformations → select zone → enable. Sources default to "allowed origins" and "Cloudflare will always allow source images from the same zone"; other origins are added explicitly (root domains do not cover subdomains; `*.b.com` does). ([Transform images](https://developers.cloudflare.com/images/transform-images/), [Sources](https://developers.cloudflare.com/images/transform-images/sources/))
- URL: `https://<ZONE>/cdn-cgi/image/<OPTIONS>/<SOURCE-IMAGE>` where the source is "an absolute path on the origin server or an absolute URL"; options include `width`, `height`, `fit` (`scale-down` default, `cover`, `contain`, …), `format=auto|avif|webp|jpeg`, `quality` (default 85). ([Transform via URL](https://developers.cloudflare.com/images/transform-images/transform-via-url/))
- Relevance here: the builder already produces AVIF/WebP/JPEG at fixed widths, so transformations are not needed for the main path. They are a cheap fallback for one-off sizes (OG card crops, newsletter thumbnails) and would run against `images.hastcode.com` as a same-zone source. **[inferred]** that a subdomain of the zone counts as "the same zone"; the Sources page phrases the exception around root vs subdomains for *other* origins only.

---

## 3. R2 free tier and pricing

Standard storage, as published ([R2 pricing](https://developers.cloudflare.com/r2/pricing/)):

| Item | Price | Free tier (per month) |
|---|---|---|
| Storage | $0.015 / GB-month | 10 GB-month |
| Class A operations | $4.50 / million | 1 million |
| Class B operations | $0.36 / million | 10 million |
| Egress | Free | — |
| Infrequent Access | $0.01 / GB-month, Class A $9.00/M, Class B $0.90/M, retrieval $0.01/GB, 30-day minimum; **no free tier** | — |

- Class A: `ListBuckets`, `PutBucket`, `ListObjects`, `PutObject`, `CopyObject`, `CompleteMultipartUpload`, `CreateMultipartUpload`, `LifecycleStorageTierTransition`, `ListMultipartUploads`, `UploadPart`, `UploadPartCopy`, `ListParts`, `PutBucketEncryption`, `PutBucketCors`, `PutBucketLifecycleConfiguration`.
- Class B: `HeadBucket`, `HeadObject`, `GetObject`, `UsageSummary`, `GetBucketEncryption`, `GetBucketLocation`, `GetBucketCors`, `GetBucketLifecycleConfiguration`.
- Free: `DeleteObject`, `DeleteBucket`, `AbortMultipartUpload`. (`DeleteObjects`, the batch form, is not named — see Open questions.)
- Minimum billing: no minimum storage duration for Standard; "Cloudflare rounds up your usage to the next billing unit" ("If you have used 1.1 GB-month, you will be billed for 2 GB-month"; 1,000,001 ops bills as 2 million). "A GB-month is calculated by averaging the *peak* storage per day over a billing period (30 days)."
- Worked estimate for this site: 2,000 articles × ~12 variants × ~150 KB ≈ 3.6 GB stored → inside 10 GB free; per build ~24,000 HEADs (Class B) + only new PUTs (Class A); monthly reads come from the CDN cache, not R2, after the first hit per PoP. Expect $0/month. The multipliers are **[unsourced]** assumptions about the site; the prices are not.

---

## 4. Deleting orphans

- Listing: `ListObjectsV2` is implemented (Cloudflare: "it is recommended that you use `ListObjectsV2` instead" of `ListObjects`); it is a **Class A** op. S3 semantics: "By default, the action returns up to 1,000 key names", paginate with `IsTruncated`/`NextContinuationToken`. R2's Workers API `list()` likewise caps at 1,000 per call. Cost of a full sweep: ⌈N/1000⌉ Class A calls — for 25,000 objects that is 25 calls, i.e. ~$0.0001 and inside the free tier. ([R2 S3 API](https://developers.cloudflare.com/r2/api/s3/api/), [R2 pricing](https://developers.cloudflare.com/r2/pricing/), [S3 ListObjectsV2](https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjectsV2.html), [R2 Workers API](https://developers.cloudflare.com/r2/api/workers/workers-api-reference/))
- Deleting: `DeleteObject` is free. `DeleteObjects` (batch) is implemented; S3 allows "a list of up to 1,000 keys" per request and requires `Content-MD5` (the SDK handles it, see §1.3). ([R2 pricing](https://developers.cloudflare.com/r2/pricing/), [R2 S3 API](https://developers.cloudflare.com/r2/api/s3/api/), [S3 DeleteObjects](https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObjects.html))
- Algorithm **[inferred]**: after a successful deploy, list the build's prefix (`heroes/`), subtract the set of keys the build just referenced (the builder already knows every emitted variant), and `DeleteObjects` the remainder in batches of ≤1,000. Guard rails that are engineering judgement, **[unsourced]**: only ever delete under the prefix the builder owns; run the sweep as a separate step/flag (`--prune`) rather than on every build; keep a grace window (e.g. skip objects whose `LastModified` is < 7 days old) because old HTML in RSS readers, link previews and the Vercel deployment that is still live during a Cloudflare cutover may still reference the previous hashes. Every orphan left behind costs $0.015/GB-month, so pruning is hygiene, not savings.
- In aws-sdk-go-v2 the paginator is `s3.NewListObjectsV2Paginator` ([pkg.go.dev s3](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3)); the name was not re-read in source today.

---

## 5. Cloudflare deploy from GitHub Actions

### 5.1 Pages vs Workers Static Assets — what Cloudflare says

- Pages landing page, first callout, verbatim: **"Are you sure you want to use Pages?** Workers supports most Pages use cases and offers a broader feature set. It is Cloudflare's primary platform for building applications. Start new projects with Workers." ([Pages](https://developers.cloudflare.com/pages/), page stamped "Last updated Aug 25, 2026")
- Compatibility matrix: "Unlike Pages, Workers has a distinctly broader set of features available to it." "`_headers` and `_redirects` files are supported natively in Workers with static assets." Pages-only: "Custom domains outside Cloudflare zones", full Early Hints, full Branch Deploy Controls. Workers-only: non-root routes, Gradual Deployments, Workers Logs/Logpush/Tail Workers, Cron Triggers, Queue consumers, Image Resizing binding, Rate Limiting. ([Compatibility matrix](https://developers.cloudflare.com/workers/static-assets/compatibility-matrix/), [Migrate from Pages](https://developers.cloudflare.com/workers/static-assets/migration-guides/migrate-from-pages/))
- The migration guide says the Pages commands are replaced by `wrangler dev` and `wrangler deploy`, `pages_build_output_dir` becomes `assets.directory`, and `not_found_handling` must be set explicitly "since Workers no longer infers this". ([Migrate from Pages](https://developers.cloudflare.com/workers/static-assets/migration-guides/migrate-from-pages/))

### 5.2 Workers Static Assets configuration

From the Wrangler configuration reference ([Wrangler config](https://developers.cloudflare.com/workers/wrangler/configuration/)) and the binding page ([Binding](https://developers.cloudflare.com/workers/static-assets/binding/)):

| Key | Values | Default |
|---|---|---|
| `assets.directory` | path | — |
| `assets.binding` | name, e.g. `ASSETS`; "Only useful when a Worker script is set with `main`" | — |
| `assets.html_handling` | `auto-trailing-slash`, `force-trailing-slash`, `drop-trailing-slash`, `none` | `auto-trailing-slash` |
| `assets.not_found_handling` | `none`, `404-page`, `single-page-application` | `none` |
| `assets.run_worker_first` | `true`/`false` or array of patterns each starting with `/` or `!/`; `*` globs; "Negative patterns have precedence"; max 100 entries | `false` |

- `html_handling` behaviour ([HTML handling](https://developers.cloudflare.com/workers/static-assets/routing/advanced/html-handling/)): under `auto-trailing-slash` and `force-trailing-slash`, `/folder` → 307 to `/folder/`, `/folder/` served from `folder/index.html`, `/folder/index.html` → 307 to `/folder/`. They differ only for standalone `file.html`: auto serves it at `/file` and redirects `/file/` → `/file`; force serves it at `/file/`. This site writes every page as `dir/index.html` and canonicalises to `/dir/` (`builder/render.go`, `internal/vercel/output.go`), with Vercel issuing 308s — Workers issues **307s** instead, which is fine for GETs. Choose `force-trailing-slash` so a stray `*.html` can never become a slash-less canonical. **[inferred]** choice from the table.
- `not_found_handling = "404-page"` serves "the nearest `404.html` file with a `404 Not Found` status". ([SSG routing](https://developers.cloudflare.com/workers/static-assets/routing/static-site-generation/))
- Routing order: "Cloudflare will first attempt to serve static assets if one matches the incoming request. ... If an appropriate static asset is not found, Cloudflare will invoke your Worker script." With `run_worker_first` as an array, only matching routes invoke the Worker before assets. ([Worker script routing](https://developers.cloudflare.com/workers/static-assets/routing/worker-script/))
- `env.ASSETS.fetch(request)` takes a Request/URL/string, matches on pathname only, and "Requests made through this method have `html_handling` and `not_found_handling` configuration applied to them." ([Binding](https://developers.cloudflare.com/workers/static-assets/binding/))
- `_headers`: plain-text file in the asset directory, not served; up to 100 rules, 2,000 chars per line; "Custom headers defined in the `_headers` file are not applied to responses generated by your Worker code, even if the request URL matches a rule" (so the 410 handler must set its own headers). ([Workers `_headers`](https://developers.cloudflare.com/workers/static-assets/headers/))
- `_redirects`: `[source] [destination] [code?]`, default 302; "2,000 static redirects and 100 dynamic redirects"; 1,000 chars per rule; "Redirects defined in the `_redirects` file are not applied to requests served by your Worker code". Supported: 301/302/303/307/308 and proxy 200. `/home / 301` covers the existing `/home` → `/` rule. ([Workers `_redirects`](https://developers.cloudflare.com/workers/static-assets/redirects/))
- Limits ([Workers limits](https://developers.cloudflare.com/workers/platform/limits/)): files per Worker version **20,000 Free / 100,000 Paid** ("To use the increased file count limits in Wrangler, you must use version 4.34.0 or higher"), 25 MiB per file, Worker bundle 64 MiB uncompressed (free and paid). Free plan: 100,000 Worker requests/day (Error 1027 beyond).
- Billing ([Billing and limitations](https://developers.cloudflare.com/workers/static-assets/billing-and-limitations/), [Workers pricing](https://developers.cloudflare.com/workers/platform/pricing/)): "Requests to static assets are free and unlimited." "Requests to the Worker script (for example, in the case of SSR content) are billed according to Workers pricing." "There is no additional cost for storing Assets." On Free, if `run_worker_first` requests exceed the daily limit "these requests will receive a 429 (Too Many Requests) response instead of falling back to static asset serving." Workers Paid is $5/month with 10 M requests included, +$0.30/M.
- Custom domain: `routes = [{ pattern = "hastcode.com", custom_domain = true }]`; `workers_dev` "Defaults to `true` when the configuration has no `route` or `routes`, and `false` otherwise." ([Wrangler config](https://developers.cloudflare.com/workers/wrangler/configuration/))
- Previews on Workers: two mechanisms. **Version URLs** (`wrangler versions upload`, `<version-prefix>-<worker-name>.<subdomain>.workers.dev`, toggled by `preview_urls`; the config key kept its old name) and **Previews** (per branch: "Running `npx wrangler preview` creates or updates a Preview for your current branch under the same Worker"; URL `<preview-name>-<worker-name>.<subdomain>.workers.dev` or a custom-domain `<preview-name>.app.example.com`; "requires `Wrangler 4.135.0` or later"; limits 100 Previews/Worker Free, 500 Paid; `workers.dev` preview URLs carry `X-Robots-Tag: noindex`). ([Version URLs](https://developers.cloudflare.com/workers/configuration/previews/), [Previews](https://developers.cloudflare.com/workers/previews/))

### 5.3 Pages, for comparison

- Deploy: `npx wrangler pages deploy <OUTPUT_DIRECTORY> --project-name=<p> --branch=<b>`; the production branch is chosen at `pages project create --production-branch`; "To deploy assets to a preview environment, run ... `--branch=<BRANCH_NAME>`"; "For every branch you create, a branch alias will be available to you at `<BRANCH_NAME>.<PROJECT_NAME>.pages.dev`"; "If you are in a Git workspace, Wrangler will automatically pull the branch information for you. Otherwise, you will need to specify your branch" — in a tag checkout (detached HEAD) pass `--branch` explicitly (**[inferred]**). Direct Upload projects cannot set production-branch controls in the dashboard; use the Update Project API. ([Direct upload](https://developers.cloudflare.com/pages/get-started/direct-upload/), [Pages commands](https://developers.cloudflare.com/workers/wrangler/commands/pages/))
- Preview URLs: `<hash>.<project>.pages.dev` (permanent) and `<branch>.<project>.pages.dev` (latest commit of the branch); public by default, `X-Robots-Tag: noindex`. ([Preview deployments](https://developers.cloudflare.com/pages/configuration/preview-deployments/))
- Limits: files per site 20,000 Free / 100,000 Pro & Business (paid plans need `PAGES_WRANGLER_MAJOR_VERSION=4`), 25 MiB per file, 100 projects per account, 100/250/500 custom domains, builds time out at 20 min, "Requests to Pages functions count towards your quota for Workers plans". ([Pages limits](https://developers.cloudflare.com/pages/platform/limits/))
- `_headers` (100 rules, 2,000 chars; not applied to Functions responses) and `_redirects` (same table as Workers; not applied to Functions). ([Pages headers](https://developers.cloudflare.com/pages/configuration/headers/), [Pages redirects](https://developers.cloudflare.com/pages/configuration/redirects/))
- Pages Functions: file-based routing under `functions/`, `[param].js` one segment, `[[path]].js` "any depth of route after" the directory; `_routes.json` (`version`, `include`, `exclude`; `exclude` wins; ≤100 rules) controls which paths invoke Functions; `context.next()` "Passes the request through to the next Function or to the asset server"; `env.ASSETS.fetch()` fetches a static asset by its pretty path and "will run the header and redirect rules". Advanced mode `_worker.js` takes over all requests and "is required to make or forward requests to your project's static assets" via `env.ASSETS.fetch(request)`. ([Functions routing](https://developers.cloudflare.com/pages/functions/routing/), [Functions API](https://developers.cloudflare.com/pages/functions/api-reference/), [Advanced mode](https://developers.cloudflare.com/pages/functions/advanced-mode/))

### 5.4 `cloudflare/wrangler-action`

- Current major: **v4**; "The action now defaults to Wrangler v4"; Wrangler v1 and Global API key auth are gone. Inputs: `apiToken` (required; never printed), `accountId`, `command` (default is `deploy`; multi-line allowed), `workingDirectory`, `wranglerVersion` ("any version format NPM understands"), `packageManager`, `preCommands`/`postCommands`, `secrets` + `environment`, `gitHubToken` ("Enable this if you want to have GitHub Deployments triggered"). Outputs: `command-output`, `command-stderr`, `deployment-url`, `pages-deployment-alias-url` (Wrangler ≥ 3.78.0), and for Workers Previews `preview-url`, `preview-deployment-url`, `preview-name`, `preview-id` (Wrangler ≥ 4.136.3). ([wrangler-action README](https://github.com/cloudflare/wrangler-action))
- API token for Workers: Account API Tokens → Create Token → Custom → "Under **Permission policies**, open the **Custom** dropdown and select **Edit Cloudflare Workers**"; "We recommend scoping these down as much as possible". Secrets: `CLOUDFLARE_API_TOKEN`, `CLOUDFLARE_ACCOUNT_ID`. ([Workers GitHub Actions guide](https://developers.cloudflare.com/workers/ci-cd/external-cicd/github-actions/))
- API token for Pages: "Under **Permissions**, select *Account*, *Cloudflare Pages* and *Edit*". ([Pages direct upload + CI](https://developers.cloudflare.com/pages/how-to/use-direct-upload-with-continuous-integration/))
- Neither Cloudflare page lists the individual scopes the "Edit Cloudflare Workers" template expands to; treat the template as the documented answer.

### 5.5 Proposed `wrangler.toml` for this site

```toml
name = "hastcode"
main = "cf/worker.js"            # only runs for the 410 prefixes and asset misses
compatibility_date = "2026-10-22"
workers_dev = true               # keep the *.workers.dev URL for smoke tests
preview_urls = true

[[routes]]
pattern = "hastcode.com"
custom_domain = true

[assets]
directory = "./dist/cloudflare"  # the builder's Cloudflare output: site + _headers + _redirects + 404.html
binding = "ASSETS"
html_handling = "force-trailing-slash"
not_found_handling = "404-page"
run_worker_first = ["/video/*", "/community/*", "/group/*"]
```

Keys and defaults per [Wrangler config](https://developers.cloudflare.com/workers/wrangler/configuration/); whether `/video/*` also matches the bare `/video` is not specified on the binding page — add `"/video"` etc. if the old site had bare collection URLs (**[unsourced]**).

---

## 6. 410 Gone on Cloudflare

- Confirmed: `_redirects` on Workers and on Pages supports only "Redirects (301, 302, 303, 307, 308)" plus proxy 200; "Rewrites (other status codes)" is ❌ (the doc's own example of the unsupported form is `/blog/* /blog/404.html 404`). No 410. ([Workers `_redirects`](https://developers.cloudflare.com/workers/static-assets/redirects/), [Pages `_redirects`](https://developers.cloudflare.com/pages/configuration/redirects/))
- Workers pattern (ES module Worker; `fetch` handler, `env.ASSETS` binding per [Binding](https://developers.cloudflare.com/workers/static-assets/binding/); invoked for the three prefixes via `run_worker_first` and for any asset miss per [Worker script routing](https://developers.cloudflare.com/workers/static-assets/routing/worker-script/)):

```js
// cf/worker.js
const GONE = ["/video/", "/community/", "/group/"];

export default {
  async fetch(request, env) {
    const { pathname } = new URL(request.url);
    if (GONE.some((p) => pathname === p.slice(0, -1) || pathname.startsWith(p))) {
      // Body = the site's 404 page, status 410 (matches the Vercel route today).
      // A path that no asset matches goes through not_found_handling => 404.html.
      const page = await env.ASSETS.fetch(new URL("/__gone__/", request.url));
      return new Response(page.body, {
        status: 410,
        headers: {
          "content-type": page.headers.get("content-type") ?? "text/html; charset=utf-8",
          "cache-control": "public, max-age=86400",
          "x-robots-tag": "noindex",
        },
      });
    }
    // Not one of ours and no asset matched: let not_found_handling answer (404.html).
    return env.ASSETS.fetch(request);
  },
};
```

  Only `/video/*`, `/community/*`, `/group/*` and genuine 404s cost a Worker invocation; every real page and image is a free static-asset hit. `_headers` rules are not applied to this response, hence the explicit headers. ([Billing and limitations](https://developers.cloudflare.com/workers/static-assets/billing-and-limitations/), [Workers `_headers`](https://developers.cloudflare.com/workers/static-assets/headers/)) The `/__gone__/` trick relies on "`not_found_handling` configuration applied" to `ASSETS.fetch` ([Binding](https://developers.cloudflare.com/workers/static-assets/binding/)); fetching `/404.html` directly would be rewritten by `html_handling`.
- Pages pattern (if Pages were chosen): `functions/video/[[path]].js`, `functions/community/[[path]].js`, `functions/group/[[path]].js` each exporting `export function onRequest() { return new Response("Gone", { status: 410 }); }`, plus `_routes.json` `{"version":1,"include":["/video/*","/community/*","/group/*"],"exclude":[]}` so nothing else invokes Functions. Whether `[[path]]` matches the bare `/video` is not stated ("any depth of route after `/users/`"); add `functions/video/index.js` if needed. ([Functions routing](https://developers.cloudflare.com/pages/functions/routing/), [Functions API](https://developers.cloudflare.com/pages/functions/api-reference/))
- Cloudflare's default edge TTL for 410 is 3 minutes, but only for URLs with a cacheable extension; these paths have none, so the Worker runs per request. ([Default cache behavior](https://developers.cloudflare.com/cache/concepts/default-cache-behavior/))

---

## 7. GitHub Actions tag triggers

- Filters: "The `branches`, `branches-ignore`, `tags`, and `tags-ignore` keywords accept glob patterns that use characters like `*`, `**`, `+`, `?`, `!` and others". "If you define only `tags`/`tags-ignore` or only `branches`/`branches-ignore`, the workflow won't run for events affecting the undefined Git ref." "If you define neither ... the workflow will run for events affecting either branches or tags." You cannot combine `tags` with `tags-ignore` on one event; a `!` pattern must follow at least one positive pattern. Cheat sheet: `v2*` matches `v2`, `v2.0`, `v2.9`; `feature/*` does not cross `/`, `feature/**` does; "If you start a pattern with `*`, `[`, or `!`, you must enclose the pattern in quotes." `PUBCF-*` starts with a letter, so quotes are optional. ([Workflow syntax](https://docs.github.com/en/actions/writing-workflows/workflow-syntax-for-github-actions))
- Consequence: the existing Vercel workflow (`on: push: branches: [main]`) does **not** run when `PUBCF-2026.10.22` is pushed, and the Cloudflare workflow (`on: push: tags: ['PUBCF-*']`) does not run on branch pushes. ([Workflow syntax](https://docs.github.com/en/actions/writing-workflows/workflow-syntax-for-github-actions))
- Reading the tag: `github.ref` = "`refs/tags/<tag_name>`"; `github.ref_name` = "The short ref name of the branch or tag that triggered the workflow run. This value matches the branch or tag name shown on GitHub."; `github.ref_type` = "`branch` or `tag`"; `github.sha` for `push` = "Tip commit pushed to the ref." ([github context](https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/accessing-contextual-information-about-workflow-runs), [Default variables](https://docs.github.com/en/actions/reference/workflows-and-actions/variables), [Events: push](https://docs.github.com/en/actions/writing-workflows/choosing-when-your-workflow-runs/events-that-trigger-workflows#push))
- Injection hygiene: GitHub's guidance is to pass context values through an intermediate environment variable rather than inline in `run:` ("the preferred approach to handling untrusted input is to set the value of the expression to an intermediate environment variable"). The doc's example is a PR title, but tag names are also collaborator-controlled strings, so use `env: TAG: ${{ github.ref_name }}` and `"$TAG"` in shell. ([Security hardening](https://docs.github.com/en/actions/security-for-github-actions/security-guides/security-hardening-for-github-actions)) That ref names may contain shell-significant characters such as `$`, `;`, `|` is **[unsourced]** here (git's ref-name rules were not re-read today).
- Tag must point at a commit on `main`: `actions/checkout@v7` with `fetch-depth: 0` ("0 indicates all history for all branches and tags"; default 1), then `git merge-base --is-ancestor A B` which exits "with status 0 if true, or with status 1 if not". For annotated tags, peel with `git rev-parse "$GITHUB_REF^{commit}"` ("dereference the object at `<rev>` recursively until an object of type `<type>` is found") rather than trusting `GITHUB_SHA`, whose value for annotated tags is not specified by GitHub. ([actions/checkout](https://github.com/actions/checkout), [git merge-base](https://git-scm.com/docs/git-merge-base), [gitrevisions](https://git-scm.com/docs/gitrevisions))
- Concurrency: "there can be at most one running job or workflow in a concurrency group at any time"; by default one pending run waits and newer pending runs replace it (`queue: single`); `cancel-in-progress: true` also cancels the running one; `queue: max` allows up to 100 waiting but cannot be combined with `cancel-in-progress: true`; group names are case-insensitive. Use one group per deploy target so a Cloudflare deploy never cancels or races a Vercel deploy. ([Concurrency](https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/control-the-concurrency-of-workflows-and-jobs))

### 7.1 Workflow sketch

```yaml
name: Deploy to Cloudflare
on:
  push:
    tags: ['PUBCF-*']            # does not fire for branch pushes
concurrency:
  group: deploy-cloudflare       # the Vercel workflow uses group: deploy-vercel
  cancel-in-progress: false
permissions:
  contents: read
jobs:
  deploy:
    runs-on: ubuntu-latest
    timeout-minutes: 30
    steps:
      - uses: actions/checkout@v7
        with: { fetch-depth: 0 }
      - name: Tag must point at a commit on main
        env: { REF: "${{ github.ref }}" }
        run: |
          git fetch --no-tags origin main
          commit="$(git rev-parse "${REF}^{commit}")"
          git merge-base --is-ancestor "$commit" origin/main
      - uses: actions/setup-go@v5
        with: { go-version-file: go.mod }
      - name: Build site and upload images to R2
        env:
          R2_ACCOUNT_ID: ${{ secrets.CLOUDFLARE_ACCOUNT_ID }}
          R2_ACCESS_KEY_ID: ${{ secrets.R2_ACCESS_KEY_ID }}
          R2_SECRET_ACCESS_KEY: ${{ secrets.R2_SECRET_ACCESS_KEY }}
          IMAGE_BASE_URL: https://images.hastcode.com
          TAG: ${{ github.ref_name }}
        run: make build-cloudflare
      - uses: cloudflare/wrangler-action@v4
        with:
          apiToken: ${{ secrets.CLOUDFLARE_API_TOKEN }}
          accountId: ${{ secrets.CLOUDFLARE_ACCOUNT_ID }}
          command: deploy --message "$TAG"
```

`--message` is a documented `wrangler deploy` flag ("A descriptive message for this Worker version and deployment"); `--tag` also exists. ([Workers commands](https://developers.cloudflare.com/workers/wrangler/commands/workers/)) `actions/setup-go@v5` is the conventional action; its version was not verified today (**[unsourced]**).

---

## Open questions

1. **Does R2 accept aws-sdk-go-v2's default CRC32 on single-part `PutObject`?** The compatibility table lists CRC32 as `COMPOSITE`-only and no checksum headers under `PutObject`; Cloudflare's Go example sets nothing. Resolution: set `WhenRequired` (safe either way) and keep a one-object integration test. ([R2 S3 API](https://developers.cloudflare.com/r2/api/s3/api/))
2. **Is a `PutObject` that fails with `412 PreconditionFailed` billed as a Class A op?** Not stated on the pricing page. Using `HeadObject` first sidesteps it. ([R2 pricing](https://developers.cloudflare.com/r2/pricing/))
3. **Is batch `DeleteObjects` free like `DeleteObject`?** Only `DeleteObject`, `DeleteBucket`, `AbortMultipartUpload` are named. At ⌈N/1000⌉ calls it is immaterial, but confirm in the billing dashboard after the first prune. ([R2 pricing](https://developers.cloudflare.com/r2/pricing/))
4. **Does the custom-domain path emit the object's `Cache-Control` header?** Inferred, not quoted. Verify with `curl -I` and `cf-cache-status: HIT` on the second request. ([Public buckets](https://developers.cloudflare.com/r2/buckets/public-buckets/))
5. **Does `run_worker_first: ["/video/*"]` match bare `/video`?** Pattern semantics beyond "`*` for deep matching" and `!` negation are not documented. ([Binding](https://developers.cloudflare.com/workers/static-assets/binding/))
6. **307 vs 308 for trailing-slash redirects.** Workers `html_handling` issues 307; the Vercel config today issues 308. Both are fine for GET; decide whether the link-gate tests need to accept both. ([HTML handling](https://developers.cloudflare.com/workers/static-assets/routing/advanced/html-handling/))
7. **`GITHUB_SHA` for annotated tags.** GitHub says "Tip commit pushed to the ref"; whether that is the tag object or the peeled commit is not stated. The workflow above peels explicitly. ([Events: push](https://docs.github.com/en/actions/writing-workflows/choosing-when-your-workflow-runs/events-that-trigger-workflows#push))
8. **Absolute image URLs and Vercel.** Once HTML references `https://images.hastcode.com/...`, the Vercel deployment also depends on the R2 bucket and the Cloudflare zone being up; decide whether Vercel builds keep a copy of variants in `.vercel/output/static` as a fallback. Product decision, not a documentation question.
9. **Workers Free vs Paid.** Static assets are free and unlimited on both; the 410 Worker and 404 misses count against 100,000 requests/day on Free. If crawler traffic to the ~1,000 retired `/video/*` URLs plus 404 noise could approach that, budget the $5/month Paid plan. ([Workers pricing](https://developers.cloudflare.com/workers/platform/pricing/))

---

## Sources (all accessed 2026-10-22)

Cloudflare R2
- https://developers.cloudflare.com/r2/api/s3/api/
- https://developers.cloudflare.com/r2/api/s3/extensions/
- https://developers.cloudflare.com/r2/api/tokens/
- https://developers.cloudflare.com/r2/api/workers/workers-api-reference/
- https://developers.cloudflare.com/r2/examples/
- https://developers.cloudflare.com/r2/examples/aws/aws-sdk-go/
- https://developers.cloudflare.com/r2/examples/aws/aws-sdk-js-v3/
- https://developers.cloudflare.com/r2/pricing/
- https://developers.cloudflare.com/r2/platform/limits/
- https://developers.cloudflare.com/r2/buckets/public-buckets/
- https://developers.cloudflare.com/r2/buckets/cors/
- https://developers.cloudflare.com/r2/objects/multipart-objects/

Cloudflare cache and Images
- https://developers.cloudflare.com/cache/concepts/default-cache-behavior/
- https://developers.cloudflare.com/cache/concepts/cache-control/
- https://developers.cloudflare.com/images/pricing/
- https://developers.cloudflare.com/images/transform-images/
- https://developers.cloudflare.com/images/transform-images/transform-via-url/
- https://developers.cloudflare.com/images/transform-images/sources/

Cloudflare Workers (static assets, Wrangler, CI)
- https://developers.cloudflare.com/workers/static-assets/
- https://developers.cloudflare.com/workers/static-assets/binding/
- https://developers.cloudflare.com/workers/static-assets/routing/worker-script/
- https://developers.cloudflare.com/workers/static-assets/routing/static-site-generation/
- https://developers.cloudflare.com/workers/static-assets/routing/advanced/html-handling/
- https://developers.cloudflare.com/workers/static-assets/headers/
- https://developers.cloudflare.com/workers/static-assets/redirects/
- https://developers.cloudflare.com/workers/static-assets/billing-and-limitations/
- https://developers.cloudflare.com/workers/static-assets/compatibility-matrix/
- https://developers.cloudflare.com/workers/static-assets/migration-guides/migrate-from-pages/
- https://developers.cloudflare.com/workers/platform/limits/
- https://developers.cloudflare.com/workers/platform/pricing/
- https://developers.cloudflare.com/workers/wrangler/configuration/
- https://developers.cloudflare.com/workers/wrangler/commands/workers/
- https://developers.cloudflare.com/workers/wrangler/commands/pages/
- https://developers.cloudflare.com/workers/configuration/previews/
- https://developers.cloudflare.com/workers/previews/
- https://developers.cloudflare.com/workers/ci-cd/external-cicd/github-actions/
- https://github.com/cloudflare/wrangler-action

Cloudflare Pages
- https://developers.cloudflare.com/pages/
- https://developers.cloudflare.com/pages/get-started/direct-upload/
- https://developers.cloudflare.com/pages/how-to/use-direct-upload-with-continuous-integration/
- https://developers.cloudflare.com/pages/configuration/preview-deployments/
- https://developers.cloudflare.com/pages/configuration/headers/
- https://developers.cloudflare.com/pages/configuration/redirects/
- https://developers.cloudflare.com/pages/platform/limits/
- https://developers.cloudflare.com/pages/functions/routing/
- https://developers.cloudflare.com/pages/functions/api-reference/
- https://developers.cloudflare.com/pages/functions/advanced-mode/

AWS SDK for Go v2 and S3 API
- https://github.com/aws/aws-sdk-go-v2/discussions/2960
- https://docs.aws.amazon.com/sdkref/latest/guide/feature-dataintegrity.html
- https://github.com/aws/aws-sdk-go-v2/blob/main/service/s3/api_op_PutObject.go
- https://github.com/aws/aws-sdk-go-v2/blob/main/service/s3/types/errors.go
- https://github.com/aws/aws-sdk-go-v2/blob/main/aws/checksum.go
- https://github.com/aws/aws-sdk-go-v2/blob/main/config/load_options.go
- https://github.com/aws/aws-sdk-go-v2/blob/main/config/env_config.go
- https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/s3/manager
- https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3
- https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListObjectsV2.html
- https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObjects.html
- https://github.com/minio/minio-go

GitHub Actions and git
- https://docs.github.com/en/actions/writing-workflows/workflow-syntax-for-github-actions
- https://docs.github.com/en/actions/writing-workflows/choosing-when-your-workflow-runs/events-that-trigger-workflows
- https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/accessing-contextual-information-about-workflow-runs
- https://docs.github.com/en/actions/reference/workflows-and-actions/variables
- https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/control-the-concurrency-of-workflows-and-jobs
- https://docs.github.com/en/actions/security-for-github-actions/security-guides/security-hardening-for-github-actions
- https://github.com/actions/checkout
- https://git-scm.com/docs/git-merge-base
- https://git-scm.com/docs/gitrevisions

Web standards
- https://html.spec.whatwg.org/multipage/urls-and-fetching.html#cors-settings-attributes
- https://www.w3.org/TR/css-fonts-4/#font-fetching-requirements

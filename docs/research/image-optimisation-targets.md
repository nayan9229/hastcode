# Hero-image optimisation targets: codec settings, width set, fallbacks, metadata, social crops, LCP

Research for hastcode.com (Go static builder; heroes resized with `golang.org/x/image/draw` CatmullRom to 480/800/1200/1600 and encoded with `gen2brain/avif` (Q60, speed 8, 4:2:0), `gen2brain/webp` (Q80, method 4) and stdlib `image/jpeg` (Q82); three 1200-px JPEG crops for OG/schema.org; `<picture sizes="(min-width: 48rem) 46rem, 100vw">`; LCP element; target LCP <= 1.5 s on mid-range Android over 4G).
Researched 2026-10-22. Primary-source quotes were pulled from the live pages/source files on that date. Secondary sources are marked **[secondary]**, opinions **[opinion]**, and numbers from the local experiment (see section 0) **[local]**. Anything not backed by a link is marked **[unsourced]**.

---

## TL;DR

- **AVIF Q60 / speed 8 / 4:2:0 is a sound default for photos; speed 6 buys nothing on photos and speed 10 costs 2-4 %.** [local] At Q60, speed 6 vs 8 changed bytes by -3 %..+1 % on three photos while taking 3.2-3.9x longer (800 px: 0.8 s vs 0.23 s; 1600 px: 2.6 s vs 0.75 s single-threaded WASM on an M5); speed 10 was 2.5-4x faster than 8 for +2-5 % bytes. libavif maps speeds 7-10 to libaom's *realtime* usage mode and 0-6 to *good-quality* mode, so 6 and 8 are different encoder paths, which is why the step is not linear. At Q60 the AVIF is 46-58 % of the stdlib-JPEG-Q82 bytes (800 px: 32-65 KB vs 65-115 KB). Published matched-quality studies (DSSIM, ~600 photos) put AVIF at 50 % of MozJPEG and WebP at 69 %.
- **Use 4:4:4 for illustrations/diagrams, 4:2:0 for photos.** [local] On a synthetic text+line-art diagram AVIF 4:4:4 cost +0.8-1.2 % bytes over 4:2:0 and raised chroma SSIM from 0.988-0.993 to 0.999; for jpegli 4:4:4 costs +20 %. Lossy WebP cannot do 4:4:4 at all (libwebp: colorspace "should be YUV420 for now"), and `gen2brain/webp` does not expose libwebp's `use_sharp_yuv`. On flat art, AVIF speed 6 did matter (-11 % bytes vs speed 8), so use speed 6 + 4:4:4 for illustrations.
- **The 1200-px variant, not the 800-px one, is the mobile LCP image for most modern phones** (DPR 2.625-3 at 360-430 CSS px -> 945-1290 device px). [local] AVIF Q60 at 1200 px is 66-135 KB; Q55 is 54-110 KB. On Lighthouse's "Slow 4G" (1.6 Mbps, 150 ms RTT) 100 KB is ~0.5 s of pure transfer, so **drop to Q55 for widths >= 1200** [opinion] or verify field LCP first.
- **Width set: add a ~1000 step and consider 1472 instead of 1600.** Chrome picks a candidate by the geometric mean of adjacent densities, WebKit picks the first candidate whose density >= DPR. With 480/800/1200/1600 and a 736-px slot: 1.25x laptops (1536x864, 10 % of desktops) get 800 upscaled in Chrome and 1200 in Safari; 414-px 2x iPhones (13 % of mobiles) get 800 upscaled in Chrome and 1200 in Safari. A 1000-px candidate serves both exactly (every iOS browser uses WebKit's rule). 1600 is only ever chosen for the 2x desktop column (needs 1472) and for 3x iPhones >= 402 CSS px in Safari; 1472 satisfies both (-15 % pixels). Lighthouse itself caps the DPR it expects at 2.
- **Stdlib `image/jpeg` is baseline-only, 4:2:0-only, writes no metadata, and reads none.** The encoder emits SOI, DQT, SOF0, DHT, SOS, EOI (no JFIF APP0, EXIF or ICC); the decoder ignores APP1/APP2, so EXIF orientation is not applied (golang/go#4341, open since 2012) and ICC profiles are dropped (golang/go#23936, open). `gen2brain/jpegli` (CGO-free, WASM) adds progressive, 4:4:4/4:2:2, Huffman optimisation and adaptive quantisation; at matched luma-SSIM it ranged from 11 % smaller to 15 % larger than stdlib [local] (Google's "35 % at high quality" claim is from a crowd-sourced rater study; the gain is perceptual and invisible to SSIM). Gotcha: passing an `EncodingOptions` struct turns OptimizeCoding/AdaptiveQuantization **off** and selects 4:4:4 unless you set them explicitly.
- **JPEG stays; WebP is optional.** caniuse (StatCounter Sep 2026): AVIF 96.35 %, WebP 97.18 % (+0.1 partial), so WebP adds ~0.9 points of reach over AVIF, JPEG ~2.8 more. `<picture>` requires an `<img src>` anyway (MDN), Google asks for a `src` fallback for crawlers, and OG/social scrapers only document JPEG/PNG (X: also WebP/GIF **[secondary]**; Google Images: also WebP/AVIF). JPEG XL shipped on by default in Chrome 155 (2026-10-06) and Firefox 158; caniuse 60 % -> candidate second format for 2027, not now.
- **Resize: keep CatmullRom.** x/image/draw ships only NearestNeighbor, ApproxBiLinear, BiLinear and CatmullRom; kernels widen their support by the shrink factor so CatmullRom is alias-free when downscaling, ApproxBiLinear is not. [local] 3072->800: ApproxBiLinear 2.7 ms (visibly aliased), BiLinear 26 ms (soft), CatmullRom 48 ms, a custom Lanczos3 `draw.Kernel{Support: 3}` 69 ms and marginally sharper (SSIM 0.997 vs CatmullRom).
- **Outputs are implicitly sRGB (AVIF CICP unspecified -> Chrome assumes BT.709/sRGB/BT.601; WebP without ICCP -> "sRGB SHOULD be assumed"; untagged JPEG -> WebKit assumes sRGB).** That is correct only if sources are sRGB; Display P3 / Adobe RGB sources will desaturate. Enforce sRGB inputs or convert before resizing (mandykoh/prism can read the ICC description and convert P3/Adobe RGB -> sRGB in pure Go). Also: never decode WebP sources with x/image/webp or gen2brain/webp - both return libwebp's limited-range planes as full-range `image.YCbCr` (golang/go#40173, open): [local] 247 -> 228, 5 -> 17.
- **Social/structured crops:** og:image 1200x630 (Facebook: ">= 1200 x 630", "as close to 1.91:1 as possible", <= 8 MB, give `og:image:width/height`), LinkedIn >= 1200x627 1.91:1 <= 5 MB, X 2:1 min 300x157 max 4096x4096 < 5 MB **[secondary]**; Google Discover >= 1200 px wide, 16x9 recommended, `max-image-preview:large`; Article structured data 16x9 + 4x3 + 1x1, >= 50 K pixels, Google-Images-supported format. Keep them JPEG.
- **LQIP cannot improve LCP:** Chrome 112 (back-ported to 109) ignores images < 0.05 bits per displayed pixel as LCP candidates, so a blurred placeholder never counts and may even push LCP later. `fetchpriority="high"` (Google Flights: LCP 2.6 s -> 1.9 s) and not lazy-loading are the levers; preload the AVIF only if the hero markup sits late in the HTML. Cheapest placeholder: a solid dominant colour (7 bytes) or an inline 16-24-px WebP as `background-image` (136-220 base64 chars) [local]; BlurHash needs a 1.6 KB (gzip) JS decoder. All are < 1 % of the 50 KB page budget.

---

## Recommended settings

| Item | Current | Recommended | Evidence |
|---|---|---|---|
| AVIF, photos | Q60, speed 8, 4:2:0 | **Q60 / speed 8 / 4:2:0 for widths <= 800; Q55 / speed 8 / 4:2:0 for >= 1200** [opinion on the Q55 split] | Section 1.3 [local]; LCP byte budget section 2.3 |
| AVIF, illustrations / diagrams / AI line art | same as photos | **Q60 / speed 6 / 4:4:4** (set `ChromaSubsampling: image.YCbCrSubsampleRatio444` explicitly) | Section 1.6 [local]; avifenc itself defaults `--yuv auto` to 444 for non-JPEG input |
| WebP | Q80, method 4 | **Keep Q80 / method 4, or drop WebP entirely** if analytics show Safari <= 16.3 under ~0.5 % | Section 1.4 [local], section 4 (caniuse) |
| JPEG fallback | stdlib Q82 | **Keep stdlib Q82** (fallback-only cohort ~3 %). If bytes matter, trial `gen2brain/jpegli` Q85, 4:2:0, `ProgressiveLevel: 2`, `OptimizeCoding: true`, `AdaptiveQuantization: true` and judge it with SSIMULACRA 2 or by eye (break-even on luma-SSIM; the claimed gain is perceptual) | Section 3 [local] |
| Width set | 480 / 800 / 1200 / 1600 | **480 / 800 / 1000 / 1200 / 1472** (1600 is fine if you prefer round numbers; the 2x desktop column needs 1472) | Section 2 |
| Resize filter | CatmullRom | **CatmullRom** (optional custom Lanczos3 kernel: +40 % time, marginal sharpness). Never ApproxBiLinear for final output | Section 5 [local] |
| Metadata | nothing written | Fine as is. **Add input checks:** apply EXIF orientation before resizing; require sRGB (convert P3/Adobe RGB or fail the build); reject WebP sources | Section 6 |
| OG / schema crops | 3 x 1200-px JPEG (16:9, 4:3, 1:1) | Keep JPEG. Add **1200x630 for og:image** (or accept FB's ~7 % crop of the 16:9 one) and emit `og:image:width`, `og:image:height`, `og:image:alt`, `og:image:type` | Section 7 |
| LCP hints | `fetchpriority="high"`, `decoding="async"`, responsive preload already in `partial_hero.templ` | Keep `fetchpriority="high"`; never `loading="lazy"`; preload only the AVIF candidate (`type="image/avif"`); `width`/`height` on `<img>` | Section 8 |
| Placeholder | none | **Solid dominant-colour background** (or <= 0.25 KB data-URI WebP). No JS BlurHash | Section 8 [local] |

---

## 0. Local experiment (method and caveats) [local]

Scratch harness: `/private/tmp/claude-502/-Users-nayanmevada-Desktop-Genuin-demos-go-articles/6a14c6a9-3136-4351-a0c8-500378c90d6f/scratchpad/imgopt/exp/` (`main.go`, `visual.go`, `webpfix.go`, `codecs.log`, `webpfix.log`, `illo.log`, `out/*.png` crops). Not part of the repo.

- Apple M5, Go 1.26.4, `gen2brain/avif v0.6.0`, `gen2brain/webp v0.6.4`, `gen2brain/jpegli v0.4.2`, `golang.org/x/image v0.47.0`, `wazero v1.12.0`, built with `-tags nodynamic` (pure WASM path, i.e. what CI without native libavif/libwebp gets). The AVIF encoder is single-threaded (`lib/avif.c`: `encoder->maxThreads = 1`), so all times are one core.
- Sources: Wikimedia Commons "Hopetoun falls" (3072x2048, CC BY-SA 3.0; high-frequency foliage), "Tokyo Tower at night" (2592x1944, public domain; dark, noisy), "Lion waiting in Namibia" (1280x960, CC BY 2.0), plus a synthetic 2400x1350 architecture diagram (flat fills, 2-3-px lines, 22-64-px Go Regular text) rendered with `x/image/font`.
- Pipeline identical to the site's: `draw.CatmullRom` resize to 480/800/1200/1600, then encode; decode the result and compare with the resized source.
- Metric: single-scale SSIM (11-tap Gaussian, sigma 1.5) on luma and separately on Cb/Cr. **Plain SSIM is not a perceptual metric**: it rewards JPEG for preserving sensor noise and punishes the denoising/deblocking that AVIF and WebP do (visible in the night photo crops), and it says nothing about ringing or banding. It is used here for *within-codec* deltas (speed and quality steps) and for sanity; the cross-codec "equal quality" numbers come from the DSSIM / SSIMULACRA 2 literature in 1.1. Nothing here is Butteraugli or SSIMULACRA 2 (no CGO-free Go implementation; see Open questions).
- The first pass of WebP SSIM numbers was wrong because both Go WebP decoders return limited-range planes as full-range `image.YCbCr` (section 6.4). All WebP numbers below use a range-corrected decode (`webpfix.go`), verified against macOS ImageIO.

## 1. Codec quality settings

### 1.1 Published matched-quality measurements

- **DSSIM study, ~600 photos x 6 widths (96-1080 px), target DSSIM 0.0025, binary search per image; MozJPEG 3.3.1 vs libwebp 1.1.0 (sharp YUV, slowest method) vs cavif/libavif (4:4:4):** median reduction vs JPEG: **WebP 31.5 %, AVIF 50.3 %**; AVIF beat JPEG and WebP on every image; WebP needed a *higher* nominal quality than JPEG to reach the same DSSIM (mean WebP 93 vs JPEG 84). (ctrl.blog, 2020-07-12, updated 2021-08-12) **[secondary, measured]**
- **SSIMULACRA 2 Pareto fronts (Cloudinary, 2024-02-28; libjpeg-turbo 2.1.5.1, mozjpeg 4.1.5, jpegli/libjxl 0.10, libavif 1.0.3 + aom 3.8.1, libwebp 1.3.2, M3 Pro, 8 threads):** at medium quality (SSIMULACRA 2 ~60) AVIF has better density than JPEG and WebP "at the cost of slower encoding"; at visually lossless quality libavif at default speed is ~20 % smaller than libjpeg-turbo but "takes an order of magnitude longer to encode"; lossy WebP "cannot reach visually lossless quality ... because 4:2:0 chroma subsampling is mandatory"; jpegli "compresses better than WebP and even high-speed AVIF" at medium quality while producing standard JPEG. **[secondary, measured]**
- **Jake Archibald, "AVIF has landed" (2020), judged by eye at 50 % zoom:** F1 photo at acceptable quality: JPEG 74.4 KB (65.5 KB tuned MozJPEG), WebP 43 KB, AVIF 18.2 KB; flat illustration: AVIF 8.69 KB vs WebP lossy 8.65 KB vs PNG 16.3 KB; he disabled chroma subsampling for the illustrations "to keep the colours sharp" and notes lossy WebP has no such option (its "Sharp YUV" setting caused colour artefacts). Squoosh effort 10 "can take over 10 minutes" per image in WASM; libavif CLI is "orders of magnitude faster". **[secondary, by eye]**
- **Cloudflare (2020-10-03):** "images can reduce to half the size of JPEG and WebP"; no metric, no speed numbers; encoding "may take even a few seconds to create a single image". **[secondary, unmeasured]**
- **WebP docs:** "WebP images are about 30% smaller in size compared to PNG and JPEG images at equivalent visual quality" (no methodology on that page).
- Netflix's AVIF study (2020) could not be fetched (HTTP 403 twice); not cited.

### 1.2 Bytes at the site's current settings [local]

| Image @ width | JPEG std Q82 | WebP Q80 m4 | AVIF Q60 s8 | AVIF / JPEG |
|---|---|---|---|---|
| falls @480 | 42.3 KB | 35.6 KB | 24.7 KB | 58 % |
| falls @800 | 112.1 KB | 92.3 KB | 64.6 KB | 58 % |
| tokyo @800 | 64.7 KB | 39.6 KB | 31.5 KB | 49 % |
| lion @800 | 115.1 KB | 92.3 KB | 65.0 KB | 56 % |
| falls @1200 | 233.8 KB | 188.4 KB | 134.6 KB | 58 % |
| tokyo @1200 | 142.5 KB | 78.3 KB | 66.4 KB | 47 % |
| lion @1200 | 226.2 KB | 167.4 KB | 122.6 KB | 54 % |
| falls @1600 | 391.5 KB | 303.9 KB | 222.2 KB | 57 % |
| tokyo @1600 | 249.3 KB | 126.8 KB | 114.7 KB | 46 % |
| diagram @800 | 30.3 KB | 11.6 KB | 9.9 KB (4:2:0) / 10.0 KB (4:4:4) | 33 % |
| diagram @1600 | 84.4 KB | 27.8 KB | 20.8 KB / 21.0 KB | 25 % |

Luma-SSIM at these settings (800 px): JPEG Q82 0.954-0.969, WebP Q80 0.964-0.966, AVIF Q60 0.941-0.960. On this metric AVIF Q60 sits slightly below JPEG Q82 and needs Q65-70 to match it (then at 69-82 % of the JPEG bytes), while WebP Q80 is within 0.004 of it at 61-82 % of the bytes (matching it exactly needs Q85 on two of the three photos). That is the opposite ranking to the DSSIM/SSIMULACRA 2 studies above and is the expected failure mode of plain SSIM on denoising codecs (section 0). Visual check of 2x crops (`out/*-800-crops-*.png`) [opinion]: AVIF Q60 and WebP Q80 are both acceptable for a 2x-density hero; the night sky loses grain in both (WebP more), AVIF Q50 starts to smear fern texture.

### 1.3 AVIF quality and speed [local]

Quality ladder (speed 8, 4:2:0), KB / luma-SSIM:

| | Q45 | Q50 | Q55 | Q60 | Q65 | Q70 |
|---|---|---|---|---|---|---|
| falls @800 | 35.4 / 0.880 | 43.6 / 0.904 | 52.5 / 0.923 | 64.6 / 0.941 | 73.1 / 0.951 | 90.6 / 0.967 |
| tokyo @800 | 17.7 / 0.937 | 21.8 / 0.947 | 25.8 / 0.954 | 31.5 / 0.960 | 35.7 / 0.964 | 44.9 / 0.970 |
| lion @800 | 36.7 / 0.906 | 44.7 / 0.926 | 53.4 / 0.941 | 65.0 / 0.955 | 73.4 / 0.962 | 90.2 / 0.974 |
| falls @1200 | | 91.9 / 0.912 | 110.0 / 0.929 | 134.6 / 0.945 | | |
| tokyo @1200 | | 44.5 / 0.942 | 53.6 / 0.948 | 66.4 / 0.954 | | |
| lion @1200 | | 85.8 / 0.946 | 101.8 / 0.957 | 122.6 / 0.966 | | |
| falls @1600 | 124.6 / 0.902 | 153.1 / 0.922 | 182.1 / 0.938 | 222.2 / 0.951 | 250.2 / 0.958 | 308.3 / 0.970 |
| tokyo @1600 | 57.9 / 0.929 | 74.8 / 0.938 | 91.1 / 0.945 | 114.7 / 0.951 | 133.9 / 0.955 | 173.5 / 0.962 |

Each 5 quality points is roughly +15-20 % bytes. libavif maps quality to the AV1 quantizer as `quantizer = ((100 - quality) * 63 + 50) / 100` (Q60 -> qp 25, Q55 -> 28, Q50 -> 32, Q70 -> 19).

Speed (Q60, 4:2:0), KB / encode ms, one core, WASM:

| | speed 4 | speed 6 | speed 8 | speed 10 |
|---|---|---|---|---|
| falls @800 | 67.8 / 5279 | 65.1 / 803 | 64.6 / 228 | 66.1 / 61 |
| tokyo @800 | 30.5 / 3770 | 31.2 / 567 | 31.5 / 147 | 33.1 / 58 |
| lion @800 | 63.3 / 5693 | 62.9 / 826 | 65.0 / 261 | 66.1 / 66 |
| falls @1600 | 223.8 / 17963 | 220.7 / 2633 | 222.2 / 750 | 227.3 / 222 |
| tokyo @1600 | 104.6 / 13902 | 112.9 / 1853 | 114.7 / 509 | 119.0 / 202 |
| diagram @800 (Q60, 4:4:4) | | 8.9 / 387 | 10.0 / 96 | |

Luma-SSIM differences across speeds are <= 0.005 (e.g. falls @800 Q60: 0.952 / 0.946 / 0.941 / 0.939). Read: speed 6 -> 8 is free on photos (3.2-3.9x faster, -1..+3 % bytes), speed 8 -> 10 costs +2-5 % bytes for another 2.5-4x; speed 4 is 20-25x slower than 8 for -9..+5 % bytes at fixed Q (it buys ~0.01 luma-SSIM instead; only the night photo at 1600 got materially smaller). Why there is a cliff rather than a slope: libavif's `codec_aom.c` sets `aomUsage = AOM_USAGE_GOOD_QUALITY` and switches to `AOM_USAGE_REALTIME` when `aomCpuUsed >= 7` ("Speed 6: GoodQuality CpuUsed 6", "Speed 7: RealTime CpuUsed 7", "Speed 10: RealTime CpuUsed 9"). avifenc's own default is speed 6; `gen2brain/avif`'s `DefaultSpeed` is 10 and `DefaultQuality` 60.

Build-time estimate per article at speed 8 (4 widths, one core): AVIF ~1.0-1.5 s, WebP m4 ~0.7 s, JPEG ~0.05 s, plus 3 crops; at speed 6 AVIF becomes ~4-5 s. Encoders are safe to run per-variant in parallel (the repo already does).

### 1.4 WebP [local]

Range-corrected decode (see 6.4), KB / luma-SSIM:

| | Q70 m4 | Q75 m4 | Q80 m4 | Q85 m4 | Q90 m4 | Q80 m6 | Q80 m0 |
|---|---|---|---|---|---|---|---|
| falls @800 | 70.6 / 0.944 | 75.9 / 0.950 | 92.3 / 0.964 | 112.3 / 0.977 | 142.1 / 0.987 | 87.9 / 0.962 | 102.1 / - |
| tokyo @800 | 30.0 / 0.955 | 32.1 / 0.957 | 39.6 / 0.964 | 48.7 / 0.970 | 70.7 / 0.977 | 37.4 / 0.963 | 58.7 / - |
| lion @800 | 70.6 / 0.949 | 74.8 / 0.953 | 92.3 / 0.966 | 113.4 / 0.976 | 145.9 / 0.986 | 88.2 / 0.963 | 113.8 / - |
| falls @1600 | 236 / 0.947 | 252 / 0.952 | 304 / 0.964 | 381 / 0.975 | 485 / 0.985 | 290 / 0.961 | 332 / - |

Encode time at 800 px: m0 16-21 ms, m4 85-110 ms, m6 100-154 ms; at 1600 px: m0 57-73 ms, m4 320-400 ms, m6 376-525 ms. Method 6 saves 4-7 % bytes for +40-50 % time (cwebp docs: `-m` "Possible values range from 0 to 6. Default value is 4"). Method 0 (the repo's `Fast` mode) is +10-48 % bytes. `gen2brain/webp` passes only `quality`, `method`, `lossless`, `exact` to `WebPConfig`; `use_sharp_yuv` ("use sharp (and slow) RGB->YUV conversion") is not exposed, so the sharp-YUV chroma fix used in the DSSIM study is unavailable.

### 1.5 JPEG (stdlib) [local]

`image/jpeg` Q75 / 80 / 82 / 85 / 90 at 800 px (falls): 91.6 / 105.2 / 112.1 / 124.1 / 156.0 KB, luma-SSIM 0.939 / 0.949 / 0.954 / 0.961 / 0.975; ~6 ms per encode at 800 px, 21-25 ms at 1600 px. Lighthouse's "Efficiently encode images" audit re-encodes JPEGs at quality 85 and flags >= 4 KiB savings, so Q82 passes it. Q82 is a fallback for the ~3 % of users without AVIF/WebP (section 4); it is not on the LCP path for anyone else, so there is little to gain from tuning it. See section 3 for jpegli.

### 1.6 Chroma subsampling for illustrations / text [local]

Synthetic diagram (flat fills, hairlines, 22-64-px red/black/white text), KB / SSIM Y / Cb / Cr:

| @800 | bytes | Y | Cb | Cr |
|---|---|---|---|---|
| AVIF Q60 s8 4:2:0 | 9.9 KB | 0.998 | 0.993 | 0.988 |
| AVIF Q60 s8 4:4:4 | 10.0 KB | 0.998 | 0.999 | 0.999 |
| AVIF Q50 s8 4:4:4 | 8.6 KB | 0.997 | 0.998 | 0.997 |
| AVIF Q60 s6 4:4:4 | 8.9 KB | 0.999 | 0.999 | 0.999 |
| WebP Q80 m4 (4:2:0 only) | 11.6 KB | 0.998 | 0.991 | 0.986 |
| JPEG std Q82 (4:2:0 only) | 30.3 KB | 0.995 | 0.985 | 0.977 |
| jpegli Q82 4:2:0 | 25.5 KB | 0.995 | 0.987 | 0.986 |
| jpegli Q82 4:4:4 | 30.5 KB | 0.997 | 0.989 | 0.995 |

At 1600 px: AVIF 4:2:0 20.8 KB vs 4:4:4 21.0 KB (+0.8 %); jpegli 66.2 vs 79.3 KB (+20 %). So for AVIF, 4:4:4 on illustrations is effectively free and removes the chroma bleed on coloured text/lines (2x crops in `out/illo-800-text-crops-*.png`); the stdlib JPEG cannot do it; WebP cannot do it. For photos 4:2:0 is the right choice [opinion: not measured here on photos; the DSSIM study used 4:4:4 for AVIF and suggested per-image 4:2:2/4:4:4 switching could shrink files further, so the photo cost of 4:4:4 is unquantified]. Note avifenc's `--yuv auto` default is 4:4:4 for any non-JPEG input and 4:2:0 only when the input JPEG already is; `gen2brain/avif` defaults to 4:2:0 only when `Encode` is called with **no** options; if you pass an `Options` value its zero `ChromaSubsampling` is `image.YCbCrSubsampleRatio444` (the repo sets it explicitly, which is correct).

Also seen on the diagram: WebP shifts flat near-white (247,248,250) to (228,230,231) **only through the Go decoders** (section 6.4); the file itself decodes correctly in macOS ImageIO.

## 2. Width set and `sizes` math

### 2.1 How browsers pick a candidate

- Spec: each width descriptor is normalised to a density by "the width descriptor value divided by source size" (HTML spec "normalize the source densities"), then the UA chooses "in an implementation-defined manner" and may consider "the user's network conditions".
- **Chrome** (`html_srcset_parser.cc`, `SelectionLogic`): walks candidates in density order; when the next candidate's density >= DPR it computes `geometric_mean = sqrt(current_density * next_density)` and picks the *next* only if `device_scale_factor >= geometric_mean` (or DPR <= 1 and current < DPR); otherwise it keeps the smaller one. Then `AvoidDownloadIfHigherDensityResourceIsInCache` swaps in any larger candidate already in the memory cache.
- **WebKit** (`HTMLSrcsetParser.cpp`, `pickBestImageCandidate`): sorts by density and breaks at the first candidate with `density >= deviceScaleFactor`; largest if none. (Firefox not verified; assumed the same **[unsourced]**.)
- web.dev: "in most browsers, an `img` using the `srcset` or `sizes` syntax will never bother requesting a source with smaller dimensions than one that the user already has in their browser's cache"; MDN: the browser loads "the first image that is bigger than the chosen slot size".
- Lighthouse "Properly size images" computes the expected size as displayed size x `quantizeDpr(DPR)` where DPR >= 2 is quantised to 2 ("This is to relax the required size of the image"), with a 0.75 factor for large images; it fails an image when the rendered size is "at least 4KiB smaller than the actual size".

### 2.2 What the current set serves (slot 46 rem = 736 px at the 16-px default root; 100vw below 48 rem)

Candidate densities for the desktop slot: 480 -> 0.65x, 800 -> 1.09x, 1000 -> 1.36x, 1200 -> 1.63x, 1472 -> 2.00x, 1600 -> 2.17x. Cohort shares are StatCounter worldwide, September 2026 (desktop: 1920x1080 28.1 %, 1536x864 10.0 %, 1366x768 8.0 %, 1280x720 4.7 %, 2560x1440 4.0 %; mobile: 414x896 13.4 %, 360x800 7.5 %, 384x832 7.1 %, 390x844 6.1 %, 393x873 4.2 %, 360x780 3.4 %). StatCounter reports CSS-pixel sizes, so 1536x864 is a 1920x1080 panel at 125 % scaling **[inference]**.

| Cohort | Needed px | Chrome picks (now) | Safari picks (now) | With 1000 added | With 1472 instead of 1600 |
|---|---|---|---|---|---|
| Desktop 1x (1920x1080, 1366x768) | 736 | 800 | 800 | 800 | - |
| Laptop 1.25x (1536x864) | 920 | 800 (upscaled 15 %) | 1200 | **1000** (both) | - |
| Laptop 1.5x (1280x720 CSS) | 1104 | 1200 | 1200 | 1200 | - |
| 2x desktop (Mac, 2560x1440 CSS @2) | 1472 | 1600 | 1600 | 1600 | **1472** |
| Phone 360 @2x | 720 | 800 | 800 | 800 | - |
| Phone 360 @3x | 1080 | 1200 | 1200 | 1200 | - |
| Pixel 384 @2.625x | 1008 | 1200 | 1200 | 1200 (1000 is 2.60x < 2.625) | - |
| iPhone 390/393 @3x | 1170/1179 | 1200 | 1200 | 1200 | - |
| iPhone 414 @2x (all iOS browsers follow the WebKit column) | 828 | 800 (upscaled 3.5 %) | 1200 | Chrome 800, **Safari 1000** | - |
| iPhone 414/430 @3x | 1242/1290 | 1200 | 1600 | 1200 / 1600 | Safari **1472** |

Conclusions:
- **Is 1600 needed?** Only the 2x desktop column (1472) and Safari on >= 402-CSS-px 3x iPhones ever select it. Nothing needs more than 1472; renaming the top step to 1472 saves 15 % of pixels (and roughly bytes) on those two cohorts and still satisfies Lighthouse's DPR-2 cap exactly. Keeping 1600 costs little because that cohort is on desktop/Wi-Fi, not 4G [opinion].
- **Is a ~1000 step useful?** Yes for two sizeable cohorts: 125 %-scaled Windows laptops (10 % of desktops; today Chrome upscales 800 by 15 %, Safari/Firefox download 1200 = +44 % pixels over 1000) and 2x 414-px iPhones in Safari (today 1200, +44 % pixels). It does not help 3x phones. Cost: one more encode per format (+25-30 % build time for the AVIF ladder) and ~50 bytes of HTML per `<source>`.
- **480** serves only 1x phones <= 480 CSS px (rare) and nothing on desktop; keep it as the cheap `src` fallback / Lighthouse-friendly minimum, or drop it [opinion].
- `sizes` is in `rem`, so a non-16-px root font size or browser zoom changes the slot; browsers include zoom in DPR, so the math above still holds. `sizes="auto"` exists only for `loading="lazy"` images (MDN), so it is not an option for the hero.

### 2.3 LCP byte budget on 4G

- Lighthouse "Slow 4G" = 1.6 Mbps down / 750 Kbps up / 150 ms RTT, 4x CPU slowdown ("roughly the bottom 25% of 4G connections and top 25% of 3G connections"; CPU target "mid-tier mobile", Moto G4 class). WebPageTest "4G" = 9 Mbps / 170 ms RTT; "LTE" = 12 Mbps / 70 ms; "3G Fast" = 1.6 Mbps / 150 ms.
- At 1.6 Mbps, 100 KB is ~0.5 s of transfer before slow-start and sharing with the 50 KB HTML and CSS; at 9 Mbps it is ~0.09 s and RTTs dominate. web.dev's guideline split of LCP is TTFB ~40 %, resource load delay < 10 %, load duration ~40 %, render delay < 10 % ("guidelines, not strict rules"), i.e. ~0.6 s of a 1.5 s target for the image download.
- Which variant is the mobile LCP image: for 3x phones (iPhone 12+, Galaxy S at 3x, Pixel at 2.625x) it is the **1200-px** file; 800 px only serves 2x phones <= 400 CSS px. Measured 1200-px AVIF: Q60 66-135 KB, Q55 54-110 KB, Q50 45-92 KB [local]. On Slow 4G a 135 KB hero alone is ~0.7 s; that is why the table recommends Q55 for >= 1200 [opinion]. A per-image byte cap (e.g. re-encode at Q55/Q50 if the 1200-px AVIF exceeds ~100 KB) is the more robust rule than a fixed quality [opinion].
- HTML cost of the markup: `<picture>` with 3 formats x 5 candidates ~ 0.9 KB; the responsive preload ~0.3 KB; dropping WebP saves ~0.3 KB and one third of encode time.

## 3. JPEG encoders in Go (CGO-free)

- **stdlib `image/jpeg`**: "Encode writes the Image m to w in JPEG 4:2:0 baseline format"; `Options` has only `Quality` (1-100, default 75). Source (`writer.go`, `Encode`): writes SOI, DQT, SOF0, DHT, SOS, EOI - no progressive scans, no 4:4:4, no JFIF/EXIF/ICC segments, no Huffman optimisation. Decoder supports progressive input (added by Nigel Tao, ~10 % slower baseline decode).
- **`github.com/gen2brain/jpegli` v0.4.2 (2026-07-10)**: jpegli (libjxl's JPEG encoder) compiled to WASM under wazero, "CGo-free"; on arm64 the `wasm2go` transpiled backend is always used. `EncodingOptions`: `Quality` (default 75), `ChromaSubsampling` "444|440|422|420", `ProgressiveLevel` 0-2, `OptimizeCoding` and `AdaptiveQuantization` ("Enabled by default"), `StandardQuantTables`, `FancyDownsampling`, `DCTMethod`. **Gotcha (source, `Encode`)**: the "enabled by default" booleans are only true when you pass *no* options; with an `EncodingOptions` value they are copied verbatim, and the zero `ChromaSubsampling` is 4:4:4. Output decodes with stdlib `image/jpeg` (verified [local]). Google's claim: "35% compression ratio improvement at high quality compression settings", measured with crowd-sourced raters (ELO) on the Cloudinary Image Dataset '22 vs libjpeg-turbo and MozJPEG, XYB disabled, decode via libjpeg-turbo; speed "comparable" to libjpeg-turbo/MozJPEG; output "fully interoperable" with existing decoders. [local] At the same nominal Q, jpegli is ~20 % smaller than stdlib (Q82 @800: 88.7 vs 112.1 KB) but lower luma-SSIM (0.944 vs 0.954); at matched luma-SSIM it is 11 % smaller on the foliage photo (Q85, 99.6 KB = 89 %) and 6-15 % larger on the other two (they need Q90 = 106-115 %). `ProgressiveLevel: 2` saves another ~3 % (85.9 vs 88.7 KB); turning AQ+Huffman optimisation off (the gotcha) gives 106.9 KB. Encode 19-21 ms @800, 70-75 ms @1600 (3x stdlib; negligible). Conclusion: the rater-study gains are perceptual (adaptive quantisation, jpegli tables) and do not show up in SSIM; adopt it only if you can judge by eye or with SSIMULACRA 2, and only for the JPEG fallback/OG crops where ~10-20 % of a 150-250 KB file is worth ~30 KB per crop.
- **`github.com/gen2brain/jpegn`** (new, pure Go, SIMD on amd64/arm64/riscv64, "No CGo, no dependencies"): baseline/progressive, 4:4:4/4:4:0/4:2:2/4:2:0, optimised Huffman, jpegli-style adaptive dead zone; README claims: within 0.2 % of stdlib at the same settings, Huffman optimisation -4 % (Q50) to -20 % (Q95), progressive -9 % on flat art, AQ -10 %, all three -17 % at 1.6x encode time, encode 4.5-6.3x faster than stdlib, EXIF/ICC pass-through. Not benchmarked here; maturity unknown (no release history checked). Worth a look if the WASM dependency of jpegli is unwelcome.
- **`github.com/gen2brain/jpegxl`**: libjxl via WASM/wazero (CGO-free), decode and encode; the `wasm2go` build encodes only losslessly (zune-jpegxl). Relevant once JPEG XL reach is high enough (section 4).
- **`github.com/kolesa-team/go-webp`**: "C binding for Google libwebp", requires `brew install webp` / `libwebp-dev`, i.e. **CGO**, not an option for a CGO-free build. (`gen2brain/webp` is libwebp transpiled to pure Go with wasm2go; `gen2brain/avif` is libavif+aom+dav1d under wazero, with an optional `purego` path to a system shared library; `gen2brain/gav1d` is a pure-Go AV1/AVIF port of dav1d with an 8-bit, 4:2:0-only all-intra encoder.)

## 4. Should the JPEG fallback exist? Should WebP?

- caniuse, usage from StatCounter September 2026: **AVIF 96.35 %** (Chrome 85+, Edge 121+, Safari 16.4+ with 16.1-16.3 partial, Firefox 93+, Samsung Internet 14+); **WebP 97.18 % + 0.1 % partial** (Safari 16.0+ full, 14-15.6 partial, iOS 14+; not in IE, Safari <= 13.1, Firefox <= 64, Chrome <= 8). So the WebP-but-not-AVIF cohort is ~0.9 % (Safari 14-16.3, Edge 114-120, Firefox 77-92, Samsung < 14) and the JPEG-only remainder ~2.8 %.
- JPEG must stay regardless of share: MDN - inside `<picture>` "you must provide an `<img>` element, with `src` and `alt`, right before `</picture>`, otherwise no images will appear"; Google Images - "some browsers and crawlers don't understand these attributes, so always specify a fallback URL using the `src` attribute"; web.dev - "You don't need to specify a second `<source>` element with `type="image/jpeg"`", the inner `<img>` is the default. The OG/schema crops need a scraper-safe format anyway (section 7), so the JPEG encoder stays in the build either way.
- WebP is the only negotiable format: it costs ~0.1-0.4 s per variant of build time, storage, and ~0.3 KB of HTML, to serve < 1 % of users a file that is ~1.4x the AVIF. Decision rule [opinion]: keep WebP until analytics show Safari <= 16.3 share below ~0.5 %, then drop it. Counter-argument to check: AVIF software decode cost on low-end Android is higher than WebP's and lands in LCP's render-delay slice; no primary measurement found (Open questions).
- **JPEG XL** (for 2027 planning): Chrome blog, 2026-10-06: "Chrome is shipping decoding support for the JPEG XL (`.jxl`) image format starting from Chrome 155", "30-50% better compression than JPEG", "we recommend trying both AVIF and JPEG XL"; decoder is jxl-rs (Rust). caniuse: 60.19 % (43.26 % full: Chrome 155+, Firefox 158+, Chrome for Android 156+; 16.93 % partial: Safari 17+). Not yet a replacement for anything; revisit when it passes ~90 %.

## 5. Resize filter

- `golang.org/x/image/draw` (v0.47.0) provides `NearestNeighbor` ("very fast, but usually gives very low quality results"), `ApproxBiLinear` ("fast, but usually gives medium quality results"; bilinear when upscaling, "4 nearest neighbors when downscaling"), `BiLinear` ("slow, but usually gives high quality results") and `CatmullRom` ("very slow, but usually gives very high quality results"; cubic BC-spline B=0, C=0.5, Mitchell-Netravali). No Lanczos, no Box; `Kernel` is a public struct (`Support`, `At`) so a Lanczos3 is 6 lines. `scale.go` `newDistrib`: "When shrinking, broaden the effective kernel support so that we still visit every source pixel" (`halfWidth *= scale`), i.e. kernel scalers are area-aware when downscaling; `ApproxBiLinear` is not, hence aliasing.
- [local] 3072 -> 800 / 1600, min of 3 runs: ApproxBiLinear 2.7 / 10.7 ms, BiLinear 26.1 / 39.6 ms, CatmullRom 47.6 / 63.5 ms, Lanczos3 custom kernel 68.5 / 87.8 ms. Mean-gradient sharpness index at 800: 34.3 (ApproxBiLinear, aliasing), 25.9 (BiLinear), 29.5 (CatmullRom), 31.5 (Lanczos3); luma-SSIM of CatmullRom vs Lanczos3 output 0.997. `Kernel.NewScaler` (pre-computed weights) only pays off when source and destination sizes repeat, which they do not here (48.7 ms, no gain).
- Other encoders' documentation: ImageMagick's filter guide calls Lanczos "the most popular" windowed-sinc, Catrom "a reasonably sharp edge" that "some consider far too sharp" and prone to a slight negative halo, and recommends Box only for integer "binning". Pillow's filter table rates downscaling quality LANCZOS 4/5, BICUBIC 3/5, HAMMING 2/5, BILINEAR 1/5, BOX 1/5, with performance in the inverse order.
- Recommendation: keep CatmullRom; it is one notch below Lanczos3 in sharpness, alias-free, and 30 % faster. Do not use ApproxBiLinear for the final hero (the repo's `Fast` mode does not; it only lowers codec effort) [opinion].

## 6. Metadata and colour management

### 6.1 What the Go encoders write
- `image/jpeg`: SOI, DQT, SOF0, DHT, SOS, EOI only (source, `Encode`) - no APP0/JFIF, no EXIF, no ICC, no comment. Confirmed.
- `image/png`: IHDR, PLTE/tRNS (palette only), IDAT, IEND (source, `writer.go`) - no iCCP, sRGB, gAMA, eXIf.
- `gen2brain/webp` (`lib/webp.c`): `WebPConfigInit`, sets `quality`, `method`, `lossless`, `exact`, imports RGBA, encodes; writes no ICCP/EXIF/XMP chunks (all three are optional per the container spec, and the Go API has no way to pass them). cwebp's own CLI default is also `-metadata none`.
- `gen2brain/avif` (`lib/avif.c`): `avifImageCreate` + `avifImageRGBToYUV` + `avifEncoderCreate` with `quality`, `qualityAlpha`, `speed`; no ICC/EXIF/XMP; CICP left at libavif's `avifImageSetDefaults` values `AVIF_COLOR_PRIMARIES_UNSPECIFIED` / `AVIF_TRANSFER_CHARACTERISTICS_UNSPECIFIED` / `AVIF_MATRIX_COEFFICIENTS_UNSPECIFIED`, `yuvRange = AVIF_RANGE_FULL` (lossless mode switches to identity matrix). Chrome's `avif_image_decoder.cc` maps unspecified to BT.709 primaries, sRGB transfer, BT.601 matrix ("we must choose some defaults and these are the most reasonable defaults to choose. We also advocate that all AVIF decoders choose these defaults: av1-avif issue 84"), i.e. the file is treated as sRGB. (The libavif version bundled in v0.6.0 is not stated; defaults checked against libavif main.)

So stripping is already the behaviour; the hero pipeline never leaks EXIF (GPS, camera) and never ships an ICC profile. The cost is 6.2 and 6.3.

### 6.2 Colour: everything is assumed sRGB
- WebP container spec: "If this chunk [ICCP] is not present, sRGB SHOULD be assumed." WebKit: "If an image doesn't have a tagged profile, WebKit assumes that it is sRGB"; tagged images are colour-matched on iOS/macOS. Chrome AVIF: as above.
- Go's `image/jpeg` decoder processes only APP0 (JFIF) and APP14 (Adobe) segments; APP1 (EXIF) and APP2 (ICC) are skipped (`reader.go` marker switch). golang/go#23936 "image/jpeg: encoding with RGB profile causing loss of image saturation" (opened 2018-02-20, still open, NeedsInvestigation) is exactly this: the profile is dropped and the pixels are re-interpreted as sRGB. Impact: iPhone exports (Display P3), Adobe RGB camera/Photoshop output and wide-gamut AI-tool exports come out desaturated or hue-shifted in *all three* formats, since the resize happens on the raw decoded values.
- Mitigation options: (a) editorial rule "export sRGB" plus a build-time check that rejects sources whose APP2 ICC description is not sRGB; (b) convert in Go: `github.com/mandykoh/prism` (pure Go, v0.35.3 2024-06-27) reads the embedded ICC profile description from JPEG/PNG/WebP and converts between sRGB, Display P3, Adobe RGB and ProPhoto (named spaces only, no arbitrary-profile CMM). Either way the conversion must happen before resizing [opinion].

### 6.3 EXIF orientation is not applied
golang/go#4341 "image/jpeg: correct for EXIF orientation?" (2012, open). Phone photos with orientation 6/8 will be resized sideways. `gen2brain/avif` and `gen2brain/webp` only auto-rotate their *own* formats on decode (`AutoRotate`). Read the orientation tag (prism's metadata reader or an EXIF library **[unsourced choice]**) and rotate before `cover()`.

### 6.4 Do not decode WebP with Go
Both `golang.org/x/image/webp` (returns `*image.YCbCr`) and `gen2brain/webp` (returns `*image.NYCbCrA` from `MODE_YUVA`) hand libwebp's limited-range BT.601 planes to Go's full-range `YCbCr` conversion. golang/go#40173 "x/image/webp: WebP decode contrast issue" (2020-07-12, open; Nigel Tao, 2024-11-06: a proper fix "is a lot of work"). [local] Diagram background (247,248,250) decodes as (228,230,231) and (5,150,105) as (17,146,106) in both Go decoders; macOS ImageIO gives (246,248,250). The site registers only `image/jpeg` and `image/png` decoders (`internal/images/images.go`), so this cannot currently corrupt a hero, but it will bite any test or tool that round-trips the WebP output through Go (e.g. a SSIM check), and it would bite if WebP sources were ever accepted.

## 7. Open Graph, X, LinkedIn, Discover and schema.org crops

| Consumer | Size / ratio | Limits | Formats documented | Source |
|---|---|---|---|---|
| Facebook / Open Graph | "at least 1200 x 630 pixels for the best display on high resolution devices"; "as close to 1.91:1 aspect ratio as possible"; min 600x315 for the large layout; absolute min 200x200 | <= 8 MB; add `og:image:width`/`og:image:height` so the first share can "render the image immediately" | not stated | developers.facebook.com sharing/webmasters/images |
| ogp.me | `og:image:url`, `:secure_url`, `:type` (MIME), `:width`, `:height`, `:alt` ("If the page specifies an og:image it should specify og:image:alt") | - | example uses image/jpeg | ogp.me |
| X summary_large_image | 2:1; min 300x157; max 4096x4096 | < 5 MB | JPG, PNG, WEBP, GIF; no SVG | **[secondary]** The SEO Framework KB quoting X docs; developer.x.com returned HTTP 402 / client-rendered page |
| LinkedIn | "Minimum image dimensions: 1200 (w) x 627 (h) pixels"; "Recommended ratio: 1.91:1" | "Max file size: 5 MB" | JPG/PNG/GIF listed only for ads | linkedin.com/help answer a521928 |
| Google Discover | "At least 1200 px wide"; "16x9 aspect ratio" recommended; "avoid automatically applying an aspect ratio" | needs `max-image-preview:large` (or AMP) | not stated (Google Images rules apply) | developers.google.com google-discover |
| Google Article structured data | "16x9, 4x3, and 1x1" recommended; "minimum of 50K pixels when multiplying width and height"; "multiple high-resolution images" | crawlable and indexable | "BMP, GIF, JPEG, PNG, WebP, SVG, and AVIF" (Google Images) | developers.google.com structured-data/article, google-images |

Recommendations:
- Keep all share/structured crops as **JPEG** (stdlib Q80-82, or jpegli): JPEG/PNG is the only format every scraper documents; WebP is documented for Google and (secondary) X; AVIF only for Google. These files are fetched by bots, not on the user's LCP path, so bytes matter less than compatibility; a 1200x675 JPEG Q82 of a detailed photo is ~120-200 KB [local, scaled from the 1200x800 / 1200x900 rows: 142-234 KB], far under every limit.
- Four crops cover everyone: **1200x630** (og:image for Facebook/LinkedIn/X/Slack; 1.905:1), **1200x675** (16x9, first schema.org `image`, matches Discover's recommendation), **1200x900** (4x3), **1200x1200** (1x1). If you want to stay at three, use the 16x9 crop as og:image and accept a ~7 % crop on Facebook/LinkedIn and a ~11 % crop on X [opinion]. All exceed Discover's "> 300,000 px" (1200x630 = 756,000).
- Emit `og:image:width`, `og:image:height`, `og:image:type="image/jpeg"`, `og:image:alt`. Discover and Google Images also index the in-article AVIF hero (AVIF is a supported format), so the structured-data `image` array can list the hero's JPEG 1600 or 1472 as well.

## 8. LQIP, `fetchpriority` and preload for the LCP hero

Evidence:
- web.dev "Fetch Priority": images "start at Low priority" until layout confirms viewport presence; `fetchpriority="high"` lets them "start at High immediately and load much faster"; Google Flights "LCP improved from 2.6s to 1.9s"; support Chrome 102, Edge 102, Firefox 132, Safari 17.2; "if a preload for an LCP image is already among the first items in the `<head>`, then a `high` Fetch Priority might not improve LCP significantly" (the two are substitutes); offscreen carousel images boosted to `high` "load even if a `loading=lazy` attribute is added".
- web.dev "Optimize LCP": "Never lazy-load your LCP image"; "setting a high priority on more than one or two images makes priority setting unhelpful"; when the LCP `<img>`'s `src`/`srcset` "are present in the initial HTML markup" the preload scanner discovers it, preload is for images referenced from CSS/JS; "Our research shows that the load duration portion of LCP tends not to be a significant bottleneck for most sites".
- web.dev "Preload responsive images": `<link rel="preload" as="image" href=... imagesrcset="... 400w, ... 800w, ... 1600w" imagesizes="50vw" fetchpriority="high">`; responsive preload support Chrome 73, Edge 79, Firefox 78, Safari 17.2; preloading `<picture>` has "technical issues": with several `type` preloads "all the images will be preloaded", so preload only the preferred format (AVIF); demo: Speed Index -273 ms, Last Painted Hero -1.2 s for images that were otherwise discovered late.
- Chromium LCP changelog (Chrome 112, rolled out ~2023-04-06, back to 109): "The threshold is currently 0.05 bits of image data per displayed pixel"; images below it "will not be considered for LCP"; this may *raise* reported LCP "if a low-content placeholder was being used to reserve space for a more contentful image which loaded later". web.dev LCP: Chromium excludes "Placeholder images or other images with a low entropy", opacity-0 elements and full-viewport elements.
- csswizardry (2023-09-28, updated 2025-09-05): a deliberately >= 0.055-bpp "LQIP" (his masthead: 24 KB) that *is* the LCP, preloaded with `fetchpriority=high`; he calls it gaming the metric and warns against any placeholder "you wouldn't want a user to look at for a while". **[secondary]**
- Placeholder payloads: BlurHash "only 20-30 characters", recommended 4x3 components, decoded client-side (`blurhash` 2.0.5: 3,182 B minified / 1,592 B gzip); ThumbHash ~21 bytes "roughly the same size as a BlurHash", also JS-decoded; CSS-only LQIP (Lean Rada, 2025-03-30): one integer per image ("20 bits of information"), rendered with layered gradients, needs `round()`, `mod()`, `pow()`, `oklab()`, no JS. [local] Inline data-URI of the falls photo: 16 px WebP 102 B (136 base64 chars), 24 px 164 B (220), 32 px 228 B (304); AVIF 345-407 B; stdlib JPEG 641-761 B (fixed table overhead); PNG 553-1721 B.

Implications for a 50 KB page budget and the 0.05-bpp rule:
- A 24-px WebP placeholder displayed at 390x219 CSS px is 164 x 8 / 85,410 = 0.015 bpp -> never an LCP candidate; the 800-px AVIF Q60 hero is ~1.25 bpp. So any blur placeholder is purely cosmetic: it improves perceived loading and FCP, not LCP, and if it is implemented as a JS swap it can delay the real request. The hero partial already has `fetchpriority="high"`, `decoding="async"` and a responsive preload; add `width`/`height` (CLS) and make sure the preload lists only the AVIF `imagesrcset` [opinion].
- Cheapest acceptable placeholder: `style="background:#5b6a4c"` (dominant colour, ~25 bytes) or `background-image:url(data:image/webp;base64,...)` on the `<img>` (~0.2 KB, zero JS, no swap). BlurHash costs 1.6 KB gzip of JS plus canvas work on the main thread of the very device that is already CPU-throttled 4x; not worth it here [opinion].
- Preload: the hero `<picture>` is in the initial HTML, so the preload scanner finds it; the preload only buys time if the markup is late in the document (after nav/ads) or if the browser would otherwise start it Low. Keep it to one `<link>` (AVIF only, with `imagesrcset`/`imagesizes` identical to the `<source>`), ~0.3 KB.

## Open questions

1. **Perceptual verification.** All cross-codec quality claims here rest on the DSSIM/SSIMULACRA 2 literature and an eyeball check; the local SSIM cannot rank AVIF vs WebP vs JPEG. Run SSIMULACRA 2 (libjxl `ssimulacra2` or the `ssimulacra2_rs` CLI; nothing CGO-free in Go) on 20-30 real hero sources at AVIF Q50/55/60/65, WebP Q75/80, JPEG Q80/82 before freezing numbers. Target bands used in that literature: SSIMULACRA 2 >= 70 "high quality", >= 85 "very high" **[unsourced thresholds]**.
2. **AVIF decode cost on mid-range Android** (LCP render delay, 4x CPU throttle): no primary measurement found; worth a WebPageTest run on a real Moto G-class device with AVIF vs WebP vs JPEG heroes.
3. **Drop WebP?** Needs the site's own Safari <= 16.3 / old-Edge share.
4. **JPEG XL** as second format in 2027 (Chrome 155+, Firefox 158+, Safari partial; 60 % today). `gen2brain/jpegxl` exists; lossy encode is only available on the wazero path.
5. **Firefox srcset selection** rule not verified (hg.mozilla.org raw file 404); assumed WebKit-like.
6. **X card requirements** are from a secondary source (developer.x.com returns 402 / client-rendered); verify in X's Card Validator.
7. **Netflix AVIF study** (4:4:4 vs 4:2:0 on photos, VMAF/SSIM) not fetchable (403); the photo cost of 4:4:4 AVIF is unmeasured.
8. **libavif version** bundled in `gen2brain/avif` v0.6.0 is not stated; the CICP defaults and speed->usage mapping were checked against libavif `main`.
9. **Wide-gamut sources**: how often do editors supply Display P3 / Adobe RGB files? Decide between "reject" and "convert with prism".
10. **1472 vs 1600 and the 1000 step** change cache keys and URLs for every hero; confirm the cost of re-encoding the back catalogue (AVIF speed 8 ~1-1.5 s per article per core) before switching.

## Sources (accessed 2026-10-22)

**Go packages and sources**
- https://pkg.go.dev/github.com/gen2brain/avif - Options, DefaultQuality=60, DefaultSpeed=10; module cache `avif.go`, `lib/avif.c` (encoder defaults, maxThreads=1, CICP)
- https://github.com/gen2brain/avif (README: libavif + dav1d + aom via wazero, CGo-free; purego dynamic-library path; `nodynamic`, `wasm2go` tags)
- https://pkg.go.dev/github.com/gen2brain/webp - Options, DefaultQuality=75, DefaultMethod=4; module cache `webp.go`, `lib/webp.c`
- https://pkg.go.dev/github.com/gen2brain/jpegli - EncodingOptions (v0.4.2, 2026-07-10); module cache `jpegli.go` (`Encode` option handling)
- https://github.com/gen2brain/jpegli (README), https://github.com/gen2brain/jpegn (README), https://github.com/gen2brain/jpegxl (README), https://github.com/gen2brain/gav1d (README)
- https://github.com/kolesa-team/go-webp (README: C binding, libwebp required)
- https://pkg.go.dev/golang.org/x/image/draw ; `golang.org/x/image@v0.47.0/draw/scale.go` (`newDistrib`, `CatmullRom`, `Kernel`)
- https://pkg.go.dev/image/jpeg ; Go 1.26.4 `src/image/jpeg/writer.go` (`Encode`), `reader.go` (marker switch), `src/image/png/writer.go`
- https://github.com/golang/go/issues/4341 - EXIF orientation (open, 2012)
- https://github.com/golang/go/issues/23936 - ICC profile dropped, saturation loss (open, 2018)
- https://github.com/golang/go/issues/40173 - x/image/webp limited-range decode (open, 2020)
- https://pkg.go.dev/github.com/mandykoh/prism - pure-Go colour conversion (v0.35.3, 2024-06-27)

**Encoder documentation and source**
- https://github.com/AOMediaCodec/libavif/blob/main/doc/avifenc.1.md - `--speed` 0..10 default 6, `--qcolor`, `--yuv auto` (444 default for non-JPEG input), `--cicp`, `--icc`, `--sharpyuv`
- https://github.com/AOMediaCodec/libavif/blob/main/src/codec_aom.c - speed -> AOM usage/cpu-used table, `aomQualityToQuantizer`
- https://github.com/AOMediaCodec/libavif/blob/main/src/avif.c - `avifImageSetDefaults` CICP unspecified, full range
- https://chromium.googlesource.com/chromium/src/+/main/third_party/blink/renderer/platform/image-decoders/avif/avif_image_decoder.cc - unspecified CICP -> BT.709 / sRGB / BT.601
- https://developers.google.com/speed/webp/docs/cwebp - `-q` default 75, `-m` default 4, `-sharp_yuv`, `-metadata none`
- https://github.com/webmproject/libwebp/blob/main/src/webp/encode.h - `method`, `use_sharp_yuv`, "colorspace should be YUV420 for now"
- https://developers.google.com/speed/webp/docs/riff_container - "If this chunk is not present, sRGB SHOULD be assumed"; EXIF/XMP optional
- https://developers.google.com/speed/webp/docs/compression - "about 30% smaller ... at equivalent visual quality"
- https://opensource.googleblog.com/2024/04/introducing-jpegli-new-jpeg-coding-library.html - 35 % claim, rater methodology
- https://developer.chrome.com/blog/jpeg-xl-in-chrome - Chrome 155 ships JPEG XL (2026-10-06)

**Measurements (secondary, labelled)**
- https://www.ctrl.blog/entry/webp-avif-comparison.html - DSSIM study (2020-07-12; updated 2021-08-12)
- https://cloudinary.com/blog/jpeg-xl-and-the-pareto-front - SSIMULACRA 2 Pareto fronts (2024-02-28)
- https://jakearchibald.com/2020/avif-has-landed/ - per-image sizes, chroma notes, effort/time
- https://blog.cloudflare.com/generate-avif-images-with-image-resizing/ (2020-10-03)
- https://netflixtechblog.com/avif-for-next-generation-image-coding-b1d75675fe4 - not fetchable (403)

**Browser selection, support, metrics**
- https://html.spec.whatwg.org/multipage/images.html#selecting-an-image-source ; https://html.spec.whatwg.org/multipage/images.html#normalise-the-source-densities
- https://chromium.googlesource.com/chromium/src/+/main/third_party/blink/renderer/core/html/parser/html_srcset_parser.cc - `SelectionLogic`, `AvoidDownloadIfHigherDensityResourceIsInCache`
- https://github.com/WebKit/WebKit/blob/main/Source/WebCore/html/parser/HTMLSrcsetParser.cpp - `pickBestImageCandidate`
- https://web.dev/learn/images/descriptive ; https://web.dev/articles/serve-responsive-images ; https://web.dev/learn/images/prescriptive ; https://web.dev/learn/images/avif ; https://web.dev/learn/images/webp
- https://developer.mozilla.org/en-US/docs/Web/HTML/Guides/Responsive_images ; https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/img (`fetchpriority`, `decoding`, `loading`, `sizes="auto"`)
- https://github.com/GoogleChrome/lighthouse/blob/main/core/audits/image-size-responsive.js - `quantizeDpr` caps at 2
- https://developer.chrome.com/docs/lighthouse/performance/uses-responsive-images ; https://developer.chrome.com/docs/lighthouse/performance/uses-optimized-images
- https://github.com/GoogleChrome/lighthouse/blob/main/docs/throttling.md - Slow 4G 1.6 Mbps / 150 ms, 4x CPU
- https://github.com/catchpoint/WebPageTest/blob/master/www/settings/connectivity.ini.sample - 4G 9 Mbps / 170 ms, LTE 12 Mbps / 70 ms
- https://caniuse.com/avif ; https://caniuse.com/webp ; https://caniuse.com/jpegxl (usage: StatCounter September 2026)
- https://gs.statcounter.com/screen-resolution-stats/mobile/worldwide ; https://gs.statcounter.com/screen-resolution-stats/desktop/worldwide (September 2026)
- https://web.dev/articles/optimize-lcp ; https://web.dev/articles/fetch-priority ; https://web.dev/articles/preload-responsive-images ; https://web.dev/articles/lcp
- https://chromium.googlesource.com/chromium/src/+/main/docs/speed/metrics_changelog/2023_04_lcp.md - 0.05 bpp rule
- https://csswizardry.com/2023/09/the-ultimate-lqip-lcp-technique/ **[secondary]**
- https://github.com/woltapp/blurhash ; https://bundlephobia.com/api/size?package=blurhash ; https://evanw.github.io/thumbhash/ ; https://leanrada.com/notes/css-only-lqip/ **[secondary]**
- https://webkit.org/blog/6682/improving-color-on-the-web/ (2016-07-01)

**Social and search**
- https://developers.facebook.com/docs/sharing/webmasters/images/
- https://ogp.me/
- https://www.linkedin.com/help/linkedin/answer/a521928
- https://developer.x.com/en/docs/x-for-websites/cards/overview/summary-card-with-large-image - HTTP 402 / client-rendered; values from https://kb.theseoframework.com/kb/twitter-cards-and-x-sharing/ **[secondary]**
- https://developers.google.com/search/docs/appearance/google-discover
- https://developers.google.com/search/docs/appearance/structured-data/article
- https://developers.google.com/search/docs/appearance/google-images

**Resize filter references**
- https://usage.imagemagick.org/filter/ (Lanczos, Catrom, Box notes)
- https://pillow.readthedocs.io/en/stable/handbook/concepts.html#filters (filter quality/performance table)

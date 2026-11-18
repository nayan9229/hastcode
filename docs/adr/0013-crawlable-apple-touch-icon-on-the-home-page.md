---
status: proposed
date: 2026-11-18
amends: ADR-0006
---

# A crawlable apple-touch-icon on the home page; the data: favicon stays

ADR-0006 embedded the favicon as a `data:` URI to take a request out of first paint, so the site had no icon file. Google's "Favicon in Search" says "Googlebot-Image must be able to crawl the favicon file" and asks for "a <link> tag to the header of your home page". It wants a square file, recommends one "larger than 48x48px", requires a stable URL, and accepts `rel` values `icon`, `shortcut icon`, `apple-touch-icon` and `apple-touch-icon-precomposed`. It says nothing about `data:` URIs. Google News also takes a publication's logo from the favicon (`docs/research/indexing-submission-channels.md` §1.7; ticket `.scratch/hastcode-indexing/issues/04-crawlable-favicon.md`).

The Author decided on 2026-11-18 (option (a) of the ticket):

- **One PNG file at `/apple-touch-icon.png`**, 192x192 (a multiple of 48), drawn by the build from the site mark (`/mark.svg`'s own path data, light palette, opaque, no rounded corners: iOS fills transparent pixels with black and masks the corners itself). The same bytes on every build; no content hash, so it sits at the site root, outside `/assets/`.
- **Linked from the home page only**, `<link rel="apple-touch-icon" href="/apple-touch-icon.png">`, the page Google reads the favicon from.
- **The `data:` URI stays every page's browser favicon** (ADR-0006, decision 2, unchanged).

Why this option:

- Google accepts `apple-touch-icon`, so the file is crawlable without touching the browser favicon.
- Desktop browsers do not fetch `apple-touch-icon` while loading a page, so it adds no request before first paint and no LCP. Measured: Lighthouse made no request for it (table below).
- Explainers gain zero bytes. The largest page had about 316 bytes of headroom under the 52 KB `html-size` gate on 2026-11-17.
- iPhones already request `/apple-touch-icon.png` on their own; those requests now get the icon instead of a 404.

The Author decided without a Lighthouse CI run on a Preview. The cost was measured locally instead (below); the Preview run happens on the next pushed Preview.

## Measured cost (local, 2026-11-18)

Preview build of the real `content/`, served on localhost by `cmd/subscribe` (the local router over `config.json`). `lighthouserc.json` settings (mobile, simulated throttling, 3 runs), @lhci/cli 0.15.1, Lighthouse 12.6.1. "Before" is the same build with the one link removed from the home page.

| Page | Mobile LCP median, before | After | CLS (worst run) | Icon requested? |
|---|---|---|---|---|
| Home `/` | 0.77 s (0.75, 0.77, 0.98) | 0.76 s (0.75, 0.76, 0.77) | 0 | no |
| Explainer `/languages-runtimes/g1-region-based-garbage-collection/` | 1.05 s | 1.05 s | 0 | no |

Both are under the 1.5 s budget, and the home page made 12 requests in both builds. The difference is run-to-run noise, not the link.

Size gate (production build of `content/`, HTML without the inline stylesheet, limit 53,248 bytes): the home page grew by 58 bytes, to 16,346 (36,902 left). No other page changed. The largest page is still the G1 Explainer at 52,917 bytes (331 left).

## Considered options

- **(b) a second `rel="icon"` file link on the home page**: a browser may pick the file over the `data:` icon and fetch it during load (our inference; not measured). That is ADR-0006's known LCP failure mode.
- **(c) the link on every page**: the same file, but about 58 bytes on every page, and the G1 Explainer has about 330 bytes left. Google reads the home page.
- **(d) no file**: keeps the generic icon in results and Google News.
- **`/mark.svg` itself as the file**: it would need a `rel="icon"` link, which is option (b)'s risk. iPhones also want a PNG at `/apple-touch-icon.png`, and one PNG serves both.

## Consequences

- Every build writes `/apple-touch-icon.png` into the static tree, in `r2` mode too. It is the one image outside the Image Store (ADR-0004), because its URL must stay fixed on the site's own host.
- On both Targets it gets the same cache headers as `/mark.svg` (the revalidating default, never `immutable`). It is not in `sitemap.xml`, and `robots.txt` does not block it. `builder/touch_icon_test.go` pins all of this.
- The icon is drawn from `markStrokePath` (`internal/views/touch_icon.go`), which `markSVG` also uses, so changing the mark changes both files. A path command the rasteriser does not know panics (a programming error on a constant, like `regexp.MustCompile`), so `make test` and every build crash before anything ships; it is not reported as a `build-step` error.
- Lighthouse CI on a Preview, the Search Console URL Inspection live test after the next Production Publish, and the later `site:` check remain the Author's checks (ticket 04).

# 04: Crawlable favicon file

**What to build:** The Author decides whether to keep the `data:` URI favicon for browsers (ADR-0006) and add a crawlable icon file next to it, so Google can show the site's mark beside its results. Google's favicon doc says "Googlebot-Image must be able to crawl the favicon file". It asks for a square icon, recommends one "larger than 48x48px", requires a stable URL, and accepts `rel` values `icon`, `shortcut icon`, `apple-touch-icon` and `apple-touch-icon-precomposed` (note §1.7). It also says "Add a <link> tag to the header of your home page". That line is quoted from Google's "Favicon in Search" page itself, which the note cites but does not quote; checked 2026-11-17. The page says nothing about `data:` URIs. Today the icon is only a `data:` URI and `/favicon.ico` is 404 (observed 2026-11-16; `docs/research/indexing-submission-channels.md` TL;DR "Three small fixes" 2, §1.7, "Code-side follow-ups" 1). Google News also takes a publication's logo from its favicon (`docs/research/google-news-discover-requirements.md`, technical checklist row 15 and §1).

This is a trade-off against ADR-0006. In ADR-0006's local measurements, embedding the favicon took about 0.12 s off simulated mobile LCP on the home page (0.90 s to 0.78 s). Any new request before first paint shows up as LCP in Lighthouse CI, against a 1.5 s budget. So the decision rests on two inputs: whether Google ignores our `data:` icon today, and what a file actually costs in Lighthouse CI.

Options to measure. Each serves a file of at least 48 px, square, at a stable path, and keeps the `data:` icon as every page's browser favicon:

- (a) an `apple-touch-icon` PNG linked from the home page only
- (b) a second `rel="icon"` file link on the home page only
- (c) the same as (a) or (b), on every page
- (d) no file: accept the generic icon

The existing `/mark.svg` (the Organization logo) may be enough as the file.

**Blocked by:** None (can start immediately)

**Status:** ready-for-human

Decision (Author):
- [ ] **Evidence for open question 2** ("Does Google use a `data:` URI favicon at all?"). Once the home page is indexed (ticket 01), search `site:hastcode.com` on a phone and record whether results show our mark or a generic globe, with the date.
- [ ] **Cost measured before deciding.** On a throwaway Preview branch, for each option being considered, record the Lighthouse CI median mobile LCP on the home page and one Explainer. Compare them with ADR-0006's table and the 1.5 s budget. Also record whether any page gets closer to the 52 KB HTML size gate (rule `html-size`; the spec's 50 KB, raised by ADR-0011). An agent can build the branch and collect the numbers; the Author makes the call.
- [ ] **Decision recorded.** The Author makes the call. An agent may draft the text for the Author to approve. The record is either an amendment to ADR-0006 (the option, its measured cost and the reason) or, for (d), a note here saying why not and what would change the answer.

If the decision is to add a file (an agent may do this part once the ADR-0006 amendment exists):
- [x] The build writes the icon file at one stable path (no content hash), square, at least 48 px and ideally a multiple of 48. It is linked from the pages the amendment names with a `rel` value Google supports, and the `data:` icon stays as the browser favicon.
- [x] Build-seam tests: the linked file exists in the output with the stated size and type. The link is on exactly the pages the amendment names. The `data:` icon link is still on every page.
- [x] The file is not in `/sitemap.xml` and is not blocked by robots.txt. Its name is stable, so it sits outside `/assets/`: that folder's `immutable` rule is only for content-hashed files. On each Target it gets the same cache rule as `/mark.svg` (README "Routes and caching" and "Deploying"). A test or the routes check shows the header on each Target.
- [ ] Lighthouse CI passes on the Preview (mobile LCP ≤ 1.5 s median, CLS 0); the 52 KB size gate passes; the real `content/` builds in production mode
- [x] README (the files the build writes, and the Production cache rules) mentions the file. ADR-0006's Consequences point to the amendment.
- [x] `make test` passes; preview build passes
- [ ] After the next Production Publish, the Author checks it with the Search Console URL Inspection live test on the home page, then repeats the `site:` check a few weeks later and records what it shows

## Comments

**Built, 2026-11-18** (branch `impl/fav-1118`). The Author chose option (a) on 2026-11-18, without a Preview Lighthouse run. The record is ADR-0013 (`docs/adr/0013-crawlable-apple-touch-icon-on-the-home-page.md`), which amends ADR-0006. ADR-0006's frontmatter and Consequences point to it. An agent drafted the ADR text, so the "Decision recorded" box stays open until the Author approves that text.

What the build does now:
- It writes `/apple-touch-icon.png`, a 192x192 opaque PNG (1,570 bytes) drawn from the mark's own path data (`internal/views/touch_icon.go`). The bytes are the same on every build, and it is written in `r2` mode too.
- Only the home page links it: `<link rel="apple-touch-icon" href="/apple-touch-icon.png">`. Every page keeps the one `data:` favicon.
- Tests in `builder/touch_icon_test.go` cover the size and type, the same bytes on every build, the home-page-only link, and the `data:` icon on every page. They also check that the icon gets the same headers as `/mark.svg` on both Targets (Cloudflare `_headers`; the Vercel routes through the local router), that it is not in `sitemap.xml` or blocked by `robots.txt`, and that it sits outside `/assets/`.

Measured locally (a stand-in for the Preview run, which is still the Author's):
- A preview build of `content/` was served on localhost by `cmd/subscribe` and run with `lighthouserc.json` (@lhci/cli 0.15.1, Lighthouse 12.6.1, 3 runs each). `lhci assert` passed.
- Home: median mobile LCP 0.76 s, against 0.77 s for the same build with the link removed. 12 requests in both.
- G1 Explainer: 1.05 s in both. CLS was 0 everywhere.
- Lighthouse never requested the icon.
- A production build of `content/` passed with no errors (930 warnings). The home page grew by 58 bytes, to 16,346 (36,902 under the 53,248-byte gate). No other page changed. The largest page is still the G1 Explainer at 52,917 bytes, 331 under the gate.

Still open, and the Author's: the `site:` evidence, Lighthouse CI on a pushed Preview, and the Search Console live test after the next Production Publish, then the `site:` check a few weeks later.

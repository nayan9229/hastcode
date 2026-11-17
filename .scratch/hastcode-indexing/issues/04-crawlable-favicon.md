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
- [ ] The build writes the icon file at one stable path (no content hash), square, at least 48 px and ideally a multiple of 48. It is linked from the pages the amendment names with a `rel` value Google supports, and the `data:` icon stays as the browser favicon.
- [ ] Build-seam tests: the linked file exists in the output with the stated size and type. The link is on exactly the pages the amendment names. The `data:` icon link is still on every page.
- [ ] The file is not in `/sitemap.xml` and is not blocked by robots.txt. Its name is stable, so it sits outside `/assets/`: that folder's `immutable` rule is only for content-hashed files. On each Target it gets the same cache rule as `/mark.svg` (README "Routes and caching" and "Deploying"). A test or the routes check shows the header on each Target.
- [ ] Lighthouse CI passes on the Preview (mobile LCP ≤ 1.5 s median, CLS 0); the 52 KB size gate passes; the real `content/` builds in production mode
- [ ] README (the files the build writes, and the Production cache rules) mentions the file. ADR-0006's Consequences point to the amendment.
- [ ] `make test` passes; preview build passes
- [ ] After the next Production Publish, the Author checks it with the Search Console URL Inspection live test on the home page, then repeats the `site:` check a few weeks later and records what it shows

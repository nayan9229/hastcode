# 02: Feed link in page head

**What to build:** A feed reader, a browser or a feed-intake service (such as Naver or Flipboard) that is given any hastcode.com page URL finds the site's RSS feed without being told its address. Today the feeds are linked only from the footer, and no page's `<head>` advertises one (observed 2026-11-16: `docs/research/indexing-submission-channels.md` TL;DR item "Three small fixes" (3), §7 "RSS readers", and "Code-side follow-ups" 2). The note gives Naver and Flipboard only as examples of feed intake ("Naver/Flipboard-style"). That a head link is how readers and browsers find a feed is the note's inference, not a vendor quote.

- Every page's `<head>` carries a feed autodiscovery link: `rel="alternate"`, type `application/rss+xml`, pointing at `/feed.xml`, with the title "HastCode" (the site feed's channel title).
- Every page of a Vertical hub's listing also advertises that Vertical's own feed, `/<vertical>/feed.xml`, titled the same as that feed's channel ("HastCode: <Vertical name>"). The Briefs listing (`/briefs/` and its later pages) also advertises `/briefs/feed.xml`, titled "HastCode: Briefs". On those pages the site feed comes first.
- Nothing else changes. The footer feed links stay. The feeds themselves, the Markdown-twin `alternate` link (ADR-0011), the inline stylesheet and the embedded favicon (ADR-0006) all stay as they are. Production, the QA Target and Previews carry the same links, because the build is the same. A browser does not fetch an `alternate` link while loading the page, so this adds no request before first paint and ADR-0006's known LCP failure mode does not apply. Lighthouse CI on the Preview confirms it. The link does add bytes to every page, and head chrome counts toward the 52 KB HTML size gate (the spec's 50 KB, raised by ADR-0011 after its own head link pushed three live pages close to the limit). So the real `content/` must still pass the gate.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] Build-seam test: a built home page, Explainer, Brief, static page and the 404 page each carry exactly one site-feed autodiscovery link with the right type, href and title
- [ ] Build-seam test: every listing page of a Vertical hub carries the site feed and then that Vertical's feed. The Briefs listing carries the site feed and then the Briefs feed. No other page carries a Vertical or Briefs feed link.
- [ ] Each advertised href is a feed the same build actually wrote (tested against the build output, so a renamed feed fails the test)
- [ ] The real `content/` still builds in production mode with the 52 KB HTML size gate (rule `html-size`) unchanged. If a live page would go over, stop and ask the Author; do not raise the gate. Lighthouse CI passes on the Preview.
- [ ] README "Feeds, sitemaps, robots and IndexNow" says pages advertise the feeds in `<head>` (site feed everywhere, Vertical and Briefs feeds on their listings)
- [ ] `make test` passes; preview build passes

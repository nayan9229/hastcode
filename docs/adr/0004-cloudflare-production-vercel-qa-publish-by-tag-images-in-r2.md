---
status: accepted
date: 2026-10-22
supersedes: the hosting half of ADR-0002 (the static-first Go build stays)
---

# Cloudflare serves Production, Vercel is QA, Publish by tag, images in R2

ADR-0002 put the static site and its one Go function on Vercel Pro. Vercel's CDN request tiers would cost $40–120/month at the traffic we are building for, the free Cloudflare plan allows commercial use with unmetered bandwidth, and hastcode.com's DNS already lives at Cloudflare. We decided:

- **Two Targets, one build.** The builder emits one static tree plus a manifest per Target. Cloudflare (Workers Static Assets) is the **Production Target** and the only host hastcode.com ever points at. Vercel (Hobby, access-protected, `noindex`, at `qa.hastcode.com`) is the **QA Target** for checking a release before it goes live. Branch Previews deploy to Cloudflare preview URLs.
- **Publish by tag, not by push.** Pushing `main` runs tests only. A tag `PUBCF-*` publishes that commit to Production; `PUBVL-*` publishes it to QA, built in production mode so QA sees the exact bytes. Both tags on one commit run independently. IndexNow is pinged only by a Production publish.
- **Images live in R2, not in the site output.** Every raster (hero variants, crops, author photos, OG badge) is uploaded at build time to an R2 bucket served at `images.hastcode.com` and referenced by absolute URL. Content-hashed keys mean uploads skip what already exists. Local builds keep writing images into the output (`images.store: local`); only CI uses `r2`. Orphans are pruned by a separate manual command, never by the build.
- **The signup function is ported to a Cloudflare Worker.** It is the only code that runs on Production; every other path is a static asset.

## Considered options

- **Stay on Vercel Pro only**: built and tested, but $20–140/month for a static site, and a second CDN layer under Cloudflare DNS.
- **Cloudflare Pages instead of Workers Static Assets**: same limits; Cloudflare's investment goes to Workers, so new projects start there.
- **Keep the Go function on Vercel and post cross-origin from Cloudflare**: one implementation, but Production signups would depend on the QA host.
- **Images in the site output on Cloudflare**: 16 files per Explainer against a 20,000-file cap (~1,200 Explainers on the free plan). R2 reduces each Explainer to one file.
- **410 for retired community URLs on Cloudflare**: needs Worker code on those paths; 404 is near-equivalent for Google (ADR-0001 amended).

## Consequences

- Vercel Hobby's terms reserve it for non-commercial use. QA shows no ads, takes no public traffic and is access-protected; we accept the grey area and will move to Pro if Vercel objects.
- A production build with `images.store: r2` and no credentials fails; it never silently falls back to local images.
- Absolute image URLs appear in HTML, JSON-LD, Open Graph and feeds, so changing the image hostname is a full rebuild and republish.
- The subscribe logic exists twice (Go for QA, JavaScript for Production) and both need tests at the HTTP seam.

---
status: accepted
date: 2026-10-22
---

# Repo Profiles are an Explainer kind with a fetched Fact Pack, not a product

The content plan (grilled 2026-10-22) adds Repo Profiles: pages that evaluate one trending open-source repository. The candidate brief ("RepoScope") described them as a separate, dynamic product: an empty-directory build with a database, six-hourly trending ingestion, an LLM pipeline, an admin area with login, a 0–100 health score with an SVG badge endpoint, leaderboards, tools pages, a public JSON API and an exchange API for a video project. We decided Repo Profiles are a **fourth Explainer kind inside hastcode.com**, produced the way every other Piece is produced, and that the product half of the brief is **out of scope** until Profiles have proven they earn traffic or citations.

Concretely:

- **URL and shape.** A Repo Profile lives at `/{vertical}/{slug}/` like every Explainer and inherits the Vertical hub, kind page, feeds, sitemaps, 90-day Review, House Byline and Editorial Policy disclosure. There is no `/repos/`, `/{owner}/{repo}`, `/compare/`, `/alternatives/`, `/topics/`, `/trending`, `/rising`, `/tools/*` or `/submit`. "X vs Y" pages are Comparison Explainers; the weekly "what trended and why" roundup is a Brief, and that Brief is the newsletter issue.
- **Curated, never exhaustive.** An Author chooses roughly three to five repositories a week through the Reader's lens (the existing Verticals) against written selection floors. Trending discovery is a shortlist aid for the Author, never a page generator.
- **Facts and opinion live in different files.** One deterministic, re-runnable Go command fetches a repository's Fact Pack from the GitHub REST API into a JSON file beside the Profile's Markdown. The builder renders the fact sections (header, Signals, Changelog) from the Fact Pack; the opinion sections (Architecture, Use it if / Don't use it if, Alternatives, FAQ, Verdict) are Markdown, AI-drafted in an editing session and verified by the Author. No LLM runs in the builder.
- **Refresh is manual, Publish stays a tag.** The Author runs `make refresh` locally (weekly), reviews the diff, commits and pushes; the next `PUBCF-*` tag ships it (ADR-0004). GitHub Actions never fetches from GitHub: it tests, builds Previews and publishes. Nothing is published without a human. (Amended 2026-10-22: the first version was a scheduled workflow opening a Refresh pull request; the Author chose a local run so CI stays deploy-only.)
- **No health score yet.** The Signals section shows raw facts without a number. A branded score, its methodology page and a badge are decided at the four-week checkpoint after the first five Profiles.
- **Hedge for the planned video project.** The Architecture section keeps one stable slot (Mermaid today, diagram or video later), and Fact Pack field names follow the brief's exchange JSON where they overlap, so a future export is a rename. No exchange API is built now.

## Why

- **Google's scaled-content-abuse policy (2024-03-05) and Raptive's quality review.** Hundreds of generated profiles are the shape that policy names, however they are made. A handful a week, each verified under the Editorial Policy, is not.
- **Raptive eligibility is per domain.** 25k pageviews/month, 50% Tier-1 traffic, domain age. A second site splits the pageviews and restarts the clock.
- **ADR-0002 and ADR-0004 bind.** The site is a static build with one Worker. The brief's admin area, cron ingestion and dynamic endpoints would reintroduce a server, a database and a login for a review gate the repository already provides (branch Preview, PR, tag).
- **Format risk before code.** The first five Profiles are written before any pipeline exists, so the code that follows is built for a format that has been seen to work.
- **GitHub's population is not Raptive's.** GitHub's own Innovation Graph puts US+UK+CA+AU+NZ at ~23% of developer accounts and falling, India at 14% and growing fastest (see `docs/research/github-trending-geography-and-data-access.md`). A page whose readers mirror GitHub would miss Raptive's 50% Tier-1 rule by half. Profiles therefore stay a minority of the site, and the checkpoint measures their Tier-1 share, not just their pageviews.
- **The API is the sanctioned path.** GitHub's Acceptable Use Policies define and police scraping of github.com; the ToS forbids pooling tokens or accounts to exceed rate limits, and there is no trending API. Fact Packs and the shortlist aid use the REST API (search with `pushed:`/`created:` qualifiers, `stargazers/count` and `stargazers/history`) under one token, never a scrape of the trending page.

## Considered options

- **Separate product on its own domain**, as the brief asked: a second brand, a second pipeline, a second ad-network application, and the Reader split in two. Rejected unless a second brand is wanted on purpose.
- **Repo Profiles as a new top-level Piece type beside Explainer and Brief**: a clean break, but it would re-implement hubs, feeds, Review and JSON-LD that Explainers already have. The glossary definition of Explainer was widened to "or examines one real system" instead.
- **A `/repos/{owner}/{repo}` URL space** for a predictable badge target: rejected because a badge can carry the canonical URL, and Vertical URLs keep Profiles inside the hubs Readers browse.
- **Ship the health score from day one**: the brief's "backlink engine". Deferred because a published score must be recomputable and documented, which is a feature, not a tedious step, and the format has not been validated yet.

## Consequences

- `internal/content` gains a `repo-profile` kind and a Fact Pack loader; the builder gains fact-section rendering and a generated hero/OG card per Profile (wordmark, repository name, language, stars, license) so no Profile needs a sourced image; a new command fetches Fact Packs and re-fetches them all (`make refresh`). Nothing else in ADR-0004's deploy shape changes.
- Selection floors (OSI license, commit in the last 90 days, fits a Vertical, ≥1,000 stars or clearly rising, not already profiled, not primarily UI/mobile/game) live in `docs/editorial/repo-profile-selection.md`. A GitHub license of `NOASSERTION` means "check by hand", not "fails the floor".
- Weekly rhythm: the Author runs a Refresh locally on Monday and commits it; roundup Brief on Friday, and that Brief is the newsletter issue, sent by hand from Buttondown. Repo Profiles are Explainers and count toward the cutover threshold of ADR-0001.
- The four-week checkpoint, measured at each Profile's day 28, compares Plausible pageviews and Tier-1 share with the median Explainer of the same age, plus Search Console impressions for the templated queries, Bing Webmaster AI Performance citations and inbound links. Health score and badge proceed only if the median Profile matches the median Explainer's pageviews **and** its Tier-1 share is within 10 points of the site median, or if any Profile is cited by an AI answer engine. Otherwise cadence drops and the product half stays shelved.
- Repo Profiles are capped at ten in total until the checkpoint (Author, 2026-10-22), so the format is proven on a small set first.
- The FAQ section is visible prose only. FAQPage JSON-LD was built and then dropped (Author, 2026-10-22): Google retired FAQ results in 2026 and `docs/research/aeo-citation-signals.md` found no citation benefit. `SoftwareSourceCode` stays.
- The checkpoint reopens exactly three questions: health score and badge, product pages (leaderboards, tools, API), and the video exchange. It does not reopen the static-first shape.

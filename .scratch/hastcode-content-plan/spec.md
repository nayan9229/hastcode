# Spec: hastcode.com content plan v1 — Repo Profiles as an Explainer kind

Status: implemented on integration/content-plan-v1 (2026-10-22); tickets 06 and 08 ready-for-human
Created: 2026-10-22
Glossary: ../../GLOSSARY.md (Content section: Explainer, Repo Profile, Fact Pack, Refresh, Review, Brief, Vertical, Reader)
ADRs: ../../docs/adr/0005-repo-profiles-as-an-explainer-kind-with-a-fetched-fact-pack.md (governs), ../../docs/adr/0004-cloudflare-production-vercel-qa-publish-by-tag-images-in-r2.md (deploy shape, unchanged), ../../docs/adr/0002-static-first-go-build-on-vercel.md (build half still stands), ../../docs/adr/0001-replace-genuin-community-with-static-publication.md (cutover threshold)
Editorial: ../../docs/editorial/repo-profile-selection.md (floors, cadence, first five), ../../content/pages/editorial-policy.md (AI-assistance note, Sources rule)
Research: ../../docs/research/github-trending-geography-and-data-access.md, ../../docs/research/aeo-citation-signals.md, ../../docs/research/ad-rpm-tech-audience.md
Shortlist: shortlist-2026-10-22.md
Predecessor: ../hastcode-publish-v2/spec.md (tickets 04 and 09 still open; this spec's tickets queue behind them and never touch the deploy shape)

## Problem Statement

hastcode.com has a build, a deploy shape and an Editorial Policy, but no content and no plan for which content earns traffic first. The Author's candidate idea ("RepoScope") was a separate dynamic product: profile every trending GitHub repository, with ingestion jobs, a database, an admin area, an LLM pipeline, a health score and badge, leaderboards and APIs. Grilling on 2026-10-22 settled that the valuable part is the page, not the product: a structured evaluation of one repository that is more useful than its README and that search and AI answer engines want to cite. Everything else in the idea either contradicts the static-first decisions, duplicates the review gate the repository already provides, or is a feature that should wait until the page format has proven itself.

Today the builder cannot express such a page: Explainers have three kinds, every piece needs a sourced hero image, nothing renders from fetched data, and nothing re-fetches anything on a schedule.

## Solution

Add **Repo Profile** as the fourth Explainer kind. A Repo Profile is written like any other Explainer (Markdown in the repository, AI-drafted in an editing session, verified by the Author, published by tag), but its fact sections (header card, Signals, Changelog, generated hero) render from a **Fact Pack**: a JSON file beside the Markdown that one deterministic command fetches from the GitHub REST API. A weekly **Refresh** workflow re-fetches every Fact Pack and opens one pull request when anything changed; the Author merges it and the next `PUBCF-*` tag ships it. Nothing publishes without a human.

The content plan around it: three to five Repo Profiles a week beside five other Explainers and five Briefs, chosen by an Author against written selection floors, spread across Verticals, roughly three rising repositories to two evergreen anchors. "X vs Y" pages are Comparison Explainers. The Friday "what trended this week and why" roundup is a Brief and doubles as the newsletter issue. The first five repositories are chosen (see the selection guide). A written checkpoint at each Profile's day 28 decides whether the deferred product half (health score, badge, leaderboards, tools, APIs, video exchange) is reopened.

## User Stories

### Reading a Repo Profile

1. As a Reader, I want a Repo Profile at the same kind of URL as every other Explainer (`/{vertical}/{slug}/`), so that it sits in the Vertical I browse and links like any other piece.
2. As a Reader, I want a header card at the top with what the repository is, who it is for, stars, license, primary language, latest release and last commit, so that I can judge it in five seconds.
3. As a Reader, I want every Repo Profile to run the same sections in the same order (TL;DR, In 60 seconds, Architecture, Signals, Use it if, Don't use it if, Alternatives, Changelog, FAQ, Verdict, Sources), so that I know where to look on every Profile.
4. As a Reader, I want the "In 60 seconds" section to show the install or quickstart commands as the README gives them, so that I can try the thing without leaving the page.
5. As a Reader, I want an Architecture section with a diagram and a short explanation of how the repository is built, so that I understand its shape before adopting it.
6. As a Reader, I want a Signals section that shows plain facts (last commit, releases in the last year, contributor count and the lead contributor's share of the listed top-ten contributions, open issues, security policy present, documentation files present) without a made-up number, so that I can weigh them myself.
7. As a Reader, I want "Use it if" and "Don't use it if" lists that take a position and ground it in the facts above, so that I get an opinion, not a tour.
8. As a Reader, I want an Alternatives section naming three to five alternatives with one line each, linking to Comparison Explainers or other Profiles where they exist, so that I can keep evaluating.
9. As a Reader, I want a Changelog section with the latest releases and their dates, kept current, so that I know the project is alive and what changed.
10. As a Reader, I want an FAQ section answering five to eight questions phrased the way I would search ("is X production ready?", "does X support Y?"), so that my question is answered on the page.
11. As a Reader, I want a one-paragraph Verdict at the end, so that I get the recommendation in one place.
12. As a Reader, I want a Sources section listing the repository, README and release pages the Profile drew from, so that every claim is checkable, as the Editorial Policy promises.
13. As a Reader, I want the page to carry the same AI-assistance note, byline, published and last-reviewed dates as other Explainers, so that I know how it was made.
14. As a Reader, I want Repo Profiles to appear in the Vertical hub, the home page, the kind page for Repo Profiles, feeds and search, so that I find them the way I find everything else.
15. As a Reader on a slow phone, I want the Profile to be as light as any Explainer (no scripts beyond the site's, no layout shift from the header card or ad slots), so that it loads fast.
16. As a Reader using dark mode, I want the header card and generated hero to respect the theme, so that nothing glares.

### Writing a Repo Profile

17. As an Author, I want to run one command with `owner/repo` that writes a Fact Pack beside my Markdown and a copy of the README where I draft, so that I never type stars, dates or file lists by hand.
18. As an Author, I want the fetch command to be deterministic and safe to re-run, so that running it twice on an unchanged repository changes nothing but the fetched-at timestamp.
19. As an Author, I want the fetch command to tell me when GitHub could not classify the license (`NOASSERTION`), so that I read the LICENSE file by hand instead of failing the floor.
20. As an Author, I want the fetch command to use one token from the environment and stop with a clear message when the rate limit is hit, so that I never pool tokens or scrape.
21. As an Author, I want a Markdown template and a short writing guide for Repo Profiles, so that a draft starts with the right headings in the right order.
22. As an Author, I want the build to fail with a named rule when a Repo Profile's headings are missing, extra or out of order, or when I wrote a Signals or Changelog section by hand, so that fact sections are never hand-maintained.
23. As an Author, I want the build to fail when a Repo Profile's frontmatter lacks its repository, or when a non-Profile Explainer carries one, so that the kind and the Fact Pack always agree.
24. As an Author, I want the build to fail when the Fact Pack is missing, malformed or names a different repository than the frontmatter, so that a page never renders stale or mismatched facts.
25. As an Author, I want the hero and Open Graph image generated from the Fact Pack (wordmark, repository name, description, language, stars, license) when I give no hero, so that a Profile needs no sourced image.
26. As an Author, I want to be able to supply my own hero for a Profile when I have a better image, so that the generated card is a default, not a rule.
27. As an Author, I want the Fact Pack's source URLs merged into the Sources section the way frontmatter sources are, so that the Sources rule holds without me copying URLs.
28. As an Author, I want a Comparison Explainer between two profiled repositories to be an ordinary Comparison, so that nothing new is needed to write "X vs Y".
29. As an Author, I want the Friday roundup to be an ordinary Brief whose "What happened" lists the week's repositories and whose `explainers` list links the Profiles, so that the newsletter issue needs no new piece type.
30. As an Author, I want the selection floors, cadence and first five repositories written down in the editorial docs, so that every week's choice follows the same rule.

### Refreshing facts

31. As an Author, I want a scheduled job every Monday that re-fetches every Fact Pack and opens one pull request if anything changed, so that Changelogs and header cards stay current without me remembering.
32. As an Author, I want that job to run the full test suite before opening the pull request, so that a changed Fact Pack can never break the next build unnoticed.
33. As an Author, I want the job to reuse its open pull request rather than open a second one, so that re-runs and the manual trigger are idempotent.
34. As an Author, I want the pull request to list which repositories changed and what (new release, star count, files), so that I can merge it in a minute.
35. As an Author, I want a Refresh to never touch prose, so that a merged Refresh cannot alter what I wrote.
36. As an Author, I want the Review cycle (90 days) to keep applying to a Profile's prose independently of Refreshes, so that opinion is re-read on schedule.

### Measuring

37. As the publisher, I want each Profile's day-28 numbers (Plausible pageviews and Tier-1 share vs the median Explainer of the same age, Search Console impressions for the templated queries, Bing Webmaster AI Performance citations, inbound links) collected in one note per Profile, so that the checkpoint rule can be applied as written.
38. As the publisher, I want the checkpoint rule written in the ADR and the selection guide, so that the health-score and product decisions are made by the rule, not by mood.

### Hedging

39. As the publisher planning a video project, I want the Architecture section to have one stable slot (a Mermaid diagram today) and Fact Pack field names that match the exchange JSON of the original brief where they overlap, so that a future export is a rename, not a redesign.

## Implementation Decisions

### Content model

- `repo-profile` is a fourth Explainer kind. Its frontmatter is the Explainer frontmatter plus a required `repo` key (`owner/name`). `repo` on any other kind is a build error; a `repo-profile` without it is a build error. `hero`/`hero_alt` become optional for this kind only; when absent the builder generates the hero.
- The Fact Pack is a JSON file beside the Profile's Markdown, named after the Markdown file. The builder loads it with the Explainer, fails the build when it is missing, unparsable, or its repository differs from the frontmatter `repo`, and warns (`fact-pack-stale`) when the Fact Pack has been unchanged upstream for over 30 days (a Refresh keeps an unchanged Fact Pack's fetched-at date, so its age is how long the facts have not changed): the repository may be dormant, and the Author checks it at the next Review.
- Fact Pack v1 fields (one GitHub REST call each, at most ten calls per repository): repository (full name, URL, description, homepage, topics, license SPDX id plus a `needs_review` flag for `NOASSERTION`, default branch, stars, forks, watchers, open issues, created and pushed dates), languages by bytes, latest ten releases (tag, name, date, URL, prerelease flag), contributors (total count and the top ten with contribution counts, from the contributors list, not the stats endpoint), repository files present at the root and under `.github` (README, LICENSE, SECURITY, CONTRIBUTING, CODE_OF_CONDUCT, CI configuration, Dockerfile, package manifests, `docs/` and `examples/` directories, as booleans or paths), the README URL and blob SHA, the stargazer count, and a `sources` list of the github.com pages a Reader can open to check what was fetched (the repository page, the README blob page, the releases page, the contributors graph page and the `.github` tree page; each page once, never an API endpoint). Field names follow the original brief's exchange JSON where they overlap (`slug`, `url`, `name`, `language`, `license`, `stars`, `latest_release`, `manifests`). Deferred until the checkpoint: vulnerabilities, dependents, download counts, pull-request merge rate, issue response time, star history.
- The README text is written to the builder's cache directory (not the content tree) for drafting; the Fact Pack carries only its URL and SHA. No README copy is committed.
- The Fact Pack renders: the header card (what it is from the frontmatter description, who it is for from a new optional `who_for` frontmatter key, stars, license, language, latest release, last commit), the Signals section, the Changelog section, and its `sources` merged into the Sources section once per URL.

### Page composition and gates

- The Markdown body holds only the opinion sections, in this fixed order of H2 headings: TL;DR (existing rule: first, three to five bullets), In 60 seconds, Architecture, Use it if, Don't use it if, Alternatives, FAQ, Verdict, Sources (existing rule: last, with a link). A new skeleton rule for this kind fails the build on a missing, extra or reordered heading, and on a hand-written Signals or Changelog heading with a message saying they render from the Fact Pack.
- The builder composes the page: header card, then the body with the Signals section inserted after Architecture and the Changelog section inserted after Alternatives. FAQ questions are H3 headings under FAQ; the rule requires five to eight.
- Signals shows facts only: last commit date, releases in the last 12 months and latest release date, contributor count and the lead contributor's share of the listed top-ten contributions, open issues, security policy present, license, stars, forks, watchers, and which documentation files and directories exist. No score, no colour grading.
- The Architecture section's first Mermaid fence is the stable diagram slot (rendered as today). A later diagram or video replaces that slot without changing the surrounding markup.
- JSON-LD: the page keeps the Article graph every Explainer emits and adds `SoftwareSourceCode` (name, code repository URL, programming language, license URL, description) and `FAQPage` built from the FAQ H3s and their answers. The allowed-types test is extended accordingly.
- Kind page and badge: "Repo Profile" / "Repo Profiles", described with the glossary sentence.
- Generated hero: a 1600×900 PNG rendered in Go with the bundled Go font family already in the module graph, showing the wordmark, repository full name, description (truncated), language, stars and license on the brand palette, light theme (the hero pipeline then produces the usual variants and the 1200×630 Open Graph crop). It enters the existing hero pipeline as if it were a sourced file, so every hero rule still holds. Alt text is generated ("{repo}: {description}").

### Fetch command

- One new command alongside the builder, prune and adapt commands: given `owner/name` and the target Markdown path, it fetches the Fact Pack and writes it; with a refresh flag it re-fetches every Fact Pack found under the content tree, rewrites those that changed (ignoring fetched-at), prints a one-line summary per changed repository, and exits with a distinct code when anything changed.
- It reads one token from the environment (the same variable name GitHub Actions provides), sends an identifying User-Agent, respects `Retry-After` and the primary rate-limit headers with a bounded backoff, and stops with a clear message rather than retrying past the limit. Unauthenticated use is refused with a message explaining why.
- The GitHub client is a small interface over `net/http`; the command is built against it so tests can serve recorded JSON from an in-process server.
- Output JSON is written with stable key order and indentation so Refresh diffs are readable.

### Refresh workflow

- A scheduled GitHub Actions workflow (Monday 06:00 UTC, plus manual trigger) checks out `main`, runs the fetch command in refresh mode, and stops if nothing changed. If something changed it runs the full check target and a production build of the content tree (the Publish build with the Image Store local, so no secret is needed), then commits to a `refresh/<date>` branch and opens a pull request titled with the date, listing the changed repositories. There is one open Refresh pull request at a time: when any open pull request's branch starts with `refresh/`, the job resets that branch onto `main`, commits, force-pushes it with a lease and edits the pull request's title and body instead of opening a second. The job uses the repository token; because pull requests opened by that token do not trigger the test workflow, the job runs the tests itself before opening it, and the README says so.
- Publish stays by tag (ADR-0004). The workflow never tags.

### Editorial docs

- A Repo Profile Markdown template and a one-page writing guide (what goes in each section, how to cite the Fact Pack, the quickstart-verbatim rule, title and description patterns that match search intent: "{repo}: …", "{repo} vs {alternative}", "is {repo} production ready", "{repo} alternatives") live in the editorial docs beside the selection guide.
- The roundup Brief: "What happened" lists the week's repositories with one line each and links their Profiles; "Why it matters" is the commentary; `source_url` is the week's trending page; `explainers` lists the linked Profiles (existing gate). The Buttondown issue is sent by hand from the published Brief.
- The checkpoint note: a short Markdown file per Profile under the feature's scratch directory, created at day 28 with the numbers in User Story 37.

### What does not change

- The deploy shape of ADR-0004, the Image Store, the signup Worker, routes and caching, the 50 KB HTML gate, the section-length guide, the Review-overdue rule.
- No LLM runs in the builder or in any workflow. Drafting happens in editing sessions.

## Testing Decisions

- A good test drives the highest seam and asserts on what the Reader or Author sees: built HTML, JSON-LD, build errors and warnings by rule name, files written. It never asserts on internal structs or private functions.
- **Builder seam (existing).** Every page behaviour above is tested through the build-fixture harness: a content tree with a valid Repo Profile fixture and a valid Fact Pack fixture goes in; the output directory and build result come out; assertions use the page, JSON-LD and rule helpers that Explainer and Brief tests already use. Prior art: the Explainer, Brief gate, hero, diagram, JSON-LD, listing and kind-page tests in the builder package. The fixture helpers gain a valid Repo Profile and a valid Fact Pack.
- **Fetch seam (new).** The fetch command is tested against an in-process HTTP server serving recorded GitHub JSON: the written Fact Pack matches a golden file; a second run on the same responses is byte-identical except fetched-at; refresh mode reports exactly the changed repositories and the right exit code; `NOASSERTION` sets the review flag; a 403 with rate-limit headers stops with the documented message; missing token is refused. The command's main stays thin and untested.
- **Generated hero.** Tested through the builder seam: a Profile without a hero produces the five widths, the Open Graph crop and the preload the hero tests already require; a Profile with a hero uses it; the card's alt text is set.
- **Refresh workflow.** Not unit-tested. The workflow lint job already in CI covers syntax; the first run is a manual trigger checked by hand and recorded in the README.
- **Editorial docs and the roundup Brief.** A roundup Brief fixture passes the existing Brief gates in the builder tests, proving no new piece type is needed.

## Out of Scope

- Health score, badge endpoint, methodology page, star-history sparklines (checkpoint).
- Leaderboards, tools pages, topic hubs, submit page, public JSON API, internal exchange API, webhook, offline export, admin area, login, database, any server-side ingestion (checkpoint or never; ADR-0005).
- Automatic shortlist generation. The shortlist is produced by a person (or an editing session) from the sources the selection guide names; a shortlist command may follow once the weekly routine has been run by hand a few times.
- Vulnerability, dependents and download data in the Fact Pack.
- Any change to the deploy shape, Targets, Image Store or signup Worker.
- Monetisation plumbing beyond the existing ads flag and newsletter: affiliate registry, sponsor slot UI, paid fast-track, job board.
- Video embedding and video structured data.

## Further Notes

- Sequencing: these tickets wait behind publish-v2 tickets 04 and 09 and the merge of `integration/publish-v2`. The first ticket lands the kind and the Fact Pack shape; the fetch command, skeleton rule, fact sections and hero card fan out from it in parallel; the Refresh workflow follows the fetch command; the first Profile (apache/iggy) closes the loop. Tickets: `issues/`.
- The first five: apache/iggy (distributed-systems), xerj-org/xerj (databases-storage), NVIDIA/OpenShell (infrastructure-cloud), deepseek-ai/DeepInfer (ai-ml-systems), php/frankenphp (languages-runtimes). Alternate cilium/cilium. Evidence in the shortlist file.
- Geography is the risk to watch: GitHub's population is about a quarter Tier-1 and Raptive wants half. Profiles stay a minority of output and the checkpoint measures Tier-1 share, not only pageviews.
- GitHub removed public access to per-user stargazer lists in July 2026 and added `stargazers/count` and `stargazers/history`; the Fact Pack reserves a place for history but v1 stores only the count.
- The Cloudflare free plan's cron limits (10 ms CPU, 50 subrequests) rule out a Worker for Refresh; GitHub Actions is the right place and already holds the token.

## Review fixes

Spec text changed by the review fixes on `cp/review-fixes` (2026-10-22):

- **User Story 6 and "Page composition and gates"**: Signals shows the contributor count and the lead contributor's share of the listed top-ten contributions (what the page has always shown), not "the share held by the top ten".
- **"Content model", Fact Pack fields**: `sources` lists the github.com pages a Reader can open to check what was fetched (repository, README blob, releases, contributors graph, `.github` tree), each once, never an API endpoint; an Author's citation of the same page merges with it.
- **"Content model", stale warning**: `fact-pack-stale` means the Fact Pack has been unchanged upstream for over 30 days (a Refresh keeps an unchanged Fact Pack's fetched-at), so the repository may be dormant and the Author checks it at the next Review; it no longer advises a Refresh.
- **"Refresh workflow"**: after the check target the job runs a production build of the real content tree (the Publish build with the Image Store local, no secret), so a changed Fact Pack that breaks the build fails the job before the pull request; and there is one open Refresh pull request at a time: any open pull request from a `refresh/` branch is reused (branch reset onto `main`, force-pushed with a lease, title and body edited) instead of opening a second.

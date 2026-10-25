# Spec: hastcode.com content intake — Source Watchlist, Leads, and the verified gate

Status: ready-for-agent
Created: 2026-10-26
Glossary: ../../GLOSSARY.md (Sourcing: Source, Source Watchlist, Lead, Intake; also Piece, Explainer kinds, Brief, Vertical, Author, House Byline)
ADRs: ../../docs/adr/0007-content-intake-produces-leads-not-pieces.md (governs), 0005 (Repo Profiles, curation, scaled-content concern), 0001 (launch gate), 0004 (Actions deploy-only, Publish by tag)
Research: ../../docs/research/content-source-intake.md (77 verified feeds, Google policy quotes, copyright limits, intake mechanics)
Prior art: the Repo Profile Fact Pack fetcher (fetch and Refresh, run locally, GitHub REST with the repo token, tested against an in-process fake GitHub)

## Problem Statement

HastCode has ten Repo Profiles and nothing else: no Case Studies, Concepts, Comparisons or Briefs. The launch gate needs at least 30 Explainers across all six Verticals and 20 Briefs, and the plan is 5 Explainers and 5 Briefs a week after that. The Author works alone and has no systematic way to find what to write about: candidate material is scattered across dozens of engineering blogs, release feeds and postmortem pages. The obvious shortcut (generate pages from feeds) is what Google's spam policies name, and the ten Profiles already showed that a verification step can be skipped under time pressure.

## Solution

A committed Source Watchlist (6 first-party Sources per Vertical plus incident feeds, each with a default Vertical and kind) and a local `make intake` that checks it each morning and lists the new Leads: title, Source, date, link, suggested kind and Vertical. Intake never writes prose. The Author picks at most two Leads a day; `make lead` turns a pick into a draft skeleton that a drafting session fills in the house-byline voice; `make lead-dismiss` hides the rest; undecided Leads expire after 7 days. Briefs come only from Leads at most 3 days old. Production builds refuse any new non-draft Piece without `verified: <date>`, which only the Author sets.

Order of work (refining ADR-0007's "format before pipeline"): the gate and the listing half of Intake come first because they do not depend on article format; then the Author writes one Case Study, one Concept, one Comparison and one Brief by hand from real Leads; only then are the per-kind templates, writing guides and the `make lead` skeleton built from those pieces.

## User Stories

### Source Watchlist

1. As the Author, I want a committed Source Watchlist, so that the Sources I follow are reviewed and versioned like code.
2. As the Author, I want each Source to carry a name, its site, its feed or GitHub repository, its Vertical, the kind it usually feeds, and a one-line reason, so that I remember why it is there.
3. As the Author, I want to start with 6 Sources per Vertical plus the incident feeds, so that the daily Lead list stays readable.
4. As the Author, I want a build-time or intake-time check that every Source names a known Vertical and kind, so that a typo cannot misfile Leads.
5. As the Author, I want to add a Source with one entry and a reason, so that growing the Watchlist is a deliberate commit.
6. As the Author, I want GitHub repositories listed as Sources by owner/name, so that their releases come through the REST API rather than scraped pages.

### Intake

7. As the Author, I want `make intake` to check every Source and list only Leads I have not seen before, so that each morning shows what is new.
8. As the Author, I want each Lead to show title, Source, publication date, link, suggested kind and Vertical, and its age, so that I can decide in seconds.
9. As the Author, I want the suggested kind and Vertical copied from the Source, never guessed by AI, so that suggestions are predictable.
10. As the Author, I want Intake to read RSS 2.0 and Atom feeds, so that every verified Source works.
11. As the Author, I want Intake to send conditional requests with the validators each Source last returned, so that unchanged feeds cost almost nothing and Sources are not hammered.
12. As the Author, I want Sources checked one at a time with a clear User-Agent naming HastCode and a contact address, so that publishers see a polite, identifiable client.
13. As the Author, I want GitHub releases fetched through the REST API with my existing token, so that I stay within GitHub's terms and rate limits.
14. As the Author, I want a Source that fails (network error, 4xx/5xx, unparseable feed) reported by name without stopping the run, so that one broken feed never hides the rest.
15. As the Author, I want a Source that has published nothing for 90 days flagged in the run summary, so that I can retire dead Sources.
16. As the Author, I want the Lead queue and seen-state kept locally and out of git, so that daily working notes never churn the repository.
17. As the Author, I want Leads older than 7 days that I never decided on dismissed automatically, so that the queue never becomes a backlog.
18. As the Author, I want the run summary to show counts per Vertical and per suggested kind, so that I can see which Vertical is starving.
19. As the Author, I want a Lead whose link points at a page HastCode already cites in a Piece marked as covered, so that I do not write the same story twice.
20. As the Author, I want Intake to never write anything under the content tree, so that it cannot publish by accident.

### Picking and dismissing

21. As the Author, I want `make lead ID=…` to turn a Lead into a draft skeleton at the right content path with frontmatter (kind, Vertical, author mira-vance, draft: true, published left for me), the Source link already in Sources, and the kind's empty section headings, so that drafting starts from structure, not a blank file.
22. As the Author, I want to override the suggested kind or Vertical when picking, so that a Source's default never forces the wrong format.
23. As the Author, I want picking a Lead older than 3 days as a Brief refused with a clear message, so that Briefs stay timely.
24. As the Author, I want `make lead-dismiss ID=…` to hide a Lead for good, so that the queue only holds open decisions.
25. As the Author, I want picking to refuse a slug that already exists, so that an existing Piece is never overwritten.
26. As the Author, I want the skeleton to contain no generated prose, so that every sentence that ships is drafted and verified on purpose.

### Verification gate

27. As the Author, I want a production build to fail on any non-draft Piece published after 2026-10-26 that has no `verified:` date, naming the file, so that nothing ships unchecked.
28. As the Author, I want `verified:` to be a date no earlier than `published` and not in the future, so that the field cannot be filled carelessly.
29. As the Author, I want preview builds to only warn on a missing `verified:`, so that drafts and Previews still work while I check.
30. As the Author, I want the ten earlier Repo Profiles to keep publishing until I clear their review queue, so that the gate does not take the site down today.
31. As the Reader, I want every new Piece to have been checked against its Sources by a person, so that I can trust what I read.
32. As an agent, I want a standing rule that I never set `verified:`, so that the field always means a human checked.

### Format before pipeline

33. As the Author, I want a one-off Intake run to propose three Leads per kind (Case Study, Concept, Comparison) for the first hand-written pieces, so that the format is proven on real material.
34. As the Author, I want one Case Study, one Concept, one Comparison and then one Brief linking one of them written and verified before templates are fixed, so that templates describe what worked.
35. As the Author, I want a template and a writing guide per kind (Case Study, Concept, Comparison, Brief), drawn from those pieces and matching the Repo Profile documents' shape, so that every drafting session follows the same rules.
36. As the Author, I want every writing guide to say: first-party Sources only, every claim linked, short credited quotes only as long as the point needs, diagrams always redrawn in Mermaid with the original linked, code quoted only under the repository's licence, so that copyright and Google's policies are respected by default.
37. As the Author, I want the Brief guide to require a "why it matters" analysis that is HastCode's own, never a restatement of the Source, so that Briefs add original value.

### Cadence

38. As the Author, I want the documented daily rhythm (run Intake each morning, pick at most two Leads, dismiss the rest), so that the pipeline matches my time.
39. As the Author, I want the planned pace documented as 5 Explainers (2 Case Studies, 2 Concepts, 1 Comparison) and 5 Briefs a week, with every Vertical getting a new Explainer in the first two weeks, so that the launch gate is reached in about four weeks.

## Implementation Decisions

### Source Watchlist

- A committed YAML file at the repo root holds the Watchlist: a list of Sources, each with name, site URL, either a feed URL (RSS or Atom) or a GitHub repository (owner/name), Vertical (slug from site config), kind (case-study, concept, comparison, brief), and reason. Incident feeds are ordinary Sources whose kind is case-study.
- The initial file is drafted from the research tables: 6 Sources per Vertical chosen for first-party signal and cadence, plus the 11 incident feeds. The Author reviews it before it is committed; the drafting agent does not commit it unreviewed.
- Validation (unknown Vertical or kind, both or neither of feed and repository, duplicate names) runs at the start of every Intake run and fails the run.

### Intake

- A new Intake package with three entry points: run (Watchlist, state directory, clock, HTTP client and GitHub base URL in; Lead queue and run summary out), pick (Lead ID and optional kind/Vertical override in; draft skeleton path out), dismiss (Lead ID in). Thin commands wrap them as `make intake`, `make lead` and `make lead-dismiss`.
- Feeds: RSS 2.0 and Atom only (research saw no JSON Feed). Conditional GET with stored ETag and Last-Modified per Source (RFC 9110); a Source without validators is fetched in full. Sequential requests, a fixed User-Agent naming HastCode and the contact address, a per-request timeout.
- GitHub: releases through the REST API with the token the Fact Pack fetcher already uses, reusing its client; never `releases.atom`.
- Lead identity: a stable ID derived from the Source and the item's guid or link, short enough to type. Seen-state records every Lead ID ever listed, with its decision (open, picked, dismissed, expired) and the date first seen.
- Queue: a local Markdown file the Author reads (newest first, grouped by Vertical), regenerated each run from state. Both queue and state live in a gitignored local directory; Intake writes nothing under the content tree.
- Expiry: open Leads first seen more than 7 days before the run's clock become expired.
- Covered detection: a Lead whose link (normalised) already appears in any Piece's Sources or body is marked covered in the queue.
- Run summary: new Leads per Vertical and kind, Sources that failed (and why), Sources silent for 90 days.

### Picking

- Pick writes a skeleton at the kind's content path with frontmatter (title from the Lead, kind, Vertical, author mira-vance, draft: true, no published date, no verified), a Sources section holding the Lead's link, and the kind's section headings from its template, then marks the Lead picked. Refuses an existing slug, and refuses kind brief when the Lead is older than 3 days by the run clock.
- Pick lands only after the per-kind templates exist (see order of work); until then it is out of reach by design.

### Verification gate

- Frontmatter gains an optional `verified` date. A new build gate: in production builds, every non-draft Piece whose `published` is after 2026-10-26 must have `verified`, not earlier than `published` and not after the build clock; otherwise a build error naming the file (rule `unverified`). In preview builds the same condition is a warning. Pieces published on or before 2026-10-26 are outside the rule.
- CLAUDE.md already states that an agent never sets `verified:`.

### Documents

- Per kind (Case Study, Concept, Comparison, Brief): a template and a writing guide under the editorial docs, in the same shape as the Repo Profile template and writing guide, written after the corresponding hand-written piece. Every guide carries the sourcing, quoting, diagram and code-licence rules from the research; the Brief guide carries the 3-day freshness rule and the original-analysis requirement.
- README gains an Intake section (commands, daily rhythm, where state lives, how to add a Source) and documents the `verified` gate.

## Testing Decisions

- A good test exercises an agreed seam and asserts on observable output: the queue file, the state, the skeleton file, the build result. No tests of feed parsing, ID hashing or HTTP helpers in isolation.
- **Intake seam** (new; prior art: the Fact Pack fetcher's tests with an in-process fake GitHub): tests start `httptest` servers serving RSS 2.0 and Atom fixtures and a fake GitHub releases endpoint, run Intake with a temp state directory and a fixed clock, and assert: new Leads listed with the Source's kind and Vertical; a second run lists nothing new; conditional requests sent with stored validators and a 304 handled; a failing Source reported while others succeed; a silent Source flagged; 7-day expiry; covered Leads marked; nothing written under the content tree; Watchlist validation errors. Pick tests assert the skeleton's frontmatter, Sources link and headings, the slug-exists refusal, the 3-day Brief refusal, and kind/Vertical override; dismiss removes a Lead from later queues.
- **Build seam** (existing): fixtures for the `verified` gate: new non-draft Piece without it fails production naming the file; with a valid date passes; a date before `published` or in the future fails; a draft is unaffected; a Piece published on or before 2026-10-26 is unaffected; preview builds warn instead of failing.
- Go-only tests (Author decision), `make test` runs them.

## Out of Scope

- Any AI step inside Intake (classification, summarisation, drafting).
- Scheduled Intake in GitHub Actions or anywhere else; Intake runs only when the Author runs it.
- JSON Feed, scraping HTML pages, `releases.atom`, Sources without a feed or GitHub repository (they stay manual).
- Changing the Repo Profile pipeline, the Profile cap, or the ten existing Profiles' content.
- Auto-publishing, auto-tagging, setting `verified:` by anything but the Author.
- Newsletter assembly from Leads.

## Further Notes

- The ADR-0007 order of work is refined, not changed: the listing half of Intake is format-independent and is needed to propose the first Leads (round-2 decision), so it precedes the hand-written pieces; the skeleton and templates follow them.
- The hand-written pieces are the Author's work with drafting help; their tickets end ready-for-human at the verification step.
- Research open questions that matter here: whether any Source's terms forbid automated feed reading (none found among those checked), and that the Go blog sends no validators (fetched in full daily is acceptable).

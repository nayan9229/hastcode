# HastCode

Ad-monetised publication at hastcode.com: evergreen software-architecture explainers, fed by short commentary on tech events. English only. Written for mid-to-senior backend engineers.

## Language

### Content

**Explainer**:
An evergreen, long-form article that teaches one architectural idea or examines one real system. The spine of the site; everything else exists to feed it. Comes in four kinds: Case Study, Concept, Comparison, Repo Profile.
_Avoid_: Post, blog, tutorial, deep dive

**Case Study**:
An Explainer that reconstructs how one real, named system is built, from that company's public material. Flagship kind.
_Avoid_: Teardown, breakdown

**Concept**:
An Explainer that defines and teaches one named architectural idea or pattern.
_Avoid_: Guide, primer, 101

**Comparison**:
An Explainer that sets two or more alternatives side by side to help the Reader choose.
_Avoid_: Versus, showdown, review

**Repo Profile**:
An Explainer that evaluates one open-source repository: what it is, how it is built, and whether the Reader should adopt it. Its fact sections render from a Fact Pack; its opinion sections are written.
_Avoid_: Deep dive, review, profile page, listing, RepoScope

**Fact Pack**:
The machine-fetched facts about a repository (description, stars, license, languages, releases, contributors, which files exist) that a Repo Profile's fact sections render from and its opinion sections are drafted against. Fetched, never hand-written.
_Avoid_: Metadata, data file, facts.json, scrape

**Refresh**:
A re-fetch of every Repo Profile's Fact Pack, run by the Author, that rewrites only the facts that changed upstream. Touches facts only, never prose; a Review is the prose counterpart.
_Avoid_: Sync, update, auto-update, crawl

**Brief**:
A short, dated commentary on something that just happened in tech, with the architectural "why it matters", linking into Explainers. Curated from public sources; never original reporting.
_Avoid_: News, article, update, digest

**Review**:
A scheduled re-read of a published Explainer by its Author, every 90 days, that either changes it (and moves its last-reviewed date) or confirms it is still right (and does not). A Refresh is not a Review: it changes facts, not prose.
_Avoid_: Refresh (that is the Fact Pack re-fetch), update, audit

**Piece**:
Any published Explainer or Brief. Used when something applies to both.
_Avoid_: Post, article, content item, entry

**Tag**:
A free-form label on an Explainer, used only to find Related Explainers across Verticals. Never shown as a browsing structure.
_Avoid_: Topic, keyword, category

**Vertical**:
One of the site's top-level subject areas that every Explainer belongs to. At launch: Distributed Systems, Databases & Storage, APIs & Networking, Infrastructure & Cloud, AI/ML Systems, Languages & Runtimes. Planned: Security, Frontend Architecture.
_Avoid_: Category, topic, section, tag

### People

**Reader**:
The person the site is written for: a mid-to-senior backend engineer. Juniors are welcome but not the target.
_Avoid_: User, visitor, audience

**Author**:
The named person, or the team publishing under one byline, who has an author page and is accountable for an Explainer or Brief. At launch the only Author is the House Byline.
_Avoid_: Writer, contributor, editor

**House Byline**:
A disclosed editorial voice the team publishes under, with its own name and author page, never presented as a real individual. HastCode's is Mira Vance.
_Avoid_: Persona, pen name, ghostwriter, fake author

### Publishing

**Target**:
One of the two hosts a build can be deployed to. The Production Target serves hastcode.com; the QA Target is for checking a release and never serves hastcode.com.
_Avoid_: Environment, stage, primary, standby

**Publish**:
Deploying one commit of the main branch to a Target by tagging it. The tag names the Target.
_Avoid_: Release, ship, push live

**Preview**:
A build of a branch, drafts included, deployed to a temporary address for review before merging.
_Avoid_: Staging, QA (that is a Target)

**Image Store**:
Where a Piece's image variants live and are served from, separate from the pages.
_Avoid_: CDN, bucket, media library

---
status: accepted
date: 2026-10-20
---

# Replace the Genuin community at hastcode.com with a static publication

hastcode.com currently serves a Next.js Genuin community app (~1,000 `/video/*` pages, 11 `/community/*`, 24 `/group/*`, all unrelated to software architecture). We decided the new publication takes the apex domain outright: the community is retired, every old URL returns `410 Gone` (not a 301 to the homepage, which Google treats as a soft 404 anyway), and `/home` 301s to `/`. We chose this over a subdomain or a `/articles/*` path because the old content carries no topical authority worth inheriting and a split domain would dilute the new site's SEO signals from day one.

## Consequences

- Cutover is one-way: the Genuin community cannot be restored at this domain without evicting the publication.
- The new site has no reader-facing community features at launch (no comments, reactions, or groups). The newsletter is the only reader relationship.

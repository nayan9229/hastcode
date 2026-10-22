# 01: Repo Profile kind renders from a Fact Pack

**What to build:** An Author can add an Explainer with `kind: repo-profile` and a `repo: owner/name` frontmatter key, drop a Fact Pack JSON file beside it, and get a page at `/{vertical}/{slug}/` whose header card shows what the repository is, who it is for (new optional `who_for` key), stars, license, primary language, latest release and last commit, all taken from the Fact Pack. The piece appears in its Vertical hub, the home page, feeds, search and a new "Repo Profiles" kind page, with the kind badge in listings. The Fact Pack shape (spec "Content model") is defined in this ticket and shared with ticket 05; the test fixtures gain a valid Repo Profile and a valid Fact Pack. A sourced hero is still required here (ticket 04 removes that).

Spec: `../spec.md` ("Content model", User Stories 1, 2, 13, 14, 23, 24). Glossary: Repo Profile, Fact Pack. ADR-0005 governs.

**Blocked by:** None (can start immediately; queue behind publish-v2 tickets 04 and 09 and the merge of `integration/publish-v2`)

**Status:** ready-for-agent

- [ ] `kind: repo-profile` is accepted; the frontmatter error for unknown kinds names four kinds
- [ ] `repo` is required for a Repo Profile and a build error on any other kind; `who_for` is optional on a Repo Profile and rejected elsewhere
- [ ] A Fact Pack beside the Markdown, named after it, is loaded; missing, unparsable, or naming a different repository than `repo` is a build error naming the file and rule; fetched more than 30 days ago is a warning
- [ ] The page renders at the Vertical URL with the header card fields above, no layout shift (fixed-height card), and the existing byline, dates and AI-assistance note
- [ ] Kind page "Repo Profiles" lists Profiles across Verticals newest first; listings show the "Repo Profile" badge; the kind description is the glossary sentence
- [ ] Fact Pack field names follow the spec, including the exchange-JSON names where they overlap (`slug`, `url`, `name`, `language`, `license`, `stars`, `latest_release`, `manifests`) and `needs_review` on the license
- [ ] Fixture helpers provide a valid Repo Profile and a valid Fact Pack; every existing builder test still passes

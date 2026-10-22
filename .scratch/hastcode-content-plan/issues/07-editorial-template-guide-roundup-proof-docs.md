# 07: Editorial template, writing guide, roundup proof and docs

**What to build:** An Author starting a Repo Profile copies a Markdown template that already has the frontmatter keys and the nine headings in the order the skeleton rule demands, and follows a one-page writing guide: what belongs in each section, how to ground every claim in the Fact Pack and README, the quickstart-verbatim rule for "In 60 seconds", the three-to-five Alternatives rule, the five-to-eight FAQ rule phrased as searches, and title and description patterns that match search intent ("{repo}: …", "{repo} vs {alternative}", "is {repo} production ready", "{repo} alternatives"). A checkpoint-note template captures the day-28 numbers the ADR's rule needs. The README and the agent instructions document the kind, the Fact Pack, the fetch command and the rule that nothing publishes without the Author. A roundup-shaped Brief fixture ("What happened" lists five repositories linking their Profiles, "Why it matters" is commentary, `source_url` the trending page, `explainers` the Profiles) passes the existing Brief gates in the builder tests, proving the Friday roundup needs no new piece type.

Spec: `../spec.md` ("Editorial docs", User Stories 21, 29, 30, 37, 38). Editorial: `docs/editorial/repo-profile-selection.md` already holds floors, cadence and the first five.

**Blocked by:** 02 (Skeleton rule for Repo Profiles), 05 (Fetch command)

**Status:** ready-for-agent

- [ ] Template file passes the skeleton rule when filled with placeholder prose (a builder test builds it)
- [ ] Writing guide covers every section, the citing rule, quickstart verbatim, Alternatives and FAQ counts, title/description patterns, and links the selection guide and Editorial Policy
- [ ] Checkpoint-note template lists the metrics and the pass rule from ADR-0005
- [ ] README and agent instructions updated; agent instructions state the no-publish-without-approval rule for Profiles and Refresh pull requests
- [ ] Roundup Brief fixture builds clean under the existing Brief gates and is written at its Briefs URL

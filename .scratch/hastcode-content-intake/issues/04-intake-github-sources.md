# 04: Intake: GitHub repository Sources

**What to build:** A Source can be a GitHub repository; its new releases become Leads. Releases come through the GitHub REST API with the token the Fact Pack fetcher already uses, reusing its client (and its rate-limit handling); never the `releases.atom` web endpoint. Conditional requests apply as for feeds; an authenticated 304 does not count against the rate limit.

Spec: `../spec.md` ("Intake"). Research section 5 (REST limits, `releases.atom` concern).

**Blocked by:** 02 (Intake tracer)

**Status:** ready-for-agent

- [ ] Intake-seam tests against a fake GitHub (reuse the Fact Pack fake's approach): new releases listed as Leads with the Source's kind and Vertical; second run lists nothing new; a missing token gives a clear error for repository Sources only; prereleases handled per a documented rule
- [ ] No request ever goes to github.com web pages

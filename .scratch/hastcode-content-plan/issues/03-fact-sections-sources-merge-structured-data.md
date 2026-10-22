# 03: Fact sections, Sources merge and structured data

**What to build:** A Repo Profile page shows two sections the Author never writes. "Signals", inserted after Architecture, lists plain facts from the Fact Pack: last commit date, releases in the last twelve months and the latest release date, contributor count and the share of contributions held by the top ten, open issues, security policy present, license, stars, forks, watchers, and which documentation files and directories exist. No score, no colour grading. "Changelog", inserted after Alternatives, lists the latest releases with tag, name, date and link. The Fact Pack's source URLs join the Sources section once per URL, the way frontmatter sources do. The page's JSON-LD adds SoftwareSourceCode (name, repository URL, programming language, license URL, description) and FAQPage built from the FAQ H3s and their answers, and the allowed-types test admits them. The header card and Signals respect the dark theme.

Spec: `../spec.md` ("Page composition and gates", User Stories 6, 9, 12, 16, 27). Prior art: frontmatter sources merge test, Explainer JSON-LD tests, Mermaid figure insertion.

**Blocked by:** 01 (Repo Profile kind renders from a Fact Pack)

**Status:** ready-for-agent

- [ ] Signals appears after Architecture with every fact above, rendered from the fixture Fact Pack; a fact absent from the Fact Pack is shown as "not found", never invented
- [ ] Changelog appears after Alternatives with the latest releases newest first, prereleases marked, linking to the release page
- [ ] Sources lists each Fact Pack source URL once, even when the Author also cited it in the body or frontmatter
- [ ] JSON-LD carries Article as before plus SoftwareSourceCode and FAQPage; FAQPage questions equal the FAQ H3s in order; the no-unsupported-types test passes
- [ ] Header card and Signals have no hard-coded colours; both themes verified by the stylesheet rules used for other components
- [ ] Rendered HTML stays under the 50 KB gate for the fixture Profile

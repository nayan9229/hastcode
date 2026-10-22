# 08: First Repo Profile: apache/iggy on a Preview

**What to build:** The whole path once, for real. Run the fetch command for apache/iggy, start from the template, draft the opinion sections from the Fact Pack and README (Architecture with a Mermaid diagram of the server's shape, Use it if / Don't use it if, three to five Alternatives, five to eight FAQ questions phrased as searches, Verdict), pass the production build, and push a branch so a Preview is deployed. The agent drafts; the Author verifies every claim against its source, edits, and decides whether it merges. The remaining four of the first five (xerj-org/xerj, NVIDIA/OpenShell, deepseek-ai/DeepInfer, php/frankenphp) follow the same path as editorial routine, not as tickets.

Spec: `../spec.md` ("Further Notes": first five). Editorial: `docs/editorial/repo-profile-selection.md`, the writing guide from 07, `content/pages/editorial-policy.md`.

**Blocked by:** 03 (Fact sections, Sources merge and structured data), 04 (Generated hero card), 05 (Fetch command), 07 (Editorial template, writing guide, roundup proof and docs)

**Status:** ready-for-agent

- [ ] Fact Pack fetched with the command; `needs_review` false (Apache-2.0)
- [ ] Draft follows the template; every factual claim in prose cites a Fact Pack field, the README or a release page; quickstart commands are verbatim from the README
- [ ] Production build green: skeleton rule, section-length guide, 50 KB gate, link gate
- [ ] Preview deployed and its URL recorded in this ticket
- [ ] Author verification: every claim checked, draft flag removed only by the Author (ready-for-human step; the agent stops at the Preview)
- [ ] Day-28 checkpoint note scheduled (date written here) using the template from 07

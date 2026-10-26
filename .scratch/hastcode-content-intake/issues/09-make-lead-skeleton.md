# 09: make lead: from Lead to draft skeleton

**What to build:** Picking a Lead starts a draft with structure, not prose. `make lead ID=…` writes a skeleton at the kind's content path: frontmatter (title from the Lead, kind, Vertical, author mira-vance, `draft: true`, a provisional `published` (and `reviewed` for an Explainer) at the pick time for the Author to reset when publishing, no `verified`, no `hero`), a Sources section holding the Lead's link, and the kind's section headings from its template (ticket 08), then marks the Lead picked. Kind and Vertical can be overridden at pick time. A Brief names the Explainer(s) it links (`EXPLAINER=slug[,slug]`): written to `explainers` and linked in "Why it matters"; refused without one or for one that does not exist. It refuses an existing slug and refuses kind Brief when the Lead is more than 3 days old. README documents the daily rhythm (run Intake each morning, pick at most two Leads, dismiss the rest) and the planned pace (5 Explainers: 2 Case Studies, 2 Concepts, 1 Comparison; 5 Briefs a week).

Spec: `../spec.md` ("Picking", "Cadence").

**Blocked by:** 05 (dismiss, expiry, covered), 08 (templates and writing guides)

**Status:** resolved

- [x] Intake-seam tests: skeleton frontmatter, Sources link and headings per kind; override applied; existing slug refused; Brief from a 4-day-old Lead refused; the Lead marked picked and absent from later queues
- [x] A skeleton contains no generated prose and is a draft (absent from production builds)
- [x] README Intake section completed with the rhythm and pace

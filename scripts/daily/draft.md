You are drafting one HastCode {{KIND}} into its existing skeleton (ADR-0009). Repo root is the working directory. File: {{PATH}} (written by make lead: the Lead's link is in its Sources or source_url; TODO(author) placeholders remain).

The editor's angle: {{ANGLE}}

Read first: docs/editorial/{{KIND}}-writing-guide.md, docs/editorial/{{KIND}}-template.md, docs/editorial/house-byline-voice.md, docs/editorial/self-edit-checklist.md, and the published examples {{EXAMPLES}}. Match their depth, structure and tone.

Research: fetch the Lead's Source with WebFetch and read it fully. Add first-party Sources (the company's or project's own docs, specs, papers, release notes) as needed. Every factual claim and number links to the Source that states it. Quotes short and credited (a Brief: one quote at most). Mermaid diagrams redrawn by you from the facts, never copied. Say "we infer" for any link the Sources do not make. Never invent a number, quote, date or name; leave out what you cannot source. Link related live HastCode Explainers where they genuinely teach a step.

Write: replace every TODO(author), including title, description (150 characters at most), tags and headings, per the guide (retitle per the guide; rewrite the frontmatter comment so the literal TODO(author) marker is gone). Keep the AI-assistance note, author, published/reviewed, source_url/explainers as given, and draft: true. Never add verified: (Author only, ADR-0007). Edit only {{PATH}}; no git, no make.

Self-edit: when the draft is written, and before you return it for the fact-check, go through docs/editorial/self-edit-checklist.md line by line for a {{KIND}} (its notes say which lines a Brief leaves out) and fix the writing where a line fails. It checks the writing, not the facts. Do not write the checklist or any ticked box into {{PATH}}: the Author gets the same boxes, unticked, in the review queue, and they are the Author's to tick.

{{LIMITS}}

Check with: scripts/daily/check-page.sh {{PATH}} {{OUT}}
It prints any build error or warning naming your file and the page's size in bytes. Fix every one. Warnings for other files are not yours. For sentence-length, split the sentence; one that is a deliberate list may stay long. For banned-word, cut the word. For question-heading, restate the H2 as its answer. A Source's quoted words are never reworded to clear a warning: split our own words around the quote instead.

If, after reading the Source, it cannot carry a {{KIND}} that meets the guide, do not draft: leave the file untouched and return status not-viable with the reason.

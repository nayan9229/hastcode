---
# Case Study template, drawn from netflix-single-tenant-titus. `make lead`
# writes it as a skeleton at content/explainers/<vertical>/<slug>.md; to
# start by hand, copy it there. Follow
# docs/editorial/case-study-writing-guide.md. Replace every TODO(author)
# line: a production build fails while one is left. `verified:` is added
# by the Author alone, after checking every claim against its Source
# (ADR-0007); an agent never writes it.
title: "TODO(author): {system}: how {company} does it, in the words a Reader searches"
description: "TODO(author): one sentence, 150 characters at most: what the system does and the trade-off that matters."
kind: case-study
vertical: "TODO(author): the one Vertical slug it fits, e.g. infrastructure-cloud"
tags: ["TODO(author): two to five lower-case tags"]
published: "TODO(author): RFC 3339 with offset, e.g. 2026-10-23T09:00:00Z"
reviewed: "TODO(author): the same as published"
author: mira-vance
# Optional: an image HastCode may use, beside this file, replaces the
# generated title card. Never a Source's image.
# hero: system-name.png
# hero_alt: "What the image shows"
# Optional: related Explainer slugs and extra sources (title + url).
# related: [metastable-failures]
draft: true
---

TODO(author): the hook, one or two sentences: a sourced fact about {system} (a number, a design choice, an incident) that makes a backend engineer look twice, linked to the Lead.

> Drafted with AI assistance, researched, verified and edited by the HastCode editorial team. We publish only as many pieces as the Author can verify this way.

## TL;DR

- TODO(author): what {system} is and what it does, in one plain sentence, linked.
- TODO(author): why {company} built it: the problem the old design could not solve.
- TODO(author): how it works, naming the parts the diagram shows, in order.
- TODO(author): the price it pays, and the result that surprised us.

## TODO(author): {system} at a glance, named for what the diagram shows

TODO(author): say the diagram is ours, redrawn from the Source, and link the original figure. Then one paragraph on the path the diagram follows, top to bottom, and which section walks each part.

```mermaid
%% title: TODO(author): what the diagram shows (redrawn from {company}'s description)
flowchart TD
  Request --> Decide --> Run
```

## TODO(author): the problem {company} had to solve

TODO(author): what the old design did, what broke, and who paid for it, with the Source's own evidence linked.

## TODO(author): the options {company} weighed

TODO(author): each alternative the Source names, why it was rejected in the Source's words, linked, and the argument that decided it.

## TODO(author): the first part of the diagram and what it does

TODO(author): how this part works, linked claim by claim. Say "we infer" wherever the Source leaves a gap.

## TODO(author): the next part of the diagram and what it does

TODO(author): how this part works, linked claim by claim. Add one H2 per part of the diagram, in the diagram's order.

## TODO(author): what it costs

TODO(author): the trade-offs the Source admits, plainly, and how {company} manages each. Label any number of our own as ours.

## TODO(author): the evidence it works, and what the Source leaves out

TODO(author): what {company} measured or reports, linked, and what it does not publish (no benchmark numbers, no failure counts).

## TODO(author): what it makes possible

TODO(author): who uses it now, what is next, and the lesson for a Reader running their own systems: what could go right, without hiding the cost.

## Sources

- TODO(author): the Lead first, as [{title}](https://example.com/the-lead), its authors, publisher and date; then every other page cited, one per line.

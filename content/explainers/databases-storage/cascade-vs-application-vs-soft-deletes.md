---
# Skeleton from Lead ddftdwnr (PlanetScale blog), picked 2026-10-27.
# Follow docs/editorial/comparison-writing-guide.md and replace every
# TODO(author) line. published and reviewed are the pick time, provisional: set them
# when you publish. verified: is the Author's alone (ADR-0007).
title: "I don't trust ON DELETE CASCADE"
description: "TODO(author): one sentence, 150 characters at most: the options, the question they answer differently, and who should pick which."
kind: comparison
vertical: databases-storage
tags: ["TODO(author): two to five lower-case tags"]
published: "2026-10-27T17:47:19Z"
reviewed: "2026-10-27T17:47:19Z"
author: mira-vance
# Optional: an image HastCode may use, beside this file, replaces the
# generated title card. Never a Source's image.
# hero: option-a-vs-option-b.png
# hero_alt: "What the image shows"
# Optional: related Explainer slugs and extra sources (title + url).
# related: [metastable-failures]
draft: true
---

TODO(author): the hook, a few sentences: one sourced fact (a measurement, a release, a design change) that makes the choice between {option A} and {option B} matter, linked to the Lead. End on the question this piece answers for every option.

> Drafted with AI assistance, researched, verified and edited by the HastCode editorial team. We publish only as many pieces as the Author can verify this way.

## TL;DR

- TODO(author): {option A} in one plain sentence: how it does the thing being compared.
- TODO(author): {option B} in one plain sentence, the same way.
- TODO(author): the dimension where they differ most.
- TODO(author): which to pick when, labelled "we infer" unless a Source says it.

## TODO(author): the options side by side, named for what the diagram shows

TODO(author): what the options have in common, then one paragraph walking the diagram: the same request or workload through each option.

```mermaid
%% title: TODO(author): the same workload through each option
flowchart LR
  subgraph A["Option A"]
    A1["Request"] --> A2["Path A"]
  end
  subgraph B["Option B"]
    B1["Request"] --> B2["Path B"]
  end
```

## TODO(author): {option A}: how it works

TODO(author): how {option A} works, from its specification or documentation, linked claim by claim.

## TODO(author): {option B}: how it works

TODO(author): how {option B} works, the same way. Add one H2 per option.

## TODO(author): the first dimension, answered for every option

TODO(author): one question (a failure mode, a limit, a cost) answered for each option in turn, in the same order every time. Add one H2 per dimension; each becomes a row of the table.

## Trade-offs

Each cell summarises a section above, where its source is linked.

| Question | {option A} | {option B} |
|---|---|---|
| TODO(author): the first dimension | TODO(author): short answer | TODO(author): short answer |
| TODO(author): the next dimension | TODO(author): short answer | TODO(author): short answer |

## TODO(author): when each one wins

TODO(author): one paragraph per option: the situation where it is the right choice and what to watch. Say plainly that this section is our inference unless a Source states it.

## Sources

- [I don't trust ON DELETE CASCADE](https://planetscale.com/blog/i-dont-like-delete-cascade), PlanetScale blog

---
# Brief template, drawn from github-dependabot-fan-out. `make lead
# KIND=brief` writes it as a skeleton at content/briefs/<slug>.md, only
# from a Lead at most 3 days old; to start by hand, copy it there. Follow
# docs/editorial/brief-writing-guide.md. Replace every TODO(author) line:
# a production build fails while one is left. `verified:` is added by the
# Author alone, after checking every claim against its Source (ADR-0007);
# an agent never writes it.
title: "TODO(author): {company}'s {event}: the architectural lesson, in a few words"
description: "TODO(author): one sentence, 150 characters at most: what happened and the lesson."
tags: ["TODO(author): two to four lower-case tags"]
published: "TODO(author): RFC 3339 with offset, at most 3 days after the Lead, e.g. 2026-10-26T18:00:00Z"
author: mira-vance
source_url: "TODO(author): the Lead's link, the first-party page the Brief is about"
explainers: ["TODO(author): the slug of every Explainer the body links, and no other"]
# Optional: an image HastCode may use, beside this file, replaces the
# generated title card. Never a Source's image.
# hero: event-name.png
# hero_alt: "What the image shows"
# Optional: related Explainer slugs and extra sources (title + url).
# related: [metastable-failures]
draft: true
---

> Drafted with AI assistance, researched, verified and edited by the HastCode editorial team. We publish only as many pieces as the Author can verify this way.

## What happened

TODO(author): the facts in our own words, each linked to the Source: who, what, when (with time zone). One short credited quote at most, only where the exact words matter.

## Why it matters

TODO(author): our own analysis, never the Source restated: the architectural cause, the lesson for a Reader's systems, and a link to the Explainer that teaches it, as [the idea]({explainer}). Say "we infer" for anything the Source does not state.

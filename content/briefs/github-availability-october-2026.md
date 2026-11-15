---
# From Lead lh66psks (GitHub availability reports), picked 2026-11-15.
# published is the pick time, provisional: set it when you publish.
# verified: is the Author's alone (ADR-0007).
title: "GitHub's October incidents: clients that never backed off"
description: "In October, retries, reconnects and an import burst that did not back off worsened three GitHub faults. Backoff and rate limits break the loop."
tags: [retries, backoff, incidents, github]
published: "2026-11-15T07:09:16Z"
author: mira-vance
source_url: "https://github.blog/news-insights/company-news/github-availability-report-october-2026/"
explainers: [metastable-failures]
# Optional: an image HastCode may use, beside this file, replaces the
# generated title card. Never a Source's image.
# hero: event-name.png
# hero_alt: "What the image shows"
# Optional: related Explainer slugs and extra sources (title + url).
# related: [metastable-failures]
# sources: [{title: "The page's title", url: "https://example.com/page"}]
draft: true
---

> Drafted with AI assistance, researched, verified and edited by the HastCode editorial team. We publish only as many pieces as the Author can verify this way.

## What happened

GitHub's [October 2026 availability report](https://github.blog/news-insights/company-news/github-availability-report-october-2026/), posted 12 November, describes eight incidents, three involving load that did not back off. On 6 October (19:57 UTC) an Enterprise Importer burst of very large queries used all of a shared database's connections; the importer did not reduce its load. On 19 October (15:49 UTC) an Actions service lost connections to some cache servers and all its servers retried at once, saturating a shared network link. On 21 October (14:52 UTC) requests failed to event-streaming servers near their memory limits and new connections opened immediately; overloaded servers' traffic moved to others, which overloaded too.

## Why it matters

In each case, we infer, load that did not back off made the fault worse. On 21 October, shifted traffic overloaded the servers that received it. That is the loop [metastable failures](/distributed-systems/metastable-failures/) describes: load the failure creates keeps the system down after the trigger passes. We infer the link; the report does not use the term.

The [fixes](https://github.blog/news-insights/company-news/github-availability-report-october-2026/) target the loop: a cap on how fast the Actions service opens connections, randomized reconnect delays, and importer backoff, still being automated. GitHub says connection limits and other client changes cut the risk of ["a broader reconnect storm"](https://github.blog/news-insights/company-news/github-availability-report-october-2026/). Give your retrying clients capped, randomized backoff and a connection-rate limit, then test that a recovered dependency stays up.

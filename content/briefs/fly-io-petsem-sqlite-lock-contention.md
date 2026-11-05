---
# From Lead 3wlmobhd (Fly.io Infra Log), picked 2026-11-05.
# published is the pick time, provisional: set it when you publish.
# verified: is the Author's alone (ADR-0007).
title: "Fly.io's flaky Petsem: SQLite has one writer, so keep writes bounded"
description: "Extra tracking writes tipped Petsem's SQLite primary into lock contention, failing over half of writes. SQLite has one writer: bound the load."
tags: [sqlite, locking, incidents, fly-io]
published: "2026-11-05T17:02:45Z"
verified: 2026-11-05
author: mira-vance
source_url: "https://fly.io/infra-log/petsem-was-flaky/"
explainers: [cloudflare-d1]
# Optional: an image HastCode may use, beside this file, replaces the
# generated title card. Never a Source's image.
# hero: event-name.png
# hero_alt: "What the image shows"
# Optional: related Explainer slugs and extra sources (title + url).
# related: [metastable-failures]
sources: [{title: "Appropriate Uses For SQLite", url: "https://www.sqlite.org/whentouse.html"}, {title: "Write-Ahead Logging", url: "https://www.sqlite.org/wal.html"}]
---

> Drafted with AI assistance, researched, verified and edited by the HastCode editorial team. We publish only as many pieces as the Author can verify this way.

## What happened

Fly.io's [infra log entry of 4 November 2026](https://fly.io/infra-log/feed.xml), headed [15 October, 09:16 UTC](https://fly.io/infra-log/petsem-was-flaky/), says its secrets store, Petsem, runs on a SQLite primary in IAD, its [Ashburn, Virginia region](https://docs.fly.io/reference/regions). A deployed change, groundwork for a possible move off SQLite, [added record tracking](https://fly.io/infra-log/petsem-was-flaky/). The extra writes, [Fly.io says](https://fly.io/infra-log/petsem-was-flaky/), "tipped the database into sustained lock contention": a little over half of write requests failed, while replicas and reads from the primary stayed healthy. A [new build without the tracking code](https://fly.io/infra-log/petsem-was-flaky/) stopped the errors.

## Why it matters

SQLite allows [unlimited readers but one writer at a time](https://www.sqlite.org/whentouse.html). Writers queue behind one lock, and the queue drains only while they arrive slower than it serves them. Readers and writers [do not block each other in write-ahead log (WAL) mode](https://www.sqlite.org/wal.html), which fits healthy reads beside failing writes; we infer the mode; the post does not say. Fly.io's fix was fewer writes, not a bigger primary, and that is the design rule: a SQLite primary's write load must stay small and bounded. Cloudflare's D1 builds on the same rule: each database is single-threaded, so an application [scales by adding databases](/databases-storage/cloudflare-d1/), never by writing more to one. If one SQLite file backs your service, treat its write rate as a budget and cost every change against it: tracking rows are writes too.

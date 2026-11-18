---
# From Lead 6oimuyp5 (Fly.io Infra Log), picked 2026-11-18.
# published is the pick time (make daily may move it to a Gap Day);
# make release keeps it (ADR-0012). verified: is the Author's alone (ADR-0007).
title: "Fly.io's connection exhaustion: many pools, no shared limit"
description: "Fly.io's RDS database ran out of connection slots, so some dashboard and API calls failed. Cap connections in a pooler, not only in each process."
tags: [connection-pooling, backpressure, incidents, fly-io]
published: "2026-11-18T06:32:27Z"
author: mira-vance
source_url: "https://fly.io/infra-log/database-connection-exhaustion/"
explainers: [cloudflare-hyperdrive-mysql-caching]
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

Fly.io's [infra log entry of 16 November 2026](https://fly.io/infra-log/feed.xml), headed [27 October, 13:12 UTC](https://fly.io/infra-log/database-connection-exhaustion/), says the RDS database behind its main web app ran out of connection slots. Some dashboard and API requests returned 500 and 503 errors; existing Machines and apps were unaffected. Many Rails processes each kept their own connection pool, and background worker machines held [around 2,000 connections](https://fly.io/infra-log/database-connection-exhaustion/), most of the total. Moving container registry authentication to a read replica freed [a little over 10%](https://fly.io/infra-log/database-connection-exhaustion/) of the primary's connections.

## Why it matters

A pool inside a process limits that process, not the fleet. Peak demand is processes times pool size, so adding processes for a spike uses more of the primary's slots. [Fly.io says](https://fly.io/infra-log/database-connection-exhaustion/) only the database itself enforces any limit. That cap refuses whoever asks next, a user's request or a periodic job. Fly.io [plans safe limits](https://fly.io/infra-log/database-connection-exhaustion/) for its workers.

Put the limit in a shared pooler instead. In transaction mode, PgBouncer [assigns a database connection for one transaction only](https://www.pgbouncer.org/features.html), so many clients share few slots. [RDS Proxy queues or throttles](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/rds-proxy.html) connections it cannot serve at once. Fly.io, [looking into a pooler](https://fly.io/infra-log/database-connection-exhaustion/), names both. [Cloudflare's Hyperdrive](cloudflare-hyperdrive-mysql-caching) shows the pattern and its cost: a [`SET` lasts only for its transaction or query](https://developers.cloudflare.com/hyperdrive/concepts/how-hyperdrive-works/#pooling-mode). Multiply your processes by pool size at full scale, then compare with your database's connection limit.

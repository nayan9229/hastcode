---
# From Lead b7bibqjl (AWS Database Blog), picked 2026-11-10. Drafted per
# docs/editorial/brief-writing-guide.md. published is the pick time,
# provisional: set it when you publish. verified: is the Author's alone
# (ADR-0007).
title: "Topstep's Redis to Valkey cutover: the DNS TTL is the write window"
description: "Topstep cut about 53 million keys over to ElastiCache for Valkey with one Route 53 change. A DNS flip freezes nothing: the TTL is the write window."
tags: [valkey, elasticache, dns, migration]
published: "2026-11-10T06:56:04Z"
verified: 2026-11-10
author: mira-vance
source_url: "https://aws.amazon.com/blogs/database/topsteps-migration-journey-self-managed-redis-oss-to-amazon-elasticache-for-valkey/"
explainers: [planned-database-switchover]
# Optional: an image HastCode may use, beside this file, replaces the
# generated title card. Never a Source's image.
# hero: event-name.png
# hero_alt: "What the image shows"
# Optional: related Explainer slugs and extra sources (title + url).
# related: [metastable-failures]
sources:
  - {title: "RedisShake documentation: Sync Reader", url: "https://tair-opensource.github.io/RedisShake/en/reader/sync_reader.html"}
  - {title: "Amazon Route 53 Developer Guide: Values specific for simple records, TTL (seconds)", url: "https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/resource-record-sets-values-basic.html#rrsets-values-basic-ttl"}
---

> Drafted with AI assistance, researched, verified and edited by the HastCode editorial team. We publish only as many pieces as the Author can verify this way.

## What happened

On 9 November, AWS's Database Blog [described how Topstep](https://aws.amazon.com/blogs/database/topsteps-migration-journey-self-managed-redis-oss-to-amazon-elasticache-for-valkey/) moved about 53 million keys from one self-managed Redis OSS 7.4 node with no replica to Multi-AZ ElastiCache for Valkey 7.2.6. The post says Topstep's source ["falls outside the Redis OSS ≤ 7.2 range that the managed migration supports"](https://aws.amazon.com/blogs/database/topsteps-migration-journey-self-managed-redis-oss-to-amazon-elasticache-for-valkey/), so [RedisShake](https://tair-opensource.github.io/RedisShake/en/reader/sync_reader.html), in PSYNC mode, replicated it [until lag was near zero](https://aws.amazon.com/blogs/database/topsteps-migration-journey-self-managed-redis-oss-to-amazon-elasticache-for-valkey/). Then [one Route 53 change](https://aws.amazon.com/blogs/database/topsteps-migration-journey-self-managed-redis-oss-to-amazon-elasticache-for-valkey/) turned an A record into a CNAME with a 30-second TTL; services followed within seconds, and the old node stayed up as a warm standby.

## Why it matters

A DNS flip is a planned switchover with no write freeze. Resolvers [cache a record for its TTL](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/resource-record-sets-values-basic.html#rrsets-values-basic-ttl), so for up to 30 seconds some clients could still write to the old primary (we infer; the post does not say). Our comparison of [planned switchovers](/databases-storage/planned-database-switchover/) calls that window unavoidable: stop the primary, pause writes, or buffer them at a proxy. Topstep's flip leaves it open for one TTL; we infer stragglers follow only while replication runs, then they are lost. [Two choices made the cutover a record edit](https://aws.amazon.com/blogs/database/topsteps-migration-journey-self-managed-redis-oss-to-amazon-elasticache-for-valkey/): a hostname already in front of the cache, and a warm rollback target left running. Add the hostname before you need it, and price a write in the window before you flip.

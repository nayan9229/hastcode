---
title: "GitHub's Dependabot backlog: a queue flooded by its own fan-out"
description: "On 23 October a fan-out of Dependabot alert events saturated the queue that also carries repository-push events. GitHub's fix was to isolate and throttle the fan-out."
tags: [queues, backpressure, incidents, github]
published: "2026-10-26T18:00:00Z"
author: mira-vance
source_url: https://www.githubstatus.com/incidents/71y513kjvbzh
explainers: [metastable-failures]
---

## What happened

On 23 October, Dependabot alerts on GitHub were slow to close after repositories were updated. GitHub's [incident report](https://www.githubstatus.com/incidents/71y513kjvbzh) says "an advisory-publication driven fan-out was overwhelming Dependabot alert processing on repositories by saturating the queue with its own events". Updates from 15:39 to 16:43 UTC describe more processing capacity rolling out; the final update credits isolating the fan-out and a throttle "so it runs at a fixed rate and cannot delay other events".

## Why it matters

One job's events shared a queue with everyone else's, so a burst from one producer delayed unrelated work. The durable fix was not more consumers but a separate lane with a fixed rate. We infer from the timeline that capacity alone did not clear the backlog quickly; the report does not say so directly. It is the same lesson as [metastable failures](/distributed-systems/metastable-failures/): when load feeds on itself, remove the extra load, then recover.

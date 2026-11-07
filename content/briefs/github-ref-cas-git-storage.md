---
# From Lead wrui62wm (GitHub Engineering), picked 2026-11-07.
# published is the pick time, provisional: set it when you publish.
# verified: is the Author's alone (ADR-0007).
title: "GitHub's Project Rubicon: only the pointers need consensus"
description: "GitHub's new Git storage keeps packs in Blob Storage and votes only on reference pointers: 7,200 pushes a second to one repository in its tests."
tags: [git, consensus, object-storage, github]
published: "2026-11-07T06:09:09Z"
author: mira-vance
source_url: "https://github.blog/engineering/architecture-optimization/inside-githubs-new-git-storage-architecture/"
explainers: [github-agent-scale-git]
# Optional: an image HastCode may use, beside this file, replaces the
# generated title card. Never a Source's image.
# hero: event-name.png
# hero_alt: "What the image shows"
# Optional: related Explainer slugs and extra sources (title + url).
# related: [metastable-failures]
sources: [{title: "Pro Git: Git Objects", url: "https://git-scm.com/book/en/v2/Git-Internals-Git-Objects"}]
draft: true
---

> Drafted with AI assistance, researched, verified and edited by the HastCode editorial team. We publish only as many pieces as the Author can verify this way.

## What happened

On 5 November, GitHub [described Project Rubicon](https://github.blog/engineering/architecture-optimization/inside-githubs-new-git-storage-architecture/), its next Git storage architecture, now in limited private preview. A virtual replica [indexes each push's pack and uploads it to Azure Blob Storage](https://github.blog/engineering/architecture-optimization/inside-githubs-new-git-storage-architecture/), then commits one transaction to Ref-CAS, a replicated log holding references and the pack manifest, never object data. Each update is [a compare-and-swap, visible once a majority of log replicas record it](https://github.blog/engineering/architecture-optimization/inside-githubs-new-git-storage-architecture/); the latest log sequence number is the repository's version. Compaction and garbage collection [run on dedicated background nodes, never on a virtual replica](https://github.blog/engineering/architecture-optimization/inside-githubs-new-git-storage-architecture/). In internal tests one repository [took 7,200 pushes a second from 5,000 concurrent agents](https://github.blog/engineering/architecture-optimization/inside-githubs-new-git-storage-architecture/).

## Why it matters

Git [names objects by content](https://git-scm.com/book/en/v2/Git-Internals-Git-Objects), so a pack can be durable before anyone agrees on it; a vote only settles which commit a reference names. That is the design: ["Only reference updates need agreement"](https://github.blog/engineering/architecture-optimization/inside-githubs-new-git-storage-architecture/), as GitHub's Karthik Raman and Brian Celenza put it. Our Explainer on [the per-reference quorum](/infrastructure-cloud/github-agent-scale-git/), where each vote still covered a full copy, left open what a write pays when the blob-backed replica is slow. Rubicon half answers it: the upload [precedes the commit](https://github.blog/engineering/architecture-optimization/inside-githubs-new-git-storage-architecture/), so we infer a slow one delays only that push, not its neighbours; it gives no latency figures. If your store votes on bulk data, find what your consistency rule actually names, and vote on that alone.

---
# Skeleton from Lead fd3v3vtj (Meta Engineering, ML Applications), picked 2026-10-27.
# Follow docs/editorial/brief-writing-guide.md. published is the pick time,
# provisional: set it when you publish. verified: is the Author's alone (ADR-0007).
title: "Meta's OCP 2026 networking: three fabrics for three kinds of traffic"
description: "At the 2026 OCP Global Summit, Meta showed separate scale-up, scale-out and scale-across networks, each built for the traffic it carries."
tags: [datacenter-networking, ai-infrastructure, ethernet, meta]
published: "2026-10-28T12:54:21Z"
verified: 2026-10-28
author: mira-vance
source_url: "https://engineering.fb.com/2026/10/26/open-source/ocp-summit-2026-meta-networking-talks/"
explainers: [meta-prometheus-supercluster]
# Optional: an image HastCode may use, beside this file, replaces the
# generated title card. Never a Source's image.
# hero: event-name.png
# hero_alt: "What the image shows"
# Optional: related Explainer slugs and extra sources (title + url).
# related: [metastable-failures]
sources: [{title: "Inside Prometheus: Meta's Journey to a Gigawatt-Scale Supercluster", url: "https://engineering.fb.com/2026/10/19/data-center-engineering/inside-prometheus-meta-gigawatt-scale-supercluster/"}]
---

> Drafted with AI assistance, researched, verified and edited by the HastCode editorial team. We publish only as many pieces as the Author can verify this way.

## What happened

On 26 October, Meta [published the slides and videos](https://engineering.fb.com/2026/10/26/open-source/ocp-summit-2026-meta-networking-talks/) of its networking talks at the 2026 OCP Global Summit. For scale-up inside the rack, it showed Tahan, an Ethernet switch for its MTIA 300 system. For scale-out, its Flattened Planar Fabric connects each GPU to eight network planes and flattens a three-tier Clos to two. For scale-across, its BAG super-spine joins buildings, campuses and regions. Its MetaRoCE RDMA transport is going open source.

## Why it matters

We infer that Meta runs three networks, not one, because training traffic is not uniform; the talks post does not say so. Meta's [Prometheus post](https://engineering.fb.com/2026/10/19/data-center-engineering/inside-prometheus-meta-gigawatt-scale-supercluster/) says latency- and bandwidth-sensitive operations concentrate within a model replica, while traffic between replicas tolerates more latency. On our reading, that lets each tier trade bandwidth against distance differently: densest in the rack, flatter across a building, longest between sites. The [MetaRoCE talk](https://engineering.fb.com/2026/10/26/open-source/ocp-summit-2026-meta-networking-talks/) shows the stakes: collective operations run only as fast as the slowest node. Our Case Study on [Prometheus](meta-prometheus-supercluster) walks through each tier and its cost. For your own systems, map where traffic actually flows first, then spend on the hot paths and oversubscribe the cold ones.

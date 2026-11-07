---
# From Lead svnk2r4u (Fly.io blog, published 2026-11-06 per its feed; the
# page shows no date), picked 2026-11-07. Drafted per
# docs/editorial/brief-writing-guide.md. published is the pick time,
# provisional: set it when you publish. verified: is the Author's alone
# (ADR-0007).
title: "Fly.io's fleet reboot: a VM is only as strong as the host's parser"
description: "A kernel bug in IPv4 checksum offload let a guest rewrite a routed header. Fly.io rebooted its fleet in waves; isolation ends at the host's parser."
tags: [isolation, virtualization, incidents, fly-io]
published: "2026-11-07T06:09:09Z"
author: mira-vance
source_url: "https://fly.io/blog/well-that-wasnt-fun/"
explainers: [netflix-single-tenant-titus, openshell]
# Optional: an image HastCode may use, beside this file, replaces the
# generated title card. Never a Source's image.
# hero: event-name.png
# hero_alt: "What the image shows"
# Optional: related Explainer slugs and extra sources (title + url).
# related: [metastable-failures]
sources: [{title: "CVE-2026-99999: ipv4: validate checksum_start before completing checksum", url: "https://www.cve.org/CVERecord?id=CVE-2026-99999"}, {title: "Fly.io architecture", url: "https://docs.fly.io/reference/architecture"}]
draft: true
---

> Drafted with AI assistance, researched, verified and edited by the HastCode editorial team. We publish only as many pieces as the Author can verify this way.

## What happened

On 6 November, Fly.io [described 16 October](https://fly.io/blog/well-that-wasnt-fun/): every host patched and rebooted in 10 hours. [CVE-2026-99999](https://www.cve.org/CVERecord?id=CVE-2026-99999), published a day earlier, is a wrong IPv4 check on checksum offload, the metadata a guest's virtual network card hands the host saying where a checksum reads and writes. A guest aiming it at a header the host had already routed and filtered [rewrites two bytes](https://fly.io/blog/well-that-wasnt-fun/): a firewall bypass, or 32 bytes of host memory, a neighbouring VM's packets included. Hosts take [around 15 to 20 minutes](https://fly.io/blog/well-that-wasnt-fun/) to reboot, in waves so no Managed Postgres cluster rebooted both nodes at once; `iad` needed 17.

## Why it matters

Fly.io's wall between customers is hardware virtualization, [Firecracker microVMs](https://docs.fly.io/reference/architecture). But a VM still reaches the host through device paths, and here the host [trusted guest metadata after routing and filtering](https://fly.io/blog/well-that-wasnt-fun/): the wall is only as strong as the host's parsing of what the guest hands it. Fly.io's own lesson: ["Assume every boundary leaks a little"](https://fly.io/blog/well-that-wasnt-fun/). Netflix's [single-tenant Titus](netflix-single-tenant-titus) chose dedicated VMs for resource isolation, not hostile tenants; hostile tenants add the host's device emulation to the attack surface. [OpenShell](openshell) keeps real credentials outside the sandbox and freezes the agent if the trusted side disconnects, so a 32-byte leak finds placeholders. A reboot patch is a planned outage: stage waves so no pair's two halves fall together.

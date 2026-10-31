---
# From Lead 4xabvz3v (Cloudflare blog), picked 2026-10-31. Drafted per
# docs/editorial/brief-writing-guide.md. verified: is the Author's alone
# (ADR-0007).
title: "eMudhra's weak CA keys: trust stores that cannot drop a broken CA"
description: "Cloudflare rejected eMudhra's weak CA keys within hours; a week on, most large-network probes still accepted them. Trust stores must update fast."
tags: [tls, pki, trust-anchors, cloudflare]
published: "2026-10-31T22:05:29+05:30"
verified: 2026-10-31
author: mira-vance
source_url: "https://blog.cloudflare.com/emudhra-patching/"
explainers: [dnssec-root-key-rollover]
# Optional: an image HastCode may use, beside this file, replaces the
# generated title card. Never a Source's image.
# hero: event-name.png
# hero_alt: "What the image shows"
# Optional: related Explainer slugs and extra sources (title + url).
# related: [metastable-failures]
# sources: [{title: "The page's title", url: "https://example.com/page"}]
---

> Drafted with AI assistance, researched, verified and edited by the HastCode editorial team. We publish only as many pieces as the Author can verify this way.

## What happened

On 19 October 2026, Cloudflare [learned](https://blog.cloudflare.com/emudhra-patching/) that eMudhra's intermediate certificates had eight weak keys; one was the number 4. Anyone could [recover](https://blog.cloudflare.com/emudhra-patching/) those private keys from the public certificates and forge certificates that unpatched clients accept. Cloudflare [patched its Pingora origin-facing proxies within hours](https://blog.cloudflare.com/emudhra-patching/). Paired probes from some of its challenge pages, [50 million over a week](https://blog.cloudflare.com/emudhra-patching/), found more than half on large networks, and 30–40% of embedded browsers such as Android WebView, still accepting the broken chain.

## Why it matters

A server operator can distrust a CA with one deploy. A client that carries its own trust store waits for its vendor's update, so the window to forge stays open long after the fix. Cloudflare [points to](https://blog.cloudflare.com/emudhra-patching/) projects like upki that understand "a trust store can’t be a static list of root certificate authorities".

The post never mentions DNSSEC, but we infer this is the same stale-anchor long tail the [root KSK rollover](dnssec-root-key-rollover) shows: an old anchor costs nothing until the day it must change. DNSSEC's answer is in-band, authenticated updates, where a key the resolver already trusts vouches for the change. If you ship clients, give their trust store that kind of update path, and measure acceptance from the client side, as Cloudflare's probes did, not from your servers' patch status.

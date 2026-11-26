---
# From Lead 3wnuachg (GitHub status), picked 2026-11-26. published is the
# pick time (make daily may move it to a Gap Day); make release keeps it
# (ADR-0012). verified: is the Author's alone (ADR-0007).
title: "GitHub's Git LFS overload: retries hid one region's shortfall"
description: "On 25 November automated unauthenticated traffic overwhelmed Git LFS servers in one GitHub region. Shed load by class before requests time out."
tags: [load-shedding, retries, incidents, github]
published: "2026-11-26T07:37:21Z"
verified: 2026-11-26
author: mira-vance
source_url: "https://www.githubstatus.com/incidents/m0wvf9pt3l0l"
explainers: [metastable-failures]
# Optional: an image HastCode may use, beside this file, replaces the
# generated title card. Never a Source's image.
# hero: event-name.png
# hero_alt: "What the image shows"
# Optional: related Explainer slugs and extra sources (title + url).
# related: [metastable-failures]
# sources: [{title: "The page's title", url: "https://example.com/page"}]
---

## What happened

GitHub's [incident report](https://www.githubstatus.com/incidents/m0wvf9pt3l0l) says Git Large File Storage (LFS) operations on github.com were degraded from 21:30 UTC on 25 November to 00:28 UTC the next day. On average [2% of Git LFS requests failed](https://www.githubstatus.com/incidents/m0wvf9pt3l0l), 4% at the peak. Many more [failed first and passed only when Git LFS clients retried](https://www.githubstatus.com/incidents/m0wvf9pt3l0l) automatically, slowing uploads and downloads. The report blames heavy automated, unauthenticated traffic, which overwhelmed one region's Git LFS servers and left them ["unable to respond to requests in time"](https://www.githubstatus.com/incidents/m0wvf9pt3l0l). GitHub [moved most Git LFS traffic](https://www.githubstatus.com/incidents/m0wvf9pt3l0l) to regions with more capacity.

## Why it matters

The overload was regional. Other regions had [more capacity](https://www.githubstatus.com/incidents/m0wvf9pt3l0l), yet we infer it did not help until GitHub moved traffic to it.

Retries hid the overload's size. We infer the 2% [leaves out every request a retry saved](https://www.githubstatus.com/incidents/m0wvf9pt3l0l), so an error rate counted after retries understates an overload. We also infer each retry added load to servers already past their limit.

The report [describes late responses](https://www.githubstatus.com/incidents/m0wvf9pt3l0l) and mentions no refusals. We infer the servers admitted more work than they could finish. A timed-out request wastes the server's work, then returns as a retry: the loop behind [metastable failures](metastable-failures). The defence is to shed load: admit only what you can finish before the client gives up.

GitHub is working to [limit what heavy unauthenticated traffic can do to Git LFS](https://www.githubstatus.com/incidents/m0wvf9pt3l0l) and to add capacity in that region. Count first-try failures per region, and decide which traffic class each region refuses first.

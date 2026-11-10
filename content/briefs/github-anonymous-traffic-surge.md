---
# From Lead iswkkkhv (GitHub status), picked 2026-11-10. published is the
# pick time, provisional: set it when you publish. verified: is the
# Author's alone (ADR-0007).
title: "GitHub's anonymous-traffic surge: shed the cheapest class first"
description: "On 9 November anonymous requests to GitHub got 504s while logged-in traffic stayed healthy. Admit by traffic class and shed the cheapest first."
tags: [load-shedding, admission-control, incidents, github]
published: "2026-11-10T06:34:59Z"
author: mira-vance
source_url: "https://www.githubstatus.com/incidents/c9px7cw2ghrs"
explainers: [metastable-failures]
# Optional: an image HastCode may use, beside this file, replaces the
# generated title card. Never a Source's image.
# hero: event-name.png
# hero_alt: "What the image shows"
# Optional: related Explainer slugs and extra sources (title + url).
# related: [metastable-failures]
sources: [{title: "Increase in timeouts for some unauthenticated users", url: "https://www.githubstatus.com/incidents/pyynvrm1vskw"}]
draft: true
---

> Drafted with AI assistance, researched, verified and edited by the HastCode editorial team. We publish only as many pieces as the Author can verify this way.

## What happened

On 9 November GitHub [investigated degraded performance from 05:59 UTC](https://www.githubstatus.com/incidents/c9px7cw2ghrs) and by 06:30 was [debugging a spike in anonymous traffic](https://www.githubstatus.com/incidents/c9px7cw2ghrs): logged-out users could see errors, and CI jobs and release downloads making anonymous calls could get a 504. It [reported mitigation at 06:51](https://www.githubstatus.com/incidents/c9px7cw2ghrs), then saw [504s climb again at 07:43](https://www.githubstatus.com/incidents/c9px7cw2ghrs). At 08:16 it said ["Logged-in users and authenticated requests are unaffected."](https://www.githubstatus.com/incidents/c9px7cw2ghrs) Errors eased by 08:47, GitHub was [preparing targeted traffic controls](https://www.githubstatus.com/incidents/c9px7cw2ghrs), and the incident closed at 09:09 UTC. The day before, a [separate incident](https://www.githubstatus.com/incidents/pyynvrm1vskw) had degraded API requests, mainly for unauthenticated users, and GitHub said it was working to redeploy a problematic component.

## Why it matters

Anonymous requests are the cheapest traffic class to refuse. Authenticated traffic stayed healthy, so we infer the two classes were already kept apart. That is admission control by class: shed the cheapest class at the edge so overload never reaches authenticated work. Mitigation held under an hour before the 504s returned, and [metastable failures](metastable-failures) offers one reading of that comeback: the trigger was handled, not what sustains the load. The way out is shedding load until the system can carry the rest; we infer the targeted traffic controls do that. Neither page connects the two days. Teach your edge to tell classes apart before a surge forces the choice.

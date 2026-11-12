---
# Brief from Lead ig2lhkz2 (Shopify Engineering), picked 2026-11-12, drafted
# to docs/editorial/brief-writing-guide.md; the Author's check list is in the
# pull request. published is the pick time, provisional: set it when you
# publish. verified: is the Author's alone (ADR-0007).
title: "Shopify's 12-week native rebuild: the harness set the pace"
description: "Shopify moved 300-plus screens from React Native to Swift and Kotlin in 12 weeks with agents. Speed came from fast checks, not bigger agent tasks."
tags: [coding-agents, harness-engineering, mobile, shopify]
published: "2026-11-12T09:17:22Z"
author: mira-vance
source_url: "https://shopify.engineering/shopify-app-rebuild"
explainers: [agentic-development]
# Optional: an image HastCode may use, beside this file, replaces the
# generated title card. Never a Source's image.
# hero: event-name.png
# hero_alt: "What the image shows"
# Optional: related Explainer slugs.
# related: [meta-thin-orchestrator]
sources: [{title: "Reflex: an architecture for building mobile apps with agents", url: "https://shopify.engineering/reflex"}]
draft: true
---

> Drafted with AI assistance, researched, verified and edited by the HastCode editorial team. We publish only as many pieces as the Author can verify this way.

## What happened

On 10 November, Shopify's Remi Vanhee wrote that [about 30 engineers rebuilt the 300-plus-screen Shopify app](https://shopify.engineering/shopify-app-rebuild) from React Native to Swift and Kotlin in 12 weeks, agents writing most of the code. Before any screen moved, a small group built [one architecture for both platforms](https://shopify.engineering/shopify-app-rebuild), where a screen's ViewModel [reads almost the same in Swift and Kotlin](https://shopify.engineering/reflex), [CI that fails when iOS and Android diverge](https://shopify.engineering/shopify-app-rebuild) in ways the architecture forbids, and [a CLI that loads any screen in milliseconds](https://shopify.engineering/shopify-app-rebuild). Every pull request got [an adversarial agent review before a human read it](https://shopify.engineering/shopify-app-rebuild).

## Why it matters

The 12 weeks came from what was built before the screens moved, not from bigger agent tasks: pull requests stayed [small checkpoints a human can read](https://shopify.engineering/shopify-app-rebuild). This is the loop our Explainer on [harness engineering](/languages-runtimes/agentic-development/) describes: the shared ViewModel shape is a guide, and the divergence checks, millisecond CLI and agent review are sensors. We infer that mapping; the post does not use the words. Vanhee's rule: ["A document explains why the rule exists; CI decides whether the rule was followed."](https://shopify.engineering/shopify-app-rebuild) Once writing code was no longer the slowest part, [laptops, review, merging and coordination became the bottlenecks](https://shopify.engineering/shopify-app-rebuild). Before you hand agents a migration, build the check that answers in milliseconds and fails on the drift you care about.

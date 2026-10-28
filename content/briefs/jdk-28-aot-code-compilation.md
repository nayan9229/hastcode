---
# Skeleton from Lead 7zptspgb (Inside Java), picked 2026-10-27.
# Follow docs/editorial/brief-writing-guide.md and replace every
# TODO(author) line. published is the pick time, provisional: set it
# when you publish. verified: is the Author's alone (ADR-0007).
title: "OpenJDK's JEP 544: cached compiled code, with the JIT as fallback"
description: "JEP 544, targeted to JDK 28, caches compiled code from a training run. The lesson: ship precomputed work, keep the slow path as the fallback."
tags: [jvm, aot-compilation, project-leyden, startup]
published: "2026-10-28T12:54:21Z"
verified: 2026-10-28
author: mira-vance
source_url: "https://inside.java/2026/10/25/jep544-target-jdk28/"
explainers: [jvm-warmup-and-aot-compilation]
# Optional: an image HastCode may use, beside this file, replaces the
# generated title card. Never a Source's image.
# hero: event-name.png
# hero_alt: "What the image shows"
# Optional: related Explainer slugs and extra sources (title + url).
# related: [metastable-failures]
sources: [{title: "JEP 544: Ahead-of-Time Code Compilation", url: "https://openjdk.org/jeps/544"}, {title: "JDK 28", url: "https://openjdk.org/projects/jdk/28/"}]
---

> Drafted with AI assistance, researched, verified and edited by the HastCode editorial team. We publish only as many pieces as the Author can verify this way.

## What happened

On 25 October 2026, Inside Java announced that [JEP 544, Ahead-of-Time Code Compilation, is targeted to JDK 28](https://inside.java/2026/10/25/jep544-target-jdk28/), due for General Availability on [9 March 2027, as proposed](https://openjdk.org/projects/jdk/28/). The AOT cache that a training run writes will now hold [native code from the C1 and C2 compilers](https://openjdk.org/jeps/544#Description) beside the classes and profiles it already stored, with no new options. On a two-core machine, the cache cut five framework benchmarks' startup by [around 65% to 80% with that code, and 50% to 70% without](https://openjdk.org/jeps/544#Performance).

## Why it matters

Static compilation ends warmup but optimises for [only one set of hot spots](https://openjdk.org/jeps/544#What-about-static-compilation?); the JIT adapts but [redoes profiling and compiling on every start](https://openjdk.org/jeps/544#Shifting-compilation-work-to-training-runs). JEP 544 ships the precomputed answer and keeps the dynamic path behind it. If production's CPU features or garbage collector differ from training, HotSpot [warns and skips the AOT code](https://openjdk.org/jeps/544#Consistency-of-training-and-production-runs) and still runs correctly, only slower to start and warm up; if the workload shifts, it [recompiles](https://openjdk.org/jeps/544#Shifting-compilation-work-to-training-runs). We infer the general lesson: a cache of expensive work is safe to ship when a miss falls back to the slow path, not to an error. [Our Concept on JVM warmup](jvm-warmup-and-aot-compilation) walks the interpreter-to-C2 ladder this shortens. Make your training run resemble production, and test the cache with `-XX:AOTMode=required`, which [exits when a constraint is violated](https://openjdk.org/jeps/544#Observing-AOT-cache-usage-in-production).

---
# From Lead cklyppp4 (vLLM blog), picked 2026-11-26.
# published is the pick time (make daily may move it to a Gap Day);
# make release keeps it (ADR-0012). verified: is the Author's alone (ADR-0007).
title: "vLLM's ERNIE 5.5 Flash plugin: model code outside the core"
description: "Baidu's ERNIE 5.5 Flash Preview runs on unmodified vLLM 0.32.0 on release day, through a plugin. A new model need not wait for an engine release."
tags: [llm-inference, vllm, mixture-of-experts, expert-parallelism]
published: "2026-11-26T07:37:21Z"
verified: 2026-11-26
author: mira-vance
source_url: "https://vllm.ai/blog/2026-11-24-ernie-5-5-flash-preview"
explainers: [vllm-vs-deepinfer, expert-parallelism]
# Optional: an image HastCode may use, beside this file, replaces the
# generated title card. Never a Source's image.
# hero: event-name.png
# hero_alt: "What the image shows"
# Optional: related Explainer slugs and extra sources (title + url).
# related: [metastable-failures]
# sources: [{title: "The page's title", url: "https://example.com/page"}]
---

## What happened

On 24 November 2026, the Baidu ERNIE and vLLM teams [announced release-day support](https://vllm.ai/blog/2026-11-24-ernie-5-5-flash-preview) for ERNIE 5.5 Flash Preview. It is a mixture-of-experts (MoE) model that selects the [top 10 of 512 routed experts](https://vllm.ai/blog/2026-11-24-ernie-5-5-flash-preview). Support ships as [vLLM-ERNIE](https://github.com/baidu/vLLM-ERNIE), a standalone package installed beside the official `vllm==0.32.0` wheel and found through the [`vllm.general_plugins` entry point](https://docs.vllm.ai/en/latest/design/plugin_system.html). The post says the plugin ["mainly contains the ERNIE model code and a set of kernel optimizations"](https://vllm.ai/blog/2026-11-24-ernie-5-5-flash-preview). The post's fused kernels, each merging several operations into one, run [1.55× to 6.53× faster than naive implementations](https://vllm.ai/blog/2026-11-24-ernie-5-5-flash-preview), measured per kernel.

## Why it matters

[vLLM vs DeepInfer](vllm-vs-deepinfer) weighs a general-purpose engine against a model's own. This release sits between them. ERNIE's model code lives in Baidu's plugin, while the server, KV cache manager and parallelism [mostly come from vLLM itself](https://vllm.ai/blog/2026-11-24-ernie-5-5-flash-preview). The plugin has added features [without waiting for a vLLM release](https://vllm.ai/blog/2026-11-24-ernie-5-5-flash-preview).

We infer that this is how a general-purpose engine keeps pace: per-model speed work moves outside the core. We also infer that the fused kernels are among the plugin's kernel optimizations; the post does not say so. The cost is coupling. The plugin [applies a few hooks to shared vLLM code](https://vllm.ai/blog/2026-11-24-ernie-5-5-flash-preview) and [runs on one vLLM wheel, 0.32.0](https://github.com/baidu/vLLM-ERNIE). So each vLLM upgrade, we infer, needs a retest.

The serve command pairs [data-parallel attention with expert parallelism](https://vllm.ai/blog/2026-11-24-ernie-5-5-flash-preview): each GPU holds whole experts, as [expert parallelism](expert-parallelism) explains. If you maintain a custom architecture, try a plugin before a fork, and benchmark whole requests, not kernels.

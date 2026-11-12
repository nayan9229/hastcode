---
# Skeleton from Lead zwpcrk6k (vLLM blog), picked 2026-11-12; drafted
# 2026-11-12 per docs/editorial/brief-writing-guide.md. published is the
# pick time, provisional: set it when you publish. verified: is the
# Author's alone (ADR-0007).
title: "vLLM's encoder disaggregation: a third GPU pool must earn its keep"
description: "vLLM v0.32.0 puts multimodal encoders on their own GPUs. On GB200, EPD beat PD per GPU at concurrency 64–128 and lost at 16. The split must pay."
tags: [llm-inference, multimodal, gpu, vllm]
published: "2026-11-12T09:17:21Z"
author: mira-vance
source_url: "https://vllm.ai/blog/2026-11-11-encoder-disaggregation-epd"
explainers: [disaggregated-inference]
# Optional: an image HastCode may use, beside this file, replaces the
# generated title card. Never a Source's image.
# hero: event-name.png
# hero_alt: "What the image shows"
# Optional: related Explainer slugs and extra sources (title + url).
# related: [metastable-failures]
sources: [{title: "Disaggregated Encoder", url: "https://docs.vllm.ai/en/latest/features/disagg_encoder/"}]
draft: true
---

> Drafted with AI assistance, researched, verified and edited by the HastCode editorial team. We publish only as many pieces as the Author can verify this way.

## What happened

On 11 November, Inferact, InferMatrix and IBM [described](https://vllm.ai/blog/2026-11-11-encoder-disaggregation-epd) encoder/prefill/decode (EPD) disaggregation, now in [vLLM v0.32.0](https://vllm.ai/blog/2026-11-11-encoder-disaggregation-epd): the multimodal encoder gets its own instances, which [run no language model](https://vllm.ai/blog/2026-11-11-encoder-disaggregation-epd). Embeddings reach prefill over an [EC connector](https://vllm.ai/blog/2026-11-11-encoder-disaggregation-epd), NIXL or Mooncake; a metadata-only request carries the [layout the encoder computed](https://vllm.ai/blog/2026-11-11-encoder-disaggregation-epd), so placeholders are rebuilt downstream. On GB200 with DeepSeek V4.1 Flash, four encoder GPUs beside a four-plus-four PD baseline gave [6.6–9.5% more throughput per GPU at concurrency 64–128](https://vllm.ai/blog/2026-11-11-encoder-disaggregation-epd), all 12 counted, and 44–47% lower P90 time to first token; at [16 it lost per GPU, at 32 roughly tied](https://vllm.ai/blog/2026-11-11-encoder-disaggregation-epd).

## Why it matters

Encoder load moves with media size and modality, [independently of decode](https://vllm.ai/blog/2026-11-11-encoder-disaggregation-epd), so one GPU cannot be sized for both. [Disaggregated inference](disaggregated-inference) argues the same for prefill and decode: separate pools, each tuned to its own target, bought with a copy. EPD adds a third pool and [a second copy](https://vllm.ai/blog/2026-11-11-encoder-disaggregation-epd), and the GB200 run prices them. Vision encoders are [small beside the language model](https://docs.vllm.ai/en/latest/features/disagg_encoder/), so four encoder-only GPUs earn their share of the denominator only once enough media arrives to fill them, which we infer is why EPD lost at concurrency 16. The authors agree: ["The benefit depends on load and resource allocation, rather than separation alone."](https://vllm.ai/blog/2026-11-11-encoder-disaggregation-epd) Before a stage gets its own tier, count every GPU, new ones included, at your real concurrency.

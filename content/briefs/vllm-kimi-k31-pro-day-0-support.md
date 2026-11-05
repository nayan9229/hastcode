---
# From Lead 5tqgafyy (vLLM blog), picked 2026-11-05. published is the
# pick time, provisional: set it when you publish. verified: is the
# Author's alone (ADR-0007).
title: "vLLM's Kimi K3.1 Pro day 0: merge shape and cache layout set the cost"
description: "vLLM's Kimi K3.1 Pro support packs each DCP top-k candidate into 32 bits and the KV cache by block: merge shape and cache layout set the cost."
tags: ["llm-inference", "kv-cache", "context-parallelism", "vllm"]
published: "2026-11-05T17:02:45Z"
verified: 2026-11-05
author: mira-vance
source_url: "https://vllm.ai/blog/2026-11-04-kimi-k31-pro"
explainers: [decode-context-parallelism, disaggregated-inference]
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

On 4 November, Inferact and the vLLM team [announced day-0 support](https://vllm.ai/blog/2026-11-04-kimi-k31-pro) for Moonshot AI's Kimi K3.1 Pro, whose sparse attention [keeps the top 2,048 tokens per query](https://vllm.ai/blog/2026-11-04-kimi-k31-pro#sparse-attention-in-k31-pro). All figures are vendor-reported. Under decode context parallelism, the generic merge sends every GPU's full local top-2,048 as FP32 score-and-id pairs; vLLM's [sends fewer candidates, each id packed with its score into 32 bits, about 10× fewer bytes](https://vllm.ai/blog/2026-11-04-kimi-k31-pro#an-exact-top-k-merge-with-10-less-traffic), and the attention module runs 1.27–1.51× faster at DCP8 on eight GB300s. The KV cache is [packed by block, not by layer](https://vllm.ai/blog/2026-11-04-kimi-k31-pro#a-kv-cache-for-three-kinds-of-state), so, the post says, a prefill-to-decode transfer uses ["one large RDMA read instead of dozens of small per-layer reads"](https://vllm.ai/blog/2026-11-04-kimi-k31-pro#a-kv-cache-for-three-kinds-of-state).

## Why it matters

Once the cache is sharded by token, [every attention layer of every decode step](https://vllm.ai/blog/2026-11-04-kimi-k31-pro#long-context-decode-context-parallelism-with-a-sparse-indexer) must reconcile what the GPUs hold, and the shape of that exchange, we infer, sets more of the cost than the GPU count. The cache's layout sets the other bill: per-layer pages turn one transfer into [dozens of reads](https://vllm.ai/blog/2026-11-04-kimi-k31-pro#a-kv-cache-for-three-kinds-of-state). We infer a rule the post does not state: design the per-layer merge and the cache layout together, before adding GPUs. [Decode context parallelism](decode-context-parallelism) explains why the merge exists; [disaggregated inference](disaggregated-inference) explains why the transfer does. If your cache is sharded, measure bytes per merge and reads per transfer first.

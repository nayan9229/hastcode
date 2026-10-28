---
# Skeleton from Lead xtsbvhnf (vLLM blog), picked 2026-10-27.
# published and reviewed are the pick time, provisional: set them
# when you publish. verified: is the Author's alone (ADR-0007).
title: "Expert parallelism: why each GPU holds only some of the experts"
description: "Expert parallelism serves mixture-of-experts models by spreading whole experts across GPUs and routing tokens to them. How it works, what it costs."
kind: concept
vertical: ai-ml-systems
tags: ["expert-parallelism", "mixture-of-experts", "llm-inference", "gpu", "load-balancing"]
published: "2026-10-27T17:47:17Z"
verified: 2026-10-28
reviewed: "2026-10-27T17:47:17Z"
author: mira-vance
# Optional: an image HastCode may use, beside this file, replaces the
# generated title card. Never a Source's image.
# hero: concept-name.png
# hero_alt: "What the image shows"
# Optional: related Explainer slugs and extra sources (title + url).
# related: [http-connection-reuse]
---

On 19 October 2026, the vLLM team put a number on what it costs to add four GPUs to a running model. Restarting vLLM to scale DeepSeek V4 Pro ["from 16 to 20 GB200 GPUs stopped serving for 197 seconds"](https://vllm.ai/blog/2026-10-19-elastic-ep-near-zero-downtime). The project's [first version of Elastic Expert Parallelism](https://vllm.ai/blog/2026-05-14-elastic-expert-parallelism), from May, avoided the restart, yet still [paused serving for 100 seconds scaling up and 56 seconds scaling back down](https://vllm.ai/blog/2026-10-19-elastic-ep-near-zero-downtime). The new version cut the pause to [0.82 seconds up and 1.96 seconds down](https://vllm.ai/blog/2026-10-19-elastic-ep-near-zero-downtime#scale-up-and-scale-down-timings). Why should four new GPUs disturb the sixteen already serving? Because in this deployment no GPU holds the whole model. Each holds a share of the experts, and every token's route through the model depends on who holds which.

> Drafted with AI assistance, researched, verified and edited by the HastCode editorial team. We publish only as many pieces as the Author can verify this way.

## TL;DR

- **Expert parallelism** (EP) serves a mixture-of-experts model by giving each GPU [whole experts rather than a slice of every expert](https://vllm.ai/blog/2026-05-14-elastic-expert-parallelism#background-expert-parallelism-and-dp-attention), and sending each token only to the GPUs that own the experts it was routed to.
- Each MoE layer becomes route, **dispatch**, compute, **combine**: two [all-to-all exchanges](https://github.com/deepseek-ai/DeepEP) that move tokens to their experts and the outputs back.
- It is how the largest open models are served: DeepSeek-V3 activates [8 of 256 routed experts per layer](https://arxiv.org/html/2412.19437v2#S4.SS2), and its paper describes a decoding unit that spreads the MoE part over [320 GPUs, one expert per GPU](https://arxiv.org/html/2412.19437v2#S3.SS4.SSS2).
- The practices are balancing the router, [replicating hot experts](https://docs.vllm.ai/en/latest/serving/expert_parallel_deployment/#expert-parallel-load-balancer-eplb), hiding the all-to-all behind compute and, now, [resizing the group while it serves](https://vllm.ai/blog/2026-10-19-elastic-ep-near-zero-downtime). Each costs memory, GPU cores or engineering.

## Each GPU owns some of the experts, not a slice of all of them

A **mixture-of-experts** (MoE) model keeps its attention layers dense but replaces most feed-forward layers with sparse **expert layers** that [route each token to a selected set of experts](https://vllm.ai/blog/2026-05-14-elastic-expert-parallelism#background-expert-parallelism-and-dp-attention). A **router**, or gate, scores the experts for each token and keeps the top few. DeepSeek-V3 has [1 shared expert and 256 routed experts per MoE layer and activates 8 routed experts per token](https://arxiv.org/html/2412.19437v2#S4.SS2), so only [37B of its 671B parameters](https://arxiv.org/html/2412.19437v2) work on any one token.

vLLM defines **expert parallelism** this way: ["Instead of sharding each expert across GPUs, experts are distributed across different GPUs, and tokens are dispatched only to the GPUs that own the selected experts."](https://vllm.ai/blog/2026-05-14-elastic-expert-parallelism#background-expert-parallelism-and-dp-attention) The GPUs, or **ranks**, that share the expert layers form the **EP group**; by default each rank holds [the total number of experts divided by the number of ranks](https://docs.vllm.ai/en/latest/serving/expert_parallel_deployment/#expert-distribution-formula).

Google's GShard paper gave the reason in 2020. Experts are ["infeasible to be replicated in all the devices due to its sheer size"](https://arxiv.org/pdf/2006.16668), so ["the only viable strategy is to shard experts into many devices"](https://arxiv.org/pdf/2006.16668). Attention, by contrast, had its weights [replicated to every device](https://arxiv.org/pdf/2006.16668).

What does not count: splitting each expert's matrices across GPUs. In vLLM, without `--enable-expert-parallel`, MoE layers [use tensor parallelism, as a dense model would](https://docs.vllm.ai/en/latest/serving/expert_parallel_deployment/#layer-behavior-with-ep-enabled). Running whole copies of the model side by side, which the Lead calls [scaling by replica](https://vllm.ai/blog/2026-10-19-elastic-ep-near-zero-downtime), is not EP either.

## Route, dispatch, compute, combine: one MoE layer across GPUs

```mermaid
%% title: One MoE layer under expert parallelism: tokens travel to their experts' GPUs and the outputs travel back
flowchart LR
  A["Attention on each GPU, for its own requests"] --> R["Router picks each token's top experts"]
  R --> D["Dispatch: all-to-all sends each token to the GPUs that own its experts"]
  D --> E["Each GPU runs only its local experts on the tokens it received"]
  E --> C["Combine: all-to-all returns each output to the token's home GPU"]
  C --> W["Outputs weighted by the gate values and summed"]
  W --> N["Next layer"]
```

Walk it from the left. In vLLM's layout, [attention runs independently on each data-parallel worker](https://vllm.ai/blog/2026-05-14-elastic-expert-parallelism#background-expert-parallelism-and-dp-attention), each with its own requests and KV cache. Then DeepEP, DeepSeek's communication library, does the two exchanges: dispatch ["route tokens to the corresponding experts across all ranks"](https://github.com/deepseek-ai/DeepEP), and combine ["reduce expert outputs back to their original ranks"](https://github.com/deepseek-ai/DeepEP). The gate value is [multiplied with each expert's output](https://arxiv.org/html/2412.19437v2#S2.SS1.SSS2). GShard's verdict on the whole step: ["The most significant part of gating is communication."](https://arxiv.org/pdf/2006.16668)

DeepSeek-V3's paper puts numbers on each box for decoding. The minimum unit is [40 nodes with 320 GPUs](https://arxiv.org/html/2412.19437v2#S3.SS4.SSS2). Attention runs as four-way tensor parallelism with sequence parallelism, times 80-way data parallelism. The MoE part is EP320: each GPU hosts one expert, and 64 GPUs host redundant and shared experts. Each token selects [9 experts](https://arxiv.org/html/2412.19437v2#S3.SS4.SSS2), the shared one always among them, so, we infer, its dispatch reaches nine GPUs. Dispatch and combine are [direct point-to-point transfers over InfiniBand](https://arxiv.org/html/2412.19437v2#S3.SS4.SSS2). The batch per expert is [usually within 256 tokens](https://arxiv.org/html/2412.19437v2#S3.SS4.SSS2), so the bottleneck is memory access, not arithmetic.

## Where expert parallelism shows up

- **GShard (training).** Its best-quality model had [2,048 experts per MoE layer and 600 billion parameters](https://arxiv.org/pdf/2006.16668), trained on 2,048 TPU v3 cores in 4 days. The authors ["tied the number of devices used for training to the number of experts per MoE layer for simplicity"](https://arxiv.org/pdf/2006.16668), one expert per device.
- **DeepSeek-V3 (serving).** Prefill uses a smaller group: [4 nodes, 32 GPUs, EP32](https://arxiv.org/html/2412.19437v2#S3.SS4.SSS1), each GPU hosting 8 experts plus one redundant copy. Tokens cross nodes over InfiniBand first, then hop over NVLink to the GPU that owns their expert. DeepSeek's February 2025 overview of its production service gives a different decode layout from the paper's: [EP144 across 18 nodes](https://github.com/deepseek-ai/open-infra-index/blob/main/202502OpenSourceWeek/day_6_one_more_thing_deepseekV3R1_inference_system_overview.md#large-scale-cross-node-expert-parallelism-ep), with 2 routed experts and 1 shared expert per GPU.
- **vLLM.** EP size is [tensor-parallel size times data-parallel size](https://docs.vllm.ai/en/latest/serving/expert_parallel_deployment/#configuration). The documented single-node example serves [DeepSeek-V3-0324 with 8-way data and expert parallelism](https://docs.vllm.ai/en/latest/serving/expert_parallel_deployment/#example-command) on one 8-GPU H200 node. The all-to-all is pluggable: a [general-purpose default and DeepEP kernels tuned for prefill or decode](https://docs.vllm.ai/en/latest/serving/expert_parallel_deployment/#backend-selection-guide), among others.
- **Elastic EP**, the variant in the Lead, lets vLLM [change the number of workers at runtime](https://vllm.ai/blog/2026-05-14-elastic-expert-parallelism) with a single API call. It gets its own section below.

## Why sparse models end up spread across many GPUs

Sparsity is efficient per token but hungry for batch. DeepSeek explains that because ["only 8 out of 256 experts per layer are activated"](https://github.com/deepseek-ai/open-infra-index/blob/main/202502OpenSourceWeek/day_6_one_more_thing_deepseekV3R1_inference_system_overview.md#large-scale-cross-node-expert-parallelism-ep), the model needs an extremely large overall batch for each expert to get enough tokens. Spreading experts wide fixes both sides at once, in DeepSeek's account: EP [scales the batch size](https://github.com/deepseek-ai/open-infra-index/blob/main/202502OpenSourceWeek/day_6_one_more_thing_deepseekV3R1_inference_system_overview.md#system-design-principles), which keeps matrix units busy, and each GPU reads only a few experts' weights, which lowers latency. The paper chose EP32 for prefill because it [gives each expert a large enough batch](https://arxiv.org/html/2412.19437v2#S3.SS4.SSS1).

Memory pushes the same way. Data-parallel attention avoids the [duplicated KV cache](https://vllm.ai/blog/2026-05-14-elastic-expert-parallelism#background-expert-parallelism-and-dp-attention) that tensor parallelism would cause in some architectures, and wide EP deployments [maximise KV cache capacity](https://vllm.ai/blog/2026-05-14-elastic-expert-parallelism) for high concurrency or long contexts.

DeepSeek names the price: EP [introduces cross-node communication](https://github.com/deepseek-ai/open-infra-index/blob/main/202502OpenSourceWeek/day_6_one_more_thing_deepseekV3R1_inference_system_overview.md#system-design-principles) that must be overlapped with compute, and load balancing between data-parallel instances. In training, cross-node EP left DeepSeek-V3 with a computation-to-communication ratio of [roughly 1:1](https://arxiv.org/html/2412.19437v2#S3.SS2.SSS1) until the team redesigned its pipeline. The practices below each pay part of that price.

## Practice one: train the router to spread the load

A router left alone picks favourites. GShard found that without balancing, a few busy experts get most tokens, ["leaving other experts untrained, slowing down the training"](https://arxiv.org/pdf/2006.16668). Its answer had two parts. An **auxiliary loss** added to the training objective rewards even routing. An **expert capacity** caps each expert's tokens per batch, and a token whose chosen experts are both full becomes ["an overflowed token"](https://arxiv.org/pdf/2006.16668): it skips the MoE layer and passes on through the residual connection.

Under expert parallelism, both matter for speed as well as quality. DeepSeek-V3's paper says an unbalanced load will [diminish computational efficiency in scenarios with expert parallelism](https://arxiv.org/html/2412.19437v2#S2.SS1.SSS2), but that too large an auxiliary loss hurts model quality. So DeepSeek adds a per-expert **bias** to the routing scores only: after each step it lowers the bias of overloaded experts and raises it for underloaded ones, while the gate value still comes from the original score. With that, DeepSeek-V3 [drops no tokens in training or in inference](https://arxiv.org/html/2412.19437v2#S2.SS1.SSS2).

The cost lands on the model, not the cluster. Capacity limits drop work on hot batches, an auxiliary loss trades some quality for balance, and DeepSeek still keeps a [small sequence-wise balance loss](https://arxiv.org/html/2412.19437v2#S2.SS1.SSS2) to stop extreme imbalance inside any single sequence. A balanced training run does not guarantee balanced traffic in production.

## Practice two: replicate the hot experts at serving time

Production traffic is skewed anyway. vLLM's docs note that although MoE models are trained for an even spread, [in practice the distribution can be highly skewed](https://docs.vllm.ai/en/latest/serving/expert_parallel_deployment/#expert-parallel-load-balancer-eplb), and DeepSeek reports [inherently high-load experts](https://github.com/deepseek-ai/open-infra-index/blob/main/202502OpenSourceWeek/day_6_one_more_thing_deepseekV3R1_inference_system_overview.md#achieving-optimal-load-balancing). One overloaded GPU ["becomes a performance bottleneck, slowing the entire system while leaving other GPUs idle"](https://github.com/deepseek-ai/open-infra-index/blob/main/202502OpenSourceWeek/day_6_one_more_thing_deepseekV3R1_inference_system_overview.md#achieving-optimal-load-balancing).

The fix is **redundant experts**: copy the busiest experts onto extra slots. DeepSeek [detects high-load experts from online statistics](https://arxiv.org/html/2412.19437v2#S3.SS4.SSS1), adjusts the set periodically ("e.g., every 10 minutes") and rearranges experts within a node so cross-node traffic does not grow. It open-sourced the planner as [EPLB](https://github.com/deepseek-ai/EPLB), under the MIT licence, with a [hierarchical policy for prefill](https://github.com/deepseek-ai/EPLB#hierarchical-load-balancing) and a global one for decode's larger groups.

vLLM's EPLB [collects load statistics on every forward pass](https://docs.vllm.ai/en/latest/serving/expert_parallel_deployment/#eplb-parameters), looks back over 1,000 engine steps and rebalances every 3,000 by default. The cost is GPU memory that would otherwise hold KV cache: about [2.4 GB for DeepSeek-V3 per redundant expert per rank](https://docs.vllm.ai/en/latest/serving/expert_parallel_deployment/#memory-footprint-overhead), and vLLM [recommends 32 redundant experts](https://docs.vllm.ai/en/latest/serving/expert_parallel_deployment/#example-command_1) at large scale.

## Practice three: hide the all-to-all behind compute

If dispatch and combine run while the GPU waits, the network sets the speed. The common answer is to split each batch into two **micro-batches** and let one compute while the other communicates. In DeepSeek-V3's prefill, the attention and MoE of one micro-batch [overlap the dispatch and combine of the other](https://arxiv.org/html/2412.19437v2#S3.SS4.SSS1); in decode, where attention takes more of the time, attention overlaps [dispatch, MoE and combine](https://arxiv.org/html/2412.19437v2#S3.SS4.SSS2). The production service splits decode attention further into a [5-stage pipeline](https://github.com/deepseek-ai/open-infra-index/blob/main/202502OpenSourceWeek/day_6_one_more_thing_deepseekV3R1_inference_system_overview.md#computation-communication-overlapping). vLLM offers the same idea as [`--enable-dbo`](https://docs.vllm.ai/en/latest/serving/expert_parallel_deployment/#performance-optimization), dual batch overlap.

The router can help too. In training, DeepSeek-V3 sends each token to [at most 4 nodes](https://arxiv.org/html/2412.19437v2#S3.SS2.SSS2), because NVLink inside a node offers 160 GB/s, about 3.2 times InfiniBand's 50 GB/s. At the same communication cost, a token could reach up to 13 experts.

The costs are real. Communication kernels run on the GPU's own cores: DeepSeek gave [20 of the H800's 132 streaming multiprocessors](https://arxiv.org/html/2412.19437v2#S3.SS5.SSS1) to communication and notes that this limits compute throughput, and DeepEP states that ["EP dispatch and combine require GPU SMs"](https://github.com/deepseek-ai/DeepEP). Kernels tuned for one phase can hurt another: vLLM warns DeepEP's two kernel families [may perform poorly on mixed workloads](https://docs.vllm.ai/en/latest/serving/expert_parallel_deployment/#performance-optimization).

## Practice four: resize the expert group while it serves

Before Elastic EP, EP in vLLM was static: once a deployment started, [its serving capacity was fixed](https://vllm.ai/blog/2026-05-14-elastic-expert-parallelism), and the only option was a restart that was slow and could drop traffic. Resizing is hard because a change in EP size [invalidates several pieces of state](https://vllm.ai/blog/2026-05-14-elastic-expert-parallelism#the-challenge-what-state-needs-to-change): communication groups with fixed rank sets, the expert-to-rank mapping, model weights, and compiled CUDA graphs.

Elastic EP treats resizing as a sequence. To scale up, new workers start with [placeholder weights](https://vllm.ai/blog/2026-05-14-elastic-expert-parallelism#2-new-engine-core-initialization), receive the non-expert weights [GPU to GPU from existing ranks](https://vllm.ai/blog/2026-05-14-elastic-expert-parallelism#4-expert-mapping-and-weight-transfer), and join the active group. At first they [run attention but own no experts](https://vllm.ai/blog/2026-05-14-elastic-expert-parallelism#5-the-switch); an EPLB reshuffle then [redistributes experts across all ranks](https://vllm.ai/blog/2026-05-14-elastic-expert-parallelism#6-eplb-reshuffle). To scale down, the reshuffle comes first, so experts [move off the departing ranks](https://vllm.ai/blog/2026-05-14-elastic-expert-parallelism#scale-down-flow) before they leave.

We infer that this is why EP makes elasticity harder than replica scaling. A new replica holds every expert and can serve alone, while a new EP rank changes which GPU every token visits, so every rank has to agree on the new layout at once.

## What near-zero downtime took

The October work moved almost everything out of the pause. vLLM now runs [background preparation and a short commit](https://vllm.ai/blog/2026-10-19-elastic-ep-near-zero-downtime#overlapping-topology-preparation-with-serving): joining workers start, receive weights over NVLink or RDMA and warm up while existing workers keep serving, then vLLM "briefly pauses its schedulers, activates the prepared topology, and resumes serving."

The hard part was the existing workers' CUDA graphs. Switching to a new EP communicator [changes the buffers that the dispatch and combine kernels use](https://vllm.ai/blog/2026-10-19-elastic-ep-near-zero-downtime#expand-ep-communication-in-place), which invalidates the captured graphs. [NIXL EP](https://github.com/ai-dynamo/nixl/blob/main/examples/device/ep/README.md) instead adds or removes ranks in place, keeping the same buffers and connections, and vLLM keeps [tensor addresses and layouts stable](https://vllm.ai/blog/2026-10-19-elastic-ep-near-zero-downtime#preserve-cuda-graphs-on-existing-workers) so the old graphs stay valid. Joining workers [mask every remote peer](https://vllm.ai/blog/2026-10-19-elastic-ep-near-zero-downtime#prepare-joining-workers-independently) to warm up alone.

On DeepSeek V4 Pro, scaling from 16 to 20 GPUs paused serving for [99.62 s when blocking, 58.04 s with background preparation and 0.82 s with graph reuse too](https://vllm.ai/blog/2026-10-19-elastic-ep-near-zero-downtime#scale-up-and-scale-down-timings). Under load rising from about 8 to 16 requests per second, Elastic EP [completed all 11,391 requests](https://vllm.ai/blog/2026-10-19-elastic-ep-near-zero-downtime#maintaining-the-slo-while-adding-gpus) and kept rolling p95 time to first token under 5 seconds.

The cost: the full scale-up still took [92.34 s](https://vllm.ai/blog/2026-10-19-elastic-ep-near-zero-downtime#scale-up-and-scale-down-timings), so new capacity arrives slowly, and graph reuse is supported only with NIXL EP.

## What expert parallelism is not

The nearest neighbour is **tensor parallelism** applied to the MoE layers. There, every expert is cut into pieces and the same GPUs share each one; vLLM does this when EP is off, forming a [tensor-parallel group across all the GPUs](https://docs.vllm.ai/en/latest/serving/expert_parallel_deployment/#layer-behavior-with-ep-enabled), "similar to dense models". With EP, each expert lives whole on one rank and tokens move instead of weights. The two combine: vLLM's EP size is [tensor-parallel size times data-parallel size](https://docs.vllm.ai/en/latest/serving/expert_parallel_deployment/#configuration), and DeepSeek uses tensor parallelism for attention [alongside EP for experts](https://arxiv.org/html/2412.19437v2#S3.SS4.SSS1).

The second confusion is between balancing experts and balancing requests. DeepSeek runs [three load balancers](https://github.com/deepseek-ai/open-infra-index/blob/main/202502OpenSourceWeek/day_6_one_more_thing_deepseekV3R1_inference_system_overview.md#achieving-optimal-load-balancing). Two spread requests, sequence lengths and KV cache across data-parallel instances. Only the third is about experts: it aims to minimise the maximum dispatch receive load on any GPU. A router can send most tokens to one GPU even when every instance has the same number of requests.

Prefill-decode disaggregation, which DeepSeek [also uses](https://github.com/deepseek-ai/open-infra-index/blob/main/202502OpenSourceWeek/day_6_one_more_thing_deepseekV3R1_inference_system_overview.md#large-scale-cross-node-expert-parallelism-ep), is a separate idea that runs the two phases on different GPUs. It deserves its own Explainer.

## What understanding expert parallelism makes possible

Once you see the dispatch and combine, a slow MoE deployment asks sharper questions. Is one rank receiving far more tokens than the rest? vLLM will log that as [balancedness, average tokens per expert divided by the maximum](https://docs.vllm.ai/en/latest/serving/expert_parallel_deployment/#eplb-parameters). Is the all-to-all hidden behind compute, or is the GPU waiting on the network? Is a redundant expert worth 2.4 GB of KV cache on your traffic?

It also changes what capacity planning can be. DeepSeek already runs its service on [all nodes during the day and fewer at night](https://github.com/deepseek-ai/open-infra-index/blob/main/202502OpenSourceWeek/day_6_one_more_thing_deepseekV3R1_inference_system_overview.md#statistics-of-deepseeks-online-service), moving the rest to research and training: over one 24-hour window, a peak of 278 nodes and an average of 226.75. We infer that a group which can grow and shrink while serving lets more operators follow demand like this, in steps smaller than a whole replica.

You can try it at small scale. vLLM's [Elastic EP test](https://github.com/vllm-project/vllm/blob/v0.30.0/tests/distributed/test_elastic_ep.py#L224) scales a deployment from two GPUs to four and back while sending requests, and the team plans to [connect Elastic EP to orchestrator autoscalers](https://vllm.ai/blog/2026-10-19-elastic-ep-near-zero-downtime#connect-elastic-ep-to-dynamo-and-llm-d-orchestration). The expert group, once fixed at launch, is becoming something you can resize like any other pool.

## Sources

- Itay Alroy (NVIDIA), ["Achieving near-zero-downtime scaling in Elastic Expert Parallelism"](https://vllm.ai/blog/2026-10-19-elastic-ep-near-zero-downtime), vLLM blog, 19 October 2026.
- Itay Alroy, Yongji Wu, Rui Qiao, Tyler Michael Smith, Moein Khazraee, Omri Kahalon, Tzu-Ling Kan, Ron Tourgeman, ["Elastic Expert Parallelism in vLLM"](https://vllm.ai/blog/2026-05-14-elastic-expert-parallelism), vLLM blog, 14 May 2026.
- vLLM, ["Expert Parallel Deployment"](https://docs.vllm.ai/en/latest/serving/expert_parallel_deployment/), vLLM documentation.
- vLLM, [Elastic EP test, `cuda_graphs_heavy_nixl_ep` case](https://github.com/vllm-project/vllm/blob/v0.30.0/tests/distributed/test_elastic_ep.py#L224), v0.30.0.
- NIXL, [NIXL EP README](https://github.com/ai-dynamo/nixl/blob/main/examples/device/ep/README.md).
- DeepSeek-AI, ["DeepSeek-V3 Technical Report"](https://arxiv.org/html/2412.19437v2), arXiv:2412.19437.
- DeepSeek-AI, ["Day 6: One More Thing, DeepSeek-V3/R1 Inference System Overview"](https://github.com/deepseek-ai/open-infra-index/blob/main/202502OpenSourceWeek/day_6_one_more_thing_deepseekV3R1_inference_system_overview.md), February 2025.
- DeepSeek-AI, [DeepEP](https://github.com/deepseek-ai/DeepEP), MIT License.
- DeepSeek-AI, [EPLB: Expert Parallelism Load Balancer](https://github.com/deepseek-ai/EPLB), MIT License.
- Dmitry Lepikhin, HyoukJoong Lee, Yuanzhong Xu, Dehao Chen, Orhan Firat, Yanping Huang, Maxim Krikun, Noam Shazeer, Zhifeng Chen, ["GShard: Scaling Giant Models with Conditional Computation and Automatic Sharding"](https://arxiv.org/pdf/2006.16668), arXiv:2006.16668, 2020 ([abstract](https://arxiv.org/abs/2006.16668)).

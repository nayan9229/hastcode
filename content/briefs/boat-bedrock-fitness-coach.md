---
# From Lead waovhlae (AWS Architecture Blog), picked 2026-10-29.
# verified: is the Author's alone (ADR-0007).
title: "boAt's AI fitness coach: the prompt is not the guardrail"
description: "boAt runs Llama 3.3 on Bedrock behind Lambda. Prompts alone did not stop bad output, so code checks every reply. Guardrails belong in code."
tags: [llm, guardrails, aws, reliability]
published: "2026-10-29T19:00:46Z"
verified: 2026-10-29
author: mira-vance
source_url: "https://aws.amazon.com/blogs/architecture/how-boat-delivers-personalized-fitness-coaching-with-amazon-bedrock/"
explainers: [meta-thin-orchestrator]
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

On 29 October, AWS's Architecture Blog [described boAt's AI Coach](https://aws.amazon.com/blogs/architecture/how-boat-delivers-personalized-fitness-coaching-with-amazon-bedrock/), which writes fitness plans, and nudges sent through the day, from boAt smartwatch data. It runs Meta's Llama 3.3 70B Instruct on [Amazon Bedrock](https://aws.amazon.com/blogs/architecture/how-boat-delivers-personalized-fitness-coaching-with-amazon-bedrock/), called from AWS Lambda. Lambda renames terse database keys for the model (`weight` becomes `weight_in_kg`) and parses each reply for the app. The team says the feature has [served millions of users for over six months](https://aws.amazon.com/blogs/architecture/how-boat-delivers-personalized-fitness-coaching-with-amazon-bedrock/).

## Why it matters

The model is the one component whose output cannot be trusted to keep its shape. boAt's engineers found that ["prompt optimization alone did not eliminate all errors"](https://aws.amazon.com/blogs/architecture/how-boat-delivers-personalized-fitness-coaching-with-amazon-bedrock/). So they added a parsing layer with safe defaults in code, and they list the activities the watch supports in the prompt. The Lambda layer also keeps the model's quirks away from the app: the post says [future input changes need no backend change](https://aws.amazon.com/blogs/architecture/how-boat-delivers-personalized-fitness-coaching-with-amazon-bedrock/).

We infer a general rule from this. Put the contract between the model and the rest of the system in deterministic code, and treat the prompt as a hint, not a guarantee. Meta reached the same split at larger scale in its [thin orchestrator](/ai-ml-systems/meta-thin-orchestrator/), where a deterministic router of about 1,300 lines steers 14 specialist agents. If an LLM feeds your app, validate every field it returns, and choose the default before you ship.

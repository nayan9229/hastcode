# Release plan, from 2026-10-30 (ADR-0010: two Pieces a day)

Live on 2026-10-30: `cloudflare-analytics-pipeline`, `zgc-adaptive-heap-sizing`.

The eight below are drafts again. On each release day, for its two Pieces: remove `draft: true`, add `verified: <that day>` under `published`, set `published` (and `reviewed`) to the release time, then Publish. Any fresh Brief takes a slot first, because it expires 72 hours after its Lead; push the plan back a day when one does.

| Day | Pieces |
|---|---|
| 2026-10-30 | `cloudflare-analytics-pipeline`, `zgc-adaptive-heap-sizing` (live) |
| 2026-10-31 | `dnssec-root-key-rollover` (00:09 IST), Brief `emudhra-weak-ca-keys` (22:05 IST; took Meta PrivacyLib's slot) |
| 2026-11-01 | nothing released |
| 2026-11-02 | `meta-privacylib`, `snapshot-restore-hydration` (18:13 UTC) |
| 2026-11-03 | nothing released |
| 2026-11-04 | nothing released |
| 2026-11-05 | Briefs `fly-io-petsem-sqlite-lock-contention`, `vllm-kimi-k31-pro-day-0-support` (17:02 UTC; both expire 2026-11-07 00:00 UTC) |
| 2026-11-06 | nothing released |
| 2026-11-07 | Brief `github-ref-cas-git-storage` (expires 2026-11-08 20:23 UTC), `workers-ai-kv-cache-exposure` (06:26 UTC) |
| 2026-11-08 | nothing released; Brief `fly-io-ipv4-checksum-offload-cve` expired 2026-11-09 (draft kept, Author's choice 2026-11-12) |
| 2026-11-09 | nothing released |
| 2026-11-10 | Briefs `github-anonymous-traffic-surge`, `topstep-redis-to-valkey-migration` (06:56 UTC) |
| 2026-11-11 | `cloudflare-origin-range-requests`, `shopify-overexposed` (released 2026-11-12 with `published` backdated to 11 Nov: the Author chose to clear the backlog in one day; `verified` 2026-11-12) |
| 2026-11-12 | `volatile-expressions-sharded-sql`, `bundled-vs-platform-trust-stores` |

Backlog cleared 2026-11-12; `make daily` drafts Explainers again from 2026-11-13. Intake queue cut the same day to Leads under 3 days old (84 dismissed).

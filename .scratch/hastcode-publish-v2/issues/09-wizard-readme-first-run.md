# 09: Wizard, README and first-run checklist

**What to build:** The Author can set up and run the whole publish path without this conversation. `scripts/setup-deploy.sh` is rewritten with stages for: tools and logins; GitHub repository; Cloudflare account, Workers project, R2 bucket with `images.hastcode.com` custom domain, two API tokens (Workers edit; R2 object read/write scoped to the bucket) and Worker secrets; Buttondown keys; Vercel Hobby project, `qa.hastcode.com`, Vercel Authentication; GitHub secrets and the IndexNow gate variable; first `publish vl` with the QA checklist; first branch Preview with its checklist; the Production DNS flip (apex and `www` to the Worker) behind the ADR-0001 content gate; first `publish cf` with the live checklist drawn from the spec's "Further Notes". Vercel Pro, Vercel domain and bypass stages are removed. README is updated end to end (publish, QA, Preview, Image Store modes, prune, secrets, Worker). Call the Skill tool with `wizard` for the script; the stages above are pre-approved.

Spec: `../spec.md` ("Wizard and docs", "Further Notes"). ADR-0004.

**Blocked by:** 04 (Prune command), 08 (Workflows and the publish command)

**Status:** ready-for-agent

- [ ] Wizard runs through every stage in dry conversation (no real accounts touched here), writes `.env` and GitHub secrets via the library helpers, and `confirm`s before the DNS flip and before `publish cf`
- [ ] Each checklist item maps to a spec story or a research open question
- [ ] README has no remaining reference to Vercel Pro, Vercel production, push-to-main deploys or the Vercel preview workflow
- [ ] First deploy criteria that only a human can complete are left unticked and the ticket ends `ready-for-human` with a "Remaining human steps" list

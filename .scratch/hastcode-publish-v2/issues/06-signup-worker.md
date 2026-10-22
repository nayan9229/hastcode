# 06: Signup Worker for the Production Target

**What to build:** On Production, a Reader submits the newsletter form and it behaves exactly as on QA. A Cloudflare Worker module at `worker/index.js` exports `fetch`: for `POST` on the signup path it validates email syntax, drops honeypot hits silently, applies a best-effort per-IP rate limit (in-isolate memory, documented as such), calls Buttondown's subscriber API with the page `source` as a tag, and answers 303 to the confirmation page for form posts or JSON for `Accept: application/json`; wrong methods get 405; upstream failures get a friendly 502-style error; the Buttondown key never appears in any response. Every other request is forwarded to the assets binding. The key and base URL come from Worker secrets/vars; the base URL is overridable so a stub Buttondown can be used.

Tests stay Go-only (Author's decision): no Vitest, no JavaScript test toolchain. Instead ship `scripts/check-worker.sh`, which starts a stub Buttondown, runs the Worker with `wrangler dev` against the built output, and curls every case the Go subscribe tests cover, printing pass/fail per case. Run it by hand here and record the results in this ticket; ticket 09 puts the same cases on the first-Preview checklist.

Spec: `../spec.md` ("Worker", "Testing Decisions"). Mirror the Go handler's contract exactly; read its tests first.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] `scripts/check-worker.sh` run locally passes every case: form 303 with correct `Location`; JSON 200 `{"ok":true}`; bad email 400; honeypot filled → success with no upstream call; rate limit trips after N from one IP; upstream 5xx → error without the key; GET → 405; a non-signup path → served from assets
- [ ] The stub Buttondown logs the email, the `Authorization: Token` header and the source tag, matching the Go tests' expectations
- [ ] `make test` unchanged (Go only, still green); `wrangler` pinned as a devDependency for the check script and deploys
- [ ] README "Subscribe function" documents both implementations, the Worker's secrets, and how to run the check script

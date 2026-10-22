# 06: Signup Worker for the Production Target

**What to build:** On Production, a Reader submits the newsletter form and it behaves exactly as on QA. A Cloudflare Worker module at `worker/index.js` exports `fetch`: for `POST` on the signup path it validates email syntax, drops honeypot hits silently, applies a best-effort per-IP rate limit (in-isolate memory, documented as such), calls Buttondown's subscriber API with the page `source` as a tag, and answers 303 to the confirmation page for form posts or JSON for `Accept: application/json`; wrong methods get 405; upstream failures get a friendly 502-style error; the Buttondown key never appears in any response. Every other request is forwarded to the assets binding. The key and base URL come from Worker secrets/vars; the base URL is injectable so tests use a stub. Tests run with Vitest and Cloudflare's Workers test pool, with a stubbed global `fetch` as Buttondown, one case per existing Go case. `make test` runs them after the Go tests; CI already has Node.

Spec: `../spec.md` ("Worker", "Testing Decisions"). Mirror the Go handler's contract exactly; read its tests first.

**Blocked by:** None (can start immediately)

**Status:** ready-for-agent

- [ ] Worker-seam tests: form 303 with correct `Location`; JSON 200 `{"ok":true}`; bad email 400; honeypot filled → success with no upstream call; rate limit trips after N from one IP; upstream 5xx → error response without the key; GET → 405; any non-signup path → forwarded to assets
- [ ] The stubbed Buttondown receives the email, the `Authorization: Token` header and the source tag, matching the Go tests' expectations
- [ ] `make test` runs Go and Worker tests and fails on either; `npm ci` installs the test pool pinned in `package.json`
- [ ] README "Subscribe function" documents both implementations and the Worker's secrets

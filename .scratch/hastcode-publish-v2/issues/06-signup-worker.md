# 06: Signup Worker for the Production Target

**What to build:** On Production, a Reader submits the newsletter form and it behaves exactly as on QA. A Cloudflare Worker module at `worker/index.js` exports `fetch`: for `POST` on the signup path it validates email syntax, drops honeypot hits silently, applies a best-effort per-IP rate limit (in-isolate memory, documented as such), calls Buttondown's subscriber API with the page `source` as a tag, and answers 303 to the confirmation page for form posts or JSON for `Accept: application/json`; wrong methods get 405; upstream failures get a friendly 502-style error; the Buttondown key never appears in any response. Every other request is forwarded to the assets binding. The key and base URL come from Worker secrets/vars; the base URL is overridable so a stub Buttondown can be used.

Tests stay Go-only (Author's decision): no Vitest, no JavaScript test toolchain. Instead ship `scripts/check-worker.sh`, which starts a stub Buttondown, runs the Worker with `wrangler dev` against the built output, and curls every case the Go subscribe tests cover, printing pass/fail per case. Run it by hand here and record the results in this ticket; ticket 09 puts the same cases on the first-Preview checklist.

Spec: `../spec.md` ("Worker", "Testing Decisions"). Mirror the Go handler's contract exactly; read its tests first.

**Blocked by:** None (can start immediately)

**Status:** resolved

- [x] `scripts/check-worker.sh` run locally passes every case: form 303 with correct `Location`; JSON 200 `{"ok":true}`; bad email 400; honeypot filled → success with no upstream call; rate limit trips after N from one IP; upstream 5xx → error without the key; GET → 405; a non-signup path → served from assets
- [x] The stub Buttondown logs the email, the `Authorization: Token` header and the source tag, matching the Go tests' expectations
- [x] `make test` unchanged (Go only, still green); `wrangler` pinned as a devDependency for the check script and deploys
- [x] README "Subscribe function" documents both implementations, the Worker's secrets, and how to run the check script

## Result

`worker/index.js` (port of `subscribe/`), `scripts/check-worker.sh`, `wrangler` 4.149.0 pinned in `package.json` devDependencies, README "Subscribe function" split into Go function / Worker sections. Worker secret `BUTTONDOWN_API_KEY` (required), var `BUTTONDOWN_BASE_URL` (optional). Client IP is `CF-Connecting-IP`. Rate limit 5 per IP per 10 min in isolate memory (best-effort). Known divergence: Go's `net/mail` accepts quoted local parts (`"a b"@x.com`); the Worker rejects them.

`scripts/check-worker.sh` run 2026-10-22 (wrangler 4.149.0, Node 24.10.0, macOS, assets = `make build MODE=preview FLAGS=-skip-search-index` output):

```
PASS  form post -> 303 Location /subscribed/, upstream gets email, Token key, source tag
PASS  scripted post -> 200 {"ok":true}
PASS  invalid email -> 400 (JSON and form), no upstream call
PASS  honeypot filled -> success, no upstream call
PASS  rate limit trips after 5 from one IP, other IP unaffected
PASS  upstream 5xx -> 502 friendly message, key never in response
PASS  upstream unreachable -> 502
PASS  already subscribed -> same success
PASS  Buttondown rejects email -> 400
PASS  GET -> 405 Allow: POST
PASS  missing BUTTONDOWN_API_KEY -> 502, no upstream call
PASS  non-signup paths served from assets (200, 307, 404 page)
PASS  key never appears in any response (all requests above)
stub Buttondown received (form-post case): POST /v1/subscribers auth=Token check-worker-key-120512377-0123456789 email=reader@example.com tags=["/distributed-systems/raft-leader-election/"] ip=203.0.113.1
13 passed, 0 failed
```

Mutation check: with 303 changed to 302, the limit raised to 50, the `Token` header altered and email rejection mapped to 502, the script reported 5 FAIL cases naming each break. `make test` green, Go only.

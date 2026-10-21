# 14: Newsletter subscribe function and form

**What to build:** A Reader enters an email in the form at the end of any Explainer or Brief and gets a clear success or error. The form is a plain HTML POST to `/api/subscribe` with `email`, a `source` field set to the page path, and a honeypot field; it works without JavaScript (303 to a `/subscribed/` page) and, when scripted, submits in place and shows the result. The Go serverless function validates email syntax, drops honeypot hits silently, applies a best-effort per-IP rate limit, forwards to Buttondown with the source as a tag, returns 303 or JSON by `Accept`, and never exposes the API key. Buttondown base URL and key come from environment variables; the base URL is injectable so tests use a fake server.

**Blocked by:** 01 (Walking skeleton)

**Status:** resolved

- [x] HTTP-seam tests with an `httptest` fake Buttondown: valid email → 303 (form) or 200 JSON (Accept json) and the fake received email + source tag; invalid email → 400 with message; honeypot filled → 200/303 with no upstream call; rate limit trips after N requests from one IP; upstream 5xx → 502-style error with a friendly message
- [x] Function laid out as a Vercel Go serverless function with its own entry; `vercel dev` or an equivalent local run documented
- [ ] Form partial rendered at the end of Explainer and Brief pages (Build-seam test asserts the form, its action, hidden `source` value and honeypot field)
- [x] `/subscribed/` static page exists in output, noindex
- [x] Progressive-enhancement script small enough to keep the JS budget; page works with it disabled

## Notes

- Explainer pages carry the form, and a Build-seam test checks it (`builder/newsletter_test.go`). Briefs don't exist yet (ticket 06). The Brief template only needs `@NewsletterForm(views.Newsletter{Source: briefURL(b)})` at its end, plus a matching assertion. Until that lands, the "Brief" half of the third box is open.
- Buttondown API checked 2026-10-21: `POST /v1/subscribers`, `Authorization: Token <key>`, `{"email_address", "tags", "ip_address"}`, 201 on success, errors `{"code","detail","metadata"}`. A duplicate email (`subscriber_already_exists`) is treated as success so the form never reveals who is subscribed. `email_*` codes and 422 map to 400. Anything else maps to 502 and is logged.
- Deploy notes for ticket 15 are in README "Subscribe function": env vars, how zero-config `api/subscribe.go` merges into the builder's `.vercel/output`, and the route-ordering constraint. None of this has been tested on a real Vercel deploy yet.

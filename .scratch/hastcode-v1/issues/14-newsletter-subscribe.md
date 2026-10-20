# 14: Newsletter subscribe function and form

**What to build:** A Reader enters an email in the form at the end of any Explainer or Brief and gets a clear success or error. The form is a plain HTML POST to `/api/subscribe` with `email`, a `source` field set to the page path, and a honeypot field; it works without JavaScript (303 to a `/subscribed/` page) and, when scripted, submits in place and shows the result. The Go serverless function validates email syntax, drops honeypot hits silently, applies a best-effort per-IP rate limit, forwards to Buttondown with the source as a tag, returns 303 or JSON by `Accept`, and never exposes the API key. Buttondown base URL and key come from environment variables; the base URL is injectable so tests use a fake server.

**Blocked by:** 01 (Walking skeleton)

**Status:** ready-for-agent

- [ ] HTTP-seam tests with an `httptest` fake Buttondown: valid email → 303 (form) or 200 JSON (Accept json) and the fake received email + source tag; invalid email → 400 with message; honeypot filled → 200/303 with no upstream call; rate limit trips after N requests from one IP; upstream 5xx → 502-style error with a friendly message
- [ ] Function laid out as a Vercel Go serverless function with its own entry; `vercel dev` or an equivalent local run documented
- [ ] Form partial rendered at the end of Explainer and Brief pages (Build-seam test asserts the form, its action, hidden `source` value and honeypot field)
- [ ] `/subscribed/` static page exists in output, noindex
- [ ] Progressive-enhancement script small enough to keep the JS budget; page works with it disabled

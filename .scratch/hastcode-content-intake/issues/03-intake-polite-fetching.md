# 03: Intake: polite fetching and failure reporting

**What to build:** Intake is a good citizen and survives broken Sources. Each Source's ETag and Last-Modified are stored and sent back as conditional request headers (RFC 9110); a 304 costs nothing and lists nothing. Requests go one at a time with a fixed User-Agent naming HastCode and the contact address, and a per-request timeout. A Source that errors (network, 4xx/5xx, unparseable) is reported by name in the run summary while every other Source still runs. A Source whose newest item is more than 90 days old is flagged as silent.

Spec: `../spec.md` ("Intake"). Research section 5 (validators verified on Kubernetes, Cloudflare; Go blog sends none).

**Blocked by:** 02 (Intake tracer)

**Status:** resolved

- [x] Intake-seam tests: second request carries the stored validators; a 304 is handled; a Source without validators is fetched in full; User-Agent asserted; a failing Source reported while others list Leads; a silent Source flagged; a hanging Source times out
- [x] Run exits non-zero only on Watchlist errors, not on individual Source failures

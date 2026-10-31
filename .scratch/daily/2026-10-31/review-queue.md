# Daily review queue, 2026-10-31

Drafted by `make daily` (ADR-0009). Every Piece is `draft: true` and unverified. Verify every claim against its Source, then set `verified:` and remove `draft: true` yourself (ADR-0007).

## Run 20261031-123258

### eMudhra's weak CA keys: trust stores that cannot drop a broken CA

`content/briefs/emudhra-weak-ca-keys.md` · brief · 238 words · fact-checker: 18 claims checked, 6 fixed, 0 removed

**Publish within 72 hours of the Lead's date**, with or after its Explainer: dnssec-root-key-rollover

**Drafter asks you to check:**
- [ ] Date 19 October 2026: Source gives no time or time zone, so none is stated
- [ ] 'eight weak keys for its intermediate certificates': Source says 'eight shockingly weak keys' and separately 'affected intermediate private keys'
- [ ] 'one was the number 4': Source 'Not even four digits, just the number 4'
- [ ] Pingora origin-facing proxies patched within hours: Source 'These were patched in hours to reject the vulnerable eMudhra keys'
- [ ] '50 million over a week': Source 'With 50 million probes over the last week'
- [ ] 'more than half on large networks': Source 'Across the large networks, more than half of probes still accept the compromised chain'
- [ ] '30–40% of Android WebViews': Source says embedded browsers, like Android WebView, 30–40% still accepting
- [ ] Quote: 'a trust store can't be a static list of root certificate authorities': check exact wording and apostrophe
- [ ] Description's 'a week on': inferred, since the post is dated 30 Oct, 11 days after 19 Oct
- [ ] Inference: DNSSEC stale-anchor analogy is labelled 'we infer'; the post never mentions DNSSEC
- [ ] The Source was read with WebFetch (a model summary plus quotes), not raw page text, because curl was denied: check every quote against the live page

**Fact-checker could not confirm:**
- [ ] Quote and excerpts were checked with WebFetch (the model returned verbatim snippets), not raw page text, because Bash/curl was unavailable. The Author should check 'a trust store can’t be a static list of root certificate authorities' (curly apostrophe) against the live page
- [ ] 'forge certificates that unpatched clients accept': the Source says the keys 'appear to have been used to forge certificates' and that unpatched clients still accept the chain. Joining the two is a light synthesis
- [ ] Description's 'a week on': based on 50 million probes 'over the last week' in a post dated 30 Oct; 11 days after 19 Oct
- [ ] Date 19 October 2026 has no time or time zone in the Source, so no UTC time is stated
- [ ] Why it matters: the DNSSEC/RFC 5011 claim ('a key the resolver already trusts vouches for the change') is supported by the linked Explainer content/explainers/apis-networking/dnssec-root-key-rollover.md, not by the Source. It is labelled 'we infer'

<details><summary>What the fact-checker changed</summary>

- Fixed: Description: 'Cloudflare patched eMudhra's weak CA keys in hours' -> 'rejected eMudhra's weak CA keys within hours' (Cloudflare patched its own systems to reject the keys; it did not patch eMudhra's keys)
- Fixed: 'eMudhra had generated eight weak keys for its intermediate certificates' -> 'eMudhra's intermediate certificates had eight weak keys' (the Source says 'eMudhra's eight shockingly weak keys' and 'affected intermediate private keys'; it never says 'generated')
- Fixed: Linked the 'recover the private keys' claim to the Source; it had no link (Source: 'could be recovered almost instantly by anyone from the public certificates themselves')
- Fixed: 'Paired probes from its challenge pages' -> 'from some of its challenge pages' (Source: 'two probe connections on a fraction of challenge pages')
- Fixed: '30–40% of Android WebViews' -> '30–40% of embedded browsers such as Android WebView' (Source: 'embedded browsers, like Android WebView, lag behind patching with 30–40% still accepting')
- Fixed: Quote was framed as 'Cloudflare's conclusion'. Reframed: Cloudflare 'points to projects like upki that understand "a trust store can’t be a static list..."'. The apostrophe is now curly to match the Source

</details>

### Not drafted

- none

<details><summary>Leads the editor rejected</summary>

- `grxwqfpc`: Lead page is only a pointer to Thomas Schatzl's 10 Aug post; the G1 heap-free-ratio default change (fits zgc-adaptive-heap-sizing) shipped with JDK 27 in Sept, so it is not news that just happened and the source_url page carries no claims.
- `rkcfnfgr`: GitHub Dependabot status page: no root cause stated, RCA pending.
- `gkdt53er`: GitHub Actions/Pages/Webhooks status page: mitigation not described, no root cause stated.
- `6sjomtrx`: Cause only 'a temporarily overloaded component' with a 9-minute error burst; too thin, and no live Explainer teaches shared-component blast radius (metastable-failures doesn't fit since it recovered at once).
- `3bofscfa`: Discord status page gives no cause or mechanism.
- `covxn3zy`: Copilot agent sessions incident: symptoms only, no root cause.
- `b6h3yjql`: Copilot disruption blamed on an unnamed upstream model provider with no mechanism; no live Explainer fits.
- `3dk6zt3m`: E-Maj 5.0.1 notice names a superuser privilege escalation but gives no mechanism; no live Explainer fits.
- `6qhqbqua`: Node 24.22.2: V8 backports and a stream fix, no CVEs, no design change.
- `rt7w4cwg`: Node 26.13.0: LTS promotion with only metadata changes.
- `yeldokkh`: pgAdmin point release, a tool update with no design change.
- `3f4kc4bb`: Product GA announcement, not an event with a lesson.
- `mm72sfrp`: Link points to the AWS Health Dashboard root, with no incident-specific cause.
- `hf74t2u5`: Quarterly threat report, not an event at a named system.
- `yomhrs3n`: RFC publication (intra-domain SAV) is not an event at a named system, and no live Explainer teaches it.
- `uzpknzyx`: Customer success story, not an event; Explainers capped at 0 today.

</details>


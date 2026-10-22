# 07: QA Target hardening

**What to build:** Nothing on the QA Target can be indexed or mistaken for Production. Because Vercel is only ever the QA Target now (ADR-0004), the Vercel manifest always adds `X-Robots-Tag: noindex` to every response, and the Vercel deploy adapter overlays a disallow-all `robots.txt` in place of the real one. Canonical, Open Graph and JSON-LD URLs keep pointing at `https://hastcode.com/...` because `site_url` is unchanged. The v1 Lighthouse bypass secret and its wizard stage are dropped (Lighthouse now runs against Cloudflare Previews).

Spec: `../spec.md` ("QA Target").

**Blocked by:** 05 (Cloudflare manifest and deploy-adapter rule)

**Status:** resolved

- [x] Build-seam test: the Vercel manifest carries a header rule applying `X-Robots-Tag: noindex` to every path, and the local router (`make dev`) sends it
- [x] The Vercel adapter's file set includes the disallow-all `robots.txt` overlay and a test asserts its content; the Cloudflare set keeps the real `robots.txt`
- [x] Canonical and `og:url` on a QA build equal the production URLs
- [x] Bypass secret references removed from workflows, Makefile, README and wizard

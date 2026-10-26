# 01: Google Analytics 4 replaces Plausible

**What to build:** Production pages report to Google Analytics 4 instead of Plausible, so Raptive can review the site's traffic and AdSense revenue can be read per page. Decisions: ADR-0008 (amends ADR-0006). Grilled 2026-10-27; every recommendation accepted.

- `site.yaml` loses `plausible_domain` and gains `ga_measurement_id` (`G-…`). A Production build (Production and QA Targets) fails, rule `analytics-config`, when it is blank or malformed. Previews and `make dev` never carry analytics and never need the id.
- Every Production page's `<head>` carries a small inline snippet and `<script defer src="https://www.googletagmanager.com/gtag/js?id=…">`. The snippet sets consent defaults (denied for EEA, UK and Switzerland, granted elsewhere) and configures the property only when the page is served from the site's own host. On any other host (the QA Target) it sets Google's opt-out flag, so QA never reports.
- The 60 KB JS budget bounds first-party script only. gtag.js and the AdSense library are named third-party exemptions.
- A `/privacy/` page, linked from every page's footer, covers AdSense, Google Analytics cookies, consent, the newsletter provider and contact.
- Lighthouse CI blocks Google's analytics hosts instead of Plausible. The setup wizard's Plausible step becomes the Google Analytics 4 steps (property, data stream, Measurement ID, AdSense link, consent mode for analytics, Raptive Viewer later).

**Blocked by:** None

**Status:** resolved

- [x] Build-seam tests: Production page carries the snippet and the deferred gtag.js with the configured id; snippet guards on the site host, sets regional consent defaults and configures the id; Preview build carries none of it; Production build with a blank or malformed id fails `analytics-config`; Preview with a blank id passes
- [x] JS budget counts first-party script only; unknown third-party scripts still fail
- [x] `/privacy/` exists, is in the footer of every page, and the real `content/` builds
- [x] No Plausible reference left in code, config, wizard or README (ADRs keep history)
- [x] `make test` passes; preview build passes

## Comments

Author steps, not agent steps: create the property and web data stream in analytics.google.com, paste the Measurement ID into `site.yaml`, link AdSense, turn on "consent mode for analytics purposes" in AdSense Privacy & messaging, review `/privacy/`, Publish, cancel any Plausible account afterwards.

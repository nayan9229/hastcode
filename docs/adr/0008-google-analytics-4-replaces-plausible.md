---
status: accepted
date: 2026-10-27
amends: ADR-0006
---

# Google Analytics 4 replaces Plausible

The first revenue milestone is Raptive, which will not review an application without read-only Google Analytics access and later asks for Viewer on the GA4 property; Plausible alone cannot get us in. GA4 also answers what Plausible does (pageviews per Piece, referrers, country share for Raptive's 50% Tier-1 gate) and, linked from AdSense, ad revenue per page. Running both would double third-party weight for no gain. We decided:

- **GA4 replaces Plausible on Production builds.** `plausible_domain` leaves `site.yaml`; `ga_measurement_id` takes its place, and a Production build (Production and QA Targets) fails while it is blank or malformed, because a missing id would silently lose the history Raptive reviews. Previews and `make dev` still carry no analytics (ADR-0006) and need no id. On any host other than hastcode.com the snippet sets Google's opt-out flag instead of configuring the property, so the QA Target never reports.
- **gtag.js loads with `defer`, not Google's default `async`**, so 157 KB of third-party script never runs while the page is still parsing; the inline snippet queues the pageview until it does.
- **The JS budget splits.** gtag.js is about 449 KB uncompressed (157 KB gzipped), seven times the 60 KB page budget. The budget now bounds first-party script only; third-party analytics and ad scripts are named exemptions, loaded without blocking parsing. AdSense was already in that position.
- **Readers in the EEA, UK and Switzerland default to analytics denied.** The snippet sets `analytics_storage` defaults per region, nothing else; the AdSense Privacy & messaging GDPR message, with "consent mode for analytics purposes" turned on, grants it. Until then GA4 sends cookieless pings only. Ad consent stays with AdSense's own message, untouched by the snippet, so ad serving is unchanged.
- **A privacy page is required** and linked from every page's footer: AdSense, GA4 cookies, consent, the newsletter provider, contact.
- **Search ranking and Core Web Vitals are not GA4's job.** Search Console and CrUX cover them; no web-vitals library ships.

## Considered options

- **Keep Plausible, add GA4 beside it**: two scripts, two dashboards, Plausible's fee, nothing GA4 lacks.
- **Hand-written ~1 KB beacon to GA4's collect endpoint instead of gtag.js**: fits the budget, but unofficial, cookieless (users and sessions wrong) and liable to fail Raptive's "Analytics correctly set up" check.
- **Load gtag.js after first interaction**: fits the lab metrics, but drops every bounce pageview and undercounts the traffic Raptive judges.

## Consequences

- The site is no longer cookieless: the "no cookies anywhere" test and the cookieless comments go.
- PageSpeed Insights against Production will show more third-party JS and a slower lab LCP. Lighthouse CI's blocklist swaps Plausible for Google's analytics hosts; its budgets are unchanged.
- Raptive reviews GA4 history, so the 30 days it judges must fall after this is Published.
- ADR-0005's day-28 checkpoint reads its pageviews and Tier-1 share from GA4. GA4 has no history before its first Publish, so a Profile whose day 28 falls earlier is measured over the days GA4 has.

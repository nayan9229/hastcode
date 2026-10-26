# Checkpoint: TODO(author): owner/name

<!--
Copy to .scratch/hastcode-content-plan/checkpoints/<slug>.md on the
Profile's day 28 (published date + 28 days) and fill in every TODO(author).
One note per Profile. The decision is taken once the first five notes exist,
by the rule below (ADR-0005), not by mood.
-->

- **Profile:** TODO(author): /{vertical}/{slug}/
- **Published:** TODO(author): YYYY-MM-DD
- **Day 28:** TODO(author): YYYY-MM-DD (published + 28 days)
- **Measured on:** TODO(author): YYYY-MM-DD, by TODO(author): who
- **Comparison set:** TODO(author): the Explainers of the same age (published within a week of this Profile), measured over their own first 28 days

## Numbers at day 28

| Metric | This Profile | Median Explainer, same age | Where it comes from |
|---|---|---|---|
| Pageviews, days 1–28 | TODO(author) | TODO(author) | Google Analytics 4: Reports → Engagement → Pages and screens, Views, filtered to the page path, custom range (ADR-0008) |
| Tier-1 share (US, UK, CA, AU, NZ visitors ÷ all) | TODO(author) | TODO(author) | Google Analytics 4: same report, add the Country dimension, same filter and range |
| Site median Tier-1 share, same range | — | TODO(author) | Google Analytics 4: Reports → User attributes → Country, Views, whole site |
| Bing Webmaster AI Performance citations | TODO(author) | — | Bing Webmaster Tools, AI Performance; name each answer engine and cited page |
| Inbound links (referring pages) | TODO(author) | — | Search Console, Links; list them below |

### Search Console impressions for the templated queries

Search Console, Performance, page filter on the Profile, days 1–28.

| Query | Impressions | Clicks | Average position |
|---|---|---|---|
| `{repo}` / "what is {repo}" | TODO(author) | TODO(author) | TODO(author) |
| "is {repo} production ready" | TODO(author) | TODO(author) | TODO(author) |
| "{repo} alternatives" | TODO(author) | TODO(author) | TODO(author) |
| "{repo} vs {alternative}" (each seen) | TODO(author) | TODO(author) | TODO(author) |

### Inbound links

- TODO(author): referring URL, or "none"

### AI answer engine citations

- TODO(author): engine, query or prompt if known, date, or "none"

## The rule (ADR-0005, copied verbatim; read its "Plausible" as Google Analytics 4 since ADR-0008)

> The four-week checkpoint, measured at each Profile's day 28, compares Plausible pageviews and Tier-1 share with the median Explainer of the same age, plus Search Console impressions for the templated queries, Bing Webmaster AI Performance citations and inbound links. Health score and badge proceed only if the median Profile matches the median Explainer's pageviews **and** its Tier-1 share is within 10 points of the site median, or if any Profile is cited by an AI answer engine. Otherwise cadence drops and the product half stays shelved.

## This Profile against the rule

- Pageviews at or above the median Explainer of the same age: TODO(author): yes / no
- Tier-1 share within 10 points of the site median: TODO(author): yes / no (difference in points)
- Cited by an AI answer engine: TODO(author): yes / no

The rule reads the **median** Profile across the first five, so one note does not decide it. Once all five notes exist, record the medians and the decision (proceed, or drop cadence and keep the product half shelved) in the last note and in the feature's scratch directory.

## Notes

TODO(author): anything that skews the numbers (a Hacker News spike, a launch, an outage), with links.

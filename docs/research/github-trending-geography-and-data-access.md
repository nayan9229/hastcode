# GitHub developer geography, trending-data access, repo-health data sources, and Cloudflare job limits

Research for hastcode.com (static Go build, Cloudflare Workers Static Assets + one Worker, ad-monetised; Raptive target requires 50% of traffic from US/UK/CA/NZ/AU — see `ad-rpm-tech-audience.md`). Question: should we add "repo profile" pages for selected trending GitHub repos, and what can we legitimately fetch to build them?

All sources accessed **2026-10-22**. Labels: **[P]** primary (GitHub/Cloudflare/Google/PSF-owned page, official API reference, raw dataset) / **[S]** secondary (trade press, third-party blog). **[observed]** = a read-only request I made against the live API on 2026-10-22 (the response, not a doc). **[inferred]** = my arithmetic or deduction from primary facts. **UNSOURCED** = could not be tied to a page. GitHub and Cloudflare pages carry their own "last updated" dates where shown; re-check before committing to a number.

---

## TL;DR

- **Tier-1 (US+UK+CA+AU+NZ) accounts are ~23% of GitHub developer accounts and falling** — 23.3% in 2026 Q1 vs 25.7% in 2023 Q1, computed from GitHub's own Innovation Graph CSV. India+China+Brazil alone are 24.1%. India is 13.9% and growing 39%/yr vs the US at 16.9% and 25%/yr. A page whose audience mirrors GitHub's population would fail Raptive's 50% rule by a factor of two. [P: raw CSV] + [inferred]
- **Stack Overflow's 2025 raw data skews far more Western** — tier-1 = 32.8% of the 35,437 respondents who gave a country (US 20.4%, UK 5.8%, CA 3.7%, AU 2.3%, NZ 0.7%); India 7.2%, China 0.7%, Brazil 2.3%. But it is a self-selected sample recruited through Stack Overflow's own channels, not a population estimate. [P: raw CSV]
- **GitHub says India is the fastest-growing large community** and will overtake the US in developer count — "by 2027" (Octoverse 2023), "by 2028" (Octoverse 2024), "within the next few years" (Octoverse 2025, which also says India "overtook the United States as the largest contributor base to public and open source projects"). Octoverse 2025: India "added more than 5 million developers this year (over 14% of all new accounts)". [P]
- **There is no official API for github.com/trending.** The REST API root lists no trending resource, the GraphQL `Repository` reference has no trending field, and GitHub's docs describe Trending only as a web page. GitHub tried to retire it on 30 Sep 2022, reversed after a 1,200-upvote thread, and said the backend "has stagnated" in usage and needs re-investment. [P + observed]
- **Scraping github.com is defined and policed by the Acceptable Use Policies, not banned outright**: "Scraping refers to extracting information from our Service via an automated process... Scraping does not refer to the collection of information through our API." Use of scraped data is restricted to researchers (open-access publications only) and archivists; the ToS adds "You may not share API tokens to exceed GitHub's rate limitations" and "One person or legal entity may maintain no more than one free Account". `robots.txt` does not disallow `/trending`. A commercial ad-funded site has no explicit scraping allowance. [P]
- **Rate limits (REST): 60/h unauthenticated per IP, 5,000/h per PAT, 5,000/h base per GitHub App installation (+50 per repo/user beyond 20, cap 12,500), secondary limits 100 concurrent, 900 points/min, 90 s CPU per 60 s.** GraphQL: 5,000 points/h; unauthenticated GraphQL is refused (observed `x-ratelimit-limit: 0`). Search endpoints: 30/min authenticated, 10/min unauthenticated (observed). [P + observed]
- **The legitimate trending substitute is two new endpoints, not stargazer lists.** Since July 2026 `GET /repos/{o}/{r}/stargazers` is "limited to admins and collaborators" (observed: 401 unauthenticated). GitHub added `GET /repos/{o}/{r}/stargazers/history` (weekly totals + per-day array, max 30 weeks/page, 100 pages) and `/stargazers/count`, both usable "without authentication ... if only public resources are requested". Pair with `GET /search/repositories?q=pushed:>DATE stars:>N&sort=stars` (1,000-result cap). [P]
- **Health-score data is mostly free and unauthenticated**: OSV.dev ("Currently there are no limits on the API"), deps.dev (no key; Google API ToS; `GetProject` returns stars, forks, license, OpenSSF Scorecard; `v3alpha` `:dependents` gives direct/indirect dependent counts for npm, Cargo, Maven, PyPI), ecosyste.ms (no key; 300 req/h on the pricing page, but observed `x-ratelimit-limit: 15000` on packages and `5000` on repos; CC BY-SA 4.0 data), npm (`api.npmjs.org`, no documented limit), pypistats.org (PSF-operated, observed limit 60, 180-day retention), crates.io (1 req/s, identifying User-Agent with contact). GitHub "dependents" exists only as HTML (`/network/dependents`) and is "approximate"; the community-profile endpoint does **not** detect `SECURITY.md` — use GraphQL `isSecurityPolicyEnabled`. `stats/contributors` returns **202** while computing and zeroes additions/deletions above 10,000 commits. [P + observed]
- **Cloudflare Free plan bounds a scheduled fetcher hard**: 5 Cron Triggers per account, 1-minute cron granularity, **10 ms CPU per Cron Trigger invocation** (network wait excluded), 50 subrequests per invocation, 100,000 requests/day, KV 1,000 writes + 100,000 reads/day and 1 GB, D1 100,000 rows written + 5 M rows read/day and 500 MB per database. Paid ($5/mo) lifts CPU to 30 s (<1 h interval) or 15 min (≥1 h interval), subrequests to 10,000, and makes KV/D1 metered. Cron Triggers are ✅ on Workers (❌ on Pages) and `triggers.crons` + `assets` are sibling keys of one Worker config, so one Worker can serve the site and run the job — [inferred] from those pages; no Cloudflare page states the combination in one sentence. [P + inferred]

---

## 1. Developer geography

### 1.1 GitHub Innovation Graph (GitHub's own raw quarterly data) — the best primary source

GitHub publishes developer counts per economy per quarter, Q1 2020 → Q1 2026, as CSV under CC0-1.0. [P] https://innovationgraph.github.com/global-metrics/developers ; data file https://github.com/github/innovationgraph/blob/main/data/developers.csv (raw: https://raw.githubusercontent.com/github/innovationgraph/main/data/developers.csv, 5,657 rows, downloaded 2026-10-22).

Definitions from the datasheet (v1.0.10, dated 2026-05-08) [P] https://github.com/github/innovationgraph/blob/main/docs/datasheet.md :

- Developers = "the number of developer accounts located in a given economy based on mode daily location."
- "This count excludes users that are bots or otherwise flagged as "spammy" within internal systems."
- Location is "determined by their IP address"; "the last known location of the developer is carried forward on a daily basis."
- "VPNs and other means of altering or hiding IP addresses distort the metrics."
- Economies are "only reported when there are 100 or more unique developers performing the relevant activity".
- "The dataset will be updated quarterly"; "made available under a CC0-1.0 license."

My computation from the CSV (sum over the 229 reported economies, excluding the `EU` aggregate row; the CSV has no world-total row, so "share" = share of the reported-economy sum) [P data, inferred arithmetic]:

| Economy | 2026 Q1 developers | Share | YoY vs 2025 Q1 |
|---|---:|---:|---:|
| United States | 32,410,529 | 16.9% | +25.4% |
| India | 26,727,004 | 13.9% | +38.8% |
| China | 11,351,939 | 5.9% | +9.3% |
| Brazil | 8,317,960 | 4.3% | +36.4% |
| United Kingdom | 5,559,767 | 2.9% | +24.5% |
| Japan | 5,433,912 | 2.8% | +36.5% |
| Indonesia | 5,270,032 | 2.7% | +32.4% |
| Germany | 5,255,625 | 2.7% | +31.2% |
| Russia | 4,627,668 | 2.4% | — |
| France | 4,034,001 | 2.1% | — |
| Canada | 3,991,305 | 2.1% | +24.3% |
| Australia | 2,368,357 | 1.2% | +26.6% |
| New Zealand | 441,201 | 0.2% | +23.6% |
| **US+UK+CA+AU+NZ** | **44,771,159** | **23.3%** | |
| **India+China+Brazil** | **46,396,903** | **24.1%** | |

Trend of the tier-1 share (same method): 2023 Q1 **25.7%** → 2024 Q1 **25.1%** → 2025 Q1 **24.3%** → 2026 Q1 **23.3%**. India alone: 11.7% → 12.4% → 13.1% → 13.9%. China: 9.0% → 7.7% → 7.1% → 5.9%. Fastest YoY growth among the larger economies in 2026 Q1: Pakistan +45.3%, Philippines +41.3%, India +38.8%, Japan +36.5%, Brazil +36.4%, Nigeria +35.8%, Vietnam +35.2%, Bangladesh +34.5%. Reported-economy total: 98.1 M (2023 Q1) → 122.3 M → 147.2 M → 192.3 M (2026 Q1). Top 20 in 2026 Q1 in order: US, IN, CN, BR, GB, JP, ID, DE, RU, FR, CA, HK, KR, VN, MX, PH, ES, TR, PK, AU.

Caveat: "developer" is any non-spam account, not an active or professional developer, and location is IP-based (VPN-distorted). Treat these as audience-shape, not audience-size, numbers.

### 1.2 Octoverse reports (narrative numbers)

**Octoverse 2023** — [P] https://github.blog/news-insights/research/the-state-of-open-source-and-ai/ (published 2023-11-08, updated 2024-07-30; data 2022-10-01 → 2023-09-30)

- "since last year we've seen a nearly 26% increase in all global developer accounts on GitHub."
- US: "With 20.2 million developers and a 21% increase in developer growth over the past year" the U.S. "continues to have the largest developer community globally."
- India: over 13.2 M, +36% — "3.5 million new developers joining GitHub in 2023". Brazil over 4.3 M, +30%. UK over 3.4 M, +21%. Germany over 2.9 M, +22%. Japan over 2.8 M, +31%. Australia over 1.4 M, +21%. Singapore over 1 M, +39%. Nigeria over 868 K, +45% — "the largest worldwide increase". Canada and New Zealand: no figures given.
- Projection: India expected to overtake the United States as the largest developer community on GitHub by 2027 (the text assumes "linear growth to forecast").
- Method: a "developer" is "anyone with a GitHub account"; geography "is based on the last known network location of individual users and organization profiles."

**Octoverse 2024** — [P] https://github.blog/news-insights/octoverse/octoverse-2024/ (published 2024-10-29, updated 2025-10-28)

- "The United States continues to have the most developers worldwide on GitHub." No US count given.
- Chart caption: "The United States consistently holds the top position, followed by India, China, and Brazil. Other countries in the top 10 include the United Kingdom, Russia, Germany, Indonesia, Japan, and Canada, with France dropping out of the top 10 in 2022."
- Year-over-year table (developers, growth): India >17 M, 28%; China >9 M, 10%; Brazil >5.4 M, 27%; United Kingdom >4 M, 19%; Germany >3.5 M, 21%; Japan >3.5 M, 23%; Indonesia >3.5 M, 23%; France >2.8 M, 20%; Mexico >1.9 M, 21%; Spain >1.8 M, 24%; Philippines >1.7 M, 29%; Turkey >1.7 M, 19%; Nigeria >1.1 M, 28%; Argentina >1.1 M, 22%; Colombia >1 M, 25%; Egypt >990 K, 25%; South Africa >664 K, 23%; Peru >583 K, 27.5%; Morocco >556 K, 25%; Kenya >393 K, 33%; UAE "32% year-over-year increase". Canada, Australia, New Zealand: no counts; "the Philippines (#18) overtook Australia (#19)".
- "India is now on track to surpass the United States in number of developers on GitHub by 2028 (previously, we had predicted it would overtake the United States by 2027...)."
- "we have seen greater growth outside the United States every year since 2013".
- Open-source contributions by region: US, India, Germany, Brazil, UK, China, France, Canada, Hong Kong SAR, Japan.
- Method note: economies reported only with "100 or more unique developers"; the page does not say how location is determined.

**Octoverse 2025** — [P] https://github.blog/news-insights/octoverse/octoverse-a-new-developer-joins-github-every-second-as-ai-leads-typescript-to-1/ (published 2025-10-28, updated 2026-02-28; data 2024-09-01 → 2025-08-31)

- "180 million-plus developers now work and build on GitHub." Over 36 million joined "in the past year" (+23% YoY). Per minute: "~25 developers joined from APAC, ~12 from Europe, ~6.5 from Africa and the Middle East, and ~6 from LATAM".
- "India alone added more than 5 million developers this year (over 14% of all new accounts)"; India is "on track to account for one in every three new developers on GitHub by 2030"; "India overtook the United States as the largest contributor base to public and open source projects on GitHub"; it is "poised to overtake the US developer population within the next few years."
- Chart alt text: "The United States remains #1 with 10.2 million to 28 million developers, followed by India rising to #2 with 4.5 million to 21.9 million" (2020 → 2025). Brazil "6.89 million developers ... (now #4 globally)". Indonesia "4.37 million developers ... (up from 0.9 million in 2020)".
- 2020→2025 multiples: "the USA, UK, and Canada have more than doubled"; "Japan and Germany have more than tripled"; "India, Brazil, and Indonesia have more than quadrupled".
- 2030 projections: India "57.5 million developers by 2030"; "The United States will be the second-largest community with more than 40 million developers expected" (the chart alt text says 54.7 M — the page is internally inconsistent); Brazil 19.6 M, Japan 11.7 M, UK 11 M. Egypt, Nigeria, Kenya, Morocco projected to add "millions".
- Method: "Developers are mapped to countries via self-reported locations, standardized to ISO country codes" (note: a different method from the Innovation Graph's IP-based assignment and the 2023 report's "network location"); "Developer: Anyone with a GitHub account"; forecasts backtested at "less than 30% Mean Absolute Percentage Error".

### 1.3 Stack Overflow Developer Survey (raw data)

**2025** — results page [P] https://survey.stackoverflow.co/2025/ ; methodology [P] https://survey.stackoverflow.co/2025/methodology/ : "The survey was fielded from May 29, 2025 to June 23, 2025." "49,009 responses from 177 countries are used in these survey results." "Respondents were recruited primarily through channels owned by Stack Overflow" (onsite messaging, blog, newsletter, banner ads, social; a Reddit ad campaign under 2%); "highly-engaged users on Stack Overflow were more likely to notice the prompts to take the survey." Licence: "Survey data licensed under ODbL 1.0" [P] https://survey.stackoverflow.co/ . Raw data: https://github.com/StackExchange/Survey/tree/main/packages/archive/2025 (`results.csv` is a Git LFS object; actual bytes at https://media.githubusercontent.com/media/StackExchange/Survey/main/packages/archive/2025/results.csv, 140,893,245 bytes, downloaded 2026-10-22).

My computation from `results.csv` (`Country` column; 49,191 rows, 35,437 with a country, 13,754 `NA`; shares are of the 35,437) [P data, inferred arithmetic]:

| Country | n | Share |
|---|---:|---:|
| United States of America | 7,233 | 20.4% |
| Germany | 3,025 | 8.5% |
| India | 2,547 | 7.2% |
| United Kingdom | 2,042 | 5.8% |
| France | 1,409 | 4.0% |
| Canada | 1,305 | 3.7% |
| Ukraine | 964 | 2.7% |
| Poland | 888 | 2.5% |
| Netherlands | 867 | 2.4% |
| Italy | 835 | 2.4% |
| Brazil | 825 | 2.3% |
| Australia | 804 | 2.3% |
| New Zealand | 255 | 0.7% |
| China | 255 | 0.7% |
| **US+UK+CA+AU+NZ** | **11,639** | **32.8%** |
| **India+China+Brazil** | **3,627** | **10.2%** |

These match the shares printed on the 2025 results page (US 20.4%, Germany 8.6%, India 7.2%, UK 5.8%, France 4%, Canada 3.7%, Ukraine 2.7%, Poland 2.5%, Netherlands 2.5%, Italy 2.4%). [P] https://survey.stackoverflow.co/2025/

**2024** — [P] https://survey.stackoverflow.co/2024/ : ~65,437 respondents from 185 countries (May 2024). US 18.9%, Germany 8.4%, India 7.2%, UK 5.5%, Ukraine 4.6%, Canada 3.6%, France 3.6%, Poland 2.6%, Netherlands 2.5%, Brazil 2.3%. Raw data: https://github.com/StackExchange/Survey/tree/main/packages/archive/2024 (not recomputed).

### 1.4 What it means for the Raptive 50% rule [inferred]

- GitHub's own population: tier-1 ≈ 23% and declining ~0.8 pt/year. The Stack Overflow sample (English-language, Western-skewed) is ≈ 33%. Neither approaches 50%.
- Content whose natural audience is "people who look at GitHub trending" will inherit GitHub's shape unless the pages are written and ranked for tier-1 search queries. India alone is ~14% of GitHub accounts and the fastest-growing large cohort; China is 6% of accounts but 0.7% of SO respondents, consistent with Chinese developers rarely reaching English-language sites.
- The decision therefore hinges less on GitHub's geography than on what fraction of *search* traffic for "repo X" queries comes from tier-1 countries — not answerable from these sources (UNSOURCED; needs Search Console data once pages exist).

---

## 2. GitHub trending data access

### 2.1 Is there an official API for github.com/trending? No.

- The REST API root (`GET https://api.github.com/`) enumerates resource URL templates (search, starred, etc.); it contains no `trending` resource. [observed 2026-10-22]
- The GraphQL `Repository` object reference contains no field with "trending" in its name or description. [P, checked 2026-10-22] https://docs.github.com/en/graphql/reference/repos
- GitHub's docs describe Trending only as a UI: you can "browse popular repositories of the day by visiting Trending"; no API or feed is mentioned, and the page does not say how the ranking is computed. [P] https://docs.github.com/en/get-started/exploring-projects-on-github/finding-ways-to-contribute-to-open-source-on-github
- 2022 deprecation and reversal, primary thread [P] https://github.com/community/community/discussions/31644 (marked "Answered" by martinwoodward; 303 comments, 1,238 👍 on the OP):
  - OP (daiyam, 2022-09-01) quoting the site banner: "Heads up! This Trending tab is being deprecated. Due to low usage of Trending Repositories and Trending Developers, this tab will no longer be available beginning September 30, 2022."
  - Martin Woodward, GitHub (2022-09-02): "Given the feedback here, it's clear we need to look again at the plan to deprecate this. The team is going to re-evaluate and see if we can come up with some other options. We'll get back to you before the end of the month. In the meantime we'll take the deprecation warning off the page until we have a clearer idea of what the long term plan is." and "from all our telemetry, usage of the Trending page has been very low since it was created and usage has stagnated, despite more and more people using GitHub."
  - The follow-up plan ("Keep trending projects — Leave it where it is for now but investigate ways of making it more prominent and more useful. From our side we'll need to re-invest to make the backend that powers it more stable and more maintainable.") appears in the thread only as a user's quotation of the staff answer (tommywalkie, 2023-07-10); the full staff answer is collapsed in the fetched HTML, so treat the wording as **[S]**. DevClass reported the same statements and that "the data pipeline that powers the page has historically proven problematic to maintain" [S] https://www.devclass.com/ci-cd/2022/09/05/github-trending-tab-gets-temporary-reprieve-after-users-perplexed-and-upset-by-threatened-removal/1624832
  - Nothing in the thread or in GitHub's docs promises the page's continued existence or a stable HTML structure.
- `robots.txt` [P] https://github.com/robots.txt : no `Disallow` line mentions `/trending` or `/explore` for any user-agent. The `User-agent: *` group (495 Disallow lines in total) does disallow `/search$`, `/search/advanced`, `/*/*/pulse*`, `/*/archive/`, `/*/tree/`, `/*/raw/`, `/*/blame/`, `/*/*/commits/`, `/*.atom$`, `/gist/` and others; `Bytespider` gets `Disallow: /`. robots.txt is a crawl-etiquette signal; the AUP below governs use.

### 2.2 Terms of Service and Acceptable Use Policies — exact clauses

**Terms of Service** [P] https://docs.github.com/en/site-policy/github-terms/github-terms-of-service — "Effective date: April 27, 2026"

- B.3 Account Requirements: "You must be a human to create an Account. Accounts registered by "bots" or other automated methods are not permitted. We do permit machine accounts:" — "A machine account is an Account set up by an individual human who accepts the Terms on behalf of the Account, provides a valid email address, and is responsible for its actions. A machine account is used exclusively for performing automated tasks. Multiple users may direct the actions of a machine account, but the owner of the Account is ultimately responsible for the machine's actions. You may maintain no more than one free machine account in addition to your free Personal Account." — "One person or legal entity may maintain no more than one free Account (if you choose to control a machine account as well, that's fine, but it can only be used for running a machine)."
- H. API Terms: "Abuse or excessively frequent requests to GitHub via the API may result in the temporary or permanent suspension of your Account's access to the API. GitHub, in our sole discretion, will determine abuse or excessive usage of the API. We will make a reasonable attempt to warn you via email prior to suspension." — "You may not share API tokens to exceed GitHub's rate limitations." — "You may not use the API to download data or Content from GitHub for spamming purposes, including for the purposes of selling GitHub users' personal information, such as to recruiters, headhunters, and job boards."
- D.9 Access Reciprocity (new): applies to anyone "using automated means to access, collect, or otherwise use ... any publicly accessible Content from the Service for the purpose of developing or training any commercially available artificial intelligence model"; it "does not apply to Access solely for the purpose of academic research or if ... the number of monthly active users of the products or services made available by you is less than 700 million". Not relevant to a repo-profile fetcher, noted for completeness.

**Acceptable Use Policies** [P] https://docs.github.com/en/site-policy/acceptable-use-policies/github-acceptable-use-policies — the page shows no effective date (its only date is an HTML comment "localization update 2026-03-11")

- §4 Spam and Inauthentic Activity: GitHub does not allow content or activity that is "automated excessive bulk activity and coordinated inauthentic activity, such as spamming; ... inauthentic interactions, such as fake accounts and automated inauthentic activity; rank abuse, such as automated starring or following; ... using our servers for any form of excessive automated bulk activity, to place undue burden on our servers through automated means, or to relay any form of unsolicited advertising or solicitation through our servers".
- §6 Services Usage Limits: "You will not reproduce, duplicate, copy, sell, resell or exploit any portion of the Service, use of the Service, or access to the Service without our express written permission."
- §7 Information Usage Restrictions (verbatim, complete): "You may use information from our Service for the following reasons, regardless of whether the information was scraped, collected through our API, or obtained otherwise: Researchers may use public, non-personal information from the Service for research purposes, only if any publications resulting from that research are open access. Archivists may use public information from the Service for archival purposes. Scraping refers to extracting information from our Service via an automated process, such as a bot or webcrawler. Scraping does not refer to the collection of information through our API. Please see Section H of our Terms of Service for our API Terms. You may not use information from the Service (whether scraped, collected through our API, or obtained otherwise) for spamming purposes, including for the purposes of sending unsolicited emails to users or selling personal information, such as to recruiters, headhunters, and job boards. Your use of information from the Service must comply with the GitHub Privacy Statement."
- §9 Excessive Bandwidth Use: "If we determine your bandwidth usage to be significantly excessive in relation to other users of similar features, we reserve the right to suspend your Account, throttle your file hosting, or otherwise limit your activity until you can reduce your bandwidth consumption."
- Neither document contains a clause that literally says "you may not use multiple accounts to circumvent rate limits"; the operative text is the one-free-Account rule (B.3), the token-sharing ban (H), and §4's "fake accounts". Multiple PATs on one account share that account's 5,000/h budget anyway (see 2.3).

Reading [inferred]: the AUP's enumerated allowances for scraped data are "researchers" (open-access only) and "archivists". An ad-funded publication fits neither. Fetching through the REST/GraphQL API under one account's token, within the published limits, is the path the policy explicitly routes you to ("Please see Section H").

### 2.3 Rate limits (current)

**REST primary limits** [P] https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api

| Caller | Limit |
|---|---|
| Unauthenticated | "The primary rate limit for unauthenticated requests is 60 requests per hour." (per originating IP) — [observed] `/rate_limit` → `core.limit: 60` |
| PAT / user-authorised app | 5,000 requests per hour per user (15,000 when the app is owned by a GitHub Enterprise Cloud org and you are a member). All of a user's PATs and apps draw on the same budget. |
| GitHub App installation | 5,000/h base; +50/h per repository beyond 20 and +50/h per org user beyond 20; "cannot exceed 12,500 per hour"; 15,000/h on Enterprise Cloud orgs |
| OAuth app client-id/secret (public data only) | 5,000/h per app (15,000 Enterprise Cloud) |
| `GITHUB_TOKEN` in Actions | 1,000/h per repository (15,000 Enterprise Cloud) |

**Secondary limits** (same page; "may change without notice"): "No more than 100 concurrent requests are allowed. This limit is shared across the REST API and GraphQL API." — "No more than 900 points per minute are allowed for REST API endpoints" and "no more than 2,000 points per minute are allowed for the GraphQL API endpoint." — "No more than 90 seconds of CPU time per 60 seconds of real time is allowed" ("No more than 60 seconds of this CPU time may be for the GraphQL API"). — "no more than 80 content-generating requests per minute" and "no more than 500 content-generating requests per hour". — "No more than 2,000 OAuth access token requests per hour". Points: REST GET/HEAD/OPTIONS 1, REST POST/PATCH/PUT/DELETE 5, GraphQL query 1, GraphQL mutation 5. On exceeding: 403 or 429; honour `retry-after`, else `x-ratelimit-reset`, else "wait at least one minute"; continuing "may result in the banning of your integration".

**Best practices** [P] https://docs.github.com/en/rest/using-the-rest-api/best-practices-for-using-the-rest-api : "you should make requests serially instead of concurrently"; for mutations "wait at least one second between each request"; "Making a conditional request does not count against your primary rate limit if a `304` response is returned" (authenticated requests only).

**GraphQL** [P] https://docs.github.com/en/graphql/overview/rate-limits-and-node-limits-for-the-graphql-api : "5,000 points per hour per user" (PATs and user-authorised apps); 10,000 for Enterprise Cloud-owned apps; installations 5,000 base, +50 per repo/user beyond 20, cap 12,500 (10,000 on Enterprise Cloud); `GITHUB_TOKEN` 1,000 points/h per repository. Cost = number of requests needed to satisfy every connection at its `first`/`last` size ÷ 100, rounded, minimum 1; every connection must pass `first` or `last` between 1 and 100; max 500,000 nodes per call. Unauthenticated GraphQL is not available: [observed] `POST /graphql` without a token → HTTP 403, `x-ratelimit-limit: 0`, `x-ratelimit-resource: graphql`; `/rate_limit` shows `"graphql": {"limit": 0}`.

### 2.4 Search repositories as a trending substitute

[P] https://docs.github.com/en/rest/search/search

- `GET /search/repositories` — `sort` "Can be one of: `stars`, `forks`, `help-wanted-issues`, `updated`" (default best match); `order` `desc`/`asc` ("ignored unless `sort` is provided"); `per_page` "max 100"; the API "provides up to 1,000 results for each search".
- Rate limit section: authenticated "you can make up to 30 requests per minute for all search endpoints except for the Search code endpoint" (code search 10/min); unauthenticated "up to 10 requests per minute". Search draws on its own `search` resource, separate from `core` — [observed] `GET /search/repositories?q=stars:>1000+pushed:>2026-10-15&sort=stars&order=desc&per_page=1` unauthenticated → 200 with `x-ratelimit-limit: 10`, `x-ratelimit-resource: search`. Secondary limits (§2.3) also apply.
- Qualifiers [P] https://docs.github.com/en/search-github/searching-on-github/searching-for-repositories : `created:` (example `created:<2011-01-01`), `pushed:` ("The `pushed` qualifier will return a list of repositories, sorted by the most recent commit"; examples `pushed:>2013-02-01`, `pushed:>=2013-03-06`), `stars:` (`stars:500`, `stars:10..20`, `stars:>=500`), `language:`, `archived:true|false`, `fork:true|only`. "Date formatting must follow the ISO8601 standard, which is `YYYY-MM-DD`"; optional `THH:MM:SS+00:00`.
- So `q=created:>2026-10-15 stars:>50&sort=stars` ("new repos that got stars fast") and `q=pushed:>2026-10-15 stars:>5000&sort=updated` are legitimate, documented, 1-point queries. What search cannot give is star *velocity* over a window — for that see 2.5.

### 2.5 Star velocity: stargazer lists are closed, star history is open

- Starring docs [P] https://docs.github.com/en/rest/activity/starring : "In July 2026, to ensure responsible use of this data and to protect users from misuse, we are introducing new access restrictions to the endpoints related to starring. Access to the stargazers listing endpoints will be limited to admins and collaborators." [observed] `GET /repos/golang/go/stargazers?per_page=1` unauthenticated → HTTP 401 `"Requires authentication"`.
- Changelog 2026-06-30 [P] https://github.blog/changelog/2026-06-30-upcoming-access-restrictions-to-public-api-endpoints-and-ui-views/ : affected `/repos/{owner}/{repo}/stargazers` and `/repos/{owner}/{repo}/subscribers` ("limited to admins and collaborators"), `/users/{username}/subscriptions` (deprecated; "will remain accessible but will return empty responses"), and the `/stargazers`, `/stargazers/you_know`, `/watchers` UI views; reason: the information "has increasingly been misused to collect user data for spam activities"; "Some users may begin receiving empty responses or a 403 Forbidden status."
- Changelog 2026-09-04 [P] https://github.blog/changelog/2026-09-04-new-api-endpoint-provides-privacy-safe-star-history-data : "Earlier this year, stargazer listing endpoints were restricted to admins and collaborators to protect user privacy." Introduces the star history endpoint "without exposing stargazer identities."
- New endpoints (starring docs, API version 2026-03-10) [P] https://docs.github.com/en/rest/activity/starring :
  - `GET /repos/{owner}/{repo}/stargazers/count` — "Gets the current number of users who have starred the repository. Users who previously starred the repository but later removed their star are not included." Fine-grained token needs "Metadata" read; "This endpoint can be used without authentication or the aforementioned permissions if only public resources are requested."
  - `GET /repos/{owner}/{repo}/stargazers/history` — "Returns repository stars grouped by calendar weeks, most recent first. Pages move backward toward the repository's creation week, and weeks within a page are ordered newest to oldest, so concatenating pages produces one continuous series. Weeks without stars contain zero counts. Week and day boundaries are not guaranteed to align with UTC. The `days` array contains the number of stars created on each day of the week, starting on Sunday." `per_page` "max 30", `page` "max 100"; 422 = "Validation failed, or the endpoint has been spammed." Same no-auth-for-public rule. Example item: `{"week": 1754784000, "total": 19, "days": [0,12,7,0,0,0,0]}`.
  - I could not probe these unauthenticated (my IP's 60/h budget was spent on earlier probes); the doc text above is the claim.
- Cost model [inferred]: one `history` call per candidate repo (1 point) gives the last 30 weeks of daily star counts, which is exactly what a "trending" score needs; 1,000 candidates/day from search = ~1,000 points/day, well inside a PAT's 5,000/h and the 900-points/min secondary limit if sent serially.

---

## 3. Data sources for a repo health score

| Source | URL | Auth | Rate limit | Returns |
|---|---|---|---|---|
| GitHub Advisory DB (REST) | `GET https://api.github.com/advisories` | none for public data ([observed] 200 unauthenticated, `x-ratelimit-limit: 60`) | core limits (§2.3) | advisories filtered by `ecosystem`, `affects=pkg[@ver]` (≤1,000), `severity`, `published`, `updated`, `type` |
| OSV.dev | `POST https://api.osv.dev/v1/query`, `/v1/querybatch`, `GET /v1/vulns/{id}` | none | "Currently there are no limits on the API." | OSV records for a package (+version) or commit |
| deps.dev | `https://api.deps.dev/v3/...`, `/v3alpha/...` | none | none stated; Google API ToS | package/version metadata, resolved deps, `GetProject` (stars, forks, issues, license, Scorecard), `:dependents` counts (v3alpha) |
| ecosyste.ms packages | `https://packages.ecosyste.ms/api/v1/packages/lookup?repository_url=…` | none | pricing page: 300/h free; [observed] header 15,000/h ("polite") | downloads, downloads_period, dependent_packages_count, dependent_repos_count, latest release, maintainers, advisories |
| ecosyste.ms repos | `https://repos.ecosyste.ms/api/v1/hosts/GitHub/repositories/{owner}/{repo}` | none | [observed] header 5,000/h ("anonymous") | stars/forks/issues, pushed_at, license, `metadata.files` (readme, license, security, funding…), `commit_stats`, last_synced_at |
| npm downloads | `https://api.npmjs.org/downloads/point/{period}/{pkg}` | none | not documented | download totals; bulk ≤128 unscoped packages |
| PyPI (official) | BigQuery `bigquery-public-data.pypi.file_downloads` | Google account | 1 TB queries/month free | raw per-file download logs |
| pypistats.org (PSF-operated) | `https://pypistats.org/api/packages/{pkg}/recent` | none | "IP-based rate limiting"; [observed] `x-ratelimit-limit: 60` | last_day / last_week / last_month; 180-day history |
| crates.io | `https://crates.io/api/v1/crates/{name}` | none | 1 request/second; identifying User-Agent | `downloads` (all-time), `recent_downloads` (90 days), `/downloads` per-day 90 days, `/reverse_dependencies` |
| GitHub dependents | `https://github.com/{o}/{r}/network/dependents` | web only | n/a | HTML only; counts "approximate" |
| GitHub community profile | `GET /repos/{o}/{r}/community/profile` | none for public ([observed] 200) | core | health_percentage + files (no SECURITY.md) |
| GitHub security policy | GraphQL `Repository.isSecurityPolicyEnabled`, `securityPolicyUrl` | token required (GraphQL) | 5,000 points/h | boolean + URL |
| GitHub releases | `GET /repos/{o}/{r}/releases`, `/releases/latest` | none for published | core | tag, published_at, prerelease, assets[].download_count |
| GitHub stats | `GET /repos/{o}/{r}/stats/contributors` etc. | none for public | core; **202 while computing** ([observed]) | per-contributor weekly a/d/c |

Details and exact wording:

**GitHub Advisory Database** [P] https://docs.github.com/en/rest/security-advisories/global-advisories — `GET /advisories`; `affects`: "If specified, only return advisories that affect any of package or package@version", up to 1,000 packages; `ecosystem` values such as `npm`, `pip`, `maven`, `go`, `rust`; `type` `reviewed` (default), `malware`, `unreviewed` — "By default, all responses will exclude advisories for malware"; `per_page` "max 100". Each advisory carries `ghsa_id`, `cve_id`, `severity`, `cvss_severities.cvss_v3/cvss_v4` (`vector_string`, `score`), `vulnerabilities[]` with `package{ecosystem,name}`, `vulnerable_version_range`, `first_patched_version`, plus `epss` ([observed] in a live response). Database background [P] https://docs.github.com/en/code-security/security-advisories/working-with-global-security-advisories-from-the-github-advisory-database/about-the-github-advisory-database : ecosystems Composer, Erlang, Go, GitHub Actions, Maven, Npm, NuGet, Pip, Pub, RubyGems, Rust, Swift; "Security advisories are published as JSON files in the Open Source Vulnerability (OSV) format"; unreviewed advisories are "published automatically ... directly from the National Vulnerability Database feed". The repo export is CC-BY 4.0: "This project is licensed under the terms of the CC-BY 4.0 open source license." [P] https://github.com/github/advisory-database

**OSV.dev** [P] https://google.github.io/osv.dev/api/ : "Currently there are no limits on the API." "The API has a response size limit of 32MiB when using HTTP/1.1." No authentication is mentioned. Query spec [P] https://google.github.io/osv.dev/post-v1-query/ : body `commit` ("If specified, `version` should not be set"), `version`, `package` — "Package Objects can be described by package name AND ecosystem OR by the package URL"; "pagination will occur when there more than 1,000 vulnerabilities in the response" or the query exceeds 20 s, continued via `page_token` = `next_page_token`; "API requests are case-sensitive" (use `PyPI`, not `pypi`). [observed] `{"package":{"name":"jinja2","ecosystem":"PyPI"},"version":"2.4.1"}` → 200 with `vulns[]`; the same query **without** `version` also returns 200 with all known vulns for the package (behaviour not documented on the page).

**deps.dev** [P] https://docs.deps.dev/api/v3/ and https://docs.deps.dev/api/v3alpha/ : base `https://api.deps.dev/v3`; no key; "Use of the deps.dev API is subject to the Google API Terms of Service"; "Clients are expressly permitted to cache data served by the API"; no rate limit stated; generated data CC-BY 4.0. Systems: `GO`, `RUBYGEMS`, `NPM`, `CARGO`, `MAVEN`, `PYPI`, `NUGET` (Go covers only modules seen by proxy.golang.org). `GetProject` (`GET /v3/projects/github.com%2Fowner%2Frepo`) returns `starsCount`, `forksCount` ("Only available for GitHub and GitLab"), `openIssuesCount`, `license` ("The license reported by the project host"), `scorecard` (OpenSSF, `overallScore` 0–10 + checks) — [observed] expressjs/express: stars 69,575, forks 25,253, open issues 244, MIT, scorecard dated 2026-08-24. `GetDependents` (v3alpha, `GET /v3alpha/systems/{sys}/packages/{name}/versions/{ver}:dependents`): "GetDependents returns information about the number of distinct packages known to depend on the given package version." "Dependent counts are currently available for npm, Cargo, Maven and PyPI." Fields `dependentCount`, `directDependentCount`, `indirectDependentCount` — [observed] npm/express/4.18.2 → 2,184 / 861 / 1,373. Batch endpoints: "Up to 5000 requests are allowed in a single batch."

**ecosyste.ms** — OpenAPI specs define no security scheme [P] https://packages.ecosyste.ms/docs/api/v1/openapi.yaml (server `https://packages.ecosyste.ms/api/v1`) and [P] https://repos.ecosyste.ms/docs/api/v1/openapi.yaml (server `https://repos.ecosyste.ms/api/v1`). Packages: `GET /packages/lookup` — "lookup a single package by repository URL, purl or ecosystem+name"; `Package` fields `downloads`, `downloads_period`, `dependent_packages_count`, `dependent_repos_count`, `repository_url`, `latest_release_number`, `latest_release_published_at`, `maintainers`, `advisories`, `rankings` (the spec gives no description for `downloads`/`downloads_period`). Gotcha [observed]: `lookup?repository_url=https://github.com/expressjs/express` returned `star-express` (another npm package pointing at the same repo) first — filter by `name`. Repos: `GET /repositories/lookup?url=` or `/hosts/{host}/repositories/{owner/repo}`; `Repository` fields include `stargazers_count`, `forks_count`, `subscribers_count`, `open_issues_count`, `archived`, `fork`, `license`, `pushed_at`, `last_synced_at`, `topics`, `metadata` and `commit_stats` (both untyped objects in the spec). [observed] for expressjs/express `metadata.files` keys include `readme`, `license`, `code_of_conduct`, `contributing`, `security`, `funding`, `changelog`, `governance`, `codeowners`, `citation`; `commit_stats` = `{total_commits, total_committers, mean_commits, dds, last_synced_commit}`; `last_synced_at` 2026-10-21. Limits: pricing page [P] https://ecosyste.ms/pricing — Free "300 requests per hour", access "Common or polite pools", no key; Develop "$200 per month", "1,000 requests per hour", API key; Scale "$1000 per month", "5,000 requests per hour". Blog 2025-09-01 [P] https://blog.ecosyste.ms/2025/09/01/rate-limiting-the-right-way.html — serving "500GB of data and 40m requests daily"; get into the polite pool by adding a `mailto` parameter or "`mailto:you@exmaple.com` somewhere in your user-agent request header"; the common pool "will have less consistent response times". [observed] with a User-Agent containing an email: packages API `x-ratelimit-limit: 15000`, `x-ratelimit-tier: polite`, `x-ratelimit-consumer: anonymous`; repos API `x-ratelimit-limit: 5000`, `x-ratelimit-tier: anonymous` — i.e. live limits are far above the pricing page's 300/h; treat the pricing page as the policy and the headers as today's behaviour. Data licence "Code: AGPL-3 — Data: CC BY-SA 4.0" [P] https://ecosyste.ms/ (share-alike applies to published derived data).

**npm** [P] https://github.com/npm/registry/blob/main/docs/download-counts.md : base `https://api.npmjs.org/downloads`; `GET /downloads/point/{period}[/{package}]` (total) and `/downloads/range/{period}[/{package}]` (per day); periods `last-day`, `last-week`, `last-month`, `last-year`, or `YYYY-MM-DD:YYYY-MM-DD`; "All other queries are limited to at most 18 months of data"; bulk = comma-separated up to 128 packages and 365 days, "Scoped packages are not yet supported in bulk queries"; per-version counts (`/versions/{pkg}/last-week`) cover the previous 7 days only. No auth and no rate limit are mentioned. [observed] `point/last-week/express` → `{"downloads":87449208,"start":"2026-10-12","end":"2026-10-18"}`, no rate-limit headers. What counts as a download: the doc defers to an npm blog post that now returns 404 — **UNSOURCED**.

**PyPI** — no official download-count API. PyPA guide [P] https://packaging.python.org/en/latest/guides/analyzing-pypi-package-downloads/ : "PyPI does not display download statistics for a number of reasons"; logs stream to BigQuery table `bigquery-public-data.pypi.file_downloads`; "You can run up to 1TB of queries per month" free; caveats "Internal or unofficial mirrors (can both raise or lower download counts)", "pip's download cache (lowers download counts)", "Unofficial scripts or attempts at download count inflation (raises download counts)". pypistats.org [P] https://pypistats.org/api/ and https://pypistats.org/about : "hosted and operated by The Python Software Foundation" (created by Christopher Flynn); endpoints `recent` (`last_day`, `last_week`, `last_month`), `overall` (with/without mirrors), `python_major`, `python_minor`, `system`; "Time series data is retained only for 180 days"; "All download stats exclude known mirrors (such as bandersnatch) unless noted otherwise"; "All download data is updated once daily"; "IP-based rate limiting is imposed application-wide" (no number) — [observed] `x-ratelimit-limit: 60` with a reset about one minute out.

**crates.io** — OpenAPI [P] https://crates.io/api/openapi.json : "Experimental API documentation"; "Before using this API, please read the crates.io data access policy and ensure that your use of the API complies with the policy."; "Most parts of the API do not require authentication." `GET /api/v1/crates/{name}` → `downloads` "The total number of downloads for this crate." and `recent_downloads` "The total number of downloads for this crate in the last 90 days."; `GET /api/v1/crates/{name}/downloads` — "This includes the per-day downloads for the last 90 days and for the latest 5 versions plus the sum of the rest."; `GET /api/v1/crates/{name}/reverse_dependencies`. [observed] serde: 1,569,459,593 total, 374,786,156 recent. The live policy page https://crates.io/data-access is a JavaScript app (fetch returns only the title), so the wording below is from the RFC that defined it [P] https://rust-lang.github.io/rfcs/3463-crates-io-policy-update.html (RFC 3463, start date 2023-07-24): API users must stay at or below one request per second; the User-Agent must uniquely identify the application and should include contact information, e.g. "my_bot (help@my_bot.com)"; "We reserve the right to block traffic from any client that we determine to be in violation of this policy"; a database dump is published every 24 hours at `https://static.crates.io/db-dump.tar.gz`, with no stability guarantee on table layouts.

**GitHub "dependents"** — HTML only. [P] https://docs.github.com/en/code-security/supply-chain-security/understanding-your-software-supply-chain/exploring-the-dependencies-of-a-repository : "For public repositories, the dependents view shows how the repository is used by other repositories"; "The dependent counts are approximate and may not always match the dependents listed"; GitHub "currently only determines dependents for public repositories"; reached via Insights → Dependency graph → Dependents. The REST Dependency graph category lists only dependency review, dependency submission and SBOM export [P] https://docs.github.com/en/rest/dependency-graph ; the GraphQL `Repository` reference has no field containing "dependents" [P, checked 2026-10-22] https://docs.github.com/en/graphql/reference/repos . Substitutes: deps.dev `:dependents` (package-level, four ecosystems), ecosyste.ms `dependent_repos_count` / `GET /usage/{ecosystem}/{package}/dependent_repositories`.

**Community profile** [P] https://docs.github.com/en/rest/metrics/community : `GET /repos/{owner}/{repo}/community/profile` — "Returns all community profile metrics for a repository. The repository cannot be a fork."; `health_percentage` is "a percentage of how many of the recommended community health files are present"; `files` = `code_of_conduct`, `code_of_conduct_file`, `license`, `contributing`, `readme`, `issue_template`, `pull_request_template` — **no security policy**. [observed] golang/go unauthenticated → 200, `health_percentage: 87`.

**Security policy detection** — GraphQL `Repository.isSecurityPolicyEnabled` ("Returns true if this repository has a security policy.") and `securityPolicyUrl` ("The security policy URL."); same object exposes `fundingLinks`, `hasVulnerabilityAlertsEnabled`, `codeOfConduct`, `licenseInfo`, `latestRelease`, `stargazerCount`, `vulnerabilityAlerts`, `dependencyGraphManifests` [P] https://docs.github.com/en/graphql/reference/repos . File-location rules if you check via the contents API instead: "Store all other supported files in the root of the repository, the `.github` folder, or the `docs` folder"; org-level defaults from a `.github` repository apply only when a repo "does not contain its own file of that type" [P] https://docs.github.com/en/communities/setting-up-your-project-for-healthy-contributions/creating-a-default-community-health-file . Published repository advisories: `GET /repos/{owner}/{repo}/security-advisories` — "Anyone can access any published security advisory on a public repository." [P] https://docs.github.com/en/rest/security-advisories/repository-advisories

**Releases** [P] https://docs.github.com/en/rest/releases/releases : `GET /repos/{owner}/{repo}/releases` — "Information about published releases are available to everyone. Only users with push access will receive listings for draft releases."; "This returns a list of releases, which does not include regular Git tags that have not been associated with a release."; `per_page` max 100. `GET .../releases/latest` — "The latest release is the most recent non-prerelease, non-draft release, sorted by the created_at attribute. The created_at attribute is the date of the commit used for the release, and not the date when the release was drafted or published." Fields: `tag_name`, `published_at` (null when unpublished), `prerelease`, `draft`, `assets[].download_count`.

**Statistics** [P] https://docs.github.com/en/rest/metrics/statistics : "Computing repository statistics is an expensive operation, so we try to return cached data whenever possible." If not cached the API returns **202** and starts a background job; "You should allow the job a short time to complete, and then submit the request again." — [observed] first `GET /repos/hashicorp/terraform/stats/contributors` → 202. Contributors: "This endpoint will return 0 values for all addition and deletion counts in repositories with 10,000 or more commits." (`code_frequency` returns 422 for such repos). Response per contributor: `total`, `weeks[]` of `{w, a, d, c}` (week start, additions, deletions, commits). Other endpoints: `commit_activity`, `code_frequency`, `participation`, `punch_card`.

**Get a repository** [P] https://docs.github.com/en/rest/repos/repos#get-a-repository : `stargazers_count`, `watchers_count`, `subscribers_count`, `forks_count`, `open_issues_count`, `pushed_at`, `updated_at`, `archived`, `disabled`, `license`, `topics`, `default_branch` (the starring page notes `watchers`/`watchers_count`/`stargazers_count` all mean stars; `subscribers_count` is watchers).

---

## 4. Cloudflare limits that bound a scheduled fetch job

**Workers plan limits** [P] https://developers.cloudflare.com/workers/platform/limits/ (last updated 2026-10-12)

| Row (verbatim) | Free | Paid |
|---|---|---|
| Requests | 100,000/day | No limit |
| CPU time | 10 ms | 5 min |
| CPU time per Cron Trigger | 10 ms | 30 seconds (< 1 hour interval) / 15 min (>= 1 hour interval) |
| Memory | 128 MB | 128 MB |
| Subrequests | 50/request | 10,000/request |
| Simultaneous outgoing connections/request | 6 | 6 |
| Number of Cron Triggers per account | 5 | 250 |
| Number of Static Asset files per Worker version | 20,000 | 100,000 |
| Individual Static Asset file size | 25 MiB | 25 MiB |

Also: "Waiting on network requests (such as `fetch()` calls, KV reads, or database queries) does **not** count toward CPU time."; Duration table "Cron Trigger | 15 min" (wall clock); "Accounts on the Workers Free plan have a daily request limit of 100,000 requests, resetting at midnight UTC. When a Worker exceeds this limit, Cloudflare returns **Error 1027**."

**Cron Triggers** [P] https://developers.cloudflare.com/workers/configuration/cron-triggers/ (last updated 2026-10-17): "Cron Triggers allow users to map a cron expression to a Worker using a `scheduled()` handler"; "Cron Triggers execute on UTC time."; the Minute field accepts `0-59` with `* , - /`, and `* * * * *` ("At every minute") is a listed example — so the minimum interval is one minute; changes "may take several minutes (up to 15 minutes) to propagate"; "Cron Events stores the 100 most recent invocations"; local test by sending a request to `/cdn-cgi/local/scheduled` under `wrangler dev` (Local Explorer needs Wrangler 4.139.0+). Its Limits section only says: "Refer to Workers limits for the maximum number of Cron Triggers per account."

**Pricing** [P] https://developers.cloudflare.com/workers/platform/pricing/ (last updated 2026-10-02): Free — "100,000 per day" requests, "10 milliseconds of CPU time per invocation"; Paid — "a minimum charge of $5 USD per month", "10 million included per month, +$0.30 per additional million" requests, "30 million CPU milliseconds included per month, +$0.02 per additional million CPU milliseconds", "Max of 15 minutes of CPU time per Cron Trigger or Queue Consumer invocation"; footnote 1 "Inbound requests to your Worker. Cloudflare does not bill for subrequests"; footnote 3 "Requests to static assets are free and unlimited." Example 3 bills an hourly Cron Worker as "720 requests/month", so scheduled invocations count as requests [inferred from the example; not stated as a rule].

**Workers KV** — limits [P] https://developers.cloudflare.com/kv/platform/limits/ (2026-10-08): Reads "100,000 reads per day" Free / Unlimited Paid; Writes to different keys "1,000 writes per day" / Unlimited; Writes to same key "1 per second" (both); Operations/Worker invocation 1000 (both; "A bulk request to Workers KV counts for 1 request"); Storage/account "1 GB" / Unlimited; Key size 512 bytes; Value size 25 MiB; Minimum `cacheTtl` 30 seconds. Pricing [P] https://developers.cloudflare.com/kv/platform/pricing/ (2026-04-21): Keys read "100,000 / day" vs "10 million/month, + $0.50/million"; Keys written / deleted / List requests "1,000 / day" vs "1 million/month, + $5.00/million"; Stored data "1 GB" vs "1 GB, + $0.50/ GB-month"; "All limits reset daily at 00:00 UTC. If you exceed any one of these limits, further operations of that type will fail with an error." Consistency [P] https://developers.cloudflare.com/kv/concepts/how-kv-works/ : "Changes may take up to 60 seconds or more to be visible in other global network locations."

**D1** — limits [P] https://developers.cloudflare.com/d1/platform/limits/ (2026-04-21): Databases per account "50,000 (Workers Paid) / 10 (Free)"; Maximum database size "10 GB (Workers Paid) / 500 MB (Free)"; Maximum storage per account "1 TB (Workers Paid) / 5 GB (Free)"; Time Travel "30 days (Workers Paid) / 7 days (Free)"; Queries per Worker invocation "1000 (Workers Paid) / 50 (Free)"; max statement 100 KB; max row/BLOB 2 MB; max query duration 30 seconds. Pricing [P] https://developers.cloudflare.com/d1/platform/pricing/ (2026-04-21): Rows read "5 million / day" Free vs "First 25 billion / month included + $0.001 / million rows"; Rows written "100,000 / day" vs "First 50 million / month included + $1.00 / million rows"; Storage "5 GB (total)" vs "First 5 GB included + $0.75 / GB-mo"; "Rows read measure how many rows a query reads (scans), regardless of the size of each row" (a full-table scan of 5,000 rows = 5,000 reads); "Free limits reset daily at 00:00 UTC."

**Can a Cron Trigger live in the same Worker as Static Assets?** No single Cloudflare sentence says so; the evidence is:
- Compatibility matrix [P] https://developers.cloudflare.com/workers/static-assets/compatibility-matrix/ (2026-10-12): "Cron Triggers | ✅ | ❌" (Workers | Pages); KV, D1, R2 ✅ on both; Workers Logs ✅ Workers only.
- Wrangler configuration [P] https://developers.cloudflare.com/workers/wrangler/configuration/ : `triggers.crons` ("An array of `cron` expressions"; "Commenting out the `crons` key will not disable a Cron Trigger") and `assets` (`directory`, `binding`, `run_worker_first`, `not_found_handling`; "You can only configure one collection of assets per Worker") are both top-level keys of the same Worker configuration, alongside `kv_namespaces` and `d1_databases`; nothing on the page excludes combining them.
- Billing [P] https://developers.cloudflare.com/workers/static-assets/billing-and-limitations/ (2026-04-23): "Requests to static assets are free and unlimited."; requests that invoke the script (including `run_worker_first` matches) are billed as Workers requests, and on the Free tier over-limit requests "will receive a 429 (Too Many Requests) response."
- Conclusion **[inferred]**: one Worker with `assets` + a `fetch` handler (the existing 410 router) + a `scheduled` handler + `triggers.crons` is supported; verify locally with `wrangler dev` and a request to `/cdn-cgi/local/scheduled` before relying on it.

**Free-plan envelope for the job [inferred from the rows above]**: ≤5 cron schedules per account; one hourly schedule = 24 invocations/day, each ≤50 subrequests and ≤10 ms CPU (JSON parsing of large GitHub/ecosyste.ms responses can exceed 10 ms; network wait does not count) → ≤1,200 upstream fetches/day; KV allows only 1,000 writes/day (≈41/hour) and D1 100,000 rows written/day; GitHub's 5,000/h PAT budget is not the binding constraint, Cloudflare Free's CPU and subrequest caps are. Workers Paid ($5/mo) raises CPU to 30 s (<1 h cadence) or 15 min (≥1 h cadence), subrequests to 10,000/invocation, and removes the KV/D1 daily caps in favour of metered usage. Because the site is built by Go in GitHub Actions, the fetch step can alternatively run there and commit/upload a data file, side-stepping Worker CPU limits entirely — out of scope here, noted as an option.

---

## Open items / unsourced

- What fraction of *search* traffic for repo-name queries comes from tier-1 countries — UNSOURCED; only Search Console on live pages can answer it.
- How github.com/trending ranks repositories — GitHub publishes no method (the docs say only "popular repositories of the day").
- The full text of Martin Woodward's follow-up "Keep trending projects" plan — only available as a user's quotation in the thread [S].
- npm's definition of a "download" — the registry doc's blog link is dead (404) — UNSOURCED.
- Whether `stargazers/history` and `stargazers/count` really serve unauthenticated callers — stated in the docs, not probed (quota exhausted).
- ecosyste.ms: pricing page says 300 req/h free, live headers show 5,000–15,000/h; which one is enforced is unknown.

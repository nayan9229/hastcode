# Author runbook: running HastCode day to day

Plain words, copy-paste commands. For the Author (you). The long technical
version is `README.md`; the rules behind it are in `docs/adr/`.

## 1. The whole thing in one paragraph

Robots watch a list of websites (`sources.yaml`) and collect new articles
and incidents as **Leads**. A Claude session turns the best Leads into
**drafts**. You check every claim in a draft against its Source, mark it
**verified**, and **publish**. At most **two pieces a day** go live. Nothing
goes live without your check. That is the entire system.

## 2. The five kinds of piece

| Kind | What it is | Length | Lives in | Special rule |
|---|---|---|---|---|
| Brief | Short note on something that just happened, with "why it matters" | 250 words max | `content/briefs/` | Lead must be at most 3 days old when published. Must link a live Explainer. |
| Case Study | How a named company built or fixed one thing | ~2,400 words | `content/explainers/<vertical>/` | |
| Concept | One idea explained | ~2,400 words | `content/explainers/<vertical>/` | |
| Comparison | Two or more ways to do one thing, compared | ~2,400 words | `content/explainers/<vertical>/` | |
| Repo Profile | A GitHub project, with a fetched Fact Pack | | `content/explainers/<vertical>/` | Facts come from `make fetch`, never typed by hand. |

Case Study, Concept, Comparison and Repo Profile together are called
**Explainers**. A Brief is not an Explainer.

The six **Verticals** (sections) are: `distributed-systems`,
`databases-storage`, `apis-networking`, `infrastructure-cloud`,
`ai-ml-systems`, `languages-runtimes`.

## 3. One-time setup

You need these installed and working once:

- Go and Node 22+. Then run `npm ci` in the repo.
- Claude Code command line (`claude`), logged in. `make daily` uses it.
- A `.env` file in the repo root with `GITHUB_TOKEN=...` (GitHub
  personal token, read-only is enough). `make intake`, `make fetch` and
  `make refresh` read it. `scripts/setup-deploy.sh` can write it for you.
- `gh` (GitHub command line), logged in. Only used to watch deploys.

Quick health check: `make check` must end with `ok` lines and no `FAIL`.

Getting found by search engines (one time, by hand, about 40 minutes):
follow "Runbook: get indexed, in order" in
`docs/research/indexing-submission-channels.md`. In short: let
PerplexityBot through in Cloudflare, turn on Always Use HTTPS, submit both
sitemaps in Google Search Console, import the site into Bing Webmaster
Tools. IndexNow already pings on every publish; don't send it by hand.

Optional: also tell Google's WebSub hub when a feed changes. It is off
until you turn it on. First, check once that the hub still answers (this
sends it one notice about `/feed.xml`; nothing else):

```sh
curl -sS -o /dev/null -w '%{http_code}\n' -d hub.mode=publish -d hub.url=https://hastcode.com/feed.xml https://pubsubhubbub.appspot.com/
```

It should print `204`. Write that and today's date into README "WebSub
hub ping" (the "Live check" line). If it prints anything else, stop here
and leave the ping off. On `204`, turn it on:

```sh
gh variable set WEBSUB_ENABLED --body true
```

After a `scripts/publish.sh cf`, check what the hub answered (`answered
204` is good; "no feed changed" means nothing to send):

```sh
gh run view "$(gh run list --workflow publish.yml --limit 20 --json databaseId,displayTitle --jq '[.[] | select(.displayTitle | test("^Publish PUBCF-[0-9]+-[0-9]+$"))][0].databaseId')" --log | grep -E 'changed feed: |answered|no artefact: no feed'
```

Weekly check (5 minutes, read-only, changes nothing). Run it on your
laptop; nothing in GitHub runs it:

```sh
make crawl-check
```

It prints one `PASS` or `FAIL` line per check and a summary at the end.
All `PASS`: nothing to do. Each `FAIL` line says what is wrong. The steps
below are in "Runbook: get indexed, in order" in
`docs/research/indexing-submission-channels.md`:

- "Cloudflare is blocking PerplexityBot" (or another crawler): step 1,
  set that crawler to Allow in AI Crawl Control.
- "answers 200 instead of redirecting": step 2, turn on Always Use HTTPS.
- "robots.txt stops ... crawling /": if the line it names is in
  Cloudflare's managed block (the top of the file), step 1 item 5;
  otherwise ask Claude, it is in our build.
- "has lost our ... allow group" or "has no Sitemap: line": something in
  front of the site rewrote robots.txt; ask Claude to look.
- "https://www.hastcode.com/ answers ...": the www-to-hastcode.com
  redirect rule in Cloudflare is gone; put it back.
- "could not reach": your internet or the site is down; run it again.

It sends crawler names from your own IP address, so it shows what
Cloudflare does with each name, not what the real crawler sees. Then look
at Search Console **Pages** and **Performance**, and Bing **IndexNow
Insights**.

## 4. A normal morning

### Step 1: draft today's pieces

```sh
make daily
```

What it does, in order (it prints `1/5` to `5/5`):

1. Fetches every Source and lists new Leads in `.intake/queue.md`.
2. A Claude session chooses: up to 2 Briefs, and up to 2 Explainers
   (but 0 Explainers while 4 or more Explainer drafts are waiting, see
   section 6).
3. Starts a skeleton file for each pick.
4. One Claude session writes each draft, a second one fact-checks it.
5. Commits the drafts with your check list at
   `.scratch/daily/<today>/review-queue.md`.

It never publishes, never pushes, never marks anything verified. It
refuses to start if you have uncommitted changes, so commit or stash
first. It prints its Claude cost at the end.

Want to see the picks without writing anything? `make daily DRY=1`.

### Step 2: verify

Open `.scratch/daily/<today>/review-queue.md`. For each draft it lists:

- the file,
- a tick list "Drafter asks you to check",
- a list "Fact-checker could not confirm" (check these hardest),
- a tick list "Self-edit before you verify": the
  [self-edit checklist](editorial/self-edit-checklist.md) for that draft
  (a Brief gets only the lines that fit a Brief).

First, read the draft once and tick the self-edit boxes as you go. They
check the writing (is the point clear, do the headings tell the story,
are the sentences short), not the facts. Fix the writing in the draft
yourself where a line fails. Ticking them is not verifying, and it never
stands in for `verified:`.

Then open the draft and its Source side by side and check every line of
the other two lists against the Source. Fix wording in the draft yourself
if needed. If a claim is wrong and cannot be fixed, delete the draft file
(`git rm`).

To see a draft as a web page: `make dev`, then open
http://localhost:3000. Drafts show in this local preview.

### Step 3: release

Once every draft in the check list is verified:

```sh
make release                 # every draft; or FILES="content/briefs/a.md content/explainers/v/b.md"
git push origin main
scripts/publish.sh cf        # go live on hastcode.com
```

`make release` does, for each draft: removes `draft: true`, sets
`verified:` to today, sets `reviewed:` to now on Explainers, keeps
`published:` as the draft has it, runs the production build, and commits
with a "Release <day>: ..." message. Only you run it; an agent never does.
`make release DRY=1` shows what it would do.

It refuses, before changing anything, when a day would end up with more
than two pieces live. Then release fewer with `FILES=`, or delete a stale
draft (`git rm`).

`scripts/publish.sh vl` first sends the same commit to the QA copy
(https://qa.hastcode.com/) if you want to look before going live.

Watch the deploy: `gh run list --workflow=publish.yml --limit 2`. When it
says `success`, the page is at `https://hastcode.com/briefs/<slug>/` or
`https://hastcode.com/<vertical>/<slug>/`.

### If `make build` fails

| Error says | Meaning | Fix |
|---|---|---|
| `unverified` | A non-draft piece has no `verified:` line, or the date is before `published`, or in the future | Add or fix `verified:` |
| `publish-pace` | More than 2 pieces share one `published` day | Put `draft: true` back on the extras, or give them another day |
| `html-size` | Page too big (limit about 51 KB) | Cut words or a diagram |
| `broken-link` | A link in the piece is dead | Fix or remove it |
| `brief-length` | Brief over 250 words | Cut words |
| `brief-explainer-link` | Brief does not link a live Explainer | Link one in "Why it matters" and in `explainers:` |
| `ai-note` | A piece or page says "AI assistance" or "drafted with AI" | Delete that line or those words. No AI note is ever published |
| `review-overdue` | An Explainer's `reviewed:` date is old | Re-read it, update `reviewed:` |

**Only want the errors?** Every build also prints hundreds of writing
warnings (below), and the errors come first, so they scroll off the top.
To see just the errors and the last line:

```sh
make build 2>&1 | grep -E '^(error|build failed)'
```

### Writing warnings (they never block anything)

The build points at three writing habits. They are warnings: they never
fail `make build`, never stop `make release`, and never block a Publish.
Fix them when you edit; nothing already live has to change.

| Warning says | Meaning | Fix |
|---|---|---|
| `sentence-length` | A sentence over 30 words. The message names the section, the first six words and the count | Split the sentence. A deliberate list may stay long. If it quotes a Source, split our own words around the quote; never change the quote |
| `banned-word` | A hype or filler word from the banned list in voice rule 5 (`docs/editorial/house-byline-voice.md`) | Cut the word. Words inside quotation marks never warn. A product's name is not hype (Dropbox's "Magic Pocket"): leave that warning |
| `question-heading` | A `##` heading that ends in "?" | Restate it as the answer ("Why do terms only go up?" → "Terms only go up"). FAQ questions in a Repo Profile are fine |

Code, tables, quotes in `>` boxes and link addresses are never checked.

### Missed a day? `make daily` fills it

If you skip a morning, the next `make daily` looks back up to 3 days for
days with fewer than two pieces (a **Gap Day**) and drafts Explainers for
them, oldest day first, two per day, up to 6 drafts in one run. Such a
draft already carries that past day in `published:`; the check list says
"Fills Gap Day 2026-11-13". `make release` keeps that day. Gap Days never
get Briefs: a Brief is news and goes stale. (ADR-0012)

## 5. Why did `make daily` make only Briefs today?

Because **four Explainer drafts were already waiting**. The rule: while
4 or more Explainer drafts sit unpublished, `make daily` drafts 0 new
Explainers (Briefs still, because they expire in 72 hours). This stops
the backlog growing. Clear it by releasing two Explainers a day.

## 6. How to get the other kinds

**Option A, let `make daily` draft Explainers anyway:**

```sh
make daily EXPLAINERS=2            # force up to 2 Explainers today
make daily BRIEFS=0 EXPLAINERS=2   # Explainers only, no Briefs
make daily BACKLOG=8               # raise the waiting limit from 4 to 8
```

**Option B, pick a specific Lead yourself:**

```sh
make intake                 # refresh the Lead list, no AI, cheap
open .intake/queue.md       # find the Lead's 8-letter ID
make lead ID=abcd1234 KIND=case-study VERTICAL=databases-storage
```

`KIND` is `brief`, `case-study`, `concept` or `comparison`. A Brief also
needs `EXPLAINER=<slug>` naming the live Explainer it links. This writes a
skeleton file with empty headings. Then open Claude Code in the repo and
say: "Draft `content/explainers/<vertical>/<slug>.md` following
`docs/editorial/<kind>-writing-guide.md` and `docs/editorial/house-byline-voice.md`,
self-edit it with `docs/editorial/self-edit-checklist.md`,
then fact-check every claim against its Source. Keep `draft: true`."
Then go through the same checklist yourself while you read the draft,
before you check its claims.
Check the page size with
`scripts/daily/check-page.sh content/explainers/<vertical>/<slug>.md /tmp/check`.

Hide a Lead you never want: `make lead-dismiss ID=abcd1234`. Leads
vanish by themselves after 7 days.

**Option C, write from scratch, no Lead:** copy the template
`docs/editorial/<kind>-template.md` into the right folder, keep
`draft: true`, follow the writing guide, then self-edit with
`docs/editorial/self-edit-checklist.md` before you check the claims.

**Repo Profiles** are different: `make fetch REPO=owner/name
MD=content/explainers/<vertical>/<slug>.md` fetches the Fact Pack from
GitHub. `make refresh` re-fetches all of them; it exits with `Error 3`
when anything changed, which is normal: review the diff, commit, push.
If you ran the quickstart, add one line right after its command block
saying when, on what machine, and whether it worked as written, e.g. "We
ran the quickstart on 2026-11-17 on macOS 26, at v0.5.0: it worked as
written" (leave the commands as they are, and skip the line if you did
not run them).

## 7. Rules to remember

- Two pieces a day, counted by `published` day. The build enforces it.
- `verified` is yours alone (`make release` writes it when you run it).
  Same day as `published` or later, never before, never in the future.
- A Brief must be published within 72 hours of its Lead's date. A Brief
  draft older than that cannot be published; delete it or turn the topic
  into an Explainer later.
- Every claim links its Source. Quotes short and credited. Diagrams
  redrawn in Mermaid, never copied.
- No AI note. Nothing a reader sees says a piece was drafted with AI
  (ADR-0013). The build stops on it (`ai-note`). Writing about AI as a
  topic is fine.
- Nothing from GitHub Actions ever fetches or drafts. It only tests,
  builds previews and deploys a tag.

## 8. Folders you will touch

| Path | What is in it |
|---|---|
| `content/briefs/`, `content/explainers/<vertical>/` | The pieces |
| `.scratch/daily/<date>/review-queue.md` | Your check list for that day's drafts |
| `.intake/queue.md` | Current Leads with IDs (local, not committed) |
| `.intake/daily/<run>/` | Logs and cost of each `make daily` run |
| `docs/editorial/` | Templates and writing guides per kind, and the self-edit checklist |
| `sources.yaml` | The watched websites; add one with a reason in the commit |

## 9. Words

- **Lead**: a new article or incident found by Intake, not yet a piece.
- **Piece**: one page on the site, any kind.
- **Explainer**: Case Study, Concept, Comparison or Repo Profile.
- **Vertical**: one of the six sections.
- **Source**: the first-party page a claim comes from.
- **Fact Pack**: the fetched JSON beside a Repo Profile.
- **Draft**: `draft: true` in the file; shows locally and in previews, never in production.
- **Verified**: `verified: YYYY-MM-DD`; your record that you checked every claim.
- **Publish**: a git tag made by `scripts/publish.sh`; `cf` is production, `vl` is QA.

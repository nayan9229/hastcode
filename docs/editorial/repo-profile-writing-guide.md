# Writing a Repo Profile

A Repo Profile (`GLOSSARY.md`, ADR-0005) tells the Reader, a mid-to-senior backend engineer, what one open-source repository is, how it is built and whether to adopt it. The build renders the facts (header card, Signals, Changelog, hero card) from the Fact Pack. You write the opinion. Choose the repository with [the selection guide](repo-profile-selection.md), write in [the house byline voice](house-byline-voice.md), and hold every sentence to the [Editorial Policy](../../content/pages/editorial-policy.md).

## The flow

1. **Choose.** The repository passes every floor in the selection guide.
2. **Copy** [the template](repo-profile-template.md) to `content/explainers/<vertical>/<slug>.md`, on a branch. The slug is the repository name in lower case (`apache-iggy` when the name alone is ambiguous). Keep `draft: true`.
3. **Fetch.** `make fetch REPO=owner/name MD=content/explainers/<vertical>/<slug>.md` (needs `GITHUB_TOKEN`) writes the Fact Pack beside the Markdown and the README text to `node_modules/.cache/hastcode/readmes/<owner>/<name>/`. A `note: … NOASSERTION` line means GitHub could not classify the license: read the LICENSE file before you go on. Commit the Fact Pack; never edit it by hand.
4. **Draft** against the Fact Pack, the README and the release pages. AI assistance is fine; it is a draft, not a source.
5. **Preview.** `make dev`, then push the branch: the Preview workflow posts a preview URL in its run summary. Check the header card, Signals, Changelog and the hero card read right.
6. **The Author verifies** every claim, number, command and diagram against its source, and rewrites what is wrong. Nothing publishes until the Author has done this.
7. **Publish.** Remove `draft: true`, merge, then tag (`scripts/publish.sh vl`, check QA, `scripts/publish.sh cf`).
8. **Day 28.** Fill in a [checkpoint note](checkpoint-note-template.md).

After that, the Monday Refresh keeps the facts current through a pull request a human merges, and the 90-day Review re-reads the prose. Neither rewrites the other.

## Ground every claim

- **Three sources, in this order:** the Fact Pack, the README, the release pages (then the repository's own docs and code). No star trackers, no other outlets' summaries, no "people say".
- **Link inline** where you make the claim. Sources lists the README and every release page or doc you cited; the build adds the Fact Pack's API URLs itself, once each.
- **Do not retype what the Fact Pack carries.** Stars, license, language, release dates and contributor counts render on the card and in Signals, and a Refresh keeps them current. A number in your prose goes stale on the next Refresh. Point at it ("the release cadence in Signals") or say what it means ("releases every few weeks since 2025").
- **Say "we infer"** when you reconstruct something the sources leave open.

## Section by section

The build fails unless the nine H2s appear exactly as in the template, in that order. Headings match byte for byte, with a straight apostrophe in "Don't".

- **Opening (before TL;DR).** One hook sentence: a sourced fact that makes the Reader look twice. Then the AI-assistance note as a blockquote, exactly as the template has it.
- **TL;DR.** Three to five bullets. Answer first: what it is, what sets it apart, who should adopt it.
- **In 60 seconds.** The README's install or quickstart commands, **verbatim**: same commands, same order, same flags, in a fenced block, with a link to the README section. Never improve, shorten or merge them. If the README has none, say so and copy the shortest install the official docs give, linked.
- **Architecture.** How it is built: the components, how data moves, where state lives. The first Mermaid fence is the diagram slot; walk it in order and link the code or docs behind each claim. *Signals renders after this section: do not write it.*
- **Use it if / Don't use it if.** Two to four bullets each. Take a position and ground it in a fact from above. Each "Don't" says what to use instead.
- **Alternatives.** Three to five, one line each on how it differs. Link our Comparison or Profile where one exists, as a full `/{vertical}/{slug}/` path; otherwise the alternative's own repository or site. *Changelog renders after this section: do not write it.*
- **FAQ.** Five to eight `###` questions phrased the way the Reader types them into a search box: "Is {repo} production ready?", "Does {repo} support {feature}?", "What are the alternatives to {repo}?". Each answer is one to three sentences that stand alone: name the repository, never write "as above". Answers become the page's FAQPage data and get quoted out of context.
- **Verdict.** One paragraph: the recommendation, who it is for, and what would change it.
- **Sources.** Last, with at least one link.

## Titles, descriptions and search intent

| Search | Where it goes |
|---|---|
| `{repo}`, "what is {repo}" | Profile title `{repo}: {what it is}`, e.g. "Apache Iggy: a persistent message stream in Rust" |
| "is {repo} production ready" | The first FAQ question, answered plainly |
| "{repo} alternatives" | The Alternatives section, plus an FAQ question "What are the alternatives to {repo}?" |
| "{repo} vs {alternative}" | An FAQ question here; when the pair earns a full piece, a Comparison Explainer titled "{repo} vs {alternative}: …", linked from Alternatives |

The description is one sentence of 150 characters or fewer that names the repository, says what it is and states the trade-off that matters. It is the card's "what it is" line and the search snippet. `who_for` is one line naming the Reader who should adopt it.

## What the build checks

`repo-profile-sections` (headings missing, extra or out of order), `repo-profile-fact-section` (a hand-written Signals or Changelog), `repo-profile-faq` (five to eight questions), `tldr`, `diagram`, `sources`, `fact-pack` and `fact-pack-repo` (Fact Pack missing, unparsable or for another repository), `author-placeholder` (a `TODO(author)` left in a production build). Warnings: `fact-pack-stale` (fetched over 30 days ago), `section-length` (an H2 over 250 words).

## The Friday roundup

The week's roundup is an ordinary Brief, and that Brief is the newsletter issue (sent by hand from Buttondown). "What happened" lists the week's repositories, one line each, linking each Profile by its bare slug (`[Apache Iggy](apache-iggy)`). "Why it matters" is the commentary. `source_url` is the trending page you read (e.g. `https://github.com/trending?since=weekly`), and `explainers` lists every Profile slug you linked. The Brief gates apply as usual: 250 words at most.

# The craft of article writing, from scratch, for hastcode.com

The Author asked: "What is the art of article writing?" This note answers it in plain words, from the people who studied or set the rules: usability researchers who watched people read, government plain-language teams, the Google and Microsoft style guides, newsroom standards, classic essays, explainer writers on their own method, and Google Search on what it rewards. It then applies all of that to our five kinds of Piece and compares it with our own editorial guides.

It does not repeat what is already researched elsewhere:

- What answer engines cite (page structure, section length, question headings, schema): [`aeo-citation-signals.md`](aeo-citation-signals.md).
- How much of a Source we may quote, and copyright: [`content-source-intake.md`](content-source-intake.md) §4.5.
- The house voice and the drafting pipeline: [`../editorial/house-byline-voice.md`](../editorial/house-byline-voice.md).

Researched **2026-11-16**. Labels: **[P]** primary (the owner's own page, fetched that day) / **[observed]** what I counted in our repo / **[inferred]** my deduction / **[opinion]** / **[secondary]** a source that reports someone else's work. Where a page could not be fetched, it says so, and Open questions lists it. The before/after examples are illustrations I wrote: their facts are made up unless they carry a link.

---

## TL;DR

- **The craft is mostly one habit: decide the single point, put it first, then make every later sentence earn its place.** Newsrooms call this the inverted pyramid. NN/g found it is "perfectly suited for the web" [P]. Orwell's version is four questions, starting with "What am I trying to say?" [P].
- **Readers scan; they do not read.** In NN/g's 1997 tests, "79 percent of our test users always scanned any new page", and "only 16 percent read word-by-word" [P]. Their 2020 review found the same: "People still primarily scan" [P]. GOV.UK says most people "only read 20 to 28% of text on a webpage" [P]. So headings, first sentences and the first two words of each carry most of the load.
- **Plain language works for experts too.** "The more educated the person and the more specialist their knowledge, the greater their preference for plain English" (GOV.UK) [P]. NN/g found the same with domain experts [P]. Our Reader is a senior engineer. That is a reason to write plainly, not an exception to it.
- **Concrete before abstract.** Start with a real case, then name the idea. This is what Julia Evans, Nicky Case and a US Department of Education practice guide all say, each from their own work [P]. Our "lead with a concrete hook" rule already does this.
- **Search engines reward the same craft, and punish its absence.** Google asks whether content provides "original information, reporting, research, or analysis". Google says it has no preferred word count: "(No, we don't.)". Its rater guidelines rate copied or paraphrased content Lowest "even if the page assigns credit" [P]. For AI-assisted drafting, Google says "manually fact-check and review all AI-generated content" [P].
- **We already do much of this well** (§5): a sourced hook, answer-first TL;DR, headings that state the point, every claim linked, short credited quotes, our own analysis, redrawn diagrams.
- **Top gaps** (§5.3): (1) our word-count "aim" reads like a target, which Google calls a warning sign; (2) "short sentences" has no number, and about a third of our sentences run past 25 words [observed]; (3) the "edit" step in the pipeline has no checklist (§4 gives one); (4) hooks often skip the "nut graf", the one sentence that says what this piece will give the Reader; (5) Repo Profiles show no first-hand use, which is the "Experience" in E-E-A-T.

---

## 1. What the craft is: 12 principles for a beginner

Each principle has: what to do, why (with the source), and a tiny before/after in our subject area.

### 1.1 Know the one point and the one Reader before you write

**Do:** Write one sentence that says what the piece will tell the Reader. If you can't, you are not ready to draft.

**Why:**
- Orwell's "scrupulous writer" asks, for every sentence: "What am I trying to say? What words will express it? What image or idiom will make it clearer? Is this image fresh enough to have an effect?" [P, [Orwell, "Politics and the English Language", 1946](https://www.orwellfoundation.com/the-orwell-foundation/orwell/essays-and-other-works/politics-and-the-english-language/)].
- NN/g: "When you write for everyone, you write for no one" [P, [Loranger, "Plain Language Is for Everyone, Even Experts", 2017-10-08](https://www.nngroup.com/articles/plain-language-experts/)].
- Diátaxis says an explanation should be "closely bounded" by a "why" question [P, [Diátaxis: Explanation](https://diataxis.fr/explanation/)].
- Paul Graham: writing "tells people something true and important that they didn't already know" [P, ["How to Write Usefully", 2020-02](https://paulgraham.com/useful.html)].

**Before:** "This piece looks at Kafka: its history, the alternatives, and the future of streaming."
**After:** "Kafka keeps messages after they are read, and that one choice explains most of its design."

### 1.2 Put the answer first

**Do:** State the conclusion in the first paragraph (the TL;DR). Then give the details, most important first.

**Why:**
- The inverted pyramid: "the most important information (or what might even be considered the conclusion) is presented first". It is "perfectly suited for the web" because skimmers who stop early still leave with the main point [P, [Schade, NN/g, 2018-02-11](https://www.nngroup.com/articles/inverted-pyramid/)].
- Its newsroom origin: "The inverted pyramid puts the most newsworthy information at the top". Its critics call it "the anti-story" [P, [Scanlan, Poynter, 2003-06-20](https://www.poynter.org/reporting-editing/2003/writing-from-the-top-down-pros-and-cons-of-the-inverted-pyramid/)].
- Google's raters are told to place "the most helpful and essential MC near the top of the page". They rate a page Low if a large amount of "filler" makes the helpful content hard to find [P, [Search Quality Rater Guidelines](https://static.googleusercontent.com/media/guidelines.raterhub.com/en//searchqualityevaluatorguidelines.pdf), version of 2025-09-11, §5.2.2].
- Microsoft: "Get to the point fast. Lead with what's most important." [P, [Top 10 tips](https://learn.microsoft.com/en-us/style-guide/top-10-tips-style-voice), 2026-07-02].

**Before:** "There are many factors to consider when thinking about connection reuse in HTTP."
**After:** "HTTP/2 sends many requests over one TCP connection, so one lost packet can stall all of them."

### 1.3 Write for people who scan: headings that say the point, first words that carry it

**Do:**
- Make each H2 a plain statement of what the section shows.
- Start headings, sentences and list items with the words that carry the meaning.
- Never use "Background", "Overview" or "Introduction" as a heading.

**Why:**
- The F-pattern: people read across the top, then a shorter band, then down the left edge. It happens when text has little formatting and the reader wants to be efficient. NN/g's advice: "Start headings and subheadings with the words carrying most information" [P, [NN/g, F-pattern, 2017-11-12](https://www.nngroup.com/articles/f-shaped-pattern-reading-web-content/)].
- The layer-cake pattern: readers fixate "mostly on the page's headings and subheadings". So headings must be "Descriptive of all topics in the section, and only topics in the section" [P, [Pernice, NN/g, 2019-08-04](https://www.nngroup.com/articles/layer-cake-pattern-scanning/)].
- Users "see about 2 words" of most list items and links. In NN/g's test, for 35% of links, people who saw only the first 11 characters "hadn't the faintest idea" where the link went [P, [Nielsen, "First 2 Words", 2009-04-06](https://www.nngroup.com/articles/first-2-words-a-signal-for-scanning/)].
- GOV.UK: headings should be "descriptive – avoid generic headings like 'Introduction'". They should not "be questions – they're hard to frontload and users want answers, not questions" [P, [GOV.UK: clear structure](https://guidance.publishing.service.gov.uk/writing-to-gov-uk-standards/writing-guidelines/clear-structure/)].
- Microsoft: "If readers don't read the headings, they probably won't read the text that follows, either." Use "parallel sentence structure for all headings at the same level" [P, [Microsoft: Headings](https://learn.microsoft.com/en-us/style-guide/scannable-content/headings)].

**Before:** `## Background`
**After:** `## Retries tripled the load in four minutes`

### 1.4 Use plain words, even for experts

**Do:**
- Use the everyday word when one exists.
- Keep a technical term when it is the right term. Define it at first use.
- Cut buzzwords.

**Why:**
- GOV.UK cites research on legal documents: "80% of people preferred sentences written in clear English". "The more educated the person and the more specialist their knowledge, the greater their preference for plain English." Also: "Where you need to use specialist terms, you can. They're not always jargon." [P, [GOV.UK: clear language](https://guidance.publishing.service.gov.uk/writing-to-gov-uk-standards/writing-guidelines/clear-language/)].
- NN/g's study with domain experts (an IT manager, a professor and others) found that "even highly educated online readers crave succinct information that is easy to scan". NN/g suggests a 10th–12th grade reading level for experts. "Text beyond the 12th grade reading level requires too much mental effort, even for highly educated people." [P, [NN/g, 2017](https://www.nngroup.com/articles/plain-language-experts/)].
- Orwell's rules ii and v: "Never use a long word where a short one will do." "Never use a foreign phrase, a scientific word or a jargon word if you can think of an everyday English equivalent." [P, Orwell].
- Paul Graham: "the easier something is to read, the more deeply readers will engage with it". Also: "you can't assume that writing about a difficult topic means you can use difficult words" [P, ["Write Simply", 2021-03](https://paulgraham.com/simply.html)].
- Julia Evans keeps some jargon on purpose: "otherwise it's impossible for the reader to search & learn more about the concept" [P, ["Simple explanations without sounding condescending", 2020-11-15](https://jvns.ca/blog/2020/11/15/simple-explanations-without-sounding-condescending/)].

**Before:** "Leverage idempotency semantics to facilitate exactly-once processing guarantees."
**After:** "Give each request an ID, so the server can spot a repeat and skip it. This is called an *idempotency key*."

### 1.5 Short sentences, one idea each, with some variety

**Do:**
- Aim for 15–20 words on average.
- Split any sentence over about 25 words, unless it is a deliberate list or a deliberate long sentence.
- Keep each paragraph to one topic, starting with the sentence that states it.

**Why:**
- GOV.UK: "Try to split up sentences that are over 25 words long". "Paragraphs should have no more than 5 sentences each" [P, GOV.UK clear language].
- NN/g: "Sentences should be no more than 15–20 words" [P, NN/g 2017].
- Cloudflare's docs go further: "Keep sentences between 8 to 12 words." "One sentence should be one idea." [P, [Cloudflare writing guidelines](https://developers.cloudflare.com/style-guide/documentation-content-strategy/writing-guidelines/), updated 2026-08-20].
- US federal plain-language guidance: "Express only one idea in each sentence". "Long, complicated sentences often mean that you aren't sure about what you want to say." Also: "Limit each paragraph or section to one topic." [P, [Digital.gov: Clear and short](https://digital.gov/guides/plain-language/writing/clear-short)].
- The counterweight, from Poynter's Roy Peter Clark: "Fear not the long sentence". A long sentence is fine if it is easy to follow. But "Use shorter words, sentences and paragraphs at points of complexity" [P, ["Fifty Writing Tools: Quick List", 2006-06-30](https://www.poynter.org/reporting-editing/2006/fifty-writing-tools-quick-list/), tools 7 and 11].

**Before:** "When the primary fails, the replica, which has been receiving the write-ahead log asynchronously and may therefore be a few seconds behind, is promoted, which means some writes acknowledged to clients can be lost." (36 words)
**After:** "When the primary fails, a replica takes over. That replica copies the write-ahead log asynchronously, so it can be a few seconds behind. Writes from those seconds are lost, even though clients were told they succeeded."

### 1.6 Active voice: name who did what

**Do:** Put the actor first and use a strong verb. Use the passive only when the actor is unknown or does not matter.

**Why:**
- Orwell rule iv: "Never use the passive where you can use the active." [P].
- Strunk (1918), rule 10: "The active voice is usually more direct and vigorous than the passive" [P, [*The Elements of Style*, Project Gutenberg](https://www.gutenberg.org/files/37134/37134-h/37134-h.htm)].
- Digital.gov: "Active voice makes it clear who should do what." [P, [Digital.gov: Writing](https://digital.gov/guides/plain-language/writing)].
- Google's developer style guide: "Use active voice" and "Use second person" [P, [Highlights](https://developers.google.com/style/highlights)].
- Clark adds the exception: use the passive "to showcase the 'victim' of action" [P, tool 4].

**Before:** "The cache was invalidated and requests were sent to the origin."
**After:** "The deploy cleared the cache, so every request went to the origin."

### 1.7 Concrete before abstract: example first, then the idea

**Do:**
- Open with a real system, a real number or a real failure. Then name the idea.
- Give at least one worked example that puts numbers on each step.

**Why:**
- Julia Evans lists "starting out abstract" and "no examples" among the patterns that make explanations confusing. "Leaving out examples makes it harder for the reader to relate the new terminology" [P, ["Patterns in confusing explanations", 2021-08-19](https://jvns.ca/blog/confusing-explanations/)].
- She says the way to good examples is "to start with real code that you wrote and then remove irrelevant details" [P, ["Write good examples by starting with real code", 2021-07-08](https://jvns.ca/blog/2021/07/08/writing-great-examples/)].
- Nicky Case's method for explaining: "start on the ground", then "move up, step by step" up the ladder of abstraction [P, ["How I Make an Explorable Explanation", 2017-09-20](https://blog.ncase.me/how-i-make-an-explorable-explanation/)].
- The US Institute of Education Sciences reviewed the learning research in 2007. Two of its seven recommendations, both rated "Moderate" evidence, were "Connect and integrate abstract and concrete representations of concepts" and "Interleave worked example solutions with problem-solving exercises" [P, [IES Practice Guide NCER 2007-2004](https://ies.ed.gov/ncee/WWC/Docs/PracticeGuide/20072004.pdf), Pashler et al., 2007-09]. The studies were of students, so applying them to adult engineers reading online is [inferred].
- Strunk, rule 12: "Prefer the specific to the general, the definite to the vague, the concrete to the abstract." [P]. Clark: "Get the name of the dog." [P, tool 14].

**Before:** "Backpressure is a flow-control mechanism in which downstream consumers signal capacity constraints upstream."
**After:** "A consumer can handle 100 messages a second, and the producer sends 500. Backpressure is how the consumer says 'slow down' before its queue overflows."

### 1.8 Show the mechanism with a picture, and walk the reader through it in words

**Do:**
- Draw the one diagram that shows how the thing works.
- In the text, walk the diagram in order and name its boxes.
- Never write "see the diagram" and stop there.

**Why:**
- IES recommendation 3 ("Moderate" evidence): "Combine graphics with verbal descriptions". This integration "leads to better learning than simply presenting text alone" [P, IES 2007].
- Google: keep the key content in text, because answer engines read text (see [`aeo-citation-signals.md`](aeo-citation-signals.md) §1) [P via that note].
- Julia Evans warns against "fun illustrations on dry explanations". A picture does not rescue an unclear text [P, Evans 2021].

**Before:** "The diagram below shows the flow." (then nothing)
**After:** "Follow the dashed arrow in the diagram: the timed-out request goes back into the same queue. That loop is why the load keeps growing after the trigger is gone."

### 1.9 Say why, not only what

**Do:** For every feature, design choice or step, say what problem it solves or what it costs.

**Why:**
- Evans calls "'what' without 'why'" a pattern of confusing explanations: "a weak 'why' is much better than no 'why'" [P, Evans 2021].
- Cloudflare: "Explain the 'why' before the 'how'." [P, Cloudflare guidelines].
- Diátaxis: explanation provides context ("design decisions, history, and technical constraints") and "can and must consider alternatives" [P, Diátaxis].

**Before:** "Raft has leader election, log replication and snapshots."
**After:** "Raft elects one leader so that only one node decides the order of writes. Replication copies that order to the others, and snapshots keep the log from growing forever."

### 1.10 Show where every fact came from; quote exactly; say what you didn't check

**Do:**
- Link each claim to the page that supports it.
- Use quotation marks only for exact words.
- Give a number its unit, its date and its caveat.
- Say "we infer" when you go beyond the Source.

**Why:**
- NPR's standard: "Attribute, attribute and attribute some more." "'Media reports' or 'sources say' is not good enough. Be specific." Verbatim words go in quotation marks. "We don't lightly edit quotes just to avoid putting them in quotes" [P, [NPR Ethics Handbook: Attribution](https://www.npr.org/about-npr/688424402/special-section-attribution); fetched by curl, as the fetch tool timed out].
- BBC Editorial Guidelines: "If content producers have been unable to verify material, they should usually say so and attribute the information" (3.4.6). Be "reluctant to rely on a single source" (3.4.4). Statistics "should be explained clearly, put into context, weighed, interpreted", with "important caveats and limitations explained" (3.4.8). Serious errors are "corrected promptly, clearly and appropriately" (3.4.34) [P, [BBC Editorial Guidelines §3.4](https://www.bbc.co.uk/editorialguidelines/guidelines/accuracy/guidelines); fetched by curl, as the fetch tool was blocked].
- The Reuters Handbook's line "It is our job to get it first but it is above all our job to get it right" reached me only through search snippets and mirrors [secondary, unverified] (Open questions).
- Google: "Of these aspects, trust is most important." [P, [Creating helpful, reliable, people-first content](https://developers.google.com/search/docs/fundamentals/creating-helpful-content), updated 2026-10-05].

**Before:** "It's widely known that Postgres struggles with lots of connections."
**After:** "Postgres uses a 'process per user' model ([docs](https://www.postgresql.org/docs/current/connect-estab.html)), so 5,000 open connections means 5,000 server processes, even when most are idle."

### 1.11 Add something the Sources don't have

**Do:** Every Piece needs its own contribution. That can be the connection between Sources, the trade-off nobody stated, the worked example, the comparison or the lesson for the Reader's own system. Restating the Source, even with credit, is not a contribution.

**Why:**
- Google's self-assessment asks: "Does the content provide original information, reporting, research, or analysis?" If the content draws on other sources, "does it avoid simply copying or rewriting those sources, and instead provide substantial additional value and originality?" [P, helpful content page].
- The rater guidelines rate a page Lowest if almost all of it is "copied, paraphrased, embedded, auto or AI generated, or reposted from other sources with little to no effort, little to no originality, and little to no added value ... even if the page assigns credit for the content to another source". Paraphrase includes "restating or summarizing the content on other pages" [P, QRG §4.6.6].
- NPR: "Do not quote or paraphrase another organization's material so much that you effectively make reading, watching or listening to their reports unnecessary." [P, NPR Attribution].

**Before (a Brief's "Why it matters" that only restates):** "GitHub said the outage lasted three hours and was caused by a database failover."
**After:** "A failover that takes minutes can turn into hours when every client retries at once. That is the loop in [metastable failures](metastable-failures), and the fix is a retry budget, not more database capacity."

### 1.12 Cut, then read it aloud, then cut again

**Do:**
- Draft fast. Then do two passes: first cut whole sections and paragraphs that don't serve the point, then cut words.
- Read it aloud to find where you stumble.

**Why:**
- Strunk, rule 13: "Vigorous writing is concise. A sentence should contain no unnecessary words, a paragraph no unnecessary sentences ... This requires not that the writer make all his sentences short ... but that he make every word tell." [P].
- Orwell rule iii: "If it is possible to cut a word out, always cut it out." Rule vi: "Break any of these rules sooner than say anything outright barbarous." [P].
- Paul Graham: "Write a bad version 1 as fast as you can; rewrite it over and over; cut out everything unnecessary". Also "read your essays out loud to see (a) where you stumble over awkward phrases and (b) which bits are boring" [P, ["Writing, Briefly", 2005-03](https://paulgraham.com/writing44.html)].
- Clark: "Cut big, then small." [P, tool 10].
- NN/g's 1997 study measured it: against a promotional control, the concise version scored 58% better on usability, the scannable one 47%, the objective one 27%, and all three combined 124% [P, [Morkes & Nielsen, 1997](https://www.nngroup.com/articles/concise-scannable-and-objective-how-to-write-for-the-web/)].
- Microsoft: "Avoid weak phrasing like *there is*, *there are*, and *there were*" [P].

**Before:** "It is important to note that, in order to be able to scale, the system basically needs to shard its data."
**After:** "To scale, the system shards its data."

---

## 2. Anatomy of an article: each part, its job, and common mistakes

In our Pieces, each part maps to a field or section. Kind-specific rules are in the writing guides. This section explains why each part exists.

| Part | Where it lives for us | Its one job | Common mistakes |
|---|---|---|---|
| **Headline** | `title` (also the H1 and the search title link) | Say what the piece is about, in words the Reader would search for, with the key term first | Puns, teasers, "Everything you need to know about…", the key term at the end, over ~65–70 characters |
| **Dek** (subheading or summary) | `description` (the meta description, the card line and often the snippet) | Add the "so what" that the title had no room for | Repeating the title, being vague, no verb, running long |
| **Lede / hook** | The paragraph(s) before the AI note | Make the Reader care, with one real, specific moment | Opening on a definition, on history, or on "In today's fast-paced world" |
| **Nut graf** | The last sentence of the hook, then the TL;DR | Say what this piece will give the Reader and why it matters now | Skipping it, so the Reader doesn't know where the story is going; or giving away so little that it reads as a tease |
| **Body sections** | The H2s | Walk the Reader through the mechanism, one point per section, in a logical order (often the diagram's order) | Generic headings; two ideas in one section; sections only a writer could love (the history of the field) |
| **Evidence and quotes** | Inline links; short quotes | Make each claim checkable; let exact words carry a point only where the wording matters | Vague attribution ("experts say"); a quote doing the analysis for you; a number with no date, unit or caveat |
| **Diagrams and code** | One Mermaid fence or more; code only under its licence | Show the mechanism; let the Reader run or recognise the real thing | A diagram the text never walks; toy code that nobody writes; commands "improved" from the README |
| **Ending** | "What it makes possible", the Verdict, the close of "Why it matters" | Tell the Reader what they can now do, measure or decide | A summary of what they just read; "only time will tell"; a new idea introduced in the last line |
| **Sources** | `## Sources` / `source_url` | Let anyone retrace the work | Listing pages the text never used; missing the Lead |

What the sources say, part by part:

- **Headline.** "Headline text has to stand on its own and make sense when the rest of the content is not available". Write it "in plain language, with no puns, and no 'cute' or 'clever' wordings". "Avoid teasers (click bait)" [P, [Loranger & Nielsen, NN/g, "Microcontent", 2017-01-29](https://www.nngroup.com/articles/microcontent-how-to-write-headlines-page-titles-and-subject-lines/)]. Front-load the key terms (NN/g First 2 Words) [P]. Google: "Write descriptive and concise text for your `<title>` elements" and "Avoid keyword stuffing". Google may build the title link from the `<title>`, the "Main visual title shown on the page", `<h1>` and `og:title` [P, [Influencing title links](https://developers.google.com/search/docs/appearance/title-link)]. GOV.UK: "65 characters or less (including spaces)", and find "what the public calls your content by looking at common search terms" [P, [GOV.UK: titles](https://guidance.publishing.service.gov.uk/writing-to-gov-uk-standards/writing-guidelines/clear-titles/)].
- **Dek.** GOV.UK: "Keep all summaries to 160 characters", reuse as few title words as possible, and "Keep summaries active and include a verb" [P, [GOV.UK: summaries](https://guidance.publishing.service.gov.uk/writing-to-gov-uk-standards/writing-guidelines/summaries/)]. Google: "Make sure your descriptions are truly descriptive", unique per page. "Snippets are primarily created from the page content itself", so the first paragraph and TL;DR matter as much as the description [P, [Snippets](https://developers.google.com/search/docs/appearance/snippet)].
- **Lede and nut graf.** Scanlan (Poynter) on the nut graf: it "tells the reader what the writer is up to; it delivers a promise of the story's content and message". It "justifies the story by telling readers why they should care", "provides a transition from the lead", and "often tells readers why the story is timely". Ken Wells of the *Wall Street Journal* calls it "a paragraph that says what this whole story is about and why you should read it". James B. Stewart warns against giving away too much and losing the curiosity the lead built [P, [Scanlan, "The nut graf…", Poynter, 2003-05-19](https://www.poynter.org/archive/2003/the-nut-graf-part-i/)]. NN/g: put the most important points in the first two paragraphs [P, F-pattern].
- **Body.** One topic per paragraph (Digital.gov); headings parallel at each level (Microsoft); "Work from a plan" and "Place gold coins along the path", meaning reward the Reader in the middle, not only at the ends (Clark tools 24, 32) [P].
- **Evidence.** §1.10 above. For quote length and copyright, see [`content-source-intake.md`](content-source-intake.md) §4.5.
- **Code.** Google's developer style guide: "precede a code sample with an introductory sentence or paragraph". "Wrap lines at 80 characters." "If a code block contains an omission, don't format the block as click-to-copy." [P, [Code samples](https://developers.google.com/style/code-samples)]. Evans: start from real code and strip it down; avoid "unrealistic examples" [P].
- **Ending.** Clark: "Write toward an ending. Help readers close the circle of meaning." [P, tool 39]. Nicky Case ends on an open question the learner can explore [P]. Our "what could go right" close does both.

---

## 3. How the craft differs across our five kinds

The craft is the same everywhere. What changes is the reader's question and the shape that answers it.

| Kind | Reader's question | Closest classic form | Shape that works | Where the craft is most at risk |
|---|---|---|---|---|
| **Brief** | "What just happened, and should I care?" | Hard news, inverted pyramid, with a short analysis | Title = the lesson. "What happened" = the lede facts (who, what, when in UTC), 2–4 sentences. "Why it matters" = the nut graf grown into analysis | Restating the Source (QRG §4.6.6; NPR "make reading… unnecessary"). In 250 words every sentence must work (Clark 37: "In short pieces of writing, don't waste a syllable") |
| **Concept** | "What is this idea, and how does it work?" | Explanation (Diátaxis); teaching | Concrete story → definition from whoever named it → mechanism walked through a diagram → worked example with numbers → where it shows up → what it is not → what to do | Starting abstract; too many new terms at once (Evans); a strained analogy; skipping the misconception ("what it is not") |
| **Case Study** | "How did they build it, and why that way?" | Feature story with a nut graf; reconstruction | A scale fact or reversal → nut graf → the problem → options weighed → the system walked step by step → costs → evidence and its gaps | Becoming a summary of the company's post (no added value); hiding the costs; glossing beyond the Source's words |
| **Comparison** | "Which one should I choose?" | Buyer's guide; parallel analysis | Same questions for every option, in the same order (Microsoft: parallel structure); a table; "when each wins", labelled as inference | Uneven depth across options; cells with no source behind them; a verdict the evidence doesn't support |
| **Repo Profile** | "What is it, and should I adopt it?" | Review | Answer first; quickstart verbatim; architecture; take a position ("Use it if / Don't use it if"); FAQ in the Reader's search words; verdict | No first-hand experience: Google's raters ask "which would you trust: a product review from someone who has personally used the product or a 'review' by someone who has not?" [P, QRG §3.4]. Stale numbers in prose |

Notes per kind:

- **Brief = news.** The Brief is the only kind where the inverted pyramid applies strictly. The title and description do the work of the lede, because the body opens with the AI note. "Why it matters" is the nut graf grown into analysis: why now, why you, what to learn. Scanlan's four jobs of the nut graf are a good test for it [P, Scanlan] [inferred].
- **Concept = explanation.** Diátaxis separates explanation from tutorials, how-to guides and reference. An explanation "permits *reflection*", is "understanding-oriented", and should not turn into step-by-step instruction or a reference dump [P, Diátaxis]. Our Concept template already follows this shape [inferred].
- **Case Study = feature.** The feature form came from the *Wall Street Journal*, where Barney Kilgore in 1941 said "It doesn't have to have happened today to be news" [P, quoted in Scanlan 2003]. A Case Study is news-adjacent but not tied to a date. Its value is the "options weighed" and "what it costs" sections, which a company's own post rarely gives [inferred].
- **Comparison = parallel structure.** Each option answers the same questions in the same order. That is Microsoft's parallel-heading rule applied to the whole piece, and it is why "the table wrote itself" in our first Comparison [P + observed in `comparison-writing-guide.md`].
- **Repo Profile = review.** E-E-A-T's first E is "Experience: Consider the extent to which the content creator has the necessary first-hand or life experience for the topic" [P, QRG §3.4]. A Profile built only from the README and Fact Pack has expertise and sources, but no first-hand use [inferred]. See gap 5.3 #5.

---

## 4. One-page self-edit checklist

Use this after the draft is done and **before** the claim-by-claim check against Sources. It checks the writing, not the facts. The Author's verification stays exactly as the guides describe it. Copy it into the ticket or PR.

```text
SELF-EDIT: {slug}

POINT
[ ] I can say the point in one sentence: ______________________________
[ ] The piece covers one idea; anything else is a Related link.

TOP OF THE PAGE
[ ] Title: key term in the first two words, plain words, no pun, under ~65 chars.
[ ] Description: adds the "so what", has a verb, doesn't repeat the title, <=150 chars.
[ ] Hook: one real, linked moment (system, number, failure), not a definition or history.
[ ] Nut graf: the hook's last sentence says what this piece gives the Reader.
[ ] TL;DR: the answer in plain words; a reader who stops here is not misled.

STRUCTURE
[ ] Read only the headings: do they tell the story in order?
[ ] No heading says "Background", "Overview", "Introduction" or "Conclusion".
[ ] No question headings outside a Repo Profile FAQ.
[ ] Each paragraph has one topic, stated in its first sentence.
[ ] Each H2 is under 250 words; cut, don't pad.

SENTENCES
[ ] Average about 20 words or fewer; every sentence over ~30 words is split or earns its length.
[ ] Active voice with a named actor ("the client retried", not "requests were retried").
[ ] Hardest parts use the shortest sentences.

WORDS
[ ] Every technical term is defined at first use, and kept where it is the right term.
[ ] One name per thing: the same term every time (no synonyms for the same component).
[ ] No hype or filler: revolutionary, game-changing, magic, simply, just, easy,
    obviously, "it is important to note", "please note", "there is/are" openings.
[ ] At most one analogy, and it says where it stops being true.
[ ] No idioms or pop-culture references a non-native reader would miss.

TEACHING
[ ] A concrete example comes before the abstract statement.
[ ] At least one worked example with real numbers (Concept, Case Study, Comparison).
[ ] Every "what" has a "why" or a cost.
[ ] The common confusion is named and corrected ("what it is not").
[ ] The text walks the diagram box by box; the diagram is never left to speak alone.
[ ] Code: introduced by a sentence, realistic, under its licence; README commands verbatim.

EVIDENCE (writing side; the Author's verification checks the facts)
[ ] Every claim has its link where the claim is made.
[ ] Every number has a unit, a date or version, and its caveat.
[ ] Quotes are exact, short, credited, and followed by our own analysis.
[ ] Each inference says "we infer" in its own sentence.

OUR CONTRIBUTION
[ ] Name what this Piece adds that no Source says: ____________________
[ ] Delete every paragraph that only restates a Source.

ENDING
[ ] The last section tells the Reader what they can now do, measure or decide.
[ ] No new idea in the last paragraph; no recap of what they just read.

LAST PASS
[ ] Read it aloud once. Fix every place you stumbled or got bored.
[ ] Cut big (sections, paragraphs), then small (words). Then cut 10% more.
```

Where each line comes from: the Orwell questions, Graham, Strunk and Clark (point, cuts, read-aloud); NN/g, GOV.UK and Microsoft (top of page, headings, sentences); Google's developer style guide, Microsoft's global-English tips and Evans (words, analogies); IES, Evans, Case and Diátaxis (teaching); NPR and BBC (evidence); Google Search and the QRG (contribution). All are cited in §1–§2.

---

## 5. Gap analysis: research vs our guides and house voice

These are suggestions only. No editorial file was changed. Files read: `house-byline-voice.md` and the five `*-writing-guide.md` files with their templates. Counts are from 74 published Markdown Pieces (13 Briefs, 61 Explainers, Fact Packs excluded), checked on 2026-11-16 [observed].

### 5.1 What we already do well (and the research that backs it)

| Our rule | Backed by |
|---|---|
| "Lead with a concrete hook… not a definition" (voice rule 1) | Evans "starting out abstract"; Nicky Case "start on the ground"; IES rec. 4; Strunk rule 12 |
| "Answer first" TL;DR (voice rule 2) | Inverted pyramid (NN/g, Poynter); QRG §5.2.2; Microsoft "Get to the point fast" |
| "One idea per piece" (voice rule 3) | Orwell's questions; Diátaxis "closely bounded"; Evans "too many concepts at a time" |
| "Diagram first", sections walk it (voice rule 4) | IES rec. 3 (graphics with verbal descriptions); Clark "Work from a plan" |
| "Define every term at first use… No hype words" (voice rule 5) | GOV.UK; NN/g experts; Cloudflare "Define acronyms and technical terms"; Morkes & Nielsen "objective" +27% |
| "Name the trade-offs plainly" (voice rule 6) | Evans "'what' without 'why'"; Diátaxis "consider alternatives"; Graham's "strength" without becoming false |
| "Show the receipts… Say 'we infer'" (voice rule 7) | NPR "Attribute everything"; BBC 3.4.6, 3.4.20; Google "trust is most important" |
| Second person, "we" for the team (voice rule 8) | Google style "Use second person"; Microsoft "write like you speak" |
| Headings renamed to "a plain sentence that states its point" (all guides) | NN/g layer-cake; GOV.UK "descriptive"; Microsoft headings |
| "Quote short and credited… follow each quote with our own analysis" | NPR (exact words in quotes, don't make the original unnecessary); QRG §4.6.6 |
| Brief: "Why it matters is HastCode's own analysis" | Google "original… analysis"; QRG paraphrase = Lowest even with credit |
| Concept "what it is not" section | Evans "explaining the 'wrong' way… without saying it's wrong"; Diátaxis boundaries |
| Repo Profile FAQ "phrased the way the Reader types them" | GOV.UK titles "common search terms" |
| Title ~70 chars and description ≤150 chars, with past overruns noted | GOV.UK 65 / 160; Google "descriptive and concise". Observed: 2 of 74 titles over 70 characters, and no description over 150 [observed] |
| AI-assistance note on every Piece | Google: disclosures are "useful for content where someone might think 'How was this created?'"; BBC 3.4.29 |

### 5.2 What is missing or in tension with the evidence

1. **Word counts read as targets.** The Concept, Case Study and Comparison guides say "Aim for 2,000 to 2,800 words". Google lists "writing to a particular word count" as a sign of search-engine-first content: "(No, we don't.)". The raters mark down "filler" that "can artificially inflate content" [P]. [`aeo-citation-signals.md`](aeo-citation-signals.md) found word count has about zero correlation with AI Overview citation. As a planning range it is harmless. As an "aim" it invites padding [inferred].
2. **"Short sentences" has no number.** Across our Explainers, the mean sentence is about 21.7 words and the median 19. 32% of sentences run over 25 words and 12% over 35. Briefs: mean 21.5, 31% over 25 words [observed, rough script: it splits on full stops and removes code, tables, headings and link URLs]. GOV.UK splits at 25 words, NN/g advises 15–20, Cloudflare 8–12 [P]. Our median is fine. The long tail is the problem: one sentence in eight over 35 words is hard going on a phone [inferred].
3. **"Edit" is a step with no content.** The voice pipeline ends "verify every claim → edit → publish", and none of the guides says what editing means. The Author's verification is about facts. Nothing checks the prose (cuts, sentence length, headings read alone, read-aloud). §4 fills this gap.
4. **No nut graf.** The Comparison guide says the opening ends "on the question the piece answers". The Concept and Case Study openings end on a turn ("What held DynamoDB down for hours was its own retry traffic"), which is a good lede, but no sentence promises what the piece will give the Reader. The AI note and TL;DR follow straight away. Scanlan's nut graf "provides a transition from the lead and explains the lead and its connection to the rest of the story" [P]. The TL;DR partly does this job. A one-sentence bridge would make the turn explicit [inferred].
5. **Repo Profiles have no first-hand experience.** The guide grounds Profiles in "the Fact Pack, the README, the release pages", and says to copy the quickstart verbatim, but never asks anyone to run it. The QRG's first E is first-hand Experience, and its example is a product review by someone who used the product [P]. A line like "We ran the quickstart on {date} on {OS}; it worked as written / it needed X" is the cheapest first-hand signal we could add [inferred]. It costs the Author's time, so it is a choice, not a fix.
6. **Analogies have no rule.** The voice guide says nothing about analogies. Evans names "strained analogies" as a confusion pattern. Google's developer style guide says to avoid "figurative language, which includes metaphors". Orwell rule i bans worn-out figures [P]. At the same time, Orwell asks "What image or idiom will make it clearer?" Good explainers use analogies with care. Suggested rule: at most one, and say where it breaks [opinion].
7. **No "one name per thing" rule.** Microsoft: "Use one word for a concept, and use it consistently. Avoid using synonyms to refer to the same concept" [P, [Microsoft global writing tips](https://learn.microsoft.com/en-us/style-guide/global-communications/writing-tips)]. The Case Study guide learned this the hard way ("Terms drift between Sources… use the Source's own term"), but the voice guide has no general rule.
8. **The hype-word list is short.** We ban "revolutionary", "game-changing" and "magic". Google's developer style guide also flags "simply", "It's easy", "quickly", "please note", "at this time", exclamation marks and starting every sentence with the same phrase [P, [Google tone](https://developers.google.com/style/tone), updated 2026-05-27]. Evans notes that words like "simply" make a reader who is stuck feel stupid [secondary: from a compilation of her advice, not fetched from her page].
9. **Question headings outside FAQs.** 64 of our 816 headings end in "?" [observed]. Most are Repo Profile FAQ `###` questions, which the guide requires and which match how people search. GOV.UK says headings should not be questions because "users want answers, not questions". The AEO note's SE Ranking study found question headings did worse than plain ones [P; STUDY via aeo note]. Keep FAQ questions; avoid them in H2s [opinion].
10. **The AI note may sit in the most valuable spot.** In Explainers the AI note is the second block, after the hook and before the TL;DR. In a Brief it is the very first thing the Reader sees. NN/g says the first two paragraphs get the most attention [P]. Google wants the disclosure, but its page doesn't say where it should go [P]. Whether moving it lower would be just as honest is the Author's and the Editorial Policy's call. Open question, not a recommendation.
11. **"TL;DR" is slang to Google's docs style guide.** It lists "tl;dr" among "Internet slang" to avoid [P]. That guide is for documentation, while "TL;DR" is a common, well-understood blog heading, and our build requires it. No change suggested; noted only so it is a conscious choice [opinion].
12. **Numbers need context, not only links.** The Comparison guide asks for "a measurement quoted with its caveat". The voice guide doesn't. BBC 3.4.8 asks that numbers be "explained clearly, put into context… with important caveats" [P]. A general rule would help: every number carries its unit, date or version, and what it is compared with.
13. **Pseudonymous house byline vs "who wrote this".** Google asks: "Is it self-evident to your visitors who authored your content?" [P]. The QRG accepts an "alias or username" when "there is good reason for anonymity" (§2.5.2, §5.5) [P]. Mira Vance is disclosed as a house byline, and that decision is settled. The research only suggests the Editorial Policy and Author page say clearly who stands behind the verification [inferred]. Open question.

### 5.3 Suggested edits (draft wording, for the Author to accept or reject)

Ranked by how much they would improve the writing for the effort. **The top five are in bold.**

1. **Concept, Case Study and Comparison guides, "Aim for 2,000 to 2,800 words".** Change it to: *"Most land between 2,000 and 2,800 words. This is a ceiling to check, not a target: stop when the point is made. Google has no preferred word count, and padding is the 'filler' its raters mark down."*
2. **`house-byline-voice.md`, voice rule 5.** Change "Short sentences" to: *"Short sentences: about 20 words on average. Split anything over 30 unless it is a deliberate list. Use the shortest sentences at the hardest points."*
3. **`house-byline-voice.md`, drafting pipeline.** Change "→ draft in this voice → verify every claim → edit →" to *"→ draft in this voice → self-edit with the checklist (`docs/research/article-writing-craft.md` §4, or a copy in `docs/editorial/`) → verify every claim → edit →"*. Also put the checklist in the PR template or ticket the Author already uses [inferred: I did not check whether one exists].
4. **Concept and Case Study guides, "Opening (before TL;DR)".** Add: *"End the hook with one sentence that says what this piece will show the Reader (the 'nut graf'): 'This is how that loop forms, and the three defences that break it.'"* The Comparison guide already does this ("ending on the question the piece answers").
5. **Repo Profile guide, "In 60 seconds".** Add: *"When the Author has run the quickstart, say so in one line after the block: when, on what, and whether it worked as written. Never change the commands; report what happened."* It stays optional, because it costs the Author time.
6. `house-byline-voice.md`: add voice rule 9, *"One name per thing. Pick the Source's term and use it every time; mention a synonym once, in brackets, if Readers search for it."*
7. `house-byline-voice.md`: extend rule 5's banned list to *"revolutionary, game-changing, magic, simply, just (as in 'just add'), easy, obviously, 'it is important to note', 'please note'"*.
8. `house-byline-voice.md`: add *"Analogies: at most one per piece, and say where it stops being true."*
9. All guides, "The skeleton": add *"Only Repo Profile FAQ headings may be questions; every H2 is a statement."*
10. `house-byline-voice.md` rule 7: add *"Every number carries its unit, its date or version, and what it is compared with."*

---

## 6. Further reading: the best primary sources (start here)

1. **NN/g, ["How Users Read on the Web"](https://www.nngroup.com/articles/how-users-read-on-the-web/) (1997) and ["How People Read Online: New and Old Findings"](https://www.nngroup.com/articles/how-people-read-online/) (2020).** The evidence that people scan, and what to do about it, from the people who ran the studies.
2. **GOV.UK, ["Writing to GOV.UK standards"](https://guidance.publishing.service.gov.uk/writing-to-gov-uk-standards/writing-guidelines/) (structure, clear language, titles, summaries).** The clearest short rules for plain web writing, with numbers you can apply.
3. **Google, [developer documentation style guide](https://developers.google.com/style/highlights) ([tone](https://developers.google.com/style/tone), [code samples](https://developers.google.com/style/code-samples)).** Voice, word choice and code rules written for developer readers.
4. **Microsoft, [Writing Style Guide: top 10 tips](https://learn.microsoft.com/en-us/style-guide/top-10-tips-style-voice) and [global writing tips](https://learn.microsoft.com/en-us/style-guide/global-communications/writing-tips).** Short before/after pairs, and the rules for readers whose first language isn't English.
5. **Julia Evans, ["Patterns in confusing explanations"](https://jvns.ca/blog/confusing-explanations/).** Thirteen named mistakes in technical explanation, from someone known for clear ones. The closest thing to a checklist for our Concepts.
6. **Google Search Central, ["Creating helpful, reliable, people-first content"](https://developers.google.com/search/docs/fundamentals/creating-helpful-content) and the [Search Quality Rater Guidelines](https://static.googleusercontent.com/media/guidelines.raterhub.com/en//searchqualityevaluatorguidelines.pdf) §3.2, §3.4, §4.6.6, §5.2.2.** What "quality" means to the search engine that sends us most traffic, including for AI-assisted content.
7. **George Orwell, ["Politics and the English Language"](https://www.orwellfoundation.com/the-orwell-foundation/orwell/essays-and-other-works/politics-and-the-english-language/) (1946).** The six rules and four questions, still the shortest course in clear prose.
8. **Paul Graham, ["Writing, Briefly"](https://paulgraham.com/writing44.html) and ["Write Simply"](https://paulgraham.com/simply.html).** A working writer's whole process in a page: draft fast, cut, read aloud.
9. **Chip Scanlan (Poynter), ["The nut graf…"](https://www.poynter.org/archive/2003/the-nut-graf-part-i/), and Roy Peter Clark, ["Fifty Writing Tools: Quick List"](https://www.poynter.org/reporting-editing/2006/fifty-writing-tools-quick-list/).** The newsroom view of ledes, nut grafs, cutting and endings.
10. **NPR Ethics Handbook, ["Attribution"](https://www.npr.org/about-npr/688424402/special-section-attribution), and BBC Editorial Guidelines, ["Accuracy" §3.4](https://www.bbc.co.uk/editorialguidelines/guidelines/accuracy/guidelines).** How serious newsrooms attribute, quote, handle numbers and correct errors.

---

## 7. Open questions and unverifiable items

- **AP Stylebook / AP News Values:** `ap.org` refused the fetch tool, and the Stylebook is paywalled. No AP claim is made here.
- **Reuters Handbook of Journalism:** I found only search snippets and third-party mirrors (mediareform.org.uk PDF, accountablejournalism.org). The "get it right" line in §1.10 is [secondary, unverified]. Reuters' own current copy was not found.
- **BBC and NPR:** the fetch tool was blocked (BBC) or timed out (NPR). Both pages were downloaded with `curl` and read as text on 2026-11-16, so the quotes are from the pages themselves [P]. Neither page showed a last-updated date.
- **Plain Writing Act of 2010:** congress.gov returned 403. Only Digital.gov's description ("content for the public is written for its specific audience") was read. The statute's own definition was not.
- **plainlanguage.gov** now redirects to Digital.gov's plain-language guides. The old Federal Plain Language Guidelines are, per Digital.gov, "archived in a GitHub repository" I did not open. Digital.gov's "Test for understanding" page link returned 404 at the URL I tried.
- **GOV.UK pages** show no last-updated dates. GOV.UK's "20 to 28%" figure gives no citation on the page. It matches NN/g's 2008 "How little do users read?" finding, but I did not fetch that article [unverified].
- **Bartosz Ciechanowski:** I found no first-party writing by him about his method. Only third parties describe his work, so none of it is used. **Stripe:** no first-party writing guidance was found or used. **Cloudflare:** only its docs style guide was used; the 2025 Content Style Guide PDF linked from search was not opened.
- **Zinsser, *On Writing Well*:** no primary text is legally reachable online, so it is not cited. **Strunk** is cited from the public-domain 1918 edition, not the later Strunk & White.
- **Analogies and misconceptions research:** I did not find a primary study that fits writing for engineers. The analogy rule in §5.3 #8 is [opinion], based on Evans and the Google and Orwell rules.
- **IES practice guide** evidence comes from classroom studies of students. Applying it to experienced engineers reading online is [inferred].
- **Our sentence-length numbers** come from a rough script (it splits sentences on full stops, removes code, tables and headings, and keeps TL;DR bullets). They show the scale of the problem, not exact values. Rerun it after any change if the Author wants a before/after.
- **For the Author to decide:** where the AI-assistance note should sit (§5.2 #10). Whether the Editorial Policy and Author page should name who verifies (§5.2 #13). Whether to add first-hand runs to Repo Profiles (§5.3 #5), given the time each one costs.

---

## Sources (all accessed 2026-11-16)

**How people read online (Nielsen Norman Group)**
- https://www.nngroup.com/articles/how-users-read-on-the-web/ (Nielsen, 1997-09-30)
- https://www.nngroup.com/articles/concise-scannable-and-objective-how-to-write-for-the-web/ (Morkes & Nielsen, 1997)
- https://www.nngroup.com/articles/how-people-read-online/ (Moran, 2020-04-05)
- https://www.nngroup.com/articles/f-shaped-pattern-reading-web-content/ (2017-11-12, reviewed 2026-08-19)
- https://www.nngroup.com/articles/layer-cake-pattern-scanning/ (Pernice, 2019-08-04)
- https://www.nngroup.com/articles/inverted-pyramid/ (Schade, 2018-02-11)
- https://www.nngroup.com/articles/first-2-words-a-signal-for-scanning/ (Nielsen, 2009-04-06)
- https://www.nngroup.com/articles/plain-language-experts/ (Loranger, 2017-10-08)
- https://www.nngroup.com/articles/microcontent-how-to-write-headlines-page-titles-and-subject-lines/ (Loranger & Nielsen, 2017-01-29)

**Plain language (government)**
- https://guidance.publishing.service.gov.uk/writing-to-gov-uk-standards/writing-guidelines/clear-structure/
- https://guidance.publishing.service.gov.uk/writing-to-gov-uk-standards/writing-guidelines/clear-language/
- https://guidance.publishing.service.gov.uk/writing-to-gov-uk-standards/writing-guidelines/clear-titles/
- https://guidance.publishing.service.gov.uk/writing-to-gov-uk-standards/writing-guidelines/summaries/
- https://digital.gov/guides/plain-language ; https://digital.gov/guides/plain-language/writing ; https://digital.gov/guides/plain-language/writing/clear-short
- Not fetched: https://www.congress.gov/bill/111th-congress/house-bill/946/text (403); https://digital.gov/guides/plain-language/testing (404)

**Technical writing style guides**
- https://developers.google.com/style/highlights ; https://developers.google.com/style/tone (updated 2026-05-27) ; https://developers.google.com/style/code-samples
- https://learn.microsoft.com/en-us/style-guide/top-10-tips-style-voice (2026-07-02)
- https://learn.microsoft.com/en-us/style-guide/global-communications/writing-tips (2023-06-14)
- https://learn.microsoft.com/en-us/style-guide/scannable-content/headings (2018-07-24)
- https://developers.cloudflare.com/style-guide/documentation-content-strategy/writing-guidelines/ (updated 2026-08-20)

**News craft and newsroom standards**
- https://www.poynter.org/archive/2003/the-nut-graf-part-i/ (Scanlan, 2003-05-19)
- https://www.poynter.org/reporting-editing/2003/writing-from-the-top-down-pros-and-cons-of-the-inverted-pyramid/ (Scanlan, 2003-06-20)
- https://www.poynter.org/reporting-editing/2006/fifty-writing-tools-quick-list/ (Clark, 2006-06-30)
- https://www.npr.org/about-npr/688424402/special-section-attribution (read via curl)
- https://www.bbc.co.uk/editorialguidelines/guidelines/accuracy/guidelines (§3.4; read via curl)
- Secondary only: https://www.mediareform.org.uk/wp-content/uploads/2015/12/Reuters_Handbook_of_Journalism.pdf and https://accountablejournalism.org/ethics-codes/standards-and-values (Reuters Handbook mirrors, not opened)
- Not fetched: https://www.ap.org/about/news-values-and-principles/ (blocked)

**Classic essays**
- https://www.orwellfoundation.com/the-orwell-foundation/orwell/essays-and-other-works/politics-and-the-english-language/ (1946)
- https://www.gutenberg.org/files/37134/37134-h/37134-h.htm (Strunk, *The Elements of Style*, 1918)
- https://paulgraham.com/writing44.html (2005-03) ; https://paulgraham.com/simply.html (2021-03) ; https://paulgraham.com/useful.html (2020-02) ; https://paulgraham.com/words.html (2022-02)

**Explanation craft**
- https://jvns.ca/blog/confusing-explanations/ (2021-08-19)
- https://jvns.ca/blog/2021/07/08/writing-great-examples/ (2021-07-08)
- https://jvns.ca/blog/2020/11/15/simple-explanations-without-sounding-condescending/ (2020-11-15)
- https://blog.ncase.me/how-i-make-an-explorable-explanation/ (2017-09-20)
- https://diataxis.fr/explanation/
- https://ies.ed.gov/ncee/WWC/Docs/PracticeGuide/20072004.pdf (Pashler et al., IES NCER 2007-2004, 2007-09)

**Search engines**
- https://developers.google.com/search/docs/fundamentals/creating-helpful-content (updated 2026-10-05)
- https://developers.google.com/search/docs/fundamentals/using-gen-ai-content (updated 2026-10-01)
- https://developers.google.com/search/blog/2023/02/google-search-and-ai-content (2023-02-08)
- https://static.googleusercontent.com/media/guidelines.raterhub.com/en//searchqualityevaluatorguidelines.pdf (General Guidelines, 2025-09-11: §2.5.2, §3.2, §3.4, §4.6.6, §5.2.2, §5.5)
- https://developers.google.com/search/docs/appearance/title-link
- https://developers.google.com/search/docs/appearance/snippet

**Example link used in §1.10**
- https://www.postgresql.org/docs/current/connect-estab.html ("PostgreSQL implements a “process per user” client/server model"; checked 2026-11-16)

**Observed (2026-11-16)**: `docs/editorial/house-byline-voice.md`, the five `*-writing-guide.md` files and the concept and brief templates; `content/explainers/distributed-systems/durable-execution.md` (read as a sample); `content/pages/editorial-policy.md` (Corrections section); sentence-length, heading and title/description counts over `content/briefs/*.md` and `content/explainers/*/*.md`, from a throwaway script in the session scratchpad.

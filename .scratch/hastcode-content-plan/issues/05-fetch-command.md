# 05: Fetch command

**What to build:** The Author runs one command with `owner/name` and the path of the Profile's Markdown. It fetches the Fact Pack from the GitHub REST API in at most ten calls (repository, languages, latest ten releases, contributors list with total count, root and `.github` trees for the file flags, README blob reference, stargazer count), writes the JSON beside the Markdown with stable key order and indentation, and writes the README text to the builder's cache directory for drafting. GitHub's `NOASSERTION` license sets `needs_review` instead of failing. The command reads one token from the environment (the variable GitHub Actions provides), sends an identifying User-Agent, honours `Retry-After` and the primary rate-limit headers with a bounded backoff, and stops with a clear message instead of retrying past the limit; without a token it refuses and says why. A refresh flag walks the content tree, re-fetches every Fact Pack, rewrites only those whose content changed (ignoring fetched-at), prints one line per changed repository saying what changed, and exits with a distinct code when anything changed.

Spec: `../spec.md` ("Fetch command", User Stories 17–20, 31 prerequisites). Research: `docs/research/github-trending-geography-and-data-access.md` sections 2 and 3 (rate limits, endpoints, 202 behaviour of the stats endpoint, security policy detection via tree).

**Blocked by:** 01 (Repo Profile kind renders from a Fact Pack) for the Fact Pack shape

**Status:** ready-for-agent

- [ ] Against an in-process fake GitHub serving recorded JSON, the written Fact Pack equals a golden file
- [ ] A second run against the same responses changes only fetched-at
- [ ] Refresh mode over a tree with two Profiles, one changed upstream, rewrites exactly that one, prints its one-line summary, and exits with the "changed" code; nothing changed exits zero
- [ ] `NOASSERTION` yields `needs_review: true` and a printed note; a normal SPDX id yields false
- [ ] A 403 with rate-limit headers stops with the documented message and non-zero exit; a 202 from any endpoint is retried a bounded number of times
- [ ] Missing token is refused with a message naming the variable
- [ ] README text lands in the cache directory, never in the content tree; the Fact Pack carries only its URL and SHA
- [ ] Command documented in the README (usage, token, what it writes) and in the agent instructions

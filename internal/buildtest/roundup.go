package buildtest

import "strings"

// The Friday roundup (ADR-0005, spec "Editorial docs") is an ordinary
// Brief: "What happened" lists the week's repositories, one line each,
// linking their Repo Profiles; "Why it matters" is the commentary;
// source_url is the week's trending page; `explainers` lists the linked
// Profiles. These fixtures prove it needs no new piece type.

// roundupRepos are the first five Repo Profiles (docs/editorial/
// repo-profile-selection.md): repository, Vertical, title and description.
var roundupRepos = []struct{ repo, vertical, title, description string }{
	{"apache/iggy", "distributed-systems", "Apache Iggy: a persistent message stream in Rust",
		"Persistent message streaming platform written in Rust"},
	{"xerj-org/xerj", "databases-storage", "xerj: a single-binary search engine",
		"Local-first search engine with BM25, vector and hybrid search"},
	{"NVIDIA/OpenShell", "infrastructure-cloud", "OpenShell: a sandboxed runtime for agents",
		"Sandboxed runtime with policy-enforced containment"},
	{"deepseek-ai/DeepInfer", "ai-ml-systems", "DeepInfer: DeepSeek's inference framework",
		"Large-scale inference framework"},
	{"php/frankenphp", "languages-runtimes", "FrankenPHP: a PHP application server on Caddy",
		"PHP application server built on Caddy with a worker mode"},
}

// RoundupProfiles returns five valid Repo Profiles, one per repository the
// roundup Brief links, each in its own Vertical with its Fact Pack (a copy
// of ValidFactPack naming that repository). Their slugs are owner-name,
// lower case (apache-iggy, xerj-org-xerj, nvidia-openshell,
// deepseek-ai-deepinfer, php-frankenphp). Add them with WithRepoProfile.
func RoundupProfiles() []Explainer {
	profiles := make([]Explainer, 0, len(roundupRepos))
	for _, r := range roundupRepos {
		owner, name, _ := strings.Cut(r.repo, "/")
		p := ValidRepoProfile()
		p.Vertical = r.vertical
		p.Slug = strings.ToLower(owner + "-" + name)
		p.Title = r.title
		p.Repo = r.repo
		fp := ValidFactPack()
		fp.Slug = r.repo
		fp.Name = name
		fp.URL = "https://github.com/" + r.repo
		fp.Description = r.description
		p.FactPack = &fp
		profiles = append(profiles, p)
	}
	return profiles
}

// RoundupBrief returns the Friday roundup Brief for RoundupProfiles: it
// meets every Brief gate when the five Profiles are in the tree (it links
// each by bare slug and lists the five in `explainers`).
func RoundupBrief() Brief {
	return Brief{
		Slug:        "trending-repositories-2026-10-23",
		Title:       "Five repositories worth a look this week",
		Description: "A message stream, a search engine, an agent sandbox, an inference server and a PHP runtime: what trended and why it matters.",
		Tags:        []string{"open-source", "roundup"},
		Published:   "2026-10-23T09:00:00Z",
		Author:      DefaultAuthor().Slug,
		Hero:        "trending-repositories-2026-10-23.png",
		HeroAlt:     "Five repository cards side by side",
		HeroImage:   HeroPNG(1600, 900),
		SourceURL:   "https://github.com/trending?since=weekly",
		Explainers:  []string{"apache-iggy", "xerj-org-xerj", "nvidia-openshell", "deepseek-ai-deepinfer", "php-frankenphp"},
		Body:        RoundupBriefBody,
	}
}

// RoundupBriefBody is the Markdown body of RoundupBrief.
const RoundupBriefBody = `## What happened

Five repositories we profiled this week:

- [Apache Iggy](apache-iggy): Kafka-style streams in one Rust binary.
- [xerj](xerj-org-xerj): BM25, vector and hybrid search in a single binary.
- [OpenShell](nvidia-openshell): NVIDIA's sandbox for autonomous agents.
- [DeepInfer](deepseek-ai-deepinfer): DeepSeek's inference framework, open-sourced.
- [FrankenPHP](php-frankenphp): a PHP application server built on Caddy.

## Why it matters

Three of the five are written in Rust, one in Go and one in Python, and three of them run as a single binary rather than a cluster. Each Profile says who should adopt it today and who should wait.
`

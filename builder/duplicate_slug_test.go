package builder_test

import (
	"testing"

	"hastcode.com/site/internal/buildtest"
)

func TestDuplicateSlug_AcrossVerticalsFailsBothExplainers(t *testing.T) {
	a := buildtest.ValidExplainer()
	b := buildtest.ValidExplainer()
	b.Vertical = "databases-storage"

	out, res := buildtest.Build(t, buildtest.Config(), buildtest.NewTree().WithExplainer(a).WithExplainer(b))

	buildtest.RequireErrorRule(t, res, "explainers/distributed-systems/raft-leader-election.md", "duplicate-slug", "explainers/databases-storage/raft-leader-election.md")
	buildtest.RequireErrorRule(t, res, "explainers/databases-storage/raft-leader-election.md", "duplicate-slug", "explainers/distributed-systems/raft-leader-election.md")
	for _, url := range []string{"/distributed-systems/raft-leader-election/", "/databases-storage/raft-leader-election/"} {
		if buildtest.Exists(out, url) {
			t.Errorf("page %s was written for a duplicate slug", url)
		}
	}
}

func TestDuplicateSlug_ExplainerAndBriefSharingASlugFail(t *testing.T) {
	brief := `---
title: "Raft lands in a new database"
description: "A new database ships Raft for leader election."
tags: ["raft"]
published: "2026-10-20T09:00:00Z"
author: "ada"
source_url: "https://example.com/raft-news"
explainers: ["raft-leader-election"]
---

## What happened

A database adopted Raft. See [how Raft elects a leader](/distributed-systems/raft-leader-election/).

## Why it matters

Leader election is now off the shelf.
`
	tree := buildtest.ValidTree().With("briefs/raft-leader-election.md", []byte(brief))

	out, res := buildtest.Build(t, buildtest.Config(), tree)

	buildtest.RequireErrorRule(t, res, "explainers/distributed-systems/raft-leader-election.md", "duplicate-slug", "briefs/raft-leader-election.md")
	buildtest.RequireErrorRule(t, res, "briefs/raft-leader-election.md", "duplicate-slug", "explainers/distributed-systems/raft-leader-election.md")
	for _, url := range []string{"/distributed-systems/raft-leader-election/", "/briefs/raft-leader-election/"} {
		if buildtest.Exists(out, url) {
			t.Errorf("page %s was written for a duplicate slug", url)
		}
	}
}

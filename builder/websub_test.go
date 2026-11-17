package builder_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"hastcode.com/site/builder"
	"hastcode.com/site/internal/buildtest"
)

// webSubHub is the hub every feed names: Google's public hub, whose own
// page (https://pubsubhubbub.appspot.com/, read 2026-11-17) asks for
// <atom:link rel="hub" href="https://pubsubhubbub.appspot.com/"> under
// rss/channel. README "Feeds, sitemaps, robots and IndexNow".
const webSubHub = "https://pubsubhubbub.appspot.com/"

// hub is the channel's Atom rel="hub" link (WebSub discovery).
func (f rss) hub() string {
	for _, l := range f.Channel.Links {
		if l.XMLName.Space == "http://www.w3.org/2005/Atom" && l.Rel == "hub" {
			return l.Href
		}
	}
	return ""
}

// allFeeds is every feed the launch Verticals give a build.
var allFeeds = []string{
	"/feed.xml",
	"/briefs/feed.xml",
	"/distributed-systems/feed.xml",
	"/databases-storage/feed.xml",
	"/apis-networking/feed.xml",
	"/infrastructure-cloud/feed.xml",
	"/ai-ml-systems/feed.xml",
	"/languages-runtimes/feed.xml",
}

// webSubBody is the changed-feed artefact, <out>/websub.json beside
// indexnow.json: the hub to notify and the topics (absolute feed URLs)
// whose bytes changed since the previous build of the same mode. The
// post-deploy step sends hub.mode=publish with one hub.url per topic.
type webSubBody struct {
	Hub    string   `json:"hub"`
	Topics []string `json:"topics"`
}

func readWebSub(t *testing.T, out string) (webSubBody, bool) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(out, "websub.json"))
	if os.IsNotExist(err) {
		return webSubBody{}, false
	}
	if err != nil {
		t.Fatal(err)
	}
	var b webSubBody
	if err := json.Unmarshal(data, &b); err != nil {
		t.Fatalf("websub.json does not decode: %v\n%s", err, data)
	}
	return b, true
}

// changedFeeds builds tree with cfg's cache (the previous build's
// manifest) and returns the feed paths the artefact lists, site-relative,
// or nil when the build wrote no artefact.
func changedFeeds(t *testing.T, cfg builder.Config, tree buildtest.Tree) []string {
	t.Helper()
	out := buildtest.MustBuild(t, cfg, tree)
	body, ok := readWebSub(t, out)
	if !ok {
		return nil
	}
	if body.Hub != webSubHub {
		t.Errorf("websub.json hub = %q, want %q", body.Hub, webSubHub)
	}
	if len(body.Topics) == 0 {
		t.Error("websub.json lists no topics; no change must mean no artefact")
	}
	var paths []string
	for _, u := range body.Topics {
		p, ok := strings.CutPrefix(u, siteURL)
		if !ok {
			t.Errorf("topic %q is not an absolute URL on %s (the feed's rel=self)", u, siteURL)
		}
		paths = append(paths, p)
	}
	return paths
}

func webSubConfig(t *testing.T) builder.Config {
	cfg := buildtest.Config()
	cfg.Now = buildClock
	cfg.CacheDir = t.TempDir() // the build cache CI keeps between Production builds
	return cfg
}

func TestWebSub_FirstBuildListsEveryFeedAndAnUnchangedRebuildListsNone(t *testing.T) {
	cfg := webSubConfig(t)

	if got := changedFeeds(t, cfg, feedTree()); !reflect.DeepEqual(got, sorted(allFeeds)) {
		t.Errorf("first build (no manifest yet) lists %v, want every feed %v", got, sorted(allFeeds))
	}
	if got := changedFeeds(t, cfg, feedTree()); got != nil {
		t.Errorf("unchanged rebuild lists %v, want no artefact", got)
	}
}

func TestWebSub_ANewOrChangedBriefListsTheSiteAndBriefsFeeds(t *testing.T) {
	cfg := webSubConfig(t)
	explainers := func() buildtest.Tree {
		return buildtest.ValidTree().WithExplainer(buildtest.ValidComparison())
	}
	want := []string{"/briefs/feed.xml", "/feed.xml"}

	changedFeeds(t, cfg, explainers())
	if got := changedFeeds(t, cfg, explainers().WithBrief(buildtest.ValidBrief())); !reflect.DeepEqual(got, want) {
		t.Errorf("new Brief lists %v, want %v", got, want)
	}

	revised := buildtest.ValidBrief()
	revised.Body = strings.Replace(revised.Body, "halves election latency", "roughly halves election latency", 1)
	if revised.Body == buildtest.ValidBrief().Body {
		t.Fatal("fixture: the revision changed nothing")
	}
	if got := changedFeeds(t, cfg, explainers().WithBrief(revised)); !reflect.DeepEqual(got, want) {
		t.Errorf("changed Brief lists %v, want %v", got, want)
	}
}

func TestWebSub_AChangedExplainerListsTheSiteFeedAndItsVerticalsFeed(t *testing.T) {
	cfg := webSubConfig(t)
	changedFeeds(t, cfg, feedTree())

	revised := buildtest.ValidComparison()
	revised.Body = strings.Replace(revised.Body, "## Sources", "One more sentence before the sources.\n\n## Sources", 1)
	if revised.Body == buildtest.ValidComparison().Body {
		t.Fatal("fixture: the revision changed nothing")
	}
	want := []string{"/databases-storage/feed.xml", "/feed.xml"}
	if got := changedFeeds(t, cfg, feedTree().WithExplainer(revised)); !reflect.DeepEqual(got, want) {
		t.Errorf("changed Explainer lists %v, want %v", got, want)
	}

	added := explainerFixture("distributed-systems", "paxos", "concept", 5)
	want = []string{"/distributed-systems/feed.xml", "/feed.xml"}
	if got := changedFeeds(t, cfg, feedTree().WithExplainer(revised).WithExplainer(added)); !reflect.DeepEqual(got, want) {
		t.Errorf("new Explainer lists %v, want %v", got, want)
	}
}

// Feeds carry the newest 50 pieces each. With 51 databases-storage
// Explainers (plus raft-leader-election), the oldest is in no feed, so
// revising it changes no feed's bytes.
func TestWebSub_AChangedExplainerOlderThanEveryFeedsNewestFiftyListsNoFeed(t *testing.T) {
	cfg := webSubConfig(t)
	tree := func(oldest buildtest.Explainer) buildtest.Tree {
		tree := buildtest.ValidTree().WithExplainer(oldest)
		for i := 2; i <= 51; i++ {
			e := explainerFixture("databases-storage", fmt.Sprintf("db-%02d", i), "concept", 1)
			e.Published = fmt.Sprintf("2026-08-%02dT%02d:00:00Z", 1+i/24, i%24)
			tree.WithExplainer(e)
		}
		return tree
	}
	oldest := explainerFixture("databases-storage", "db-01", "concept", 1)
	oldest.Published = "2026-07-01T09:00:00Z"

	changedFeeds(t, cfg, tree(oldest))
	out := buildtest.MustBuild(t, webSubConfig(t), tree(oldest))
	for _, path := range []string{"/feed.xml", "/databases-storage/feed.xml"} {
		if links := itemLinks(readFeed(t, out, path)); len(links) != 50 || slices.Contains(links, siteURL+"/databases-storage/db-01/") {
			t.Fatalf("fixture: %s has %d items (want 50, without db-01)", path, len(links))
		}
	}

	revised := oldest
	revised.Body = strings.Replace(revised.Body, "split votes are rare", "split votes are rare and short", 1)
	if got := changedFeeds(t, cfg, tree(revised)); got != nil {
		t.Errorf("revising an Explainer in no feed lists %v, want no feeds", got)
	}
}

func sorted(s []string) []string {
	s = slices.Clone(s)
	slices.Sort(s)
	return s
}

func TestWebSub_EveryFeedNamesTheHubBesideItsOwnSelfLink(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), feedTree())

	for _, path := range allFeeds {
		f := readFeed(t, out, path) // also checks it parses as RSS 2.0
		if got := f.hub(); got != webSubHub {
			t.Errorf("%s atom:link rel=hub = %q, want %q", path, got, webSubHub)
		}
		if _, self := f.link(); self != siteURL+path {
			t.Errorf("%s atom:link rel=self = %q, want the feed's own absolute URL (the WebSub topic)", path, self)
		}
		var hubs, selfs int
		for _, l := range f.Channel.Links {
			switch l.Rel {
			case "hub":
				hubs++
			case "self":
				selfs++
			}
		}
		if hubs != 1 || selfs != 1 {
			t.Errorf("%s has %d rel=hub and %d rel=self links, want one each", path, hubs, selfs)
		}
	}
}

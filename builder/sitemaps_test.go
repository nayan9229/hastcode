package builder_test

import (
	"encoding/xml"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"hastcode.com/site/internal/buildtest"
)

// urlset is a sitemap (standard or news) as the sitemap tests read it.
type urlset struct {
	XMLName xml.Name `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
	URLs    []struct {
		Loc     string `xml:"loc"`
		LastMod string `xml:"lastmod"`
		News    *struct {
			Name            string `xml:"publication>name"`
			Language        string `xml:"publication>language"`
			PublicationDate string `xml:"publication_date"`
			Title           string `xml:"title"`
		} `xml:"http://www.google.com/schemas/sitemap-news/0.9 news"`
	} `xml:"url"`
}

func readSitemap(t *testing.T, out, urlPath string) urlset {
	t.Helper()
	data, err := os.ReadFile(buildtest.StaticPath(out, urlPath))
	if err != nil {
		t.Fatalf("sitemap %s: %v", urlPath, err)
	}
	var u urlset
	if err := xml.Unmarshal(data, &u); err != nil {
		t.Fatalf("sitemap %s does not parse: %v\n%s", urlPath, err, data)
	}
	return u
}

func (u urlset) locs() []string {
	var locs []string
	for _, x := range u.URLs {
		locs = append(locs, x.Loc)
	}
	return locs
}

func (u urlset) lastmod(loc string) string {
	for _, x := range u.URLs {
		if x.Loc == loc {
			return x.LastMod
		}
	}
	return ""
}

// sitemapTree has an Explainer per kind, a Brief, the transparency pages
// and a draft that production must leave out.
func sitemapTree() buildtest.Tree {
	d := buildtest.ValidExplainer()
	d.Slug = "raft-membership-changes"
	d.Draft = true
	return feedTree().WithExplainer(d)
}

func TestSitemap_ListsEveryIndexablePageExactlyOnce(t *testing.T) {
	out := buildtest.MustBuild(t, buildtest.Config(), sitemapTree())
	locs := readSitemap(t, out, "/sitemap.xml").locs()

	// Every directory-index page the build wrote, except the noindex
	// "check your inbox" page; /404.html is noindex and not a page URL.
	var want []string
	for _, p := range buildtest.URLPaths(t, out) {
		if strings.HasSuffix(p, "/") && p != "/subscribed/" {
			want = append(want, siteURL+p)
		}
	}
	got := slices.Clone(locs)
	slices.Sort(got)
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("sitemap locs =\n%s\nwant every indexable page once:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}

	for _, must := range []string{
		"/", "/distributed-systems/raft-leader-election/", "/databases-storage/btree-vs-lsm-tree/",
		"/briefs/etcd-raft-rewrite/", "/briefs/", "/distributed-systems/", "/explainers/comparison/",
		"/authors/" + buildtest.DefaultAuthor().Slug + "/", "/about/",
	} {
		if !slices.Contains(locs, siteURL+must) {
			t.Errorf("sitemap lacks %s", must)
		}
	}
	for _, loc := range locs {
		if !strings.HasSuffix(loc, "/") {
			t.Errorf("sitemap loc %s has no trailing slash", loc)
		}
		if strings.Contains(loc, "raft-membership-changes") || strings.Contains(loc, "404") || strings.Contains(loc, "subscribed") {
			t.Errorf("sitemap lists %s (draft or noindex page)", loc)
		}
	}
}

func TestSitemap_LastmodIsTheContentDate(t *testing.T) {
	cfg := buildtest.Config()
	cfg.Now = time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)
	out := buildtest.MustBuild(t, cfg, sitemapTree())
	sm := readSitemap(t, out, "/sitemap.xml")

	for path, want := range map[string]string{
		// An Explainer changes when it is Reviewed.
		"/distributed-systems/raft-leader-election/": "2026-10-15T18:30:00Z",
		// A Brief is never revised in place.
		"/briefs/etcd-raft-rewrite/": "2026-10-20T08:15:00-04:00",
		// A listing changes when its newest piece appears.
		"/":                   "2026-10-20T08:15:00-04:00",
		"/databases-storage/": "2026-10-02T09:00:00Z",
		// A page with no dates of its own: when the build first saw it.
		"/about/": "2026-10-21T12:00:00Z",
	} {
		if got := sm.lastmod(siteURL + path); got != want {
			t.Errorf("lastmod of %s = %q, want %q", path, got, want)
		}
	}
	for _, u := range sm.URLs {
		if u.LastMod == "" {
			t.Errorf("%s has no lastmod", u.Loc)
		}
	}
}

// buildClock is 2026-10-21 12:00 UTC: ValidBrief (2026-10-20 08:15 -04:00)
// is 23 h 45 m old.
var buildClock = time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)

func TestNewsSitemap_ListsOnlyBriefsFromTheLast48Hours(t *testing.T) {
	old := buildtest.ValidBrief()
	old.Slug = "kafka-drops-zookeeper"
	old.Title = "Kafka drops ZooKeeper"
	old.Published = "2026-10-19T11:59:00Z" // 48 h 1 m before the build clock
	cfg := buildtest.Config()
	cfg.Now = buildClock

	out := buildtest.MustBuild(t, cfg, feedTree().WithBrief(old))
	sm := readSitemap(t, out, "/news-sitemap.xml")

	if got := sm.locs(); len(got) != 1 || got[0] != siteURL+"/briefs/etcd-raft-rewrite/" {
		t.Fatalf("news sitemap locs = %v, want only the Brief from the last 48 h", got)
	}
	n := sm.URLs[0].News
	if n == nil {
		t.Fatal("news sitemap entry has no news:news block")
	}
	if n.Name != "Hast Code" || n.Language != "en" || n.PublicationDate != "2026-10-20T08:15:00-04:00" || n.Title != "etcd ships a rewritten Raft library" {
		t.Errorf("news:news = %+v", *n)
	}
	if !buildtest.Exists(out, "/briefs/kafka-drops-zookeeper/") {
		t.Error("the older Brief's own page was not written")
	}
}

func TestNewsSitemap_NeverListsExplainers(t *testing.T) {
	fresh := buildtest.ValidExplainer()
	fresh.Slug = "raft-log-compaction"
	fresh.Published = "2026-10-21T09:00:00Z"
	fresh.Reviewed = "2026-10-21T09:00:00Z"
	cfg := buildtest.Config()
	cfg.Now = buildClock

	out := buildtest.MustBuild(t, cfg, feedTree().WithExplainer(fresh))
	for _, loc := range readSitemap(t, out, "/news-sitemap.xml").locs() {
		if !strings.HasPrefix(loc, siteURL+"/briefs/") {
			t.Errorf("news sitemap lists %s, which is not a Brief", loc)
		}
	}
}

func TestNewsSitemap_HoldsAtMost1000Entries(t *testing.T) {
	if testing.Short() {
		t.Skip("builds 1,001 Briefs")
	}
	tree := buildtest.ValidTree()
	for i := range 1001 {
		b := buildtest.ValidBrief()
		b.Slug = "brief-" + strconv.Itoa(i)
		b.Title = "Brief number " + strconv.Itoa(i)
		b.Published = buildClock.Add(-time.Duration(i) * time.Minute).Format(time.RFC3339)
		tree.WithBrief(b)
	}
	cfg := buildtest.Config()
	cfg.Now = buildClock

	out := buildtest.MustBuild(t, cfg, tree)
	locs := readSitemap(t, out, "/news-sitemap.xml").locs()
	if len(locs) != 1000 {
		t.Fatalf("news sitemap has %d entries, want the 1,000 newest", len(locs))
	}
	if locs[0] != siteURL+"/briefs/brief-0/" || slices.Contains(locs, siteURL+"/briefs/brief-1000/") {
		t.Errorf("news sitemap keeps the wrong 1,000: first %s", locs[0])
	}
}

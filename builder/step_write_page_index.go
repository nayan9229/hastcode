package builder

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func init() {
	// First in the write phase: the sitemap and IndexNow steps read it.
	registerStep(phaseWrite, 5, "index written pages", indexWrittenPages)
}

// pageIndex is every indexable page this build wrote (read back from the
// output after the link gate, so a page the gate removed is not in it),
// and how it differs from the previous successful build's.
type pageIndex struct {
	// pages are the indexable pages, sorted by canonical URL.
	pages []indexedPage
	// changed are the canonical URLs that are new, whose main content
	// changed since the previous build, or that the previous build had and
	// this one does not (IndexNow takes deletions too). Sorted.
	changed []string
	// manifest is what savePageManifest stores for the next build.
	manifest pageManifest
}

// indexedPage is one indexable HTML page.
type indexedPage struct {
	canonical string // absolute, from <link rel="canonical">
	// lastmod is when the page's content last changed, as written in the
	// page (RFC 3339): an article's content date, a listing's newest
	// piece, or, for a page with no dates, when its content last changed
	// between builds.
	lastmod string
}

// pageManifest is the per-page record kept between builds in the build
// cache (Config.CacheDir, persisted by Vercel under node_modules/.cache;
// never in the public output).
type pageManifest struct {
	SiteURL string                   `json:"site_url"`
	Pages   map[string]manifestEntry `json:"pages"` // by canonical URL
}

type manifestEntry struct {
	// Hash is the SHA-256 of the page's <main> element: site chrome (nav,
	// stylesheet hash) changes on every page at once and is not a content
	// change worth pinging.
	Hash string `json:"hash"`
	// Changed is when the build first saw this Hash (RFC 3339).
	Changed string `json:"changed"`
}

// manifestPath is the previous build's page manifest. Production and
// preview builds keep separate manifests so a preview never hides a change
// from the next production build.
func (s *site) manifestPath() string {
	return filepath.Join(s.cfg.CacheDir, "pages", string(s.cfg.Mode)+".json")
}

func indexWrittenPages(s *site) error {
	prev, err := s.loadPageManifest()
	if err != nil {
		return err
	}
	contentDates := s.contentDates()
	now := s.cfg.Now.UTC().Format(time.RFC3339)

	idx := &pageIndex{manifest: pageManifest{SiteURL: s.cfg.SiteURL, Pages: map[string]manifestEntry{}}}
	root := filepath.Join(s.out.Dir, "static")
	err = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".html") {
			return err
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		h := readHead(data)
		if h.canonical == "" || h.noindex {
			return nil
		}
		if _, dup := idx.manifest.Pages[h.canonical]; dup {
			return nil
		}
		sum := sha256.Sum256(mainElement(data))
		entry := manifestEntry{Hash: hex.EncodeToString(sum[:]), Changed: now}
		if old, ok := prev.Pages[h.canonical]; ok && old.Hash == entry.Hash {
			entry.Changed = old.Changed
		} else {
			idx.changed = append(idx.changed, h.canonical)
		}
		idx.manifest.Pages[h.canonical] = entry

		lastmod := contentDates[h.canonical]
		if lastmod == "" {
			lastmod = h.newestTime
		}
		if lastmod == "" {
			lastmod = entry.Changed
		}
		idx.pages = append(idx.pages, indexedPage{canonical: h.canonical, lastmod: lastmod})
		return nil
	})
	if err != nil {
		return err
	}
	for u := range prev.Pages {
		if _, ok := idx.manifest.Pages[u]; !ok {
			idx.changed = append(idx.changed, u)
		}
	}
	slices.SortFunc(idx.pages, func(a, b indexedPage) int { return strings.Compare(a.canonical, b.canonical) })
	slices.Sort(idx.changed)
	s.index = idx
	return nil
}

// contentDates are the absolute URLs of articles with the date their
// content last changed: an Explainer's Review (or publication, if later),
// a Brief's publication. The page itself also shows other pieces' dates
// (reading lists), so these are not read from the HTML.
func (s *site) contentDates() map[string]string {
	dates := map[string]string{}
	for _, e := range s.publishedExplainers() {
		t := e.Published
		if e.Reviewed.After(t) {
			t = e.Reviewed
		}
		dates[s.absURL(explainerURL(e))] = t.Format(time.RFC3339)
	}
	for _, b := range s.publishedBriefs() {
		dates[s.absURL(briefURL(b))] = b.Published.Format(time.RFC3339)
	}
	return dates
}

// loadPageManifest reads the previous build's manifest; none (first
// build, cleared cache) or one for another site URL is an empty manifest.
func (s *site) loadPageManifest() (pageManifest, error) {
	var m pageManifest
	data, err := os.ReadFile(s.manifestPath())
	if errors.Is(err, fs.ErrNotExist) {
		return pageManifest{}, nil
	}
	if err != nil {
		return m, err
	}
	if err := json.Unmarshal(data, &m); err != nil || m.SiteURL != s.cfg.SiteURL {
		return pageManifest{}, nil // unreadable or foreign: start afresh
	}
	return m, nil
}

// pageHead is what the index reads from a written page.
type pageHead struct {
	canonical string
	noindex   bool
	// newestTime is the latest <time datetime> on the page, as written.
	newestTime string
}

func readHead(doc []byte) pageHead {
	var h pageHead
	var newest time.Time
	z := html.NewTokenizer(bytes.NewReader(doc))
	for {
		switch z.Next() {
		case html.ErrorToken:
			return h
		case html.StartTagToken, html.SelfClosingTagToken:
			name, hasAttr := z.TagName()
			if !hasAttr {
				continue
			}
			attrs := map[string]string{}
			for more := true; more; {
				var k, v []byte
				k, v, more = z.TagAttr()
				attrs[string(k)] = string(v)
			}
			switch string(name) {
			case "link":
				if attrs["rel"] == "canonical" {
					h.canonical = attrs["href"]
				}
			case "meta":
				if attrs["name"] == "robots" && strings.Contains(attrs["content"], "noindex") {
					h.noindex = true
				}
			case "time":
				if t, err := time.Parse(time.RFC3339, attrs["datetime"]); err == nil && t.After(newest) {
					newest, h.newestTime = t, attrs["datetime"]
				}
			}
		}
	}
}

// mainElement is the page's <main>…</main>, or the whole page if it has
// none.
func mainElement(doc []byte) []byte {
	start := bytes.Index(doc, []byte("<main"))
	end := bytes.LastIndex(doc, []byte("</main>"))
	if start < 0 || end < start {
		return doc
	}
	return doc[start:end]
}

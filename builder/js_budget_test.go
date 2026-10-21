package builder_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"hastcode.com/site/internal/buildtest"
)

// Spec story 61: total client JS ≤ 60 KB per page. This test is the
// authority for that budget, not Lighthouse CI: Lighthouse has to block
// Plausible on previews (it would otherwise send the Vercel bypass secret
// to a third party) and its resource-summary counts compressed transfer
// bytes. Here every script a page loads initially is summed uncompressed:
// inline code, same-origin files from the output, and an upper bound for
// each third-party script. Pagefind loads only on interaction, so it is
// not initial JS.
const jsBudget = 60 * 1024

// thirdPartyScriptBytes are upper bounds for the uncompressed size of the
// third-party scripts the site may load, since their bytes are not in the
// build output. Plausible's script.js measured 2,841 bytes uncompressed
// (1,283 gzipped) on 2026-10-21; 4 KB leaves headroom for its updates.
var thirdPartyScriptBytes = map[string]int{
	"https://plausible.io/js/script.js": 4 * 1024,
}

func TestJSBudget_EveryPageStaysWithin60KB(t *testing.T) {
	cfg := buildtest.Config() // v1 production shape: analytics on, ads off
	cfg.PlausibleDomain = "hastcode.com"
	out := buildtest.MustBuild(t, cfg, buildtest.ValidTree())

	pages := 0
	err := filepath.WalkDir(filepath.Join(out, "static"), func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".html") {
			return err
		}
		pages++
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if total := pageJSBytes(t, out, data); total > jsBudget {
			t.Errorf("%s ships %d bytes of JS, over the %d-byte budget", p, total, jsBudget)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if pages == 0 {
		t.Fatal("no HTML pages in output")
	}
}

// pageJSBytes sums a page's executable scripts: inline code, same-origin
// files from the output, and the known size of each third-party script.
func pageJSBytes(t *testing.T, outDir string, html []byte) int {
	t.Helper()
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	doc.Find("script").Each(func(_ int, s *goquery.Selection) {
		if s.AttrOr("type", "") == "application/ld+json" {
			return
		}
		src, external := s.Attr("src")
		switch {
		case !external:
			total += len(s.Text())
		case strings.HasPrefix(src, "/"):
			info, err := os.Stat(buildtest.StaticPath(outDir, src))
			if err != nil {
				t.Errorf("script %s not in output: %v", src, err)
				return
			}
			total += int(info.Size())
		default:
			n, known := thirdPartyScriptBytes[src]
			if !known {
				t.Errorf("third-party script %s has no size budget", src)
			}
			total += n
		}
	})
	return total
}

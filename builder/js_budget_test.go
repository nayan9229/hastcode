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

// Spec story 61: client JS ≤ 60 KB per page. This test is the authority
// for that budget, not Lighthouse CI (which blocks third-party hosts on
// previews and counts compressed transfer bytes). Here every first-party
// script a page loads initially is summed uncompressed: inline code and
// same-origin files from the output. Since ADR-0008 the budget bounds our
// own code only: gtag.js alone is ~449 KB, and the third-party scripts in
// exemptThirdPartyScripts are named exemptions; any other third-party
// script fails. Pagefind loads only on interaction, so it is not initial JS.
const jsBudget = 60 * 1024

// exemptThirdPartyScripts are the URL prefixes of the third-party scripts a
// page may load outside the budget (ADR-0008), each without blocking parsing.
var exemptThirdPartyScripts = []string{
	"https://www.googletagmanager.com/gtag/js?id=",
	"https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js?client=",
}

func TestJSBudget_EveryPageStaysWithin60KB(t *testing.T) {
	cfg := buildtest.Config()
	cfg.AdsEnabled = true
	cfg.AdSenseClient = "ca-pub-1234567890123456"
	cfg.AdSenseSlots = map[string]string{"after-tldr": "1", "after-what-happened": "2", "after-body": "3"}
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

// pageJSBytes sums a page's first-party executable scripts: inline code and
// same-origin files from the output. Exempt third-party scripts add
// nothing; any other third-party script is an error.
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
			for _, prefix := range exemptThirdPartyScripts {
				if strings.HasPrefix(src, prefix) {
					return
				}
			}
			t.Errorf("third-party script %s is neither first-party nor a named exemption (ADR-0008)", src)
		}
	})
	return total
}

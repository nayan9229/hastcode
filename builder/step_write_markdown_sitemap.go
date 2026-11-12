package builder

import (
	"fmt"
	"strings"
	"time"

	"hastcode.com/site/internal/views"
)

func init() {
	registerStep(phaseWrite, 33, "write sitemap.md", writeMarkdownSitemap)
}

// The Agent Surface (ADR-0011): sitemap.md lists every live Piece whose
// page this build wrote as a link to its Markdown twin with its kind and
// published day, Explainers grouped by Vertical in configured order, newest
// first, then Briefs.
func writeMarkdownSitemap(s *site) error {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s: every live Piece\n\n", views.PublicationName)
	b.WriteString("Each link is the Piece's Markdown twin; drop `.md` and add a trailing slash for the page.\n")
	explainers := newestExplainers(s.publishedExplainers())
	for _, v := range s.sortedVerticals() {
		var lines []string
		for _, e := range explainers {
			if e.Vertical != v.Slug || !s.isWritten(explainerURL(e)) {
				continue
			}
			lines = append(lines, fmt.Sprintf("- [%s](%s) · %s · %s\n",
				e.Title, s.absURL(twinURL(explainerURL(e))), kindOf(e.Kind).name, e.Published.Format(time.DateOnly)))
		}
		if len(lines) > 0 {
			fmt.Fprintf(&b, "\n## %s\n\n%s", v.Name, strings.Join(lines, ""))
		}
	}
	var briefs []string
	for _, br := range newestBriefs(s.publishedBriefs()) {
		if !s.isWritten(briefURL(br)) {
			continue
		}
		briefs = append(briefs, fmt.Sprintf("- [%s](%s) · %s · %s\n",
			br.Title, s.absURL(twinURL(briefURL(br))), briefLabel, br.Published.Format(time.DateOnly)))
	}
	if len(briefs) > 0 {
		fmt.Fprintf(&b, "\n## Briefs\n\n%s", strings.Join(briefs, ""))
	}
	return s.out.WriteStatic(markdownSitemapURL, []byte(b.String()))
}

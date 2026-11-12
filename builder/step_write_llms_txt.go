package builder

import (
	"fmt"
	"strings"

	"hastcode.com/site/internal/views"
)

func init() {
	registerStep(phaseWrite, 32, "write llms.txt", writeLlmsTxt)
}

// The Agent Surface (ADR-0011): llms.txt orients an Agent in one short
// page (llmstxt.org: H1, a blockquote summary, H2 link sections) and never
// lists Pieces; sitemap.md enumerates them.
func writeLlmsTxt(s *site) error {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", views.PublicationName)
	fmt.Fprintf(&b, "> %s Tech news and architecture Explainers for mid-to-senior backend engineers, every claim linked to its first-party Source.\n\n", views.Tagline)
	fmt.Fprintf(&b, "Every Piece has a Markdown twin: replace the URL's trailing slash with `.md` (%s -> %s).\n\n",
		s.absURL(briefsIndexURL+"x/"), s.absURL(twinURL(briefsIndexURL+"x/")))
	b.WriteString("## Verticals\n\n")
	for _, v := range s.sortedVerticals() {
		fmt.Fprintf(&b, "- [%s](%s): %s\n", v.Name, s.absURL(verticalURL(v.Slug)), v.Description)
	}
	b.WriteString("\n## Briefs\n\n")
	fmt.Fprintf(&b, "- [Briefs](%s): short, dated commentary on something that just happened, with the architectural \"why it matters\", linking into Explainers.\n", s.absURL(briefsIndexURL))
	b.WriteString("\n## Indexes and feeds\n\n")
	fmt.Fprintf(&b, "- [Every live Piece, with its Markdown twin](%s)\n", s.absURL(markdownSitemapURL))
	fmt.Fprintf(&b, "- [RSS, all Pieces](%s)\n", s.absURL(siteFeedURL))
	fmt.Fprintf(&b, "- [RSS, Briefs](%s)\n", s.absURL(briefsFeedURL))
	fmt.Fprintf(&b, "- [XML sitemap](%s)\n", s.absURL(sitemapURL))
	b.WriteString("\n## Optional\n\n")
	fmt.Fprintf(&b, "- [About](%s)\n", s.absURL(aboutURL))
	fmt.Fprintf(&b, "- [Editorial Policy](%s)\n", s.absURL(editorialPolicyURL))
	return s.out.WriteStatic(llmsTxtURL, []byte(b.String()))
}

package builder

import (
	"maps"
	"regexp"
	"slices"
	"strings"

	"hastcode.com/site/internal/content"
)

func init() {
	registerStep(phaseValidate, 13, "validate no AI note", validateNoAINote)
}

// ruleAINote: the Author decided (ADR-0013) that no Reader sees an
// AI-assistance note or AI-assistance wording anywhere on the site.
const ruleAINote = "ai-note"

// aiNoteWording matches the disclosure's wording, never the word "AI" on
// its own: Pieces about AI as a topic ("AI-assisted development", "AI
// agents", "an AI assistant") must not trip it, and neither does the old
// note's second sentence ("We publish only as many pieces as the Author
// can verify this way"), which is true of the process and says nothing
// about AI. Any whitespace, line breaks and blockquote markers included,
// may sit between the words.
var aiNoteWording = regexp.MustCompile(`(?i)\bAI[\s>-]+assistance\b|\bdrafted[\s>]+with[\s>]+AI\b`)

// readerText is one Reader-visible field of a content file.
type readerText struct{ where, text string }

// validateNoAINote fails every Piece, static page and Author file whose
// Reader-visible text holds AI-assistance wording: the frontmatter fields
// a page, feed, llms.txt or structured data shows, and the whole body
// source, which the Markdown twin carries verbatim (HTML comments and
// code included). Frontmatter YAML comments are never published and are
// not read.
func validateNoAINote(s *site) error {
	for _, e := range s.explainers {
		texts := pieceTexts(e.Title, e.Description, e.Tags, e.HeroAlt, e.Sources, e.Body.Source)
		texts = append(texts, readerText{"who_for", e.WhoFor})
		checkNoAINote(s, e.File, texts)
	}
	for _, b := range s.briefs {
		checkNoAINote(s, b.File, pieceTexts(b.Title, b.Description, b.Tags, b.HeroAlt, b.Sources, b.Body.Source))
	}
	for _, p := range s.pages {
		checkNoAINote(s, p.File, []readerText{{"title", p.Title}, {"description", p.Description}, {"body", string(p.Body.Source)}})
	}
	for _, slug := range slices.Sorted(maps.Keys(s.authors)) {
		a := s.authors[slug]
		checkNoAINote(s, a.File, []readerText{{"name", a.Name}, {"photo_alt", a.PhotoAlt}, {"bio", string(a.Bio.Source)}})
	}
	return nil
}

func pieceTexts(title, description string, tags []string, heroAlt string, sources []content.Source, body []byte) []readerText {
	texts := []readerText{{"title", title}, {"description", description}, {"tags", strings.Join(tags, " ")}, {"hero_alt", heroAlt}}
	for _, src := range sources {
		texts = append(texts, readerText{"sources", src.Title})
	}
	return append(texts, readerText{"body", string(body)})
}

func checkNoAINote(s *site, file string, texts []readerText) {
	for _, t := range texts {
		if m := aiNoteWording.FindString(t.text); m != "" {
			s.fail(file, ruleAINote, "%s carries AI-assistance wording (%q): no AI note is published (ADR-0013); remove it", t.where, strings.Join(strings.Fields(strings.ReplaceAll(m, ">", " ")), " "))
		}
	}
}

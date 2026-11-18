package builder

import (
	"maps"
	"regexp"
	"slices"
	"strings"
)

func init() {
	registerStep(phaseValidate, 13, "validate no AI note", validateNoAINote)
}

// ruleAINote: the Author decided (ADR-0014) that no Reader sees an
// AI-assistance note, or any wording that says AI made this Piece or the
// site's Pieces, anywhere on the site.
const ruleAINote = "ai-note"

// The ai-note wording, built from parts so one change reaches every
// pattern. Between words it allows any space (Unicode ones too: Go's \s is
// ASCII only), line breaks, blockquote markers and emphasis; around "AI"
// in a compound it also allows any dash (hyphen, non-breaking hyphen, en
// dash). The match is case-insensitive.
const (
	aiSep  = `[\s\p{Z}\x{200B}-\x{200D}\x{2060}\x{FEFF}>*_]+`
	aiDash = `[\s\p{Z}\p{Pd}\x{200B}-\x{200D}\x{2060}\x{FEFF}>*_]+`
	// aiTerm names the AI itself: AI, a model family or an LLM.
	aiTerm = `(?:\bAI\b|\bClaude\b|\bChatGPT\b|\bGPT[\w.\-]*|\bGemini\b|\bCopilot\b|\bLLMs?\b|\b(?:large` + aiSep + `)?language` + aiSep + `models?\b|\bartificial` + aiSep + `intelligence\b|\bgenerative` + aiSep + `AI\b)`
	// aiHelp is "the help of", "assistance from", "an", "the" before aiTerm.
	aiHelp = `(?:(?:the` + aiSep + `)?(?:help|assistance|support)` + aiSep + `(?:of|from)` + aiSep + `)?(?:(?:an?|the)` + aiSep + `)?`
	aiWith = `(?:with|by|using|via|through)` + aiSep + aiHelp + aiTerm
	// aiMade is a participle of making a text.
	aiMade = `(?:drafted|written|generated|produced|created|composed|prepared|authored|assisted|edited|researched|summari[sz]ed|translated|made)`
	// aiMaking is any form of a making verb ("produces", "wrote").
	aiMaking = `\b(?:draft|writ|wrote|generat|produc|creat|compos|prepar|author|assist|edit|research|summari[sz]|translat|made|make)\w*`
	aiLead   = `^(?:[\s\p{Z}>*_#(\["'“‘]|<!--|-` + aiSep + `)*`
)

var (
	// aiNoteAlways is the note's own wording: a disclosure wherever it sits.
	aiNoteAlways = regexp.MustCompile(`(?i)\bdrafted` + aiSep + aiWith + `|\bAI` + aiDash + `(?:assistance|use|disclosure)` + aiSep + `(?:note|notice|disclosure|statement|label)\b`)
	// aiNoteOpening is a sentence that opens with how it was made
	// ("Generated with Claude.", "Written using an LLM, ..."), or is only
	// an "AI-generated" label.
	aiNoteOpening = regexp.MustCompile(`(?i)` + aiLead + `(?:` + aiMade + aiSep + aiWith + `|` + aiTerm + aiDash + aiMade + `(?:` + aiSep + `(?:content|text|article|piece|draft|summary))?[\s\p{Z}>*_.!)]*$)`)
	// aiMaking is AI making text: "written with the help of AI",
	// "AI-generated", "AI assistance". It counts only in a sentence that
	// is about this Piece or the site's own work (aiSelf).
	aiNoteMaking = regexp.MustCompile(`(?i)` + aiMaking + `(?:` + aiSep + `\S+){0,4}?` + aiSep + aiWith + `|` + aiTerm + aiDash + `(?:` + aiMade + `|assistance|help)\b|\bwith` + aiSep + aiTerm + aiDash + `(?:assistance|help)\b`)
	// aiSelf points at this Piece or at the site's own Pieces and team.
	aiSelf = regexp.MustCompile(`(?i)\b(?:this|these)` + aiSep + `(?:piece|article|post|page|brief|explainer|profile|repo` + aiSep + `profile|case` + aiSep + `study|concept|comparison|summary|story|text|draft|site|write-?up|content)s?\b|\bour` + aiSep + `(?:pieces|articles|posts|briefs|explainers|profiles|drafts|content|writing|work|site|editorial)\b|\bHastCode\b|\beditorial` + aiSep + `team\b|\bMira` + aiSep + `Vance\b|\bhouse` + aiSep + `byline\b|\bthis` + aiSep + `(?:name|byline)\b|` + aiLead + `(?:drafts|pieces|articles|briefs|explainers|posts)\b`)
	// aiSentenceEnd splits text into sentences: at . ! ? before a space,
	// and at a blank line (a blockquote's ">" lines count as one block).
	aiSentenceEnd = regexp.MustCompile(`[.!?]+["'”’)\]*_]*[\s\p{Z}]+|\n[ \t>]*\n`)
)

// aiNoteWording returns the AI-note wording in text, or "". It never
// matches AI as a topic ("AI-assisted development", "Copilot's AI
// assistance in code review", "Was pgrust written by AI?"): only the
// note's own wording, a sentence that opens with how it was made, or a
// sentence about this Piece or the site's own work that says AI made it.
func aiNoteWording(text string) string {
	if m := aiNoteAlways.FindString(text); m != "" {
		return m
	}
	for _, sentence := range aiSentenceEnd.Split(text, -1) {
		if aiNoteOpening.MatchString(sentence) {
			return sentence
		}
		if m := aiNoteMaking.FindString(sentence); m != "" && aiSelf.MatchString(sentence) {
			return m
		}
	}
	return ""
}

// readerText is one Reader-visible field of a content file.
type readerText struct{ where, text string }

// pieceFields are the Reader-visible fields an Explainer and a Brief
// share. Source titles are left out: they are other publishers' titles,
// copied as they are, so they describe that publisher's work, not ours.
type pieceFields struct {
	title, description, heroAlt string
	tags                        []string
	body                        []byte
}

func (p pieceFields) texts() []readerText {
	return []readerText{{"title", p.title}, {"description", p.description}, {"tags", strings.Join(p.tags, " ")}, {"hero_alt", p.heroAlt}, {"body", string(p.body)}}
}

// validateNoAINote fails every Piece, static page and Author file whose
// Reader-visible text holds AI-note wording: the frontmatter fields a
// page, feed, llms.txt or structured data shows, and the whole body
// source, which the Markdown twin carries verbatim (HTML comments and
// code included). Frontmatter YAML comments are never published and are
// not read. site.yaml's Vertical names and descriptions are read too.
func validateNoAINote(s *site) error {
	for _, e := range s.explainers {
		texts := pieceFields{e.Title, e.Description, e.HeroAlt, e.Tags, e.Body.Source}.texts()
		checkNoAINote(s, e.File, append(texts, readerText{"who_for", e.WhoFor}))
	}
	for _, b := range s.briefs {
		checkNoAINote(s, b.File, pieceFields{b.Title, b.Description, b.HeroAlt, b.Tags, b.Body.Source}.texts())
	}
	for _, p := range s.pages {
		checkNoAINote(s, p.File, []readerText{{"title", p.Title}, {"description", p.Description}, {"body", string(p.Body.Source)}})
	}
	for _, slug := range slices.Sorted(maps.Keys(s.authors)) {
		a := s.authors[slug]
		checkNoAINote(s, a.File, []readerText{{"name", a.Name}, {"photo_alt", a.PhotoAlt}, {"bio", string(a.Bio.Source)}})
	}
	for _, v := range s.cfg.Verticals {
		checkNoAINote(s, "", []readerText{
			{"site.yaml Vertical " + v.Slug + " name", v.Name},
			{"site.yaml Vertical " + v.Slug + " description", v.Description},
		})
	}
	return nil
}

func checkNoAINote(s *site, file string, texts []readerText) {
	for _, t := range texts {
		if m := aiNoteWording(t.text); m != "" {
			s.fail(file, ruleAINote, "%s carries AI-assistance wording (%q): no AI note is published (ADR-0014); remove it", t.where, strings.Join(strings.Fields(strings.ReplaceAll(m, ">", " ")), " "))
		}
	}
}

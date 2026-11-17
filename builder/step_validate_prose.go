package builder

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/yuin/goldmark/ast"

	"hastcode.com/site/internal/markdown"
)

func init() {
	registerStep(phaseValidate, 41, "warn on prose (sentence length, banned words, question headings)", warnProse)
}

const (
	ruleSentenceLength = "sentence-length"
	ruleBannedWord     = "banned-word"
	// ruleQuestionHeading is an H2 ending in "?". H3s are never read, so
	// a Repo Profile's FAQ questions (### under ## FAQ) never warn.
	ruleQuestionHeading = "question-heading"
)

// banned is house-byline-voice.md rule 5's banned list, lower case, as
// whole words. "just" is banned only as filler: see fillerJust.
var banned = []string{
	"revolutionary", "game-changing", "magic", "simply", "easy", "obviously",
	"it is important to note", "please note",
}

// fillerJust are the verbs that make "just" filler ("just add", "just
// run"): the instruction that hides the work. "just" before anything else
// ("just as", "just under", "just before", "just in time", "just
// pushed") means exactly, only or recently, and stays.
var fillerJust = map[string]bool{
	"add": true, "apply": true, "ask": true, "build": true, "bump": true, "call": true,
	"change": true, "check": true, "click": true, "configure": true, "copy": true,
	"create": true, "delete": true, "deploy": true, "disable": true, "do": true,
	"drop": true, "edit": true, "enable": true, "flip": true, "follow": true,
	"import": true, "increase": true, "install": true, "load": true, "make": true,
	"move": true, "open": true, "paste": true, "pick": true, "plug": true,
	"point": true, "pull": true, "push": true, "put": true, "read": true,
	"rebuild": true, "remove": true, "rename": true, "replace": true, "rerun": true,
	"restart": true, "retry": true, "rewrite": true, "run": true, "scale": true,
	"send": true, "set": true, "ship": true, "skip": true, "start": true,
	"stop": true, "swap": true, "switch": true, "toggle": true, "try": true,
	"turn": true, "type": true, "update": true, "upgrade": true, "use": true,
	"wrap": true, "write": true,
}

// maxSentenceWords is house-byline-voice.md rule 5: split anything over
// 30 words unless it is a deliberate list.
const maxSentenceWords = 30

// excerptWords is how many of a sentence's first words a warning quotes,
// enough for the Author to find it.
const excerptWords = 6

// proseSection is one stretch of a Piece's body the prose rules read: the
// opening before the first H2, or one H2 section.
type proseSection struct {
	where  string // "the opening" or "## Heading"
	blocks []ast.Node
}

func pieceSections(doc *markdown.Document) []proseSection {
	intro, sections := doc.Sections()
	out := []proseSection{{where: "the opening", blocks: intro}}
	for _, sec := range sections {
		out = append(out, proseSection{where: "## " + sec.Heading, blocks: sec.Blocks})
	}
	return out
}

// warnProse warns (never fails) on every Explainer's and Brief's prose
// that breaks the voice rules' numbers. Frontmatter, code, tables,
// blockquotes and headings are never read as prose (markdown.Prose). A
// block still holding a TODO(author) hint is a skeleton's placeholder,
// not prose: author-placeholder already reports it.
func warnProse(s *site) error {
	check := func(file string, doc *markdown.Document) {
		for _, sec := range pieceSections(doc) {
			if strings.HasSuffix(sec.where, "?") {
				s.warn(file, ruleQuestionHeading, `%s is a question; only Repo Profile FAQ headings (### under ## FAQ) may be questions: restate it as the answer`, sec.where)
			}
			for _, p := range doc.Prose(sec.blocks) {
				if strings.Contains(p.Text, authorPlaceholder) {
					continue
				}
				for _, sn := range sentences(p) {
					if sn.words > maxSentenceWords {
						s.warn(file, ruleSentenceLength, `%s: "%s" is %d words, over the %d-word limit; split it (a deliberate list may stay long)`, sec.where, sn.excerpt, sn.words, maxSentenceWords)
					}
				}
				for _, w := range bannedWords(p) {
					s.warn(file, ruleBannedWord, `%s: "%s" is on the banned list (house-byline-voice.md rule 5); cut it`, sec.where, w)
				}
			}
		}
	}
	for _, e := range s.explainers {
		check(e.File, e.Body)
	}
	for _, b := range s.briefs {
		check(b.File, b.Body)
	}
	return nil
}

type sentence struct {
	excerpt string
	words   int
}

// sentences splits a prose block at sentence ends: a full stop, "!" or "?"
// (outside inline code) with any closing quotes or brackets, then space
// and a capital, digit, quote, bracket or code, or the end of the block.
// "e.g.", "i.e.", "vs." and initials do not end a sentence.
func sentences(p markdown.Prose) []sentence {
	var out []sentence
	text, literal := p.Text, p.Literal
	start := 0
	flush := func(end int) {
		if w := countWords(text[start:end], literal[start:end]); w > 0 {
			out = append(out, sentence{excerpt: excerpt(text[start:end]), words: w})
		}
		start = end
	}
	for i := 0; i < len(text); i++ {
		if literal[i] || !strings.ContainsRune(".!?", rune(text[i])) {
			continue
		}
		j := i + 1
		for j < len(text) && !literal[j] && strings.ContainsRune(`"')]*_`, rune(text[j])) {
			j++
		}
		for strings.HasPrefix(text[j:], "”") || strings.HasPrefix(text[j:], "’") {
			j += len("”")
		}
		if j == len(text) {
			break
		}
		if text[j] != ' ' {
			continue
		}
		k := j
		for k < len(text) && text[k] == ' ' {
			k++
		}
		if k == len(text) {
			break
		}
		if text[i] == '.' && abbreviation(text[:i]) {
			continue
		}
		if r, _ := utf8.DecodeRuneInString(text[k:]); literal[k] || unicode.IsUpper(r) || unicode.IsDigit(r) || strings.ContainsRune(`"'(“‘[*_`, r) {
			flush(j)
		}
	}
	flush(len(text))
	return out
}

// abbreviation reports whether the word ending text (before its full
// stop) is one that does not end a sentence.
func abbreviation(text string) bool {
	word := text[strings.LastIndexAny(text, " (")+1:]
	switch strings.ToLower(word) {
	case "e.g", "i.e", "vs", "cf", "approx":
		return true
	}
	r, size := utf8.DecodeRuneInString(word)
	return size == len(word) && unicode.IsUpper(r) // an initial
}

// countWords counts the words in a sentence: space-separated tokens with a
// letter or digit in them, an inline code span counting as one word.
func countWords(text string, literal []bool) int {
	n, inWord, counted := 0, false, false
	for i, r := range text {
		if r == ' ' && !literal[i] {
			inWord = false
			continue
		}
		if !inWord {
			inWord, counted = true, false
		}
		if !counted && (literal[i] || unicode.IsLetter(r) || unicode.IsDigit(r)) {
			n++
			counted = true
		}
	}
	return n
}

// bannedWords returns each banned word or phrase in a prose block, in
// order, as it appears in the banned list ("just" with its verb). Words
// inside quotation marks are a Source's and never count, nor does inline
// code.
func bannedWords(p markdown.Prose) []string {
	skip := quoted(p)
	for i, lit := range p.Literal {
		skip[i] = skip[i] || lit
	}
	lower := []byte(p.Text)
	for i, c := range lower {
		if 'A' <= c && c <= 'Z' {
			lower[i] = c + 'a' - 'A'
		}
	}
	text := string(lower)
	type hit struct {
		at   int
		word string
	}
	var hits []hit
	find := func(phrase string, keep func(end int) (string, bool)) {
		for at := 0; ; {
			i := strings.Index(text[at:], phrase)
			if i < 0 {
				return
			}
			i += at
			end := i + len(phrase)
			at = end
			if wordByte(text, i-1) || wordByte(text, end) || slicesAny(skip[i:end]) {
				continue
			}
			if w, ok := keep(end); ok {
				hits = append(hits, hit{i, w})
			}
		}
	}
	for _, w := range banned {
		find(w, func(int) (string, bool) { return w, true })
	}
	find("just", func(end int) (string, bool) {
		rest := strings.TrimLeft(text[end:], " ")
		next := rest[:len(rest)-len(strings.TrimLeftFunc(rest, func(r rune) bool { return 'a' <= r && r <= 'z' }))]
		at := len(text) - len(rest)
		if !fillerJust[next] || wordByte(text, at+len(next)) || slicesAny(skip[at:at+len(next)]) {
			return "", false
		}
		return "just " + next, true
	})
	slices.SortStableFunc(hits, func(a, b hit) int { return cmp.Compare(a.at, b.at) })
	var out []string
	for _, h := range hits {
		out = append(out, h.word)
	}
	return out
}

// wordByte reports whether text[i] exists and is a letter, digit or
// underscore, so a match there is inside a longer word ("easy" in
// "uneasy"). A hyphen or dash is a boundary: "easy-to-use" warns.
func wordByte(text string, i int) bool {
	if i < 0 || i >= len(text) {
		return false
	}
	c := text[i]
	return c == '_' || ('a' <= c && c <= 'z') || ('0' <= c && c <= '9')
}

func slicesAny(bs []bool) bool { return slices.Contains(bs, true) }

// quoted marks the bytes of a prose block inside quotation marks: a pair
// of straight double quotes, or “curly” ones. An unmatched opening quote
// marks nothing.
func quoted(p markdown.Prose) []bool {
	mark := make([]bool, len(p.Text))
	open := -1
	closing := ""
	for i := 0; i < len(p.Text); i++ {
		if p.Literal[i] {
			continue
		}
		rest := p.Text[i:]
		switch {
		case open < 0 && strings.HasPrefix(rest, `"`):
			open, closing = i, `"`
		case open < 0 && strings.HasPrefix(rest, "“"):
			open, closing = i, "”"
		case open >= 0 && strings.HasPrefix(rest, closing):
			end := i + len(closing)
			for j := open; j < end; j++ {
				mark[j] = true
			}
			open = -1
			i = end - 1
		}
	}
	return mark
}

// excerpt is the sentence's first excerptWords words, with "…" when it
// runs on.
func excerpt(text string) string {
	words := strings.Fields(text)
	if len(words) <= excerptWords {
		return strings.Join(words, " ")
	}
	return fmt.Sprintf("%s…", strings.Join(words[:excerptWords], " "))
}

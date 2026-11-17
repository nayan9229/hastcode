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
// whole words; rule 5 is the list's one home in the docs, so change both
// together. "just" is banned only as filler (see fillerJust), and "magic"
// in a technical term is not hype (see technicalMagic).
var banned = []string{
	"revolutionary", "game-changing", "magic", "simply", "easy", "obviously",
	"it is important to note", "please note",
}

// fillerJust are the verbs that make "just" filler ("just add", "just
// run"): the instruction that hides the work. "just" before anything else
// ("just as", "just under", "just before", "just in time", "just
// pushed") means exactly, only or recently, and stays, and so does "just"
// after a form of "be" ("the job is just run once"), which is passive,
// not an instruction (beForms).
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

// beForms are the forms of "be" that make a following "just <verb>"
// passive ("is just run once") rather than an instruction.
var beForms = map[string]bool{
	"am": true, "is": true, "are": true, "was": true, "were": true,
	"be": true, "been": true, "being": true,
}

// technicalMagic are the words after "magic" that make it a technical
// term ("magic number", "magic bytes"), not hype.
var technicalMagic = map[string]bool{
	"number": true, "numbers": true, "byte": true, "bytes": true,
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
	where   string // for messages: "the opening" or "## Heading"
	heading string // the H2's text; empty for the opening
	blocks  []ast.Node
}

func pieceSections(doc *markdown.Document) []proseSection {
	intro, sections := doc.Sections()
	out := []proseSection{{where: "the opening", blocks: intro}}
	for _, sec := range sections {
		out = append(out, proseSection{where: "## " + sec.Heading, heading: sec.Heading, blocks: sec.Blocks})
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
			if strings.HasSuffix(sec.heading, "?") {
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
// and a capital, digit, quote, bracket, code or a lower-case brand with a
// capital in it ("iOS", "gRPC"), or the end of the block. "e.g.", "i.e.",
// "vs.", "Fig." and initials do not end a sentence, nor does "No." before
// a number ("No. 5"). Known limits, which only move a warning: a sentence
// ending on one capital letter ("Plan B.") is read as an initial, and one
// starting with an all lower-case word ("npm", "etcd") is not split off.
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
		if text[i] == '.' && abbreviation(text[:i], text[k:]) {
			continue
		}
		if r, _ := utf8.DecodeRuneInString(text[k:]); literal[k] || unicode.IsUpper(r) || unicode.IsDigit(r) || strings.ContainsRune(`"'(“‘[*_`, r) || lowerBrand(text[k:]) {
			flush(j)
		}
	}
	flush(len(text))
	return out
}

// abbreviation reports whether the word ending text (before its full
// stop) is one that does not end a sentence; next is the text after the
// space.
func abbreviation(text, next string) bool {
	word := text[strings.LastIndexAny(text, " (")+1:]
	switch strings.ToLower(word) {
	case "e.g", "i.e", "vs", "cf", "approx", "fig":
		return true
	case "no", "nos":
		r, _ := utf8.DecodeRuneInString(next)
		return unicode.IsDigit(r)
	}
	r, size := utf8.DecodeRuneInString(word)
	return size == len(word) && unicode.IsUpper(r) // an initial
}

// lowerBrand reports whether text starts with a word that begins lower
// case and has a capital later ("iOS", "gRPC", "eBPF"): a name that can
// open a sentence.
func lowerBrand(text string) bool {
	word, _, _ := strings.Cut(text, " ")
	r, _ := utf8.DecodeRuneInString(word)
	return unicode.IsLower(r) && strings.IndexFunc(word, unicode.IsUpper) > 0
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
	// nextWord returns the lower-case word starting at or after at, and
	// where it starts.
	nextWord := func(at int) (string, int) {
		rest := strings.TrimLeft(text[at:], " ")
		start := len(text) - len(rest)
		end := start + len(rest) - len(strings.TrimLeftFunc(rest, func(r rune) bool { return 'a' <= r && r <= 'z' }))
		return text[start:end], start
	}
	find := func(phrase string, keep func(start, end int) (string, bool)) {
		for at := 0; ; {
			i := strings.Index(text[at:], phrase)
			if i < 0 {
				return
			}
			i += at
			end := i + len(phrase)
			at = end
			if wordByte(text, i-1) || wordByte(text, end) || slices.Contains(skip[i:end], true) {
				continue
			}
			if w, ok := keep(i, end); ok {
				hits = append(hits, hit{i, w})
			}
		}
	}
	for _, w := range banned {
		find(w, func(_, end int) (string, bool) {
			if next, _ := nextWord(end); w == "magic" && technicalMagic[next] {
				return "", false
			}
			return w, true
		})
	}
	find("just", func(start, end int) (string, bool) {
		next, at := nextWord(end)
		if !fillerJust[next] || wordByte(text, at+len(next)) || slices.Contains(skip[at:at+len(next)], true) {
			return "", false
		}
		before := strings.Fields(text[:start])
		if len(before) > 0 && beForms[before[len(before)-1]] {
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

// quoted marks the bytes of a prose block inside quotation marks: a pair
// of straight double quotes, “curly” ones or ‘curly singles’. A ’ followed
// by a letter is an apostrophe ("it’s"), never a closing quote, and a
// straight ' is never read as a quote. An unmatched opening quote marks
// nothing.
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
		case open < 0 && strings.HasPrefix(rest, "‘"):
			open, closing = i, "’"
		case open >= 0 && strings.HasPrefix(rest, closing) && !apostrophe(rest, closing):
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

// apostrophe reports whether the ’ starting rest is an apostrophe: a
// letter follows it.
func apostrophe(rest, closing string) bool {
	if closing != "’" {
		return false
	}
	r, _ := utf8.DecodeRuneInString(rest[len(closing):])
	return unicode.IsLetter(r)
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

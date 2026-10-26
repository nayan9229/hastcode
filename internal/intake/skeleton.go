package intake

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// skeleton fills the template tmpl for p.
func (p pick) skeleton(tmpl string) ([]byte, error) {
	rest, ok := strings.CutPrefix(tmpl, "---\n")
	fm, body, ok2 := strings.Cut(rest, "\n---\n")
	if !ok || !ok2 {
		return nil, errors.New("no frontmatter block")
	}
	frontmatter, err := p.frontmatter(fm)
	if err != nil {
		return nil, err
	}
	if p.isBrief() {
		body, err = p.briefBody(body)
	} else {
		body, err = p.explainerBody(body)
	}
	if err != nil {
		return nil, err
	}
	return []byte("---\n" + frontmatter + "\n---\n" + body), nil
}

// frontmatter is the template's frontmatter with its opening comment
// replaced by the skeleton's own and the Lead's values set; every other
// line, TODO(author) values and commented options included, is kept.
func (p pick) frontmatter(fm string) (string, error) {
	published := quote(p.now.Format(time.RFC3339))
	values := map[string]string{
		"title":     quote(p.lead.Title),
		"published": published,
	}
	dates := "published is the pick time, provisional: set it"
	if p.isBrief() {
		values["source_url"] = quote(p.lead.Link)
		values["explainers"] = "[" + strings.Join(p.explainers, ", ") + "]"
	} else {
		values["vertical"] = p.vertical
		values["reviewed"] = published
		dates = "published and reviewed are the pick time, provisional: set them"
	}
	lines := []string{
		fmt.Sprintf("# Skeleton from Lead %s (%s), picked %s.", p.lead.ID, p.lead.Source, p.now.Format(time.DateOnly)),
		fmt.Sprintf("# Follow docs/editorial/%s-writing-guide.md and replace every", p.kind),
		"# TODO(author) line. " + dates,
		"# when you publish. verified: is the Author's alone (ADR-0007).",
	}
	opening := true
	set := map[string]bool{}
	for _, line := range strings.Split(fm, "\n") {
		comment := strings.HasPrefix(line, "#")
		if opening && comment {
			continue // the template's own note on how to use it
		}
		opening = false
		key, _, isKey := strings.Cut(line, ":")
		if comment || !isKey {
			lines = append(lines, line)
			continue
		}
		switch value, known := values[key]; {
		case known:
			line = key + ": " + value
			set[key] = true
		case key == "kind" && strings.TrimSpace(line) != "kind: "+p.kind:
			return "", fmt.Errorf("its kind line %q is not kind: %s", line, p.kind)
		case key == "verified" || key == "hero" || key == "hero_alt":
			return "", fmt.Errorf("it sets %s, which a skeleton never does", key)
		}
		lines = append(lines, line)
	}
	for key := range values {
		if !set[key] {
			return "", fmt.Errorf("it has no %s line to fill", key)
		}
	}
	return strings.Join(lines, "\n"), nil
}

// explainerBody is the template body with ## Sources holding only the
// Lead's link.
func (p pick) explainerBody(body string) (string, error) {
	before, sources, ok := strings.Cut(body, "\n## Sources\n")
	if !ok {
		return "", errors.New("no ## Sources section")
	}
	if strings.Contains(sources, "\n## ") {
		return "", errors.New("## Sources is not the last section")
	}
	lead := fmt.Sprintf("- [%s](%s), %s\n", escapeLinkText(p.lead.Title), p.lead.Link, p.lead.Source)
	return before + "\n## Sources\n\n" + lead, nil
}

// explainerToken is where the Brief template links its Explainer.
const explainerToken = "{explainer}"

// briefBody is the template body with "Why it matters" linking each of
// the Brief's Explainers: the template's link names the first, and a
// TODO(author) line links each other one.
func (p pick) briefBody(body string) (string, error) {
	if strings.Count(body, explainerToken) != 1 {
		return "", fmt.Errorf("want one %s link in the body", explainerToken)
	}
	body = strings.Replace(body, explainerToken, p.explainers[0], 1)
	var extra strings.Builder
	for _, slug := range p.explainers[1:] {
		fmt.Fprintf(&extra, "\nTODO(author): link [the idea](%s) where the analysis uses it.\n", slug)
	}
	return strings.TrimRight(body, "\n") + "\n" + extra.String(), nil
}

// writeNewFile creates file with data, refusing one that exists.
func writeNewFile(file string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if errors.Is(err, fs.ErrExist) {
		return fmt.Errorf("intake: %s exists; a pick never overwrites a Piece", file)
	}
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		os.Remove(file)
		return err
	}
	return f.Close()
}

// quote is s as a YAML double-quoted scalar (JSON's string syntax is a
// subset of it).
func quote(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

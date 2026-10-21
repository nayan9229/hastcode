// Package mermaid renders Mermaid diagram sources to SVG ready to inline in
// a page. The renderer is beautiful-mermaid (npm, pinned in the repo's
// package.json; pure JS, no browser), run as one Node process per batch.
// Results are cached on disk by a hash of the source, so unchanged diagrams
// never reach Node.
//
// The SVG is post-processed: the renderer's <style> block (a Google Fonts
// @import plus page-global selectors) is removed — the site stylesheet
// carries those rules once, scoped to .diagram — and data-* attributes are
// dropped. Colours stay CSS variables and currentColor, so the page's CSS
// decides light and dark mode.
package mermaid

import (
	"bytes"
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

// Package is the npm package (and version pin lives in package.json).
const Package = "beautiful-mermaid"

// cacheVersion changes whenever render options or post-processing change,
// invalidating every cached SVG.
const cacheVersion = "1"

//go:embed render.mjs
var script string

// Supported lists the Mermaid diagram types the renderer handles with
// theme-able output. Others (gantt, pie, xychart-beta…) are rejected.
var Supported = []string{"flowchart", "graph", "sequenceDiagram", "stateDiagram", "stateDiagram-v2", "classDiagram", "erDiagram"}

// Renderer renders batches of Mermaid sources.
type Renderer struct {
	// Node is the node executable ("node" if empty).
	Node string
	// ToolsDir holds node_modules/beautiful-mermaid; Node runs there.
	ToolsDir string
	// CacheDir is where rendered SVGs are kept between builds ("" disables
	// caching).
	CacheDir string
}

// Result is one diagram's outcome: SVG, or Err if the source is not a
// diagram this renderer can draw.
type Result struct {
	SVG string
	Err error
}

// Render renders every source, in order. Problems with a single diagram
// are in its Result.Err; the returned error is for the environment (Node
// or the npm package missing, the renderer crashing).
func (r Renderer) Render(ctx context.Context, sources []string) ([]Result, error) {
	results := make([]Result, len(sources))
	version := r.packageVersion()
	var missIdx []int
	var missSrc []string
	for i, src := range sources {
		if err := check(src); err != nil {
			results[i].Err = err
			continue
		}
		if svg, ok := r.cached(version, src); ok {
			results[i].SVG = svg
			continue
		}
		missIdx = append(missIdx, i)
		missSrc = append(missSrc, src)
	}
	if len(missSrc) == 0 {
		return results, nil
	}
	rendered, err := r.run(ctx, missSrc)
	if err != nil {
		return nil, err
	}
	for j, i := range missIdx {
		res := rendered[j]
		if res.Error != "" {
			results[i].Err = errors.New(res.Error)
			continue
		}
		svg := postProcess(res.SVG)
		if !strings.Contains(svg, "<text") {
			results[i].Err = errors.New("diagram rendered with no nodes; check the Mermaid syntax")
			continue
		}
		results[i].SVG = svg
		r.store(version, sources[i], svg)
	}
	return results, nil
}

// check rejects, before they reach the renderer, empty sources, diagram
// types outside Supported, and flowchart lines the renderer would silently
// drop (it is lenient: "A[[[x --> B" renders A alone).
func check(src string) error {
	lines := strings.Split(src, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "%%") {
			continue
		}
		typ := strings.Fields(line)[0]
		if !slices.Contains(Supported, typ) {
			return fmt.Errorf("diagram type %q is not supported (use one of %s)", typ, strings.Join(Supported, ", "))
		}
		if typ == "flowchart" || typ == "graph" {
			return checkFlowchart(lines[i+1:])
		}
		return nil
	}
	return errors.New("diagram is empty")
}

var (
	// Text that may hold brackets legitimately: "quoted labels", |edge
	// labels|, and asymmetric shapes id>label].
	quoted       = regexp.MustCompile(`"[^"]*"|\|[^|]*\||\w>[^\]]*\]`)
	danglingEdge = regexp.MustCompile(`(-->|---|==>|===|\.->|-\.-|~~~)$`)
	closer       = map[rune]rune{')': '(', ']': '[', '}': '{'}
)

func checkFlowchart(lines []string) error {
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "%%") {
			continue
		}
		bare := quoted.ReplaceAllString(line, `""`)
		if danglingEdge.MatchString(bare) {
			return fmt.Errorf("edge with no target node in %q", line)
		}
		var open []rune
		for _, c := range bare {
			switch c {
			case '(', '[', '{':
				open = append(open, c)
			case ')', ']', '}':
				if len(open) == 0 || open[len(open)-1] != closer[c] {
					return fmt.Errorf("unbalanced brackets in %q", line)
				}
				open = open[:len(open)-1]
			}
		}
		if len(open) > 0 {
			return fmt.Errorf("unbalanced brackets in %q", line)
		}
	}
	return nil
}

type rendered struct {
	SVG   string `json:"svg"`
	Error string `json:"error"`
}

func (r Renderer) run(ctx context.Context, sources []string) ([]rendered, error) {
	node := r.Node
	if node == "" {
		node = "node"
	}
	in, err := json.Marshal(sources)
	if err != nil {
		return nil, err
	}
	cmd := exec.CommandContext(ctx, node, "--input-type=module", "--eval", script)
	cmd.Dir = r.ToolsDir
	cmd.Stdin = bytes.NewReader(in)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("mermaid renderer (%s with npm package %s in %q; run `npm ci` there): %v\n%s",
			node, Package, r.ToolsDir, err, strings.TrimSpace(stderr.String()))
	}
	var out []rendered
	if err := json.Unmarshal(stdout.Bytes(), &out); err != nil || len(out) != len(sources) {
		return nil, fmt.Errorf("mermaid renderer: unexpected output (%v): %.200s", err, stdout.String())
	}
	return out, nil
}

var (
	styleBlock = regexp.MustCompile(`(?s)<style>.*?</style>\s*`)
	dataAttr   = regexp.MustCompile(` data-[a-z-]+="[^"]*"`)
)

func postProcess(svg string) string {
	svg = styleBlock.ReplaceAllString(svg, "")
	return dataAttr.ReplaceAllString(svg, "")
}

var (
	idAttr = regexp.MustCompile(`\bid="([^"]+)"`)
	idRef  = regexp.MustCompile(`url\(#([^)]+)\)|href="#([^"]+)"`)
)

// PrefixIDs prefixes every element id in svg, and every url(#…) and
// href="#…" reference to one, so several diagrams can share a page.
func PrefixIDs(svg, prefix string) string {
	svg = idAttr.ReplaceAllString(svg, `id="`+prefix+`$1"`)
	return idRef.ReplaceAllStringFunc(svg, func(m string) string {
		if strings.HasPrefix(m, "url(#") {
			return "url(#" + prefix + m[len("url(#"):]
		}
		return `href="#` + prefix + m[len(`href="#`):]
	})
}

// packageVersion is the installed renderer version, part of the cache key
// so an upgrade re-renders everything.
func (r Renderer) packageVersion() string {
	data, err := os.ReadFile(filepath.Join(r.ToolsDir, "node_modules", Package, "package.json"))
	if err != nil {
		return "unknown"
	}
	var pkg struct {
		Version string `json:"version"`
	}
	if json.Unmarshal(data, &pkg) != nil {
		return "unknown"
	}
	return pkg.Version
}

func (r Renderer) cachePath(version, src string) string {
	sum := sha256.Sum256([]byte(cacheVersion + "\x00" + version + "\x00" + src))
	return filepath.Join(r.CacheDir, hex.EncodeToString(sum[:])+".svg")
}

func (r Renderer) cached(version, src string) (string, bool) {
	if r.CacheDir == "" {
		return "", false
	}
	data, err := os.ReadFile(r.cachePath(version, src))
	return string(data), err == nil
}

// store writes a cache entry; failures only cost a re-render next time.
func (r Renderer) store(version, src, svg string) {
	if r.CacheDir == "" {
		return
	}
	if err := os.MkdirAll(r.CacheDir, 0o755); err != nil {
		return
	}
	path := r.cachePath(version, src)
	tmp, err := os.CreateTemp(r.CacheDir, ".tmp-*")
	if err != nil {
		return
	}
	_, werr := tmp.WriteString(svg)
	cerr := tmp.Close()
	if werr != nil || cerr != nil || os.Rename(tmp.Name(), path) != nil {
		_ = os.Remove(tmp.Name())
	}
}

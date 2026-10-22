package factpack

import (
	"bytes"
	"encoding/json"
	"os"

	"hastcode.com/site/internal/content"
)

// Marshal renders a Fact Pack as its file: keys in content.FactPack's
// declaration order, two-space indentation, no HTML escaping, every list
// present ([] rather than null) and the languages map present ({}), and a
// trailing newline. Equal facts give identical bytes, so a Refresh diff
// shows values, never shape.
func Marshal(fp *content.FactPack) ([]byte, error) {
	cp := *fp
	if cp.Topics == nil {
		cp.Topics = []string{}
	}
	if cp.Languages == nil {
		cp.Languages = map[string]int64{}
	}
	if cp.Releases == nil {
		cp.Releases = []content.Release{}
	}
	if cp.Contributors.Top == nil {
		cp.Contributors.Top = []content.Contributor{}
	}
	if cp.Files.CI == nil {
		cp.Files.CI = []string{}
	}
	if cp.Files.Manifests == nil {
		cp.Files.Manifests = []string{}
	}
	if cp.Sources == nil {
		cp.Sources = []string{}
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(&cp); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeFactPack(file string, fp *content.FactPack) error {
	data, err := Marshal(fp)
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0o644)
}

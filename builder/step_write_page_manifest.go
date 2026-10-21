package builder

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func init() {
	// Last: the manifest is the baseline for the next build's changed
	// URLs, so it is saved only once this build is known to have passed.
	registerStep(phaseWrite, 900, "save page manifest", savePageManifest)
}

// savePageManifest stores this build's manifest (built by
// indexWrittenPages, step_write_page_index.go) for the next build, but
// only when this build passed: a failed build is not deployed, so its
// changes must still count as changes next time.
func savePageManifest(s *site) error {
	if s.result.Failed() || s.index == nil {
		return nil
	}
	data, err := json.MarshalIndent(s.index.manifest, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.manifestPath()), 0o755); err != nil {
		return err
	}
	return os.WriteFile(s.manifestPath(), append(data, '\n'), 0o644)
}

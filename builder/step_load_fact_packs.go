package builder

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"hastcode.com/site/internal/content"
)

func init() {
	registerStep(phaseLoad, 22, "load fact packs", loadFactPacks)
}

// ruleFactPack is a Repo Profile whose Fact Pack is missing or does not
// parse; the Profile's page is not written.
const ruleFactPack = "fact-pack"

// loadFactPacks reads every Repo Profile's Fact Pack (the JSON named after
// its Markdown, beside it: content.FactPackFile) into s.factPacks. A
// Profile without a usable Fact Pack fails here, attributed to its Markdown
// file, with the message naming the JSON file.
func loadFactPacks(s *site) error {
	for _, e := range s.explainers {
		if e.Kind != content.RepoProfile {
			continue
		}
		file := content.FactPackFile(e.File)
		data, err := os.ReadFile(filepath.Join(s.contentDir, filepath.FromSlash(file)))
		switch {
		case errors.Is(err, fs.ErrNotExist):
			s.fail(e.File, ruleFactPack, "missing Fact Pack %s: a Repo Profile renders from a Fact Pack fetched beside its Markdown", file)
			continue
		case err != nil:
			return err
		}
		fp, err := content.ParseFactPack(data)
		if err != nil {
			s.fail(e.File, ruleFactPack, "Fact Pack %s: %v", file, err)
			continue
		}
		s.factPacks[e.File] = fp
	}
	return nil
}

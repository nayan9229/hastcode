package builder

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

func init() {
	registerStep(phaseWrite, 55, "write WebSub changed feeds", writeWebSub)
}

// WebSub: every feed names the hub (webSubHub, step_write_feeds.go) next
// to its rel="self" topic URL. The build never pings: it writes the feeds
// whose bytes changed since the previous successful build of the same
// mode (manifest Feeds, the manifest the IndexNow diff uses) to
// <out>/websub.json, beside indexnow.json and outside static/, and the
// Production publish job POSTs hub.mode=publish with one hub.url per topic
// once the deploy is live. README "Feeds, sitemaps, robots and IndexNow".
const webSubArtefact = "websub.json"

// webSubNotice is the artefact: the hub and the changed topics, sorted.
type webSubNotice struct {
	Hub    string   `json:"hub"`
	Topics []string `json:"topics"`
}

// writeWebSub writes the artefact when any feed is new or changed. With
// no changed feed there is no artefact, so the post-deploy step skips the
// ping.
func writeWebSub(s *site) error {
	artefact := filepath.Join(s.out.Dir, webSubArtefact)
	if err := os.Remove(artefact); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	if s.index == nil || len(s.index.changedFeeds) == 0 {
		return nil
	}
	topics := slices.Sorted(slices.Values(s.index.changedFeeds))
	data, err := json.MarshalIndent(webSubNotice{Hub: webSubHub, Topics: topics}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(artefact, append(data, '\n'), 0o644)
}

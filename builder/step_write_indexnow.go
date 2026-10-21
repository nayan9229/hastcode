package builder

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
)

func init() {
	registerStep(phaseValidate, 6, "validate IndexNow key", validateIndexNowKey)
	registerStep(phaseWrite, 50, "write IndexNow key and changed URLs", writeIndexNow)
}

// IndexNow (https://www.indexnow.org/documentation, checked 2026-10-21):
// the key file is /<key>.txt at the site root holding the key; a JSON POST
// {host, key, keyLocation, urlList} to an engine's /indexnow endpoint
// submits up to 10,000 URLs of that host. The build never pings: it writes
// the request body to <out>/indexnow.json (beside static/, so it is not
// served) and a post-deploy step sends it. See README "Feeds, sitemaps,
// robots and IndexNow".
const (
	ruleIndexNowKey     = "indexnow-key"
	indexNowArtefact    = "indexnow.json"
	indexNowMaxURLs     = 10000
	ruleIndexNowTooMany = "indexnow-urls"
)

var indexNowKeyPattern = regexp.MustCompile(`^[a-zA-Z0-9-]{8,128}$`)

func validateIndexNowKey(s *site) error {
	if k := s.cfg.IndexNowKey; k != "" && !indexNowKeyPattern.MatchString(k) {
		s.fail("", ruleIndexNowKey, "indexnow_key %q must be 8 to 128 characters of a-z, A-Z, 0-9 and -", k)
	}
	return nil
}

// indexNowRequest is the IndexNow JSON POST body.
type indexNowRequest struct {
	Host        string   `json:"host"`
	Key         string   `json:"key"`
	KeyLocation string   `json:"keyLocation"`
	URLList     []string `json:"urlList"`
}

// writeIndexNow writes the key file and, when any indexable page is new,
// changed or gone since the previous successful build (s.index.changed),
// the request body listing them. With nothing to submit there is no
// artefact, so the post-deploy step skips the ping.
func writeIndexNow(s *site) error {
	artefact := filepath.Join(s.out.Dir, indexNowArtefact)
	if err := os.Remove(artefact); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	key := s.cfg.IndexNowKey
	if key == "" || !indexNowKeyPattern.MatchString(key) {
		return nil
	}
	keyURL := "/" + key + ".txt"
	if err := s.out.WriteStatic(keyURL, []byte(key+"\n")); err != nil {
		return err
	}
	changed := s.index.changed
	if len(changed) == 0 {
		return nil
	}
	if len(changed) > indexNowMaxURLs {
		s.warn("", ruleIndexNowTooMany, "%d changed URLs; IndexNow takes %d per request, so only the first %d are listed", len(changed), indexNowMaxURLs, indexNowMaxURLs)
		changed = changed[:indexNowMaxURLs]
	}
	site, err := url.Parse(s.cfg.SiteURL)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(indexNowRequest{
		Host:        site.Host,
		Key:         key,
		KeyLocation: s.absURL(keyURL),
		URLList:     changed,
	}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(artefact, append(data, '\n'), 0o644)
}

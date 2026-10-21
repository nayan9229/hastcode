package builder

import "regexp"

func init() {
	registerStep(phaseValidate, 6, "validate IndexNow key", validateIndexNowKey)
}

const ruleIndexNowKey = "indexnow-key"

// indexNowKeyPattern is IndexNow's key format: 8 to 128 characters of
// a-z, A-Z, 0-9 and "-" (see step_write_indexnow.go).
var indexNowKeyPattern = regexp.MustCompile(`^[a-zA-Z0-9-]{8,128}$`)

func validateIndexNowKey(s *site) error {
	if k := s.cfg.IndexNowKey; k != "" && !indexNowKeyPattern.MatchString(k) {
		s.fail("", ruleIndexNowKey, "indexnow_key %q must be 8 to 128 characters of a-z, A-Z, 0-9 and -", k)
	}
	return nil
}

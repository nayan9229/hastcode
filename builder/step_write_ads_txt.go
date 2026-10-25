package builder

import "strings"

func init() {
	registerStep(phaseWrite, 41, "write ads.txt", writeAdsTxt)
}

// writeAdsTxt writes /ads.txt from site.yaml ads_txt, one record per line,
// when any are configured.
func writeAdsTxt(s *site) error {
	if len(s.cfg.AdsTxt) == 0 {
		return nil
	}
	return s.out.WriteStatic(adsTxtURL, []byte(strings.Join(s.cfg.AdsTxt, "\n")+"\n"))
}

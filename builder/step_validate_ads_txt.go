package builder

import "strings"

// /ads.txt records are checked before anything is written: a malformed line
// would quietly stop an ad system from being recognised as an authorised
// seller of the site's inventory.
func init() {
	registerStep(phaseValidate, 6, "validate ads.txt", validateAdsTxt)
}

const ruleAdsTxt = "ads-txt"

func validateAdsTxt(s *site) error {
	for _, line := range s.cfg.AdsTxt {
		if !validAdsTxtLine(line) {
			s.fail("", ruleAdsTxt, "site.yaml ads_txt record %q is not \"<ad system domain>, <account id>, DIRECT|RESELLER[, <cert id>]\"", line)
		}
	}
	return nil
}

// validAdsTxtLine accepts a data record, a "#" comment or a "key=value"
// variable line (IAB ads.txt 1.1).
func validAdsTxtLine(line string) bool {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return line != ""
	}
	if !strings.Contains(line, ",") && strings.Contains(line, "=") {
		return true
	}
	fields := strings.Split(line, ",")
	if len(fields) < 3 || len(fields) > 4 {
		return false
	}
	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
	}
	domain, account, relation := fields[0], fields[1], strings.ToUpper(fields[2])
	return strings.Contains(domain, ".") && account != "" && (relation == "DIRECT" || relation == "RESELLER")
}

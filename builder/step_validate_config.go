package builder

import (
	"regexp"
	"strings"

	"hastcode.com/site/internal/views"
)

func init() {
	registerStep(phaseValidate, 5, "ads config", checkAdsConfig)
}

const ruleAdsConfig = "ads-config"

// adUnitID is an AdSense ad unit id (data-ad-slot): digits only.
var adUnitID = regexp.MustCompile(`^[0-9]+$`)

// checkAdsConfig: turning ads on needs the AdSense account and an ad unit
// id for every slot position, or slots would render units that never fill.
// Preview builds carry no ads (adsOn), so they are not checked: a missing
// ad unit id fails the Production and QA Publish, never a Preview.
func checkAdsConfig(s *site) error {
	if !s.adsOn() {
		return nil
	}
	if s.cfg.AdSenseClient == "" {
		s.fail("", ruleAdsConfig, "ads_enabled is true but adsense_client is blank; set the AdSense publisher id (ca-pub-...)")
	}
	var missing []string
	for _, position := range views.AdPositions {
		id := s.cfg.AdSenseSlots[position]
		switch {
		case id == "":
			missing = append(missing, position)
		case !adUnitID.MatchString(id):
			s.fail("", ruleAdsConfig, "adsense_slots %s = %q is not an AdSense ad unit id (digits, the data-ad-slot value)", position, id)
		}
	}
	if len(missing) > 0 {
		s.fail("", ruleAdsConfig, "ads_enabled is true but adsense_slots has no ad unit id for %s; create a display ad unit per position in AdSense (Ads > By ad unit) and copy its data-ad-slot", strings.Join(missing, ", "))
	}
	return nil
}

package builder

import (
	"regexp"
	"strings"

	"hastcode.com/site/internal/views"
)

func init() {
	registerStep(phaseValidate, 5, "ads config", checkAdsConfig)
	registerStep(phaseValidate, 6, "analytics config", checkAnalyticsConfig)
}

const (
	ruleAdsConfig       = "ads-config"
	ruleAnalyticsConfig = "analytics-config"
)

// gaMeasurementID is a Google Analytics 4 web data stream id.
var gaMeasurementID = regexp.MustCompile(`^G-[A-Z0-9]{4,}$`)

// checkAnalyticsConfig: every Production build (the Production and QA
// Targets) needs a Google Analytics 4 Measurement ID (ADR-0008), or the
// pageviews Raptive reviews would silently go unrecorded. Preview builds
// carry no analytics, so they are not checked.
func checkAnalyticsConfig(s *site) error {
	if s.cfg.Mode == Preview {
		return nil
	}
	switch id := s.cfg.GAMeasurementID; {
	case id == "":
		s.fail("", ruleAnalyticsConfig, "ga_measurement_id is blank; create a GA4 web data stream for the site (analytics.google.com > Admin > Data streams) and copy its Measurement ID (G-...)")
	case !gaMeasurementID.MatchString(id):
		s.fail("", ruleAnalyticsConfig, "ga_measurement_id = %q is not a GA4 Measurement ID (G- followed by letters and digits)", id)
	}
	return nil
}

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

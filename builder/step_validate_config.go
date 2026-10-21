package builder

func init() {
	registerStep(phaseValidate, 5, "ads config", checkAdsConfig)
}

const ruleAdsConfig = "ads-config"

// checkAdsConfig: turning ads on needs the AdSense account, or every slot
// would render a unit with no publisher.
func checkAdsConfig(s *site) error {
	if s.cfg.AdsEnabled && s.cfg.AdSenseClient == "" {
		s.fail("", ruleAdsConfig, "ads_enabled is true but adsense_client is blank; set the AdSense publisher id (ca-pub-...)")
	}
	return nil
}

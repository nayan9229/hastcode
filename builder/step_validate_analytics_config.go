package builder

import "regexp"

func init() {
	registerStep(phaseValidate, 7, "analytics config", checkAnalyticsConfig)
}

const ruleAnalyticsConfig = "analytics-config"

// gaMeasurementID is a Google Analytics 4 web data stream id.
var gaMeasurementID = regexp.MustCompile(`^G-[A-Z0-9]{4,}$`)

// checkAnalyticsConfig: every Production build (the Production and QA
// Targets) needs a Google Analytics 4 Measurement ID (ADR-0008), or the
// pageviews Raptive reviews would silently go unrecorded, and a site_url
// with a host, the only host that reports. Preview builds carry no
// analytics (analyticsOn), so they are not checked.
func checkAnalyticsConfig(s *site) error {
	if !s.analyticsOn() {
		return nil
	}
	switch id := s.cfg.GAMeasurementID; {
	case id == "":
		s.fail("", ruleAnalyticsConfig, "ga_measurement_id is blank; create a GA4 web data stream for the site (analytics.google.com > Admin > Data streams) and copy its Measurement ID (G-...)")
	case !gaMeasurementID.MatchString(id):
		s.fail("", ruleAnalyticsConfig, "ga_measurement_id = %q is not a GA4 Measurement ID (G- followed by letters and digits)", id)
	}
	if s.siteHost() == "" {
		s.fail("", ruleAnalyticsConfig, "site_url %q has no host; analytics reports only on site_url's host", s.cfg.SiteURL)
	}
	return nil
}

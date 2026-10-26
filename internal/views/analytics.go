package views

import (
	"encoding/json"
	"net/url"
	"strings"
)

// consentDeniedRegions are where analytics and ad storage default to
// denied until the AdSense Privacy & messaging consent message grants them
// (ADR-0008): the EEA (EU, Iceland, Liechtenstein, Norway), the UK and
// Switzerland, as ISO 3166-1 alpha-2 codes.
var consentDeniedRegions = []string{
	"AT", "BE", "BG", "HR", "CY", "CZ", "DK", "EE", "FI", "FR", "DE", "GR", "HU", "IE",
	"IT", "LV", "LT", "LU", "MT", "NL", "PL", "PT", "RO", "SK", "SI", "ES", "SE",
	"IS", "LI", "NO", "GB", "CH",
}

// gtagURL is the Google tag library for one Measurement ID.
func gtagURL(id string) string {
	return "https://www.googletagmanager.com/gtag/js?id=" + url.QueryEscape(id)
}

// analyticsSnippet is the inline code before gtag.js: Google's dataLayer
// queue, then, on the site's own host only, regional consent defaults
// (denied in consentDeniedRegions, granted elsewhere; wait_for_update gives
// the consent message time to restore a stored choice) and the property's
// config. On any other host (the QA Target) it sets Google's documented
// opt-out flag instead, so gtag.js loads but never reports.
func analyticsSnippet(a Analytics) string {
	id := jsString(a.MeasurementID)
	regions, _ := json.Marshal(consentDeniedRegions)
	var b strings.Builder
	b.WriteString(`window.dataLayer=window.dataLayer||[];function gtag(){dataLayer.push(arguments);}`)
	b.WriteString(`if(location.hostname!==` + jsString(a.Host) + `){window["ga-disable-"+` + id + `]=true;}else{`)
	b.WriteString(`gtag("consent","default",{ad_storage:"denied",ad_user_data:"denied",ad_personalization:"denied",analytics_storage:"denied",wait_for_update:500,region:` + string(regions) + `});`)
	b.WriteString(`gtag("consent","default",{ad_storage:"granted",ad_user_data:"granted",ad_personalization:"granted",analytics_storage:"granted"});`)
	b.WriteString(`gtag("js",new Date());gtag("config",` + id + `);}`)
	return b.String()
}

// jsString is s as a JavaScript string literal, safe inside <script>
// (encoding/json escapes <, > and &).
func jsString(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

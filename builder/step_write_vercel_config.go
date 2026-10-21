package builder

import "hastcode.com/site/internal/vercel"

func init() {
	registerStep(phaseWrite, 100, "write config.json", writeVercelConfig)
}

// writeVercelConfig writes the Build Output API v3 config.json. Routes,
// overrides and headers (ticket 10) are added to vercel.Config.
func writeVercelConfig(s *site) error {
	return s.out.WriteConfig(vercel.Config{Version: 3})
}

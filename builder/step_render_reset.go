package builder

import "hastcode.com/site/internal/cloudflare"

func init() {
	// First in the render phase: every later step writes into a clean tree.
	registerStep(phaseRender, 0, "reset output", resetOutput)
}

// resetOutput clears static/, config.json and wrangler.json (_headers and
// _redirects live in static/).
func resetOutput(s *site) error {
	if err := s.out.Reset(); err != nil {
		return err
	}
	return cloudflare.RemoveConfig(s.out.Dir)
}

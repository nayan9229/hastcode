package builder

func init() {
	// First in the render phase: every later step writes into a clean tree.
	registerStep(phaseRender, 0, "reset output", resetOutput)
}

func resetOutput(s *site) error { return s.out.Reset() }

package builder

func init() {
	registerStep(phaseValidate, 40, "warn on explainer section length", warnSectionLength)
}

const ruleSectionLength = "section-length"

// maxSectionWords is the soft cap on an H2 section's prose (code
// excluded). Answer engines extract 120–180-word sections best.
const maxSectionWords = 250

// warnSectionLength warns (never fails) on any H2 section over
// maxSectionWords words.
func warnSectionLength(s *site) error {
	for _, e := range s.explainers {
		_, sections := e.Body.Sections()
		for _, sec := range sections {
			if n := e.Body.Words(sec.Blocks); n > maxSectionWords {
				s.warn(e.File, ruleSectionLength, "## %s has %d words, over the %d-word guide; split it", sec.Heading, n, maxSectionWords)
			}
		}
	}
	return nil
}

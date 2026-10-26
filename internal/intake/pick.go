package intake

import (
	"context"
	"errors"
)

// PickOptions override what a pick copies from the Lead's Source.
type PickOptions struct {
	// Kind replaces the Lead's suggested kind when set (one of Kinds).
	Kind string
	// Vertical replaces the Lead's Vertical when set (a site.yaml slug).
	Vertical string
}

// Pick turns the Lead id into a draft skeleton under ContentDir and marks
// it Picked; it returns the skeleton's path. `make lead ID=…` runs it.
//
// Not built yet: ticket 09 (.scratch/hastcode-content-intake/issues/09-*),
// after the per-kind templates (ticket 08).
func (in *Intake) Pick(ctx context.Context, id string, opts PickOptions) (string, error) {
	return "", errors.New("intake: picking a Lead is not built yet (ticket 09)")
}

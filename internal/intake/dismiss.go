package intake

import (
	"context"
	"errors"
)

// Dismiss hides the Lead id for good: its decision becomes Dismissed and
// later queues leave it out. `make lead-dismiss ID=…` runs it.
//
// Not built yet: ticket 05 (.scratch/hastcode-content-intake/issues/05-*).
func (in *Intake) Dismiss(ctx context.Context, id string) error {
	return errors.New("intake: dismissing a Lead is not built yet (ticket 05)")
}

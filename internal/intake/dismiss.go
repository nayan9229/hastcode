package intake

import (
	"context"
	"fmt"
	"strings"
)

// Dismiss hides the Lead id for good: its decision becomes Dismissed in the
// seen-state, so every later queue leaves it out and no run lists it
// again. Only an open Lead can be dismissed. `make lead-dismiss ID=…` runs
// it. It changes the state only; the next run rewrites the queue.
func (in *Intake) Dismiss(ctx context.Context, id string) error {
	id = strings.ToLower(strings.TrimSpace(id))
	if err := in.checkStateDir(); err != nil {
		return err
	}
	st, err := loadState(in.StateDir)
	if err != nil {
		return err
	}
	lead, ok := st.Leads[id]
	if !ok {
		return fmt.Errorf("intake: no Lead %q in %s (IDs are in the queue, %s)", id, in.StateDir, QueueFile)
	}
	if lead.Decision != Open {
		return fmt.Errorf("intake: Lead %s (%q) is already %s; only an open Lead can be dismissed", id, lead.Title, lead.Decision)
	}
	lead.Decision = Dismissed
	st.Leads[id] = lead
	if err := st.save(in.StateDir); err != nil {
		return err
	}
	in.logf("dismissed %s: %s (the next run's queue leaves it out)", id, lead.Title)
	return nil
}

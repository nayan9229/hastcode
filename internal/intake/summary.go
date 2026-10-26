package intake

import (
	"fmt"
	"io"
)

// WriteSummary prints the run summary the command shows after a run: the
// new and open counts, where the queue is, and every Source skipped.
func (r *Report) WriteSummary(w io.Writer) {
	fmt.Fprintf(w, "%s, %d open: %s\n", plural(len(r.New), "new Lead"), r.Open, r.QueueFile)
	for _, s := range r.Skipped {
		fmt.Fprintf(w, "skipped %s: %s\n", s.Source, s.Reason)
	}
	r.writeSourceHealth(w) // failed and silent Sources (source_health.go)
}

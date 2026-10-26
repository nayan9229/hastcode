package intake

import (
	"fmt"
	"io"
	"strings"
)

// WriteSummary prints the run summary the command shows after a run: the
// new and open counts, where the queue is, the new Leads per Vertical and
// per kind, how many Leads expired, then the failed and silent Sources.
func (r *Report) WriteSummary(w io.Writer) {
	fmt.Fprintf(w, "%s, %d open: %s\n", plural(len(r.New), "new Lead"), r.Open, r.QueueFile)
	fmt.Fprintf(w, "new by Vertical: %s\n", formatCounts(r.NewByVertical))
	fmt.Fprintf(w, "new by kind: %s\n", formatCounts(r.NewByKind))
	if r.Expired > 0 {
		fmt.Fprintf(w, "expired %s left open %d days\n", plural(r.Expired, "Lead"), wholeDays(ExpireAfter))
	}
	r.writeSourceHealth(w) // failed and silent Sources (source_health.go)
}

// countNew counts the new Leads per Vertical (display order) and per kind
// (Kinds order), zeros included.
func (in *Intake) countNew(leads []Lead) (byVertical, byKind []Count) {
	for _, v := range in.Verticals {
		c := Count{Name: v.Name}
		for _, l := range leads {
			if l.Vertical == v.Slug {
				c.New++
			}
		}
		byVertical = append(byVertical, c)
	}
	for _, k := range Kinds {
		c := Count{Name: k}
		for _, l := range leads {
			if l.Kind == k {
				c.New++
			}
		}
		byKind = append(byKind, c)
	}
	return byVertical, byKind
}

// formatCounts is "Distributed Systems 2, Databases & Storage 0".
func formatCounts(counts []Count) string {
	parts := make([]string, len(counts))
	for i, c := range counts {
		parts[i] = fmt.Sprintf("%s %d", c.Name, c.New)
	}
	return strings.Join(parts, ", ")
}

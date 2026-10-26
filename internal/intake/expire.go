package intake

import "time"

// ExpireAfter is how long a Lead stays open undecided: a run whose clock is
// this long or more after the Lead was first seen marks it Expired, so the
// queue never becomes a backlog.
const ExpireAfter = 7 * 24 * time.Hour

// expire marks every open Lead first seen ExpireAfter or more before now as
// Expired and returns how many it marked.
func (st *state) expire(now time.Time) int {
	n := 0
	for id, l := range st.Leads {
		if l.Decision == Open && now.Sub(l.FirstSeen) >= ExpireAfter {
			l.Decision = Expired
			st.Leads[id] = l
			n++
		}
	}
	return n
}

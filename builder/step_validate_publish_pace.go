package builder

import (
	"sort"
	"strings"
	"time"
)

func init() {
	registerStep(phaseValidate, 46, "check publish pace", checkPublishPace)
}

// rulePublishPace: more than maxPublishedPerDay non-draft Pieces carry the
// same `published` day (ADR-0010). Google's scaled-content policy applies
// from the first page, so the site grows at a steady pace, never in
// bursts. An error in production, a warning in preview (a Preview may show
// a whole backlog of drafts).
const rulePublishPace = "publish-pace"

// maxPublishedPerDay is the daily Publish cap, Explainers and Briefs
// together.
const maxPublishedPerDay = 2

// paceStart is the first publication day the cap applies to: the 32
// Explainers published on 28 and 29 October 2026, before ADR-0010, stay.
var paceStart = time.Date(2026, 10, 30, 0, 0, 0, 0, time.UTC)

// checkPublishPace counts non-draft Pieces per `published` day, in
// published's own offset, and reports each day over the cap with its files.
func checkPublishPace(s *site) error {
	report := s.fail
	if s.cfg.Mode == Preview {
		report = s.warn
	}
	byDay := map[time.Time][]string{}
	add := func(draft bool, published time.Time, file string) {
		if draft || published.IsZero() {
			return
		}
		if day := calendarDayUTC(published); !day.Before(paceStart) {
			byDay[day] = append(byDay[day], file)
		}
	}
	for _, e := range s.explainers {
		add(e.Draft, e.Published, e.File)
	}
	for _, b := range s.briefs {
		add(b.Draft, b.Published, b.File)
	}
	days := make([]time.Time, 0, len(byDay))
	for d := range byDay {
		days = append(days, d)
	}
	sort.Slice(days, func(i, j int) bool { return days[i].Before(days[j]) })
	for _, d := range days {
		files := byDay[d]
		if len(files) <= maxPublishedPerDay {
			continue
		}
		sort.Strings(files)
		report("", rulePublishPace, "%d Pieces are published on %s, over the cap of %d a day (ADR-0010): keep draft: true on all but %d and give the rest later published days: %s",
			len(files), d.Format(time.DateOnly), maxPublishedPerDay, maxPublishedPerDay, strings.Join(files, ", "))
	}
	return nil
}

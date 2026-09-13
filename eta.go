package readable

import (
	"time"
)

// ETA estimates the time left for a task that has done of total units after
// elapsed, assuming a constant rate, as "~3m left". The estimate keeps the
// largest unit only and is at least one second. It returns "done" when
// done >= total > 0 and "" when there is no rate yet (done, total or elapsed
// not positive).
//
//	ETA(25, 100, time.Minute)   // "~3m left"
//	ETA(100, 100, time.Minute)  // "done"
//	ETA(0, 100, time.Minute)    // ""
func ETA(done, total int64, elapsed time.Duration) string {
	switch {
	case done <= 0 || total <= 0 || elapsed <= 0:
		return ""
	case done >= total:
		return "done"
	}
	left := float64(elapsed) / float64(done) * float64(total-done)
	d := time.Duration(maxDurationFloat)
	if left < maxDurationFloat {
		d = max(time.Duration(left), time.Second)
	}
	return "~" + DurationWithOptions(d, DurationOptions{Units: 1}) + " left"
}

// maxDurationFloat is the largest float64 below math.MaxInt64 nanoseconds.
const maxDurationFloat = float64(1<<63 - 1024)

// Calendar describes t for schedules and feeds, relative to the calendar day
// of now in t's location:
//
//	same day          "today 14:30"
//	previous day      "yesterday 14:30"
//	next day          "tomorrow 14:30"
//	within 6 days     "Monday 14:30"
//	same year         "Sep 2, 14:30"
//	otherwise         "Sep 2, 2025, 14:30"
func Calendar(now, t time.Time) string {
	now = now.In(t.Location())
	clock := t.Format("15:04")
	switch diff := dateCivilDay(t) - dateCivilDay(now); {
	case diff == 0:
		return "today " + clock
	case diff == -1:
		return "yesterday " + clock
	case diff == 1:
		return "tomorrow " + clock
	case diff > -7 && diff < 7:
		return t.Weekday().String() + " " + clock
	case t.Year() == now.Year():
		return t.Format("Jan 2, 15:04")
	}
	return t.Format("Jan 2, 2006, 15:04")
}

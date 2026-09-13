package readable

import "time"

// Date describes the calendar day of t relative to the current time.
// It calls time.Now; prefer DateFrom for deterministic output.
func Date(t time.Time) string {
	return DateFrom(time.Now(), t)
}

// DateFrom describes the calendar day of t relative to now. Calendar days are
// compared in t's location: now is converted to t.Location() first. Any
// year, including years before 1970 and after 9999, is supported.
//
//	same day          "Today"
//	previous day      "Yesterday"
//	next day          "Tomorrow"
//	same year         "Sep 13"
//	otherwise         "Sep 13, 2025"
func DateFrom(now, t time.Time) string {
	now = now.In(t.Location())
	switch dateCivilDay(t) - dateCivilDay(now) {
	case 0:
		return "Today"
	case -1:
		return "Yesterday"
	case 1:
		return "Tomorrow"
	}
	if t.Year() == now.Year() {
		return t.Format("Jan 2")
	}
	return t.Format("Jan 2, 2006")
}

// Time describes t relative to the current time followed by its absolute
// date and time. It calls time.Now; prefer TimeFrom for deterministic output.
func Time(t time.Time) string {
	return TimeFrom(time.Now(), t)
}

// TimeFrom returns RelativeTimeFrom(now, t), a middle dot and t's date and
// 24-hour time in t's location. The year is included only when it differs
// from now's year in that location.
//
//	"2 hours ago · Sep 13, 07:20"
//	"1 year ago · Sep 13, 2025, 07:20"
func TimeFrom(now, t time.Time) string {
	layout := "Jan 2, 15:04"
	if t.Year() != now.In(t.Location()).Year() {
		layout = "Jan 2, 2006, 15:04"
	}
	var arr [64]byte
	buf := append(arr[:0], RelativeTimeFrom(now, t)...)
	buf = append(buf, " · "...)
	buf = t.AppendFormat(buf, layout)
	return string(buf)
}

// TimeRange formats the interval between start and end relative to the
// current time. It calls time.Now; prefer TimeRangeFrom for deterministic
// output.
func TimeRange(start, end time.Time) string {
	return TimeRangeFrom(time.Now(), start, end)
}

// TimeRangeFrom formats the interval between start and end using a 24-hour
// clock. If end is before start the two are swapped. Both ends are shown in
// start's location, and "today" and "this year" are judged against now
// converted to that location.
//
//	same day, today        "09:30–11:45"
//	same day               "Sep 20, 09:30–11:45"
//	different days         "Sep 13, 23:30 → Sep 14, 02:15"
//	different years        "Dec 31, 2026, 23:30 → Jan 1, 2027, 02:15"
//
// The year is shown on both dates when start and end fall in different years
// or when either differs from now's year. The en dash is U+2013.
func TimeRangeFrom(now, start, end time.Time) string {
	if end.Before(start) {
		start, end = end, start
	}
	loc := start.Location()
	end = end.In(loc)
	now = now.In(loc)

	var arr [64]byte
	buf := arr[:0]
	dayLayout := "Jan 2, "
	if start.Year() != now.Year() || end.Year() != now.Year() {
		dayLayout = "Jan 2, 2006, "
	}
	if sd := dateCivilDay(start); sd == dateCivilDay(end) {
		if sd != dateCivilDay(now) {
			buf = start.AppendFormat(buf, dayLayout)
		}
		buf = start.AppendFormat(buf, "15:04")
		buf = append(buf, "–"...)
		buf = end.AppendFormat(buf, "15:04")
		return string(buf)
	}
	buf = start.AppendFormat(buf, dayLayout+"15:04")
	buf = append(buf, " → "...)
	buf = end.AppendFormat(buf, dayLayout+"15:04")
	return string(buf)
}

// dateCivilDay returns the day number of t's calendar date in t's location.
func dateCivilDay(t time.Time) int64 {
	y, m, d := t.Date()
	// Midnight UTC is an exact multiple of secondsPerDay, so the truncating
	// division is exact for dates before 1970 too.
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix() / secondsPerDay
}

package readable

import (
	"strconv"
	"time"
)

// RelativeTime describes t relative to the current time, e.g. "2 hours ago".
// It calls time.Now; prefer RelativeTimeFrom for deterministic output.
func RelativeTime(t time.Time) string {
	return RelativeTimeFrom(time.Now(), t)
}

// relativeJustNowSeconds is the magnitude below which RelativeTimeFrom
// reports "just now".
const relativeJustNowSeconds = 5

// RelativeTimeFrom describes t relative to now: "in N units" when t is after
// now and "N units ago" otherwise. Counts are truncated towards zero, the
// sub-second part of the difference is ignored and units are pluralised
// unless the count is 1.
//
//	|diff| < 5s     "just now"
//	< 60s           "N seconds ago"   / "in N seconds"
//	< 60m           "N minutes ago"
//	< 24h           "N hours ago"
//	< 7d            "N days ago"
//	< 30d           "N weeks ago"
//	< 365d          "N months ago"    (30-day months)
//	otherwise       "N years ago"     (365-day years)
//
// The difference is computed with time.Time.Sub, which saturates at
// math.MaxInt64 and math.MinInt64 nanoseconds, so any gap of more than about
// 292 years reports "292 years ago" or "in 292 years".
// Time zones do not affect the result; only the instants are compared.
func RelativeTimeFrom(now, t time.Time) string {
	n, unit, future, justNow := relativeParts(now, t)
	if justNow {
		return "just now"
	}
	name := relativeUnitNames[unit]
	var arr [32]byte
	buf := arr[:0]
	if future {
		buf = append(buf, "in "...)
	}
	buf = strconv.AppendUint(buf, n, 10)
	buf = append(buf, ' ')
	buf = append(buf, name...)
	if n != 1 {
		buf = append(buf, 's')
	}
	if !future {
		buf = append(buf, " ago"...)
	}
	return string(buf)
}

// Unit indexes returned by relativeParts.
const (
	relSecond = iota
	relMinute
	relHour
	relDay
	relWeek
	relMonth
	relYear
)

var relativeUnitNames = [...]string{"second", "minute", "hour", "day", "week", "month", "year"}

// relativeParts splits t - now into a truncated count, a unit index
// (relSecond...relYear) and a direction, as documented on RelativeTimeFrom.
func relativeParts(now, t time.Time) (n uint64, unit int, future, justNow bool) {
	diff := t.Sub(now)
	future = diff > 0
	_, abs := absInt64(int64(diff))
	secs := abs / uint64(time.Second)
	days := secs / secondsPerDay
	switch {
	case secs < relativeJustNowSeconds:
		return 0, relSecond, future, true
	case secs < secondsPerMinute:
		return secs, relSecond, future, false
	case secs < secondsPerHour:
		return secs / secondsPerMinute, relMinute, future, false
	case secs < secondsPerDay:
		return secs / secondsPerHour, relHour, future, false
	case days < daysPerWeek:
		return days, relDay, future, false
	case days < daysPerMonth:
		return days / daysPerWeek, relWeek, future, false
	case days < daysPerYear:
		return days / daysPerMonth, relMonth, future, false
	}
	return days / daysPerYear, relYear, future, false
}

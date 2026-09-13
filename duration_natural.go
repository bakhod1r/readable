package readable

import (
	"strconv"
	"time"
)

// DurationNatural formats d as an English phrase with full unit names,
// joining the last two components with "and". Days, hours, minutes and
// seconds are used, zero components are omitted and the sub-second remainder
// is truncated. Durations below one second use a single millisecond,
// microsecond or nanosecond component, like DurationLong. Negative durations
// are prefixed with "-"; math.MinInt64 is handled without overflow.
//
//	DurationNatural(49*time.Hour + 32*time.Minute) // "2 days, 1 hour and 32 minutes"
//	DurationNatural(time.Hour + 5*time.Second)     // "1 hour and 5 seconds"
//	DurationNatural(125 * time.Millisecond)        // "125 milliseconds"
//	DurationNatural(0)                             // "0 seconds"
func DurationNatural(d time.Duration) string {
	neg, abs := absInt64(int64(d))
	if abs < uint64(time.Second) {
		return formatDuration(d, 0, true)
	}
	var vals [4]uint64
	var idx [4]int
	n := 0
	for i, p := range durationParts {
		if v := abs / p.size; v != 0 {
			vals[n], idx[n] = v, i
			n++
		}
		abs %= p.size
	}
	var arr [96]byte
	buf := arr[:0]
	if neg {
		buf = append(buf, '-')
	}
	for i := range n {
		switch i {
		case 0:
		case n - 1:
			buf = append(buf, " and "...)
		default:
			buf = append(buf, ", "...)
		}
		buf = appendPart(buf, vals[i], durationParts[idx[i]], true)
	}
	return string(buf)
}

// DurationApprox describes the magnitude of d as a rough English phrase.
// The sign of d is ignored. The unit is chosen from the exact magnitude and
// the count is rounded to the nearest whole unit, halves rounding up; a count
// that rounds up to a full next unit is reported in that unit.
//
//	< 1 minute   "less than a minute"
//	< 1 hour     "about N minutes"  (59m30s is "about 1 hour")
//	< 1 day      "about N hours"
//	< 30 days    "about N days"
//	< 365 days   "about N months"   (30-day months)
//	otherwise    "about N years"    (365-day years)
//
// Examples:
//
//	DurationApprox(100 * time.Minute) // "about 2 hours"
//	DurationApprox(49 * time.Hour)    // "about 2 days"
//	DurationApprox(30 * time.Second)  // "less than a minute"
func DurationApprox(d time.Duration) string {
	_, abs := absInt64(int64(d))
	const (
		minute = uint64(secondsPerMinute * time.Second)
		hour   = uint64(secondsPerHour * time.Second)
		day    = uint64(secondsPerDay * time.Second)
		month  = daysPerMonth * day
		year   = daysPerYear * day
	)
	var n uint64
	var name string
	switch {
	case abs < minute:
		return "less than a minute"
	case abs < hour:
		n, name = durationRoundHalfUp(abs, minute), "minute"
		if n == 60 {
			n, name = 1, "hour"
		}
	case abs < day:
		n, name = durationRoundHalfUp(abs, hour), "hour"
		if n == 24 {
			n, name = 1, "day"
		}
	case abs < month:
		n, name = durationRoundHalfUp(abs, day), "day"
		if n == 30 {
			n, name = 1, "month"
		}
	case abs < year:
		n, name = durationRoundHalfUp(abs, month), "month"
	default:
		n, name = durationRoundHalfUp(abs, year), "year"
	}
	var arr [32]byte
	buf := append(arr[:0], "about "...)
	buf = strconv.AppendUint(buf, n, 10)
	buf = append(buf, ' ')
	buf = append(buf, name...)
	if n != 1 {
		buf = append(buf, 's')
	}
	return string(buf)
}

// durationRoundHalfUp returns abs/size rounded half up.
func durationRoundHalfUp(abs, size uint64) uint64 {
	q, _ := divRound(abs, size, 0)
	return q
}

// DurationShort is equivalent to Duration, provided for symmetry with
// DurationNatural and DurationApprox.
//
//	DurationShort(90 * time.Second) // "1m 30s"
func DurationShort(d time.Duration) string {
	return Duration(d)
}

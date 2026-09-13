package readable

import (
	"strconv"
	"time"
)

// Calendar units in seconds and days. DurationApprox and RelativeTimeFrom
// share them: a month is 30 days and a year is 365 days.
const (
	secondsPerMinute = 60
	secondsPerHour   = 60 * secondsPerMinute
	secondsPerDay    = 24 * secondsPerHour
	daysPerWeek      = 7
	daysPerMonth     = 30
	daysPerYear      = 365
)

type durationPart struct {
	size         uint64
	short        string
	long, plural string
}

var durationParts = []durationPart{
	{uint64(secondsPerDay * time.Second), "d", "day", "days"},
	{uint64(time.Hour), "h", "hour", "hours"},
	{uint64(time.Minute), "m", "minute", "minutes"},
	{uint64(time.Second), "s", "second", "seconds"},
}

var subSecondParts = []durationPart{
	{uint64(time.Millisecond), "ms", "millisecond", "milliseconds"},
	{uint64(time.Microsecond), "µs", "microsecond", "microseconds"},
	{uint64(time.Nanosecond), "ns", "nanosecond", "nanoseconds"},
}

// Duration formats d as space-separated days, hours, minutes and seconds,
// omitting zero components. Durations of one second or more truncate the
// sub-second remainder; shorter non-zero ones use a single ms, µs or ns
// component. Negative durations are prefixed with "-"; math.MinInt64 is
// handled without overflow.
//
//	Duration(3*time.Hour + 25*time.Minute) // "3h 25m"
//	Duration(90 * time.Second)             // "1m 30s"
//	Duration(125 * time.Millisecond)       // "125ms"
//	Duration(0)                            // "0s"
func Duration(d time.Duration) string {
	return formatDuration(d, 0, false)
}

// DurationWithPrecision is like Duration but shows at most units components,
// counted from the largest non-zero one; smaller components are truncated,
// never rounded. Zero components inside that window still count towards
// units. units <= 0 shows all components. It has no effect below one second.
//
//	DurationWithPrecision(3*24*time.Hour+12*time.Hour+32*time.Minute, 2) // "3d 12h"
//	DurationWithPrecision(time.Hour+5*time.Second, 2)                    // "1h"
func DurationWithPrecision(d time.Duration, units int) string {
	return formatDuration(d, units, false)
}

// DurationLong is like Duration but uses full unit names, pluralised unless
// the count is 1, separated by ", ". Zero is "0 seconds".
//
//	DurationLong(3*time.Hour + 25*time.Minute + 12*time.Second)
//	// "3 hours, 25 minutes, 12 seconds"
func DurationLong(d time.Duration) string {
	return formatDuration(d, 0, true)
}

// DurationLongWithPrecision is DurationLong limited to units components, as
// in DurationWithPrecision.
//
//	DurationLongWithPrecision(3*time.Hour+25*time.Minute+12*time.Second, 2) // "3 hours, 25 minutes"
func DurationLongWithPrecision(d time.Duration, units int) string {
	return formatDuration(d, units, true)
}

func formatDuration(d time.Duration, units int, long bool) string {
	neg, abs := absInt64(int64(d))
	var arr [96]byte
	buf := arr[:0]
	if neg {
		buf = append(buf, '-')
	}
	if abs == 0 {
		if long {
			return "0 seconds"
		}
		return "0s"
	}
	if abs < uint64(time.Second) {
		for _, p := range subSecondParts {
			if abs >= p.size {
				buf = appendPart(buf, abs/p.size, p, long)
				break
			}
		}
		return string(buf)
	}

	start := 0
	for abs < durationParts[start].size {
		start++
	}
	end := len(durationParts)
	if units > 0 && start+units < end {
		end = start + units
	}
	first := true
	for _, p := range durationParts[start:end] {
		v := abs / p.size
		abs %= p.size
		if v == 0 {
			continue
		}
		if !first {
			if long {
				buf = append(buf, ", "...)
			} else {
				buf = append(buf, ' ')
			}
		}
		first = false
		buf = appendPart(buf, v, p, long)
	}
	return string(buf)
}

func appendPart(buf []byte, v uint64, p durationPart, long bool) []byte {
	buf = strconv.AppendUint(buf, v, 10)
	if !long {
		return append(buf, p.short...)
	}
	buf = append(buf, ' ')
	if v == 1 {
		return append(buf, p.long...)
	}
	return append(buf, p.plural...)
}

var latencyUnits = []unit{
	{uint64(time.Nanosecond), "ns"},
	{uint64(time.Microsecond), "µs"},
	{uint64(time.Millisecond), "ms"},
	{uint64(time.Second), "s"},
	{uint64(time.Minute), "m"},
	{uint64(time.Hour), "h"},
}

// Latency formats d using the single most appropriate unit (ns, µs, ms, s, m
// or h) with at most one fraction digit. It suits logs and metrics.
//
//	Latency(350 * time.Microsecond)  // "350µs"
//	Latency(42 * time.Millisecond)   // "42ms"
//	Latency(1400 * time.Millisecond) // "1.4s"
//	Latency(120 * time.Second)       // "2m"
func Latency(d time.Duration) string {
	neg, abs := absInt64(int64(d))
	return formatUnits(neg, abs, latencyUnits, 1, false, ".", "")
}

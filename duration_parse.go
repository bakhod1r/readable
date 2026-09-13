package readable

import (
	"math/big"
	"strings"
	"time"
)

var durationUnitsByName = map[string]time.Duration{
	"ns": time.Nanosecond, "nanosecond": time.Nanosecond, "nanoseconds": time.Nanosecond,
	"µs": time.Microsecond, "μs": time.Microsecond, "us": time.Microsecond,
	"microsecond": time.Microsecond, "microseconds": time.Microsecond,
	"ms": time.Millisecond, "millisecond": time.Millisecond, "milliseconds": time.Millisecond,
	"s": time.Second, "sec": time.Second, "second": time.Second, "seconds": time.Second,
	"m": time.Minute, "min": time.Minute, "minute": time.Minute, "minutes": time.Minute,
	"h": time.Hour, "hr": time.Hour, "hour": time.Hour, "hours": time.Hour,
	"d": 24 * time.Hour, "day": 24 * time.Hour, "days": 24 * time.Hour,
}

// ParseDuration parses a duration as printed by Duration, DurationLong,
// DurationNatural or Latency: "3d 12h", "2h30m", "1.5ms", "350µs",
// "2 days, 1 hour and 32 minutes". Components are a number and a unit (ns,
// µs/us, ms, s, m, h, d or their English names, case-insensitive), separated
// by optional spaces, commas or "and". A day is 24 hours. A leading "-"
// negates the whole duration. Results are rounded to the nearest nanosecond,
// half away from zero.
//
//	ParseDuration("3h 25m") // 3h25m0s, nil
//	ParseDuration("1.5ms")  // 1.5ms, nil
//
// Errors wrap ErrSyntax, ErrUnit or ErrRange (outside time.Duration).
func ParseDuration(s string) (time.Duration, error) {
	d, err := parseDuration(s)
	if err != nil {
		return 0, parseError("ParseDuration", s, err)
	}
	return d, nil
}

func parseDuration(s string) (time.Duration, error) {
	t := strings.TrimSpace(s)
	neg := strings.HasPrefix(t, "-")
	if neg || strings.HasPrefix(t, "+") {
		t = t[1:]
	}
	total := new(big.Rat)
	parts := 0
	for i := 0; ; parts++ {
		i = skipDurationSeparators(t, i)
		if i == len(t) {
			break
		}
		r, j, ok := scanNumber(t, i)
		if !ok {
			return 0, ErrSyntax
		}
		for j < len(t) && t[j] == ' ' {
			j++
		}
		k := j
		for k < len(t) && !strings.ContainsRune("0123456789., ", rune(t[k])) {
			k++
		}
		unit, ok := durationUnitsByName[strings.ToLower(t[j:k])]
		if !ok {
			return 0, ErrUnit
		}
		total.Add(total, r.Mul(r, new(big.Rat).SetInt64(int64(unit))))
		i = k
	}
	if parts == 0 {
		return 0, ErrSyntax
	}
	if neg {
		total.Neg(total)
	}
	n := roundRat(total)
	if !fitsInt64(n) {
		return 0, ErrRange
	}
	return time.Duration(n.Int64()), nil
}

// skipDurationSeparators skips spaces, commas and the word "and".
func skipDurationSeparators(t string, i int) int {
	for i < len(t) {
		switch {
		case t[i] == ' ' || t[i] == ',':
			i++
		case strings.HasPrefix(t[i:], "and "):
			i += len("and ")
		default:
			return i
		}
	}
	return i
}

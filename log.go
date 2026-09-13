package readable

import (
	"strconv"
	"time"
)

// LogDuration formats d for machine-readable logs. It is exactly
// time.Duration.String: stable, no spaces, full precision.
//
//	LogDuration(1532 * time.Millisecond) // "1.532s"
//	LogDuration(time.Hour)               // "1h0m0s"
func LogDuration(d time.Duration) string {
	return d.String()
}

// LogBytes formats n bytes for logs with IEC units, up to DefaultPrecision
// trimmed fraction digits and no space.
//
//	LogBytes(12345678) // "11.77MiB"
func LogBytes(n uint64) string {
	return formatUnits(false, n, iecUnits, DefaultPrecision, false, ".", "")
}

// LogNumber formats n as a plain base-10 integer with no grouping.
//
//	LogNumber(1234567) // "1234567"
func LogNumber(n int64) string {
	return strconv.FormatInt(n, 10)
}

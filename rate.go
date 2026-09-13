package readable

import (
	"math"
	"math/bits"
	"time"
)

// Rate formats count events observed over period as a per-second rate like
// NumberFloat followed by "/s".
//
//	Rate(125000, time.Second) // "125K/s"
//	Rate(300, time.Minute)    // "5/s"
//
// A non-positive period yields "NaN/s".
func Rate(count float64, period time.Duration) string {
	return RateWithLabel(count, period, "")
}

// RateWithLabel is like Rate with label and a space between the value and
// "/s". An empty label is omitted.
//
//	RateWithLabel(1200, time.Second, "req") // "1.2K req/s"
func RateWithLabel(count float64, period time.Duration, label string) string {
	r := math.NaN()
	if period > 0 {
		r = count / period.Seconds()
	}
	return perUnit(r, label, "/s")
}

// BytesRate formats n bytes transferred over period as binary units per
// second like Bytes followed by "/s". The per-second value is computed
// exactly in integer arithmetic, truncated, and saturates at MaxUint64 bytes.
//
//	BytesRate(5662310, time.Second) // "5.4 MB/s"
//
// A non-positive period yields "NaN/s".
func BytesRate(n uint64, period time.Duration) string {
	if period <= 0 {
		return "NaN/s"
	}
	var arr [64]byte
	buf := appendUnits(arr[:0], false, bytesPerSecond(n, period), binaryUnits, DefaultPrecision, false, ".", " ")
	return string(append(buf, "/s"...))
}

// bytesPerSecond returns floor(n*1e9/period) for a positive period,
// saturating at math.MaxUint64.
func bytesPerSecond(n uint64, period time.Duration) uint64 {
	hi, lo := bits.Mul64(n, uint64(time.Second))
	d := uint64(period)
	if hi >= d { // quotient needs more than 64 bits
		return math.MaxUint64
	}
	q, _ := bits.Div64(hi, lo, d)
	return q
}

// perUnit renders v like NumberFloat, then " "+label when label is not
// empty, then suffix.
func perUnit(v float64, label, suffix string) string {
	var arr [64]byte
	buf := appendNumberFloat(arr[:0], v)
	if label != "" {
		buf = append(buf, ' ')
		buf = append(buf, label...)
	}
	return string(append(buf, suffix...))
}

package readable

import (
	"errors"
	"math"
)

// ErrUndefined is returned when a result has no meaningful value, such as a
// percentage change from zero. Test for it with errors.Is.
var ErrUndefined = errors.New("readable: undefined result")

// Percent formats a ratio as a percentage with up to DefaultPrecision fraction
// digits, rounded to nearest and trimmed of trailing zeros: 0.15 means 15%.
// Results that round to zero lose their sign.
//
//	Percent(0.1534) // "15.34%"
//	Percent(1)      // "100%"
//
// NaN and infinities format as "NaN%", "+Inf%" and "-Inf%".
func Percent(ratio float64) string {
	return PercentWithPrecision(ratio, DefaultPrecision)
}

// PercentWithPrecision is like Percent with at most precision fraction digits
// (clamped to [0, MaxPrecision]).
//
//	PercentWithPrecision(0.123456, 0) // "12%"
//	PercentWithPrecision(0.123456, 1) // "12.3%"
func PercentWithPrecision(ratio float64, precision int) string {
	x := ratio * 100
	switch {
	case math.IsNaN(x):
		return "NaN%"
	case math.IsInf(x, 1):
		return "+Inf%"
	case math.IsInf(x, -1):
		return "-Inf%"
	}
	var arr [64]byte
	buf := appendFloat(arr[:0], x, clampPrecision(precision))
	return string(append(buf, '%'))
}

// PercentChange formats the relative change from from to to, measured against
// the magnitude of from, like Percent.
//
//	PercentChange(100, 120) // "20%", nil
//	PercentChange(100, 80)  // "-20%", nil
//
// It returns ErrUndefined if from is zero, any input is NaN or infinite, or
// the change overflows.
func PercentChange(from, to float64) (string, error) {
	if from == 0 || math.IsNaN(from) || math.IsNaN(to) || math.IsInf(from, 0) || math.IsInf(to, 0) {
		return "", ErrUndefined
	}
	change := (to - from) / math.Abs(from)
	if math.IsInf(change, 0) {
		return "", ErrUndefined
	}
	return Percent(change), nil
}

package readable

import (
	"bytes"
	"math"
	"strconv"
)

var decimalUnits = []unit{
	{1, ""},
	{1e3, "K"},
	{1e6, "M"},
	{1e9, "B"},
	{1e12, "T"},
}

// NumberOptions configures NumberWithOptions.
type NumberOptions struct {
	// Precision is the maximum number of fraction digits (clamped to
	// [0, MaxPrecision]).
	Precision int
	// FixedPrecision keeps trailing zeros ("1.50M" instead of "1.5M").
	FixedPrecision bool
	// Decimal is the decimal separator. Empty means ".".
	Decimal string
}

// Number formats n using compact decimal units K, M, B and T with up to
// DefaultPrecision fraction digits, rounding half away from zero.
//
//	Number(999)     // "999"
//	Number(12500)   // "12.5K"
//	Number(1234567) // "1.23M"
//	Number(999999)  // "1M" (rounding promotes to the next unit)
//
// Values of a quadrillion or more stay in T ("9223372.04T" for MaxInt64).
func Number(n int64) string {
	return NumberWithPrecision(n, DefaultPrecision)
}

// NumberWithPrecision is like Number with at most precision fraction digits.
//
//	NumberWithPrecision(1234567, 1) // "1.2M"
//	NumberWithPrecision(1234567, 3) // "1.235M"
func NumberWithPrecision(n int64, precision int) string {
	neg, abs := absInt64(n)
	return formatUnits(neg, abs, decimalUnits, precision, false, ".", "")
}

// NumberWithOptions is like Number with explicit options.
func NumberWithOptions(n int64, opts NumberOptions) string {
	dec := opts.Decimal
	if dec == "" {
		dec = "."
	}
	neg, abs := absInt64(n)
	return formatUnits(neg, abs, decimalUnits, opts.Precision, opts.FixedPrecision, dec, "")
}

// NumberFloat formats f like Number with up to DefaultPrecision fraction
// digits, rounded to nearest by strconv and trimmed of trailing zeros. Values
// that round to 1000 of a unit move to the next one ("1K" for 999.995).
// Results that round to zero lose their sign ("0" for -0.001). Values of a
// quadrillion or more stay in T with all integer digits. NaN and infinities
// format as "NaN", "+Inf" and "-Inf".
//
//	NumberFloat(0.5)     // "0.5"
//	NumberFloat(1500.25) // "1.5K"
func NumberFloat(f float64) string {
	var arr [64]byte
	return string(appendNumberFloat(arr[:0], f))
}

// appendNumberFloat appends f formatted like NumberFloat.
func appendNumberFloat(buf []byte, f float64) []byte {
	switch {
	case math.IsNaN(f):
		return append(buf, "NaN"...)
	case math.IsInf(f, 1):
		return append(buf, "+Inf"...)
	case math.IsInf(f, -1):
		return append(buf, "-Inf"...)
	}
	abs := math.Abs(f)
	i := 0
	for i+1 < len(decimalUnits) && abs >= float64(decimalUnits[i+1].size) {
		i++
	}
	start := len(buf)
	if f < 0 {
		buf = append(buf, '-')
	}
	digits := len(buf)
	buf = strconv.AppendFloat(buf, abs/float64(decimalUnits[i].size), 'f', DefaultPrecision, 64)
	// Rounding may reach 1000 of the unit (999.995 -> "1000.00"). With
	// DefaultPrecision > 0 the output always has a decimal point.
	if i+1 < len(decimalUnits) && bytes.IndexByte(buf[digits:], '.') > len("999") {
		i++
		buf = strconv.AppendFloat(buf[:digits], abs/float64(decimalUnits[i].size), 'f', DefaultPrecision, 64)
	}
	buf = trimFloatBytes(buf, digits)
	if f < 0 && string(buf[digits:]) == "0" {
		buf = append(buf[:start], '0')
	}
	return append(buf, decimalUnits[i].suffix...)
}

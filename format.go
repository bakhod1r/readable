package readable

import (
	"bytes"
	"math/bits"
	"strconv"
)

// DefaultPrecision is the maximum number of fraction digits used by
// formatters that do not take an explicit precision.
const DefaultPrecision = 2

// MaxPrecision is the largest precision honoured; larger values are clamped.
const MaxPrecision = 9

var pow10 = [MaxPrecision + 1]uint64{1, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9}

type unit struct {
	size   uint64
	suffix string
}

func clampPrecision(p int) int {
	return min(max(p, 0), MaxPrecision)
}

// absInt64 splits n into sign and magnitude without overflowing on MinInt64.
func absInt64(n int64) (neg bool, abs uint64) {
	if n < 0 {
		return true, uint64(-(n + 1)) + 1
	}
	return false, uint64(n)
}

// divRound returns abs/size rounded half up to prec fraction digits as an
// integer part q and fraction digits frac (0 <= frac < 10^prec).
func divRound(abs, size uint64, prec int) (q, frac uint64) {
	q, r := abs/size, abs%size
	if prec == 0 {
		if r != 0 && r >= size-r {
			q++
		}
		return q, 0
	}
	pow := pow10[prec]
	// frac = floor((2*r*pow + size) / (2*size)), computed in 128 bits.
	hi, lo := bits.Mul64(r, 2*pow)
	var carry uint64
	lo, carry = bits.Add64(lo, size, 0)
	hi += carry
	frac, _ = bits.Div64(hi, lo, 2*size)
	if frac == pow {
		q++
		frac = 0
	}
	return q, frac
}

// formatUnits scales abs to the largest fitting unit and renders it.
// units must be ascending, each size an exact multiple of the previous one.
func formatUnits(neg bool, abs uint64, units []unit, prec int, fixed bool, dec, sep string) string {
	var arr [64]byte
	return string(appendUnits(arr[:0], neg, abs, units, prec, fixed, dec, sep))
}

// appendUnits is like formatUnits but appends to buf.
func appendUnits(buf []byte, neg bool, abs uint64, units []unit, prec int, fixed bool, dec, sep string) []byte {
	prec = clampPrecision(prec)
	i := 0
	for i+1 < len(units) && abs >= units[i+1].size {
		i++
	}
	p := prec
	if units[i].size == 1 {
		p = 0
	}
	q, frac := divRound(abs, units[i].size, p)
	// Rounding may reach the next unit (999.999K -> 1M).
	if i+1 < len(units) && q >= units[i+1].size/units[i].size {
		i++
		p = prec
		q, frac = divRound(abs, units[i].size, p)
	}

	if neg {
		buf = append(buf, '-')
	}
	buf = strconv.AppendUint(buf, q, 10)
	buf = appendFraction(buf, frac, p, fixed, dec)
	if units[i].suffix != "" {
		buf = append(buf, sep...)
		buf = append(buf, units[i].suffix...)
	}
	return buf
}

// appendFraction appends dec and frac zero-padded to p digits. Trailing zeros
// are trimmed unless fixed; nothing is appended if no digits remain.
func appendFraction(buf []byte, frac uint64, p int, fixed bool, dec string) []byte {
	if p == 0 {
		return buf
	}
	if !fixed {
		for p > 0 && frac%10 == 0 {
			frac /= 10
			p--
		}
		if p == 0 {
			return buf
		}
	}
	buf = append(buf, dec...)
	var digits [MaxPrecision]byte
	for j := p - 1; j >= 0; j-- {
		digits[j] = byte('0' + frac%10)
		frac /= 10
	}
	return append(buf, digits[:p]...)
}

// appendGrouped appends u with ',' between groups of three digits.
func appendGrouped(buf []byte, u uint64) []byte {
	var tmp [20]byte
	s := strconv.AppendUint(tmp[:0], u, 10)
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			buf = append(buf, ',')
		}
		buf = append(buf, c)
	}
	return buf
}

// trimFloat removes trailing fraction zeros and a dangling point from a
// strconv 'f' formatted number, and normalises "-0" to "0".
func trimFloat(s string) string {
	var arr [64]byte
	return string(trimFloatBytes(append(arr[:0], s...), 0))
}

// trimFloatBytes applies trimFloat to buf[start:], which holds a strconv 'f'
// formatted number.
func trimFloatBytes(buf []byte, start int) []byte {
	if bytes.IndexByte(buf[start:], '.') >= 0 {
		buf = bytes.TrimRight(buf, "0")
		buf = bytes.TrimSuffix(buf, []byte("."))
	}
	if string(buf[start:]) == "-0" {
		buf = append(buf[:start], '0')
	}
	return buf
}

// appendFloat appends f formatted with prec fraction digits and trimmed like
// trimFloat.
func appendFloat(buf []byte, f float64, prec int) []byte {
	start := len(buf)
	return trimFloatBytes(strconv.AppendFloat(buf, f, 'f', prec, 64), start)
}

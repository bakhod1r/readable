package readable

import "strconv"

// Ordinal formats n as an English ordinal number: the suffix is "th" when the
// last two digits of |n| are 11, 12 or 13, otherwise "st", "nd" or "rd" for a
// last digit of 1, 2 or 3 and "th" for anything else. Negative numbers keep
// their sign; math.MinInt64 is handled without overflow.
//
//	Ordinal(1)   // "1st"
//	Ordinal(12)  // "12th"
//	Ordinal(23)  // "23rd"
//	Ordinal(-2)  // "-2nd"
//	Ordinal(111) // "111th"
func Ordinal(n int64) string {
	neg, abs := absInt64(n)
	suffix := "th"
	if m := abs % 100; m < 11 || m > 13 {
		switch abs % 10 {
		case 1:
			suffix = "st"
		case 2:
			suffix = "nd"
		case 3:
			suffix = "rd"
		}
	}
	var arr [24]byte
	buf := arr[:0]
	if neg {
		buf = append(buf, '-')
	}
	buf = strconv.AppendUint(buf, abs, 10)
	buf = append(buf, suffix...)
	return string(buf)
}

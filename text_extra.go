package readable

import (
	"math"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var romanNumerals = []struct {
	value  int
	symbol string
}{
	{1000, "M"}, {900, "CM"}, {500, "D"}, {400, "CD"}, {100, "C"}, {90, "XC"},
	{50, "L"}, {40, "XL"}, {10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"},
}

// Roman formats n as an upper-case Roman numeral. Values outside 1–3999
// have no standard form and return "".
//
//	Roman(2026) // "MMXXVI"
//	Roman(0)    // ""
func Roman(n int) string {
	if n < 1 || n > 3999 {
		return ""
	}
	var b strings.Builder
	for _, r := range romanNumerals {
		for ; n >= r.value; n -= r.value {
			b.WriteString(r.symbol)
		}
	}
	return b.String()
}

var siPrefixes = []struct {
	scale  float64
	prefix string
}{
	{1e12, "T"}, {1e9, "G"}, {1e6, "M"}, {1e3, "k"}, {1, ""}, {1e-3, "m"}, {1e-6, "µ"}, {1e-9, "n"},
}

// SI formats v with a metric prefix (n, µ, m, k, M, G, T) for unit, with up
// to DefaultPrecision trimmed fraction digits. Values below 1n stay in n;
// zero, NaN and infinities have no prefix.
//
//	SI(0.0012, "A") // "1.2 mA"
//	SI(4700, "Ω")   // "4.7 kΩ"
//	SI(0, "V")      // "0 V"
func SI(v float64, unit string) string {
	a := math.Abs(v)
	if a == 0 || math.IsNaN(v) || math.IsInf(v, 0) {
		return strconv.FormatFloat(v, 'f', -1, 64) + " " + unit
	}
	p := siPrefixes[len(siPrefixes)-1]
	for _, c := range siPrefixes {
		// Compare after rounding so 999.999 m becomes 1 (not 1000 m).
		if math.Round(a/c.scale*100)/100 >= 1 {
			p = c
			break
		}
	}
	s := strconv.FormatFloat(v/p.scale, 'f', DefaultPrecision, 64)
	s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
	if s == "-0" {
		s = "0"
	}
	return s + " " + p.prefix + unit
}

// Initials returns the upper-case first letter of up to the first two
// words of name.
//
//	Initials("john ronald tolkien") // "JR"
//	Initials("  ")                  // ""
func Initials(name string) string {
	var b strings.Builder
	for i, w := range strings.Fields(name) {
		if i == 2 {
			break
		}
		r, _ := utf8.DecodeRuneInString(w)
		b.WriteRune(unicode.ToUpper(r))
	}
	return b.String()
}

// Slug turns s into a lower-case URL slug of ASCII letters and digits joined
// by single "-". Other characters separate words and are dropped.
//
//	Slug("Hello, World! 2026") // "hello-world-2026"
func Slug(s string) string {
	var b strings.Builder
	sep := false
	for _, r := range s {
		r = unicode.ToLower(r)
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			if sep && b.Len() > 0 {
				b.WriteByte('-')
			}
			b.WriteRune(r)
			sep = false
			continue
		}
		sep = true
	}
	return b.String()
}

// Ellipsis shortens s to at most max runes by replacing its middle with "…".
// Strings that fit are returned unchanged; max < 1 returns "".
//
//	Ellipsis("/usr/local/share/readable/docs", 16) // "/usr/loc…le/docs"
func Ellipsis(s string, max int) string {
	n := utf8.RuneCountInString(s)
	switch {
	case max < 1:
		return ""
	case n <= max:
		return s
	}
	r := []rune(s)
	head := max / 2
	tail := max - 1 - head
	return string(r[:head]) + "…" + string(r[n-tail:])
}

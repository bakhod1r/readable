package readable

import (
	"math/big"
	"strings"
)

// ParseNumber parses a compact number as printed by Number: "12.5K",
// "1.23M", "-1.5B", "2T". Plain integers with "," group separators such as
// "1,234,567" are accepted. Suffixes are case-insensitive and may follow a
// space. Results are rounded to the nearest integer, half away from zero.
//
//	ParseNumber("12.5K")     // 12500, nil
//	ParseNumber("1,234,567") // 1234567, nil
//
// Errors wrap ErrSyntax, ErrUnit or ErrRange (outside int64).
func ParseNumber(s string) (int64, error) {
	n, err := parseScaled(s, numberUnit)
	if err == nil && !fitsInt64(n) {
		err = ErrRange
	}
	if err != nil {
		return 0, parseError("ParseNumber", s, err)
	}
	return n.Int64(), nil
}

func numberUnit(u string) (*big.Int, bool) {
	for _, du := range decimalUnits {
		if strings.EqualFold(u, du.suffix) {
			return new(big.Int).SetUint64(du.size), true
		}
	}
	return nil, false
}

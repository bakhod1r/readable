package readable

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
)

// Errors returned (wrapped) by the Parse functions. Test with errors.Is.
var (
	// ErrSyntax reports input that is not a number followed by a unit.
	ErrSyntax = errors.New("invalid syntax")
	// ErrUnit reports a missing or unknown unit.
	ErrUnit = errors.New("unknown unit")
	// ErrRange reports a value that does not fit the result type.
	ErrRange = errors.New("value out of range")
)

func parseError(fn, s string, err error) error {
	return fmt.Errorf("readable: %s %q: %w", fn, s, err)
}

// scanNumber reads an unsigned decimal number starting at s[i]: digits with
// optional "," group separators and at most one ".". It returns the value
// and the index after it, or ok=false if no digit was found.
func scanNumber(s string, i int) (r *big.Rat, next int, ok bool) {
	var b strings.Builder
	digits, dot := 0, false
	for ; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= '0' && c <= '9':
			b.WriteByte(c)
			digits++
		case c == ',' && digits > 0 && !dot:
		case c == '.' && !dot:
			b.WriteByte(c)
			dot = true
		default:
			return finishNumber(b.String(), digits, i)
		}
	}
	return finishNumber(b.String(), digits, i)
}

func finishNumber(num string, digits, next int) (*big.Rat, int, bool) {
	if digits == 0 {
		return nil, next, false
	}
	if strings.HasSuffix(num, ".") {
		num += "0"
	}
	r, _ := new(big.Rat).SetString(num) // digits and one dot always parse
	return r, next, true
}

// parseScaled parses "<sign><number><space><unit>", looks the unit up with
// unitOf and returns the product rounded half away from zero.
func parseScaled(s string, unitOf func(string) (*big.Int, bool)) (*big.Int, error) {
	t := strings.TrimSpace(s)
	neg := false
	if t != "" && (t[0] == '-' || t[0] == '+') {
		neg = t[0] == '-'
		t = t[1:]
	}
	r, i, ok := scanNumber(t, 0)
	if !ok {
		return nil, ErrSyntax
	}
	mult, ok := unitOf(strings.TrimSpace(t[i:]))
	if !ok {
		return nil, ErrUnit
	}
	r.Mul(r, new(big.Rat).SetInt(mult))
	if neg {
		r.Neg(r)
	}
	return roundRat(r), nil
}

var half = big.NewRat(1, 2)

// roundRat rounds r to the nearest integer, half away from zero.
func roundRat(r *big.Rat) *big.Int {
	a := new(big.Rat).Abs(r)
	a.Add(a, half)
	n := new(big.Int).Quo(a.Num(), a.Denom())
	if r.Sign() < 0 {
		n.Neg(n)
	}
	return n
}

var (
	bigMaxInt64  = big.NewInt(math.MaxInt64)
	bigMinInt64  = big.NewInt(math.MinInt64)
	bigMaxUint64 = new(big.Int).SetUint64(math.MaxUint64)
)

func fitsInt64(n *big.Int) bool {
	return n.Cmp(bigMinInt64) >= 0 && n.Cmp(bigMaxInt64) <= 0
}

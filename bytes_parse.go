package readable

import (
	"math/big"
	"strings"
)

// ParseBytes parses a byte size such as "1.5 KB", "512B", "2 GiB" or
// "1,024 bytes". KB, MB, GB, TB, PB and EB are 1024-based, as printed by
// Bytes; the IEC forms KiB...EiB are accepted too. Units are
// case-insensitive, the space before them is optional and a bare number is
// bytes. Results are rounded to the nearest byte, half up.
//
//	ParseBytes("1.5 KB") // 1536, nil
//	ParseBytes("1 MiB")  // 1048576, nil
//
// Errors wrap ErrSyntax, ErrUnit or ErrRange (negative or above MaxUint64).
func ParseBytes(s string) (uint64, error) {
	return parseBytes("ParseBytes", s, 1024)
}

// ParseBytesSI is like ParseBytes but reads kB, MB, GB, TB, PB and EB as
// 1000-based, as printed by BytesSI. KiB...EiB stay 1024-based.
//
//	ParseBytesSI("1 kB")  // 1000, nil
//	ParseBytesSI("1 KiB") // 1024, nil
func ParseBytesSI(s string) (uint64, error) {
	return parseBytes("ParseBytesSI", s, 1000)
}

func parseBytes(fn, s string, base int64) (uint64, error) {
	n, err := parseScaled(s, func(u string) (*big.Int, bool) { return byteUnit(u, base) })
	if err == nil && (n.Sign() < 0 || n.Cmp(bigMaxUint64) > 0) {
		err = ErrRange
	}
	if err != nil {
		return 0, parseError(fn, s, err)
	}
	return n.Uint64(), nil
}

func byteUnit(u string, base int64) (*big.Int, bool) {
	u = strings.ToLower(u)
	switch u {
	case "", "b", "byte", "bytes":
		return big.NewInt(1), true
	}
	exp := strings.IndexByte("kmgtpe", u[0]) + 1
	switch {
	case exp == 0:
		return nil, false
	case u[1:] == "ib":
		base = 1024
	case u[1:] != "b" && u[1:] != "":
		return nil, false
	}
	b := big.NewInt(base)
	return b.Exp(b, big.NewInt(int64(exp)), nil), true
}

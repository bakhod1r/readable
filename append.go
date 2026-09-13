package readable

import "time"

// The AppendX functions append the output of the matching X formatter to dst
// and return the extended slice. They never allocate when dst has enough
// spare capacity, which makes them suitable for hot logging and metrics
// paths:
//
//	buf := make([]byte, 0, 64)
//	buf = readable.AppendNumber(buf, 1234567) // "1.23M"
//	buf = append(buf, ' ')
//	buf = readable.AppendBytes(buf, 1536)     // "1.23M 1.5 KB"
//
// Output is byte-for-byte identical to the string-returning formatter.

// AppendNumber appends n formatted like Number.
func AppendNumber(dst []byte, n int64) []byte {
	neg, abs := absInt64(n)
	return appendUnits(dst, neg, abs, decimalUnits, DefaultPrecision, false, ".", "")
}

// AppendNumberFloat appends f formatted like NumberFloat.
func AppendNumberFloat(dst []byte, f float64) []byte {
	return appendNumberFloat(dst, f)
}

// AppendBytes appends n formatted like Bytes.
func AppendBytes(dst []byte, n uint64) []byte {
	return appendUnits(dst, false, n, binaryUnits, DefaultPrecision, false, ".", " ")
}

// AppendBytesIEC appends n formatted like BytesIEC.
func AppendBytesIEC(dst []byte, n uint64) []byte {
	return appendUnits(dst, false, n, iecUnits, DefaultPrecision, false, ".", " ")
}

// AppendBytesSI appends n formatted like BytesSI.
func AppendBytesSI(dst []byte, n uint64) []byte {
	return appendUnits(dst, false, n, siUnits, DefaultPrecision, false, ".", " ")
}

// AppendDuration appends d formatted like Duration.
func AppendDuration(dst []byte, d time.Duration) []byte {
	return appendDuration(dst, d, 0, false)
}

// AppendPercent appends ratio formatted like Percent.
func AppendPercent(dst []byte, ratio float64) []byte {
	return appendPercent(dst, ratio, DefaultPrecision)
}

// AppendMoney appends amount formatted like Money.
func AppendMoney(dst []byte, amount int64, currency string) []byte {
	c, _ := lookupCurrency(currency)
	return appendMoney(dst, amount, currency, c.exponent, "")
}

// AppendMoneySymbol appends amount formatted like MoneySymbol.
func AppendMoneySymbol(dst []byte, amount int64, currency string) []byte {
	c, _ := lookupCurrency(currency)
	return appendMoney(dst, amount, currency, c.exponent, c.symbol)
}

// AppendOrdinal appends n formatted like Ordinal.
func AppendOrdinal(dst []byte, n int64) []byte {
	return appendOrdinal(dst, n)
}

// Range formats the inclusive span between lo and hi like Number, joined by
// an en dash. The bounds are swapped if hi < lo, and equal bounds collapse to
// a single value.
//
//	Range(1200, 1800) // "1.2K–1.8K"
//	Range(900, 1.5e6) // "900–1.5M"
//	Range(1200, 1200) // "1.2K"
func Range(lo, hi int64) string {
	if hi < lo {
		lo, hi = hi, lo
	}
	var arr [64]byte
	buf := AppendNumber(arr[:0], lo)
	if lo != hi {
		buf = append(buf, "–"...)
		buf = AppendNumber(buf, hi)
	}
	return string(buf)
}

// ByteSize is a byte count that prints like Bytes, for use with fmt and
// structured loggers.
//
//	fmt.Println(readable.ByteSize(1536)) // 1.5 KB
type ByteSize uint64

// String formats b like Bytes.
func (b ByteSize) String() string { return Bytes(uint64(b)) }

// Percentage is a ratio that prints like Percent: 0.15 means 15%.
//
//	fmt.Println(readable.Percentage(0.9234)) // 92.34%
type Percentage float64

// String formats p like Percent.
func (p Percentage) String() string { return Percent(float64(p)) }

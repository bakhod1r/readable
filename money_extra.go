package readable

// moneyCompactUnits returns the K/M/B/T units scaled to minor units for a
// currency with exp minor digits. exp is at most 4 for known currencies, so
// 1e12*10^exp fits in uint64.
func moneyCompactUnits(exp int) [4]unit {
	pow := pow10[exp]
	return [4]unit{
		{1e3 * pow, "K"},
		{1e6 * pow, "M"},
		{1e9 * pow, "B"},
		{1e12 * pow, "T"},
	}
}

// MoneyCompact formats amount, expressed in the currency's minor units, with
// the major value compacted like Number (K, M, B, T, at most DefaultPrecision
// trimmed fraction digits, rounded half up using integer math). A known
// symbol is prefixed; otherwise the code is appended. Major values below
// 1000 are formatted like MoneySymbol.
//
//	MoneyCompact(150000000, "USD") // "$1.5M"
//	MoneyCompact(1500000, "UZS")   // "1.5M UZS"
//	MoneyCompact(99999, "USD")     // "$999.99"
func MoneyCompact(amount int64, currency string) string {
	c, _ := lookupCurrency(currency)
	neg, abs := absInt64(amount)
	units := moneyCompactUnits(c.exponent)
	if abs < units[0].size {
		return formatMoney(amount, currency, c.exponent, c.symbol)
	}
	var arr [64]byte
	buf := arr[:0]
	if neg {
		buf = append(buf, '-')
	}
	buf = append(buf, c.symbol...)
	buf = appendUnits(buf, false, abs, units[:], DefaultPrecision, false, ".", "")
	if c.symbol == "" {
		buf = appendCode(buf, currency)
	}
	return string(buf)
}

// MoneyAccounting is like Money but renders negative amounts in accounting
// style, wrapped in parentheses without a minus sign.
//
//	MoneyAccounting(-1500000, "UZS") // "(1,500,000 UZS)"
//	MoneyAccounting(1500000, "UZS")  // "1,500,000 UZS"
func MoneyAccounting(amount int64, currency string) string {
	if amount >= 0 {
		return Money(amount, currency)
	}
	c, _ := lookupCurrency(currency)
	_, abs := absInt64(amount)
	var arr [64]byte
	buf := appendMoneyUnsigned(append(arr[:0], '('), abs, currency, c.exponent)
	return string(append(buf, ')'))
}

// MoneyChange is like Money but always shows the direction of change:
// positive amounts get "+", negative amounts get the typographic minus sign
// U+2212 ("−"), and zero has no sign.
//
//	MoneyChange(1500000, "UZS")  // "+1,500,000 UZS"
//	MoneyChange(-1500000, "UZS") // "−1,500,000 UZS"
func MoneyChange(amount int64, currency string) string {
	c, _ := lookupCurrency(currency)
	neg, abs := absInt64(amount)
	var arr [64]byte
	buf := arr[:0]
	switch {
	case neg:
		buf = append(buf, "−"...)
	case abs > 0:
		buf = append(buf, '+')
	}
	return string(appendMoneyUnsigned(buf, abs, currency, c.exponent))
}

// appendMoneyUnsigned appends a magnitude formatted like Money.
func appendMoneyUnsigned(buf []byte, abs uint64, code string, exp int) []byte {
	pow := pow10[exp]
	buf = appendGrouped(buf, abs/pow, ',')
	buf = appendFraction(buf, abs%pow, exp, true, ".")
	return appendCode(buf, code)
}

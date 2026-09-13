package readable

type currencyInfo struct {
	exponent int
	symbol   string
}

// lookupCurrency returns the minor digits and symbol of a known ISO 4217 code,
// matched ASCII case-insensitively without allocating. Every known exponent
// must stay <= 2 so moneyCompactUnits cannot overflow.
func lookupCurrency(code string) (currencyInfo, bool) {
	if len(code) != 3 {
		return currencyInfo{}, false
	}
	var k [3]byte
	for i := range k {
		c := code[i]
		if 'a' <= c && c <= 'z' {
			c -= 'a' - 'A'
		}
		k[i] = c
	}
	switch string(k[:]) {
	case "USD":
		return currencyInfo{2, "$"}, true
	case "EUR":
		return currencyInfo{2, "€"}, true
	case "GBP":
		return currencyInfo{2, "£"}, true
	case "UZS":
		return currencyInfo{0, ""}, true
	case "RUB":
		return currencyInfo{2, "₽"}, true
	case "KZT":
		return currencyInfo{2, "₸"}, true
	case "CNY":
		return currencyInfo{2, "CN¥"}, true
	case "JPY":
		return currencyInfo{0, "¥"}, true
	case "KRW":
		return currencyInfo{0, "₩"}, true
	}
	return currencyInfo{}, false
}

// Money formats amount, expressed in the currency's minor units, with
// thousands separators followed by the currency code. Using integer minor
// units avoids floating-point error. Codes are matched ASCII
// case-insensitively but printed as given.
//
// Known currencies use their conventional number of minor digits: USD, EUR,
// GBP, RUB, KZT and CNY have 2; UZS, JPY and KRW have 0. Any other code (for
// example "USDT") is accepted and treated as having 0 minor digits; use
// MoneyWithPrecision to choose explicitly.
//
//	Money(150000000, "UZS") // "150,000,000 UZS"
//	Money(150050, "USD")    // "1,500.50 USD"
func Money(amount int64, currency string) string {
	c, _ := lookupCurrency(currency)
	return formatMoney(amount, currency, c.exponent, "")
}

// MoneySymbol is like Money but prefixes the currency symbol when one is
// known ("$1,500.50"). Currencies without a symbol fall back to the code.
func MoneySymbol(amount int64, currency string) string {
	c, _ := lookupCurrency(currency)
	return formatMoney(amount, currency, c.exponent, c.symbol)
}

// MoneyWithPrecision is like Money but amount has exactly minorDigits minor
// digits (clamped to [0, MaxPrecision]), whatever the currency.
//
//	MoneyWithPrecision(1234567, "USDT", 4) // "123.4567 USDT"
func MoneyWithPrecision(amount int64, currency string, minorDigits int) string {
	return formatMoney(amount, currency, clampPrecision(minorDigits), "")
}

// formatMoney renders amount minor units with exp minor digits, prefixed by
// symbol or, when symbol is empty, followed by code.
func formatMoney(amount int64, code string, exp int, symbol string) string {
	neg, abs := absInt64(amount)
	var arr [64]byte
	buf := arr[:0]
	if neg {
		buf = append(buf, '-')
	}
	buf = append(buf, symbol...)
	pow := pow10[exp]
	buf = appendGrouped(buf, abs/pow, ',')
	buf = appendFraction(buf, abs%pow, exp, true, ".")
	if symbol == "" {
		buf = appendCode(buf, code)
	}
	return string(buf)
}

// appendCode appends " "+code, or nothing for an empty code.
func appendCode(buf []byte, code string) []byte {
	if code == "" {
		return buf
	}
	return append(append(buf, ' '), code...)
}

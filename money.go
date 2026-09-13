package readable

type currencyInfo struct {
	exponent int
	symbol   string
}

// iso4217 lists the active ISO 4217 codes by number of minor digits. Every
// exponent must stay <= 4 so moneyCompactUnits cannot overflow. UZS is listed
// with 0 digits: the tiyin is not used in practice.
var iso4217 = [...]string{
	0: "BIF CLP DJF GNF ISK JPY KMF KRW PYG RWF UGX UYI UZS VND VUV XAF XOF XPF",
	2: "AED AFN ALL AMD ANG AOA ARS AUD AWG AZN BAM BBD BDT BGN BMD BND BOB BOV " +
		"BRL BSD BTN BWP BYN BZD CAD CDF CHE CHF CHW CNY COP COU CRC CUP CVE CZK " +
		"DKK DOP DZD EGP ERN ETB EUR FJD FKP GBP GEL GHS GIP GMD GTQ GYD HKD HNL " +
		"HTG HUF IDR ILS INR IRR JMD KES KGS KHR KPW KYD KZT LAK LBP LKR LRD LSL " +
		"MAD MDL MGA MKD MMK MNT MOP MRU MUR MVR MWK MXN MXV MYR MZN NAD NGN NIO " +
		"NOK NPR NZD PAB PEN PGK PHP PKR PLN QAR RON RSD RUB SAR SBD SCR SDG SEK " +
		"SGD SHP SLE SOS SRD SSP STN SVC SYP SZL THB TJS TMT TOP TRY TTD TWD TZS " +
		"UAH USD USN UYU VED VES WST XCD XCG YER ZAR ZMW ZWG",
	3: "BHD IQD JOD KWD LYD OMR TND",
	4: "CLF UYW",
}

// currencies maps an upper-case ISO 4217 code to its minor digits and symbol.
var currencies = func() map[string]currencyInfo {
	symbols := map[string]string{
		"USD": "$", "EUR": "€", "GBP": "£", "RUB": "₽", "KZT": "₸",
		"CNY": "CN¥", "JPY": "¥", "KRW": "₩",
	}
	m := make(map[string]currencyInfo, 200)
	for exp, codes := range iso4217 {
		for i := 0; i+3 <= len(codes); i += 4 {
			code := codes[i : i+3]
			m[code] = currencyInfo{exp, symbols[code]}
		}
	}
	return m
}()

// lookupCurrency returns the minor digits and symbol of a known ISO 4217 code,
// matched ASCII case-insensitively without allocating.
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
	c, ok := currencies[string(k[:])]
	return c, ok
}

// Money formats amount, expressed in the currency's minor units, with
// thousands separators followed by the currency code. Using integer minor
// units avoids floating-point error. Codes are matched ASCII
// case-insensitively but printed as given.
//
// Active ISO 4217 codes use their standard number of minor digits (USD and
// CHF have 2, JPY has 0, BHD has 3), except UZS, which has 0. Any other code
// (for example "USDT") is accepted and treated as having 0 minor digits; use
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
	var arr [64]byte
	return string(appendMoney(arr[:0], amount, code, exp, symbol))
}

// appendMoney appends amount formatted like formatMoney.
func appendMoney(buf []byte, amount int64, code string, exp int, symbol string) []byte {
	neg, abs := absInt64(amount)
	if neg {
		buf = append(buf, '-')
	}
	buf = append(buf, symbol...)
	pow := pow10[exp]
	buf = appendGrouped(buf, abs/pow)
	buf = appendFraction(buf, abs%pow, exp, true, ".")
	if symbol == "" {
		buf = appendCode(buf, code)
	}
	return buf
}

// appendCode appends " "+code, or nothing for an empty code.
func appendCode(buf []byte, code string) []byte {
	if code == "" {
		return buf
	}
	return append(append(buf, ' '), code...)
}

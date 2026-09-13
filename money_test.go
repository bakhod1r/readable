package readable

import (
	"math"
	"testing"
)

func TestMoney(t *testing.T) {
	tests := []struct {
		amount int64
		code   string
		want   string
	}{
		{150050, "USD", "1,500.50 USD"},
		{0, "USD", "0.00 USD"},
		{5, "USD", "0.05 USD"},
		{-150050, "USD", "-1,500.50 USD"},
		{150000000, "UZS", "150,000,000 UZS"},
		{1000, "JPY", "1,000 JPY"},
		{1234567, "USDT", "1,234,567 USDT"},
		{150050, "usd", "1,500.50 usd"},
		{100, "", "100"},
		{math.MinInt64, "USD", "-92,233,720,368,547,758.08 USD"},
		{math.MaxInt64, "UZS", "9,223,372,036,854,775,807 UZS"},
	}
	for _, tt := range tests {
		if got := Money(tt.amount, tt.code); got != tt.want {
			t.Errorf("Money(%d, %q) = %q, want %q", tt.amount, tt.code, got, tt.want)
		}
	}
}

func TestMoneySymbol(t *testing.T) {
	tests := []struct {
		amount int64
		code   string
		want   string
	}{
		{150050, "USD", "$1,500.50"},
		{-150050, "USD", "-$1,500.50"},
		{150050, "eur", "€1,500.50"},
		{1000, "JPY", "¥1,000"},
		{1000, "UZS", "1,000 UZS"},
		{1000, "USDT", "1,000 USDT"},
		{100, "CNY", "CN¥1.00"},
	}
	for _, tt := range tests {
		if got := MoneySymbol(tt.amount, tt.code); got != tt.want {
			t.Errorf("MoneySymbol(%d, %q) = %q, want %q", tt.amount, tt.code, got, tt.want)
		}
	}
}

func TestMoneyWithPrecision(t *testing.T) {
	runStrCases(t, []strCase{
		{"4", MoneyWithPrecision(1234567, "USDT", 4), "123.4567 USDT"},
		{"0", MoneyWithPrecision(1234567, "USD", 0), "1,234,567 USD"},
		{"neg clamp", MoneyWithPrecision(5, "X", -3), "5 X"},
		{"clamp 9", MoneyWithPrecision(1, "BTC", 20), "0.000000001 BTC"},
		{"min", MoneyWithPrecision(math.MinInt64, "X", 9), "-9,223,372,036.854775808 X"},
	})
}

func TestLookupCurrency(t *testing.T) {
	tests := []struct {
		code   string
		ok     bool
		exp    int
		symbol string
	}{
		{"USD", true, 2, "$"}, {"usd", true, 2, "$"}, {"uSd", true, 2, "$"},
		{"EUR", true, 2, "€"}, {"GBP", true, 2, "£"}, {"UZS", true, 0, ""},
		{"RUB", true, 2, "₽"}, {"KZT", true, 2, "₸"}, {"CNY", true, 2, "CN¥"},
		{"jpy", true, 0, "¥"}, {"KRW", true, 0, "₩"},
		{"", false, 0, ""}, {"US", false, 0, ""}, {"USDT", false, 0, ""},
		{"U$D", false, 0, ""}, {"ＵSD", false, 0, ""}, {"US\x00", false, 0, ""},
		{"USD\x00", false, 0, ""}, {"\xd5SD", false, 0, ""},
	}
	for _, tt := range tests {
		c, ok := lookupCurrency(tt.code)
		if ok != tt.ok || c.exponent != tt.exp || c.symbol != tt.symbol {
			t.Errorf("lookupCurrency(%q) = %+v,%v want {%d %q},%v", tt.code, c, ok, tt.exp, tt.symbol, tt.ok)
		}
	}
}

// TestKnownCurrencyExponents enumerates every three-letter code so that a new
// currency with more than two minor digits, which would overflow
// moneyCompactUnits (1e12*10^exp), cannot be added unnoticed.
func TestKnownCurrencyExponents(t *testing.T) {
	const maxCompactExponent = 2
	known := 0
	for i := range 26 * 26 * 26 {
		code := string([]byte{byte('A' + i/676), byte('A' + i/26%26), byte('A' + i%26)})
		c, ok := lookupCurrency(code)
		if !ok {
			continue
		}
		known++
		if c.exponent < 0 || c.exponent > maxCompactExponent {
			t.Errorf("%s exponent %d outside [0, %d]", code, c.exponent, maxCompactExponent)
		}
	}
	if known != 9 {
		t.Errorf("known currencies = %d, want 9 (update Money docs)", known)
	}
}

func TestLookupCurrencyAllocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { allocCurrency, _ = lookupCurrency("usd") }); n != 0 {
		t.Errorf("lookupCurrency allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { allocSink = Money(150050, "usd") }); n != 1 {
		t.Errorf("Money allocs = %v, want 1", n)
	}
}

package readable

import (
	"math"
	"strings"
	"testing"
)

func TestMoneyCompact(t *testing.T) {
	tests := []struct {
		amount int64
		cur    string
		want   string
	}{
		{1_500_000, "UZS", "1.5M UZS"},
		{150_000_000, "USD", "$1.5M"},
		{-150_000_000, "USD", "-$1.5M"},
		{99_999, "USD", "$999.99"},
		{100_000, "USD", "$1K"},
		{999, "UZS", "999 UZS"},
		{0, "UZS", "0 UZS"},
		{999_999_999, "UZS", "1B UZS"},
		{123_456, "USDT", "123.46K USDT"},
		{math.MinInt64, "USD", "-$92233.72T"},
		{math.MaxInt64, "UZS", "9223372.04T UZS"},
	}
	for _, tt := range tests {
		if got := MoneyCompact(tt.amount, tt.cur); got != tt.want {
			t.Errorf("MoneyCompact(%d, %q) = %q, want %q", tt.amount, tt.cur, got, tt.want)
		}
	}
}

func TestMoneyAccounting(t *testing.T) {
	tests := []struct {
		amount int64
		cur    string
		want   string
	}{
		{-1_500_000, "UZS", "(1,500,000 UZS)"},
		{1_500_000, "UZS", "1,500,000 UZS"},
		{-150050, "USD", "(1,500.50 USD)"},
		{0, "USD", "0.00 USD"},
	}
	for _, tt := range tests {
		if got := MoneyAccounting(tt.amount, tt.cur); got != tt.want {
			t.Errorf("MoneyAccounting(%d, %q) = %q, want %q", tt.amount, tt.cur, got, tt.want)
		}
	}
}

func TestMoneyChange(t *testing.T) {
	tests := []struct {
		amount int64
		cur    string
		want   string
	}{
		{1_500_000, "UZS", "+1,500,000 UZS"},
		{-1_500_000, "UZS", "−1,500,000 UZS"},
		{0, "UZS", "0 UZS"},
		{math.MinInt64, "UZS", "−9,223,372,036,854,775,808 UZS"},
	}
	for _, tt := range tests {
		if got := MoneyChange(tt.amount, tt.cur); got != tt.want {
			t.Errorf("MoneyChange(%d, %q) = %q, want %q", tt.amount, tt.cur, got, tt.want)
		}
	}
}

func FuzzMoneyCompact(f *testing.F) {
	f.Add(int64(150_000_000), "USD")
	f.Add(int64(math.MinInt64), "UZS")
	f.Add(int64(-999), "xyz")
	f.Fuzz(func(t *testing.T, amount int64, cur string) {
		s := MoneyCompact(amount, cur)
		if s == "" {
			t.Fatal("empty")
		}
		if amount < 0 && !strings.HasPrefix(s, "-") {
			t.Fatalf("missing sign: %q", s)
		}
	})
}

func BenchmarkMoneyCompact(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = MoneyCompact(150_000_000, "USD")
	}
}

func TestMoneyExtraAllocs(t *testing.T) {
	for name, f := range map[string]func(){
		"MoneyCompact":      func() { allocSink = MoneyCompact(-150_000_000, "usd") },
		"MoneyCompact code": func() { allocSink = MoneyCompact(1_500_000, "UZS") },
		"MoneyAccounting":   func() { allocSink = MoneyAccounting(-150050, "USD") },
		"MoneyChange":       func() { allocSink = MoneyChange(-150050, "USD") },
	} {
		if n := testing.AllocsPerRun(100, f); n != 1 {
			t.Errorf("%s allocs = %v, want 1", name, n)
		}
	}
}

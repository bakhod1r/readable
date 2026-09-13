package readable

import (
	"testing"
	"time"
)

func BenchmarkNumber(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = Number(1234567)
	}
}

func BenchmarkBytes(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = Bytes(1536)
	}
}

func BenchmarkDuration(b *testing.B) {
	b.ReportAllocs()
	d := 3*time.Hour + 25*time.Minute + 12*time.Second
	for b.Loop() {
		_ = Duration(d)
	}
}

func BenchmarkLatency(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = Latency(1400 * time.Millisecond)
	}
}

func BenchmarkMoney(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = Money(150050, "USD")
	}
}

func BenchmarkPercent(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = Percent(0.1534)
	}
}

func BenchmarkRelativeTime(b *testing.B) {
	b.ReportAllocs()
	now := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	t := now.Add(-2 * time.Hour)
	for b.Loop() {
		_ = RelativeTimeFrom(now, t)
	}
}

func BenchmarkOrdinal(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = Ordinal(123)
	}
}

func BenchmarkMoneyLowercase(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = Money(150050, "usd")
	}
}

func BenchmarkLookupCurrency(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = lookupCurrency("eur")
	}
}

func BenchmarkNumberFloat(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = NumberFloat(1234.5678)
	}
}

func BenchmarkBytesRate(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = BytesRate(5662310, time.Second)
	}
}

func BenchmarkRateWithLabel(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = RateWithLabel(1200, time.Second, "req")
	}
}

func BenchmarkMoneyChange(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = MoneyChange(-1500000, "UZS")
	}
}

func BenchmarkPercentChange(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PercentChange(100, 120)
	}
}

package readable

import (
	"math"
	"testing"
	"time"
)

func TestDurationNatural(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{0, "0 seconds"},
		{125 * time.Millisecond, "125 milliseconds"},
		{time.Second, "1 second"},
		{48 * time.Hour, "2 days"},
		{time.Hour + 5*time.Second, "1 hour and 5 seconds"},
		{49*time.Hour + 32*time.Minute, "2 days, 1 hour and 32 minutes"},
		{49*time.Hour + 32*time.Minute + 7*time.Second, "2 days, 1 hour, 32 minutes and 7 seconds"},
		{-90 * time.Second, "-1 minute and 30 seconds"},
		{-3 * time.Microsecond, "-3 microseconds"},
		{time.Minute + 999*time.Millisecond, "1 minute"},
		{math.MinInt64, "-106751 days, 23 hours, 47 minutes and 16 seconds"},
		{math.MaxInt64, "106751 days, 23 hours, 47 minutes and 16 seconds"},
	}
	for _, tt := range tests {
		if got := DurationNatural(tt.d); got != tt.want {
			t.Errorf("DurationNatural(%d) = %q, want %q", int64(tt.d), got, tt.want)
		}
	}
}

func TestDurationApprox(t *testing.T) {
	day := 24 * time.Hour
	tests := []struct {
		d    time.Duration
		want string
	}{
		{0, "less than a minute"},
		{59 * time.Second, "less than a minute"},
		{-59 * time.Second, "less than a minute"},
		{time.Minute, "about 1 minute"},
		{90 * time.Second, "about 2 minutes"},
		{89 * time.Second, "about 1 minute"},
		{59*time.Minute + 29*time.Second, "about 59 minutes"},
		{59*time.Minute + 30*time.Second, "about 1 hour"},
		{time.Hour, "about 1 hour"},
		{100 * time.Minute, "about 2 hours"},
		{-3 * time.Hour, "about 3 hours"},
		{23*time.Hour + 30*time.Minute, "about 1 day"},
		{49 * time.Hour, "about 2 days"},
		{29*day + 12*time.Hour, "about 1 month"},
		{45 * day, "about 2 months"},
		{364 * day, "about 12 months"},
		{365 * day, "about 1 year"},
		{548 * day, "about 2 years"},
		{math.MinInt64, "about 292 years"},
	}
	for _, tt := range tests {
		if got := DurationApprox(tt.d); got != tt.want {
			t.Errorf("DurationApprox(%v) = %q, want %q", tt.d, got, tt.want)
		}
	}
}

func TestDurationShort(t *testing.T) {
	for _, d := range []time.Duration{0, 125 * time.Millisecond, 90 * time.Second, -49 * time.Hour, math.MinInt64} {
		if got, want := DurationShort(d), Duration(d); got != want {
			t.Errorf("DurationShort(%v) = %q, want %q", d, got, want)
		}
	}
}

func BenchmarkDurationNatural(b *testing.B) {
	b.ReportAllocs()
	d := 49*time.Hour + 32*time.Minute + 7*time.Second
	for b.Loop() {
		_ = DurationNatural(d)
	}
}

func BenchmarkDurationApprox(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = DurationApprox(49 * time.Hour)
	}
}

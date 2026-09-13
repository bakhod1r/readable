package readable

import (
	"math"
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	tests := []struct {
		d          time.Duration
		short, lng string
	}{
		{0, "0s", "0 seconds"},
		{1, "1ns", "1 nanosecond"},
		{999, "999ns", "999 nanoseconds"},
		{time.Microsecond, "1µs", "1 microsecond"},
		{1500 * time.Microsecond, "1ms", "1 millisecond"},
		{999 * time.Millisecond, "999ms", "999 milliseconds"},
		{time.Second, "1s", "1 second"},
		{1999 * time.Millisecond, "1s", "1 second"},
		{59 * time.Second, "59s", "59 seconds"},
		{60 * time.Second, "1m", "1 minute"},
		{90 * time.Second, "1m 30s", "1 minute, 30 seconds"},
		{59 * time.Minute, "59m", "59 minutes"},
		{60 * time.Minute, "1h", "1 hour"},
		{3*time.Hour + 25*time.Minute + 12*time.Second, "3h 25m 12s", "3 hours, 25 minutes, 12 seconds"},
		{24*time.Hour + time.Second, "1d 1s", "1 day, 1 second"},
		{-90 * time.Second, "-1m 30s", "-1 minute, 30 seconds"},
		{-5 * time.Millisecond, "-5ms", "-5 milliseconds"},
		{math.MaxInt64, "106751d 23h 47m 16s", "106751 days, 23 hours, 47 minutes, 16 seconds"},
		{math.MinInt64, "-106751d 23h 47m 16s", "-106751 days, 23 hours, 47 minutes, 16 seconds"},
	}
	for _, tt := range tests {
		if got := Duration(tt.d); got != tt.short {
			t.Errorf("Duration(%d) = %q, want %q", tt.d, got, tt.short)
		}
		if got := DurationLong(tt.d); got != tt.lng {
			t.Errorf("DurationLong(%d) = %q, want %q", tt.d, got, tt.lng)
		}
	}
}

func TestDurationWithOptions(t *testing.T) {
	d := 3*24*time.Hour + 12*time.Hour + 32*time.Minute + 5*time.Second
	runStrCases(t, []strCase{
		{"2", DurationWithOptions(d, DurationOptions{Units: 2}), "3d 12h"},
		{"1", DurationWithOptions(d, DurationOptions{Units: 1}), "3d"},
		{"0 all", DurationWithOptions(d, DurationOptions{Units: 0}), "3d 12h 32m 5s"},
		{"-1 all", DurationWithOptions(d, DurationOptions{Units: -1}), "3d 12h 32m 5s"},
		{"99", DurationWithOptions(d, DurationOptions{Units: 99}), "3d 12h 32m 5s"},
		{"gap", DurationWithOptions(time.Hour+5*time.Second, DurationOptions{Units: 2}), "1h"},
		{"subsecond", DurationWithOptions(5*time.Millisecond, DurationOptions{Units: 1}), "5ms"},
		{"long", DurationWithOptions(d, DurationOptions{Units: 2, Long: true}), "3 days, 12 hours"},
		{"long zero", DurationWithOptions(0, DurationOptions{Units: 2, Long: true}), "0 seconds"},
	})
}

func TestLatency(t *testing.T) {
	tests := []struct {
		d    time.Duration
		want string
	}{
		{0, "0ns"}, {999, "999ns"}, {1000, "1µs"}, {1050, "1.1µs"},
		{350 * time.Microsecond, "350µs"}, {999949, "999.9µs"}, {999950, "1ms"},
		{42 * time.Millisecond, "42ms"}, {999 * time.Millisecond, "999ms"},
		{time.Second, "1s"}, {1400 * time.Millisecond, "1.4s"},
		{59 * time.Second, "59s"}, {60 * time.Second, "1m"}, {120 * time.Second, "2m"},
		{59 * time.Minute, "59m"}, {60 * time.Minute, "1h"},
		{-42 * time.Millisecond, "-42ms"},
		{math.MaxInt64, "2562047.8h"}, {math.MinInt64, "-2562047.8h"},
	}
	for _, tt := range tests {
		if got := Latency(tt.d); got != tt.want {
			t.Errorf("Latency(%d) = %q, want %q", tt.d, got, tt.want)
		}
	}
}

// TestDurationDocExamples keeps the GoDoc examples in duration.go and
// duration_natural.go honest.
func TestDurationDocExamples(t *testing.T) {
	tests := []struct{ got, want string }{
		{Duration(3*time.Hour + 25*time.Minute), "3h 25m"},
		{Duration(90 * time.Second), "1m 30s"},
		{Duration(125 * time.Millisecond), "125ms"},
		{Duration(0), "0s"},
		{DurationWithOptions(3*24*time.Hour+12*time.Hour+32*time.Minute, DurationOptions{Units: 2}), "3d 12h"},
		{DurationWithOptions(time.Hour+5*time.Second, DurationOptions{Units: 2}), "1h"},
		{DurationLong(3*time.Hour + 25*time.Minute + 12*time.Second), "3 hours, 25 minutes, 12 seconds"},
		{DurationLong(0), "0 seconds"},
		{DurationWithOptions(3*time.Hour+25*time.Minute+12*time.Second, DurationOptions{Units: 2, Long: true}), "3 hours, 25 minutes"},
		{Latency(350 * time.Microsecond), "350µs"},
		{Latency(42 * time.Millisecond), "42ms"},
		{Latency(1400 * time.Millisecond), "1.4s"},
		{Latency(120 * time.Second), "2m"},
		{DurationNatural(49*time.Hour + 32*time.Minute), "2 days, 1 hour and 32 minutes"},
		{DurationNatural(time.Hour + 5*time.Second), "1 hour and 5 seconds"},
		{DurationNatural(125 * time.Millisecond), "125 milliseconds"},
		{DurationNatural(0), "0 seconds"},
		{DurationApprox(100 * time.Minute), "about 2 hours"},
		{DurationApprox(49 * time.Hour), "about 2 days"},
		{DurationApprox(30 * time.Second), "less than a minute"},
		{DurationApprox(59*time.Minute + 30*time.Second), "about 1 hour"},
		{DurationShort(90 * time.Second), "1m 30s"},
	}
	for i, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("example %d: got %q, want %q", i, tt.got, tt.want)
		}
	}
}

func TestCalendarConstants(t *testing.T) {
	if secondsPerDay*time.Second != 24*time.Hour || secondsPerHour*time.Second != time.Hour ||
		secondsPerMinute*time.Second != time.Minute {
		t.Fatal("calendar constants disagree with package time")
	}
}

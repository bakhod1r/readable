package readable

import (
	"errors"
	"math"
	"strings"
	"testing"
	"time"
)

func TestParseBytes(t *testing.T) {
	tests := []struct {
		in   string
		want uint64
		err  error
	}{
		{"512 B", 512, nil},
		{"1.5 KB", 1536, nil},
		{"1 MB", 1 << 20, nil},
		{"1mib", 1 << 20, nil},
		{"2 GiB", 2 << 30, nil},
		{"1,024 bytes", 1024, nil},
		{"1 byte", 1, nil},
		{"42", 42, nil},
		{"  3k ", 3072, nil},
		{"+1.5 KB", 1536, nil},
		{"1. KB", 1024, nil},
		{".5 KB", 512, nil},
		{"0.0001 B", 0, nil},
		{"0.5 B", 1, nil},
		{"-0 B", 0, nil},
		{"15.99 EB", 18435214858663483146, nil},
		{"16 EiB", 0, ErrRange},
		{"-1 B", 0, ErrRange},
		{"", 0, ErrSyntax},
		{"KB", 0, ErrSyntax},
		{"-", 0, ErrSyntax},
		{"1 XB", 0, ErrUnit},
		{"1 KBB", 0, ErrUnit},
		{"1 KiBs", 0, ErrUnit},
		{"1.2.3 B", 0, ErrUnit},
		{"1 B B", 0, ErrUnit},
	}
	for _, tt := range tests {
		got, err := ParseBytes(tt.in)
		if got != tt.want || !errors.Is(err, tt.err) {
			t.Errorf("ParseBytes(%q) = %d, %v; want %d, %v", tt.in, got, err, tt.want, tt.err)
		}
	}
}

func TestParseBytesSI(t *testing.T) {
	tests := []struct {
		in   string
		want uint64
	}{
		{"1 kB", 1000}, {"1.5 MB", 1_500_000}, {"1 KiB", 1024}, {"7 B", 7}, {"1 EB", 1e18},
	}
	for _, tt := range tests {
		if got, err := ParseBytesSI(tt.in); got != tt.want || err != nil {
			t.Errorf("ParseBytesSI(%q) = %d, %v; want %d", tt.in, got, err, tt.want)
		}
	}
	if _, err := ParseBytesSI("x"); err == nil || !strings.Contains(err.Error(), `readable: ParseBytesSI "x": invalid syntax`) {
		t.Errorf("error = %v", err)
	}
}

func TestParseNumber(t *testing.T) {
	tests := []struct {
		in   string
		want int64
		err  error
	}{
		{"999", 999, nil},
		{"12.5K", 12500, nil},
		{"1.23M", 1_230_000, nil},
		{"-1.5b", -1_500_000_000, nil},
		{"2 T", 2e12, nil},
		{"1,234,567", 1234567, nil},
		{"-0.5", -1, nil},
		{"0.4", 0, nil},
		{"9223372.036854775807T", math.MaxInt64, nil},
		{"-9223372.036854775808T", math.MinInt64, nil},
		{"9223372.036854775808T", 0, ErrRange},
		{"1 KB", 0, ErrUnit},
		{"abc", 0, ErrSyntax},
		{",5", 0, ErrSyntax},
	}
	for _, tt := range tests {
		got, err := ParseNumber(tt.in)
		if got != tt.want || !errors.Is(err, tt.err) {
			t.Errorf("ParseNumber(%q) = %d, %v; want %d, %v", tt.in, got, err, tt.want, tt.err)
		}
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		in   string
		want time.Duration
		err  error
	}{
		{"3h 25m", 3*time.Hour + 25*time.Minute, nil},
		{"2h30m", 2*time.Hour + 30*time.Minute, nil},
		{"3d 12h", 84 * time.Hour, nil},
		{"125ms", 125 * time.Millisecond, nil},
		{"350µs", 350 * time.Microsecond, nil},
		{"350μs", 350 * time.Microsecond, nil},
		{"1.5us", 1500 * time.Nanosecond, nil},
		{"7ns", 7, nil},
		{"1.4s", 1400 * time.Millisecond, nil},
		{"0s", 0, nil},
		{"-1m 30s", -90 * time.Second, nil},
		{"+1H", time.Hour, nil},
		{"2 days, 1 hour, 32 minutes", 49*time.Hour + 32*time.Minute, nil},
		{"2 days, 1 hour and 32 minutes", 49*time.Hour + 32*time.Minute, nil},
		{"1 min 5 sec", 65 * time.Second, nil},
		{"0.5ns", 1, nil},
		{"-0.5ns", -1, nil},
		{"-106751d 23h 47m 16.854775808s", math.MinInt64, nil},
		{"106751d 23h 47m 16.854775808s", 0, ErrRange},
		{"", 0, ErrSyntax},
		{" , and ", 0, ErrSyntax},
		{"h", 0, ErrSyntax},
		{"5", 0, ErrUnit},
		{"5 weeks", 0, ErrUnit},
		{"1h -5m", 0, ErrSyntax},
	}
	for _, tt := range tests {
		got, err := ParseDuration(tt.in)
		if got != tt.want || !errors.Is(err, tt.err) {
			t.Errorf("ParseDuration(%q) = %v, %v; want %v, %v", tt.in, got, err, tt.want, tt.err)
		}
	}
}

// Formatting then parsing returns the value within the formatter's
// documented rounding.
func TestParseRoundTrip(t *testing.T) {
	for _, n := range []uint64{0, 1, 512, 1024, 1536, 1 << 20, 5 << 30, 1 << 60} {
		for _, f := range []func(uint64) string{Bytes, BytesIEC} {
			if got, err := ParseBytes(f(n)); got != n || err != nil {
				t.Errorf("ParseBytes(%q) = %d, %v; want %d", f(n), got, err, n)
			}
		}
	}
	for _, n := range []uint64{0, 999, 1000, 1500, 25e6, 3e18} {
		if got, err := ParseBytesSI(BytesSI(n)); got != n || err != nil {
			t.Errorf("ParseBytesSI(%q) = %d, %v; want %d", BytesSI(n), got, err, n)
		}
	}
	for _, n := range []int64{0, 999, 12500, -1500, 1_230_000, 2e12} {
		if got, err := ParseNumber(Number(n)); got != n || err != nil {
			t.Errorf("ParseNumber(%q) = %d, %v; want %d", Number(n), got, err, n)
		}
	}
	for _, d := range []time.Duration{0, 7, 350 * time.Microsecond, 125 * time.Millisecond,
		49*time.Hour + 32*time.Minute + 5*time.Second, -90 * time.Second} {
		for _, f := range []func(time.Duration) string{Duration, DurationLong, DurationNatural} {
			if got, err := ParseDuration(f(d)); got != d || err != nil {
				t.Errorf("ParseDuration(%q) = %v, %v; want %v", f(d), got, err, d)
			}
		}
	}
}

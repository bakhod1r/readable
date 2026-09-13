package readable

import (
	"math"
	"testing"
	"time"
)

func TestRate(t *testing.T) {
	runStrCases(t, []strCase{
		{"125K", Rate(125000, time.Second), "125K/s"},
		{"per minute", Rate(300, time.Minute), "5/s"},
		{"fraction", Rate(1, 2*time.Second), "0.5/s"},
		{"zero period", Rate(1, 0), "NaN/s"},
		{"neg period", Rate(1, -time.Second), "NaN/s"},
		{"neg count", Rate(-1500, time.Second), "-1.5K/s"},
		{"label", RateWithLabel(1200, time.Second, "req"), "1.2K req/s"},
		{"label nan", RateWithLabel(1, 0, "req"), "NaN req/s"},
		{"inf", Rate(math.Inf(1), time.Second), "+Inf/s"},
	})
}

func TestBytesRate(t *testing.T) {
	runStrCases(t, []strCase{
		{"5.4MB", BytesRate(5662310, time.Second), "5.4 MB/s"},
		{"zero", BytesRate(0, time.Second), "0 B/s"},
		{"zero period", BytesRate(1, 0), "NaN/s"},
		{"neg period", BytesRate(1, -1), "NaN/s"},
		{"overflow clamp", BytesRate(math.MaxUint64, time.Nanosecond), "16 EB/s"},
		{"max per sec", BytesRate(math.MaxUint64, time.Second), "16 EB/s"},
	})
}

func TestBytesPerSecondExact(t *testing.T) {
	const big = 1<<53 + 1 // not representable as float64
	const ns = uint64(time.Second)
	tests := []struct {
		name   string
		bytes  uint64
		period time.Duration
		want   uint64
	}{
		{"above 2^53", big, time.Second, big},
		{"above 2^53 halved", 2 * big, 2 * time.Second, big},
		{"truncates", 10, 3 * time.Second, 3},
		{"max per second", math.MaxUint64, time.Second, math.MaxUint64},
		{"just below saturation", math.MaxUint64 / ns, time.Nanosecond, math.MaxUint64 / ns * ns},
		{"saturates", math.MaxUint64/ns + 1, time.Nanosecond, math.MaxUint64},
		{"saturates max", math.MaxUint64, time.Nanosecond, math.MaxUint64},
	}
	for _, tt := range tests {
		if got := bytesPerSecond(tt.bytes, tt.period); got != tt.want {
			t.Errorf("%s: bytesPerSecond(%d, %v) = %d, want %d", tt.name, tt.bytes, tt.period, got, tt.want)
		}
	}
	if v := uint64(big); uint64(float64(v)) == v {
		t.Fatal("test value must not be exactly representable as float64")
	}
}

func TestRateAllocs(t *testing.T) {
	for name, f := range map[string]func(){
		"BytesRate":     func() { allocSink = BytesRate(5662310, time.Second) },
		"RateWithLabel": func() { allocSink = RateWithLabel(1200, time.Second, "req") },
		"Rate":          func() { allocSink = Rate(1200, time.Second) },
	} {
		if n := testing.AllocsPerRun(100, f); n > 1 {
			t.Errorf("%s allocs = %v, want <= 1", name, n)
		}
	}
}

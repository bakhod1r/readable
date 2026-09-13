package readable

import (
	"errors"
	"math"
	"testing"
)

func TestPercent(t *testing.T) {
	tests := []struct {
		r    float64
		want string
	}{
		{0, "0%"}, {0.00001, "0%"}, {-0.00001, "0%"}, {0.001, "0.1%"}, {0.01, "1%"},
		{0.1, "10%"}, {1, "100%"}, {0.1534, "15.34%"}, {-0.25, "-25%"}, {12, "1200%"},
		{math.NaN(), "NaN%"}, {math.Inf(1), "+Inf%"}, {math.Inf(-1), "-Inf%"},
		{math.MaxFloat64, "+Inf%"},
	}
	for _, tt := range tests {
		if got := Percent(tt.r); got != tt.want {
			t.Errorf("Percent(%v) = %q, want %q", tt.r, got, tt.want)
		}
	}
}

func TestPercentWithPrecision(t *testing.T) {
	runStrCases(t, []strCase{
		{"0", PercentWithPrecision(0.123456, 0), "12%"},
		{"1", PercentWithPrecision(0.123456, 1), "12.3%"},
		{"-1", PercentWithPrecision(0.123456, -1), "12%"},
		{"clamp", PercentWithPrecision(0.1234567891234, 50), "12.345678912%"},
	})
}

func TestPercentChange(t *testing.T) {
	ok := []struct {
		from, to float64
		want     string
	}{
		{100, 120, "20%"}, {100, 80, "-20%"}, {100, 100, "0%"},
		{-100, -50, "50%"}, {-100, -150, "-50%"}, {50, 0, "-100%"},
	}
	for _, tt := range ok {
		got, err := PercentChange(tt.from, tt.to)
		if err != nil || got != tt.want {
			t.Errorf("PercentChange(%v,%v) = %q,%v want %q", tt.from, tt.to, got, err, tt.want)
		}
	}
	bad := [][2]float64{
		{0, 1}, {math.Copysign(0, -1), 1}, {0, 0}, {math.NaN(), 1}, {1, math.NaN()},
		{math.Inf(1), 1}, {1, math.Inf(-1)}, {5e-324, math.MaxFloat64},
	}
	for _, b := range bad {
		got, err := PercentChange(b[0], b[1])
		if !errors.Is(err, ErrUndefined) || got != "" {
			t.Errorf("PercentChange(%v,%v) = %q,%v want ErrUndefined", b[0], b[1], got, err)
		}
	}
}

func TestPercentAllocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { allocSink = Percent(0.1534) }); n != 1 {
		t.Errorf("Percent allocs = %v, want 1", n)
	}
}

func TestPercentNegativeZero(t *testing.T) {
	runStrCases(t, []strCase{
		{"-0.0001 p0", PercentWithPrecision(-0.001, 0), "0%"},
		{"-0", Percent(math.Copysign(0, -1)), "0%"},
		{"huge", PercentWithPrecision(1e15, 0), "100000000000000000%"},
	})
}

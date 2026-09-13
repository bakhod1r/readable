package readable

import (
	"math"
	"testing"
)

type strCase struct {
	name string
	got  string
	want string
}

func runStrCases(t *testing.T, cases []strCase) {
	t.Helper()
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.name, c.got, c.want)
		}
	}
}

func TestClampPrecision(t *testing.T) {
	for in, want := range map[int]int{-5: 0, 0: 0, 3: 3, 9: 9, 10: 9, 100: 9} {
		if got := clampPrecision(in); got != want {
			t.Errorf("clampPrecision(%d) = %d, want %d", in, got, want)
		}
	}
}

func TestAbsInt64(t *testing.T) {
	tests := []struct {
		n   int64
		neg bool
		abs uint64
	}{
		{0, false, 0}, {5, false, 5}, {-5, true, 5},
		{math.MaxInt64, false, math.MaxInt64},
		{math.MinInt64, true, 1 << 63},
	}
	for _, tt := range tests {
		neg, abs := absInt64(tt.n)
		if neg != tt.neg || abs != tt.abs {
			t.Errorf("absInt64(%d) = %v,%d want %v,%d", tt.n, neg, abs, tt.neg, tt.abs)
		}
	}
}

func TestDivRound(t *testing.T) {
	tests := []struct {
		abs, size uint64
		prec      int
		q, frac   uint64
	}{
		{1500, 1000, 0, 2, 0},
		{1499, 1000, 0, 1, 0},
		{1005, 1000, 2, 1, 1},
		{1004, 1000, 2, 1, 0},
		{1995, 1000, 2, 2, 0},
		{math.MaxUint64, 1 << 60, 9, 16, 0},
		{math.MaxUint64, 1, 9, math.MaxUint64, 0},
	}
	for _, tt := range tests {
		q, f := divRound(tt.abs, tt.size, tt.prec)
		if q != tt.q || f != tt.frac {
			t.Errorf("divRound(%d,%d,%d) = %d,%d want %d,%d", tt.abs, tt.size, tt.prec, q, f, tt.q, tt.frac)
		}
	}
}

func TestAppendGrouped(t *testing.T) {
	runStrCases(t, []strCase{
		{"0", string(appendGrouped(nil, 0)), "0"},
		{"999", string(appendGrouped(nil, 999)), "999"},
		{"1000", string(appendGrouped(nil, 1000)), "1,000"},
		{"max", string(appendGrouped(nil, math.MaxUint64)), "18,446,744,073,709,551,615"},
	})
}

func TestTrimFloat(t *testing.T) {
	runStrCases(t, []strCase{
		{"1.500", trimFloat("1.500"), "1.5"},
		{"1.000", trimFloat("1.000"), "1"},
		{"10", trimFloat("10"), "10"},
		{"100", trimFloat("100"), "100"},
		{"-0.00", trimFloat("-0.00"), "0"},
		{"-1.20", trimFloat("-1.20"), "-1.2"},
	})
}

// Sinks force results to escape so testing.AllocsPerRun counts real allocations.
var (
	allocSink     string
	allocCurrency currencyInfo
)

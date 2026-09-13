package readable

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

func TestNumber(t *testing.T) {
	tests := []struct {
		n    int64
		want string
	}{
		{0, "0"}, {1, "1"}, {-1, "-1"},
		{999, "999"}, {1000, "1K"}, {1001, "1K"}, {1005, "1.01K"},
		{12500, "12.5K"}, {999500, "999.5K"}, {999994, "999.99K"},
		{999995, "1M"}, {999999, "1M"}, {1000000, "1M"}, {1234567, "1.23M"},
		{-1234567, "-1.23M"}, {-999999, "-1M"},
		{1e9, "1B"}, {1e12, "1T"}, {1e15, "1000T"},
		{math.MaxInt64, "9223372.04T"}, {math.MinInt64, "-9223372.04T"},
	}
	for _, tt := range tests {
		if got := Number(tt.n); got != tt.want {
			t.Errorf("Number(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

func TestNumberWithPrecision(t *testing.T) {
	runStrCases(t, []strCase{
		{"p1", NumberWithPrecision(1234567, 1), "1.2M"},
		{"p3", NumberWithPrecision(1234567, 3), "1.235M"},
		{"p0", NumberWithPrecision(1500, 0), "2K"},
		{"p-1", NumberWithPrecision(1500, -1), "2K"},
		{"p100 clamp", NumberWithPrecision(1234567891, 100), "1.234567891B"},
		{"p9", NumberWithPrecision(1234567891, 9), "1.234567891B"},
		{"small ignores prec", NumberWithPrecision(999, 5), "999"},
	})
}

func TestNumberWithOptions(t *testing.T) {
	runStrCases(t, []strCase{
		{"fixed", NumberWithOptions(1500000, NumberOptions{Precision: 2, FixedPrecision: true}), "1.50M"},
		{"fixed whole", NumberWithOptions(1000000, NumberOptions{Precision: 2, FixedPrecision: true}), "1.00M"},
		{"comma", NumberWithOptions(1500000, NumberOptions{Precision: 2, Decimal: ","}), "1,5M"},
		{"zero value", NumberWithOptions(1500000, NumberOptions{}), "2M"},
		{"small fixed", NumberWithOptions(5, NumberOptions{Precision: 2, FixedPrecision: true}), "5"},
	})
}

func TestNumberFloat(t *testing.T) {
	tests := []struct {
		f    float64
		want string
	}{
		{0, "0"}, {math.Copysign(0, -1), "0"}, {0.5, "0.5"}, {-0.001, "0"},
		{999.99, "999.99"}, {999.999, "1K"}, {1500, "1.5K"}, {-1500, "-1.5K"},
		{999999.9, "1M"}, {1e15, "1000T"},
		{999.995, "1K"}, {999.9949, "999.99"}, {-999.999, "-1K"}, {999999.995, "1M"},
		{0.005, "0.01"}, {-0.004, "0"}, {1e3, "1K"}, {999.994, "999.99"},
		{math.SmallestNonzeroFloat64, "0"}, {-math.SmallestNonzeroFloat64, "0"},
		{math.NaN(), "NaN"}, {math.Inf(1), "+Inf"}, {math.Inf(-1), "-Inf"},
	}
	for _, tt := range tests {
		if got := NumberFloat(tt.f); got != tt.want {
			t.Errorf("NumberFloat(%v) = %q, want %q", tt.f, got, tt.want)
		}
	}
}

func TestNumberFloatHuge(t *testing.T) {
	// Beyond the largest unit the value stays in T with every integer digit.
	got := NumberFloat(1e300)
	if !strings.HasPrefix(got, "1000000000000000007630") || !strings.HasSuffix(got, "T") || strings.ContainsRune(got, '.') {
		t.Errorf("NumberFloat(1e300) = %q", got)
	}
	if got, want := NumberFloat(-math.MaxFloat64), "-"+strconv.FormatFloat(math.MaxFloat64/1e12, 'f', 0, 64)+"T"; got != want {
		t.Errorf("NumberFloat(-MaxFloat64) = %q, want %q", got, want)
	}
}

func TestNumberFloatAllocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { allocSink = NumberFloat(-1234.5678) }); n != 1 {
		t.Errorf("NumberFloat allocs = %v, want 1", n)
	}
}

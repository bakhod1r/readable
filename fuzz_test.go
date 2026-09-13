package readable

import (
	"math"
	"math/big"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

func checkFuzzOutput(t *testing.T, f func() string) {
	t.Helper()
	a, b := f(), f()
	if a == "" {
		t.Fatal("empty output")
	}
	if !utf8.ValidString(a) {
		t.Fatalf("invalid UTF-8: %q", a)
	}
	if a != b {
		t.Fatalf("nondeterministic: %q vs %q", a, b)
	}
}

func FuzzNumber(f *testing.F) {
	for _, s := range []int64{0, 999, 1000, 999999, -1, 1 << 62, -1 << 63} {
		f.Add(s, 2)
	}
	f.Fuzz(func(t *testing.T, n int64, p int) {
		checkFuzzOutput(t, func() string { return Number(n) })
		checkFuzzOutput(t, func() string { return NumberWithPrecision(n, p) })
		checkFuzzOutput(t, func() string { return NumberWithOptions(n, NumberOptions{Precision: p, FixedPrecision: true}) })
	})
}

func FuzzBytes(f *testing.F) {
	for _, s := range []uint64{0, 1023, 1024, 1048575, 1<<64 - 1} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, n uint64) {
		checkFuzzOutput(t, func() string { return Bytes(n) })
		checkFuzzOutput(t, func() string { return BytesIEC(n) })
		checkFuzzOutput(t, func() string { return BytesSI(n) })
		checkFuzzOutput(t, func() string { return FileSize(int64(n)) })
	})
}

func FuzzDuration(f *testing.F) {
	for _, s := range []int64{0, 1, 999999999, 1e9, 1<<63 - 1, -1 << 63} {
		f.Add(s, 2)
	}
	f.Fuzz(func(t *testing.T, n int64, u int) {
		d := time.Duration(n)
		checkFuzzOutput(t, func() string { return Duration(d) })
		checkFuzzOutput(t, func() string { return DurationLong(d) })
		checkFuzzOutput(t, func() string { return DurationWithPrecision(d, u) })
		checkFuzzOutput(t, func() string { return DurationLongWithPrecision(d, u) })
	})
}

func FuzzLatency(f *testing.F) {
	for _, s := range []int64{0, 999, 999950, 1<<63 - 1, -1 << 63} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, n int64) {
		checkFuzzOutput(t, func() string { return Latency(time.Duration(n)) })
	})
}

func FuzzMoney(f *testing.F) {
	f.Add(int64(150050), "USD", 2)
	f.Add(int64(-1<<63), "usdt", 9)
	f.Add(int64(0), "", -1)
	f.Fuzz(func(t *testing.T, n int64, code string, p int) {
		if !utf8.ValidString(code) {
			return
		}
		checkFuzzOutput(t, func() string { return Money(n, code) })
		checkFuzzOutput(t, func() string { return MoneySymbol(n, code) })
		checkFuzzOutput(t, func() string { return MoneyWithPrecision(n, code, p) })
	})
}

func FuzzOrdinal(f *testing.F) {
	for _, s := range []int64{0, 11, 111, 21, -1 << 63} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, n int64) {
		checkFuzzOutput(t, func() string { return Ordinal(n) })
	})
}

func FuzzPercent(f *testing.F) {
	f.Add(0.1534, 100.0, 2)
	f.Add(0.0, 0.0, -1)
	f.Add(1e308, -1e308, 20)
	f.Fuzz(func(t *testing.T, r, from float64, p int) {
		checkFuzzOutput(t, func() string { return Percent(r) })
		checkFuzzOutput(t, func() string { return PercentWithPrecision(r, p) })
		s, err := PercentChange(from, r)
		if (err == nil) == (s == "") {
			t.Fatalf("PercentChange(%v,%v) = %q,%v", from, r, s, err)
		}
	})
}

func FuzzBytesPerSecond(f *testing.F) {
	f.Add(uint64(1<<53+1), int64(time.Second))
	f.Add(uint64(1<<64-1), int64(1))
	f.Add(uint64(0), int64(-1))
	f.Fuzz(func(t *testing.T, n uint64, ns int64) {
		checkFuzzOutput(t, func() string { return BytesRate(n, time.Duration(ns)) })
		if ns <= 0 {
			return
		}
		want := new(big.Int).SetUint64(n)
		want.Mul(want, big.NewInt(int64(time.Second)))
		want.Quo(want, big.NewInt(ns))
		if !want.IsUint64() {
			want.SetUint64(math.MaxUint64)
		}
		if got := bytesPerSecond(n, time.Duration(ns)); got != want.Uint64() {
			t.Fatalf("bytesPerSecond(%d, %d) = %d, want %d", n, ns, got, want)
		}
	})
}

// numberFloatReference is the straightforward string-based NumberFloat.
func numberFloatReference(f float64) string {
	switch {
	case math.IsNaN(f):
		return "NaN"
	case math.IsInf(f, 0):
		return map[bool]string{true: "+Inf", false: "-Inf"}[f > 0]
	}
	abs := math.Abs(f)
	i := 0
	for i+1 < len(decimalUnits) && abs >= float64(decimalUnits[i+1].size) {
		i++
	}
	s := strconv.FormatFloat(abs/float64(decimalUnits[i].size), 'f', DefaultPrecision, 64)
	if v, _ := strconv.ParseFloat(s, 64); v >= 1000 && i+1 < len(decimalUnits) {
		i++
		s = strconv.FormatFloat(abs/float64(decimalUnits[i].size), 'f', DefaultPrecision, 64)
	}
	s = strings.TrimRight(s, "0")
	s = strings.TrimSuffix(s, ".")
	if f < 0 && s != "0" {
		s = "-" + s
	}
	return s + decimalUnits[i].suffix
}

func FuzzNumberFloat(f *testing.F) {
	for _, s := range []float64{0, -0.001, 999.995, 999999.995, 1e300, -math.MaxFloat64, 5e-324} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, x float64) {
		if got, want := NumberFloat(x), numberFloatReference(x); got != want {
			t.Fatalf("NumberFloat(%v) = %q, want %q", x, got, want)
		}
		checkFuzzOutput(t, func() string { return Percent(x) })
	})
}

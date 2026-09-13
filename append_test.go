package readable

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestAppendMatchesString(t *testing.T) {
	prefix := []byte("x=")
	app := func(f func([]byte) []byte) string {
		b := f(append([]byte(nil), prefix...))
		return string(b[len(prefix):])
	}
	durations := []time.Duration{0, 125 * time.Millisecond, 90 * time.Second, -3 * time.Hour, math.MinInt64}
	ints := []int64{0, 999, 12500, -1234567, math.MinInt64, math.MaxInt64}
	floats := []float64{0, 0.5, 1500.25, -0.001, math.NaN(), math.Inf(1), math.Inf(-1), 1e20}
	for _, n := range ints {
		runStrCases(t, []strCase{
			{"Number", app(func(b []byte) []byte { return AppendNumber(b, n) }), Number(n)},
			{"Ordinal", app(func(b []byte) []byte { return AppendOrdinal(b, n) }), Ordinal(n)},
			{"Money", app(func(b []byte) []byte { return AppendMoney(b, n, "USD") }), Money(n, "USD")},
			{"MoneySymbol", app(func(b []byte) []byte { return AppendMoneySymbol(b, n, "EUR") }), MoneySymbol(n, "EUR")},
			{"MoneyNoCode", app(func(b []byte) []byte { return AppendMoney(b, n, "") }), Money(n, "")},
		})
		u := uint64(n)
		runStrCases(t, []strCase{
			{"Bytes", app(func(b []byte) []byte { return AppendBytes(b, u) }), Bytes(u)},
			{"BytesIEC", app(func(b []byte) []byte { return AppendBytesIEC(b, u) }), BytesIEC(u)},
			{"BytesSI", app(func(b []byte) []byte { return AppendBytesSI(b, u) }), BytesSI(u)},
		})
	}
	for _, f := range floats {
		runStrCases(t, []strCase{
			{"NumberFloat", app(func(b []byte) []byte { return AppendNumberFloat(b, f) }), NumberFloat(f)},
			{"Percent", app(func(b []byte) []byte { return AppendPercent(b, f) }), Percent(f)},
		})
	}
	for _, d := range durations {
		runStrCases(t, []strCase{
			{"Duration", app(func(b []byte) []byte { return AppendDuration(b, d) }), Duration(d)},
		})
	}
}

func TestAppendAllocs(t *testing.T) {
	buf := make([]byte, 0, 128)
	cases := map[string]func(){
		"AppendNumber":      func() { buf = AppendNumber(buf[:0], 1234567) },
		"AppendNumberFloat": func() { buf = AppendNumberFloat(buf[:0], 1500.25) },
		"AppendBytes":       func() { buf = AppendBytes(buf[:0], 1536) },
		"AppendBytesIEC":    func() { buf = AppendBytesIEC(buf[:0], 1536) },
		"AppendBytesSI":     func() { buf = AppendBytesSI(buf[:0], 1536) },
		"AppendDuration":    func() { buf = AppendDuration(buf[:0], 3723*time.Second) },
		"AppendPercent":     func() { buf = AppendPercent(buf[:0], 0.9234) },
		"AppendMoney":       func() { buf = AppendMoney(buf[:0], 123450, "USD") },
		"AppendMoneySymbol": func() { buf = AppendMoneySymbol(buf[:0], 123450, "usd") },
		"AppendOrdinal":     func() { buf = AppendOrdinal(buf[:0], 23) },
	}
	for name, f := range cases {
		if n := testing.AllocsPerRun(100, f); n != 0 {
			t.Errorf("%s allocs = %v, want 0", name, n)
		}
	}
}

func TestRange(t *testing.T) {
	runStrCases(t, []strCase{
		{"compact", Range(1200, 1800), "1.2K–1.8K"},
		{"plain", Range(5, 10), "5–10"},
		{"mixed", Range(900, 1_500_000), "900–1.5M"},
		{"negative", Range(-1500, 1500), "-1.5K–1.5K"},
		{"equal", Range(1200, 1200), "1.2K"},
		{"reversed", Range(1800, 1200), "1.2K–1.8K"},
		{"extremes", Range(math.MinInt64, math.MaxInt64), "-9223372.04T–9223372.04T"},
	})
}

func TestStringerTypes(t *testing.T) {
	runStrCases(t, []strCase{
		{"ByteSize", ByteSize(1536).String(), "1.5 KB"},
		{"ByteSize fmt", fmt.Sprint(ByteSize(1 << 20)), "1 MB"},
		{"Percentage", Percentage(0.9234).String(), "92.34%"},
		{"Percentage fmt", fmt.Sprintf("%v", Percentage(1)), "100%"},
	})
	var _ fmt.Stringer = ByteSize(0)
	var _ fmt.Stringer = Percentage(0)
}

func FuzzAppendNumberFloat(f *testing.F) {
	for _, v := range []float64{0, -0.001, 999.995, 1e300, math.NaN(), math.Inf(-1)} {
		f.Add(v)
	}
	f.Fuzz(func(t *testing.T, v float64) {
		if got, want := string(AppendNumberFloat(nil, v)), NumberFloat(v); got != want {
			t.Fatalf("AppendNumberFloat(%v) = %q, want %q", v, got, want)
		}
		if got, want := string(AppendPercent(nil, v)), Percent(v); got != want {
			t.Fatalf("AppendPercent(%v) = %q, want %q", v, got, want)
		}
	})
}

func FuzzAppendInt(f *testing.F) {
	for _, v := range []int64{0, 1, -1, 999_999, math.MinInt64, math.MaxInt64} {
		f.Add(v)
	}
	f.Fuzz(func(t *testing.T, v int64) {
		checks := [][2]string{
			{string(AppendNumber(nil, v)), Number(v)},
			{string(AppendOrdinal(nil, v)), Ordinal(v)},
			{string(AppendMoney(nil, v, "BHD")), Money(v, "BHD")},
			{string(AppendDuration(nil, time.Duration(v))), Duration(time.Duration(v))},
			{string(AppendBytesIEC(nil, uint64(v))), BytesIEC(uint64(v))},
		}
		for i, c := range checks {
			if c[0] != c[1] {
				t.Fatalf("check %d for %d: append %q, string %q", i, v, c[0], c[1])
			}
		}
	})
}

func BenchmarkAppendNumber(b *testing.B) {
	buf := make([]byte, 0, 64)
	b.ReportAllocs()
	for b.Loop() {
		buf = AppendNumber(buf[:0], 1234567)
	}
}

func BenchmarkAppendBytes(b *testing.B) {
	buf := make([]byte, 0, 64)
	b.ReportAllocs()
	for b.Loop() {
		buf = AppendBytes(buf[:0], 1536000)
	}
}

func BenchmarkAppendDuration(b *testing.B) {
	buf := make([]byte, 0, 64)
	b.ReportAllocs()
	for b.Loop() {
		buf = AppendDuration(buf[:0], 3723*time.Second)
	}
}

func BenchmarkAppendPercent(b *testing.B) {
	buf := make([]byte, 0, 64)
	b.ReportAllocs()
	for b.Loop() {
		buf = AppendPercent(buf[:0], 0.9234)
	}
}

func BenchmarkAppendMoney(b *testing.B) {
	buf := make([]byte, 0, 64)
	b.ReportAllocs()
	for b.Loop() {
		buf = AppendMoney(buf[:0], 123450, "USD")
	}
}

func BenchmarkRange(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		allocSink = Range(1200, 1800)
	}
}

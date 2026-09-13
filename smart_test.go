package readable

import (
	"math"
	"testing"
)

func TestNumberWords(t *testing.T) {
	cases := map[int64]string{
		0: "0", 999: "999", 12500: "12.5 thousand", 1234567: "1.23 million",
		1e9: "1 billion", 1_500_000_000_000: "1.5 trillion", 999_999: "1 million",
		-1500: "-1.5 thousand", 1005: "1.01 thousand", math.MinInt64: "-9223372.04 trillion",
	}
	for in, want := range cases {
		if got := NumberWords(in); got != want {
			t.Errorf("NumberWords(%d) = %q, want %q", in, got, want)
		}
	}
}

func BenchmarkNumberWords(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = NumberWords(1234567)
	}
}

package readable

import (
	"math"
	"testing"
	"time"
)

func TestLogDuration(t *testing.T) {
	cases := map[time.Duration]string{1532 * time.Millisecond: "1.532s", time.Hour: "1h0m0s", 0: "0s"}
	for in, want := range cases {
		if got := LogDuration(in); got != want {
			t.Errorf("LogDuration(%v) = %q, want %q", in, got, want)
		}
	}
}

func TestLogBytes(t *testing.T) {
	cases := map[uint64]string{0: "0B", 1024: "1KiB", 12345678: "11.77MiB"}
	for in, want := range cases {
		if got := LogBytes(in); got != want {
			t.Errorf("LogBytes(%d) = %q, want %q", in, got, want)
		}
	}
}

func TestLogNumber(t *testing.T) {
	cases := map[int64]string{1234567: "1234567", -5: "-5", math.MinInt64: "-9223372036854775808"}
	for in, want := range cases {
		if got := LogNumber(in); got != want {
			t.Errorf("LogNumber(%d) = %q, want %q", in, got, want)
		}
	}
}

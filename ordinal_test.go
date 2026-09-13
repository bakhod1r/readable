package readable

import (
	"math"
	"testing"
)

func TestOrdinal(t *testing.T) {
	tests := []struct {
		n    int64
		want string
	}{
		{0, "0th"}, {1, "1st"}, {2, "2nd"}, {3, "3rd"}, {4, "4th"},
		{11, "11th"}, {12, "12th"}, {13, "13th"}, {21, "21st"}, {22, "22nd"}, {23, "23rd"},
		{101, "101st"}, {111, "111th"}, {112, "112th"}, {113, "113th"}, {1001, "1001st"},
		{-1, "-1st"}, {-2, "-2nd"}, {-11, "-11th"},
		{math.MaxInt64, "9223372036854775807th"}, {math.MinInt64, "-9223372036854775808th"},
	}
	for _, tt := range tests {
		if got := Ordinal(tt.n); got != tt.want {
			t.Errorf("Ordinal(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

package readable

import (
	"math"
	"testing"
)

func TestBytes(t *testing.T) {
	tests := []struct {
		n            uint64
		bin, iec, si string
	}{
		{0, "0 B", "0 B", "0 B"},
		{999, "999 B", "999 B", "999 B"},
		{1000, "1000 B", "1000 B", "1 kB"},
		{1023, "1023 B", "1023 B", "1.02 kB"},
		{1024, "1 KB", "1 KiB", "1.02 kB"},
		{1025, "1 KB", "1 KiB", "1.03 kB"},
		{1536, "1.5 KB", "1.5 KiB", "1.54 kB"},
		{1048575, "1 MB", "1 MiB", "1.05 MB"},
		{1048576, "1 MB", "1 MiB", "1.05 MB"},
		{1 << 30, "1 GB", "1 GiB", "1.07 GB"},
		{1 << 60, "1 EB", "1 EiB", "1.15 EB"},
		{math.MaxUint64, "16 EB", "16 EiB", "18.45 EB"},
	}
	for _, tt := range tests {
		if got := Bytes(tt.n); got != tt.bin {
			t.Errorf("Bytes(%d) = %q, want %q", tt.n, got, tt.bin)
		}
		if got := BytesIEC(tt.n); got != tt.iec {
			t.Errorf("BytesIEC(%d) = %q, want %q", tt.n, got, tt.iec)
		}
		if got := BytesSI(tt.n); got != tt.si {
			t.Errorf("BytesSI(%d) = %q, want %q", tt.n, got, tt.si)
		}
	}
}

func TestFileSize(t *testing.T) {
	tests := []struct {
		n    int64
		want string
	}{
		{0, "0 B"}, {-1, "-1 B"}, {1023, "1023 B"}, {1024, "1 KB"},
		{5242880, "5 MB"}, {1590000, "1.5 MB"}, {-1590000, "-1.5 MB"},
		{1048524, "1023.9 KB"}, {1048525, "1 MB"},
		{math.MaxInt64, "8 EB"}, {math.MinInt64, "-8 EB"},
	}
	for _, tt := range tests {
		if got := FileSize(tt.n); got != tt.want {
			t.Errorf("FileSize(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}
